## AlignedToken

## Requirements

- Foundry

## Local deploying

Install the dependencies:

```
make deps
```

### Token

To deploy the token, modify the file `script-config/config.example.json` and set the following values:

- `safe`: The address of the safe that's going to own the Proxy admin that in turn owns the token and airdrop contracts.
- `salt`: The salt used to generate CREATE2 addresses.
- `deployer`: The address of the account that's going to deploy the contracts.
- `foundation`: The address of the Aligned Foundation account.
- `claimSupplier`: The address of the account that's going to supply the tokens.
- `limitTimestampToClaim`: The claim time limit timestamp.
- `claimMerkleRoot`: The merkle root of all valid token claims.

Also, set the following environment variables:

- `PRIVATE_KEY`: The private key of the account that's going to deploy the contracts. This MUST be the same account as the one in `deployer`.
- `RPC_URL`: The url of the network to deploy to.

Then run the script:

```
make deploy-token
```

> [!TIP]
> You can create another config file on `script-config/config.custom_name.json` and run the script with
> ```
> make deploy-token CONFIG=custom_name
> ```

### Airdrop

TBD
