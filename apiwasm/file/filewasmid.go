package main

import (
	"github.com/iansmith/parigot/apishared/id"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
)

//
//  Begin Boilerplate for declaring our two types of Ids
//

type FileErrCode id.IdRootErrorCode
type FileIdDef struct{}

func (f FileIdDef) ShortString() string { return "file" }
func (f FileIdDef) Letter() byte        { return 0x66 } //f
func (f FileIdDef) IsError() bool       { return false }

type FileId id.IdRoot[FileIdDef]

type FileErrIdDef struct{}

func (f FileErrIdDef) ShortString() string { return "fileerr" }
func (f FileErrIdDef) Letter() byte        { return 0x46 } //F
func (f FileErrIdDef) IsError() bool       { return true }

type FileErrId id.IdRoot[FileErrIdDef]

// NoErr is a convenience for referring to an errorId that is
// "everything is ok".
var NoErr = NewFileErrId(NoError)

func NewFileErrId(code FileErrCode) FileErrId {
	e := id.NewIdRootError[FileErrIdDef](FileErrIdDef{}, id.IdRootErrorCode(NoError))
	return FileErrId(e)
}

func (f FileId) Marshal() *protosupportmsg.BaseId {
	return id.MarshalProtobuf[FileIdDef](id.IdRoot[FileIdDef](f))
}
func (f FileErrId) Marshal() *protosupportmsg.BaseId {
	return id.MarshalProtobuf[FileErrIdDef](id.IdRoot[FileErrIdDef](f))
}
func ZeroValueFileId() FileId {
	return FileId(id.ZeroValue[FileIdDef](FileIdDef{}))
}
func ZeroValueFileErrId() FileErrId {
	return FileErrId(id.ZeroValue[FileErrIdDef](FileErrIdDef{}))
}

func UnmarshalFileId(b *protosupportmsg.BaseId) (FileId, id.IdErr) {
	fid, err := id.UnmarshalProtobuf[FileIdDef](FileIdDef{}, b)
	if err != id.NoIdErr {
		return ZeroValueFileId(), err
	}
	return FileId(fid), id.NoIdErr
}

func NewFileId() FileId {
	idroot := id.NewIdRoot[FileIdDef](FileIdDef{})
	return FileId(idroot)
}

//
// End Boilerplate for declaring our two types of Ids
//
