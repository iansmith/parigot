package lib

import (
	"bytes"
	"fmt"
)

type Stmt interface {
	IndentedStringer
	StmtType() StmtT
}

type Container interface {
	AddStmt(s Stmt)
}

type StmtT int

const (
	OpStmtT    StmtT = 1
	BlockStmtT StmtT = 2
	IfStmtT    StmtT = 3
)

type BlockStmt struct {
	PreviousContainer Container
	nestingLevel      int
	Result            *ResultDef
	Block             []Stmt
}

func (b *BlockStmt) StmtType() StmtT {
	return BlockStmtT
}

func (b *BlockStmt) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(controlStmtToString("block", indented, b.nestingLevel, b.Result))
	buf.WriteString(stmtsToString(b.Block, indented))
	buf.WriteString("\nend")
	return buf.String()
}

func (b *BlockStmt) AddStmt(s Stmt) {
	b.Block = append(b.Block, s)
}

type IfStmt struct {
	nestingLevel      int
	PreviousContainer Container
	Result            *ResultDef
	IfPart            []Stmt
	ElsePart          []Stmt
}

func (i *IfStmt) StmtType() StmtT {
	return IfStmtT
}

func (i *IfStmt) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(controlStmtToString("if", indented, i.nestingLevel, i.Result))
	buf.WriteString(stmtsToString(i.IfPart, indented+2))
	if i.ElsePart != nil {
		buf.WriteString(controlStmtToString("else", indented, i.nestingLevel, nil))
		buf.WriteString(stmtsToString(i.ElsePart, indented+2))
	}
	buf.WriteString("\nend")
	return buf.String()
}

func (i *IfStmt) AddStmt(s Stmt) {
	if i.ElsePart != nil {
		i.ElsePart = append(i.ElsePart, s)
	} else {
		i.IfPart = append(i.IfPart, s)
	}
}

func stmtsToString(stmt []Stmt, indented int) string {
	var buf bytes.Buffer
	for _, s := range stmt {
		// weird, they prefix the \n not suffix it
		buf.WriteString("\n")
		buf.WriteString(s.IndentedString(indented + 2))
	}
	return buf.String()
}

func controlStmtToString(name string, indented int, nestingLevel int, result *ResultDef) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(name)
	if result != nil {
		buf.WriteString(" " + result.String())
	}
	buf.WriteString(fmt.Sprintf(";; label = @%d\n",
		nestingLevel))
	return buf.String()
}
