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
   > The first 10 characters must be the same the signature hash obtained in the previous step.
   > Make sure the `alignedLayerServiceManager` address is the same as the one you deployed in the [Deploy Contracts Guide](./2_deploy_contracts.md).
   > Make sure the `alignedServiceManagerImplementation` address is the same as the one you deployed in this guide.

    Else, if you are upgrading the BatcherPaymentService, you can validate the calldata by running:

    ```
    cast calldata-decode "upgradeTo(address)" <calldata>
    ```
   
    This will show the `batcherPaymentServiceImplementation` address.

    > [!NOTE]
    > The first 10 characters must be the same the signature hash obtained in the previous step.
    > Make sure the `batcherPaymentServiceImplementation` address is the same as the one you deployed in this guide. 

3. Verify the contract bytecode running the following command:

   ```
   TODO
   ```
   
4. Once the calldata and the proposed upgrade are validated, you can create the upgrade transaction on [Safe](https://app.safe.global/home)

5. Click on `New transaction` -> `Transaction Builder`
   
   ![New transaction](./images/3_b_2_multisig_1.png)

   ![Transaction Builder](./images/3_b_2_multisig_2.png)

6. Enable `Custom data`

7. Get the `ProxyAdmin` address, and paste it on `Enter Address or ENS Name`

   To get the `ProxyAdmin` address the following command will copy the address to the clipboard:

    ```bash
    # SEPOLIA
    jq -r ".addresses.alignedLayerProxyAdmin" contracts/script/output/sepolia/alignedlayer_deployment_output.json | pbcopy
    ```

    ```bash
   # HOLESKY
   jq -r ".addresses.alignedLayerProxyAdmin" contracts/script/output/holesky/alignedlayer_deployment_output.json | pbcopy
    ```
   
    ```bash
    # MAINNET
    jq -r ".addresses.alignedLayerProxyAdmin" contracts/script/output/mainnet/alignedlayer_deployment_output.json | pbcopy
    ```
   
>    [!NOTE]
>    Make sure to set the path to the correct deployment output file.

8. Once you paste the `ProxyAdmin` address, the ABI should be automatically filled.

    ![Enter Address or ENS Name](./images/3_b_2_multisig_3.png)

    ![Enter ABI](./images/3_b_2_multisig_4.png)

9. Set the `ETH Value` as 0

    ![ETH Value](./images/3_b_2_multisig_5.png)

10. Paste the calldata obtained from the deployment of the new implementation on the `Data` box and click on `+ Add new transaction`.

    ![Data](./images/3_b_2_multisig_6.png)

   You should see the new transaction to be executed on the right side.

11. Click on `Create batch` to create the transaction.

    ![Create batch](./images/3_b_2_multisig_7.png)

12. Review and confirm the transaction.
   
   To make sure everything is fine, simulate the transaction by clicking on `Simulate batch`

   Once the simulation is successful, click on `Send Batch` to send the transaction.

    ![Simulate batch](./images/3_b_2_multisig_8.png)

13. Confirm the transaction checking the function being called is correct and the contract address is the one you deployed.

   If everything is correct, click on `Sign` to send the transaction.

    ![Confirm transaction](./images/3_b_2_multisig_9.png)

14. Now in your transactions, you should be able to see the newly created transaction.

    ![New transaction](./images/3_b_2_multisig_10.png)

15. If the transaction is correctly created, you have to wait until the required Multisig member signs the transaction to send it.
