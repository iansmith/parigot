package file

import (
	"bytes"
	"io"
	"log"
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
func openHookForStrings(str string) io.ReadCloser {
	return NewStringReaderWrapper(strings.NewReader(str))
}

func openHookForFiles(path string) io.ReadCloser {
	realPath := getRealPath(path)
	f, err := os.Open(realPath)
	if err != nil {
		log.Fatal("Error form opening a file: ", err)
	}
	return f
}

type OpenHook func(pathOrString string) io.ReadCloser

var defaultOpenHook OpenHook = openHookForFiles

// Create hook
func createHookForStrings(str string) io.WriteCloser {
	return NewBytesBufferWrapper(&bytes.Buffer{})
}

func createHookForFiles(path string) io.WriteCloser {
	realPath := getRealPath(path)
	// If it does not exist, recursively create the directory
	err := os.MkdirAll(filepath.Dir(realPath), 0755)
	if err != nil {
		log.Fatal("Error creating directory:", err)
	}

	f, err := os.Create(realPath)
	if err != nil {
		log.Fatal("Error from creating a file: ", err)
	}

	return f
}

type CreateHook func(path string) io.WriteCloser

var defaultCreateHook CreateHook = createHookForFiles
