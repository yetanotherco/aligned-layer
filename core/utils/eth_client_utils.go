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

func WaitForTransactionReceiptRetryable(client eth.InstrumentedClient, ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	receipt_func := func() (*types.Receipt, error) {
		return client.TransactionReceipt(ctx, txHash)
	}
	return retry.RetryWithData(receipt_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
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

// Simple algorithm to calculate the gasPrice bump based on:
// the currentGasPrice, a base bump percentage, a retry percentage, and the retry count.
// Formula: currentGasPrice + (currentGasPrice * (baseBumpPercentage + retryCount * incrementalRetryPercentage) / 100)
func CalculateGasPriceBumpBasedOnRetry(currentGasPrice *big.Int, baseBumpPercentage int, retryAttemptPercentage int, retryCount int) *big.Int {
	// Incremental percentage increase for each retry attempt (i*5%)
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
