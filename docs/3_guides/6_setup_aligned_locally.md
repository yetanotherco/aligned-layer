# Aligned Infrastructure Deployment Guide

## Dependencies

Ensure you have the following installed:

- [Go](https://go.dev/doc/install)
- [Rust](https://www.rust-lang.org/tools/install)
- [Foundry](https://book.getfoundry.sh/getting-started/installation)
- [jq](https://jqlang.github.io/jq/)
- [yq](https://github.com/mikefarah/yq)

Then run:

```shell
make deps
```

This will:

- Initialize git submodules
- Install: `eigenlayer-cli`, `zap-pretty` and `abigen`
- Build ffis for your os.

## Contracts and eth node

To start anvil, a local Ethereum devnet with all necessary contracts already deployed and ready to be interacted with, run:

```shell
make anvil_start_with_block_time
```

<details>
<summary>More information on deploying the smart contracts on anvil:</summary>

### EigenLayer Contracts

If EigenLayer contracts change, the anvil state needs to be updated with:

```bash
make anvil_deploy_eigen_contracts
```

You will also need to redeploy the MockStrategy & MockERC20 contracts:

```bash
make anvil_deploy_mock_strategy
```

### Aligned Contracts

When changing Aligned contracts, the anvil state needs to be updated with:

```bash
make anvil_deploy_aligned_contracts
```

Note that when changing the contracts, you must also re-generate the Go smart contract bindings:

```bash
make bindings
```

</details>

## Aggregator

To start the [Aggregator](../2_architecture/components/5_aggregator.md):

```bash
make aggregator_start
```

or with a custom config:

```bash
make aggregator_start CONFIG_FILE=<path_to_config_file>
```

## Operator

To setup an [Operator](../2_architecture/components/4_operator.md) run:

```bash
make operator_register_and_start
```

or with a custom config:

```bash
make operator_register_and_start CONFIG_FILE=<path_to_config_file>
```

Different configs for operators can be found in `config-files/config-operator`.

<details>
<summary>More information about Operator registration:</summary>

If you wish to only register an operator you can run:

```bash
make operator_full_registration CONFIG_FILE<path_to_config_file>
```

and to start it once it has been registered:

```bash
make operator_start CONFIG_FILE=<path_to_config_file>
```

</details>

## Batcher

To start the [Batcher](../2_architecture/components/1_batcher.md) locally:

```bash
make batcher_start_local
```

This starts a [localstack](https://www.localstack.cloud/) to act as a replacement for S3.

If you want to use the batcher under a real `S3` connection you'll need to specify the environment variables under `batcher/aligned-batcher/.env` and then run:

```bash
make batcher_start
```

---

## Send test proofs

To send proofs quickly you can run any of the targets that have the prefix `batcher_send` for example:

Send a single plonk proof:

```shell
make batcher_send_plonk_bn254_task
```

Send a burst of `<N>` risc0 proofs:

```shell
make batcher_send_risc0_burst BURST_SIZE=<N>
```

Send an infinite stream of groth_16 proofs:

```shell
make batcher_send_burst_groth16 BURST_SIZE=2
```

Feel free to explore the rest of targets.
