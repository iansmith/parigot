package lib

import (
	"bytes"
)

type perrorImpl struct {
	underlying error
	message    string
	id         Id
}

// Error is the type that is return from all parigot calls. If there is
// a well-known error code, it will be in the Id() method of this value
// otherwise the Id() method returns nil.
type Error interface {
	Error() string
	Unwrap() error
	Id() Id
}

// NewPerrorFromError wraps a given error with a PError and a message.
func NewPerrorFromError(msg string, err error) *perrorImpl {
	return &perrorImpl{
		underlying: err,
		message:    msg,
	}
}

// NewPerrorFromId creates a PError that contains a given id that must
// be an error code. Note that this message panics if you supply either
// an id that is not of an error code type or is an error code, but has
// a value that is "no error" (IsError() returns false).
func NewPerrorFromId(msg string, idv Id) *perrorImpl {
	if idv.IsError() == false {
		panic("unexpected usage of id; tried to use it as an error id but it contains no error")
	}
	return &perrorImpl{
		message: msg + idv.Short(),
	}
}

// NewPerror just creates a Perror object from the given string.  This function
// should be used as little as possible, with either of the above functions
// being far better since their inside values are computer friendly.
func NewPerror(msg string) *perrorImpl {
	return &perrorImpl{message: msg}
}

// Unwrap returns the error that this PError is wrapped around.
func (e *perrorImpl) Unwrap() error {
	return e.underlying
}

// Error returns a human readable string about this error. If there is
// an internal, wrapped error it's string is included. If there is an
// internal id, its Short() value is included.
func (e *perrorImpl) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.message)
	if e.underlying != nil {
		buf.WriteString(e.underlying.Error())
	} else if e.id != nil {
		buf.WriteString(e.id.Short())
	}
	buf.WriteString("\n")
	return buf.String()
}
