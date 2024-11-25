# Deploy New Implementation

To deploy a new implementation, you can follow the steps below.

## Prerequisites

1. Make sure you have set variables as specified in the [Deploy Contracts Guide](./2_deploy_contracts.md).

2. Set ```MULTISIG=true``` on the ```.env``` used to deploy. ```contracts/scripts/.env.mainnet``` or ```contracts/scripts/.env.holesky``` or ```contracts/scripts/.env.sepolia```

## What contracts can be upgraded?

You can deploy the new implementation of the following contracts:

- AlignedLayerServiceManager
- BatcherPaymentService

## Deploy New Implementation for AlignedLayerServiceManager

1. Deploy the new implementation by running:

   For **Mainnet** deployment:
   
   ```bash
    make upgrade_aligned_contracts NETWORK=mainnet
   ```

    For **Holesky** deployment:
    
    ```bash
     make upgrade_aligned_contracts NETWORK=holesky
    ```
   
    For **Sepolia** deployment:
    
    ```bash
     make upgrade_aligned_contracts NETWORK=sepolia
    ```

   If the new implementation is correctly deployed, the terminal will show the following message:

   ```sh
   The new aligned layer service manager implementation is <new_aligned_layer_service_manager_implementation>
   
   You can propose the upgrade transaction with the multisig using this calldata
    <calldata>
    ```

   You should save this `calldata` for later use.

2. Get the `upgrade` function signature:

   ```sh
   cast sig "upgrade(address, address)"
   ```

   This will show the `upgrade` signature hash: `0x99a88ec4`.

3. Validate the `calldata` by running:

   ```sh
   cast calldata-decode "upgrade(address, address)" <calldata>
   ```

   This will show two addresses. The first one is the `alignedLayerServiceManager` address, and the second one is the new implementation address of `alignedServiceManagerImplementation`.

> [!NOTE]
> Make sure the `alignedLayerServiceManager` address is the same as the one you deployed in the [Deploy Contracts Guide](./2_deploy_contracts.md).
>
> Make sure the `alignedServiceManagerImplementation` address is the same as the one you deployed in this guide.

## Deploy New Implementation for BatcherPaymentService

1. Deploy the new implementation by running:

    For **Mainnet** deployment:
    
    ```bash
    make upgrade_batcher_payment_service NETWORK=mainnet
    ```
   
    For **Holesky** deployment:
     
     ```bash
     make upgrade_batcher_payment_service NETWORK=holesky
     ```
   
    For **Sepolia** deployment:
        
     ```bash  
     make upgrade_batcher_payment_service NETWORK=sepolia
     ```

   If the new implementation is correctly deployed, the script will show the following message:

   `You can propose the upgrade transaction with the multisig using this calldata <calldata>`

   You should save this `calldata` for later use.

2. Get the `upgradeTo` function signature:

   ```bash
    cast sig "upgradeTo(address)"
    ```

   This will show the `upgradeTo` signature hash: `0x3659cfe6`.

3. Validate the `calldata` by running:

   ```shell
   cast calldata-decode "upgradeTo(address)" <calldata>
   ```

   This will show the `batcherPaymentServiceImplementation` address.

> [!NOTE]
> Make sure the `batcherPaymentServiceImplementation` address is the same as the one you deployed in this guide.


## Next Steps

Once you have deployed the new implementation of the contract you want to upgrade, you need to propose the upgrade transaction to the mulstisig, following this [guide](./3_b_2_propose_upgrade.md).

You must also send the contracts addresses and `calldata` you gathered from this guide, to all the multisig participants.
