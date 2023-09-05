# parigot-example/hello-world

Quickstart:

* This is for people that know golang.
* Create a codespace.
* Make sure you are in the `helloworld` directory, not `parigot-example`.

#### The most important files
* `main.go`: the hello world program that calls the greeting service
* `greeting.go`: the service that can return a greeting in a few languages
* `proto/greeting/v1/greeting.proto`: the protocol schema for calling the greeting service
* `helloworld.toml`: deployment descriptor
* `g/`: generated code goes here, contents can be deleted at will.

#### Makefile targets
All Makefile targets are `.PHONY`, they contain no dependencies.
* `make tools`: get all the tools for a specific commit it of parigot
    * You need to do this once at startup. After that, you only need this if you change the parigot version.
* `make generate`: convert the `proto/greeting/v1/greeting.proto` spec into go code.
    * You need to use this when you change the `.proto` file.
* `make` builds the source.
    * This uses golang's build tooling, which is exact, so it doesn't hurt to run it anytime.  No changes, nothing done.

#### Run it
* `runner helloworld.toml`




