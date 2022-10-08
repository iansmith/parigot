### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory of vvv (the dir containing this file)

PARIGOT_SOURCE=../..

all: proto/gen/demo/vvv/locator.go

# toml file just representative of all the generated code
# note, tricky: proto/vvv/business.proto is not proto/DEMO/vvv/business despite the package names
# note, tricky: what is most important is that _gen_ is laid out "as you would expect"
proto/gen/demo/vvv/business.p.toml: proto/demo/vvv/business.proto
	pushd proto >& /dev/null && buf generate && popd >& /dev/null

# should depend on all the toml files, which should in turn depend 1-to-1 with proto files
proto/gen/demo/vvv/locator.go: proto/gen/demo/vvv/business.p.toml
	$(PARIGOT_SOURCE)/build/findservices .

clean:
	rm -rf proto/gen
