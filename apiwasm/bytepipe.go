package apiwasm

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	pcontext "github.com/iansmith/parigot/context"
	"google.golang.org/protobuf/proto"
)

type BytePipeIn[T proto.Message] struct {
	rd       io.Reader
	syncLost bool
	ctx      context.Context
}

const maxProtobufSizeInBytes = 4 * 4096

var ErrTooLarge = errors.New("input to read is to large")
var ErrSyncLost = errors.New("unable to find message boundaries")
var ErrPipeClosed = errors.New("pipe is closed")
var ErrUnmarshal = errors.New("unable to unmarshal msg from input ")
var ErrMarshal = errors.New("unable to marshal msg for output ")
var ErrTimeout = errors.New("unable to read next byte of input, timeout expired")
var ErrUnexpectNum = errors.New("byte read is not a hex digit")
var ErrSignalExit = errors.New("you should shut down")

type nilableProto interface {
	proto.Message
	comparable
}

// NewBytePipeIn creates a new bytePipeIn that reads on the given reader.
// NewBytePipeIn creates a goroutine so that the rest of the bytePipeIn
// can use channels to read the bytes and do timeouts.
func NewBytePipeIn[T nilableProto](ctx context.Context, rd io.Reader) *BytePipeIn[T] {
	bpi := &BytePipeIn[T]{rd: rd, ctx: ctx}
	return bpi
}

func (b *BytePipeIn[T]) ReadProto(msg T, errIdPtr *int32) error {
	if b.syncLost {
		b.syncLost = false
		b.rd = nil
		return ErrSyncLost
	}
	if b.rd == nil {
		return ErrPipeClosed
	}

	var err error
	sizeBuf := make([]byte, 4)
	if err := readConstantSize(b.ctx, b.rd, 5, sizeBuf); err != nil {
		return err
	}

	ctx := pcontext.CallTo(b.ctx, "toInt")
	size, err := b.toInt(ctx, sizeBuf)
	if err != nil {
		return b.lostSync(err)
	}
	log.Printf("-----> size read was %d", size)
	if size == 0xffff {
		return ErrSignalExit
	}

	isErr := false
	if size&0x8000 != 0 {
		isErr = true
		size &= 0x7fff
		msg = *new(T) //compilicated way of saying nil
	}

	if size >= maxProtobufSizeInBytes {
		return b.lostSync(ErrTooLarge)
	}
	c := make([]byte, 1)
	_, err = b.rd.Read(c)
	if err != nil {
		return err
	}
	if c[0] != 32 {
		return b.lostSync(err)
	}
	if size == 0 {
		return nil
	}
	pcontext.Infof(b.ctx, "reading next payload, size is %d", size)

	result := make([]byte, size)
	if isErr {
		rez := binary.LittleEndian.Uint32(result)
		*errIdPtr = int32(rez)
		return nil
	}
	err = proto.Unmarshal(result, msg)
	if err != nil {
		return ErrUnmarshal
	}
	return nil
}

func (b *BytePipeIn[T]) toInt(ctx context.Context, sizeBuf []byte) (int, error) {
	total := 0
	for i := 0; i < 4; i++ {
		curr := sizeBuf[i]
		var val int
		switch curr {
		case 0x30:
			val = 0
		case 0x31:
			val = 1
		case 0x32:
			val = 2
		case 0x33:
			val = 3
		case 0x34:
			val = 4
		case 0x35:
			val = 5
		case 0x36:
			val = 6
		case 0x37:
			val = 7
		case 0x38:
			val = 8
		case 0x39:
			val = 9
		case 0x41, 0x61:
			val = 10
		case 0x42, 0x62:
			val = 11
		case 0x43, 0x63:
			val = 12
		case 0x44, 0x64:
			val = 13
		case 0x45, 0x65:
			val = 14
		case 0x46, 0x66:
			val = 15
		default:
			return 0, ErrUnexpectNum
		}
		// add this hex digit to total
		total += val << ((3 - i) * 4)

	}
	return total, nil
}

func (b *BytePipeIn[T]) lostSync(err error) error {
	b.syncLost = true
	b.rd = nil
	if err != nil {
		return err
	}
	return ErrSyncLost
}

// NewBytePipeOut creates a new BytePipeOut that writes on the given writer.
func NewBytePipeOut[T nilableProto](ctx context.Context, wr io.Writer) *BytePipeOut[T] {
	bpi := &BytePipeOut[T]{wr: wr, ctx: ctx}
	return bpi
}

type BytePipeOut[T nilableProto] struct {
	wr  io.Writer
	ctx context.Context
}

func (b *BytePipeOut[T]) WriteProto(resp T, err int32) error {
	if resp == *new(T) { // test for nil
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(err))
		if err := writeConstantSize(b.ctx, b.wr, true, 4, buf); err != nil {
			return err
		}
		return nil
	}
	buf, merr := proto.Marshal(resp)
	if merr != nil {
		return merr
	}
	if err := writeConstantSize(b.ctx, b.wr, false, int32(len(buf)), buf); err != nil {
		return err
	}
	return nil
}

func readConstantSize(ctx context.Context, rd io.Reader, size int32, buf []byte) error {
	count := int32(0)
	for count < size {
		print("xxxx read constant size %d, bytes %+x", size, buf)
		rd, err := rd.Read(buf[count:])
		if err != nil {
			return err
		}
		count += int32(rd)
	}
	return nil
}

func writeConstantSize(ctx context.Context, wr io.Writer, errorBit bool, size int32, buf []byte) error {
	sizeInBytes := size
	if errorBit {
		sizeInBytes |= 0x8000
	}
	str := fmt.Sprintf("%04x ", sizeInBytes)

	count := int32(0)
	for count < size {
		n, err := wr.Write([]byte(str)[count:])
		if err != nil {
			return err
		}
		count += int32(n)
	}
	return nil
}
