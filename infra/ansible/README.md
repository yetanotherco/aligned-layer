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

Now you have to set those variables.

Deploy the Batcher:

```shell
make ansible_batcher_deploy INVENTORY=</path/to/inventory> KEYSTORE=<path/to/keystore>
```

## Operator

To deploy the Operator you need to set some variables and the run the Operator playbook

Create the variables files:

```shell
make ansible_operator_create_env
```

This will create the following files in `infra/ansible/playbooks/ini`:

- `config-operator.ini`

Deploy the Operator:

```shell
make ansible_operator_deploy INVENTORY=</path/to/inventory> ECDSA_KEYSTORE=</path/to/ecdsa/keystore> BLS_KEYSTORE=</path/to/bls/keystore>
```
