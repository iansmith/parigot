package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getRealPath(path string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory:", err)
	}
	return filepath.Join(wd[:len(wd)-23], path)
}

// A valid path should be a shortest path name equivalent to path by purely lexical processingand.
// Specifically, it should start with "/parigot/app/", also, any use of '.', '..', in the path is
// not allowed.
func isValidPath(fpath string) (string, bool) {
	fileName := filepath.Base(fpath)
	dir := strings.ReplaceAll(fpath, fileName, "")
	if !strings.HasPrefix(dir, pathPrefix) || strings.Contains(dir, ".") {
		return fpath, false
	}
	cleanPath := filepath.Clean(fpath)

	return cleanPath, true
}

func isValidBuf(buf []byte) bool { return len(buf) <= maxBufSize }

// Deletes the specified file and its parent directories if they're empty
func deleteFileAndParentDirIfNeeded(path string) {
	realPath := getRealPath(path)

	// Delete the file
	err := os.Remove(realPath)
	if err != nil {
		log.Fatalf("Failed to delete file: %s. Error: %v", path, err)
	}

	// Walk up the directory tree and remove any empty directories.
	dir := filepath.Dir(realPath)
	for {
		// Read the directory.
		entries, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal("Failed to read dir: ", err)
		}

		// If the directory is not empty, we're done.
		if len(entries) > 0 {
			break
		}

		// Delete the directory and move to its parent.
		if err := os.Remove(dir); err != nil {
			log.Fatal("Failed to remove dir: ", err)
		}
		dir = filepath.Dir(dir)
	}
}
