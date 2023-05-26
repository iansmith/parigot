package apishared

import (
	"errors"
	"io/fs"

	"github.com/iansmith/parigot/apishared/id"
)

// HostData represents data passed from the guest to the host.  Be careful
// that all the types specify their size because host and guest may have
// different sized words (int).
type HostData struct {
	ptr, length uint32
	idErr       [2]uint64
}

var ErrNotAvailable = errors.New("unable to provide implementation of request")

type HostDataWriter interface {
	fs.File
	WriteHostData(*HostData) error
}

func NewHostDataNoErr(ptr, length uint32) *HostData {
	hd := &HostData{ptr: ptr, length: length}
	noerr := id.Unmarshal(id.NoKernelError())
	hd.idErr[0] = noerr.High()
	hd.idErr[1] = noerr.Low()
	return hd
}

func (h *HostData) Ptr() uint32 {
	return h.ptr
}
func (h *HostData) Len() uint32 {
	return h.length
}
func (h *HostData) IdErr() id.Id {
	return id.NewIdCopy(h.idErr[0], h.idErr[1])
}
func NewHostDataErr(idErr id.Id) *HostData {
	hd := &HostData{ptr: 0, length: 0}
	hd.idErr[0] = idErr.High()
	hd.idErr[1] = idErr.Low()
	return hd

}
