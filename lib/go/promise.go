package lib

type ErrorType interface {
	~int32
}

type BaseType interface {
	~int32 | ~int64 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
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

func NewFutureValue[T BaseType](t T) *Future {
	return &Future{
		resolvedValue: t,
	}
}

func NewFutureError[U ErrorType](u U) *Future {
	return &Future{
		rejectedConst: true,
		rejectedValue: int32(u),
	}
}

func (f *Future) SuccessAny(fn func(any)) {
	if f.resolvedValue != nil {
		fn(f.resolvedValue)
		return
	}
	newFuture := NewFuture[any, int32](fn, nil)
	f.downstream = append(f.downstream, newFuture)
}

func (f *Future) FailureAny(fn func(int32)) {
	if f.rejectedConst {
		fn(f.rejectedValue)
		return
	}
	newFuture := NewFuture[any, int32](nil, fn)
	f.downstream = append(f.downstream, newFuture)
}
