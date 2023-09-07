package queue

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for QueueId
//

import (
	
    "github.com/iansmith/parigot/api/shared/id"

	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for Queue
//

type DefQueue struct{}

func (f DefQueue) ShortString() string { return "queue" }
func (f DefQueue) Letter() byte        { return 0x71 } 

type QueueId id.IdRoot[DefQueue]


func NewQueueId() QueueId {
	return QueueId(id.NewIdRoot[DefQueue]())
}

func (f QueueId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func QueueIdZeroValue() QueueId {
	return QueueId(id.NewIdTyped[DefQueue](0xffffffffffffff,0xffffffffffffffff))
}
func QueueIdEmptyValue() QueueId {
	return QueueId(id.NewIdTyped[DefQueue](0,0))
}

func (f QueueId) Equal(other QueueId) bool{
	return id.IdRoot[DefQueue](f).Equal(id.IdRoot[DefQueue](other))
}
func (f QueueId) String() string{
	return id.IdRoot[DefQueue](f).String()
}
func (f QueueId) Short() string{
	return id.IdRoot[DefQueue](f).Short()
}

func (f QueueId) IsZeroValue() bool{
	return id.IdRoot[DefQueue](f).IsZeroValue()
}
func (f QueueId) IsEmptyValue() bool{
	return id.IdRoot[DefQueue](f).IsEmptyValue()
}
func (f QueueId) IsZeroOrEmptyValue() bool{
	return id.IdRoot[DefQueue](f).IsZeroOrEmptyValue()
}

func (f QueueId) High() uint64{
	return id.IdRoot[DefQueue](f).High()
}
func (f QueueId) Low() uint64{
	return id.IdRoot[DefQueue](f).Low()
}

func UnmarshalQueueId(b *protosupport.IdRaw) QueueId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return QueueId(id.NewIdTyped[DefQueue](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func QueueIdFromPair(high, low uint64) QueueId {
	return QueueId(id.NewIdTyped[DefQueue](high,low))
}

//
// End Boilerplate for Queue
//
