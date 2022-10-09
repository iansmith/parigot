package util

import (
	"log"
	"strings"
)

const (
	WasmServiceName = "WasmServiceName"
	WasmFuncName    = "WasmFuncName"
)

var ParigotPrefixes = []string{
	"//parigot:",
	"// parigot:",
}

func ParigotCommentLine(allLines string) string {
	lines := strings.Split(allLines, "\n") // xxx break on windows?
	if len(lines) == 0 {
		return ""
	}
	line := ""
	for l := len(lines) - 1; l >= 0; l-- {
		if lines[l] != "" {
			line = lines[l]
			break
		}
	}
	return line
}

func ParigotCommentSettings(line string) map[string]string {
	if line == "" {
		return nil
	}
	found := ""
	for _, p := range ParigotPrefixes {
		if strings.HasPrefix(line, p) {
			found = strings.TrimPrefix(line, p)
			break
		}
	}
	if found == "" {
		return nil
	}
	parts := strings.Split(found, ",")
	if len(parts) == 0 {
		panic("cant parse comment line: " + line)
	}
	settings := make(map[string]string)
	for _, setting := range parts {
		parts := strings.Split(setting, "=")
		if len(parts) != 2 {
			log.Printf("can't parse parigot comment setting on line:" + line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		settings[key] = val
	}
	return settings
}
