#!/bin/bash

source .env

export ENVIRONMENT=$ENVIRONMENT
export RPC_URL=$RPC_URL
ELIXIR_HOSTNAME=$(elixir -e 'IO.puts(:inet.gethostname() |> elem(1))')
export ELIXIR_HOSTNAME=$ELIXIR_HOSTNAME
export ALIGNED_CONFIG_FILE=$ALIGNED_CONFIG_FILE
export OPERATOR_FETCHER_WAIT_TIME_MS=$OPERATOR_FETCHER_WAIT_TIME_MS

if [ "$#" -eq 0 ]; then
    echo "Error, No arguments provided."
    echo "Try running the make target with FROM_BLOCK=<n>"
    exit 1
elif [ "$#" -eq 1 ]; then
    # argument provided, use it
    FROM=$1
else
    echo "Please provide 1 arguments."
    exit 1
fi

echo "Running fetch_operators_metadata.sh from block: $FROM"

mix compile --force #force recompile to get the latest .env values

echo "You will now connect to the Telemetry Node, make sure you run the following command inside it:"
echo "Scripts.FetchOperatorsMetadata.run($FROM)"

iex --sname fetch_operators_metadata --remsh telemetry@$ELIXIR_HOSTNAME
