# all builds the necessary parts of the two wasm files and the
# parigot system call library. because it is go, it doesn't
# build things it doesn't need to.
all:  build/frontdoor.p.wasm build/simple.p.wasm

GO_VERSION=1.21.0

# run unit test... this is tricky because you have to run this inside a wasm host
# and this approach will not work if you have native code in your server
# implementation
.PHONY: test
test:
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} test -c -o tester . ./simple
	wasmtime -- tester -test.v

# first WASM file, compiled with GOOS and GOARCH set 
build/frontdoor.p.wasm: ./main.go
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build -o build/frontdoor.p.wasm ./main.go

# second WASM file, compiled with GOOS and GOARCH set
build/simple.p.wasm: ./simple/main.go
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build -o build/simple.p.wasm ./simple/main.go

# this needs to be called to regenerate the generated files in
# g/ that are derived from the .proto files.  You only need this
# if you change the .proto schema.
.PHONY: generate
generate: proto/simple/v1/simple.proto proto/frontdoor/v1/frontdoor.proto
	buf lint
	buf generate

# clean
.PHONY:clean
clean:
	rm -rf build/frontdoor.p.wasm build/simple.p.wasm
	rm -rf g/frontdoor g/simple
	

# this pulls down copies of the parigot tools from the internet ... it is
# use only if you are building outside the parigot root.  it is here mostly
# for reference.
.PHONY:tools
tools: generate
	go get github.com/iansmith/parigot@${PARIGOT_VERSION}
	go get github.com/iansmith/parigot/command/runner@${PARIGOT_VERSION}
	go get github.com/iansmith/parigot/command/protoc-gen-parigot@${PARIGOT_VERSION}
	go install github.com/iansmith/parigot/command/runner
	go install github.com/iansmith/parigot/command/protoc-gen-parigot
#	rm -f /home/parigot/tools/bin/runner
#	ln -s /home/parigot/go/bin/runner /home/parigot/tools/bin/runner
#	rm -f /home/parigot/tools/bin/protoc-gen-parigot
#	ln -s /home/parigot/go/bin/protoc-gen-parigot /home/parigot/tools/bin/protoc-gen-parigot
	go get github.com/iansmith/parigot/api/plugin/syscall/main@${PARIGOT_VERSION}
	go build -o plugin/syscall.so -buildmode=plugin github.com/iansmith/parigot/api/plugin/syscall/main

	