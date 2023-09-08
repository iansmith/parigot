---
title: Hello World.
weight: 1
date: 2023-07-04
description: This is the classic, prepared in the style of microservices and finished with drizzle of farm-raised unit testing.
---

If you want to just get straight to the _stuff you gotta type_ use the nav
bar on the right to jump to [Source Code Files]({{< ref "#source-code-files" >}}) or even
the [hello world code proper]({{< ref "#hello-world-code-proper" >}}).

---------------

{{% alert title="Warning" color="warning" %}}

This explanation of parigot's hello world program is too long.
This effort was begun as a snack and ended up being a week's worth of
rations for a camping trip.
{{%/alert %}}

One of the purposes of the ur-version of hello world was to show readers that a 
textually small program could be self-contained and do something useful.

```b
	main( ) {
 		extrn a, b, c;
 		putchar(a); putchar(b); putchar(c); putchar(’!*n’);
	}
	a ’hell’;
	b ’o, w’;
	c ’orld’;

```

> A Tutorial Introduction to the Programming Language B,
 B. Kernighan, 1973.

We are going to walk through all the files in the
[helloworld](https://github.com/iansmith/parigot/tree/master/example/helloworld)
example. We will pay special attention to the make targets as these are
re-usable and in many cases explain "how to do things" with parigot.

## Boilerplate files

We'll start with the boilerplate files that any project will need for doing
development. We will describe the purpose of each of these files, even
if their content can be largely ignored for now.

### Makefile

The project's [Makefile](https://github.com/iansmith/parigot/blob/master/example/helloworld/Makefile) is really just a repository for useful commands, not
something that does builds by looking at files and computing new files based
on recipies. Because golang has a much better build ability than most other
languages, typically you can just do `make` anytime and let the Makefile try
to engage the go build process. Go will do nothing if no compilation
is needed.

#### Make commands ("targets")

##### all

`all` builds the two wasm files that make up the program, `hello.p.wasm` and
 `greeting.p.wasm`.  This is the default command, so it is run if you type
 just `make`.  The two wasm files are built into the `build` subdirectory.
 The build is sufficiently fast that doing `make clean; make` is also ok.

{{% alert title="Deep cut" color="info" %}}
parigot uses the traditional go compiler, version 1.21.x, from google's go team for the parigot
infrastructure. However, sometimes the infrastructure needs to build programs
for the local host (native code) such as  the `runner` program. Other times it builds 
wasm files that 
are libraries that your code can use in _wasmland_.  When building programs
like example/helloworld it is used only to build wasm code.
{{% /alert %}}

The `all` command invokes the compiler to build the wasm (guest) code.  This
code sets the `GOOS` and `GOARCH` environment variables to convince the google
go compiler to emit WASM code. As was stated before, because the go compiler
does dependency analysis of the go source, no attempt is made to elaborate
dependencies in the Makefile.

This `all` target's command in the Makefile would be the part that would be
replaced by a call to other go compilers (gccgo, tinygo) or compilers for other
compiled programming lanugages like C#, Java, or Rust.  The output files are
suffixed with ".p.wasm" to indicate that they are intended to be dynamically linked
against parigot; without parigot, these programs will not execute in any Host.

{{% alert title="Deep Cut" color="info" %}} 

The `GOOS` setting for the google go
compiler is `wasipi` for [WASI](https://wasi.dev) preview version 1.  The
parigot project currently doesn't support WASI--although it's easily possible to
use it--and expects that commercial ventures will do so in the future.  

{{%/alert %}}

##### test

The `test` make target builds an runs the example test.  `test` uses the standard go testing
framework. This test command is useful if you want to run unit tests that are
completely guest-side code.  Because WASM isn't directly executable,
`go test` will
not work, you have build a "program" that is the test and then run that program
via a Host.  This test uses [wazero CLI](https://wazero.io) as the
Host, since it is already present in our dev container, but any host should
work.

{{% alert title="Penguinism" color="warning" %}}
On linux it is possible to run WASM programs directly. Effectively, and
as always on Linux, one is required to (carefully!) configure your system with obscure
commands, troubling, out of date config files, references to teenage, scandinavian developers,
and inscrutable yet hopeful kernel changes.

```bash
	sudo bash -c 'echo ":wasm:M::\x00\x61\x73\x6d\x01\x00::/opt/lucet/bin/lucet-ondemand:" > /proc/sys/fs/binfmt_misc/register'
```

Most scenarios for the WASM case involve [configuring your
kernel](https://gist.github.com/jedisct1/8ce91d746e09c913ee0d0f33b0ba7981) to
"delegate" the execution responsibility to a WASM Host. In the above case, it
delegates to `lucet-ondemand` which is, naturally, an open source project that
is no longer maintained.  
{{% /alert %}}

{{< blocks/section color="light" type="text" >}}

"generate" is the most key-to-parigot make target, as it generates stubs. These
stubs implement the golang-specific type safe code that makes the parigot
programming model work for golang developers.  You need to run this anytime you
change the `greeting.proto` protobuf schema definition contained in the `proto`
directory.

The default make target all runs `generate` just to be safe, even if it is not
needed.

{{< /blocks/section  >}}

-----

##### generate

The `generate` target is a key part of the Makefile
and the build process more generally.  This make target runs three key commands.
First, it uses [buf](https://buf.build) to run a lint pass over your .proto files.
Then it uses buf to generate the typical golang types that derive from your
protobuf schema (contained in the `proto` directory).  Third, it uses 
`protoc-gen-parigot` which is a plugin for
the `protoc` protobuf compiler that generates parigot's stubs.  

In the case of `helloworld` you'll notice that the file 
[proto/greeting/v1/greeting.proto]((https://github.com/iansmith/parigot-example/blob/master/helloworld/proto/greeting/v1/greeting.proto))
which is the protobuf schema for the greeting service in our hello world program.
Running `buf generate` results in the files [g/greeting/v1/greetingserver.p.go](https://github.com/iansmith/parigot-example/blob/master/helloworld/g/greeting/v1/greetingserver.p.go),
[g/greeting/v1/greetingserviceddecl.p.go](https://github.com/iansmith/parigot-example/blob/master/helloworld/g/greeting/v1/greetingservicedecl.p.go), and
[g/greeting/v1/greetingserver.p.go](https://github.com/iansmith/parigot-example/blob/master/helloworld/g/greeting/v1/greeting.pb.go).

The lint settings that are enforced with `buf lint` are set in [proto/buf.yaml](https://github.com/iansmith/parigot-example/blob/master/helloworld/proto/buf.yaml) and they are largely what is
recommended by the buf team (who know their stuff).  The places where it deviates are
largely because of the repetitiveness of names that result.  For example,
`FileServiceServer` and `KernelErrError` seem like a bridge too far. Changing
these settings in the `buf.yaml` file is not recommended if you are working with
parigot.

The `buf generate` step generates the standard, probuf golang types based on
`protoc-gen-go`.

{{% alert title="Deep Cut" color="info" %}}

[Tinygo](https://tinygo.org) is a "go compiler for small spaces" and can generate
WASM code because it is based [llvm]((https://github.com/iansmith/parigot-example/blob/master/helloworld/proto/greeting/v1/greeting.proto)). Tinygo is known
to produce small executables.  Tinygo produces an appoximately 5K byte output
file for the golang hello world program; the google go compiler's output is
about 5M byte.  However, tinygo achieves this primarily because it does not
support enough of the golang standard library's
[reflect](`https://pkg.go.dev/reflect`) package to allow folks to call functions
or create objects from string names. 

In the google go compiler, one can get input
from the user and then turn that text into a function call inside a program! Because
of this limitation, tinygo can do dead code elimination more aggresively than
the google go compiler.

This also means that the standard `protoc-gen-go` protobuf definitions will not work
with tinygo as they use package `reflect` quite heavily.  There have been reports
that other protobuf implementations (perhaps [gogobuf](https://github.com/gogo/protobuf)?
[molecule](https://pkg.go.dev/github.com/richardartoul/molecule)?) will
generate code that does not use reflection capabilities that are excluded from Tinygo.

Possibly of interest for dodging the protobuf problem:
* [go-plugin](https://github.com/knqyf263/go-plugin)
* [vtprotobuf](https://github.com/planetscale/vtprotobuf)


{{%/alert %}}

##### clean

The `clean` make target removes all the generated files in the `g/` directory
and the compiled binaries (`.p.wasm` files) in `build`.  The 'g' stands for generated and only
automated tool's output belongs there. `make clean` does not remove the tools
installed by `make tools` in the next section.


### buf.gen.yaml, buf.work.yaml
Both of
[buf.gen.yaml](https://github.com/iansmith/parigot-example/blob/master/helloworld/buf.gen.yaml)
and
[buf.work.yaml](https://github.com/iansmith/parigot-example/blob/master/helloworld/buf.gen.yaml)
are configuration files that most folks do not need to change.  The
`buf.gen.yaml` file tells `buf` what code generators to use.  These are currently the
golang (standard) one and the parigot one.  The `buf.work.yml` file tells `buf`
where to look for your protobuf schema files, but is recommended that
all your protobuf files be in a top level directory called `proto`.

## Source code files

### That `g` directory

We covered the `g` directory previously in our discussion of the [generate]({{<
ref "#generate" >}}) target.  

The `g` directory's content is all
machine generated, and thus it always ok to delete.
You can always create the content again with `make generate`.  The `g`
directory, sadly, is probably something you _should_ check into your repo,
despite the fact that the content is generated and that practice is generally
discouraged.  The reason for this is that someone might be using your project as
a go module (as you are doing with parigot) and needs to be able to compile your
code as part of a their build when they don't have your module's source
installed.  If `g` is not checked into  __your__ repository, __the other party's__ build
will fail because the generated code in `g` is needed.

{{% alert title="Deep Cut" color="info" %}}
The origin of the name `g` is that during parigot proper's development, the
names of imports could become quite large.

```go
	import "github.com/iansmith/parigot/generated/api/guest/file/v1"
```

These names were later shorted further by some reorgazations, but the shortening
of `generated` to `g` stuck.

{{% /alert %}}

### go.mod and go.sum

These are the standard files used by the go compiler to do dependency management
and versioning.  The contents of
[go.mod](https://github.com/iansmith/parigot-example/blob/master/helloworld/go.mod)
is probably enough for most projects when they are starting out.

### helloworld.toml

This file is the [delpoyment
descriptor](https://github.com/iansmith/parigot/blob/master/example/helloworld/helloworld.toml)
for the hello world program.  Its contents tell the [runner]({{< ref "#tools"
>}}) program the information about the microservices, tests, and programs that
need to be deployed to make the entire application run.  The contents as of the
time of writing are:

```toml
## Note that because this file is expected to be consumed in a dev container
## so the files have absolute paths that are correct for the container.

## dev configuration in general
[config.dev]
## Load parigot code from shared object
ParigotLibPath="/workspaces/parigot/build/syscall.so"
ParigotLibSymbol="ParigotInitialize"
Timezone = "US/Eastern"
[config.dev.Timeout]
Startup = 100 # millis
Complete = 20 # millis


# greeting service
[config.dev.microservice.greet]
WasmPath="/workspaces/parigot/example/helloworld/build/greeting.p.wasm"
Arg=[]
Env=[]

# helloworld, it has no services that it implements, it just consumes greet and
# runs to the end of main.
[config.dev.microservice.helloworld]
WasmPath="/workspaces/parigot/example/helloworld/build/hello.p.wasm"
Arg=[]
Env=[]
# this is the crucial line for parigot. "this is just a client and should run to completion".
Main=true
```

### main.go

This is the meat and cheese of the helloworld code. This is the main program that
"drives" the hello world application. Unlike most applications of parigot, which
just "exist", this program runs to completion.

#### imports

```go
	package main

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/iansmith/parigot-example/helloworld/g/greeting/v1"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/exit"
)
```

The imports are quite straightforward and should be supplied automatically by
any sensible IDE (especially if your IDE uses
[gopls](https://pkg.go.dev/golang.org/x/tools/gopls#section-readme)). Of mild
interest are the aliases `syscallguest` and `lib`.

* syscallguest:  This alias is necessary to not conflict with the import of
`github.com/iansmith/parigot/g/syscall/v1` which is imported as `syscall`.
[syscallguest](https://github.com/iansmith/parigot/blob/master/api/guest/syscall/syscallguest.go)
has the implementation to the system calls (of parigot) for guest code in
golang. The syscall one is the definitions of the system calls generated from
[the proto
spec](https://github.com/iansmith/parigot/blob/master/api/proto/syscall/v1/syscall.proto).

* lib: The golang guest side library uses the import name "lib" despite being
`.../lib/go`.  The alias is actually superfluous but was inserted by VSCode.

{{% alert title="Change coming" color="warning" %}}

There is a plan already in place for the development
of the golang guest library it can imported at `.../lib/goroutine` and also
expose the name "lib". This addition is to utilize golang's built-in concurrency
mechanisms rather than [continuations]({{< ref"continuations" >}}).  This
second library will be more comfortable for go developers; the existing
library will remain as a model for other implementors of language bindings.

{{% /alert %}}

#### hello world code proper

In the initialization part of `helloworld` there is one key thing to notice:

[futures in launch](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/main.go#L20)

The call to `LaunchClient()` in theory would block the code from proceeding until all of it is
dependencies listed above `MustInitClient` are up and running.  Since in practice
we cannot block, the returned value is a future that we attach success and failure
functions to.

The [mainloop](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/main.go#L42)
of hello world is the critical portion, centered around `ReadOneAndCallClient`.
This is primarily for exposition, the call [MustRunClient](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/lib/go/callhelper.go#L37)
is more commonly used in such a situation.

[afterLaunch](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/main.go#L53)

`afterLaunch()` is called from the success branch of the launch future.  It does two things. One, it "locates" 
some service that implements the greeting service protocol.  It then calls the method `FetchGreeting()` on the located
service, and associates success and failure functions to the future resulting from call to `FetchGreeting()`.

"Locate" is a key operation in parigot. Locate is the one that one turns the abstraction of a interface name, like
"file.v1.File" into an object which obeys the protocol defined to be implemented by that interface.  Typically, the
file `file.proto` is going to have the functions and data associated with all the operations of a "file.v1.File".  The
true implementation of the interface "file.v1.File" may be in a different program, a different container, or a different
machine.  This does not matter to the caller, as it is only concerned that the object that it has a reference to and names a 
"file.v1.File" is something that can understands the methods defined in the file.proto specification.
has the functions

Per our exposition above, this version of the library uses [continuations]({{<
ref"continuations" >}}) and you seem them in action with the method call of
`FetchGreeting` which is defined by [greeting's protobuf
specification](https://github.com/iansmith/parigot/blob/master/example/helloworld/proto/greeting/v1/greeting.proto).
This continuation, in short, means that we do not yet know the outcome of the
call, so you need to handle both the `Success` and `Failure` cases.  In this
example, it just prints out a message on the terminal.

#### Greeting service implementation

##### "main" of Greeting

Every service--and every program like helloworld--has a `main()` function to initialize
data structures and the like at startup.

[MustInitClient](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/lib/go/callhelper.go#L80)

The service implementation of greet uses the generated function [Init()](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/g/greeting/v1/greetingserver.p.go#L54) 
to initialize and
launch itself.  As we saw with the `main()` function of helloworld, we are returned
a future that represents the success or failure of launching the program (and waiting
on its dependencies). An additional parameter returned from `Init()` is a set of
method bindings called `bindings`.  This set of bindings is only useful if you want
to manipulate the set of methods that this service responds to.  It is not something
that many programs will ever need.  This service method map will passed to
the generated [Run()](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/greeting/main.go#L41) method.

[fetchGreeting](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/greeting/main.go#L67) and
[FetchGreeting](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/greeting/main.go#L83) are
now defined by the implementation of the greeting service.  Because the greeting service has a 
known method, `FetchGreeting`, no futher initialization work is needed and the method can be implemented as `FetchGreeting`
(as seen in the greeting.proto file).  The of the "split" versions of "FetchGreeting" is so that the true call implementation,
`fetchGreeting()` can be called directly from a test.  The return value of `FetchGreeting` is naturally a future and these
can be hard to test without running the "main loop" of parigot.  See [greeting_test](https://github.com/iansmith/parigot/blob/master/example/helloworld/greeting/greeting_test.go)
 for how a method like `fetchGreeting`
is tested without using futures.  

## The Big Trick

You may have noticed a contradicton above: Both the code for the
hello world main program and the code for the greeting service have a call
to a "Run function" and use futures.  The run loop in addition to futures (continuations)
means that both of these are running single-threaded.... but, wait
we cannot have two programs _both_ running singly threaded can we?!?

parigot makes each program that has a "main" function a single-threaded guest
program, in WASM terminology.  In parigot terminology, each of these two singly threaded programs
is a _process_.  They behave much like processes in Unix/Linux for several reasons:

* Each process runs until it finishes or exits due to an error (perhaps due to `panic`).
* Two or more processes can and do run at the same time.
* Programs are memory and code isolated from one another.  
* These processes use an InterProcess Communication  (IPC) mechanism to exchange
information. In parigot, these are the calls to services.

The careful reader will have noticed that this four point definition above could be
generalized slightly from "processes" to "processes that run on different machines separated
by a network" since all the same properties apply. Parigot does indeed generalize
this and decides how to _deploy_ the application based on the
[deployment descriptor](https://github.com/iansmith/parigot/blob/master/example/helloworld/helloworld.toml).
With the code remaining unchanged but with a different deployment descriptor,
you can run your application as these "guest" processes inside a single WASM Host
(and on a single host!), or it will create multiple WASM host programs on multiple
machines, each with a single service. In this latter case, all the calls between services
are network calls. That's a microservice architecture!


##### greeting_test.go

The greeting test is quite simple: send one request (not using any
parigot machinery) to the private, implementation function 
[fetchGreeting](https://github.com/iansmith/parigot/blob/master/example/helloworld/greeting/greeting_test.go).

We test that we got back what we expected in terms of the response object
and the error code.  Then we repeat the test, but with an out of bounds
language number (defined as zero in the enum definition).

This test can run in pure WASM as a guest program, as we mentioned before.
From a quality and ease of testing standpoint, the fact that this test does
not require  *any* booted-up service to be running makes it a true unit test.

## Decorative files

Files that are not strictly necessary.

### README.md

Explanation of parigot's hello world program, in very terse form, for those
that are browsing on github.  Almost by definition, these fools (!) are not
reading this document!

### helloworld.code-workspace

The file [helloworld.code-workspace](https://github.com/iansmith/parigot-example/blob/master/helloworld/helloworld.code-workspace) contains the project level settings
for development of the project in VSCode.  The file is a jsonc document.
