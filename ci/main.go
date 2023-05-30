package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

const (
	apiVersion = "v1"

	goToWASM   = "go1.21"
	goToHost   = "go1.19.9"
	goToPlugin = "go1.19.9"

	// EXTRA ARGS FOR BUILDING (placed after the "go build")
	// use -x for more details from a go compiler
	extraWASMCompArgs = ""
	extraHostArgs     = ""

	syscallClientSide = "apiwasm/syscall/*.go"

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
		"GOROOT": "/home/parigot/deps/go1.19.9",
	}
	goEnvVarsPlugin = map[string]string{
		"GOROOT": "/home/parigot/deps/go1.19.9",
	}

	// protobuf files
	apiProto  string
	testProto string

	// protoc plugin
	template    string
	eneratorSrc string
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
			dagger.ContainerWithDirectoryOpts{Exclude: []string{".devcontainer/", "ci/", "build/"}},
		).
		WithWorkdir("/workspaces/parigot").
		WithUser("root")

	// build client side of api
	if img, err = buildClientSideOfApi(ctx, img); err != nil {
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

func buildClientSideOfApi(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	// set up HOST env variables
	for key, value := range goEnvVarsHost {
		img = img.WithEnvVariable(key, value)
	}
	img, err := buildProtocGenParigot(ctx, img)
	if err != nil {
		return img, err
	}
	img, err = generateRep(ctx, img)
	if err != nil {
		return img, err
	}

	// set up WASM env variables
	for key, value := range goEnvVarsWASM {
		img = img.WithEnvVariable(key, value)
	}

	img, err = buildFilePWasm(ctx, img)
	if err != nil {
		return img, err
	}

	// change environment variable
	img = img.WithEnvVariable("GOOS", "js")

	img, err = buildTestPWasm(ctx, img)
	if err != nil {
		return img, err
	}
	img, err = buildQueuePWasm(ctx, img)
	if err != nil {
		return img, err
	}

	return img, nil
}

func buildProtocGenParigot(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function `buildProtocGenParigot` is derived from this Makefile code:
	 *
	 *	TEMPLATE=$(shell find command/protoc-gen-parigot -type f -regex ".*\.tmpl")
	 *	GENERATOR_SRC=$(shell find command/protoc-gen-parigot -type f -regex ".*\.go")
	 *	build/protoc-gen-parigot: $(TEMPLATE) $(GENERATOR_SRC)
	 *		@rm -f $@
	 *		$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/protoc-gen-parigot
	 */

	var err error
	template, err = img.WithExec([]string{"bash", "-c", `find command/protoc-gen-parigot -type f -regex ".*\.tmpl"`}).Stdout(ctx)
	if err != nil {
		return img, err
	}
	eneratorSrc, err = img.WithExec([]string{"bash", "-c", `find command/protoc-gen-parigot -type f -regex ".*\.go"`}).Stdout(ctx)
	if err != nil {
		return img, err
	}

	target := "build/protoc-gen-parigot"
	packagePath := "github.com/iansmith/parigot/command/protoc-gen-parigot"
	img = img.WithExec([]string{"rm", "-f", target})
	img = img.WithExec([]string{goToHost, "build", "-o", target, packagePath})

	return img, nil
}

func generateRep(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function `generateRep` is derived from this Makefile code:
	 *
	 *	API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
	 *	TEST_PROTO=$(shell find test -type f -regex ".*\.proto")
	 *
	 *	## we just use a single representative file for all the generated code from
	 *	REP=g/file/$(API_VERSION)/file.pb.go
	 *	$(REP): $(API_PROTO) $(TEST_PROTO) build/protoc-gen-parigot
	 *		@rm -rf g/*
	 *		buf lint
	 *		buf generate
	 */

	var err error
	apiProto, err = img.WithExec([]string{"bash", "-c", `find api/proto -type f -regex ".*\.proto"`}).Stdout(ctx)
	if err != nil {
		return img, err
	}
	testProto, err = img.WithExec([]string{"bash", "-c", `find test -type f -regex ".*\.proto"`}).Stdout(ctx)
	if err != nil {
		return img, err
	}

	img = img.WithExec([]string{"rm", "-rf", "g/*"})
	// switch user from root to parigot to be able to use buf package
	img = img.WithUser("parigot").WithExec([]string{"buf", "lint"})
	// switch user back to root
	img = img.WithExec([]string{"buf", "generate"}).WithUser("root")

	return img, nil
}

func buildFilePWasm(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function `buildFilePWasm` is derived from this Makefile code:
	 *
	 *	FILE_SERVICE=$(shell find apiwasm/file -type f -regex ".*\.go")
	 *	build/file.p.wasm: $(FILE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	 *		@rm -f $@
	 *		$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/file
	 */
	var err error
	_, err = img.WithExec([]string{"bash", "-c", `find apiwasm/file -type f -regex ".*\.go"`}).Stdout(ctx)
	if err != nil {
		return img, err
	}

	target := "build/file.p.wasm"
	packagePath := "github.com/iansmith/parigot/apiwasm/file"
	img = img.WithExec([]string{"rm", "-f", target})
	img = img.WithExec([]string{goToWASM, "build", "-tags", `"buildvcs=false"`, "-o", target, packagePath})

	return img, nil
}

func buildTestPWasm(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function `buildTestPWasm` is derived from this Makefile code:
	 *
	 *	TEST_SERVICE=$(shell find apiwasm/test -type f -regex ".*\.go")
	 *	build/test.p.wasm: export GOOS=js
	 *	build/test.p.wasm: export GOARCH=wasm
	 *	build/test.p.wasm: $(TEST_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	 *		@rm -f $@
	 *		$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/test
	 */
	var err error
	_, err = img.WithExec([]string{"bash", "-c", `find apiwasm/test -type f -regex ".*\.go"`}).Stdout(ctx)
	if err != nil {
		return img, err
	}

	target := "build/test.p.wasm"
	packagePath := "github.com/iansmith/parigot/apiwasm/test"
	img = img.WithExec([]string{"rm", "-f", target})
	img = img.WithExec([]string{goToWASM, "build", "-tags", `"buildvcs=false"`, "-o", target, packagePath})

	return img, nil
}

func buildQueuePWasm(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function `buildQueuePWasm` is derived from this Makefile code:
	 *
	 *	QUEUE_SERVICE=$(shell find apiwasm/test -type f -regex ".*\.go")
	 *	build/queue.p.wasm: export GOOS=js
	 *	build/queue.p.wasm: export GOARCH=wasm
	 *	build/queue.p.wasm: $(QUEUE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	 *		@rm -f $@
	 *		$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/queue
	 */
	var err error
	_, err = img.WithExec([]string{"bash", "-c", `find apiwasm/queue -type f -regex ".*\.go"`}).Stdout(ctx)
	if err != nil {
		return img, err
	}

	target := "build/queue.p.wasm"
	packagePath := "github.com/iansmith/parigot/apiwasm/queue"
	img = img.WithExec([]string{"rm", "-f", target})
	img = img.WithExec([]string{goToWASM, "build", "-tags", `"buildvcs=false"`, "-o", target, packagePath})

	return img, nil
}
