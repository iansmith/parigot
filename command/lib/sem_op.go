package lib

import "fmt"

type Op interface {
	IndentedStringer
	OpType() OpT
}
type OpT int

const (
	ZeroT         OpT = 1
	Int1T         OpT = 2
	I32LoadStoreT OpT = 3
	I64LoadStoreT OpT = 4
	Id1T          OpT = 5
)

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

type Int1Op struct {
	Op  string
	Arg int
}

type BrIfOp struct {
	*Int1Op
	BranchTarget int
}

func (i *Int1Op) OpType() OpT {
	return Int1T
}

func (i *Int1Op) StmtType() StmtT {
	return OpStmtT
}

func (i *Int1Op) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("%s %d", i.Op, i.Arg))
	return "\n" + buf.String()

}

type I64LoadStore struct {
	IsStore bool
	Align   *int
	Offset  *int
}

func (i *I64LoadStore) OpType() OpT {
	return I64LoadStoreT
}

func (i *I64LoadStore) StmtType() StmtT {
	return OpStmtT
}

func (i *I64LoadStore) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	op := "i64.load"
	if i.IsStore {
		op = "i64.store"
	}
	buf.WriteString(op)
	if i.Offset != nil {
		buf.WriteString(fmt.Sprintf(" offset=%d", i.Offset))
	}
	if i.Align != nil {
		buf.WriteString(fmt.Sprintf(" align=%d", i.Align))
	}
	return "\n" + buf.String()
}
func (i *I64LoadStore) SetOffset(offset int) {
	i.Offset = new(int)
	*i.Offset = offset
}
func (i *I64LoadStore) SetAlign(align int) {
	i.Align = new(int)
	*i.Align = align
}

type I32LoadStore struct {
	IsStore bool
	Offset  *int
}

func (i *I32LoadStore) OpType() OpT {
	return I32LoadStoreT
}

func (i *I32LoadStore) StmtType() StmtT {
	return OpStmtT
}

func (i *I32LoadStore) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	op := "i32.load"
	if i.IsStore {
		op = "i32.store"
	}
	buf.WriteString(op)
	if i.Offset != nil {
		buf.WriteString(fmt.Sprintf(" offset=%d", i.Offset))
	}
	return "\n" + buf.String()
}

func (i *I32LoadStore) SetOffset(offset int) {
	i.Offset = new(int)
	*i.Offset = offset
}

type Id1Op struct {
	Op  string
	Arg string
}

func (i *Id1Op) OpType() OpT {
	return Id1T
}

func (i *Id1Op) StmtType() StmtT {
	return OpStmtT
}

func (i *Id1Op) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("%s %s", i.Op, i.Arg))
	return "\n" + buf.String()
}
