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
) -> Result<Vec<VerificationDataCommitment>, SubmitError> {
    let mut sent_verification_data = Vec::new();

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

        sent_verification_data.push(verification_data.clone());
    }

    // This vector is reversed so that when responses are received, the commitments corresponding
    // to that response can simply be popped of this vector.
    let verification_data_commitments_rev: Vec<VerificationDataCommitment> =
        sent_verification_data
            .into_iter()
            .map(|vd| vd.into())
            .rev()
            .collect();

    Ok(verification_data_commitments_rev)
}

pub async fn receive(
    response_stream: Arc<Mutex<ResponseStream>>,
    total_messages: usize,
    verification_data_commitments_rev: &mut Vec<VerificationDataCommitment>,
) -> Result<Vec<AlignedVerificationData>, SubmitError> {
    // Responses are filtered to only admit binary or close messages.
    let mut response_stream = response_stream.lock().await;
    let mut aligned_verification_data: Vec<AlignedVerificationData> = Vec::new();
    let mut num_responses: usize = 0;

    while let Some(Ok(msg)) = response_stream.next().await {
        if let Message::Close(close_frame) = msg {
            if let Some(close_msg) = close_frame {
                return Err(SubmitError::WebSocketClosedUnexpectedlyError(
                    close_msg.to_owned(),
                ));
            }
            return Err(SubmitError::GenericError(
                "Connection was closed without close message before receiving all messages"
                    .to_string(),
            ));
        }

        process_batch_inclusion_data(
            msg,
            &mut aligned_verification_data,
            verification_data_commitments_rev,
        )
        .await?;
        
        num_responses += 1;
        if num_responses == total_messages {
            info!("All message responses received");
            return Ok(aligned_verification_data);
        }
    }

    Err(SubmitError::GenericError(
        "Connection was closed without close message before receiving all messages".to_string(),
    ))
}

async fn process_batch_inclusion_data(
    msg: Message,
    aligned_verification_data: &mut Vec<AlignedVerificationData>,
    verification_data_commitments_rev: &mut Vec<VerificationDataCommitment>,
) -> Result<(), SubmitError> {

    let data = msg.into_data();
    match cbor_deserialize(data.as_slice()) {
        Ok(ResponseMessage::BatchInclusionData(batch_inclusion_data)) => { //OK case. Proofs was valid and it was included in this batch.
            let _ = handle_batch_inclusion_data(
                batch_inclusion_data,
                aligned_verification_data,
                verification_data_commitments_rev,
            );
        }
        Ok(ResponseMessage::InvalidNonce) => {
            return Err(SubmitError::InvalidNonce);
        }
        Ok(ResponseMessage::InvalidSignature) => {
            return Err(SubmitError::InvalidSignature);
        }
        Ok(ResponseMessage::ProofTooLarge) => {
            return Err(SubmitError::ProofTooLarge);
        }
        Ok(ResponseMessage::InvalidMaxFee) => {
            return Err(SubmitError::InvalidMaxFee);
        }
        Ok(ResponseMessage::InsufficientBalance) => {
            return Err(SubmitError::InsufficientBalance);
        }
        Ok(ResponseMessage::InvalidChainId) => {
            return Err(SubmitError::InvalidChainId);
        }
        Ok(ResponseMessage::InvalidReplacementMessage) => {
            return Err(SubmitError::InvalidReplacementMessage);
        }
        Ok(ResponseMessage::AddToBatchError) => {
            return Err(SubmitError::AddToBatchError);
        }
        Ok(ResponseMessage::EthRpcError) => {
            return Err(SubmitError::EthereumProviderError(
                "Batcher experienced Eth RPC connection error".to_string(),
            ));
        }
        Ok(ResponseMessage::InvalidPaymentServiceAddress(received_addr, expected_addr)) => {
            return Err(SubmitError::InvalidPaymentServiceAddress(
                received_addr,
                expected_addr,
            ));
        }
        Ok(ResponseMessage::InvalidProof(reason)) => { 
            return Err(SubmitError::InvalidProof(reason));
        }
        Ok(ResponseMessage::CreateNewTaskError(merkle_root, error)) => {
            return Err(SubmitError::BatchSubmissionFailed(
                "Could not create task with merkle root ".to_owned() + &merkle_root + ", failed with error: " + &error, 
            ));
        }
        Ok(ResponseMessage::ProtocolVersion(_)) => {
            return Err(SubmitError::UnexpectedBatcherResponse(
                "Batcher responded with protocol version instead of batch inclusion data"
                    .to_string(),
            ));
        }
        Ok(ResponseMessage::BatchReset) => {
            return Err(SubmitError::ProofQueueFlushed);
        }
        Ok(ResponseMessage::Error(e)) => {
            error!("Batcher responded with error: {}", e);
        }
        Err(e) => {
            return Err(SubmitError::SerializationError(e));
        }
    }

    Ok(())
}
