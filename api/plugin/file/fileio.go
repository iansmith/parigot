package file

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// String readcloser
type stringReaderWrapper struct {
	io.Reader
}

func newStringReaderWrapper(r io.Reader) *stringReaderWrapper {
	return &stringReaderWrapper{r}
}

func (s *stringReaderWrapper) Close() error {
	// We do nothing here, since the string is some data in memory
	// and the error is initialized to no-error
	return nil
}

// Buffer writecloser
type bytesBufferWrapper struct {
	*bytes.Buffer
}

func newBytesBufferWrapper(b *bytes.Buffer) *bytesBufferWrapper {
	return &bytesBufferWrapper{b}
}

func (b *bytesBufferWrapper) Close() error { return nil }

// Open hook
func openHookForStrings(str string) (io.ReadCloser, error) {
	return newStringReaderWrapper(strings.NewReader(str)), nil
}

func openHookForFiles(path string) (io.ReadCloser, error) {
	realPath, err := getRealPath(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(realPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type OpenHook func(pathOrString string) (io.ReadCloser, error)

// Create hook
func createHookForStrings(str string) (io.WriteCloser, error) {
	return newBytesBufferWrapper(&bytes.Buffer{}), nil
}

func createHookForFiles(path string) (io.WriteCloser, error) {
	realPath, err := getRealPath(path)
	if err != nil {
		return nil, err
	}
	// If it does not exist, recursively create the directory
	err = os.MkdirAll(filepath.Dir(realPath), 0755)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(realPath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type createHook func(path string) (io.WriteCloser, error)
