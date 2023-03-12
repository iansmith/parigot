//go:build !js
// +build !js

package go_

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/iansmith/parigot/apiimpl/splitutil"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/backdoor"
	"github.com/iansmith/parigot/sys/jspatch"

	"github.com/psanford/memfs"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FileSvcImpl struct {
	mem             *jspatch.WasmMem
	fs              *memfs.FS
	idToFilePointer map[string] /*really string version of lib.Id*/ int64
	idToMemPath     map[string]string
}

// This is the native code side of the file service.  It reads the payload sent by the wasm world.

// FileSvcOpen opens a file which passed (via the trap mechanism) to this function. Note that the path
// will be checked for some kinds of controls.
//
//go:noinline
func (l *FileSvcImpl) FileSvcOpen(sp int32) {
	req := filemsg.OpenRequest{}
	errId, errDetail := splitutil.StackPointerToRequest(l.mem, sp, &req)
	if errId != nil {
		splitutil.ErrorResponse(l.mem, sp, errId, errDetail)
	}
	logger(logmsg.LogLevel_LOG_LEVEL_INFO, "FileSvcOpen path to file %s", req.GetPath())
	newPath, err := ValidatePathForParigot(req.GetPath(), "open")
	if err != nil {
		splitutil.ErrorResponse(l.mem, sp, lib.NewFileError(lib.FileBadPath),
			"invalid path:"+req.GetPath())
		return
	}
	// newpath can be quite different if there is something like /app/foo/bar/../baz as the parameter
	_, err = fs.ReadFile(l.fs, newPath)
	if err != nil {
		splitutil.ErrorResponse(l.mem, sp, lib.NewFileError(lib.FileBadPath),
			fmt.Sprintf("read file failed on %s: %v", req.GetPath(), err))
		return
	}
	fileId := lib.NewId[*protosupportmsg.FileId]()
	marshaledId := lib.Marshal[protosupportmsg.FileId](fileId)
	resp := filemsg.OpenResponse{Path: req.GetPath(), Id: marshaledId}
	if l.idToFilePointer == nil {
		l.idToFilePointer = make(map[string]int64)
		l.idToMemPath = make(map[string]string)
	}
	l.idToFilePointer[fileId.String()] = 0
	l.idToMemPath[fileId.String()] = newPath
	splitutil.RespondSingleProto(l.mem, sp, &resp)
}

func (l *FileSvcImpl) SetWasmMem(ptr uintptr) {
	l.mem = jspatch.NewWasmMem(ptr)
}

func logger(level logmsg.LogLevel, spec string, rest ...interface{}) {
	req := &logmsg.LogRequest{
		Stamp:   timestamppb.New(time.Now()),
		Level:   level,
		Message: fmt.Sprintf(spec, rest...),
	}
	backdoor.Log(req, false, true, false, nil)
}

// ValidatePathForParigot checks for and avoids many of the common pitfalls in pathnames but is not
// a substitute for chroot() or similar types jailing of user code implemented by the underlying operating
// systems.
func ValidatePathForParigot(path string, op string) (string, error) {

	cleaned := filepath.Clean(path)
	noFrontSlash := cleaned
	if strings.HasPrefix(cleaned, "/") {
		noFrontSlash = noFrontSlash[1:]
	} else {
		return "", &fs.PathError{
			Op:   op,
			Path: path,
			Err:  errors.New("parigot requires fully qualified path names"),
		}
	}
	if !fs.ValidPath(noFrontSlash) {
		return "", &fs.PathError{
			Op:   op,
			Path: path,
			Err:  errors.New("failed ValidPath() test"),
		}
	}
	part := filepath.SplitList(cleaned)
	if len(part) == 0 {
		return "", &fs.PathError{
			Op:   op,
			Path: path,
			Err:  errors.New("empty path"), //can this actually happen?
		}
	}
	for _, element := range part {
		for _, r := range element {
			if unicode.IsControl(r) {
				return "", &fs.PathError{
					Op:   op,
					Path: path,
					Err:  errors.New("control characters not allowed"),
				}
			}
		}
	}
	if !strings.HasPrefix(cleaned, "/") {
		panic("unable to understand path " + cleaned + " because not fully qualified?")
	}
	return cleaned[1:], nil
}

//go:noinline
func (l *FileSvcImpl) FileSvcLoad(sp int32) {
	req := filemsg.LoadTestRequest{}

	errId, errDetail := splitutil.StackPointerToRequest(l.mem, sp, &req)
	if errId != nil {
		splitutil.ErrorResponse(l.mem, sp, errId, errDetail)
		return
	}
	l.fs = memfs.New()

	// implement semantics
	resp, err := l.loadLocal(&req)
	if err != nil {
		splitutil.ErrorResponse(l.mem, sp,
			lib.NewFileError(lib.FileNotFound),
			fmt.Sprintf("reading in-memory file %s:%v", req.GetPath(),
				err))
		return
	}
	// send the result home
	splitutil.RespondSingleProto(l.mem, sp, resp)
}

func (l *FileSvcImpl) loadLocal(req *filemsg.LoadTestRequest) (*filemsg.LoadTestResponse, error) {
	stat, err := os.Stat(req.GetPath())
	if err != nil {
		backdoor.Log(&logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_ERROR,
			Stamp:   timestamppb.Now(),
			Message: "load local file (" + req.Path + ") failed:" + err.Error(),
		}, false, true, false, nil)
		return nil, err
	}
	if !stat.IsDir() {
		return nil, &os.PathError{
			Op:   "read",
			Path: req.GetPath(),
			Err:  errors.New("path is a not a directory"),
		}
	}
	memoryPrefix := "app" //no first slash for any call that uses io.fs.ValidPath()
	// start the import
	p := filepath.Join(memoryPrefix, req.GetMountLocation())

	err = l.fs.MkdirAll(p, 0777)
	backdoor.Log(&logmsg.LogRequest{
		Level: logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Stamp: timestamppb.Now(),
		Message: fmt.Sprintf("loadLocal: considering file %s (GetMountPoint %s)",
			p, req.GetMountLocation()),
	}, false, true, false, nil)
	if err != nil {
		if !req.ReturnOnFail {
			panic(fmt.Sprintf("tried to mkdir all %s (from %s): %s ",
				p, req.GetMountLocation(), err.Error()))
		}
		return nil, err
	}

	// children is a flattened list of all child files,but does include directories
	children, badpath, err := l.readDirContents(req.GetPath(), req.GetReturnOnFail())
	if err != nil {
		if !req.ReturnOnFail {
			panic("tried to load start at path, but got: " + err.Error())
		}
		return nil, err
	}
	// walk child list
	for _, child := range children {
		backdoor.Log(&logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
			Stamp:   timestamppb.Now(),
			Message: fmt.Sprintf("loadLocal: considering file %s", child),
		}, false, true, false, nil)

		//check for dir, we need to create those, not copy
		stat, err := os.Stat(child)
		if err != nil {
			backdoor.Log(&logmsg.LogRequest{
				Level:   logmsg.LogLevel_LOG_LEVEL_ERROR,
				Stamp:   timestamppb.Now(),
				Message: fmt.Sprintf("loadLocal: error with file %s: %v", child, err.Error()),
			}, false, true, false, nil)
			return nil, err
		}
		if stat.IsDir() {
			err := l.fs.MkdirAll(filepath.Join(memoryPrefix, req.MountLocation, child), 0777)
			if err != nil {
				backdoor.Log(&logmsg.LogRequest{
					Level: logmsg.LogLevel_LOG_LEVEL_ERROR,
					Stamp: timestamppb.Now(),
					Message: fmt.Sprintf("loadLocal: error with file %s trying to stat %s: %v",
						child, filepath.Join(memoryPrefix, req.MountLocation, child), err.Error()),
				}, false, true, false, nil)
			}
			continue
		}
		// make sure in memory FS has the directory(ies) we need
		memPath := filepath.Join(memoryPrefix, req.MountLocation, child)
		backdoor.Log(&logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
			Stamp:   timestamppb.Now(),
			Message: fmt.Sprintf("loadLocal: memPath %s => %s", child, memPath),
		}, false, true, false, nil)

		fp, err := os.Open(child)
		if err != nil {
			return nil, err
		}
		// because this is only available in test, we don't bother trying to be clever here to limit size of readAll()
		all, err := io.ReadAll(fp)
		if err != nil {
			return nil, err
		}
		stat, err = os.Stat(child)
		if err != nil {
			return nil, err
		}
		//perm := stat.Mode()
		err = l.fs.WriteFile(memPath, all, 0755)
		if err != nil {
			return nil, err
		}
	}

	var resp filemsg.LoadTestResponse
	resp.ErrorPath = badpath
	return &resp, nil
}

// readDirContents is run recursively to read all the contents of all nested directories of path.  It returns
// the list of paths, list of failed paths and nil, or nil, nil and error.  If return on Fail is false, then
// we stop at the first error and return it with no error paths.  If return on Fail is true, we ignore errors
// and return the paths that generated errors.
func (l *FileSvcImpl) readDirContents(path string, returnOnFail bool) ([]string, []string, error) {
	stat, err := os.Stat(path) // sanity check
	if err != nil {
		if returnOnFail {
			return nil, []string{path}, nil
		}
		return nil, nil, err
	}
	if !stat.IsDir() { //sanity check
		if returnOnFail {
			return nil, []string{path}, nil
		}
		return nil, nil, &os.PathError{
			Op:   "read",
			Path: path,
			Err:  errors.New("path is a not a directory"),
		}
	}
	entry, err := os.ReadDir(path)
	if err != nil {
		if returnOnFail {
			return nil, []string{path}, nil
		}
		return nil, nil, err
	}
	localResult := []string{}
	localBadPath := []string{}
	for _, e := range entry {
		stat, err := os.Stat(filepath.Join(path, e.Name()))
		if err != nil {
			if returnOnFail {
				localBadPath = append(localBadPath, filepath.Join(path, e.Name()))
				continue
			}
			return nil, nil, err
		}
		if !stat.IsDir() {
			localResult = append(localResult, filepath.Join(path, e.Name()))
			continue
		}
		children, failure, err := l.readDirContents(filepath.Join(path, e.Name()), returnOnFail)
		if err != nil {
			return nil, nil, err // we already have bailed out somewhere below
		}
		localResult = append(localResult, children...)
		localBadPath = append(localBadPath, failure...)
	}
	return localResult, localBadPath, nil
}