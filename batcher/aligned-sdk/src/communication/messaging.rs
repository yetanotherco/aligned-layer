use ethers::signers::Signer;
use ethers::types::Address;
use futures_util::{stream::SplitStream, SinkExt, StreamExt};
use log::{debug, error, info};
use std::sync::Arc;
use tokio::{net::TcpStream, sync::Mutex};

use ethers::{core::k256::ecdsa::SigningKey, signers::Wallet, types::U256};
use futures_util::future::Ready;
use futures_util::stream::{SplitSink, TryFilter};
use tokio_tungstenite::{tungstenite::Message, MaybeTlsStream, WebSocketStream};

use crate::communication::serialization::{cbor_deserialize, cbor_serialize};
use crate::{
    communication::batch::handle_batch_inclusion_data,
    core::{
        errors::SubmitError,
        types::{
            AlignedVerificationData, ClientMessage, NoncedVerificationData, ResponseMessage,
            VerificationData, VerificationDataCommitment,
        },
    },
};

pub type ResponseStream = TryFilter<
    SplitStream<WebSocketStream<MaybeTlsStream<TcpStream>>>,
    Ready<bool>,
    fn(&Message) -> Ready<bool>,
>;

pub async fn send_messages(
    ws_write: Arc<Mutex<SplitSink<WebSocketStream<MaybeTlsStream<TcpStream>>, Message>>>,
    payment_service_addr: Address,
    verification_data: &[VerificationData],
    max_fees: &[U256],
    wallet: Wallet<SigningKey>,
    mut nonce: U256,
    sender_channel: tokio::sync::mpsc::Sender<VerificationDataCommitment>,
) -> Result<(), SubmitError> {
    let chain_id = U256::from(wallet.chain_id());
    let mut ws_write = ws_write.lock().await;

    for (idx, verification_data) in verification_data.iter().enumerate() {
        // Build each message to send
        let verification_data = NoncedVerificationData::new(
            verification_data.clone(),
            nonce,
            max_fees[idx],
            chain_id,
            payment_service_addr,
        );

        nonce += U256::one();
        let msg = ClientMessage::new(verification_data.clone(), wallet.clone()).await;
        let msg_bin = cbor_serialize(&msg).map_err(SubmitError::SerializationError)?;

        // Send the message
        ws_write
            .send(Message::Binary(msg_bin.clone()))
            .await
            .map_err(SubmitError::WebSocketConnectionError)?;

        debug!("{:?} Message sent", idx);

        match sender_channel.send(verification_data.into()).await {//.map_err(|e| SubmitError::GenericError(e.to_string()))?;
            Ok(_) => {
                debug!("Message sent to channel");
            }
            Err(e) if e.to_string() == "channel closed" => { // happens when receive has exited, because batcher replied with error
                error!("Error sending message because batcher has previously replied with an error");
                // return Err(SubmitError::GenericError(("Batcher has previously replied with an error").to_string()));
                return Ok(());
            }
            Err(e) => {
                error!("Error sending message to channel: {:?}", e.to_string());
                return Err(SubmitError::GenericError(e.to_string()));
            }
        }
    }

    //sender_channel will be closed as it falls out of scope, sending a 'None' to the receiver
    info!("All messages sent");
    Ok(()) 
}

pub async fn receive(
    response_stream: Arc<Mutex<ResponseStream>>,
    mut receiver_channnel: tokio::sync::mpsc::Receiver<VerificationDataCommitment>,
) -> Result<Vec<AlignedVerificationData>, SubmitError> {
    // Responses are filtered to only admit binary or close messages.
    let mut response_stream = response_stream.lock().await;
    let mut aligned_verification_data: Vec<AlignedVerificationData> = Vec::new();

    while let Some(verification_data_commitment) = receiver_channnel.recv().await { //while there are messages in the channel
        // Read from WS
        let msg = response_stream.next().await.unwrap().map_err(SubmitError::WebSocketConnectionError)?;

        // If websocket was closed prematurely:
        if let Message::Close(close_frame) = msg {
            if let Some(close_msg) = close_frame {
                return Err(SubmitError::WebSocketClosedUnexpectedlyError(
                    close_msg.to_owned(),
                ));
            }
            return Err(SubmitError::GenericError(
                "Connection was closed before receive() processed all sent messages "
                    .to_string(),
            ));
        }

        process_batch_inclusion_data(
            msg,
            &mut aligned_verification_data,
            verification_data_commitment,
        ).await?; // If batcher returned an error, this will close the channel and return the error
    }

    info!("All message responses handled succesfully");
    Ok(aligned_verification_data)
}

async fn process_batch_inclusion_data(
    msg: Message,
    aligned_verification_data: &mut Vec<AlignedVerificationData>,
    verification_data_commitment: VerificationDataCommitment,
) -> Result<(), SubmitError> {

    let data = msg.into_data();
    match cbor_deserialize(data.as_slice()) {
        Ok(ResponseMessage::BatchInclusionData(batch_inclusion_data)) => { //OK case. Proofs was valid and it was included in this batch.
            let _ = handle_batch_inclusion_data(
                batch_inclusion_data,
                aligned_verification_data,
                verification_data_commitment,
            );
        }
        Ok(ResponseMessage::ReplacementMessageReceived) => {
            // This message is not processed, it is only used to signal the client that the replacement message was received by the batcher.
            // This is because the sender expects to receive the same amount of messages as it has sent.
        }
        Ok(ResponseMessage::InvalidNonce) => {
            error!("Batcher responded with invalid nonce");
            return Err(SubmitError::InvalidNonce);
        }
        Ok(ResponseMessage::InvalidSignature) => {
            error!("Batcher responded with invalid signature");
            return Err(SubmitError::InvalidSignature);
        }
        Ok(ResponseMessage::ProofTooLarge) => {
            error!("Batcher responded with proof too large");
            return Err(SubmitError::ProofTooLarge);
        }
        Ok(ResponseMessage::InvalidMaxFee) => {
            error!("Batcher responded with invalid max fee");
            return Err(SubmitError::InvalidMaxFee);
        }
        Ok(ResponseMessage::InsufficientBalance(addr)) => {
            error!("Batcher responded with insufficient balance");
            return Err(SubmitError::InsufficientBalance(addr));
        }
        Ok(ResponseMessage::InvalidChainId) => {
            error!("Batcher responded with invalid chain id");
            return Err(SubmitError::InvalidChainId);
        }
        Ok(ResponseMessage::InvalidReplacementMessage) => {
            error!("Batcher responded with invalid replacement message");
            return Err(SubmitError::InvalidReplacementMessage);
        }
        Ok(ResponseMessage::AddToBatchError) => {
            error!("Batcher responded with add to batch error");
            return Err(SubmitError::AddToBatchError);
        }
        Ok(ResponseMessage::EthRpcError) => {
            error!("Batcher experienced Eth RPC connection error");
            return Err(SubmitError::EthereumProviderError(
                "Batcher experienced Eth RPC connection error".to_string(),
            ));
        }
        Ok(ResponseMessage::InvalidPaymentServiceAddress(received_addr, expected_addr)) => {
            error!(
                "Batcher responded with invalid payment service address: {:?}, expected: {:?}",
                received_addr, expected_addr
            );
            return Err(SubmitError::InvalidPaymentServiceAddress(
                received_addr,
                expected_addr,
            ));
        }
        Ok(ResponseMessage::InvalidProof(reason)) => { 
            error!("Batcher responded with invalid proof: {}", reason);
            return Err(SubmitError::InvalidProof(reason));
        }
        Ok(ResponseMessage::CreateNewTaskError(merkle_root, error)) => {
            error!("Batcher responded with create new task error: {}", error);
            return Err(SubmitError::BatchSubmissionFailed(
                "Could not create task with merkle root ".to_owned() + &merkle_root + ", failed with error: " + &error, 
            ));
        }
        Ok(ResponseMessage::ProtocolVersion(_)) => {
            error!("Batcher responded with protocol version instead of batch inclusion data");
            return Err(SubmitError::UnexpectedBatcherResponse(
                "Batcher responded with protocol version instead of batch inclusion data"
                    .to_string(),
            ));
        }
        Ok(ResponseMessage::BatchReset) => {
            error!("Batcher responded with batch reset");
            return Err(SubmitError::ProofQueueFlushed);
        }
        Ok(ResponseMessage::Error(e)) => {
            error!("Batcher responded with error: {}", e);
            error!("Batcher responded with error: {}", e);
        }
        Err(e) => {
            error!("Error while deserializing batch inclusion data: {}", e);
            return Err(SubmitError::SerializationError(e));
        }
    }

    Ok(())
}
