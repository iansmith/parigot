package apiwasm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/iansmith/parigot/apishared/id"
	"google.golang.org/protobuf/proto"
)

type setupFunc func(t *testing.T) io.Reader
type checkFunc func(t *testing.T, m proto.Message, err error) any

var checker = []checkFunc{check0, check1, check2, check3, check4}

var check1Var = make([]id.Id, 2)

func TestBytePipeIn(t *testing.T) {

	var bpi *BytePipeIn
	for i, setup := range []setupFunc{setup0, nil, setup1, nil, nil} {
		if setup != nil {
			rd := setup(t)
			bpi = NewBytePipeIn(context.Background(), rd)
		}
		m, err := bpi.NextBlockUntilCall()
		if checker[i](t, m, err) == nil {
			return
		}
	}
}

func check0(t *testing.T, result proto.Message, err error) any {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error in test 0: %v", err)
	}
	bucr, ok := result.(*syscall.BlockUntilCallResponse)
	if !ok {
		t.Errorf("unexpected proto type in test 0: %v", bucr)
	}
	// doesn't matter the type
	return 0
}
func check1(t *testing.T, result proto.Message, err error) any {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error in test 0 (part2): %v", err)
	}
	bucr, ok := result.(*syscall.BlockUntilCallResponse)
	if !ok {
		t.Errorf("unexpected proto type in test 0 (part2): %v", bucr)
	}
	call := id.Unmarshal[*protosupport.CallId](bucr.GetCall())
	meth := id.Unmarshal[*protosupport.MethodId](bucr.GetMethod())
	if !call.Equal(check1Var[0]) {
		t.Errorf("expected %s for call id, int test 0 (part2) but got %s", check1Var[0], call)
	}
	if !meth.Equal(check1Var[1]) {
		t.Errorf("expected %s for method id in test 0 (part2), but got %s", check1Var[1], meth)
	}
	// doesn't matter the type
	return 0
}

func expectError(t *testing.T, actual, expected error) bool {
	t.Helper()
	if actual == nil {
		t.Errorf("expected error in test but found none")
		return false
	}
	if actual != expected {
		t.Errorf("expected error in test to be ErrTooLarge but got %T", actual)
		return false
	}
	return true

}

func check2(t *testing.T, result proto.Message, err error) any {
	if !expectError(t, err, ErrTooLarge) {
		return nil
	}
	return 0
}
func check3(t *testing.T, result proto.Message, err error) any {
	if !expectError(t, err, ErrSyncLost) {
		return nil
	}
	return 0
}
func check4(t *testing.T, result proto.Message, err error) any {
	if !expectError(t, err, ErrPipeClosed) {
		return nil
	}
	return 0
}

func setup0(t *testing.T) io.Reader {
	t.Helper()

	// empty
	blockMsg := syscall.BlockUntilCallResponse{}
	marshaled, err := proto.Marshal(&blockMsg)
	if err != nil {
		t.Errorf("unable to marshal proto in test 0")
		return nil
	}
	len1 := fmt.Sprintf("%04x ", len(marshaled))
	buf1 := mergeBuffer(t, len1, marshaled, nil)

	check1Var[0] = id.NewCallId()
	check1Var[1] = id.NewMethodId()
	blockMsg2 := syscall.BlockUntilCallResponse{
		Method: id.Marshal[protosupport.MethodId](check1Var[1]),
		Call:   id.Marshal[protosupport.CallId](check1Var[0]),
	}

	marshaled2, err := proto.Marshal(&blockMsg2)
	len2 := fmt.Sprintf("%04x ", len(marshaled2))
	buf2 := mergeBuffer(t, len2, marshaled2, buf1)

	rd := readerFromBuf(t, buf2)
	return rd
}

func setup1(t *testing.T) io.Reader {
	t.Helper()

	blockMsg := &syscall.BlockUntilCallResponse{}
	marshaled, err := proto.Marshal(blockMsg)
	if err != nil {
		t.Errorf("unable to marshal proto")
		return nil
	}
	test1 := fmt.Sprintf("%04x ", maxProtobufSizeInBytes+1)
	rd := readerFromBuf(t, mergeBuffer(t, test1, marshaled, nil))
	return rd
}
func mergeBuffer(t *testing.T, s string, b []byte, previous *bytes.Buffer) *bytes.Buffer {
	t.Helper()
	var buffer *bytes.Buffer
	if previous == nil {
		buffer = &bytes.Buffer{}
	} else {
		buffer = previous
	}

	_, err := buffer.Write([]byte(s))
	if err != nil {
		t.Errorf("bad write to bytes.Buffer")
		t.FailNow()
	}
	_, err = buffer.Write(b)
	if err != nil {
		t.Errorf("bad write to bytes.Buffer")
		t.FailNow()
	}
	return buffer
}

func readerFromBuf(t *testing.T, b *bytes.Buffer) io.Reader {
	return b
}
