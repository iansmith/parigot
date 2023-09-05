# all builds the necessary parts of the two wasm files and the
# parigot system call library. because it is go, it doesn't
# build things it doesn't need to.
all:  build/methodcalltest.p.wasm \
build/methodcallfoo.p.wasm build/methodcallbar.p.wasm


PARIGOT_VERSION=0.3.0
GO_VERSION=1.21.0

# run unit test... this is tricky because you have to run this inside a wasm host
# and this approach will not work if you have native code in your server
# implementation
.PHONY: test
test:
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} test -c -o tester ./greeting
	wasmtime -- tester -test.v



#
# METHODCALL TEST
#
g/id/methodcallid.go: ../../command/boilerplateid/main.go ../../command/boilerplateid/template/*.tmpl
	go${GO_VERSION} run ../../command/boilerplateid/main.go -p methodcall Methodcall m methcall > g/id/methodcallid.go

## methodcall test code
METHODCALLTEST_SRC=*.go
METHODCALLTEST_SVC=build/methodcallbar.p.wasm build/methodcallfoo.p.wasm 
METHODCALL_TOML=methodcall.toml
build/methodcalltest.p.wasm: $(METHODCALLTEST_SRC) $(API_CLIENT_SIDE) $(METHODCALLTEST_SVC) $(METHODCALL_TOML) g/methodcall/v1/methodcallid.go
	@rm -f $@
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/example/methodcall


## methodcall service impl: methodcall.FooService
FOO_SERVICE=impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) $(API_CLIENT_SIDE) g/id/methodcallid.go proto/methodcall/foo/v1/foo.proto
	@rm -f $@
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build  $(EXTRA_WASM_COMP_ARGS) -o $@ $@ github.com/iansmith/parigot/example/methodcall/impl/foo

## methodcall service impl: methodcall.BarService
BAR_SERVICE=impl/bar/*.go 
build/methodcallbar.p.wasm: $(BAR_SERVICE) $(API_CLIENT_SIDE) g/id/methodcallid.go proto/methodcall/bar/v1/bar.proto
	@rm -f $@
	GOROOT=/home/parigot/deps/go${GO_VERSION} GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/example/methodcall/impl/bar


.PHONY: generate
generate: 
	buf lint
	buf generate
