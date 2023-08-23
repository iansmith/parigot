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
	go networkCallReader(single, net)
	return k, true
}

func networkCallReader(single *SingleApproach, net chan proto.Message) {
	for {
		req := <-net
		desc := req.ProtoReflect().Descriptor()
		var err syscall.KernelErr
		switch desc.FullName() {
		case "syscall.v1.DispatchRequest":
			err = single.HandleDispatch((*syscall.DispatchRequest)(req.(*syscall.DispatchRequest)))
		case "syscall.v1.ReturnValueRequest":
			err = single.HandleReturnValue((*syscall.ReturnValueRequest)(req.(*syscall.ReturnValueRequest)))
		default:
			klog.Errorf("unexpected type received from network: %s, ignoring", desc.FullName())
		}
		if err != syscall.KernelErr_NoError {
			klog.Errorf("error returned from attempt to handle an incoming request on the network: %s",
				syscall.KernelErr_name[int32(err)])
		}
	}

}
