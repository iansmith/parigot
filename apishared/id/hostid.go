package id

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for HostId
//

import (


	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Host
//

type DefHost struct{}

func (f DefHost) ShortString() string { return "host" }
func (f DefHost) Letter() byte        { return 0x68 } 

type HostId IdRoot[DefHost]


func NewHostId() HostId {
	return HostId(NewIdRoot[DefHost]())
}

func (f HostId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}
func HostIdZeroValue() HostId {
	return HostId(NewIdTyped[DefHost](0xffffffffffffff,0xffffffffffffffff))
}
func HostIdEmptyValue() HostId {
	return HostId(NewIdTyped[DefHost](0,0))
}

func (f HostId) Equal(other HostId) bool{
	return IdRoot[DefHost](f).Equal(IdRoot[DefHost](other))
}
func (f HostId) String() string{
	return IdRoot[DefHost](f).String()
}
func (f HostId) Short() string{
	return IdRoot[DefHost](f).Short()
}

func (f HostId) IsZeroValue() bool{
	return IdRoot[DefHost](f).IsZeroValue()
}
func (f HostId) IsEmptyValue() bool{
	return IdRoot[DefHost](f).IsEmptyValue()
}
func (f HostId) IsZeroOrEmptyValue() bool{
	return IdRoot[DefHost](f).IsZeroOrEmptyValue()
}

func (f HostId) High() uint64{
	return IdRoot[DefHost](f).High()
}
func (f HostId) Low() uint64{
	return IdRoot[DefHost](f).Low()
}

func UnmarshalHostId(b *protosupport.IdRaw) HostId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return HostId(NewIdTyped[DefHost](h,l))
}

//
// End Boilerplate for Host
//
