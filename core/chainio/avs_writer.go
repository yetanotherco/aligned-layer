package chainio

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	servicemanager "github.com/yetanotherco/aligned_layer/contracts/bindings/AlignedLayerServiceManager"
	connection "github.com/yetanotherco/aligned_layer/core"
	"github.com/yetanotherco/aligned_layer/core/config"
)

type AvsWriter struct {
	*avsregistry.ChainWriter
	AvsContractBindings *AvsServiceBindings
	logger              logging.Logger
	Signer              signer.Signer
	Client              eth.InstrumentedClient
	ClientFallback      eth.InstrumentedClient
}

func NewAvsWriterFromConfig(baseConfig *config.BaseConfig, ecdsaConfig *config.EcdsaConfig) (*AvsWriter, error) {

	buildAllConfig := clients.BuildAllConfig{
		EthHttpUrl:                 baseConfig.EthRpcUrl,
		EthWsUrl:                   baseConfig.EthWsUrl,
		RegistryCoordinatorAddr:    baseConfig.AlignedLayerDeploymentConfig.AlignedLayerRegistryCoordinatorAddr.String(),
		OperatorStateRetrieverAddr: baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr.String(),
		AvsName:                    "AlignedLayer",
		PromMetricsIpPortAddress:   baseConfig.EigenMetricsIpPortAddress,
	}

	clients, err := clients.BuildAll(buildAllConfig, ecdsaConfig.PrivateKey, baseConfig.Logger)

	if err != nil {
		baseConfig.Logger.Error("Cannot build signer config", "err", err)
		return nil, err
	}

	avsServiceBindings, err := NewAvsServiceBindings(baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr, baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr, baseConfig.EthRpcClient, baseConfig.EthRpcClientFallback, baseConfig.Logger)

	if err != nil {
		baseConfig.Logger.Error("Cannot create avs service bindings", "err", err)
		return nil, err
	}

	privateKeySigner, err := signer.NewPrivateKeySigner(ecdsaConfig.PrivateKey, baseConfig.ChainId)
	if err != nil {
		baseConfig.Logger.Error("Cannot create signer", "err", err)
		return nil, err
	}

	chainWriter := clients.AvsRegistryChainWriter

	return &AvsWriter{
		ChainWriter:         chainWriter,
		AvsContractBindings: avsServiceBindings,
		logger:              baseConfig.Logger,
		Signer:              privateKeySigner,
		Client:              baseConfig.EthRpcClient,
		ClientFallback:      baseConfig.EthRpcClientFallback,
	}, nil
}

func (w *AvsWriter) SendAggregatedResponse(batchIdentifierHash [32]byte, batchMerkleRoot [32]byte, senderAddress [20]byte, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*common.Hash, error) {
	txOpts := *w.Signer.GetTxOpts()
	txOpts.NoSend = true // simulate the transaction
	tx, err := w.RespondToTaskV2Retryable(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
	if err != nil {
		return nil, err
	}

	err = w.checkRespondToTaskFeeLimit(tx, txOpts, batchIdentifierHash, senderAddress)
	if err != nil {
		return nil, err
	}

	// Send the transaction
	txOpts.NoSend = false
	txOpts.GasLimit = tx.Gas() * 110 / 100 // Add 10% to the gas limit
	tx, err = w.RespondToTaskV2Retryable(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
	if err != nil {
		return nil, err
	}

	txHash := tx.Hash()

	return &txHash, nil
}

func (w *AvsWriter) checkRespondToTaskFeeLimit(tx *types.Transaction, txOpts bind.TransactOpts, batchIdentifierHash [32]byte, senderAddress [20]byte) error {
	aggregatorAddress := txOpts.From
	simulatedCost := new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasPrice())
	w.logger.Info("Simulated cost", "cost", simulatedCost)

	// Get RespondToTaskFeeLimit
	batchState, err := w.BatchesStateRetryable(batchIdentifierHash)
	if err != nil {
		// Fallback also failed
		// Proceed to check values against simulated costs
		w.logger.Error("Failed to get batch state", "error", err)
		w.logger.Info("Proceeding with simulated cost checks")
		return w.compareBalances(simulatedCost, aggregatorAddress, senderAddress)
	}
	// At this point, batchState was successfully retrieved
	// Proceed to check values against RespondToTaskFeeLimit
	respondToTaskFeeLimit := batchState.RespondToTaskFeeLimit
	w.logger.Info("Batch RespondToTaskFeeLimit", "RespondToTaskFeeLimit", respondToTaskFeeLimit)

	if respondToTaskFeeLimit.Cmp(simulatedCost) < 0 {
		return fmt.Errorf("cost of transaction is higher than Batch.RespondToTaskFeeLimit")
	}

	return w.compareBalances(respondToTaskFeeLimit, aggregatorAddress, senderAddress)
}

func (w *AvsWriter) compareBalances(amount *big.Int, aggregatorAddress common.Address, senderAddress [20]byte) error {
	if err := w.compareAggregatorBalance(amount, aggregatorAddress); err != nil {
		return err
	}
	if err := w.compareBatcherBalance(amount, senderAddress); err != nil {
		return err
	}
	return nil
}

func (w *AvsWriter) compareAggregatorBalance(amount *big.Int, aggregatorAddress common.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	aggregatorBalance, err := w.BalanceAtRetryable(ctx, aggregatorAddress, nil)
	if err != nil {
		// Ignore and continue.
		w.logger.Error("failed to get aggregator balance: %v", err)
		return nil
	}
	w.logger.Info("Aggregator balance", "balance", aggregatorBalance)
	if aggregatorBalance.Cmp(amount) < 0 {
		return fmt.Errorf("cost is higher than Aggregator balance")
	}
	return nil
}

func (w *AvsWriter) compareBatcherBalance(amount *big.Int, senderAddress [20]byte) error {
	// Get batcher balance
	batcherBalance, err := w.BatcherBalancesRetryable(senderAddress)
	if err != nil {
		// Ignore and continue.
		w.logger.Error("Failed to get batcherBalance", "error", err)
		return nil
	}
	w.logger.Info("Batcher balance", "balance", batcherBalance)
	if batcherBalance.Cmp(amount) < 0 {
		return fmt.Errorf("cost is higher than Batcher balance")
	}
	return nil
}

// |---RETRYABLE---|

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
				// Note return type will be nil
				if err.Error() == "not found" {
					err = connection.TransientError{Inner: err}
					return tx, err
				}
				if strings.Contains(err.Error(), "connect: connection refused") {
					err = connection.TransientError{Inner: err}
					return tx, err
				}
				if strings.Contains(err.Error(), "read: connection reset by peer") {
					return tx, connection.TransientError{Inner: err}
				}
				err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
			}
		}
		return tx, err
	}
	return connection.RetryWithData(respondToTaskV2_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func (w *AvsWriter) BatchesStateRetryable(arg0 [32]byte) (struct {
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
			// If error is not nil throw out result in state!
			if err != nil {
				// Note return type will be nil
				if err.Error() == "not found" {
					err = connection.TransientError{Inner: err}
					return state, err
				}
				if strings.Contains(err.Error(), "connect: connection refused") {
					err = connection.TransientError{Inner: err}
					return state, err
				}
				if strings.Contains(err.Error(), "read: connection reset by peer") {
					return state, connection.TransientError{Inner: err}
				}
				err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
			}
		}
		return state, err
	}
	return connection.RetryWithData(batchesState_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
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
				// Note return type will be nil
				if err.Error() == "not found" {
					err = connection.TransientError{Inner: err}
					return batcherBalance, err
				}
				if strings.Contains(err.Error(), "connect: connection refused") {
					err = connection.TransientError{Inner: err}
					return batcherBalance, err
				}
				if strings.Contains(err.Error(), "read: connection reset by peer") {
					return batcherBalance, connection.TransientError{Inner: err}
				}
				err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
			}
		}
		return batcherBalance, err
	}
	return connection.RetryWithData(batcherBalances_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}

func (w *AvsWriter) BalanceAtRetryable(ctx context.Context, aggregatorAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	var (
		aggregatorBalance *big.Int
		err               error
	)
	balanceAt_func := func() (*big.Int, error) {
		aggregatorBalance, err = w.Client.BalanceAt(ctx, aggregatorAddress, blockNumber)
		//aggregatorBalance, err := connection.RetryWithData(balanceAt_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries)
		if err != nil {
			aggregatorBalance, err = w.ClientFallback.BalanceAt(ctx, aggregatorAddress, blockNumber)
			//aggregatorBalance, err = connection.RetryWithData(balanceAt_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries)
			if err != nil {
				// Note return type will be nil
				if err.Error() == "not found" {
					err = connection.TransientError{Inner: err}
					return aggregatorBalance, err
				}
				if strings.Contains(err.Error(), "connect: connection refused") {
					err = connection.TransientError{Inner: err}
					return aggregatorBalance, err
				}
				if strings.Contains(err.Error(), "read: connection reset by peer") {
					return aggregatorBalance, connection.TransientError{Inner: err}
				}
				err = connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
			}
		}
		return aggregatorBalance, err
	}
	return connection.RetryWithData(balanceAt_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
}
