package lib

import (
	log2 "github.com/iansmith/parigot/g/parigot/log"
	"github.com/iansmith/parigot/lib/interface_"
)

func NewProtoLogger(line []*log2.LogRequest) interface_.Log {
	return &protoLogger{
		line:         line,
		abortOnFatal: true,
	}
}

type protoLogger struct {
	line         []*log2.LogRequest
	abortOnFatal bool
}

func (p *protoLogger) AbortOnFatal() bool {
	return p.abortOnFatal
}
func (p *protoLogger) SetAbortOnFatal(a bool) {
	p.abortOnFatal = a
}

func (p *protoLogger) Log(prefix string, level interface_.LogLevel, msg string) {
	in := log2.LogRequest{
		Stamp:   nil,
		Level:   int32(level),
		Message: prefix + ":" + msg,
	}
	p.line = append(p.line, &in)
}
