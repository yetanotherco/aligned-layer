# Deploying Aligned Contracts

This guide will walk you through the deployment of the Aligned Layer contracts.

Also, you will be able to deploy the Batcher Payment Service contract.

## Prerequisites

- You need to have installed `git` and `make`.

- Clone the repository
   ```
   git clone https://github.com/yetanotherco/aligned_layer.git
   ```

- Install foundry
    ```shell
    make install_foundry
    foundryup -v nightly-a428ba6ad8856611339a6319290aade3347d25d9
    ```
  
## AlignedServiceManager Contracts

1 You should have a keystore for the Aggregator, with the seedphrase saved on paper
1. If you don't have an API Key on Etherscan, create one following this guide  [guide](https://docs.etherscan.io/getting-started/viewing-api-usage-statistics).
2. If you don't have it, create a wallet to deploy Aligned with ```cast wallet new```. and keep the ```Address``` and ```Private Key```values
3. Transfer around 35M gas to the ```Address```
4. ```CD``` into the Aligned repo
5. Set the ```PRIVATE_KEY``` and ```ETHERSCAN_API_KEY```on the ```.env``` used to deploy. ```contracts/scripts/.env.mainnet``` or ```contracts/scripts/.env.holesky``` or ```contracts/scripts/.env.sepolia```
6. Set `owner`, `upgrader`, `churner`, `ejector` and `deployer` values of ```contracts/script/deploy/config/mainnet/aligned.holesky.config.json``` to the multisig controlling this processes, or the owner of the deploying address, depending on what you want

You need to set the following address in the config file:

- `aggregator`: Address that will be the aggregator of Aligned.



- `owner`: Address that will be the owner of the contracts. This address can be a multisig or a normal wallet.
- `upgrader`: Address that will be the upgrader of the contracts. This address can be a multisig or a normal wallet.
- `churner`: Address that will be the churner of the contracts. This address can be a multisig or a normal wallet.
- `ejector`: Address that will be the ejector of the contracts. This address can be a multisig or a normal wallet.
- `deployer`: Address that will be the deployer of the contracts. This address can be a multisig or a normal wallet.

For **Mainnet** deployment, complete the following [config file](../../contracts/script/deploy/config/mainnet/aligned.mainnet.config.json).

For **Holesky** deployment, complete the following [config file](../../contracts/script/deploy/config/holesky/aligned.holesky.config.json).

For **Sepolia** deployment, complete the following [config file](../../contracts/script/deploy/config/sepolia/aligned.sepolia.config.json).

> [!NOTE]
> 
> Strategies are already filled in the config files. Also, you can find them in the following links
> 
> You can find the list of Mainnet strategies for the `strategyWeights` field [here](https://github.com/Layr-Labs/eigenlayer-contracts?tab=readme-ov-file#current-mainnet-deployment)
> 
> You can find the list of Holesky strategies for the `strategyWeights` field [here](https://github.com/Layr-Labs/eigenlayer-contracts?tab=readme-ov-file#current-testnet-deployment)

### Deploy the contracts

Once you have configured the `.env` and `config` files, you can run the following command:

For **Mainnet** deployment:

```bash
make deploy_mainnet_aligned_contracts
```

For **Holesky** deployment:

```bash
make deploy_holesky_aligned_contracts
```

For **Sepolia** deployment:

```bash
make deploy_sepolia_aligned_contracts
```

If the deployment is correct, you will find the deployment information here:

- [Mainnet deployment](../../contracts/script/output/mainnet/alignedlayer_deployment_output.json)
- [Holesky deployment](../../contracts/script/output/holesky/alignedlayer_deployment_output.json)
- [Sepolia deployment](../../contracts/script/output/sepolia/alignedlayer_deployment_output.json)


## Batcher Payments Service Contracts

This section will guide you through the deployment of the Aligned Layer Batcher Payment Service contract.

After finishing the deployment, you will have the deployed contract addresses.

### Set .env variables

To deploy the Batcher Payment Service contract, you will need to set environment variables in a `.env` file in the same
directory as the deployment script (`contracts/scripts/`).

The necessary environment variables are:

| Variable Name                         | Description                                                             | Sepolia                                                                    | Holesky                                                                    | Mainnet                             |
|---------------------------------------|-------------------------------------------------------------------------|----------------------------------------------------------------------------|----------------------------------------------------------------------------|-------------------------------------|
| `RPC_URL`                             | The RPC URL of the network you want to deploy to.                       | https://ethereum-sepolia-rpc.publicnode.com                                | https://ethereum-holesky-rpc.publicnode.com                                | https://ethereum-rpc.publicnode.com |
| `PRIVATE_KEY`                         | The private key of the account you want to deploy the contracts with.   | <your_private_key>                                                         | <your_private_key>                                                         | <your_private_key>                  |
| `EXISTING_DEPLOYMENT_INFO_PATH`       | The path to the file containing the deployment info about EigenLayer.   | ./script/output/sepolia/eigenlayer_deployment_output.json                  | ./script/output/holesky/eigenlayer_deployment_output.json                  | TBD                                 |
| `DEPLOY_CONFIG_PATH`                  | The path to the deployment config file for the Service Manager.         | ./script/deploy/config/sepolia/aligned.sepolia.config.json                 | ./script/deploy/config/holesky/batcher_payment_service.holesky.config.json | TBD                                 |
| `BATCHER_PAYMENT_SERVICE_CONFIG_PATH` | The path to the deployment config file for the Batcher Payment Service. | ./script/deploy/config/sepolia/batcher-payment-service.sepolia.config.json | ./script/deploy/config/holesky/batcher-payment-service.holesky.config.json | TBD                                 |
| `OUTPUT_PATH`                         | The path to the file where the deployment info will be saved.           | ./script/output/sepolia/alignedlayer_deployment_output.json                | ./script/output/holesky/alignedlayer_deployment_output.json                | TBD                                 |
| `ETHERSCAN_API_KEY`                   | API KEY to verify the contracts in Etherscan.                           | <your_etherscan_api_key>                                                   | <your_etherscan_api_key>                                                   | <your_etherscan_api_key>            |

You can find an example `.env` file in [.env.example.holesky](../../contracts/scripts/.env.example.holesky)

### Set BATCHER_PAYMENT_SERVICE_CONFIG_PATH file

You need to complete the `BATCHER_PAYMENT_SERVICE_CONFIG_PATH` file with the following information:

```json
{
  "address": {
    "batcherWallet": "<batcher_wallet_address>",
    "alignedLayerServiceManager": "<aligned_layer_service_manager_address>"
  },
  "permissions": {
    "owner": "<owner_address>"
  },
  "eip712": {
    "noncedVerificationDataTypeHash": "0x41817b5c5b0c3dcda70ccb43ba175fdcd7e586f9e0484422a2c6bba678fdf4a3"
  }
}
```

#### Multisig configuration

If you are using a Multisig for the contracts management (like upgrades or pauses), you need to set the Multisig address in the `permissions` sections.

For the batcher payment service, you can set the Multisig address in the `owner` field. This will allow the Multisig to upgrade and pause the contract with the Multisig.

### Deploy the contracts

Once you have configured the `.env` and `BATCHER_PAYMENT_SERVICE_CONFIG_PATH` files, you can run the following command:

```bash
make deploy_batcher_payment_service
```

Once the contracts are deployed, you will see the following output at `OUTPUT_PATH` file:

```json
{
  "addresses": {
    "alignedLayerProxyAdmin": "<aligned_layer_proxy_admin_address>",
    "alignedLayerServiceManager": "<aligned_layer_service_manager_address>",
    "alignedLayerServiceManagerImplementation": "<aligned_layer_service_manager_implementation_address>",
    "blsApkRegistry": "<bls_apk_registry_address>",
    "blsApkRegistryImplementation": "<bls_apk_registry_implementation_address>",
    "indexRegistry": "<index_registry_address>",
    "indexRegistryImplementation": "<index_registry_implementation_address>",
    "operatorStateRetriever": "<operator_state_retriever_address>",
    "pauserRegistry": "<pauser_registry_address>",
    "registryCoordinator": "<registry_coordinator_address>",
    "registryCoordinatorImplementation": "<registry_coordinator_implementation_address>",
    "serviceManagerRouter": "<service_manager_router_address>",
    "stakeRegistry": "<stake_registry_address>",
    "stakeRegistryImplementation": "<stake_registry_implementation_address>",
    "batcherPaymentService": "<batcher_payment_service_address>",
    "batcherPaymentServiceImplementation": "<batcher_payment_service_implementation_address>"
  },
  "chainInfo": {
    "chainId": 17000,
    "deploymentBlock": 1628199
  },
  "permissions": {
    "alignedLayerAggregator": "<aligned_layer_aggregator_address>",
    "alignedLayerChurner": "<aligned_layer_churner_address>",
    "alignedLayerEjector": "<aligned_layer_ejector_address>",
    "alignedLayerOwner": "<aligned_layer_owner_address>",
    "alignedLayerPauser": "<aligned_layer_pauser_address>",
    "alignedLayerUpgrader": "<aligned_layer_upgrader_address>"
  }
}
```

