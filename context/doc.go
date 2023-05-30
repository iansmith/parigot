// Package context provides to crticial utilities for both user and system code.  Inside a context.Context we expect
// to have two objects that can be manipulated.
//
// CurrentTime: This is the parigot current time.  You can access this with CurrentTime(ctx) and returns a time object.
// It is critical that this be used in favor of time.Now() or similar because this value can and will be manipulated
// to "fake" times for testing. E.g. if you want to test a function that runs only at midnight, you can use the
// SetCurrentTimeTest(ctx,time) function to test at 23:59:59, 00:00:00, and 00:00:01.  This value, when not being
// returns the current time and that time is fixed at the start of a method call.  Thus, all the log lines in the
// same function will occur "at the same time" which makes it easier to debug things in a complex system.
//
// LogContainer:  This is where all the parigot logs get collected.  A log container can be created and inserted
// into the context to create a container that contains only a particular bit of log information.  The function
// Dump() dumps the contents of the log container to the terminal.   Log containers can also be transferred over
// the network so as create a "unified log" that includes log messages from many different microservices.
//
// The current LogContainer implementation is quite slow and this is because it is intended to be used for development
// and the clarity of the logs is the most important thing.  We expect to have high performance log containers
// and log containers that connect to particular logging services.
package context
