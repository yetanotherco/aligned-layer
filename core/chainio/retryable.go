package chainio

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	servicemanager "github.com/yetanotherco/aligned_layer/contracts/bindings/AlignedLayerServiceManager"
	retry "github.com/yetanotherco/aligned_layer/core"
)

// |---AVS_WRITER---|

func RespondToTaskV2(w *AvsWriter, opts *bind.TransactOpts, batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) func() (*types.Transaction, error) {
	respondToTaskV2_func := func() (*types.Transaction, error) {
		// Try with main connection
		tx, err := w.AvsContractBindings.ServiceManager.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			// If error try with fallback
			tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		}
		return tx, err
	}
	return respondToTaskV2_func
}

/*
RespondToTaskV2Retryable
Send a transaction to the AVS contract to respond to a task.
- All errors are considered Transient Errors
- Retry times (3 retries): 12 sec (1 Blocks), 24 sec (2 Blocks), 48 sec (4 Blocks)
- NOTE: Contract call reverts are not considered `PermanentError`'s as block reorg's may lead to contract call revert in which case the aggregator should retry.
*/
func (w *AvsWriter) RespondToTaskV2Retryable(opts *bind.TransactOpts, batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return retry.RetryWithData(RespondToTaskV2(w, opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature), retry.ChainRetryConfig())
}

func BatchesState(w *AvsWriter, opts *bind.CallOpts, arg0 [32]byte) func() (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	batchesState_func := func() (struct {
		TaskCreatedBlock      uint32
		Responded             bool
		RespondToTaskFeeLimit *big.Int
	}, error) {
		// Try with main connection
		state, err := w.AvsContractBindings.ServiceManager.BatchesState(opts, arg0)
		if err != nil {
			// If error try with fallback connection
			state, err = w.AvsContractBindings.ServiceManagerFallback.BatchesState(opts, arg0)
		}
		return state, err
	}
	return batchesState_func
}

/*
BatchesStateRetryable
Get the state of a batch from the AVS contract.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec
*/
func (w *AvsWriter) BatchesStateRetryable(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	return retry.RetryWithData(BatchesState(w, opts, arg0), retry.EthCallRetryConfig())
}

func BatcherBalances(w *AvsWriter, opts *bind.CallOpts, senderAddress common.Address) func() (*big.Int, error) {
	batcherBalances_func := func() (*big.Int, error) {
		// Try with main connection
		batcherBalance, err := w.AvsContractBindings.ServiceManager.BatchersBalances(opts, senderAddress)
		if err != nil {
			// If error try with fallback connection
			batcherBalance, err = w.AvsContractBindings.ServiceManagerFallback.BatchersBalances(opts, senderAddress)
		}
		return batcherBalance, err
	}
	return batcherBalances_func
}

/*
BatcherBalancesRetryable
Get the balance of a batcher from the AVS contract.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec
*/
func (w *AvsWriter) BatcherBalancesRetryable(opts *bind.CallOpts, senderAddress common.Address) (*big.Int, error) {
	return retry.RetryWithData(BatcherBalances(w, opts, senderAddress), retry.EthCallRetryConfig())
}

func BalanceAt(w *AvsWriter, ctx context.Context, aggregatorAddress common.Address, blockNumber *big.Int) func() (*big.Int, error) {
	balanceAt_func := func() (*big.Int, error) {
		// Try with main connection
		aggregatorBalance, err := w.Client.BalanceAt(ctx, aggregatorAddress, blockNumber)
		if err != nil {
			// If error try with fallback connection
			aggregatorBalance, err = w.ClientFallback.BalanceAt(ctx, aggregatorAddress, blockNumber)
		}
		return aggregatorBalance, err
	}
	return balanceAt_func
}

/*
BalanceAtRetryable
Get the balance of aggregatorAddress at blockNumber.
If blockNumber is nil, it gets the latest balance.
TODO: it gets the balance from an Address, not necessarily an aggregator. The name of the parameter should be changed.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (w *AvsWriter) BalanceAtRetryable(ctx context.Context, aggregatorAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	return retry.RetryWithData(BalanceAt(w, ctx, aggregatorAddress, blockNumber), retry.EthCallRetryConfig())
}

// |---AVS_SUBSCRIBER---|

func BlockNumber(s *AvsSubscriber, ctx context.Context) func() (uint64, error) {
	latestBlock_func := func() (uint64, error) {
		// Try with main connection
		latestBlock, err := s.AvsContractBindings.ethClient.BlockNumber(ctx)
		if err != nil {
			// If error try with fallback connection
			latestBlock, err = s.AvsContractBindings.ethClientFallback.BlockNumber(ctx)
		}
		return latestBlock, err
	}
	return latestBlock_func
}

/*
BlockNumberRetryable
Get the latest block number from Ethereum
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) BlockNumberRetryable(ctx context.Context) (uint64, error) {
	return retry.RetryWithData(BlockNumber(s, ctx), retry.EthCallRetryConfig())
}

func FilterBatchV2(s *AvsSubscriber, opts *bind.FilterOpts, batchMerkleRoot [][32]byte) func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
		logs, err := s.AvsContractBindings.ServiceManager.FilterNewBatchV2(opts, batchMerkleRoot)
		if err != nil {
			logs, err = s.AvsContractBindings.ServiceManagerFallback.FilterNewBatchV2(opts, batchMerkleRoot)
		}
		return logs, err
	}
	return filterNewBatchV2_func
}

/*
FilterBatchV2Retryable
Get NewBatchV2 logs from the AVS contract.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) FilterBatchV2Retryable(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
	return retry.RetryWithData(FilterBatchV2(s, opts, batchMerkleRoot), retry.EthCallRetryConfig())
}

func FilterBatchV3(s *AvsSubscriber, opts *bind.FilterOpts, batchMerkleRoot [][32]byte) func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
	filterNewBatchV3_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
		logs, err := s.AvsContractBindings.ServiceManager.FilterNewBatchV3(opts, batchMerkleRoot)
		if err != nil {
			logs, err = s.AvsContractBindings.ServiceManagerFallback.FilterNewBatchV3(opts, batchMerkleRoot)
		}
		return logs, err
	}
	return filterNewBatchV3_func
}

/*
FilterBatchV3Retryable
Get NewBatchV3 logs from the AVS contract.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) FilterBatchV3Retryable(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
	return retry.RetryWithData(FilterBatchV3(s, opts, batchMerkleRoot), retry.EthCallRetryConfig())
}

func BatchState(s *AvsSubscriber, opts *bind.CallOpts, arg0 [32]byte) func() (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	batchState_func := func() (struct {
		TaskCreatedBlock      uint32
		Responded             bool
		RespondToTaskFeeLimit *big.Int
	}, error) {
		state, err := s.AvsContractBindings.ServiceManager.ContractAlignedLayerServiceManagerCaller.BatchesState(opts, arg0)
		if err != nil {
			state, err = s.AvsContractBindings.ServiceManagerFallback.BatchesState(opts, arg0)
		}
		return state, err
	}
	return batchState_func
}

/*
BatchesStateRetryable
Get the state of a batch from the AVS contract.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec
*/
func (s *AvsSubscriber) BatchesStateRetryable(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {

	return retry.RetryWithData(BatchState(s, opts, arg0), retry.EthCallRetryConfig())
}

func SubscribeNewHead(s *AvsSubscriber, ctx context.Context, c chan<- *types.Header) func() (ethereum.Subscription, error) {
	subscribeNewHead_func := func() (ethereum.Subscription, error) {
		// Try with main connection
		sub, err := s.AvsContractBindings.ethClient.SubscribeNewHead(ctx, c)
		if err != nil {
			// If error try with fallback connection
			sub, err = s.AvsContractBindings.ethClientFallback.SubscribeNewHead(ctx, c)
		}
		return sub, err
	}
	return subscribeNewHead_func
}

/*
SubscribeNewHeadRetryable
Subscribe to new heads from the Ethereum node.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) SubscribeNewHeadRetryable(ctx context.Context, c chan<- *types.Header) (ethereum.Subscription, error) {
	return retry.RetryWithData(SubscribeNewHead(s, ctx, c), retry.EthCallRetryConfig())
}

func SubscribeToNewTasksV2(
	opts *bind.WatchOpts,
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2,
	batchMerkleRoot [][32]byte,
) func() (event.Subscription, error) {
	subscribe_func := func() (event.Subscription, error) {
		return serviceManager.WatchNewBatchV2(opts, newTaskCreatedChan, batchMerkleRoot)
	}
	return subscribe_func
}

/*
SubscribeToNewTasksV2Retryable
Subscribe to NewBatchV2 logs from the AVS contract.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func SubscribeToNewTasksV2Retryable(
	opts *bind.WatchOpts,
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2,
	batchMerkleRoot [][32]byte,
) (event.Subscription, error) {
	return retry.RetryWithData(SubscribeToNewTasksV2(opts, serviceManager, newTaskCreatedChan, batchMerkleRoot), retry.EthCallRetryConfig())
}

func SubscribeToNewTasksV3(
	opts *bind.WatchOpts,
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3,
	batchMerkleRoot [][32]byte,
) func() (event.Subscription, error) {
	subscribe_func := func() (event.Subscription, error) {
		return serviceManager.WatchNewBatchV3(opts, newTaskCreatedChan, batchMerkleRoot)
	}
	return subscribe_func
}

/*
SubscribeToNewTasksV3Retryable
Subscribe to NewBatchV3 logs from the AVS contract.
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func SubscribeToNewTasksV3Retryable(
	opts *bind.WatchOpts,
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3,
	batchMerkleRoot [][32]byte,
) (event.Subscription, error) {
	return retry.RetryWithData(SubscribeToNewTasksV3(opts, serviceManager, newTaskCreatedChan, batchMerkleRoot), retry.EthCallRetryConfig())
}
