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
	"github.com/yetanotherco/aligned_layer/metrics"
)

type AvsWriter struct {
	*avsregistry.ChainWriter
	AvsContractBindings *AvsServiceBindings
	logger              logging.Logger
	Signer              signer.Signer
	Client              eth.InstrumentedClient
	ClientFallback      eth.InstrumentedClient
	metrics             *metrics.Metrics
}

func NewAvsWriterFromConfig(baseConfig *config.BaseConfig, ecdsaConfig *config.EcdsaConfig, metrics *metrics.Metrics) (*AvsWriter, error) {

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
		metrics:             metrics,
	}, nil
}

// Sends AggregatedResponse and waits for the receipt for three blocks, if not received
// it will try again bumping the last tx gas price based on `CalculateGasPriceBump`
// This process happens indefinitely until the transaction is included.
func (w *AvsWriter) SendAggregatedResponse(batchIdentifierHash [32]byte, batchMerkleRoot [32]byte, senderAddress [20]byte, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature, gasBumpPercentage uint, gasBumpIncrementalPercentage uint, timeToWaitBeforeBump time.Duration, onGasPriceBumped func(*big.Int)) (*types.Receipt, error) {
	txOpts := *w.Signer.GetTxOpts()
	txOpts.NoSend = true // simulate the transaction
	sim_tx, err := w.RespondToTaskV2Retryable(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
	if err != nil {
		return nil, err
	}

	// Set the nonce, as we might have to replace the transaction with a higher gas price
	txNonce := big.NewInt(int64(sim_tx.Nonce()))
	txOpts.Nonce = txNonce
	txOpts.GasPrice = sim_tx.GasPrice()
	txOpts.NoSend = false
	i := 0

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
			w.logger.Infof("Checking batch state again before sending bumped transaction")
			batchState, _ := w.BatchesStateRetryable(&bind.CallOpts{}, batchIdentifierHash)
			if batchState.Responded {
				return nil, nil
			}
			onGasPriceBumped(txOpts.GasPrice)
		}

		// We compare both Aggregator funds and Batcher balance in Aligned against respondToTaskFeeLimit
		// Both are required to have some balance, more details inside the function
		err = w.checkAggAndBatcherHaveEnoughBalance(sim_tx, txOpts, batchIdentifierHash, senderAddress)
		if err != nil {
			w.logger.Errorf("Permanent error when checking respond to task fee limit, err %v", err)
			return nil, retry.PermanentError{Inner: err}
		}

		w.logger.Infof("Sending RespondToTask transaction with a gas price of %v", txOpts.GasPrice)
		real_tx, err := w.RespondToTaskV2Retryable(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			w.logger.Errorf("Respond to task transaction err, %v", err)
			return nil, err
		}
		w.logger.Infof("Transaction sent, waiting for receipt")
		receipt, err := utils.WaitForTransactionReceiptRetryable(w.Client, w.ClientFallback, real_tx.Hash(), timeToWaitBeforeBump)
		if receipt != nil {
			w.checkIfAggregatorHadToPaidForBatcher(real_tx, batchIdentifierHash)
			return receipt, nil
		}

		// if we are here, it means we have reached the receipt waiting timeout
		// we increment the i here to add an incremental percentage to increase the odds of being included in the next blocks
		i++

		w.logger.Infof("RespondToTask receipt waiting timeout has passed, will try again...", "merkle_root", batchMerkleRoot)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("transaction failed")
	}

	return retry.RetryWithData(respondToTaskV2Func, retry.MinDelay, retry.RetryFactor, 0, retry.MaxInterval, 0)
}

// Calculates the transaction cost from the receipt and compares it with the batcher respondToTaskFeeLimit
// if the tx cost was higher, then it means the aggregator has paid the difference for the batcher (txCost - respondToTaskFeeLimit) and so metrics are updated accordingly.
// otherwise nothing is done.
func (w *AvsWriter) checkIfAggregatorHadToPaidForBatcher(tx *types.Transaction, batchIdentifierHash [32]byte) {
	batchState, err := w.BatchesStateRetryable(&bind.CallOpts{}, batchIdentifierHash)
	if err != nil {
		return
	}
	respondToTaskFeeLimit := batchState.RespondToTaskFeeLimit

	// NOTE we are not using tx.Cost() because tx.Cost() includes tx.Value()
	txCost := new(big.Int).Mul(big.NewInt(int64(tx.Gas())), tx.GasPrice())

	if respondToTaskFeeLimit.Cmp(txCost) < 0 {
		aggregatorDifferencePaid := new(big.Int).Sub(txCost, respondToTaskFeeLimit)
		aggregatorDifferencePaidInEth := utils.WeiToEth(aggregatorDifferencePaid)
		w.metrics.AddAggregatorGasPaidForBatcher(aggregatorDifferencePaidInEth)
		w.metrics.IncAggregatorPaidForBatcher()
		w.logger.Warnf("cost of transaction was higher than Batch.RespondToTaskFeeLimit, aggregator has paid the for the difference, aprox: %vethers", aggregatorDifferencePaidInEth)
	}
}

func (w *AvsWriter) checkAggAndBatcherHaveEnoughBalance(tx *types.Transaction, txOpts bind.TransactOpts, batchIdentifierHash [32]byte, senderAddress [20]byte) error {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("Recovered from panic", "error", r)
		}
		// return fmt.Errorf("Recovered from panic") // TODO can't do this
	}()

	w.logger.Info("Checking if aggregator and batcher have enough balance for the transaction")
	aggregatorAddress := txOpts.From
	txGasAsBigInt := new(big.Int).SetUint64(tx.Gas())
	txGasPrice := txOpts.GasPrice
	w.logger.Info("Transaction Gas Cost", "cost", txGasAsBigInt)
	w.logger.Info("Transaction Gas Price", "cost", txGasPrice)

	txCost := new(big.Int).Mul(txGasAsBigInt, txGasPrice)

	// txCost := new(big.Int).Mul(new(big.Int)tx.Gas()), txOpts.GasPrice)
	w.logger.Info("Transaction cost", "cost", txCost)

	batchState, err := w.BatchesStateRetryable(&bind.CallOpts{}, batchIdentifierHash)
	if err != nil {
		w.logger.Error("Failed to get batch state", "error", err)
		w.logger.Info("Proceeding to check balances against transaction cost")
		return w.compareBalances(txCost, aggregatorAddress, senderAddress)
	}
	respondToTaskFeeLimit := batchState.RespondToTaskFeeLimit
	w.logger.Info("Checking balance against Batch RespondToTaskFeeLimit", "RespondToTaskFeeLimit", respondToTaskFeeLimit)
	// Note: we compare both Aggregator funds and Batcher balance in Aligned against respondToTaskFeeLimit
	// Batcher will pay up to respondToTaskFeeLimit, for this he needs that amount of funds in Aligned
	// Aggregator will pay any extra cost, for this he needs at least respondToTaskFeeLimit in his balance
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
