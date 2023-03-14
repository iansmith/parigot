package helper

import (
	"os"
	"path/filepath"
)

// RelativePath computes the correct path, given that the
// provided path is relative to file
func RelativePath(path, file string) string {
	dir := filepath.Dir(file)
	result := filepath.Join(dir, path)
	return filepath.Clean(result)
}

// FindProtobufFile returns "" if it cannot find the protocol buffer
// specified by the FQName.  If the prefix is not "" that
// that is a directory that should be added to the search path.
func FindProtobufFile(name, prefix string) string {
	checked := name
	found := ""
	if len(checked) > 0 && ((checked[0] != '/') || (checked[0] != '\\')) {
		importPath := ProtobufSearchPath(prefix)
		for _, candidate := range importPath {
			path := filepath.Join(candidate, checked)
			_, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
			}
			found = path
			break
		}
	} else {
		// fully qualified path
		_, err := os.Stat(name)
		if err != nil {
			if os.IsNotExist(err) {
				return ""
			}
		}
		found = name
	}
	return found
}
