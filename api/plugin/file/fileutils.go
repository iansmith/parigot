package file

import (
	"os"
	"path/filepath"
	"strings"

	apishared "github.com/iansmith/parigot/api/shared"
)

const pathPrefix = apishared.FileServicePathPrefix

func getRealPath(path string) (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	realPath := ""
	for _, part := range strings.Split(wd, "/") {
		if part == "parigot" {
			break
		}
		realPath += part + "/"
	}

	return filepath.Join(realPath, path), nil
}

// Checks if the given string contains any illegal characters.
func containsIllegalChars(s string) bool {
	illegalChar := "*?><|&;$\\`\"'"
	return strings.ContainsAny(s, illegalChar)
}

// Checks if the given path starts with the expected prefix.
func startsWithPrefix(path, prefix string) bool {
	return strings.HasPrefix(path, prefix)
}

// Checks if the given path contains "." or "..".
func containsDisallowedPathComponents(path string) bool {
	fileName := filepath.Base(path)
	dir := path[:len(path)-len(fileName)-1]
	return strings.Contains(dir, ".") || strings.Contains(fileName, "..") || fileName == "."
}

// Checks if the given path exceeds the maximum number of allowed parts.
func exceedsMaxParts(path string) bool {
	parts := strings.Split(path, "/")
	return len(parts) > apishared.FileServiceMaxPathPart
}

// Validates if a given file path complies with specific rules.
func isValidPath(fpath string) (string, bool) {
	if containsIllegalChars(fpath) || !startsWithPrefix(fpath, pathPrefix) {
		return fpath, false
	}

	if containsDisallowedPathComponents(fpath) {
		return fpath, false
	}

	cleanPath := filepath.Clean(fpath)
	if exceedsMaxParts(cleanPath) {
		return fpath, false
	}

	return cleanPath, true
}

func isValidBuf(buf []byte) bool { return len(buf) <= maxBufSize }

// Deletes the specified file and its parent directories if they're empty
func deleteFileAndParentDirIfNeeded(path string) error {
	realPath, err := getRealPath(path)
	if err != nil {
		return err
	}

	// Delete the file
	err = os.Remove(realPath)
	if err != nil {
		return err
	}

	// Walk up the directory tree and remove any empty directories.
	dir := filepath.Dir(realPath)
	for {
		// Read the directory.
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		// If the directory is not empty, we're done.
		if len(entries) > 0 {
			break
		}

		// Delete the directory and move to its parent.
		if err := os.Remove(dir); err != nil {
			return err
		}
		dir = filepath.Dir(dir)
	}
	return nil
}

// Check whether a directory exists and is valid.
func isValidDirOnHost(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		// Path is not a directory
		return false, err
	}

	// Directory exists and is valid
	return true, nil
}
