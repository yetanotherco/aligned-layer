#!/bin/bash

if  [ -z "$1" ]; then
    echo "Usage: $0 <num> [<num> ...]"
    echo "or"
    echo "Usage: $0 all"
    exit 1
fi

if [[ "$1" == "all" ]]; then
    echo "Unpausing whole contract"
    cast send $ALIGNED_SERVICE_MANAGER \
        "unpause(uint256)()" 0 \
        --rpc-url $RPC_URL \
        --private-key $PRIVATE_KEY
    return
fi


result=0

for num in "$@"; do
    result=$((result | (1 << (num - 1))))
done

echo "New pause state: $result"

cast send $ALIGNED_SERVICE_MANAGER \
    "unpause(uint256)()" "$result" \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY

    
