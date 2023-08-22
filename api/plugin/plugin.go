package plugin

import (
	"context"
	"fmt"
	"log"
	"os"
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
	if pluginPath == "" {
		return nil, fmt.Errorf("cannot load a plugin with no path provided")
	}
	if pluginSymbol == "" {
		return nil, fmt.Errorf("cannot load a plugin with no symbol provided")
	}
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

type Plog struct {
	*log.Logger
}

func (k *Plog) Errorf(spec string, rest ...interface{}) {
	k.Logger.Printf("ERR : "+spec, rest...)
}
func (k *Plog) Warnf(spec string, rest ...interface{}) {
	k.Logger.Printf("WARN: "+spec, rest...)
}
func (k *Plog) Infof(spec string, rest ...interface{}) {
	k.Logger.Printf("INFO: "+spec, rest...)
}

func NewPluginLogger(pluginName string) *Plog {
	return &Plog{
		Logger: log.New(os.Stdout, pluginName+":", log.Default().Flags()),
	}
}
