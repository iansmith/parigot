package file

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/file/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tetratelabs/wazero/api"
)

var filelogger = slog.Default().With("source", "file", "plugin", "true")

// make sure edit the makefile so you can have FileId and FileErrId, just like queue
// you'll need to pick the short names and letters for them... I would
// recommend f for FileId and F for FileErrId, but you can choose
// others if you want.

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
	id         file.FileId
	path       string
	size       int // length of the file content
	content    string
	status     FileStatus
	createTime time.Time
	modTime    time.Time

	rdClose io.ReadCloser
	wrClose io.WriteCloser
}

// for now, create a map of FileId -> to myFileInfo
// var fileDataCache = make(map[Fileid])*myFileInfo

type fileSvcImpl struct {
	fileDataCache *map[file.FileId]*fileInfo
	// track fid based on file path
	fpathTofid *map[string]file.FileId
	isTesting  bool

	defaultOpenHook   openHook
	defaultCreateHook createHook
}

// enum for file status
type FileStatus int

const (
	Fs_Write FileStatus = iota
	Fs_Read
	Fs_Close
)

func (fs FileStatus) String() string {
	return []string{"WRITE", "READ", "CLOSE"}[fs]
}

func (*FilePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "file", "open_", openFileHost) // this should call the "wrapper"
	e.AddSupportedFunc(ctx, "file", "create_", createFileHost)
	e.AddSupportedFunc(ctx, "file", "close_", closeFileHost)
	e.AddSupportedFunc(ctx, "file", "read_", readFileHost)
	e.AddSupportedFunc(ctx, "file", "delete_", deleteFileHost)
	e.AddSupportedFunc(ctx, "file", "write_", writeFileHost)
	e.AddSupportedFunc(ctx, "file", "load_test_data_", loadTestDataHost)
	e.AddSupportedFunc(ctx, "file", "stat_", statHost)

	_ = newFileSvc(ctx)

	return true
}

var parigotTime = "parigot_time"
var parigotTimeIndex = "parigot_time_index"

var localTimeZone *time.Location

// xxx fixme(iansmith) should this be here? is it clear we want
// to use the same mechanism as lib.CurrentTime()?  note this does
// not deal with local timezone!  This is running on the server side, so xxx?
func currentTimeHost(ctx context.Context) (context.Context, time.Time) {
	if localTimeZone == nil {
		var err error
		localTimeZone, err = time.LoadLocation("")
		if err != nil {
			panic("unable to load UTC time")
		}
	}
	if ctx == nil {
		return ctx, time.Now()
	}
	raw := ctx.Value(parigotTime)
	if raw != nil {
		t, ok := raw.([]time.Time)
		if ok {
			i := ctx.Value(parigotTimeIndex).(int)
			next := (i + 1) % (len(t))
			result := context.WithValue(ctx, parigotTimeIndex, next)
			now := t[i]
			return result, now
		}
	}
	// we got no value, just use clock
	return ctx, time.Now().In(localTimeZone)
}

func (fi *fileInfo) Read(p []byte) (n int, err error) {
	return fi.rdClose.Read(p)
}

func (fi *fileInfo) Write(p []byte) (n int, err error) {
	n, err = fi.wrClose.Write(p)
	if err != nil {
		return 0, err
	}

	fi.size += n

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

func loadTestDataHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.LoadTestDataRequest{}
	resp := &file.LoadTestDataResponse{}

	hostBase(ctx, "[file]loadTestData", fileSvc.loadTestData, m, stack, req, resp)
}

func statHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &file.StatRequest{}
	resp := &file.StatResponse{}

	hostBase(ctx, "[file]stat", fileSvc.stat, m, stack, req, resp)
}

func newFileSvc(ctx context.Context) *fileSvcImpl {
	f := &fileSvcImpl{
		fileDataCache: &map[file.FileId]*fileInfo{},
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
		filelogger.Error("File path is not valid", "path", fpath)
		return int32(file.FileErr_InvalidPathError)
	}
	resp.Path = cleanPath

	// Check if the file exists
	fid, exist := (*f.fpathTofid)[cleanPath]
	if !exist {
		filelogger.Error("File does not exist and cannot be opened", "path", fpath)
		return int32(file.FileErr_NotExistError)
	}

	// If file exists, fetch its information from fileDataCache
	myFileInfo := (*f.fileDataCache)[fid]

	// Check file status
	if myFileInfo.status != Fs_Close {
		filelogger.Error("File is already in use", "path", fpath)
		return int32(file.FileErr_AlreadyInUseError)
	}

	// If the file is closed, update the file information and open it for reading
	resp.Id = fid.Marshal()

	ctx, myFileInfo.modTime = currentTimeHost(ctx)
	myFileInfo.status = Fs_Read

	var err error
	myFileInfo.rdClose, err = f.defaultOpenHook(cleanPath)
	if err != nil {
		filelogger.Error("Failed to open file for reading", "path", fpath)
		return int32(file.FileErr_OpenError)
	}

	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) create(ctx context.Context, req *file.CreateRequest,
	resp *file.CreateResponse) int32 {

	fpath := req.GetPath()
	// Validate the file path
	cleanPath, valid := isValidPath(fpath)
	if !valid {
		filelogger.Error("Invalid file path", "path", fpath)
		return int32(file.FileErr_InvalidPathError)
	}

	resp.Path = cleanPath
	resp.Truncated = false

	// Validate the content buffer
	content := req.GetContent()
	buf := []byte(content)
	if !isValidBuf(buf) {
		filelogger.Error("Content size exceeds the maximum buffer size", "content size", len(buf), "maximum buffer size", maxBufSize)
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
			filelogger.Error("file is in use", "path", fpath)
			return int32(file.FileErr_AlreadyInUseError)
		}

		// Extend the file
		myFileInfo.content += content
		ctx, myFileInfo.modTime = currentTimeHost(ctx)
		myFileInfo.status = Fs_Write
		myFileInfo.Write(buf)

	} else {
		// If file/path does not exist, create a new file
		fid, err := f.createANewFile(ctx, cleanPath, content)
		if err != nil {
			filelogger.Error("Failed to create a new file", "path", fpath)
			return int32(file.FileErr_CreateError)
		}
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
		filelogger.Error("File does not exist, cannot be closed", "file number", fid)
		return int32(file.FileErr_NotExistError)
	}

	// Check if the file is already closed
	if fileData.status == Fs_Close {
		filelogger.Error("File is already closed, cannot be closed again", "file number", fid)
		return int32(file.FileErr_FileClosedError)
	}

	status := fileData.status

	// If the file exists and is not closed, change its status to "closed"
	(*f.fileDataCache)[fid].status = Fs_Close
	ctx, (*f.fileDataCache)[fid].modTime = currentTimeHost(ctx)

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
		filelogger.Error("File does not exist, cannot be read", "file number", fid)
		return int32(file.FileErr_NotExistError)
	}

	myFileInfo := (*f.fileDataCache)[fid]

	// Verify file status to prevent errors during operation.
	// Only files with "read" status can be processed.
	switch myFileInfo.status {
	// If file status is "closed", log an error and return a file closed error code.
	case Fs_Close:
		filelogger.Error("Operation aborted. File is closed.", "file number", fid)
		return int32(file.FileErr_FileClosedError)
	// If file status is "write", meaning it is currently being written by others,
	// log an error and return a file already in use error code.
	case Fs_Write:
		filelogger.Error("Operation aborted. File is being written by others.", "file number", fid)
		return int32(file.FileErr_AlreadyInUseError)
	}

	// Check if the file's reader is initialized
	if myFileInfo.rdClose == nil {
		filelogger.Error("Internal Error in file service", "file naumber", fid)
		return int32(file.FileErr_InternalError)
	}

	// Validate the requested buffer size
	buf := req.GetBuf()
	if !isValidBuf(buf) {
		filelogger.Error("The expected buffer size exceeds the maximum buffer",
			"buffer size", len(buf), "max buffer size", maxBufSize)
		return int32(file.FileErr_LargeBufError)
	}

	n, err := myFileInfo.Read(buf)
	if err != nil && err != io.EOF {
		filelogger.Error("Failed to read", "file number", fid)
		return int32(file.FileErr_ReadError)
	}

	ctx, myFileInfo.modTime = currentTimeHost(ctx)

	resp.Id = req.GetId()
	resp.NumRead = int32(n)

	if err == io.EOF {
		return int32(file.FileErr_EOFError)
	}
	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) delete(ctx context.Context, req *file.DeleteRequest,
	resp *file.DeleteResponse) int32 {

	fpath := req.GetPath()
	// Validate the file path
	cleanPath, valid := isValidPath(fpath)
	if !valid {
		filelogger.Error("Invalid file path", "path", fpath)
		return int32(file.FileErr_InvalidPathError)
	}

	resp.Path = cleanPath

	// Validate if the file exists in the cache
	fid, exist := (*f.fpathTofid)[cleanPath]
	if !exist {
		filelogger.Error("File does not exist, cannot be deleted", "path", cleanPath)
		return int32(file.FileErr_NotExistError)
	}

	// Ensure the file is not currently in use
	if (*f.fileDataCache)[fid].status != Fs_Close {
		filelogger.Error("File is in use, cannot be deleted currently", "file number", fid)
		return int32(file.FileErr_AlreadyInUseError)
	}

	// If testing, simply delete the file from the cache
	// otherwise, delete it from the disk
	if f.isTesting {
		delete(*f.fileDataCache, fid)
		delete(*f.fpathTofid, fpath)
	} else {
		err := deleteFileAndParentDirIfNeeded(fpath)
		if err != nil {
			filelogger.Error("Error deleting file", "error", err)
			return int32(file.FileErr_DeleteError)
		}
	}

	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) write(ctx context.Context, req *file.WriteRequest,
	resp *file.WriteResponse) int32 {

	fid := file.UnmarshalFileId(req.GetId())
	fileDataCache := *f.fileDataCache

	// Validate if the file exists in the cache
	if _, exist := fileDataCache[fid]; !exist {
		filelogger.Error("File does not exist, cannot be written to", "file number", fid)
		return int32(file.FileErr_NotExistError)
	}

	myFileInfo := (*f.fileDataCache)[fid]

	// Verify file status to prevent errors during operation.
	// Only files with "write" status can be processed.
	switch myFileInfo.status {
	// If file status is "closed", log an error and return a file closed error code.
	case Fs_Close:
		filelogger.Error("Operation aborted, file with ID, file is closed", "file id", fid.Short())
		return int32(file.FileErr_FileClosedError)
	// If file status is "read", meaning it is currently being read by others,
	// log an error and return a file already in use error code.
	case Fs_Read:
		filelogger.Error("Operation aborted, file is being read by others", "file id", fid.Short())
		return int32(file.FileErr_AlreadyInUseError)
	}

	// Check writer initialization, we cannot write a file without writer
	if myFileInfo.wrClose == nil {
		filelogger.Error("Internal Error in file service", "file id", fid.Short())
		return int32(file.FileErr_InternalError)
	}

	// Validate the requested buffer size
	buf := req.GetBuf()
	// Check is size of buf exceeds the maximum buffer size allowed
	if !isValidBuf(buf) {
		filelogger.Error("the buffer size exceeds the maximum buffer"+
			"size allowed", "buffer size", len(buf), "max buffer size", maxBufSize)
		return int32(file.FileErr_LargeBufError)
	}

	n, err := myFileInfo.Write(buf)
	if err != nil {
		filelogger.Error("Failed to write", "file id", fid.Short(), "error", err)
		return int32(file.FileErr_WriteError)
	}

	ctx, myFileInfo.modTime = currentTimeHost(ctx)

	resp.Id = req.GetId()
	resp.NumWrite = int32(n)
	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) loadTestData(ctx context.Context, req *file.LoadTestDataRequest,
	resp *file.LoadTestDataResponse) int32 {

	hostPath, dir := req.GetDirPath(), req.GetMountLocation()

	// Validate the directory path in the host machine
	if valid, err := isValidDirOnHost(hostPath); !valid {
		filelogger.Error("Invalid directory path", "path", hostPath, "error", err)
		return int32(file.FileErr_NotExistError)
	}

	// Validate the directory path in the mount location
	cleanPath, valid := isValidPath(dir)
	if !valid {
		filelogger.Error("Invalid directory path", "path", dir)
		return int32(file.FileErr_InvalidPathError)
	}

	resp.ErrorPath = []string{}

	if !f.loadFilesFromHost(ctx, hostPath, cleanPath, req, resp) {
		return int32(file.FileErr_NoDataFoundError)
	}

	return int32(file.FileErr_NoError)
}

func (f *fileSvcImpl) stat(ctx context.Context, req *file.StatRequest,
	resp *file.StatResponse) int32 {

	fpath := req.GetPath()
	// Validate the file path, host can only operate on files have specific prefix
	cleanPath, valid := isValidPath(fpath)
	if !valid {
		filelogger.Error("Invalid file path", "path", fpath)
		return int32(file.FileErr_InvalidPathError)
	}

	// Validate if the file/dir exists in the cache
	fid, exist := (*f.fpathTofid)[cleanPath]
	if exist {
		f.populateStat(resp, cleanPath, fid, false)
	} else {
		// If the file does not exist, check if the directory exists
		// If the directory exists, populate the directory stat
		for path, newFid := range *f.fpathTofid {
			if strings.HasPrefix(path, cleanPath) {
				// Populate directory stat
				f.populateStat(resp, cleanPath, newFid, true)
				exist = true
			}
		}
	}
	if !exist {
		filelogger.Error("Path does not exist, cannot compute stat call", "path", cleanPath)
		return int32(file.FileErr_NotExistError)
	}

	return int32(file.FileErr_NoError)
}

//
// helper functions
//

// A helper function to populate file/dir stat.
func (f *fileSvcImpl) populateStat(resp *file.StatResponse, cleanPath string, fid file.FileId, isDir bool) {
	fileStat := resp.GetFileInfo()

	fileStat.Path = cleanPath
	fileStat.Size += int32((*f.fileDataCache)[fid].size)
	fileStat.ModTime = updateTimestamp(fileStat.ModTime, (*f.fileDataCache)[fid].modTime, true)
	fileStat.CreateTime = updateTimestamp(fileStat.CreateTime, (*f.fileDataCache)[fid].createTime, false)
	fileStat.IsDir = isDir
}

func updateTimestamp(t1 *timestamppb.Timestamp, t2 time.Time, isLatest bool) *timestamppb.Timestamp {
	if isLatest == t1.AsTime().After(t2) {
		return t1
	}
	return timestamppb.New(t2)
}

// A helper function to create a new file in the file service.
func (f *fileSvcImpl) createANewFile(ctx context.Context, fpath string, fcontent string) (file.FileId, error) {
	if f.isTesting {
		f.defaultCreateHook = createHookForStrings
	}

	_, currentTime := currentTimeHost(ctx)
	fid := file.NewFileId()

	newFileInfo := fileInfo{
		id:         fid,
		path:       fpath,
		content:    fcontent,
		status:     Fs_Write,
		createTime: currentTime,
		modTime:    currentTime,
	}
	var err error
	newFileInfo.wrClose, err = f.defaultCreateHook(fpath)
	if err != nil {
		return file.FileIdEmptyValue(), err
	}

	newFileInfo.Write([]byte(fcontent))

	(*f.fileDataCache)[fid] = &newFileInfo
	(*f.fpathTofid)[fpath] = fid

	return fid, nil
}

// recursively load all files from the directory in the host machine
// if such file already exists in the cache, rewrite it
// otherwise, create a new file in the cache
func (f *fileSvcImpl) loadFilesFromHost(ctx context.Context, hostPath string, cleanPath string, req *file.LoadTestDataRequest,
	resp *file.LoadTestDataResponse) bool {

	getAnyData := false

	err := filepath.WalkDir(hostPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return handleLoadTestFileError(path, req, resp, err)
		}
		if !d.IsDir() {
			getAnyData = true
			if err := f.processFile(ctx, path, cleanPath, req); err != nil {
				return handleLoadTestFileError(path, req, resp, err)
			}
		}
		return nil
	})

	// If returnOnFail is false and there is an error in reading files,
	// program will be panic immediately and there is no need to return the error to the client.
	if err != nil {
		panic(fmt.Sprintf("error loading test data: %v", err))
	}

	return getAnyData
}
func (f *fileSvcImpl) processFile(ctx context.Context, hostPath string, cleanPath string, req *file.LoadTestDataRequest) error {
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
	f.createANewFile(ctx, workPath, string(fcontent))

	return nil
}

func handleLoadTestFileError(path string, req *file.LoadTestDataRequest, resp *file.LoadTestDataResponse, err error) error {
	if req.GetReturnOnFail() {
		resp.ErrorPath = append(resp.ErrorPath, path)
		return nil
	}
	return err
}
