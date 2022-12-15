package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/iansmith/parigot/api/netconst"
	pb "github.com/iansmith/parigot/api/proto/g/pb/log"

	"google.golang.org/protobuf/proto"
)

const socketDir = "PARIGOT_SOCKET_DIR"
const sockName = "logviewer.sock"

var envVar string
var sockAddr string

func main() {

	if os.Getenv(socketDir) == "" {
		log.Printf("Unable to find environment variable: %s", socketDir)
		log.Printf("This environment variable should point the _directory_ that contains or will contain ")
		log.Printf("the unix domain sockets.  This directory will be mapped to /var/run/parigot inside ")
		log.Printf("the dev container, if you are using it.  We recommend using a directory like")
		log.Printf("'~/parigot/socket' and we recommend using a fully qualified path.  The directory ")
		log.Printf("passed as %s  must already exist.", socketDir)
		os.Exit(1)
	}
	dir := os.Getenv(socketDir)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("The directory '%s' which is the value of the ", dir)
			log.Printf("environment variable '%s' does not exist.", socketDir)
			log.Printf("If this directory is the one you intended, you need to make sure the directory ")
			log.Printf("exists with 'mkdir' or similar before running a parigot command that ")
			log.Printf("uses the '%s' environment variable.", socketDir)
			os.Exit(1)
		}
		log.Fatalf("%v", err)
	}
	if !info.IsDir() {
		log.Printf("The value of the environment variable '%s' is not ", socketDir)
		log.Printf("a directory.  Make sure that the value of that variable both is a directory and exists.")
		os.Exit(1)
	}
	envVar = os.Getenv(socketDir)
	sockAddr = filepath.Join([]string{envVar, sockName}...)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cleanupAndExit()
		}
	}()

	for {
		l, e := net.Listen("unix", sockAddr)
		if e != nil {
			log.Printf("%v", e)
			cleanupAndExit()
		}
		conn, e := l.Accept()
		if e != nil {
			log.Printf("%v", e)
			cleanupAndExit()
		}
		go handler(conn)
	}

}

func handler(conn net.Conn) {
	dataBuffer := make([]byte, netconst.FrontMatterSize+netconst.TrailerSize+netconst.ReadBufferSize)
	magicBuffer := dataBuffer[0:8]
	lenBuffer := dataBuffer[8:12]
	crcBuffer := dataBuffer[len(dataBuffer)-4:]

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
		read, err = conn.Read(magicBuffer)
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
		var req pb.LogRequest
		err = proto.Unmarshal(dataBuffer[:length], &req)
		if err != nil {
			internalMessage(fmt.Sprintf("[unable to unmarshal data from socket: %v, closing %s]", err, sockAddr))
			return
		}
		result := crc32.Checksum(dataBuffer[:netconst.FrontMatterSize+netconst.ReadBufferSize],
			netconst.KoopmanTable)
		expected := binary.LittleEndian.Uint32(crcBuffer)
		if expected != result {
			internalMessage(fmt.Sprintf("[bad crc: expected %x but got %x, closing %s]", expected, result, sockAddr))
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
	log.Printf("removing unix domain socket: %s", sockAddr)
	err := os.Remove(sockAddr)
	if err != nil {
		log.Printf("unable to remove our unix domain socket %s:%v", sockAddr, err)
	}
	os.Exit(1)
}
