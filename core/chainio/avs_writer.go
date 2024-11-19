package chainio

import (
	"context"
	"fmt"
	"math/big"
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
	retry "github.com/yetanotherco/aligned_layer/core"
	"github.com/yetanotherco/aligned_layer/core/config"
	"github.com/yetanotherco/aligned_layer/core/utils"
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
func (w *AvsWriter) SendAggregatedResponse(batchIdentifierHash [32]byte, batchMerkleRoot [32]byte, senderAddress [20]byte, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature, gasBumpPercentage uint, gasBumpIncrementalPercentage uint, timeToWaitBeforeBump time.Duration, onGasPriceBumped func(*big.Int)) (*types.Receipt, error) {
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

	// Set the nonce, as we might have to replace the transaction with a higher gas price
	txNonce := big.NewInt(int64(tx.Nonce()))
	txOpts.Nonce = txNonce
	txOpts.GasPrice = tx.GasPrice()
	txOpts.NoSend = false
	i := 0

	// Set Retry config for RespondToTaskV2
	respondToTaskV2Config := retry.DefaultRetryConfig()
	respondToTaskV2Config.MaxNumRetries = 0
	respondToTaskV2Config.MaxElapsedTime = 0

	// Set Retry config for WaitForTxRetryable
	waitForTxConfig := retry.DefaultRetryConfig()
	waitForTxConfig.MaxInterval = 2 * time.Second
	waitForTxConfig.MaxNumRetries = 0
	waitForTxConfig.MaxElapsedTime = timeToWaitBeforeBump

	respondToTaskV2Func := func() (*types.Receipt, error) {
		gasPrice, err := utils.GetGasPriceRetryable(w.Client, w.ClientFallback)
		if err != nil {
			return nil, err
		}

		bumpedGasPrice := utils.CalculateGasPriceBumpBasedOnRetry(gasPrice, gasBumpPercentage, gasBumpIncrementalPercentage, i)
		// new bumped gas price must be higher than the last one (this should hardly ever happen though)
		if bumpedGasPrice.Cmp(txOpts.GasPrice) > 0 {
			txOpts.GasPrice = bumpedGasPrice
		} else {
			// bump the last tx gas price a little by `gasBumpIncrementalPercentage` to replace it.
			txOpts.GasPrice = utils.CalculateGasPriceBumpBasedOnRetry(txOpts.GasPrice, gasBumpIncrementalPercentage, 0, 0)
		}

		if i > 0 {
			onGasPriceBumped(txOpts.GasPrice)
		}

		err = w.checkRespondToTaskFeeLimit(tx, txOpts, batchIdentifierHash, senderAddress)
		if err != nil {
			return nil, retry.PermanentError{Inner: err}
		}

		w.logger.Infof("Sending RespondToTask transaction with a gas price of %v", txOpts.GasPrice)

		tx, err = w.RespondToTaskV2Retryable(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			return nil, err
		}

		receipt, err := utils.WaitForTransactionReceiptRetryable(w.Client, w.ClientFallback, tx.Hash(), waitForTxConfig)
		if receipt != nil {
			return receipt, nil
		}

		// if we are here, it means we have reached the receipt waiting timeout
		// we increment the i here to add an incremental percentage to increase the odds of being included in the next blocks
		i++

		w.logger.Infof("RespondToTask receipt waiting timeout has passed, will try again...")
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("transaction failed")
	}

	return retry.RetryWithData(respondToTaskV2Func, respondToTaskV2Config)
}

func (w *AvsWriter) checkRespondToTaskFeeLimit(tx *types.Transaction, txOpts bind.TransactOpts, batchIdentifierHash [32]byte, senderAddress [20]byte) error {
	aggregatorAddress := txOpts.From
	simulatedCost := new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasPrice())
	w.logger.Info("Simulated cost", "cost", simulatedCost)

	// Get RespondToTaskFeeLimit
	batchState, err := w.BatchesStateRetryable(&bind.CallOpts{}, batchIdentifierHash)
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
	batcherBalance, err := w.BatcherBalancesRetryable(&bind.CallOpts{}, senderAddress)
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
