package chainio

import (
	"context"
	"math/big"

	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	servicemanager "github.com/yetanotherco/aligned_layer/contracts/bindings/AlignedLayerServiceManager"
	retry "github.com/yetanotherco/aligned_layer/core"
)

// |---AVS_WRITER---|

func (w *AvsWriter) RespondToTaskV2Retryable(opts *bind.TransactOpts, batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	respondToTaskV2_func := func() (*types.Transaction, error) {
		return w.AvsContractBindings.ServiceManager.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
	}
	return retry.RetryWithData(respondToTaskV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

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
		return w.AvsContractBindings.ServiceManager.BatchesState(&bind.CallOpts{}, arg0)
	}
	return retry.RetryWithData(batchesState_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (w *AvsWriter) BatcherBalancesRetryable(senderAddress common.Address) (*big.Int, error) {
	batcherBalances_func := func() (*big.Int, error) {
		return w.AvsContractBindings.ServiceManager.BatchersBalances(&bind.CallOpts{}, senderAddress)
	}
	return retry.RetryWithData(batcherBalances_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (w *AvsWriter) BalanceAtRetryable(ctx context.Context, aggregatorAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	balanceAt_func := func() (*big.Int, error) {
		return w.Client.BalanceAt(ctx, aggregatorAddress, blockNumber)
	}
	return retry.RetryWithData(balanceAt_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

// |---AVS_SUBSCRIBER---|

func (s *AvsSubscriber) BlockNumberRetryable(ctx context.Context) (uint64, error) {
	latestBlock_func := func() (uint64, error) {
		return s.AvsContractBindings.ethClient.BlockNumber(ctx)
	}
	return retry.RetryWithData(latestBlock_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (s *AvsSubscriber) FilterBatchV2Retryable(fromBlock uint64, ctx context.Context) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
		return s.AvsContractBindings.ServiceManager.FilterNewBatchV2(&bind.FilterOpts{Start: fromBlock, End: nil, Context: ctx}, nil)
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (s *AvsSubscriber) FilterBatchV3Retryable(fromBlock uint64, ctx context.Context) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {

	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
		return s.AvsContractBindings.ServiceManager.FilterNewBatchV3(&bind.FilterOpts{Start: fromBlock, End: nil, Context: ctx}, nil)
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

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

	return retry.RetryWithData(batchState_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (s *AvsSubscriber) SubscribeNewHeadRetryable(ctx context.Context, c chan<- *types.Header) (ethereum.Subscription, error) {
	subscribeNewHead_func := func() (ethereum.Subscription, error) {
		return s.AvsContractBindings.ethClient.SubscribeNewHead(ctx, c)
	}
	return retry.RetryWithData(subscribeNewHead_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func SubscribeToNewTasksV2Retrayable(
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2,
	logger sdklogging.Logger,
) (event.Subscription, error) {
	subscribe_func := func() (event.Subscription, error) {
		return serviceManager.WatchNewBatchV2(&bind.WatchOpts{}, newTaskCreatedChan, nil)
	}
	return retry.RetryWithData(subscribe_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func SubscribeToNewTasksV3Retryable(
	serviceManager *servicemanager.ContractAlignedLayerServiceManager,
	newTaskCreatedChan chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3,
	logger sdklogging.Logger,
) (event.Subscription, error) {
	subscribe_func := func() (event.Subscription, error) {
		return serviceManager.WatchNewBatchV3(&bind.WatchOpts{}, newTaskCreatedChan, nil)
	}
	return retry.RetryWithData(subscribe_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}
