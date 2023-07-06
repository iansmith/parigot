package file

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	apishared "github.com/iansmith/parigot/api/shared"
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
const pathPrefix = apishared.FileServicePathPrefix
const maxBufSize = apishared.FileServiceMaxBufSize

var (
	fileSvc *fileSvcImpl
	_       = unsafe.Sizeof([]byte{})
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
	createDate     time.Time
	lastAccessTime time.Time

	rdClose io.ReadCloser
	wrClose io.WriteCloser
}

// for now, create a map of FileId -> to myFileInfo
// var fileDataCache = make(map[Fileid])*myFileInfo

type fileSvcImpl struct {
	fileDataCache *map[file.FileId]*fileInfo
	ctx           context.Context
	// track fid based on file path
	fpathTofid *map[string]file.FileId
	isTesting  bool

	defaultOpenHook   OpenHook
	defaultCreateHook CreateHook
}

// enum for file status
type FileStatus int

const (
	Fs_Write FileStatus = iota
	Fs_Read
	Fs_Close
)

func (fs FileStatus) String() string {
	return []string{"Write", "Read", "Close"}[fs]
}

func (*FilePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "file", "open_file_", openFileHost) // this should call the "wrapper"
	e.AddSupportedFunc(ctx, "file", "create_file_", createFileHost)
	e.AddSupportedFunc(ctx, "file", "close_file_", closeFileHost)
	e.AddSupportedFunc(ctx, "file", "read_file_", readFileHost)
	e.AddSupportedFunc(ctx, "file", "delete_file_", deleteFileHost)
	e.AddSupportedFunc(ctx, "file", "write_file_", writeFileHost)

	_ = newFileSvc(ctx)

	return true
}

func (fi *fileInfo) Read(p []byte) (n int, err error) {
	n, err = fi.rdClose.Read(p)
	if err == io.EOF {
		log.Printf("We read %d bytes and the file has exhausted its content", n)
	} else if err != nil {
		log.Fatal("Error reading from a file: ", err)
	} else {
		log.Printf("We read %d bytes", n)
	}

	return n, err
}

func (fi *fileInfo) Write(p []byte) (n int, err error) {
	n, err = fi.wrClose.Write(p)
	if err != nil {
		log.Fatal("Error writing to a file: ", err)
	}

	fi.length += n
	log.Printf("We write %d bytes", n)

	return n, nil
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

func readFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.ReadRequest{}
	resp := &file.ReadResponse{}

	hostBase(ctx, "[file]read", fileSvc.read, m, stack, req, resp)
}

func deleteFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.DeleteRequest{}
	resp := &file.DeleteResponse{}

	hostBase(ctx, "[file]delete", fileSvc.delete, m, stack, req, resp)
}

func writeFileHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.WriteRequest{}
	resp := &file.WriteResponse{}

	hostBase(ctx, "[file]write", fileSvc.write, m, stack, req, resp)
}

func newFileSvc(ctx context.Context) *fileSvcImpl {
	newCtx := pcontext.ServerGoContext(ctx)

	f := &fileSvcImpl{
		fileDataCache: &map[file.FileId]*fileInfo{},
		ctx:           newCtx,
		fpathTofid:    &map[string]file.FileId{},
		isTesting:     false,

		defaultOpenHook:   openHookForFiles,
		defaultCreateHook: createHookForFiles,
	}

	return f
}

func (f *fileSvcImpl) open(ctx context.Context, req *file.OpenRequest,
	resp *file.OpenResponse) int32 {

	// If we are in test mode, set the default OpenHook to openHookForStrings.
	// In this mode, we do not operate on the real files on the disk
	if f.isTesting {
		f.defaultOpenHook = openHookForStrings
	}

	fpath := req.GetPath()
	cleanPath, valid := isValidPath(fpath)

	// Validate the file path
	if !valid {
		pcontext.Errorf(ctx, "File path is not valid: %s", fpath)
		return int32(file.FileErr_InvalidPathError)
	}
	resp.Path = cleanPath

	// Check if the file exists
	fid, exist := (*f.fpathTofid)[cleanPath]
	if !exist {
		pcontext.Errorf(ctx, "File does not exist and cannot be opened: %s", fpath)
		return int32(file.FileErr_NotExistError)
	}

	// If file exists, fetch its information from fileDataCache
	myFileInfo := (*f.fileDataCache)[fid]

	// Check file status
	if myFileInfo.status != Fs_Close {
		pcontext.Errorf(ctx, "File is already in use: %s", fpath)
		return int32(file.FileErr_AlreadyInUseError)
	}

	// If the file is closed, update the file information and open it for reading
	resp.Id = fid.Marshal()

	myFileInfo.lastAccessTime = pcontext.CurrentTime(ctx)
	myFileInfo.status = Fs_Read
	myFileInfo.rdClose = f.defaultOpenHook(cleanPath)

	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) create(ctx context.Context, req *file.CreateRequest,
	resp *file.CreateResponse) int32 {

	fpath := req.GetPath()
	// Validate the file path
	cleanPath, valid := isValidPath(fpath)
	if !valid {
		pcontext.Errorf(ctx, "Invalid file path: %s", fpath)
		return int32(file.FileErr_InvalidPathError)
	}

	resp.Path = cleanPath
	resp.Truncated = false

	// Validate the content buffer
	content := req.GetContent()
	buf := []byte(content)
	if !isValidBuf(buf) {
		pcontext.Errorf(ctx, "Content size %d exceeds the maximum buffer size (%d) "+
			"allowed", len(buf), maxBufSize)
		return int32(file.FileErr_LargeBufError)
	}

	// If the file/path already exists, truncate it
	if fid, exist := (*f.fpathTofid)[fpath]; exist {
		// Update the response
		resp.Id = fid.Marshal()
		resp.Truncated = true

		// Fetch existing file information
		myFileInfo := (*f.fileDataCache)[fid]

		// The create request only applies to a closed file
		if myFileInfo.status != Fs_Close {
			pcontext.Errorf(ctx, "file is in use: %s", fpath)
			return int32(file.FileErr_AlreadyInUseError)
		}

		// Extend the file
		myFileInfo.content += content
		myFileInfo.lastAccessTime = pcontext.CurrentTime(ctx)
		myFileInfo.status = Fs_Write
		myFileInfo.Write(buf)

	} else {
		// If file/path does not exist, create a new file
		fid := f.createANewFile(cleanPath, content)
		resp.Id = fid.Marshal()
	}

	return int32(file.FileErr_NoError)
}

// turn the file status to close and close the file
func (f *fileSvcImpl) close(ctx context.Context, req *file.CloseRequest,
	resp *file.CloseResponse) int32 {

	fid := file.UnmarshalFileId(req.GetId())

	// Fetch file data from cache
	fileData, exist := (*f.fileDataCache)[fid]

	// Validate if the file exists in the cache
	if !exist {
		pcontext.Errorf(ctx, "File does not exist, cannot be closed: %d", fid)
		return int32(file.FileErr_NotExistError)
	}

	// Check if the file is already closed
	if fileData.status == Fs_Close {
		pcontext.Errorf(ctx, "File is already closed, cannot be closed again: %d", fid)
		return int32(file.FileErr_FileClosedError)
	}

	status := fileData.status

	// If the file exists and is not closed, change its status to "closed"
	(*f.fileDataCache)[fid].status = Fs_Close
	(*f.fileDataCache)[fid].lastAccessTime = pcontext.CurrentTime(ctx)

	switch status {
	case Fs_Read:
		(*f.fileDataCache)[fid].rdClose.Close()
	default:
		(*f.fileDataCache)[fid].wrClose.Close()
	}

	resp.Id = req.GetId()

	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) read(ctx context.Context, req *file.ReadRequest,
	resp *file.ReadResponse) int32 {

	fid := file.UnmarshalFileId(req.GetId())

	// Validate if the file exists in the cache
	if _, exist := (*f.fileDataCache)[fid]; !exist {
		pcontext.Errorf(ctx, "File does not exist, cannot be read: %d", fid)
		return int32(file.FileErr_NotExistError)
	}

	myFileInfo := (*f.fileDataCache)[fid]

	// Verify file status to prevent errors during operation.
	// Only files with "read" status can be processed.
	switch myFileInfo.status {
	// If file status is "closed", log an error and return a file closed error code.
	case Fs_Close:
		pcontext.Errorf(ctx, "Operation aborted. File with ID: %d is closed.", fid)
		return int32(file.FileErr_FileClosedError)
	// If file status is "write", meaning it is currently being written by others,
	// log an error and return a file already in use error code.
	case Fs_Write:
		pcontext.Errorf(ctx, "Operation aborted. File with ID: %d is being written by others.", fid)
		return int32(file.FileErr_AlreadyInUseError)
	}

	// Check if the file's reader is initialized
	if myFileInfo.rdClose == nil {
		pcontext.Errorf(ctx, "Internal Error in file service: %d", fid)
		return int32(file.FileErr_InternalError)
	}

	// Validate the requested buffer size
	buf := req.GetBuf()
	if !isValidBuf(buf) {
		pcontext.Errorf(ctx, "the expected buffer size %d exceeds the maximum buffer"+
			"size (%d) allowed", len(buf), maxBufSize)
		return int32(file.FileErr_LargeBufError)
	}

	n, _ := myFileInfo.Read(buf)

	myFileInfo.lastAccessTime = pcontext.CurrentTime(ctx)

	resp.Id = req.GetId()
	resp.NumRead = int32(n)

	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) delete(ctx context.Context, req *file.DeleteRequest,
	resp *file.DeleteResponse) int32 {

	// fid := file.UnmarshalFileId(req.GetId())

	fpath := req.GetPath()
	// Validate the file path
	cleanPath, valid := isValidPath(fpath)
	if !valid {
		pcontext.Errorf(ctx, "Invalid file path: %s", fpath)
		return int32(file.FileErr_InvalidPathError)
	}

	resp.Path = cleanPath

	// Validate if the file exists in the cache
	fid, exist := (*f.fpathTofid)[cleanPath]
	if !exist {
		pcontext.Errorf(ctx, "File does not exist, cannot be deleted: %s", cleanPath)
		return int32(file.FileErr_NotExistError)
	}

	// Ensure the file is not currently in use
	if (*f.fileDataCache)[fid].status != Fs_Close {
		pcontext.Errorf(ctx, "File is in use, cannot be deleted currently: %d", fid)
		return int32(file.FileErr_AlreadyInUseError)
	}

	// If testing, simply delete the file from the cache
	// otherwise, delete it from the disk
	if f.isTesting {
		delete(*f.fileDataCache, fid)
		delete(*f.fpathTofid, fpath)
	} else {
		deleteFileAndParentDirIfNeeded(fpath)
	}

	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) write(ctx context.Context, req *file.WriteRequest,
	resp *file.WriteResponse) int32 {

	fid := file.UnmarshalFileId(req.GetId())
	fileDataCache := *f.fileDataCache

	// Validate if the file exists in the cache
	if _, exist := fileDataCache[fid]; !exist {
		pcontext.Errorf(ctx, "File does not exist, cannot be write: %d", fid)
		return int32(file.FileErr_NotExistError)
	}

	myFileInfo := (*f.fileDataCache)[fid]

	// Verify file status to prevent errors during operation.
	// Only files with "write" status can be processed.
	switch myFileInfo.status {
	// If file status is "closed", log an error and return a file closed error code.
	case Fs_Close:
		pcontext.Errorf(ctx, "Operation aborted. File with ID: %d is closed.", fid)
		return int32(file.FileErr_FileClosedError)
	// If file status is "read", meaning it is currently being read by others,
	// log an error and return a file already in use error code.
	case Fs_Read:
		pcontext.Errorf(ctx, "Operation aborted. File with ID: %d is being read by others.", fid)
		return int32(file.FileErr_AlreadyInUseError)
	}

	// Check writer initialization, we cannot write a file without writer
	if myFileInfo.wrClose == nil {
		pcontext.Errorf(ctx, "Internal Error in file service: %d", fid)
		return int32(file.FileErr_InternalError)
	}

	// Validate the requested buffer size
	buf := req.GetBuf()
	// Check is size of buf exceeds the maximum buffer size allowed
	if !isValidBuf(buf) {
		pcontext.Errorf(ctx, "the buffer size %d exceeds the maximum buffer"+
			"size (%d) allowed", len(buf), maxBufSize)
		return int32(file.FileErr_LargeBufError)
	}

	n, _ := myFileInfo.Write(buf)

	myFileInfo.lastAccessTime = pcontext.CurrentTime(ctx)

	resp.Id = req.GetId()
	resp.NumWrite = int32(n)
	return int32(file.FileErr_NoError)
}

//
// helper functions
//

// A helper function to create a new file in the file service.
func (f *fileSvcImpl) createANewFile(fpath string, fcontent string) file.FileId {
	if f.isTesting {
		f.defaultCreateHook = createHookForStrings
	}

	currentTime := pcontext.CurrentTime(f.ctx)
	fid := file.NewFileId()

	newFileInfo := fileInfo{
		id:             fid,
		path:           fpath,
		content:        fcontent,
		status:         Fs_Write,
		createDate:     currentTime,
		lastAccessTime: currentTime,

		wrClose: f.defaultCreateHook(fpath),
	}

	newFileInfo.Write([]byte(fcontent))

	(*f.fileDataCache)[fid] = &newFileInfo
	(*f.fpathTofid)[fpath] = fid

	return fid
}

func (f *fileSvcImpl) loadTestData(ctx context.Context, req *file.LoadTestDataRequest,
	resp *file.LoadTestDataResponse) int32 {

	hostPath, dir := req.GetDirPath(), req.GetMountLocation()

	// Validate the directory path in the host machine
	if !isValidDirectory(hostPath) {
		pcontext.Errorf(ctx, "Invalid directory path: %s", hostPath)
		return int32(file.FileErr_NotExistError)
	}

	// Validate the directory path in the mount location
	cleanPath, valid := isValidPath(dir)
	if !valid {
		pcontext.Errorf(ctx, "Invalid directory path: %s", dir)
		return int32(file.FileErr_InvalidPathError)
	}

	resp.ErrorPath = []string{}

	if f.handleLoadFilesFromHost(hostPath, cleanPath, req, resp) {
		return int32(file.FileErr_NoDataFoundError)
	}

	return int32(file.FileErr_NoError)
}

// helper function

// recursively load all files from the directory in the host machine
// if such file already exists in the cache, rewrite it
// otherwise, create a new file in the cache
func (f *fileSvcImpl) handleLoadFilesFromHost(hostPath string, cleanPath string, req *file.LoadTestDataRequest,
	resp *file.LoadTestDataResponse) bool {

	isEmptyRoot := true

	err := filepath.WalkDir(hostPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return handleLoadTestFileError(path, req, resp, err)
		}
		if !d.IsDir() {
			isEmptyRoot = false
			if err := f.processFile(path, cleanPath, req); err != nil {
				return handleLoadTestFileError(path, req, resp, err)
			}
		}
		return nil
	})

	// If returnOnFail is false and there is an error in reading files,
	// program will be panic immediately and there is no need to return the error to the client.
	if err != nil {
		log.Fatalf("Error in loading test data: %v", err)
	}

	return isEmptyRoot
}
func (f *fileSvcImpl) processFile(hostPath string, cleanPath string, req *file.LoadTestDataRequest) error {
	file, err := os.Open(hostPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fcontent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	workPath := filepath.Join(cleanPath, hostPath[len(req.GetDirPath()):])

	// If the file already exists in the cache, rewrite it
	if fid, exist := (*f.fpathTofid)[workPath]; exist {
		delete(*f.fileDataCache, fid)
	}
	f.createANewFile(workPath, string(fcontent))

	return nil
}

func handleLoadTestFileError(path string, req *file.LoadTestDataRequest, resp *file.LoadTestDataResponse, err error) error {
	if req.GetReturnOnFail() {
		resp.ErrorPath = append(resp.ErrorPath, path)
		return nil
	}
	return err
}
