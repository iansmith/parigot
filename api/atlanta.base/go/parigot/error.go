package parigot

import (
	"bytes"
)

type ErrorImpl struct {
	underlying error
	message    string
}

type Error interface {
	Error() string
	Unwrap() error
}

func NewErrorFromError(msg string, err error) *ErrorImpl {
	return &ErrorImpl{
		underlying: err,
		message:    msg,
	}
}

func NewErrorFromId(msg string, id AnyId) *ErrorImpl {
	return &ErrorImpl{
		message: msg + AsId(id),
	}
}

func NewError(msg string) *ErrorImpl {
	return &ErrorImpl{message: msg}
}

func (e *ErrorImpl) Unwrap() error {
	return e.underlying
}

func (e *ErrorImpl) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.message + "\n")
	if e.underlying != nil {
		buf.WriteString(e.underlying.Error())
	}
	return buf.String()
}
