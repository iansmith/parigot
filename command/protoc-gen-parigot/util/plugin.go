package util

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"log"
	"os"
)

const (
	oneMB      = 1048576
	bufferSize = 2 * oneMB
	readSize   = 1024
)

var buffer [bufferSize]byte

func ReadStdinIntoBuffer() *pluginpb.CodeGeneratorRequest {
	curr := 0
	ok := false
	for readSize < bufferSize-curr {
		r, err := os.Stdin.Read(buffer[curr : curr+readSize])
		if err != nil {
			if err == io.EOF {
				curr += r
				ok = true
				break
			}
			log.Fatalf("failed reading stdin: %v", err)
		}
		curr += r
	}
	if !ok {
		log.Fatalf("input message on stdin was larger than %0x bytes", bufferSize)
	}
	var req pluginpb.CodeGeneratorRequest
	err := proto.Unmarshal(buffer[:curr], &req)
	if err != nil {
		log.Fatalf("unable to understand generator request:%v", err)
	}
	return &req
}
