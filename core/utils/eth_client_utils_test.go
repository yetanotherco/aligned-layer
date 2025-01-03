package utils_test

import (
	"math/big"
	"testing"

	"github.com/yetanotherco/aligned_layer/core/utils"
)

func TestCalculateGasPriceBumpBasedOnRetry(t *testing.T) {
	baseBumpPercentage := uint(20)
	incrementalRetryPercentage := uint(5)

	gasPrices := [5]*big.Int{
		big.NewInt(3000000000),
		big.NewInt(3000000000),
		big.NewInt(4000000000),
		big.NewInt(4000000000),
		big.NewInt(5000000000)}

	expectedBumpedGasPrices := [5]*big.Int{
		big.NewInt(3600000000),
		big.NewInt(3750000000),
		big.NewInt(5200000000),
		big.NewInt(5400000000),
		big.NewInt(7000000000)}

	for i := 0; i < len(gasPrices); i++ {
		currentGasPrice := gasPrices[i]
		bumpedGasPrice := utils.CalculateGasPriceBumpBasedOnRetry(currentGasPrice, baseBumpPercentage, incrementalRetryPercentage, 100, i)
		expectedGasPrice := expectedBumpedGasPrices[i]

		if bumpedGasPrice.Cmp(expectedGasPrice) != 0 {
			t.Errorf("Bumped gas price does not match expected gas price, expected value %v, got: %v", expectedGasPrice, bumpedGasPrice)
		}
	}
}

func TestCalculateGasPriceBumpBasedOnRetryPercentageLimit(t *testing.T) {
	baseBumpPercentage := uint(20)
	incrementalRetryPercentage := uint(5)

	gasPrices := [5]*big.Int{
		big.NewInt(3000000000),
		big.NewInt(3000000000),
		big.NewInt(4000000000),
		big.NewInt(4000000000),
		big.NewInt(5000000000)}

	expectedBumpedGasPrices := [5]*big.Int{
		big.NewInt(3600000000),
		big.NewInt(3750000000),
		big.NewInt(5200000000),
		big.NewInt(5200000000),
		big.NewInt(6500000000)}

	for i := 0; i < len(gasPrices); i++ {
		currentGasPrice := gasPrices[i]
		bumpedGasPrice := utils.CalculateGasPriceBumpBasedOnRetry(currentGasPrice, baseBumpPercentage, incrementalRetryPercentage, 30, i)
		expectedGasPrice := expectedBumpedGasPrices[i]

		if bumpedGasPrice.Cmp(expectedGasPrice) != 0 {
			t.Errorf("Bumped gas price does not match expected gas price, expected value %v, got: %v", expectedGasPrice, bumpedGasPrice)
		}
	}
}
