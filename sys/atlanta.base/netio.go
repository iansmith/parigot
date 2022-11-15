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

	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const nsioVerbose = true

var ParigotProtoNameServer = []string{"quic-parigot-ns"}
var ParigotProtoRPC = []string{"quic-parigot-rpc"}

func NetSend(any *anypb.Any, stream quic.Stream) error {
	buf, err := proto.Marshal(any)
	if err != nil {
		return err
	}
	netioPrint("NETSEND ", " outgoing type is %s", any.TypeUrl)
	full := make([]byte, frontMatterSize+trailerSize+len(buf))
	copy(full[frontMatterSize:frontMatterSize+len(buf)], buf[:])
	attention := uint64(magicStringOfBytes)
	binary.LittleEndian.PutUint64(full[0:8], attention)
	l := uint32(len(buf))
	binary.LittleEndian.PutUint32(full[8:12], l)
	result := crc32.Checksum(full[:frontMatterSize+len(buf)], koopmanTable)
	binary.LittleEndian.PutUint32(full[frontMatterSize+len(buf):], result)
	stream.SetWriteDeadline(time.Now().Add(writeTimeout))

	pos := 0
	for pos < len(full)-frontMatterSize-trailerSize {
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
	readBuffer := make([]byte, readBufferSize)
	stream.SetReadDeadline(time.Now().Add(timeout))
	pos := 0
	for pos < frontMatterSize {
		read, err := stream.Read(readBuffer[pos:frontMatterSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	netioPrint("NETRECEIVE ", "read start buffer of size %d", pos)
	val := binary.LittleEndian.Uint64(readBuffer[0:8])
	if val != magicStringOfBytes {
		return nil, fmt.Errorf("unexpected prefix on bundle, must have lost sync")
	}
	sentLen := binary.LittleEndian.Uint32(readBuffer[8:12])
	if sentLen > uint32(readBufferSize) {
		return nil, fmt.Errorf("read packet was of size %d, but only had %d bytes to store the data", sentLen, readBufferSize)
	}
	netioPrint("NETRECEIVE ", "read size info %d", sentLen)
	for pos < frontMatterSize+trailerSize+int(sentLen) {
		read, err := stream.Read(readBuffer[pos : frontMatterSize+int(sentLen)+trailerSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	netioPrint("NETRECEIVE ", "completed read with size of  %d", pos)

	a := anypb.Any{}
	err := proto.Unmarshal(readBuffer[frontMatterSize:frontMatterSize+int(sentLen)], &a)
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
