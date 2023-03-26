API_VERSION=v1

all: allprotos \
	commands \
	apiimpl \
	methodcalltest \
	sqlc \
	static/t1.wasm

allprotos: g/file/$(API_VERSION)/file.pb.go 
methodcalltest: build/methodcalltest.p.wasm build/methodcallfoo.p.wasm build/methodcallbar.p.wasm
apiimpl: build/file.p.wasm build/log.p.wasm build/test.p.wasm build/queue.p.wasm build/dom.p.wasm build/dom.p.wasm
commands: 	build/protoc-gen-parigot build/runner build/wcl build/pbmodel
sqlc: apiimpl/queue/go_/db.go

GO_CMD=GOOS=js GOARCH=wasm go
GO_LOCAL=go

API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

SPLIT_UTIL=$(shell find apiimpl/splitutil -type f -regex ".*\.go")

## we just use a single representative file for all the generated code
REP=g/file/$(API_VERSION)/file.pb.go
$(REP): $(API_PROTO) $(TEST_PROTO) build/protoc-gen-parigot 
	rm -rf g/*
	buf lint
	buf generate

## running the ANTLR code for the protobuf3 code
pbmodel/protobuf3_parser.gom: pbmodel/protobuf3.g4 
	cd pbmodel;./generate.sh

## running the ANTLR code for the WCL parsers/lexers
ui/parser/wcl_parser.go: ui/parser/wcl.g4 
	cd ui/parser;./generate.sh

ui/parser/wcllex_lexer.go: ui/parser/wcllex.g4 
	cd ui/parser;./generate.sh

## running the ANTLR code for the css3
ui/css/css3_lexer.go: ui/css/css3.g4
	cd ui/css;./generate.sh

pbmodel/protobuf3_parser.go: pbmodel/protobuf3.g4
	cd pbmodel; ./generate.sh
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
build/runner: $(RUNNER_SRC) $(REP) $(SYS_SRC) $(SPLIT_UTIL)
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/runner

# implementation of the file service
FILE_SERVICE=$(shell find apiimpl/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) $(REP) $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/apiimpl/file

# implementation of the test service
TEST_SERVICE=$(shell find apiimpl/test -type f -regex ".*\.go")
build/test.p.wasm: $(TEST_SERVICE) $(REP) $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/apiimpl/test

# implementation of the log service
LOG_SERVICE=$(shell find apiimpl/log -type f -regex ".*\.go")
build/log.p.wasm: $(LOG_SERVICE) $(REP) $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/apiimpl/log

# queue service impl
QUEUE_SERVICE=$(shell find apiimpl/queue -type f -regex ".*\.go")
build/queue.p.wasm: $(QUEUE_SERVICE) $(REP) $(SPLIT_UTIL) apiimpl/queue/go_/db.go 
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/apiimpl/queue

# dom service impl
QUEUE_SERVICE=$(shell find apiimpl/dom -type f -regex ".*\.go")
build/dom.p.wasm: $(DOM_SERVICE) $(REP) apiimpl/dom/*.go 
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/apiimpl/dom

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
SYSCALL_CLIENT_SIDE=apiimpl/syscall/*.go

build/methodcalltest.p.wasm: $(METHODCALLTEST) $(SYSCALL_CLIENT_SIDE) g/file/$(API_VERSION)/file.pb.go build/runner $(METHODCALL_TEST_SVC) 
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall

# methodcall service impl: methodcall.FooService
FOO_SERVICE=test/func/methodcall/impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

# methodcall service impl: methodcall.BarService
BAR_SERVICE=test/func/methodcall/impl/bar/*.go
build/methodcallbar.p.wasm: $(BAR_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml 
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/bar

# sqlc for queue
QUEUE_SQL=$(shell find apiimpl/queue/go_/ -type f -regex ".*\.sql")
apiimpl/queue/go_/db.go: $(QUEUE_SQL)
	# sql.yaml has some relative paths in it, must be in correct dir
	cd apiimpl/queue/go_/sqlc && sqlc generate

# build a wasm binary for the test program
static/t1.wasm: command/t1/nap.go command/t1/main.go
	GOOS=js GOARCH=wasm go build -tags browser -o static/t1.wasm command/t1/*.go

command/t1/nap.go:  ui/testdata/event_test.wcl ui/driver/template/go.tmpl ui/parser/*.go apiimpl/dom/*.go build/wcl 
	build/wcl -o command/t1/nap.go ui/testdata/event_test.wcl

#
# TEST
#
test: methodcalltest test/func/methodcall/methodcall.toml all
	go test github.com/iansmith/parigot/apiimpl/queue/go_
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
	rm -f apiimpl/queue/go_/db.go apiimpl/queue/go_/models.go apiimpl/queue/go_/query.sql.go

.PHONY: clean
clean: protoclean parserclean
	rm -f build/* static/t1.wasm

