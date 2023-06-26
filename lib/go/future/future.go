package future

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Completer is the interface that means that a given
// type can be "completed" at a later time.  This is used
// only for Methods.
type Completer interface {
	CompleteCall(a *anypb.Any, resultErr int32)
}

type ErrorType interface {
	~int32
}

// Method is a special type of future that is used frequently
// in Parigot because all the methods of a service, and the
// methods of clients that use that same service, must return
// this type.  It has the special behavior that when CompletMethod
// is called on this method, the error value is compared to zero
// and this determines if the Success (error value is 0)
// or Failure (error is not 0) function is called.
//
// It is thus impossible to have a Method that can behave in
// a failed way (call to Failure) based on the return value
// being 0.  In this case, use a BaseFuture[int32].
type Method[T proto.Message, U ErrorType] struct {
	resolveFn     func(t T)
	rejectFn      func(U)
	resolved      bool
	resolvedValue T
	rejected      bool
	rejectedValue U
}

// NewMethod return as method future with two types given.  The T
// type (first) must be a proto.Message and typically is a Response object
// from a previous call to the Method.  The error value, U, is typically
// a named enum that is used for error signaling by the method called.
func NewMethod[T proto.Message, U ErrorType](resolve func(T), reject func(U)) *Method[T, U] {
	return &Method[T, U]{
		resolveFn: resolve,
		rejectFn:  reject,
	}
}

// CompleteMethod is called to indicate that the outcome, or value,
// of the future is now known.  This method is typically called by
// the infrastructure of Parigot, but it can be useful to call this
// method directly in tests.  Calling CompleteMethod() on an already
// completed method future will panic as such a call to CompleteMethod()
// makes no sense (a method can only be completed once).
func (f *Method[T, U]) CompleteMethod(result T, resultErr U) {
	if f.resolved || f.rejected {
		panic("cannot call CompleteMethod on a future that is already completed")
	}
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

// Success provides a function to be called if the Method returns a
// success.  Unlike Base[T], the given function fn will replace the
// function, if any, that was previously given as the success function.
// Calling Success() on an already completed method is useless and
// causes a panic (a method can only be 'completed' once).
func (f *Method[T, U]) Success(fn func(T)) {
	if f.resolved || f.rejected {
		panic("Success() called on already completed Method")
	}
	f.resolveFn = fn
}

// Failure provides a function to be called if the Method returns a
// non zero error value.  Unlike Base[T], the given function fn will
// replace any previous function that was given as the failure function.
// Calling Failure() on an already completed method is useless and
// causes a panic (a method can only be 'completed' once).
func (f *Method[T, U]) Failure(fn func(U)) {
	if f.resolved || f.rejected {
		panic("Success() called on already completed Method")
	}
	f.rejectFn = fn
}

// Completed returns true if the Method has been completed.
func (f *Method[T, U]) Completed() bool {
	return f.rejected || f.resolved
}

//
// BaseFuture (single value with Handle())
//

// Base[T] represents a future computation resulting in
// a T.  This is useful for simple value types such bool and int64.
// It only has one user-visible method which is Handle() and
// that is used to set a handler for when the value of type T
// actually arrives.
type Base[T any] struct {
	resolveFn     func(T)
	resolved      bool
	resolvedValue T
	successor     *Base[T]
}

// NewBase returns a new pointer at a Base[T].  The value will
// be the zero value of T and the future is not marked completed.
// If you wish to make the zero value the result of the future
// you should use NewBaseWithValue or Set().
func NewBase[T any]() *Base[T] {
	return &Base[T]{
		resolved: false,
	}
}

// NewBaseWithValue creates a new pointer a Base[T] with the given
// value and with future marked as completed.
func NewBaseWithValue[T any](t T) *Base[T] {
	return &Base[T]{
		resolved:      true,
		resolvedValue: t,
	}
}

// Set causes the future it is called on be marked as completed
// with the given value.  This will cause all registered
// Handle() functions to run. Set can be called multiple times
// and the effect is that only the as yet incomplete
// Handle() functions will be executed as a result. These
// previously unexecuted handle functions will be marked and
// and have their result values set to the value of *this* call to
// set.  It is thus possible that different handlers will
// run with different values as their parameters.  Set
// returns true if any Handle functions were run.
func (f *Base[T]) Set(t T) bool {
	runAny := false
	value := t
	curr := f
	for curr != nil {
		if !curr.resolved {
			if curr.resolveFn != nil {
				curr.resolveFn(value)
			}
			curr.resolved = true
			curr.resolvedValue = t
			runAny = true
		}
		curr = curr.successor
	}
	return runAny
}

// Handle sets up the handler for the future it is called on.
// This function may create a new future as part of its operation
// and the new future is returned as the result; if no new future
// has been created, the result will be the same as the future
// Handle() was called on.  It is both allowed and useful to
// call Handle() multiple times on the same future if there are
// handle functions that remain uncalled because they were added
// after the previous call to Set was executed. When the future
// is completed, all of the incomplete, supplied Handle() functions
// will run and in the order they were called on the future.
//
// If Handle() is called on a future which has completed running
// all of its registered handlers, the Handle() function provided
// is run immediately.  This is usually what you want.
// If you wish to delay the execution until the next Set() call use HandleLater.

func (f *Base[T]) Handle(fn func(T)) *Base[T] {
	if f.alreadyRunAll() {
		fn(f.resolvedValue)
		return f
	}
	return f._handle(fn)
}

// HnadleLater is used in the (relatively rare) can where you
// have a future that has already completed all of
// its Handle() functions and you wish to delay the excution
// of fn until the next Set() call.  Note that the default
// behavior of Handle() would be to run fn immediately and thus
// you only need this function if you call Set multiple times
// on the same future.
func (f *Base[T]) HandleLater(fn func(T)) *Base[T] {
	return f._handle(fn)
}

// create a new future and put it at the end of the chain
// rooted by f.
func (f *Base[T]) _handle(fn func(T)) *Base[T] {
	last := f.findLast()
	end := NewBase[T]()
	last.successor = end
	end.resolveFn = fn
	return end
}

// findLast is a utility function which finds the last
// future in a sequence of futures, such as those chains of
// futures created by Handle().
func (f *Base[T]) findLast() *Base[T] {
	current := f.successor
	prev := f
	for current != nil {
		prev = current
		current = current.successor
	}
	return prev
}

// alreadyRunAll returns true if there are no currently unresolved
// functions in the sequence of futures.
func (f *Base[T]) alreadyRunAll() bool {
	current := f
	for current != nil {
		if !current.resolved {
			return false
		}
		current = current.successor
	}
	return true
}
