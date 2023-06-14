package file

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/file/v1"

	"github.com/tetratelabs/wazero/api"
)

// make sure edit the makefile so you can have FileId and FileErrId, just like queue
// you'll need to pick the short names and letters for them... I would
// recommend f for FileId and F for FileErrId, but you can choose
// others if you want.

type FilePlugin struct{}

var _ = unsafe.Sizeof([]byte{})

// RULE: All files opened by a user program have to have a
// RULE: pathname that looks like /app/...  also, any
// RULE: use of . or .. in the path is not allowed.

// You may want to look at the utilities in filepath and
// strings.Split(). You can ignore  os.FileSeparator, we
// always assume the separator in filenames is /

// this should contain all the internal things you need to
// track of... path, readers, writers, etc.  Make sure to include
// a lastAccessTime that is derived from our context package
// via CurrentTime()... later on we will be expiring entries
// in fileDataCache

// for now, create a map of FileId -> to myFileInfo
//var fileDataCache = make(map[Fileid])*myFileInfo

func (*FilePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "file", "open_", open) // this should call the "wrapper"
	return true
}

// true native implementation of open... assume this is read only
func openImpl(ctx context.Context, in *file.OpenRequest, out *file.OpenResponse) int32 {
	// use Os
	return int32(file.FileErr_NoError)
}

// the wrappers always look like this.. notice where openImpl is in this function
func open(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.OpenRequest{}
	resp := &file.OpenResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[file]open", m, stack, openImpl, req, resp)
}

//  add two more functions: create and close.  Create is
// like open, but WRITE only.  Close frees up items from the
// table fileDataCache
