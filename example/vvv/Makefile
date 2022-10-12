### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS.
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory of vvv (the dir containing this file)

PARIGOT_SOURCE=../..

all: proto/gen/demo/vvv/servicedecl.go

# should depend on all the toml files, which should in turn depend 1-to-1 with proto files
proto/gen/demo/vvv/servicedecl.go: proto/gen/demo/vvv/business.p.toml
	$(PARIGOT_SOURCE)/build/findservices .

clean:
	rm -rf proto/gen
