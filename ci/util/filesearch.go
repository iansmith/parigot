package util

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
)

func FindFilesWithSuffixRecursively(path string, suffix string) (bool, error) {
	/*
	 *	This is a helper function that recursively finds all files in the
	 *	current folder and subfolders based on their suffix names
	 */
	exist := false
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			file := filepath.Join(path, suffix)
			return errors.New(file + " does not exist")
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), suffix) {
			exist = true
		}
		return nil

	})

	return exist, err
}

func FindFilesWithPattern(pattern string) (bool, error) {
	/*
	 *	This is a helper function that finds all files based on
	 *	wildcard matching
	 */
	files, err := filepath.Glob(pattern)
	exist := false
	if err != nil {
		return exist, err
	}
	if len(files) == 0 {
		return exist, errors.New(pattern + " does not exist")
	}

	exist = true

	return exist, nil
}
