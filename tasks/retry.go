package tasks

import (
	"time"
)

const (
	// OperationSuccess means operation was successful, according to you
	OperationSuccess = iota
	// PolicyViolation means you told me not to try any more times
	PolicyViolation = iota
	// UserCancelled means you cancelled
	UserCancelled = iota
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
			return OperationSuccess
		}
		count++
		shouldWait, waitDuration := retryPolicy(count)
		if shouldWait {
			select {
			case <-time.After(waitDuration):
			case <-token:
				return UserCancelled
			}
		} else {
			return PolicyViolation
		}
	}
}
