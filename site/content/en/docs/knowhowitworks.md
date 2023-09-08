+++
title= "I know how this works"
description= """Quick start for folks that know about distributed systems."""
date='2023-08-08'
weight= 2
+++


### I know how this works

This assumes you have a lot of experience and significant knowlege of computer science.
It is not for beginners.

### What it is
parigot is an RPC toolkit for building microservices.  By virtue of its API, it ran run a collection of microservices
as a single process (e.g. a debugger works) or a separate processes connected by a network.

#### Setup
* Start with the [hello world example](https://github.com/iansmith/parigot/tree/master/example/helloworld)
* Copy all the files in helloworld to a new directory under `examples`, like `examples/myprog`.

#### a simple program
* "program" here means a program that runs to completion.
* Define the protocols between your services as a .proto file. You can look at 
[the greeting service](https://github.com/iansmith/parigot/tree/master/example/helloworld).
* The .proto should be placed in the same place as they are in helloworld, the
top level directory "proto".
* Generate code via the [make generate](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/Makefile#L30) approach 
from the helloworld Makefile.
* Define a `main()` that is your program, like [helloworld](https://github.com/iansmith/parigot-example/blob/master/helloworld/main.go).   It is expected that this
program will use `ExitSelf()` to end at some point.
* Use `LaunchClient()` to start your program.  It'll return a future.
* On `Success()` proceed to your program proper.
* Call methods on other services via the generated guest code like [this](https://github.com/iansmith/parigot/tree/master/example/helloworld/g/greeting/v1) for greeting service.
* You _must_ run the loop that checks for events. For programs that run to completion, 
[MustRunClient](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/lib/go/callhelper.go#L131) is the usual approach.


#### a simple service
* A simple service should have a main that sets up the service. Here is the
example from [greenting/main.go](https://github.com/iansmith/parigot-example/blob/master/helloworld/greeting/main.go)
* As above, you need to launch the service and deal with the future on startup. Although,
the generated `Init()` function for your type is easier than doing a bunch of setup.
Here is the [generated Init](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/g/greeting/v1/greetingserver.p.go#L54) from greeting service.
* If you need references to other services to implement your services, do that
in the [Ready](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/greeting/main.go#L101) method that is called just after launch. 
	* Use `Locate()`to find the other service. Locate is defined by the generated code of
	the other service,e.g. ServiceINeed.Locate().
* Implement the methods from the .proto.  It's best to do them with the 
two parts that differ by a capital first letter.  It makes for easier unit testing.
An example is [here](https://github.com/iansmith/parigot/blob/19e3202376bb2298a389186cc6fd8ce388bfd4e2/example/helloworld/greeting/main.go#L67C1-L93C2)
* Methods defined in the .proto will automatically be hooked up to the event loop.
* You can call [Run()] that has been generated for your type. 
The implementations of your methods will be called when other services or
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


