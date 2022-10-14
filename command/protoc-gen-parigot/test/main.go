package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	flag.Parse()

	if flag.NArg() != 4 {
		log.Fatalf("arguments are <path to binary> <inputfile> <comparison file> <expected output> ")
	}

	cmd := exec.Cmd{
		Path: flag.Arg(0),
		Args: []string{"-t", "-l", flag.Arg(1)},
		Env:  nil,
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("FAIL:\n%s\n", out)
		log.Printf("ERROR\n%v", err)
	}
	outBuffer := bufio.NewScanner(bytes.NewBuffer(out))
	var testBuffer *bufio.Scanner
	if flag.Arg(2) != "-" { // for debugging
		fp, err := os.Open(flag.Arg(3))
		if err != nil {
			log.Fatalf("%v", err)
		}
		buf, err := io.ReadAll(fp)
		if err != nil {
			log.Fatalf("%v", err)
		}
		testBuffer = bufio.NewScanner(bytes.NewBuffer(buf))
	}

	// cue up out buffer
	for {
		filename := readUntilEOFOrNextFile(outBuffer, false)
		if filename == "" {
			log.Fatalf("unexpected eof in output looking for file %s", flag.Arg(3))
		}
		if strings.TrimSpace(filename) == strings.TrimSpace(flag.Arg(3)) {
			break
		}
	}
	if testBuffer == nil { // useful for debugging, just dump it
		_ = readUntilEOFOrNextFile(outBuffer, true)
		os.Exit(0)
	}
	var currentTestLine, currentOutLine string
	for {
		nextOut, ok := nextNonBlankLine(outBuffer)
		if !ok {
			//off chance of perfect sync
			_, ok := nextNonBlankLine(testBuffer)
			if !ok {
				os.Exit(0)
			}
			log.Fatalf("unexpected eof in output buffer after line %s", currentOutLine)
		}
		currentOutLine = nextOut
		nextTest, ok := nextNonBlankLine(testBuffer)
		if !ok {
			// all the lines matched until eof of test file
			os.Exit(0)
		}
		currentTestLine = nextTest
		if strings.TrimSpace(currentOutLine) != strings.TrimSpace(currentTestLine) {
			log.Printf("failed to match output")
			log.Printf("output:")
			log.Printf("%s", currentOutLine)
			log.Printf("test:")
			log.Printf("%s", currentTestLine)
			os.Exit(1)
		}
		// line is ok
	}
}

func nextNonBlankLine(s *bufio.Scanner) (string, bool) {
	for {
		ok := s.Scan()
		if !ok {
			if s.Err() == nil {
				return "", false
			} else {
				log.Fatalf("err cueing up output %v", s.Err())
			}
		}
		line := s.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		return line, true
	}
}

func readUntilEOFOrNextFile(scanner *bufio.Scanner, dumpToStdout bool) string {
	for {
		line, ok := nextNonBlankLine(scanner)
		if !ok {
			return ""
		}
		if dumpToStdout {
			fmt.Fprint(os.Stdout, line+"\n")
		}
		// look for marker
		if strings.HasPrefix(line, "***") {
			f := strings.TrimPrefix(line, "***")
			return f
		}
	}

}
