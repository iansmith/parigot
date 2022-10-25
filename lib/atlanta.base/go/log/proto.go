package log

import (
	log2 "github.com/iansmith/parigot/g/parigot/log"
)

func NewProtoLogger(line []*log2.LogRequest) log2.Log {
	return &protoLogger{
		line: line,
	}
}

type protoLogger struct {
	line []*log2.LogRequest
}

func (p *protoLogger) Log(in *log2.LogRequest) error {
	p.line = append(p.line, in)
	return nil
}
