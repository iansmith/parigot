package main

import (
	"log"
	"strconv"
	"strings"
)

type SexprBuilder struct {
	*BasesexprListener
	Filename string
}

type Atom struct {
	Symbol     string
	String     string
	Number     int
	CommentNum bool
}

type Item struct {
	Atom *Atom
	List []*Item
}

var _ sexprListener = &SexprBuilder{}

func NewSexprBuilder(f string) *SexprBuilder {
	base := &BasesexprListener{}
	return &SexprBuilder{BasesexprListener: base, Filename: f}
}

func (s *SexprBuilder) ExitAtom(ctx *AtomContext) {
	if ctx.STRING() != nil {
		a := &Atom{String: ctx.STRING().GetText()}
		ctx.SetAtom_(a)
		return
	}
	if ctx.SYMBOL() != nil {
		a := &Atom{Symbol: ctx.SYMBOL().GetText()}
		ctx.SetAtom_(a)
		return
	}
	if ctx.NUMBER() != nil {
		t := ctx.NUMBER().GetText()
		a, err := strconv.Atoi(t)
		if err != nil {
			panic("badly formed number:" + t)
		}
		ctx.SetAtom_(&Atom{Number: a})
		return
	}
	if ctx.COMMENT_NUM() != nil {
		t := ctx.COMMENT_NUM().GetText()
		t = strings.TrimPrefix(strings.TrimSuffix(t, ";"), ";")
		a, err := strconv.Atoi(t)
		if err != nil {
			panic("badly formed number:" + t)
		}
		ctx.SetAtom_(&Atom{Number: a, CommentNum: true})
		return
	}
}
func (s *SexprBuilder) ExitList_(ctx *List_Context) {
	raw := ctx.AllItem()
	item_ := make([]*Item, len(raw))
	for i, item := range raw {
		item_[i] = item.GetItem_()
	}
	ctx.SetList(item_)
}

func (s *SexprBuilder) ExitItem(ctx *ItemContext) {
	item := &Item{}
	found := false
	if ctx.Atom() != nil {
		item.Atom = ctx.Atom().GetAtom_()
		found = true
	}
	if ctx.List_() != nil {
		item.List = ctx.List_().GetList()
		found = true
	}
	if !found {
		panic("unable to understand item (not atom or list)")
	}

	ctx.SetItem_(item)
}

func (s *SexprBuilder) ExitSexpr(ctx *SexprContext) {
	raw := ctx.AllItem()
	item := make([]*Item, len(raw))
	log.Printf("sexpr size %d", len(raw))
	for i, c := range raw {
		item[i] = c.GetItem_()
	}
	ctx.SetItem_(item)
}
