# all builds the necessary parts of the two wasm files and the
# parigot system call library. because it is go, it doesn't
# build things it doesn't need to.
all: generate build/frontdoor.p.wasm build/simple.p.wasm

GO_VERSION=1.21.0
PARIGOT_VERSION=0.3-atlanta

# run unit test... this is tricky because you have to run this inside a wasm host
# and this approach will not work if you have native code in your server
# implementation
.PHONY: test
test:
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} test -c -o tester . ./simple
	wasmtime -- tester -test.v

# first WASM file, compiled with GOOS and GOARCH set 
.PHONY: build/frontdoor.p.wasm
build/frontdoor.p.wasm:
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build -o build/frontdoor.p.wasm ./frontdoor.go

# second WASM file, compiled with GOOS and GOARCH set
.PHONY: build/simple.p.wasm
build/simple.p.wasm: 
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build -o build/simple.p.wasm ./simple/main.go

# this needs to be called to regenerate the generated files in
# g/ that are derived from the .proto files.  
#.PHONY: generate
# this is kinda crap. we need to know the path of the parigot code for this
PARIGOT_PROTO=$(shell find /workspaces/parigot/api/proto -type f -regex ".*\.proto")
generate: proto/simple/v1/simple.proto proto/frontdoor/v1/frontdoor.proto $(PARIGOT_PROTO)
	buf lint
	buf generate
	./removesingle.sh

# clean
.PHONY:clean
clean:
	rm -rf build/frontdoor.p.wasm build/simple.p.wasm
	rm -rf g/*
	

# this pulls down copies of the parigot tools from the internet ... it is
# use only if you are building outside the parigot root.  it is here mostly
# for reference.
.PHONY:tools
tools: 
	go get github.com/iansmith/parigot@${PARIGOT_VERSION}
	go get github.com/iansmith/parigot/command/runner@${PARIGOT_VERSION}
	go get github.com/iansmith/parigot/command/protoc-gen-parigot@${PARIGOT_VERSION}
	go install github.com/iansmith/parigot/command/runner
	go install github.com/iansmith/parigot/command/protoc-gen-parigot
	go get github.com/iansmith/parigot/api/plugin/syscall/main@${PARIGOT_VERSION}
	go build -o plugin/syscall.so -buildmode=plugin github.com/iansmith/parigot/api/plugin/syscall/main

	