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

var core = sys.NewNSCore()

type waitInfo struct {
	waitFor   string
	ch        chan bool
	stream    quic.Stream
	waitStart time.Time
	waitId    int
}

var waitCounter = 0
var waitLock sync.Mutex
var typeToHostMap = make(map[string]string)
var requireWaitingList = make(map[string][]*waitInfo)

type waitPair struct {
	waitingFor string
	waitInfo   *waitInfo
}

func main() {
	go timeoutHandler()

	addr := fmt.Sprintf("0.0.0.0:%d", parigotNSPort)
	// right now this is not concurrent, but it should be
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		panic(fmt.Sprintf("error in ListenAddr:%v", err))
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("error in Accept:%v", err)
			continue
		}
		log.Printf("xxx1 accept stream")
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("error in accept stream: %v", err)
			continue
		}
		log.Printf("xxx2 read stream bundle")
		msg, err := sys.NSReceive(stream)
		if err != nil {
			log.Printf("error in read: %v", err)
			continue
		}
		log.Printf("xxx2+ unmarshal new of any")
		p, err := msg.UnmarshalNew()
		if err != nil {
			log.Printf("error in trying to unmarshal result of readStream: %v", err)
			continue
		}

		// if there is a long wait, we will use the channel and queue above
		log.Printf("xxx3 switch on type of p")

		switch m := p.(type) {
		case *ns.CloseServiceRequest:
			err = closeService(m, stream)
		case *ns.ExportRequest:
			err = export(m, stream)
		case *ns.LocateRequest:
			err = locate(m, stream)
		case *ns.RequireRequest:
			err = require(m, stream)
		default:
			panic(fmt.Sprintf("nameserver received a bundle from a client that it could not understand the type of:%T", p))
		}
		if err != nil {
			// close the stream to force a retry
			stream.Close()
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
	log.Printf("xxx close service 1")
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

	for _, export := range m.GetExport() {
		pkg := export.GetPackagePath()
		svc := export.GetService()
		name := fmt.Sprintf("%s.%s", pkg, svc)
		typeToHostMap[name] = export.GetAddr()

		// if anybody is waiting on this, we better notify them
		list := requireWaitingList[name]
		for _, info := range list {
			info.ch <- true
			// note that the channel receiver does the work of removing himself from the requireWaitingList
		}
		// now can tell the client we are cool
		resp := &ns.ExportResponse{
			KernelErr: lib.MarshalKernelErrId(lib.NoKernelErr()),
		}
		err := sendResponse(resp, stream)
		if err != nil {
			return err
		}
	}
	return nil
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

// func runBlock(m *ns.RunBlockRequest, stream quic.Stream) error {
// 	if !ok {
// 		waitLock.Lock()
// 		defer waitLock.Unlock()
// 		info := &waitInfo{
// 			ch:        make(chan bool),
// 			stream:    stream,
// 			waitStart: time.Now(),
// 			waitFor:   name,
// 			waitId:    waitCounter,
// 		}
// 		waitCounter++
// 		list, ok := requireWaitingList[name]
// 		if !ok {
// 			list = []*waitInfo{}
// 			requireWaitingList[name] = list
// 		}
// 		list = append(list, info)
// 		// this is ok because it runs on another goroutine and thus we
// 		// release the lock on this goroutine with this return
// 		go waitForExport(info)
// 		return nil
// 	}

// }

func require(m *ns.RequireRequest, stream quic.Stream) error {
	waitLock.Lock()
	defer waitLock.Unlock()

	for _, require := range m.GetRequire() {
		pkg := require.GetPackagePath()
		svc := require.GetService()
		//name := fmt.Sprintf("%s.%s", pkg, svc)

		// where should we put this?
		if id := core.Require(nil, pkg, svc); id != nil {
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

// waitForExport runs on a different goroutine
func waitForExport(info *waitInfo) {
	log.Printf("blocking client who required %s", info.waitFor)

	t := <-info.ch
	resp := &ns.RequireResponse{}
	waitLock.Lock()
	defer waitLock.Unlock()

	if !t {
		log.Printf("telling client waiting on %s that he timed out (id %d)", info.waitFor, info.waitId)
		resp.KernelErr = lib.MarshalKernelErrId(lib.NewKernelError(lib.KernelNotFound))
	} else {
		resp.KernelErr = lib.MarshalKernelErrId(lib.NoKernelErr())
	}
	sendResponse(resp, info.stream)

	// remove from the map before we return and release the lock
	list := requireWaitingList[info.waitFor]
	if list == nil {
		log.Printf("unable to find the client info waiting for export %s (id %d) ", info.waitFor, info.waitId)
	} else {
		if len(list) == 1 {
			delete(requireWaitingList, info.waitFor)
		} else {
			result := []*waitInfo{}
			for i := 0; i < len(list); i++ {
				if list[i].waitId == info.waitId {
					continue
				}
				result = append(result, list[i])
			}
			requireWaitingList[info.waitFor] = result
		}
	}
	// we've sent the client the response and removed them from the map, so we are done
}

func timeoutHandler() {
	for {
		time.Sleep(sleepAmount * time.Second)

		// note that we assert the lock up here because it is not safe to read the map requireWaitingList
		// with others modifying it
		waitLock.Lock()
		for _, list := range requireWaitingList {
			for _, info := range list {
				if time.Now().Sub(info.waitStart) > time.Duration(timeoutClient*time.Second) {
					info.ch <- false
					// note that the channel receiver does the work of removing himself from the requireWaitingList
				}
			}
		}
		waitLock.Unlock()
	}
}
