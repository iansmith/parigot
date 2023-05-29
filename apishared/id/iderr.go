package id

type IdErrorIdDef struct{}

func (f IdErrorIdDef) ShortString() string { return "!!err!!" }
func (f IdErrorIdDef) Letter() byte        { return 0x21 } //!
func (f IdErrorIdDef) IsError() bool       { return true }

type IdRootErrorCode uint16

const (
	IdNoError         IdRootErrorCode = 0
	IdErrTypeMismatch IdRootErrorCode = 1
	IdErrNoType       IdRootErrorCode = 2
)

type IdErr = IdRoot[IdErrorIdDef]

func NewIdErr(code IdRootErrorCode) IdErr {
	return NewIdRootError[IdErrorIdDef](code)
}

var NoIdErr = NewIdErr(IdNoError)
