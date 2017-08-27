package debug

import "testing"
import "fmt"

// This is not a test
func TestThrow(t *testing.T) {
	exception := Problem().(*Exception)
	fmt.Printf("%v\n", exception.String())
}

func Problem() error {
	return Throw(Drill(), "Problem %s", "Meh")
}

func Drill() *Exception {
	return Throw(RootCause(), "Drill %s", "Yeah")
}

func RootCause() *Exception {
	return Throw(nil, "RootCause %s", "Exception")
}
