//go:build noplugin

package sys

import (
	"context"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/plugin/file"
	"github.com/iansmith/parigot/api/plugin/httpconnector"
	"github.com/iansmith/parigot/api/plugin/queue"
	"github.com/iansmith/parigot/api/plugin/syscall"
)

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
	}
	panic("unknown name for LoadPlugin:" + name)
}
