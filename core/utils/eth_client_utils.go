package utils

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	eigentypes "github.com/Layr-Labs/eigensdk-go/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	connection "github.com/yetanotherco/aligned_layer/core"
)

/*
Errors:
- "not found": (Transient) Call successfully returns but the tx receipt was not found.
- "connect: connection refused": (Transient) Could not connect.
*/
func WaitForTransactionReceiptRetryable(client eth.InstrumentedClient, ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	// For if no receipt and no error TransactionReceipt return "not found" as an error catch all ref: https://github.com/ethereum/go-ethereum/blob/master/ethclient/ethclient.go#L313
	receipt_func := func() (*types.Receipt, error) {
		tx, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			// Note return type will be nil
			if err.Error() == "not found" {
				return nil, connection.TransientError{Inner: err}
			}
			if strings.Contains(err.Error(), "connect: connection refused") {
				return nil, connection.TransientError{Inner: err}
			}
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				return nil, connection.TransientError{Inner: err}
			}
			return nil, connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return tx, err
	}
	return connection.RetryWithData(receipt_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries, connection.MaxInterval)
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
