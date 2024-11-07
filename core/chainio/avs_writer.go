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
	"github.com/yetanotherco/aligned_layer/core/utils"
)

const (
	// How much to bump every retry (constant)
	GasBaseBumpPercentage int = 20
	// An extra percentage to bump every retry i*5 (linear)
	GasBumpIncrementalPercentage int = 5
	// Wait as much as 3 blocks time for the receipt
	BlocksToWaitBeforeBump time.Duration = time.Second * 36
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

// Sends AggregatedResponse and waits for the receipt for three blocks, if not received
// it will try again bumping the last tx gas price based on `CalculateGasPriceBump`
// This process happens indefinitely until the transaction is included.
func (w *AvsWriter) SendAggregatedResponse(batchIdentifierHash [32]byte, batchMerkleRoot [32]byte, senderAddress [20]byte, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature, onRetry func()) (*types.Receipt, error) {
	txOpts := *w.Signer.GetTxOpts()
	txOpts.NoSend = true // simulate the transaction

	respondToTaskV2SimulationFunc := func() (*types.Transaction, error) {
		tx, err := w.AvsContractBindings.ServiceManager.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
			if err != nil {
				// check if reverted only else transient error
				err = connection.PermanentError{Inner: fmt.Errorf("transaction reverted")}
			}
		}
		return tx, err
	}
	tx, err := connection.RetryWithData(respondToTaskV2SimulationFunc, connection.MinDelay, connection.RetryFactor, 0, connection.MaxInterval)
	if err != nil {
		return nil, err
	}

	err = w.checkRespondToTaskFeeLimit(tx, txOpts, batchIdentifierHash, senderAddress)
	if err != nil {
		return nil, connection.PermanentError{Inner: err}
	}

	// Set the nonce, as we might have to replace the transaction with a higher gas price
	txNonce := big.NewInt(int64(tx.Nonce()))
	txOpts.Nonce = txNonce
	txOpts.NoSend = false
	i := 0

	shouldBump := false

	respondToTaskV2Func := func() (*types.Receipt, error) {
		// We bump only when the timeout for waiting for the receipt has passed
		// not when an rpc failed
		if shouldBump {
			gasPrice, err := w.ClientFallback.SuggestGasPrice(context.Background())
			if err != nil {
				shouldBump = false
				return nil, connection.TransientError{Inner: err}
			}
			bumpedGasPrice := utils.CalculateGasPriceBumpBasedOnRetry(gasPrice, GasBaseBumpPercentage, GasBumpIncrementalPercentage, i)

			if gasPrice.Cmp(txOpts.GasPrice) > 0 {
				txOpts.GasPrice = bumpedGasPrice
			} else {
				txOpts.GasPrice = new(big.Int).Mul(txOpts.GasPrice, big.NewInt(1))
			}
		}

		err = w.checkRespondToTaskFeeLimit(tx, txOpts, batchIdentifierHash, senderAddress)
		if err != nil {
			shouldBump = false
			// We bump the fee so much that the transaction cost is more expensive than the batcher fee limit
			return nil, connection.PermanentError{Inner: err}
		}

		tx, err = w.AvsContractBindings.ServiceManager.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
			if err != nil {
				// check if reverted only to be permanent
				err = connection.TransientError{Inner: err}
				shouldBump = false
			}
		}

		receipt, err := utils.WaitForTransactionReceiptRetryable(w.Client, context.Background(), tx.Hash())
		if receipt != nil {
			return receipt, nil
		}

		// if we are here, it means we have reached the receipt waiting timeout
		// so in the next iteration we try again by bumping the fee to make sure its included
		if i > 0 {
			onRetry()
		}
		i++

		shouldBump = true
		if err != nil {
			return nil, connection.TransientError{Inner: err}
		}
		return nil, connection.TransientError{Inner: fmt.Errorf("transaction failed")}
	}

	return connection.RetryWithData(respondToTaskV2Func, 1000, 2, 0, 60)
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
		if err != nil {
			aggregatorBalance, err = w.ClientFallback.BalanceAt(ctx, aggregatorAddress, blockNumber)
			if err != nil {
				// Note return type will be nil
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
