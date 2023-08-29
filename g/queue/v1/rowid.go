package queue

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for RowId
//

import (
	
    "github.com/iansmith/parigot/api/shared/id"

	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Row
//

type DefRow struct{}

func (f DefRow) ShortString() string { return "row" }
func (f DefRow) Letter() byte        { return 0x72 } 

type RowId id.IdRoot[DefRow]


func NewRowId() RowId {
	return RowId(id.NewIdRoot[DefRow]())
}

func (f RowId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func RowIdZeroValue() RowId {
	return RowId(id.NewIdTyped[DefRow](0xffffffffffffff,0xffffffffffffffff))
}
func RowIdEmptyValue() RowId {
	return RowId(id.NewIdTyped[DefRow](0,0))
}

func (f RowId) Equal(other RowId) bool{
	return id.IdRoot[DefRow](f).Equal(id.IdRoot[DefRow](other))
}
func (f RowId) String() string{
	return id.IdRoot[DefRow](f).String()
}
func (f RowId) Short() string{
	return id.IdRoot[DefRow](f).Short()
}

func (f RowId) IsZeroValue() bool{
	return id.IdRoot[DefRow](f).IsZeroValue()
}
func (f RowId) IsEmptyValue() bool{
	return id.IdRoot[DefRow](f).IsEmptyValue()
}
func (f RowId) IsZeroOrEmptyValue() bool{
	return id.IdRoot[DefRow](f).IsZeroOrEmptyValue()
}

func (f RowId) High() uint64{
	return id.IdRoot[DefRow](f).High()
}
func (f RowId) Low() uint64{
	return id.IdRoot[DefRow](f).Low()
}

func UnmarshalRowId(b *protosupport.IdRaw) RowId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return RowId(id.NewIdTyped[DefRow](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func RowIdFromPair(high, low uint64) RowId {
	return RowId(id.NewIdTyped[DefRow](high,low))
}

//
// End Boilerplate for Row
//
