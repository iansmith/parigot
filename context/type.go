package context

import (
	"context"
	"fmt"
	"log"
	"time"
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
)

type LogLevel int

const (
	UnknownLL LogLevel = 0
	Debug     LogLevel = 1
	Info      LogLevel = 2
	Warn      LogLevel = 3
	Error     LogLevel = 4
	Fatal     LogLevel = 5

	// stackTraceInternal is really for internal use only.
	stackTraceInternal LogLevel = 6
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
	case stackTraceInternal:
		return " STR"
	}
	return "FATL"
}

func (l LogLevel) Integer() int {
	return int(l)
}

type Source int

const (
	UnknownS   Source = 0
	Client     Source = 1
	ServerGo   Source = 2
	ServerWasm Source = 3
	Parigot    Source = 4
	Wazero     Source = 5
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
	}
	return "Parigot"
}

func (s Source) Integer() int {
	return int(s)
}

func detailPrefix(ctx context.Context, level LogLevel, source Source, funcName string) string {
	tString := CurrentTimeString(ctx)
	lString := level.String()
	sString := source.String()
	if source == UnknownS {
		possibleS := ctx.Value(ParigotSource)
		if possibleS != nil {
			sString = possibleS.(Source).String()
		}
	}
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
}

type LogContainer interface {
	StackTrace(ctx context.Context, detailPrefix, header, footer, funcName string)
	AddLogLine(ctx context.Context, l LogLine)
	Dump()
}

const stackBufferSize = 4096

func StackTrace(ctx context.Context, funcName string) {
	detail := fmt.Sprintf("StackTrace (%s) ", funcName)
	header := detail + ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
	footer := "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	b := make([]byte, stackBufferSize) // adjust buffer size to be larger than expected stack
	s := string(b[:stackBufferSize])
	cont := GetContainer(ctx)
	if cont == nil {
		log.Println(header)
		log.Print(s)
		log.Println(footer)
		return
	} else {
		cont.StackTrace(ctx, detail, header, footer, s)
	}
}

func newContext(orig context.Context, src Source, name string) context.Context {
	if orig == nil {
		orig = context.Background()
		Errorf(orig, "the use of nil context to newContext() is discouraged")
		StackTrace(orig, "newContext")
	}
	cont := newLogContainer()
	sanity := LogContainer(cont)
	ctx := context.WithValue(orig, ParigotTime, time.Now())
	ctx = context.WithValue(ctx, ParigotFunc, name)
	ctx = context.WithValue(ctx, ParigotSource, src)
	ctx = context.WithValue(ctx, ParigotLogContainer, sanity)
	return ctx
}
