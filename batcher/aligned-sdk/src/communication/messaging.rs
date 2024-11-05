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
use crate::core::types::BatchInclusionData;
use crate::{
    communication::batch::process_batcher_response,
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

// Sends the proofs to the batcher via WS
// Stores the proofs sent in an array
// Returns the array
pub async fn send_messages(
    ws_write: Arc<Mutex<SplitSink<WebSocketStream<MaybeTlsStream<TcpStream>>, Message>>>,
    payment_service_addr: Address,
    verification_data: &[VerificationData],
    max_fees: &[U256],
    wallet: Wallet<SigningKey>,
    mut nonce: U256,
) -> Result<Vec<NoncedVerificationData>, SubmitError> {
    let chain_id = U256::from(wallet.chain_id());
    let mut ws_write = ws_write.lock().await;
    let mut sent_verification_data: Vec<NoncedVerificationData> = Vec::new();

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
        
        // Save the verification data commitment to read its response later
        sent_verification_data.push(verification_data);
    }

    info!("All messages sent");
    Ok(sent_verification_data) 
}


// Instead of using a channel, use a storage.
// Using a 
// From there, you can match received messages to the ones you sent.

// TODO missing analyzing which is the last expected nonce.
// When received message of last expected nonce, i can exit this function

// Receives the array of proofs sent
// Reads the WS responses
// Matches each response with the corresponding proof sent
// finishes when the last proof sent receives its response
pub async fn receive(
    response_stream: Arc<Mutex<ResponseStream>>,
    mut sent_verification_data: Vec<NoncedVerificationData>,
) -> Result<Vec<AlignedVerificationData>, SubmitError> {
    // Responses are filtered to only admit binary or close messages.
    let mut response_stream = response_stream.lock().await;
    let mut aligned_submitted_data: Vec<AlignedVerificationData> = Vec::new();
    let last_proof_nonce = get_biggest_nonce(&sent_verification_data);

    // read from WS
    while let Some(Ok(msg)) = response_stream.next().await {
        // unexpected WS close:
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
 
        let batch_inclusion_data_message = handle_batcher_response(
            msg,
        ).await?;

        let related_verification_data = match_batcher_response_with_stored_verification_data(
            &batch_inclusion_data_message,
            &mut sent_verification_data,
        )?;
            
        let aligned_verification_data = process_batcher_response(
            &batch_inclusion_data_message,
            &related_verification_data,
        )?;

        aligned_submitted_data.push(aligned_verification_data);
        info!("Message response handled succesfully");

        if batch_inclusion_data_message.user_nonce == last_proof_nonce {
            break;
        }
    }

    debug!("All message responses handled succesfully");
    Ok(aligned_submitted_data)
}

async fn handle_batcher_response(
    msg: Message,
) -> Result<BatchInclusionData, SubmitError> {

    let data = msg.into_data();
    match cbor_deserialize(data.as_slice()) {
        Ok(ResponseMessage::BatchInclusionData(batch_inclusion_data)) => { //OK case. Proofs was valid and it was included in this batch.
            return Ok(batch_inclusion_data);
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
            return Err(SubmitError::GenericError(e))
        }
        Err(e) => {
            error!("Error while deserializing batch inclusion data: {}", e);
            return Err(SubmitError::SerializationError(e));
        }
    }
}

// Used to match the message received from the batcher,
// with the NoncedVerificationData you sent
// This is used to verify the proof you sent was indeed included in the batch
fn match_batcher_response_with_stored_verification_data(
    batch_inclusion_data: &BatchInclusionData,
    sent_verification_data: &mut Vec<NoncedVerificationData>,
) -> Result<VerificationDataCommitment, SubmitError> {
    debug!("Matching verification data with batcher response ...");
    let mut index = None;
    for (i, sent_nonced_verification_data) in sent_verification_data.iter_mut().enumerate() {
        if sent_nonced_verification_data.nonce == batch_inclusion_data.user_nonce {
            debug!("local nonced verification data matched with batcher response");
            index = Some(i);
            break;
        }
    }

    if let Some(i) = index {
        let verification_data = sent_verification_data.swap_remove(i); //TODO maybe only remove?
        return Ok(verification_data.verification_data.clone().into());
    }

    Err(SubmitError::InvalidProofInclusionData)
}

// Returns the biggest nonce from the sent verification data
// Used to know which is the last proof sent to the Batcher,
// to know when to stop reading the WS for responses
fn get_biggest_nonce(sent_verification_data: &Vec<NoncedVerificationData>) -> U256 {
    let mut biggest_nonce = U256::zero();
    for verification_data in sent_verification_data.iter() {
        if verification_data.nonce > biggest_nonce {
            biggest_nonce = verification_data.nonce;
        }
    }
    biggest_nonce
}
