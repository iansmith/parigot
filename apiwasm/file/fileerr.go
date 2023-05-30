package main

import (
	"github.com/iansmith/parigot/g/file/v1"
)

const (
	// FileBadPath means that the given path (filename) is not valid.
	FileBadPath file.FileErrIdCode = iota + file.FileErrIdGuestStart
	// FileNotFound means that the given path could point to a file (it is valid)
	// but the path given could not be found on the filesystem.
	FileNotFound
)
