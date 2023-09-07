package file

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for FileId
//

import (
	
    "github.com/iansmith/parigot/api/shared/id"

	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for File
//

type DefFile struct{}

func (f DefFile) ShortString() string { return "file" }
func (f DefFile) Letter() byte        { return 0x66 } 

type FileId id.IdRoot[DefFile]


func NewFileId() FileId {
	return FileId(id.NewIdRoot[DefFile]())
}

func (f FileId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func FileIdZeroValue() FileId {
	return FileId(id.NewIdTyped[DefFile](0xffffffffffffff,0xffffffffffffffff))
}
func FileIdEmptyValue() FileId {
	return FileId(id.NewIdTyped[DefFile](0,0))
}

func (f FileId) Equal(other FileId) bool{
	return id.IdRoot[DefFile](f).Equal(id.IdRoot[DefFile](other))
}
func (f FileId) String() string{
	return id.IdRoot[DefFile](f).String()
}
func (f FileId) Short() string{
	return id.IdRoot[DefFile](f).Short()
}

func (f FileId) IsZeroValue() bool{
	return id.IdRoot[DefFile](f).IsZeroValue()
}
func (f FileId) IsEmptyValue() bool{
	return id.IdRoot[DefFile](f).IsEmptyValue()
}
func (f FileId) IsZeroOrEmptyValue() bool{
	return id.IdRoot[DefFile](f).IsZeroOrEmptyValue()
}

func (f FileId) High() uint64{
	return id.IdRoot[DefFile](f).High()
}
func (f FileId) Low() uint64{
	return id.IdRoot[DefFile](f).Low()
}

func UnmarshalFileId(b *protosupport.IdRaw) FileId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return FileId(id.NewIdTyped[DefFile](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func FileIdFromPair(high, low uint64) FileId {
	return FileId(id.NewIdTyped[DefFile](high,low))
}

//
// End Boilerplate for File
//
