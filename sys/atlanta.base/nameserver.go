package sys

import (
	"sync"

	"github.com/iansmith/parigot/lib"
)

const MaxService = 127

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
	methodIdToProcess map[lib.Id]*Process
}

func newServiceData() *serviceData {
	return &serviceData{
		serviceId:         nil,
		method:            make(map[string]lib.Id),
		methodIdToProcess: make(map[lib.Id]*Process),
	}
}

type nameServer struct {
	packageRegistry        map[string]*packageData
	serviceCounter         int
	serviceIdToServiceData map[lib.Id]*serviceData // accelerator only, we could walk to find this
	lock                   *sync.RWMutex
}

func newNameServer() *nameServer {
	return &nameServer{
		lock:                   new(sync.RWMutex),
		packageRegistry:        make(map[string]*packageData),
		serviceIdToServiceData: make(map[lib.Id]*serviceData),
		serviceCounter:         7, //first 8 are reserved
	}
}

// RegisterClientService connects a packagePath.service with a particular service id.
// This is used by client side code and is ususally called via the init() method in
// generated code.
func (n *nameServer) RegisterClientService(packagePath string, service string) (lib.Id, lib.Id) {
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
		n.serviceIdToServiceData[sData.serviceId] = sData
		pData.service[service] = sData
	}
	return sData.serviceId, nil
}

// FindMethodByName is called by the client side when doing a dispatch.  It needs to find the
// process that can handle the request.
func (n *nameServer) FindMethodByName(serviceId lib.Id, name string) (lib.Id, *Process) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	sData, ok := n.serviceIdToServiceData[serviceId]
	if !ok {
		return nil, nil
	}
	mid, ok := sData.method[name]
	if !ok {
		return nil, nil
	}
	return mid, sData.methodIdToProcess[mid]
}

// HandleMethod is called by the server side to indicate that it will handle a particular
// method call on a particular service.
func (n *nameServer) HandleMethod(pkgPath, service, method string, proc *Process) (lib.Id, lib.Id) {
	n.lock.Lock()
	defer n.lock.Unlock()

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
		n.serviceIdToServiceData[sData.serviceId] = sData
	}
	result := lib.NewMethodId()
	sData.method[method] = result
	sData.methodIdToProcess[result] = proc

	// xxx fix me, should be able to realize that a method does not exist and reject the attempt to
	// handle it

	//xxx fix me, there should be a limit on the number of methods per service
	return result, nil
}

// GetService can be called by either a client or a server. If this returns without error, the resulting
// serviceId can be used to be a client of the requested service.
func (n *nameServer) GetService(pkgPath, service string) (lib.Id, lib.Id) {
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
