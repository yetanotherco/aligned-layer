#!/bin/bash

if [ -z "$OUTPUT_PATH" ]; then
    echo "OUTPUT_PATH env var is not set"
    exit 1
fi

if  [ -z "$RPC_URL" ]; then
    echo "RPC_URL env var is not set"
    exit 1
fi

if [ -z "$PRIVATE_KEY" ]; then
    echo "PRIVATE_KEY env var is not set"
    exit 1
fi

if [ -z "$INDICES_TO_REMOVE" ]; then
    echo "INDICES_TO_REMOVE env var is not set"
    exit 1
fi
if [[ ! "$INDICES_TO_REMOVE" =~ ^\[[0-9]+(,[0-9]+)*\]$ ]]; then
  echo "The INDICES_TO_REMOVE doesn't match the required format: [0,1,...,n]"
  exit 1
fi

STAKE_REGISTRY=$(jq -r '.addresses.stakeRegistry' "$OUTPUT_PATH")

## Using in this cast call:

    # /**
    #  * @notice This function is used for removing strategies and their associated weights from the
    #  * mapping strategyParams for a specific @param quorumNumber.
    #  * @dev higher indices should be *first* in the list of @param indicesToRemove, since otherwise
    #  * the removal of lower index entries will cause a shift in the indices of the other strategiesToRemove
    #  */
    # function removeStrategies(uint8 quorumNumber, uint256[] calldata indicesToRemove) external;

QUORUM_NUMBER=0 #Aligned has only 1 quorum for now

cast send $STAKE_REGISTRY "removeStrategies(uint8, uint256[])()" $QUORUM_NUMBER $INDICES_TO_REMOVE --private-key $PRIVATE_KEY --rpc-url $RPC_URL
