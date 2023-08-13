package kernel

import (
	"fmt"
	"strings"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

var _ Nameserver = &ns{}
var _ Registrar = &ns{}

type serviceSet map[string]struct{}

type hostSet map[string]struct{}

type ns struct {
	hostStrToId    map[string]id.HostId
	serviceStrToId map[string]id.ServiceId
	serviceToHost  map[string]id.HostId

	hostSvcMethodToName map[string]string
	hostSvcNameSet      map[string]id.MethodId

	host    hostSet
	service serviceSet

	matcher matcher

	net chan proto.Message
}

// NewSimpleNameServer returns a fully initialized instance of the NameServer
// interface.  This version is designed for the single situation.
func NewSimpleNameServer(net chan proto.Message) *ns {
	return &ns{
		hostStrToId:         make(map[string]id.HostId),
		host:                make(map[string]struct{}),
		service:             make(map[string]struct{}),
		serviceStrToId:      make(map[string]id.ServiceId),
		serviceToHost:       make(map[string]id.HostId),
		hostSvcMethodToName: make(map[string]string),
		hostSvcNameSet:      make(map[string]id.MethodId),
		matcher:             *newCallMatcher(),
		net:                 net,
	}
}

// AllHosts() returns all the known hosts.  This is useful for a
// broadcast message, like Exit.
func (n *ns) AllHosts() []id.HostId {
	result := make([]id.HostId, len(n.hostStrToId))
	count := 0
	for _, h := range n.hostStrToId {
		result[count] = h
		count++
	}
	return result
}

// AddHost can be called by hand if there is a host we need to know about.
// More frequently, it is automatically called by the use the Register() call.
func (n *ns) AddHost(h id.HostId) {
	s := h.String()
	_, ok := n.hostStrToId[s]
	if !ok {
		n.hostStrToId[s] = h
	}
}

// FindHost will return the Host that is running a particular service.
// If the service cannot be found, it returns the Zero value of host.
func (n *ns) FindHost(sid id.ServiceId) id.HostId {
	hid, ok := n.serviceToHost[sid.String()]
	if !ok {
		return id.HostIdZeroValue()
	}
	return hid
}

// FindHostChan will return the channel that allows messages to
// be written to the host.
func (n *ns) FindHostChan(hid id.HostId) chan<- proto.Message {
	return n.net
}

// In() returns the channel for reading input requests.
func (n *ns) In() chan proto.Message {
	return n.net
}

// Register is called when any service registers. We copy in the host provided.
func (n *ns) Register(hid id.HostId, _ id.ServiceId, _ string) syscall.KernelErr {
	n.AddHost(hid)
	return syscall.KernelErr_NoError
}

// BindMethod is called when any service registers.  We are building the list
// of known methods of a given service so we copy much of this info.
func (n *ns) Bind(hid id.HostId, sid id.ServiceId, mid id.MethodId, methName string) syscall.KernelErr {
	methName = strings.TrimSpace(methName)
	if methName == "" {
		return syscall.KernelErr_IdDispatch
	}
	key1 := keyOfThreeIds(hid, sid, mid)
	key2 := keyOfTwoIdsAndName(hid, sid, methName)
	n.hostSvcMethodToName[key1] = methName
	n.hostSvcNameSet[key2] = mid
	return syscall.KernelErr_NoError
}

func (n *ns) FindMethod(hid id.HostId, sid id.ServiceId, mid id.MethodId) string {
	key := keyOfThreeIds(hid, sid, mid)
	return n.hostSvcMethodToName[key]
}

func keyOfThreeIds(host id.HostId, sid id.ServiceId, mid id.MethodId) string {
	return fmt.Sprintf("%s.%s.%s", host.String(), sid.String(), mid.String())
}

func keyOfTwoIdsAndName(host id.HostId, sid id.ServiceId, methName string) string {
	return fmt.Sprintf("%s.%s.%s", host.String(), sid.String(), methName)
}
