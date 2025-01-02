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
