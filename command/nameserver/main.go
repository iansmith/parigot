package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash/crc32"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/iansmith/parigot/g/pb/ns"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys"
	"github.com/iansmith/parigot/sys/dep"

	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// XXX these are duplicates of the values in the netnameserver.go code that is the client
var koopmanTable = crc32.MakeTable(crc32.Koopman)
var writeTimeout = 250 * time.Millisecond
var readTimeout = writeTimeout
var readBufferSize = 4096
var magicStringOfBytes = uint64(0x1789071417760704)

const frontMatter = 12
const trailer = 4

const parigotNSPort = 13330

const timeoutClient = 6 // amount of time after which we decide your require aint gonna happen
const sleepAmount = 2   //two secs between checks for timeout

var core = sys.NewNSCore(false)

type waitInfo struct {
	ch        chan bool
	stream    quic.Stream
	waitStart time.Time
	waitId    int
	// we keep this just we can reverse the mapping without walking
	waitKey dep.DepKey
}

var waitCounter = 0
var waitLock sync.Mutex
var typeToHostMap = make(map[string]string)
var runBlockWaitingList = make(map[string] /*dep.DepKey*/ *waitInfo)

type waitPair struct {
	waitingFor string
	waitInfo   *waitInfo
}

func main() {
	go timeoutHandler()
	addr := fmt.Sprintf("0.0.0.0:%d", parigotNSPort)

	log.Printf("main: waiting on client connection...")
	for {
		listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
		if err != nil {
			// this failure occurs regulararly because of timeouts in the ListenAddr
			continue
		}

		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("main: error in Accept:%v", err)
			continue
		}
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("error in accept stream: %v", err)
			continue
		}
		go singleClientLoop(stream)
	}
}

func singleClientLoop(stream quic.Stream) {
	for {
		msg, err := sys.NSReceive(stream, readTimeout)
		if err != nil {
			log.Printf("singleClientLoop:error in read: %v, closing stream", err)
			stream.Close()
			return
		}
		log.Printf("singleClientLoop: received bundle of type %s", msg.TypeUrl)
		p, err := msg.UnmarshalNew()
		if err != nil {
			log.Printf("error in trying to unmarshal result of readStream: %v, closing stream", err)
			stream.Close()
			return
		}

		switch m := p.(type) {
		case *ns.CloseServiceRequest:
			err = closeService(m, stream)
		case *ns.ExportRequest:
			err = export(m, stream)
		case *ns.LocateRequest:
			err = locate(m, stream)
		case *ns.RequireRequest:
			err = require(m, stream)
		case *ns.RunBlockRequest:
			err = runBlock(m, stream)
		default:
			panic(fmt.Sprintf("nameserver received a bundle from a client that it could not understand the type of:%T", p))
		}
		if err != nil {
			log.Printf("got error from client stream: %v, closing stream", err)
			// close the stream to force a retry
			stream.Close()
			return
		}
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-parigot-ns"},
	}
}

func closeService(m *ns.CloseServiceRequest, stream quic.Stream) error {
	log.Printf("xxx close service , making sure it exists in our graph")
	core.CloseService(m.GetPackagePath(), m.GetService())

	// at the moment, we really don't track this in the network case
	resp := &ns.CloseServiceResponse{
		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	log.Printf("xxx close service 3\n")
	return sendResponse(resp, stream)
}

func export(m *ns.ExportRequest, stream quic.Stream) error {
	// we lock here ONLY because the map requireWaitingList is not safe to read when others
	// are updating it
	waitLock.Lock()
	defer waitLock.Unlock()
	var failure lib.Id

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
		resp := &ns.ExportResponse{
			KernelErr: lib.MarshalKernelErrId(failure),
		}
		return sendResponse(resp, stream)
	}
	resp := &ns.ExportResponse{
		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	return sendResponse(resp, stream)
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

func locate(m *ns.LocateRequest, stream quic.Stream) error {
	pkg := m.GetPackagePath()
	svc := m.GetService()
	name := fmt.Sprintf("%s.%s", pkg, svc)
	addr, ok := typeToHostMap[name]
	if !ok {
		log.Printf("shouldn't have many locate() failures due to ordering of startup...%s", name)
	}
	resp := &ns.LocateResponse{
		Addr:      addr,
		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	return sendResponse(resp, stream)
}

func runBlock(m *ns.RunBlockRequest, stream quic.Stream) error {
	waitLock.Lock()
	defer waitLock.Unlock()

	log.Printf("runblock: %s,%v", m.GetAddr(), m.GetWaiter())
	key := sys.NewDepKeyFromAddr(m.GetAddr())
	info := &waitInfo{
		ch:        make(chan bool),
		stream:    stream,
		waitStart: time.Now(),
		waitId:    waitCounter,
		waitKey:   key,
	}
	waitCounter++
	if runBlockWaitingList[key.String()] != nil {
		log.Printf("RUNBLOCK Found addr %s in the list already, ignoring", m.GetAddr())
		return nil
	}
	log.Printf("runblock: 2 checked wait list")
	if m.GetAddr() == "" {
		log.Printf("RUNBLOCK Run block called with empty address, ignoring")
		return nil

	}
	// this creates the record in the waiting list so we can block on it with a different goroutine
	runBlockWaitingList[key.String()] = info
	go waitForExport(info)

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
			return
		}
		info.ch <- true
	})

	log.Printf("runblock: 4 done")
	return nil
}

func require(m *ns.RequireRequest, stream quic.Stream) error {
	waitLock.Lock()
	defer waitLock.Unlock()

	for _, require := range m.GetRequire() {
		pkg := require.GetPackagePath()
		svc := require.GetService()
		key := sys.NewDepKeyFromAddr(require.GetAddr())

		// where should we put this?
		if id := core.Require(key, pkg, svc); id != nil {
			return lib.NewPerrorFromId("require failed", id)
		}
	}

	// tell them it's ok
	resp := &ns.RequireResponse{
		KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
	}
	return sendResponse(resp, stream)
}

func sendResponse(resp proto.Message, stream quic.Stream) error {
	var any anypb.Any
	err := any.MarshalFrom(resp)
	if err != nil {
		return err
	}
	kerr := sys.NSSend(&any, stream)
	if kerr.IsError() {
		return lib.NewPerrorFromId("failed to send response to client", kerr)
	}
	return nil
}

// waitForExport runs on a different goroutine, it just waits for somebody to send a message
// through the channel info.ch
func waitForExport(info *waitInfo) {
	log.Printf("blocking client called WaitForExport")
	t := <-info.ch
	log.Printf("blocking client finished waiting in WaitForExport")

	resp := &ns.RequireResponse{}
	waitLock.Lock()
	defer waitLock.Unlock()

	if !t {
		log.Printf("we've been told that we were timed out (id %d)", info.waitId)
		resp.KernelErr = lib.MarshalKernelErrId(lib.NewKernelError(lib.KernelNotFound))
	} else {
		resp.KernelErr = lib.MarshalKernelErrId(lib.NoKernelErr())
	}
	sendResponse(resp, info.stream)

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
