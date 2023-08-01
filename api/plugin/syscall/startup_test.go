package syscall

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
)

func TestHappyPath(t *testing.T) {
	ctx := context.Background()
	coord := makeCoordUnderTest()

	pkg := "foo"
	name := "Foo"

	expectNoService(ctx, t, coord, pkg, name)

	msg := "expected service to not exist yet: %s.%s"
	svc, created := coord.SetService(ctx, pkg, name, false)
	if !created {
		t.Errorf(msg, pkg, name)
	}
	msg = "returned service id is either zero or empty valuet: %s.%s"
	if svc.Id().IsZeroOrEmptyValue() {
		t.Errorf(msg, pkg, name)
	}

	msg = "expected service to already exist: %s.%s"
	_, created2 := coord.SetService(ctx, pkg, name, false)
	if created2 {
		t.Errorf(msg, pkg, name)
	}

	expectConsistentResult(ctx, t, coord, svc.Id(), pkg, name)
}

func TestStartupGates(t *testing.T) {
	ctx := pcontext.DevNullContext(context.Background())
	coord := makeCoordUnderTest()

	foo, _ := coord.SetService(ctx, "foo", "Foo", false)
	bar, _ := coord.SetService(ctx, "bar", "Bar", false)
	baz, _ := coord.SetService(ctx, "baz", "Baz", false)
	fleazil, _ := coord.SetService(ctx, "fleazil", "Fleazil", true)

	msg := "failed to have client 'service' already exported:%s.%s"
	if !fleazil.Exported() {
		t.Errorf(msg, "fleazil", "Fleazil")
	}

	var fooStart, barStart, bazStart, fleazilStart bool

	// fleazil->foo
	if kerr := coord.Import(ctx, fleazil.Id(), foo.Id()); kerr != syscall.KernelErr_NoError {
		log.Fatalf("bad attempt to add edge")
	}
	// bar -> baz
	if kerr := coord.Import(ctx, bar.Id(), baz.Id()); kerr != syscall.KernelErr_NoError {
		log.Fatalf("bad attempt to add edge")
	}
	// foo->baz
	if kerr := coord.Import(ctx, foo.Id(), baz.Id()); kerr != syscall.KernelErr_NoError {
		log.Fatalf("bad attempt to add edge")
	}
	// topo is: baz, (foo|bar), fleazil
	if fooStart || barStart || bazStart || fleazilStart {
		t.Errorf("should not have any service started, only export so far is a client")
		t.FailNow()
	}

	// check dep graph reasonable
	if !coord.PathExists(ctx, fleazil.Id().String(), baz.Id().String()) {
		t.Errorf("should have found a dependency edge from fleazil to baz (via foo)")
	}
	if coord.PathExists(ctx, bar.Id().String(), foo.Id().String()) {
		t.Errorf("should not have found a dependency edge from bar to foo")
	}
	if coord.PathExists(ctx, foo.Id().String(), bar.Id().String()) {
		t.Errorf("should not have found a dependency edge from foo to bar")
	}
	// if coord.PathExists(ctx, baz.Id().String(), bar.Id().String()) {
	// 	t.Errorf("should not have found a dependency edge from baz to bar (edge is backwards)")
	// }

	coord.Export(ctx, foo.Id())
	coord.Export(ctx, bar.Id())
	coord.Export(ctx, baz.Id())
	if fooStart || barStart || bazStart || fleazilStart {
		t.Errorf("should not have any service started, nobody ready to run yet")
		t.FailNow()
	}

	// these two have deps, so no way to run yet
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		foo.Run(ctx)
		fooStart = true
		wg.Add(-1)
	}()
	go func() {
		fleazil.Run(ctx)
		fleazilStart = true
		wg.Add(-1)
	}()
	if fooStart || barStart || bazStart || fleazilStart {
		t.Errorf("should not have any service started, because baz not ready")
		t.FailNow()
	}
	go func() {
		baz.Run(ctx)
		bazStart = true
		wg.Add(-1)
	}()
	wg.Wait()
	if !fooStart || barStart || !bazStart || !fleazilStart {
		t.Errorf("only bar not started")
		t.FailNow()
	}
	bar.Run(ctx) // should happen immediately, the test should finish normally
}

//
// Helpers
//

func makeCoordUnderTest() SyscallData {
	return newSyscallDataImpl()
}

func expectConsistentResult(ctx context.Context, t *testing.T, coord SyscallData, sid id.ServiceId, pkg, name string) {
	t.Helper()
	svc2, _ := coord.SetService(ctx, pkg, name, false)

	msg := "expected service id to be consistent in two uses of SetService: %s.%s"
	if !sid.Equal(svc2.Id()) {
		t.Errorf(msg, pkg, name)
	}

	msg = "expected service id to be consistent betweenf SetService and ServiceById: %s.%s"
	if svc3 := coord.ServiceById(ctx, sid); !svc3.Id().Equal(sid) {
		t.Errorf(msg, pkg, name)
	}
	msg = "expected service id to be consistent betweenf SetService and ServiceByIdString: %s.%s"
	if svc3 := coord.ServiceByIdString(ctx, sid.String()); !svc3.Id().Equal(sid) {
		t.Errorf(msg, pkg, name)
	}
	msg = "expected service id to be consistent betweenf SetService and ServiceByName: %s.%s"
	if svc3 := coord.ServiceByName(ctx, pkg, name); !svc3.Id().Equal(sid) {
		t.Errorf(msg, pkg, name)
	}
}

func expectNoService(ctx context.Context, t *testing.T, coord SyscallData, pkg, name string) {
	t.Helper()
	msg := "found service when not expecting it, %s.%s"
	if svc := coord.ServiceByName(ctx, pkg, name); svc != nil {
		t.Errorf(msg, pkg, name)
	}
	msg = "should never match an empty service id: %s.%s"
	if svc := coord.ServiceById(ctx, id.ServiceIdEmptyValue()); svc != nil {
		t.Errorf(msg, pkg, name)
	}
	msg = "should never match an zero service id: %s.%s"
	if svc := coord.ServiceById(ctx, id.ServiceIdZeroValue()); svc != nil {
		t.Errorf(msg, pkg, name)
	}

	msg = "should not match a new, random service id: %s.%s"
	if svc := coord.ServiceById(ctx, id.NewServiceId()); svc != nil {
		t.Logf(msg, pkg, name)
	}
}
