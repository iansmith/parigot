package eng

import (
	"context"
	"io"
	"io/fs"
	"log"

	"github.com/iansmith/parigot/apishared"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
)

type chanFS struct {
}

func newChanFs() fs.FS {
	return &chanFS{}
}

func (c *chanFS) Open(path string) (fs.File, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{Err: fs.ErrNotExist}
	}
	log.Printf("Open: %s", path)
	if path == "chan/bind" {
		return NewChanFileOut("chan/bind"), nil
	}
	return nil, &fs.PathError{Err: fs.ErrNotExist}
}

type ChanFile struct {
	path         string
	toGuestPipe  *io.PipeReader
	fromHostPipe *io.PipeWriter
}

func NewChanFileOut(path string) apishared.HostDataWriter {
	cf := &ChanFile{path: path}
	cf.toGuestPipe, cf.fromHostPipe = io.Pipe()
	return cf
}

func (c *ChanFile) Stat() (fs.FileInfo, error) {
	return nil, apishared.ErrNotAvailable
}

func (c *ChanFile) Close() error {
	e := id.NewKernelError(id.KernelClosedErr)
	hd := apishared.NewHostDataErr(e)
	ctx := pcontext.CallTo(pcontext.ServerGoContext(context.Background()), "Close")
	pcontext.Debugf(ctx, "About to send close message, %+v", hd)
	return nil
}

func (c *ChanFile) WriteHost(hd *apishared.HostData) error {
	ctx := pcontext.CallTo(pcontext.ServerGoContext(context.Background()), "WriteHost")
	pcontext.Debugf(ctx, "THD: length %d, ptr %x-- %s", hd.Len(), hd.Ptr(), hd.IdErr())
	return nil
}

func (c *ChanFile) Read(buf []byte) (int, error) {
	ctx := pcontext.CallTo(pcontext.ServerGoContext(context.Background()), "WriteHost")
	pcontext.Debugf(ctx, "THD: read on pipe")
	return c.toGuestPipe.Read(buf)
}

func (c *ChanFile) WriteHostData(hd *apishared.HostData) error {
	panic("not implemented WriteHostData")
}
