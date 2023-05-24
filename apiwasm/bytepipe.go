package apiwasm

import (
	"context"
	"errors"
	"io"
	"time"

	pcontext "github.com/iansmith/parigot/context"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"google.golang.org/protobuf/proto"
)

type BytePipeIn struct {
	rd       io.Reader
	ch       chan byte
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

type allowedToRead interface {
	proto.Message
}

// NewBytePipeIn creates a new bytePipeIn that reads on the given reader.
// NewBytePipeIn creates a goroutine so that the rest of the bytePipeIn
// can use channels to read the bytes and do timeouts.
func NewBytePipeIn(ctx context.Context, rd io.Reader) *BytePipeIn {
	bpi := &BytePipeIn{rd: rd, ctx: ctx, ch: make(chan byte)}
	go func(ctx context.Context) {
		bpi.reader(pcontext.CallTo(ctx, "bytePipeIn.reader"))
	}(ctx)
	return bpi
}

func (b *BytePipeIn) reader(ctx context.Context) {
	if b.rd == nil {
		close(b.ch)
		return
	}
	for {
		buf := make([]byte, 1)
		_, err := b.rd.Read(buf)
		if err != io.EOF && err != nil {
			pcontext.Errorf(b.ctx, "failed to read from pipe: %v", err)
			close(b.ch)
			return
		}
		if err == io.EOF {
			close(b.ch)
			return
		}
		b.ch <- buf[0]
	}
}
func (b *BytePipeIn) NextBlockUntilCall() (*syscallmsg.BlockUntilCallResponse, error) {
	buc := &syscallmsg.BlockUntilCallResponse{}
	if err := readNext(b, buc); err != nil {
		return nil, err
	}
	return buc, nil
}

func readNext[U allowedToRead](b *BytePipeIn, msg U) error {
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
	// first step is to block waiting for the first byte, can wait forever
	sizeBuf[0] = <-b.ch
	sizeBuf[1], err = b.readByteShortWait()
	if err != nil {
		return b.lostSync(err)
	}
	sizeBuf[2], err = b.readByteShortWait()
	if err != nil {
		return b.lostSync(err)
	}
	sizeBuf[3], err = b.readByteShortWait()
	if err != nil {
		return b.lostSync(err)
	}
	ctx := pcontext.ServerWasmContext(b.ctx)
	size, err := b.toInt(pcontext.CallTo(ctx, "toInt"), sizeBuf)
	if err != nil {
		return b.lostSync(err)
	}
	if size >= maxProtobufSizeInBytes {
		return b.lostSync(ErrTooLarge)
	}
	space, err := b.readByteShortWait()
	if space != 32 {
		return b.lostSync(err)
	}
	count := 0
	result := make([]byte, size)
	for count < size {
		result[count], err = b.readByteShortWait()
		if err != nil {
			return b.lostSync(err)
		}
		count++
	}
	err = proto.Unmarshal(result, msg)
	if err != nil {
		return ErrUnmarshal
	}
	return nil
}

func (b *BytePipeIn) readByteShortWait() (byte, error) {

	select {
	case data := <-b.ch:
		return data, nil
	case <-time.After(time.Duration(timeoutInMillis) * time.Millisecond):
		return 0, ErrTimeout
	}
}
func (b *BytePipeIn) toInt(ctx context.Context, sizeBuf []byte) (int, error) {
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

func (b *BytePipeIn) lostSync(err error) error {
	b.syncLost = true
	b.rd = nil
	if err != nil {
		return err
	}
	return ErrSyncLost
}
