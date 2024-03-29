
all: commands \
	guest \
	sqlc \
	plugins 
		
#
# GROUPS OF TARGETS
#
protos: g/file/$(API_VERSION)/file.pb.go # only need one file to trigger all being built
guest: build/file.p.wasm  build/queue.p.wasm build/nutsdb.p.wasm
commands: 	build/protoc-gen-parigot build/runner build/pdep
plugins: build/queue.so build/file.so build/syscall.so build/httpconn.so build/nutsdb.so
sqlc: api/plugin/queue/db.go
helloworld: build/greeting.p.wasm build/helloworld.p.wasm

#
# EXTRA ARGS FOR BUILDING (placed after the "go build")
# use -x for more details from a go compiler
#
#EXTRA_WASM_COMP_ARGS=-target=wasi -opt=1 -x -scheduler=none
#EXTRA_WASM_COMP_ARGS=

#no plugin setup
#EXTRA_HOST_ARGS=-tags noplugin  
#EXTRA_PLUGIN_ARGS=

# with plugins (build/*.so) setup
EXTRA_HOST_ARGS=
EXTRA_PLUGIN_ARGS=-buildmode=plugin 

SHARED_SRC=$(shell find api/shared -type f -regex ".*\.go")
SYSCALL_CLIENT_SIDE=api/guest/syscall/*.go 
LIB_SRC=$(shell find lib -type f -regex ".*\.go")
API_CLIENT_SIDE=$(LIB_SRC) $(CTX_SRC) $(SHARED_SRC) $(API_ID)

#we are the systemb library makefile so we want the libraries generated
GEN_SYS_LIB=1
SHARED_SRC=$(shell find api/shared -type f -regex ".*\.go")

#
# GO
#
GO_TO_WASM=GOOS=wasip1 GOARCH=wasm go
GO_TO_HOST=go
GO_TO_PLUGIN=go

#
# PROTOBUF FILES
#
API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")

## we just use a single representative file for all the protobuf generated code from
REP=g/file/$(API_VERSION)/file.pb.go
$(REP): $(API_PROTO) build/protoc-gen-parigot
	@rm -rf g/*
	buf lint
	GEN_SYS_LIB=$(GEN_SYS_LIB) buf generate

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
STATIC_LINK= -tags netgo,osusergob -ldflags "-linkmode 'external' -extldflags '-static'"
STATIC_LINK_SO=-tags netgo,osusergob
PLUGIN= build/queue.so build/file.so build/syscall.so build/nutsdb.so
build/runner: $(PLUGIN) $(RUNNER_SRC) $(REP) $(ENG_SRC) $(SYS_SRC) $(SHARED_SRC)
	@rm -f $@
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS)  -o $@ github.com/iansmith/parigot/command/runner

PDEP_SRC=$(shell find command/pdep -type f -regex ".*\.go")
PDEP_TEMPL_SRC=$(shell find command/pdep -type f -regex ".*\.template")
build/pdep: $(PDEP_SRC) $(PDEP_TEMPL_SRC)
	@rm -f $@
	$(GO_TO_HOST) build $(EXTRA_HOST_ARGS)  -o $@ github.com/iansmith/parigot/command/pdep/cmd/pdep


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
	g/nutsdb/v1/nutsdbid.go \
	g/test/v1/testid.go 

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
g/nutsdb/v1/nutsdbid.go: api/shared/id/id.go command/boilerplateid/main.go command/boilerplateid/template/*.tmpl
	$(GO_TO_HOST) run command/boilerplateid/main.go -p nutsdb NutsDB n nutsdb  > g/nutsdb/v1/nutsdbid.go

## client side of the nutsdb service
NUTSDB_SERVICE=$(shell find api/guest/nutsdb -type f -regex ".*\.go")
build/nutsdb.p.wasm: $(NUTSDB_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) g/nutsdb/v1/nutsdbid.go $(API_ID)
	@rm -f $@
	$(GO_TO_WASM) build  $(EXTRA_WASM_COMP_ARGS) -o $@ github.com/iansmith/parigot/api/guest/nutsdb

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

# ## client side of service impl (httpconn)
# HTTPCON_SERVICE=$(shell find api/guest/httpconnector -type f -regex ".*\.go")
# build/httpconn.p.wasm: $(HTTPCON_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) $(API_ID)
# 	@rm -f $@
# 	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/api/guest/httpconnector

# ## client side of service impl (http)
# HTTP_SERVICE=$(shell find api/guest/http -type f -regex ".*\.go")
# build/http.p.wasm: $(HTTP_SERVICE) $(REP) $(SYSCALL_CLIENT_SIDE) $(API_ID)
# 	@rm -f $@
# 	$(GO_TO_WASM) build $(EXTRA_WASM_COMP_ARGS) -tags "buildvcs=false" -o $@ github.com/iansmith/parigot/api/guest/http



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
build/queue.so: $(QUEUE_PLUGIN)  $(ENG_SRC) $(SHARED_SRC) $(API_ID) api/plugin/queue/db.go 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS)  -o $@ github.com/iansmith/parigot/api/plugin/queue/main

NUTSDB_PLUGIN=$(shell find api/plugin/nutsdb -type f -regex ".*\.go")
build/nutsdb.so: $(NUTSDB_PLUGIN)  $(ENG_SRC) $(SHARED_SRC) $(API_ID) 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS)  -o $@ github.com/iansmith/parigot/api/plugin/nutsdb/main

FILE_PLUGIN=$(shell find api/plugin/file -type f -regex ".*\.go")
build/file.so: $(FILE_PLUGIN) $(SYS_SRC) $(ENG_SRC) $(SHARED_SRC) $(API_ID) 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS)  -o $@ github.com/iansmith/parigot/api/plugin/file/main

HTTPCON_PLUGIN=$(shell find api/plugin/httpconnector -type f -regex ".*\.go")
build/httpconn.so: $(HTTPCON_PLUGIN) $(SYS_SRC) $(ENG_SRC) $(SHARED_SRC) $(API_ID) 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS)  -o $@ github.com/iansmith/parigot/api/plugin/httpconnector/main

SYSCALL_PLUGIN=$(shell find api/plugin/syscall -type f -regex ".*\.go")
build/syscall.so: $(SYSCALL_PLUGIN) $(SYS_SRC) $(ENG_SRC)  $(SHARED_SRC) $(API_ID) 
	@rm -f $@
	$(GO_TO_PLUGIN) build $(EXTRA_PLUGIN_ARGS)  -o $@ github.com/iansmith/parigot/api/plugin/syscall/main


## this is only used for building the (multi-arch) container that is used for the base
## of (non-multiarch) apps, most people don't need to care about this
.PHONY: koyeb
koyeb:
	docker buildx build --platform linux/amd64,linux/arm64 -t iansmith/parigot-koyeb-base-0.3.1 --push -f deploy/Dockerfile.buildbase --push deploy 


#
# deploy using our scripts
#
HTTPSIMPLE_GO_SRC=$(shell find example/httpsimple -type f -regex ".*\.go")
HTTPSIMPLE_PROTO_SRC=$(shell find example/httpsimple -type f -regex ".*\.proto")
deploy-httpsimple: build/pdep deploy/Dockerfile.buildbase command/pdep/cmd/pdep/*.template $(HTTPSIMPLE_GO_SRC) $(HTTPSIMPLE_PROTO_SRC)
	mkdir -p /workspaces/parigot/deploy/build
	rm -f /workspaces/parigot/deploy/build/*
	cp /home/parigot/deps/caddy build/*.so build/runner /workspaces/parigot/deploy/build
	cd example/httpsimple && make 
	@### if you get
	@### ERROR: Multiple platforms feature is currently not supported for docker driver. Please switch to a different driver (eg. "docker buildx create --use")
	@## then type the command suggested at the end and then the makefile will work
	pdep -u iansmith example/httpsimple
#
# TEST
#
.PHONY: test
test: sqlc helloworldtest
	go test -v github.com/iansmith/parigot/api/plugin/queue
	go test -v github.com/iansmith/parigot/api/plugin/file
	go test -v github.com/iansmith/parigot/api/plugin/nutsdb
	go test -v github.com/iansmith/parigot/lib/go/future
	go test -v github.com/iansmith/parigot/api/plugin/syscall
	cd example/kv && make clean && make && runner kv.toml
	cd example/helloworld && make clean && make && runner helloworld.toml
###
### DOCS
###

DOCROOT=site/content/en/docs
PROTOROOT=site/content/en/docs/reference/api/proto
GUESTROOT=site/content/en/docs/reference/api/guest
PLUGINROOT=site/content/en/docs/reference/api/plugin

# this should probably be written with make recipies, dependencies
.PHONY: docs
docs:
	#protobufs
	protoc --doc_out=$(PROTOROOT) --doc_opt=markdown,api.md -I api/proto \
	api/proto/file/v1/file.proto \
	api/proto/syscall/v1/syscall.proto \
	api/proto/queue/v1/queue.proto \
	api/proto/protosupport/v1/protosupport.proto \
	api/proto/queue/v1/queue.proto \
	api/proto/test/v1/test.proto

	# fix date on the frontmatter
	cat $(PROTOROOT)/frontmatter.tmpl| sed -e "s/_date_/`date -I`/" > \
		$(PROTOROOT)/frontmatter_date.md
	cat $(PROTOROOT)/frontmatter_date.md $(PROTOROOT)/api.md > $(PROTOROOT)/_index.md
	rm $(PROTOROOT)/frontmatter_date.md
	rm $(PROTOROOT)/api.md

	## guest docs
	#github.com/iansmith/parigot/api/guest/test 
	GOFLAGS= gomarkdoc -o $(GUESTROOT)/guest.md \
	github.com/iansmith/parigot/api/guest/queue/lib \
	github.com/iansmith/parigot/api/guest/syscall \
	github.com/iansmith/parigot/lib/go \
	github.com/iansmith/parigot/lib/go/client \
	github.com/iansmith/parigot/lib/go/future \
	github.com/iansmith/parigot/api/shared \
	github.com/iansmith/parigot/api/shared/id

	## fix date on frontmatter
	cat $(GUESTROOT)/frontmatter.tmpl| sed -e "s/_date_/`date -I`/" > $(GUESTROOT)/frontmatter_date.md
	cat $(GUESTROOT)/frontmatter_date.md $(GUESTROOT)/guest.md > $(GUESTROOT)/_index.md
	rm $(GUESTROOT)/frontmatter_date.md
	rm $(GUESTROOT)/guest.md

	## plugin docs
	GOFLAGS= gomarkdoc -o $(PLUGINROOT)/plugin.md \
	github.com/iansmith/parigot/api/plugin \
	github.com/iansmith/parigot/api/plugin/file \
	github.com/iansmith/parigot/api/plugin/queue \
	github.com/iansmith/parigot/api/plugin/syscall 
	cat $(PLUGINROOT)/frontmatter.tmpl| sed -e "s/_date_/`date -I`/" > $(PLUGINROOT)/frontmatter_date.md
	cat $(PLUGINROOT)/frontmatter_date.md $(PLUGINROOT)/plugin.md > $(PLUGINROOT)/_index.md
	rm $(PLUGINROOT)/frontmatter_date.md
	rm $(PLUGINROOT)/plugin.md

##
## HELLOWORLD
##
build/helloworld.p.wasm: example/helloworld/main.go example/helloworld/g/greeting/v1/greetingserver.p.go
	cd example/helloworld && ${GO_TO_WASM} build -o ../../build/helloworld.p.wasm ./main.go 

build/greeting.p.wasm: example/helloworld/greeting/main.go example/helloworld/g/greeting/v1/greetingserver.p.go
	cd example/helloworld && ${GO_TO_WASM} build -o ../../build/greeting.p.wasm ./greeting/main.go 

example/helloworld/g/greeting/v1/greetingserver.p.go: example/helloworld/proto/greeting/v1/greeting.proto 
	cd example/helloworld && make generate

.PHONY: helloworldtest
helloworldtest: 
	GOOS=wasip1 GOARCH=wasm go test -exec 'wasmtime --' -v github.com/iansmith/parigot/example/helloworld/greeting

build/tester: example/helloworld/greeting/greeting_test.go
	${GO_TO_WASM} test -c -o build/tester github.com/iansmith/parigot/example/helloworld/greeting
	


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

