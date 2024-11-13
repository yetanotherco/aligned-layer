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

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 12 sec (1 Blocks), 24 sec (2 Blocks), 48 sec (4 Blocks)
- NOTE: Contract call reverts are not considered `PermanentError`'s as block reorg's may lead to contract call revert in which case the aggregator should retry.
*/
func (w *AvsWriter) RespondToTaskV2Retryable(opts *bind.TransactOpts, batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	respondToTaskV2_func := func() (*types.Transaction, error) {
		// Try with main connection
		tx, err := w.AvsContractBindings.ServiceManager.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			// If error try with fallback
			tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		}
		return tx, err
	}
	return retry.RetryWithData(respondToTaskV2_func, retry.MinDelayChain, retry.RetryFactor, retry.NumRetries, retry.MaxIntervalChain, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 12 sec (1 Blocks), 24 sec (2 Blocks), 48 sec (4 Blocks)
*/
func (w *AvsWriter) BatchesStateRetryable(opts *bind.CallOpts, arg0 [32]byte) (struct {
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
	return retry.RetryWithData(batchesState_func, retry.MinDelayChain, retry.RetryFactor, retry.NumRetries, retry.MaxIntervalChain, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 12 sec (1 Blocks), 24 sec (2 Blocks), 48 sec (4 Blocks)
*/
func (w *AvsWriter) BatcherBalancesRetryable(opts *bind.CallOpts, senderAddress common.Address) (*big.Int, error) {
	batcherBalances_func := func() (*big.Int, error) {
		// Try with main connection
		batcherBalance, err := w.AvsContractBindings.ServiceManager.BatchersBalances(opts, senderAddress)
		if err != nil {
			// If error try with fallback connection
			batcherBalance, err = w.AvsContractBindings.ServiceManagerFallback.BatchersBalances(opts, senderAddress)
		}
		return batcherBalance, err
	}
	return retry.RetryWithData(batcherBalances_func, retry.MinDelayChain, retry.RetryFactor, retry.NumRetries, retry.MaxIntervalChain, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (w *AvsWriter) BalanceAtRetryable(ctx context.Context, aggregatorAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	balanceAt_func := func() (*big.Int, error) {
		// Try with main connection
		aggregatorBalance, err := w.Client.BalanceAt(ctx, aggregatorAddress, blockNumber)
		if err != nil {
			// If error try with fallback connection
			aggregatorBalance, err = w.ClientFallback.BalanceAt(ctx, aggregatorAddress, blockNumber)
		}
		return aggregatorBalance, err
	}
	return retry.RetryWithData(balanceAt_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

// |---AVS_SUBSCRIBER---|

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) BlockNumberRetryable(ctx context.Context) (uint64, error) {
	latestBlock_func := func() (uint64, error) {
		// Try with main connection
		latestBlock, err := s.AvsContractBindings.ethClient.BlockNumber(ctx)
		if err != nil {
			// If error try with fallback connection
			latestBlock, err = s.AvsContractBindings.ethClientFallback.BlockNumber(ctx)
		}
		return latestBlock, err
	}
	return retry.RetryWithData(latestBlock_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) FilterBatchV2Retryable(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
		return s.AvsContractBindings.ServiceManager.FilterNewBatchV2(opts, batchMerkleRoot)
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) FilterBatchV3Retryable(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
		return s.AvsContractBindings.ServiceManager.FilterNewBatchV3(opts, batchMerkleRoot)
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 12 sec (1 Blocks), 24 sec (2 Blocks), 48 sec (4 Blocks)
*/
func (s *AvsSubscriber) BatchesStateRetryable(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	batchState_func := func() (struct {
		TaskCreatedBlock      uint32
		Responded             bool
		RespondToTaskFeeLimit *big.Int
	}, error) {
		return s.AvsContractBindings.ServiceManager.ContractAlignedLayerServiceManagerCaller.BatchesState(opts, arg0)
	}

	return retry.RetryWithData(batchState_func, retry.MinDelayChain, retry.RetryFactor, retry.NumRetries, retry.MaxIntervalChain, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (s *AvsSubscriber) SubscribeNewHeadRetryable(ctx context.Context, c chan<- *types.Header) (ethereum.Subscription, error) {
	subscribeNewHead_func := func() (ethereum.Subscription, error) {
		// Try with main connection
		sub, err := s.AvsContractBindings.ethClient.SubscribeNewHead(ctx, c)
		if err != nil {
			// If error try with fallback connection
			sub, err = s.AvsContractBindings.ethClientFallback.SubscribeNewHead(ctx, c)
		}
		return sub, err
	}
	return retry.RetryWithData(subscribeNewHead_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func SubscribeToNewTasksV2Retryable(
	opts *bind.WatchOpts,
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2,
	batchMerkleRoot [][32]byte,
) (event.Subscription, error) {
	subscribe_func := func() (event.Subscription, error) {
		return serviceManager.WatchNewBatchV2(opts, newTaskCreatedChan, batchMerkleRoot)
	}
	return retry.RetryWithData(subscribe_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func SubscribeToNewTasksV3Retryable(
	opts *bind.WatchOpts,
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3,
	batchMerkleRoot [][32]byte,
) (event.Subscription, error) {
	subscribe_func := func() (event.Subscription, error) {
		return serviceManager.WatchNewBatchV3(opts, newTaskCreatedChan, batchMerkleRoot)
	}
	return retry.RetryWithData(subscribe_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

// |---AVS_READER---|

// TODO: These functions are being copied from AvsSubscriber and should be refactorized
// we don't actually need access to the AvsReader, AvsSubscriber or AbsWriter, but instead to the AvsContractBindings

// TODO: We should also add the fallback calls to the functions which are missing it

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (r *AvsReader) BlockNumberRetryable(ctx context.Context) (uint64, error) {
	latestBlock_func := func() (uint64, error) {
		// Try with main connection
		latestBlock, err := r.AvsContractBindings.ethClient.BlockNumber(ctx)
		if err != nil {
			// If error try with fallback connection
			latestBlock, err = r.AvsContractBindings.ethClientFallback.BlockNumber(ctx)
		}
		return latestBlock, err
	}
	return retry.RetryWithData(latestBlock_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 1 sec, 2 sec, 4 sec.
*/
func (r *AvsReader) FilterBatchV3Retryable(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
		// Try with main connection
		batch, err := r.AvsContractBindings.ServiceManager.FilterNewBatchV3(opts, batchMerkleRoot)
		if err != nil {
			// If error try with fallback connection
			batch, err = r.AvsContractBindings.ServiceManagerFallback.FilterNewBatchV3(opts, batchMerkleRoot)
		}
		return batch, err
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

/*
- All errors are considered Transient Errors
- Retry times (3 retries): 12 sec (1 Blocks), 24 sec (2 Blocks), 48 sec (4 Blocks)
*/
func (r *AvsReader) BatchesStateRetryable(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	batchState_func := func() (struct {
		TaskCreatedBlock      uint32
		Responded             bool
		RespondToTaskFeeLimit *big.Int
	}, error) {
		// Try with main connection
		state, err := r.AvsContractBindings.ServiceManager.ContractAlignedLayerServiceManagerCaller.BatchesState(opts, arg0)
		if err != nil {
			// If error try with fallback connection
			state, err = r.AvsContractBindings.ServiceManagerFallback.ContractAlignedLayerServiceManagerCaller.BatchesState(opts, arg0)
		}
		return state, err
	}

	return retry.RetryWithData(batchState_func, retry.MinDelayChain, retry.RetryFactor, retry.NumRetries, retry.MaxIntervalChain, retry.MaxElapsedTime)
}
