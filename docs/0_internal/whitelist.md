# Whitelist

Whitelisting is a functionality added by the Aligned dev team to the `eigenlayer-middleware` contracts.

This functionality is added in `contracts/lib/eigenlayer-middleware/src/Whitelist.sol`, and its functionality ends up present in the `RegistryCoordinator` contract.

The reason behind this functionality is to control which Operators can operate in Aligned, making necessary for the team to manually whitelist an Operator before it can join the network.

## Interacting with whitelisting

### Adding an Operator

There are 2 ways of adding Operators to the whitelist.

To add only 1 Operator:
```
make operator_whitelist OPERATOR_ADDRESS=<operator_address>
```

To add a list of Operators:
```
export contracts/scripts/.env
contracts/scripts/operator_whitelist.sh <operator_address_1> <operator_address_2> ... <operator_address_n>
```

### Removing an Operator

There are 2 ways of removing Operators from the whitelist.

To remove only 1 Operator:
```
make operator_remove_from_whitelist OPERATOR_ADDRESS=<operator_address>
```

To remove a list of Operators:
```
export contracts/scripts/.env
contracts/scripts/operator_remove_from_whitelist.sh <operator_address_1> <operator_address_2> ... <operator_address_n>
```

### Querying the state of an Operator

To view the whitelist state of an Operator you can:

```
cast call <aligned_registry_coordinator_address> "isWhitelisted(address)(bool)" <operator_address>
```

or in Etherscan:

1. Locate the Aligned `registryCoordinator` contract address, in `aligned_deployment_output.json`
2. `Read as Proxy` in Etherscan
3. Find the `isWhitelisted` function, and put the Operator's address as the parameter
