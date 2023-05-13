package runner

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"plugin"
	"strings"

	"github.com/BurntSushi/toml"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
)

var deployVerbose = false || os.Getenv("PARIGOT_VERBOSE") != ""

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
	WasmPath   string
	PluginPath string

	Arg []string
	Env []string

	Server bool
	Main   bool
	Test   bool

	// stuff we add
	name   string
	remote bool
	module eng.Module

	plug *plugin.Plugin
}

func (m *Microservice) Name() string {
	return m.name
}

func (m *Microservice) Module() eng.Module {
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
		m.WasmPath = strings.TrimSpace(m.WasmPath)
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
		if m.WasmPath == "" {
			return nil, fmt.Errorf("bad microservice configuration (%s): Path is a required field", name)
		}
		_, err := pathExists(m.name, m.WasmPath, false)
		if err != nil {
			return nil, err
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
		// load plugin if necessary
		if !m.Server && m.PluginPath != "" {
			return nil, fmt.Errorf("bad microservice configuration (%s): PluginPath is only allowed for microservices that are servers", m.name)
		}
		if m.Server && m.PluginPath != "" {
			_, err := pathExists(m.name, m.PluginPath, true)
			if err != nil {
				return nil, err
			}
			m.plug, err = plugin.Open(m.PluginPath)
			if err != nil {
				return nil, err
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

func (c *DeployConfig) loadSingleModule(ctx context.Context, engine eng.Engine, m *Microservice) (eng.Module, error) {
	mod, err := engine.NewModuleFromFile(ctx, m.WasmPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load microservice (%s): cannot convert %s into a module: %v",
			m.name, m.WasmPath, err)
	}
	deployPrint(ctx, pcontext.Debug, "loadSingleModule", "loading module %s (%s)", m.name, m.WasmPath)
	return mod, nil
}

func deployPrint(ctx context.Context, ll pcontext.LogLevel, method string, spec string, rest ...interface{}) {
	if deployVerbose {
		pcontext.LogFullf(ctx, ll, pcontext.UnknownS, method, spec, rest...)
	}
}

func (c *DeployConfig) LoadAllModules(ctx context.Context, engine eng.Engine) error {
	for _, m := range c.Microservice {
		mod, err := c.loadSingleModule(ctx, engine, m)
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

func (c *DeployConfig) Module(name string) eng.Module {
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
func (m *Microservice) GetWasmPath() string {
	return m.WasmPath
}
func (m *Microservice) GetPluginPath() string {
	return m.PluginPath
}
func (m *Microservice) GetPlugin() *plugin.Plugin {
	return m.plug
}
func (m *Microservice) GetModule() eng.Module {
	return m.module
}

func pathExists(serviceName, path string, isPlugin bool) (fs.FileInfo, error) {
	pathType := "wasm path"
	if isPlugin {
		pathType = "plugin path"
	}
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("bad microservice configuration (%s): %s '%s' does not exist",
				serviceName, pathType, path)
		}
		return nil, fmt.Errorf("bad microservice configuration (%s): %s '%s': %v",
			serviceName, pathType, path, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("bad microservice configuration (%s): %s '%s' cannot be a directory",
			serviceName, pathType, path)

	}
	return info, nil
}
