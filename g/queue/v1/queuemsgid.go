package queue

//
// DO NOT EDIT.  This file was machine genarted by boilerplateid for QueueMsgId
//

import (
	
    "github.com/iansmith/parigot/api/shared/id"

	protosupport "github.com/iansmith/parigot/g/protosupport/v1"
)
//
//  Begin Boilerplate for QueueMsg
//

type DefQueueMsg struct{}

func (f DefQueueMsg) ShortString() string { return "msg" }
func (f DefQueueMsg) Letter() byte        { return 0x6d } 

type QueueMsgId id.IdRoot[DefQueueMsg]


func NewQueueMsgId() QueueMsgId {
	return QueueMsgId(id.NewIdRoot[DefQueueMsg]())
}

func (f QueueMsgId) Marshal() *protosupport.IdRaw {
	raw:=&protosupport.IdRaw{}
	raw.High = f.High()
	raw.Low = f.Low()
	return raw
}

func QueueMsgIdZeroValue() QueueMsgId {
	return QueueMsgId(id.NewIdTyped[DefQueueMsg](0xffffffffffffff,0xffffffffffffffff))
}
func QueueMsgIdEmptyValue() QueueMsgId {
	return QueueMsgId(id.NewIdTyped[DefQueueMsg](0,0))
}

func (f QueueMsgId) Equal(other QueueMsgId) bool{
	return id.IdRoot[DefQueueMsg](f).Equal(id.IdRoot[DefQueueMsg](other))
}
func (f QueueMsgId) String() string{
	return id.IdRoot[DefQueueMsg](f).String()
}
func (f QueueMsgId) Short() string{
	return id.IdRoot[DefQueueMsg](f).Short()
}

func (f QueueMsgId) IsZeroValue() bool{
	return id.IdRoot[DefQueueMsg](f).IsZeroValue()
}
func (f QueueMsgId) IsEmptyValue() bool{
	return id.IdRoot[DefQueueMsg](f).IsEmptyValue()
}
func (f QueueMsgId) IsZeroOrEmptyValue() bool{
	return id.IdRoot[DefQueueMsg](f).IsZeroOrEmptyValue()
}

func (f QueueMsgId) High() uint64{
	return id.IdRoot[DefQueueMsg](f).High()
}
func (f QueueMsgId) Low() uint64{
	return id.IdRoot[DefQueueMsg](f).Low()
}

func UnmarshalQueueMsgId(b *protosupport.IdRaw) QueueMsgId {
	l:=b.GetLow()
	h:=b.GetHigh()
	return QueueMsgId(id.NewIdTyped[DefQueueMsg](h,l))
}

// FromPair is probably not something you want to use unless you
// are pulling values from external storage or files.  If you pulling
// values from the network, use the Marshal() ad Unmarshal()
// functions to work with Ids.  Absolutely no checking is done
// on the values provided, so much caution is advised.
func QueueMsgIdFromPair(high, low uint64) QueueMsgId {
	return QueueMsgId(id.NewIdTyped[DefQueueMsg](high,low))
}

//
// End Boilerplate for QueueMsg
//
