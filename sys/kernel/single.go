package kernel

import (
	"sync"
	"time"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

const timeout = 50

type SingleApproach struct {
	lock               sync.Mutex
	init               bool
	seen               map[string]struct{}
	ns                 Nameserver
	dispatchChannel    map[string]chan *syscall.ReadOneResponse
	returnValueChannel map[string]chan *syscall.ReadOneResponse
}

func NewSingleApproach(nameserver *ns) *SingleApproach {
	return &SingleApproach{
		ns: nameserver,
	}
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
	s.seen = make(map[string]struct{})
	s.dispatchChannel = make(map[string]chan *syscall.ReadOneResponse)
	s.returnValueChannel = make(map[string]chan *syscall.ReadOneResponse)
	s.init = true
	return true
}

// HandleDispatch is called to generate a ReadOneResponse given an input dispatch request.
// It handles the exit call as a special case.
func (s *SingleApproach) HandleDispatch(req *syscall.DispatchRequest,
	dispChan chan *syscall.ReadOneResponse) syscall.KernelErr {
	s.lock.Lock()
	defer s.lock.Unlock()

	resp := &syscall.ReadOneResponse{}
	resp.Bundle = &syscall.MethodBundle{}
	sid := id.UnmarshalServiceId(req.Bundle.ServiceId)
	mid := id.UnmarshalMethodId(req.Bundle.MethodId)
	//
	// special case for exit
	//
	if sid.Equal(apishared.ExitService) && mid.Equal(id.MethodId(apishared.ExitMethod)) {
		pair := syscall.ExitPair{}
		if e := req.Param.UnmarshalTo(&pair); e != nil {
			klog.Errorf("unable to unmarshal into exit pair: %v", e)
		}
		resp.Bundle = &syscall.MethodBundle{}
		proto.Merge(resp.Bundle, req.Bundle)
		resp.Exit = &pair
		resp.ParamOrResult = nil
		resp.Timeout = false
		dispChan <- resp
		return syscall.KernelErr_NoError
	}
	// anything other than exit comes through here

	proto.Merge(resp.Bundle, req.Bundle)
	resp.Timeout = false
	resp.ParamOrResult = req.Param
	resp.Resolved = nil
	resp.Exit = nil
	dispChan <- resp
	return syscall.KernelErr_NoError
}

func (s *SingleApproach) HandleReturnValue(req *syscall.ReturnValueRequest,
	finishChan chan *syscall.ReadOneResponse) syscall.KernelErr {
	s.lock.Lock()
	defer s.lock.Unlock()

	resp := &syscall.ReadOneResponse{}
	resp.Bundle = &syscall.MethodBundle{}
	proto.Merge(resp.Bundle, req.Bundle)
	resp.Timeout = false
	resp.ParamOrResult = req.Result
	resp.ResultErr = req.ResultError
	resp.Resolved = nil
	resp.Exit = nil
	finishChan <- resp
	return syscall.KernelErr_NoError
}

func (s *SingleApproach) Register(hid id.HostId, sid id.ServiceId, debugName string) syscall.KernelErr {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
		go s.ReadNetworkMessages(hid)
	}()

	_, ok := s.seen[hid.String()]
	if !ok {
		d := make(chan *syscall.ReadOneResponse)
		f := make(chan *syscall.ReadOneResponse)
		s.dispatchChannel[hid.String()] = d
		s.returnValueChannel[hid.String()] = f

		K.AddReceiver(newRecvAdapter(hid, d, f, true))
		K.AddReceiver(newRecvAdapter(hid, d, f, false))
		s.seen[hid.String()] = struct{}{}
	}
	return syscall.KernelErr_NoError
}

func (s *SingleApproach) DispatchChan(hid id.HostId) chan *syscall.ReadOneResponse {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.dispatchChannel[hid.String()]
}
func (s *SingleApproach) FinishChan(hid id.HostId) chan *syscall.ReadOneResponse {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.returnValueChannel[hid.String()]
}

// This is used to send messages to the kernel that we have received
// something from the network.
type recvAdapter struct {
	host id.HostId
	ch   chan *syscall.ReadOneResponse
}

func newRecvAdapter(h id.HostId, dispatch, finish chan *syscall.ReadOneResponse, isDispatch bool) *recvAdapter {
	ch := dispatch
	if !isDispatch {
		ch = finish
	}
	adapter := &recvAdapter{
		//inner: s,
		host: h,
		ch:   ch,
	}

	return adapter
}

func (r *recvAdapter) Ch() chan *syscall.ReadOneResponse {
	return r.ch
}

func (r *recvAdapter) HostId() id.HostId {
	return r.host
}

func (r *recvAdapter) TimeoutInMillis() int {
	return timeout
}

// ReadNetworkMessages where messags to ths host "come in" from the
// outside world.  Note that this method does not lock because it
// calls methods that do lock.
func (s *SingleApproach) ReadNetworkMessages(host id.HostId) {
	var in chan proto.Message
	count := 0
	for {
		time.Sleep(100 * time.Millisecond)
		in = s.ns.FindHostChan(host)
		if in == nil && count > 10 {
			klog.Errorf("unable to find way to send to host", "host", host.Short())
			time.Sleep(500 * time.Millisecond)
			count++
		} else {
			break
		}
	}

	for {
		req := <-in
		desc := req.ProtoReflect().Descriptor()
		var err syscall.KernelErr
		switch desc.FullName() {
		case "syscall.v1.DispatchRequest":
			dispChan := s.DispatchChan(host)
			err = s.HandleDispatch((*syscall.DispatchRequest)(req.(*syscall.DispatchRequest)), dispChan)
		case "syscall.v1.ReturnValueRequest":
			finishChan := s.FinishChan(host)
			err = s.HandleReturnValue((*syscall.ReturnValueRequest)(req.(*syscall.ReturnValueRequest)), finishChan)
		default:
			klog.Errorf("unexpected type received from network: %s, ignoring", desc.FullName())
		}
		if err != syscall.KernelErr_NoError {
			klog.Errorf("error returned from attempt to handle an incoming request on the network: %s",
				syscall.KernelErr_name[int32(err)])
		}
	}
}
