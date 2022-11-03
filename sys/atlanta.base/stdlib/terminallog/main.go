package main

import (
	"fmt"

	pblog "github.com/iansmith/parigot/g/pb/log"
	"google.golang.org/protobuf/proto"
)

func init() {
}

type terminalLog struct {
	prefix string
}

func main() {
	run(&terminalLog{})
}

func newTerminalLogImpl() *terminalLog {
	return &terminalLog{}
}

func (l *terminalLog) Log(request proto.Message) error {
	req := request.(*pblog.LogRequest)

	fmt.Printf("%s%s: %s", l.prefix, req.GetLevel().Descriptor().Name(), req.GetMessage())
	msgLen := len(req.GetMessage())
	if msgLen == 0 {
		fmt.Printf("\n")
	} else {
		last := req.GetMessage()[msgLen-1:]
		if last != "\n" {
			fmt.Printf("\n")
		}
	}
	return nil
}

func (l *terminalLog) SetPrefix(request proto.Message) error {
	req := request.(*pblog.SetPrefixRequest)
	p := req.GetPrefix()
	prefixLen := len(p)
	if prefixLen > 0 {
		if req.GetPrefix()[prefixLen-1:] != " " && req.GetPrefix()[prefixLen-1:] != ":" {
			p += ":"
		}
	}
	l.prefix = p
	return nil
}
