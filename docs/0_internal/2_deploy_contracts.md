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

## Eigenlayer Contracts

To deploy Aligned contracts, you need previously deployed EigenLayer contracts.

These contracts are not deployed by Aligned. These are the current EigenLayer contracts:

- [Holesky Contracts](https://github.com/Layr-Labs/eigenlayer-contracts/blob/testnet-holesky/script/configs/holesky/eigenlayer_addresses_testnet.config.json)
- [Mainnet Contracts](https://github.com/Layr-Labs/eigenlayer-contracts/blob/mainnet/script/configs/mainnet/mainnet-addresses.config.json)

## Aligned Contracts

This section will guide you through the deployment of the Aligned Layer contracts.

After finishing the deployment, you will have the deployed contract addresses.

### Set .env variables

To deploy the AlignedLayer contracts, you will need to set environment variables in a `.env` file in the same
directory as the deployment script (`contracts/scripts/`).

The necessary environment variables are:

| Variable Name                   | Description                                                           | Holesky                                                     | Mainnet                             |
|---------------------------------|-----------------------------------------------------------------------|-------------------------------------------------------------|-------------------------------------|
| `RPC_URL`                       | The RPC URL of the network you want to deploy to.                     | https://ethereum-holesky-rpc.publicnode.com                 | https://ethereum-rpc.publicnode.com |
| `PRIVATE_KEY`                   | The private key of the account you want to deploy the contracts with. | <your_private_key>                                          | <your_private_key>                  |
| `EXISTING_DEPLOYMENT_INFO_PATH` | The path to the file containing the deployment info about EigenLayer. | ./script/output/holesky/eigenlayer_deployment_output.json   | TBD                                 |
| `DEPLOY_CONFIG_PATH`            | The path to the deployment config file for the Service Manager.       | ./script/deploy/config/holesky/aligned.holesky.config.json  | TBD                                 |
| `OUTPUT_PATH`                   | The path to the file where the deployment info will be saved.         | ./script/output/holesky/alignedlayer_deployment_output.json | TBD                                 |
| `ETHERSCAN_API_KEY`             | API KEY to verify the contracts in Etherscan.                         | <your_etherscan_api_key>                                    | <your_etherscan_api_key>            |

You can find an example `.env` file in [.env.example.holesky](../../contracts/scripts/.env.example.holesky)

> [!WARNING]
> All file paths must be inside the `script/` directory, as shown in `.env.example.holesky` because of `foundry`'s permissions to read and write files.

### Set DEPLOY_CONFIG_PATH file

You need to complete the `DEPLOY_CONFIG_PATH` file with the following information:

```json
{
  "chainInfo": {
    "chainId": "<chain_id>"
  },
  "permissions": {
    "owner": "<owner_address>",
    "aggregator": "<aggregator_address>",
    "upgrader": "<upgrader_address>",
    "churner": "<churner_address>",
    "ejector": "<ejector_address>",
    "deployer": "<deployer_address>",
    "initalPausedStatus": 0
  },
  "minimumStakes": [],
  "strategyWeights": [],
  "operatorSetParams": [],
  "uri": ""
}
```

You can find an example config file in [aligned.holesky.config.json](../../contracts/script/deploy/config/holesky/aligned.holesky.config.json).

> [!NOTE]
> 
> You can find the list of Mainnet strategies for the `strategyWeights` field [here](https://github.com/Layr-Labs/eigenlayer-contracts?tab=readme-ov-file#current-mainnet-deployment)
> 
> You can find the list of Holesky strategies for the `strategyWeights` field [here](https://github.com/Layr-Labs/eigenlayer-contracts?tab=readme-ov-file#current-testnet-deployment)

#### Multisig configuration

If you are using a Multisig for the contracts management (like upgrades or pauses), you need to set the Multisig address in the `permissions` sections.

For example, if you are using a Multisig for the `upgrader` permission, you need to set the Multisig address in the `upgrader` field.

### Deploy the contracts

Once you have configured the `.env` and `DEPLOY_CONFIG_PATH` files, you can run the following command:

```bash
make deploy_aligned_contracts
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
    "stakeRegistryImplementation": "<stake_registry_implementation_address>"
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

## Batcher Payments Service Contracts

This section will guide you through the deployment of the Aligned Layer Batcher Payment Service contract.

After finishing the deployment, you will have the deployed contract addresses.

### Set .env variables

To deploy the Batcher Payment Service contract, you will need to set environment variables in a `.env` file in the same
directory as the deployment script (`contracts/scripts/`).

The necessary environment variables are:

| Variable Name                         | Description                                                             | Holesky                                                                    | Mainnet                             |
|---------------------------------------|-------------------------------------------------------------------------|----------------------------------------------------------------------------|-------------------------------------|
| `RPC_URL`                             | The RPC URL of the network you want to deploy to.                       | https://ethereum-holesky-rpc.publicnode.com                                | https://ethereum-rpc.publicnode.com |
| `PRIVATE_KEY`                         | The private key of the account you want to deploy the contracts with.   | <your_private_key>                                                         | <your_private_key>                  |
| `EXISTING_DEPLOYMENT_INFO_PATH`       | The path to the file containing the deployment info about EigenLayer.   | ./script/output/holesky/eigenlayer_deployment_output.json                  | TBD                                 |
| `DEPLOY_CONFIG_PATH`                  | The path to the deployment config file for the Service Manager.         | ./script/deploy/config/holesky/batcher_payment_service.holesky.config.json | TBD                                 |
| `BATCHER_PAYMENT_SERVICE_CONFIG_PATH` | The path to the deployment config file for the Batcher Payment Service. | ./script/deploy/config/holesky/batcher-payment-service.holesky.config.json | TBD                                 |
| `OUTPUT_PATH`                         | The path to the file where the deployment info will be saved.           | ./script/output/holesky/alignedlayer_deployment_output.json                | TBD                                 |
| `ETHERSCAN_API_KEY`                   | API KEY to verify the contracts in Etherscan.                           | <your_etherscan_api_key>                                                   | <your_etherscan_api_key>            |

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

