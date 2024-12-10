#!/bin/bash

# if [ "$#" -ne 7 ]; then
#   echo "Usage: $0 <safe> <owner1> <owner2> <owner3> <mint_amount> <rpc_url> <merkle_root>"
#   exit 1
# fi

jq --version >/dev/null 2>&1 && forge --version >/dev/null 2>&1
if [ $? != 0 ]; then
    echo "Error: Please make sure you have jq and forge installed and in your PATH"
    exit 2
fi

safe=${SAFE_ADDRESS:-$1}
owner1=${OWNER1_ADDRESS:-$2}
owner2=${OWNER2_ADDRESS:-$3}
owner3=${OWNER3_ADDRESS:-$4}
mint_amount=${MINT_AMOUNT:-$5}
rpc_url=${RPC_URL:-$6}
claim_time_limit=${CLAIM_TIME_LIMIT:-2733247661}
merkle_root=${MERKLE_ROOT:-$7}

cd script && forge script DeployL1Script $safe $owner1 $owner2 $owner3 $mint_amount $claim_time_limit $merkle_root --sig "run(address,address,address,address,uint256,uint256,bytes32)" --fork-url $rpc_url --broadcast
