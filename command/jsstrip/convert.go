package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	wasmSuffix      = ".wasm"
	watSuffix       = ".wat"
	wasmProgram     = "wasm2wat"
	watProgram      = "wat2wasm"
	parigotFilename = "parigot_transformed.wat"
)

// wat2wasm takes an argument for it's output file, puts errors on stderr
func convertWatToWasm(tmpDir, source string, target *string) error {
	extraArgs := []string{}
	var tempTargetFile string
	if *target == "" {
		var tempFd *os.File
		var err error
		tempFd, tempTargetFile, err = createFileInTmpdir(tmpDir, watProgram, false)
		if err != nil {
			log.Printf("unable to create temporary file for wasm output: %v", err)
			return err
		}
		_ = tempFd.Close()
		extraArgs = []string{source, "-o", tempTargetFile}
	} else {
		extraArgs = []string{source, "-o", *target}
	}
	errFp, errFile, err := createFileInTmpdir(tmpDir, watProgram, true)
	if err != nil {
		log.Printf("unable to create temporary file for wasm output: %v", err)
		return err
	}
	cmd := exec.Command(watProgram, extraArgs...)
	defer func() {
		_ = errFp.Close()
	}()
	cmd.Stderr = errFp
	err = cmd.Run()
	if err != nil {
		if *target == "" {
			_ = os.Remove(*target) // don't want to confuse make
		}
		log.Printf("conversion of %s to wasm failed,\n\t errors in %s: %v", source, errFile, err)
		return err
	}
	if *target == "" {
		fp, err := os.Open(tempTargetFile)
		if err != nil {
			log.Printf("unable to read file %s to print to stdout: %v", tempTargetFile, err)
		}
		_, err = io.Copy(os.Stdout, fp)
		if err != nil {
			log.Printf("unable to copy file %s content tto stdout: %v", tempTargetFile, err)
		}
	}
	return nil
}

// wasm2wat puts its output on stdout for the converted result
// does it put errors on stdout or stderr?
func convertWasmToWat(tmpDir string, source string) (string, error) {
	targetFp, target, err := createFileInTmpdir(tmpDir, wasmNameToWatName(basename(source)), false)
	if err != nil {
		log.Printf("converting input file ("+source+") failed, cannot create temp file: %v", err)
		return "", err
	}
	defer func() {
		_ = targetFp.Close()
	}()
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("converting input file ("+source+") failed, input file does not exist: %v", err)
		return "", err
	}
	cmd := exec.Command(wasmProgram, source)
	cmd.Stdout = targetFp
	errFp, errPath, err := createFileInTmpdir(tmpDir, wasmProgram, true)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = errFp.Close()
	}()
	cmd.Stderr = errFp

	err = cmd.Run()
	if err != nil {
		log.Printf("conversion of %s to wat failed,\n\terrors in %s (or possibly %s): %v", source, target, errPath, err)
		return "", err
	}
	return target, nil
}

func basename(path string) string {
	_, b := filepath.Split(path)
	return b
}

func wasmNameToWatName(name string) string {
	if len(name) <= len(wasmSuffix) || !strings.HasSuffix(name, wasmSuffix) {
		return name + ".wat"
	}
	return strings.TrimSuffix(name, wasmSuffix) + watSuffix
}

func createFileInTmpdir(tmpDir, base string, isErrorFile bool) (*os.File, string, error) {
	// stderr file
	path := base
	if isErrorFile {
		path = path + "-errors"
	}
	errFile := filepath.Join(tmpDir, path)
	errFp, err := os.Create(errFile)
	if err != nil {
		log.Printf("converting input file cannot create temporary error file: %v", err)
		return nil, errFile, err
	}
	return errFp, errFile, nil
}

func parigotProcessing(inputFilename, tmpDir string) (string, error) {
	mod := parse(inputFilename)
	transformation(mod)
	fp, path, err := createFileInTmpdir(tmpDir, parigotFilename, false)
	if err != nil {
		log.Printf("unable to create intermediate output file: %v", err)
		return "", err
	}
	defer func() {
		_ = fp.Close()
	}()
	_, err = fp.WriteString(mod.IndentedString(0))
	if err != nil {
		log.Printf("unable to write to intermediate output file: %v", err)
		return "", err
	}
	return path, nil
}
