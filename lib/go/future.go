package lib

import (
	"github.com/iansmith/parigot/g/protosupport/v1"
)

type ErrorType interface {
	~int32
}

type Future struct {
	resolveFn     func(any)
	rejectFn      func(int32)
	resolved      bool
	resolvedValue any
	rejected      bool
	rejectedValue int32
	rejectedConst bool
	downstream    []*Future
}

func NewFuture[T any, U ErrorType](resolve func(T), reject func(U)) *Future {
	resolveWrap := func(a any) {
		resolve(a.(T))
	}
	rejectWrap := func(a int32) {
		reject(U(a))
	}

	return &Future{
		resolveFn: resolveWrap,
		rejectFn:  rejectWrap,
	}
}

func NewFutureError[U ErrorType](u U) *Future {
	return &Future{
		rejectedConst: true,
		rejectedValue: int32(u),
	}
}

func (f *Future) CompleteCall(result any, resultErr int32) {
	if resultErr != 0 {
		if f.rejectFn != nil {
			f.rejectFn(resultErr)
		}
		f.rejected = true
		f.rejectedConst = true
		f.rejectedValue = resultErr
	} else {
		if result != nil {
			if f.resolveFn != nil {
				f.resolveFn(result)
			}
		}
		f.resolved = true
		f.resolvedValue = result
	}
	// float downstream
	for _, d := range f.downstream {
		d.CompleteCall(result, resultErr)
	}
}

func (f *Future) SuccessAny(fn func(any)) {
	if f.rejected {
		return // no way this fn call ever be called
	}
	// call immediately if resolved
	if f.resolved {
		fn(f.resolvedValue)
		return
	}
	newFuture := NewFuture[any, int32](fn, nil)
	f.downstream = append(f.downstream, newFuture)
}

func (f *Future) FailureAny(fn func(int32)) {
	if f.resolved {
		return // cannot be reached
	}
	if f.rejected {
		fn(f.rejectedValue)
		return
	}
	newFuture := NewFuture[any, int32](nil, fn)
	f.downstream = append(f.downstream, newFuture)
}

func NewFutureId() *Future {
	return NewFuture[*protosupport.IdRaw, int32](nil, nil)
}

type BaseFuture[T any] struct {
	resolveFn     func(T)
	resolved      bool
	resolvedValue T
	downstream    []*Future
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
