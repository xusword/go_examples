package tasks

import (
	"fmt"
	"time"
)

// ErrorMaxRetryReached means maximum number of retry reached
func ErrorMaxRetryReached(numTry int) error {
	return fmt.Errorf("Maximum number of retry reached %d", numTry)
}

// ErrorRetryTimeout means maximum number of retry reached
func ErrorRetryTimeout(duration time.Duration) error {
	return fmt.Errorf("Retry timeout %f seconds reached", duration.Seconds())
}

// ErrorTaskCancelled means task is cancelled with a token
var ErrorTaskCancelled = fmt.Errorf("Task is cancelled")

// CancellationToken => Technically any value (include nil) can be sent
// to the channel but just don't. Use Cancel() instead because nil error
// sent to cancellation token will cause dependent to misbehave
type CancellationToken chan error

// NewCancellationToken makes things look more neat, you can use the raw type if you want
func NewCancellationToken() CancellationToken {
	return make(CancellationToken)
}

// Cancel function is also here just to make things look more neat
func (c CancellationToken) Cancel() {
	c <- ErrorTaskCancelled
}

type fixedDurationRetryPolicy struct {
	retryPeriod time.Duration
	maxRetry    int
}

func (r *fixedDurationRetryPolicy) GetDelay(itrNum int, lastErr error) (time.Duration, error) {
	var err error
	if itrNum >= r.maxRetry {
		err = ErrorMaxRetryReached(r.maxRetry)
	}
	return r.retryPeriod, err
}

// LinearRetry is the basic retry policy where you always
// wait for the same amount of time for a fix number of times
func LinearRetry(retryPeriod time.Duration, maxRetry int) RetryPolicy {
	return &fixedDurationRetryPolicy{retryPeriod: retryPeriod, maxRetry: maxRetry}
}

type timeoutRetryPolicy struct {
	retryPolicy RetryPolicy
	maxDuration time.Duration
	endTime     time.Time
}

func (r *timeoutRetryPolicy) GetDelay(itrNum int, lastErr error) (time.Duration, error) {
	if time.Now().After(r.endTime) {
		return time.Hour, ErrorRetryTimeout(r.maxDuration)
	}
	return r.retryPolicy.GetDelay(itrNum, lastErr)
}

// WithTimeout creates a policy that will expire maxDuration form now
// Note that the returned object is not designed to be reuseable as it records
// the time that the method is called
func WithTimeout(maxDuration time.Duration, retryPolicy RetryPolicy) RetryPolicy {
	endTime := time.Now().Add(maxDuration)
	return &timeoutRetryPolicy{
		retryPolicy: retryPolicy,
		maxDuration: maxDuration,
		endTime:     endTime,
	}
}

// RetryPolicy => given the number of iteration already tried,
// determine whether more retry is needed and for how long
type RetryPolicy interface {
	GetDelay(itrNum int, lastErr error) (time.Duration, error)
}

// RetryOperation help you make your code look cleaner but I am not here
// to protect you from infinite loops
func RetryOperation(operation func() error, retryPolicy RetryPolicy, token CancellationToken) error {
	count := 0
	for {
		err := operation()
		if err == nil {
			return nil
		}
		count++
		duration, policyViolation := retryPolicy.GetDelay(count, err)
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
