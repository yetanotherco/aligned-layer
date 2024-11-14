# Guide to Deploy

## Batcher

> [!IMPORTANT]
> You need to have previously created an ECDSA keystore with at least 1ETH.
> You can create keystore following this [guide](#How-to-Create-Keystores)

To deploy the Batcher you need to set some variables and then run the Batcher playbook

Create the variables files:

```shell
make ansible_batcher_create_env
```

This will create the following files in `infra/ansible/playbooks/ini`

- `config-batcher.ini`
- `env-batcher.ini`
- `caddy-batcher.ini`

The `config-batcher.ini` contains the following variables:

| Variable                                  | Description                                                | Stage                                                                                                           | Testnet                                                                                                    | Mainnet                             |
|-------------------------------------------|------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|-------------------------------------|
| aligned_layer_deployment_config_file_path | JSON with Aligned contracts addresses                      | /home/app/repos/batcher/aligned_layer/contracts/script/output/holesky/alignedlayer_deployment_output.stage.json | /home/app/repos/operator/aligned_layer/contracts/script/output/holesky/alignedlayer_deployment_output.json | TBD                                 |
| eigen_layer_deployment_config_file_path   | JSON with EigenLayer contracts addresses                   | /home/app/repos/batcher/aligned_layer/contracts/script/output/holesky/eigenlayer_deployment_output.json         | /home/app/repos/operator/aligned_layer/contracts/script/output/holesky/eigenlayer_deployment_output.json   | TBD                                 |
| eth_rpc_url                               | HTTP RPC url                                               | <your_rpc_http_provider>                                                                                        | <your_rpc_http_provider>                                                                                   | <your_rpc_http_provider>            |
| eth_rpc_url_fallback                      | HTTP RPC fallback url. Must be different than eth_rpc_url  | https://ethereum-holesky-rpc.publicnode.com                                                                     | https://ethereum-holesky-rpc.publicnode.com                                                                | https://ethereum-rpc.publicnode.com |
| eth_ws_url                                | WS RPC url                                                 | <your_rpc_ws_provider>                                                                                          | <your_rpc_ws_provider>                                                                                     | <your_rpc_ws_provider>              |
| eth_ws_url_fallback                       | WS RPC fallback url. Must be different than eth_ws_rpc_url | wss://ethereum-holesky-rpc.publicnode.com                                                                       | wss://ethereum-holesky-rpc.publicnode.com                                                                  | wss://ethereum-rpc.publicnode.com   |
| ecdsa_private_key_store_password          | Password of the ECDSA keystore                             | <your_ecdsa_keystore_password>                                                                                  | <your_ecdsa_keystore_password>                                                                             | <your_ecdsa_keystore_password>      |
| telemetry_ip_port_address                 | IP:Port of the Telemetry service                           | <IP>:4001                                                                                                       | <IP>:4001                                                                                                  | TBD                                 |
| batcher_replacement_private_key           | This is the private key for the non-paying users           | -                                                                                                               | -                                                                                                          | -                                   |

The `env-batcher.ini` contains the following variables:

| Variable          | Description                                                | Stage                               | Testnet                               | Mainnet                       |
|-------------------|------------------------------------------------------------|-------------------------------------|---------------------------------------|-------------------------------|
| secret_access_key | Secret access key for user with access to the Storage (S3) | <your_secret_access_key>            | <your_secret_access_key>              | <your_secret_access_key>      |
| region            | Region of the Storage                                      | <us-east-1>                         | <us-east-1>                           | <us-east-1>                   |
| access_key_id     | Access key for the user with access to the Storage (S3)    | <your_access_key_id>                | <your_access_key_id>                  | <your_access_key_id>          |
| bucket_name       | Bucket name                                                | <stage.storage.example.com>         | <holesky.storage.example.com>         | <storage.example.com>         |
| download_endpoint | Public endpoint to download batcher                        | <https://stage.storage.example.com> | <https://holesky.storage.example.com> | <https://storage.example.com> |
| log_level         | Log level                                                  | info                                | info                                  | info                          |

The `caddy-batcher.ini` contains the following variables:

| Variable       | Description                          | Stage                       | Testnet                       | Mainnet               |
|----------------|--------------------------------------|-----------------------------|-------------------------------|-----------------------|
| batcher_domain | Domain of the Batcher to send proofs | <stage.batcher.example.com> | <holesky.batcher.example.com> | <batcher.example.com> |

> [!WARNING]
> You need to previously set the `batcher_domain` in your DNS provider to point to the Batcher IP.

Deploy the Batcher:

```shell
make ansible_batcher_deploy INVENTORY=</path/to/inventory> KEYSTORE=<path/to/keystore/in/your/local>
```

## Operator

> [!IMPORTANT]
> You need to have previously created an ECDSA keystore with at least 1ETH and a BLS keystore.
> You can create keystore following this [guide](#How-to-Create-Keystores)
> The ECDSA keystore for the Operator must be created with the Eigenlayer CLI

> [!CAUTION]
> To register the Operator in Aligned successfully, you need to have been whitelisted by the Aligned team previously.

To deploy the Operator you need to set some variables and then run the Operator playbook.

Create the variables files:

```shell
make ansible_operator_create_env
```

This will create the following files in `infra/ansible/playbooks/ini`:

- `config-operator.ini`
- `config-register-operator.ini`

The `config-register-operator.ini` contains the variables to register the Operator in EigenLayer:

| Variable                      | Description                                                                                                                                                                        | Stage                                                                                                                         | Testnet                                                                                                                       | Mainnet                                              |
|-------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| address                       | ECDSA address of the Operator                                                                                                                                                      | <your_ecdsa_operator_address>                                                                                                 | <your_ecdsa_operator_address>                                                                                                 | <your_ecdsa_operator_address>                        |
| metadata_url                  | Operator Metadata. You can create one following this [guide](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation#operator-configuration-and-registration) | <your_metadata_url>                                                                                                           | <your_metadata_url>                                                                                                           | <your_metadata_url>                                  |
| el_delegation_manager_address | Delegation Manager Address                                                                                                                                                         | `0xA44151489861Fe9e3055d95adC98FbD462B948e7`                                                                                  | `0xA44151489861Fe9e3055d95adC98FbD462B948e7`                                                                                  | TBD                                                  |
| eth_rpc_url                   | HTTP RPC url                                                                                                                                                                       | `https://ethereum-holesky-rpc.publicnode.com`                                                                                 | `https://ethereum-holesky-rpc.publicnode.com`                                                                                 | `https://ethereum-rpc.publicnode.com`                |
| private_key_store_path        | Path to the ECDSA keystore in the Operator host                                                                                                                                    | `/home/app/.keystores/operator.ecdsa`                                                                          | `/home/app/.keystores/operator.ecdsa`                                                                          | `/home/app/.keystores/operator.ecdsa` |
| private_key_store_password    | Password of the ECDSA keystore                                                                                                                                                     | <your_ecdsa_keystore_password>                                                                                                | <your_ecdsa_keystore_password>                                                                                                | <your_ecdsa_keystore_password>                       |
| chain_id                      | Chain ID                                                                                                                                                                           | 17000                                                                                                                         | 17000                                                                                                                         | 1                                                    |
| weth_address                  | Address of wETH token                                                                                                                                                              | [0x94373a4919B3240D86eA41593D5eBa789FEF3848](https://holesky.etherscan.io/address/0x94373a4919B3240D86eA41593D5eBa789FEF3848) | [0x94373a4919B3240D86eA41593D5eBa789FEF3848](https://holesky.etherscan.io/address/0x94373a4919B3240D86eA41593D5eBa789FEF3848) | TBD                                                  |
| weth_strategy_address         | Address of wETH token strategy                                                                                                                                                     | [0x80528D6e9A2BAbFc766965E0E26d5aB08D9CFaF9](https://holesky.eigenlayer.xyz/restake/WETH)                                     | [0x80528D6e9A2BAbFc766965E0E26d5aB08D9CFaF9](https://holesky.eigenlayer.xyz/restake/WETH)                                     | TBD                                                  |

The `config-operator.ini` contains the variables to run the Operator in Aligned:

| Variable                                  | Description                                                                                                                                                                        | Stage                                                                                                              | Testnet                                                                                                      | Mainnet                                        |
|-------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|------------------------------------------------|
| aligned_layer_deployment_config_file_path | JSON with Aligned contracts addresses                                                                                                                                              | `/home/app/repos/operator/aligned_layer/contracts/script/output/holesky/alignedlayer_deployment_output.stage.json` | `/home/app/repos/operator/aligned_layer/contracts/script/output/holesky/alignedlayer_deployment_output.json` | TBD                                            |
| eigen_layer_deployment_config_file_path   | JSON with EigenLayer contracts addresses                                                                                                                                           | `/home/app/repos/operator/aligned_layer/contracts/script/output/holesky/eigenlayer_deployment_output.json`         | `/home/app/repos/operator/aligned_layer/contracts/script/output/holesky/eigenlayer_deployment_output.json`   | TBD                                            |
| eth_rpc_url                               | HTTP RPC url                                                                                                                                                                       | <your_rpc_http_provider>                                                                                           | <your_rpc_http_provider>                                                                                     | <your_rpc_http_provider>                       |
| eth_rpc_url_fallback                      | HTTP RPC fallback url. Must be different than `eth_rpc_url`                                                                                                                        | `https://ethereum-holesky-rpc.publicnode.com`                                                                      | `https://ethereum-holesky-rpc.publicnode.com`                                                                | `https://ethereum-rpc.publicnode.com`          |
| eth_ws_url                                | WS RPC url                                                                                                                                                                         | <your_rpc_ws_provider>                                                                                             | your_rpc_ws_provider>                                                                                        | <your_rpc_ws_provider>                         |
| eth_ws_url_fallback                       | WS RPC fallback url. Must be different than `eth_ws_rpc_url`                                                                                                                       | `wss://ethereum-holesky-rpc.publicnode.com`                                                                        | `wss://ethereum-holesky-rpc.publicnode.com`                                                                  | `wss://ethereum-rpc.publicnode.com`            |
| ecdsa_private_key_store_path              | Path to the ECDSA keystore in the Operator host                                                                                                                                    | `/home/app/.keystores/operator.ecdsa`                                                                              | `/home/app/.keystores/operator.ecdsa`                                                                        | `/home/app/.keystores/operator.ecdsa`          |
| ecdsa_private_key_store_password          | Password of the ECDSA keystore                                                                                                                                                     | <your_ecdsa_keystore_password>                                                                                     | <your_ecdsa_keystore_password>                                                                               | <your_ecdsa_keystore_password>                 |
| bls_private_key_store_path                | Path to the BLS keystore in the Operator host                                                                                                                                      | `/home/app/.keystores/operator.bls`                                                                                | `/home/app/.keystores/operator.bls`                                                                          | `/home/app/.keystores/operator.bls`            |
| bls_private_key_store_password            | Password of the BLS keystore                                                                                                                                                       | <your_bls_keystore_password>                                                                                       | <your_bls_keystore_password>                                                                                 | <your_bls_keystore_password>                   |
| aggregator_rpc_server_ip_port_address     | Aggregator url                                                                                                                                                                     | `stage.aggregator.alignedlayer.com:8090`                                                                           | `aggregator.alignedlayer.com:8090`                                                                           | TBD                                            |
| operator_tracker_ip_port_address          | Telemetry service url                                                                                                                                                              | `https://stage.telemetry.alignedlayer.com`                                                                         | `https://holesky.telemetry.alignedlayer.com`                                                                 | TBD                                            |
| address                                   | ECDSA address of the Operator                                                                                                                                                      | <your_ecdsa_operator_address>                                                                                      | <your_ecdsa_operator_address>                                                                                | <your_ecdsa_operator_address>                  |
| metadata_url                              | Operator Metadata. You can create one following this [guide](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation#operator-configuration-and-registration) | <your_metadata_url>                                                                                                | <your_metadata_url>                                                                                          | <your_metadata_url>                            |
| enable_metrics                            | Expose or not prometheus metrics                                                                                                                                                   | `true`                                                                                                             | `true`                                                                                                       | `true`                                         |
| metrics_ip_port_address                   | Where to expose prometheus metrics if enabled                                                                                                                                      | `localhost:9092`                                                                                                   | `localhost:9092`                                                                                             | `localhost:9092`                               |
| last_processed_batch_filepath             | Where to store the last processed batch for system recovery                                                                                                                        | `/home/app/operator.last_processed_batch.json`                                                                     | `/home/app/operator.last_processed_batch.json`                                                               | `/home/app/operator.last_processed_batch.json` |

Deploy the Operator:

```shell
make ansible_operator_deploy INVENTORY=</path/to/inventory> ECDSA_KEYSTORE=</path/to/ecdsa/keystore> BLS_KEYSTORE=</path/to/bls/keystore>
```
## Aggregator

> [!IMPORTANT]
> You need to have previously created an ECDSA keystore with at least 1ETH and a BLS keystore.
> You can create keystore following this [guide](#How-to-Create-Keystores)

To deploy the Aggregator you need to set some variables and then run the Aggregator playbook

This will create the following file in `infra/ansible/playbooks/ini`:

- `config-aggregator.ini`

The `config-aggregator.ini` contains the variables to run the Operator in Aligned:

| Variable                                  | Description                                                                                                                                                                        | Stage                                                                                                                               | Testnet                                                                                                                       | Mainnet                                                       |
|-------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------|
| aligned_layer_deployment_config_file_path | JSON with Aligned contracts addresses                                                                                                                                              | `/home/app/repos/aggregator/aligned_layer/contracts/script/output/holesky/alignedlayer_deployment_output.stage.json` | `/home/app/repos/aggregator/aligned_layer/contracts/script/output/holesky/alignedlayer_deployment_output.json` | TBD                                                           |
| eigen_layer_deployment_config_file_path   | JSON with EigenLayer contracts addresses                                                                                                                                           | `/home/app/repos/aggregator/aligned_layer/contracts/script/output/holesky/eigenlayer_deployment_output.json`         | `/home/app/repos/aggregator/aligned_layer/contracts/script/output/holesky/eigenlayer_deployment_output.json`   | TBD                                                           |
| eth_rpc_url                               | HTTP RPC url                                                                                                                                                                       | <your_rpc_http_provider>                                                                                                            | <your_rpc_http_provider>                                                                                                      | <your_rpc_http_provider>                                      |
| eth_rpc_url_fallback                      | HTTP RPC fallback url. Must be different than `eth_rpc_url`                                                                                                                        | `https://ethereum-holesky-rpc.publicnode.com`                                                                                       | `https://ethereum-holesky-rpc.publicnode.com`                                                                                 | `https://ethereum-rpc.publicnode.com`                         |
| eth_ws_url                                | WS RPC url                                                                                                                                                                         | <your_rpc_ws_provider>                                                                                                              | your_rpc_ws_provider>                                                                                                         | <your_rpc_ws_provider>                                        |
| eth_ws_url_fallback                       | WS RPC fallback url. Must be different than `eth_ws_rpc_url`                                                                                                                       | `wss://ethereum-holesky-rpc.publicnode.com`                                                                                         | `wss://ethereum-holesky-rpc.publicnode.com`                                                                                   | `wss://ethereum-rpc.publicnode.com`                           |
| ecdsa_private_key_store_path              | Path to the ECDSA keystore in the Operator host                                                                                                                                    | `/home/app/.keystores/aggregator.ecdsa`                                                                              | `/home/app/.keystores/aggregator.ecdsa`                                                                        | `/home/app/.keystores/operator.ecdsa`          |
| ecdsa_private_key_store_password          | Password of the ECDSA keystore                                                                                                                                                     | <your_ecdsa_keystore_password>                                                                                                      | <your_ecdsa_keystore_password>                                                                                                | <your_ecdsa_keystore_password>                                |
| bls_private_key_store_path                | Path to the BLS keystore in the Operator host                                                                                                                                      | `/home/app/.keystores/aggregator.bls`                                                                                | `/home/app/.keystores/aggregator.bls`                                                                          | `/home/app/.keystores/operator.bls`            |
| bls_private_key_store_password            | Password of the BLS keystore                                                                                                                                                       | <your_bls_keystore_password>                                                                                                        | <your_bls_keystore_password>                                                                                                  | <your_bls_keystore_password>                                  |
| enable_metrics                            | Expose or not prometheus metrics                                                                                                                                                   | `true`                                                                                                                              | `true`                                                                                                                        | `true`                                                        |
| metrics_ip_port_address                   | Where to expose prometheus metrics if enabled                                                                                                                                      | `localhost:9092`                                                                                                                    | `localhost:9092`                                                                                                              | `localhost:9092`                                              |

Deploy the Aggregator:

```shell
make ansible_aggregator_deploy INVENTORY=</path/to/inventory> ECDSA_KEYSTORE=<path/to/keystore> BLS_KEYSTORE=</path/to/bls/keystore>
```
> [!Note]
> ECDSA_KEYSTORE and BLS_KEYSTORE are the paths of the keystores in your machine.

# How to Create Keystores

## Create ECDSA Keystore

Make sure you have installed:

- [Foundry](https://book.getfoundry.sh/getting-started/installation)

Now you can create the ECDSA keystore using the following command:

```shell
cast wallet new .
```

It will prompt you for a password and will save the keystore in your current directory.

If everything is okay, you will get the following output:

```
Created new encrypted keystore file: /your/current/path/f2e73ef1-d365-43b5-8818-07d6f7a254d4
Address: 0x...
```

Refer to this link for more details about keystore creation https://book.getfoundry.sh/reference/cast/cast-wallet-new

## Create ECDSA for Operator

Make sure you have installed:

- [Eigenlayer CLI](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation)

Now you can create the ECDSA keystore using the following command:

```shell
eigenlayer operator keys create --key-type ecdsa <keyname>
```
It will prompt for a password and will save the keystore in `$HOME/.eigenlayer/operator_keys/`.

If everything is okay, you will get the following output:

```
ECDSA Private Key (Hex):

////////////////////////////////////////////////////////////////////////////
//                                                                        //
//    ...                                                                 //
//                                                                        //
////////////////////////////////////////////////////////////////////////////

🔐 Please backup the above private key hex in a safe place 🔒
```

And then,

```

Key location: $HOME/.eigenlayer/operator_keys/<keyname>.ecdsa.key.json
Public Key hex: ...
Ethereum Address: 0x...

```

Refer to this link for more details about keystore creation https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation#create-keys

## Create BLS Keystore

Make sure you have installed:

- [Eigenlayer CLI](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation)

Now you can create the BLS keystore using the following command:

```shell
eigenlayer operator keys create --key-type bls <keyname>
```
It will prompt for a password and will save the keystore in `$HOME/.eigenlayer/operator_keys/`.

If everything is okay, you will get the following output:

```
BLS Private Key (Hex):

/////////////////////////////////////////////////////////////////////////////////////////
//                                                                                     //
//    ...                                                                              //
//                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////

🔐 Please backup the above private key hex in a safe place 🔒
```

And then,

```
Key location: $HOME/.eigenlayer/operator_keys/<keyname>.bls.key.json
Public Key: E([...,...])
```

Refer to this link for more details about keystore creation https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-installation#create-keys
