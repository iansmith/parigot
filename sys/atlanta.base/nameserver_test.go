package sys

import (
	"testing"
)

func dummyProc(name string, number int) *Process {
	return &Process{
		path: "./" + name,
		id:   number,
	}
}

func TestSimpleDead(t *testing.T) {
	ns := NewLocalNameServer()

	dummyA := dummyProc("a.wasm", 1)
	dummyB := dummyProc("b.wasm", 2)

	a := &EdgeHolder{
		proc:    dummyA,
		require: []string{"foo.b"},
	}
	b := &EdgeHolder{
		proc:    dummyB,
		require: []string{"foo.a"},
	}
	ns.dependencyGraph = map[string]*EdgeHolder{
		a.proc.String(): a,
		b.proc.String(): b,
	}
	if ns.GetLoopContent() != "" {
		t.Errorf("expected no loop from graph with no exports!")
	}
	if ns.GetDeadNodeContent() == "" {
		t.Errorf("expected loop with nodes a and b but got nothing")
	}
}
func TestSimpleLoop(t *testing.T) {
	ns := NewLocalNameServer()

	redh := dummyProc("redherring.wasm", 3)
	dummyA := dummyProc("a.wasm", 1)
	dummyB := dummyProc("b.wasm", 2)

	a := &EdgeHolder{
		proc:    dummyA,
		require: []string{"foo.b"},
		export:  []string{"foo.a"},
	}
	b := &EdgeHolder{
		proc:    dummyB,
		require: []string{"foo.a"},
		export:  []string{"foo.b"},
	}
	r := &EdgeHolder{
		proc:    redh,
		require: []string{"foo.b"},
	}
	ns.dependencyGraph = map[string]*EdgeHolder{
		a.proc.String(): a,
		b.proc.String(): b,
		r.proc.String(): r,
	}
	t.Logf("loop is:%s", ns.GetLoopContent())

	if ns.GetLoopContent() == "" {
		t.Errorf("expected loop of a and b but got nothing!")
	}
	if ns.GetDeadNodeContent() != "" {
		t.Errorf("expected no dead nodes in simple loop bc red herring would get freed by b:\n" + ns.GetDeadNodeContent())
	}
	t.Fail()
}

func TestLongChain(t *testing.T) {
	ns := NewLocalNameServer()

	redh := dummyProc("redherring.wasm", 3)
	dummyA := dummyProc("a.wasm", 1)
	dummyB := dummyProc("b.wasm", 2)
	dummyC := dummyProc("c.wasm", 3)
	dummyD := dummyProc("d.wasm", 4)
	dummyE := dummyProc("e.wasm", 5)

	a := &EdgeHolder{
		proc:    dummyA,
		require: []string{"foo.b"},
		export:  []string{"foo.a"},
	}
	b := &EdgeHolder{
		proc:    dummyB,
		require: []string{"foo.c"},
		export:  []string{"foo.b"},
	}
	r := &EdgeHolder{
		proc:    redh,
		require: []string{"foo.b"},
	}
	c := &EdgeHolder{
		proc:    dummyC,
		require: []string{"foo.d"},
		export:  []string{"foo.c"},
	}
	d := &EdgeHolder{
		proc:    dummyD,
		require: []string{"foo.e"},
		export:  []string{"foo.d"},
	}
	e := &EdgeHolder{
		proc:   dummyE,
		export: []string{"foo.e"},
	}
	ns.dependencyGraph = map[string]*EdgeHolder{
		a.proc.String(): a,
		b.proc.String(): b,
		r.proc.String(): r,
		c.proc.String(): c,
		d.proc.String(): d,
		e.proc.String(): e,
	}

	//log.Printf("loop ---> '%s'", ns.GetLoopContent())

	//no loop no dead
	if ns.GetLoopContent() != "" {
		t.Errorf("did not expect loop of long chain!")
	}
	if ns.GetDeadNodeContent() != "" {
		t.Errorf("expected no dead nodes in long chain")
	}

	// now we are going to add an edge that will create a loop
	e.require = []string{"foo.a"}
	if ns.GetLoopContent() == "" {
		t.Errorf("expected loop after adding edge!")
	}
	if ns.GetDeadNodeContent() != "" {
		t.Errorf("expected no dead nodes in long chain")
	}

}
