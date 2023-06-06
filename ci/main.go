package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

const (
	apiVersion = "v1"

	goToWASM   = "go1.21"
	goToHost   = "go1.20.4"
	goToPlugin = "go1.20.4"

	rep = "g/file/" + apiVersion + "/file.pb.go"
)

var (
	// go environment variables
	goEnvVarsWASM = map[string]string{
		"GOROOT": "/home/parigot/deps/go1.21",
		"GOOS":   "wasip1",
		"GOARCH": "wasm",
	}
	goEnvVarsHost = map[string]string{
		"GOROOT": "/home/parigot/deps/go1.20.4",
		"GOOS":   "",
		"GOARCH": "",
	}
	goEnvVarsPlugin = map[string]string{
		"GOROOT": "/home/parigot/deps/go1.20.4",
		"GOOS":   "",
		"GOARCH": "",
	}

	// EXTRA ARGS FOR BUILDING (placed after the "go build")
	// use -x for more details from a go compiler
	extraWASMCompArgs = []string{}
	extraHostArgs     = []string{}
	extraPluginArgs   = []string{"-buildmode=plugin"}
)

func main() {
	if err := build(context.Background()); err != nil {
		log.Fatalf("Cannot build dagger pipeline %v\n", err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	// set the current(root) directory on the host as the working directory
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local Dockerfile
	dockerDir := client.Host().Directory("./ci/")

	// initialize image fron dockerfile
	// mount a host directory in the container at the '/workspaces/parigot' path
	img := client.Container().
		Build(dockerDir).
		WithDirectory(
			"/workspaces/parigot",
			client.Host().Directory("."),
			dagger.ContainerWithDirectoryOpts{Exclude: []string{".devcontainer/", "ci/", "build/", "g/"}},
		).
		WithWorkdir("/workspaces/parigot").
		WithUser("root")

	// build/protoc-gen-parigot
	img, err = buildProtocGenParigot(ctx, img)
	if err != nil {
		return err
	}

	// generate rep
	img, err = generateRep(ctx, img)
	if err != nil {
		return err
	}

	img, err = generateApiID(ctx, img)
	if err != nil {
		return err
	}

	img, err = buildPlugins(ctx, img)
	if err != nil {
		return err
	}

	img, err = buildClientSideOfAPIs(ctx, img)
	if err != nil {
		return err
	}

	img, err = buildRunner(ctx, img)
	if err != nil {
		return err
	}

	// get reference to build output directory in container
	output := img.Directory("build")
	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, "build")
	if err != nil {
		return err
	}

	return nil
}

func buildRunner(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to build runner
	 */
	exist, _ := findFilesWithSuffixRecursively("command/runner", ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory command/runner and its subdirectories")
	}
	exist, _ = findFilesWithPattern("apishared/id/*.go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory apishared/id")
	}
	exist, _ = findFilesWithPattern("apiwasm/*.go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory apiwasm")
	}
	exist, _ = findFilesWithPattern("apiplugin/*")
	if !exist {
		log.Fatalf("There is no such file: apiplugin/*")
	}
	exist, _ = findFilesWithPattern("wazero-src-1.1/fauxfd.go")
	if !exist {
		log.Fatalf("There is no such file: wazero-src-1.1/fauxfd.go")
	}

	// set up HOST env variables
	for key, value := range goEnvVarsHost {
		img = img.WithEnvVariable(key, value)
	}

	target := "build/runner"
	packagePath := "github.com/iansmith/parigot/command/runner"
	img = img.WithExec([]string{"rm", "-f", target})
	// go build
	goCmd := []string{goToHost, "build"}
	for _, arg := range extraHostArgs {
		goCmd = append(goCmd, arg)
	}
	goCmd = append(goCmd, "-o", target, packagePath)
	img = img.WithExec(goCmd)

	return img, nil
}

func buildPlugins(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to build plugins
	 */
	// set up Plugin env variables
	for key, value := range goEnvVarsPlugin {
		img = img.WithEnvVariable(key, value)
	}

	// SYS_SRC
	exist, _ := findFilesWithSuffixRecursively("sys", ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory sys and its subdirectories")
	}
	// ENG_SRC
	exist, _ = findFilesWithSuffixRecursively("eng", ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory eng and its subdirectories")
	}
	// CTX_SRC
	exist, _ = findFilesWithSuffixRecursively("context", ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory context and its subdirectories")
	}
	// SHARED_SRC
	exist, _ = findFilesWithSuffixRecursively("apishared", ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory apishared and its subdirectories")
	}
	// check apiplugin/*.go
	exist, _ = findFilesWithPattern("apiplugin/*.go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory apiplugin")
	}
	// check wazero-src-1.1/fauxfd.go
	exist, _ = findFilesWithPattern("wazero-src-1.1/fauxfd.go")
	if !exist {
		log.Fatalf("There are no such file: wazero-src-1.1/fauxfd.go")
	}

	// apiplugin/queue/db.go
	img, err := sqlcForQueue(ctx, img)
	if err != nil {
		return img, err
	}

	// build/syscall.so
	dir := "apiplugin/syscall"
	target := "build/syscall.so"
	packagePath := "github.com/iansmith/parigot/apiplugin/syscall"
	img, err = buildAPlugin(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/queue.so
	dir = "apiplugin/queue"
	target = "build/queue.so"
	packagePath = "github.com/iansmith/parigot/apiplugin/queue"
	img, err = buildAPlugin(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/file.so
	dir = "apiplugin/file"
	target = "build/file.so"
	packagePath = "github.com/iansmith/parigot/apiplugin/file"
	img, err = buildAPlugin(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}

	return img, nil
}

func buildClientSideOfAPIs(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to build some client side of APIs:
	 *		build/file.p.wasm,
	 *		build/test.p.wasm,
	 *		build/queue.p.wasm
	 */
	syscallClientSide := "apiwasm/syscall/*.go"
	exist, _ := findFilesWithPattern(syscallClientSide)
	if !exist {
		log.Fatalf("There are no such files %s", syscallClientSide)
	}

	// set up WASM env variables
	for key, value := range goEnvVarsWASM {
		img = img.WithEnvVariable(key, value)
	}

	// build/file.p.wasm
	dir := "apiwasm/file"
	target := "build/file.p.wasm"
	packagePath := "github.com/iansmith/parigot/apiwasm/file"
	img, err := buildAClientService(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/test.p.wasm
	dir = "apiwasm/test"
	target = "build/test.p.wasm"
	packagePath = "github.com/iansmith/parigot/apiwasm/test"
	img, err = buildAClientService(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/queue.p.wasm
	dir = "apiwasm/queue"
	target = "build/queue.p.wasm"
	packagePath = "github.com/iansmith/parigot/apiwasm/queue"
	img, err = buildAClientService(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}

	return img, nil
}

func buildProtocGenParigot(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to build protoc-gen-parigot
	 */
	dir := "command/protoc-gen-parigot"
	exist, _ := findFilesWithSuffixRecursively(dir, ".tmpl")
	if !exist {
		log.Fatalf("There are no such files *.tmpl in the directory %s and its subdirectories", dir)
	}
	exist, _ = findFilesWithSuffixRecursively(dir, ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory %s and its subdirectories", dir)
	}

	// set up HOST env variables
	for key, value := range goEnvVarsHost {
		img = img.WithEnvVariable(key, value)
	}

	target := "build/protoc-gen-parigot"
	packagePath := "github.com/iansmith/parigot/command/protoc-gen-parigot"
	img = img.WithExec([]string{"rm", "-f", target}).
		WithExec([]string{goToHost, "build", "-o", target, packagePath})

	return img, nil
}

func generateRep(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	generate a single representative file for all the protobuf generated code
	 */
	exist, _ := findFilesWithSuffixRecursively("api/proto", ".proto")
	if !exist {
		log.Fatalf("There are no such files *.proto in the directory api/proto and its subdirectories")
	}
	exist, _ = findFilesWithSuffixRecursively("test", ".proto")
	if !exist {
		log.Fatalf("There are no such files *.proto in the directory test and its subdirectories")
	}

	img = img.WithExec([]string{"rm", "-rf", "g/*"}).
		WithExec([]string{"buf", "lint"}).
		WithExec([]string{"buf", "generate"})

	// export g dir to host
	output := img.Directory("g")
	_, err := output.Export(ctx, "g")
	if err != nil {
		return img, err
	}

	return img, nil
}

func generateApiID(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to generate some id cruft for a couple of types built by parigot:
	 *		apishared/id/kernelerrid.go,	apiwasm/bytepipeid.go
	 *		apishared/id/serviceid.go,		apishared/id/methodid.go
	 *		apishared/id/callid.go,			g/queue/v1/queueid.go,
	 *		g/queue/v1/rowid.go,			g/queue/v1/queuemsgid.go,
	 *		g/file/v1/fileid.go,			g/test/v1/testid.go,
	 *		g/methodcall/v1/methodcallid.go
	 */
	apiID := "apishared/id/id.go"
	boilerplateid := "command/boilerplateid/main.go"
	goCmd := []string{goToHost, "run", boilerplateid}

	exist, _ := findFilesWithPattern(apiID)
	if !exist {
		log.Fatalf("There are no such file: %s", apiID)
	}
	exist, _ = findFilesWithPattern(boilerplateid)
	if !exist {
		log.Fatalf("There are no such file: %s", boilerplateid)
	}
	exist, _ = findFilesWithPattern("command/boilerplateid/template/*.tmpl")
	if !exist {
		log.Fatalf("There are no such files *.tmpl in the directory command/boilerplateid/template")
	}

	target := "apishared/id/kernelerrid.go"
	kernelCmd := append(goCmd, "-i", "-e", "id", "KernelErr", "k", "errkern")
	img, err := generateIdFile(ctx, img, target, kernelCmd)
	if err != nil {
		return img, err
	}
	target = "apiwasm/bytepipeid.go"
	bytepipCmd := append(goCmd, "-e", "apiwasm", "BytePipeErr", "b", "errbytep")
	img, err = generateIdFile(ctx, img, target, bytepipCmd)
	if err != nil {
		return img, err
	}
	target = "apishared/id/serviceid.go"
	serviceCmd := append(goCmd, "-i", "-p", "id", "Service", "s", "svc")
	img, err = generateIdFile(ctx, img, target, serviceCmd)
	if err != nil {
		return img, err
	}
	target = "apishared/id/methodid.go"
	methodCmd := append(goCmd, "-i", "-p", "id", "Method", "m", "method")
	img, err = generateIdFile(ctx, img, target, methodCmd)
	if err != nil {
		return img, err
	}
	target = "apishared/id/callid.go"
	callCmd := append(goCmd, "-i", "-p", "id", "Call", "c", "call")
	img, err = generateIdFile(ctx, img, target, callCmd)
	if err != nil {
		return img, err
	}

	target = "g/queue/v1/queueid.go"
	queueCmd := append(goCmd, "queue", "Queue", "QueueErr", "q", "queue", "Q", "errqueue")
	img, err = generateIdFile(ctx, img, target, queueCmd)
	if err != nil {
		return img, err
	}
	target = "g/queue/v1/rowid.go"
	rowCmd := append(goCmd, "-p", "queue", "Row", "r", "row")
	img, err = generateIdFile(ctx, img, target, rowCmd)
	if err != nil {
		return img, err
	}
	target = "g/queue/v1/queuemsgid.go"
	queuemsgCmd := append(goCmd, "-p", "queue", "QueueMsg", "m", "msg")
	img, err = generateIdFile(ctx, img, target, queuemsgCmd)
	if err != nil {
		return img, err
	}
	target = "g/file/v1/fileid.go"
	fileCmd := append(goCmd, "file", "File", "FileErr", "f", "file", "F", "errfile")
	img, err = generateIdFile(ctx, img, target, fileCmd)
	if err != nil {
		return img, err
	}
	target = "g/test/v1/testid.go"
	testCmd := append(goCmd, "test", "Test", "TestErr", "t", "\\test", "T", "errtest")
	img, err = generateIdFile(ctx, img, target, testCmd)
	if err != nil {
		return img, err
	}

	// methodcall
	file := "command/boilerplateid/template/idanderr.tmpl"
	exist, _ = findFilesWithPattern(file)
	if !exist {
		log.Fatalf("There are no such file: %s", file)
	}
	target = "g/methodcall/v1/methodcallid.go"
	methodcallCmd := append(goCmd, "methodcall", "Methodcall", "MethodcallErr", "m", "methodcall", "M", "errmeth")
	img, err = generateIdFile(ctx, img, target, methodcallCmd)
	if err != nil {
		return img, err
	}

	return img, nil
}

func generateIdFile(ctx context.Context, img *dagger.Container, filePath string, goCmd []string) (*dagger.Container, error) {
	/*
	 *	A helper function for func generateApiID
	 */
	fileContent, err := img.WithExec(goCmd).Stdout(ctx)
	if err != nil {
		return img, err
	}

	// write the output to the target file
	newfile := dagger.ContainerWithNewFileOpts{
		Contents:    fileContent,
		Permissions: 0644,
	}
	img = img.WithNewFile(filePath, newfile)

	// export the file to the host
	dir := path.Dir(filePath)
	output := img.Directory(dir)
	if _, err := output.Export(ctx, dir); err != nil {
		return img, err
	}

	return img, nil
}

func buildAPlugin(ctx context.Context, img *dagger.Container, fileDir string, target string, packagePath string) (*dagger.Container, error) {
	/*
	 *	A helper function for func buildPlugins
	 */
	exist, _ := findFilesWithSuffixRecursively(fileDir, ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory %s and its subdirectories", fileDir)
	}

	img = img.WithExec([]string{"rm", "-f", target})
	// go build
	goCmd := []string{goToPlugin, "build"}
	for _, arg := range extraPluginArgs {
		goCmd = append(goCmd, arg)
	}
	goCmd = append(goCmd, "-o", target, packagePath)
	img = img.WithExec(goCmd)

	return img, nil
}

func buildAClientService(ctx context.Context, img *dagger.Container, fileDir string, target string, packagePath string) (*dagger.Container, error) {
	/*
	 *	A helper function for func buildClientSideOfAPIs
	 */
	exist, _ := findFilesWithSuffixRecursively(fileDir, ".go")
	if !exist {
		log.Fatalf("There are no such files *.go in the directory %s and its subdirectories", fileDir)
	}

	img = img.WithExec([]string{"rm", "-f", target})
	// go build
	goCmd := []string{goToWASM, "build"}
	for _, arg := range extraWASMCompArgs {
		goCmd = append(goCmd, arg)
	}
	goCmd = append(goCmd, "-tags", `"buildvcs=false"`, "-o", target, packagePath)
	img = img.WithExec(goCmd)

	return img, nil
}

func sqlcForQueue(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to generate sqlc for queue: apiplugin/queue/db.go
	 */
	dir := "apiplugin/queue"
	exist, _ := findFilesWithSuffixRecursively(dir, ".sql")
	if !exist {
		log.Fatalf("There are no such files *.sql in the directory %s and its subdirectories", dir)
	}
	yamlName := dir + "/sqlc/sqlc.yaml"
	exist, _ = findFilesWithPattern(yamlName)
	if !exist {
		log.Fatalf("There are no such file: %s", yamlName)
	}

	img, err := img.WithWorkdir("apiplugin/queue/sqlc").
		WithExec([]string{"sqlc", "generate"}).
		WithWorkdir("/workspaces/parigot").Sync(ctx)
	if err != nil {
		return img, err
	}

	return img, nil
}

func findFilesWithSuffixRecursively(path string, suffix string) (bool, error) {
	/*
	 *	This is a helper function that recursively finds all files in the
	 *	current folder and subfolders based on their suffix names
	 */
	exist := false
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("Cannot find such files: %s", err)
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), suffix) {
			exist = true
		}
		return nil
	})

	return exist, err
}

func findFilesWithPattern(pattern string) (bool, error) {
	/*
	 *	This is a helper function that finds all files based on
	 *	wildcard matching
	 */
	files, err := filepath.Glob(pattern)
	exist := false
	if err != nil {
		log.Fatalf("Cannot find such files %s: %s", pattern, err)
	}

	if len(files) != 0 {
		exist = true
	}
	return exist, nil
}
