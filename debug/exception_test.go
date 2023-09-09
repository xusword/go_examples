package debug

import (
	"fmt"
	"testing"
)

// This is not a test
func TestThrow(t *testing.T) {
	fmt.Printf("%s\n", Problem())
}

func Problem() error {
	return Throw(Drill(), "Problem %s", "Meh")
}

func Drill() error {
	return Throw(RootCause(), "Drill %s", "Yeah")
}

func RootCause() error {
	return Throw(nil, "RootCause %s", "Exception")
}
