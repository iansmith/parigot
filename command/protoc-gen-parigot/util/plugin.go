package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	oneMB      = 1048576
	bufferSize = 2 * oneMB
	readSize   = 4096

	pathArg    = "path"
	trueValue  = "true"
	abiArg     = "abi"
	locatorArg = "locator"
)

var buffer [bufferSize]byte

func ReadStdinIntoBuffer(reader io.Reader, saveTemp bool) *pluginpb.CodeGeneratorRequest {
	curr := 0
	ok := false

	for readSize < bufferSize-curr {
		r, err := reader.Read(buffer[curr : curr+readSize])
		if err != nil {
			if err == io.EOF {
				curr += r
				ok = true
				break
			}
		}
		curr += r
	}
	if !ok {
		log.Fatalf("input message on stdin was larger than %0x bytes (%0x)", bufferSize, curr)
	}
	var req pluginpb.CodeGeneratorRequest
	err := proto.Unmarshal(buffer[:curr], &req)
	if err != nil {
		log.Fatalf("unable to understand generator request:%v", err)
	}
	if saveTemp {
		dir, err := os.MkdirTemp("/tmp/", "parse")
		if err != nil {
			log.Fatalf("%v", err)
		}
		protoGuess := req.GetProtoFile()[len(req.GetProtoFile())-1].GetName()
		path := fmt.Sprintf("%s", protoGuess)
		out := filepath.Join(dir, path)
		parentDir, _ := filepath.Split(out)
		if err := os.MkdirAll(parentDir, 0700); err != nil {
			log.Fatalf("%v", err)
		}
		outfp, err := os.Create(out)
		if err != nil {
			log.Fatalf("%v", err)
		}
		count, err := io.Copy(outfp, bytes.NewBuffer(buffer[:curr]))
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Printf("saved copy of input to %s (0x%0x bytes) from %s", out, count,
			protoGuess)
		outfp.Close()
	}
	return &req
}

func OutputTerminal(output []*OutputFile) {
	for _, file := range output {
		fmt.Fprintf(os.Stdout, "*** %s\n", file.name)
		_, err := io.Copy(os.Stdout, &file.buf)
		if err != nil {
			log.Fatalf("unable to copy to output:%v", err)
		}
	}
}

func MarshalResponseAndExit(message proto.Message) {
	b, err := proto.Marshal(message)
	if err != nil {
		panic("unable to marshal protobuf response:" + err.Error())
	}
	fmt.Fprintf(os.Stdout, "%s", string(b))
	os.Exit(0) // by the spec, must be zero
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
