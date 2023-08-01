package syscall

import (
	"context"

	syscall "github.com/iansmith/parigot/g/syscall/v1"
	future "github.com/iansmith/parigot/lib/go/future"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type LaunchCompleter struct {
	fut *LaunchFuture
}

func (l *LaunchCompleter) CompleteMethod(ctx context.Context, a proto.Message, e int32) syscall.KernelErr {
	out := &syscall.LaunchResponse{}
	if err := a.(*anypb.Any).UnmarshalTo(out); err != nil {
		return syscall.KernelErr_UnmarshalFailed
	}
	l.fut.CompleteMethod(ctx, out, syscall.KernelErr(e))
	return syscall.KernelErr_NoError
}

func (l *LaunchCompleter) Success(succFunc func(proto.Message)) {
	x := func(m *syscall.LaunchResponse) {
		succFunc(m)
	}
	l.fut.Success(x)

}
func (l *LaunchCompleter) Failure(failFunc func(int32)) {
	x := func(err syscall.KernelErr) {
		failFunc(int32(err))
	}
	l.fut.Failure(x)
}

func NewLaunchCompleter(f *LaunchFuture) future.Completer {
	return &LaunchCompleter{f}
}

type ExitCompleter struct {
	fut *ExitFuture
}

func (l *ExitCompleter) CompleteMethod(ctx context.Context, a proto.Message, e int32) syscall.KernelErr {
	out := &syscall.ExitResponse{}
	if err := a.(*anypb.Any).UnmarshalTo(out); err != nil {
		return syscall.KernelErr_UnmarshalFailed
	}
	l.fut.CompleteMethod(ctx, out, syscall.KernelErr(e))
	return syscall.KernelErr_NoError
}

func (l *ExitCompleter) Success(succFunc func(proto.Message)) {
	x := func(m *syscall.ExitResponse) {
		succFunc(m)
	}
	l.fut.Success(x)

}
func (l *ExitCompleter) Failure(failFunc func(int32)) {
	x := func(err syscall.KernelErr) {
		failFunc(int32(err))
	}
	l.fut.Failure(x)
}

func NewExitCompleter(f *ExitFuture) future.Completer {
	return &ExitCompleter{f}
}
