package nutsdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"unicode"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/nutsdb/v1"

	nuts "github.com/nutsdb/nutsdb"
	"github.com/tetratelabs/wazero/api"
)

var nutsdblogger = slog.Default().With("plugin", "nutsdb")

type NutsDBPlugin struct{}

var nutsdbSvc *nutsdbSvcImpl

const maxBucketPathLen = 4096
const maxKeySize = 4096
const maxValueSize = 4096 << 3

type nutsdbSvcImpl struct {
	datadir string
	idToDB  map[string]*nuts.DB
	lock    sync.Mutex
	parent  map[string]struct{}
}

type wrappedParigotErr struct {
	err int32
	e   error
}

func (w *wrappedParigotErr) Error() string {
	return w.e.Error()
}

func (*NutsDBPlugin) Init(ctx context.Context, e eng.Engine, _ id.HostId) bool {
	e.AddSupportedFunc(ctx, "nutsdb", "open_", openNutsDBHost)            // call the wrapper
	e.AddSupportedFunc(ctx, "nutsdb", "close_", closeNutsDBHost)          // call the wrapper
	e.AddSupportedFunc(ctx, "nutsdb", "read_pair_", readPairNutsDBHost)   // call the wrapper
	e.AddSupportedFunc(ctx, "nutsdb", "write_pair_", writePairNutsDBHost) // call the wrapper

	// he is the one, true service
	nutsdbSvc = newNutsDBImpl()

	return true
}

func (n *nutsdbSvcImpl) open(ctx context.Context, req *nutsdb.OpenRequest,
	resp *nutsdb.OpenResponse) int32 {

	n.lock.Lock()
	defer n.lock.Unlock()

	//k.matcher().Dispatch(targetHid, hid, cid, mid, hostFunc, w)

	name := req.GetDbName()
	for i, n := range name {
		if !unicode.IsLetter(n) && !(i > 0 && unicode.IsNumber(n)) {
			return int32(nutsdb.NutsDBErr_BadDBName)
		}
	}
	dbDir := filepath.Join(n.datadir, name)
	stat, err := os.Stat(dbDir)
	if err != nil && req.GetErrIfNotFound() {
		return int32(nutsdb.NutsDBErr_DBNotFound)
	}
	// stat is nil if it does not exist
	if stat != nil && !stat.IsDir() {
		return int32(nutsdb.NutsDBErr_InternalError)
	}

	id := nutsdb.NewNutsDBId()

	opts := nuts.DefaultOptions
	//opts.SyncEnable = false
	db, e := nuts.Open(
		opts,
		nuts.WithDir(dbDir),
	)
	if e != nil {
		nutsdblogger.Error("error opening nutsdb", "error", e)
		return int32(nutsdb.NutsDBErr_InternalError)
	}
	n.idToDB[id.String()] = db
	resp.NutsdbId = id.Marshal()

	return int32(nutsdb.NutsDBErr_NoError)
}

func newNutsDBImpl() *nutsdbSvcImpl {
	impl := &nutsdbSvcImpl{
		datadir: "nutsdb",
		idToDB:  make(map[string]*nuts.DB),
		parent:  make(map[string]struct{}),
	}
	impl.parent["/"] = struct{}{}
	return impl
}

func (n *nutsdbSvcImpl) close(ctx context.Context, req *nutsdb.CloseRequest,
	resp *nutsdb.CloseResponse) int32 {

	n.lock.Lock()
	defer n.lock.Unlock()

	nid := nutsdb.UnmarshalNutsDBId(req.GetNutsdbId())
	if nid.IsZeroOrEmptyValue() {
		return int32(nutsdb.NutsDBErr_BadId)
	}

	db, ok := n.idToDB[nid.String()]
	if !ok {
		return int32(nutsdb.NutsDBErr_DBNotFound)
	}

	if e := db.Close(); e != nil {
		nutsdblogger.Error("error closing nutsdb", "error", e)
	}
	delete(n.idToDB, nid.String())

	return int32(nutsdb.NutsDBErr_NoError)
}

func (n *nutsdbSvcImpl) writePair(ctx context.Context, req *nutsdb.WritePairRequest,
	resp *nutsdb.WritePairResponse) int32 {

	nid := nutsdb.UnmarshalNutsDBId(req.GetNutsdbId())
	if nid.IsZeroOrEmptyValue() {
		return int32(nutsdb.NutsDBErr_BadId)
	}
	ok, _ := n.isValidBucketPath(ctx, req.Pair.GetBucketPath())
	if !ok {
		return int32(nutsdb.NutsDBErr_BadBucketPath)
	}
	bpath := req.Pair.GetBucketPath()
	if bpath == "" {
		bpath = "/"
	}
	db, ok := n.idToDB[nid.String()]
	if !ok {
		return int32(nutsdb.NutsDBErr_BadId)
	}
	l := checkLength(req.Pair.GetBucketPath(),
		req.Pair.GetKey(), req.Pair.GetValue())
	if l != nutsdb.NutsDBErr_NoError {
		return int32(l)
	}
	// end of preamble
	err := db.Update(
		func(tx *nuts.Tx) error {
			k := make([]byte, len(req.Pair.GetKey()))
			copy(k, req.Pair.GetKey())
			v := make([]byte, len(req.Pair.GetValue()))
			copy(v, req.Pair.GetValue())
			b := req.Pair.GetBucketPath()
			err := tx.Put(b, k, v, nuts.Persistent)
			return err
		})
	if err != nil {
		perr, ok := err.(*wrappedParigotErr)
		if !ok {
			inner := errors.Unwrap(err)
			nutsdblogger.Error("unable to understand returned error from nutsdb, not a parigot error", "error", err, "type", fmt.Sprintf("%T", err), "inner", inner, "inner type", fmt.Sprintf("%T", inner))
			return int32(nutsdb.NutsDBErr_InternalError)
		}
		return perr.err
	}
	resp.Pair = &nutsdb.Pair{}
	resp.Pair.Key = make([]byte, len(req.GetPair().GetKey()))
	copy(resp.GetPair().GetKey(), req.GetPair().GetKey())
	resp.Pair.Value = make([]byte, len(req.GetPair().GetValue()))
	copy(resp.GetPair().GetValue(), req.GetPair().GetValue())
	resp.GetPair().BucketPath = req.GetPair().GetBucketPath()

	return int32(nutsdb.NutsDBErr_NoError)
}

func (n *nutsdbSvcImpl) readPair(ctx context.Context, req *nutsdb.ReadPairRequest,
	resp *nutsdb.ReadPairResponse) int32 {

	nid := nutsdb.UnmarshalNutsDBId(req.GetNutsdbId())
	if nid.IsZeroOrEmptyValue() {
		return int32(nutsdb.NutsDBErr_BadId)
	}
	ok, _ := n.isValidBucketPath(ctx, req.Pair.GetBucketPath())
	if !ok {
		return int32(nutsdb.NutsDBErr_BadBucketPath)
	}
	bpath := req.Pair.GetBucketPath()
	if bpath == "" {
		bpath = "/"
	}

	db, ok := n.idToDB[nid.String()]
	if !ok {
		return int32(nutsdb.NutsDBErr_BadId)
	}
	// end of preamble

	wasNotFound := false
	l := checkLength(req.Pair.GetBucketPath(), req.Pair.GetKey(), req.Pair.GetValue())
	if l != nutsdb.NutsDBErr_NoError {
		return int32(l)
	}
	resp.Pair = &nutsdb.Pair{}
	raw := db.View(
		func(tx *nuts.Tx) error {
			k := req.Pair.GetKey()
			b := bpath

			value, err := tx.Get(b, k)
			if err == nil {
				resp.Pair.Value = make([]byte, len(value.Value))
				copy(resp.Pair.Value, value.Value)
				resp.Pair.BucketPath = req.Pair.GetBucketPath()
				resp.Pair.Key = make([]byte, len(req.Pair.GetKey()))
				copy(resp.Pair.Key, req.Pair.Key)
				return nil
			}
			if nuts.IsKeyNotFound(err) {
				// arg the returned error is actually a Rollback error because
				// we are in a transaction and nutsdb is not using Wrap/Unwrap
				// so we use this hack, ugh
				wasNotFound = true
			}
			return err
		})
	if raw == nil {
		return int32(nutsdb.NutsDBErr_NoError)
	}
	// this is awful, see above
	if wasNotFound {
		if len(req.Pair.Value) != 0 {
			resp.Pair.Value = req.Pair.Value
			return int32(nutsdb.NutsDBErr_NoError)
		}
		return int32(nutsdb.NutsDBErr_PairNotFound)
	}

	return int32(nutsdb.NutsDBErr_InternalError)
}

func (n *nutsdbSvcImpl) createBucketParent(_ context.Context, path string) bool {
	_, ok := n.parent[path]
	if ok {
		nutsdblogger.Warn("creating directory as parent multiple times", "directory", path)
		return true
	}
	n.parent[path] = struct{}{}
	return true
}

func (n *nutsdbSvcImpl) isValidBucketPath(ctx context.Context, path string) (bool, string) {
	//special case
	if path == "" {
		return true, ""
	}

	for _, c := range path {
		if c != '/' && !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			nutsdblogger.Warn("bad character in bucket path", "path", path)
			return false, ""
		}
	}
	path = filepath.Clean(path)
	dir := filepath.Dir(path)
	if dir[0] != '/' {
		nutsdblogger.Warn("not fully qualified bucket path (must start with /)", "path", path)
		return false, ""
	}
	if _, ok := n.parent[dir]; !ok {
		if !n.createBucketParent(ctx, dir) {
			nutsdblogger.Error("unable to create parent bucket", "directory", dir)
			return false, ""
		}
	}
	return true, dir
}

// boiler plate to hook up the native functions

// openFileHost is a wrapper around the guest interaction code in hostBase that
// ends up calling open
func openNutsDBHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &nutsdb.OpenRequest{}
	resp := &nutsdb.OpenResponse{}
	apiplugin.HostBase(ctx, "[nutsdb]open", nutsdbSvc.open, m, stack, req, resp)
}

func closeNutsDBHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &nutsdb.CloseRequest{}
	resp := &nutsdb.CloseResponse{}

	apiplugin.HostBase(ctx, "[nutsdb]close", nutsdbSvc.close, m, stack, req, resp)
}

func readPairNutsDBHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &nutsdb.ReadPairRequest{}
	resp := &nutsdb.ReadPairResponse{}

	apiplugin.HostBase(ctx, "[nutsdb]readPair", nutsdbSvc.readPair, m, stack, req, resp)
}

func writePairNutsDBHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &nutsdb.WritePairRequest{}
	resp := &nutsdb.WritePairResponse{}

	apiplugin.HostBase(ctx, "[nutsdb]writePair", nutsdbSvc.writePair, m, stack, req, resp)
}

func checkLength(path string, k []byte, v []byte) nutsdb.NutsDBErr {
	if len(path) > maxBucketPathLen {
		return nutsdb.NutsDBErr_BucketPathTooLong
	}
	if len(k) > maxKeySize {
		return nutsdb.NutsDBErr_KeyTooLarge
	}
	if len(v) > maxValueSize {
		return nutsdb.NutsDBErr_ValueTooLarge
	}
	return nutsdb.NutsDBErr_NoError
}
