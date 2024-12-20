# Airdrop Claim Contracts

## Deployment

### Local

#### Requisites

- Foundry

#### Run

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
6. Approve the claimable contract to spend the token from the distributor
   ```
   make approve-claimable TOKEN=<token-proxy-address> AIRDROP=<claimable-proxy-address> PRIVATE_KEY=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
   ```

### Testnet (Sepolia) Deployment

#### Requisites

- Foundry
- Etherscan API key

#### Run

1. Create a file `script-config/config.sepolia.json` following the example in `script-config/config.sepolia.example.json`.
2. Deploy the token
   ```
   make deploy-token-testnet RPC_URL=<sepolia-rpc-url> PRIVATE_KEY=<sepolia-funded-account-private-key>
   ```
3. Write down the `token-proxy-address` that is printed in the console output. Do this in the `config.sepolia.json` file, under the `tokenProxy` key.
4. Deploy the claimable contract
   ```
   make deploy-claimable-testnet RPC_URL=<sepolia-rpc-url> PRIVATE_KEY=<sepolia-funded-account-private-key> ETHERSCAN_API_KEY=<etherscan-api-key>
   ```
5. Write down the `claimable-proxy-address` that is printed in the console output.
6. Approve the claimable contract to spend the token from the distributor
   ```
   make approve-claimable TOKEN=<token-proxy-address> AIRDROP=<claimable-proxy-address> PRIVATE_KEY=<sepolia-funded-distributor-private-key>
   ```

## Enabling Claimability

### Local

1. Deploy the claimable contract as explained above.
2. Set the correct merkle root
   ```
   make claimable-update-root RPC_URL=http://localhost:8545 PRIVATE_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a MERKLE_ROOT=<claims-merkle-root>
   ```
3. Set the correct claim time limit
   ```
   make claimable-update-timestamp RPC_URL=http://localhost:8545 PRIVATE_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a TIMESTAMP=2733427549
   ```
4. Approve the claimable contract to spend the token from the distributor
   ```
   make approve-claimable RPC_URL=http://localhost:8545 TOKEN=<token-proxy-address> AIRDROP=<claimable-proxy-address> PRIVATE_KEY=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
   ```
5. Unpause the claimable contract
   ```
   make unpause RPC_URL=http://localhost:8545 PRIVATE_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a AIRDROP=<claimable-proxy-address>
   ```
