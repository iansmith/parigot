//go:build wasip1

package exit

import (
	"context"
)

// exitFunc is the type of the function used during atExit cleanup.
type exitFunc func(ctx context.Context)

// atExit is the list of exit funcs to run, registered with AtExit()
var atExit = []exitFunc{}

// Call AtExit() with a function to call to clean up resources. AtExit functions
// are called in reverse order to their registration.  Note than AtExit function
// should not run for a long period of time (say, more than 50 milliseconds) because
// the program's termination is imminent and the number of AtExit functions that
// need to run can be large.
func AtExit(fn exitFunc) {
	atExit = append(atExit, fn)
}

// ExecuteAtExit runs the previously registered AtExit functions. This call is
// used by the infrastructure to run the AtExit functions at exit.  User code
// should never need to call this function.
func ExecuteAtExit(ctx context.Context) {
	for i := len(atExit) - 1; i >= 0; i-- {
		atExit[i](ctx)
	}
}
