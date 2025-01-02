# Aligned CLI

The Aligned CLI serves as a interface for users to interact with and retrieve information from Alinged Layer. This document serves as a reference for commands of the Aligned CLI.

## Installation:

1. Download and install Aligned from the Aligned github repo `https://github.com/yetanotherco/aligned_layer`:

```bash
curl -L https://raw.githubusercontent.com/yetanotherco/aligned_layer/main/batcher/aligned/install_aligned.sh | bash
```

2. A source command will be printed in your terminal after installation. Execute that command to update your shell environment.

3. Verify that the installation was successful:
```bash
aligned --help
```

## Help:

To see the available commands run:
```bash
aligned --help
```

To see the usage of a command run:
```bash
aligned [COMMAND] --help
```

## Commands:

- `submit [OPTIONS] --proving_system <Proving system> --proof <Proof file path>`: Submit a proof or (repetitions of a proof) to the Aligned Layer batcher.

    - Options:
        - `--batcher_url <Batcher connection address>`: Websocket URL for the Aligned Layer batcher.  
            - **default**: 
                - devnet: `ws://localhost:8080`
            - **possible values**: 
                - mainnet: `wss://mainnet.batcher.alignedlayer.com`
                - holesky: `wss://batcher.alignedlayer.com`
        - `--rpc_url <Ethereum RPC provider connection address>`: User's Ethereum RPC Node URL.
        - `--proving_system <Proving system>`: Proof system of the submitted proof.
            - **possible values**: `GnarkPlonkBls12_381`, `GnarkPlonkBn254`, `Groth16Bn254`, `SP1`, `Risc0`
        - `--proof <Proof file path>`: Path to the file the proof being submitted for verification is written.
        - `--public_input <Public input file name>`: Path to a file where the public inputs of the proof being submitted for verification is written.
        - `--vk <Verification key file name>`: Path to the file where the verification key of the proof being submitted for verification is written. Note, the following proof systems require a verification key to be supplied: `GnarkPlonkBls12_381`, `GnarkPlonkBn254`, `Groth16Bn254` 
        - `--vm_program <VM prgram code file name>`: Path to the file where the vm program code (ELF File) of the proof being submitted for verification is written. Note, the following proof systems require a vm program code (ELF File) to be supplied: `SP1`, `Risc0`. 
        - `--repetitions <Number of repetitions>`:
        Number of repetitions of this proof to be submitted.
            - **default**: 1
        - `--proof_generator_addr <Proof generator address>`: An optional parameter that can be used in some applications to avoid front-running.
            - **default**: `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`
        - `--aligned_verification_data_path <Aligned verification data directory Path>`: Path to the location the verification data of the proof will be written after it is submitted.
            - **default**: `./aligned_verification_data/`
        - `--keystore_path <Path to local keystore>`: Path to the local keystore of the user's wallet.
        - `--private_key <Private key>`: User's wallet private key.
        - `--max_fee <Max Fee (ether)>`: Amount of Ethereum in Ether to deposit into the `BatcherPaymentService.sol` contract to pay for proof submission.
        - `--nonce <Nonce>`: Nonce of the proof.
        - `--network <The working network's name>`: The Ethereum network of the Aligned Layer Network the user is submitting to.
            - **default**: `devnet`
            - **possible values**: `devnet`, `holesky`, `holesky-stage`, `mainnet`

- `verify-proof-onchain [OPTIONS] --aligned-verification-data <Aligned verification data>`:  Verify the proof was included in a verified batch on Ethereum.

    - Options:
        - `--aligned-verification-data <Aligned verification data>`: Path to verification data file generated from proof submission with the Aligned CLI.
        - `--rpc-url <Ethereum RPC provider address>`: User's Ethereum RPC node URL. 
        - `--network <The working network's name>`: The Ethereum network of the Aligned Layer network the user is querying their balance on.
            - **default**: `devnet`
            - **possible values**: `devnet`, `holesky`, `holesky-stage`, `mainnet`
- `get-vk-commitment [OPTIONS] --verification_key_file <Verification key file path> --proving_system <Proving system>`: Computes the verification data commitment specifically [provingSystemAuxDataCommitment](../2_architecture/components/3_service_manager_contract.md#verify-batch-inclusion) from the verification data file generated from submitting a proof to the batcher.

    - Options:
        - `--verification_key_file <Verification key file path>`: Path to verification data file generated from proof submission with the Aligned CLI.
        - `--proving_system <Proving system>`: Proof system of the supplied verification data file.
            - **possible values**: `GnarkPlonkBls12_381`, `GnarkPlonkBn254`, `Groth16Bn254`, `SP1`, `Risc0`
        - `--output <Output file>`: File path for output to be written to.

- `deposit-to-batcher [OPTIONS] --keystore_path <Path to local keystore> --amount <Amount to deposit>`: Deposits Ethereum into the Aligned Layer's `BatcherPaymentService.sol` contract to pay for users proof submission.

    - Options:
        - `--keystore_path <Path to local keystore>`: Path to the local keystore of the user's wallet. 
        - `--rpc-url <Ethereum RPC provider address>`: User's Ethereum RPC Node URL. 
        - `--network <The working network's name>`: The Ethereum network of the Aligned Layer Network the user is querying their balance on.
            - **default**: `devnet`
            - **possible values**: `devnet`, `holesky`, `holesky-stage`, `mainnet`
        - `--amount <Amount to deposit>`: Amount of Ethereum in Ether to deposit into the `BatcherPaymentService.sol` contract to pay for proof submission.

- `get-user-balance [OPTIONS] --user_addr <The user's Ethereum address>`: Retrieves the user's current balance in Aligned Layer's `BatcherPaymentService.sol` contract.

    - Options:
        - `--network <The working network's name>`: The Ethereum network of the Aligned Layer Network the user is querying their balance on.
            - **default**: `devnet`
            - **possible values**: `devnet`, `holesky`, `holesky-stage`, `mainnet`
        - `--rpc_url <Ethereum RPC provider address>`: User's Ethereum RPC Node URL.
            - **default**: `http://localhost:8545`
        - `--user_addr <The user's Ethereum address>`: User's Ethereum address on the provided network. 

- `get-user-nonce [OPTIONS] --user_addr <The user's Ethereum address>`:  Retrieves the user's current nonce from the batcher.

    - Options:
        - `--batcher_url <Batcher connection address>`: Websocket URL for the Aligned Layer batcher.  
            - **default**: 
                - devnet: `ws://localhost:8080`
            - **possible values**: 
                - mainnet: `wss://mainnet.batcher.alignedlayer.com`
                - holesky: `wss://batcher.alignedlayer.com`
        - `--user_addr <The user's Ethereum address>`: User's Ethereum address on Ethereum Mainnet.

## Example:

An example workflow for using the Aligned CLI (Runnable from the root directory of the Aligned Layer Repository):

1. Deposit funds to the batcher.
```bash
aligned deposit-to-batcher --network holesky --amount 0.5ether --keystore_path <KEYSTORE_PATH>
```

2. Check your balance in the aligned batcher.
```bash
aligned get-user-balance --user_addr <WALLET_ADDRESS> --network holesky --batcher_url wss://batcher.alignedlayer.com
```

3. Submit a Proof.
```bash
aligned submit  --proving_system Risc0 --proof ./scripts/test_files/risc_zero/fibonacci_proof_generator/risc_zero_fibonacci.proof --vm_program ./scripts/test_files/risc_zero/fibonacci_proof_generator/fibonacci_id.bin --public_input ./scripts/test_files/risc_zero/fibonacci_proof_generator/risc_zero_fibonacci.pub --repetitions <BURST_SIZE> --keystore_path <KEYSTORE_PATH> --batcher_url wss://batcher.alignedlayer.com --network holesky --max_fee 1300000000
```

4. Verify that your proof has been found on chain.
```bash
aligned verify-proof-onchain --aligned-verification-data ./aligned_verification_data/<VERIFICATION_DATA_FILE> --network holesky 
```

5. Check that the number of proofs you have submitted is incremented.
```bash
aligned get-user-nonce --user_addr <USER_ETH_ADDRESS> --batcher_url wss://holesky.batcher.alignedlayer.com
```
