#!/bin/bash

source .env

export ENVIRONMENT=$ENVIRONMENT
export RPC_URL=$RPC_URL
export PHX_HOST=$PHX_HOST
export DB_NAME=$DB_NAME
export DB_USER=$DB_USER
export DB_PASS=$DB_PASS
export DB_HOST=$DB_HOST
export ALIGNED_CONFIG_FILE=$ALIGNED_CONFIG_FILE
export DEBUG_ERRORS=$DEBUG_ERRORS
export TRACKER_API_URL=$TRACKER_API_URL

mix compile --force #force recompile to get the latest .env values

elixir --sname explorer -S mix phx.server
