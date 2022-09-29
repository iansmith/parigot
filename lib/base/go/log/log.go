package log

import (
	"bytes"
	"fmt"
	"os"
)

//import "go.uber.org/zap"

// type Log zap.Logger

type LogLevel int64

const (
	DebugLevel LogLevel = 1 << 1
	InfoLevel  LogLevel = 1 << 2
	WarnLevel  LogLevel = 1 << 3
	ErrorLevel LogLevel = 1 << 4
	FatalLevel LogLevel = 1 << 5

	DevMask  LogLevel = DebugLevel | InfoLevel | WarnLevel | ErrorLevel | FatalLevel
	ProdMask LogLevel = WarnLevel | ErrorLevel | FatalLevel
)

func (lvl LogLevel) String() string {
	var buf bytes.Buffer
	prefix := ""

	if lvl == 0 {
		outputString("WARN", "uninitialized log level, zero found as the allow mask, assuming no logs printed")
		return "[no logs]"
	}

	for _, s := range []LogLevel{DebugLevel, InfoLevel, WarnLevel, ErrorLevel,
		FatalLevel} {
		switch s {
		case DebugLevel:
			buf.WriteString(prefix + "DebugLevel")
		case InfoLevel:
			buf.WriteString(prefix + "InfoLevel")
		case WarnLevel:
			buf.WriteString(prefix + "WarnLevel")
		case ErrorLevel:
			buf.WriteString(prefix + "ErrorLevel")
		case FatalLevel:
			buf.WriteString(prefix + "FatalLevel")
		}
		prefix = " || "
	}
	return buf.String()
}

type T interface {
	LogMask() LogLevel
	SetLogMask(LogLevel)
	AbortOnFatal()
	SetAbortOnFatal(bool)
	Debug(string, ...string)
	Info(string, ...string)
	Warn(string, ...string)
	Error(string, ...string)
	Fatal(string, ...string)
}

type LocalT struct {
	levelMask    LogLevel
	abortOnFatal bool
}

func (l *LocalT) SetLogMask(mask LogLevel) {
	l.levelMask = mask
}

func (l *LocalT) LogMask() LogLevel {
	return l.levelMask
}

func (l *LocalT) SetAbortOnFatal(f bool) {
	l.abortOnFatal = f
}

func (l *LocalT) AbortOnFatal() bool {
	return l.abortOnFatal
}

func (l *LocalT) Debug(f string, rest ...string) {
	if l.levelMask&DebugLevel != 0 {
		outputString("DEBUG", f, rest...)
	}
}

func (l *LocalT) Info(f string, rest ...string) {
	if l.levelMask&InfoLevel != 0 {
		outputString("INFO ", f, rest...)
	}
}

func (l *LocalT) Warn(f string, rest ...string) {
	if l.levelMask&WarnLevel != 0 {
		outputString("WARN ", f, rest...)
	}
}

func (l *LocalT) Error(f string, rest ...string) {
	if l.levelMask&ErrorLevel != 0 {
		outputString("ERROR", f, rest...)
	}
}

func (l *LocalT) Fatal(f string, rest ...string) {
	if !l.abortOnFatal {
		if l.levelMask&FatalLevel != 0 { // probably a bad idea
			outputString("FATAL", f, rest...)
		}
	} else {
		// if you are aborting, then you probably want to print something
		outputString("FATAL", f, rest...)
		exit(1)
	}
}

func exit(code int) {
	os.Exit(code)
}

func outputString(prefix, f string, rest ...string) {
	fmt.Printf(prefix+":"+f, rest)
}
