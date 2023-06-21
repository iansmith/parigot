package lib

import (
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

type ErrorType interface {
	int32
}

type Promise[T proto.Message, U ErrorType] struct {
	ready bool
	msg   T
	err   U
	cs    *ClientSideService
}

func NewPromise[T proto.Message, U ErrorType](cs *ClientSideService) Promise[T, U] {
	return Promise[T, U]{
		cs: cs,
	}
}

func (p *Promise[T, U]) Resolve(fn func(T, U) syscall.KernelErr) syscall.KernelErr {
	if p.ready {
		return fn(p.msg, p.err)
	}
	for {
		p.cs.WaitResolve(notificationCh)
		if p.ready {
			break
		}
	}
	return fn(p.msg, p.err)
}
