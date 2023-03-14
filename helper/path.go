package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// RelativePath computes the correct path, given that the
// provided path is relative to file
func RelativePath(rel, origin, pkg string) string {
	if pkg == "" {
		dir := filepath.Dir(origin)
		result := filepath.Join(dir, rel)
		log.Printf("computing relative (no package) path of %s to %s => %s", rel, origin, result)
		return filepath.Clean(result)
	}
	dir := StripEndsOfPathForPkg(pkg, origin)
	log.Printf("strip ends called (%s,%s)=>%s", pkg, rel, dir)
	result := filepath.Join(dir, rel)

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

func StripEndsOfPathForPkg(pkg, path string) string {
	pkgPart := strings.Split(pkg, ".")
	dir := filepath.Dir(path)
	dir = filepath.Clean(dir)
	elem := strings.Split(dir, string(os.PathSeparator))
	if elem == nil {
		panic("empty path given to strip ends")
	}
	for len(pkgPart) > 0 {
		lastPart := pkgPart[len(pkgPart)-1]
		lastElem := elem[len(elem)-1]
		if lastElem != lastPart {
			break
		}
		pkgPart = pkgPart[:len(pkgPart)-1]
		elem = elem[:len(elem)-1]
	}
	if len(pkgPart) != 0 {
		panic(fmt.Sprintf("had package parts left: %s", filepath.Join(pkgPart...)))
	}
	return filepath.Join(elem...)

}
