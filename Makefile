API_VERSION=v1

all: commands \
	guest \
	methodcalltest \
	sqlc \
	plugins 
	
#
# GROUPS OF TARGETS
#
protos: g/file/$(API_VERSION)/file.pb.go # only need one file to trigger all being built
methodcalltest: build/methodcalltest.p.wasm build/methodcallfoo.p.wasm build/methodcallbar.p.wasm
guest: build/file.p.wasm build/test.p.wasm build/queue.p.wasm 
commands: 	build/protoc-gen-parigot build/runner 
plugins: build/queue.so build/file.so build/syscall.so
sqlc: api/plugin/queue/db.go

#
# EXTRA ARGS FOR BUILDING (placed after the "go build")
# use -x for more details from a go compiler
#
#EXTRA_WASM_COMP_ARGS=-target=wasi -opt=1 -x -scheduler=none
EXTRA_WASM_COMP_ARGS=
EXTRA_HOST_ARGS=
EXTRA_PLUGIN_ARGS=-buildmode=plugin

SHARED_SRC=$(shell find api/shared -type f -regex ".*\.go")
SYSCALL_CLIENT_SIDE=api/guest/syscall/*.go 
LIB_SRC=$(shell find lib -type f -regex ".*\.go")
API_CLIENT_SIDE=$(LIB_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID)


CC=/usr/lib/llvm-15/bin/clang
CTX_SRC=$(shell find context -type f -regex ".*\.go")
SHARED_SRC=$(shell find api/shared -type f -regex ".*\.go")


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
build/runner: $(PLUGIN) $(RUNNER_SRC) $(REP) $(ENG_SRC) $(SYS_SRC) $(CTX_SRC) $(SHARED_SRC)
	@rm -f $@
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS) -o $@ github.com/iansmith/parigot/command/runner


#
# CLIENT SIDE OF API
#

## generate some id cruft for a couple of types built-in to parigot and
## things that are part of the public api
API_ID= \
	api/shared/id/serviceid.go \
	api/shared/id/methodid.go \
	api/shared/id/callid.go \
	g/queue/v1/queueid.go \
	g/queue/v1/rowid.go \
	g/queue/v1/queuemsgid.go \
	g/file/v1/fileid.go \
	g/test/v1/testid.go \
	g/methodcall/v1/methodcallid.go  

api/shared/id/serviceid.go:api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Service s svc > api/shared/id/serviceid.go	
api/shared/id/methodid.go:api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Method m method > api/shared/id/methodid.go	
api/shared/id/callid.go:api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Call c call > api/shared/id/callid.go	
api/shared/id/hostid.go:api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -i -p id Host h host > api/shared/id/hostid.go	

#id cruft
g/file/v1/fileid.go: api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -p file File f file  > g/file/v1/fileid.go

## client side of the file service
FILE_SERVICE=$(shell find api/guest/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) g/file/v1/fileid.go $(API_ID)
	@rm -f $@
	$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/api/guest/file

#id cruft
g/test/v1/testid.go: api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p test Test t test > g/test/v1/testid.go

## client side of the test service
TEST_SERVICE=$(shell find api/guest/test -type f -regex ".*\.go")
build/test.p.wasm: $(TEST_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) g/test/v1/testid.go $(API_ID)
	@rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/api/guest/test

#id cruft
g/queue/v1/queueid.go: api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl $(REP) 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p queue Queue q queue  > g/queue/v1/queueid.go
g/queue/v1/rowid.go: api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl $(REP) 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p queue Row r row > g/queue/v1/rowid.go
g/queue/v1/queuemsgid.go: api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl $(REP) 
	GOOS= GOARCH= $(GO_TO_HOST) run command/boilerplateid/main.go -p queue QueueMsg m msg > g/queue/v1/queuemsgid.go

## client side of service impl
QUEUE_SERVICE=$(shell find api/guest/queue -type f -regex ".*\.go")
build/queue.p.wasm: $(QUEUE_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) $(API_ID)
	@rm -f $@
	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/api/guest/queue


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

#
# SQL generator 
#

## sqlc for queue
QUEUE_SQL=$(shell find api/plugin/queue -type f -regex ".*\.sql")
api/plugin/queue/db.go: $(QUEUE_SQL) api/plugin/queue/sqlc/sqlc.yaml
	# sql.yaml has some relative paths in it, must be in correct dir
	cd api/plugin/queue/sqlc && sqlc generate

#
# PLUGINS
# 
QUEUE_PLUGIN=$(shell find api/plugin/queue -type f -regex ".*\.go")
build/queue.so: $(QUEUE_PLUGIN)  $(ENG_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID) api/plugin/queue/db.go 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/api/plugin/queue/main

FILE_PLUGIN=$(shell find api/plugin/file -type f -regex ".*\.go")
build/file.so: $(FILE_PLUGIN) $(SYS_SRC) $(ENG_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID) 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/api/plugin/file/main

SYSCALL_PLUGIN=$(shell find api/plugin/syscall -type f -regex ".*\.go")
build/syscall.so: $(SYSCALL_PLUGIN) $(SYS_SRC) $(ENG_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID) 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS) -o $@ github.com/iansmith/parigot/api/plugin/syscall/main

#
# TEST
#
#test: methodcalltest test/func/methodcall/methodcall.toml all
test:
	go test github.com/iansmith/parigot/api/plugin/queue
	go test github.com/iansmith/parigot/api/plugin/file
	go test github.com/iansmith/parigot/lib/go/future
#	build/runner -t test/func/methodcall/methodcall.toml 

#
# CLEAN
#
.PHONY: protoclean
protoclean: 
	rm -rf g/*

.PHONY: sqlclean
sqlclean:
	rm -f api/plugin/queue/db.go api/plugin/queue/models.go api/plugin/queue/query.sql.go

.PHONY: idclean
idclean:
	rm -f $(API_ID)

.PHONY: binclean
binclean:
	rm -f build/*


.PHONY: clean
clean: protoclean sqlclean idclean binclean

