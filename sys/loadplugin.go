//go:build !noplugin

package sys

import (
	"context"

	"github.com/iansmith/parigot/apiplugin"
)

func LoadPlugin(ctx context.Context, plugin, symbol, name string) (apiplugin.ParigotInit, error) {
	return apiplugin.LoadAndReturnInit(ctx, plugin, symbol, name)
}
