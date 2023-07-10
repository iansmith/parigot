package file

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	apishared "github.com/iansmith/parigot/api/shared"
)

const pathPrefix = apishared.FileServicePathPrefix

func getRealPath(path string) string {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatal("Error getting working directory:", err)
	}

	realPath := ""
	for _, part := range strings.Split(wd, "/") {
		if part == "parigot" {
			break
		}
		realPath += part + "/"
	}

	return filepath.Join(realPath, path)
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
// Rules are:
//  1. The separator should be "/"
//  2. It should start with specific prefix
//  3. It should not contain any "." or ".." in the path
//  4. It should not exceed a specific value for the number of parts in the path
//  5. It should avoid certain special characters, including:
//     Asterisk (*)					Question mark (?)		Greater than (>)
//     Less than (<)				Pipe symbol (|)			Ampersand (&)
//     Semicolon (;)				Dollar sign ($)			Backtick (`)
//     Double quotation marks (")	Single quotation mark (')
//
// Invalid example:
//
//	'/parigot/app/..' -> '..' is not allowed
//	'/parigot/app/./' -> '.' is not allowed
//	'/parigot/app/foo\bar' -> '\' is not allowed
//	'//parigot/app/foo', '/parigot/app' -> prefix should be '/parigot/app/'
func isValidFilePath(fpath string) (string, bool) {
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
func deleteFileAndParentDirIfNeeded(path string) {
	realPath := getRealPath(path)

	log.Println("Deleting file: ", realPath)

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

// Check whether a directory exists and is valid.
func isValidDirOnHost(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		// Path is not a directory
		return false
	}

	// Directory exists and is valid
	return true
}
