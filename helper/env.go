package helper

import (
	"os"
	"path/filepath"
)

const parigotImportVar = "PARIGOT_IMPORT_PATH"

var cachedImportPath []string

func ParigotImportPath() []string {
	if cachedImportPath != nil {
		return cachedImportPath
	}
	raw, ok := os.LookupEnv(parigotImportVar)
	if !ok {
		cachedImportPath = []string{}
	} else {
		cachedImportPath = filepath.SplitList(raw)
	}
	wd, _ := os.Getwd() // if getwd doesn't work, you have bigger problems
	cachedImportPath = append(cachedImportPath, wd)

	return cachedImportPath
}

func ProtobufSearchPath(prefix string) []string {
	currentPlusImportPath := []string{}
	// if prefix != "" {
	// 	currentPlusImportPath = []string{prefix}
	// }
	currentPlusImportPath = append(currentPlusImportPath, ParigotImportPath()...)
	currentPlusImportPath = append(currentPlusImportPath, "")
	return currentPlusImportPath
}
