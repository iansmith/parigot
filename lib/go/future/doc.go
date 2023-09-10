//go:build wasip1

// The concept of a "future" is important in parigot.  For
// any programming language that can create closures,
// this concept allows one to write parigot programs
// that are singly-threaded.  A future is similar to
// a promise in Javascript and related languages.
//
// A future represents a computation that has not
// completed yet, but is expected to be completed in
// the near future.  A future is "completed" when its
// value or values are known because the computation that
// was being waited on has finished.  In the case of
// parigot, frequently that computation has been completed
// by another process or another process on another machine.
//
// parigot offers two types of futures, Base and Method.
// Base futures represent a single value and the result
// can be acted upon in the method Handle().  Base futures
// are used when a single value, such as bool, is being
// computed but it is not completed as-of yet.  The method
// Ready() on a service returns a Base future to inidicate
// if a service can start or not (bool).
//
// Method futures are a special case of future that is used frequently
// in parigot.  Method futures represent the value of a
// remote method call that has not completed yet.  The value of
// a method call is computed by some other program and
// then set to indicate the value(s) of that this method
// has returned.  In the simple case of all the programs
// being run in a single address space (process) this
// other program is another WASM binary executing inside
// the same WASM engine.
package future
