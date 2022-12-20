package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/iansmith/parigot/api/fileimpl/go_"
	"github.com/iansmith/parigot/api/proto/g/file"
	"github.com/iansmith/parigot/api/proto/g/log"
	pb "github.com/iansmith/parigot/api/proto/g/pb/file"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"

	"github.com/psanford/memfs"
	"google.golang.org/protobuf/proto"
)

func main() {
	// we export and require services before the call to file.Run()... our call to the Run() system call is in Ready()
	if _, err := lib.Export1("file", "File"); err != nil {
		panic("ready: error in attempt to export api.Log: " + err.Error())
	}
	if _, err := lib.Require1("log", "Logger"); err != nil {
		print("ready: error in attempt to export api.Log: " + err.Error())
	}

	file.Run(&myFileServer{})
}

type myFileServer struct {
	logger log.Log
}

func (m *myFileServer) Ready() bool {
	if _, err := lib.Run(false); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	return true

}

// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.
// This code receives an encoded single protobuf object from the wasm-side client code.
func (m *myFileServer) Open(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	u, err := splitutil.SendSingleProto(inProto)
	if err != nil {
		return nil, err
	}
	go_.FileSvcOpen(int32(u))
	return nil, nil
}

func (m *myFileServer) Close(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*pb.CloseRequest)
	panic("Close")
}
func (m *myFileServer) Create(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*pb.CreateRequest)
	panic("Create")
}
func (m *myFileServer) log(pctx *protosupport.Pctx, spec string, rest ...interface{}) {
	s := fmt.Sprintf(spec, rest...)
	req := pblog.LogRequest{
		Stamp: pctx.GetNow(),

		Level:   pblog.LogLevel_LOG_LEVEL_DEBUG,
		Message: "myFileServer:" + s,
	}
	err := m.logger.Log(&req)
	if err != nil {
		print("unable to log ", s, ":", err.Error(), "\n")
	}
}

func (m *myFileServer) Load(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	data := inProto.(*pb.LoadRequest)
	fs := memfs.New()

	stat, err := os.Stat(data.GetPath())
	if err != nil {
		return nil, err
	}
	m.log(pctx, "opened path on host machine: %s", data.GetPath())
	if !stat.IsDir() {
		return nil, &os.PathError{
			Op:   "read",
			Path: data.GetPath(),
			Err:  errors.New("path is a not a directory"),
		}
	}
	memoryPrefix := "/app"
	// start the import
	fs.MkdirAll(memoryPrefix, os.FileMode(os.ModeDir))

	// children is a flattened list of all child files,but does include directories
	m.log(pctx, "about to run recursive call: %s", data.GetPath())
	children, badpath, err := m.readDirContents(pctx, data.GetPath(), data.GetReturnOnFail())
	if err != nil {
		return nil, err
	}
	// walk child list
	for _, child := range children {
		//check for dir
		stat, err := os.Stat(data.GetPath())
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			continue
		}
		// make sure in memory FS has the directory(ies) we need
		memPath := filepath.Join(memoryPrefix, child)
		dir, _ := filepath.Split(memPath)
		err = fs.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return nil, err
		}
		// read and copy bytes
		fp, err := os.Open(child)
		if err != nil {
			return nil, err
		}
		all, err := io.ReadAll(fp)
		stat, err = os.Stat(child)
		if err != nil {
			return nil, err
		}
		perm := stat.Mode()
		err = fs.WriteFile(memPath, all, perm)
		if err != nil {
			return nil, err
		}
	}
	var resp pb.LoadResponse
	resp.ErrorPath = badpath
	return &resp, nil
}

// readDirContents is run recursively to read all the contents of all nested directories of path.  It returns
// the list of paths, list of failed paths and nil, or nil, nil and error.  Only one error object is returned
// no matter how many bad paths there are.
func (m *myFileServer) readDirContents(pctx *protosupport.Pctx, path string, returnOnFail bool) ([]string, []string, error) {
	m.log(pctx, "reached readDirContents: %s", path)
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
		children, failure, err := m.readDirContents(pctx, filepath.Join(path, e.Name()), returnOnFail)
		if err != nil {
			return nil, nil, err // we already have bailed out somewhere below
		}
		localResult = append(localResult, children...)
		localBadPath = append(localBadPath, failure...)
	}
	return localResult, localBadPath, nil
}

func (m *myFileServer) validatePathForParigot(path, op string) (string, error) {
	return go_.ValidatePathForParigot(path, op)
}
