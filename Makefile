API_VERSION=v1

all: allprotos \
	commands \
	apiimpl \
	methodcalltest \
	sqlc

allprotos: g/file/$(API_VERSION)/file.pb.go 
methodcalltest: build/methodcalltest.p.wasm build/methodcallfoo.p.wasm build/methodcallbar.p.wasm
apiimpl: build/file.p.wasm build/log.p.wasm build/test.p.wasm build/queue.p.wasm
commands: 	build/protoc-gen-parigot build/runner 
sqlc: api_impl/queue/go_/db.go

GO_CMD=GOOS=js GOARCH=wasm go

API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

SPLIT_UTIL=$(shell find api_impl/splitutil -type f -regex ".*\.go")

## we just use a single representative file for all the generated code
REP=g/file/$(API_VERSION)/file.pb.go
$(REP): $(API_PROTO) $(TEST_PROTO) build/protoc-gen-parigot 
	rm -rf g/*
	buf lint
	buf generate

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
FILE_SERVICE=$(shell find api_impl/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) $(REP) $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/file

# implementation of the test service
TEST_SERVICE=$(shell find api_impl/test -type f -regex ".*\.go")
build/test.p.wasm: $(TEST_SERVICE) $(REP) $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/test

# implementation of the log service
LOG_SERVICE=$(shell find api_impl/log -type f -regex ".*\.go")
build/log.p.wasm: $(LOG_SERVICE) $(REP) $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/log

# queue service impl
QUEUE_SERVICE=$(shell find api_impl/queue -type f -regex ".*\.go")
build/queue.p.wasm: $(QUEUE_SERVICE) $(REP) $(SPLIT_UTIL) api_impl/queue/go_/db.go 
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/queue

# methodcall test code
METHODCALLTEST=test/func/methodcall/*.go
METHODCALL_TEST_SVC=build/methodcallbar.p.wasm build/methodcallfoo.p.wasm 
SYSCALL_CLIENT_SIDE=api_impl/syscall/*.go

build/methodcalltest.p.wasm: $(METHODCALLTEST) $(SYSCALL_CLIENT_SIDE) g/file/$(API_VERSION)/file.pb.go build/runner $(METHODCALL_TEST_SVC) $(api_impl)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall

# methodcall service impl: methodcall.FooService
FOO_SERVICE=test/func/methodcall/impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml $(api_impl)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

# methodcall service impl: methodcall.BarService
BAR_SERVICE=test/func/methodcall/impl/bar/*.go
build/methodcallbar.p.wasm: $(BAR_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml $(api_impl)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/bar

# sqlc for queue
QUEUE_SQL=$(shell find api_impl/queue/go_/ -type f -regex ".*\.sql")
api_impl/queue/go_/db.go: $(QUEUE_SQL)
	# sql.yaml has some relative paths in it, must be in correct dir
	cd api_impl/queue/go_/sqlc && sqlc generate

#
# TEST
#
test: methodcalltest test/func/methodcall/methodcall.toml all
	go test github.com/iansmith/parigot/api_impl/queue/go_
	build/runner -t test/func/methodcall/methodcall.toml 
#
# CLEAN
#
.PHONY: protoclean
protoclean: 
	rm -rf g/*
	rm api_impl/queue/go_/db.go api_impl/queue/go_/models.go api_impl/queue/go_/query.sql.go

.PHONY: clean
clean: protoclean
	rm -f build/*

