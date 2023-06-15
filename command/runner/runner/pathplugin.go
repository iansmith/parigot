//go:build !noplugin

package runner

import (
	"io/fs"
	"os"
)

func pathExistsPlugin(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}
