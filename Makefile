API_VERSION=v1

all: allprotos \
	commands \
	apiwasm \
	methodcalltest \
	sqlc \
	build/runner 
	
#static/t1.wasm

allprotos: g/file/$(API_VERSION)/file.pb.go 
methodcalltest: build/methodcalltest.p.wasm build/methodcallfoo.p.wasm build/methodcallbar.p.wasm
apiwasm: build/file.p.wasm build/log.p.wasm build/test.p.wasm build/queue.p.wasm build/dom.p.wasm
commands: 	build/protoc-gen-parigot build/runner build/wcl build/pbmodel
sqlc: apigo/queue/go_/db.go

## -x make this empty to not show the commands on a build of wasm code
EXTRA_WASM_COMP_ARGS=-x -target=wasi 

#this command can be useful if you want to run tinygo in a container but otherwise use your host machine
#GO_CMD=docker run --rm --env CC=/usr/bin/clang --env GOFLAGS="-buildvcs=false" --mount type=bind,source=`pwd`,target=/home/tinygo/parigot --workdir=/home/tinygo/parigot parigot-tinygo:0.27 tinygo 

GO_CMD=tinygo
GO_LOCAL=go

API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

SPLIT_UTIL=$(shell find apiwasm/splitutil -type f -regex ".*\.go")
SYSCALL_CLIENT_SIDE=apiwasm/syscall/*.go $(SPLIT_UTIL)

## we just use a single representative file for all the generated code
REP=g/file/$(API_VERSION)/file.pb.go
$(REP): $(API_PROTO) $(TEST_PROTO) build/protoc-gen-parigot 
	rm -rf g/*
	buf lint
	buf generate

## running the ANTLR code for the protobuf3 code
pbmodel/protobuf3_parser.go: pbmodel/protobuf3.g4 
	cd pbmodel;./generate.sh

## running the ANTLR code for the WCL parsers/lexers
ui/parser/wcl_parser.go: ui/parser/wcl.g4 
	cd ui/parser;./generate.sh

ui/parser/wcllex_lexer.go: ui/parser/wcllex.g4 
	cd ui/parser;./generate.sh

## running the ANTLR code for the css3
ui/css/css3_lexer.go: ui/css/css3.g4
	cd ui/css;./generate.sh

#ui/css/css3_parser.go: ui/css/css3.g4
#	cd css;./generate.sh

# protoc plugin
TEMPLATE=$(shell find command/protoc-gen-parigot -type f -regex ".*\.tmpl")
GENERATOR_SRC=$(shell find command/protoc-gen-parigot -type f -regex ".*\.go")
build/protoc-gen-parigot: $(TEMPLATE) $(GENERATOR_SRC)
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/protoc-gen-parigot


# launch a deployment based on a config file with runner
RUNNER_SRC=$(shell find command/runner -type f -regex ".*\.go")
SYS_SRC=$(shell find sys -type f -regex ".*\.go")
LIB_SRC=$(shell find lib -type f -regex ".*\.go")
build/runner: $(RUNNER_SRC) $(REP) $(SYS_SRC) $(SYSCALL_CLIENT_SIDE) $(LIB_SRC) apiwasm
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/runner


# impl of the code to parse a wat from emscripten
PEM_SRC=$(shell find command/parse-emscripten-wat -type f -regex ".*\.go")
build/parse-emscripten-wat: $(PEM_SRC) 
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/parse-emscripten-wat

# implementation of the file service
FILE_SERVICE=$(shell find apiwasm/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_CMD) build  $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/file

# implementation of the test service
TEST_SERVICE=$(shell find apiwasm/test -type f -regex ".*\.go")
build/test.p.wasm: $(TEST_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_CMD) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/test

# implementation of the log service
LOG_SERVICE=$(shell find apiwasm/log -type f -regex ".*\.go")
build/log.p.wasm: $(LOG_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_CMD) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/log

# queue service impl
QUEUE_SERVICE=$(shell find apiwasm/queue -type f -regex ".*\.go")
build/queue.p.wasm: $(QUEUE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) 
	rm -f $@
	$(GO_CMD) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/queue

# dom service impl
DOM_SERVICE=$(shell find apiwasm/dom -type f -regex ".*\.go")
build/dom.p.wasm: $(DOM_SERVICE) $(REP) apiwasm/dom/*.go 
	rm -f $@
	$(GO_CMD) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/apiwasm/dom

# wcl compiler
WCL_COMPILER=$(shell find ui/parser -type f -regex ".*\.go")
CSS_COMPILER=$(shell find ui/css -type f -regex ".*\.go")
WCL_DRIVER=$(shell find ui/driver -type f -regex ".*\.go")
PB_COMPILER=$(shell find pbmodel -type f -regex ".*\.go")

build/wcl: $(WCL_COMPILER) $(CSS_COMPILER) $(WCL_DRIVER) $(PB_COMPILER) $(REP)\
	ui/driver/template/*.tmpl helper/antlr/antlr.go\
	ui/parser/wcl_parser.go ui/parser/wcllex_lexer.go pbmodel/protobuf3_parser.go ui/css/css3_lexer.go\
	helper/*.go helper/antlr/*.go
	rm -f $@
	$(GO_LOCAL) build  -o $@ github.com/iansmith/parigot/command/wcl

## model compiler
build/pbmodel: pbmodel/protobuf3_parser.go command/pbmodel/*.go pbmodel/*.go helper/*.go
	$(GO_LOCAL) build  -o $@ github.com/iansmith/parigot/command/pbmodel

# methodcall test code
METHODCALLTEST=test/func/methodcall/*.go
METHODCALL_TEST_SVC=build/methodcallbar.p.wasm build/methodcallfoo.p.wasm 
build/methodcalltest.p.wasm: $(METHODCALLTEST) $(SYSCALL_CLIENT_SIDE) g/file/$(API_VERSION)/file.pb.go build/runner $(METHODCALL_TEST_SVC) 
	rm -f $@
	$(GO_CMD) build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall

# methodcall service impl: methodcall.FooService
FOO_SERVICE=test/func/methodcall/impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_CMD) build  $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

# methodcall service impl: methodcall.BarService
BAR_SERVICE=test/func/methodcall/impl/bar/*.go
build/methodcallbar.p.wasm: $(BAR_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml $(SYSCALL_CLIENT_SIDE)
	rm -f $@
	$(GO_CMD) build $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/bar

# sqlc for queue
QUEUE_SQL=$(shell find apigo/queue/go_/ -type f -regex ".*\.sql")
apiwasm/queue/go_/db.go: $(QUEUE_SQL)
	# sql.yaml has some relative paths in it, must be in correct dir
	cd apiwasm/queue/go_/sqlc && sqlc generate

# build a wasm binary for the test program
# static/t1.wasm: command/t1/mvc.go command/t1/main.go
# 	$(GO_CMD) build -tags browser -o static/t1.wasm github.com/iansmith/parigot/command/t1

# command/t1/nap.go:  ui/testdata/event_test.wcl ui/driver/template/go.tmpl ui/parser/*.go apiwasm/dom/*.go build/wcl 
# 	build/wcl -o command/t1/nap.go ui/testdata/event_test.wcl

# command/t1/mvc.go:  ui/testdata/model_test.wcl ui/driver/template/go.tmpl ui/parser/*.go apiwasm/dom/*.go build/wcl 
# 	build/wcl -o command/t1/mvc.go ui/testdata/model_test.wcl

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
	rm -f apigo/queue/go_/db.go apigo/queue/go_/models.go apigo/queue/go_/query.sql.go

.PHONY: clean
clean: protoclean parserclean
	rm -f build/* static/t1.wasm

