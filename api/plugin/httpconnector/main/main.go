//go:build !noplugin

package main

import (
	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/plugin/httpconnector"
)

var ParigotInitialize apiplugin.ParigotInit = &httpconnector.HttpConnectorPlugin{}

// This is the "entrypoint" of the .so when it is loaded by the dynamic
// loader. The symbol that parigot searches for is above.  That symbol
// is used directly when in "noplugin" mode.
func main() {
	// we are connected to the runner via a symbol lookup
	// the guest is running the event loop
}
