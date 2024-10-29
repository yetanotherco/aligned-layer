package utils

import (
	"context"
	"fmt"
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
			return nil, connection.PermanentError{Inner: fmt.Errorf("Permanent error: Unexpected Error while retrying: %s\n", err)}
		}
		return tx, err
	}
	return connection.RetryWithData(receipt_func, connection.MinDelay, connection.RetryFactor, connection.NumRetries)
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
