//go:build !js
// +build !js

package go_

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	pb "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/sys/jspatch"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LogViewerImpl struct {
	mem *jspatch.WasmMem
}

const dialPathToLogViewer = "host.docker.internal:4004"

func init() {
	go channelProcessor(dialPathToLogViewer)
}

// This is the native code side of the logviewer.  It reads the payload from the WASM world and either
// dumps it to the terminal or sends it through the UD socket to the GUI.

//go:noinline
func (l *LogViewerImpl) LogRequestHandler(sp int32) {
	req := pb.LogRequest{}
	// xxxfixme: StackPointerToRequest really should return (or provide access to) the buffer used
	// xxxfixme: to read the request because we end up regenerating it when we send this to the
	// xxxfixme: logviewer.
	print(fmt.Sprintf("reached LogRequestHandler about to convert stack pointer\n"))
	err := splitutil.StackPointerToRequest(l.mem, sp, &req)
	if err != nil {
		return // already set the error code
	}
	print(fmt.Sprintf("reached LogRequestHandler about to process request %#v, %s", req, req.Message))
	ProcessLogRequest(&req, false, false, nil)
	splitutil.RespondEmpty(l.mem, sp)
	return
}

func intToLogLevel(i int) string {
	switch {
	case pb.LogLevel(i) == pb.LogLevel_LOG_LEVEL_UNSPECIFIED:
		return "UNKNOWN"
	case pb.LogLevel(i) == pb.LogLevel_LOG_LEVEL_DEBUG:
		return "DEBUG"
	case pb.LogLevel(i) == pb.LogLevel_LOG_LEVEL_INFO:
		return "INFO "
	case pb.LogLevel(i) == pb.LogLevel_LOG_LEVEL_WARNING:
		return "WARN "
	case pb.LogLevel(i) == pb.LogLevel_LOG_LEVEL_ERROR:
		return "ERROR"
	case pb.LogLevel(i) == pb.LogLevel_LOG_LEVEL_FATAL:
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

// logTuple is what is actually passed through logChannel to the implementation portion of the log
// machinery.
type logTuple struct {
	buffer    []byte
	req       *pb.LogRequest
	isKernel  bool
	isBackend bool
}

// logChannel is here to allow the LogRequests to be processed serially and without locks.
var logChannel = make(chan *logTuple, 32)

// ProcessLogRequest can come from two sources, 1) the "normal" path coming from some WASM program.  That goes
// through the "trap" style interface and is handled by LogRequestHandler who calls this function.  2) The
// other path is from some other part of the *go* infrastructure.  Note that this is not referring to the
// "normal" server side of a user program, but rather the _implementation_ of some system function that
// is defined in go.  This function can be called by both paths simultaneously thus a channel is used to serialize.
// The isKernel and isBackend should be set to true only if this is called by some part of the kernel itself
// or some "backend" implementation of a function, respectively.  If the caller does not have the already serialized
// version of req, buffer can be passed as nil and this function will create the buffer itself.
func ProcessLogRequest(req *pb.LogRequest, isKernel, isBackend bool, buffer []byte) {
	tuple := &logTuple{buffer, req, isKernel, isBackend}
	logChannel <- tuple
}

// channelProcessor is started by the init() function on its own goroutine. Its job is to take each
// message out of the logChannel (serially) and send it to the logSingleMessage function.
func channelProcessor(dialPath string) {
	ls := newLogState(dialPath)
	for tuple := range logChannel {
		ls.logSingleMessage(tuple)
	}
}

// logState is used for storing information about the way we are currently logging such as the connection
// to the log UI, if there is one.  It is used by by logSingleMessage().
type logState struct {
	connection net.Conn
	terminal   bool
	path       string
}

// newLogState creates a new logState object and checks to see if we can connect with the given dialPath.  It
// sets the terminal flag to true if we cannot connect to the given dial path.
func newLogState(dialPath string) *logState {
	ls := &logState{}
	conn, err := net.Dial("tcp", dialPath)
	if err != nil {
		log.Printf("unable to connect to logViewer on  %s, defaulting to terminal output: %v", dialPath, err)
		ls.terminal = true
	} else {
		ls.connection = conn
		ls.path = dialPath
		ls.terminal = false
	}
	return ls
}

// logSingleMessage should *only* be called by the channelProcessor and on the channel processors goroutine.
// It either sends the data to the logviewer (if it can be connected to) or the the terminal
func (l *logState) logSingleMessage(tuple *logTuple) {
	if tuple.req == nil && tuple.buffer == nil {
		log.Printf("badly formed logTuple, cannot be processed with the log service")
		return
	}
	if tuple.req == nil {
		tmp := pb.LogRequest{}
		err := splitutil.DecodeSingleProto(tuple.buffer, &tmp)
		if err != nil {
			tmp.Stamp = timestamppb.New(time.Now())
			tmp.Level = pb.LogLevel_LOG_LEVEL_ERROR
			tmp.Message = fmt.Sprintf("unable to extract LogRequest from data buffer: %v", err)
			tuple.isKernel = false
			tuple.isBackend = true
		}
		tuple.req = &tmp
	}
	if l.connection == nil && !l.terminal {
		log.Printf("badly formed logState, not attempting to process log message")
		return
	}
	if l.terminal {
		suffix := ""
		if !strings.HasSuffix(tuple.req.Message, "\n") {
			suffix = "\n"
		}
		prefix := ""
		if tuple.isBackend {
			prefix = ">>"
		}
		if tuple.isKernel {
			prefix = "**"
		}
		fmt.Printf("%s%s:%s:%s%s",
			prefix,
			tuple.req.Stamp.AsTime().Format(time.RFC3339),
			intToLogLevel(int(tuple.req.Level)),
			tuple.req.Message,
			suffix)
		return
	}
	// we need to send the GUI log viewer in the serialized format
	if tuple.buffer == nil {
		// we may need to flatten the req if the caller did not provide the bytes
		buff, err := proto.Marshal(tuple.req)
		if err != nil {
			log.Printf("unable to marshal log message for log viewer: %v", err)
			return
		}
		tuple.buffer = buff
	}
	// write to the connection
	written := 0
	for written < len(tuple.buffer) {
		w, err := l.connection.Write(tuple.buffer[written:])
		if err != nil {
			log.Printf("warning: error writing message (%s) to the log viewer program over socket %s: %v", tuple.req.Message, l.path, err)
			log.Printf("switching to terminal output for subsequent log messages")
			l.connection = nil
			l.terminal = true
			l.path = ""
			return
		}
		written += w
	}
}
