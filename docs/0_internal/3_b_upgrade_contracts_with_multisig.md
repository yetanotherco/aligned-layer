# Upgrade Contracts using a Multisig

> [!WARNING]  
> Safe Multisig Wallet is not currently supported in Holesky Testnet.
> For this reason, we deployed EigenLayer contracts in Sepolia to test the upgrade on AlignedLayer Contracts.

> [!NOTE]
> EigenLayer Sepolia contracts information is available in `contracts/script/output/sepolia/eigenlayer_deployment_output.json`.

> [!NOTE]
> You can find a guide on how to Deploy the contracts [here](./2_deploy_contracts.md).


This guide is for upgrading contracts deployed using a Multisig wallet. If you deployed the contract without a Multisig wallet, you can follow the [Upgrade Contracts](./3_a_upgrade_contracts.md) tutorial.

In this guide we make an upgrade of Aligned Layer Service Manager contract using a multisig wallet. This is important to ensure not one party can do an upgrade, and helps ensure the team is not locked out of upgrading the network due to a loss of keys.

## Prerequisites

- To upgrade any of the contracts, you need to have set the `.env`, `DEPLOY_CONFIG_PATH` and `BATCHER_PAYMENT_SERVICE_CONFIG_PATH`

- You need to have installed git and make.

- Clone the repository
   ```
   git clone https://github.com/yetanotherco/aligned_layer.git
   ```

- Install foundry
    ```shell
    make install_foundry
    foundryup -v nightly-a428ba6ad8856611339a6319290aade3347d25d9
    ```

## Deploy the new implementation

The first step is to deploy the new implementation of the contract. 

You need to 
