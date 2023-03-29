package tree

type ProtobufMsgElem struct {
	Name  string
	Field *ProtobufMsgField
	Map   *ProtobufMapField
}

func NewProtobufMsgElem(f *ProtobufMsgField, m *ProtobufMapField) *ProtobufMsgElem {
	result := &ProtobufMsgElem{Field: f}
	if f != nil {
		result.Name = f.Name
	} else {
		result.Name = m.Name
	}
	return result
}

type ProtobufMsgBody struct {
	Elem []*ProtobufMsgElem
}

func NewProtobufMsgBody(e []*ProtobufMsgElem) *ProtobufMsgBody {
	return &ProtobufMsgBody{e}
}

type ProtobufMsgField struct {
	IsOptional, IsRepeated bool
	Name                   string
	Type                   string
	TypeBase               bool
	Message                *ProtobufMessage
}

func NewProtobufMsgField(opt, rep bool) *ProtobufMsgField {
	return &ProtobufMsgField{IsOptional: opt, IsRepeated: rep}

}

type ProtobufMapField struct {
	KeyType       string
	ValueType     string
	Name          string
	ValueTypeBase bool
}

func NewProtobufMapField(key, value, name string) *ProtobufMapField {
	return &ProtobufMapField{KeyType: key, ValueType: value, Name: name}

}
