package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/iansmith/parigot/api/netconst"
	pb "github.com/iansmith/parigot/api/proto/g/pb/log"

	"google.golang.org/protobuf/proto"
)

var envVar string
var sockAddr string

func main() {

	if os.Getenv(netconst.SocketEnvVar) == "" {
		log.Printf("Unable to find environment variable: %s", netconst.SocketEnvVar)
		log.Printf("This environment variable should point the _directory_ that contains or will contain ")
		log.Printf("the unix domain sockets.  This directory will be mapped to /var/run/parigot inside ")
		log.Printf("the dev container, if you are using it.  We recommend using a directory like")
		log.Printf("'~/parigot/socket' and we recommend using a fully qualified path.  The directory ")
		log.Printf("passed as %s  must already exist.", netconst.SocketEnvVar)
		os.Exit(1)
	}
	dir := os.Getenv(netconst.SocketEnvVar)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("The directory '%s' which is the value of the ", dir)
			log.Printf("environment variable '%s' does not exist.", netconst.SocketEnvVar)
			log.Printf("If this directory is the one you intended, you need to make sure the directory ")
			log.Printf("exists with 'mkdir' or similar before running a parigot command that ")
			log.Printf("uses the '%s' environment variable.", netconst.SocketEnvVar)
			os.Exit(1)
		}
		log.Fatalf("%v", err)
	}
	if !info.IsDir() {
		log.Printf("The value of the environment variable '%s' is not ", netconst.SocketEnvVar)
		log.Printf("a directory.  Make sure that the value of that variable both is a directory and exists.")
		os.Exit(1)
	}
	envVar = os.Getenv(netconst.SocketEnvVar)
	sockAddr = ":4004"

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cleanupAndExit()
		}
	}()

	l, e := net.Listen("tcp", sockAddr)
	if e != nil {
		log.Printf("%v", e)
		cleanupAndExit()
	}
	internalMessage(fmt.Sprintf("waiting for connection on %s", sockAddr))
	for {
		conn, e := l.Accept()
		if e != nil {
			log.Printf("%v", e)
			cleanupAndExit()
		}
		internalMessage(fmt.Sprintf("handling new connection on %s", sockAddr))
		go handler(conn)
	}

}

func handler(conn net.Conn) {
	dataBuffer := make([]byte, netconst.FrontMatterSize+netconst.TrailerSize+netconst.ReadBufferSize)
	magicBuffer := dataBuffer[0:8]
	lenBuffer := dataBuffer[8:12]

	for {
		read, err := conn.Read(magicBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		if read != 8 {
			internalMessage(fmt.Sprintf("[bad size of read for magic number: %d expected 8, closing %s]", read, sockAddr))
			return
		}
		magicNum := binary.LittleEndian.Uint64(magicBuffer)
		if magicNum != netconst.MagicStringOfBytes {
			internalMessage(fmt.Sprintf("[bad magic number, got %x expected %x, closing %s]", magicNum, netconst.MagicStringOfBytes, sockAddr))
		}
		read, err = conn.Read(lenBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		if read != 4 {
			internalMessage(fmt.Sprintf("[bad size of read for length: %d expected 4, closing %s]", read, sockAddr))
			return
		}
		l := binary.LittleEndian.Uint32(lenBuffer)
		length := int(l)
		if length >= len(dataBuffer) {
			internalMessage(fmt.Sprintf("[data is too large: %d expected no more than %d, closing %s]", length, netconst.ReadBufferSize, sockAddr))
			return
		}
		objBuffer := dataBuffer[netconst.FrontMatterSize : length+netconst.FrontMatterSize]
		read, err = conn.Read(objBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		var req pb.LogRequest
		err = proto.Unmarshal(objBuffer, &req)
		if err != nil {
			internalMessage(fmt.Sprintf("[unable to unmarshal data from socket: %v, closing %s]", err, sockAddr))
			return
		}
		crcBuffer := dataBuffer[length+netconst.FrontMatterSize : length+netconst.FrontMatterSize+4]

		read, err = conn.Read(crcBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		result := crc32.Checksum(objBuffer, netconst.KoopmanTable)
		expected := binary.LittleEndian.Uint32(crcBuffer)
		if expected != result {
			internalMessage(fmt.Sprintf("[bad crc: expected %x but got %x, with CRC over %d bytes, closing %s]", expected, result, length, sockAddr))
			return
		}
		logMessage(&req)
	}
}

func internalMessage(s string) {
	log.Printf("xxx internal message %s", s)
}

func logMessage(req *pb.LogRequest) {
	s := fmt.Sprintf("%s:%d:%s", req.Stamp.AsTime().Format(time.RFC3339), req.Level, req.Message)
	log.Printf("xxx log request: %s", s)
}

func cleanupAndExit() {
	log.Printf("closing socket: %s", sockAddr)
	// err := os.Remove(sockAddr)
	// if err != nil {
	// 	log.Printf("unable to remove our unix domain socket %s:%v", sockAddr, err)
	// }
	os.Exit(1)
}
