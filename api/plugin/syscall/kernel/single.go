package kernel

import (
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

const timeout = 50

type SingleApproach struct {
	init   bool
	chRecv chan *syscall.ReadOneResponse
	chFin  chan *syscall.ReadOneResponse
}

func NewSingleApproach() *SingleApproach {
	return &SingleApproach{}
}

// maxConcRecv is the number of concurret attempts to send us
// a method call or a resolution at a time
const maxConcRecv = 64

// Init() does initialization of this object and returns if it succeeds
// or not.  Note that in this case of Single this method will be called
// three times, once for each Sender, Receiver, Finisher
func (s *SingleApproach) Init() bool {
	if s.init {
		return true
	}
	// no need to worry about blocking since we are ultimately
	// writing to the chRecv and if that blocks we are hosed in
	// great many ways

	s.chRecv = make(chan *syscall.ReadOneResponse, maxConcRecv)
	s.chFin = make(chan *syscall.ReadOneResponse, maxConcRecv)

	s.init = true
	return true
}

func (s *SingleApproach) Sender() syscall.KernelErr {
	return syscall.KernelErr_NoError
}

func (s *SingleApproach) TimeoutInMillis() int {
	return timeout
}

func (s *SingleApproach) Cancel(k Kernel) {
	k.CancelRead()
}

func (s *SingleApproach) HandleDispatch(req *syscall.DispatchRequest) syscall.KernelErr {
	resp := &syscall.ReadOneResponse{}
	resp.Bundle = &syscall.MethodBundle{}
	proto.Merge(resp.Bundle, req.Bundle)
	resp.Timeout = false
	resp.ParamOrResult = req.Param
	resp.Resolved = nil
	resp.Exit = nil
	s.chRecv <- resp
	return syscall.KernelErr_NoError
}

func (s *SingleApproach) HandleReturnValue(req *syscall.ReturnValueRequest) syscall.KernelErr {
	resp := &syscall.ReadOneResponse{}
	resp.Bundle = &syscall.MethodBundle{}
	proto.Merge(resp.Bundle, req.Bundle)
	resp.Timeout = false
	resp.ParamOrResult = req.Result
	resp.ResultErr = req.ResultError
	resp.Resolved = nil
	resp.Exit = nil
	s.chFin <- resp
	return syscall.KernelErr_NoError
}

//
// Adapters
//

// There are two methods that overlap with the Receiver case, Ch() and
// Handle.  We use an adapter for each one.
type recvAdapter struct {
	inner *SingleApproach
}

func (r *recvAdapter) Ch() chan *syscall.ReadOneResponse {
	return r.inner.chRecv
}

func (r *recvAdapter) TimeoutInMillis() int {
	return timeout
}

type retvalAdapter struct {
	inner *SingleApproach
}

func (r *retvalAdapter) Ch() chan *syscall.ReadOneResponse {
	return r.inner.chFin
}

func (r *retvalAdapter) TimeoutInMillis() int {
	return timeout
}
