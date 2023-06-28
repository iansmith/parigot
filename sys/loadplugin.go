//go:build !noplugin

package sys

import (
	"context"

	apiplugin "github.com/iansmith/parigot/api/plugin"
)

func LoadPlugin(ctx context.Context, plugin, symbol, name string) (apiplugin.ParigotInit, error) {
	return apiplugin.LoadAndReturnInit(ctx, plugin, symbol, name)
}
