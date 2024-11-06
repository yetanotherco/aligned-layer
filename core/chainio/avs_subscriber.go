package chainio

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	servicemanager "github.com/yetanotherco/aligned_layer/contracts/bindings/AlignedLayerServiceManager"
	"github.com/yetanotherco/aligned_layer/core/config"

	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/crypto"
	connection "github.com/yetanotherco/aligned_layer/core"
)

const (
	MaxRetries                        = 100
	RetryInterval                     = 1 * time.Second
	BlockInterval              uint64 = 1000
	PollLatestBatchInterval           = 5 * time.Second
	RemoveBatchFromSetInterval        = 5 * time.Minute
)

// NOTE(marian): Leaving this commented code here as it may be useful in the short term.
// type AvsSubscriberer interface {
// 	SubscribeToNewTasks(newTaskCreatedChan chan *cstaskmanager.ContractAlignedLayerTaskManagerNewTaskCreated) event.Subscription
// 	SubscribeToTaskResponses(taskResponseLogs chan *cstaskmanager.ContractAlignedLayerTaskManagerTaskResponded) event.Subscription
// 	ParseTaskResponded(rawLog types.Log) (*cstaskmanager.ContractAlignedLayerTaskManagerTaskResponded, error)
// }

// Subscribers use a ws connection instead of http connection like Readers
// kind of stupid that the geth client doesn't have a unified interface for both...
// it takes a single url, so the bindings, even though they have watcher functions, those can't be used
// with the http connection... seems very very stupid. Am I missing something?
type AvsSubscriber struct {
	AvsContractBindings            *AvsServiceBindings
	AlignedLayerServiceManagerAddr ethcommon.Address
	logger                         sdklogging.Logger
}

func NewAvsSubscriberFromConfig(baseConfig *config.BaseConfig) (*AvsSubscriber, error) {
	avsContractBindings, err := NewAvsServiceBindings(
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr,
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr,
		baseConfig.EthWsClient, baseConfig.EthWsClientFallback, baseConfig.Logger)

	if err != nil {
		baseConfig.Logger.Errorf("Failed to create contract bindings", "err", err)
		return nil, err
	}

	return &AvsSubscriber{
		AvsContractBindings:            avsContractBindings,
		AlignedLayerServiceManagerAddr: baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr,
		logger:                         baseConfig.Logger,
	}, nil
}

func (s *AvsSubscriber) SubscribeToNewTasksV2(newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2) (chan error, error) {
	// Create a new channel to receive new tasks
	internalChannel := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2)

	// Subscribe to new tasks
	// NOTE: Should we consolidate these subscribers behind two interfaces
	sub, err := SubscribeToNewTasksV2Retrayable(s.AvsContractBindings.ServiceManager, internalChannel, s.logger)
	if err != nil {
		s.logger.Error("Failed to subscribe to new AlignedLayer tasks after %d retries", connection.NumRetries, "err", err)
		return nil, err
	}
	s.logger.Info("Subscribed to new AlignedLayer V2 tasks")

	subFallback, err := SubscribeToNewTasksV2Retrayable(s.AvsContractBindings.ServiceManagerFallback, internalChannel, s.logger)
	if err != nil {
		s.logger.Error("Failed to subscribe to new AlignedLayer tasks %d retries", connection.NumRetries, "err", err)
		return nil, err
	}
	s.logger.Info("Subscribed to new AlignedLayer V2 tasks")

	// create a new channel to foward errors
	errorChannel := make(chan error)

	pollLatestBatchTicker := time.NewTicker(PollLatestBatchInterval)

	// Forward the new tasks to the provided channel
	go func() {
		defer pollLatestBatchTicker.Stop()
		newBatchMutex := &sync.Mutex{}
		batchesSet := make(map[[32]byte]struct{})
		for {
			select {
			case newBatch := <-internalChannel:
				s.processNewBatchV2(newBatch, batchesSet, newBatchMutex, newTaskCreatedChan)
			case <-pollLatestBatchTicker.C:
				latestBatch, err := s.getLatestNotRespondedTaskFromEthereumV2()
				if err != nil {
					s.logger.Debug("Failed to get latest task from blockchain", "err", err)
					continue
				}
				if latestBatch != nil {
					s.processNewBatchV2(latestBatch, batchesSet, newBatchMutex, newTaskCreatedChan)
				}
			}
		}

	}()

	// Handle errors and resubscribe
	go func() {
		for {
			select {
			case err := <-sub.Err():
				s.logger.Warn("Error in new task subscription", "err", err)
				sub.Unsubscribe()
				sub, err = SubscribeToNewTasksV2Retrayable(s.AvsContractBindings.ServiceManager, internalChannel, s.logger)
				if err != nil {
					errorChannel <- err
				}
			case err := <-subFallback.Err():
				s.logger.Warn("Error in fallback new task subscription", "err", err)
				subFallback.Unsubscribe()
				subFallback, err = SubscribeToNewTasksV2Retrayable(s.AvsContractBindings.ServiceManagerFallback, internalChannel, s.logger)
				if err != nil {
					errorChannel <- err
				}
			}
		}
	}()

	return errorChannel, nil
}

func (s *AvsSubscriber) SubscribeToNewTasksV3(newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3) (chan error, error) {
	// Create a new channel to receive new tasks
	internalChannel := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3)

	// Subscribe to new tasks
	sub, err := SubscribeToNewTasksV3Retryable(s.AvsContractBindings.ServiceManager, internalChannel, s.logger)
	if err != nil {
		s.logger.Error("Failed to subscribe to new AlignedLayer tasks after %d retries", MaxRetries, "err", err)
		return nil, err
	}
	s.logger.Info("Subscribed to new AlignedLayer V3 tasks")

	//TODO: collapse this
	subFallback, err := SubscribeToNewTasksV3Retryable(s.AvsContractBindings.ServiceManagerFallback, internalChannel, s.logger)
	if err != nil {
		s.logger.Error("Failed to subscribe to new AlignedLayer %d tasks", MaxRetries, "err", err)
		return nil, err
	}
	s.logger.Info("Subscribed to new AlignedLayer V3 tasks")

	// create a new channel to foward errors
	errorChannel := make(chan error)

	pollLatestBatchTicker := time.NewTicker(PollLatestBatchInterval)

	// Forward the new tasks to the provided channel
	go func() {
		defer pollLatestBatchTicker.Stop()
		newBatchMutex := &sync.Mutex{}
		batchesSet := make(map[[32]byte]struct{})
		for {
			select {
			case newBatch := <-internalChannel:
				s.processNewBatchV3(newBatch, batchesSet, newBatchMutex, newTaskCreatedChan)
			case <-pollLatestBatchTicker.C:
				latestBatch, err := s.getLatestNotRespondedTaskFromEthereumV3()
				if err != nil {
					s.logger.Debug("Failed to get latest task from blockchain", "err", err)
					continue
				}
				if latestBatch != nil {
					s.processNewBatchV3(latestBatch, batchesSet, newBatchMutex, newTaskCreatedChan)
				}
			}
		}

	}()

	// Handle errors and resubscribe
	go func() {
		for {
			select {
			case err := <-sub.Err():
				s.logger.Warn("Error in new task subscription", "err", err)
				sub.Unsubscribe()
				sub, err = SubscribeToNewTasksV3Retryable(s.AvsContractBindings.ServiceManager, internalChannel, s.logger)
				if err != nil {
					errorChannel <- err
				}
			case err := <-subFallback.Err():
				s.logger.Warn("Error in fallback new task subscription", "err", err)
				subFallback.Unsubscribe()
				subFallback, err = SubscribeToNewTasksV3Retryable(s.AvsContractBindings.ServiceManagerFallback, internalChannel, s.logger)
				if err != nil {
					errorChannel <- err
				}
			}
		}
	}()

	return errorChannel, nil
}

func (s *AvsSubscriber) processNewBatchV2(batch *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2, batchesSet map[[32]byte]struct{}, newBatchMutex *sync.Mutex, newTaskCreatedChan chan<- *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2) {
	newBatchMutex.Lock()
	defer newBatchMutex.Unlock()

	batchIdentifier := append(batch.BatchMerkleRoot[:], batch.SenderAddress[:]...)
	var batchIdentifierHash = *(*[32]byte)(crypto.Keccak256(batchIdentifier))

	if _, ok := batchesSet[batchIdentifierHash]; !ok {
		s.logger.Info("Received new task",
			"batchMerkleRoot", hex.EncodeToString(batch.BatchMerkleRoot[:]),
			"senderAddress", hex.EncodeToString(batch.SenderAddress[:]),
			"batchIdentifierHash", hex.EncodeToString(batchIdentifierHash[:]))

		batchesSet[batchIdentifierHash] = struct{}{}
		newTaskCreatedChan <- batch

		// Remove the batch from the set after RemoveBatchFromSetInterval time
		go func() {
			time.Sleep(RemoveBatchFromSetInterval)
			newBatchMutex.Lock()
			delete(batchesSet, batchIdentifierHash)
			newBatchMutex.Unlock()
		}()
	}
}

func (s *AvsSubscriber) processNewBatchV3(batch *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3, batchesSet map[[32]byte]struct{}, newBatchMutex *sync.Mutex, newTaskCreatedChan chan<- *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3) {
	newBatchMutex.Lock()
	defer newBatchMutex.Unlock()

	batchIdentifier := append(batch.BatchMerkleRoot[:], batch.SenderAddress[:]...)
	var batchIdentifierHash = *(*[32]byte)(crypto.Keccak256(batchIdentifier))

	if _, ok := batchesSet[batchIdentifierHash]; !ok {
		s.logger.Info("Received new task",
			"batchMerkleRoot", hex.EncodeToString(batch.BatchMerkleRoot[:]),
			"senderAddress", hex.EncodeToString(batch.SenderAddress[:]),
			"batchIdentifierHash", hex.EncodeToString(batchIdentifierHash[:]))

		batchesSet[batchIdentifierHash] = struct{}{}
		newTaskCreatedChan <- batch

		// Remove the batch from the set after RemoveBatchFromSetInterval time
		go func() {
			time.Sleep(RemoveBatchFromSetInterval)
			newBatchMutex.Lock()
			delete(batchesSet, batchIdentifierHash)
			newBatchMutex.Unlock()
		}()
	}
}

// getLatestNotRespondedTaskFromEthereum queries the blockchain for the latest not responded task using the FilterNewBatch method.
func (s *AvsSubscriber) getLatestNotRespondedTaskFromEthereumV2() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2, error) {

	latestBlock, err := s.BlockNumberRetryable(context.Background())
	if err != nil {
		return nil, err
	}

	var fromBlock uint64

	if latestBlock < BlockInterval {
		fromBlock = 0
	} else {
		fromBlock = latestBlock - BlockInterval
	}

	logs, err := s.FilterBatchV2Retryable(fromBlock, context.Background())
	if err != nil {
		return nil, err
	}

	var lastLog *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2

	// Iterate over the logs until the end
	for logs.Next() {
		lastLog = logs.Event
	}

	if err := logs.Error(); err != nil {
		return nil, err
	}

	if lastLog == nil {
		return nil, nil
	}

	batchIdentifier := append(lastLog.BatchMerkleRoot[:], lastLog.SenderAddress[:]...)
	batchIdentifierHash := *(*[32]byte)(crypto.Keccak256(batchIdentifier))
	state, err := s.BatchesStateRetryable(nil, batchIdentifierHash)
	if err != nil {
		return nil, err
	}

	if state.Responded {
		return nil, nil
	}

	return lastLog, nil
}

// getLatestNotRespondedTaskFromEthereum queries the blockchain for the latest not responded task using the FilterNewBatch method.
func (s *AvsSubscriber) getLatestNotRespondedTaskFromEthereumV3() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3, error) {
	latestBlock, err := s.BlockNumberRetryable(context.Background())
	if err != nil {
		return nil, err
	}

	var fromBlock uint64

	if latestBlock < BlockInterval {
		fromBlock = 0
	} else {
		fromBlock = latestBlock - BlockInterval
	}

	logs, err := s.FilterBatchV3Retryable(fromBlock, context.Background())
	if err != nil {
		return nil, err
	}

	var lastLog *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3

	// Iterate over the logs until the end
	for logs.Next() {
		lastLog = logs.Event
	}

	if err := logs.Error(); err != nil {
		return nil, err
	}

	if lastLog == nil {
		return nil, nil
	}

	batchIdentifier := append(lastLog.BatchMerkleRoot[:], lastLog.SenderAddress[:]...)
	batchIdentifierHash := *(*[32]byte)(crypto.Keccak256(batchIdentifier))
	state, err := s.BatchesStateRetryable(nil, batchIdentifierHash)
	if err != nil {
		return nil, err
	}

	if state.Responded {
		return nil, nil
	}

	return lastLog, nil
}

func (s *AvsSubscriber) WaitForOneBlock(startBlock uint64) error {
	currentBlock, err := s.BlockNumberRetryable(context.Background())
	if err != nil {
		return err
	}

	if currentBlock <= startBlock { // should really be == but just in case
		// Subscribe to new head
		c := make(chan *types.Header)
		sub, err := s.SubscribeNewHeadRetryable(context.Background(), c)
		if err != nil {
			return err
		}

		// Read channel for the new block
		<-c
		(sub).Unsubscribe()
	}

	return nil
}

// func (s *AvsSubscriber) SubscribeToTaskResponses(taskResponseChan chan *cstaskmanager.ContractAlignedLayerTaskManagerTaskResponded) event.Subscription {
// 	sub, err := s.AvsContractBindings.TaskManager.WatchTaskResponded(
// 		&bind.WatchOpts{}, taskResponseChan,
// 	)
// 	if err != nil {
// 		s.logger.Error("Failed to subscribe to TaskResponded events", "err", err)
// 	}
// 	s.logger.Infof("Subscribed to TaskResponded events")
// 	return sub
// }

// func (s *AvsSubscriber) ParseTaskResponded(rawLog types.Log) (*cstaskmanager.ContractAlignedLayerTaskManagerTaskResponded, error) {
// 	return s.AvsContractBindings.TaskManager.ContractAlignedLayerTaskManagerFilterer.ParseTaskResponded(rawLog)
// }

// |---RETRYABLE---|

func (s *AvsSubscriber) BlockNumberRetryable(ctx context.Context) (uint64, error) {
	var (
		latestBlock uint64
		err         error
	)
	latestBlock_func := func() (uint64, error) {
		latestBlock, err = s.AvsContractBindings.ethClient.BlockNumber(ctx)
		if err != nil {
			latestBlock, err = s.AvsContractBindings.ethClientFallback.BlockNumber(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "connect: connection refused") {
					err = connection.TransientError{Inner: err}
					return latestBlock, err
				}
				if strings.Contains(err.Error(), "read: connection reset by peer") {
					return latestBlock, connection.TransientError{Inner: err}
				}
				err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
			}
		}
		return latestBlock, err
	}
	return connection.RetryWithData(latestBlock_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func (s *AvsSubscriber) FilterBatchV2Retryable(fromBlock uint64, ctx context.Context) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
	var (
		logs *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator
		err  error
	)

	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
		logs, err = s.AvsContractBindings.ServiceManager.FilterNewBatchV2(&bind.FilterOpts{Start: fromBlock, End: nil, Context: ctx}, nil)
		if err != nil {
			// Note return type will be nil
			if strings.Contains(err.Error(), "connect: connection refused") {
				err = connection.TransientError{Inner: err}
				return logs, err
			}
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				return logs, connection.TransientError{Inner: err}
			}
			err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return logs, err
	}
	return connection.RetryWithData(filterNewBatchV2_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func (s *AvsSubscriber) FilterBatchV3Retryable(fromBlock uint64, ctx context.Context) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
	var (
		logs *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator
		err  error
	)
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
		logs, err = s.AvsContractBindings.ServiceManager.FilterNewBatchV3(&bind.FilterOpts{Start: fromBlock, End: nil, Context: ctx}, nil)
		if err != nil {
			// Note return type will be nil
			if strings.Contains(err.Error(), "connect: connection refused") {
				err = connection.TransientError{Inner: err}
				return logs, err
			}
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				return logs, connection.TransientError{Inner: err}
			}
			err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return logs, err
	}
	return connection.RetryWithData(filterNewBatchV2_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func (s *AvsSubscriber) BatchesStateRetryable(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	var (
		state struct {
			TaskCreatedBlock      uint32
			Responded             bool
			RespondToTaskFeeLimit *big.Int
		}
		err error
	)
	batchState_func := func() (struct {
		TaskCreatedBlock      uint32
		Responded             bool
		RespondToTaskFeeLimit *big.Int
	}, error) {
		state, err = s.AvsContractBindings.ServiceManager.ContractAlignedLayerServiceManagerCaller.BatchesState(opts, arg0)
		if err != nil {
			// Note return type will be nil
			if strings.Contains(err.Error(), "connect: connection refused") {
				err = connection.TransientError{Inner: err}
				return state, err
			}
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				return state, connection.TransientError{Inner: err}
			}
			err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return state, err
	}

	return connection.RetryWithData(batchState_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func (s *AvsSubscriber) SubscribeNewHeadRetryable(ctx context.Context, c chan<- *types.Header) (ethereum.Subscription, error) {
	var (
		sub ethereum.Subscription
		err error
	)
	subscribeNewHead_func := func() (ethereum.Subscription, error) {
		sub, err = s.AvsContractBindings.ethClient.SubscribeNewHead(ctx, c)
		if err != nil {
			sub, err = s.AvsContractBindings.ethClientFallback.SubscribeNewHead(ctx, c)
			if err != nil {
				// Note return type will be nil
				if strings.Contains(err.Error(), "connect: connection refused") {
					err = connection.TransientError{Inner: err}
					return sub, err
				}
				if strings.Contains(err.Error(), "read: connection reset by peer") {
					return sub, connection.TransientError{Inner: err}
				}
				err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
			}
		}
		return sub, err
	}
	return connection.RetryWithData(subscribeNewHead_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func SubscribeToNewTasksV2Retrayable(
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2,
	logger sdklogging.Logger,
) (event.Subscription, error) {
	var (
		sub event.Subscription
		err error
	)
	subscribe_func := func() (event.Subscription, error) {
		sub, err = serviceManager.WatchNewBatchV2(&bind.WatchOpts{}, newTaskCreatedChan, nil)
		if err != nil {
			// Note return type will be nil
			if strings.Contains(err.Error(), "connect: connection refused") {
				err = connection.TransientError{Inner: err}
				return sub, err
			}
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				return sub, connection.TransientError{Inner: err}
			}
			err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return sub, err
	}
	return connection.RetryWithData(subscribe_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func SubscribeToNewTasksV3Retryable(
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3,
	logger sdklogging.Logger,
) (event.Subscription, error) {
	var (
		sub event.Subscription
		err error
	)
	subscribe_func := func() (event.Subscription, error) {
		sub, err = serviceManager.WatchNewBatchV3(&bind.WatchOpts{}, newTaskCreatedChan, nil)
		if err != nil {
			if strings.Contains(err.Error(), "connect: connection refused") {
				err = connection.TransientError{Inner: err}
				return sub, err
			}
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				return sub, connection.TransientError{Inner: err}
			}
			err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return sub, err
	}
	return connection.RetryWithData(subscribe_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}
