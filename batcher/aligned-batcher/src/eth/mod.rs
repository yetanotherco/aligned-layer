use std::str::FromStr;
use std::sync::Arc;

use ethers::prelude::k256::ecdsa::SigningKey;
use ethers::prelude::*;
use stream::EventStream;

use crate::config::ECDSAConfig;

abigen!(
    AlignedLayerServiceManagerContract,
    "./src/eth/abi/AlignedLayerServiceManager.json",
);

abigen!(
    BatcherPaymentServiceContract,
    "./src/eth/abi/BatcherPaymentService.json",
);

#[derive(Debug, Clone, EthEvent)]
pub struct BatchVerified {
    pub batch_merkle_root: [u8; 32],
}

pub type AlignedLayerServiceManager =
    AlignedLayerServiceManagerContract<SignerMiddleware<Provider<Http>, Wallet<SigningKey>>>;

pub type BatchVerifiedEventStream<'s> = EventStream<
    's,
    FilterWatcher<'s, Http, Log>,
    BatchVerifiedFilter,
    ContractError<SignerMiddleware<Provider<Http>, Wallet<SigningKey>>>,
>;

pub type BatcherPaymentService =
    BatcherPaymentServiceContract<SignerMiddleware<Provider<Http>, Wallet<SigningKey>>>;

pub fn get_provider(eth_rpc_url: String) -> Result<Provider<Http>, anyhow::Error> {
    Provider::<Http>::try_from(eth_rpc_url).map_err(|err| anyhow::anyhow!(err))
}

pub async fn get_service_manager(
    provider: Provider<Http>,
    ecdsa_config: ECDSAConfig,
    contract_address: String,
) -> Result<AlignedLayerServiceManager, anyhow::Error> {
    let chain_id = provider.get_chainid().await?;

    // get private key from keystore
    let wallet = Wallet::decrypt_keystore(
        &ecdsa_config.private_key_store_path,
        &ecdsa_config.private_key_store_password,
    )?
    .with_chain_id(chain_id.as_u64());

    let signer = Arc::new(SignerMiddleware::new(provider, wallet));

    let service_manager =
        AlignedLayerServiceManager::new(H160::from_str(contract_address.as_str())?, signer);

    Ok(service_manager)
}

pub async fn create_new_task(
    payment_service: &BatcherPaymentService,
    batch_merkle_root: [u8; 32],
    batch_data_pointer: String,
    leaves: Vec<[u8; 32]>,
    signatures: Vec<Signature>,
    gas_for_aggregator: U256,
    gas_per_proof: U256,
) -> Result<TransactionReceipt, anyhow::Error> {
    let signatures = signatures
        .iter()
        .map(|s| SignatureData {
            v: s.v as u8,
            r: s.r.into(),
            s: s.s.into(),
        })
        .collect::<Vec<SignatureData>>();

    let call = payment_service.create_new_task(
        batch_merkle_root,
        batch_data_pointer,
        leaves,
        signatures,
        gas_for_aggregator,
        gas_per_proof,
    );
    let pending_tx = call.send().await?;

    match pending_tx.await? {
        Some(receipt) => Ok(receipt),
        None => Err(anyhow::anyhow!("Receipt not found")),
    }
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
