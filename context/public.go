package context

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	_ "time/tzdata"
)

// XXX Should be passed in the config or someting like that
var localTimeZone *time.Location

func init() {
	var err error
	localTimeZone, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Printf("unable to load time zone: %s", err.Error())
	}
}

// ServerGoContext returns a new context based on ctx with source ServerGo and the given function name.
// This should be called before entering plugins that define host functions.
func ServerGoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParigotSource, HostGo)
}

// GuestContext returns a new context based on ctx with source Guest and the given function name.
// This should be called before entering code that is client wasm code, like the main
// of an application or the start of test.
func GuestContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParigotSource, Guest)
}

// WazeroContext returns a context that has the source sent to Wazero.
func WazeroContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParigotSource, Wazero)
}

// CurrentTime should *always* be used in both host and guest code to get the current
// time.  This should be done because it allows testing code that needs to run at a
// particular time (hourly, midnight, etc) to be tested more easily.
func CurrentTime(ctx context.Context) time.Time {
	t := ctx.Value(ParigotTime)
	if t != nil && !t.(time.Time).IsZero() {
		if localTimeZone != nil {
			return t.(time.Time).In(localTimeZone)
		}
		return t.(time.Time)
	}
	if localTimeZone != nil {
		return time.Now().In(localTimeZone)
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
	cont.Dump(ctx)
}

// Logf is a shorthand for a call to LogFull with the currently set source in ctx.
func Logf(ctx context.Context, level LogLevel, spec string, rest ...interface{}) {
	LogFullf(ctx, level, PullSource(ctx, UnknownS), pullFunc(ctx, ""),
		spec, rest...)
}

// Raw creates a new log line with no prefix other than the source.
// This should be used when dumping lines of text to the logger
// such as stack traces, lines from stdout, etc.  Most users
// will never need this function.
func Raw(ctx context.Context, src Source, line string) {
	addToContainerOrPrint(ctx, NewLogLine(ctx, PullSource(ctx, src), Debug, pullFunc(ctx, ""), true, line))
}

// Fatalf is a shorthand for a call to LogFull with the currently set source in ctx.
func Fatalf(ctx context.Context, spec string, rest ...interface{}) {
	LogFullf(ctx, Fatal, PullSource(ctx, UnknownS), pullFunc(ctx, ""),
		spec, rest...)
}

// Errorf is a shorthand for a call to LogFull with the currently set source in ctx, and
// log level being Error.
func Errorf(ctx context.Context, spec string, rest ...interface{}) {
	LogFullf(ctx, Error, PullSource(ctx, UnknownS), pullFunc(ctx, ""),
		spec, rest...)
}

// Warnf is a shorthand for a call to LogFull with the currently set source in ctx, and
// log level being Warn.
func Warnf(ctx context.Context, spec string, rest ...interface{}) {
	LogFullf(ctx, Warn, PullSource(ctx, UnknownS), pullFunc(ctx, ""),
		spec, rest...)
}

// Debugf is a shorthand for a call to LogFull with the currently set source in ctx, and
// log level being Debug.  The function name is the one associated with the given ctx.
func Debugf(ctx context.Context, spec string, rest ...interface{}) {
	LogFullf(ctx, Debug, PullSource(ctx, UnknownS),
		pullFunc(ctx, ""), spec, rest...)
}

// Infof is a shorthand for a call to LogFull with the currently set source in ctx, and
// log level being Infof.  The function name is the one associated with the given ctx.
func Infof(ctx context.Context, spec string, rest ...interface{}) {
	LogFullf(ctx, Info, PullSource(ctx, UnknownS),
		pullFunc(ctx, ""), spec, rest...)
}

// DebugFuncf is a shorthand for a call to LogFull with the currently set source in ctx, and
// log level being Debug.  The function name must be supplied here and overrides any function
// name found in the ctx.
func DebugFuncf(ctx context.Context, funcName string, spec string, rest ...interface{}) {
	LogFullf(ctx, Debug, UnknownS, funcName, spec, rest...)
}

// InternalParigot is for internal use only.  It creates a log line attributed to Parigot.
func InternalParigot(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParigotSource, Source(Parigot))
}

// CurrentTimeString is a wrapper around CurrentTime that returns a string representation
// of the current time in a standard form (RFC822Z).
func CurrentTimeString(ctx context.Context, rfc822 bool) string {
	if rfc822 {
		return CurrentTime(ctx).Format(time.RFC822)
	}
	return CurrentTime(ctx).Format(time.Kitchen)
}

// CallTo returns a new context with the current function updated to s.
func CallTo(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, ParigotFunc, s)
}

// LogFullf creates a new log line in the current log container.  This function allows all
// the varibels to be specified or overriden.  The level parameter can be set to UnknownLL to
// indicate the log level is not known.  The source can be set to UnknownS which indicates
// that the caller doesn't know the source and is willing to accept whatever value (if any)
// that is contained in the context ctx. The funcName can be "" to indicate that the caller
// is willing to accept whatever value (if any) is inside the current context.  The spec and
// rest arguments work like fmt.Printf().
func LogFullf(ctx context.Context, level LogLevel, source Source,
	funcName string, spec string, rest ...interface{}) {
	logLine := NewLogLine(ctx, source, level, funcName,
		false, spec, rest...)

	addToContainerOrPrint(ctx, logLine)
}

func addToContainerOrPrint(ctx context.Context, line LogLine) {
	cont := GetContainer(ctx)
	if line.IsDevNull() {
		return
	}
	if cont == nil {
		line.Print()
		return
	}
	container := cont.(*logContainer)
	container.AddLogLine(ctx, line)
}

// SourceContext lets you build a context that contains any source value.  Most user
// code will be better off using GoWasmContext() or similar. This is only
// useful when the source is a variable.
func SourceContext(ctx context.Context, source Source) context.Context {
	return context.WithValue(ctx, ParigotSource, source)
}

func pullLineAndFile(ctx context.Context) string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return "(unknown)"
	}
	part := strings.Split(file, "/")
	if len(part) > 1 {
		part = part[len(part)-2:]
	}
	file = filepath.Join(part...)

	return fmt.Sprintf("(%s:%d)", file, line)
}
