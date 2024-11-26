#!/bin/bash

# cd to the directory of this script so that this can be run from anywhere
parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
# At this point we are in contracts/scripts
cd "$parent_path"

# At this point we are in contracts
cd ../

# Check if the number of arguments is correct
if [ "$#" -lt 1 ]; then
    echo "Usage: operator_remove_from_whitelist.sh <OPERATOR_ADDRESS>"
    echo "or"
    echo "Usage: operator_remove_from_whitelist.sh <OPERATOR_ADDRESS_1> <OPERATOR_ADDRESS_2> ... <OPERATOR_ADDRESS_N>"
    exit 1
fi

OPERATOR_ADDRESS=$1

# Read the registry coordinator address from the JSON file
REGISTRY_COORDINATOR=$(jq -r '.addresses.registryCoordinator' "$OUTPUT_PATH")

# Check if the registry coordinator address is empty
if [ -z "$REGISTRY_COORDINATOR" ]; then
    echo "Registry coordinator address is empty"
    exit 1
fi

# Check if the Ethereum RPC URL is empty
if [ -z "$RPC_URL" ]; then
    echo "Ethereum RPC URL is empty"
    exit 1
fi

# Check if the private key is empty
if [ -z "$PRIVATE_KEY" ]; then
    echo "Private key is empty"
    exit 1
fi

if [ "$#" -gt 1 ]; then
    OPERATORS=$(echo "$@" | sed 's/ /,/g') # separate the operators with a comma
    echo "Removing many operators from whitelist: $@"
    cast send \
        --rpc-url=$RPC_URL \
        --private-key=$PRIVATE_KEY \
        $REGISTRY_COORDINATOR 'remove_multiple(address[])' \
        "[$OPERATORS]"
else
    echo "Removing operator from whitelist: $OPERATOR_ADDRESS"
    cast send \
        --rpc-url=$RPC_URL \
        --private-key=$PRIVATE_KEY \
        $REGISTRY_COORDINATOR 'remove(address)' \
        $OPERATOR_ADDRESS
fi
