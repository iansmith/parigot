package test

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for TestId
//

import (
	
    "github.com/iansmith/parigot/api/shared/id"

	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Test
//

type DefTest struct{}

func (f DefTest) ShortString() string { return "test" }
func (f DefTest) Letter() byte        { return 0x74 } 

type TestId id.IdRoot[DefTest]


func NewTestId() TestId {
	return TestId(id.NewIdRoot[DefTest]())
}

func (f TestId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func TestIdZeroValue() TestId {
	return TestId(id.NewIdTyped[DefTest](0xffffffffffffff,0xffffffffffffffff))
}
func TestIdEmptyValue() TestId {
	return TestId(id.NewIdTyped[DefTest](0,0))
}

func (f TestId) Equal(other TestId) bool{
	return id.IdRoot[DefTest](f).Equal(id.IdRoot[DefTest](other))
}
func (f TestId) String() string{
	return id.IdRoot[DefTest](f).String()
}
func (f TestId) Short() string{
	return id.IdRoot[DefTest](f).Short()
}

func (f TestId) IsZeroValue() bool{
	return id.IdRoot[DefTest](f).IsZeroValue()
}
func (f TestId) IsEmptyValue() bool{
	return id.IdRoot[DefTest](f).IsEmptyValue()
}
func (f TestId) IsZeroOrEmptyValue() bool{
	return id.IdRoot[DefTest](f).IsZeroOrEmptyValue()
}

func (f TestId) High() uint64{
	return id.IdRoot[DefTest](f).High()
}
func (f TestId) Low() uint64{
	return id.IdRoot[DefTest](f).Low()
}

func UnmarshalTestId(b *protosupport.IdRaw) TestId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return TestId(id.NewIdTyped[DefTest](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func TestIdFromPair(high, low uint64) TestId {
	return TestId(id.NewIdTyped[DefTest](high,low))
}

//
// End Boilerplate for Test
//
