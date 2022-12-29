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

	ilog "github.com/iansmith/parigot/api/logimpl/go_"
	pb "github.com/iansmith/parigot/api/proto/g/pb/file"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"
	"github.com/psanford/memfs"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FileSvcImpl struct {
	mem      *jspatch.WasmMem
	fs       *memfs.FS
	idToFile map[string] /*really lib.Id*/ *os.File
}

// This is the native code side of the file service.  It reads the payload sent by the wasm world.

// FileSvcOpen opens a file which passed (via the trap mechanism) to this function. Note that the path
// will be checked for some kinds of controls.
//
//go:noinline
func (l *FileSvcImpl) FileSvcOpen(sp int32) {
	req := pb.OpenRequest{}
	err := splitutil.StackPointerToRequest(l.mem, sp, &req)
	if err != nil {
		return
	}
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "FileSvcOpen path to file %s", req.GetPath())
	newPath, err := ValidatePathForParigot(req.GetPath(), "open")
	if err != nil {
		splitutil.ErrorResponse(l.mem, sp, lib.KernelBadPath)
		return
	}
	// newpath can be different if there is something like /app/foo/bar/../baz as the parameter
	fp, err := os.Open(newPath)
	if err != nil {
		print("XXX file svc open failed,", err.Error(), "\n")
		splitutil.ErrorResponse(l.mem, sp, lib.KernelNotFound)
		return
	}
	fileId := lib.NewId[*protosupport.FileId]()
	marshaledId := lib.Marshal[protosupport.FileId](fileId)
	resp := pb.OpenResponse{Path: req.GetPath(), Id: marshaledId}
	if l.idToFile == nil {
		l.idToFile = make(map[string]*os.File)
	}
	l.idToFile[fileId.String()] = fp
	splitutil.RespondSingleProto(l.mem, sp, &resp)
	return
}

func (l *FileSvcImpl) SetWasmMem(ptr uintptr) {
	l.mem = jspatch.NewWasmMem(ptr)
}

func logger(level pblog.LogLevel, spec string, rest ...interface{}) {
	req := &pblog.LogRequest{
		Stamp:   timestamppb.New(time.Now()),
		Level:   level,
		Message: fmt.Sprintf(spec, rest...),
	}
	ilog.ProcessLogRequest(req, false, true, nil)
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
	return cleaned, nil
}

//go:noinline
func (l *FileSvcImpl) FileSvcLoad(sp int32) {
	req := pb.LoadRequest{}
	err := splitutil.StackPointerToRequest(l.mem, sp, &req)
	if err != nil {
		return
	}
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx FileSvcLoad path to LOCAL file %s", req.GetPath())
	l.fs = memfs.New()

	// implement semantics
	resp, err := l.loadLocal(&req)
	if err != nil {
		splitutil.ErrorResponse(l.mem, sp, lib.KernelNotFound /* xxxfixme, this error code is poor*/)
		return
	}
	// send the result home
	splitutil.RespondSingleProto(l.mem, sp, resp)
}

func (l *FileSvcImpl) loadLocal(req *pb.LoadRequest) (*pb.LoadResponse, error) {
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "file server load -1: %v -- %s\n", l.fs == nil, req.GetPath())
	stat, err := os.Stat(req.GetPath())
	if err != nil {
		return nil, err
	}
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "opened path on host machine: %s", req.GetPath())
	if !stat.IsDir() {
		return nil, &os.PathError{
			Op:   "read",
			Path: req.GetPath(),
			Err:  errors.New("path is a not a directory"),
		}
	}
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server load -3\n")
	memoryPrefix := "app" //no first slash for any call that uses io.fs.ValidPath()
	// start the import
	p := filepath.Join(memoryPrefix, req.GetPath())
	err = l.fs.MkdirAll(p, os.ModeDir)
	if err != nil {
		if !req.ReturnOnFail {
			panic("tried to mkdir all of " + p + ": " + err.Error())
		}
		return nil, err
	}

	// children is a flattened list of all child files,but does include directories
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "about to run recursive call: %s", req.GetPath())
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server load -4\n")
	children, badpath, err := l.readDirContents(req.GetPath(), req.GetReturnOnFail())
	if err != nil {
		if !req.ReturnOnFail {
			panic("tried to load start at path, but got: " + err.Error())
		}
		return nil, err
	}
	// walk child list
	for _, child := range children {
		logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server load -5, child:%s\n", child)
		//check for dir, we need to create those, not copy
		stat, err := os.Stat(child)
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server creating dir -5a: %s\n", child)
			l.fs.MkdirAll(filepath.Join(memoryPrefix, child), os.ModeDir)
			continue
		}
		// make sure in memory FS has the directory(ies) we need
		memPath := filepath.Join(memoryPrefix, child)
		logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server load copying -6 %s, and mem path %s\n", child, memPath)
		fp, err := os.Open(child)
		if err != nil {
			return nil, err
		}
		// because this is only available in test, we don't bother trying to be clever here to limit size of readAll()
		all, err := io.ReadAll(fp)
		if err != nil {
			return nil, err
		}
		logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server load copying -7 and read len is %d\n", len(all))
		stat, err = os.Stat(child)
		if err != nil {
			return nil, err
		}
		perm := stat.Mode()
		logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server about to write %s with perm %s, we are trying %s ", memPath, perm.String(),
			os.FileMode(0).String())
		err = l.fs.WriteFile(memPath, all, perm)
		if err != nil {
			logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "failed to write file %s: %v", memPath, err)
			return nil, err
		}
		logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server wrote all bytes")
	}
	var resp pb.LoadResponse
	resp.ErrorPath = badpath
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "xxx-file server .. number of bad paths: %d, %+v", len(resp.ErrorPath), resp.ErrorPath)
	return &resp, nil
}

// readDirContents is run recursively to read all the contents of all nested directories of path.  It returns
// the list of paths, list of failed paths and nil, or nil, nil and error.  If return on Fail is false, then
// we stop at the first error and return it with no error paths.  If return on Fail is true, we ignore errors
// and return the paths that generated errors.
func (l *FileSvcImpl) readDirContents(path string, returnOnFail bool) ([]string, []string, error) {
	logger(pblog.LogLevel_LOG_LEVEL_DEBUG, "reached readDirContents: %s", path)
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
