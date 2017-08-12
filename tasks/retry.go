package tasks

import (
	"fmt"
	"time"
)

var MaxRetryReachedError = fmt.Errorf("Maximum number of retry reached")
var TaskCancelledError = fmt.Errorf("Task is cancelled")

// RetryPolicy => given the number of iteration already tried,
// determine whether more retry is needed and for how long
type RetryPolicy func(int) (time.Duration, error)

// CancellationToken => send any value to signify cancellation
type CancellationToken chan error

// NewCancellationToken makes things look more neat, you can use the raw type if you want
func NewCancellationToken() CancellationToken {
	return make(CancellationToken)
}

// Cancel function is also here just to make things look more neat
func (c CancellationToken) Cancel() {
	c <- TaskCancelledError
}

// FixedDuration is the basic retry policy where you always
// wait for the same amount of time for a fix number of times
func FixedDuration(retryPeriod time.Duration, maxRetry int) RetryPolicy {
	return func(itr int) (time.Duration, error) {
		var err error
		if itr >= maxRetry {
			err = MaxRetryReachedError
		}
		return retryPeriod, err
	}
}

// RetryOperation help you make your code look cleaner but I am not here
// to protect you from infinite loops
func RetryOperation(operation func() bool, retryPolicy RetryPolicy, token CancellationToken) error {
	count := 0
	for {
		success := operation()
		if success {
			return nil
		}
		count++
		duration, policyViolation := retryPolicy(count)
		if policyViolation != nil {
			return policyViolation
		}
		if err := Wait(duration, token); err != nil {
			return err
		}
	}
}

// Wait for a specific amount of time, but cancel any time
func Wait(duration time.Duration, token CancellationToken) error {
	select {
	case <-time.After(duration):
		return nil
	case err := <-token:
		return err
	}
}
