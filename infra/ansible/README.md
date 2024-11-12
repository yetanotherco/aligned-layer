# Guide to Deploy

## Batcher

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
| ecdsa_private_key_store_path              | Path to the ECDSA keystore in the Operator host            | /home/app/.keystores/batcher.ecdsa                                                                              | /home/app/.keystores/batcher.ecdsa                                                                         | /home/app/.keystores/batcher.ecdsa  |
| ecdsa_private_key_store_password          | Password of the ECDSA keystore                             | <your_ecdsa_keystore_password>                                                                                  | <your_ecdsa_keystore_password>                                                                             | <your_ecdsa_keystore_password>      |
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
make ansible_batcher_deploy INVENTORY=</path/to/inventory> KEYSTORE=<path/to/keystore>
```