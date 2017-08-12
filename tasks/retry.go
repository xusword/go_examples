package tasks

import (
	"time"
)

const (
	// Success means operation was successful, according to you
	Success = iota
	// PolicyViolation means you told me not to try any more times
	PolicyViolation = iota
	// Cancelled means you cancelled
	Cancelled = iota
)

// RetryPolicy => given the number of iteration already tried,
// determine whether more retry is needed and for how long
type RetryPolicy func(int) (bool, time.Duration)

// CancellationToken => send any value to signify cancellation
type CancellationToken chan int

// NewCancellationToken makes things look more neat, you can use the raw type if you want
func NewCancellationToken() CancellationToken {
	return make(CancellationToken)
}

// Cancel function is also here just to make things look more neat
func (c CancellationToken) Cancel() {
	c <- 1
}

// FixedDuration is the basic retry policy where you always
// wait for the same amount of time for a fix number of times
func FixedDuration(retryPeriod time.Duration, maxRetry int) RetryPolicy {
	return func(itr int) (bool, time.Duration) {
		return itr < maxRetry, retryPeriod
	}
}

// RetryOperation help you make your code look cleaner but I am not here
// to protect you from infinite loops
func RetryOperation(operation func() bool, retryPolicy RetryPolicy, token CancellationToken) int {
	count := 0
	for {
		success := operation()
		if success {
			return Success
		}
		count++
		shouldWait, duration := retryPolicy(count)
		if shouldWait {
			if Wait(duration, token) == Cancelled {
				return Cancelled
			}
			// else keep trying
		} else {
			return PolicyViolation
		}
	}
}

// Wait for a specific amount of time, but cancel any time
func Wait(duration time.Duration, token CancellationToken) int {
	select {
	case <-time.After(duration):
		return Success
	case <-token:
		return Cancelled
	}
}
