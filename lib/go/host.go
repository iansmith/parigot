package lib

import (
	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
)

var _hostId = id.NewHostId()

func CurrentHostId() id.HostId {
	return _hostId
}

func NewInMemoryHost(name string) *syscall.Host {
	return &syscall.Host{
		HostId: id.NewHostId().Marshal(),
		Name:   name,
	}
}
