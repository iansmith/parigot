+++
title= "I know how this works"
description= """Quick start for folks that know about distributed systems."""
date='2023-07-31'
weight= 2
+++


### I know how this works

This assumes you have a lot of experience and significant knowlege of computer science.
It is not for beginners.

### What it is
parigot is an RPC toolkit for building microservices.  By virtue of its API, it ran run a collection of microservices
as a single process (e.g. a debugger works) or a separate processes connected by a network.

#### Setup
* Start with the hello world [Makefile](https://github.com/iansmith/parigot-example/blob/master/helloworld/Makefile)
* Use `make tools` from the Makefile to download a copy of program `runner`, a copy of
the protobuf plugin `protoc-gen-parigot` and the library `syscall.so`.

#### a simple program
* Define the protocols between your services as a .proto file
* Generate code via the `make generate` approach from the Makefile.
* Define a `main()` that is your program, like [helloworld](https://github.com/iansmith/parigot-example/blob/master/helloworld/main.go).  
	* Use `Launch()` to start your program.  It'll return a future.
	* On `Success()` proceed to your program proper.
	* Call methods on other services via the [generated guest code](https://github.com/iansmith/parigot-example/tree/master/helloworld/g/greeting/v1).
	* You _must_ run the loop that checks for events. For clients, 
	[MustRunClient](https://github.com/iansmith/parigot-example/blob/ddb4801f62167aff79e9d36005b21280f2e378b2/helloworld/main.go#L55) is the usual approach.
	* Use the `Exit()` and its returned future to exit.


#### a simple service
* A simple service should probably start with [greenting/main.go](https://github.com/iansmith/parigot-example/blob/master/helloworld/greeting/main.go)
* As above, you need to launch the service and deal with the future on startup.
* If you need references to other services to implement your services, do that
in the [Ready](https://github.com/iansmith/parigot-example/blob/ddb4801f62167aff79e9d36005b21280f2e378b2/helloworld/greeting/main.go#L87) method that is called just after launch. 
	* Use `Locate()`to find the other service. Locate is defined by the generated code of
	the other service.
* Implement the methods from the .proto.  It's best to do them with the 
[two parts](https://github.com/iansmith/parigot-example/blob/ddb4801f62167aff79e9d36005b21280f2e378b2/helloworld/greeting/main.go#L52) that different by a capital first letter. 
* Methods defined in the .proto will automatically be hooked up to the event loop.
* You can call [Run()](https://github.com/iansmith/parigot-example/blob/ddb4801f62167aff79e9d36005b21280f2e378b2/helloworld/greeting/main.go#L29) and
the implementations of your methods will be called when other services or
programs call them.


#### notes
##### parallelism
* A service or program is single threaded, but _different_ services can and do run in parallel.
* We use [futures](https://github.com/iansmith/parigot/blob/master/lib/go/future/doc.go) to deal
with the single-threadedness.  Within a service or program you should not block.  Nor
should you bother with locks.
* Methods in a service need to return in <50ms or it should it return an id
that be "polled" later.  Services that have methods that take too long may be
killed by parigot.
##### WASM
* Client programs are WASM executable programs 
	* Currently only golang support, but that's not not for long.
	* No built in support for WASI or WASIX.
##### ids
* Numerous built-in types are `<BLAH>Id` like `HostId`,`ServiceId`, etc. These
	are typed, random, 112 bit values.  Two 64bit unsigned ints - 2 bytes for
	typing.
* There is a tool, `boilerplateid` that can create new types of Ids if you want
your own.  See the [Makefile](https://github.com/iansmith/parigot/blob/master/Makefile) 
of parigot for examples.


