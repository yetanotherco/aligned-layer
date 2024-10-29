use std::str::FromStr;
use std::sync::Arc;

use crate::config::ECDSAConfig;
use crate::retry::batcher_retryables::{get_current_nonce_retryable, get_gas_price_retryable};
use crate::retry::retry_function;
use crate::types::errors::BatcherError;
use aligned_sdk::core::constants::{
    DEFAULT_BACKOFF_FACTOR, DEFAULT_MAX_RETRIES, DEFAULT_MIN_RETRY_DELAY,
    OVERRIDE_GAS_PRICE_MULTIPLIER, PERCENTAGE_DIVIDER,
};
use ethers::prelude::*;
use ethers::providers::{Http, Provider};
use log::error;

use super::payment_service::SignerMiddlewareT;

pub fn get_provider(eth_rpc_url: String) -> Result<Provider<Http>, anyhow::Error> {
    let provider = Http::from_str(eth_rpc_url.as_str())
        .map_err(|e| anyhow::Error::msg(format!("Failed to create provider: {}", e)))?;
    Ok(Provider::new(provider))
}

pub async fn get_batcher_signer(
    provider: Provider<Http>,
    ecdsa_config: ECDSAConfig,
) -> anyhow::Result<Arc<SignerMiddlewareT>> {
    let chain_id = provider.get_chainid().await?;

    // get private key from keystore
    let wallet = Wallet::decrypt_keystore(
        &ecdsa_config.private_key_store_path,
        &ecdsa_config.private_key_store_password,
    )?
    .with_chain_id(chain_id.as_u64());

    let signer = Arc::new(SignerMiddleware::new(provider, wallet));
    Ok(signer)
}

/// Calculates an increased gas price for retrying a transaction override.
/// The gas price rises with each retry by applying a multiplier based on the iteration count.
pub fn get_bumped_gas_price(
    previous_gas_price: U256,
    current_gas_price: U256,
    iteration: usize,
) -> U256 {
    let override_gas_multiplier = U256::from(OVERRIDE_GAS_PRICE_MULTIPLIER) + (5 * iteration);
    let bumped_previous_gas_price =
        previous_gas_price * override_gas_multiplier / U256::from(PERCENTAGE_DIVIDER);

    let bumped_current_gas_price =
        current_gas_price * override_gas_multiplier / U256::from(PERCENTAGE_DIVIDER);
    bumped_current_gas_price.max(bumped_previous_gas_price)
}

pub async fn get_current_nonce(
    eth_http_provider: &Provider<Http>,
    eth_http_provider_fallback: &Provider<Http>,
    addr: H160,
) -> Result<U256, String> {
    retry_function(
        || get_current_nonce_retryable(eth_http_provider, eth_http_provider_fallback, addr),
        DEFAULT_MIN_RETRY_DELAY,
        DEFAULT_BACKOFF_FACTOR,
        DEFAULT_MAX_RETRIES,
    )
    .await
    .map_err(|e| {
        error!("Could't get nonce: {:?}", e);
        e.to_string()
    })
}

/// Gets the current gas price from Ethereum using exponential backoff.
pub async fn get_gas_price(
    eth_http_provider: &Provider<Http>,
    eth_http_provider_fallback: &Provider<Http>,
) -> Result<U256, BatcherError> {
    retry_function(
        || get_gas_price_retryable(eth_http_provider, eth_http_provider_fallback),
        DEFAULT_MIN_RETRY_DELAY,
        DEFAULT_BACKOFF_FACTOR,
        DEFAULT_MAX_RETRIES,
    )
    .await
    .map_err(|e| {
        error!("Could't get gas price: {e}");
        BatcherError::GasPriceError
    })
}
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_bumped_gas_price_initial_iteration() {
        let previous_gas_price = U256::from(1000);
        let current_gas_price = U256::from(1200);
        let iteration = 0;

        let expected = U256::from(1440); // (1200 * (120 + 0)) / 100

        assert_eq!(
            get_bumped_gas_price(previous_gas_price, current_gas_price, iteration),
            expected
        );
    }

    #[test]
    fn test_get_bumped_gas_price_with_iteration() {
        let previous_gas_price = U256::from(1000);
        let current_gas_price = U256::from(1200);
        let iteration = 2;

        let expected = U256::from(1560); // (1200 * (120 + 10) / 100

        assert_eq!(
            get_bumped_gas_price(previous_gas_price, current_gas_price, iteration),
            expected
        );
    }

    #[test]
    fn test_get_bumped_gas_price_previous_higher() {
        let previous_gas_price = U256::from(1500);
        let current_gas_price = U256::from(1200);
        let iteration = 1;

        let expected = U256::from(1875); // (1500 * (120 + 5) / 100

        assert_eq!(
            get_bumped_gas_price(previous_gas_price, current_gas_price, iteration),
            expected
        );
    }
}
