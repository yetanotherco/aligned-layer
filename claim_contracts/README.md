#

## Local

### Requisites

- Foundry

### Run

1. Run anvil in one terminal:
   ```
   anvil
   ```
2. Deploy the token
   ```
   make deploy-token
   ```
3. Write down the token proxy address that is printed in the console output. Do this in the `config.example.json` file, under the `tokenProxy` key.
4. Deploy the claimable contract
   ```
   make deploy-claimable-local
   ```
5. Write down the claimable contract proxy address that is printed in the console output.

## Testnet (Sepolia)

### Requisites

- Foundry
- Etherscan API key

### Run

1. Create a file `script-config/config.sepolia.json` following the example in `script-config/config.sepolia.example.json`.
2. Deploy the token
   ```
   make deploy-token-testnet RPC_URL=<sepolia-rpc-url> PRIVATE_KEY=<sepolia-funded-account-private-key>
   ```
3. Write down the `token-proxy-address` that is printed in the console output. Do this in the `config.sepolia.json` file, under the `tokenProxy` key.
4. Deploy the claimable contract
   ```
   make deploy-claimable-testnet RPC_URL=<sepolia-rpc-url> DEPLOYER_PRIVATE_KEY=<sepolia-funded-account-private-key> ETHERSCAN_API_KEY=<etherscan-api-key>
   ```
5. Write down the `claimable-proxy-address` that is printed in the console output.

## Enabling Claimability

### By Calldata

> [!IMPORTANT]
>
> - This step-by-step **assumes** that the claimable proxy contract **is already deployed** and that **is already paused**. If it is not paused, the first transaction should be to pause it using this calldata `cast calldata "pause()"`.
> - This method **only** generates the necessary calldata to call the methods through transactions. It does **not** actually call the methods. This method is useful for copy-pasting the calldata into a multisig wallet.
> - Steps 1, 2, and 4 can be batched into a single transaction in a multisig wallet. This multisig must be the `ClaimableAirdrop` contract owner.
> - Step 3 must be done by the token distributor multisig as it is the one that has the tokens to be claimed.

> [!WARNING]
> - Double-check the data you passing into the commands, any mistake can lead to undesired behavior.

1. Update the merkle root
   ```
   // Example merkle_root = 0x97619aea42a289b94acc9fb98f5030576fa7449f1dd6701275815a6e99441927
   cast calldata "updateMerkleRoot(bytes32)" <merkle_root>
   ```
2. Update the claim time limit
   ```
   // Example timestamp = 2733427549
   cast calldata "extendClaimPeriod(uint256)" <timestamp>
   ```
3. Approve the claimable proxy contract to spend the token from the distributor (_2.6B, taking into account the 18 decimals_)
   ```
   // Example claim_proxy_address = 0x0234947ce63d1a5E731e5700b911FB32ec54C3c6
   cast calldata "approve(address,uint256)" <claim_proxy_address> 2600000000000000000000000000
   ```
4. Unpause the claimable contract (it is paused by default)
   ```
   cast calldata "unpause()"
   ```

### Local

1. Deploy the claimable contract as explained above.
2. Set the correct merkle root
   ```
   make claimable-update-root MERKLE_ROOT=<claims-merkle-root>
   ```
3. Set the correct claim time limit
   ```
   make claimable-update-timestamp TIMESTAMP=2733427549
   ```
4. Approve the claimable contract to spend the token from the distributor
   ```
   make approve-claimable
   ```
5. Unpause the claimable contract
   ```
   make claimable-unpause
   ```

or

```
make deploy-example MERKLE_ROOT=<claims-merkle-root> TIMESTAMP=2733427549
```
