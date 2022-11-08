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
	print("NSSend xxx5 \n")
	full := make([]byte, frontMatterSize+trailerSize+len(buf))
	copy(full[frontMatterSize:frontMatterSize+len(buf)], buf[:])
	attention := uint64(magicStringOfBytes)
	binary.LittleEndian.PutUint64(full[0:8], attention)
	l := uint32(len(buf))
	binary.LittleEndian.PutUint32(full[8:12], l)
	result := crc32.Checksum(full[:frontMatterSize+len(buf)], koopmanTable)
	binary.LittleEndian.PutUint32(full[frontMatterSize+len(buf):], result)
	print("NSSend xxx6 \n")
	stream.SetWriteDeadline(time.Now().Add(writeTimeout))
	print("NSSend xxx7 \n")

	// xxx fixme: we are swallowing the error details here because we may not have a way to log them
	pos := 0
	for pos < len(full)-frontMatterSize-trailerSize {
		print("NSSend xxx8 ", pos, "\n")
		written, err := stream.Write(full[pos:])
		if err != nil {
			lib.NewKernelError(lib.KernelNameserverLost)
		}
		pos += written
	}
	print("NSSend xxx9\n")
	return lib.NoKernelErr()
}

func NSReceive(stream quic.Stream) (*anypb.Any, error) {
	print("NSReceive xxx3 \n")
	readBuffer := make([]byte, readBufferSize)
	stream.SetReadDeadline(time.Now().Add(readTimeout))
	pos := 0
	for pos < frontMatterSize {
		print("NSReceive xxx4 ", pos, "\n")
		read, err := stream.Read(readBuffer[pos:frontMatterSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	print("NSReceive xxx4a", pos, "\n")
	val := binary.LittleEndian.Uint64(readBuffer[0:8])
	if val != magicStringOfBytes {
		return nil, fmt.Errorf("unexpected prefix on bundle, must have lost sync")
	}
	sentLen := binary.LittleEndian.Uint32(readBuffer[8:12])
	if sentLen > uint32(readBufferSize) {
		return nil, fmt.Errorf("read packet was of size %d, but only had %d bytes to store the data", sentLen, readBufferSize)
	}
	print("NSReceive xxx5 ", sentLen, "\n")
	for pos < frontMatterSize+trailerSize+int(sentLen) {
		print("NSReceive xxx5a ", pos, "\n")
		read, err := stream.Read(readBuffer[pos : frontMatterSize+int(sentLen)+trailerSize])
		if err != nil {
			return nil, err
		}
		pos += read
	}
	print("NSReceive xxx7 ", pos, "\n")

	a := anypb.Any{}
	err := proto.Unmarshal(readBuffer[frontMatterSize:frontMatterSize+int(sentLen)], &a)
	if err != nil {
		return nil, err
	}
	print("NSReceive xxx8 successfully read a value from client ", a.TypeUrl, "\n")
	return &a, nil
}

func NSRoundTrip(any *anypb.Any, stream quic.Stream) (*anypb.Any, lib.Id) {
	print("net NSRoundTrip xxx3 \n")
	kerr := NSSend(any, stream)
	if kerr.IsError() {
		return nil, kerr
	}
	print("net NSRoundTrip xxx4\n")
	result, err := NSReceive(stream)
	print("net NSRoundTrip xxx5 \n")
	if err != nil {
		print("error from receive ", err.Error(), "\n")
		return nil, lib.NewKernelError(lib.KernelNameserverFailed)
	}
	print("net NSRoundTrip xxx6 \n")
	return result, nil
}
