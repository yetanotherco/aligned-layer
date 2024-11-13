package utils

import (
	"context"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	eigentypes "github.com/Layr-Labs/eigensdk-go/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	retry "github.com/yetanotherco/aligned_layer/core"
)

/*
WaitForTransactionReceipt
Wait for the transaction receipt on Ethereum using exponential backoff
- All errors are considered Transient Errors
- Retry times (3 retries): 12 sec (1 Blocks), 24 sec (2 Blocks), 48 sec (4 Blocks)
*/
func WaitForTransactionReceipt(client eth.InstrumentedClient, ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	receipt_func := func() (*types.Receipt, error) {
		return client.TransactionReceipt(ctx, txHash)
	}
	return retry.RetryWithData(receipt_func, retry.MinDelayChain, retry.RetryFactor, retry.NumRetries, retry.MaxIntervalChain, retry.MaxElapsedTime)
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
