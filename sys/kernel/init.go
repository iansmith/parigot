package kernel

import (
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

var _ Kernel = Kernel(&kdata{})

func InitSingle() (*kdata, bool) {
	net := make(chan proto.Message)

	nameserver := NewSimpleNameServer(net)
	starter := NewStarter()
	single := NewSingleApproach(nameserver)

	if !single.Init() {
		klog.Errorf("unable to initialize approach 'Single'")
		return nil, false
	}
	k := newKData()

	if kerr := k.SetApproach(nameserver, starter,
		single, nil, nil); kerr != syscall.KernelErr_NoError {
		klog.Errorf("unable to install approach 'Single'")
		return nil, false
	}
	return k, true
}
