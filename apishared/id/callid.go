package id

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for CallId
//

import (


	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
)
//
//  Begin Boilerplate for Call
//

type CallCode IdRootErrorCode
const CallNoError CallCode = 0
type defCall struct{}

func (f defCall) ShortString() string { return "call" }
func (f defCall) Letter() byte        { return 0x63 } 
func (f defCall) IsError() bool       { return true }

type CallId IdRoot[defCall]


func NewCallIdFromRaw(in IdRaw) CallId {
	return CallId(NewIdRootFromRaw[defCall](in))
}

func (f CallId) Marshal() *protosupportmsg.IdRaw {
	return MarshalProtobuf(IdRoot[defCall](f))
}
func ZeroValueCallId() CallId {
	return CallId(ZeroValue[defCall]())
}
func (f CallId) Raw() IdRaw {
	return IdRoot[defCall](f).Raw()
}
func (f CallId) Equal(other CallId) bool{
	return IdRoot[defCall](f).Equal(IdRoot[defCall](other))
}
func (f CallId) String() string{
	return IdRoot[defCall](f).String()
}
func (f CallId) Short() string{
	return IdRoot[defCall](f).Short()
}
func (f CallId) IsError() bool{
	return IdRoot[defCall](f).IsError()
}
func (f CallId) IsZeroValue() bool{
	return IdRoot[defCall](f).IsZeroValue()
}

func UnmarshalCallId(b *protosupportmsg.IdRaw) (CallId, IdErr) {
	fid, err := UnmarshalProtobuf[defCall](b)
	if err.IsError() {
		return ZeroValueCallId(), err
	}
	return CallId(fid), NoIdErr
}

//
// End Boilerplate for Call
//
