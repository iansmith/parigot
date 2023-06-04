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
type DefCall struct{}

func (f DefCall) ShortString() string { return "call" }
func (f DefCall) Letter() byte        { return 0x63 } 
func (f DefCall) IsError() bool       { return true }

type CallId IdRoot[DefCall]


func NewCallIdFromRaw(in IdRaw) CallId {
	return CallId(NewIdRootFromRaw[DefCall](in))
}

func (f CallId) Marshal() *protosupportmsg.IdRaw {
	return MarshalProtobuf(IdRoot[DefCall](f))
}
func ZeroValueCallId() CallId {
	return CallId(ZeroValue[DefCall]())
}
func (f CallId) Raw() IdRaw {
	return IdRoot[DefCall](f).Raw()
}
func (f CallId) Equal(other CallId) bool{
	return IdRoot[DefCall](f).Equal(IdRoot[DefCall](other))
}
func (f CallId) String() string{
	return IdRoot[DefCall](f).String()
}
func (f CallId) Short() string{
	return IdRoot[DefCall](f).Short()
}
func (f CallId) IsError() bool{
	return IdRoot[DefCall](f).IsError()
}
func (f CallId) IsZeroValue() bool{
	return IdRoot[DefCall](f).IsZeroValue()
}
func (f CallId) IsEmptyValue() bool{
	return IdRoot[DefCall](f).IsEmptyValue()
}

func (f CallId) High() uint64{
	return IdRoot[DefCall](f).High()
}
func (f CallId) Low() uint64{
	return IdRoot[DefCall](f).Low()
}

func UnmarshalCallId(b *protosupportmsg.IdRaw) (CallId, IdErr) {
	fid, err := UnmarshalProtobuf[DefCall](b)
	if err.IsError() {
		return ZeroValueCallId(), err
	}
	return CallId(fid), NoIdErr
}

func MustUnmarshalCallId(b *protosupportmsg.IdRaw) CallId{
	result, err:=UnmarshalCallId(b)
	if err.IsError() {
		panic("unable to unmarshal CallId from raw value: "+err.String())
	}
	return result
}

func NewCallIdFromProto(in *protosupportmsg.IdRaw) CallId {
	raw:=MustUnmarshalCallId(in)
	return CallId(raw)
}

func NewCallId() CallId {
	idroot := NewIdRoot[DefCall]()
	return CallId(idroot)
}


//
// End Boilerplate for Call
//
