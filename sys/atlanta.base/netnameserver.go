package sys

import (
	"context"
	"crypto/tls"
	"fmt"
	"hash/crc32"
	"os"
	"time"

	"github.com/iansmith/parigot/g/pb/ns"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"

	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/protobuf/types/known/anypb"
)

var netnameserverVerbose = true

type NetNameServer struct {
	*NSCore
	local        *LocalNameServer
	stream       quic.Stream
	remoteNSAddr string
	port         int
}

// servicePort is the port that a service listens on for incoming requests.
// This port is fixed because each service has its own container, thus its
// own udp portspace.
const servicePort = 13331

func NewNetNameserver(loc *LocalNameServer, addr string) *NetNameServer {
	stream, err := setupConnection(addr)
	if err != nil {
		panic("unable to connect to our network nameserver:" + err.Error())
	}
	return &NetNameServer{
		NSCore:       NewNSCore(),
		local:        loc,
		stream:       stream,
		remoteNSAddr: addr,
		port:         servicePort,
	}
}

var koopmanTable = crc32.MakeTable(crc32.Koopman)
var writeTimeout = 250 * time.Millisecond
var readTimeout = writeTimeout
var readBufferSize = 4096

// if your message doesn't start with this, you have lost sync and should close the connection
// so we can try to reconnect
var magicStringOfBytes = uint64(0x1789071417760704)
var frontMatterSize = 8 + 4
var trailerSize = 4

const parigotNSPort = 13330

func (n *NetNameServer) MakeRequest(any *anypb.Any) (*anypb.Any, lib.Id) {
	var err error
	if n.stream == nil {
		netnameserverPrint("MakeRequest ", "net nameserver xxx1 \n")
		n.stream, err = setupConnection(n.remoteNSAddr)
		if err != nil {
			// xxx swallow error
			return nil, lib.NewKernelError(lib.KernelNameserverFailed)
		}
	}
	return NSRoundTrip(any, n.stream)
}

func setupConnection(addr string) (quic.Stream, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-p"},
	}
	netnameserverPrint("setupConnection ", "XXXX NET NAMESERVER XXXX dialing %s ", addr)
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return nil, err
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	netnameserverPrint("setupConnection ", "XXXX NET NAMESERVER XXXX dial and open stream was successful to  %s", addr)
	return stream, nil
}

func (n *NetNameServer) Export(key dep.DepKey, packagePath, service string) lib.Id {
	host, err := os.Hostname()
	if err != nil {
		panic("unable to get my own hostname: " + err.Error())
	}
	expInfo := &ns.ExportInfo{
		PackagePath: packagePath,
		Service:     service,
		Addr:        fmt.Sprintf("%s:%d", host, n.port)}
	expReq := &ns.ExportRequest{
		Export: []*ns.ExportInfo{expInfo},
	}
	var any anypb.Any
	err = any.MarshalFrom(expReq)
	if err != nil {
		return lib.NewKernelError(lib.KernelMarshalFailed)
	}
	expResult, respErr := n.MakeRequest(&any)
	if respErr != nil {
		return respErr
	}
	netnameserverPrint("EXPORT ", "xxx export remote 4\n")
	expResp := ns.ExportResponse{}
	err = expResult.UnmarshalTo(&expResp)
	if err != nil {
		return lib.NewKernelError(lib.KernelUnmarshalFailed)
	}
	respOk := lib.UnmarshalKernelErrorId(expResp.KernelErr)
	if respOk.IsError() {
		return respOk
	}
	netnameserverPrint("EXPORT ", "xxx export remote 5\n")
	netnameserverPrint("EXPORT ", "result was %s", expResult.TypeUrl)
	return lib.NoKernelErr()
}

func (n *NetNameServer) CloseService(packagePath, service string) lib.Id {
	req := &ns.CloseServiceRequest{PackagePath: packagePath, Service: service}
	any := &anypb.Any{}
	err := any.MarshalFrom(req)
	netnameserverPrint("CLOSESERVICE ", "xxx export remote 1")
	if err != nil {
		return lib.NewKernelError(lib.KernelMarshalFailed)
	}
	result, kerr := n.MakeRequest(any)
	if kerr != nil {
		return kerr
	}
	netnameserverPrint("CLOSESERVICE ", " xxx 2")
	resp := ns.CloseServiceResponse{}
	err = result.UnmarshalTo(&resp)
	if err != nil {
		return lib.NewKernelError(lib.KernelUnmarshalFailed)
	}
	respErr := lib.UnmarshalKernelErrorId(resp.GetKernelErr())
	if respErr.IsError() {
		return respErr
	}
	netnameserverPrint("CLOSESERVICE ", "xxx export remote 3")
	netnameserverPrint("CLOSESERVICE ", "result was %s", result.TypeUrl)
	return lib.NoKernelErr()
}

func (n *NetNameServer) HandleMethod(p *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	panic("HandleMethod")
}

func (n *NetNameServer) RunNotify(key dep.DepKey) {
	panic("should never call run notify (process) on net nameserver")
}

func (n *NetNameServer) RunBlock(key dep.DepKey) bool {
	panic("should be calling the the network to inform the remote NS")
}

func (n *NetNameServer) RunIfReady(key dep.DepKey) {
	panic("we need to talk to the network to do this")
}

// StartFailedInfo is supposed to return details about why the
// startup failed (e.g. a loop of dependencies). For now, we don't
// have a way to calculate this in the network case.
func (n *NetNameServer) StartFailedInfo() string {
	return n.NSCore.StartFailedInfo()
}

func netnameserverPrint(method, spec string, arg ...interface{}) {
	if netnameserverVerbose {
		part1 := fmt.Sprintf("NetNameServer:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
