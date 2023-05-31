package id

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for MethodId
//

import (


	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
)
//
//  Begin Boilerplate for Method
//

type MethodCode IdRootErrorCode
const MethodNoError MethodCode = 0
type DefMethod struct{}

func (f DefMethod) ShortString() string { return "method" }
func (f DefMethod) Letter() byte        { return 0x6d } 
func (f DefMethod) IsError() bool       { return true }

type MethodId IdRoot[DefMethod]


func NewMethodIdFromRaw(in IdRaw) MethodId {
	return MethodId(NewIdRootFromRaw[DefMethod](in))
}

func (f MethodId) Marshal() *protosupportmsg.IdRaw {
	return MarshalProtobuf(IdRoot[DefMethod](f))
}
func ZeroValueMethodId() MethodId {
	return MethodId(ZeroValue[DefMethod]())
}
func (f MethodId) Raw() IdRaw {
	return IdRoot[DefMethod](f).Raw()
}
func (f MethodId) Equal(other MethodId) bool{
	return IdRoot[DefMethod](f).Equal(IdRoot[DefMethod](other))
}
func (f MethodId) String() string{
	return IdRoot[DefMethod](f).String()
}
func (f MethodId) Short() string{
	return IdRoot[DefMethod](f).Short()
}
func (f MethodId) IsError() bool{
	return IdRoot[DefMethod](f).IsError()
}
func (f MethodId) IsZeroValue() bool{
	return IdRoot[DefMethod](f).IsZeroValue()
}

func UnmarshalMethodId(b *protosupportmsg.IdRaw) (MethodId, IdErr) {
	fid, err := UnmarshalProtobuf[DefMethod](b)
	if err.IsError() {
		return ZeroValueMethodId(), err
	}
	return MethodId(fid), NoIdErr
}

func NewMethodId() MethodId {
	idroot := NewIdRoot[DefMethod]()
	return MethodId(idroot)
}

//
// End Boilerplate for Method
//
