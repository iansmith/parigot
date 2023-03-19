package tree

type MVCSectionNode struct {
	Program   *ProgramNode
	ModelDecl []*ModelDecl
	ViewDecl  []*ViewDecl
}

func NewMvcSectionNode(p *ProgramNode) *MVCSectionNode {
	return &MVCSectionNode{Program: p}
}

type ModelDecl struct {
	Name string
	Path []string
	File []*ProtobufFileNode
}

func NewModelDecl() *ModelDecl {
	m := &ModelDecl{}
	GCurrentModel = m
	return m
}

type ViewDecl struct {
	Name      string
	ModelName string
	DocFn     *DocFuncNode
}

func NewViewDecl(view string) *ViewDecl {
	return &ViewDecl{Name: view}
}

type ProtobufFileNode struct {
	PackageName string
	FileName    string
	GoPkg       string
	LocalGoPkg  string
	ImportFile  []string
	Import      []*ProtobufFileNode
	Message     []*ProtobufMessage
}

func NewProtobufFileNode() *ProtobufFileNode {
	return &ProtobufFileNode{}
}

type ProtobufMessage struct {
	Name string
}

func NewProtobufMessage(name string) *ProtobufMessage {
	return &ProtobufMessage{Name: name}
}

type FullIdent struct {
	Part []string
}

func NewFullIdent(c []string) *FullIdent {
	return &FullIdent{Part: c}
}

type OptionTriple struct {
	Name         *FullIdent
	Value, Extra string
}

func NewOptionTriple() *OptionTriple {
	return &OptionTriple{}
}
