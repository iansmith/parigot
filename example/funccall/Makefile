#
# METHODCALL TEST
#
g/methodcall/v1/methodcallid.go: api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -p methodcall Methodcall m methcall > g/methodcall/v1/methodcallid.go

## methodcall test code
METHODCALLTEST_SRC=test/func/methodcall/*.go
METHODCALLTEST_SVC=build/methodcallbar.p.wasm build/methodcallfoo.p.wasm 
METHODCALL_TOML=test/func/methodcall/methodcall.toml
build/methodcalltest.p.wasm: $(METHODCALLTEST_SRC) $(API_CLIENT_SIDE) $(METHODCALLTEST_SVC) $(METHODCALL_TOML) g/methodcall/v1/methodcallid.go
	@rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall


## methodcall service impl: methodcall.FooService
FOO_SERVICE=test/func/methodcall/impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) $(API_CLIENT_SIDE) g/methodcall/v1/methodcallid.go test/func/methodcall/proto/methodcall/foo/v1/foo.proto
	@rm -f $@
	$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

## methodcall service impl: methodcall.BarService
BAR_SERVICE=test/func/methodcall/impl/bar/*.go 
build/methodcallbar.p.wasm: $(BAR_SERVICE) $(API_CLIENT_SIDE) g/methodcall/v1/methodcallid.go test/func/methodcall/proto/methodcall/bar/v1/bar.proto
	@rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/bar
