API_VERSION=v1

all: commands \
	apiwasm \
	methodcalltest \
	sqlc 
	
#
# GROUPS OF TARGETS
#
protos: g/file/$(API_VERSION)/file.pb.go 
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
#EXTRA_WASM_COMP_ARGS=-target=wasi -opt=1 -x -scheduler=none
EXTRA_WASM_COMP_ARGS=
EXTRA_HOST_ARGS=
EXTRA_PLUGIN_ARGS=-buildmode=plugin

SYSCALL_CLIENT_SIDE=apiwasm/syscall/*.go 
LIB_SRC=$(shell find lib -type f -regex ".*\.go")
API_CLIENT_SIDE=build/test.p.wasm build/file.p.wasm build/queue.p.wasm $(LIB_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID)


CC=/usr/lib/llvm-15/bin/clang
CTX_SRC=$(shell find context -type f -regex ".*\.go")
SHARED_SRC=$(shell find apishared -type f -regex ".*\.go")


#this command can be useful if you want to run tinygo in a container but otherwise use your host machine
#GO_TO_WASM=docker run --rm --env CC=/usr/bin/clang --env GOFLAGS="-buildvcs=false" --mount type=bind,source=`pwd`,target=/home/tinygo/parigot --workdir=/home/tinygo/parigot parigot-tinygo:0.27 tinygo 

#
# GO
#
GO_TO_WASM=GOROOT=/home/parigot/deps/go1.21 GOOS=wasip1 GOARCH=wasm go1.21
GO_TO_HOST=GOROOT=/home/parigot/deps/go1.20.4 go1.20.4
GO_TO_PLUGIN=GOROOT=/home/parigot/deps/go1.20.4 go1.20.4

#
# PROTOBUF FILES
#
API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

## we just use a single representative file for all the protobuf generated code from
REP=g/file/$(API_VERSION)/file.pb.go
$(REP): $(API_PROTO) $(TEST_PROTO) build/protoc-gen-parigot
	@rm -rf g/*
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
	@rm -f $@
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/protoc-gen-parigot


#
# RUNNER
#
RUNNER_SRC=$(shell find command/runner -type f -regex ".*\.go")
SYS_SRC=$(shell find sys -type f -regex ".*\.go")
ENG_SRC=$(shell find eng -type f -regex ".*\.go")
PLUGIN= build/queue.so build/file.so build/syscall.so
build/runner: $(PLUGIN) $(RUNNER_SRC) $(REP) $(ENG_SRC) $(SYS_SRC) $(CTX_SRC) $(SHARED_SRC) apishared/id/*.go apiplugin/* 
	@rm -f $@
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/runner

#
# patchreflectproto
#
# PATCHRP_SRC=$(shell find command/patchreflectproto -type f -regex ".*\.go")
# build/patchreflectproto: $(PATCHRP_SRC)
# 	rm -f $@
# 	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/patchreflectproto


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

## generate some id cruft for a couple of types built by parigot
API_ID= \
	apishared/id/serviceid.go \
	apishared/id/methodid.go \
	apishared/id/callid.go \
	g/queue/v1/queueid.go \
	g/queue/v1/rowid.go \
	g/queue/v1/queuemsgid.go \
	g/file/v1/fileid.go \
	g/test/v1/testid.go \
	g/methodcall/v1/methodcallid.go 

apishared/id/serviceid.go:apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Service s svc > apishared/id/serviceid.go	
apishared/id/methodid.go:apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Method m method > apishared/id/methodid.go	
apishared/id/callid.go:apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Call c call > apishared/id/callid.go	
apishared/id/hostid.go:apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Host h host > apishared/id/hostid.go	

#id cruft
g/file/v1/fileid.go: apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -p file File f file  > g/file/v1/fileid.go

## client side of the file service
FILE_SERVICE=$(shell find apiwasm/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) g/file/v1/fileid.go $(API_ID)
	@rm -f $@
	$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/file

#id cruft
g/test/v1/testid.go: apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p test Test t test > g/test/v1/testid.go

## client side of the test service
TEST_SERVICE=$(shell find apiwasm/test -type f -regex ".*\.go")
build/test.p.wasm: $(TEST_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) g/test/v1/testid.go $(API_ID)
	@rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/test

#id cruft
g/queue/v1/queueid.go: apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl $(REP) 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p queue Queue q queue  > g/queue/v1/queueid.go
g/queue/v1/rowid.go: apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl $(REP) 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p queue Row r row > g/queue/v1/rowid.go
g/queue/v1/queuemsgid.go: apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl $(REP) 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p queue QueueMsg m msg > g/queue/v1/queuemsgid.go

## client side of service impl
QUEUE_SERVICE=$(shell find apiwasm/queue -type f -regex ".*\.go")
build/queue.p.wasm: $(QUEUE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) $(API_ID)
	@rm -f $@
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
g/methodcall/v1/methodcallid.go: apishared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/idanderr.tmpl
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
build/methodcallfoo.p.wasm: $(FOO_SERVICE) $(API_CLIENT_SIDE) g/methodcall/v1/methodcallid.go
	@rm -f $@
	$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

## methodcall service impl: methodcall.BarService
BAR_SERVICE=test/func/methodcall/impl/bar/*.go
build/methodcallbar.p.wasm: $(BAR_SERVICE) $(API_CLIENT_SIDE) g/methodcall/v1/methodcallid.go
	@rm -f $@
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
build/queue.so: $(QUEUE_PLUGIN)  $(ENG_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID) apiplugin/queue/db.go build/syscall.so  wazero-src-1.1/fauxfd.go
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/apiplugin/queue/main

FILE_PLUGIN=$(shell find apiplugin/file -type f -regex ".*\.go")
build/file.so: $(FILE_PLUGIN) $(SYS_SRC) $(ENG_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID) build/syscall.so  wazero-src-1.1/fauxfd.go
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/apiplugin/file/main

SYSCALL_PLUGIN=$(shell find apiplugin/syscall -type f -regex ".*\.go")
build/syscall.so: $(SYSCALL_PLUGIN) $(SYS_SRC) $(ENG_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID) apiplugin/*.go  wazero-src-1.1/fauxfd.go
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/apiplugin/syscall/main

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

