package file

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// String readcloser
type StringReaderWrapper struct {
	io.Reader
}

func NewStringReaderWrapper(r io.Reader) *StringReaderWrapper {
	return &StringReaderWrapper{r}
}

func (s *StringReaderWrapper) Close() error {
	// We do nothing here, since the string is some data in memory
	// and the error is initialized to no-error
	return nil
}

// Buffer writecloser
type BytesBufferWrapper struct {
	*bytes.Buffer
}

func NewBytesBufferWrapper(b *bytes.Buffer) *BytesBufferWrapper {
	return &BytesBufferWrapper{b}
}

func (b *BytesBufferWrapper) Close() error { return nil }

// Open hook
func openHookForStrings(str string) (io.ReadCloser, error) {
	return NewStringReaderWrapper(strings.NewReader(str)), nil
}

func openHookForFiles(path string) (io.ReadCloser, error) {
	realPath, err := getRealPath(path)
	if err != nil {
		// log.Fatal("Error opening a file: ", file.FileErr_InternalError.String())
		return nil, err
	}
	f, err := os.Open(realPath)
	if err != nil {
		// log.Fatal("Error form opening a file: ", err)
		return nil, err
	}
	return f, nil
}

type OpenHook func(pathOrString string) (io.ReadCloser, error)

// Create hook
func createHookForStrings(str string) (io.WriteCloser, error) {
	return NewBytesBufferWrapper(&bytes.Buffer{}), nil
}

func createHookForFiles(path string) (io.WriteCloser, error) {
	realPath, err := getRealPath(path)
	if err != nil {
		// log.Fatal("Error creating a file: ", file.FileErr_InternalError.String())
		return nil, err
	}
	// If it does not exist, recursively create the directory
	err = os.MkdirAll(filepath.Dir(realPath), 0755)
	if err != nil {
		// log.Fatal("Error creating directory:", file.FileErr_InternalError.String())
		return nil, err
	}

	f, err := os.Create(realPath)
	if err != nil {
		// log.Fatal("Error from creating a file: ", file.FileErr_InternalError.String())
		return nil, err
	}

	return f, nil
}

type CreateHook func(path string) (io.WriteCloser, error)
