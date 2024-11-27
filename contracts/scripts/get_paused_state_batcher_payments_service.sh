#!/bin/bash

if  [ -z "$BATCHER_PAYMENT_SERVICE" ]; then
    echo "BATCHER_PAYMENT_SERVICE env var is not set"
    exit 1
fi

if  [ -z "$RPC_URL" ]; then
    echo "RPC_URL env var is not set"
    exit 1
fi

number=$(cast call $BATCHER_PAYMENT_SERVICE "paused()(bool)" --rpc-url $RPC_URL)
echo Batcher Payments Paused state: $number
