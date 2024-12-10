#!/bin/bash

source .env

export ENVIRONMENT=$ENVIRONMENT
export RPC_URL=$RPC_URL
ELIXIR_HOSTNAME=$(elixir -e 'IO.puts(:inet.gethostname() |> elem(1))')
export ELIXIR_HOSTNAME=$ELIXIR_HOSTNAME
export ALIGNED_CONFIG_FILE=$ALIGNED_CONFIG_FILE
export OPERATOR_FETCHER_WAIT_TIME_MS=$OPERATOR_FETCHER_WAIT_TIME_MS

mix compile --force #force recompile to get the latest .env values

echo "You will now connect to the Telemetry Node, make sure you run the following command inside it:"
echo "Scripts.FetchOperatorsMetadata.run()"

iex --sname fetch_operators_metadata --remsh telemetry@$ELIXIR_HOSTNAME
