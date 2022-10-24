package context

import "github.com/iansmith/parigot/g/parigot/log"

type Pctx interface {
	Log() log.Log
	Entry(group string, name string) (string, bool)
	SetEntry(group string, name string, value string) bool
}

type pctx struct {
	log   log.Log
	entry map[string]map[string]string
}

func NewPctx(l log.Log) Pctx {
	return &pctx{
		log:   l,
		entry: make(map[string]map[string]string),
	}
}
func (p *pctx) Log() log.Log {
	return p.log
}

func (p *pctx) Entry(group, name string) (string, bool) {
	g, present := p.entry[group]
	if !present {
		g = make(map[string]string)
		p.entry[group] = g
	}
	result, found := g[name]
	return result, found
}
func (p *pctx) SetEntry(group, name, value string) bool {
	g, present := p.entry[group]
	if !present {
		g = make(map[string]string)
		p.entry[group] = g
	}
	_, found := g[name]
	g[name] = value
	return found
}
