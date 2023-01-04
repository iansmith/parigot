API_VERSION=v1

all: g/file/$(API_VERSION)/file.pb.go build/protoc-gen-parigot

API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
EXAMPLE_PROTO=$(shell find example -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")
## we just use a single representative file for all the generated code
g/file/$(API_VERSION)/file.pb.go: $(API_PROTO) $(EXAMPLE_PROTO) $(TEST_PROTO) build/protoc-gen-parigot
	rm -rf g/*
	buf lint
	buf generate

.PHONY: protoclean
protoclean: 
	rm -rf g/*

.PHONY: clean
clean: protoclean
	rm -f build/*

TEMPLATE=$(shell find command/protoc-gen-parigot -type f -regex ".*\.tmpl")
GENERATOR_SRC=$(shell find command/protoc-gen-parigot -type f -regex ".*\.go")
build/protoc-gen-parigot: $(TEMPLATE) $(GENERATOR_SRC)
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/protoc-gen-parigot

