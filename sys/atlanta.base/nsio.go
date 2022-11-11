package sys

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/iansmith/parigot/lib"

	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const nsioVerbose = true

func NSSend(any *anypb.Any, stream quic.Stream) lib.Id {
	buf, err := proto.Marshal(any)
	if err != nil {
		return lib.NewKernelError(lib.KernelMarshalFailed)
	}
	nsioPrint("NSSEND ", " outgoing type is %s", any.TypeUrl)
	full := make([]byte, frontMatterSize+trailerSize+len(buf))
	copy(full[frontMatterSize:frontMatterSize+len(buf)], buf[:])
	attention := uint64(magicStringOfBytes)
	binary.LittleEndian.PutUint64(full[0:8], attention)
	l := uint32(len(buf))
	binary.LittleEndian.PutUint32(full[8:12], l)
	result := crc32.Checksum(full[:frontMatterSize+len(buf)], koopmanTable)
	binary.LittleEndian.PutUint32(full[frontMatterSize+len(buf):], result)
	stream.SetWriteDeadline(time.Now().Add(writeTimeout))

	// xxx fixme: we are swallowing the error details here because we may not have a way to log them
	pos := 0
	for pos < len(full)-frontMatterSize-trailerSize {
		written, err := stream.Write(full[pos:])
		if err != nil {
			lib.NewKernelError(lib.KernelNameserverLost)
		}
		pos += written
	}
	nsioPrint("NSSend ", "wrote buffer of size %d", pos)
	return lib.NoKernelErr()
}

func NSReceive(stream quic.Stream, timeout time.Duration) (*anypb.Any, error) {
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
	nsioPrint("NSReceive ", "read start buffer of size %d", pos)
	val := binary.LittleEndian.Uint64(readBuffer[0:8])
	if val != magicStringOfBytes {
		return nil, fmt.Errorf("unexpected prefix on bundle, must have lost sync")
	}
	sentLen := binary.LittleEndian.Uint32(readBuffer[8:12])
	if sentLen > uint32(readBufferSize) {
		return nil, fmt.Errorf("read packet was of size %d, but only had %d bytes to store the data", sentLen, readBufferSize)
	}
	nsioPrint("NSReceive ", "read size info %d", sentLen)
	for pos < frontMatterSize+trailerSize+int(sentLen) {
		read, err := stream.Read(readBuffer[pos : frontMatterSize+int(sentLen)+trailerSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	nsioPrint("NSReceive ", "completed read with size of  %d", pos)

	a := anypb.Any{}
	err := proto.Unmarshal(readBuffer[frontMatterSize:frontMatterSize+int(sentLen)], &a)
	if err != nil {
		return nil, err
	}
	nsioPrint("NSReceive ", "successfully read a value from client %s", a.TypeUrl)
	return &a, nil
}

func NSRoundTrip(any *anypb.Any, stream quic.Stream, timeout time.Duration) (*anypb.Any, lib.Id) {
	nsioPrint("NSRoundTrip ", "outgoing value of type %s with timeout %d", any.TypeUrl, timeout)
	kerr := NSSend(any, stream)
	if kerr != nil && kerr.IsError() {
		return nil, kerr
	}
	nsioPrint("NSRound Trip ", "send completed with type %s", any.TypeUrl)
	result, err := NSReceive(stream, timeout)
	if err != nil {
		print("error from (recv following send of ", any.TypeUrl, "): ", err.Error(), "\n")
		return nil, lib.NewKernelError(lib.KernelNameserverFailed)
	}
	nsioPrint("NSRoundTrip DONE! read value of type", result.TypeUrl)
	return result, nil
}

func nsioPrint(method, spec string, arg ...interface{}) {
	if nsioVerbose {
		part1 := fmt.Sprintf("NSIO:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
