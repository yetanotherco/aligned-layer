# Batcher

The Batcher receives proofs from different Users, bundles them in a batch of proofs, builds a Merkle Root from these, uploads the batch to a data service (like an S3 bucket), and submits this information to the [Aligned Service Manager](./3_service_manager_contract.md).

To ensure that the User is sure that their proof was included in a batch, the Batcher will send each User their Merkle Proof (or Merkle Path). With this, the User can rebuild the Merkle Root starting from their proof, thus verifying it was actually included in the batch.

Also, to avoid unnecessary proof submissions, the Batcher performs preliminary verifications of the submitted proofs in to minimize the submission of false proofs in a batch.

However, each proof has a cost of verification, so each batch must contain some sort of payment for it to be verified. To handle the payment for each batch, the Batcher submits the batch through its [Batcher Payment Service](./2_payment_service_contract.md).

To send the batch of proofs to the [Aligned Service Manager](./3_service_manager_contract.md), the Batcher stores the batch of proofs in an S3 for 1 week, and sends the link to the file to the [Aligned Service Manager](./3_service_manager_contract.md).

To view how to submit your own batch, without the use of this Batcher, you may follow [the following guide](../../3_guides/8_submitting_batch_without_batcher.md)


### Max fee priority queue

The batcher queue is now ordered by max fee signed by users in their proof messages - the ones willing to pay more will be prioritized in the batch.

Because of this, a user can't have a proof with higher nonce set with a higher fee included in the batch. For example, consider this situation in a batch. Let the two entries in the batch be from the same address:

	[(nonce: 1, max_fee: 5), (nonce: 2, max_fee: 10)]

This cannot happen because it will make the message with higher nonce be processed earlier than the one with a lower nonce, hence raising an invalid nonce error.

When a user submits a proof for the first time in the batch, its max fee is cached and set as the user min fee for the rest of the proofs to be included in the batch.
If a later message with a higher max fee is received, the message is rejected and not included in the queue, while if a message with a lower max fee is received,
the message is included in the queue and the user min fee is updated to that value. In summary, **no messages with a higher max fee than this user's min fee will be included**.

In case a message has a max fee that is too low - making it to be stuck in the batcher's *mempool*,  the message can be re-sent and be replaced with a higher fee.
To this end, a validation is done first. We check that when the max fee for the message with that nonce is updated, there is no message with a lower nonce and a lower max fee too, because this would lead to the problem
of messages with higher nonce processed earlier than messages with lower nonce, as discussed earlier.
As an example, consider all these messages in a batch from the same address:

	[(nonce: 1, max_fee: 10), (nonce: 2, max_fee: 5), (nonce: 3, max_fee: 3)]

If the user wants to send a replace message for the one with nonce 2, updating the max fee to something greater than 10 would not be valid
But it could update the max fee to, for example, a value of 8 and that would work.

## Batch finalization algorithm

There are some analogies in the processing of the batch with respect to how Ethereum handles transactions in the mempool and builds blocks.
We can consider the Batcher priority queue as a sort of *mempool* in the Ethereum case. Once certain conditions are met, the Batcher will start grabbing proofs and try to make a *finalized batch*, which in the analogy is like assembling a block in Ethereum.

When the conditions to build a batch are met, a batch finalization algorithm runs to create a batch of proofs from the priority queue.

This algorithm starts by calculating the **batch size** in bytes by adding the verification data bytes of each proof of the queue. This is needed in order to compute the **fee per proof** of the batch later. The next step is to build a new **resulting priority queue**, which will replace the current priority queue when this algorithm ends. On this new queue, all proofs which are not suited for the current batch will be stored.

In order for the batch to be considered valid, two conditions have to be met:
* The **batch size** in bytes must be less or equal to a certain established limit.
* All proofs found in the batch must have a **max fee** equal or higher to the calculated **fee per proof** of the batch.

The **fee per proof** is calculated with a formula, which depends on the **batch length**, calculated as the amount of proofs that the batch contains:

```
gas_per_proof = (constant_gas_cost + additional_submission_gas_cost_per_proof * batch_len) / batch_len
fee_per_proof = gas_per_proof * gas_price
```

Since the priority queue is sorted in ascending order by proof **max fee**, we can be certain that if the proof with the smallest **max fee** complies with the **fee per proof** rule, then all remaining proofs in the queue will do so

```
priority_queue = [(proof_a, 87), (proof_b, 90), (proof_c, 99)]
```

The algorithm will try to build a new batch by iterating on each proof, starting with the one with the smallest **max fee** in the queue. On each iteration, the **batch size** and **fee per proof** will be recalculated and both conditions reevaluated. When both conditions are met, all proofs contained in the queue will be used to build the new batch. The remaining proofs, stored in the **resulting priority queue**, will be candidates to the next batch finalization algorithm execution.

There is an edge case for this algorithm: If the fee per proof is too high even for the first entry, the algorithm will iterate over every entry until the **priority queue** empties. If this happens, the finalization of the batch is suspended and all the process will start again when a new block is received.

Let's see a very simple example:

The algorithm starts with the following elements:

```
priority_queue = [(E, 74), (D, 75), (C, 90), (B, 95), (A, 100)]
resulting_priority_queue = []
max_batch_size = 1000 # Defined constant
```

On the first iteration, the proof with the smallest **max fee** is taken and the **fee per proof** is calculated

```
current_proof = (E, 74)
fee_per_proof = calculate_fee_per_proof(priority_queue) # Result: 70
batch_size_bytes = calculate_batch_size(priority_queue) # Result: 1150
```

This batch can't be finalized, since it exceeds the maximum batch limit of 1000 bytes. This proof will be discarded for the current batch and stored in the **resulting priority queue**

```
priority_queue = [(D, 75), (C, 90), (B, 95), (A, 100)]
resulting_priority_queue = [(E, 74)]
current_proof = (D, 75)
fee_per_proof = calculate_fee_per_proof(priority_queue) # Result: 76
batch_size_bytes = calculate_batch_size(priority_queue) # Result: 990
```

This batch won't be finalized either, since the **fee per proof** of the batch is higher than the **max fee** of the current proof. This proof will be discarded for the current batch and stored in the **resulting priority queue**

```
priority_queue = [(C, 90), (B, 95), (A, 100)]
resulting_priority_queue = [(E, 74), (D, 75)]
current_proof = (C, 90)
fee_per_proof = calculate_fee_per_proof(priority_queue) # Result: 90
batch_size_bytes = calculate_batch_size(priority_queue) # Result: 850
```

All proofs in this batch comply with the **max fee** and **fee per proof** condition, and the batch size is lower than the established limit, so this batch will be finalized!

The execution ends with the following state for the batcher: a new batch is created , and the **priority queue** is replaced by the **resulting priority queue**

```
new_batch = [A, B, C] # Batch to send
priority_queue = [(E, 74), (D, 75)] # New priority queue
```
