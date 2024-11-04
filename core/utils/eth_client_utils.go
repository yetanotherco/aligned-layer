package utils

import (
	"context"
	"math/big"
	"time"

	"fmt"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	eigentypes "github.com/Layr-Labs/eigensdk-go/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const maxRetries = 25
const sleepTime = 1000 * time.Millisecond

func WaitForTransactionReceipt(client eth.InstrumentedClient, ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	for i := 0; i < maxRetries; i++ {
		receipt, _ := client.TransactionReceipt(ctx, txHash)
		if receipt != nil {
			return receipt, nil
		}
		// if context has timed out, return
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		time.Sleep(sleepTime)
	}
	return nil, fmt.Errorf("transaction receipt not found for txHash: %s", txHash.String())
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
