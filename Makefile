API_VERSION=v1

all: allprotos \
	commands \
	apiwasm \
	methodcalltest \
	sqlc \
	plugins \
	build/runner 
	
#static/t1.wasm


#
# GROUPS OF TARGETS
#
allprotos: g/file/$(API_VERSION)/file.pb.go 
methodcalltest: build/methodcalltest.p.wasm build/methodcallfoo.p.wasm build/methodcallbar.p.wasm
apiwasm: build/file.p.wasm build/test.p.wasm build/queue.p.wasm 
#commands: 	build/protoc-gen-parigot build/runner build/wcl build/pbmodel
commands: 	build/protoc-gen-parigot build/runner
plugins: build/queue.so build/file.so build/syscall.so
sqlc: apiplugin/queue/db.go

#
# EXTRA ARGS FOR BUILDING (placed after the "go build")
# use -x for more details from a go compiler
#
EXTRA_WASM_COMP_ARGS=-target=wasi -opt=1
EXTRA_HOST_ARGS=
EXTRA_PLUGIN_ARGS=-buildmode=plugin

#this command can be useful if you want to run tinygo in a container but otherwise use your host machine
#GO_TO_WASM=docker run --rm --env CC=/usr/bin/clang --env GOFLAGS="-buildvcs=false" --mount type=bind,source=`pwd`,target=/home/tinygo/parigot --workdir=/home/tinygo/parigot parigot-tinygo:0.27 tinygo 

#
# GO
#
GO_TO_WASM=tinygo
GO_TO_HOST=go
GO_TO_PLUGIN=go 

#
# PROTOBUF FILES
#
API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

SYSCALL_CLIENT_SIDE=apiwasm/syscall/*.go 

## we just use a single representative file for all the generated code
REP=g/file/$(API_VERSION)/file.pb.go
$(REP): $(API_PROTO) $(TEST_PROTO) build/protoc-gen-parigot 
	rm -rf g/*
	buf lint
	buf generate

#
# ANTLR
#
## running the ANTLR code for the protobuf3 code
#pbmodel/protobuf3_parser.go: pbmodel/protobuf3.g4 
#	cd pbmodel;./generate.sh

## running the ANTLR code for the WCL parsers/lexers
#ui/parser/wcl_parser.go: ui/parser/wcl.g4 
#	cd ui/parser;./generate.sh

#ui/parser/wcllex_lexer.go: ui/parser/wcllex.g4 
#	cd ui/parser;./generate.sh

## running the ANTLR code for the css3
#ui/css/css3_lexer.go: ui/css/css3.g4
#	cd ui/css;./generate.sh

#ui/css/css3_parser.go: ui/css/css3.g4
#	cd css;./generate.sh

#
# PROTOC EXTENSION
#
# protoc plugin
TEMPLATE=$(shell find command/protoc-gen-parigot -type f -regex ".*\.tmpl")
GENERATOR_SRC=$(shell find command/protoc-gen-parigot -type f -regex ".*\.go")
build/protoc-gen-parigot: $(TEMPLATE) $(GENERATOR_SRC)
	rm -f $@
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/protoc-gen-parigot


#
# RUNNER
#
RUNNER_SRC=$(shell find command/runner -type f -regex "command/runner/.*\.go")
RUNNER_RUNNER_SRC=$(shell find command/runner -type f -regex "command/runner/runner/.*\.go")
PLUGIN_SRC=$(shell find lib -type f -regex "apiplugin/.*/.*\.go")
build/runner: $(RUNNER_SRC) $(REP) $(RUNNER_RUNNER_SRC) plugins 
	rm -f $@
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/runner


#
# EMSCRIPTEN UTIL
#
# impl of the code to parse a wat from emscripten
# PEM_SRC=$(shell find command/parse-emscripten-wat -type f -regex ".*\.go")
# build/parse-emscripten-wat: $(PEM_SRC) 
# 	rm -f $@
# 	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/parse-emscripten-wat

#
# CLIENT SIDE OF API
#

## client side of the file service
FILE_SERVICE=$(shell find apiwasm/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/file

## client side of the test service
TEST_SERVICE=$(shell find apiwasm/test -type f -regex ".*\.go")
build/test.p.wasm: $(TEST_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/test

## client side of service impl
QUEUE_SERVICE=$(shell find apiwasm/queue -type f -regex ".*\.go")
build/queue.p.wasm: $(QUEUE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) 
	rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/queue

## dom service impl
#DOM_SERVICE=$(shell find apiwasm/dom -type f -regex ".*\.go")
#build/dom.p.wasm: $(DOM_SERVICE) $(REP) apiwasm/dom/*.go 
#	rm -f $@
#	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/dom

#
# WCL
#
# wcl compiler
# WCL_COMPILER=$(shell find ui/parser -type f -regex ".*\.go")
# CSS_COMPILER=$(shell find ui/css -type f -regex ".*\.go")
# WCL_DRIVER=$(shell find ui/driver -type f -regex ".*\.go")
# PB_COMPILER=$(shell find pbmodel -type f -regex ".*\.go")

# build/wcl: $(WCL_COMPILER) $(CSS_COMPILER) $(WCL_DRIVER) $(PB_COMPILER) $(REP)\
# 	ui/driver/template/*.tmpl helper/antlr/antlr.go\
# 	ui/parser/wcl_parser.go ui/parser/wcllex_lexer.go pbmodel/protobuf3_parser.go ui/css/css3_lexer.go\
# 	helper/*.go helper/antlr/*.go
# 	rm -f $@
# 	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS)  -o $@ github.com/iansmith/parigot/command/wcl

#
# pbmodel
#
## model compiler
build/pbmodel: pbmodel/protobuf3_parser.go command/pbmodel/*.go pbmodel/*.go helper/*.go
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/pbmodel


#
# METHODCALL TEST
#
## methodcall test code
METHODCALLTEST=test/func/methodcall/*.go
METHODCALL_TEST_SVC=build/methodcallbar.p.wasm build/methodcallfoo.p.wasm 
build/methodcalltest.p.wasm: $(METHODCALLTEST) $(SYSCALL_CLIENT_SIDE) g/file/$(API_VERSION)/file.pb.go build/runner $(METHODCALL_TEST_SVC) 
	rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall

## methodcall service impl: methodcall.FooService
FOO_SERVICE=test/func/methodcall/impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

## methodcall service impl: methodcall.BarService
BAR_SERVICE=test/func/methodcall/impl/bar/*.go
build/methodcallbar.p.wasm: $(BAR_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/bar

#
# SQL generator 
#

## sqlc for queue
QUEUE_SQL=$(shell find apiplugin/queue -type f -regex ".*\.sql")
apiplugin/queue/db.go: $(QUEUE_SQL) apiplugin/queue/sqlc/sqlc.yaml
	# sql.yaml has some relative paths in it, must be in correct dir
	cd apiplugin/queue/sqlc && sqlc generate

#
# TEST PROGRAM
#
# build a wasm binary for the test program
# static/t1.wasm: command/t1/mvc.go command/t1/main.go
# 	$(GO_TO_WASM) build -tags browser -o static/t1.wasm github.com/iansmith/parigot/command/t1

# command/t1/nap.go:  ui/testdata/event_test.wcl ui/driver/template/go.tmpl ui/parser/*.go build/wcl 
# 	build/wcl -o command/t1/nap.go ui/testdata/event_test.wcl

# command/t1/mvc.go:  ui/testdata/model_test.wcl ui/driver/template/go.tmpl ui/parser/*.go build/wcl 
# 	build/wcl -o command/t1/mvc.go ui/testdata/model_test.wcl

#
# PLUGINS
# 
QUEUE_PLUGIN=$(shell find apiplugin/queue -type f -regex ".*\.go")
build/queue.so: $(QUEUE_PLUGIN)
	rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/apiplugin/queue

FILE_PLUGIN=$(shell find apiplugin/file -type f -regex ".*\.go")
build/file.so: $(FILE_PLUGIN)
	rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/apiplugin/file

SYSCALL_PLUGIN=$(shell find apiplugin/syscall -type f -regex ".*\.go")
build/syscall.so: $(SYSCALL_PLUGIN)
	rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/apiplugin/syscall

#
# TEST
#
test: methodcalltest test/func/methodcall/methodcall.toml all
	go test github.com/iansmith/parigot/apiwasm/queue/go_
	build/runner -t test/func/methodcall/methodcall.toml 

parserclean:
	rm -f ui/parser/wcl.interp ui/parser/wcl.tokens \
	ui/parser/wcl_base_listener.go ui/parser/wcl_listener.go \
	ui/parser/wcl_parser.go ui/parser/wcllex.interp \
	ui/parser/wcllex.tokens ui/parser/wcllex_lexer.go \
	ui/parser/wcl_visitor.go ui/parser/wcl_base_visitor.go \
	ui/css/css3.interp ui/css/css3.tokens ui/css/css3Lexer.interp \
	ui/css/css3Lexer.tokens ui/css/css3_base_listener.go \
	ui/css/css3_lexer.go ui/css/css3_parser.go \
	pbmodel/protobuf3.interp pbmodel/protobuf3.tokens pbmodel/protobuf3Lexer.interp \
	pbmodel/protobuf3Lexer.tokens pbmodel/protobuf3_*.go

semfailtest: build/wcl
	build/wcl -invert ui/testdata/fail_dupparam.wcl 
	build/wcl -invert ui/testdata/fail_duplocal.wcl 
	build/wcl -invert ui/testdata/fail_duptextname.wcl 
	build/wcl -invert ui/testdata/fail_duptextnameparam.wcl 
	build/wcl -invert ui/testdata/fail_dupdocfunc.wcl 
	build/wcl -invert ui/testdata/fail_dotindocname.wcl 
	build/wcl -invert ui/testdata/fail_dotintextname.wcl 
	build/wcl -invert ui/testdata/fail_dupnamefuncdoctext.wcl
	build/wcl -invert ui/testdata/fail_dotbadformal.wcl
	build/wcl -invert ui/testdata/fail_dotbadlocal.wcl
	build/wcl -invert ui/testdata/fail_badtag.wcl
	build/wcl -invert ui/testdata/fail_baddocvar.wcl
	build/wcl -invert ui/testdata/fail_badtextvarpre.wcl
	build/wcl -invert ui/testdata/fail_badtextvarpost.wcl
	build/wcl -invert ui/testdata/fail_conflictlocalparamtext.wcl
	build/wcl -invert ui/testdata/fail_conflictlocalnametext.wcl
	build/wcl -invert ui/testdata/fail_unknowncss.wcl
	#build/wcl -invert ui/testdata/fail_badmodelmsg.wcl
	#build/wcl -invert ui/testdata/fail_badmodelderef.wcl
	build/wcl -invert ui/testdata/fail_badtextfunccall.wcl
	@echo "PASS"
	

semtest: build/wcl
	build/wcl -o /dev/null ui/testdata/textfunc_test.wcl
	build/wcl -o /dev/null ui/testdata/docfunc_test.wcl
	build/wcl -o /dev/null ui/testdata/event_test.wcl
	build/wcl -o /dev/null ui/testdata/model_test.wcl
	@echo PASS

#
# CLEAN
#
.PHONY: protoclean
protoclean: 
	rm -rf g/*
	rm -f apiplugin/queue/db.go apiplugin/queue/models.go apiplugin/queue/query.sql.go

.PHONY: clean
clean: protoclean parserclean
	rm -f build/* static/t1.wasm

