package main

import (
	"github.com/iansmith/parigot/command/transform"
)

// addToplevelToModule adds a TopLevel t into the code of a module m.  If module m
// has no TopLevel's that are the same type as t, this function panics.
func addToplevelToModule(m *transform.Module, t transform.TopLevel, atEnd bool) {
	location := findTopLevelLocation(m, t.TopLevelType(), atEnd)
	m.Code = append(m.Code[:location], append([]transform.TopLevel{t}, m.Code[location:]...)...)
}

func findTopLevelLocation(m *transform.Module, t transform.TopLevelT, atEnd bool) int {
	if !atEnd {
		// xxx assumes all the TLs of a type are together
		for i := 0; i < len(m.Code); i++ {
			if m.Code[i].TopLevelType() == t {
				return i
			}
		}
		// xxx should we just put it after the imports?
		panic("can't place new TopLevel at beginning of group, unable to find any TopLevel in module to combine with same type ")
	} else {
		if len(m.Code) == 0 || len(m.Code) == 1 {
			return len(m.Code)
		}
		// xxx assumes all the top levels of a type are together
		prevType := m.Code[0].TopLevelType()
		for i := 1; i < len(m.Code); i++ {
			if prevType == t && m.Code[i].TopLevelType() != t {
				return i
			}
			prevType = m.Code[i].TopLevelType()
		}
		if m.Code[len(m.Code)-1].TopLevelType() == t {
			return len(m.Code) - 1
		}
		panic("unable to place new top level type at end of the group, maybe there are no other top levels of that type?")
	}
}

// changeStmtCodeOnly walks all the stmts in a sequence and then calls fn at each stmt.
// If fn returns anything different from the input, that becomes the new stmt.  If fn returnns
// nil, the stmt is elided. Modified version of the sequence of statements is returned.
func changeStmtCodeOnly(code []transform.Stmt, fn func(stmt transform.Stmt) transform.Stmt) []transform.Stmt {
	newCode := []transform.Stmt{}
	for _, stmt := range code {
		result := fn(stmt)
		if result != stmt {
			if result != nil {
				newCode = append(newCode, result)
			}
		} else {
			newCode = append(newCode, result)
		}
		// recurse through the code blocks that are nested
		if stmt.StmtType() == transform.IfStmtT ||
			stmt.StmtType() == transform.BlockStmtT {
			if stmt.StmtType() == transform.IfStmtT {
				ifStmt := stmt.(*transform.IfStmt)
				if ifStmt.IfPart != nil {
					changeStmtCodeOnly(ifStmt.IfPart, fn)
				}
				if ifStmt.ElsePart != nil {
					changeStmtCodeOnly(ifStmt.ElsePart, fn)
				}
			}
			if stmt.StmtType() == transform.BlockStmtT {
				bl, ok := stmt.(*transform.BlockStmt)
				if ok {
					if bl.Code != nil {
						changeStmtCodeOnly(bl.Code, fn)
					}
				}
				loop, ok := stmt.(*transform.LoopStmt)
				if ok {
					bl = loop.BlockStmt
					if bl.Code != nil {
						changeStmtCodeOnly(bl.Code, fn)
					}
				}

			}
		}
	}
	return newCode
}

// changeStatementInModule walks all the code in a module, calling fn at each statement.
// If fn returns anything different from the input, that becomes the new stmt. If fn returns
// nil, the statement is elided.
func changeStatementInModule(m *transform.Module, fn func(stmt transform.Stmt) transform.Stmt) {
	for _, candidate := range m.Code {
		if candidate.TopLevelType() != transform.FuncDefT {
			continue
		}
		changeStmtCodeOnly(candidate.(*transform.FuncDef).Code, fn)
	}
}

// changeTopLevelInModule walks all the TopLevel entities in a module and calls fn if the candidate
// TopLevel has the same type as the given tlType.  If fn returns anything different from the input,
// that becomes the new stmt.  If fn returns nil, the TopLevel is elided.
func changeToplevelInModule(m *transform.Module, tlType transform.TopLevelT, fn func(t transform.TopLevel) transform.TopLevel) {
	newTL := []transform.TopLevel{}
	for _, tl := range m.Code {
		if tl.TopLevelType() != tlType {
			newTL = append(newTL, tl)
			continue
		}
		result := fn(tl)
		if result != tl {
			if result == nil {
				continue
			}
		}
		newTL = append(newTL, result)
	}
	m.Code = newTL
}

func findToplevelInModule(mod *transform.Module, tlType transform.TopLevelT, fn func(transform.TopLevel)) {
	for _, tl := range mod.Code {
		if tl.TopLevelType() != tlType {
			continue
		}
		fn(tl)
	}
}