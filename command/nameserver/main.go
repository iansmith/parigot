package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/iansmith/parigot/api/proto/g/pb/net"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys"
	"github.com/iansmith/parigot/sys/dep"

	"google.golang.org/protobuf/types/known/anypb"
)

const parigotNSPort = 13330
const timeoutClient = 6 // amount of time after which we decide your require(s) aint gonna happen
const sleepAmount = 2   //two secs between checks for timeout
var core = sys.NewNSCore(false)

type waitInfo struct {
	//ch        chan bool
	waitStart time.Time
	waitId    int
	// we keep this just we can reverse the mapping without walking
	waitKey dep.DepKey
	// we are still waiting to send the result back to the caller
	respChan chan *anypb.Any
}

var waitCounter = 0
var waitLock sync.Mutex
var typeToHostMap = make(map[string]string)
var runBlockWaitingList = make(map[string] /*dep.DepKey*/ *waitInfo)

func main() {
	go timeoutHandler()

	// the listener and the channel he talks to us on
	ch := make(chan *sys.NetResult)
	_ = sys.NewQuicListener(parigotNSPort, sys.ParigotProtoNameServer, ch)

	// This nameserver has three kinds of events he needs respond to.
	// 1) Incoming requests from clients.  These are read by a different
	// goroutine in the listener but it just queues all the requests in
	// our channel (ch).  2) Requests to unblock a particular waiting
	// program-- which is blocked waiting on us to respond to its
	// runBlock() request.  The _detection_ of when a program is ok
	// to run in handled by a different goroutine that then sends us
	// a pointer to the particular waitInfo of the elgible program.
	// 3) There is a timer that will cause us to timeout programs that
	// have been waiting on a runblock for too long.  The timer is another
	// go routine, but it just notifies on the timerCh.

	// we alert him to a new waiter with this empty message
	newWaiterCh := make(chan dep.DepKey)
	// he messages us with this channel
	readyCh := make(chan *waitInfo)

	// The timer thread  sends a message every sleepAmount seconds.
	timerCh := make(chan struct{})
	go func() {
		time.Sleep(sleepAmount * time.Second)
		timerCh <- struct{}{}
	}()

	go checkForReadyWaiter(newWaiterCh, readyCh)

	// process this loop until the end of time
	for {

		var nr *sys.NetResult

		select {
		case _ = <-timerCh:
			timeoutHandler()
			continue
		case ready := <-readyCh:
			sendRunBlockResponse(ready, true)
			continue
		case nr = <-ch:
		}

		p, err := nr.Data().UnmarshalNew()
		if err != nil {
			nr.RespChan() <- nil
			continue
		}

		switch m := p.(type) {
		case *net.CloseServiceRequest:
			log.Printf("dispatching to close service")
			err = closeService(m, nr.RespChan())
		case *net.ExportRequest:
			log.Printf("dispatching to export")
			err = export(m, nr.RespChan())
		case *net.GetServiceRequest:
			log.Printf("dispatching to get service")
			err = getService(m, nr.RespChan())
		case *net.RequireRequest:
			log.Printf("dispatching to get require")
			err = require(m, nr.RespChan())
		case *net.RunBlockRequest:
			log.Printf("dispatching to run block")
			err = runBlock(m, nr.RespChan(), newWaiterCh)
		default:
			panic(fmt.Sprintf("nameserver received a bundle from a client that it could not understand the type of:%T", p))
		}
		if err != nil {
			log.Printf("got error from processing client input: %v", err)
			// close the stream
			nr.RespChan() <- nil
		}
	}
}

func closeService(m *net.CloseServiceRequest, respChan chan *anypb.Any) error {
	log.Printf("close service , making sure it exists in our graph")
	key := sys.NewDepKeyFromAddr(m.GetAddr())
	id := core.CloseService(key, m.GetPackagePath(), m.GetService())
	if id != nil && id.IsError() {
		log.Printf("closeservice: can't close service %s.%s [%s] because %s",
			m.GetPackagePath(), m.GetService(), m.GetAddr(), id.Short())
		respChan <- nil
		return nil
	}

	// at the moment, we really don't track this in the network case
	resp := &net.CloseServiceResponse{
		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	log.Printf("close service, computing response")
	a := anypb.Any{}
	if err := a.MarshalFrom(resp); err != nil {
		respChan <- nil
		return err
	}
	log.Printf("close service DONE!")
	respChan <- &a
	return nil
}

func export(m *net.ExportRequest, respChan chan *anypb.Any) error {
	// we lock here ONLY because the map requireWaitingList is not safe to read when others
	// are updating it
	waitLock.Lock()
	defer waitLock.Unlock()
	var failure lib.Id
	log.Printf("export of %d elements", len(m.GetExport()))

	for _, export := range m.GetExport() {
		pkg := export.GetPackagePath()
		svc := export.GetService()
		sid := lib.UnmarshalServiceId(export.GetServiceId())
		name := fmt.Sprintf("%s.%s", pkg, svc)
		addr := sys.NewDepKeyFromAddr(export.GetAddr())
		log.Printf("EXPORT, key is %s, sid is %s calling core Export for %s",
			addr, sid, name)
		//core.CreateWithSid(addr, pkg, svc, sid)
		id := core.Export(addr, pkg, svc, sid)
		if id != nil {
			failure = id
			break
		}
	}
	if failure != nil {
		resp := &net.ExportResponse{
			KernelErr: lib.MarshalKernelErrId(failure),
		}
		a := anypb.Any{}
		if err := a.MarshalFrom(resp); err != nil {
			return err
		}
		respChan <- &a
		return nil
	}
	resp := &net.ExportResponse{
		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	a := anypb.Any{}
	if err := a.MarshalFrom(resp); err != nil {
		respChan <- nil
		return err
	}
	respChan <- &a
	return nil
}

func getService(m *net.GetServiceRequest, respChan chan *anypb.Any) error {
	pkg := m.GetPackagePath()
	svc := m.GetService()

	core.DumpSIDTables()

	sdata := core.GetSData(pkg, svc)
	if sdata == nil {
		respChan <- nil
		return lib.NewPerrorFromId("failed to find sData for "+pkg+"."+svc,
			lib.NewKernelError(lib.KernelNotFound))
	}

	resp := &net.GetServiceResponse{

		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	resp.Addr = sdata.GetKey().String()
	resp.Sid = lib.MarshalServiceId(sdata.GetServiceId())

	log.Printf("found service requested (%s.%s)=>%s on %s", pkg, svc, sdata.GetServiceId(), sdata.GetKey().String())
	a := anypb.Any{}
	if err := a.MarshalFrom(resp); err != nil {
		respChan <- nil
		return err
	}
	respChan <- &a
	return nil
}

func runBlock(m *net.RunBlockRequest, respChan chan *anypb.Any, newWaiter chan dep.DepKey) error {
	waitLock.Lock()
	defer waitLock.Unlock()

	log.Printf("runblock: %s,%v", m.GetAddr(), m.GetWaiter())
	key := sys.NewDepKeyFromAddr(m.GetAddr())
	info := &waitInfo{
		waitStart: time.Now(),
		waitId:    waitCounter,
		waitKey:   key,
		respChan:  respChan,
	}
	waitCounter++
	log.Printf("runblock: %s is waiting id %d", m.GetAddr(), info.waitId)

	if runBlockWaitingList[key.String()] != nil {
		log.Printf("runblock found addr %s in the waiting list already, ignoring", m.GetAddr())
		respChan <- nil
		return nil
	}
	if m.GetAddr() == "" {
		log.Printf("Run block called with empty address, ignoring")
		respChan <- nil
		return nil

	}
	// this creates the record in the waiting list
	runBlockWaitingList[key.String()] = info
	newWaiter <- key
	log.Printf("runblock: DONE!")
	return nil
}

func require(m *net.RequireRequest, respChan chan *anypb.Any) error {
	waitLock.Lock()
	defer waitLock.Unlock()

	for _, require := range m.GetRequire() {
		pkg := require.GetPackagePath()
		svc := require.GetService()
		key := sys.NewDepKeyFromAddr(require.GetAddr())

		// where should we put this?
		if id := core.Require(key, pkg, svc); id != nil && id.IsError() {
			respChan <- nil
			return lib.NewPerrorFromId("require failed", id)
		}
	}

	// tell them it's ok
	resp := &net.RequireResponse{
		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	a := anypb.Any{}
	if err := a.MarshalFrom(resp); err != nil {
		respChan <- nil
		return err
	}
	respChan <- &a
	return nil
}
