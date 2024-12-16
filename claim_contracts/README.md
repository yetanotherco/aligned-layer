## AlignedToken

## Requirements

- Foundry

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

## Production

This is a series of steps to deploy the token to production and upgrade it if necessary.

### Safe wallet creation

First we create a wallet in [Safe](https://app.safe.global/) to represent the foundation. We assume this is done by the user. This safe will:
- Receive part of the deployed tokens.
- Own the proxy admin contract. This is the contract that can upgrade the token contract.
- Own the proxy contract. This will point to the current implementation of the contract. All users and the safe itself will interact with this contract to mint, transfer, burn, and other erc20 operations.

### Configuration

There is an example configuration file in `script-config/config.example.json`. Before deploying to production, we need to create a `config.mainnet.json` file in the same folder with the same contents as the example, but we need to change a couple of fields:

- `foundation`: this is the address of the safe that was created in the previous step.
- `claimSupplier`: this is the address of a different safe that will provide the funds for the claim contract when it is deployed.
- `deployer`: The address of the deterministic create2 deployer as specified in the [official repo](https://github.com/Arachnid/deterministic-deployment-proxy). The address should be `0x4e59b44847b379578588920ca78fbf26c0b4956c`.
- `salt`: An arbitrary value provided by the sender. This is a 32-bytes hex string. We default to 0.

### Deployment of the Token (proxy and implementation)

```bash
make deploy-token PRIVATE_KEY=<private_key> RPC_URL=<rpc_url> CONFIG=<config>
```

This make target internally executes a forge script that:

1. Deploys the token implementation that will be used only for its logic.
2. Deploys the Transparent Proxy, which also deploys the proxy admin. In this step the safe is set as the owner of both the owner of the proxy admin and the proxy (token).

The private key does NOT correspond to the safe, it needs to represent an account with sufficient funds to deploy the token.

Arguments (env variables):
- `PRIVATE_KEY`: the private key of the deployer account. This is NOT the foundation safe, just any account with enough eth for the deployment. This operation consumes approximately `3935470` gas units. As of Dec 16 2024, the gas price for a high priority is 16 gwei, which means around `0.063` eth.
- `RPC_URL`: a gateway or node that allows for rpc calls.
- `CONFIG`: the name of the configuration file. For `config.example.json` the name would be `example`. For `config.mainnet.json`, this would be `mainnet`.

The output of the deployment will look something like this:

```
== Logs ==
Aligned Token Proxy deployed at address: 0x9eDC342ADc2B73B2E36d0e77475bCF2103F09a22 with proxy admin: 0x51D94AdA2FFBFED637e6446CC991D8C65B93e167 and owner: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC

...

##### sepolia
✅  [Success]Hash: 0x6b58e061e0c37209bf7e33a4a4705706b5accd13b8964ecf2ae919fe01f41da1
Contract Address: 0xDC1dc4e84b0FB522DBa3a14C909b90c96496830C
Block: 7293558
Paid: 0.008478255653505184 ETH (1850587 gas * 4.581387232 gwei)
```

As a sanity check, we can call the following make target with the address of the deployed proxy contract.

```bash
make test-token ADDRESS=0x9eDC342ADc2B73B2E36d0e77475bCF2103F09a22
```

This calls functions to get the name, symbol and total supply in the deploy contract and displays them on the screen.

### Contract Verification in explorer

We can also verify the implementation contract in Etherscan, which is useful for other applications like wallets to recognize the ABI.

```
forge verify-contract 0xDC1dc4e84b0FB522DBa3a14C909b90c96496830C --rpc-url $RPC_URL --etherscan-api-key $ETHERSCAN_API_KEY
```

In this case we pass the implementation contract address.
