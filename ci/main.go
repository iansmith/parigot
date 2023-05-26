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
	goToHost   = "go1.20.4"
	goToPlugin = "go1.20.4"

	// EXTRA ARGS FOR BUILDING (placed after the "go build")
	// use -x for more details from a go compiler
	extraHostArgs = ""

	syscallClientSide = "apiwasm/syscall/*.go"

	rep = "g/file/" + apiVersion + "/file.pb.go"
)

var (
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
		WithWorkdir("/workspaces/parigot")

	// // get version execute
	// golang := img.WithExec([]string{"go1.20.4", "version"})
	// version, err := golang.Stdout(ctx)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("Hello from Dagger ", version)

	// // get files in the working directory
	// contents, err := img.WithExec([]string{"ls", "/src"}).Stdout(ctx)
	// if err != nil {
	// 	return (err)
	// }
	// fmt.Println(contents)

	img, err = buildProtocGenParigot(ctx, img)
	if err != nil {
		return err
	}

	img, err = generateRep(ctx, img)
	if err != nil {
		return err
	}

	img, err = buildFilePWasm(ctx, img)
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

	// contents, err := img.WithExec([]string{"ls", "-l", "."}).Stdout(ctx)
	// if err != nil {
	// 	return img, err
	// }
	// fmt.Println(contents)

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
	img = img.WithExec([]string{"buf", "lint"})
	img = img.WithExec([]string{"buf", "generate"})

	return img, nil
}

func buildFilePWasm(ctx context.Context, img *dagger.Container) (*dagger.Container, error) {
	/*
	 *	This function `generateRep` is derived from this Makefile code:
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

// func findFilesEndWith(path string, suffix string) ([]string, error) {
// 	var files []string
// 	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if !d.IsDir() && strings.HasSuffix(d.Name(), suffix) {
// 			files = append(files, path)
// 		}
// 		return nil
// 	})
// 	fmt.Println(files)
// 	return files, err
// }
