# Whitelist Operators
This doc contains a guide on how to Whitelist Operators in Aligned.

To run the make targets specified in this guide, you must first have the relevant following env vars under `contracts/scripts/.env`:
```
export RPC_URL=<rpc_url>
export OUTPUT_PATH=<aligned_deployment_output_file_path>
```

## Without Multisig

### Adding to Whitelist

You can whitelist a single Operator as following:
```
make operator_whitelist OPERATOR_ADDRESS=<opreator_address>
```

Or you can whitelist multiple Operators as following:
```
make operator_whitelist OPERATOR_ADDRESS=<opreator_address1,operator_address2,...,operator_addressN>
```

Note how there are no spaces between the commas that separate each Operator address.

For example:
```
make operator_whitelist OPERATOR_ADDRESS=0x70997970C51812dc3A010C7d01b50e0d17dc79C8,0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC,0x90F79bf6EB2c4f870365E785982E1f101E93b906
```

### Removing from Whitelist

You can remove from whitelist a single Operator as following:
```
make operator_remove_from_whitelist OPERATOR_ADDRESS=<opreator_address>
```

Or you can whitelist multiple Operators as following:
```
make operator_remove_from_whitelist OPERATOR_ADDRESS=<opreator_address1,operator_address2,...,operator_addressN>
```

Note how there are no spaces between the commas that separate each Operator address.

For example:
```
make operator_remove_from_whitelist OPERATOR_ADDRESS=0x70997970C51812dc3A010C7d01b50e0d17dc79C8,0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC,0x90F79bf6EB2c4f870365E785982E1f101E93b906
```

### Viewing Operator Whitelist status 

TODO

## With Multisig

To add or remove Operators from the Whitelist using a Multisig, you can follow the next guide
