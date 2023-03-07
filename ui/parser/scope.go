package parser

// ScopeStack keeps track of the currently being compiled structure.
type ScopeStack struct {
	Stack []Scope
	Top   int
}

func NewScopeStack() *ScopeStack {
	raw := make([]Scope, maxScopeStackSize)
	return &ScopeStack{
		Stack: raw[:],
		Top:   0,
	}
}

func (s *ScopeStack) PushScope(scope Scope) {
	s.Stack[s.Top] = scope
	s.Top++
}

func (s *ScopeStack) PopScope() Scope {
	if s.Top == 0 {
		panic("attempt to pop a scope from an empty scope stack")
	}
	result := s.Stack[s.Top-1]
	s.Top--
	return result
}
