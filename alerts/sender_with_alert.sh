#!/bin/bash

# ENV VARIABLES:
# - REPETITIONS
# - EXPLORER_URL
# - SENDER_ADDRESS
# - BATCHER_URL
# - RPC_URL
# - EXPLORER_URL
# - NETWORK
# - PRIVATE_KEY
# - VERIFICATION_WAIT_TIME
# - SLEEP_TIME
# - PAGER_DUTY_KEY
# - PAGER_DUTY_EMAIL
# - PAGER_DUTY_SERVICE_ID

# TODO (Improvements):
# 1. This script does not account for proofs being included in different batches.
#    You can test that behavior by modifying the batcher's batch limit and sending many repetitions (REPETITIONS > BATCH_LIMIT) will throw an error on batcher
# 2. This script parses the submission and response tx hashes from the explorer html. This may easily break if the explorer is modified.
#    We should be able to parse this information from ethereum logs instead
# 3. This script waits VERIFICATION_WAIT_TIME seconds before fetching the explorer for the response tx hash.
#    We should instead poll in a loop until the batch is marked as verified

# ACKNOWLEDGMENTS
#
# Special thanks to AniV for their contribution on StackExchange regarding AWK formatting with floating point operations: 
# https://unix.stackexchange.com/questions/292087/how-do-i-get-bc-to-start-decimal-fractions-with-a-leading-zero


# Load env file from $1 path
source "$1"

# Just for debugging
#set -ex

### FUNCTIONS ###

# Function to get the tx cost from the explorer:
# @param batch_explorer_url
# @param type_of_hash
#   - Submission Transaction Hash
#   - Response Transaction Hash
function fetch_tx_cost() {
  tx_hash=$(curl -s $1 | grep -C 5 "$2" | grep -oE '0x[[:alnum:]]{64}' | head -n 1)
  # Get the tx receipt from the blockchain 
  receipt=$(cast receipt --rpc-url $RPC_URL $tx_hash)
  # Parse the gas used and gas price
  gas_price=$(echo "$receipt" | grep "effectiveGasPrice" | awk '{ print $2 }')
  gas_used=$(echo "$receipt" | grep "gasUsed" | awk '{ print $2 }')
  # Calculate fee in wei
  fee_in_wei=$(($gas_price * $gas_used))

  echo $fee_in_wei
}

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
  nonce=$(aligned get-user-nonce --batcher_url $BATCHER_URL --user_addr $SENDER_ADDRESS 2>&1 | awk '{print $9}')
  x=$((nonce + 1)) # So we don't have any issues with nonce = 0
  echo "Generating proof $x != 0"
  go run ./scripts/test_files/gnark_groth16_bn254_infinite_script/cmd/main.go $x

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
  
  echo "Waiting $VERIFICATION_WAIT_TIME seconds for verification"
  sleep $VERIFICATION_WAIT_TIME

  # Get the batch merkle root
  batch_merkle_root=$(echo "$submit" | grep "Batch merkle root: " | grep -oE "0x[[:alnum:]]{64}" | uniq | head -n 1) # TODO: Here we are only getting the first merkle root

  # Calculate the fee
  batch_explorer_url="$EXPLORER_URL/batches/$batch_merkle_root"
  submission_fee_in_wei=$(fetch_tx_cost $batch_explorer_url "Submission Transaction Hash")
  response_fee_in_wei=$(fetch_tx_cost $batch_explorer_url "Response Transaction Hash")
  total_fee_in_wei=$((submission_fee_in_wei + response_fee_in_wei))

  # Calculate the spent amount by converting the fee to ETH
  wei_to_eth_division_factor=$((10**18))
  spent_amount=$(echo "scale=30; $total_fee_in_wei / 1e18" | bc -l | awk '{printf "%.15f", $0}')

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
  eth_usd=$(curl -s https://cryptoprices.cc/ETH/)
  spent_ammount_usd=$(echo "$spent_amount * $eth_usd" | bc | awk '{printf "%.2f", $1}')
  slack_meesage="$REPETITIONS Proofs submitted and verified. Spent amount: $spent_amount ETH ($ $spent_ammount_usd) [ $batch_explorer_url ]"

  send_slack_message "$slack_meesage"

  ## Remove Proof Data
  rm -rf ./scripts/test_files/gnark_groth16_bn254_infinite_script/infinite_proofs/*
  rm -rf ./aligned_verification_data/*

  echo "Sleeping $SLEEP_TIME seconds"
  sleep $SLEEP_TIME
done
