//go:build !js
// +build !js

package ui

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/iansmith/parigot/api/netconst"
	pb "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/sys/jspatch"

	"google.golang.org/protobuf/proto"
)

type LogViewerImpl struct {
	mem        *jspatch.WasmMem
	path       string // when this is "" we have not yet tried the UD socket
	connection net.Conn
	terminal   bool
}

var decodeError = errors.New("decoding error")

// This is the native code side of the logviewer.  It reads the payload from the WASM world and either
// dumps it to the terminal or sends it through the UD socket to the GUI.

//go:noinline
func (l *LogViewerImpl) LogRequestViaSocket(sp int32) {
	wasmPtr := l.mem.GetInt64(sp + 8)

	buffer := l.ReadSlice(wasmPtr, unsafe.Offsetof(LogViewerPayload{}.Ptr),
		unsafe.Offsetof(LogViewerPayload{}.Len))

	if l.path == "" {
		dir := "/var/run/parigot" // default value
		if os.Getenv(netconst.SocketEnvVar) != "" {
			dir = os.Getenv(netconst.SocketEnvVar)
		}
		path := filepath.Join(dir, netconst.SocketName)
		l.path = path
		_, err := os.Stat(path)
		if err != nil {
			log.Printf("unable to open %s for connection to log viewer, defaulting to terminal output", path)
			l.terminal = true
		} else {
			conn, err := net.Dial("unix", path)
			if err != nil {
				log.Printf("unable to connect to logViewer on  %s, defaulting to terminal output", path)
				l.terminal = true
			} else {
				l.connection = conn
				l.terminal = false
			}
		}
	}
	if l.terminal {
		req, err := DecodeLogRequestBuffer(buffer)
		if err != nil {
			return
		}
		n := ""
		if !strings.HasSuffix(req.Message, "\n") {
			n = "\n"
		}
		fmt.Printf("%s:%s:%s%s", req.Stamp.AsTime().Format(time.RFC3339), intToLogLevel(int(req.Level)),
			req.Message, n)
		return
	} else {
		written := 0
		for written < len(buffer) {
			w, err := l.connection.Write(buffer[written:])
			if err != nil {
				log.Printf("warning: error writing to the log viewer program over socket %s: %v", l.path, err)
				return
			}
			written += w
		}
	}
}

func DecodeLogRequestBuffer(buffer []byte) (*pb.LogRequest, error) {
	m := binary.LittleEndian.Uint64(buffer[0:8])
	if m != netconst.MagicStringOfBytes {
		log.Printf("unable to print log message, bad magic number %x", m)
		return nil, decodeError
	}
	l := binary.LittleEndian.Uint32(buffer[8:12])
	if l >= uint32(netconst.ReadBufferSize) {
		log.Printf("unable to print log message, very large log message [%d bytes]", l)
		return nil, decodeError
	}
	size := int(l)
	req := pb.LogRequest{}
	if err := proto.Unmarshal(buffer[netconst.FrontMatterSize:netconst.FrontMatterSize+size], &req); err != nil {
		for i := 0; i < size; i += 8 {
			print("bytes: ", i, " ", buffer[netconst.FrontMatterSize+i+0], " ", buffer[netconst.FrontMatterSize+i+1], " ", buffer[netconst.FrontMatterSize+i+2], " ", buffer[netconst.FrontMatterSize+i+3], " ", buffer[netconst.FrontMatterSize+i+4], " ", buffer[netconst.FrontMatterSize+i+5], " ", buffer[netconst.FrontMatterSize+i+6], " ", buffer[netconst.FrontMatterSize+i+7], "\n")
		}
		print(buffer[netconst.FrontMatterSize:netconst.FrontMatterSize+size], "\n")
		log.Printf("unable to print log message, request could not be unmarshaled: %v", err)
		return nil, decodeError
	}
	result := crc32.Checksum(buffer[:size+netconst.FrontMatterSize], netconst.KoopmanTable)
	expected := binary.LittleEndian.Uint32(buffer[netconst.FrontMatterSize+size:])
	if expected != result {
		log.Printf("unable to print log message, bad checksum found on log request")
		return nil, decodeError
	}
	return &req, nil

}

func intToLogLevel(i int) string {
	switch {
	case pb.LogLevel(i) == pb.LogLevel_LOGLEVEL_UNSPECIFIED:
		return "UNKNOWN"
	case pb.LogLevel(i) == pb.LogLevel_LOGLEVEL_DEBUG:
		return "DEBUG"
	case pb.LogLevel(i) == pb.LogLevel_LOGLEVEL_INFO:
		return "INFO "
	case pb.LogLevel(i) == pb.LogLevel_LOGLEVEL_WARNING:
		return "WARN "
	case pb.LogLevel(i) == pb.LogLevel_LOGLEVEL_ERROR:
		return "ERROR"
	case pb.LogLevel(i) == pb.LogLevel_LOGLEVEL_FATAL:
		return "FATAL"
	default:
		return fmt.Sprintf("UNEXPECTED[%d]", i)
	}
}

func (s *LogViewerImpl) ReadSlice(structPtr int64, dataOffset uintptr, lenOffset uintptr) []byte {
	return s.mem.LoadSliceWithLenAddr(int32(structPtr)+int32(dataOffset),
		int32(structPtr)+int32(lenOffset))
}

func (l *LogViewerImpl) SetWasmMem(ptr uintptr) {
	l.mem = jspatch.NewWasmMem(ptr)
}
