package eng

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"syscall"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	gsys "github.com/iansmith/parigot/g/syscall/v1"

	"github.com/tetratelabs/wazero"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	ErrNotOpen       = errors.New("attempt to close file that is not open")
	ErrAlreadyExists = errors.New("attempt create a file that is already opened by others")
)

// faux fs
type AsyncClientInteraction struct {
	openFauxFile sync.Map
	origCtx      context.Context
}

func NewAsyncClientInteraction(ctx context.Context) *AsyncClientInteraction {
	a := &AsyncClientInteraction{
		openFauxFile: sync.Map{},
		origCtx:      ctx,
	}

	a.openFauxFile.Store("bind", []wazero.FauxFile{})
	return a
}

func (a *AsyncClientInteraction) String() string {
	return "asyncClientInteraction"
}

func (a *AsyncClientInteraction) Exists(path string) bool {
	_, ok := a.openFauxFile.Load(path)
	return ok
}

// faux file
func (a *AsyncClientInteraction) Create(path string, advisoryRead, advisoryWrite bool) (wazero.FauxFile, syscall.Errno) {
	ctx := pcontext.CallTo(a.origCtx, "Create")
	result, ok := a.openFauxFile.Load(path)
	if !ok {
		a.openFauxFile.Store(path, []wazero.FauxFile{})
		result = []wazero.FauxFile{}
	} else {
		if len(result.([]wazero.FauxFile)) != 0 {
			pcontext.Errorf(ctx, ErrAlreadyExists.Error())
			pcontext.Dump(ctx)
			return nil, syscall.EACCES
		}
	}
	entry := NewAsyncClientInteractionEntry(a, path, advisoryRead, advisoryWrite)
	list := result.([]wazero.FauxFile)
	list = append(list, entry)
	a.openFauxFile.Store(path, result)
	pcontext.Dump(ctx)
	return entry, 0
}

func (a *AsyncClientInteraction) Close(ff wazero.FauxFile) syscall.Errno {
	ctx := pcontext.CallTo(a.origCtx, "Close")

	ret, ok := a.openFauxFile.Load(ff.Path())
	if !ok {
		pcontext.Errorf(ctx, ErrNotOpen.Error())
		pcontext.Dump(ctx)
		return syscall.ENOENT
	}
	list := ret.([]wazero.FauxFile)
	var found wazero.FauxFile
	for i, elem := range list {
		if elem == ff {
			if len(list) == 1 {
				a.openFauxFile.Store(ff.Path(), nil)
			} else if len(list)-1 == i {
				if i == 0 {
					a.openFauxFile.Store(ff.Path(), nil)
				} else {
					list = list[:i]
				}
			} else {
				list = append(list[:i], list[i+1:]...)
			}
			found = elem
			break
		}
	}
	if found == nil {
		pcontext.Errorf(ctx, ErrNotOpen.Error())
		pcontext.Dump(ctx)
		return syscall.ENOENT
	}
	pcontext.Dump(ctx)
	return 0
}
func (a *AsyncClientInteraction) Open(path string, read, write bool) (wazero.FauxFile, syscall.Errno) {
	ctx := pcontext.CallTo(a.origCtx, "Open")

	raw, ok := a.openFauxFile.Load(path)
	pcontext.Debugf(ctx, "open %s: raw %+v, and ok %v", path, raw, ok)

	var list = []wazero.FauxFile{}
	if !ok {
		a.openFauxFile.Store(path, list)
	} else {
		list = raw.([]wazero.FauxFile)
	}
	ff := NewAsyncClientInteractionEntry(a, path, read, write)
	list = append(list, ff)
	a.openFauxFile.Store(path, list)
	return ff, 0
}

type AsyncClientInteractionEntry struct {
	parent        *AsyncClientInteraction
	path          string
	canRead       bool
	canWrite      bool
	toGuest       *io.PipeReader
	injectToGuest *io.PipeWriter
	toHost        *io.PipeWriter
	guestReader   *io.PipeReader
}

func NewAsyncClientInteractionEntry(parent *AsyncClientInteraction, path string, advisoryRead, advisoryWrite bool) wazero.FauxFile {

	a := &AsyncClientInteractionEntry{
		path:     path,
		parent:   parent,
		canRead:  advisoryRead,
		canWrite: advisoryWrite,
	}
	a.toGuest, a.injectToGuest = io.Pipe()
	a.guestReader, a.toHost = io.Pipe()
	return a
}

func (a *AsyncClientInteractionEntry) Path() string {
	return a.path
}

func (a *AsyncClientInteractionEntry) Read(buf []byte) (int, error) {
	return a.toGuest.Read(buf)
}
func (a *AsyncClientInteractionEntry) Write(buf []byte) (int, error) {
	return a.toHost.Write(buf)
}
func (a *AsyncClientInteractionEntry) Close() error {
	ctx := pcontext.CallTo(a.parent.origCtx, "Close")

	if err := a.toGuest.Close(); err != nil {
		pcontext.Dump(ctx)
		return err
	}
	pcontext.Dump(ctx)
	return a.guestReader.Close()
}

func (a *AsyncClientInteraction) Send(path string, m proto.Message) error {
	ctx := pcontext.CallTo(a.origCtx, "Send")
	raw, ok := a.openFauxFile.Load(path)
	if !ok {
		pcontext.Errorf(ctx, ErrNotOpen.Error()+":"+path)
		pcontext.Dump(ctx)
		return ErrNotFound
	}
	openedFile := raw.([]wazero.FauxFile)
	numOpen := len(openedFile)
	if numOpen == 0 {
		return fmt.Errorf("Send failed because file '%s' is not open", path)
	}
	choice := openedFile[rand.Intn(numOpen)].(*AsyncClientInteractionEntry)
	pcontext.Infof(ctx, "only writing to one listener on %s", choice.Path())
	flat, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("%04x ", len(flat))
	pcontext.Debugf(ctx, "SENDING: %s to %p", s, choice)
	transformed := []byte(s)

	n, err := choice.injectToGuest.Write(transformed)
	if n != 5 {
		return fmt.Errorf("unexpect failure to format for write: %v", err)
	}
	if err != nil {
		return err
	}
	count := 0
	pcontext.Debugf(ctx, "flattened protobuf of size %d bytes", len(flat))
	for count < len(flat) {
		written, err := choice.injectToGuest.Write(flat[count:])
		if err != nil {
			return err
		}
		count += written
	}
	return nil
}

func (a *AsyncClientInteraction) Dispatch(path string, pbAny *anypb.Any) (*anypb.Any, id.CallId, gsys.KernelErr) {
	print(">>>>>>>> path is ", path, "\n")
	return nil, id.NewCallId(), gsys.KernelErr_NoError
}
