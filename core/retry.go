package retry

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
)

/*
This retry library was inspired by and uses Cenk Alti (https://github.com/cenkalti) backoff library (https://github.com/cenkalti/backoff).
We would like to thank him for his great work.
*/

/*
Note we use a custom Permanent error type for asserting Permanent Erros within the retry library.
We do not implement an explicit Transient error type and operate under the assumption that all errors that are not Permanent are Transient.
*/
type PermanentError struct {
	Inner error
}

func (e PermanentError) Error() string { return e.Inner.Error() }
func (e PermanentError) Unwrap() error {
	return e.Inner
}
func (e PermanentError) Is(err error) bool {
	_, ok := err.(PermanentError)
	return ok
}

const (
	MinDelay               = 1 * time.Second
	MaxInterval            = 60 * time.Second
	MaxElapsedTime         = 0 * time.Second
	RetryFactor    float64 = 2
	NumRetries     uint64  = 3
)

// Same as Retry only that the functionToRetry can return a value upon correct execution
func RetryWithData[T any](functionToRetry func() (T, error), minDelay time.Duration, factor float64, maxTries uint64, maxInterval time.Duration, maxElapsedTime time.Duration) (T, error) {
	f := func() (T, error) {
		var (
			val T
			err error
		)
		func() {
			defer func() {
				if r := recover(); r != nil {
					if panic_err, ok := r.(error); ok {
						err = panic_err
					} else {
						err = fmt.Errorf("panicked: %v", panic_err)
					}
				}
			}()
			val, err = functionToRetry()
			// Convert the returned `PermanentError` (our implementation) to `backoff.PermanentError`.
			//This exits the retry loop in the `backoff` library.
			if perm, ok := err.(PermanentError); err != nil && ok {
				err = backoff.Permanent(perm.Inner)
			}
		}()
		return val, err
	}

	randomOption := backoff.WithRandomizationFactor(0)

	initialRetryOption := backoff.WithInitialInterval(minDelay)
	multiplierOption := backoff.WithMultiplier(factor)
	maxIntervalOption := backoff.WithMaxInterval(maxInterval)
	maxElapsedTimeOption := backoff.WithMaxElapsedTime(maxElapsedTime)
	expBackoff := backoff.NewExponentialBackOff(randomOption, multiplierOption, initialRetryOption, maxIntervalOption, maxElapsedTimeOption)
	var maxRetriesBackoff backoff.BackOff

	if maxTries > 0 {
		maxRetriesBackoff = backoff.WithMaxRetries(expBackoff, maxTries)
	} else {
		maxRetriesBackoff = expBackoff
	}

	return backoff.RetryWithData(f, maxRetriesBackoff)
}

// Retries a given function in an exponential backoff manner.
// It will retry calling the function while it returns an error, until the max retries.
// If maxTries == 0 then the retry function will run indefinitely until success
// from the configuration are reached, or until a `PermanentError` is returned.
// The function to be retried should return `PermanentError` when the condition for stop retrying
// is met.
func Retry(functionToRetry func() error, minDelay time.Duration, factor float64, maxTries uint64, maxInterval time.Duration, maxElapsedTime time.Duration) error {
	f := func() error {
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					if panic_err, ok := r.(error); ok {
						err = panic_err
					} else {
						err = fmt.Errorf("panicked: %v", panic_err)
					}
				}
			}()
			err = functionToRetry()
			// Convert the returned `PermanentError` (our implementation) to a `backoff.PermanentError`.
			//This exits the retry loop in the `backoff` library.
			if perm, ok := err.(PermanentError); err != nil && ok {
				err = backoff.Permanent(perm.Inner)
			}
		}()
		return err
	}

	randomOption := backoff.WithRandomizationFactor(0)

	initialRetryOption := backoff.WithInitialInterval(minDelay)
	multiplierOption := backoff.WithMultiplier(factor)
	maxIntervalOption := backoff.WithMaxInterval(maxInterval)
	maxElapsedTimeOption := backoff.WithMaxElapsedTime(maxElapsedTime)
	expBackoff := backoff.NewExponentialBackOff(randomOption, multiplierOption, initialRetryOption, maxIntervalOption, maxElapsedTimeOption)
	var maxRetriesBackoff backoff.BackOff

	if maxTries > 0 {
		maxRetriesBackoff = backoff.WithMaxRetries(expBackoff, maxTries)
	} else {
		maxRetriesBackoff = expBackoff
	}

	return backoff.Retry(f, maxRetriesBackoff)
}
