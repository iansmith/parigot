package lib

import (
	pb "github.com/iansmith/parigot/g/pb/log"
)

func NewProtoLogger(line []*pb.LogRequest) Log {
	return &protoLogger{
		line:         line,
		abortOnFatal: true,
	}
}

type protoLogger struct {
	line         []*pb.LogRequest
	abortOnFatal bool
}

func (p *protoLogger) AbortOnFatal() bool {
	return p.abortOnFatal
}
func (p *protoLogger) SetAbortOnFatal(a bool) {
	p.abortOnFatal = a
}

func (p *protoLogger) Log(prefix string, level LogLevel, msg string) {
	in := pb.LogRequest{
		Stamp:   nil,
		Level:   int32(level),
		Message: prefix + ":" + msg,
	}
	p.line = append(p.line, &in)
}
