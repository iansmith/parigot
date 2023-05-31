package eng

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"syscall"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/tetratelabs/wazero"
	"google.golang.org/protobuf/proto"
)

var ErrNotOpen = errors.New("attempt to close file that is not open")

// faux fs
type AsyncClientInteraction struct {
	openFauxFile map[string][]wazero.FauxFile
	origCtx      context.Context
}

func NewAsyncClientInteraction(ctx context.Context) *AsyncClientInteraction {
	a := &AsyncClientInteraction{
		openFauxFile: make(map[string][]wazero.FauxFile),
		origCtx:      ctx,
	}
	a.openFauxFile["bind"] = nil
	return a
}

func (a *AsyncClientInteraction) String() string {
	return "asyncClientInteraction"
}
func (a *AsyncClientInteraction) Check(path string) bool {
	_, ok := a.openFauxFile[path]
	return ok
}

// faux file
func (a *AsyncClientInteraction) Create(path string, advisoryRead, advisoryWrite bool) (wazero.FauxFile, syscall.Errno) {
	result, ok := a.openFauxFile[path]
	if !ok {
		a.openFauxFile[path] = []wazero.FauxFile{}
		result = []wazero.FauxFile{}
	}
	entry := NewAsyncClientInteractionEntry(a, path, advisoryRead, advisoryWrite)
	result = append(result, entry)
	a.openFauxFile[path] = result
	return entry, 0
}

func (a *AsyncClientInteraction) Close(ff wazero.FauxFile) syscall.Errno {
	list, ok := a.openFauxFile[ff.Path()]
	if !ok {
		pcontext.Fatalf(a.origCtx, ErrNotOpen.Error())
		pcontext.Dump(a.origCtx)
		return syscall.ENOENT
	}
	var found wazero.FauxFile
	for i, elem := range list {
		if elem == ff {
			if len(list) == 1 {
				a.openFauxFile[ff.Path()] = nil
			} else if len(list)-1 == i {
				if i == 0 {
					a.openFauxFile[ff.Path()] = nil
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
		pcontext.Fatalf(a.origCtx, ErrNotOpen.Error())
		pcontext.Dump(a.origCtx)
		return syscall.ENOTEMPTY
	}
	pcontext.Dump(a.origCtx)
	return 0
}
func (a *AsyncClientInteraction) Open(path string, read, write bool) (wazero.FauxFile, syscall.Errno) {
	list, ok := a.openFauxFile[path]
	if !ok {
		list = []wazero.FauxFile{}
		a.openFauxFile[path] = list
	}
	ff := NewAsyncClientInteractionEntry(a, path, read, write)
	list = append(list, ff)
	a.openFauxFile[path] = list
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
	if err := a.toGuest.Close(); err != nil {
		pcontext.Dump(a.parent.origCtx)
		return err
	}
	pcontext.Dump(a.parent.origCtx)
	return a.guestReader.Close()
}

func (a *AsyncClientInteraction) Send(path string, m proto.Message) error {
	ctx := pcontext.CallTo(a.origCtx, "Send")
	openedFile := a.openFauxFile[path]
	numOpen := len(openedFile)
	if numOpen == 0 {
		return fmt.Errorf("Send failed because file '%s' is not open", path)
	}
	choice := openedFile[rand.Intn(numOpen)].(*AsyncClientInteractionEntry)
	flat, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("%04x ", len(flat))
	pcontext.Debugf(a.origCtx, "SENDING: %s to %p", s, choice)
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
