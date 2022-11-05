+++
title = "Bootstrap instructions"
description = "how to get started building/using parigot"
weight = 2
disableReadmoreNav = true
toc = true
+++

There are two parts to getting set up.  The first is installing the tooling
you need to build/use parigot and the second is installing and understanding
the source code layout.

## Tool setup

### The easy way

The easy way to get your tools set up is to just use our dev container.
This container is understood by VS Code, so the only tool you need
is said editor.

* Download [vscode](https://code.visualstudio.com/download)

You should now skip to the section on using the source.

### The hard way

If you don't want to, or can't, use VSCode, you can use the local
development setup suggested by us.  The goal is to have a setup
that has an explicit set of known dependicies so you can be sure
that you have the right tools with the right versions, etc. Let's
assume you are going to put the code in `~/stuff/parigot`.  

* `mkdir ~/stuff`
* `cd ~/stuff`
* `touch enable-parigot`
* `mkdir tmp`
* `mkdir tools`
* `mkdir deps-parigot`

The directories created at the end `tmp`, `tools` and `deps-parigot`
are for temporary files create by our
tools, programs you need to have installed, and a place to put programs
that you don't already have installed on your system, respectively.

Let's install the first dependency, go 1.19.  Probably any variant of
of go version 1.19 will do, we'll use 19.3 in this example. Download 
[the go source](https://go.dev/dl/) and build it in the deps area.
If you want to try downloading a binary version, it's not that
the binary versions won't work, it's that it's not clear _where_
the go binaries get installed.  If you know or can figure that 
out, then you can use that path instead of the ones based on 
the go source in deps-parigot below.

* `cd ~/deps-parigot`
* `wget https://go.dev/dl/go1.19.3.src.tar.gz`
* `tar xzf go1.19.3.src.tar.gz`
* `mv go go1.19.3`
* `cd go1.19.3/src`
* `./all.bash`

This build takes a few minutes.

Be aware that you must have some relatively modern version of go already 
installed on your system for this build from source to work; go is 
implemented in go.  If you don't have the "bootstrap" version of go,
it would be a good time to use a package manager or download a binary
so you can build from source.  Trust us, you'll want to delete the one
you got from the package manager the instant your source is built.

#### The linking process

With your go binaries now comfortably installed in `~/stuff/deps-parigot/go1.19.3/bin`
you can start _the linking process_.  Here is a listing of the files you
will eventually want to have in tools and how they are linked:

```
tools
├── bin
│   ├── buf -> /usr/local/bin/buf
│   ├── cmake -> /usr/local/bin/cmake
│   ├── code -> /usr/local/bin/code
│   ├── dlv -> /usr/local/bin/dlv
│   ├── docker -> /usr/local/bin/docker
│   ├── gh -> /usr/local/bin/gh
│   ├── goland -> /usr/local/bin/goland
│   ├── grpcurl -> /Users/iansmith/go/bin/grpcurl
│   ├── hugo -> /usr/local/bin/hugo
│   ├── protoc -> /Users/iansmith/pb.python3/protobuf/bazel-bin/protoc
│   ├── protoc-gen-connect-go -> /Users/iansmith/go/bin/protoc-gen-connect-go
│   ├── protoc-gen-go -> /Users/iansmith/go/bin/protoc-gen-go
│   ├── protoc-gen-go-grpc -> /Users/iansmith/go/bin/protoc-gen-go-grpc
│   ├── protoc-gen-parigot -> ../../parigot/build/protoc-gen-parigot
│   ├── python -> /usr/bin/python3
│   ├── tree -> /usr/local/bin/tree
│   ├── wasm-opt -> /usr/local/bin/wasm-opt
│   ├── wasm-pack -> /usr/local/bin/wasm-pack
│   ├── wasmtime -> ../../deps-parigot/wasmtime-v1.0.1-x86_64-macos/wasmtime
│   └── wget -> /usr/local/bin/wget
├── go1.19.3 -> ../deps-parigot/go1.19.3
├── lib
│   └── antlr-4.9.3-complete.jar
├── prepath
│   └── unxz -> /usr/local/bin/unxz
└── wabt -> ../deps-parigot/wabt
```
You don't need to study the particular paths of the tools, we'll go through most of
those.  You do need to understand the strategy that this implies.

The goal is to have everything you need to do useful work with parigot linked into
the tools directory.  Later when we set the `PATH`, you will be able to set it
to something quite short without any referencs to `/usr/local/bin` or wherever the 
package manager puts its files.

Let's set up our linking area, and link go into the tools dir.

* `cd ~/stuff/tools`
* `mkdir bin`
* `ln -s ../deps-parigot/go1.19.3 .`


If you prefer, you can go into `~/stuff/tools/bin` and link
the go executable tools individually into that directory.  We prefer the above
structure because we may need to go look at the go sources.

#### Wabt

The only other package that is complex to build from source is 
[wabt](https://github.com/WebAssembly/wabt) the Web Assembly
Binary Toolkit.  You'll need `cmake`, some cmake-compatible
build tool like `make`, and `git` installed. We suggest building it in 
deps-parigot like this:

* `cd ~/stuff/deps-parigot`
* `git clone -q  --recursive https://github.com/WebAssembly/wabt`
* `cd wabt`
* `git submodule update -q --init`
* `mkdir build`
* `cd build`
* `cmake ..`
* `make`
* `cd ..`
* `rm -rf build`

Naturally, when the build has completed, you'll want to link that into
our tools area:

* `cd ~/stuff/tools`
* `ln -s ../deps-parigot/wabt .`

### The enable script
Earlier, we did `touch enable-parigot` and did not explain it&mdash;now we will.
This is a script that you can `source` into your shell (not run it!) that 
sets up that shell's world to use parigot.  We recommend you start by copying the 
text below into your enable script in `~/stuff/enable-parigot`, then going through it 
line by line to see what you might need to adjust.  With any luck, you can
use it close to "as-is".

```
### your home dir
### on linux, this frequently is /home/USERNAME but the following is from MacOS
home=/Users/YOURUSERNAME

### parigot base dir (the one with the enable script)
parigotroot=$home/stuff

### where you put parigot (the source code)
parigot=$parigotroot/parigot

### where you put your deps that are more complicated to install
deps=$parigotroot/deps-parigot

### pointer to the binary dir inside a go1.19 distribution.  
gopath=$parigotroot/tools/go1.19.3/bin
export GOROOT=$parigotroot/tools/go1.19.3

### pointer to the binary directory inside the place where you have installed the wabt
### (web assembly binary toolkit) toolkit
# get the source here -->https://github.com/WebAssembly/wabt
wabtpath=$parigotroot/tools/wabt/bin

### tools that you might or might not have, but are useful when working
### with parigot go here
binpath=$parigotroot/tools/bin

### antlr is written in java... this is only needed at when you build parigot tools
export CLASSPATH="$parigotroot/tools/antlr-4.9.3-complete.jar"
alias antlr4='java -Xmx500M -cp "/usr/local/lib/antlr-4.9.3-complete.jar:$CLASSPATH" org.antlr.v4.Tool'

#ARARRRARGH, blech! puke! java requires an installer
export JAVA_HOME=/usr/local/opt/openjdk@17

# example invocation of antlr
# antlr4 -Dlanguage=Go -package parser -o parser blahblahblah.g4

#colors
#eval $(/usr/local/bin/gcolors)

### tools that you probably do have, but you want to override something and
### pick a specific version instead of the "usual" one.  If you are ok
### with the normal installation of tools, this dir will likely be empty.
prepath=$parigotroot/tools/prepath


# the big kahuna vvvv
#
### note: this does NOT have /usr/local/bin in the PATH.
### don't use this construction:  PATH=$PATH:blahblah, make the PATH absolute and idempotent
export PATH=$prepath:/usr/bin:/bin:/usr/sbin:/sbin:$gopath:$wabtpath:$nodepath:$JAVA_HOME/bin:$binpath

##
## PARIGOT VARS
##

export P=github.com/iansmith/parigot
```

### The go tools

Source your enable script(`source ~/stuff/enable-parigot`) and then you should be able to easily 
install these tools. All of these are written go so your go compiler is sufficient to build and
install them.

* `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`
* `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`
* `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`
* `go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest`
* `go install github.com/bufbuild/buf/cmd/buf@latest`
* `go install golang.org/x/tools/gopls@latest`
* `go install github.com/go-delve/delve/cmd/dlv@latest`
* `go install honnef.co/go/tools/cmd/staticcheck@latest`
* `ln -s ~/go/bin/* ~/stuff/tools/bin`

It is worth mentioning that some of these do not specify a specific version of the tool
and thus have the suffix "~latest".  This means that we do not believe that parigot
depends on a specifc version of the tool, but YMMV.

### Protobufs and wasmtime

These two tools are critical to parigot, so we'll give some specific instructions on
how to install them.

For protoc, you'll need to adjust the URL slightly for operating systems other than
linux. The specific release of protobufs, in this case 3.15.8, is quite important. This
sequence of commands does use `curl` and `unzip` to do the download, so you'll need those
tools.  We recommend, naturally, that you link these tools from their "usual" place on
your system into `~/stuff/tools/bin` so they can be used via the enable-script.

* `cd ~/stuff/deps-parigot`
* `mkdir protoc`
* `cd protoc`
* `PB_REL="https://github.com/protocolbuffers/protobuf/releases" curl -s -LO $PB_REL/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip`
* `unzip protoc-3.15.8-linux-x86_64.zip`
* `cd ~/stuff/tools/bin`
* `ln -s ~/stuff/deps-parigot/protoc/bin/protoc .`

wasmtime is the wasm compiler/interpreter that we use to run wasm programs *and* 
how our run-time tooling gets hooked up to running wasm programs.  We use `wget`
and `unxz` to download and unpack the tarball here.  Again, this is a good time
to install these into your `~/stuff/tools/bin`.  As before, you will want to adjust
the download URL if you are not on the operating system and architecture shown here, 
MacOS on Intel silicon.  Similarly, the `unxz`, `tar` and `ln -s` commands will
require adjusting their values slighly for other operating systems.

* `cd ~/stuff/deps-parigot`
* `wget https://github.com/bytecodealliance/wasmtime/releases/download/v2.0.1/wasmtime-v2.0.1-x86_64-macos.tar.xz`
* `unxz wasmtime-v2.0.1-x86_64-macos.tar.xz`
* `tar xf wasmtime-v2.0.1-x86_64-macos.tar`
* `cd ~/stuff/tools/bin`
* `ln -s ~/stuff/deps-parigot/wasmtime-v2.0.1-x86_64-macos/wasmtime .`

At this point you have all the principal tools you'll need to work with parigot installed
and linked into the tools directory.   In any shell where you want to do parigot
work just do `source ~/stuff/enable-parigot` and you'll get exactly the same
environment everytime.  You'll notice from the original tree snapshot above that the
developers sometimes adjust the `tools` and `tools/bin` based on their preferences
and you should do the same.

## The source code

You'll want to start by cloning the parigot source code:

* `cd ~/stuff/`
* `git clone -q git@github.com:iansmith/parigot.git`

### The easy way

If you used "the easy way" above you can launch visual studio code
on a parigot workspace like this:

* `cd ~/stuff/parigot`
* `code parigot.code-workspace`

This will launch vs code and you will notice that in the bottom right corner of 
the vs code window there a small notification that says:  "Workspace contains
a Dev Container configuration file.  Reopen workspace to develop in a container."
You should click the button that says "Reopen in container" and you will get 
into a world where using the vs code terminal (control + backtick) is one where
everything you need for parigot is already in place.

The first time you launch the dev container, it may rebuild the container via
the Dockerfile.  You'll see a message that says "Starting Dev Container"  again in the
lower right. This process takes a few minutes but is only needed the first time
you launch in a dev container.  The process that the "build" of that container
follows is very similar to "the hard way" instructions above.

Open a shell (terminal) in your devcontainer with control+backtick and notice that
you are placed in `/workspaces/parigot` and your user is `parigot`.  This directory
is "mounted" into the container and is connected to your real filesystem on your
machine so changes made there will affect your local machine's copy.  You'll also
notice that `/home/parigot` has a setup that is the same as the result of going
through "the hard way" steps.

### The hard way

If you go through the hard way instructions, congratulations.  You probably don't
need any more help.  You can launch your favorite editor (!) onto any of the files
of parigot.

## Building parigot

Parigot comes with a `Makefile` that expects you to be in the root directory of the
project, or `~/stuff/parigot` if you followed the instructions above.  You can build
the tools and runtime using just `make`.  This will build a few tools that will be
explained below.   All the built binaries are placed into `~/stuff/parigot/build` and
so if you want to be sure you've got the latest, you can blow away the contents of that
directory and run make again.  Also, `make clean` does this as well.  Assuming everything 
went ok (no errors) then you can move on to building the example application.

### Building the example application

The example application is in `~/stuff/parigot/example/vvv` and is called Vinny's
Vintage Vinyl. If you do `make` in that directory, you should see it build two 
programs, the client side and the server side.  The both of these binaries will end
up in the `~/stuff/parigot/example/vvv/build` directory, just as with parigot 
above them.  Again, `make clean` removes the contents of that directory so you can
do a fresh build.

## The tools

Below is a brief descripton of the tools that are built as part of parigot:

`runner`
: This is the primary tool for invoking (running) user programs in the parigot
environment. You can just put the compiled wasm files containing your program
on the command line.  A good example might be:
* `build/runner ~/stuff/parigot/example/vvv/build/server.p.wasm ~/stuff/parigot/example/vvv/build/storeclient.p.wasm`

The command above runs the client and server of the example app.

`surgery`
: This tool is used for binary editing of WASM files.  It works by converting an input binary
to the wasm text form (a .wat file) and editing that and the converting the result back to
a wasm binary file.  For this reason it can be quite slow--even our example program has 10s of
millions of lines of text in its .wat file since each line is roughly one WASM machine instruction.
At the present time, this tool is not used heavily.  There are plans to use it more heavily in
the future to "unlink" binaries.

`jdepp`
: This is dependency analyzer for parigot source code.  This tool does for go code (at least go
code in the parigot repository) what `makedepend` does for C code.  It rewrites the Makefile of
parigot or the example app to have the appropriate dependencies.  This is necessary for parigot
because of its heavy reliance on generated code which the go compiler is not aware of and, thus,
cannot correctly compute if a binary is out of date with respect to its source code.

`protoc-gen-parigot`
: This is the code generator, sometimes called "the compiler" of parigot.  This program converts
`.proto` files, specifically the "service" and "rpc" section of proto files to go source code
useful in parigot.

### Some miscellaneous notes about the source code

The code uses a `go.work` file at the root of the repository to control what parts of the directory
structure actually contain source code.  Each of the directories mentioned in that file is
its own go package. 

The go packages, such as `github.com/iansmith/parigot/lib` are usually mapped into a dircectory
location like `lib/atlanta.base/go`.  This is done by careful construction of the `go.mod` files
that define the root of the package.

The notation `atlanta.base` that you see above and will see in many places in the code is 
a version.  parigot's version are "numbered" by the letters of the alphabet, in order, so
a is 1.  For each letter, we choose a city and so `atlanta` is the first public release
(not the 1.0 release!).  The second public release is likely `boradino` or `barcelona`.   

The second part of the name `atlanta.base` is the _variant_.  "base" here means that the
code below this directior is the standard or "base" variant.  The only other variant currently
planned is "test" and would be called `atlanta.test`.  This variant will container extra calls
that are only useful and allowed when doing testing.  A good example of this is `SetNow()`
which causes the clock available to a user program to take a specific value.  This is highly
valuable when doing testing of programs that have parts that run at specific times, but
not something that is generally allowed.

Currently, if you want to write a program, you should probably do that be either changing the
sample program or copying the `example/vvv` directory completely and adjusting the Makefile.
There is no current evidence that "user programs" can be built outside the parigot source tree.

You will note extensive use of the `buf` command to generate code from `.proto` files.  All 
of parigot and its user programs are specified by protobuf IDL (interface definiton language).
As explained above, the program `protoc-gen-parigot` is used for emitting go code needed by
parigot programs.  You will notice that the parigot system's go code interface is specified in
`~/stuff/parigot/g` with the "g" indicating generated.  All the code in that directory can be
safely deleted at will.  Similarly, in `~/stuff/parigot/example/vvv/proto/g` are the generated
files specific to that program. 

Every protobuf specification for parigot is broken into two parts.  These are usually,
"foo.proto" and "pb/foo.proto".  This split is the split between the service definition and
rpc calls&mdash;the api&mdash;and the data portion which are called "messages" in protobuf
parlance.  For example, `~/stuff/parigot/example/vvv/business.proto` and 
`~/stuff/parigot/example/vvv/business.proto` shows this split.  parigot use the standard
`protoc-gen-go` code generator for all the data portions of the protobuf specifications.

Ids, or `lib.Id` are used in many places in parigot.  Ids are composed of two 64 bit integers,
yielding 128 bits.  The uppermost 16 bits are reserved for parigot's use so the net is 
112 bits.  Usually, these 112 bits are random but sometimes are mostly zero with the lowest
few bits representing an error code.  Each id has a ascii letter in its top 8 bits represeting
what _type_ of id it is.  For example a service id is `s` and a kernel error code is `k`.
The second highest byte currently only has one bit in use, and the lowest bit is true if
the content of the id should be considered an error code.  When printed out usually you will
see the short form, like this:  `[s-4a3e71]`.  This is the type code, "s" for service id,
and the last three bytes of the id in hex. This version is useful for debugging most 
problems.  Sometimes you will also see the full id spelled out with all
the bits in hex, like this: `[s-00-0000:0000-0000-004a-3e71]`.

Some key packages are
`github.com/iansmith/parigot/lib`
: This is the library that supports interactions with parigot.  It is roughly analagous to 
`libc.a` in unix.  The library hides most of the complexity of interacting with the operating
environment at the WASM level.  The code to pack and unpack messages from the parigot kernel
is in `lib/kernel.go`.

`github.com/iansmith/parigot/sys`
: This package provides the implementation of the "system calls" between parigot programs and
and the operating environment (kernel).  Further it provides the code that handles starting
and stopping "processes" which are called "instances" in WASM.  The sample program has two
processes, for example, one that is the server (with a main()) that provides the implementation
of the service for Vinny's Vintage Vinyl.   The client side (also with a main()) is just
a consumer of the interfaces and makes "network" calls to the microservice implementation in
the server.

At some point in the future, the standard library implementation will also be part of sys in 
`~/stuff/parigot/sys/atlanta.base/stdlib`.

`github.com/iansmith/api`
: This package provides the protocol buffers definitions for all the interfaces to parigot.
It contains no go code.





