# all builds the necessary parts of the two wasm files and the
# parigot system call library. because it is go, it doesn't
# build things it doesn't need to.
all:  generate build/kv.p.wasm

PARIGOT_VERSION=0.3.0
GO_VERSION=1.21.1

# run unit test... this is tricky because you have to run this inside a wasm host
# and this approach will not work if you have native code in your server
# implementation
.PHONY: test
test:
	GOOS=wasip1 GOARCH=wasm go${GO_VERSION} test -c -o tester .
	wazero run tester -test.v

.PHONY: build/kv.p.wasm
build/kv.p.wasm:
	GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build -o build/kv.p.wasm ./main.go

# this needs to be called to regenerate the generated files in
# g/ that are derived from the .proto files.  You only need this
# if you change the .proto schema.
.PHONY: generate
generate:
	buf lint
	buf generate

# clean
.PHONY:clean
clean:
	rm -rf build/kv.p.wasm
	rm -rf g/*
	


# this needs to be called once when you open the dev container
# to initialize the tools. you will also need to call it if you
# do a "go get github.com/iansmith/parigot" or similar because the version
# affects the tools (runner and protoc-gen-parigot) and
# the lib, syscall.so.
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

	