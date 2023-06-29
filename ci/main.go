package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	fs "github.com/iansmith/parigot/ci/util"

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
	ignoreFiles = []string{".devcontainer/", "ci/", "build/", "g/", "tmp/", "ui/"}

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
		log.Fatalf("Cannot build dagger pipeline: %v\n", err)
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
	img := dockerDir.DockerBuild().
		WithDirectory("/workspaces/parigot",
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

	img, err = buildRunner(ctx, img)
	if err != nil {
		return err
	}

	img, err = buildClientSideOfAPIs(ctx, img)
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
	_, err := fs.FindFilesWithSuffixRecursively("command/runner", ".go")
	if err != nil {
		return img, err
	}
	_, err = fs.FindFilesWithPattern("api/shared/id/*.go")
	if err != nil {
		return img, err
	}
	_, err = fs.FindFilesWithPattern("api/plugin/*")
	if err != nil {
		return img, err
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
	goCmd = append(goCmd, extraHostArgs...)
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
	_, err := fs.FindFilesWithSuffixRecursively("sys", ".go")
	if err != nil {
		return img, err
	}
	// ENG_SRC
	_, err = fs.FindFilesWithSuffixRecursively("eng", ".go")
	if err != nil {
		return img, err
	}
	// CTX_SRC
	_, err = fs.FindFilesWithSuffixRecursively("context", ".go")
	if err != nil {
		return img, err
	}
	// SHARED_SRC
	_, err = fs.FindFilesWithSuffixRecursively("api/shared", ".go")
	if err != nil {
		return img, err
	}
	// check apiplugin/*.go
	_, err = fs.FindFilesWithPattern("api/plugin/*.go")
	if err != nil {
		return img, err
	}

	// apiplugin/queue/db.go
	img, err = sqlcForQueue(ctx, img)
	if err != nil {
		return img, err
	}

	// build/syscall.so
	dir := "api/plugin/syscall"
	target := "build/syscall.so"
	packagePath := "github.com/iansmith/parigot/api/plugin/syscall/main"
	img, err = buildAPlugin(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/queue.so
	dir = "api/plugin/queue"
	target = "build/queue.so"
	packagePath = "github.com/iansmith/parigot/api/plugin/queue/main"
	img, err = buildAPlugin(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/file.so
	dir = "api/plugin/file"
	target = "build/file.so"
	packagePath = "github.com/iansmith/parigot/api/plugin/file/main"
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
	syscallClientSide := "api/guest/syscall/*.go"
	_, err := fs.FindFilesWithPattern(syscallClientSide)
	if err != nil {
		return img, err
	}

	// set up WASM env variables
	for key, value := range goEnvVarsWASM {
		img = img.WithEnvVariable(key, value)
	}

	// build/file.p.wasm
	dir := "api/guest/file"
	target := "build/file.p.wasm"
	packagePath := "github.com/iansmith/parigot/api/guest/file"
	img, err = buildAClientService(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/test.p.wasm
	dir = "api/guest/test"
	target = "build/test.p.wasm"
	packagePath = "github.com/iansmith/parigot/api/guest/test"
	img, err = buildAClientService(ctx, img, dir, target, packagePath)
	if err != nil {
		return img, err
	}
	// build/queue.p.wasm
	dir = "api/guest/queue"
	target = "build/queue.p.wasm"
	packagePath = "github.com/iansmith/parigot/api/guest/queue"
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
	_, err := fs.FindFilesWithSuffixRecursively(dir, ".tmpl")
	if err != nil {
		return img, err
	}
	_, err = fs.FindFilesWithSuffixRecursively(dir, ".go")
	if err != nil {
		return img, err
	}

	// set up HOST env variables
	for key, value := range goEnvVarsHost {
		img = img.WithEnvVariable(key, value)
	}
	// set extra arguments
	goCmd := []string{goToHost, "build"}
	goCmd = append(goCmd, extraHostArgs...)
	target := "build/protoc-gen-parigot"
	packagePath := "github.com/iansmith/parigot/command/protoc-gen-parigot"
	goCmd = append(goCmd, "-o", target, packagePath)

	img = img.WithExec([]string{"rm", "-f", target}).
		WithExec(goCmd)

	return img, nil
}

func generateRep(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	generate a single representative file for all the protobuf generated code
	 */
	_, err := fs.FindFilesWithSuffixRecursively("api/proto", ".proto")
	if err != nil {
		return img, err
	}
	_, err = fs.FindFilesWithSuffixRecursively("test", ".proto")
	if err != nil {
		return img, err
	}

	img = img.WithExec([]string{"rm", "-rf", "g/*"}).
		WithExec([]string{"buf", "lint"}).
		WithExec([]string{"buf", "generate"})

	// export g dir to host
	output := img.Directory("g")
	_, err = output.Export(ctx, "g")
	if err != nil {
		return img, err
	}

	return img, nil
}

func generateApiID(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to generate some id cruft for a couple of types built by parigot:
	 *		apishared/id/serviceid.go,		apishared/id/methodid.go
	 *		apishared/id/callid.go,			g/queue/v1/queueid.go,
	 *		g/queue/v1/rowid.go,			g/queue/v1/queuemsgid.go,
	 *		g/file/v1/fileid.go,			g/test/v1/testid.go,
	 *		g/methodcall/v1/methodcallid.go
	 */
	apiID := "api/shared/id/id.go"
	boilerplateid := "command/boilerplateid/main.go"
	goCmd := []string{goToHost, "run", boilerplateid}

	_, err := fs.FindFilesWithPattern(apiID)
	if err != nil {
		return img, err
	}
	_, err = fs.FindFilesWithPattern(boilerplateid)
	if err != nil {
		return img, err
	}
	_, err = fs.FindFilesWithPattern("command/boilerplateid/template/*.tmpl")
	if err != nil {
		return img, err
	}

	target := "api/shared/id/serviceid.go"
	serviceCmd := append(goCmd, "-i", "-p", "id", "Service", "s", "svc")
	img, err = generateIdFile(ctx, img, target, serviceCmd)
	if err != nil {
		return img, err
	}
	target = "api/shared/id/methodid.go"
	methodCmd := append(goCmd, "-i", "-p", "id", "Method", "m", "method")
	img, err = generateIdFile(ctx, img, target, methodCmd)
	if err != nil {
		return img, err
	}
	target = "api/shared/id/callid.go"
	callCmd := append(goCmd, "-i", "-p", "id", "Call", "c", "call")
	img, err = generateIdFile(ctx, img, target, callCmd)
	if err != nil {
		return img, err
	}

	target = "g/queue/v1/queueid.go"
	queueCmd := append(goCmd, "-p", "queue", "Queue", "q", "queue")
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
	fileCmd := append(goCmd, "-p", "file", "File", "f", "file")
	img, err = generateIdFile(ctx, img, target, fileCmd)
	if err != nil {
		return img, err
	}
	target = "g/test/v1/testid.go"
	testCmd := append(goCmd, "-p", "test", "Test", "t", "test")
	img, err = generateIdFile(ctx, img, target, testCmd)
	if err != nil {
		return img, err
	}

	// methodcall
	file := "command/boilerplateid/template/idanderr.tmpl"
	_, err = fs.FindFilesWithPattern(file)
	if err != nil {
		return img, err
	}
	target = "g/methodcall/v1/methodcallid.go"
	methodcallCmd := append(goCmd, "-p", "methodcall", "Methodcall", "m", "methodcall")
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
	_, err := fs.FindFilesWithSuffixRecursively(fileDir, ".go")
	if err != nil {
		return img, err
	}

	img = img.WithExec([]string{"rm", "-f", target})
	// go build
	goCmd := []string{goToPlugin, "build"}
	goCmd = append(goCmd, extraPluginArgs...)
	goCmd = append(goCmd, "-o", target, packagePath)
	img = img.WithExec(goCmd)

	return img, nil
}

func buildAClientService(ctx context.Context, img *dagger.Container, fileDir string, target string, packagePath string) (*dagger.Container, error) {
	/*
	 *	A helper function for func buildClientSideOfAPIs
	 */
	_, err := fs.FindFilesWithSuffixRecursively(fileDir, ".go")
	if err != nil {
		return img, err
	}

	img = img.WithExec([]string{"rm", "-f", target})
	// go build
	goCmd := []string{goToWASM, "build"}
	goCmd = append(goCmd, extraWASMCompArgs...)
	goCmd = append(goCmd, "-tags", `"buildvcs=false"`, "-o", target, packagePath)
	img = img.WithExec(goCmd)

	return img, nil
}

func sqlcForQueue(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function is to generate sqlc for queue: apiplugin/queue/db.go
	 */
	dir := "api/plugin/queue"
	_, err := fs.FindFilesWithSuffixRecursively(dir, ".sql")
	if err != nil {
		return img, err
	}
	yamlName := dir + "/sqlc/sqlc.yaml"
	_, err = fs.FindFilesWithPattern(yamlName)
	if err != nil {
		return img, err
	}

	img = img.WithWorkdir("api/plugin/queue/sqlc").
		WithExec([]string{"sqlc", "generate"}).
		WithWorkdir("/workspaces/parigot")

	return img, nil
}
