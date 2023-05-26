package id

const (
	// FileNoError means just what it sounds like.  All Ids that are errors represent
	// no error as 0.
	FileNoError FileErrorCode = 0
	// FileBadPath means that the given path (filename) is not valid.
	FileBadPath FileErrorCode = 1
	// FileNotFound means that the given path could point to a file (it is valid)
	// but the path given could not be found on the filesystem.
	FileNotFound FileErrorCode = 2
)
