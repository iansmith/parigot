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

	path := mustMakeTempDir(t)
	impl.datadir = path // just because we are a test, not for user consumption
	defer func() {
		os.RemoveAll(path)
	}()

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
	if id.IsZeroOrEmptyValue() {
		t.Errorf("bad id returned by open")
	}

}

func TestClose(t *testing.T) {
	impl := newNutsDBImpl()
	ctx := context.Background()

	path := mustMakeTempDir(t)
	impl.datadir = path // just because we are a test, not for user consumption
	defer func() {
		os.RemoveAll(path)
	}()

	// tests

	badId := nutsdb.NutsDBIdZeroValue()
	req := &nutsdb.CloseRequest{}
	resp := &nutsdb.CloseResponse{}

	req.NutsdbId = badId.Marshal()
	e := impl.close(ctx, req, resp)
	err := nutsdb.NutsDBErr(e)
	if err != nutsdb.NutsDBErr_BadId {
		t.Errorf("tried to close the zero value for db id and got %s", nutsdb.NutsDBErr_name[e])
	}

	openReq := &nutsdb.OpenRequest{
		DbName:        "foobie",
		ErrIfNotFound: false,
	}
	openResp := &nutsdb.OpenResponse{}
	e = impl.open(ctx, openReq, openResp)
	err = nutsdb.NutsDBErr(e)
	if err != nutsdb.NutsDBErr_NoError {
		t.Errorf("unexpected error opening db: %s", nutsdb.NutsDBErr_name[e])
	}
	req.NutsdbId = openResp.GetNutsdbId()
	e = impl.close(ctx, req, resp)
	if e != 0 {
		nid := nutsdb.UnmarshalNutsDBId(req.GetNutsdbId())
		t.Errorf("unable to close db %s: %s", nid.Short(), nutsdb.NutsDBErr_name[e])
	}

	e = impl.close(ctx, req, resp)
	err = nutsdb.NutsDBErr(e)

	if err != nutsdb.NutsDBErr_DBNotFound {
		nid := nutsdb.UnmarshalNutsDBId(req.GetNutsdbId())
		t.Errorf("expected to not be able to lose db %s: %s", nid.Short(), nutsdb.NutsDBErr_name[e])
	}

}

func TestBucketPath(t *testing.T) {
	impl := newNutsDBImpl()
	ctx := context.Background()

	path := mustMakeTempDir(t)
	impl.datadir = path // just because we are a test, not for user consumption
	defer func() {
		os.RemoveAll(path)
	}()

	// tests

	fail := []string{
		"foo",
		"foo/bar",
		"/foo$",
		"/foo.bar",
	}
	good := []string{
		"/",
		"/foo",
		"//",
		"//baz//",
		"/baz/bar",
		"/baz/bar/quuux/",
		"/baz/bar/quuux",
	}
	for _, f := range fail {
		ok, _ := impl.isValidBucketPath(ctx, f)
		if ok {
			t.Errorf("did not expect %s to be a valid bucket path", f)
		}
	}
	for _, f := range good {
		ok, _ := impl.isValidBucketPath(ctx, f)
		if !ok {
			t.Errorf("expected %s to be a valid bucket path", f)
		}
	}
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

func mustMakeTempDir(t *testing.T) string {
	t.Helper()
	path, err := os.MkdirTemp("/tmp", "nutsdb")
	if err != nil {
		t.Errorf("unable to create temp directory for nutsdb: %v", err)
		t.FailNow()
	}
	return path
}
