package wazero

import (
	"io"
	"io/fs"
	"os"
	"syscall"

	"github.com/tetratelabs/wazero/internal/sysfs"
)

type FauxFs interface {
	String() string
	Exists(path string) bool
	Open(path string, read, write bool) (FauxFile, syscall.Errno)
	Create(path string, canRead, canWrite bool) (FauxFile, syscall.Errno)
	Close(f FauxFile) syscall.Errno
}

type FauxFile interface {
	io.Writer
	Read([]byte) (int, error)
	Close() error
	Path() string
}

type fauxFs struct {
	sysfs.UnimplementedFS
	provider FauxFs
}

type fauxFile struct {
	parent            *fauxFs
	canRead, canWrite bool
	path              string
	ff                FauxFile
}

func (f *fauxFs) String() string {
	return f.provider.String()
}

func newFauxFs(p FauxFs) sysfs.FS {
	return &fauxFs{provider: p}
}

func newFauxFdFile(ff FauxFile, path string, parent *fauxFs, canRead, canWrite bool) *fauxFile {
	return &fauxFile{
		parent:   parent,
		canRead:  canRead,
		canWrite: canWrite,
		ff:       ff,
	}
}

func (f *fauxFile) Stat() (fs.FileInfo, error) {
	err := syscall.ENOSYS
	return nil, err
}
func (f *fauxFile) Read(b []byte) (int, error) {
	if !f.canRead {
		return 0, syscall.EPERM
	}
	return f.ff.Read(b)
}
func (f *fauxFile) Write(b []byte) (int, error) {
	if !f.canWrite {
		return 0, syscall.EPERM
	}
	return f.ff.Write(b)
}
func (f *fauxFile) Close() error {
	return f.ff.Close()
}

func (f *fauxFs) Close(ff fs.File) syscall.Errno {
	return f.provider.Close(ff.(FauxFile))
}

func (f *fauxFs) OpenFile(path string, flag int, perm fs.FileMode) (fs.File, syscall.Errno) {
	if !fs.ValidPath(path) { // FS.OpenFile has fewer constraints than fs.FS
		return nil, syscall.EINVAL
	}

	if flag&os.O_CREATE != 0 {
		if f.provider.Exists(path) {
			return nil, syscall.EEXIST
		}
		//print(fmt.Sprintf("In wazero.fauxfd, open file %s: flag %x, filemode %o, (create?%v, exists?%v)\n", path, flag, perm, flag&os.O_CREATE != 0, f.provider.Exists(path)))

		unixBits := perm.Perm()
		owner := unixBits & 0700
		owner >>= 6
		readBit := owner & 0x4
		writeBit := owner & 0x2
		// must pass at least on of reading and writing
		if readBit == 0 && writeBit == 0 {
			return nil, syscall.EINVAL
		}
		canRead := false
		canWrite := false
		if readBit != 0 {
			canRead = true
		}
		if writeBit != 0 {
			canRead = true
		}

		fd, err := f.provider.Create(path, canRead, canWrite)
		if err != 0 {
			return nil, err
		}
		return newFauxFdFile(fd, path, f, canRead, canWrite), 0
	}
	var fd FauxFile
	var err syscall.Errno
	readOk, writeOk := false, false
	switch flag {
	case os.O_RDONLY:
		fd, err = f.provider.Open(path, true, false)
		readOk = true
	case os.O_WRONLY:
		fd, err = f.provider.Open(path, false, true)
		writeOk = true
	case os.O_RDWR:
		fd, err = f.provider.Open(path, false, true)
		writeOk = true
		readOk = true
	default:
		return nil, syscall.EINVAL
	}
	if err != 0 {
		return nil, err
	}
	return newFauxFdFile(fd, path, f, readOk, writeOk), 0
}
