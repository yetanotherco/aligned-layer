use std::iter::repeat;
use std::str::FromStr;
use std::sync::Arc;
use std::time::Duration;

use aligned_sdk::eth::batcher_payment_service::{BatcherPaymentServiceContract, SignatureData};
use ethers::prelude::k256::ecdsa::SigningKey;
use ethers::prelude::*;
use log::{error, info, warn};
use tokio::time::sleep;

const CREATE_NEW_TASK_MAX_RETRIES: usize = 15;
const CREATE_NEW_TASK_MILLISECS_BETWEEN_RETRIES: u64 = 2000;

use crate::{config::ECDSAConfig, types::errors::BatcherError};

#[derive(Debug, Clone, EthEvent)]
pub struct BatchVerified {
    pub batch_merkle_root: [u8; 32],
}

pub type BatcherPaymentService =
    BatcherPaymentServiceContract<SignerMiddleware<Provider<Http>, Wallet<SigningKey>>>;

pub fn get_provider(eth_rpc_url: String) -> Result<Provider<Http>, anyhow::Error> {
    Provider::<Http>::try_from(eth_rpc_url).map_err(|err| anyhow::anyhow!(err))
}

pub async fn create_new_task(
    payment_service: &BatcherPaymentService,
    batch_merkle_root: [u8; 32],
    batch_data_pointer: String,
    leaves: Vec<[u8; 32]>,
    signatures: Vec<SignatureData>,
    gas_for_aggregator: U256,
    gas_per_proof: U256,
) -> Result<TransactionReceipt, BatcherError> {
    // pad leaves to next power of 2
    let padded_leaves = pad_leaves(leaves);

    let call = payment_service.create_new_task(
        batch_merkle_root,
        batch_data_pointer,
        padded_leaves,
        signatures,
        gas_for_aggregator,
        gas_per_proof,
    );

    // If there was a pending transaction from a previously sent batch, the `call.send()` will
    // fail because of the nonce not being updated. We should retry sending and not returning an error
    // immediatly.
    info!("Creating task for: {:x?}", batch_merkle_root);

    for i in 0..CREATE_NEW_TASK_MAX_RETRIES {
        match call.send().await {
            Ok(pending_tx) => match pending_tx.await {
                Ok(Some(receipt)) => return Ok(receipt),
                Ok(None) => return Err(BatcherError::ReceiptNotFoundError),
                Err(_) => return Err(BatcherError::TransactionSendError),
            },
            Err(error) => {
                if i != CREATE_NEW_TASK_MAX_RETRIES - 1 {
                    warn!(
                        "Error when trying to create a task: {}\n Retrying ...",
                        error
                    );
                } else {
                    error!("Error when trying to create a task on last retry. Batch task {:x?} will be lost", batch_merkle_root);
                    return Err(BatcherError::TaskCreationError(error.to_string()));
                }
            }
        };

        sleep(Duration::from_millis(
            CREATE_NEW_TASK_MILLISECS_BETWEEN_RETRIES,
        ))
        .await;
    }

    Err(BatcherError::MaxRetriesReachedError)
}

pub async fn get_batcher_payment_service(
    provider: Provider<Http>,
    ecdsa_config: ECDSAConfig,
    contract_address: String,
) -> Result<BatcherPaymentService, anyhow::Error> {
    let chain_id = provider.get_chainid().await?;

    // get private key from keystore
    let wallet = Wallet::decrypt_keystore(
        &ecdsa_config.private_key_store_path,
        &ecdsa_config.private_key_store_password,
    )?
    .with_chain_id(chain_id.as_u64());

    let signer = Arc::new(SignerMiddleware::new(provider, wallet));

    let service_manager =
        BatcherPaymentService::new(H160::from_str(contract_address.as_str())?, signer);

    Ok(service_manager)
}

fn pad_leaves(leaves: Vec<[u8; 32]>) -> Vec<[u8; 32]> {
    let leaves_len = leaves.len();
    let last_leaf = leaves[leaves_len - 1];
    leaves
        .into_iter()
        .chain(repeat(last_leaf).take(leaves_len.next_power_of_two() - leaves_len))
        .collect()
}
