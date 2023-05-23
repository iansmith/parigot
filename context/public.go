package context

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ServerGoContext returns a new context based on ctx with source ServerGo and the given function name.
// This should be called before entering plugins that define host functions.
func ServerGoContext(ctx context.Context, funcName string) context.Context {
	return newContext(ctx, ServerGo, funcName)
}

// Client returns a new context based on ctx with source Client and the given function name.
// This should be called before entering code that is client wasm code, like the main
// of an application or the start of test.
func ClientContext(ctx context.Context, funcName string) context.Context {
	return newContext(ctx, Client, funcName)
}

// ServerWasm returns a new context based on ctx with source ServerWasm and the given function name.
// This should be called before entering code that is the _implementation_ of a service that is
// implemented solely in guest wasm.
func ServerWasmContext(ctx context.Context, funcName string) context.Context {
	return newContext(ctx, ServerWasm, funcName)
}

// CurrentTime should *always* be used in both host and guest code to get the current
// time.  This should be done because it allows testing code that needs to run at a
// particular time (hourly, midnight, etc) to be tested more easily.
func CurrentTime(ctx context.Context) time.Time {
	t := ctx.Value(ParigotTime)
	if t != nil && !t.(time.Time).IsZero() {
		return t.(time.Time)
	}
	return time.Now()
}

// SetCurrentTimeTest should be used only for testing. It allows the current time
// to be set (and not be advanced) by later code that uses CurrentTime().
func SetCurrentTimeTest(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, ParigotTime, t)
}

// Dump is called to dump the current contents of the current log container.  This
// function is intended to be used with functions that create a particular new
// context that is being used for a particular purpose.
func Dump(ctx context.Context) {
	cont := GetContainer(ctx)
	if cont == nil {
		log.Println("no log container present inside context")
		return
	}
	cont.Dump()
}

// Logf is a shorthand for a call to LogFull with the currently set source in ctx.
func Logf(ctx context.Context, level LogLevel, spec string, rest ...interface{}) {
	LogFullf(ctx, level, UnknownS, "", spec, rest...)
}

// Errorf is a shorthand for a call to LogFull with the currently set source in ctx, and
// log level being Error.
func Errorf(ctx context.Context, spec string, rest ...interface{}) {
	LogFullf(ctx, Error, UnknownS, "", spec, rest...)
}

// Debugf is a shorthand for a call to LogFull with the currently set source in ctx, and
// log level being Debug.
func Debugf(ctx context.Context, funcName string, spec string, rest ...interface{}) {
	LogFullf(ctx, Debug, UnknownS, funcName, spec, rest...)
}

// ClientLogf is just like Logf except is sets the source to be client.  This is
// useful (with no context param) because client's (like the main of an app or a test
// entry point) usually don't have a context.
func ClientLogf(level LogLevel, spec string, rest ...interface{}) {
	LogFullf(context.Background(), level, Client, "", spec, rest...)
}

// ClientLogf is just like Debugf except is sets the source to be client.  This is
// useful (with no context param) because client's (like the main of an app or a test
// entry point) usually don't have a context.
func ClientDebugf(funcName string, spec string, rest ...interface{}) {
	LogFullf(context.Background(), Debug, UnknownS, funcName, spec, rest...)
}

// Internal is for internal use only.  It creates a log line attributed to Parigot.
func Internal(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParigotSource, Source(Parigot))
}

// CurrentTimeString is a wrapper around CurrentTime that returns a string representation
// of the current time in a standard form (RFC822Z).
func CurrentTimeString(ctx context.Context) string {
	return CurrentTime(ctx).Format(time.RFC822Z)
}

// CallTo returns a new context with the current function updated to s.
func CallTo(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, ParigotFunc, s)
}

// CallTo returns a new context with the current source set to ServerGo.  This is useful
// when calling a function defined in a plugin.
func CallGo(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParigotSource, ServerGo)
}

// LogFullf creates a new log line in the current log container.  This function allows all
// the varibels to be specified or overriden.  The level parameter can be set to UnknownLL to
// indicate the log level is not known.  The source can be set to UnknownS which indicates
// that the caller doesn't know the source and is willing to accept whatever value (if any)
// that is contained in the context ctx. The funcName can be "" to indicate that the caller
// is willing to accept whatever value (if any) is inside the current context.  The spec and
// rest arguments work like fmt.Printf().
func LogFullf(ctx context.Context, level LogLevel, source Source, funcName, spec string, rest ...interface{}) {
	detailPrefix := detailPrefix(ctx, level, source, funcName)
	line := fmt.Sprintf(detailPrefix+spec, rest...)
	logLine := LogLineFromString(ctx, line, source, level)

	cont := GetContainer(ctx)
	if cont == nil {
		log.Println(line)
		return
	}
	container := cont.(*logContainer)
	container.AddLogLine(ctx, logLine)
}
