package context

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"
)

const MaxLineLen = 512
const MaxContainerSize = 256

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
)

func (s Source) String() string {
	switch s {
	case UnknownS:
		return "------"
	case Client:
		return "Client"
	case ServerGo:
		return "ServGo"
	case ServerWasm:
		return "SvWasm"
	}
	return "Prigot"
}

func (s Source) Integer() int {
	return int(s)
}

type LogLine struct {
	data [MaxLineLen]byte // c-style terminator (nul byte)
}

type LogContainer struct {
	front, back int
	line        [MaxContainerSize]*LogLine
}

func LogFullf(ctx context.Context, level LogLevel, source Source, funcName, spec string, rest ...interface{}) {
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

	detailSpec := fmt.Sprintf("%s:%s:%s:%-32s:%s", tString, lString, sString, funcName, spec)
	line := fmt.Sprintf(detailSpec, rest...)

	maxWithZero := MaxLineLen - 1
	if len(line) >= maxWithZero {
		start := len(line) - maxWithZero
		line = line[:start]
	}
	i := 0

	result := &LogLine{}
	for i < len(line) {
		result.data[i] = line[i]
		i++
	}
	result.data[i] = 0

	cont := ctx.Value(ParigotLogContainer)
	if cont == nil {
		log.Println(line)
	}

	container := cont.(*LogContainer)
	container.line[container.front] = result
	container.front = (container.front + 1) % MaxContainerSize
	container.front %= MaxContainerSize
	if container.front == container.back {
		container.back = (container.back + 1) % MaxContainerSize
	}
}

func Logf(ctx context.Context, level LogLevel, spec string, rest ...interface{}) {
	LogFullf(ctx, level, UnknownS, "", spec, rest...)
}
func Errorf(ctx context.Context, spec string, rest ...interface{}) {
	LogFullf(ctx, Error, UnknownS, "", spec, rest...)
}
func Debugf(ctx context.Context, funcName string, spec string, rest ...interface{}) {
	LogFullf(ctx, Debug, UnknownS, funcName, spec, rest...)
}

// ClientLogf is just like Logf except is sets the source to be client.  This is
// useful (with no context param) because client's usually don't have a context.
func ClientLogf(level LogLevel, spec string, rest ...interface{}) {
	LogFullf(context.Background(), level, Client, "", spec, rest...)
}

// ClientLogf is just like Debugf except is sets the source to be client.  This is
// useful (with no context param) because client's usually don't have a context.
func ClientDebugf(funcName string, spec string, rest ...interface{}) {
	LogFullf(context.Background(), Debug, UnknownS, funcName, spec, rest...)
}

// LogInternal is for internal use only.  It creates a log line attributed
// to Parigot.
func LogInternal(level LogLevel, funcName, spec string, rest ...interface{}) {
	LogFullf(context.Background(), level, Parigot, funcName, spec, rest...)
}

func dumpContainer(cont *LogContainer) {
	i := cont.front
	buf := &bytes.Buffer{}
	for i != cont.back {
		// put this line's data in the buffer
		l := cont.line[i]
		lastIsCR := false
		for k := 0; k < MaxLineLen; k++ {
			if l.data[k] == 0 {
				break
			}
			buf.WriteByte(l.data[k])
			if l.data[k] == 10 {
				lastIsCR = true
			} else {
				lastIsCR = false
			}
		}
		if !lastIsCR {
			buf.WriteString("\n")
		}
		i = (i + 1) % MaxContainerSize
	}
	log.Print(buf.String())

}

func CurrentTime(ctx context.Context) time.Time {
	t := ctx.Value(ParigotTime)
	if t != nil && !t.(time.Time).IsZero() {
		return t.(time.Time)
	}
	return time.Now()
}

func CurrentTimeString(ctx context.Context) string {
	return CurrentTime(ctx).Format(time.RFC822Z)
}

func CallTo(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, ParigotFunc, s)
}
func CallGo(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParigotSource, ServerGo)
}

func newContext(orig context.Context, src Source, name string) context.Context {
	if orig == nil {
		orig = context.Background()
	}
	cont := &LogContainer{}
	ctx := context.WithValue(orig, ParigotTime, time.Now())
	ctx = context.WithValue(ctx, ParigotFunc, name)
	ctx = context.WithValue(ctx, ParigotSource, src)
	ctx = context.WithValue(ctx, ParigotLogContainer, cont)
	return ctx
}

func ServerGoContext(ctx context.Context, funcName string) context.Context {
	return newContext(ctx, ServerGo, funcName)
}
func ClientContext(ctx context.Context, funcName string) context.Context {
	return newContext(ctx, Client, funcName)
}
func ServerWasmContext(ctx context.Context, funcName string) context.Context {
	return newContext(ctx, ServerWasm, funcName)
}

func Dump(ctx context.Context) {
	cont := ctx.Value(ParigotLogContainer)
	if cont == nil {
		log.Println("no log container present inside context")
		return
	}
	dumpContainer(cont.(*LogContainer))
}
