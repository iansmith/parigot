package apiplugin

import (
	"context"
	"fmt"
	"log"
	"plugin"

	"github.com/iansmith/parigot/eng"
)

// ParigotInit is the interface that plugins must meet to be
// initialized. It is expected that they will use the supplied
// Engine in the call to Init to register Host functions.
type ParigotInit interface {
	Init(ctx context.Context, e eng.Engine) bool
}

// LoadAndReturnInit is a utility function for plugins that
// want the default implementation.  This function accepts third
// string param (name) but ignores it.
func LoadAndReturnInit(ctx context.Context, pluginPath, pluginSymbol, _ string) (ParigotInit, error) {
	log.Printf("LoadAndRetInit %s, %s", pluginPath, pluginSymbol)
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open plugin %v: %v", plug, err)
	}
	sym, err := plug.Lookup(pluginSymbol)
	if err != nil {
		return nil, fmt.Errorf("unable to find symbol %s: %v", pluginSymbol, err)
	}
	initFn := sym.(*ParigotInit)
	return *initFn, nil
}
