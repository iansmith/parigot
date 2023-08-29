package kernel

import (
	"log"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
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

// HandleDispatch is called to generate a ReadOneResponse given an input dispatch request.
// It handles the exit call as a special case.
func (s *SingleApproach) HandleDispatch(req *syscall.DispatchRequest) syscall.KernelErr {
	resp := &syscall.ReadOneResponse{}
	resp.Bundle = &syscall.MethodBundle{}

	sid := id.UnmarshalServiceId(req.Bundle.ServiceId)
	mid := id.UnmarshalMethodId(req.Bundle.MethodId)
	//
	// special case for exit
	//
	if sid.Equal(apishared.ExitService) && mid.Equal(id.MethodId(apishared.ExitMethod)) {
		log.Printf("xxx -- got an exit request")
		pair := syscall.ExitPair{}
		if e := req.Param.UnmarshalTo(&pair); e != nil {
			klog.Errorf("unable to unmarshal into exit pair: %v", e)
		}
		resp.Bundle = &syscall.MethodBundle{}
		proto.Merge(resp.Bundle, req.Bundle)
		resp.Exit = &pair
		resp.ParamOrResult = nil
		resp.Timeout = false
		s.chRecv <- resp
		return syscall.KernelErr_NoError
	}
	// anything other than exit comes through here
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
