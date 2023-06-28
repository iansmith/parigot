package syscall

import (
	"github.com/iansmith/parigot/apishared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
)

var _hostFinder HostFinder = newHostFinder()

func finder() HostFinder {
	return _hostFinder
}

type hostFinder struct {
	name   map[string]*hostInfo
	hostId map[string]*hostInfo
}

type hostInfo struct {
	name string
	hid  id.HostId
}

func newHostFinder() *hostFinder {
	return &hostFinder{
		name:   make(map[string]*hostInfo),
		hostId: make(map[string]*hostInfo),
	}
}

func (h *hostFinder) FindByName(n string) *hostInfo {
	host, ok := h.name[n]
	if !ok {
		return nil
	}
	return host
}

func (h *hostFinder) FindById(hid id.HostId) *hostInfo {
	host, ok := h.hostId[hid.String()]
	if !ok {
		return nil
	}
	return host
}

func (h *hostFinder) AddHost(name string, hid id.HostId) syscall.KernelErr {
	hostInfo := &hostInfo{
		name: name,
		hid:  hid,
	}
	if hid.IsZeroOrEmptyValue() {
		panic("bad host id sent to the host finder")
	}
	if name == "" {
		panic("bad host name sent to host finder")
	}
	h.name[name] = hostInfo
	h.hostId[hid.String()] = hostInfo
	return syscall.KernelErr_NoError
}
