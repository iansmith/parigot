// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: msg/store/v1/store.proto

package storemsg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MediaType int32

const (
	MediaType_MEDIA_TYPE_UNSPECIFIED        MediaType = 0
	MediaType_MEDIA_TYPE_VHS                MediaType = 1
	MediaType_MEDIA_TYPE_BETA               MediaType = 2
	MediaType_MEDIA_TYPE_LASERDISC          MediaType = 4
	MediaType_MEDIA_TYPE_DVD                MediaType = 5
	MediaType_MEDIA_TYPE_CD                 MediaType = 6
	MediaType_MEDIA_TYPE_CD_SINGLE          MediaType = 7
	MediaType_MEDIA_TYPE_ATARI_CART         MediaType = 8
	MediaType_MEDIA_TYPE_INTELLIVISION_CART MediaType = 9
	MediaType_MEDIA_TYPE_CASSETTE           MediaType = 10
	MediaType_MEDIA_TYPE_8TRACK             MediaType = 11
	MediaType_MEDIA_TYPE_VINYL              MediaType = 12
)

// Enum value maps for MediaType.
var (
	MediaType_name = map[int32]string{
		0:  "MEDIA_TYPE_UNSPECIFIED",
		1:  "MEDIA_TYPE_VHS",
		2:  "MEDIA_TYPE_BETA",
		4:  "MEDIA_TYPE_LASERDISC",
		5:  "MEDIA_TYPE_DVD",
		6:  "MEDIA_TYPE_CD",
		7:  "MEDIA_TYPE_CD_SINGLE",
		8:  "MEDIA_TYPE_ATARI_CART",
		9:  "MEDIA_TYPE_INTELLIVISION_CART",
		10: "MEDIA_TYPE_CASSETTE",
		11: "MEDIA_TYPE_8TRACK",
		12: "MEDIA_TYPE_VINYL",
	}
	MediaType_value = map[string]int32{
		"MEDIA_TYPE_UNSPECIFIED":        0,
		"MEDIA_TYPE_VHS":                1,
		"MEDIA_TYPE_BETA":               2,
		"MEDIA_TYPE_LASERDISC":          4,
		"MEDIA_TYPE_DVD":                5,
		"MEDIA_TYPE_CD":                 6,
		"MEDIA_TYPE_CD_SINGLE":          7,
		"MEDIA_TYPE_ATARI_CART":         8,
		"MEDIA_TYPE_INTELLIVISION_CART": 9,
		"MEDIA_TYPE_CASSETTE":           10,
		"MEDIA_TYPE_8TRACK":             11,
		"MEDIA_TYPE_VINYL":              12,
	}
)

func (x MediaType) Enum() *MediaType {
	p := new(MediaType)
	*p = x
	return p
}

func (x MediaType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MediaType) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_store_v1_store_proto_enumTypes[0].Descriptor()
}

func (MediaType) Type() protoreflect.EnumType {
	return &file_msg_store_v1_store_proto_enumTypes[0]
}

func (x MediaType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MediaType.Descriptor instead.
func (MediaType) EnumDescriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{0}
}

type ContentType int32

const (
	ContentType_CONTENT_TYPE_UNSPECIFIED ContentType = 0
	ContentType_CONTENT_TYPE_MUSIC       ContentType = 1
	ContentType_CONTENT_TYPE_TV          ContentType = 2
	ContentType_CONTENT_TYPE_MOVIE       ContentType = 3
)

// Enum value maps for ContentType.
var (
	ContentType_name = map[int32]string{
		0: "CONTENT_TYPE_UNSPECIFIED",
		1: "CONTENT_TYPE_MUSIC",
		2: "CONTENT_TYPE_TV",
		3: "CONTENT_TYPE_MOVIE",
	}
	ContentType_value = map[string]int32{
		"CONTENT_TYPE_UNSPECIFIED": 0,
		"CONTENT_TYPE_MUSIC":       1,
		"CONTENT_TYPE_TV":          2,
		"CONTENT_TYPE_MOVIE":       3,
	}
)

func (x ContentType) Enum() *ContentType {
	p := new(ContentType)
	*p = x
	return p
}

func (x ContentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContentType) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_store_v1_store_proto_enumTypes[1].Descriptor()
}

func (ContentType) Type() protoreflect.EnumType {
	return &file_msg_store_v1_store_proto_enumTypes[1]
}

func (x ContentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContentType.Descriptor instead.
func (ContentType) EnumDescriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{1}
}

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{0}
}

func (x *Item) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Amount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Units      int32 `protobuf:"varint,1,opt,name=units,proto3" json:"units,omitempty"`
	Hundredths int32 `protobuf:"varint,2,opt,name=hundredths,proto3" json:"hundredths,omitempty"`
}

func (x *Amount) Reset() {
	*x = Amount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Amount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Amount) ProtoMessage() {}

func (x *Amount) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Amount.ProtoReflect.Descriptor instead.
func (*Amount) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{1}
}

func (x *Amount) GetUnits() int32 {
	if x != nil {
		return x.Units
	}
	return 0
}

func (x *Amount) GetHundredths() int32 {
	if x != nil {
		return x.Hundredths
	}
	return 0
}

type Boat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Creator string      `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Title   string      `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Year    int32       `protobuf:"varint,3,opt,name=year,proto3" json:"year,omitempty"`
	Media   MediaType   `protobuf:"varint,4,opt,name=media,proto3,enum=msg.store.v1.MediaType" json:"media,omitempty"`
	Price   *Amount     `protobuf:"bytes,5,opt,name=price,proto3" json:"price,omitempty"`
	Content ContentType `protobuf:"varint,6,opt,name=content,proto3,enum=msg.store.v1.ContentType" json:"content,omitempty"`
}

func (x *Boat) Reset() {
	*x = Boat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Boat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Boat) ProtoMessage() {}

func (x *Boat) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Boat.ProtoReflect.Descriptor instead.
func (*Boat) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{2}
}

func (x *Boat) GetCreator() string {
	if x != nil {
		return x.Creator
	}
	return ""
}

func (x *Boat) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Boat) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *Boat) GetMedia() MediaType {
	if x != nil {
		return x.Media
	}
	return MediaType_MEDIA_TYPE_UNSPECIFIED
}

func (x *Boat) GetPrice() *Amount {
	if x != nil {
		return x.Price
	}
	return nil
}

func (x *Boat) GetContent() ContentType {
	if x != nil {
		return x.Content
	}
	return ContentType_CONTENT_TYPE_UNSPECIFIED
}

type MediaTypesInStockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MediaTypesInStockRequest) Reset() {
	*x = MediaTypesInStockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaTypesInStockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaTypesInStockRequest) ProtoMessage() {}

func (x *MediaTypesInStockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaTypesInStockRequest.ProtoReflect.Descriptor instead.
func (*MediaTypesInStockRequest) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{3}
}

type MediaTypesInStockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InStock []MediaType `protobuf:"varint,1,rep,packed,name=in_stock,json=inStock,proto3,enum=msg.store.v1.MediaType" json:"in_stock,omitempty"`
}

func (x *MediaTypesInStockResponse) Reset() {
	*x = MediaTypesInStockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaTypesInStockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaTypesInStockResponse) ProtoMessage() {}

func (x *MediaTypesInStockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaTypesInStockResponse.ProtoReflect.Descriptor instead.
func (*MediaTypesInStockResponse) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{4}
}

func (x *MediaTypesInStockResponse) GetInStock() []MediaType {
	if x != nil {
		return x.InStock
	}
	return nil
}

type BestOfAllTimeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content ContentType `protobuf:"varint,1,opt,name=content,proto3,enum=msg.store.v1.ContentType" json:"content,omitempty"`
}

func (x *BestOfAllTimeRequest) Reset() {
	*x = BestOfAllTimeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BestOfAllTimeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BestOfAllTimeRequest) ProtoMessage() {}

func (x *BestOfAllTimeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BestOfAllTimeRequest.ProtoReflect.Descriptor instead.
func (*BestOfAllTimeRequest) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{5}
}

func (x *BestOfAllTimeRequest) GetContent() ContentType {
	if x != nil {
		return x.Content
	}
	return ContentType_CONTENT_TYPE_UNSPECIFIED
}

type BestOfAllTimeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Boat *Boat `protobuf:"bytes,1,opt,name=boat,proto3" json:"boat,omitempty"`
}

func (x *BestOfAllTimeResponse) Reset() {
	*x = BestOfAllTimeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BestOfAllTimeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BestOfAllTimeResponse) ProtoMessage() {}

func (x *BestOfAllTimeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BestOfAllTimeResponse.ProtoReflect.Descriptor instead.
func (*BestOfAllTimeResponse) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{6}
}

func (x *BestOfAllTimeResponse) GetBoat() *Boat {
	if x != nil {
		return x.Boat
	}
	return nil
}

type RevenueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Day   int32 `protobuf:"varint,1,opt,name=day,proto3" json:"day,omitempty"`
	Month int32 `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	Year  int32 `protobuf:"varint,3,opt,name=year,proto3" json:"year,omitempty"`
}

func (x *RevenueRequest) Reset() {
	*x = RevenueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RevenueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevenueRequest) ProtoMessage() {}

func (x *RevenueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevenueRequest.ProtoReflect.Descriptor instead.
func (*RevenueRequest) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{7}
}

func (x *RevenueRequest) GetDay() int32 {
	if x != nil {
		return x.Day
	}
	return 0
}

func (x *RevenueRequest) GetMonth() int32 {
	if x != nil {
		return x.Month
	}
	return 0
}

func (x *RevenueRequest) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

type RevenueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Revenue float32 `protobuf:"fixed32,1,opt,name=revenue,proto3" json:"revenue,omitempty"`
}

func (x *RevenueResponse) Reset() {
	*x = RevenueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RevenueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevenueResponse) ProtoMessage() {}

func (x *RevenueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevenueResponse.ProtoReflect.Descriptor instead.
func (*RevenueResponse) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{8}
}

func (x *RevenueResponse) GetRevenue() float32 {
	if x != nil {
		return x.Revenue
	}
	return 0
}

type SoldItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item   *Item                  `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
	Amount *Amount                `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	When   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=when,proto3" json:"when,omitempty"`
}

func (x *SoldItemRequest) Reset() {
	*x = SoldItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SoldItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SoldItemRequest) ProtoMessage() {}

func (x *SoldItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SoldItemRequest.ProtoReflect.Descriptor instead.
func (*SoldItemRequest) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{9}
}

func (x *SoldItemRequest) GetItem() *Item {
	if x != nil {
		return x.Item
	}
	return nil
}

func (x *SoldItemRequest) GetAmount() *Amount {
	if x != nil {
		return x.Amount
	}
	return nil
}

func (x *SoldItemRequest) GetWhen() *timestamppb.Timestamp {
	if x != nil {
		return x.When
	}
	return nil
}

type SoldItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SoldItemResponse) Reset() {
	*x = SoldItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SoldItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SoldItemResponse) ProtoMessage() {}

func (x *SoldItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SoldItemResponse.ProtoReflect.Descriptor instead.
func (*SoldItemResponse) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{10}
}

type GetInStockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item *Item `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
}

func (x *GetInStockRequest) Reset() {
	*x = GetInStockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInStockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInStockRequest) ProtoMessage() {}

func (x *GetInStockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInStockRequest.ProtoReflect.Descriptor instead.
func (*GetInStockRequest) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{11}
}

func (x *GetInStockRequest) GetItem() *Item {
	if x != nil {
		return x.Item
	}
	return nil
}

type GetInStockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *GetInStockResponse) Reset() {
	*x = GetInStockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_store_v1_store_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInStockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInStockResponse) ProtoMessage() {}

func (x *GetInStockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_store_v1_store_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInStockResponse.ProtoReflect.Descriptor instead.
func (*GetInStockResponse) Descriptor() ([]byte, []int) {
	return file_msg_store_v1_store_proto_rawDescGZIP(), []int{12}
}

func (x *GetInStockResponse) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_msg_store_v1_store_proto protoreflect.FileDescriptor

var file_msg_store_v1_store_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6d, 0x73, 0x67, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6d, 0x73, 0x67, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x16, 0x0a, 0x04, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x3e, 0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x75,
	0x6e, 0x69, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x75, 0x6e, 0x69, 0x74,
	0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x68, 0x75, 0x6e, 0x64, 0x72, 0x65, 0x64, 0x74, 0x68, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x68, 0x75, 0x6e, 0x64, 0x72, 0x65, 0x64, 0x74, 0x68,
	0x73, 0x22, 0xda, 0x01, 0x0a, 0x04, 0x42, 0x6f, 0x61, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x65,
	0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x12, 0x2d,
	0x0a, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e,
	0x6d, 0x73, 0x67, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x2a, 0x0a,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d,
	0x73, 0x67, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x6d, 0x73, 0x67,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x1a,
	0x0a, 0x18, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x73, 0x49, 0x6e, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4f, 0x0a, 0x19, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x73, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x69, 0x6e, 0x5f, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x6d, 0x73, 0x67, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x07, 0x69, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x22, 0x4b, 0x0a, 0x14, 0x42,
	0x65, 0x73, 0x74, 0x4f, 0x66, 0x41, 0x6c, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3f, 0x0a, 0x15, 0x42, 0x65, 0x73, 0x74,
	0x4f, 0x66, 0x41, 0x6c, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x26, 0x0a, 0x04, 0x62, 0x6f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42,
	0x6f, 0x61, 0x74, 0x52, 0x04, 0x62, 0x6f, 0x61, 0x74, 0x22, 0x4c, 0x0a, 0x0e, 0x52, 0x65, 0x76,
	0x65, 0x6e, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x64,
	0x61, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x64, 0x61, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6d, 0x6f,
	0x6e, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x22, 0x2b, 0x0a, 0x0f, 0x52, 0x65, 0x76, 0x65, 0x6e,
	0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65,
	0x76, 0x65, 0x6e, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x07, 0x72, 0x65, 0x76,
	0x65, 0x6e, 0x75, 0x65, 0x22, 0x97, 0x01, 0x0a, 0x0f, 0x53, 0x6f, 0x6c, 0x64, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d,
	0x12, 0x2c, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e,
	0x0a, 0x04, 0x77, 0x68, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x77, 0x68, 0x65, 0x6e, 0x22, 0x12,
	0x0a, 0x10, 0x53, 0x6f, 0x6c, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x3b, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x22,
	0x2a, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2a, 0xaf, 0x02, 0x0a, 0x09,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x4d, 0x45, 0x44,
	0x49, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4d, 0x45, 0x44, 0x49, 0x41, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x56, 0x48, 0x53, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x45, 0x44,
	0x49, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x45, 0x54, 0x41, 0x10, 0x02, 0x12, 0x18,
	0x0a, 0x14, 0x4d, 0x45, 0x44, 0x49, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x41, 0x53,
	0x45, 0x52, 0x44, 0x49, 0x53, 0x43, 0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x4d, 0x45, 0x44, 0x49,
	0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x56, 0x44, 0x10, 0x05, 0x12, 0x11, 0x0a, 0x0d,
	0x4d, 0x45, 0x44, 0x49, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x44, 0x10, 0x06, 0x12,
	0x18, 0x0a, 0x14, 0x4d, 0x45, 0x44, 0x49, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x44,
	0x5f, 0x53, 0x49, 0x4e, 0x47, 0x4c, 0x45, 0x10, 0x07, 0x12, 0x19, 0x0a, 0x15, 0x4d, 0x45, 0x44,
	0x49, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x54, 0x41, 0x52, 0x49, 0x5f, 0x43, 0x41,
	0x52, 0x54, 0x10, 0x08, 0x12, 0x21, 0x0a, 0x1d, 0x4d, 0x45, 0x44, 0x49, 0x41, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x4c, 0x4c, 0x49, 0x56, 0x49, 0x53, 0x49, 0x4f, 0x4e,
	0x5f, 0x43, 0x41, 0x52, 0x54, 0x10, 0x09, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x45, 0x44, 0x49, 0x41,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x41, 0x53, 0x53, 0x45, 0x54, 0x54, 0x45, 0x10, 0x0a,
	0x12, 0x15, 0x0a, 0x11, 0x4d, 0x45, 0x44, 0x49, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x38,
	0x54, 0x52, 0x41, 0x43, 0x4b, 0x10, 0x0b, 0x12, 0x14, 0x0a, 0x10, 0x4d, 0x45, 0x44, 0x49, 0x41,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x56, 0x49, 0x4e, 0x59, 0x4c, 0x10, 0x0c, 0x2a, 0x70, 0x0a,
	0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x18,
	0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x4f,
	0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x55, 0x53, 0x49, 0x43,
	0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x54, 0x56, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x4f, 0x4e, 0x54, 0x45,
	0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x4f, 0x56, 0x49, 0x45, 0x10, 0x03, 0x42,
	0x25, 0x5a, 0x23, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x76, 0x76, 0x76, 0x2f, 0x67,
	0x2f, 0x6d, 0x73, 0x67, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_store_v1_store_proto_rawDescOnce sync.Once
	file_msg_store_v1_store_proto_rawDescData = file_msg_store_v1_store_proto_rawDesc
)

func file_msg_store_v1_store_proto_rawDescGZIP() []byte {
	file_msg_store_v1_store_proto_rawDescOnce.Do(func() {
		file_msg_store_v1_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_store_v1_store_proto_rawDescData)
	})
	return file_msg_store_v1_store_proto_rawDescData
}

var file_msg_store_v1_store_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_msg_store_v1_store_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_msg_store_v1_store_proto_goTypes = []interface{}{
	(MediaType)(0),                    // 0: msg.store.v1.MediaType
	(ContentType)(0),                  // 1: msg.store.v1.ContentType
	(*Item)(nil),                      // 2: msg.store.v1.Item
	(*Amount)(nil),                    // 3: msg.store.v1.Amount
	(*Boat)(nil),                      // 4: msg.store.v1.Boat
	(*MediaTypesInStockRequest)(nil),  // 5: msg.store.v1.MediaTypesInStockRequest
	(*MediaTypesInStockResponse)(nil), // 6: msg.store.v1.MediaTypesInStockResponse
	(*BestOfAllTimeRequest)(nil),      // 7: msg.store.v1.BestOfAllTimeRequest
	(*BestOfAllTimeResponse)(nil),     // 8: msg.store.v1.BestOfAllTimeResponse
	(*RevenueRequest)(nil),            // 9: msg.store.v1.RevenueRequest
	(*RevenueResponse)(nil),           // 10: msg.store.v1.RevenueResponse
	(*SoldItemRequest)(nil),           // 11: msg.store.v1.SoldItemRequest
	(*SoldItemResponse)(nil),          // 12: msg.store.v1.SoldItemResponse
	(*GetInStockRequest)(nil),         // 13: msg.store.v1.GetInStockRequest
	(*GetInStockResponse)(nil),        // 14: msg.store.v1.GetInStockResponse
	(*timestamppb.Timestamp)(nil),     // 15: google.protobuf.Timestamp
}
var file_msg_store_v1_store_proto_depIdxs = []int32{
	0,  // 0: msg.store.v1.Boat.media:type_name -> msg.store.v1.MediaType
	3,  // 1: msg.store.v1.Boat.price:type_name -> msg.store.v1.Amount
	1,  // 2: msg.store.v1.Boat.content:type_name -> msg.store.v1.ContentType
	0,  // 3: msg.store.v1.MediaTypesInStockResponse.in_stock:type_name -> msg.store.v1.MediaType
	1,  // 4: msg.store.v1.BestOfAllTimeRequest.content:type_name -> msg.store.v1.ContentType
	4,  // 5: msg.store.v1.BestOfAllTimeResponse.boat:type_name -> msg.store.v1.Boat
	2,  // 6: msg.store.v1.SoldItemRequest.item:type_name -> msg.store.v1.Item
	3,  // 7: msg.store.v1.SoldItemRequest.amount:type_name -> msg.store.v1.Amount
	15, // 8: msg.store.v1.SoldItemRequest.when:type_name -> google.protobuf.Timestamp
	2,  // 9: msg.store.v1.GetInStockRequest.item:type_name -> msg.store.v1.Item
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_msg_store_v1_store_proto_init() }
func file_msg_store_v1_store_proto_init() {
	if File_msg_store_v1_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_store_v1_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Item); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Amount); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Boat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaTypesInStockRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaTypesInStockResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BestOfAllTimeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BestOfAllTimeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RevenueRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RevenueResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SoldItemRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SoldItemResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInStockRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_store_v1_store_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInStockResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msg_store_v1_store_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_store_v1_store_proto_goTypes,
		DependencyIndexes: file_msg_store_v1_store_proto_depIdxs,
		EnumInfos:         file_msg_store_v1_store_proto_enumTypes,
		MessageInfos:      file_msg_store_v1_store_proto_msgTypes,
	}.Build()
	File_msg_store_v1_store_proto = out.File
	file_msg_store_v1_store_proto_rawDesc = nil
	file_msg_store_v1_store_proto_goTypes = nil
	file_msg_store_v1_store_proto_depIdxs = nil
}
