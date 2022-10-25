package lib

type LogLevel int32

const (
	DebugLevel = 1
	InfoLevel  = 2
	WarnLevel  = 3
	ErrorLevel = 4
	FatalLevel = 5
)

type Log interface {
	AbortOnFatal() bool
	SetAbortOnFatal(bool)
	Log(prefix string, level LogLevel, msg string)
}
