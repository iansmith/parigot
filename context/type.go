package context

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const MaxLineLen = 512
const MaxContainerSize = 256

// size of a container in bytes is always MaxLineLen * MaxContainerSize (512 * 256 => 20k)
type ParigotKey string

const (
	ParigotTime         ParigotKey = "parigot_time"
	ParigotFunc         ParigotKey = "parigot_func"
	ParigotSource       ParigotKey = "parigot_source"
	ParigotLogContainer ParigotKey = "parigot_log_container"
	ParigotOrigin       ParigotKey = "parigot_origin"
)

type LogLevel int

const (
	UnknownLL LogLevel = 0
	Debug     LogLevel = 1
	Info      LogLevel = 2
	Warn      LogLevel = 3
	Error     LogLevel = 4
	Fatal     LogLevel = 5
)

func (l LogLevel) String() string {
	switch l {
	case UnknownLL:
		return "----"
	case Debug:
		return "DEBG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return " ERR"
	case Fatal:
		return "FATL"
	}
	panic(fmt.Sprintf("unable to understand log level %d", int(l)))
}

func (l LogLevel) Integer() int {
	return int(l)
}

type Source int

const (
	UnknownS           Source = 0
	Client             Source = 1
	ServerGo           Source = 2
	ServerWasm         Source = 3
	Parigot            Source = 4
	Wazero             Source = 5
	WasiOut            Source = 6
	WasiErr            Source = 7
	StackTraceInternal Source = 8
	SpewInternal       Source = 9
)

func (s Source) String() string {
	switch s {
	case UnknownS:
		return "-------"
	case Client:
		return " Client"
	case ServerGo:
		return "  SrvGo"
	case ServerWasm:
		return "SvrWasm"
	case Wazero:
		return " Wazero"
	case WasiOut:
		return "WasiOut"
	case WasiErr:
		return "WasiErr"
	case Parigot:
		return "Parigot"
	case StackTraceInternal:
		return "StackTr"
	case SpewInternal:
		return "   Spew"
	}
	panic(fmt.Sprintf("unknown source value %d", int(s)))
}

func (s Source) Integer() int {
	return int(s)
}

func detailPrefix(ctx context.Context, level LogLevel, source Source, fn string) string {
	var tString string
	rfc822 := true
	// if source == UnknownS {
	// 	source = PullSource(ctx, source)
	// }
	// if source == ServerWasm || source == Wazero || source == Client {
	// 	rfc822 = false
	// }
	if runtime.GOOS == "wasip1" {
		rfc822 = false
	}
	tString = CurrentTimeString(ctx, rfc822)
	lString := level.String()
	sString := source.String()
	funcName := pullFunc(ctx, fn)
	if funcName == "" {
		f := ctx.Value(ParigotFunc)
		if f == nil {
			funcName = "[-unknown-]"
		} else {
			funcName = f.(string)
		}
	}

	return fmt.Sprintf("%s:%s:%s:%-16s", tString, lString, sString, funcName)
}

type LogLine interface {
	Print(context.Context)
}

type LogContainer interface {
	StackTrace(ctx context.Context)
	AddLogLine(ctx context.Context, l LogLine)
	Dump(ctx context.Context)
}

const stackBufferSize = 4096

func StackTrace(ctx context.Context) {
	cont := GetContainer(ctx)
	if cont == nil {
		debug.PrintStack()
	} else {
		cont.StackTrace(ctx)
	}
}
func Spew(ctx context.Context, variable ...interface{}) {
	cont := GetContainer(ctx)
	if cont == nil {
		spew.Dump(variable...)
	} else {
		s := spew.Sdump(variable)
		Raw(ctx, SpewInternal, s)
	}
}

func NewContextWithContainer(orig context.Context, origin string) context.Context {
	if orig == nil {
		orig = context.Background()
		Errorf(orig, "the use of nil context to newContext() is discouraged")
		StackTrace(orig)
	}
	cont := newLogContainer(origin)
	sanity := LogContainer(cont)
	ctx := context.WithValue(orig, ParigotTime, time.Now())
	ctx = context.WithValue(ctx, ParigotLogContainer, sanity)
	ctx = context.WithValue(ctx, ParigotOrigin, origin)
	return ctx
}

func PullSource(ctx context.Context, src Source) Source {
	if src == UnknownS {
		possibleS := ctx.Value(ParigotSource)
		if possibleS != nil {
			return possibleS.(Source)
		}
	}
	return src
}

func pullFunc(ctx context.Context, fn string) string {
	if fn == "" {
		possibleFn := ctx.Value(ParigotFunc)
		if possibleFn != nil {
			return possibleFn.(string)
		}
	}
	return fn
}

func PullOrigin(ctx context.Context) string {
	o := ctx.Value(ParigotOrigin)
	if o == nil {
		return "!!BAD ORIGIN!!"
	}
	return o.(string)
}
