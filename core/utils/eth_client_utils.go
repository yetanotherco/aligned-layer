package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	eigentypes "github.com/Layr-Labs/eigensdk-go/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	retry "github.com/yetanotherco/aligned_layer/core"
)

func WaitForTransactionReceiptRetryable(client eth.InstrumentedClient, ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	// For if no receipt and no error TransactionReceipt return "not found" as an error catch all ref: https://github.com/ethereum/go-ethereum/blob/master/ethclient/ethclient.go#L313
	receipt_func := func() (*types.Receipt, error) {
		tx, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			// Note return type will be nil
			if err.Error() == "not found" {
				return nil, retry.TransientError{Inner: err}
			}
			if strings.Contains(err.Error(), "connect: connection refused") {
				return nil, retry.TransientError{Inner: err}
			}
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				return nil, retry.TransientError{Inner: err}
			}
			return nil, retry.TransientError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return tx, err
	}
	return retry.RetryWithData(receipt_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval)
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
