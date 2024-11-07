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

func (w *AvsWriter) RespondToTaskV2Retryable(opts *bind.TransactOpts, batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	respondToTaskV2_func := func() (*types.Transaction, error) {
		tx, err := w.AvsContractBindings.ServiceManager.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		}
		return tx, err
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
		state, err := w.AvsContractBindings.ServiceManager.BatchesState(opts, arg0)
		if err != nil {
			state, err = w.AvsContractBindings.ServiceManagerFallback.BatchesState(opts, arg0)
		}
		return state, err
	}
	return retry.RetryWithData(batchesState_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (w *AvsWriter) BatcherBalancesRetryable(opts *bind.CallOpts, senderAddress common.Address) (*big.Int, error) {
	batcherBalances_func := func() (*big.Int, error) {
		batcherBalance, err := w.AvsContractBindings.ServiceManager.BatchersBalances(opts, senderAddress)
		if err != nil {
			batcherBalance, err = w.AvsContractBindings.ServiceManagerFallback.BatchersBalances(opts, senderAddress)
		}
		return batcherBalance, err
	}
	return retry.RetryWithData(batcherBalances_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (w *AvsWriter) BalanceAtRetryable(ctx context.Context, aggregatorAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	balanceAt_func := func() (*big.Int, error) {
		aggregatorBalance, err := w.Client.BalanceAt(ctx, aggregatorAddress, blockNumber)
		if err != nil {
			aggregatorBalance, err = w.ClientFallback.BalanceAt(ctx, aggregatorAddress, blockNumber)
		}
		return aggregatorBalance, err
	}
	return retry.RetryWithData(balanceAt_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

// |---AVS_SUBSCRIBER---|

func (s *AvsSubscriber) BlockNumberRetryable(ctx context.Context) (uint64, error) {
	latestBlock_func := func() (uint64, error) {
		latestBlock, err := s.AvsContractBindings.ethClient.BlockNumber(ctx)
		if err != nil {
			latestBlock, err = s.AvsContractBindings.ethClientFallback.BlockNumber(ctx)
		}
		return latestBlock, err
	}
	return retry.RetryWithData(latestBlock_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (s *AvsSubscriber) FilterBatchV2Retryable(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
		return s.AvsContractBindings.ServiceManager.FilterNewBatchV2(opts, batchMerkleRoot)
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (s *AvsSubscriber) FilterBatchV3Retryable(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
		return s.AvsContractBindings.ServiceManager.FilterNewBatchV3(opts, batchMerkleRoot)
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
		sub, err := s.AvsContractBindings.ethClient.SubscribeNewHead(ctx, c)
		if err != nil {
			sub, err = s.AvsContractBindings.ethClientFallback.SubscribeNewHead(ctx, c)
		}
		return sub, err
	}
	return retry.RetryWithData(subscribeNewHead_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func SubscribeToNewTasksV2Retrayable(
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
