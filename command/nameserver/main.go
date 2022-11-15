package main

import (
	"fmt"
	"hash/crc32"
	"log"
	"sync"
	"time"

	"github.com/iansmith/parigot/g/pb/net"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys"
	"github.com/iansmith/parigot/sys/dep"

	"google.golang.org/protobuf/types/known/anypb"
)

// XXX these are duplicates of the values in the netnameserver.go code that is the client
var koopmanTable = crc32.MakeTable(crc32.Koopman)
var writeTimeout = 250 * time.Millisecond
var readTimeout = writeTimeout
var longTimeout = 10 * time.Second //ugh
var readBufferSize = 4096
var magicStringOfBytes = uint64(0x1789071417760704)

const frontMatter = 12
const trailer = 4
const readRetries = 8

const parigotNSPort = 13330

const timeoutClient = 6 // amount of time after which we decide your require aint gonna happen
const sleepAmount = 2   //two secs between checks for timeout

var core = sys.NewNSCore(false)

type waitInfo struct {
	ch        chan bool
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

	ch := make(chan *sys.NetResult)
	_ = sys.NewQuicListener(parigotNSPort, sys.ParigotProtoNameServer, ch)

	for {
		log.Printf("xxxx about to block on %p", ch)
		nr := <-ch
		log.Printf("xxxx read from %p", ch)

		p, err := nr.Data().UnmarshalNew()
		if err != nil {
			nr.RespChan() <- nil
			continue
		}

		log.Printf("dispatching new proto request")
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
			err = runBlock(m, nr.RespChan())
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
		name := fmt.Sprintf("%s.%s", pkg, svc)
		addr := sys.NewDepKeyFromAddr(export.GetAddr())
		log.Printf("EXPORT, key is %s, calling core Export for %s", addr, name)
		id := core.Export(addr, pkg, svc)
		if id != nil {
			failure = id
			break
		}
		//might be folks waiting for that
		notifyWaiter(name)
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

// notifyWaiter tells anybody on the waiting list that this new type is
// ready for consumption. This function does not lock, so the caller must
// be holding the lock.
func notifyWaiter(exportedTypeName string) {
	graph := core.DependencyGraph().AllEdge()
	candidateList := []dep.DepKey{}
	for _, eh := range graph {
		//look through the edges
		for _, req := range eh.Require() {
			if req == exportedTypeName {
				candidateList = append(candidateList, eh.Key())
				eh.RemoveRequire([]string{req})
				break
			}
		}
	}
	for _, candidate := range candidateList {
		core.RunIfReady(candidate, func(key dep.DepKey) {
			wait := runBlockWaitingList[key.String()]
			if wait == nil {
				log.Printf("NOTIFYWAITER unable to find %s on the waiting list", key)
				return // can't do anything here
			}
			// we need to tell him to hit it
			wait.ch <- true
		})
	}
}

func getService(m *net.GetServiceRequest, respChan chan *anypb.Any) error {
	pkg := m.GetPackagePath()
	svc := m.GetService()

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

	a := anypb.Any{}
	if err := a.MarshalFrom(resp); err != nil {
		respChan <- nil
		return err
	}
	respChan <- &a
	return nil
}

func runBlock(m *net.RunBlockRequest, respChan chan *anypb.Any) error {
	waitLock.Lock()
	defer waitLock.Unlock()

	log.Printf("runblock: %s,%v ", m.GetAddr(), m.GetWaiter())
	key := sys.NewDepKeyFromAddr(m.GetAddr())
	info := &waitInfo{
		ch:        make(chan bool),
		waitStart: time.Now(),
		waitId:    waitCounter,
		waitKey:   key,
		respChan:  respChan,
	}
	waitCounter++

	if runBlockWaitingList[key.String()] != nil {
		log.Printf("RUNBLOCK Found addr %s in the list already, ignoring", m.GetAddr())
		respChan <- nil
		return nil
	}
	log.Printf("runblock: 2 checked wait list")
	if m.GetAddr() == "" {
		log.Printf("RUNBLOCK Run block called with empty address, ignoring")
		respChan <- nil
		return nil

	}
	// this creates the record in the waiting list so we can block on it with a different goroutine
	runBlockWaitingList[key.String()] = info
	go waitForReady(info)

	log.Printf("runblock: 3 run in if ready")
	// its possible that there is nothing to wait on
	core.RunIfReady(key, func(key dep.DepKey) {
		info, ok := runBlockWaitingList[key.String()]
		if !ok {
			log.Printf("RUNBLOCK unable to find their key in the waiting list, even though core says ready to run")
			log.Printf("RUNBLOCK: waiting list has %d entries", len(runBlockWaitingList))
			for k, v := range runBlockWaitingList {
				log.Printf("\t%s:%#v", k, v)
			}
			log.Printf("RUNBLOCK: waiting list %#v", runBlockWaitingList)
		}
		info.ch <- true
	})

	log.Printf("runblock: 4 done....")
	return nil
}

func require(m *net.RequireRequest, respChan chan *anypb.Any) error {
	log.Printf("require reached start of func ")
	waitLock.Lock()
	defer waitLock.Unlock()
	log.Printf("require after lock ")

	for _, require := range m.GetRequire() {
		pkg := require.GetPackagePath()
		svc := require.GetService()
		key := sys.NewDepKeyFromAddr(require.GetAddr())

		log.Printf("require, about to hit nscore require %s", key)
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
	log.Printf("require done, sending response ")
	a := anypb.Any{}
	if err := a.MarshalFrom(resp); err != nil {
		respChan <- nil
		return err
	}
	respChan <- &a
	return nil
}

// waitForReady runs on a different goroutine, it just waits for somebody to send a message
// through the channel info.ch
func waitForReady(info *waitInfo) {
	log.Printf("blocking client called waitForReady (id is %d)", info.waitId)
	t := <-info.ch
	log.Printf("blocking client finished waiting in waitForReady (id is %d)", info.waitId)

	resp := &net.RunBlockResponse{}
	waitLock.Lock()
	defer waitLock.Unlock()

	if !t {
		log.Printf("we've been told that we were timed out (id %d)", info.waitId)
		resp.ErrId = lib.MarshalKernelErrId(lib.NewKernelError(lib.KernelNotFound))
	} else {
		resp.ErrId = lib.MarshalKernelErrId(lib.NoKernelErr())
	}
	log.Printf("sending req resp to original RunBlock requestor")

	a := anypb.Any{}
	err := a.MarshalFrom(resp)
	if err != nil {
		info.respChan <- nil
	} else {
		info.respChan <- &a
	}

	// remove from the map before we return and release the lock
	_, ok := runBlockWaitingList[info.waitKey.String()]
	if !ok {
		log.Printf("unable to find the client info waiting for export %s (id %d) ", info.waitKey, info.waitId)
	} else {
		delete(runBlockWaitingList, info.waitKey.String())
	}
	// we've sent the client the response and removed them from the map, so we are done
}

func timeoutHandler() {
	for {
		time.Sleep(sleepAmount * time.Second)

		// note that we assert the lock up here because it is not safe to read the map requireWaitingList
		// with others modifying it
		waitLock.Lock()
		for _, info := range runBlockWaitingList {
			if time.Now().Sub(info.waitStart) > time.Duration(timeoutClient*time.Second) {
				info.ch <- false
				// note that the channel receiver does the work of removing himself from the requireWaitingList
			}
		}
		waitLock.Unlock()
	}
}
