package lib

import (
	"google.golang.org/protobuf/types/known/anypb"
)

type Completer interface {
	CompleteCall(a *anypb.Any, resultErr int32)
}

type ErrorType interface {
	~int32
}

type Future[T any, U ErrorType] struct {
	resolveFn     func(t T)
	rejectFn      func(U)
	resolved      bool
	resolvedValue T
	rejected      bool
	rejectedValue U
}

func NewFuture[T any, U ErrorType](resolve func(T), reject func(U)) *Future[T, U] {
	return &Future[T, U]{
		resolveFn: resolve,
		rejectFn:  reject,
	}
}

func (f *Future[T, U]) CompleteCall(result T, resultErr U) {
	if resultErr != 0 || f.resolveFn == nil {
		if f.rejectFn != nil {
			f.rejectFn(resultErr)
		}
		f.rejected = true
		f.rejectedValue = resultErr
	} else {
		f.resolveFn(result)
		f.resolved = true
		f.resolvedValue = result
	}
}

func (f *Future[T, U]) Success(fn func(T)) {
	if f.rejected {
		return // no way this fn call ever be called
	}
	// call immediately if resolved
	if f.resolved {
		fn(f.resolvedValue)
		return
	}
	f.resolveFn = fn
}

func (f *Future[T, U]) Failure(fn func(U)) {
	if f.resolved {
		return // cannot be reached
	}
	if f.rejected {
		fn(f.rejectedValue)
		return
	}
	f.rejectFn = fn
}

//
// BaseFature (single value with Handle())
//

type BaseFuture[T any] struct {
	resolveFn     func(T)
	resolved      bool
	resolvedValue T
}

func NewBaseFuture[T any]() *BaseFuture[T] {
	return &BaseFuture[T]{
		resolved: false,
	}
}

func NewBaseFutureWithValue[T any](t T) *BaseFuture[T] {
	return &BaseFuture[T]{
		resolved:      true,
		resolvedValue: t,
	}
}

func (f *BaseFuture[T]) Set(t T) {
	if f.resolved {
		panic("attempt to set value for panic, but panic is already resolved")
	}
	if f.resolveFn != nil {
		f.resolveFn(t)
	} else {
		f.resolved = true
		f.resolvedValue = t
	}
}

func (f *BaseFuture[T]) Handle(fn func(T)) {
	if f.resolved {
		fn(f.resolvedValue)
		return
	}
	f.resolveFn = fn
}
