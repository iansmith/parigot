package nutsdb

import (
	"context"
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

type nutsdbSvcImpl struct {
	datadir string
	idToDB  map[string]*nuts.DB
	lock    sync.Mutex
	parent  map[string]struct{}
}

func (*NutsDBPlugin) Init(ctx context.Context, e eng.Engine, _ id.HostId) bool {
	e.AddSupportedFunc(ctx, "nutsdb", "open_", openNutsDBHost) // call the wrapper

	return true
}

// openFileHost is a wrapper around the guest interaction code in hostBase that
// ends up calling open
func openNutsDBHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &nutsdb.OpenRequest{}
	resp := &nutsdb.OpenResponse{}

	apiplugin.HostBase(ctx, "[nutsdb]open", nutsdbSvc.open, m, stack, req, resp)
}

func (n *nutsdbSvcImpl) open(ctx context.Context, req *nutsdb.OpenRequest,
	resp *nutsdb.OpenResponse) int32 {

	n.lock.Lock()
	defer n.lock.Unlock()

	name := req.GetDbName()
	for _, n := range name {
		if !unicode.IsLetter(n) {
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
		return int32(nutsdb.NutsDBErr_InternalErr)
	}

	id := nutsdb.NewNutsDBId()
	db, e := nuts.Open(
		nuts.DefaultOptions,
		nuts.WithDir(dbDir),
	)
	if e != nil {
		nutsdblogger.Error("error opening nutsdb", "error", e)
		return int32(nutsdb.NutsDBErr_InternalErr)
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
	for _, c := range path {
		if c != '/' && !unicode.IsLetter(c) {
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
