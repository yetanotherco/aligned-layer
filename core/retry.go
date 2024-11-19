package retry

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
)

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
	DefaultInitialInterval             = 1 * time.Second  // Initial delay for retry interval.
	DefaultMaxInterval                 = 60 * time.Second // Maximum interval an individual retry may have.
	DefaultMaxElapsedTime              = 0 * time.Second  // Maximum time all retries may take. `0` corresponds to no limit on the time of the retries.
	DefaultRandomizationFactor float64 = 0                // Randomization (Jitter) factor used to map retry interval to a range of values around the computed interval. In precise terms (random value in range [1 - randomizationfactor, 1 + randomizationfactor]). NOTE: This is set to 0 as we do not use jitter in Aligned.
	DefaultMultiplier          float64 = 2                // Multiplier factor computed exponential retry interval is scaled by.
	DefaultMaxNumRetries       uint64  = 3                // Total number of retries attempted.
	ChainInitialInterval               = 12 * time.Second // Initial delay for retry interval for contract calls. Corresponds to 1 ethereum block.
	ChainMaxInterval                   = 2 * time.Minute  // Maximum interval for an individual retry.
)

type RetryConfig struct {
	InitialInterval     time.Duration // Initial delay for retry interval.
	MaxInterval         time.Duration // Maximum interval an individual retry may have.
	MaxElapsedTime      time.Duration // Maximum time all retries may take. `0` corresponds to no limit on the time of the retries.
	RandomizationFactor float64
	Multiplier          float64
	MaxNumRetries       uint64
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		InitialInterval:     DefaultInitialInterval,
		MaxInterval:         DefaultMaxInterval,
		MaxElapsedTime:      DefaultMaxElapsedTime,
		RandomizationFactor: DefaultRandomizationFactor,
		Multiplier:          DefaultMultiplier,
		MaxNumRetries:       DefaultMaxNumRetries,
	}
}

func ChainRetryConfig() *RetryConfig {
	return &RetryConfig{
		InitialInterval:     ChainInitialInterval,
		MaxInterval:         ChainMaxInterval,
		MaxElapsedTime:      DefaultMaxElapsedTime,
		RandomizationFactor: DefaultRandomizationFactor,
		Multiplier:          DefaultMultiplier,
		MaxNumRetries:       DefaultMaxNumRetries,
	}
}

/*
Retry and RetryWithData are custom retry functions used in Aligned's aggregator and operator to facilitate consistent retry logic across the system.
They are interfaces for around Cenk Alti (https://github.com/cenkalti) backoff library (https://github.com/cenkalti/backoff). We would like to thank him for his great work.

The `Retry` and `RetryWithData` retry a supplied function at maximum `NumRetries` number of times. Upon execution, if the called function returns an error the retry library either re-executes the function (Transient Error) or exits and returns the error to the calling context (Permanent error) .
If the call is successful and no error is returned the library returns the result. `Permanent` errors are explicitly typed while `Transient` errors are implied by go's builtin error type.
For completeness:

Transient: The error is recoverable and the function is retried after failing. `Transient` errors do not have a defined error type and are implicitly defined by go's builtin `error` type.

Permanent: The error is not recoverable by retrying and the error to the calling context. Permanent errors are explicitly typed and defined by wrapping the err within with `PermanentError` type.

Usage of `RetryWithData` is shown in the following example:
```
	sendUserMsg_func := func() (*types.Transaction, error) {
		res, err := sendUserMessage(opts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			// Detect Permanent error by checking contents of returned error message.
			if strings.Contains(err.Error(), "client error: User not registered:") {
				err = retry.PermanentError{Inner: err}
			}
		}
		return res, err
	}
	err := retry.Retry(sendUserMsg_func, retry.MinDelay, retry.RetryFactor, retry.NumRetries, retry.MaxInterval, retry.MaxElapsedTime)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
```

# Retry Intervals:
The backoff period for each retry attempt increases using a randomization function that grows exponentially.

retryinterval =
    currentRetryinterval * (random value in range [1 - randomizationfactor, 1 + randomizationfactor]) * retryFactor

This library defaults to the use of the following parameters:

randomizationFactor = 0.5 // Randomization factor that maps the interval increase within a range around the computed retry interval.
initialRetryInterval = 1 sec // Initial value used in the retry interval
Multiplier = 2 // Multiplier used to scale the values.

# Default intervals for Retries (sec)
request     retry_interval (1 sec)     randomized_interval (0.5)		randomized_interval_scaled (2)
   1             1                			[0.5, 1.5]			 				[1, 3]
   2             2               			[1, 3]								[2, 6]
   3             4              			[2, 6]								[4, 12]

# Default intervals for Contract Calls (sec)
request     retry_interval (12 sec)    randomized_interval (0.5)		randomized_interval_scaled (2)
   1             12                			[6, 18]								[12, 36]
   2             24               			[12, 36]							[24, 72]
   3             48              			[24, 72]`							[48, 144]`

Reference: https://github.com/cenkalti/backoff/blob/v4/exponential.go#L9
*/

// Same as Retry only that the functionToRetry can return a value upon correct execution
func RetryWithData[T any](functionToRetry func() (T, error), config *RetryConfig) (T, error) {
	//func RetryWithData[T any](functionToRetry func() (T, error), minDelay time.Duration, factor float64, maxTries uint64, maxInterval time.Duration, maxElapsedTime time.Duration) (T, error) {
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

	initialRetryOption := backoff.WithInitialInterval(config.InitialInterval)
	multiplierOption := backoff.WithMultiplier(config.Multiplier)
	maxIntervalOption := backoff.WithMaxInterval(config.MaxInterval)
	maxElapsedTimeOption := backoff.WithMaxElapsedTime(config.MaxElapsedTime)
	expBackoff := backoff.NewExponentialBackOff(randomOption, multiplierOption, initialRetryOption, maxIntervalOption, maxElapsedTimeOption)
	var maxRetriesBackoff backoff.BackOff

	if config.MaxNumRetries > 0 {
		maxRetriesBackoff = backoff.WithMaxRetries(expBackoff, config.MaxNumRetries)
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
func Retry(functionToRetry func() error, config *RetryConfig) error {
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

	initialRetryOption := backoff.WithInitialInterval(config.InitialInterval)
	multiplierOption := backoff.WithMultiplier(config.Multiplier)
	maxIntervalOption := backoff.WithMaxInterval(config.MaxInterval)
	maxElapsedTimeOption := backoff.WithMaxElapsedTime(config.MaxElapsedTime)
	expBackoff := backoff.NewExponentialBackOff(randomOption, multiplierOption, initialRetryOption, maxIntervalOption, maxElapsedTimeOption)
	var maxRetriesBackoff backoff.BackOff

	if config.MaxNumRetries > 0 {
		maxRetriesBackoff = backoff.WithMaxRetries(expBackoff, config.MaxNumRetries)
	} else {
		maxRetriesBackoff = expBackoff
	}

	return backoff.Retry(f, maxRetriesBackoff)
}
