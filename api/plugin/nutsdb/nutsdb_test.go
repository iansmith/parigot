package nutsdb

import (
	"context"
	"fmt"
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
	helperBadDBName(t, ctx, impl, "7foo")

	helperDbNotExist(t, ctx, impl, "quux")
	helperDbNotExist(t, ctx, impl, "quux2")
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
		DbName:        "foobie7",
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
		"",
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

func TestReadWritePair(t *testing.T) {
	impl := newNutsDBImpl()
	ctx := context.Background()

	path := mustMakeTempDir(t)
	impl.datadir = path // just because we are a test, not for user consumption
	defer func() {
		os.RemoveAll(path)
	}()

	// tests
	id := mustOpenNutsdb(t, ctx, impl)

	// create a request to push a k,v pair
	k := "somekey"
	v := "somevalue"

	// happy path
	mustWritePair(t, ctx, impl, id, "/foo", k, v)
	result := mustReadPair(t, ctx, impl, id, "/foo", k, "")
	if result != v {
		t.Logf("failed to read back the value we put in %s", result)
	}
	// key not there
	_, err := readErr(t, ctx, impl, id, "/foo", k+"?", "")
	if err != nutsdb.NutsDBErr_PairNotFound {
		t.Logf("expected a pair not found error, but got %s", nutsdb.NutsDBErr_name[int32(err)])
	}
	// key not there, but we provided default
	def := "fleazil"
	valWithDefault, err := readErr(t, ctx, impl, id, "/foo", k+"?", def)
	if err != nutsdb.NutsDBErr_NoError {
		t.Errorf("expected a pair to be ok with default %s", nutsdb.NutsDBErr_name[int32(err)])
	}
	if valWithDefault != def {
		t.Errorf("expected a value that matches the value put in (%s): %s", def, v)
	}
	mustCloseDb(t, ctx, impl, id)
}
func TestBucketSemantics(t *testing.T) {
	impl := newNutsDBImpl()
	ctx := context.Background()

	path := mustMakeTempDir(t)
	impl.datadir = path // just because we are a test, not for user consumption
	defer func() {
		os.RemoveAll(path)
	}()

	// tests
	id := mustOpenNutsdb(t, ctx, impl)
	base := "/foo/bar"

	req := &nutsdb.WritePairRequest{}
	resp := &nutsdb.WritePairResponse{}

	req.NutsdbId = id.Marshal()
	req.Pair = &nutsdb.Pair{
		Key:   []byte("hello"),
		Value: []byte("world"),
	}

	for i := 0; i < 10; i++ {
		//normall two pairs with same key value would be an error, but the
		// buckets form namespaces
		req.Pair.BucketPath = fmt.Sprintf("%s%d", base, i)
		if err := impl.writePair(ctx, req, resp); err != int32(nutsdb.NutsDBErr_NoError) {
			t.Errorf("expected write of pair into bucket %s would be ok: %s", req.Pair.BucketPath,
				nutsdb.NutsDBErr_name[err])
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

func mustOpenNutsdb(t *testing.T, ctx context.Context, impl *nutsdbSvcImpl) nutsdb.NutsDBId {
	t.Helper()
	dbname := "test1"
	req := &nutsdb.OpenRequest{DbName: dbname}
	resp := &nutsdb.OpenResponse{}
	if err := impl.open(ctx, req, resp); err != 0 {
		t.Errorf("unable to open nutsdb: %v: %d", dbname, err)
		t.FailNow()
		return nutsdb.NutsDBIdZeroValue()
	}
	return nutsdb.UnmarshalNutsDBId(resp.GetNutsdbId())
}
func mustWritePair(t *testing.T, ctx context.Context, impl *nutsdbSvcImpl, id nutsdb.NutsDBId, bucketPath, key, value string) {
	t.Helper()
	pair := &nutsdb.Pair{
		BucketPath: bucketPath,
		Key:        []byte(key),
		Value:      []byte(value),
	}
	req := &nutsdb.WritePairRequest{NutsdbId: id.Marshal(), Pair: pair}
	resp := &nutsdb.WritePairResponse{}
	e := impl.writePair(ctx, req, resp)
	if e != 0 {
		t.Errorf("unable to write key value pair (%s,%s,%s): %s", bucketPath, key, value,
			nutsdb.NutsDBErr_name[e])
		t.FailNow()
	}
}

func mustReadPair(t *testing.T, ctx context.Context, impl *nutsdbSvcImpl, id nutsdb.NutsDBId, bucketPath, key, value string) string {
	result, e := readErr(t, ctx, impl, id, bucketPath, key, value)
	if e != nutsdb.NutsDBErr_NoError {
		t.Errorf("unable to read key value pair (%s,%s): %s", bucketPath, key,
			nutsdb.NutsDBErr_name[int32(e)])
		t.FailNow()
		return ""
	}
	return result
}

func readErr(t *testing.T, ctx context.Context, impl *nutsdbSvcImpl, id nutsdb.NutsDBId, bucketPath, key, value string) (string, nutsdb.NutsDBErr) {
	valueByte := []byte{}
	if value != "" {
		valueByte = []byte(value)
	}
	pair := &nutsdb.Pair{
		BucketPath: bucketPath,
		Key:        []byte(key),
		Value:      valueByte,
	}
	req := &nutsdb.ReadPairRequest{NutsdbId: id.Marshal(), Pair: pair}
	resp := &nutsdb.ReadPairResponse{Pair: &nutsdb.Pair{}}
	e := impl.readPair(ctx, req, resp)
	if e != 0 {
		ourErr := nutsdb.NutsDBErr(e)
		return "", ourErr
	}

	return string(resp.GetPair().GetValue()), nutsdb.NutsDBErr_NoError
}

func mustCloseDb(t *testing.T, ctx context.Context, impl *nutsdbSvcImpl, id nutsdb.NutsDBId) {
	req := &nutsdb.CloseRequest{NutsdbId: id.Marshal()}
	resp := &nutsdb.CloseResponse{}
	if err := impl.close(ctx, req, resp); err != 0 {
		t.Errorf("unable to close nutsdb: %v", nutsdb.NutsDBErr_name[err])
		t.FailNow()
	}
}
