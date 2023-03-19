package tree

type EventSectionNode struct {
	Spec    []*EventSpec
	Program *ProgramNode
}

type Selector struct {
	Id    *ValueRef
	Class *ValueRef
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

func (s *Selector) String() string {
	if s.Class.String() == "" && s.Id.String() == "" {
		return "SHOULD_NEVER_HAPPEN.UNABLE_TO_FIND_VR_VALUE"
	}
	if s.Class.String() != "" {
		return s.Class.String()
	}
	return "#" + s.Id.String()

}
func isEventName(name string) EventType {
	switch name {
	case "click", "dblclick", "mousedown", "mouseup":
		return MouseEvent
	}
	return NotDefined
}
