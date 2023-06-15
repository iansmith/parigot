//go:build !noplugin

package main

import (
	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apiplugin/file"
)

var ParigotInitialize apiplugin.ParigotInit = &file.FilePlugin{}

// This is the "entrypoint" of the .so when it is loaded by the dynamic
// loader. The symbol that parigot searches for is above.  That symbol
// is used directly when in "noplugin" mode.
func main() {
	// we are connected to the runner via a symbol lookup
}
