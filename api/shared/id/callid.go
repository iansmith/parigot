package id

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for CallId
//

import (


	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Call
//

type DefCall struct{}

func (f DefCall) ShortString() string { return "call" }
func (f DefCall) Letter() byte        { return 0x63 } 

type CallId IdRoot[DefCall]


func NewCallId() CallId {
	return CallId(NewIdRoot[DefCall]())
}

func (f CallId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func CallIdZeroValue() CallId {
	return CallId(NewIdTyped[DefCall](0xffffffffffffff,0xffffffffffffffff))
}
func CallIdEmptyValue() CallId {
	return CallId(NewIdTyped[DefCall](0,0))
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

func (f CallId) IsZeroValue() bool{
	return IdRoot[DefCall](f).IsZeroValue()
}
func (f CallId) IsEmptyValue() bool{
	return IdRoot[DefCall](f).IsEmptyValue()
}
func (f CallId) IsZeroOrEmptyValue() bool{
	return IdRoot[DefCall](f).IsZeroOrEmptyValue()
}

func (f CallId) High() uint64{
	return IdRoot[DefCall](f).High()
}
func (f CallId) Low() uint64{
	return IdRoot[DefCall](f).Low()
}

func UnmarshalCallId(b *protosupport.IdRaw) CallId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return CallId(NewIdTyped[DefCall](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func CallIdFromPair(high, low uint64) CallId {
	return CallId(NewIdTyped[DefCall](high,low))
}

//
// End Boilerplate for Call
//
