package utils

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	eigentypes "github.com/Layr-Labs/eigensdk-go/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	retry "github.com/yetanotherco/aligned_layer/core"
)

func WaitForTransactionReceipt(client eth.InstrumentedClient, fallbackClient eth.InstrumentedClient, txHash gethcommon.Hash, config *retry.RetryConfig) func() (*types.Receipt, error) {
	receipt_func := func() (*types.Receipt, error) {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			receipt, err = fallbackClient.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				return nil, err
			}
			return receipt, nil
		}
		return receipt, nil
	}
	return receipt_func
}

// WaitForTransactionReceiptRetryable repeatedly attempts to fetch the transaction receipt for a given transaction hash.
// If the receipt is not found, the function will retry with exponential backoff until the specified `waitTimeout` duration is reached.
// If the receipt is still unavailable after `waitTimeout`, it will return an error.
//
// Note: The `time.Second * 2` is set as the max interval in the retry mechanism because we can't reliably measure the specific time the tx will be included in a block.
// Setting a higher value will imply doing less retries across the waitTimeout, and so we might lose the receipt
// All errors are considered Transient Errors
// - Retry times: 0.5s, 1s, 2s, 2s, 2s, ... until it reaches waitTimeout
func WaitForTransactionReceiptRetryable(client eth.InstrumentedClient, fallbackClient eth.InstrumentedClient, txHash gethcommon.Hash, config *retry.RetryConfig) (*types.Receipt, error) {
	return retry.RetryWithData(WaitForTransactionReceipt(client, fallbackClient, txHash, config), config)
}

func BytesToQuorumNumbers(quorumNumbersBytes []byte) eigentypes.QuorumNums {
	quorumNums := make(eigentypes.QuorumNums, len(quorumNumbersBytes))
	for i, quorumNumberByte := range quorumNumbersBytes {
		quorumNums[i] = eigentypes.QuorumNum(quorumNumberByte)
	}
	return quorumNums
}

func BytesToQuorumThresholdPercentages(quorumThresholdPercentagesBytes []byte) eigentypes.QuorumThresholdPercentages {
	quorumThresholdPercentages := make(eigentypes.QuorumThresholdPercentages, len(quorumThresholdPercentagesBytes))
	for i, quorumNumberByte := range quorumThresholdPercentagesBytes {
		quorumThresholdPercentages[i] = eigentypes.QuorumThresholdPercentage(quorumNumberByte)
	}
	return quorumThresholdPercentages
}

func WeiToEth(wei *big.Int) float64 {
	weiToEth := new(big.Float).SetFloat64(1e18)
	weiFloat := new(big.Float).SetInt(wei)

	result := new(big.Float).Quo(weiFloat, weiToEth)
	eth, _ := result.Float64()

	return eth
}

// Simple algorithm to calculate the gasPrice bump based on:
// the currentGasPrice, a base bump percentage, a retry percentage, and the retry count.
// Formula: currentGasPrice + (currentGasPrice * (baseBumpPercentage + retryCount * incrementalRetryPercentage) / 100)
func CalculateGasPriceBumpBasedOnRetry(currentGasPrice *big.Int, baseBumpPercentage uint, retryAttemptPercentage uint, retryCount int) *big.Int {
	// Incremental percentage increase for each retry attempt (i*retryAttemptPercentage)
	incrementalRetryPercentage := new(big.Int).Mul(big.NewInt(int64(retryAttemptPercentage)), big.NewInt(int64(retryCount)))

	// Total bump percentage: base bump + incremental retry percentage
	totalBumpPercentage := new(big.Int).Add(big.NewInt(int64(baseBumpPercentage)), incrementalRetryPercentage)

	// Calculate the bump amount: currentGasPrice * totalBumpPercentage / 100
	bumpAmount := new(big.Int).Mul(currentGasPrice, totalBumpPercentage)
	bumpAmount = new(big.Int).Div(bumpAmount, big.NewInt(100))

	// Final bumped gas price: currentGasPrice + bumpAmount
	bumpedGasPrice := new(big.Int).Add(currentGasPrice, bumpAmount)

	return bumpedGasPrice
}

/*
GetGasPriceRetryable
Get the gas price from the client with retry logic.
- All errors are considered Transient Errors
- Retry times: 1 sec, 2 sec, 4 sec
*/
func GetGasPriceRetryable(client eth.InstrumentedClient, fallbackClient eth.InstrumentedClient) (*big.Int, error) {
	respondToTaskV2_func := func() (*big.Int, error) {
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			gasPrice, err = fallbackClient.SuggestGasPrice(context.Background())
			if err != nil {
				return nil, err
			}
		}

		return gasPrice, nil
	}
	return retry.RetryWithData(respondToTaskV2_func, retry.EthCallRetryConfig())
}
