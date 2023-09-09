package nutsdb

import (
	"context"
	"os"
	"testing"

	"github.com/iansmith/parigot/g/nutsdb/v1"
)

func TestOpen(t *testing.T) {
	impl := newNutsDBImpl()
	ctx := context.Background()

	path, err := os.MkdirTemp("/tmp", "nutsdb")
	if err != nil {
		t.Errorf("unable to create temp directory for nutsdb: %v", err)
		t.FailNow()
	}
	impl.datadir = path // just because we are a test, not for user consumption

	// tests

	helperBadDBName(t, ctx, impl, "/foo")
	helperBadDBName(t, ctx, impl, "../foo")
	helperBadDBName(t, ctx, impl, "foo/./bar")

	helperDbNotExist(t, ctx, impl, "quux")
	goodDBName := "fleazil"
	helperDbNotExist(t, ctx, impl, goodDBName)

	req := &nutsdb.OpenRequest{}
	resp := &nutsdb.OpenResponse{}

	req.DbName = goodDBName
	req.ErrIfNotFound = false

	// do a real open
	e := impl.open(ctx, req, resp)
	nutsErr := nutsdb.NutsDBErr(e)
	if nutsErr != nutsdb.NutsDBErr_NoError {
		t.Errorf("unexpected error in db open: %s", nutsdb.NutsDBErr_name[e])
	}
	id := nutsdb.UnmarshalNutsDBId(resp.GetNutsdbId())
	t.Logf("got back an id %s", id.Short())
}

//
// HELPERS
//

func helperBadDBName(t *testing.T, ctx context.Context, impl *nutsdbSvcImpl, badName string) {
	req := &nutsdb.OpenRequest{}
	resp := &nutsdb.OpenResponse{}

	t.Helper()

	req.DbName = badName
	req.ErrIfNotFound = true // doesn't matter
	e := impl.open(ctx, req, resp)
	err := nutsdb.NutsDBErr(e)
	if err != nutsdb.NutsDBErr_BadDBName {
		t.Errorf("failed to reject based on characters in DbName: %s", badName)
	}

}

func helperDbNotExist(t *testing.T, ctx context.Context, impl *nutsdbSvcImpl, name string) {
	req := &nutsdb.OpenRequest{}
	resp := &nutsdb.OpenResponse{}

	t.Helper()

	req.DbName = name
	req.ErrIfNotFound = true // doesn't matter
	e := impl.open(ctx, req, resp)
	err := nutsdb.NutsDBErr(e)
	if err != nutsdb.NutsDBErr_DBNotFound {
		t.Errorf("failed to raise correct error when checking for db not exist: %s", name)
	}
}
