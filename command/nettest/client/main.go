package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
)

var addr = "localhost:12345"
var earlyAddr = "localhost:12346"

type BufferedWriteCloser struct {
	*bufio.Writer
	*os.File
}

func (b *BufferedWriteCloser) Close() error {
	b.Writer = nil
	return b.File.Close()
}
func (b *BufferedWriteCloser) Write(buf []byte) (int, error) {
	return b.Writer.Write(buf)
}

func NewBufferedWriteCloser(f *os.File) *BufferedWriteCloser {
	bwc := &BufferedWriteCloser{
		File:   f,
		Writer: bufio.NewWriter(f),
	}
	return bwc
}

func main() {

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-parigot-test"},
	}

	var qconf quic.Config
	qconf.Tracer = qlog.NewTracer(func(_ logging.Perspective, connID []byte) io.WriteCloser {
		filename := fmt.Sprintf("client_%x.qlog", connID)
		f, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Creating qlog file %s.\n", filename)
		return NewBufferedWriteCloser(f)
	})

	tokenStore := quic.NewLRUTokenStore(4, 16)
	qconf.TokenStore = tokenStore

	conn, err := quic.DialAddr(addr, tlsConf, &qconf)
	if err != nil {
		log.Fatalf("failed to dial: %v, %T", err, err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalf("failed to open stream: %v, %T", err, err)
	}

	for {
		x := rand.Uint64()
		log.Printf("Client: Sending '%x'", x)
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, x)
		w, err := stream.Write(buf)
		if err != nil {
			log.Fatalf("failed on write: %v,%T", err, err)
		}
		if w != 8 {
			log.Fatalf("short write %d", w)
		}
		r, err := stream.Read(buf)
		if err != nil {
			log.Printf("failed on read: %v,%T", err, err)
			_, ok := err.(*quic.ApplicationError)
			if ok {
				conn, err := quic.DialAddrEarly(earlyAddr, tlsConf, &qconf)
				if err != nil {
					log.Fatalf("dial early failed: %v,%T", err, err)
				}
				stream, err := conn.OpenStreamSync(context.Background())
				if err != nil {
					log.Fatalf("unable to create stream: %v, %T", err, err)
				}
				binary.LittleEndian.PutUint64(buf, 0x0101010110101010)
				w, err := stream.Write(buf)
				if err != nil {
					log.Fatalf("failed to write bytes early: %v,%T", err, err)
				}
				if w != 8 {
					log.Fatalf("short early write")
				}
				os.Exit(6)
			}
		}
		if r != 8 {
			log.Fatalf("short read")
		}
		x = binary.LittleEndian.Uint64(buf)
		log.Printf("Client: Received '%x'", x)
		time.Sleep(1 * time.Second)
	}
}
