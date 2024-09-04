# Explorer

{% embed url="https://explorer.alignedlayer.com" %}

The Explorer keeps track of [Aligned Service Manager](./3_service_manager_contract.md).

It has an internal state of previous batches, actively listens for new batches and their responses. The Explorer then displays this information for Users to visualize the submitted batches, their states and more useful information in real time.

In the landing page, we can see information such as how many [Operators](./4_operator.md) are currently registered and running, how many Batches and how many total Proofs have been verified.

![Figure 1: Explorer Landing Page](../../images/explorer-landing-page.png)

From here, we can search for a specific batch by its Merkle Root, we can directly jump to any one of the last 5 submitted batches, and we can easily go to the `Batches` page, where we can navigate through the various pages of batches of proofs submitted to aligned, ordered by latest submission, and easily check their on-chain status, timestamp, and block number.

![Figure 2: Explorer Batches Page](../../images/explorer-latest-batches.png)

From there, we can also click any individual batch hash to view its details.

From here we can visualize:

- the whole `Batch Hash` and copy it to our clipboard
- `Status`, either `Pending` or `Verified`
- `Batcher Sender Address` which is the address of the batcher that submitted the batch
- `Number of Proofs in this Batch` the number of proofs included in the batch
- `Fee per Proof`, fee paid per proof in the batch in ETH and USD
- `Proofs in the Batch`, which when pressed will show a list of all the proof hashes included in the batch
- Ethereum's `Submission Block Number`, linked to etherscan
- `Submission Transaction Hash`, linked to etherscan
- `Submission Timestamp` of the batch
- Ethereum's `Response Block Number`, linked to etherscan
- `Response Transaction Hash`, linked to etherscan
- `Response Timestamp` of the batch

![Figure 3: Explorer Batch Details Page](../../images/explorer-batch-details.png)

