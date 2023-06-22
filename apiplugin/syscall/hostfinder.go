package syscall

import (
	"github.com/iansmith/parigot/apishared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
)

var _hostFinder HostFinder = newHostFinder

func finder() HostFinder

type hostFinder struct {
	name   map[string]*syscall.Host
	hostId map[string]*syscall.Host
}

func newHostFinder() *hostFinder {
	return &hostFinder{
		name:   make(map[string]*syscall.Host),
		hostId: make(map[string]*syscall.Host),
	}
}

func (h *hostFinder) FindByName(n string) *syscall.Host {
	host, ok := h.name[n]
	if !ok {
		return nil
	}
	return host
}

func (h *hostFinder) FindById(hid id.HostId) *syscall.Host {
	host, ok := h.hostId[hid.String()]
	if !ok {
		return nil
	}
	return host
}

func (h *hostFinder) AddHost(host *syscall.Host) syscall.KernelErr {
	hid := id.UnmarshalHostId(host.GetHostId())
	if hid.IsZeroOrEmptyValue() {
		panic("bad host id sent to the host finder")
	}
	if host.GetName() == "" {
		panic("bad host name sent to host finder")
	}
	h.name[host.GetName()] = host
	h.hostId[hid.String()] = host
	return syscall.KernelErr_NoError
}
