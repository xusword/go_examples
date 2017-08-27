package debug

import (
	"fmt"
	sysDebug "runtime/debug"
	"strings"
)

type Exception struct {
	Message    string
	Cause      error
	StackTrace []string
}

func (e *Exception) Error() string {
	stackJoined := strings.Join(e.StackTrace, "\n")
	if e.Cause == nil {
		return fmt.Sprintf("Error: %s\nStack Trace: %s", e.Message, stackJoined)
	}
	return fmt.Sprintf("Error: %s\nStack Trace: %s\nCaused by:\n%s", e.Message, stackJoined, e.Cause)
}

func Throw(cause error, msg string, args ...interface{}) *Exception {
	fullStack := string(sysDebug.Stack())
	fullStackLines := strings.Split(fullStack, "\n")
	return &Exception{
		Message:    fmt.Sprintf(msg, args...),
		Cause:      cause,
		StackTrace: fullStackLines[5:],
	}
}
