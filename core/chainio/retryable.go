package chainio

import (
	"context"
	"fmt"
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
	var (
		tx  *types.Transaction
		err error
	)
	respondToTaskV2_func := func() (*types.Transaction, error) {
		tx, err = w.AvsContractBindings.ServiceManager.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
			if err != nil {
				err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
			}
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
	var (
		state struct {
			TaskCreatedBlock      uint32
			Responded             bool
			RespondToTaskFeeLimit *big.Int
		}
		err error
	)

	batchesState_func := func() (struct {
		TaskCreatedBlock      uint32
		Responded             bool
		RespondToTaskFeeLimit *big.Int
	}, error) {
		state, err = w.AvsContractBindings.ServiceManager.BatchesState(&bind.CallOpts{}, arg0)
		if err != nil {
			state, err = w.AvsContractBindings.ServiceManagerFallback.BatchesState(&bind.CallOpts{}, arg0)
			if err != nil {
				err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
			}
		}
		return state, err
	}
	return retry.RetryWithData(batchesState_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (w *AvsWriter) BatcherBalancesRetryable(senderAddress common.Address) (*big.Int, error) {
	var (
		batcherBalance *big.Int
		err            error
	)
	batcherBalances_func := func() (*big.Int, error) {
		batcherBalance, err = w.AvsContractBindings.ServiceManager.BatchersBalances(&bind.CallOpts{}, senderAddress)
		if err != nil {
			batcherBalance, err = w.AvsContractBindings.ServiceManagerFallback.BatchersBalances(&bind.CallOpts{}, senderAddress)
			if err != nil {
				err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
			}
		}
		return batcherBalance, err
	}
	return retry.RetryWithData(batcherBalances_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (w *AvsWriter) BalanceAtRetryable(ctx context.Context, aggregatorAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	var (
		aggregatorBalance *big.Int
		err               error
	)
	balanceAt_func := func() (*big.Int, error) {
		aggregatorBalance, err = w.Client.BalanceAt(ctx, aggregatorAddress, blockNumber)
		if err != nil {
			aggregatorBalance, err = w.ClientFallback.BalanceAt(ctx, aggregatorAddress, blockNumber)
			if err != nil {
				err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
			}
		}
		return aggregatorBalance, err
	}
	return retry.RetryWithData(balanceAt_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

// |---AVS_SUBSCRIBER---|

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
				err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
			}
		}
		return latestBlock, err
	}
	return retry.RetryWithData(latestBlock_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (s *AvsSubscriber) FilterBatchV2Retryable(fromBlock uint64, ctx context.Context) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
	var (
		logs *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator
		err  error
	)

	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {
		logs, err = s.AvsContractBindings.ServiceManager.FilterNewBatchV2(&bind.FilterOpts{Start: fromBlock, End: nil, Context: ctx}, nil)
		if err != nil {
			fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
		}
		return logs, err
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}

func (s *AvsSubscriber) FilterBatchV3Retryable(fromBlock uint64, ctx context.Context) (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
	var (
		logs *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator
		err  error
	)
	filterNewBatchV2_func := func() (*servicemanager.ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {
		logs, err = s.AvsContractBindings.ServiceManager.FilterNewBatchV3(&bind.FilterOpts{Start: fromBlock, End: nil, Context: ctx}, nil)
		if err != nil {
			err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
		}
		return logs, err
	}
	return retry.RetryWithData(filterNewBatchV2_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
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
			err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
		}
		return state, err
	}

	return retry.RetryWithData(batchState_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
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
				err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
			}
		}
		return sub, err
	}
	return retry.RetryWithData(subscribeNewHead_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
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
			err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
		}
		return sub, err
	}
	return retry.RetryWithData(subscribe_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
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
			err = fmt.Errorf("Transient error: Unexpected Error while retrying: %s\n", err)
		}
		return sub, err
	}
	return retry.RetryWithData(subscribe_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
}
