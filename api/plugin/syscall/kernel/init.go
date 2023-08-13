package kernel

import (
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

var _ Kernel = Kernel(&kdata{})

func InitSingle() (*kdata, bool) {
	netOut := make(chan proto.Message)

	nameserver := NewSimpleNameServer(netOut)
	starter := NewStarter()
	single := NewSingleApproach()
	recv := &recvAdapter{single}
	retval := &retvalAdapter{single}

	if !single.Init() {
		klog.Errorf("unable to initialize approach 'Single'")
		return nil, false
	}
	k := newKData()

	if kerr := k.SetApproach(recv, retval, nameserver, starter); kerr != syscall.KernelErr_NoError {
		klog.Errorf("unable to install approach 'Single'")
		return nil, false
	}

	return k, true
}
