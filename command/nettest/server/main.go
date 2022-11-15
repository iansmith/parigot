package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	rand2 "math/rand"
	"time"

	"github.com/lucas-clemente/quic-go"
)

var proto = []string{"quic-parigot-test"}
var port = 12345
var earlyPort = 12346

var conn quic.Connection

func main() {
	allInterfaces := fmt.Sprintf("0.0.0.0:%d", port)
	allInterfacesEarly := fmt.Sprintf("0.0.0.0:%d", earlyPort)
	statelessKey := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	conf := &quic.Config{
		MaxIdleTimeout:    10 * time.Minute,
		StatelessResetKey: statelessKey,
	}
	earlyListener, err := quic.ListenAddrEarly(allInterfacesEarly, generateTLSConfig(proto), conf)
	if err != nil {
		log.Fatalf("failed to early listen on addr %s: %v, %T", allInterfaces, err, err)
	}
	go func() {
		for {
			earlyConn, err := earlyListener.Accept(context.Background())
			if err != nil {
				log.Fatalf("earlyListener accept failed: %v,%T", err, err)
			}
			log.Printf("early conn was ok")
			earlyStream, err := earlyConn.AcceptStream(context.Background())
			if err != nil {
				log.Fatalf("earlyConn accept stream failed: %v,%T", err, err)
			}
			log.Printf("early stream was ok")
			go func() {
				buf := make([]byte, 8)
				r, err := earlyStream.Read(buf)
				if err != nil {
					log.Fatalf("unable to read in early stream %v,%T", err, err)
				}
				if r != 8 {
					log.Fatalf("short read in early stream")
				}
				log.Printf("Server: got a packet %x", binary.LittleEndian.Uint64(buf))
			}()
		}
	}()

	listener, err := quic.ListenAddr(allInterfaces, generateTLSConfig(proto), conf)
	if err != nil {
		log.Fatalf("failed to listen on addr %s: %v, %T", allInterfaces, err, err)
	}
	conn, err = listener.Accept(context.Background())
	if err != nil {
		log.Fatalf("failed to accept on addr %s: %v, %T", allInterfaces, err, err)
	}
	log.Printf("connection accepted: %v, %T", conn, conn)

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		log.Fatalf("failed to accept stream on addr %s: %v, %T", allInterfaces, err, err)
	}
	log.Printf("stream accepted: %v, %T--%s,%s", stream, stream, conn.LocalAddr().String(), conn.RemoteAddr().String())
	c1 := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		c1 <- true
	}()
	for {
		select {
		case res := <-c1:
			stream.Close()
			conn.CloseWithError(quic.ApplicationErrorCode(0x101010), "migration")
			log.Printf("closed stream after %v", res)
			continue
		default:
			processStream(stream)
		}
	}
}

func processStream(stream quic.Stream) {
	b := make([]byte, 8)
	r, err := stream.Read(b)
	if err != nil {
		log.Fatalf("error in read:%v,%T", err, err)
	}
	if r != 8 {
		log.Fatalf("short read: %d", r)
	}
	x := binary.LittleEndian.Uint64(b)
	log.Printf("Server: received %x on stream", x)
	x = rand2.Uint64()
	log.Printf("Server: Sending '%x'", x)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, x)
	w, err := stream.Write(buf)
	if err != nil {
		log.Fatalf("failed on write: %v,%T", err, err)
	}
	if w != 8 {
		log.Fatalf("short write %d", w)
	}
	return
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig(proto []string) *tls.Config {
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
		NextProtos:   []string{"quic-parigot-test"},
	}
}
