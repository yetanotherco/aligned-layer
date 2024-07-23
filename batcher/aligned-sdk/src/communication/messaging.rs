use futures_util::{future, stream::SplitStream, SinkExt, StreamExt, TryStreamExt};
use log::{debug, error};
use std::sync::Arc;
use tokio::{net::TcpStream, sync::Mutex};

use ethers::{core::k256::ecdsa::SigningKey, signers::Wallet, types::U256};
use futures_util::stream::SplitSink;
use tokio_tungstenite::{tungstenite::Message, MaybeTlsStream, WebSocketStream};

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

pub async fn send_messages(
    ws_write: Arc<Mutex<SplitSink<WebSocketStream<MaybeTlsStream<TcpStream>>, Message>>>,
    verification_data: &[VerificationData],
    wallet: Wallet<SigningKey>,
    nonce: U256,
) -> Result<Vec<NoncedVerificationData>, SubmitError> {
    let mut sent_verification_data = Vec::new();

    let mut ws_write = ws_write.lock().await;

    let mut nonce = nonce.clone();
    let mut nonce_bytes = [0u8; 32];

    for verification_data in verification_data.iter() {
        nonce.to_big_endian(&mut nonce_bytes);

        let verification_data = NoncedVerificationData::new(verification_data.clone(), nonce_bytes);
        nonce += U256::one();

        let msg = ClientMessage::new(verification_data.clone(), wallet.clone());
        let msg_str = serde_json::to_string(&msg).map_err(SubmitError::SerializationError)?;
        ws_write
            .send(Message::Text(msg_str.clone()))
            .await
            .map_err(SubmitError::WebSocketConnectionError)?;
        sent_verification_data.push(verification_data.clone());
        debug!("Message sent...");
    }

    Ok(sent_verification_data)
}

pub async fn receive(
    ws_read: SplitStream<WebSocketStream<MaybeTlsStream<TcpStream>>>,
    ws_write: Arc<Mutex<SplitSink<WebSocketStream<MaybeTlsStream<TcpStream>>, Message>>>,
    total_messages: usize,
    num_responses: Arc<Mutex<usize>>,
    verification_data_commitments_rev: &mut Vec<VerificationDataCommitment>,
) -> Result<Option<Vec<AlignedVerificationData>>, SubmitError> {
    // Responses are filtered to only admit binary or close messages.
    let mut response_stream =
        ws_read.try_filter(|msg| future::ready(msg.is_binary() || msg.is_close()));

    let mut aligned_verification_data: Vec<AlignedVerificationData> = Vec::new();

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
            num_responses.clone(),
        )
        .await?;

        if *num_responses.lock().await == total_messages {
            debug!("All messages responded. Closing connection...");
            ws_write.lock().await.close().await?;
            return Ok(Some(aligned_verification_data));
        }
    }

    Ok(None)
}

async fn process_batch_inclusion_data(
    msg: Message,
    aligned_verification_data: &mut Vec<AlignedVerificationData>,
    verification_data_commitments_rev: &mut Vec<VerificationDataCommitment>,
    num_responses: Arc<Mutex<usize>>,
) -> Result<(), SubmitError> {
    let mut num_responses_lock = num_responses.lock().await;
    *num_responses_lock += 1;

    let data = msg.into_data();
    match serde_json::from_slice::<ResponseMessage>(&data) {
        Ok(ResponseMessage::BatchInclusionData(batch_inclusion_data)) => {
            let _ = handle_batch_inclusion_data(
                batch_inclusion_data,
                aligned_verification_data,
                verification_data_commitments_rev,
            );
        }
        Ok(ResponseMessage::ProtocolVersion(_)) => {
            return Err(SubmitError::UnexpectedBatcherResponse(
                "Batcher responded with protocol version instead of batch inclusion data"
                    .to_string(),
            ));
        }
        Ok(ResponseMessage::Error(e)) => {
            error!("Batcher responded with error: {}", e);
        }
        Ok(ResponseMessage::VerificationError()) => {
            error!("Invalid proof");
        }
        Ok(ResponseMessage::ProofTooLargeError()) => {
            error!("Proof is too large");
        }
        Ok(ResponseMessage::InsufficientBalanceError(address)) => {
            error!("Insufficient balance for address: {}", address);
        }
        Ok(ResponseMessage::SignatureVerificationError()) => {
            error!("Failed to verify the signature");
        }
        Ok(ResponseMessage::InvalidNonceError) => {
            error!("Invalid nonce")
        }
        Ok(ResponseMessage::CreateNewTaskError(merkle_root)) => {
            return Err(SubmitError::CreateNewTaskError(
                "Could not create task with merkle root ".to_owned() + &merkle_root,
            ));
        }
        Err(e) => {
            return Err(SubmitError::SerializationError(e));
        }
    }

    Ok(())
}
