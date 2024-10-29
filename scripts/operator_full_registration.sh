BASE_DIR="$(cd "$(dirname "$0")" && pwd)"

NUM_OPERATOR=0
LIMIT=$1
RESPOND_UNTIL=$2

if [[ -z $LIMIT ]]; then
    LIMIT=-1
fi

SHOULD_REGISTER=$3
if [[ -z $SHOULD_REGISTER || $SHOULD_REGISTER == true ]]; then
    # Remove prior configs as they will get regenerated
    rm -rf $BASE_DIR/config
    SHOULD_REGISTER=true
fi

function register_operator {
    NUM_OPERATOR=$1
    private_key=$2 
    stake=$3
    SHOULD_RESPOND=$4

    # Gen keys
    echo "Generating BLS keys"
    (echo "" | eigenlayer operator keys create  --insecure --key-type bls operator_$NUM_OPERATOR) </dev/null &>/dev/null &
    BLS_KEY_PATH=$HOME/.eigenlayer/operator_keys/operator_$NUM_OPERATOR.bls.key.json
    echo "Importing ECDSA keys"
    (echo "" | eigenlayer operator keys import --insecure --key-type ecdsa operator_$NUM_OPERATOR $private_key) </dev/null &>/dev/null &
    ECDSA_KEY_PATH=$HOME/.eigenlayer/operator_keys/operator_$NUM_OPERATOR.ecdsa.key.json
    sleep 1
    
    OPERATOR_ADDRESS=`jq -r '.address' $ECDSA_KEY_PATH`
    # Create config
    mkdir -p $BASE_DIR/config/$NUM_OPERATOR
    CONFIG_FILE=$BASE_DIR/config/$NUM_OPERATOR/config.yaml

    echo "environment: 'development'" >> $CONFIG_FILE
    echo "aligned_layer_deployment_config_file_path: './contracts/script/output/devnet/alignedlayer_deployment_output.json'" >> $CONFIG_FILE
    echo "eigen_layer_deployment_config_file_path: './contracts/script/output/devnet/eigenlayer_deployment_output.json'" >> $CONFIG_FILE
    echo "eth_rpc_url: 'http://localhost:8545'" >> $CONFIG_FILE
    echo "eth_rpc_url_fallback: 'http://localhost:8545'" >> $CONFIG_FILE
    echo "eth_ws_url: 'ws://localhost:8545'" >> $CONFIG_FILE
    echo "eth_ws_url_fallback: 'ws://localhost:8545'" >> $CONFIG_FILE
    echo "eigen_metrics_ip_port_address: 'localhost:9090'" >> $CONFIG_FILE
    echo "" >> $CONFIG_FILE
    echo "## ECDSA Configurations" >> $CONFIG_FILE
    echo "ecdsa:" >> $CONFIG_FILE
    echo "  private_key_store_path: '$ECDSA_KEY_PATH'" >> $CONFIG_FILE
    echo "  private_key_store_password: ''" >> $CONFIG_FILE
    echo "" >> $CONFIG_FILE
    echo "## BLS Configurations" >> $CONFIG_FILE
    echo "bls:" >> $CONFIG_FILE
    echo "  private_key_store_path: '$BLS_KEY_PATH'" >> $CONFIG_FILE
    echo "  private_key_store_password: ''" >> $CONFIG_FILE
    echo "" >> $CONFIG_FILE
    echo "## Operator Configurations" >> $CONFIG_FILE
    echo "operator:" >> $CONFIG_FILE
    echo "  aggregator_rpc_server_ip_port_address: 'localhost:8090'" >> $CONFIG_FILE
    echo "  operator_tracker_ip_port_address: 'http://localhost:4001'" >> $CONFIG_FILE
    echo "  address: 0x$OPERATOR_ADDRESS" >> $CONFIG_FILE
    echo "  earnings_receiver_address: 0x$OPERATOR_ADDRESS" >> $CONFIG_FILE
    echo "  delegation_approver_address: '0x0000000000000000000000000000000000000000'" >> $CONFIG_FILE
    echo "  staker_opt_out_window_blocks: 0" >> $CONFIG_FILE
    echo "  metadata_url: 'https://yetanotherco.github.io/operator_metadata/metadata.json'" >> $CONFIG_FILE
    echo "  enable_metrics: true" >> $CONFIG_FILE
    echo "  metrics_ip_port_address: 'localhost:9092'" >> $CONFIG_FILE
    echo "  max_batch_size: 268435456 # 256 MiB" >> $CONFIG_FILE
    echo "  last_processed_batch_filepath: 'config-files/operator-$NUM_OPERATOR.last_processed_batch.json'" >> $CONFIG_FILE
    echo "  should_respond: $SHOULD_RESPOND" >> $CONFIG_FILE
    echo "" >> $CONFIG_FILE
    echo "# Operators variables needed for register it in EigenLayer" >> $CONFIG_FILE
    echo "el_delegation_manager_address: '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9'" >> $CONFIG_FILE
    echo "private_key_store_path: '$ECDSA_KEY_PATH'" >> $CONFIG_FILE
    echo "bls_private_key_store_path: '$BLS_KEY_PATH'" >> $CONFIG_FILE
    echo "signer_type: local_keystore" >> $CONFIG_FILE
    echo "chain_id: 31337" >> $CONFIG_FILE

    echo "Config file generated."

    # Register operator
    echo "Registering operator with EigenLayer"
    echo "" | eigenlayer operator register $CONFIG_FILE

    # Mint mock tokens
    $BASE_DIR/mint_mock_token.sh $CONFIG_FILE $stake

    # Deposit into mock strategy
    echo "Depositing into mock strategy"
    STRATEGY_ADDRESS=$(jq -r '.addresses.strategies.MOCK' contracts/script/output/devnet/eigenlayer_deployment_output.json)
    go run operator/cmd/main.go deposit-into-strategy \
        --config $CONFIG_FILE \
        --strategy-address $STRATEGY_ADDRESS \
        --amount 100000000000000000

    # Whitelist operator
    echo "Whitelisting operator"
    RPC_URL="http://localhost:8545" PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" OUTPUT_PATH=./script/output/devnet/alignedlayer_deployment_output.json ./contracts/scripts/whitelist_operator.sh $OPERATOR_ADDRESS

    # Register with aligned layer
    echo "Registering operator with AlignedLayer"
    go run operator/cmd/main.go register \
        --config $CONFIG_FILE
}

while IFS=, read -r private_key stake should_respond; do
    # Ignore first line
    if [[ $private_key == "private_key" ]]; then
        continue
    fi

    NUM_OPERATOR=$((NUM_OPERATOR + 1))
    if [[ $LIMIT -lt $NUM_OPERATOR ]]; then
        break
    fi

    echo "NUM OPERATOR $NUM_OPERATOR"
    echo "SHOULD RESPOND $should_respond"
    echo "SHOULD REGISTER $SHOULD_REGISTER"

    if [[ $SHOULD_REGISTER == true ]]; then 
        register_operator $NUM_OPERATOR $private_key $stake $should_respond
    fi

    CONFIG_FILE=$BASE_DIR/config/$NUM_OPERATOR/config.yaml
    sed -i "s/should_respond: .*/should_respond: ${$should_respond}/g" "$CONFIG_FILE"

    # Start operator
    echo "Starting Operator..."
    (go run operator/cmd/main.go start --config $CONFIG_FILE \
        2>&1 | zap-pretty) &
done < $BASE_DIR/rich-wallets-anvil.csv

wait

