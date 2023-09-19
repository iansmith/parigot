package future

// tricky: this code runs completely on the client
// side of a parigot program.  however, to run the
// tests easily, it cannot have //go:build wasip1 because
// that prevents the tests from running locally.

import (
	"bytes"
	"context"
	"fmt"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Completer is the interface that means that a given
// type can be "completed" at a later time.  This is used
// only for Methods.
type Completer interface {
	CompleteMethod(ctx context.Context, msg proto.Message, resultErr int32, orig id.HostId) syscall.KernelErr
	Success(func(proto.Message))
	Failure(func(int32))
	Completed() bool
	Cancel()
}

// Invoker is the interface that means that a given
// type be run as an implementation of a function..
type Invoker interface {
	// Invoke has to do the work to unmarshal the msg because it knows
	// the specific type to use whereas the caller does not.
	Invoke(ctx context.Context, msg *anypb.Any) Completer
}

type ErrorType interface {
	~int32
}

// Method is a special type of future that is used frequently
// in parigot because all the methods of a service, and the
// methods of clients that use that same service, must return
// this type.  It has the special behavior that when CompleteMethod
// is called on this Method, the error value is compared to zero
// and this determines if the Success (error value is 0)
// or Failure (error is not 0) handler function is called.
//
// It is thus impossible to have a Method that can behave in
// a failed way (call to Failure) based on the return value
// being 0.  In this case, use a Base[int32], as parigot does.
type Method[T proto.Message, U ErrorType] struct {
	resolveFn               func(t T)
	rejectFn                func(U)
	resolvedValue           T
	rejectedValue           U
	completed               bool
	wasResolved             bool
	resolveSucc, rejectSucc *Method[T, U]
	waitingId               id.CallId
}

// NewMethod return as method future with two types given.  The T
// type (first) must be a proto.Message and typically is a Response object
// from a previous call to the Method.  The error value, U, is typically
// a named enum that is used for error signaling by the method called.
func NewMethod[T proto.Message, U ErrorType](resolve func(T), reject func(U)) *Method[T, U] {
	result := &Method[T, U]{
		resolveFn: resolve,
		rejectFn:  reject,
		waitingId: id.NewCallId(),
	}
	// raw := uintptr(unsafe.Pointer(result))
	// log.Printf("--------------------------------- %x\n\n", raw)
	// debug.PrintStack()
	// log.Printf("---------------------------------END %x\n\n", raw)
	return result
}

// WaitingId is useful only to the go client side library.  The WaitingId
// is a repurposing of the CallId to create a key value, a string, for use
// in a map, since Method[T,U] is not a valid key type in go.
func (m *Method[T, U]) WaitingId() string {
	return m.waitingId.String()
}

// Cancel causes a future to be marked completed and also to remove any and all
// possible calls to a Sucess() or Failure() function later.  This enforces
// that a Cancel() is permanent, even if the future is "completed" later.
// Calling Cancel() on an already completed future will be ignored.
func (m *Method[T, U]) Cancel() {
	if m.completed {
		return
	}
	m.completed = true
	m.rejectSucc = nil
	m.resolveSucc = nil
	m.resolveFn = nil
	m.rejectFn = nil
}

// findLast is a utility function which finds the last
// Method future in a sequence of futures.  Whether it finds error
// or success values based on the isSuccess paramater.
func (f *Method[T, U]) findLast(isSuccess bool) *Method[T, U] {
	var current *Method[T, U]
	if isSuccess {
		current = f.resolveSucc
	} else {
		current = f.rejectSucc
	}
	prev := f
	for current != nil {
		prev = current
		if isSuccess {
			current = current.resolveSucc
		} else {
			current = current.rejectSucc
		}
	}
	return prev
}

// CompleteMethod is called to indicate that the outcome, or value,
// of the future is now known.  This method is typically called by
// the infrastructure of Parigot, but it can be useful to call this
// method directly in tests.  Calling this method on completed Method
// future will be ignored.
func (f *Method[T, U]) CompleteMethod(ctx context.Context, result T, resultErr U) {
	if f.completed {
		return
	}
	if resultErr != 0 {
		current := f
		for current != nil {
			if current.rejectFn != nil {
				current.rejectFn(resultErr)
			}
			current.completed = true
			current.wasResolved = false
			current.rejectedValue = resultErr
			current = current.rejectSucc
		}
		return // ran them all
	}
	current := f
	for current != nil {
		if current.resolveFn != nil {
			current.resolveFn(result)
		}
		current.completed = true
		current.wasResolved = true
		current.resolvedValue = result
		current = current.resolveSucc
	}
}

// Completed returns true if this method has already completed.
func (f *Method[T, U]) Completed() bool {
	return f.completed
}

// ValueResponse may not do what you expect: This function does not
// force the world to stop and wait for the Response in question to be
// received. It can only be trusted when the function Completed() returns
// true  and the function WasSuccess() returns true.  This function
// returns the value of a response (type T) on a completed
// method
func (f *Method[T, U]) ValueResponse() T {
	return f.resolvedValue
}

// ValueErr may not do what you expect: This function does not
// force the world to stop and wait for the Error in question to be
// received. It can only be trusted when the function Completed() returns
// true and the function WasSuccess() returns false.  It returns the
// value of the error on a completed Method.
func (f *Method[T, U]) ValueErr() U {
	return f.rejectedValue
}

// String() returns a human-friendly version of this Method future.
// It shows it is resolved and if so, if the completion was an error.
func (f *Method[T, U]) String() string {
	var t T
	var u U
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("[Method:%T,%T]", t, u))
	if f.completed {
		buf.WriteString("completed:")
		if f.rejectedValue != 0 {
			buf.WriteString(fmt.Sprintf("error: %+v", f.resolvedValue))
		} else {
			buf.WriteString("success")
		}
	} else {
		haveSuccess := fmt.Sprintf("have success?%v", f.resolveFn != nil)
		haveFailure := fmt.Sprintf("have failure?%v", f.rejectFn != nil)
		buf.WriteString("waiting. " + haveSuccess + "," + haveFailure)
	}
	return buf.String()
}

// Success provides a function to be called if the Method returns a
// success.  Calling Success() on an already completed method causes the
// code supplied in the success method to be run immediately if the
// future was resolved successfully.
func (f *Method[T, U]) Success(fn func(T)) {
	if f.completed {
		if f.wasResolved {
			fn(f.resolvedValue)
		}
		// nothing to do with the error
		return
	}
	last := f.findLast(true)
	next := NewMethod[T, U](fn, nil)
	last.resolveSucc = next
}

// Failure provides a function to be called if the Method completion supplies
// a non zero error value. Calling failure on a completed Method that
// had an error causes the given function to run immediately.
func (f *Method[T, U]) Failure(fn func(U)) {
	if f.completed {
		if !f.wasResolved {
			fn(f.rejectedValue)
		}
		return
	}
	last := f.findLast(false)
	next := NewMethod[T, U](nil, fn)
	last.rejectSucc = next
}

// WasSuccess returns true if the Method is completed and finished
// as a sucess.  Before a Method is completed, it returns false.
func (f *Method[T, U]) WasSuccess() bool {
	if !f.completed {
		return false
	}
	return f.wasResolved
}

//
// Base (single value with Handle())
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
// It is both allowed and useful to
// call Handle() multiple times on the same future.  If there are
// handle functions that remain uncalled because they were added
// after the previous call to Set was executed. When the future
// is completed via Set(), all of the incomplete, supplied Handle() functions
// will run and in the order they were called on the future.
//
// If Handle() is called on a future which has completed running
// all of its registered handlers, the Handle() function provided
// is run immediately.  This is usually what you want.
// If you wish to delay the execution until the next Set() call use HandleLater.

func (f *Base[T]) Handle(fn func(T)) {
	if f.alreadyRunAll() {
		fn(f.resolvedValue)
		return
	}
	f._handle(fn)
}

// HnadleLater is used in the rare instance can where you
// have a future that has possibly completed all of
// its Handle() functions and you wish to delay the excution
// of fn until the next Set() call.  Note that the default
// behavior of Handle() would be to run fn immediately, and thus
// you only need this function if you call must Set multiple times
// on the same future with the possibility that is already completed.
func (f *Base[T]) HandleLater(fn func(T)) {
	f._handle(fn)
}

// create a new future and put it at the end of the chain
// rooted by f.
func (f *Base[T]) _handle(fn func(T)) {
	last := f.findLast()
	end := NewBase[T]()
	last.successor = end
	end.resolveFn = fn
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

// Completed returns true if all the Handle() functions
// on this future have run.  Note that is can be
// changed by the addition of new Handle() functions
// via Handle().
func (f *Base[T]) Completed() bool {
	return f.alreadyRunAll()
}

// Strings returns a human-friendly representation of this Base
// futuer.  It returns if the future is complete or not.
func (f *Base[T]) String() string {
	var t T
	if f.alreadyRunAll() {
		return fmt.Sprintf("Base[%T]:completed", t)
	}
	return fmt.Sprintf("Base[%T]:waiting", t)

}

// Cancel causes the future's state to be cleared and the future
// to be marked completed.  Any calls to
// Set() that occur after Cancel() and before any other calls
// to Handle() will have no effect.  Any existing chain
// of Handle() functions will be removed from the future by Cancel().
// Since the call to Cancel() marks the
func (f *Base[T]) Cancel() {
	f.resolved = true
	f.resolveFn = nil
	f.successor = nil
}

// All waits for all its dependent futures to complete and if they
// all complete successfully, it calls the Success function, otherwise
// the index of a failing future is sent to the Failure() method.
func All[T proto.Message, U ErrorType](dep ...*Method[T, U]) *AllFuture[T, U] {
	return NewAllFuture[T, U](dep)
}

// AllFuture is the underlying Future type for a call to the All() function.
// As with Base and Method, this future queues all the calls to Success() and
// Failure().
type AllFuture[T proto.Message, U ErrorType] struct {
	dep        []*Method[T, U]
	success    int
	successFut *Base[bool]
	failFut    *Base[int]
}

func NewAllFuture[T proto.Message, U ErrorType](dep []*Method[T, U]) *AllFuture[T, U] {
	result := &AllFuture[T, U]{
		dep: dep,
	}
	for i, fut := range result.AllDependent() {
		copyI := i
		fut.Success(func(t T) {
			result.addOneSuccess(copyI)
		})
		fut.Failure(func(u U) {
			result.fail(copyI)
		})
	}
	result.successFut = NewBase[bool]()
	result.failFut = NewBase[int]()
	return result
}

func (a *AllFuture[T, U]) AllDependent() []*Method[T, U] {
	return a.dep
}
func (a *AllFuture[T, U]) addOneSuccess(_ int) {
	a.success++
	if a.success == len(a.dep) {
		a.successFut.Set(true)
	}
}
func (a *AllFuture[T, U]) fail(index int) {
	a.failFut.Set(index)
}

func (a *AllFuture[T, U]) Success(fn func()) {
	a.successFut.Handle(func(_ bool) {
		fn()
	})
}
func (a *AllFuture[T, U]) Failure(fn func(int)) {
	a.failFut.Handle(func(index int) {
		fn(index)
	})
}
