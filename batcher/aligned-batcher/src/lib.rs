extern crate core;

use aligned_sdk::communication::serialization::{cbor_deserialize, cbor_serialize};
use aligned_sdk::eth::batcher_payment_service::SignatureData;
use config::NonPayingConfig;
use connection::{send_message, WsMessageSink};
use dotenvy::dotenv;
use ethers::contract::ContractError;
use ethers::signers::Signer;
use types::batch_state::BatchState;
use types::user_state::UserState;

use std::collections::HashMap;
use std::env;
use std::iter::repeat;
use std::net::SocketAddr;
use std::sync::Arc;

use aligned_sdk::core::types::{
    ClientMessage, NoncedVerificationData, ResponseMessage, ValidityResponseMessage,
    VerificationCommitmentBatch, VerificationData, VerificationDataCommitment,
};
use aws_sdk_s3::client::Client as S3Client;
use eth::{try_create_new_task, BatcherPaymentService, CreateNewTaskFeeParams, SignerMiddlewareT};
use ethers::prelude::{Middleware, Provider};
use ethers::providers::Ws;
use ethers::types::{Address, Signature, TransactionReceipt, U256};
use futures_util::{future, SinkExt, StreamExt, TryStreamExt};
use lambdaworks_crypto::merkle_tree::merkle::MerkleTree;
use lambdaworks_crypto::merkle_tree::traits::IsMerkleTreeBackend;
use log::{debug, error, info, warn};
use tokio::net::{TcpListener, TcpStream};
use tokio::sync::{Mutex, MutexGuard, RwLock};
use tokio_tungstenite::tungstenite::{Error, Message};
use types::batch_queue::{self, BatchQueueEntry, BatchQueueEntryPriority};
use types::errors::{BatcherError, BatcherSendError};

use crate::config::{ConfigFromYaml, ContractDeploymentOutput};

mod config;
mod connection;
mod eth;
pub mod gnark;
pub mod halo2;
pub mod risc_zero;
pub mod s3;
pub mod sp1;
pub mod types;
mod zk_utils;

const AGGREGATOR_GAS_COST: u128 = 400_000;
const BATCHER_SUBMISSION_BASE_GAS_COST: u128 = 125_000;
pub(crate) const ADDITIONAL_SUBMISSION_GAS_COST_PER_PROOF: u128 = 13_000;
pub(crate) const CONSTANT_GAS_COST: u128 =
    ((AGGREGATOR_GAS_COST * DEFAULT_AGGREGATOR_FEE_MULTIPLIER) / DEFAULT_AGGREGATOR_FEE_DIVIDER)
        + BATCHER_SUBMISSION_BASE_GAS_COST;

const DEFAULT_MAX_FEE_PER_PROOF: u128 = ADDITIONAL_SUBMISSION_GAS_COST_PER_PROOF * 100_000_000_000; // gas_price = 100 Gwei = 0.0000001 ether (high gas price)
const MIN_FEE_PER_PROOF: u128 = ADDITIONAL_SUBMISSION_GAS_COST_PER_PROOF * 100_000_000; // gas_price = 0.1 Gwei = 0.0000000001 ether (low gas price)
const RESPOND_TO_TASK_FEE_LIMIT_MULTIPLIER: u128 = 5; // to set the respondToTaskFeeLimit variable higher than fee_for_aggregator
const RESPOND_TO_TASK_FEE_LIMIT_DIVIDER: u128 = 2;
const DEFAULT_AGGREGATOR_FEE_MULTIPLIER: u128 = 3; // to set the feeForAggregator variable higher than what was calculated
const DEFAULT_AGGREGATOR_FEE_DIVIDER: u128 = 2;

pub struct Batcher {
    s3_client: S3Client,
    s3_bucket_name: String,
    download_endpoint: String,
    eth_ws_provider: Provider<Ws>,
    eth_ws_provider_fallback: Provider<Ws>,
    chain_id: U256,
    payment_service: BatcherPaymentService,
    payment_service_fallback: BatcherPaymentService,
    batch_state: Mutex<BatchState>,
    max_block_interval: u64,
    min_batch_len: usize,
    max_proof_size: usize,
    max_batch_size: usize,
    last_uploaded_batch_block: Mutex<u64>,
    pre_verification_is_enabled: bool,
    non_paying_config: Option<NonPayingConfig>,
    posting_batch: Mutex<bool>,
    user_states: RwLock<HashMap<Address, Mutex<UserState>>>,
}

impl Batcher {
    pub async fn new(config_file: String) -> Self {
        dotenv().ok();

        // https://docs.aws.amazon.com/sdk-for-rust/latest/dg/localstack.html
        let upload_endpoint = env::var("UPLOAD_ENDPOINT").ok();

        let s3_bucket_name =
            env::var("AWS_BUCKET_NAME").expect("AWS_BUCKET_NAME not found in environment");

        let download_endpoint =
            env::var("DOWNLOAD_ENDPOINT").expect("DOWNLOAD_ENDPOINT not found in environment");

        let s3_client = s3::create_client(upload_endpoint).await;

        let config = ConfigFromYaml::new(config_file);
        let deployment_output =
            ContractDeploymentOutput::new(config.aligned_layer_deployment_config_file_path);

        let eth_ws_provider =
            Provider::connect_with_reconnects(&config.eth_ws_url, config.batcher.eth_ws_reconnects)
                .await
                .expect("Failed to get ethereum websocket provider");

        let eth_ws_provider_fallback = Provider::connect_with_reconnects(
            &config.eth_ws_url_fallback,
            config.batcher.eth_ws_reconnects,
        )
        .await
        .expect("Failed to get fallback ethereum websocket provider");

        let eth_rpc_provider =
            eth::get_provider(config.eth_rpc_url.clone()).expect("Failed to get provider");

        let eth_rpc_provider_fallback = eth::get_provider(config.eth_rpc_url_fallback.clone())
            .expect("Failed to get fallback provider");

        // FIXME(marian): We are getting just the last block number right now, but we should really
        // have the last submitted batch block registered and query it when the batcher is initialized.
        let last_uploaded_batch_block = match eth_rpc_provider.get_block_number().await {
            Ok(block_num) => block_num,
            Err(e) => {
                warn!(
                    "Failed to get block number with main rpc, trying with fallback rpc. Err: {:?}",
                    e
                );
                eth_rpc_provider_fallback
                    .get_block_number()
                    .await
                    .expect("Failed to get block number with fallback rpc")
            }
        };

        let last_uploaded_batch_block = last_uploaded_batch_block.as_u64();

        let chain_id = match eth_rpc_provider.get_chainid().await {
            Ok(chain_id) => chain_id,
            Err(e) => {
                warn!("Failed to get chain id with main rpc: {}", e);
                eth_rpc_provider_fallback
                    .get_chainid()
                    .await
                    .expect("Failed to get chain id with fallback rpc")
            }
        };

        let payment_service = eth::get_batcher_payment_service(
            eth_rpc_provider,
            config.ecdsa.clone(),
            deployment_output.addresses.batcher_payment_service.clone(),
        )
        .await
        .expect("Failed to get Batcher Payment Service contract");

        let payment_service_fallback = eth::get_batcher_payment_service(
            eth_rpc_provider_fallback,
            config.ecdsa,
            deployment_output.addresses.batcher_payment_service,
        )
        .await
        .expect("Failed to get fallback Batcher Payment Service contract");

        let user_states = RwLock::new(HashMap::new());
        let non_paying_config = if let Some(non_paying_config) = config.batcher.non_paying {
            warn!("Non-paying address configuration detected. Will replace non-paying address {} with configured address.",
                non_paying_config.address);

            let non_paying_config = NonPayingConfig::from_yaml_config(non_paying_config).await;
            let nonpaying_nonce = payment_service
                .user_nonces(non_paying_config.replacement.address())
                .call()
                .await
                .expect("Could not get non-paying nonce from Ethereum");

            let non_paying_user_state = Mutex::new(UserState::new_non_paying(nonpaying_nonce));
            user_states.write().await.insert(
                non_paying_config.replacement.address(),
                non_paying_user_state,
            );

            Some(non_paying_config)
        } else {
            None
        };

        Self {
            s3_client,
            s3_bucket_name,
            download_endpoint,
            eth_ws_provider,
            eth_ws_provider_fallback,
            chain_id,
            payment_service,
            payment_service_fallback,
            batch_state: Mutex::new(BatchState::new()),
            max_block_interval: config.batcher.block_interval,
            min_batch_len: config.batcher.batch_size_interval,
            max_proof_size: config.batcher.max_proof_size,
            max_batch_size: config.batcher.max_batch_size,
            last_uploaded_batch_block: Mutex::new(last_uploaded_batch_block),
            pre_verification_is_enabled: config.batcher.pre_verification_is_enabled,
            non_paying_config,
            posting_batch: Mutex::new(false),
            user_states,
        }
    }

    pub async fn listen_connections(self: Arc<Self>, address: &str) -> Result<(), BatcherError> {
        // Create the event loop and TCP listener we'll accept connections on.
        let listener = TcpListener::bind(address)
            .await
            .map_err(|e| BatcherError::TcpListenerError(e.to_string()))?;
        info!("Listening on: {}", address);

        // Let's spawn the handling of each connection in a separate task.
        while let Ok((stream, addr)) = listener.accept().await {
            let batcher = self.clone();
            tokio::spawn(batcher.handle_connection(stream, addr));
        }
        Ok(())
    }

    pub async fn listen_new_blocks(self: Arc<Self>) -> Result<(), BatcherError> {
        let mut stream = self
            .eth_ws_provider
            .subscribe_blocks()
            .await
            .map_err(|e| BatcherError::EthereumSubscriptionError(e.to_string()))?;

        let mut stream_fallback = self
            .eth_ws_provider_fallback
            .subscribe_blocks()
            .await
            .map_err(|e| BatcherError::EthereumSubscriptionError(e.to_string()))?;

        let last_seen_block = Mutex::<u64>::new(0);

        while let Some(block) = tokio::select! {
            block = stream.next() => block,
            block = stream_fallback.next() => block,
        } {
            let batcher = self.clone();
            let block_number = block.number.unwrap_or_default();
            let block_number = u64::try_from(block_number).unwrap_or_default();

            {
                let mut last_seen_block = last_seen_block.lock().await;
                if block_number <= *last_seen_block {
                    continue;
                }
                *last_seen_block = block_number;
            }

            info!("Received new block: {}", block_number);
            tokio::spawn(async move {
                if let Err(e) = batcher.handle_new_block(block_number).await {
                    error!("Error when handling new block: {:?}", e);
                };
            });
        }

        Ok(())
    }

    async fn handle_connection(
        self: Arc<Self>,
        raw_stream: TcpStream,
        addr: SocketAddr,
    ) -> Result<(), BatcherError> {
        info!("Incoming TCP connection from: {}", addr);
        let ws_stream = tokio_tungstenite::accept_async(raw_stream).await?;

        debug!("WebSocket connection established: {}", addr);
        let (outgoing, incoming) = ws_stream.split();
        let outgoing = Arc::new(RwLock::new(outgoing));

        let protocol_version_msg = ResponseMessage::ProtocolVersion(
            aligned_sdk::communication::protocol::EXPECTED_PROTOCOL_VERSION,
        );

        let serialized_protocol_version_msg = cbor_serialize(&protocol_version_msg)
            .map_err(|e| BatcherError::SerializationError(e.to_string()))?;

        outgoing
            .write()
            .await
            .send(Message::binary(serialized_protocol_version_msg))
            .await?;

        match incoming
            .try_filter(|msg| future::ready(msg.is_binary()))
            .try_for_each(|msg| self.clone().handle_message(msg, outgoing.clone()))
            .await
        {
            Err(e) => error!("Unexpected error: {}", e),
            Ok(_) => info!("{} disconnected", &addr),
        }

        Ok(())
    }

    /// Handle an individual message from the client.
    async fn handle_message(
        self: Arc<Self>,
        message: Message,
        ws_conn_sink: WsMessageSink,
    ) -> Result<(), Error> {
        // Deserialize verification data from message
        let client_msg: ClientMessage = match cbor_deserialize(message.into_data().as_slice()) {
            Ok(msg) => msg,
            Err(e) => {
                warn!("Failed to deserialize message: {}", e);
                return Ok(());
            }
        };

        info!(
            "Received message with nonce: {}",
            client_msg.verification_data.nonce
        );

        if client_msg.verification_data.chain_id != self.chain_id {
            warn!(
                "Received message with incorrect chain id: {}",
                client_msg.verification_data.chain_id
            );

            send_message(
                ws_conn_sink.clone(),
                ValidityResponseMessage::InvalidChainId,
            )
            .await;

            return Ok(());
        }

        info!("Verifying message signature...");
        let Ok(addr) = client_msg.verify_signature() else {
            error!("Signature verification error");
            send_message(
                ws_conn_sink.clone(),
                ValidityResponseMessage::InvalidSignature,
            )
            .await;
            return Ok(());
        };
        info!("Message signature verified");

        if self.is_nonpaying(&addr) {
            println!("Handling nonpaying MESSAGEEE!");
            return self
                .handle_nonpaying_msg(ws_conn_sink.clone(), client_msg)
                .await;
        }

        if self.user_balance_is_unlocked(&addr).await {
            send_message(
                ws_conn_sink.clone(),
                ValidityResponseMessage::InsufficientBalance(addr),
            )
            .await;
            println!("UNLOCKING THE USER {addr} STATEEE!");
            return Ok(());
        }

        let nonced_verification_data = client_msg.verification_data;
        if nonced_verification_data.verification_data.proof.len() > self.max_proof_size {
            error!("Proof size exceeds the maximum allowed size.");
            send_message(ws_conn_sink.clone(), ValidityResponseMessage::ProofTooLarge).await;
            return Ok(());
        }

        // When pre-verification is enabled, batcher will verify proofs for faster feedback with clients
        if self.pre_verification_is_enabled
            && !zk_utils::verify(&nonced_verification_data.verification_data).await
        {
            error!("Invalid proof detected. Verification failed.");
            send_message(ws_conn_sink.clone(), ValidityResponseMessage::InvalidProof).await;
            return Ok(()); // Send error message to the client and return
        }

        // Nonce and max fee verification
        let msg_nonce = nonced_verification_data.nonce;
        let max_fee = nonced_verification_data.max_fee;
        if max_fee < U256::from(MIN_FEE_PER_PROOF) {
            error!("The max fee signed in the message is less than the accepted minimum fee to be included in the batch.");
            send_message(ws_conn_sink.clone(), ValidityResponseMessage::InvalidMaxFee).await;
            return Ok(());
        }

        // Check that we had a user state entry for this user and insert it if not.
        if !self.user_states.read().await.contains_key(&addr) {
            let new_user_state = Mutex::new(UserState::new());
            self.user_states
                .write()
                .await
                .insert(addr.clone(), new_user_state);
        }

        // At this point, we will have a user state for sure, since we have inserted it
        // if not already present.
        println!("LOCKING USER STATE FOR ADDR: {addr}");
        let user_state_read_lock = self.user_states.read().await;
        let mut user_state = user_state_read_lock.get(&addr).unwrap().lock().await;
        println!("USER STATE FOR {addr} LOCKED BABY!");

        println!("LOCKING BATCH STATE IN HANDLE MESSAGE");
        let batch_state_lock = self.batch_state.lock().await;
        println!("BATCH STATE LOCKEDIN HANDLE MESSAGE");

        if !self
            .check_user_balance(&addr, user_state.proofs_in_batch + 1)
            .await
        {
            send_message(
                ws_conn_sink.clone(),
                ValidityResponseMessage::InsufficientBalance(addr),
            )
            .await;
            println!("UNLOCKING THE USER {addr} STATEEE!");
            return Ok(());
        }

        let expected_nonce = if let Some(cached_nonce) = user_state.nonce {
            cached_nonce
        } else {
            let ethereum_user_nonce = match self.get_user_nonce_from_ethereum(addr).await {
                Ok(ethereum_user_nonce) => ethereum_user_nonce,
                Err(e) => {
                    error!(
                        "Failed to get user nonce from Ethereum for address {:?}. Error: {:?}",
                        addr, e
                    );
                    send_message(ws_conn_sink.clone(), ValidityResponseMessage::InvalidNonce).await;

                    return Ok(());
                }
            };

            user_state.nonce = Some(ethereum_user_nonce);
            ethereum_user_nonce
        };

        if expected_nonce < msg_nonce {
            warn!(
                "Invalid nonce for address {addr}, had nonce {:?} < {:?}",
                expected_nonce, msg_nonce
            );
            println!("UNLOCKING THE USER {addr} STATEEE!");
            send_message(ws_conn_sink.clone(), ValidityResponseMessage::InvalidNonce).await;
            return Ok(());
        } else if expected_nonce == msg_nonce {
            let max_fee = nonced_verification_data.max_fee;
            if max_fee > user_state.min_fee {
                warn!(
                    "Invalid max fee for address {addr}, had fee {:?} < {:?}",
                    user_state.min_fee, max_fee
                );
                send_message(ws_conn_sink.clone(), ValidityResponseMessage::InvalidMaxFee).await;
                return Ok(());
            }

            user_state.nonce = Some(msg_nonce + U256::one());
            user_state.min_fee = max_fee;
            user_state.proofs_in_batch += 1;

            println!(
                "USER {addr} PROOFS IN BATCH: {}",
                user_state.proofs_in_batch
            );

            println!("ADDING NORMAL ENTRY TO BATCH");
            self.add_to_batch(
                batch_state_lock,
                nonced_verification_data,
                ws_conn_sink.clone(),
                client_msg.signature,
                addr,
            )
            .await;
            println!("NORMAL ENTRY ADDED SUCCESSFULLY TO BATCH");
        } else {
            // In this case, the message might be a replacement one. If it is valid,
            // we replace the old entry with the new from the replacement message.
            if !self
                .handle_replacement_message(
                    batch_state_lock,
                    user_state,
                    nonced_verification_data,
                    ws_conn_sink.clone(),
                    client_msg.signature,
                    addr,
                    expected_nonce,
                )
                .await
            {
                println!("UNLOCKING THE USER {addr} STATEEE!");

                // message should not be added to batch
                return Ok(());
            }

            // println!("LOCKING BATCH STATE - HANDLE MESSAGE - REPLACEMENT MESSAGE");
            // let batch_state = self.batch_state.lock().await;
            // println!("BATCH STATE LOCKED - HANDLE MESSAGE - REPLACEMENT MESSAGE");
            // let updated_min_fee_in_batch = batch_state.get_user_min_fee_in_batch(&addr);
            // user_state.min_fee = updated_min_fee_in_batch;
        }
        println!("BATCH STATE UNLOCKED LOCKED - HANDLE MESSAGE");
        println!("UNLOCKING THE USER {addr} STATEEE!");
        info!("Verification data message handled");
        send_message(ws_conn_sink, ValidityResponseMessage::Valid).await;
        Ok(())
    }

    // Checks user has sufficient balance
    // If user has sufficient balance, increments the user's proof count in the batch
    async fn check_user_balance(&self, addr: &Address, user_proofs_in_batch: usize) -> bool {
        let user_balance = self.get_user_balance(addr).await;
        let min_balance = U256::from(user_proofs_in_batch) * U256::from(MIN_FEE_PER_PROOF);
        if user_balance < min_balance {
            return false;
        }
        true
    }

    /// Handles a replacement message
    /// First checks if the message is already in the batch
    /// If the message is in the batch, checks if the max fee is higher
    /// If the max fee is higher, replaces the message in the batch
    /// If the max fee is lower, sends an error message to the client
    /// If the message is not in the batch, sends an error message to the client
    /// Returns true if the message was replaced in the batch, false otherwise
    async fn handle_replacement_message(
        &self,
        mut batch_state_lock: MutexGuard<'_, BatchState>,
        mut user_state_lock: MutexGuard<'_, UserState>,
        nonced_verification_data: NoncedVerificationData,
        ws_conn_sink: WsMessageSink,
        signature: Signature,
        addr: Address,
        expected_user_nonce: U256,
    ) -> bool {
        // println!("LOCKING BATCH STATE - HANDLE MESSAGE - REPLACEMENT MESSAGE FUNCTION");
        // let mut batch_state = self.batch_state.lock().await;
        // println!("BATCH STATE LOCKED - HANDLE MESSAGE - REPLACEMENT MESSAGE FUNCTION");

        let replacement_max_fee = nonced_verification_data.max_fee;
        let nonce = nonced_verification_data.nonce;

        let Some(entry) = batch_state_lock.get_entry(addr, nonce) else {
            warn!(
                "Invalid nonce for address {addr} Expected: {:?}, got: {:?}",
                expected_user_nonce, nonce
            );
            send_message(ws_conn_sink.clone(), ValidityResponseMessage::InvalidNonce).await;
            return false;
        };

        if entry.nonced_verification_data.max_fee > replacement_max_fee {
            warn!(
                "Invalid replacement message for address {addr}, had fee {:?} < {:?}",
                entry.nonced_verification_data.max_fee, replacement_max_fee
            );
            send_message(
                ws_conn_sink.clone(),
                ValidityResponseMessage::InvalidReplacementMessage,
            )
            .await;

            return false;
        }

        info!(
            "Replacing message for address {} with nonce {} and max fee {}",
            addr, nonce, replacement_max_fee
        );

        // The replacement entry is built from the old entry and validated for then to be replaced
        let mut replacement_entry = entry.clone();
        replacement_entry.signature = signature;
        replacement_entry.verification_data_commitment =
            nonced_verification_data.verification_data.clone().into();
        replacement_entry.nonced_verification_data = nonced_verification_data;

        // Close old sink in old entry replace it with new one
        {
            if let Some(messaging_sink) = replacement_entry.messaging_sink {
                let mut old_sink = messaging_sink.write().await;
                if let Err(e) = old_sink.close().await {
                    // we dont want to exit here, just log the error
                    warn!("Error closing sink: {:?}", e);
                }
            } else {
                warn!(
                    "Old websocket sink was empty. This should only happen in testing environments"
                )
            };
        }

        replacement_entry.messaging_sink = Some(ws_conn_sink.clone());
        if !batch_state_lock.replacement_entry_is_valid(&replacement_entry) {
            warn!("Invalid max fee");
            send_message(
                ws_conn_sink.clone(),
                ValidityResponseMessage::InvalidReplacementMessage,
            )
            .await;
            println!("BATCH STATE UNLOCKED - HANDLE MESSAGE - REPLACEMENT MESSAGE FUNCTION");
            return false;
        }

        info!(
            "Replacement entry is valid, incrementing fee for sender: {:?}, nonce: {:?}, max_fee: {:?}",
            replacement_entry.sender, replacement_entry.nonced_verification_data.nonce, replacement_max_fee
        );

        // remove the old entry and insert the new one
        // note that the entries are considered equal for the priority queue
        // if they have the same nonce and sender, so we can remove the old entry
        // by calling remove with the new entry
        batch_state_lock.batch_queue.remove(&replacement_entry);
        let result = batch_state_lock.batch_queue.push(
            replacement_entry.clone(),
            BatchQueueEntryPriority::new(replacement_max_fee, nonce),
        );
        let updated_min_fee_in_batch = batch_state_lock.get_user_min_fee_in_batch(&addr);
        user_state_lock.min_fee = updated_min_fee_in_batch;

        assert!(result.is_none());

        println!("BATCH STATE UNLOCKED - HANDLE MESSAGE - REPLACEMENT MESSAGE FUNCTION");
        true
    }

    async fn get_user_nonce_from_ethereum(
        &self,
        addr: Address,
    ) -> Result<U256, ContractError<SignerMiddlewareT>> {
        match self.payment_service.user_nonces(addr).call().await {
            Ok(nonce) => Ok(nonce),
            Err(_) => self.payment_service_fallback.user_nonces(addr).call().await,
        }
    }

    /// Adds verification data to the current batch queue.
    async fn add_to_batch(
        &self,
        mut batch_state_lock: MutexGuard<'_, BatchState>,
        verification_data: NoncedVerificationData,
        ws_conn_sink: WsMessageSink,
        proof_submitter_sig: Signature,
        proof_submiter_addr: Address,
    ) {
        info!("Calculating verification data commitments...");
        let verification_data_comm = verification_data.clone().into();
        info!("Adding verification data to batch...");

        let max_fee = verification_data.max_fee;
        let nonce = verification_data.nonce;

        // println!("LOCKING BATCH STATE - ADD TO BATCH FUNCTION");
        // let mut batch_state = self.batch_state.lock().await;
        // println!("BATCH STATE LOCKED - ADD TO BATCH FUNCTION");
        batch_state_lock.batch_queue.push(
            BatchQueueEntry::new(
                verification_data,
                verification_data_comm,
                ws_conn_sink,
                proof_submitter_sig,
                proof_submiter_addr,
            ),
            BatchQueueEntryPriority::new(max_fee, nonce),
        );

        info!(
            "Current batch queue length: {}",
            batch_state_lock.batch_queue.len()
        );
        println!("UNLOCKING BATCH STATE - ADD TO BATCH FUNCTION");
    }

    /// Given a new block number listened from the blockchain, checks if the current batch is ready to be posted.
    /// There are essentially two conditions to be checked:
    ///     * Has the current batch reached the minimum size to be posted?
    ///     * Has the received block number surpassed the maximum interval with respect to the last posted batch block?
    /// Then the batch will be made as big as possible given this two conditions:
    ///     * The serialized batch size needs to be smaller than the maximum batch size
    ///     * The batch submission fee is less than the lowest `max fee` included the batch,
    ///     * And the batch submission fee is more than the highest `max fee` not included the batch.
    /// An extra sanity check is made to check if the batch size is 0, since it does not make sense to post
    /// an empty batch, even if the block interval has been reached.
    /// Once the batch meets the conditions for submission, the finalized batch is then passed to the
    /// `finalize_batch` function.
    async fn is_batch_ready(
        &self,
        block_number: u64,
        gas_price: U256,
    ) -> Option<Vec<BatchQueueEntry>> {
        println!("LOCKING BATCH STATE - IS BATCH READY FUNCTION");
        let mut batch_state = self.batch_state.lock().await;
        println!("BATCH STATE LOCKED - IS BATCH READY FUNCTION");

        let current_batch_len = batch_state.batch_queue.len();
        let last_uploaded_batch_block_lock = self.last_uploaded_batch_block.lock().await;

        // FIXME(marian): This condition should be changed to current_batch_size == 0
        // once the bug in Lambdaworks merkle tree is fixed.
        if current_batch_len < 2 {
            info!("Current batch is empty or length 1. Waiting for more proofs...");
            println!("BATCH STATE UNLOCKED - IS BATCH READY FUNCTION");
            return None;
        }

        if current_batch_len < self.min_batch_len
            && block_number < *last_uploaded_batch_block_lock + self.max_block_interval
        {
            info!(
                "Current batch not ready to be posted. Current block: {} - Last uploaded block: {}. Current batch length: {} - Minimum batch length: {}",
                block_number, *last_uploaded_batch_block_lock, current_batch_len, self.min_batch_len
            );
            println!("BATCH STATE UNLOCKED - IS BATCH READY FUNCTION");
            return None;
        }

        // Check if a batch is currently being posted
        let mut batch_posting = self.posting_batch.lock().await;
        if *batch_posting {
            info!(
                "Batch is currently being posted. Waiting for the current batch to be finalized..."
            );
            println!("BATCH STATE UNLOCKED - IS BATCH READY FUNCTION");
            return None;
        }

        // Set the batch posting flag to true
        *batch_posting = true;

        println!("LOCKING USER STATES - IS BATCH READY FUNCTION");
        let user_states = self.user_states.write().await;
        println!("USER STATES LOCKED - IS BATCH READY FUNCTION");

        let batch_queue_copy = batch_state.batch_queue.clone();
        match batch_queue::try_build_batch(batch_queue_copy, gas_price, self.max_batch_size) {
            Ok((resulting_batch_queue, finalized_batch)) => {
                if finalized_batch.len() == 1 {
                    println!("BATCH STATE UNLOCKED - IS BATCH READY FUNCTION");
                    panic!("FINALIZED BATCH IS LEN 1!!!!!!!!!!!!!!!");
                }
                batch_state.batch_queue = resulting_batch_queue;

                let updated_user_proof_count_and_min_fee =
                    batch_state.get_user_proofs_in_batch_and_min_fee();

                for addr in user_states.keys() {
                    let (proof_count, min_fee) = updated_user_proof_count_and_min_fee
                        .get(addr)
                        .unwrap_or(&(0, U256::MAX));

                    let mut user_state = user_states.get(addr).unwrap().lock().await;
                    user_state.proofs_in_batch = *proof_count;
                    user_state.min_fee = *min_fee;
                }

                // for (addr, (proof_count, min_fee)) in updated_user_proof_count_and_min_fee.iter() {
                //     let mut user_state = user_states.get(addr).unwrap().lock().await;
                //     println!("UPDATING OF PROOFS IN BATCH: {}", proof_count);
                //     user_state.proofs_in_batch = *proof_count;
                //     user_state.min_fee = *min_fee;
                // }

                println!("BATCH STATE UNLOCKED - IS BATCH READY FUNCTION");
                Some(finalized_batch)
            }
            Err(BatcherError::BatchCostTooHigh) => {
                // We can't post a batch since users are not willing to pay the needed fee, wait for more proofs
                info!("No working batch found. Waiting for more proofs...");
                *batch_posting = false;
                println!("BATCH STATE UNLOCKED - IS BATCH READY FUNCTION");
                None
            }
            // FIXME: We should refactor this code and instead of returning None, return an error.
            // See issue https://github.com/yetanotherco/aligned_layer/issues/1046.
            Err(e) => {
                error!("Unexpected error: {:?}", e);
                *batch_posting = false;
                println!("BATCH STATE UNLOCKED - IS BATCH READY FUNCTION");
                None
            }
        }
    }

    /// Takes the finalized batch as input and builds the merkle tree, posts verification data batch
    /// to s3, creates new task in Aligned contract and sends responses to all clients that added proofs
    /// to the batch. The last uploaded batch block is updated once the task is created in Aligned.
    async fn finalize_batch(
        &self,
        block_number: u64,
        finalized_batch: Vec<BatchQueueEntry>,
        gas_price: U256,
    ) -> Result<(), BatcherError> {
        let nonced_batch_verifcation_data: Vec<NoncedVerificationData> = finalized_batch
            .clone()
            .into_iter()
            .map(|entry| entry.nonced_verification_data)
            .collect();

        let batch_verification_data: Vec<VerificationData> = nonced_batch_verifcation_data
            .iter()
            .map(|vd| vd.verification_data.clone())
            .collect();

        let batch_bytes = cbor_serialize(&batch_verification_data)
            .map_err(|e| BatcherError::TaskCreationError(e.to_string()))?;

        info!("Finalizing batch. Length: {}", finalized_batch.len());
        let batch_data_comm: Vec<VerificationDataCommitment> = finalized_batch
            .clone()
            .into_iter()
            .map(|entry| entry.verification_data_commitment)
            .collect();

        let batch_merkle_tree: MerkleTree<VerificationCommitmentBatch> =
            MerkleTree::build(&batch_data_comm);

        {
            let mut last_uploaded_batch_block = self.last_uploaded_batch_block.lock().await;
            // update last uploaded batch block
            *last_uploaded_batch_block = block_number;
            info!(
                "Batch Finalizer: Last uploaded batch block updated to: {}. Lock unlocked",
                block_number
            );
        }

        let leaves: Vec<[u8; 32]> = batch_data_comm
            .iter()
            .map(VerificationCommitmentBatch::hash_data)
            .collect();

        if let Err(e) = self
            .submit_batch(
                &batch_bytes,
                &batch_merkle_tree.root,
                leaves,
                &finalized_batch,
                gas_price,
            )
            .await
        {
            for entry in finalized_batch.into_iter() {
                if let Some(ws_sink) = entry.messaging_sink {
                    let merkle_root = hex::encode(batch_merkle_tree.root);
                    send_message(
                        ws_sink.clone(),
                        ResponseMessage::CreateNewTaskError(merkle_root),
                    )
                    .await
                } else {
                    warn!("Websocket sink was found empty. This should only happen in tests");
                }
            }

            self.flush_queue_and_clear_nonce_cache().await;

            return Err(e);
        };

        connection::send_batch_inclusion_data_responses(finalized_batch, &batch_merkle_tree).await
    }

    async fn flush_queue_and_clear_nonce_cache(&self) {
        warn!("Resetting state... Flushing queue and nonces");

        println!("LOCKING BATCH STATE - FLUE QUEUE AND CLEAR NONCE FUNCTION");
        let mut batch_state = self.batch_state.lock().await;
        println!("BATCH STATE LOCKED - FLUE QUEUE AND CLEAR NONCE FUNCTION");

        for (entry, _) in batch_state.batch_queue.iter() {
            if let Some(ws_sink) = entry.messaging_sink.as_ref() {
                send_message(ws_sink.clone(), ResponseMessage::BatchReset).await;
            } else {
                warn!("Websocket sink was found empty. This should only happen in tests");
            }
        }

        batch_state.batch_queue.clear();
        self.user_states.write().await.clear();
        println!("BATCH STATE UNLOCKED - FLUE QUEUE AND CLEAR NONCE FUNCTION");
    }

    /// Receives new block numbers, checks if conditions are met for submission and
    /// finalizes the batch.
    async fn handle_new_block(&self, block_number: u64) -> Result<(), BatcherError> {
        let gas_price = match self.get_gas_price().await {
            Some(price) => price,
            None => {
                error!("Failed to get gas price");
                return Err(BatcherError::GasPriceError);
            }
        };

        if let Some(finalized_batch) = self.is_batch_ready(block_number, gas_price).await {
            let batch_finalization_result = self
                .finalize_batch(block_number, finalized_batch, gas_price)
                .await;

            // Resetting this here to avoid doing it on every return path of `finalize_batch` function
            let mut batch_posting = self.posting_batch.lock().await;
            *batch_posting = false;

            batch_finalization_result?;
        }

        Ok(())
    }

    /// Post batch to s3 and submit new task to Ethereum
    async fn submit_batch(
        &self,
        batch_bytes: &[u8],
        batch_merkle_root: &[u8; 32],
        leaves: Vec<[u8; 32]>,
        finalized_batch: &[BatchQueueEntry],
        gas_price: U256,
    ) -> Result<(), BatcherError> {
        let signatures: Vec<_> = finalized_batch
            .iter()
            .map(|entry| &entry.signature)
            .cloned()
            .collect();

        let nonces: Vec<_> = finalized_batch
            .iter()
            .map(|entry| entry.nonced_verification_data.nonce)
            .collect();

        let max_fees: Vec<_> = finalized_batch
            .iter()
            .map(|entry| entry.nonced_verification_data.max_fee)
            .collect();

        let s3_client = self.s3_client.clone();
        let batch_merkle_root_hex = hex::encode(batch_merkle_root);
        info!("Batch merkle root: 0x{}", batch_merkle_root_hex);
        let file_name = batch_merkle_root_hex.clone() + ".json";

        info!("Uploading batch to S3...");
        s3::upload_object(
            &s3_client,
            &self.s3_bucket_name,
            batch_bytes.to_vec(),
            &file_name,
        )
        .await
        .map_err(|e| BatcherError::BatchUploadError(e.to_string()))?;

        info!("Batch sent to S3 with name: {}", file_name);

        info!("Uploading batch to contract");
        let batch_data_pointer: String = "".to_owned() + &self.download_endpoint + "/" + &file_name;

        let num_proofs_in_batch = leaves.len();

        let gas_per_proof = (CONSTANT_GAS_COST
            + ADDITIONAL_SUBMISSION_GAS_COST_PER_PROOF * num_proofs_in_batch as u128)
            / num_proofs_in_batch as u128;

        let fee_per_proof = U256::from(gas_per_proof) * gas_price;
        let fee_for_aggregator = (U256::from(AGGREGATOR_GAS_COST)
            * gas_price
            * U256::from(DEFAULT_AGGREGATOR_FEE_MULTIPLIER))
            / U256::from(DEFAULT_AGGREGATOR_FEE_DIVIDER);
        let respond_to_task_fee_limit = (fee_for_aggregator
            * U256::from(RESPOND_TO_TASK_FEE_LIMIT_MULTIPLIER))
            / U256::from(RESPOND_TO_TASK_FEE_LIMIT_DIVIDER);
        let fee_params = CreateNewTaskFeeParams::new(
            fee_for_aggregator,
            fee_per_proof,
            gas_price,
            respond_to_task_fee_limit,
        );

        let signatures = signatures
            .iter()
            .enumerate()
            .map(|(i, signature)| SignatureData::new(signature, nonces[i], max_fees[i]))
            .collect();

        match self
            .create_new_task(
                *batch_merkle_root,
                batch_data_pointer,
                leaves,
                signatures,
                fee_params,
            )
            .await
        {
            Ok(_) => {
                info!("Batch verification task created on Aligned contract");
                Ok(())
            }
            Err(e) => {
                error!(
                    "Failed to send batch to contract, batch will be lost: {:?}",
                    e
                );

                Err(e)
            }
        }
    }

    async fn create_new_task(
        &self,
        batch_merkle_root: [u8; 32],
        batch_data_pointer: String,
        leaves: Vec<[u8; 32]>,
        signatures: Vec<SignatureData>,
        fee_params: CreateNewTaskFeeParams,
    ) -> Result<TransactionReceipt, BatcherError> {
        // pad leaves to next power of 2
        let padded_leaves = Self::pad_leaves(leaves);

        info!("Creating task for: 0x{}", hex::encode(batch_merkle_root));

        match try_create_new_task(
            batch_merkle_root,
            batch_data_pointer.clone(),
            padded_leaves.clone(),
            signatures.clone(),
            fee_params.clone(),
            &self.payment_service,
        )
        .await
        {
            Ok(receipt) => Ok(receipt),
            Err(BatcherSendError::TransactionReverted(err)) => {
                // Since transaction was reverted, we don't want to retry with fallback.
                warn!("Transaction reverted {:?}", err);

                Err(BatcherError::TransactionSendError)
            }
            Err(_) => {
                let receipt = try_create_new_task(
                    batch_merkle_root,
                    batch_data_pointer,
                    padded_leaves,
                    signatures,
                    fee_params,
                    &self.payment_service_fallback,
                )
                .await?;

                Ok(receipt)
            }
        }
    }

    fn pad_leaves(leaves: Vec<[u8; 32]>) -> Vec<[u8; 32]> {
        let leaves_len = leaves.len();
        let last_leaf = leaves[leaves_len - 1];
        leaves
            .into_iter()
            .chain(repeat(last_leaf).take(leaves_len.next_power_of_two() - leaves_len))
            .collect()
    }

    /// Only relevant for testing and for users to easily use Aligned
    fn is_nonpaying(&self, addr: &Address) -> bool {
        self.non_paying_config
            .as_ref()
            .is_some_and(|non_paying_config| non_paying_config.address == *addr)
    }

    /// Only relevant for testing and for users to easily use Aligned in testnet.
    async fn handle_nonpaying_msg(
        &self,
        ws_sink: WsMessageSink,
        client_msg: ClientMessage,
    ) -> Result<(), Error> {
        let Some(non_paying_config) = self.non_paying_config.as_ref() else {
            warn!("There isn't a non-paying configuration loaded. This message will be ignored");
            send_message(ws_sink.clone(), ValidityResponseMessage::InvalidNonce).await;
            return Ok(());
        };

        if client_msg.verification_data.verification_data.proof.len() > self.max_proof_size {
            error!("Proof is too large");
            send_message(ws_sink.clone(), ValidityResponseMessage::ProofTooLarge).await;
            return Ok(());
        }

        let replacement_addr = non_paying_config.replacement.address();
        let replacement_user_balance = self.get_user_balance(&replacement_addr).await;
        if replacement_user_balance == U256::from(0) {
            error!("Insufficient funds for address {:?}", replacement_addr);
            send_message(
                ws_sink.clone(),
                ValidityResponseMessage::InsufficientBalance(replacement_addr),
            )
            .await;
            return Ok(());
        }

        // When pre-verification is enabled, batcher will verify proofs for faster feedback with clients
        if self.pre_verification_is_enabled
            && !zk_utils::verify(&client_msg.verification_data.verification_data).await
        {
            error!("Invalid proof detected. Verification failed.");
            send_message(ws_sink.clone(), ValidityResponseMessage::InvalidProof).await;
            return Ok(()); // Send error message to the client and return
        }

        let user_states = self.user_states.read().await;
        let batch_state_lock = self.batch_state.lock().await;

        // Safe to call `unwrap()` here since at this point we have a non-paying configuration loaded
        // for sure and the non-paying address nonce cached in the user states.
        let non_paying_user_state_lock = user_states.get(&replacement_addr).unwrap();
        let mut non_paying_user_state = non_paying_user_state_lock.lock().await;
        let non_paying_nonce = non_paying_user_state.nonce.unwrap();

        info!("Non-paying nonce: {:?}", non_paying_nonce);

        let nonced_verification_data = NoncedVerificationData::new(
            client_msg.verification_data.verification_data.clone(),
            non_paying_nonce,
            DEFAULT_MAX_FEE_PER_PROOF.into(), // 13_000 gas per proof * 100 gwei gas price (upper bound)
            self.chain_id,
            self.payment_service.address(),
        );

        let client_msg = ClientMessage::new(
            nonced_verification_data.clone(),
            non_paying_config.replacement.clone(),
        )
        .await;

        println!("NONPAYING - ADDING ENTRY TO BATCH");
        self.add_to_batch(
            batch_state_lock,
            nonced_verification_data,
            ws_sink.clone(),
            client_msg.signature,
            non_paying_config.address,
        )
        .await;
        println!("NONPAYING - ADDED ENTRY TO BATCH SUCCESSFULLY");

        // Non-paying user nonce is updated
        (*non_paying_user_state).nonce = Some(non_paying_nonce + U256::one());
        info!("Non-paying verification data message handled");
        send_message(ws_sink, ValidityResponseMessage::Valid).await;

        Ok(())
    }

    async fn get_user_balance(&self, addr: &Address) -> U256 {
        match self.payment_service.user_balances(*addr).call().await {
            Ok(val) => val,
            Err(_) => match self
                .payment_service_fallback
                .user_balances(*addr)
                .call()
                .await
            {
                Ok(balance) => balance,
                Err(_) => {
                    warn!("Failed to get balance for address {:?}", addr);
                    U256::zero()
                }
            },
        }
    }

    async fn user_balance_is_unlocked(&self, addr: &Address) -> bool {
        let unlock_block = match self.payment_service.user_unlock_block(*addr).call().await {
            Ok(val) => val,
            Err(_) => match self
                .payment_service_fallback
                .user_unlock_block(*addr)
                .call()
                .await
            {
                Ok(unlock_block) => unlock_block,
                Err(_) => {
                    warn!("Failed to get unlock block for address {:?}", addr);
                    U256::zero()
                }
            },
        };

        unlock_block != U256::zero()
    }

    async fn get_gas_price(&self) -> Option<U256> {
        match self.eth_ws_provider.get_gas_price().await {
            Ok(gas_price) => Some(gas_price), // this is the block's max priority gas price, not the base fee
            Err(_) => match self.eth_ws_provider_fallback.get_gas_price().await {
                Ok(gas_price) => Some(gas_price),
                Err(_) => {
                    warn!("Failed to get gas price");
                    None
                }
            },
        }
    }
}
