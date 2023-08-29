package methodcall

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for MethodcallId
//

import (
	
    "github.com/iansmith/parigot/api/shared/id"

	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Methodcall
//

type DefMethodcall struct{}

func (f DefMethodcall) ShortString() string { return "methcall" }
func (f DefMethodcall) Letter() byte        { return 0x6d } 

type MethodcallId id.IdRoot[DefMethodcall]


func NewMethodcallId() MethodcallId {
	return MethodcallId(id.NewIdRoot[DefMethodcall]())
}

func (f MethodcallId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func MethodcallIdZeroValue() MethodcallId {
	return MethodcallId(id.NewIdTyped[DefMethodcall](0xffffffffffffff,0xffffffffffffffff))
}
func MethodcallIdEmptyValue() MethodcallId {
	return MethodcallId(id.NewIdTyped[DefMethodcall](0,0))
}

func (f MethodcallId) Equal(other MethodcallId) bool{
	return id.IdRoot[DefMethodcall](f).Equal(id.IdRoot[DefMethodcall](other))
}
func (f MethodcallId) String() string{
	return id.IdRoot[DefMethodcall](f).String()
}
func (f MethodcallId) Short() string{
	return id.IdRoot[DefMethodcall](f).Short()
}

func (f MethodcallId) IsZeroValue() bool{
	return id.IdRoot[DefMethodcall](f).IsZeroValue()
}
func (f MethodcallId) IsEmptyValue() bool{
	return id.IdRoot[DefMethodcall](f).IsEmptyValue()
}
func (f MethodcallId) IsZeroOrEmptyValue() bool{
	return id.IdRoot[DefMethodcall](f).IsZeroOrEmptyValue()
}

func (f MethodcallId) High() uint64{
	return id.IdRoot[DefMethodcall](f).High()
}
func (f MethodcallId) Low() uint64{
	return id.IdRoot[DefMethodcall](f).Low()
}

func UnmarshalMethodcallId(b *protosupport.IdRaw) MethodcallId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return MethodcallId(id.NewIdTyped[DefMethodcall](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func MethodcallIdFromPair(high, low uint64) MethodcallId {
	return MethodcallId(id.NewIdTyped[DefMethodcall](high,low))
}

//
// End Boilerplate for Methodcall
//
