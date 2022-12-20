//go:build !js
// +build !js

package go_

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"unsafe"

	pb "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/sys/jspatch"
)

type LogViewerImpl struct {
	mem        *jspatch.WasmMem
	path       string // when this is "" we have not yet tried the UD socket
	connection net.Conn
	terminal   bool
}

// This is the native code side of the logviewer.  It reads the payload from the WASM world and either
// dumps it to the terminal or sends it through the UD socket to the GUI.

//go:noinline
func (l *LogViewerImpl) LogRequestViaSocket(sp int32) {
	wasmPtr := l.mem.GetInt64(sp + 8)

	buffer := l.ReadSlice(wasmPtr, unsafe.Offsetof(splitutil.SplitUtilSinglePayload{}.Ptr),
		unsafe.Offsetof(splitutil.SplitUtilSinglePayload{}.Len))

	if l.path == "" {
		l.path = "host.docker.internal:4004"
		conn, err := net.Dial("tcp", "host.docker.internal:4004")
		if err != nil {
			log.Printf("unable to connect to logViewer on  %s, defaulting to terminal output: %v", "host.docker.internal:4004", err)
			l.terminal = true
		} else {
			l.connection = conn
			l.terminal = false
		}
	}
	if l.terminal {
		req, err := splitutil.DecodeProto[pb.LogRequest](buffer)
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
