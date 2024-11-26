package retry_test

import (
	"fmt"
	"testing"

	retry "github.com/yetanotherco/aligned_layer/core"
)

func DummyFunction(x uint64) (uint64, error) {
	if x == 42 {
		return 0, retry.PermanentError{Inner: fmt.Errorf("Permanent error!")}
	} else if x < 42 {
		return 0, fmt.Errorf("Transient error!")
	}
	return x, nil
}

func TestRetryWithData(t *testing.T) {
	function := func() (*uint64, error) {
		x, err := DummyFunction(43)
		return &x, err
	}
	config := &retry.RetryParams{
		InitialInterval:     1000,
		MaxInterval:         2,
		MaxElapsedTime:      3,
		RandomizationFactor: 0,
		Multiplier:          retry.NetworkMultiplier,
		NumRetries:          retry.NetworkNumRetries,
	}
	_, err := retry.RetryWithData(function, config)
	if err != nil {
		t.Errorf("Retry error!: %s", err)
	}
}

func TestRetry(t *testing.T) {
	function := func() error {
		_, err := DummyFunction(43)
		return err
	}
	config := &retry.RetryParams{
		InitialInterval:     1000,
		MaxInterval:         2,
		MaxElapsedTime:      3,
		RandomizationFactor: 0,
		Multiplier:          retry.NetworkMultiplier,
		NumRetries:          retry.NetworkNumRetries,
	}
	err := retry.Retry(function, config)
	if err != nil {
		t.Errorf("Retry error!: %s", err)
	}
}
