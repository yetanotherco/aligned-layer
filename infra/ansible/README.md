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

Now you have to set those variables.

Deploy the Batcher:

```shell
make ansible_batcher_deploy INVENTORY=</path/to/inventory> KEYSTORE=<path/to/keystore>
```

## Operator

> [!IMPORTANT]
> You need to have previously created an ECDSA keystore with at least 1ETH and a BLS keystore.
> You can create keystore following this [guide](#How-to-Create-Keystores)

> [!CAUTION]
> To register the Operator in Aligned successfully, you need to have been whitelisted by the Aligned team previously.


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
//    19041818914970832833001653774136468328626805499863326892013784940157648962638    //
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
