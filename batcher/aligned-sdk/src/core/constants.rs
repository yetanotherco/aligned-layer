/// Batcher ///
pub const GAS_PRICE_INCREMENT_PERCENTAGE_PER_ITERATION: usize = 5;
pub const AGGREGATOR_GAS_COST: u128 = 400_000;
pub const BATCHER_SUBMISSION_BASE_GAS_COST: u128 = 125_000;
pub const ADDITIONAL_SUBMISSION_GAS_COST_PER_PROOF: u128 = 13_000;
pub const CONSTANT_GAS_COST: u128 =
    ((AGGREGATOR_GAS_COST * DEFAULT_AGGREGATOR_FEE_PERCENTAGE_MULTIPLIER) / PERCENTAGE_DIVIDER)
        + BATCHER_SUBMISSION_BASE_GAS_COST;
pub const DEFAULT_MAX_FEE_PER_PROOF: u128 =
    ADDITIONAL_SUBMISSION_GAS_COST_PER_PROOF * 100_000_000_000; // gas_price = 100 Gwei = 0.0000001 ether (high gas price)
pub const CONNECTION_TIMEOUT: u64 = 30; // 30 secs

// % modifiers: (100% is x1, 10% is x0.1, 1000% is x10)
pub const RESPOND_TO_TASK_FEE_LIMIT_PERCENTAGE_MULTIPLIER: u128 = 250; // fee_for_aggregator -> respondToTaskFeeLimit modifier
pub const DEFAULT_AGGREGATOR_FEE_PERCENTAGE_MULTIPLIER: u128 = 150; // feeForAggregator modifier
pub const GAS_PRICE_PERCENTAGE_MULTIPLIER: u128 = 110; // gasPrice modifier
pub const OVERRIDE_GAS_PRICE_PERCENTAGE_MULTIPLIER: u128 = 120; // gasPrice modifier to override previous transactions
pub const PERCENTAGE_DIVIDER: u128 = 100;

/// SDK ///
/// Number of proofs we a batch for estimation.
/// This is the number of proofs in a batch of size n, where we set n = 32.
/// i.e. the user pays for the entire batch and his proof is instantly submitted.
pub const MAX_FEE_BATCH_PROOF_NUMBER: usize = 32;
/// Estimated number of proofs for batch submission.
/// This corresponds to the number of proofs to compute for a default max_fee.
pub const MAX_FEE_DEFAULT_PROOF_NUMBER: usize = 10;

/// Ethereum calls retry constants
pub const ETHEREUM_CALL_MIN_RETRY_DELAY: u64 = 500; // milliseconds
pub const ETHEREUM_CALL_MAX_RETRIES: usize = 5;
pub const ETHEREUM_CALL_BACKOFF_FACTOR: f32 = 2.0;
pub const ETHEREUM_CALL_MAX_RETRY_DELAY: u64 = 3600; // seconds

/// Ethereum transaction retry constants
pub const BUMP_MIN_RETRY_DELAY: u64 = 500; // milliseconds
pub const BUMP_MAX_RETRIES: usize = 33; // ~ 1 day
pub const BUMP_BACKOFF_FACTOR: f32 = 2.0;
pub const BUMP_MAX_RETRY_DELAY: u64 = 3600; // seconds

/// NETWORK ADDRESSES ///
/// BatcherPaymentService
pub const BATCHER_PAYMENT_SERVICE_ADDRESS_DEVNET: &str =
    "0x7bc06c482DEAd17c0e297aFbC32f6e63d3846650";
pub const BATCHER_PAYMENT_SERVICE_ADDRESS_HOLESKY: &str =
    "0x815aeCA64a974297942D2Bbf034ABEe22a38A003";
pub const BATCHER_PAYMENT_SERVICE_ADDRESS_HOLESKY_STAGE: &str =
    "0x7577Ec4ccC1E6C529162ec8019A49C13F6DAd98b";
pub const BATCHER_PAYMENT_SERVICE_ADDRESS_MAINNET: &str =
    "0xb0567184A52cB40956df6333510d6eF35B89C8de";
/// AlignedServiceManager
pub const ALIGNED_SERVICE_MANAGER_DEVNET: &str = "0x851356ae760d987E095750cCeb3bC6014560891C";
pub const ALIGNED_SERVICE_MANAGER_HOLESKY: &str = "0x58F280BeBE9B34c9939C3C39e0890C81f163B623";
pub const ALIGNED_SERVICE_MANAGER_HOLESKY_STAGE: &str =
    "0x9C5231FC88059C086Ea95712d105A2026048c39B";
pub const ALIGNED_SERVICE_MANAGER_MAINNET: &str = "0xeF2A435e5EE44B2041100EF8cbC8ae035166606c";
/// Batcher URL's
pub const BATCHER_URL_DEVNET: &str = "ws://localhost:8080";
pub const BATCHER_URL_HOLESKY: &str = "wss://holesky.batcher.alignedlayer.com";
pub const BATCHER_URL_HOLESKY_STAGE: &str = "wss://stage.batcher.alignedlayer.com";
pub const BATCHER_URL_MAINNET: &str = "wss://batcher.alignedlayer.com";
