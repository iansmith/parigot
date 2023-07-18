package syscall

import (
	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

type Wheeler interface {
	In() chan inProtoPair
}

// the only wheeler
var _wheeler = newWheeler()

// OutProtoPair is the return type of a message to the wheeler.
// it is sent through the channel given as part of the request.
// If the err != 0, the msg should be ignored.   If the err
// is 0, then the msg will be non-nil.
type OutProtoPair struct {
	Msg *proto.Message
	Err syscall.KernelErr
}

// InProtoPair is a request and the channel to send the error
// or response to.
type InProtoPair struct {
	Msg *proto.Message
	Ch  chan OutProtoPair
}

// hostBinding creates a connection between a given
// service implementation and the host that it lives
// on.
type hostBinding struct {
	service id.ServiceId
	host    id.HostId
}

// wheeler is the type that implements system calls.  It actually
// reads a channel and responds to the requests one by one.
// One can think of it as a wheel in that multiple different
// callers are all trying to get info to the wheeler and it
// is spinning around to take each one in turn.
type wheeler struct {
	ch                chan InProtoPair
	pkgToServiceName  map[string]map[string][]hostBinding
	serviceNameToImpl map[string][]hostBinding
}

// newWheeler returns a wheeler, but this should be called
// only during the init() function.
func newWheeler() *wheeler {
	w := &wheeler{
		ch:                make(chan InProtoPair),
		pkgToServiceName:  make(map[string]map[string][]hostBinding),
		serviceNameToImpl: make(map[string][]hostBinding),
	}
	go w.Run()
	return w
}

func (w *wheeler) In() chan InProtoPair {
	return w.ch
}

// export implements the export for all the given types.
// It binds each type to the hostname provided.  Note
// that this may imply the creating of a service inside
// the wheeler for a service that is not yet registered.
func (w *wheeler) export(req *syscall.ExportRequest) (*syscall.ExportResponse, syscall.KernelErr) {
	for _, fqn := range req.GetService() {

	}
}

func (w *wheeler) Run() {

}
