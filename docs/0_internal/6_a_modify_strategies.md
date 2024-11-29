# Modify Strategy Weight Multipliers Contracts
This doc contains a guide on how to modify the strategy weight multipliers on Aligned.

To run the make targets specified in this guide, you must first have the relevant following env vars under `contracts/scripts/.env`:
```
RPC_URL=<rpc_url>
ALIGNED_SERVICE_MANAGER=<aligned_contract_address>
PRIVATE_KEY=<private_key>
OUTPUT_PATH=<deployment_output_path>
```

## 

To view some relevant information you can:

### Get all available strategies:

```
make strategies_get_addresses
```

### Get weight multiplier of a specific strategy:

```
make strategies_get_weight STRATEGY_INDEX=<strategy_index>
```

### Update the weight of any amount of stratefies

```
make strategies_update_weight STRATEGY_INDICES=[0,1,...,n] NEW_MULTIPLIERS=[0,1,...,n]
```

### Remove a strategy

```
make strategies_remove INDICES_TO_REMOVE=[0,1,...,n]
```

