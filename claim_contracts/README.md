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
3. Write down the token proxy address that is printed in the console output.
4. Deploy the claimable contract
   ```
   make deploy-claimable-local
   ```
5. Write down the claimable contract proxy address that is printed in the console output.
6. Approve the claimable contract to spend the token from the distributor
   ```
   make approve-claimable TOKEN=<token-proxy-address> AIRDROP=<claimable-proxy-address> PRIVATE_KEY=<distributor-private-key>
   ```
