package sys

import (
	"fmt"
	"log"
	"sync"

	"github.com/iansmith/parigot/lib"
)

const MaxService = 127

type callContext struct {
	mid    lib.Id   // the method id this call is going to be made TO
	target *Process // the process this call is going to be made TO
	cid    lib.Id   // call id that should be be used by the caller to match results
	sender *Process // the process this call is going to be made FROM
}

type packageData struct {
	service map[string]*serviceData
}

func newPackageData() *packageData {
	return &packageData{
		service: make(map[string]*serviceData),
	}
}

type serviceData struct {
	serviceId         lib.Id
	method            map[string]lib.Id
	methodIdToProcess map[string] /*really method id*/ *Process
}

func newServiceData() *serviceData {
	return &serviceData{
		serviceId:         nil,
		method:            make(map[string]lib.Id),
		methodIdToProcess: make(map[string]*Process),
	}
}

type NameServer struct {
	packageRegistry        map[string]*packageData
	serviceCounter         int
	serviceIdToServiceData map[string] /*really service id*/ *serviceData // accelerator only, we could walk to find this
	inFlight               []*callContext
	lock                   *sync.RWMutex
}

func NewNameServer() *NameServer {
	return &NameServer{
		lock:                   new(sync.RWMutex),
		packageRegistry:        make(map[string]*packageData),
		serviceIdToServiceData: make(map[string]*serviceData),
		inFlight:               []*callContext{},
		serviceCounter:         7, //first 8 are reserved
	}
}

// RegisterClientService connects a packagePath.service with a particular service id.
// This is used by client side code and is ususally called via the init() method in
// generated code.
func (n *NameServer) RegisterClientService(packagePath string, service string) (lib.Id, lib.Id) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.serviceCounter >= MaxService {
		return nil, lib.NewKernelError(lib.KernelNamespaceExhausted)
	}

	pData, ok := n.packageRegistry[packagePath]
	if !ok {
		pData = newPackageData()
		n.packageRegistry[packagePath] = pData
	}
	sData, ok := pData.service[service]
	if !ok {
		sData = newServiceData()
		n.serviceCounter++
		sData.serviceId = lib.ServiceIdFromUint64(0, uint64(n.serviceCounter))
		n.serviceIdToServiceData[sData.serviceId.String()] = sData
		pData.service[service] = sData
	}
	return sData.serviceId, nil
}

// FindMethodByName is called by the client side when doing a dispatch.  This is where the client
// exchanges a service.id,name pair for the appropriate call context.  The call context is used
// by the calling client to 1) know where to send the message and 2) how to block waiting on
// the result.  The return result here is the property of the nameserver, don't mess with it,
// just read it.
func (n *NameServer) FindMethodByName(serviceId lib.Id, name string, caller *Process) *callContext {
	n.lock.Lock()
	defer n.lock.Unlock()

	sData, ok := n.serviceIdToServiceData[serviceId.String()]
	if !ok {
		return nil
	}
	mid, ok := sData.method[name]
	if !ok {
		return nil
	}
	target, ok := sData.methodIdToProcess[mid.String()]
	if !ok {
		return nil
	}
	cc := &callContext{
		mid:    mid,
		target: target,
		cid:    lib.NewCallId(),
		sender: caller,
	}
	return cc
}

// HandleMethod is called by the server side to indicate that it will handle a particular
// method call on a particular service.
func (n *NameServer) HandleMethod(pkgPath, service, method string, proc *Process) (lib.Id, lib.Id) {
	n.lock.Lock()
	defer n.lock.Unlock()

	log.Printf("reached the nameserver in process %+v", proc)
	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		pData = newPackageData()
		n.packageRegistry[pkgPath] = pData
	}
	sData, ok := pData.service[service]
	if !ok {
		sData = newServiceData()
		n.serviceCounter++
		sData.serviceId = lib.ServiceIdFromUint64(0, uint64(n.serviceCounter))
		pData.service[service] = sData
		n.serviceIdToServiceData[sData.serviceId.String()] = sData
	}
	result := lib.NewMethodId()
	print(fmt.Sprintf("assigning %s to the method %s in service %s\n", result, method, service))
	sData.method[method] = result
	sData.methodIdToProcess[result.String()] = proc

	// xxx fix me, should be able to realize that a method does not exist and reject the attempt to
	// handle it

	//xxx fix me, there should be a limit on the number of methods per service
	return result, nil
}

// GetService can be called by either a client or a server. If this returns without error, the resulting
// serviceId can be used to be a client of the requested service.
func (n *NameServer) GetService(pkgPath, service string) (lib.Id, lib.Id) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	sData, ok := pData.service[service]
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	return sData.serviceId, nil
}
