package tasks

import (
	"fmt"
	"testing"
	"time"
)

// Hey, I am not here to test every cases
func TestRetry(t *testing.T) {
	maxRetry := 3
	retryPolicy := LinearRetry(time.Duration(1)*time.Second, 3)
	errorMaxRetryReached := ErrorMaxRetryReached(maxRetry)
	testCase := []struct {
		durationCancel        time.Duration
		useToken              bool
		ticksToSuccess        int
		expectedTicks         int
		expectedResult        error
		expectedElapsedSecond int
	}{
		// no token, run all 3 iterations
		{
			durationCancel:        time.Duration(0),
			useToken:              false,
			ticksToSuccess:        -1,
			expectedTicks:         3,
			expectedResult:        errorMaxRetryReached,
			expectedElapsedSecond: 2,
		},
		// run all 3 iterations, token never cancelled
		{
			durationCancel:        time.Duration(10) * time.Second,
			useToken:              true,
			ticksToSuccess:        -1,
			expectedTicks:         3,
			expectedResult:        errorMaxRetryReached,
			expectedElapsedSecond: 2,
		},
		// cancel before first iteration succeed
		{
			durationCancel:        time.Duration(0),
			useToken:              true,
			ticksToSuccess:        1,
			expectedTicks:         1,
			expectedResult:        nil,
			expectedElapsedSecond: 0,
		},
		// cancel before first iteration failed
		{
			durationCancel:        time.Duration(0),
			useToken:              true,
			ticksToSuccess:        1,
			expectedTicks:         1,
			expectedResult:        nil,
			expectedElapsedSecond: 0,
		},
		// cancel right after first iteration succeed
		{
			durationCancel:        time.Duration(500) * time.Millisecond,
			useToken:              true,
			ticksToSuccess:        1,
			expectedTicks:         1,
			expectedResult:        nil,
			expectedElapsedSecond: 0,
		},
		// cancel before the 3rd operation which was supposed to succeed
		{
			durationCancel:        time.Duration(1500) * time.Millisecond,
			useToken:              true,
			ticksToSuccess:        3,
			expectedTicks:         2,
			expectedResult:        ErrorTaskCancelled,
			expectedElapsedSecond: 1,
		},
		// retry 2 times, passed, cancel after
		{
			durationCancel:        time.Duration(1500) * time.Millisecond,
			useToken:              true,
			ticksToSuccess:        2,
			expectedTicks:         2,
			expectedResult:        nil,
			expectedElapsedSecond: 1,
		},
	}

	for index, tc := range testCase {
		ticks := 0
		var token CancellationToken
		if tc.useToken {
			token = NewCancellationToken()
			go func() {
				time.Sleep(tc.durationCancel)
				token.Cancel()
			}()
		}
		time.Sleep(100)
		start := time.Now()
		result := RetryOperation(func() error {
			ticks++
			if ticks == tc.ticksToSuccess {
				return nil
			}
			return fmt.Errorf("Some error")
		}, retryPolicy, token)
		elapsed := int(time.Now().Sub(start) / time.Second)
		assertEquals(t, index, "ticks", tc.expectedTicks, ticks)
		assertEquals(t, index, "result", tc.expectedResult, result)
		assertEquals(t, index, "time", tc.expectedElapsedSecond, elapsed)
	}
}

func assertEquals(t *testing.T, index int, msg string, expected interface{}, actual interface{}) {
	if expected != actual {
		if expectedAsString, isString := expected.(fmt.Stringer); isString {
			if actualAsString, isString := actual.(fmt.Stringer); isString {
				if expectedAsString == actualAsString {
					return
				}
			}
		}
		if expectedError, isError := expected.(error); isError {
			if actualError, isError := actual.(error); isError {
				if expectedError.Error() == actualError.Error() {
					return
				}
			}
		}
		t.Errorf("%d:%s expected: %d; actual %d\n", index, msg, expected, actual)
		t.Fail()
	}
}
