package log

import (
	"bytes"
	"github.com/iansmith/parigot/g/parigot/abi"
	"time"

	"github.com/iansmith/parigot/g/parigot/kernel"
)

type LogLevel int32

const (
	DebugLevel = 1
	InfoLevel  = 2
	WarnLevel  = 3
	ErrorLevel = 4
	FatalLevel = 5
)

// init functions in parigot must be idempotent! It can and will be called many times.
func init() {
	kernel.Register("parigot", "log")
}

var Dev = NewLocalT(true)

type T interface {
	AbortOnFatal() bool
	SetAbortOnFatal(bool)
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
}

type LocalT struct {
	abortOnFatal bool
}

func (l *LocalT) SetAbortOnFatal(f bool) {
	l.abortOnFatal = f
}

func (l *LocalT) AbortOnFatal() bool {
	return l.abortOnFatal
}

func (l *LocalT) Debug(f string) {
	outputString("DEBUG", f)
}

func (l *LocalT) Info(f string) {
	outputString("INFO ", f)

}

func (l *LocalT) Warn(f string) {
	outputString("WARN ", f)
}

func (l *LocalT) Error(f string) {
	outputString("ERROR", f)
}

func (l *LocalT) Fatal(f string) {
	if !l.abortOnFatal {
		outputString("FATAL", f)
	} else {
		// if you are aborting, then you probably want to print something
		outputString("FATAL", f)
		abi.Exit(1)
	}
}

func NewLocalT(abortOnFatal bool) T {
	return &LocalT{
		abortOnFatal: true,
	}
}

func outputString(prefix, f string, rest ...string) {
	var buf bytes.Buffer
	resp := abi.Now()
	t := time.Unix(0, resp.Now)
	t = t.UTC()
	stamp := t.Format(time.Stamp)
	buf.WriteString(stamp + " UTC")
	buf.WriteString(" ")
	buf.WriteString(prefix)
	buf.WriteString(":")
	buf.WriteString(f)
	last := f
	for _, r := range rest {
		last = r
	}
	addNewlineIfDontHaveOne(last, &buf)
	abi.OutputString(buf.String())
}
func addNewlineIfDontHaveOne(s string, buf *bytes.Buffer) {
	// could do this with converting the buffer to bytes or string but this
	// seems least bad option
	if s[len(s)-1] != '\n' {
		buf.WriteString("\n")
	}
}
