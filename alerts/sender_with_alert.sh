#!/bin/bash

# ENV VARIABLES:
# - REPETITIONS
# - SLEEP
# - RPC_URL
# - SENDER_ADDRESS
# - BATCHER_URL
# - NETWORK
# - PRIVATE_KEY
# - PAGER_DUTY_KEY
# - PAGER_DUTY_EMAIL
# - PAGER_DUTY_SERVICE_ID

# Load env file from $1 path
source "$1"

# Just for debugging
#set -ex

### FUNCTIONS ###

# Function to send PagerDuty alert
# @param message
function send_pagerduty_alert() {
  . alerts/pagerduty.sh "$1"
}

# Function so send Slack message
# @param message
function send_slack_message() {
  . alerts/slack.sh "$1"
}

#################
while true
do

  ## Remove Proof Data
  rm -rf ./scripts/test_files/gnark_groth16_bn254_infinite_script/infinite_proofs/*
  rm -rf ./aligned_verification_data/*

  mkdir -p ./scripts/test_files/gnark_groth16_bn254_infinite_script/infinite_proofs

  ## Generate Proof
  x=$(aligned get-user-nonce --batcher_url $BATCHER_URL --user_addr $SENDER_ADDRESS 2>&1 | awk '{print $9}')
  echo "Generating proof $x != 0"
  go run ./scripts/test_files/gnark_groth16_bn254_infinite_script/cmd/main.go $x

  ## Get initial balance
  initial_balance=$(aligned get-user-balance --rpc_url $RPC_URL --user_addr  $SENDER_ADDRESS --network $NETWORK 2>&1 | awk '{print $7}')
  echo "Initial balance: $initial_balance"

  ## Send Proof
  echo "Submitting $REPETITIONS proofs $x != 0"
  submit=$(aligned submit \
    --proving_system Groth16Bn254 \
    --repetitions $REPETITIONS \
    --proof "./scripts/test_files/gnark_groth16_bn254_infinite_script/infinite_proofs/ineq_${x}_groth16.proof" \
    --public_input "./scripts/test_files/gnark_groth16_bn254_infinite_script/infinite_proofs/ineq_${x}_groth16.pub" \
    --vk "./scripts/test_files/gnark_groth16_bn254_infinite_script/infinite_proofs/ineq_${x}_groth16.vk" \
    --proof_generator_addr $SENDER_ADDRESS \
    --private_key $PRIVATE_KEY \
    --rpc_url $RPC_URL \
    --batcher_url $BATCHER_URL \
    --network $NETWORK \
    --max_fee 4000000000000000 \
    2>&1)

  echo "$submit"

  explorer_link=$(echo "$submit" | grep alignedlayer.com | awk '{print $4}')
  sleep 600

  ## Get final balance
  final_balance=$(aligned get-user-balance --rpc_url $RPC_URL --user_addr  $SENDER_ADDRESS --network $NETWORK 2>&1 | awk '{print $7}')
  echo "Final balance: $final_balance"

  ## Calculate the spent amount
  spent_amount=$(echo "$initial_balance - $final_balance" | bc | awk '{printf "%.9f", $1}')

  ## Verify Proofs
  echo "Verifying $REPETITIONS proofs $x != 0"
  for proof in ./aligned_verification_data/*.cbor; do
    ## Validate Proof on Chain
    verification=$(aligned verify-proof-onchain \
      --aligned-verification-data $proof \
      --rpc_url $RPC_URL \
      --network $NETWORK \
      2>&1)

    ## Send Alert is Verification Fails
    if echo "$verification" | grep -q not; then
      message="Proof verification failed for $proof [ $explorer_link ]"
      echo "$message"
      send_pagerduty_alert "$message"
      break
    elif echo "$verification" | grep -q verified; then
      echo "Proof verification succeeded for $proof"
    fi
  done

  ## Send Update to Slack
  eth_usd=$(curl https://cryptoprices.cc/ETH/)
  spent_ammount_usd=$(echo "$spent_amount * $eth_usd" | bc | awk '{printf "%.2f", $1}')
  slack_meesage="$REPETITIONS Proofs submitted and verified. Spent amount: $spent_amount ETH ($ $spent_ammount_usd) [ $explorer_link ]"
  send_slack_message "$slack_meesage"

  ## Remove Proof Data
  rm -rf ./scripts/test_files/gnark_groth16_bn254_infinite_script/infinite_proofs/*
  rm -rf ./aligned_verification_data/*

  echo "Sleeping $SLEEP seconds"
  sleep $SLEEP
done
