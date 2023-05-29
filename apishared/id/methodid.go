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
type defMethod struct{}

func (f defMethod) ShortString() string { return "method" }
func (f defMethod) Letter() byte        { return 0x6d } 
func (f defMethod) IsError() bool       { return true }

type MethodId IdRoot[defMethod]


func NewMethodIdFromRaw(in IdRaw) MethodId {
	return MethodId(NewIdRootFromRaw[defMethod](in))
}

func (f MethodId) Marshal() *protosupportmsg.IdRaw {
	return MarshalProtobuf(IdRoot[defMethod](f))
}
func ZeroValueMethodId() MethodId {
	return MethodId(ZeroValue[defMethod]())
}
func (f MethodId) Raw() IdRaw {
	return IdRoot[defMethod](f).Raw()
}
func (f MethodId) Equal(other MethodId) bool{
	return IdRoot[defMethod](f).Equal(IdRoot[defMethod](other))
}
func (f MethodId) String() string{
	return IdRoot[defMethod](f).String()
}
func (f MethodId) Short() string{
	return IdRoot[defMethod](f).Short()
}
func (f MethodId) IsError() bool{
	return IdRoot[defMethod](f).IsError()
}
func (f MethodId) IsZeroValue() bool{
	return IdRoot[defMethod](f).IsZeroValue()
}

func UnmarshalMethodId(b *protosupportmsg.IdRaw) (MethodId, IdErr) {
	fid, err := UnmarshalProtobuf[defMethod](b)
	if err.IsError() {
		return ZeroValueMethodId(), err
	}
	return MethodId(fid), NoIdErr
}

//
// End Boilerplate for Method
//
