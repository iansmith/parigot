# all builds the necessary parts of the two wasm files and the
# parigot system call library. because it is go, it doesn't
# build things it doesn't need to.
all:  generate build/$(APPNAME).p.wasm

APPNAME=myapp
PARIGOT_SRC=../../  #used for in-tree build
PARIGOT_VERSION=0.4.0
GO_VERSION=1.21.3

# run unit test... This is tricky because you BUILD it outside the 
# wasm world but you RUN it inside a wasm container.  This generates
# a file called tester that is a wasm binary with all the tests 
# baked into it.
.PHONY: test
test:
	GOOS=wasip1 GOARCH=wasm go${GO_VERSION} test -c -o tester .
	wazero run tester -test.v

.PHONY: build/$(APPNAME).p.wasm
build/$(APPNAME).p.wasm:
	GOOS=wasip1 GOARCH=wasm go${GO_VERSION} build -o build/$(APPNAME).p.wasm ./main.go

# this needs to be called to regenerate the generated files in
# g/ that are derived from the .proto files.  You only need this
# when you change the .proto schema.
.PHONY: generate
generate:
	buf lint
	buf generate

# clean
.PHONY:clean
clean:
	rm -rf build/$(APPNAME).p.wasm
	rm -rf g/*
	


# This needs to be called once when you open the dev container on a new project.
# It initializes the tools and libraries. Calling it again is probably
# a terrible idea as it may blow away things you want to keep.
.PHONY:refresh
tools: generate
	cp $(PARIGOT_SRC)/build/* build
# this is because of the restrictions about the paths of proto files you import
	mkdir -p proto
	cp -R $(PARIGOT_SRC)/api/proto/* proto
	



	