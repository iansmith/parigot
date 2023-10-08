package pdep

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

const (
	BaseImageName   = "docker.io/iansmith/parigot-koyeb-0.3"
	MaxTarEntrySize = 1024 * 1024 * 256
	Verbose         = false
)

var (
	fCodeDirectory string
	fTomlName      string
	fImageName     string
	fParigotRoot   string
	fUser          string
	fRepo          string
)

//go:embed cmd/pdep/Dockerfile.template
var Dockerfile []byte

//go:embed cmd/pdep/Caddyfile.template
var Caddyfile []byte

//go:embed cmd/pdep/startup.template
var Startup []byte

func Main() {
	flag.StringVar(&fTomlName, "c", "", "filename of the configuration file (toml format) to use")
	flag.StringVar(&fParigotRoot, "p", ".", "directory name of the root of the parigot library")
	flag.StringVar(&fImageName, "t", "", "docker tag to associate with the resulting image")
	flag.StringVar(&fUser, "u", "", "docker user name on the docker repository")
	flag.StringVar(&fRepo, "r", "docker.io", "docker repo name")

	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalf("required: directory where the code and configuration for this app can be found")
		flag.Usage()
	}
	fCodeDirectory = flag.Arg(0)
	validateArgs()

	contentDir := buildContentTarball(fCodeDirectory, fTomlName, fParigotRoot)
	log.Printf("pdep:built content tarball in %s", contentDir)
	tar, err := archive.TarWithOptions(contentDir+"/", &archive.TarOptions{})
	if err != nil {
		log.Fatalf("unable to create tar file: %v", err)
	}
	imageName := strings.Join([]string{fRepo, fUser, fImageName}, "/")
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		//Remove:    true,
		Version:  types.BuilderBuildKit,
		Platform: "linux/amd64",
		//BuildArgs: buildArg,
		Outputs: []types.ImageBuildOutput{types.ImageBuildOutput{
			Type: "registry",
			Attrs: map[string]string{
				"name": imageName,
				//"load": "true",
				"push": "true",
			},
		},
		},
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("unable to create docker client: %v", err)
	}

	res, err := cli.ImageBuild(context.Background(),
		tar, opts)
	if err != nil {
		log.Fatalf("unable to build image: %v", err)
	}
	defer res.Body.Close()

	builtId, err := print(res.Body)
	if err != nil {
		log.Fatalf("unable to build image (from docker): %s", err)
	}
	log.Printf("pdep:tagged '%s' with '%s'", builtId, fImageName)

}

func validateArgs() {
	fCodeDirectory = flag.Arg(0)
	// convert tomlName to a full path
	fTomlName = testCodeDir(fCodeDirectory, fTomlName)
	if fParigotRoot == "" {
		fParigotRoot = "."
	}
	stat, err := os.Stat(fParigotRoot)
	if err != nil {
		log.Fatalf("unable to find parigot root '%s'", fParigotRoot)
	}
	if !stat.IsDir() {
		log.Fatalf("parigot root '%s' is not a directory", fParigotRoot)
	}
	if _, err := os.Stat(filepath.Join(fParigotRoot, "build", "runner")); err != nil {
		log.Fatalf("parigot root '%s' does not appear to a parigot source installation directory (unable to find 'build/runner')", fParigotRoot)
	}
	if fImageName == "" {
		fImageName = filepath.Base(fTomlName)
	}
	if fUser == "" {
		u := os.Getenv("USER")
		if u == "" {
			log.Fatalf("must have a docker repo user, either with USER or -u flag")
		}
		log.Printf("using '%s' as docker repository username")
		fUser = u
	}
}

func testCodeDir(dir, configPath string) string {
	stat, err := os.Stat(dir)
	if err != nil {
		log.Fatalf("code directory '%s' given to deploy does not exist", dir)
	}
	if !stat.IsDir() {
		log.Fatalf("code directory '%s' given to deploy exists, but is not a directory", dir)
	}

	if configPath == "" {
		configPath = filepath.Base(dir) + ".toml"
		log.Printf("pdep:using '%s' as the configuration file", configPath)
	}

	toml := filepath.Join(dir, configPath)
	stat, err = os.Stat(toml)
	if err != nil {
		log.Fatalf("configuration file '%s' cannot be found", toml)
	}
	if stat.IsDir() || !stat.Mode().IsRegular() {
		log.Fatalf("configuration file '%s' is not a simple file", toml)
	}
	return toml
}

// by the time this is called, the arguments have been validated
func buildContentTarball(code, toml, root string) string {
	tmpdir := tempDir()
	err := os.MkdirAll(filepath.Join(tmpdir, "app", "build"), 0777)
	if err != nil {
		log.Fatalf("unable to create directories for building docker image: %v", err)
	}

	// copy burned in Dockerfile
	buf := bytes.NewBuffer(Dockerfile)
	copyFileFromReader(tmpdir, "Dockerfile", buf)

	// copy toml file from disk
	buf = readFileToBuffer(toml)
	copyFileFromReader(tmpdir, filepath.Join("app", "app.toml"), buf)

	// copy burned-in CaddyFile
	buf = bytes.NewBuffer(Caddyfile)
	copyFileFromReader(tmpdir, filepath.Join("app", "Caddyfile"), buf)

	// copy burned-in script
	buf = bytes.NewBuffer(Startup)
	copyFileFromReader(tmpdir, filepath.Join("app", "startup.sh"), buf)

	// dfile, err := os.Create(filepath.Join(tmpdir, "Dockerfile"))
	// if err != nil {
	// 	log.Fatalf("error creating new dockerfile:%v", err)
	// }
	// dockerIn := bytes.NewBuffer(Dockerfile)
	// if _, err := io.Copy(dfile, dockerIn); err != nil {
	// 	log.Fatalf("error copying tmpl file:%v", err)
	// }
	// dfile.Close()

	addEntries := func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if strings.HasSuffix(path, ".p.wasm") {
			base := filepath.Join(tmpdir, "app", "build", filepath.Base(path))
			out, err := os.Create(base)
			if err != nil {
				log.Fatalf("cannot create file in docker image: %v", err)
			}
			in, err := os.Open(path)
			if err != nil {
				log.Fatalf("cannot read file building in docker image: %v", err)
			}
			if _, err := io.Copy(out, in); err != nil {
				log.Fatalf("cannot copy file building in docker image: %v", err)
			}
			out.Close()
			in.Close()
		}
		return nil
	}
	if err := filepath.WalkDir(filepath.Join(code, "build"),
		addEntries); err != nil {
		log.Fatalf("unable to read files into tarball; %v", err)
	}
	if err := filepath.WalkDir(filepath.Join(root, "build"),
		addEntries); err != nil {
		log.Fatalf("unable to read files into tarball: %v", err)
	}
	return tmpdir
}

func tempDir() string {
	d := os.TempDir()
	if _, extraErr := os.Stat("tmp"); extraErr == nil {
		d = "tmp"
	}
	tmpdir, err := os.MkdirTemp(d, "pdep")
	if err := os.Chmod(tmpdir, os.FileMode(0755)); err != nil {
		log.Fatalf("unable to chmod temporary directory '%s':%v", tmpdir)
	}
	if err != nil {
		log.Fatalf("unable to create  :%v", err)
	}
	return tmpdir
}

// this is only used by the v1 builder
type BuiltId struct {
	Id string `json:"ID"`
}

type OutLine struct {
	// Used only by v1 builder
	Stream string `json:"stream"`
	// different setup for the v1 builder
	//Aux BuiltId
	Aux interface{} `json:"aux"`
	// used by v2, but is always moby.image.id
	Id string `json:"id"`
}
type ErrorDetail struct {
	Message string `json:"message"`
}

type ErrorLine struct {
	ErrorBase   string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

func (e *ErrorLine) Error() string {
	return fmt.Sprintf("%s [%s]", e.ErrorBase, e.ErrorDetail.Message)
}

func print(rd io.Reader) (string, error) {
	var lastLine string
	var finalId string
	scanner := bufio.NewScanner(rd)
	buf := &bytes.Buffer{}
	for scanner.Scan() {
		lastLine = scanner.Text()
		if err := scanner.Err(); err != nil {
			return "", err
		}
		dec := json.NewDecoder(strings.NewReader(lastLine))
		result := &OutLine{}
		if err := dec.Decode(result); err != nil {
			log.Fatalf("unable to decode line received from docker '%s', error %s", lastLine, err)
		}
		errLine := &ErrorLine{}
		json.Unmarshal([]byte(lastLine), errLine)
		if errLine.ErrorBase != "" {
			return "", errLine
		}
		var ok bool
		if result.Aux != nil {
			var s string
			s, ok = result.Aux.(string)
			if ok {
				buf.WriteString(s)
				continue
			}
			var m map[string]interface{}
			m, ok = result.Aux.(map[string]interface{})
			if ok {
				log.Printf("result of build: %s", m["ID"])
				continue
			}
			log.Printf("unable to understand json result in result from docker: %+v", result.Aux)
		}
	}
	return finalId, nil
}

func copyFileFromReader(tmpdir, name string, rd io.Reader) {
	dfile, err := os.Create(filepath.Join(tmpdir, name))
	if err != nil {
		log.Fatalf("error creating new file '%s':%v", name, err)
	}
	_, err = io.Copy(dfile, rd)
	if err != nil {
		log.Fatalf("error copying tmpl file:%v", err)
	}
	dfile.Close()

}

func readFileToBuffer(path string) *bytes.Buffer {
	fp, err := os.Open(path)
	if err != nil {
		log.Fatalf("unable to open '%s' for reading: %v", path, err)
	}
	rd := io.LimitReader(fp, MaxTarEntrySize)
	data, err := io.ReadAll(rd)
	if err != nil {
		log.Printf("failed to read all of '%s': %v", path, err)
	}

	return bytes.NewBuffer(data)
}
