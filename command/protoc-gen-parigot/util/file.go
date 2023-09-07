package util

import (
	"bufio"
	"bytes"
	_ "embed"
	"log"
	"os"
	"strings"
)

//go:embed syslib.txt
var syslib []byte

//go:embed pkg.txt
var ignored []byte

func IsSystemLibrary(s string) bool {
	if wantsSysLib() {
		return false
	}
	return isInEmbeddedList(s, syslib)
}

func isInEmbeddedList(s string, buf []byte) bool {
	rd := bytes.NewBuffer(buf)

	scanner := bufio.NewScanner(rd)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		//log.Printf("xxxxx testing %s vs %s (%v)", line, s, line == s)
		if line == s {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("failed reading the syslib.txt file: %v", err)
	}
	return false
}

func IsIgnoredPackage(s string) bool {
	if wantsSysLib() {
		return false
	}
	return isInEmbeddedList(s, ignored)
}

func wantsSysLib() bool {
	return os.Getenv("GEN_SYS_LIB") != ""
}
