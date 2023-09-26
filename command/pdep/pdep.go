package pdep

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
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
)

var (
	fCodeDirectory string
	fTomlName      string
	fImageName     string
	fParigotRoot   string
)

//go:embed cmd/pdep/Dockerfile.template
var Dockerfile []byte

func Main() {
	flag.StringVar(&fTomlName, "c", "", "filename of the configuration file (toml format) to use")
	flag.StringVar(&fParigotRoot, "p", ".", "directory name of the root of the parigot library")
	flag.StringVar(&fImageName, "t", "", "docker tag to associate with the resulting image")

	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalf("required: directory where the code and configuration for this app can be found")
		flag.Usage()
	}
	fCodeDirectory = flag.Arg(0)
	validateArgs()

	contentDir := buildContentTarball(fCodeDirectory, fTomlName, fParigotRoot)
	log.Printf("content tarball: %s", contentDir)
	tar, err := archive.TarWithOptions(contentDir+"/", &archive.TarOptions{})
	if err != nil {
		log.Fatalf("unable to create tar file: %v", err)
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{fImageName},
		Remove:     true,
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

	if err := print(res.Body); err != nil {
		log.Fatalf("unable to print the response from docker: %v", err)
	}

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
		log.Printf("using '%s' as the configuration file", configPath)
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

type tarFileEntry struct {
	Name   string
	Buffer string
}

// by the time this is called, the arguments have been validated
func buildContentTarball(code, toml, root string) string {
	tmpdir := tempDir()
	log.Printf("tmp dir %s", tmpdir)
	err := os.MkdirAll(filepath.Join(tmpdir, "app", "build"), 0777)
	if err != nil {
		log.Fatalf("unable to create directories for building docker image: %v", err)
	}

	// copy toml file
	tomlIn, err := os.Open(toml)
	if err != nil {
		log.Fatalf("error opening '%s':%v", toml, err)
	}
	tomlOut, err := os.Create(filepath.Join(tmpdir, "app", "app.toml"))
	if err != nil {
		log.Fatalf("error creating tmp file:%v", err)
	}
	if _, err := io.Copy(tomlOut, tomlIn); err != nil {
		log.Fatalf("error copying tmpl file:%v", err)
	}
	tomlIn.Close()
	tomlOut.Close()

	dfile, err := os.Create(filepath.Join(tmpdir, "Dockerfile"))
	if err != nil {
		log.Fatalf("error creating new dockerfile:%v", err)
	}
	dockerIn := bytes.NewBuffer(Dockerfile)
	if _, err := io.Copy(dfile, dockerIn); err != nil {
		log.Fatalf("error copying tmpl file:%v", err)
	}
	dfile.Close()

	addEntries := func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if strings.HasSuffix(path, ".p.wasm") {
			base := filepath.Join(tmpdir, "app", "build", filepath.Base(path))
			log.Printf("base %s", base)
			out, err := os.Create(base)
			if err != nil {
				log.Fatalf("cannot create file in docker image: %v", err)
			}
			log.Printf("path %s", path)
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

// func readFileToBuffer(path string) *bytes.Buffer {
// 	fp, err := os.Open(path)
// 	if err != nil {
// 		log.Fatalf("unable to open '%s' for reading: %v", path, err)
// 	}
// 	rd := io.LimitReader(fp, MaxTarEntrySize)
// 	data, err := io.ReadAll(rd)
// 	if err != nil {
// 		log.Printf("failed to read all of '%s': %v", path, err)
// 	}
// 	return bytes.NewBuffer(data)
// }

// func writeFileToTar(tarW *tar.Writer, readPath, targetPath string) {
// 	// just read it in before calling writeReaderToTar
// 	buf := readFileToBuffer(readPath)
// 	writeReaderToTar(tarW, targetPath, buf, int64(buf.Len()))
// }

// func writeReaderToTar(tarW *tar.Writer, path string, rd io.Reader, size int64) {
// 	hdr := &tar.Header{}
// 	hdr.Name = path
// 	hdr.Size = size
// 	hdr.Mode = 0644
// 	hdr.ModTime = time.Now()
// 	hdr.Uid = 0
// 	hdr.Uname = "root"
// 	hdr.Gid = 0
// 	hdr.Gname = "root"

// 	tarW.WriteHeader(hdr)
// 	if _, err := io.Copy(tarW, rd); err != nil {
// 		log.Fatalf("error copy tar file content: %v", err)
// 	}
// }

type ErrorDetail struct {
	Message string `json:"message"`
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
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

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
