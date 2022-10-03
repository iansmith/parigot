package lib

import (
	"bytes"
	"fmt"
	"math"
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
	PreviousResult    *ResultDef
	nestingLevel      int
	Result            *ResultDef
	Code              []Stmt
}

func (b *BlockStmt) StmtType() StmtT {
	return BlockStmtT
}

func (b *BlockStmt) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	controlStmtToString("block", b.nestingLevel, b.Result, buf)
	buf.WriteString(stmtsToString(b.Code, indented+2))
	outputEnd(indented, buf)
	return buf.String()
}

func (b *BlockStmt) AddStmt(s Stmt) {
	b.Code = append(b.Code, s)
}

type IfStmt struct {
	nestingLevel      int
	PreviousContainer Container
	PreviousResult    *ResultDef
	Result            *ResultDef
	IfPart            []Stmt
	ElsePart          []Stmt
}

func (i *IfStmt) StmtType() StmtT {
	return IfStmtT
}

func (i *IfStmt) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	controlStmtToString("if", i.nestingLevel, i.Result, buf)
	buf.WriteString(stmtsToString(i.IfPart, indented+2))
	if i.ElsePart != nil {
		buf.WriteString("\n") //terminate previous line
		for j := 0; j < indented; j++ {
			buf.WriteString(" ")
		}
		controlStmtToString("else", math.MinInt, nil, buf)
		buf.WriteString(stmtsToString(i.ElsePart, indented+2))
	}
	outputEnd(indented, buf)
	return buf.String()
}

func (i *IfStmt) AddStmt(s Stmt) {
	if i.ElsePart != nil {
		i.ElsePart = append(i.ElsePart, s)
	} else {
		i.IfPart = append(i.IfPart, s)
	}
}

type LoopStmt struct {
	*BlockStmt
}

func (l *LoopStmt) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	controlStmtToString("loop", l.nestingLevel, l.Result, buf)
	buf.WriteString(stmtsToString(l.Code, indented+2))
	return buf.String()
}

// helpers

func stmtsToString(stmt []Stmt, indented int) string {
	var buf bytes.Buffer
	for _, s := range stmt {
		buf.WriteString("\n" + s.IndentedString(indented))
	}
	return buf.String()
}

func controlStmtToString(name string, nestingLevel int, result *ResultDef, buf *bytes.Buffer) string {
	buf.WriteString(name)
	if result != nil {
		buf.WriteString(" " + result.String())
	}
	if nestingLevel > 0 {
		buf.WriteString(fmt.Sprintf("  ;; label = @%d",
			nestingLevel))
	}
	return buf.String()
}

func outputEnd(indented int, buf *bytes.Buffer) {
	buf.WriteString("\n")
	for i := 0; i < indented; i++ {
		buf.WriteString(" ")
	}
	buf.WriteString("end")
}
