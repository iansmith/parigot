package id

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for MethodId
//

import (


	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Method
//

type DefMethod struct{}

func (f DefMethod) ShortString() string { return "method" }
func (f DefMethod) Letter() byte        { return 0x6d } 

type MethodId IdRoot[DefMethod]


func NewMethodId() MethodId {
	return MethodId(NewIdRoot[DefMethod]())
}

func (f MethodId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func MethodIdZeroValue() MethodId {
	return MethodId(NewIdTyped[DefMethod](0xffffffffffffff,0xffffffffffffffff))
}
func MethodIdEmptyValue() MethodId {
	return MethodId(NewIdTyped[DefMethod](0,0))
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

func (f MethodId) IsZeroValue() bool{
	return IdRoot[DefMethod](f).IsZeroValue()
}
func (f MethodId) IsEmptyValue() bool{
	return IdRoot[DefMethod](f).IsEmptyValue()
}
func (f MethodId) IsZeroOrEmptyValue() bool{
	return IdRoot[DefMethod](f).IsZeroOrEmptyValue()
}

func (f MethodId) High() uint64{
	return IdRoot[DefMethod](f).High()
}
func (f MethodId) Low() uint64{
	return IdRoot[DefMethod](f).Low()
}

func UnmarshalMethodId(b *protosupport.IdRaw) MethodId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return MethodId(NewIdTyped[DefMethod](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func MethodIdFromPair(high, low uint64) MethodId {
	return MethodId(NewIdTyped[DefMethod](high,low))
}

//
// End Boilerplate for Method
//
