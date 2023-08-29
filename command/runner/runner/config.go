package runner

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"plugin"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/eng"
)

var deployVerbose = false || os.Getenv("PARIGOT_VERBOSE") != ""

const allowDeploySize = false

const envVar = "PARIGOT_DEPLOYMENT"

var chosen = os.Getenv(envVar) // should be treated as a constant

// Deployment is a mapping from names to DeployConfigs.
type Deployment struct {
	Config map[string]*DeployConfig
}

// Timeout settings are to control timeouts within the system.  All
// values in milliseconds.
type Timeout struct {
	// Startup is the amount of time to wait for services to launch and
	// have all previously launched services start.  In other words, how
	// long to wait to see if we can get everything in the deployment configuration
	// running.
	Startup int
	// Complete is the length of time before a partially completed call will
	// be considered failed. If you call method foo() on service bar, this is
	// how long we will wait for bar to fulfull the call of foo() and provide
	// the return value.
	Complete int
}

// DeployConfig represents the microservices that the user has configured for this application.
// Public fields in this struct are data that has been read from the user and has been
// sanity checked.
type DeployConfig struct {
	Microservice     map[string]*Microservice
	Flag             *DeployFlag
	ParigotLibPath   string
	ParigotLibSymbol string
	Arrangement      DeployArrangement
	ArrangementName  string
	Size             DeploySize
	SizeName         string
	Timezone         string
	TimezoneDir      string
	Timeout          Timeout
}

type DeploySize int
type DeployArrangement int

const (
	SizeNotSpecified DeploySize = 0
	SizeExtraSmall              = 1
	SizeSmall                   = 2
	SizeMedium                  = 3
	SizeLarge                   = 4
	SizeExtraLarge              = 5
	SizeLast                    = SizeExtraLarge
)
const (
	ArrangeNotSpecified       DeployArrangement = 0
	LocalProcess                                = 1
	LocalDocker                                 = 2
	ArragementRemoteMarkerBit                   = 0x10000
	RemoteProcess                               = ArragementRemoteMarkerBit | 1
	RemoteDocker                                = ArragementRemoteMarkerBit | 2
	ArrangementLocalLast                        = LocalDocker
	ArrangementRemoteLast                       = RemoteDocker
)

// DeployFlag is a structure that comes from the command line passed to the runner itself.  These
// switches have a large effect on how the runner behaves.
type DeployFlag struct {
	// TestMode being true means that microservices that are marked as "Test" will be considered the program(s)
	// to run.  If TestMode is false, programs marked as Main will be run.  Note that microservices configured
	// to be Test are ignored when TestMode is false, and vice versa with microservices marked as Main.
	TestMode bool
}

// Microservice is the unit of configuration for the DeployConfig. Public fields are data read from the user's
// configuration and these are checked for sanity before being returned.  Exactly one of Server, Main, and
// Test must be set to true.
type Microservice struct {
	WasmPath     string
	PluginPath   string
	PluginSymbol string

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
	var deployment Deployment
	md, err := toml.DecodeFile(path, &deployment)
	if err != nil {
		return nil, err
	}
	for i, j := range md.Undecoded() {
		log.Printf("undecoded %d,%+v", i, j.String())
	}

	if chosen == "" {
		chosen = "dev"
	}
	result, ok := deployment.Config[chosen]
	if !ok {
		return nil, fmt.Errorf("unable to find deployment '%s' in deployment descriptor '%s", chosen, path)
	}
	// copied from user input
	result.Flag = flag
	// loop over the configs making text -> int conversions
	for _, dc := range deployment.Config {
		var err error
		dc.Size, err = sizeToDeploySize(chosen, dc.SizeName, dc.Size)
		if err != nil {
			return nil, err
		}
		dc.Arrangement, err = arrangementToDeployArrangement(chosen, dc.ArrangementName, dc.Arrangement)
		if err != nil {
			return nil, err
		}
	}
	for name, m := range result.Microservice {
		// these are just copied to the microservice for convience
		m.name = name

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
		err := pathExists(m.name, m.WasmPath, false)
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
			err := pathExists(m.name, m.PluginPath, true)
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
	return result, nil
}

type simpleEnv struct {
	argv []string
	envp map[string]string
	hid  id.HostId
}

func (s *simpleEnv) Arg() []string                  { return s.argv }
func (s *simpleEnv) Environment() map[string]string { return s.envp }
func (s *simpleEnv) Host() id.HostId                { return s.hid }

func (c *DeployConfig) LoadSingleModule(engine eng.Engine, name string) (eng.Module, error) {
	m, ok := c.Microservice[name]
	if !ok {
		panic(fmt.Sprintf("unable to find microservice with name '%s", name))
	}
	argv := []string{name}
	argv = append(argv, m.GetArg()...)
	envp := make(map[string]string)
	raw := m.GetEnv()
	for _, s := range raw {
		part := strings.SplitN(s, "=", 2)
		if len(part) == 0 {
			continue
		}
		if len(part) == 1 {
			envp[part[0]] = ""
		} else {
			envp[part[0]] = strings.Join(part[1:], "")
		}
	}
	env := simpleEnv{
		argv: argv,
		envp: envp,
		hid:  id.NewHostId(),
	}
	mod, err := engine.NewModuleFromFile(context.Background(), m.WasmPath, &env)
	if err != nil {
		log.Printf("new module failed to create from file %s: %v", m.WasmPath, err.Error())
		return nil, fmt.Errorf("unable to load microservice (%s): cannot convert %s into a module: %v",
			m.name, m.WasmPath, err)
	}
	deployPrint("loadSingleModule", "loading module %s (with wasm code: %s)", m.name, m.WasmPath)
	return mod, nil
}

func deployPrint(spec string, rest ...interface{}) {
	if deployVerbose {
		log.Printf(spec, rest...)
	}
}

func (c *DeployConfig) LoadAllModules(engine eng.Engine) error {
	for n, m := range c.Microservice {
		mod, err := c.LoadSingleModule(engine, n)
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
func (m *Microservice) GetPluginSymbol() string {
	return m.PluginSymbol
}
func (m *Microservice) GetPlugin() *plugin.Plugin {
	return m.plug
}
func (m *Microservice) GetModule() eng.Module {
	return m.module
}

func pathExists(serviceName, path string, isPlugin bool) error {
	pathType := "wasm path"
	if isPlugin {
		pathType = "plugin path"
	}
	var info fs.FileInfo
	var err error
	if isPlugin {
		info, err = pathExistsPlugin(path)
	} else {
		info, err = os.Stat(path)
	}
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("bad microservice configuration (%s): %s '%s' does not exist",
				serviceName, pathType, path)
		}
		return fmt.Errorf("bad microservice configuration (%s): %s '%s': %v",
			serviceName, pathType, path, err)
	}
	// in the simple case where we link the plugins it returns nil
	// as the file info.
	if info != nil && info.IsDir() {
		return fmt.Errorf("bad microservice configuration (%s): %s '%s' cannot be a directory",
			serviceName, pathType, path)

	}
	return nil
}

func arrangementToDeployArrangement(name string, s string, da DeployArrangement) (DeployArrangement, error) {
	daInRange := func(da DeployArrangement) bool {
		if da&ArragementRemoteMarkerBit != 0 {
			mask := ArragementRemoteMarkerBit - 1
			da := int(da) & mask
			if da <= int(ArrangeNotSpecified) || da > ArrangementRemoteLast {
				return false
			}
			return true
		}
		if da <= ArrangeNotSpecified || da > ArrangementLocalLast {
			return false
		}
		return true
	}
	isDev := strings.ToLower(name) == "dev"
	if s == "" && da == ArrangeNotSpecified {
		if isDev {
			return LocalProcess, nil
		} else {
			return ArrangeNotSpecified, fmt.Errorf("exactly one of Arrangement or ArrangementName must be specified, neither found in deployment '%s'", name)
		}
	}
	if s != "" && daInRange(da) {
		return ArrangeNotSpecified, fmt.Errorf("exactly one of Arrangement or ArrangementName must be specified, both found in deployment '%s'", name)

	}
	if s != "" {
		switch strings.ToLower(s) {
		case "localprocess":
			return LocalProcess, nil
		case "localdocker":
			return LocalDocker, nil
		case "remoteprocess":
			return RemoteProcess, nil
		case "remotedocker":
			return RemoteDocker, nil
		}
		return ArrangeNotSpecified, fmt.Errorf("arrangement '%s' is not knwon", s)
	}
	if !daInRange(da) {
		return ArrangeNotSpecified, fmt.Errorf("arrangement number %d not known, valid values are from %d to %d and %d to %d",
			da, ArrangeNotSpecified+1, ArrangementLocalLast, (ArrangeNotSpecified|ArragementRemoteMarkerBit)+1, (ArragementRemoteMarkerBit)|ArrangementRemoteLast)
	}
	return da, nil // it's ok
}

func sizeToDeploySize(name string, s string, ds DeploySize) (DeploySize, error) {
	isDev := strings.ToLower(name) == "dev"
	if isDev {
		return SizeNotSpecified, nil
	}
	if !allowDeploySize && (s != "" || ds != SizeNotSpecified) {
		return SizeNotSpecified, fmt.Errorf("deployment size value is not permitted for this runner")
	}
	if s != "" && ds != SizeNotSpecified {
		return SizeNotSpecified, fmt.Errorf("exactly one of Size and SizeName must be specified, neither found")
	}
	if s == "" && ds == SizeNotSpecified {
		return SizeNotSpecified, fmt.Errorf("exactly one of Size and SizeName must be specified, both found")
	}
	if s != "" {
		switch strings.ToLower(s) {
		case "extrasmall":
			return SizeExtraSmall, nil
		case "small":
			return SizeSmall, nil
		case "medium":
			return SizeMedium, nil
		case "large":
			return SizeLarge, nil
		case "extralarge":
			return SizeExtraLarge, nil
		}
		return SizeNotSpecified, fmt.Errorf("unknown value for the SizeName field '%s'", s)
	}
	if ds <= SizeNotSpecified || ds > SizeLast {
		return SizeNotSpecified, fmt.Errorf("unknown value Size field %d, values from %d to %d are accepted",
			ds, SizeNotSpecified+1, SizeLast)
	}
	return ds, nil
}
