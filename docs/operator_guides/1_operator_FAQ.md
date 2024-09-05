# Operator FAQ

The EigenLayer Operator FAQ is available [here](https://docs.eigenlayer.xyz/eigenlayer/operator-guides/operator-faq).

### How do I get whitelisted to run an operator?

To get whitelisted,
you need to fill out the form available [here](https://docs.google.com/forms/d/e/1FAIpQLSdH9sgfTz4v33lAvwj6BvYJGAeIshQia3FXz36PFfF-WQAWEQ/viewform)
and wait for the Aligned team to approve your request.

### What RPC should I use?

We suggest you use your own nodes for better performance and reliability. Note that the node must support HTTP and WebSockets.

We also recommend the following RPC providers:

- [Infura](https://infura.io/)
- [Blast](https://blastapi.io/)

The default configuration uses the public nodes RPC.

Since `v0.5.2`, there is a fallback mechanism to have two RPCs, so you need to add a second RPC for redundancy.

### How can I check if the version I'm using is the latest one?

You can check the latest version available on the [releases' page](https://github.com/yetanotherco/aligned_layer/releases).

### How do I check the operator version registered on the Aligned version tracking system?

You can get the registered version of your operator by running the following command:

```bash
curl https://holesky.tracker.alignedlayer.com/versions/<operator_address>
```

Replace `<operator_address>` with your operator address.

In future releases, the explorer will show the operator version and the latest available version.

### How do I check the operator version locally?

To see the operator version, run:

```bash
./operator/build/aligned-operator --version
```

You have to be in the root `aligned-layer` directory to run this command.

This will show the version of the operator you have built in your filesystem.

### How do I get the ID of the operator?

You can get the ID of the operator by running the following command:

```bash
cast call \
	--rpc-url https://ethereum-holesky-rpc.publicnode.com \
	0xD0A725d82649f9e4155D7A60B638Fe33b3F25e3b \
	"getOperatorId(address operator)(bytes32)" <operator_address>
```

> Note: You need to have installed [Foundry](https://book.getfoundry.sh/)

This will return the operator ID if the operator is registered, otherwise it will return `0x`.
