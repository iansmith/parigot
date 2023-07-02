package id

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for ServiceId
//

import (


	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Service
//

type DefService struct{}

func (f DefService) ShortString() string { return "svc" }
func (f DefService) Letter() byte        { return 0x73 } 

type ServiceId IdRoot[DefService]


func NewServiceId() ServiceId {
	return ServiceId(NewIdRoot[DefService]())
}

func (f ServiceId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func ServiceIdZeroValue() ServiceId {
	return ServiceId(NewIdTyped[DefService](0xffffffffffffff,0xffffffffffffffff))
}
func ServiceIdEmptyValue() ServiceId {
	return ServiceId(NewIdTyped[DefService](0,0))
}

func (f ServiceId) Equal(other ServiceId) bool{
	return IdRoot[DefService](f).Equal(IdRoot[DefService](other))
}
func (f ServiceId) String() string{
	return IdRoot[DefService](f).String()
}
func (f ServiceId) Short() string{
	return IdRoot[DefService](f).Short()
}

func (f ServiceId) IsZeroValue() bool{
	return IdRoot[DefService](f).IsZeroValue()
}
func (f ServiceId) IsEmptyValue() bool{
	return IdRoot[DefService](f).IsEmptyValue()
}
func (f ServiceId) IsZeroOrEmptyValue() bool{
	return IdRoot[DefService](f).IsZeroOrEmptyValue()
}

func (f ServiceId) High() uint64{
	return IdRoot[DefService](f).High()
}
func (f ServiceId) Low() uint64{
	return IdRoot[DefService](f).Low()
}

func UnmarshalServiceId(b *protosupport.IdRaw) ServiceId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return ServiceId(NewIdTyped[DefService](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func ServiceIdFromPair(high, low uint64) ServiceId {
	return ServiceId(NewIdTyped[DefService](high,low))
}

//
// End Boilerplate for Service
//
