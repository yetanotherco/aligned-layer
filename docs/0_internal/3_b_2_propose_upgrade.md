# Propose the Transaction for Upgrade using Multisig

After deploying a new implementation candidate for the upgrade, you need to propose the upgrade transaction using the multisig wallet.

## Prerequisites

- You need to have deployed the new implementation following the [Deploy New Implementation Guide](./3_b_1_deploy_new_impl.md).

- Upgrade `calldata` obtained from the deployment of the new implementation.

## Propose transaction

To propose the upgrade transaction you can follow the steps below:

1. Get the function signature

    If you are upgrading the AlignedLayerServiceManager or the RegistryCoordinator, you can get the function signature by running:

    ```bash
    cast sig "upgrade(address, address)"
    ```

   This will show the `upgrade` signature hash: `0x99a88ec4`.

    Else, if you are upgrading the BatcherPaymentService, you can get the function signature by running:

    ```bash
    cast sig "upgradeTo(address)"
    ```
   
    This will show the `upgradeTo` signature hash: `0x3659cfe6`.
    
2. Validate the calldata

    If you are upgrading the AlignedLayerServiceManager or the RegistryCoordinator, you can validate the calldata by running:

    ```
    cast calldata-decode "upgrade(address, address)" <calldata>
    ```

    This will show two addresses. 
    
    If you are upgrading the AlignmentLayerServiceManager, the first one is the `alignedLayerServiceManager` address, and the second one is the new implementation address of `alignedServiceManagerImplementation`.

    If you are upgrading the RegistryCoordinator, the first one is the `registryCoordinator` address, and the second one is the new implementation address of `registryCoordinatorImplementation`.

   > [!NOTE]
   > 
   > The first 10 characters must be the same the signature hash obtained in the previous step.
   >
   > Make sure the `alignedLayerServiceManager` address is the same as the one you deployed in the [Deploy Contracts Guide](./2_deploy_contracts.md).
   >
   > Make sure the `alignedServiceManagerImplementation` address is the same as the one you deployed in this guide.

    Else, if you are upgrading the BatcherPaymentService, you can validate the calldata by running:

    ```
    cast calldata-decode "upgradeTo(address)" <calldata>
    ```
   
    This will show the `batcherPaymentServiceImplementation` address.

    > [!NOTE]
    > 
    > The first 10 characters must be the same the signature hash obtained in the previous step.
    >
    > Make sure the `batcherPaymentServiceImplementation` address is the same as the one you deployed in this guide. 

3. 
