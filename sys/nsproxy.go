package sys

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"
	"sync"
	"time"

	pcontext "github.com/iansmith/parigot/context"
	netmsg "github.com/iansmith/parigot/g/msg/net/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/dep"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var envVerbose = os.Getenv("PARIGOT_VERBOSE")
var netnameserverVerbose = false || envVerbose != ""

// servicePort is the port that a service listens on for incoming requests.
// This port is fixed because each service has its own container, thus its
// own udp portspace.
const servicePort = 13331

var koopmanTable = crc32.MakeTable(crc32.Koopman)
var writeTimeout = 250 * time.Millisecond
var readTimeout = writeTimeout
var longReadTimeout = 10 * time.Second
var readBufferSize = 4096

// if your message doesn't start with this, you have lost sync and should close the connection
// so we can try to reconnect
var magicStringOfBytes = uint64(0x1789071417760704)
var frontMatterSize = 8 + 4
var trailerSize = 4

const parigotNSPort = 13330
const parigotNSHost = "parigot_ns"

// NSProxy is a nameserver implementation that actually just redirects things to the network. It redirects
// things that need to do service locate to the nameserver and calls to remote implementations.  If you
// use this as a server, then it listens for incoming requests and responds to them.
type NSProxy struct {
	*NSCore
	//local *LocalNameServer

	// info about us, in case we need to tell someone who/where we are
	localKey  dep.DepKey
	localAddr string

	// this listens for other clients calling to us for an implementation
	listener *QuicListener
	inCh     chan *NetResult

	// this is when we want to call somebody else's implementation
	serviceLoc map[string]string          // sid->addr
	rpcCaller  map[string]*quicCaller     // addr -> our caller (if we've called before)
	rpcChan    map[string]chan *NetResult // addr -> result ch

	// this is for calls to the remote nameserver
	nsCall *quicCaller
	nsCh   chan *NetResult
}

func NewNSProxy(ctx context.Context, addr string) *NSProxy {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("unable to get our own hostname in swarm %v", err))
	}
	if os.Getenv("HOSTNAME") != "" {
		hostname = os.Getenv("HOSTNAME")
	}
	myAddr := fmt.Sprintf("%s:%d", hostname, servicePort)
	netnameserverPrint(ctx, "NewNetNameserver ", "our address is %s", myAddr)

	inCh := make(chan *NetResult)
	nsCh := make(chan *NetResult)
	nns := &NSProxy{
		NSCore: NewNSCore(false),
		//local:     loc,
		localAddr: myAddr,
		localKey:  NewDepKeyFromAddr(myAddr),
		inCh:      inCh,
		listener:  NewQuicListener(servicePort, ParigotProtoRPC, inCh),
		nsCh:      nsCh,
		nsCall: newQuicCaller(
			fmt.Sprintf("%s:%d", parigotNSHost, parigotNSPort),
			ParigotProtoNameServer, nsCh),
		serviceLoc: make(map[string]string),
		rpcCaller:  make(map[string]*quicCaller),
		rpcChan:    make(map[string]chan *NetResult),
	}
	return nns
}

// makeRequest is just a convenience wrapper around constructing a netResult
// and doing the write+read to the appropriate channels.  this is ONLY for the
// nameserver.
func (n *NSProxy) makeRequest(a *anypb.Any) *anypb.Any {
	nr := NetResult{}
	nr.SetData(a)
	nr.SetKey(n.localKey)
	respCh := make(chan *anypb.Any)
	nr.SetRespChan(respCh)
	n.nsCall.ch <- &nr
	result := <-respCh
	return result
}

// Export is where user code ends up when they call Export, albeit via
// a system call into the kernel (this side).  This invokes Export on the
// remote nameserver to do the startup sequence.
func (n *NSProxy) Export(ctx context.Context, key dep.DepKey, packagePath, service string) lib.Id {
	sData := n.NSCore.GetSData(packagePath, service)
	if sData == nil {
		sData = n.NSCore.create(ctx, key, packagePath, service)
	}
	expInfo := &netmsg.ExportInfo{
		PackagePath: packagePath,
		Service:     service,
		Addr:        n.localAddr,
		ServiceId:   lib.Marshal[protosupportmsg.ServiceId](sData.serviceId),
	}
	expReq := &netmsg.ExportRequest{
		Export: []*netmsg.ExportInfo{expInfo},
	}
	var any anypb.Any
	err := any.MarshalFrom(expReq)
	if err != nil {
		return lib.NewKernelError(lib.KernelMarshalFailed)
	}
	result := n.makeRequest(&any)
	if result == nil {
		return lib.NewKernelError(lib.KernelNetworkFailed)
	}
	netnameserverPrint(ctx, "EXPORT ", "remote 4, result is %s", result.TypeUrl)
	expResp := netmsg.ExportResponse{}
	err = result.UnmarshalTo(&expResp)
	if err != nil {
		netnameserverPrint(ctx, "EXPORT ", "export remote 4a:%v", err)
		return lib.NewKernelError(lib.KernelUnmarshalFailed)
	}
	respOk := lib.Unmarshal(expResp.KernelErr)
	netnameserverPrint(ctx, "EXPORT ", " export remote 4b:%s", respOk.Short())
	if respOk.IsError() {
		return respOk
	}
	netnameserverPrint(ctx, "EXPORT ", "result was %s", result.TypeUrl)
	return lib.NoError[*protosupportmsg.KernelErrorId]()
}

// ExitWhenInFlightEmpty not implemented yet.
func (n *NSProxy) ExitWhenInFlightEmpty() bool {
	panic("ExitWhenInFlightEmpty not implemented yet")
}

// Require is where user code ends up when they call Require, albeit via
// a system call into the kernel (this side).  This invokes Require on the
// remote nameserver to do the startup sequence.
func (n *NSProxy) Require(ctx context.Context, key dep.DepKey, packagePath, service string) lib.Id {
	reqInfo := &netmsg.RequireInfo{
		PackagePath: packagePath,
		Service:     service,
		Addr:        n.localAddr,
	}
	reqReq := &netmsg.RequireRequest{
		Require: []*netmsg.RequireInfo{reqInfo},
	}
	var any anypb.Any
	err := any.MarshalFrom(reqReq)
	if err != nil {
		return lib.NewKernelError(lib.KernelMarshalFailed)
	}
	reqResult := n.makeRequest(&any)
	if reqResult == nil {
		return lib.NewKernelError(lib.KernelNetworkFailed)
	}
	netnameserverPrint(ctx, "REQUIRE ", "reqResult is %s\n", reqResult.TypeUrl)
	reqResp := netmsg.RequireResponse{}
	err = reqResult.UnmarshalTo(&reqResp)
	if err != nil {
		netnameserverPrint(ctx, "REQUIRE ", "remote call failed:%v", err)
		return lib.NewKernelError(lib.KernelUnmarshalFailed)
	}
	respOk := lib.Unmarshal(reqResp.KernelErr)
	netnameserverPrint(ctx, "REQUIRE ", "remote result :%s", respOk.Short())
	if respOk.IsError() {
		return respOk
	}
	netnameserverPrint(ctx, "REQUIRE", "DONE!")
	return lib.NoError[*protosupportmsg.KernelErrorId]()
}

// CloseService is where user code ends up when they call CloseService, albeit via
// a system call into the kernel (this side).  This invokes CloseService on the
// remote nameserver to indicate that a service has completing binding its names.
func (n *NSProxy) CloseService(ctx context.Context, key dep.DepKey, packagePath, service string) lib.Id {
	req := &netmsg.CloseServiceRequest{Addr: key.String(),
		PackagePath: packagePath,
		Service:     service,
	}
	any := &anypb.Any{}
	err := any.MarshalFrom(req)
	netnameserverPrint(ctx, "CLOSESERVICE ", "xxx export remote 1 with %T", req)
	if err != nil {
		return lib.NewKernelError(lib.KernelMarshalFailed)
	}
	result := n.makeRequest(any)
	if result == nil {
		return lib.NewKernelError(lib.KernelNetworkFailed)
	}
	netnameserverPrint(ctx, "CLOSESERVICE ", " 2 with result %v", result.TypeUrl)
	resp := netmsg.CloseServiceResponse{}
	err = result.UnmarshalTo(&resp)
	if err != nil {
		return lib.NewKernelError(lib.KernelUnmarshalFailed)
	}
	respErr := lib.Unmarshal(resp.GetKernelErr())
	if respErr.IsError() {
		return respErr
	}
	netnameserverPrint(ctx, "CLOSESERVICE ", "xxx export remote 3")
	netnameserverPrint(ctx, "CLOSESERVICE ", "result was %s", result.TypeUrl)
	return lib.NoError[*protosupportmsg.KernelErrorId]()
}

// CloseService is where user code ends up when they call CloseService, albeit via
// a system call into the kernel (this side).  This invokes CloseService on the
// remote nameserver to indicate that a service has completing binding its names.
func (n *NSProxy) HandleMethod(p *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	panic("HandleMethod")
}

// GetService looks up the service that is given as a parameter (pkgPath.service) and returns
// the service id for that service.  The last two parameters will be nil, "" if
// there was no error, otherwise will contain the error details.
func (n *NSProxy) GetService(ctx context.Context, _ dep.DepKey, pkgPath, service string) (lib.Id, lib.Id, string) {
	req := &netmsg.GetServiceRequest{
		PackagePath: pkgPath,
		Service:     service,
	}
	a := &anypb.Any{}
	err := a.MarshalFrom(req)
	if err != nil {
		return nil, lib.NewKernelError(lib.KernelMarshalFailed),
			fmt.Sprintf("failed trying to marshal request: %v", err)
	}

	result := n.makeRequest(a)
	if result == nil {
		return nil, lib.NewKernelError(lib.KernelNetworkFailed), "network failed sending request"
	}
	resp := &netmsg.GetServiceResponse{}
	err = result.UnmarshalTo(resp)
	if err != nil {
		return nil, lib.NewKernelError(lib.KernelUnmarshalFailed),
			fmt.Sprintf("failed to unmarshal response (GetService): %v", err)
	}
	addr := resp.GetAddr()
	sidPtr := resp.GetSid()
	sid := lib.Unmarshal(sidPtr)
	netnameserverPrint(ctx, "GETSERVICE ", "addr is %s and sid is %s", addr, sid.Short())
	n.serviceLoc[sid.String()] = addr
	return sid, nil, ""
}

func (n *NSProxy) FindMethodByName(ctx context.Context, key dep.DepKey, serviceId lib.Id, method string) (*callContext, lib.Id, string) {
	netnameserverPrint(ctx, "FINDMETHBYNAME", "key is %s, service id %s, method %s (n.serviceLoc is nil? %v)",
		key.String(), serviceId.Short(), method, n.serviceLoc == nil)
	// we need to make sure we have the sid mapping
	_, ok := n.serviceLoc[serviceId.String()]
	if !ok {
		netnameserverPrint(ctx, "FINDMETHODBYNAME", "failed to get network addr for %s", serviceId.Short())
		return nil, lib.NewKernelError(lib.KernelNotFound),
			fmt.Sprintf("unable to find service %s", serviceId.String())
	}
	var mid lib.Id
	sData := n.NSCore.ServiceData(serviceId)
	if sData != nil {
		cachedMethodId, ok := sData.method.Load(method)
		if ok {
			mid = cachedMethodId.(lib.Id)
		}
	}
	netnameserverPrint(ctx, "FINDMETHODBYNAME", "xxx the key is %s", key.String())

	// we are going to return the callContext we'll need for the RPC
	callCtx := &callContext{
		respCh: make(chan *syscallmsg.ReturnValueRequest),
		mid:    mid,
		cid:    lib.NewId[*protosupportmsg.CallId](),
		method: method,
		sid:    serviceId,
		sender: key,
		param:  &anypb.Any{},
		pctx:   &protosupportmsg.Pctx{},
	}
	netnameserverPrint(ctx, "FINDMETHODBYNAME", "done with call context %s", callCtx.cid.Short())
	n.NSCore.addCallContextMapping(callCtx.cid, callCtx)
	return callCtx, nil, ""
}

func (n *NSProxy) GetInfoForCallId(ctx context.Context, target lib.Id) *callContext {
	return n.NSCore.getContextForCallId(ctx, target)
}

// CallService is the implementation (in the kernel) for when you have a remote (across
// the network) nameserver that you need to make an RPC call to.
func (n *NSProxy) CallService(ctx context.Context, key dep.DepKey, info *callContext) (*syscallmsg.ReturnValueRequest, lib.Id, string) {
	//this req is only used when the CallService call wants to return an error
	//or when have completed a network call successfully
	req := &syscallmsg.ReturnValueRequest{
		Call:        lib.Marshal[protosupportmsg.CallId](info.cid),
		Method:      lib.Marshal[protosupportmsg.MethodId](info.mid),
		Result:      nil,
		Pctx:        nil,
		ExecError:   "could not find the location of service " + info.sid.Short(),
		ExecErrorId: lib.Marshal[protosupportmsg.BaseId](lib.NewKernelError(lib.KernelNotFound)),
	}

	netnameserverPrint(ctx, "CALLSERVICE", "key is %s", key)
	// do we know where the service is?
	addr, ok := n.serviceLoc[info.sid.String()]
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound), "could not find the location of service " + info.sid.Short()
	}
	// have we called it before?
	//caller, ok := n.rpcCaller[addr]
	ch := n.rpcChan[addr]
	if !ok {
		ch = make(chan *NetResult)
		caller := newQuicCaller(addr, ParigotProtoRPC, ch)
		n.rpcCaller[addr] = caller
		n.rpcChan[addr] = ch
	}
	nr := NetResult{}
	rpcReq := netmsg.RPCRequest{
		Pctx:       info.pctx,
		Param:      info.param,
		ServiceId:  lib.Marshal[protosupportmsg.ServiceId](info.sid),
		MethodId:   nil,
		MethodName: info.method,
		Sender:     n.localAddr,
	}
	a := anypb.Any{}
	err := a.MarshalFrom(&rpcReq)
	if err != nil {
		return nil, lib.NewKernelError(lib.KernelMarshalFailed), "unable to marshal the return value request:" + err.Error()
	}
	nr.SetData(&a)
	nr.SetKey(key) //????
	respCh := make(chan *anypb.Any)
	nr.SetRespChan(respCh)
	ch <- &nr
	result := <-respCh
	if result == nil {
		netnameserverPrint(ctx, "CALLSERVICE", "failed in call to %s", info.sid.Short())
		return nil, lib.NewKernelError(lib.KernelNetworkFailed), "network failed"
	}
	resp := netmsg.RPCResponse{}
	err = result.UnmarshalTo(&resp)
	if err != nil {
		return nil, lib.NewKernelError(lib.KernelUnmarshalFailed), "unable to unmarshal the return value response:" + err.Error()
	}
	req.Result = resp.GetResult()
	req.Pctx = resp.GetPctx()
	return req, nil, ""
}

func (n *NSProxy) RunBlock(ctx context.Context, key dep.DepKey) (bool, lib.Id) {
	req := &netmsg.RunBlockRequest{
		Waiter: false,
		Addr:   n.localAddr,
	}
	any := &anypb.Any{}
	err := any.MarshalFrom(req)
	netnameserverPrint(ctx, "RUNBLOCK ", "my addr is %s", n.localAddr)
	if err != nil {
		return false, lib.NewKernelError(lib.KernelMarshalFailed)
	}
	result := n.makeRequest(any)
	if result == nil {
		netnameserverPrint(ctx, "RUNBLOCK ", "call to NS for runblock failed ")
		return false, lib.NewKernelError(lib.KernelNetworkFailed)
	}
	resp := netmsg.RunBlockResponse{}
	err = result.UnmarshalTo(&resp)
	if err != nil {
		return false, lib.NewKernelError(lib.KernelUnmarshalFailed)
	}
	respErr := lib.Unmarshal(resp.GetErrId())
	if respErr.IsError() {
		netnameserverPrint(ctx, "RUNBLOCK ", " read error from server %s", respErr.Short())
		return false, respErr
	}
	netnameserverPrint(ctx, "RUNBLOCK", "DONE!")
	return !resp.GetTimedOut(), nil
}

func (n *NSProxy) RunIfReady(ctx context.Context, key dep.DepKey) []dep.DepKey {
	panic("we need to talk to the network to do RunIfReady")
}

// StartFailedInfo is supposed to return details about why the
// startup failed (e.g. a loop of dependencies). For now, we don't
// have a way to calculate this in the network case.
func (n *NSProxy) StartFailedInfo(ctx context.Context) string {
	return n.NSCore.StartFailedInfo(ctx)
}

// LenSyncMap is a utility return the number of keys in a *sync.Map
func LenSyncMap(m *sync.Map) int {
	count := 0
	m.Range(func(k, v any) bool {
		count++
		return true
	})
	return count
}

// BlockUntilCall handles _incoming_ RPC requests.
func (n *NSProxy) BlockUntilCall(ctx context.Context, key dep.DepKey, _ bool) *callContext {
	netnameserverPrint(ctx, "BlockUntilCall ", " key is %s and inCh is %p", key.String(), n.inCh)
	a := <-n.inCh
	req := netmsg.RPCRequest{}
	netnameserverPrint(ctx, "BlockUntilCall ", " got a request through the channel %p, a==nil? %v", n.inCh, a == nil)
	err := a.Data().UnmarshalTo(&req)
	if err != nil {
		netnameserverPrint(ctx, "BlockUntilCall ", "error trying to unmarshal request: %v", err)
		a.RespChan() <- nil
		return nil
	}
	// the call id will be nil because the caller doesn't know anything about what is
	// going on in our address space.  the method id MIGHT be nil, if he hasn't cache the
	// method id from some previous call.
	cid := lib.NewId[*protosupportmsg.CallId]()
	sid := lib.Unmarshal(req.GetServiceId())
	var mid lib.Id // empty
	if req.GetMethodId() != nil {
		mid = lib.Unmarshal(req.GetMethodId())
		// cross check,just in case
		sData := n.NSCore.ServiceData(sid)
		if sData != nil {
			otherIdAny, _ := sData.method.Load(req.GetMethodName())
			otherId := otherIdAny.(lib.Id)
			if !mid.Equal(otherId) {
				netnameserverPrint(ctx, "BlockUntilCall ", "WARN: ignorning provided mid %s because doesn't match our method table for %s",
					mid.Short(), req.GetMethodName())
				mid = otherId
			}
		}
	} else {
		sData := n.NSCore.ServiceData(sid)
		if sData == nil {
			netnameserverPrint(ctx, "BlockUntilCall ",
				"WARN: Unable to find sData for sid %s -- number of entries in package table %d",
				sid.Short(), LenSyncMap(n.NSCore.packageRegistry))
			n.NSCore.packageRegistry.Range(func(k, v any) bool {
				key := k.(string)
				netnameserverPrint(ctx, "BlockUntilCall ", "key in pkgReg=%s", key)
				sMap := v.(*sync.Map)
				sMap.Range(func(serviceId, serviceData any) bool {
					s := serviceId.(string)
					sd := serviceData.(*ServiceData)
					netnameserverPrint(ctx, "BlockUntilCall ", "key in service table=%s, sid on sData %s", s, sd.GetServiceId().Short())
					return true
				})
				return true
			})
		} else {
			var ok bool
			midAny, ok := sData.method.Load(req.GetMethodName())
			mid = midAny.(lib.Id)
			if !ok {
				netnameserverPrint(ctx, "WARN: Unable to find method id for name %s", req.GetMethodName())
				mid = nil
			}
		}
	}

	info := &callContext{
		mid:    mid,
		cid:    cid,
		param:  req.GetParam(),
		pctx:   req.GetPctx(),
		method: req.GetMethodName(),
		sid:    sid,
		respCh: make(chan *syscallmsg.ReturnValueRequest), // xxx what should this be? a new channel? who listens to it?
		sender: NewDepKeyFromAddr(req.GetSender()),
	}
	n.NSCore.addCallContextMapping(info.cid, info)
	netnameserverPrint(ctx, "BlockUntilCall ", "finished creating info, callid=%s, method=%s and the key is %s", info.cid.Short(), info.method,
		key.String())

	go func(ctx context.Context, callId lib.Id, rinfo *callContext, netr *NetResult) {
		retReq := <-info.respCh
		if retReq == nil {
			netnameserverPrint(ctx, "BlockUntilCall [goroutine] ", "got an error signal on the call %s", callId)
			return
		}
		rpcResp := &netmsg.RPCResponse{
			Pctx:     retReq.Pctx,
			CallId:   lib.Marshal[protosupportmsg.CallId](callId),
			Result:   retReq.Result,
			KerrId:   lib.NoKernelError(),
			MethodId: lib.Marshal[protosupportmsg.MethodId](mid),
		}
		netnameserverPrint(ctx, "BlockUntilCall [goroutine] ", "got return result on call %s from %s: %d bytes of result, %d bytes of pctx",
			callId, rinfo.sender, proto.Size(retReq.Result), proto.Size(retReq.Pctx))
		aResp := &anypb.Any{}
		err := aResp.MarshalFrom(rpcResp)
		if err != nil {
			a.respCh <- nil
			netnameserverPrint(ctx, "BlockUntilCall [goroutine] ", "failed to marshal response; %v", err)
			return
		}
		netr.respCh <- aResp
	}(ctx, info.cid, info, a)
	return info
}

func netnameserverPrint(ctx context.Context, method, spec string, arg ...interface{}) {
	if netnameserverVerbose {
		pcontext.LogFullf(ctx, pcontext.Debug, pcontext.Parigot,
			method, spec, arg...)
	}
}
