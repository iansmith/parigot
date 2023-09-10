//go:build noplugin

package sys

import (
	"context"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/plugin/file"
	"github.com/iansmith/parigot/api/plugin/httpconnector"
	"github.com/iansmith/parigot/api/plugin/nutsdb"
	"github.com/iansmith/parigot/api/plugin/queue"
	"github.com/iansmith/parigot/api/plugin/syscall"
)

// LoadPlugin is kind of a hack.  This is here so that if you want to run
// WITHOUT dynamic loading you can do that and link directly to the services
// this function knows about.  This is helpful because debuggers tend to behave
// badly with .so files in go.  Note that the build tags guarantee that this
// version is only used with the noplugin tag set.
func LoadPlugin(ctx context.Context, plugin, symbol, name string) (apiplugin.ParigotInit, error) {
	switch name {
	case "queue":
		return &queue.QueuePlugin{}, nil
	case "file":
		return &file.FilePlugin{}, nil
	case "parigot":
		return &syscall.SyscallPlugin{}, nil
	case "httpconnector":
		return &httpconnector.HttpConnectorPlugin{}, nil
	case "nutsdb":
		return &nutsdb.NutsDBPlugin{}, nil
	}
	panic("unknown name for LoadPlugin:" + name)
}
