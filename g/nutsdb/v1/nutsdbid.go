package nutsdb

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for NutsDBId
//

import (
	
    "github.com/iansmith/parigot/api/shared/id"

	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for NutsDB
//

type DefNutsDB struct{}

func (f DefNutsDB) ShortString() string { return "nutsdb" }
func (f DefNutsDB) Letter() byte        { return 0x6e } 

type NutsDBId id.IdRoot[DefNutsDB]


func NewNutsDBId() NutsDBId {
	return NutsDBId(id.NewIdRoot[DefNutsDB]())
}

func (f NutsDBId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func NutsDBIdZeroValue() NutsDBId {
	return NutsDBId(id.NewIdTyped[DefNutsDB](0xffffffffffffff,0xffffffffffffffff))
}
func NutsDBIdEmptyValue() NutsDBId {
	return NutsDBId(id.NewIdTyped[DefNutsDB](0,0))
}

func (f NutsDBId) Equal(other NutsDBId) bool{
	return id.IdRoot[DefNutsDB](f).Equal(id.IdRoot[DefNutsDB](other))
}
func (f NutsDBId) String() string{
	return id.IdRoot[DefNutsDB](f).String()
}
func (f NutsDBId) Short() string{
	return id.IdRoot[DefNutsDB](f).Short()
}

func (f NutsDBId) IsZeroValue() bool{
	return id.IdRoot[DefNutsDB](f).IsZeroValue()
}
func (f NutsDBId) IsEmptyValue() bool{
	return id.IdRoot[DefNutsDB](f).IsEmptyValue()
}
func (f NutsDBId) IsZeroOrEmptyValue() bool{
	return id.IdRoot[DefNutsDB](f).IsZeroOrEmptyValue()
}

func (f NutsDBId) High() uint64{
	return id.IdRoot[DefNutsDB](f).High()
}
func (f NutsDBId) Low() uint64{
	return id.IdRoot[DefNutsDB](f).Low()
}

func UnmarshalNutsDBId(b *protosupport.IdRaw) NutsDBId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return NutsDBId(id.NewIdTyped[DefNutsDB](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func NutsDBIdFromPair(high, low uint64) NutsDBId {
	return NutsDBId(id.NewIdTyped[DefNutsDB](high,low))
}

//
// End Boilerplate for NutsDB
//
