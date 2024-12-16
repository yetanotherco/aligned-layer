## AlignedToken

## Requirements

- Foundry

## Build

Install dependencies and compile the project running:

```bash
$ make build
```

## Local deploying

To deploy the contracts, set the following environment variables:

- `DEPLOYER_PRIVATE_KEY`: The private key of the account that's going to deploy the contracts.
- `SAFE_ADDRESS`: The address of the safe that's going to own the Proxy admin that in turn owns the token and airdrop contracts.
- `OWNER1_ADDRESS`, `OWNER2_ADDRESS`, and `OWNER3_ADDRESS`: The three owners of the token.
- `MINT_AMOUNT`: The amount to mint to each account (the contract actually supports minting different amounts of the token to each owner, but in the deploy script we simplified it).
- `RPC_URL`: The url of the network to deploy to.
- `CLAIM_TIME_LIMIT`: The claim time limit timestamp.
- `MERKLE_ROOT`: The merkle root of all valid token claims.

Example:
```
export DEPLOYER_PRIVATE_KEY=<deployer_private_key>
export SAFE_ADDRESS=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
export OWNER1_ADDRESS=0x70997970C51812dc3A010C7d01b50e0d17dc79C8
export OWNER2_ADDRESS=0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
export OWNER3_ADDRESS=0x90F79bf6EB2c4f870365E785982E1f101E93b906
export MINT_AMOUNT=100
export RPC_URL=http://localhost:8545
export CLAIM_TIME_LIMIT=2733247661
export MERKLE_ROOT=0x90076b5fb9a6c81d9fce83dfd51760987b8c49e7c861ea25b328e6e63d2cd3df
```

Then run the following script:

```
./deployClaim.sh
```
