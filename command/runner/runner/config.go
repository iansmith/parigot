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

// DeployConfig represents the microservices that the user has configured for this application.
// Public fields in this struct are data that has been read from the user and has been
// sanity checked.
type DeployConfig struct {
	Microservice map[string]*Microservice
	Flag         *DeployFlag
}

// DeployFlag is a structure that comes from the command line passed to the runner itself.  These
// switches have a large effect on how the runner behaves.
type DeployFlag struct {
	// TestMode being true means that microservices that are marked as "Test" will be considered the program(s)
	// to run.  If TestMode is false, programs marked as Main will be run.  Note that microservices configured
	// to be Test are ignored when TestMode is false, and vice versa with microservices marked as Main.
	TestMode bool
	// Remote being true means that the microservices should be run in separate address spaces.  If this flag
	// is false, all the microservices are run in a single process (locally).
	Remote bool
}

// Microservice is the unit of configuration for the DeployConfig. Public fields are data read from the user's
// configuration and these are checked for sanity before being returned.  Exactly one of Server, Main, and
// Test must be set to true.
type Microservice struct {
	Path   string
	Arg    []string
	Env    []string
	Server bool
	Main   bool
	Test   bool

	// stuff we add
	name   string
	remote bool
	module *wasmtime.Module
}

func (m *Microservice) Name() string {
	return m.name
}

func (m *Microservice) Module() *wasmtime.Module {
	return m.module
}

const maxServer = 32

func Parse(path string, flag *DeployFlag) (*DeployConfig, error) {
	var result DeployConfig
	_, err := toml.DecodeFile(path, &result)
	if err != nil {
		return nil, err
	}
	result.Flag = flag
	for name, m := range result.Microservice {
		// these are just copied to the microservice for convience
		m.name = name
		m.remote = flag.Remote

		// get rid of spaces at the end and start of strings
		m.Path = strings.TrimSpace(m.Path)
		arg := make([]string, len(m.Arg))
		for i, a := range m.Arg {
			arg[i] = strings.TrimSpace(a)
		}
		m.Arg = arg
		env := make([]string, len(m.Env))
		for i, e := range m.Env {
			env[i] = strings.TrimSpace(e)
		}
		m.Env = env
		// the type of microservice
		switch {
		case !m.Server && !m.Main && !m.Test:
			m.Server = true // default
		case !m.Server && !m.Main && m.Test:
		case !m.Server && m.Main && !m.Test:
		default:
			return nil, fmt.Errorf("bad microservice configuration (%s): one of Server(%v), Test(%v), or Main(%v) must be true, or all must be false which defaults to Server=true",
				name, m.Server, m.Test, m.Main)
		}
		// path sanity check
		if m.Path == "" {
			return nil, fmt.Errorf("bad microservice configuration (%s): Path is a required field", name)
		}
		info, err := os.Stat(m.Path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("bad microservice configuration (%s): Path (%s) does not exist",
					name, m.Path)
			}
			return nil, fmt.Errorf("bad microservice configuration (%s): Path (%s): %v",
				name, m.Path, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("bad microservice configuration (%s): Path (%s) cannot be a directory",
				name, m.Path)
		}
		// sanity check env vars
		for _, envvar := range m.Env {
			index := strings.Index(envvar, "=")
			if index == -1 {
				return nil, fmt.Errorf("bad microservice configuration (%s):'%s' is not an environment variable of the form 'FOO=bar'",
					m.name, envvar)
			}
		}
		//make sure they don't have weird things in the strings we are passing to the code
		for _, s := range append(m.Arg, m.Env...) {
			for _, c := range s {
				if c < 32 || c > 126 {
					return nil, fmt.Errorf("bad microservice configuration (%s):'%s' contains non-ascii characters",
						m.name, s)
				}
			}
			// also check for empty
			if s == "" {
				return nil, fmt.Errorf("bad microservice configuration (%s): empty strings are not allowed for arguments or environment variables",
					m.name)
			}
		}
	}
	if len(result.Microservice) == 0 {
		return nil, fmt.Errorf("no microservices found in configuration %s", path)
	}
	serverCount := 0
	for _, m := range result.Microservice {
		if m.Server {
			serverCount++
		}
	}
	if serverCount >= maxServer {
		return nil, fmt.Errorf("too many server microservices found in configuration %s, limit on servers is %d", path, maxServer)
	}
	return &result, nil
}

func (c *DeployConfig) loadSingleModule(engine *wasmtime.Engine, m *Microservice) (*wasmtime.Module, error) {
	mod, err := wasmtime.NewModuleFromFile(engine, m.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to load microservice (%s): cannot convert %s into a module: %v",
			m.name, m.Path, err)
	}
	log.Printf("loading module %s (%s)", m.name, m.Path)
	return mod, nil
}

func (c *DeployConfig) LoadAllModules(engine *wasmtime.Engine) error {
	for _, m := range c.Microservice {
		mod, err := c.loadSingleModule(engine, m)
		if err != nil {
			return err
		}
		m.module = mod
	}
	return nil
}

func (c *DeployConfig) AllName() []string {
	result := []string{}
	for n := range c.Microservice {
		result = append(result, n)
	}
	return result
}

func (c *DeployConfig) Module(name string) *wasmtime.Module {
	m, ok := c.Microservice[name]
	if !ok {
		return nil
	}
	return m.module
}

func (m *Microservice) IsServer() bool {
	return m.Server
}
func (m *Microservice) IsLocal() bool {
	return !m.remote
}
func (m *Microservice) IsRemote() bool {
	return m.remote
}
func (m *Microservice) GetName() string {
	return m.name
}

func (m *Microservice) GetEnv() []string {
	return m.Env
}

func (m *Microservice) GetArg() []string {
	return m.Arg
}
func (m *Microservice) GetPath() string {
	return m.Path
}
func (m *Microservice) GetModule() *wasmtime.Module {
	return m.module
}
