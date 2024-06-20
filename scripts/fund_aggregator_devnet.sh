#!/bin/bash

# Check that AGGREGATOR_ADDRESS is not empty
if [[ "$AGGREGATOR_ADDRESS" -eq "" ]]; then
  echo "AGGREGATOR_ADDRESS is empty, using default value 0x4bbddf258606a95270e12ecd97e09fec5532378c"
  AGGREGATOR_ADDRESS=0x4bbddf258606a95270e12ecd97e09fec5532378c
fi;

# Ethereum sender address - anvil address 2
sender_address="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"

# Ethereum sender private key - anvil private key 2
sender_private_key="0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

# Recipient address
recipient_address=$AGGREGATOR_ADDRESS

# Amount of ETH to send (in Ether)
amount_in_eth="1ether"

# Send Ether transaction
cast send --from $sender_address \
    --value $amount_in_eth \
    --private-key $sender_private_key \
    --rpc-url "http://localhost:8545" \
    "$recipient_address"
