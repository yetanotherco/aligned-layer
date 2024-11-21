# Deploy New Implementation

To deploy a new implementation, you can follow the steps below.

## Prerequisites

1. Make sure you have set the `.env`, `DEPLOY_CONFIG_PATH` and `BATCHER_PAYMENT_SERVICE_CONFIG_PATH` as specified in the [Deploy Contracts Guide](./2_deploy_contracts.md).

2. Add the following variables to the `.env` file:
    
    ```makefile
    MULTISIG=true
    ```
   
## Deploy New Implementation

You can deploy the new implementation of the following contracts:

- AlignedLayerServiceManager
- BatcherPaymentService
- RegistryCoordinator

### Deploy New Implementation for AlignedLayerServiceManager

1. Deploy the new implementation by running:
          
   ```shell
   make upgrade_aligned_contracts
   ```
   
   If the new implementation is correctly deployed, the script will show the following message:
    
   ```
   The new aligned layer service manager implementation is <new_aligned_layer_service_manager_implementation>
   
   You can propose the upgrade transaction with the multisig using this calldata
    <calldata>
    ```
   
   You should save this `calldata` for later use.
   
2. Get the `upgrade` function signature:

   ```
   cast sig "upgrade(address, address)"
   ```

   This will show the `upgrade` signature hash: `0x99a88ec4`.

3. Validate the `calldata` by running:
   
   ```
   cast calldata-decode "upgrade(address, address)" <calldata>
   ```

   This will show two addresses. The first one is the `alignedLayerServiceManager` address, and the second one is the new implementation address of `alignedServiceManagerImplementation`.

   > [!NOTE]
   > 
   > Make sure the `alignedLayerServiceManager` address is the same as the one you deployed in the [Deploy Contracts Guide](./2_deploy_contracts.md).
   > 
   > Make sure the `alignedServiceManagerImplementation` address is the same as the one you deployed in this guide. 


### Deploy New Implementation for BatcherPaymentService

1. Deploy the new implementation by running:
    
    ```shell
    make upgrade_batcher_payment_service
    ```

    If the new implementation is correctly deployed, the script will show the following message:
    
    ```
    You can propose the upgrade transaction with the multisig using this calldata
    <calldata>
    ```

   You should save this `calldata` for later use.

2. Get the `upgradeTo` function signature:

   ```bash
    cast sig "upgradeTo(address)"
    ```

   This will show the `upgradeTo` signature hash: `0x3659cfe6`.

3. Validate the `calldata` by running:

   ```
   cast calldata-decode "upgradeTo(address)" <calldata>
   ```

   This will show the `batcherPaymentServiceImplementation` address.

   > [!NOTE]
   >
   > Make sure the `batcherPaymentServiceImplementation` address is the same as the one you deployed in this guide.


### Deploy New Implementation for RegistryCoordinator

1. Deploy the new implementation by running:
    
    ```shell
    make upgrade_registry_coordinator
    ```

    If the new implementation is correctly deployed, the script will show the following message:
    
    ```
    You can propose the upgrade transaction with the multisig using this calldata
    <calldata>
    ```

   You should save this `calldata` for later use.

2. Get the `upgrade` function signature:

   ```bash
   cast sig "upgrade(address, address)"
   ```

   This will show the `upgrade` signature hash: `0x99a88ec4`.

3. Validate the `calldata` by running:

   ```
   cast calldata-decode "upgrade(address, address)" <calldata>
   ```

   This will show two addresses. The first one is the `registryCoordinator` address, and the second one is the new implementation address of `registryCoordinatorImplementation`.

   > [!NOTE]
   >
   > Make sure the `registryCoordinator` address is the same as the one you deployed in the [Deploy Contracts Guide](./2_deploy_contracts.md).
   >
   > Make sure the `registryCoordinatorImplementation` address is the same as the one you deployed in this guide.

   
## Next Steps

Once you have deployed the new implementation of the contract you want to upgrade, you need to propose the upgrade transaction following this [guide](./3_b_2_propose_upgrade.md).

Send the contracts addresses and `calldata` to the multisig participants to propose the upgrade transaction.
