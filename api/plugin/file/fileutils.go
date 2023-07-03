package file

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	apishared "github.com/iansmith/parigot/api/shared"
)

func getRealPath(path string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory:", err)
	}
	return filepath.Join(wd[:len(wd)-23], path)
}

// A given file path is valid based on some specific rules:
// 1. The separator should be "/"
// 2. It should start with specific prefix
// 3. It should not contain any "." or ".." in the path
// 4. It should not exceed a specific value for the number of parts in the path
func isValidPath(fpath string) (string, bool) {
	fileName := filepath.Base(fpath)
	dir := strings.ReplaceAll(fpath, fileName, "")
	if !strings.HasPrefix(dir, pathPrefix) || strings.Contains(dir, ".") {
		return fpath, false
	}

	cleanPath := filepath.Clean(fpath)

	component := strings.Split(cleanPath, "/")
	if len(component) > apishared.FileServiceMaxPathPart {
		return fpath, false
	}

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
		entries, err := os.ReadDir(dir)
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
