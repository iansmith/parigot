package sys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"hash/crc32"
	"math/big"
	"time"

	"github.com/iansmith/parigot/apishared"

	quic "github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const nsioVerbose = false

var ParigotProtoNameServer = []string{"quic-parigot-ns"}
var ParigotProtoRPC = []string{"quic-parigot-rpc"}

func NetSend(any *anypb.Any, stream quic.Stream) error {
	buf, err := proto.Marshal(any)
	if err != nil {
		return err
	}
	netioPrint("NETSEND ", " outgoing type is %s", any.TypeUrl)
	full := make([]byte, apishared.FrontMatterSize+apishared.TrailerSize+len(buf))
	copy(full[apishared.FrontMatterSize:apishared.FrontMatterSize+len(buf)], buf[:])
	attention := uint64(apishared.MagicStringOfBytes)
	binary.LittleEndian.PutUint64(full[0:8], attention)
	l := uint32(len(buf))
	binary.LittleEndian.PutUint32(full[8:12], l)
	result := crc32.Checksum(full[:apishared.FrontMatterSize+len(buf)], apishared.KoopmanTable)
	binary.LittleEndian.PutUint32(full[apishared.FrontMatterSize+len(buf):], result)
	stream.SetWriteDeadline(time.Now().Add(apishared.WriteTimeout))

	pos := 0
	for pos < len(full)-apishared.FrontMatterSize-apishared.TrailerSize {
		written, err := stream.Write(full[pos:])
		if err != nil {
			return err
		}
		pos += written
	}
	netioPrint("NETSEND ", "wrote buffer of size %d", pos)
	return nil
}

func NetReceive(stream quic.Stream, timeout time.Duration) (*anypb.Any, error) {
	readBuffer := make([]byte, apishared.ReadBufferSize)
	stream.SetReadDeadline(time.Now().Add(timeout))
	pos := 0
	for pos < apishared.FrontMatterSize {
		read, err := stream.Read(readBuffer[pos:apishared.FrontMatterSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	netioPrint("NETRECEIVE ", "read start buffer of size %d", pos)
	val := binary.LittleEndian.Uint64(readBuffer[0:8])
	if val != apishared.MagicStringOfBytes {
		return nil, fmt.Errorf("unexpected prefix on bundle, must have lost sync")
	}
	sentLen := binary.LittleEndian.Uint32(readBuffer[8:12])
	if sentLen > uint32(apishared.ReadBufferSize) {
		return nil, fmt.Errorf("read packet was of size %d, but only had %d bytes to store the data", sentLen, apishared.ReadBufferSize)
	}
	netioPrint("NETRECEIVE ", "read size info %d", sentLen)
	for pos < apishared.FrontMatterSize+apishared.TrailerSize+int(sentLen) {
		read, err := stream.Read(readBuffer[pos : apishared.FrontMatterSize+int(sentLen)+apishared.TrailerSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	netioPrint("NETRECEIVE ", "completed read with size of  %d", pos)

	a := anypb.Any{}
	err := proto.Unmarshal(readBuffer[apishared.FrontMatterSize:apishared.FrontMatterSize+int(sentLen)], &a)
	if err != nil {
		return nil, err
	}
	netioPrint("NETRECEIVE ", "successfully read a value from client %s", a.TypeUrl)
	return &a, nil
}

func netioPrint(method, spec string, arg ...interface{}) {
	if nsioVerbose {
		part1 := fmt.Sprintf("NETIO:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}

// Setup a bare-bones TLS config for the server
func GenerateTLSConfig(protoName []string) *tls.Config {
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
		NextProtos:   protoName,
	}
}
