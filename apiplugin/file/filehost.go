package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	file "github.com/iansmith/parigot/g/file/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	"github.com/iansmith/parigot/sys"

	"google.golang.org/protobuf/proto"

	"github.com/tetratelabs/wazero/api"
)

// make sure edit the makefile so you can have FileId and FileErrId, just like queue
// you'll need to pick the short names and letters for them... I would
// recommend f for FileId and F for FileErrId, but you can choose
// others if you want.
var fileSvc *fileSvcImpl

type filePlugin struct{}

var _ = unsafe.Sizeof([]byte{})

var ParigiotInitialize sys.ParigotInit = &filePlugin{}

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

type myFileInfo struct {
	// id
	// path
	// status
	// lastAccessTime
}

// for now, create a map of FileId -> to myFileInfo
// var fileDataCache = make(map[Fileid])*myFileInfo

type fileSvcImpl struct {
	fileDataCache *map[file.FileId]*myFileInfo
	ctx           context.Context
}

func (*filePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "file", "open_file_", openFileHost) // this should call the "wrapper"
	e.AddSupportedFunc(ctx, "file", "create_file_", createFileHost)
	e.AddSupportedFunc(ctx, "file", "close_file_", closeFileHost)

	_ = newFileSvc(ctx)

	return true
}

func hostBase[T proto.Message, U proto.Message](ctx context.Context, fnName string,
	fn func(context.Context, T, U) id.IdRaw, m api.Module, stack []uint64, req T, resp U) {
	defer func() {
		if r := recover(); r != nil {
			print(">>>>>>>> Trapped recover in set up for   ", fnName, "<<<<<<<<<<\n")
		}
	}()
	apiplugin.InvokeImplFromStack(ctx, fnName, m, stack, fn, req, resp)
}

// // true native implementation of open... assume this is read only
// func openImpl(ctx context.Context, in *file.OpenRequest, out *file.OpenResponse) int32 {
// 	// use Os
// 	return int32(file.FileErr_NoError)
// }

// // the wrappers always look like this.. notice where openImpl is in this function
// func open(ctx context.Context, m api.Module, stack []uint64) {
// 	req := &file.OpenRequest{}
// 	resp := &file.OpenResponse{}
// 	apiplugin.InvokeImplFromStack(ctx, "[file]open", m, stack, openImpl, req, resp)
// }

func openFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &filemsg.OpenRequest{}
	resp := &filemsg.OpenResponse{}

	hostBase(ctx, "[file]open", fileSvc.open, m, stack, req, resp)
}

func createFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &filemsg.CreateRequest{}
	resp := &filemsg.CreateResponse{}

	hostBase(ctx, "[file]create", fileSvc.create, m, stack, req, resp)
}

func closeFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &filemsg.CloseRequest{}
	resp := &filemsg.CloseResponse{}

	hostBase(ctx, "[file]close", fileSvc.close, m, stack, req, resp)
}

func newFileSvc(ctx context.Context) *fileSvcImpl {
	newCtx := pcontext.ServerGoContext(ctx)

	f := &fileSvcImpl{
		fileDataCache: &map[file.FileId]*myFileInfo{},
		ctx:           newCtx,
	}

	return f
}

// read only, need to be implemented
func (f *fileSvcImpl) open(ctx context.Context, req *filemsg.OpenRequest, resp *filemsg.OpenResponse) id.IdRaw {
	return file.FileErrIdNoErr.Raw()
}

// write only, need to be implemented
func (f *fileSvcImpl) create(ctx context.Context, req *filemsg.CreateRequest, resp *filemsg.CreateResponse) id.IdRaw {
	return file.FileErrIdNoErr.Raw()
}

// close only, need to be implemented
func (f *fileSvcImpl) close(ctx context.Context, req *filemsg.CloseRequest, resp *filemsg.CloseResponse) id.IdRaw {
	return file.FileErrIdNoErr.Raw()
}

//  add two more functions: create and close.  Create is
// like open, but WRITE only.  Close frees up items from the
// table fileDataCache
