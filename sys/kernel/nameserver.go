package kernel

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

var _ Nameserver = &ns{}
var _ Registrar = &ns{}

type serviceSet map[string]struct{}

type hostSet map[string]struct{}

type ns struct {
	lock           sync.Mutex
	hostStrToId    map[string]id.HostId
	serviceStrToId map[string]id.ServiceId
	serviceToHost  map[string]id.HostId

	hostSvcMethodToName map[string]string
	hostSvcNameSet      map[string]id.MethodId

	serviceNameToSidSet map[string][]id.ServiceId

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
		serviceNameToSidSet: make(map[string][]id.ServiceId),
		matcher:             *newCallMatcher(),
		net:                 net,
	}
}

// AllHosts() returns all the known hosts.  This is useful for a
// broadcast message, like Exit.
func (n *ns) AllHosts() []id.HostId {
	n.lock.Lock()
	defer n.lock.Unlock()
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
	n.lock.Lock()
	defer n.lock.Unlock()
	s := h.String()
	_, ok := n.hostStrToId[s]
	if !ok {
		n.hostStrToId[s] = h
	}
}

// FindHost will return the Host that is running a particular service.
// If the service cannot be found, it returns the Zero value of host.
func (n *ns) FindHost(sid id.ServiceId) id.HostId {
	n.lock.Lock()
	defer n.lock.Unlock()
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

// PickService will return one of the services that is registered
// for a particular service name. It returns zero value if it cannot
// find the name.
func (n *ns) PickService(fqn FQName) id.ServiceId {
	n.lock.Lock()
	defer n.lock.Unlock()
	set, ok := n.serviceNameToSidSet[fqn.String()]
	if !ok {
		return id.ServiceIdZeroValue()
	}
	if len(set) == 0 {
		panic("pick service should not have a zero length set of services")
	}
	if len(set) == 1 {
		return set[0]
	}
	i := rand.Intn(len(set))
	return set[i]
}

// In() returns the channel for reading input requests.
func (n *ns) In() chan proto.Message {
	return n.net
}

// Register is called when any service registers. We copy in the host provided.
func (n *ns) Register(hid id.HostId, sid id.ServiceId, _ string) syscall.KernelErr {
	n.AddHost(hid)

	// we have wait to lock because AddHost Locks
	n.lock.Lock()
	defer n.lock.Unlock()

	n.serviceToHost[sid.String()] = hid
	return syscall.KernelErr_NoError
}

// BindMethod is called when any service registers.  We are building the list
// of known methods of a given service so we copy much of this info.
func (n *ns) Bind(hid id.HostId, sid id.ServiceId, mid id.MethodId, methName string) syscall.KernelErr {
	n.lock.Lock()
	defer n.lock.Unlock()
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

func (n *ns) Export(hid id.HostId, sid id.ServiceId, fqn FQName) syscall.KernelErr {
	n.lock.Lock()
	defer n.lock.Unlock()
	allSvc, ok := n.serviceNameToSidSet[fqn.String()]
	if !ok {
		n.serviceNameToSidSet[fqn.String()] = []id.ServiceId{}
	}
	allSvc = append(allSvc, sid)
	n.serviceNameToSidSet[fqn.String()] = allSvc
	return syscall.KernelErr_NoError
}

func (n *ns) MethodDetail(fqn FQName, methodName string) (id.HostId, id.ServiceId, id.MethodId, syscall.KernelErr) {
	n.lock.Lock()
	defer n.lock.Unlock()
	svc := n.PickService(fqn)
	if svc.IsZeroValue() {
		klog.Errorf("unable to find service for fully qualified name '%s'", fqn.String())
		return id.HostIdZeroValue(), id.ServiceIdZeroValue(), id.MethodIdZeroValue(), syscall.KernelErr_NotFound
	}
	host, ok := n.serviceToHost[svc.String()]
	if !ok {
		klog.Errorf("unable to find host for service '%s'", svc.String())
		return id.HostIdZeroValue(), id.ServiceIdZeroValue(), id.MethodIdZeroValue(), syscall.KernelErr_NotFound
	}
	key := keyOfTwoIdsAndName(host, svc, methodName)
	meth, ok := n.hostSvcNameSet[key]
	if !ok {
		klog.Errorf("unable to find method for service '%s' method '%s'", svc.Short(), methodName)
		return id.HostIdZeroValue(), id.ServiceIdZeroValue(), id.MethodIdZeroValue(), syscall.KernelErr_NotFound
	}
	return host, svc, meth, syscall.KernelErr_NoError
}

func (n *ns) FindMethod(hid id.HostId, sid id.ServiceId, mid id.MethodId) string {
	n.lock.Lock()
	defer n.lock.Unlock()
	key := keyOfThreeIds(hid, sid, mid)
	return n.hostSvcMethodToName[key]
}

func keyOfThreeIds(host id.HostId, sid id.ServiceId, mid id.MethodId) string {
	return fmt.Sprintf("%s.%s.%s", host.String(), sid.String(), mid.String())
}

func keyOfTwoIdsAndName(host id.HostId, sid id.ServiceId, methName string) string {
	return fmt.Sprintf("%s.%s.%s", host.String(), sid.String(), methName)
}
