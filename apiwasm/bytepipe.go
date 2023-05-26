package apiwasm

import (
	"context"
	"errors"
	"io"

	pcontext "github.com/iansmith/parigot/context"
	"google.golang.org/protobuf/proto"
)

type BytePipeIn[T proto.Message] struct {
	rd       io.Reader
	ch       chan T
	syncLost bool
	ctx      context.Context
}

const maxProtobufSizeInBytes = 4 * 4096

var ErrTooLarge = errors.New("input to read is to large")
var ErrSyncLost = errors.New("unable to find message boundaries")
var ErrPipeClosed = errors.New("pipe is closed")
var ErrUnmarshal = errors.New("unable to unmarshal msg from input ")
var ErrTimeout = errors.New("unable to read next byte of input, timeout expired")
var ErrUnexpectNum = errors.New("byte read is not a hex digit")

var timeoutInMillis = 50

// NewBytePipeIn creates a new bytePipeIn that reads on the given reader.
// NewBytePipeIn creates a goroutine so that the rest of the bytePipeIn
// can use channels to read the bytes and do timeouts.
func NewBytePipeIn[T proto.Message](ctx context.Context, rd io.Reader) *BytePipeIn[T] {
	bpi := &BytePipeIn[T]{rd: rd, ctx: ctx, ch: make(chan T)}
	return bpi
}

func (b *BytePipeIn[T]) Chan() chan T {
	return b.ch
}
func (b *BytePipeIn[T]) NextMessage(msg T) error {

	if b.syncLost == true {
		b.syncLost = false
		b.rd = nil
		return ErrSyncLost
	}
	if b.rd == nil {
		return ErrPipeClosed
	}

	var err error
	sizeBuf := make([]byte, 4)
	c := []byte{0}
	// first step is to block waiting for the first byte, can wait forever
	for i := 0; i < 4; i++ {
		n, err := b.rd.Read(c)
		if err != nil || n != 1 {
			return b.lostSync(err)
		}
		sizeBuf[i] = c[0]
	}
	ctx := pcontext.ServerWasmContext(b.ctx)
	size, err := b.toInt(pcontext.CallTo(ctx, "toInt"), sizeBuf)
	if err != nil {
		return b.lostSync(err)
	}

	if size >= maxProtobufSizeInBytes {
		return b.lostSync(ErrTooLarge)
	}
	_, err = b.rd.Read(c)
	if err != nil {
		return err
	}
	if c[0] != 32 {
		return b.lostSync(err)
	}
	pcontext.Logf(b.ctx, pcontext.Info, "space value %d", c[0])
	if size == 0 {
		return nil
	}
	pcontext.Infof(b.ctx, "reading next payload, size is %d", size)

	count := 0
	result := make([]byte, size)
	for count < size {
		rd, err := b.rd.Read(result[count:])
		if err != nil {
			return err
		}
		count += rd
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
