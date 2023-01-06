package runner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
)

type DeployConfig struct {
	Server   []*WasmFile
	Test     []*WasmFile
	TestMode bool
	Remote   bool
	module   map[string]*wasmtime.Module
	wasmFile map[string]*WasmFile
}

type WasmFile struct {
	Name     string
	Path     string
	remote   bool
	isServer bool
	isMain   bool
}

const maxServer = 32

func Parse(path string) (*DeployConfig, error) {
	var result DeployConfig
	_, err := toml.DecodeFile(path, &result)
	if err != nil {
		return nil, err
	}
	for _, f := range result.Server {
		f.remote = result.Remote
		f.isServer = true
	}
	for _, f := range result.Test {
		f.remote = result.Remote
		f.isMain = true

	}
	result.module = make(map[string]*wasmtime.Module)
	result.wasmFile = make(map[string]*WasmFile)
	if len(result.Server) > maxServer {
		return nil, fmt.Errorf("too many wasm modules as servers, limit on on servers is %d", maxServer)
	}
	return &result, nil
}

func (c *DeployConfig) loadSingleModule(engine *wasmtime.Engine, file *WasmFile) (*wasmtime.Module, error) {
	path := strings.TrimSpace(file.Path)
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
	log.Printf("loading module %s (%s)", file.Name, path)

	return m, nil
}

func (c *DeployConfig) loadSequenceModule(engine *wasmtime.Engine, p []*WasmFile) error {
	for _, f := range p {
		m, err := c.loadSingleModule(engine, f)
		if err != nil {
			return err
		}
		_, ok := c.wasmFile[f.Name]
		if ok {
			return fmt.Errorf("duplicate names found in deployment configuration (%s)", f.Name)
		}
		c.module[f.Name] = m
		// we don't build NameToPath until we are sure it loaded ok
		c.wasmFile[f.Name] = f
	}
	return nil
}
func (c *DeployConfig) LoadAllModules(engine *wasmtime.Engine) error {
	c.module = make(map[string]*wasmtime.Module)
	if err := c.loadSequenceModule(engine, c.Server); err != nil {
		return err
	}
	if c.TestMode {
		if err := c.loadSequenceModule(engine, c.Test); err != nil {
			return err
		}

	}
	return nil

}

func (c *DeployConfig) AllName() []string {
	result := []string{}
	for n := range c.wasmFile {
		result = append(result, n)
	}
	return result
}
