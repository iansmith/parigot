package lib

import (
	"github.com/iansmith/parigot/apishared/id"
	"google.golang.org/protobuf/proto"
)

type ErrorType interface {
	~int32
}

type Promise[T any, U ErrorType] struct {
	ready bool
	msg   T
	err   U
	cs    *ClientSideService
	cid   id.CallId
}

func NewPromiseProto[T proto.Message, U ErrorType](cs *ClientSideService, cid id.CallId) Promise[T, U] {
	return Promise[T, U]{
		cs:  cs,
		cid: cid,
	}
}
func NewPromise[T any, U ErrorType]() Promise[T, U] {
	return Promise[T, U]{}
}

func (p *Promise[T, U]) Resolved() bool {
	return p.ready
}

func (p *Promise[T, U]) Resolve(t T, u U, cid id.CallId) {
	if p.ready {
		panic("attempt to resolve an already resolved promise")
	}
	if !cid.IsZeroOrEmptyValue() {
		if !cid.Equal(p.cid) {
			panic("attempt to resolve a promise with wrong call id")
		}
	}
	p.msg = t
	p.err = u
	p.ready = true
}

func (p *Promise[T, U]) Rejected() bool {
	return p.ready && p.err != 0
}

// Error Only
type PromiseOnlyError[U ErrorType] struct {
	ready bool
	err   U
	cs    *ClientSideService
	cid   id.CallId
}

func (p *PromiseOnlyError[U]) Resolved() bool {
	return p.ready
}

func NewPromiseErrorOnly[U ErrorType](cs *ClientSideService, cid id.CallId) PromiseOnlyError[U] {
	return PromiseOnlyError[U]{
		cs:  cs,
		cid: cid,
	}
}
func (p *PromiseOnlyError[U]) Resolve(u U, cid id.CallId) {
	if p.ready {
		panic("attempt to resolve an already resolved promise")
	}
	if !cid.IsZeroOrEmptyValue() {
		if !cid.Equal(p.cid) {
			panic("attempt to resolve a promise with wrong call id")
		}
	}
	p.err = u
	p.ready = true
}
