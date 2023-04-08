package tree

type EventSectionNode struct {
	Spec    []*EventSpec
	Section *MVCSectionNode
	Program *ProgramNode
}

type Selector struct {
	Id    *ValueRef
	Class *ValueRef
}

func (s *Selector) String() string {
	if s.Id != nil {
		return s.Id.String()
	}
	return s.Class.String()
}

type EventCall struct {
	Invoc *FuncInvoc
}

type EventType int32

const (
	NotDefined   = 0
	MouseEvent   = 1
	PointerEvent = 2
)

type EventSpec struct {
	Selector  *Selector
	EventName string
	Function  *FuncInvoc
	EventType EventType
}

func NewEventSpec(s *Selector, name string, b *FuncInvoc) *EventSpec {

	t := isEventName(name)
	if t == NotDefined {
		return nil
	}
	return &EventSpec{s, name, b, t}
}

func isEventName(name string) EventType {
	switch name {
	case "click", "dblclick", "mousedown", "mouseup":
		return MouseEvent
	}
	return NotDefined
}
func NewEventSectionNode(p *ProgramNode) *EventSectionNode {
	return &EventSectionNode{
		Spec:    []*EventSpec{},
		Section: &MVCSectionNode{},
		Program: p,
	}
}
