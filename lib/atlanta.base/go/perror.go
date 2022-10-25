package lib

import (
	"bytes"
	"github.com/iansmith/parigot/lib/libint"
)

type PerrorImpl struct {
	underlying error
	message    string
}

type Error interface {
	Error() string
	Unwrap() error
}

func NewPerrorFromError(msg string, err error) *PerrorImpl {
	return &PerrorImpl{
		underlying: err,
		message:    msg,
	}
}

func NewPerrorFromId(msg string, idv libint.Id) *PerrorImpl {
	return &PerrorImpl{
		message: msg + idv.Short(),
	}
}

func NewPerror(msg string) *PerrorImpl {
	return &PerrorImpl{message: msg}
}

func (e *PerrorImpl) Unwrap() error {
	return e.underlying
}

func (e *PerrorImpl) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.message + "\n")
	if e.underlying != nil {
		buf.WriteString(e.underlying.Error())
	}
	return buf.String()
}
