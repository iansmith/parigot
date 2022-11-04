package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

func readLibList(engine *wasmtime.Engine, modToPath map[*wasmtime.Module]string) ([]*wasmtime.Module, error) {
	libFp, err := os.Open(*libFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s:%v", *libFile, err)
	}
	defer func() {
		libFp.Close()
	}()
	scanner := bufio.NewScanner(libFp)
	mod := []*wasmtime.Module{}
	for scanner.Scan() {
		line := scanner.Text()
		m, err := loadSingleModule(engine, line)
		if err != nil {
			return nil, err
		}
		modToPath[m] = line
		mod = append(mod, m)
	}
	if scanner.Err() != nil {
		log.Fatalf("failed complete reading the lib file: %v", scanner.Err())
	}
	return mod, nil
}

func loadSingleModule(engine *wasmtime.Engine, path string) (*wasmtime.Module, error) {
	path = strings.TrimSpace(path)
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("unable to find %s: %v", path, err)
		} else {
			return nil, fmt.Errorf("unable to stat %s: %v", path, err)
		}
	}
	m, err := wasmtime.NewModuleFromFile(engine, path)
	if err != nil {
		return nil, fmt.Errorf("unable to convert %s into a module: %v",
			path, err)
	}
	return m, nil
}

func walkArgs(engine *wasmtime.Engine, modToPath map[*wasmtime.Module]string) []*wasmtime.Module {
	result := []*wasmtime.Module{}
	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		m, err := loadSingleModule(engine, path)
		if err != nil {
			log.Fatalf("command line file argument #%d failed: %v", i, err)
		}
		modToPath[m] = path //for err mesgs
		result = append(result, m)
	}
	return result
}
