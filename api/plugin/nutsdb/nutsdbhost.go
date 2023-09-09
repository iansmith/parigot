package nutsdb

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"unicode"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/nutsdb/v1"
	nuts "github.com/nutsdb/nutsdb"
	"github.com/tetratelabs/wazero/api"
)

var nutsdblogger = slog.Default().With("source", "file", "plugin", "true")

type NutsDBPlugin struct{}

var nutsdbSvc *nutsdbSvcImpl

type nutsdbSvcImpl struct {
	datadir string
	idToDB  map[string]*nuts.DB
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
	return &nutsdbSvcImpl{
		datadir: "nutsdb",
		idToDB:  make(map[string]*nuts.DB),
	}
}
