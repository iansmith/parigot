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
		raw, _ = os.Getwd()
	}
	cachedImportPath = filepath.SplitList(raw)
	return cachedImportPath
}

func ProtobufSearchPath(prefix string) []string {
	currentPlusImportPath := []string{}
	if prefix != "" {
		currentPlusImportPath = []string{prefix}
	}
	currentPlusImportPath = append(currentPlusImportPath, ParigotImportPath()...)
	currentPlusImportPath = append(currentPlusImportPath, "")
	return currentPlusImportPath
}
