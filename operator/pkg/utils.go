package operator

import (
	"math/big"
	"net/url"

	"github.com/yetanotherco/aligned_layer/common"
)

func IsVerifierDisabled(disabledVerifiersBitmap *big.Int, verifierId common.ProvingSystemId) bool {
	verifierIdInt := uint8(verifierId)
	// The cast to uint64 is necessary because we need to use the bitwise AND operator.
	// This will truncate the bitmap to 64 bits, but we are not expecting to have more than 63 verifiers.
	// If we set a number that doesn't fit in 64 bits, the bitmap will be truncated and no verifier will be disabled.
	bit := disabledVerifiersBitmap.Uint64() & (1 << verifierIdInt)
	return bit != 0
}

func BaseUrlOnly(input string) (string, error) {
	// https://gobyexample.com/url-parsing
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}
