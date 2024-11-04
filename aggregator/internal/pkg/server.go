package pkg

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/rpc"
	"time"

	"github.com/yetanotherco/aligned_layer/core/types"
)

const waitForEventRetries = 50
const waitForEventSleepSeconds = 4 * time.Second

func (agg *Aggregator) ServeOperators() error {
	// Registers a new RPC server
	err := rpc.Register(agg)
	if err != nil {
		return err
	}

	// Registers an HTTP handler for RPC messages
	rpc.HandleHTTP()

	// Start listening for requests on aggregator address
	// ServeOperators accepts incoming HTTP connections on the listener, creating
	// a new service goroutine for each. The service goroutines read requests
	// and then call handler to reply to them
	agg.logger.Info("Starting RPC server on address", "address",
		agg.AggregatorConfig.Aggregator.ServerIpPortAddress)

	err = http.ListenAndServe(agg.AggregatorConfig.Aggregator.ServerIpPortAddress, nil)

	return err
}

// Waits for the arrival of task associated with signedTaskResponse and returns true on success or false on failure
// If the task is not present in the internal map, it will try to fetch it from logs and retry.
// The number of retries is specified by `waitForEventRetries`, and the waiting time between each by `waitForEventSleepSeconds`
func (agg *Aggregator) waitForTask(signedTaskResponse *types.SignedTaskResponse) bool {
	for i := 0; i < waitForEventRetries; i++ {
		// Lock
		agg.taskMutex.Lock()
		agg.AggregatorConfig.BaseConfig.Logger.Info("- Locked Resources: Starting processing of Response")
		_, ok := agg.batchesIdxByIdentifierHash[signedTaskResponse.BatchIdentifierHash]
		// Unlock
		agg.logger.Info("- Unlocked Resources: Task not found in the internal map")
		agg.taskMutex.Unlock()
		if ok {
			return true
		}

		// Task was not found in internal map, let's try to fetch it from logs
		agg.logger.Info("Trying to fetch missed task from logs...")
		batch, err := agg.avsReader.GetPendingBatchFromMerkleRoot(signedTaskResponse.BatchMerkleRoot)

		if err == nil && batch != nil {
			agg.logger.Info("Found missed task in logs with merkle root 0x%e", batch.BatchMerkleRoot)
			// Adding new task will fail only if it already exists
			agg.AddNewTask(batch.BatchMerkleRoot, batch.SenderAddress, batch.TaskCreatedBlock)
			return true
		}

		if err != nil {
			agg.logger.Warn("Error fetching task from logs: %v", err)
		}

		if batch == nil {
			agg.logger.Info("Task not found in logs")
		}

		// Task was not found, wait and retry
		time.Sleep(waitForEventSleepSeconds)
	}

	return false
}

// Aggregator Methods
// This is the list of methods that the Aggregator exposes to the Operator
// The Operator can call these methods to interact with the Aggregator
// This methods are automatically registered by the RPC server
// This takes a response an adds it to the internal. If reaching the quorum, it sends the aggregated signatures to ethereum
// Returns:
//   - 0: Success
//   - 1: Error
func (agg *Aggregator) ProcessOperatorSignedTaskResponseV2(signedTaskResponse *types.SignedTaskResponse, reply *uint8) error {
	agg.AggregatorConfig.BaseConfig.Logger.Info("New task response",
		"BatchMerkleRoot", "0x"+hex.EncodeToString(signedTaskResponse.BatchMerkleRoot[:]),
		"SenderAddress", "0x"+hex.EncodeToString(signedTaskResponse.SenderAddress[:]),
		"BatchIdentifierHash", "0x"+hex.EncodeToString(signedTaskResponse.BatchIdentifierHash[:]),
		"operatorId", hex.EncodeToString(signedTaskResponse.OperatorId[:]))

	if !agg.waitForTask(signedTaskResponse) {
		agg.logger.Warn("Task not found in the internal map, operator signature will be lost. Batch may not reach quorum")
		*reply = 1
		return nil
	}

	agg.taskMutex.Lock()
	agg.AggregatorConfig.BaseConfig.Logger.Info("- Locked Resources: Starting processing of Response")
	taskIndex, ok := agg.batchesIdxByIdentifierHash[signedTaskResponse.BatchIdentifierHash]
	if !ok {
		agg.logger.Errorf("Unexpected error fetching for task with merkle root 0x%x", signedTaskResponse.BatchMerkleRoot)
		*reply = 1
		return nil
	}
	// Unlock
	agg.logger.Info("- Unlocked Resources: Task not found in the internal map")
	agg.taskMutex.Unlock()

	agg.telemetry.LogOperatorResponse(signedTaskResponse.BatchMerkleRoot, signedTaskResponse.OperatorId)

	// Don't wait infinitely if it can't answer
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure the cancel function is called to release resources

	// Create a channel to signal when the task is done
	done := make(chan struct{})

	agg.logger.Info("Starting bls signature process")
	go func() {
		err := agg.blsAggregationService.ProcessNewSignature(
			context.Background(), taskIndex, signedTaskResponse.BatchIdentifierHash,
			&signedTaskResponse.BlsSignature, signedTaskResponse.OperatorId,
		)

		if err != nil {
			agg.logger.Warnf("BLS aggregation service error: %s", err)
			// todo shouldn't we here close the channel with a reply = 1?
		} else {
			agg.logger.Info("BLS process succeeded")
		}

		close(done)
	}()

	*reply = 1
	// Wait for either the context to be done or the task to complete
	select {
	case <-ctx.Done():
		// The context's deadline was exceeded or it was canceled
		agg.logger.Info("Bls process timed out, operator signature will be lost. Batch may not reach quorum")
	case <-done:
		// The task completed successfully
		agg.logger.Info("Bls context finished correctly")
		*reply = 0
	}

	agg.AggregatorConfig.BaseConfig.Logger.Info("- Unlocked Resources: Task response processing finished")

	return nil
}

// Dummy method to check if the server is running
// TODO: Remove this method in prod
func (agg *Aggregator) ServerRunning(_ *struct{}, reply *int64) error {
	*reply = 1
	return nil
}
