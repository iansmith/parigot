package wheeler

import (
	"log"

	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Wheeler interface {
	In() chan InProtoPair
}

// the only wheeler
var _wheeler = newWheeler()

// OutProtoPair is the return type of a message to the wheeler.
// it is sent through the channel given as part of the request.
// If the err != 0, the msg should be ignored.   If the err
// is 0, then the msg will be non-nil.
type OutProtoPair struct {
	A   *anypb.Any
	Err syscall.KernelErr
}

// InProtoPair is a request and the channel to send the error
// or response to.
type InProtoPair struct {
	Msg proto.Message
	Ch  chan OutProtoPair
}

// hostBinding creates a connection between a given
// service implementation and the host that it lives
// on.
type hostBinding struct {
	service id.ServiceId
	host    id.HostId
}

// serviceData is the set of values that we need to be
// able to update and check so we can know what when we are
// ready to start a program running.
type serviceData struct {
	reg, exp, run bool
}

// wheeler is the type that implements system calls.  It actually
// reads a channel and responds to the requests one by one.
// One can think of it as a wheel in that multiple different
// callers are all trying to get info to the wheeler and it
// is spinning around to take each one in turn.
type wheeler struct {
	ch               chan InProtoPair
	pkgToServiceImpl map[string]map[string][]hostBinding
	hostToService    map[string][]id.ServiceId
	serviceToHost    map[string]id.HostId
}

// newWheeler returns a wheeler, but this should be called
// only during the init() function.
func newWheeler() *wheeler {
	w := &wheeler{
		ch:               make(chan InProtoPair),
		pkgToServiceImpl: make(map[string]map[string][]hostBinding),
		hostToService:    make(map[string][]id.ServiceId),
		serviceToHost:    make(map[string]id.HostId),
	}
	go w.Run()
	return w
}

func In() chan InProtoPair {
	return _wheeler.ch
}

// export implements the export for all the given types.
// It binds each type to the hostname provided.  Note
// that this may imply the creating of a service inside
// the wheeler for a service that is not yet registered.
func (w *wheeler) export(req *syscall.ExportRequest) (*anypb.Any, syscall.KernelErr) {
	sid := id.UnmarshalServiceId(req.ServiceId)
	hid := id.UnmarshalHostId(req.HostId)
	w.addHost(hid, sid)
	for _, fqn := range req.GetService() {
		pkg := fqn.GetPackagePath()
		name := fqn.GetService()
		pkg2map, ok := w.pkgToServiceImpl[pkg]
		if !ok {
			pkg2map = make(map[string][]hostBinding)
			w.pkgToServiceImpl[pkg] = pkg2map
		}
		allBind, ok := pkg2map[name]
		if !ok {
			allBind = []hostBinding{}
			pkg2map[name] = allBind
		}
		allBind = append(allBind, hostBinding{
			service: sid,
			host:    hid,
		})
		pkg2map[name] = allBind
	}
	var a anypb.Any
	merr := a.MarshalFrom(&syscall.ExportResponse{})
	if merr != nil {
		return nil, syscall.KernelErr_MarshalFailed
	}
	return &a, syscall.KernelErr_NoError
}

func (w *wheeler) addHost(hid id.HostId, sid id.ServiceId) {
	w.serviceToHost[sid.String()] = hid
	allSvc, ok := w.hostToService[hid.String()]
	if !ok {
		allSvc = []id.ServiceId{}
		w.hostToService[hid.String()] = []id.ServiceId{}
	}
	allSvc = append(allSvc, sid)
	w.hostToService[hid.String()] = allSvc
	w.checkHost(hid, allSvc)
}

func (w *wheeler) checkHost(hid id.HostId, allSvc []id.ServiceId) {
	//ignored for now, should be a status check to see if the
	//alleged services are still there
}

func (w *wheeler) Run() {
	for {
		in := <-w.ch
		desc := in.Msg.ProtoReflect().Descriptor()
		var result proto.Message
		var err syscall.KernelErr
		switch desc.FullName() {
		case "syscall.v1.ExportRequest":
			log.Printf("xxx -- reached wheeler: export\n")
			result, err = w.export((*syscall.ExportRequest)(in.Msg.(*syscall.ExportRequest)))
		default:
			log.Printf("ERROR! wheeler received unknown type %s", desc.FullName())
			continue
		}
		var a anypb.Any
		e := a.MarshalFrom(result)
		if e != nil {
			in.Ch <- OutProtoPair{nil, syscall.KernelErr_MarshalFailed}
			return
		}
		outPair := OutProtoPair{
			A:   &a,
			Err: err,
		}
		log.Printf("xxx -- reached wheeler: export sending result\n")

		in.Ch <- outPair

	}
}
