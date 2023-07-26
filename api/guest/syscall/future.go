package syscall

import (
	"context"

	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"
)

// LaunchFuture is the return type of Launch() on the guest side.  The guest
// receives the future in response to a Launch() call and should attach Success()
// and Failure() funcs as needed.  If the Completed() call is true, the methods
// added in Success() or Failure() will be called immediately.
type LaunchFuture struct {
	fut *future.Method[*syscall.LaunchResponse, syscall.KernelErr]
}

// Success should be called to add a function to be called when the Launch()
// has fully completed and did so successfully.
func (l *LaunchFuture) Success(fn func(*syscall.LaunchResponse)) {
	l.fut.Success(fn)
}

// Failure should be called to add a function to be called when the Launch()
// has fully completed and was unsuccessful.
func (l *LaunchFuture) Failure(fn func(syscall.KernelErr)) {
	l.fut.Failure(fn)
}

// Completed returns if the given LaunchFuture has already been completed.
// This might be interesting for some guests if they wish to bypass the
// futures mechanism in the case where the Launch() has already finished.
// Note that quick completion of Launch() is NOT guaranteed in all cases so
// clients that use this method to bypass futures must also have a backup
// approach for when the Launch() takes some time to complete.
func (l *LaunchFuture) Completed() bool {
	return l.fut.Completed()
}

func NewLaunchFuture() *LaunchFuture {
	return &LaunchFuture{
		fut: future.NewMethod[*syscall.LaunchResponse, syscall.KernelErr](nil, nil),
	}
}

// CompleteMethod fills in the results for a Method future and it works the
// same for LaunchResponse.
func (l *LaunchFuture) CompleteMethod(ctx context.Context, lr *syscall.LaunchResponse, err syscall.KernelErr) {
	l.fut.CompleteMethod(ctx, lr, err)
}

// ExitFuture is the return type of Exit() on the guest side.  This is a future
// because it is not certain exactly when the Exit will actually occur.  Further,
// the exit itself might fail, so the program may not exit at all.
type ExitFuture struct {
	fut *future.Method[*syscall.ExitResponse, syscall.KernelErr]
}

// Success should be called to add a function to be called when the Exit()
// has fully completed and did so successfully.  Adding an exit function here
// can useful to clean up resources, with the understanding that the program
// is _about_ to exit.  The function fn should NOT exit the program with a
// call like os.Exit(1), this will happen once all the Success functions on
// the ExitFuture have been called.
func (l *ExitFuture) Success(fn func(*syscall.ExitResponse)) {
	l.fut.Success(fn)
}

// Failure should be called to add a function to be called when the Exit()
// has fully completed and was unsuccessful.  Note that this situation is
// a serious internal error in parigot when the given method  fn is called.
// It is appropriate to take drastic measures like `os.Exit(1)` to force the
// abort of the program.
func (l *ExitFuture) Failure(fn func(syscall.KernelErr)) {
	l.fut.Failure(fn)
}

// NewExitFuture returns an initialized exit future.  It is not useful to
// attempt to determine if the exit has "completed" as the program would no
// longer exit.
func NewExitFuture() *ExitFuture {
	return &ExitFuture{
		fut: future.NewMethod[*syscall.ExitResponse, syscall.KernelErr](nil, nil),
	}
}

// CompleteMethod is used to complete a previously defined future of type ExitFuture.
func (l *ExitFuture) CompleteMethod(ctx context.Context, lr *syscall.ExitResponse, err syscall.KernelErr) {
	l.fut.CompleteMethod(ctx, lr, err)
}
