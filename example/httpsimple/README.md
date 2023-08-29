# parigot/example/httpsimple

Quickstart:

* This is for people that know golang.
* Make sure you are in the `httpsimple` directory, not `parigot` or `parigot/example`.
* Open the source in vscode with `code httpsimple.code-workspace` for
just the source of this example.  Switch to the shell inside the 
devcontainer (control-backtick).  You should use this as your shell
for this project.

#### The most important files
* `main.go`: the main for the "front door" that receives the initial HTTP call
* `simple/main.go`: the service that returns something simple when contacted by the front door
* `proto/greeting/v1/greeting.proto`: the protocol schema for calling the greeting service
* `httpsimple.toml`: deployment descriptor
* `g/`: generated code goes here, contents can be deleted at will.

#### Makefile targets
* `make generate`: convert the `proto/simple/v1/simple.proto` and 
`frontdoor.proto` spec into go code that ends up in `/g/simple` and 
`g/frontdoor`.
    * You need to use this when you change the `.proto` file.
* `make` builds the source.
    * This uses golang's build tooling, which is exact, so it doesn't hurt to run it anytime.  No changes, nothing done.

#### Run it
* `runner httpsimple.toml`




