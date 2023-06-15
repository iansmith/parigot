//go:build noplugin

package sys

import (
	"context"
	"log"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apiplugin/file"
	"github.com/iansmith/parigot/apiplugin/queue"
	"github.com/iansmith/parigot/apiplugin/syscall"
)

func LoadPlugin(ctx context.Context, plugin, symbol, name string) (apiplugin.ParigotInit, error) {
	log.Printf("load plugin: %s", name)
	switch name {
	case "queue":
		return &queue.QueuePlugin{}, nil
	case "file":
		return &file.FilePlugin{}, nil
	case "parigot":
		return &syscall.SyscallPlugin{}, nil
	}
	panic("unknown name for LoadPlugin:" + name)
}
