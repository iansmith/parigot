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

func NSSend(any *anypb.Any, stream quic.Stream) lib.Id {
	buf, err := proto.Marshal(any)
	if err != nil {
		return lib.NewKernelError(lib.KernelMarshalFailed)
	}
	print("NSSend, type is ", any.TypeUrl, "\n")
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
	print("NSSend wrote buffer of size ", pos, "\n")
	return lib.NoKernelErr()
}

func NSReceive(stream quic.Stream, timeout time.Duration) (*anypb.Any, error) {
	print("NSReceive 1 \n")
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
	print("NSReceive 2, read start buffer of size ", pos, "\n")
	val := binary.LittleEndian.Uint64(readBuffer[0:8])
	if val != magicStringOfBytes {
		return nil, fmt.Errorf("unexpected prefix on bundle, must have lost sync")
	}
	sentLen := binary.LittleEndian.Uint32(readBuffer[8:12])
	if sentLen > uint32(readBufferSize) {
		return nil, fmt.Errorf("read packet was of size %d, but only had %d bytes to store the data", sentLen, readBufferSize)
	}
	print("NSReceive 3 read size info ", sentLen, "\n")
	for pos < frontMatterSize+trailerSize+int(sentLen) {
		read, err := stream.Read(readBuffer[pos : frontMatterSize+int(sentLen)+trailerSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	print("NSReceive 4 completed read with size of  ", pos, "\n")

	a := anypb.Any{}
	err := proto.Unmarshal(readBuffer[frontMatterSize:frontMatterSize+int(sentLen)], &a)
	if err != nil {
		return nil, err
	}
	print("NSReceive 5 successfully read a value from client ", a.TypeUrl, "\n")
	return &a, nil
}

func NSRoundTrip(any *anypb.Any, stream quic.Stream, timeout time.Duration) (*anypb.Any, lib.Id) {
	print("NS Round Trip 1, outgoing value ", any.TypeUrl, " with timeout ", timeout, "\n")
	kerr := NSSend(any, stream)
	if kerr.IsError() {
		return nil, kerr
	}
	print("NS Round Trip 2 sent completed", any.TypeUrl, "\n")
	result, err := NSReceive(stream, timeout)
	if err != nil {
		print("error from (recv following send of ", any.TypeUrl, "): ", err.Error(), "\n")
		return nil, lib.NewKernelError(lib.KernelNameserverFailed)
	}
	print("NSRoundTrip 3 DONE! got read value", result.TypeUrl, "\n")
	return result, nil
}
