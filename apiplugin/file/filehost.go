package file

import (
	"context"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/iansmith/parigot/apiplugin"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/file/v1"

	"google.golang.org/protobuf/proto"

	"github.com/tetratelabs/wazero/api"
)

// make sure edit the makefile so you can have FileId and FileErrId, just like queue
// you'll need to pick the short names and letters for them... I would
// recommend f for FileId and F for FileErrId, but you can choose
// others if you want.
const pathPrefix = "/parigot/app/"

var (
	fileSvc            *fileSvcImpl
	_                                        = unsafe.Sizeof([]byte{})
	ParigiotInitialize apiplugin.ParigotInit = &FilePlugin{}
	fpathTofid                               = make(map[string]file.FileId)
)

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

type FilePlugin struct{}

type fileInfo struct {
	id             file.FileId
	path           string
	content        string
	status         FileStatus
	createDate     time.Time
	lastAccessTime time.Time
}

// for now, create a map of FileId -> to myFileInfo
// var fileDataCache = make(map[Fileid])*myFileInfo

type fileSvcImpl struct {
	fileDataCache *map[file.FileId]*fileInfo
	ctx           context.Context
}

// enum for file status
type FileStatus int

const (
	Fs_Open FileStatus = iota
	Fs_Close
)

func (fs FileStatus) String() string {
	return []string{"Open", "Close"}[fs]
}

func (*FilePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "file", "open_file_", openFileHost) // this should call the "wrapper"
	e.AddSupportedFunc(ctx, "file", "create_file_", createFileHost)
	e.AddSupportedFunc(ctx, "file", "close_file_", closeFileHost)

	_ = newFileSvc(ctx)

	return true
}

func hostBase[T proto.Message, U proto.Message](ctx context.Context, fnName string,
	fn func(context.Context, T, U) int32, m api.Module, stack []uint64, req T, resp U) {
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
	req := &file.OpenRequest{}
	resp := &file.OpenResponse{}

	hostBase(ctx, "[file]open", fileSvc.open, m, stack, req, resp)
}

func createFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.CreateRequest{}
	resp := &file.CreateResponse{}

	hostBase(ctx, "[file]create", fileSvc.create, m, stack, req, resp)
}

func closeFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.CloseRequest{}
	resp := &file.CloseResponse{}

	hostBase(ctx, "[file]close", fileSvc.close, m, stack, req, resp)
}

func newFileSvc(ctx context.Context) *fileSvcImpl {
	newCtx := pcontext.ServerGoContext(ctx)

	f := &fileSvcImpl{
		fileDataCache: &map[file.FileId]*fileInfo{},
		ctx:           newCtx,
	}

	return f
}

// READ only, need to be implemented
func (f *fileSvcImpl) open(ctx context.Context, req *file.OpenRequest,
	resp *file.OpenResponse) int32 {

	fpath := req.GetPath()

	cleanPath, valid := isValidPath(fpath)
	if !valid {
		pcontext.Errorf(ctx, "file path is not valid: %s", fpath)

		return int32(file.FileErr_InvalidPathError)
	}

	resp.Path = cleanPath
	fileDataCache := *f.fileDataCache

	// if file doesn't exist, return an error
	if fid, exist := fpathTofid[cleanPath]; !exist {
		pcontext.Errorf(ctx, "file does not exist and cannot be opened: %s", fpath)

		return int32(file.FileErr_NotExistError)
	} else {
		// check file status
		if fileDataCache[fid].status == Fs_Open {
			pcontext.Errorf(ctx, "file is open, cannot be opened again: %s", fpath)

			return int32(file.FileErr_AlreadyInUseError)
		}

		resp.Id = fid.Marshal()
		fileDataCache[fid].lastAccessTime = time.Now()
		fileDataCache[fid].status = Fs_Open
	}

	return int32(file.FileErr_NoError)
}

// WRITE only
func (f *fileSvcImpl) create(ctx context.Context, req *file.CreateRequest,
	resp *file.CreateResponse) int32 {

	fpath := req.GetPath()

	cleanPath, valid := isValidPath(fpath)
	if !valid {
		pcontext.Errorf(ctx, "File path is not valid: %s", fpath)

		return int32(file.FileErr_InvalidPathError)
	}

	resp.Path = cleanPath
	resp.Truncated = false
	content := req.GetContent()
	fileDataCache := *f.fileDataCache

	// if file/path exists, truncating
	if fid, exist := fpathTofid[fpath]; exist {
		resp.Id = fid.Marshal()

		// check file status first, a opened file cannot be be written at the same time
		if fileDataCache[fid].status == Fs_Open {
			pcontext.Errorf(ctx, "file is open, cannot be created: %s", fpath)

			return int32(file.FileErr_AlreadyInUseError)
		}
		// extend a file
		resp.Truncated = true
		fileDataCache[fid].content += content
		fileDataCache[fid].lastAccessTime = time.Now()
		fileDataCache[fid].createDate = time.Now()
	} else {
		// create a file id
		fid := file.NewFileId()
		resp.Id = fid.Marshal()

		newFileInfo := fileInfo{
			id:             fid,
			path:           cleanPath,
			content:        content,
			status:         Fs_Close,
			createDate:     time.Now(),
			lastAccessTime: time.Now(),
		}
		fileDataCache[fid] = &newFileInfo
		fpathTofid[cleanPath] = fid
	}

	return int32(file.FileErr_NoError)
}

// free up item from the fileDataCache
func (f *fileSvcImpl) close(ctx context.Context, req *file.CloseRequest, resp *file.CloseResponse) int32 {
	fid := file.UnmarshalFileId(req.Id)
	fileDataCache := *f.fileDataCache

	// check if file exists. We cannot delete a file which doesn't exist
	if _, exist := fileDataCache[fid]; !exist {
		pcontext.Errorf(ctx, "file does not exist, cannot be closed: %d", fid)

		return int32(file.FileErr_NotExistError)
	}

	// remove file from the fileDataCache
	fpath := fileDataCache[fid].path
	delete(fileDataCache, fid)
	delete(fpathTofid, fpath)

	resp.Id = req.Id
	return int32(file.FileErr_NoError)
}

// A valid path should be a shortest path name equivalent to path by purely lexical processingand.
// Specifically, it should start with "/parigot/app/", also, any use of '.', '..', in the path is
// not allowed.
func isValidPath(fpath string) (string, bool) {
	fileName := filepath.Base(fpath)
	dir := strings.ReplaceAll(fpath, fileName, "")
	if !strings.HasPrefix(dir, pathPrefix) || strings.Contains(dir, ".") {
		return fpath, false
	}
	cleanPath := filepath.Clean(fpath)

	return cleanPath, true
}
