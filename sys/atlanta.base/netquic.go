package sys

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/iansmith/parigot/sys/dep"
	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	quiclistenerVerbose = true
	quiccallerVerbose   = true
)

type QuicListener struct {
	listener quic.Listener
	ch       chan *NetResult
}

type quicCaller struct {
	targetAddr string
	proto      []string
	ch         chan *NetResult
}

// NetResult has its direction depend on whether it is used for a client or server.
// In the server case (quicListener) the data is the data read from the network
// and the key is the remote address read from.  The resp channel is where the
// notified code should send the result.  The quicListener sends the NetResult
// to the channel (ch) provided at creation.
// quicCaller works the opposite way.  The quicCaller wants for netResults to be
// sent via the channel given at creation.  The data is the data to send to the
// remote server and the key is logically our address, but it is ignored.
// The resp channel is where the quicCaller should put the result read back.
// In both cases, the respCh should be sent nil to indicate a failure.
type NetResult struct {
	data   *anypb.Any
	respCh chan *anypb.Any
	key    dep.DepKey
}

func (n *NetResult) Data() *anypb.Any {
	return n.data
}
func (n *NetResult) RespChan() chan *anypb.Any {
	return n.respCh
}
func (n *NetResult) Key() dep.DepKey {
	return n.key
}
func (n *NetResult) SetData(a *anypb.Any) {
	n.data = a
}
func (n *NetResult) SetKey(key dep.DepKey) {
	n.key = key
}
func (n *NetResult) SetRespChan(ch chan *anypb.Any) {
	n.respCh = ch
}

// NewQuicListener establishes a listener on given port that is willing to
// speak the given protocol.  The received messages or errors are sent through
// the given channel.
func NewQuicListener(port int, proto []string, ch chan *NetResult) *QuicListener {
	allInterfaces := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := quic.ListenAddr(allInterfaces, GenerateTLSConfig(proto), nil)
	if err != nil {
		panic("cant establish a quic listener:" + err.Error() + " is the listener already running?")
	}
	ql := &QuicListener{listener: listener, ch: ch}
	go ql.waitForConnections()
	return ql
}

// waitForConnections tries to accept connections from remote CLIENTS
func (q *QuicListener) waitForConnections() {
	for {
		conn, err := q.listener.Accept(context.Background())
		if err != nil {
			quicListenerPrint("WAITFORCONN ", "error accepting connection:%v", err)
			continue
		}
		quicListenerPrint("WAITFORCONN ", "accepted connection from %s", conn.RemoteAddr().String())
		go q.waitForStreams(conn)
	}
}

// waitForStreams is in a loop so a client conn that gets closed can easily
// reestablish.
func (q *QuicListener) waitForStreams(conn quic.Connection) {
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			quicListenerPrint("WAITFORSTREAM", "error accepting stream: %v, closing connection", err)
			conn.CloseWithError(1, "accept stream error")
			return
		}
		go q.waitForRequests(stream, conn.RemoteAddr())
	}
}

// wait for requests actually picks up the bundle, converts to a callInfo
// and then pushes it through the channel
func (q *QuicListener) waitForRequests(stream quic.Stream, remote net.Addr) {
	for {
		a, err := NetReceive(stream, readTimeout)
		if err != nil {
			quicListenerPrint("WAITFORREQ", "error receiving: %v, closing stream", err)
			stream.Close()
			return
		}
		log.Printf("got to net result construction in listener: %s", a.TypeUrl)
		nr := NetResult{}
		nr.SetData(a)
		nr.SetKey(NewDepKeyFromAddr(remote.String()))
		nr.SetRespChan(make(chan *anypb.Any))
		log.Printf("sending through channel: %s, sending to %p", a.TypeUrl, q.ch)
		q.ch <- &nr
		log.Printf("blocking on response: %s", a.TypeUrl)
		out := <-nr.RespChan()
		if out == nil {
			log.Printf("got response after block, but it's a nil in response to %s", a.TypeUrl)
		} else {
			log.Printf("blocking completed, got response: %s in reponse to %s", out.TypeUrl, a.TypeUrl)
		}
		if out == nil {
			continue
		}
		err = NetSend(out, stream)
		if err != nil {
			quicListenerPrint("WAITFORREQ", "error sending response: %v, closing stream", err)
			stream.Close()
		}
	}
}

func newQuicCaller(addr string, proto []string, ch chan *NetResult) *quicCaller {
	q := &quicCaller{
		targetAddr: addr,
		proto:      proto,
		ch:         ch,
	}
	go q.start()
	return q
}

func (q *quicCaller) start() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         q.proto,
	}
	quicCallerPrint("CALL ", "dialing %s...", q.targetAddr)
	for {
		conn, err := quic.DialAddr(q.targetAddr, tlsConf, nil)
		if err != nil {
			panic("unable to establish connection the remote server:" + err.Error() + " are you sure it's running?")
		}
		q.useConnection(conn)
	}
}

func (q *quicCaller) useConnection(conn quic.Connection) {
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		netnameserverPrint("USENSSTREAM", "unable to establish stream: %v, will try to reconnect", err)
		return
	}
	for {
		a := <-q.ch
		err := NetSend(a.Data(), stream)
		if err != nil {
			a.RespChan() <- nil
			quicCallerPrint("USECONN ", "failed send data: %v", err)
			stream.Close()
			return
		}
		result, err := NetReceive(stream, longReadTimeout)
		if err != nil {
			a.RespChan() <- nil
			quicCallerPrint("USECONN ", "failed to receive data: %v", err)
			stream.Close()
			return
		}
		a.RespChan() <- result
	}
}

func quicListenerPrint(method, spec string, arg ...interface{}) {
	if quiclistenerVerbose {
		part1 := fmt.Sprintf("QUICLISTENER:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
func quicCallerPrint(method, spec string, arg ...interface{}) {
	if quiclistenerVerbose {
		part1 := fmt.Sprintf("QUICCALLER:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
