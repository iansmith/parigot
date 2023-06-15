//go:build noplugin

package runner

import (
	"io/fs"
)

func pathExistsPlugin(path string) (fs.FileInfo, error) {
	return nil, nil
}
