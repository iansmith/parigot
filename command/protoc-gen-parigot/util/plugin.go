package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	oneMB      = 1048576
	bufferSize = 2 * oneMB
	readSize   = 1024

	pathArg    = "path"
	trueValue  = "true"
	abiArg     = "abi"
	locatorArg = "locator"
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

func MarshalResponseAndExit(message proto.Message) {
	b, err := proto.Marshal(message)
	if err != nil {
		panic("unable to marshal protobuf response:" + err.Error())
	}
	fmt.Fprintf(os.Stdout, "%s", string(b))
	os.Exit(0) // by the spec, must be zero
}

func IsABIGeneration(param string) bool {
	for key, value := range parametersToMap(param) {
		if key == abiArg && value == trueValue {
			return true
		}
	}
	return false
}

func LocatorNames(param string) []string {
	for key, value := range parametersToMap(param) {
		if key == locatorArg {
			return strings.Split(value, ";")
		}
	}
	return []string{}
}

func parametersToMap(param string) map[string]string {
	result := make(map[string]string)
	param = strings.TrimSpace(param)
	parts := strings.Split(param, ",")
	// there is only one part, we assume it's paths=soureRelative
	if len(parts) >= 1 {
		for _, part := range parts {
			assign := strings.Split(part, "=")
			if len(assign) != 2 {
				log.Fatalf("bad assignment in parameter %s (part of %s), ignoring it",
					part, param)
				continue
			}
			key := strings.ToLower(strings.TrimSpace(assign[0]))
			value := strings.ToLower(strings.TrimSpace(assign[1]))
			result[key] = value
		}
	} else {
		log.Printf("unable to understand parameter '%s', ignoring it", param)
	}
	return result
}
