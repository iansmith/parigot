package transform

import (
	"fmt"
)

type Op interface {
	IndentedStringer
	OpType() OpT
}
type OpT int

const (
	ZeroT         OpT = 1
	ArgT          OpT = 2
	LoadStoreT    OpT = 3
	CallT         OpT = 4
	BrTableT      OpT = 5
	IndirectCallT OpT = 6
	MutOpT        OpT = 7
)

type SpecialIdT int

const (
	StackPointer SpecialIdT = 1
)

func (s *SpecialIdT) String() string {
	return "$__stack_pointer"
}

type ZeroOp struct {
	Op string
}

func (z *ZeroOp) OpType() OpT {
	return ZeroT
}

func (z *ZeroOp) StmtType() StmtT {
	return OpStmtT
}

func (z *ZeroOp) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(z.Op)
	return buf.String()
}

type ArgOp struct {
	Op         string
	IntArg     *int
	FloatArg   *string
	BranchAnno *int
	ConstAnno  *string
	Special    *SpecialIdT
}

func (i *ArgOp) OpType() OpT {
	return ArgT
}

func (i *ArgOp) StmtType() StmtT {
	return OpStmtT
}

func (i *ArgOp) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	var arg string
	switch {
	case i.Special != nil:
		arg = i.Special.String()
	case i.IntArg != nil:
		arg = fmt.Sprint(*i.IntArg)
	case i.FloatArg != nil:
		arg = fmt.Sprint(*i.FloatArg)
	}
	buf.WriteString(fmt.Sprintf("%s %s", i.Op, arg))
	if i.BranchAnno != nil {
		buf.WriteString(fmt.Sprintf(" (;@%d;)", *i.BranchAnno))
	}
	if i.ConstAnno != nil {
		buf.WriteString(fmt.Sprintf(" (;=%s;)", *i.ConstAnno))
	}
	return buf.String()
}

type LoadStoreOp struct {
	Op     string
	Align  *int
	Offset *int
}

func (m *LoadStoreOp) OpType() OpT {
	return LoadStoreT
}

func (m *LoadStoreOp) StmtType() StmtT {
	return OpStmtT
}
func (m *LoadStoreOp) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	align := ""
	offset := ""
	if m.Align != nil {
		align = fmt.Sprintf(" align=%d", *m.Align)
	}
	if m.Offset != nil {
		offset = fmt.Sprintf(" offset=%d", *m.Offset)
	}
	buf.WriteString(fmt.Sprintf("%s%s%s", m.Op, offset, align))
	return buf.String()
}

type IndirectCallOp struct {
	Type *TypeRef
}

func (i *IndirectCallOp) OpType() OpT {
	return IndirectCallT
}

func (i *IndirectCallOp) StmtType() StmtT {
	return OpStmtT
}
func (i *IndirectCallOp) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("call_indirect %s", i.Type.String()))
	return buf.String()
}

type CallOp struct {
	ArgName *string
	ArgNum  *int
}

func (i *CallOp) OpType() OpT {
	return CallT
}

func (i *CallOp) StmtType() StmtT {
	return OpStmtT
}

func (i *CallOp) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	if i.ArgNum == nil && i.ArgName == nil {
		panic("call op has neither name or number")
	}
	if i.ArgNum == nil {
		buf.WriteString(fmt.Sprintf("call %s", *i.ArgName))
	} else {
		buf.WriteString(fmt.Sprintf("call %s", *i.ArgNum))
	}
	return buf.String()
}

type BranchTarget struct {
	Num    int
	Branch int
}

type BrTableOp struct {
	Target []*BranchTarget
}

func (i *BrTableOp) OpType() OpT {
	return BrTableT
}

func (i *BrTableOp) StmtType() StmtT {
	return OpStmtT
}

func (b *BrTableOp) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString("br_table")
	for _, t := range b.Target {
		buf.WriteString(fmt.Sprintf(" %d (;@%d;)", t.Num, t.Branch))
	}
	return buf.String()
}

type MutOp struct {
	Type string
}

func (g *MutOp) OpType() OpT {
	return MutOpT
}

func (g *MutOp) StmtType() StmtT {
	return OpStmtT
}

func (g *MutOp) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString("mut " + g.Type)
	return buf.String()
}
