package file

import (
	"context"
	"io"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared"
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
const maxBufSize = apishared.MaxBufSize

var (
	fileSvc    *fileSvcImpl
	_          = unsafe.Sizeof([]byte{})
	fpathTofid = make(map[string]file.FileId)
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
	length         int // length of the file content
	content        string
	status         FileStatus
	createDate     string
	lastAccessTime string
	prevRune       int // index of previous rune; or 0
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

func (fi *fileInfo) Read(p []byte) (n int, err error) {
	// check current position is not the end of file
	if fi.prevRune == fi.length {
		return 0, io.EOF
	}

	for i, b := range []byte(fi.content)[fi.prevRune:] {
		// when length of file is larget than the size of buffer
		if i >= len(p) {
			fi.prevRune += i
			return len(p), nil
		}
		p[i] = b
	}

	// buffer is larger
	fi.prevRune = fi.length
	return fi.length, nil
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
	fid, exist := fpathTofid[cleanPath]
	if !exist {
		pcontext.Errorf(ctx, "file does not exist and cannot be opened: %s", fpath)

		return int32(file.FileErr_NotExistError)
	}

	myFileInfo := fileDataCache[fid]
	// check file status
	if myFileInfo.status == Fs_Open {
		pcontext.Errorf(ctx, "file is open, cannot be opened again: %s", fpath)

		return int32(file.FileErr_AlreadyInUseError)
	}

	resp.Id = fid.Marshal()

	myFileInfo.lastAccessTime = pcontext.CurrentTimeString(ctx, true)
	myFileInfo.status = Fs_Open

	return int32(file.FileErr_NoError)
}

// WRITE only
func (f *fileSvcImpl) create(ctx context.Context, req *file.CreateRequest,
	resp *file.CreateResponse) int32 {

	currentTime := pcontext.CurrentTimeString(ctx, true)

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

		myFileInfo := fileDataCache[fid]
		// check file status first, a opened file cannot be created at the same time
		if myFileInfo.status == Fs_Open {
			pcontext.Errorf(ctx, "file is open, cannot be created: %s", fpath)

			return int32(file.FileErr_AlreadyInUseError)
		}
		// extend a file
		resp.Truncated = true
		myFileInfo.content += content
		myFileInfo.lastAccessTime = currentTime
		myFileInfo.length += len(content)
	} else {
		// create a file id
		fid := file.NewFileId()
		resp.Id = fid.Marshal()

		newFileInfo := fileInfo{
			id:             fid,
			path:           cleanPath,
			length:         len(content),
			content:        content,
			status:         Fs_Close,
			createDate:     currentTime,
			lastAccessTime: currentTime,
			prevRune:       0,
		}
		fileDataCache[fid] = &newFileInfo
		fpathTofid[cleanPath] = fid
	}

	return int32(file.FileErr_NoError)
}

// free up item from the fileDataCache
func (f *fileSvcImpl) close(ctx context.Context, req *file.CloseRequest, resp *file.CloseResponse) int32 {
	fid := file.UnmarshalFileId(req.GetId())
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

	resp.Id = req.GetId()
	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) read(ctx context.Context, req *file.ReadRequest, resp *file.ReadResponse) int32 {
	fid := file.UnmarshalFileId(req.GetId())
	fileDataCache := *f.fileDataCache

	// check if file exists. We cannot read a file which doesn't exist
	if _, exist := fileDataCache[fid]; !exist {
		pcontext.Errorf(ctx, "file does not exist, cannot be read: %d", fid)

		return int32(file.FileErr_NotExistError)
	}

	myFileInfo := fileDataCache[fid]
	fpath := myFileInfo.path
	// check file status, we cannot read a closed file
	if myFileInfo.status == Fs_Close {
		pcontext.Errorf(ctx, "file is closed, cannot be read: %s", fpath)

		return int32(file.FileErr_FileClosedError)
	}

	bufSize := req.GetBufSize()
	// Check if bufSize exceeds the maximum buffer size allowed
	if bufSize > maxBufSize {
		pcontext.Errorf(ctx, "the expected buffer size %d exceeds the maximum buffer"+
			"size (%d) allowed", bufSize, maxBufSize)

		return int32(file.FileErr_BufferFullError)
	}

	buf := make([]byte, bufSize)
	n, _ := myFileInfo.Read(buf)
	myFileInfo.lastAccessTime = pcontext.CurrentTimeString(ctx, true)

	resp.Id = req.GetId()
	resp.Buf = buf
	resp.NumRead = int32(n)

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
