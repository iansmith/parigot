---
title: "Binary Editing"
date: 2022-10-01T04:42:48-04:00
draft: false
---

### Updated Jan 2
You can't make a service call into a function call because WASM doesn't permit 
memory sharing between wasm modules.  There are [some options here](https://docs.wasmtime.dev/api/wasmtime/struct.SharedMemory.html) 
but I don't want to turn on threads.

Yesterday I started on the ability to programmatically edit WASM binaries.  I did this by
building a simple Antlr4 grammar for the WAT format of a WASM file.  I got this far enough
to parse the real output of a 'hello world' program, which is probably pretty close to 
complete since WASM doesn't have that many opcodes and a great many of them are used in
the machinery that the compiler generates as preamble to any function call and the 
preamble to any program.

This was motivated by two issues.  First, I wanted to get rid of some extra cruft that
was being generated into the WASM binaries by the golang to WASM compiler implementations. 
This extra  cruft is the same between the gc compiler and tinygo.  This cruft is designed
to support running WASM programs in a browser, which I need to do at the moment--but the cruft
is pretty tightly bound up in the generated code because the compiler emits some calls to this
cruft to handle interactions with the environment the program is running on.  Besides not
needing browser support yet, I also wanted to be sure that I was not creating a dependency
on the browser inadvertently.  

The second issue is some future features that will need binary editing under the control
of parigot. These features, as I understand them now are:

* The ability to remove networking calls between services and replace them with simple
function calls, to simplify things for development and debugging.
* The ability to *wind* services together, creating arbitrary combinations of services
that occupy the same address space or, in my case, WASM module. 
  
The latter of these could be used to allow a developer to _think of_ a set of services as
very fine-grained but not actually be _run_ as such.  The logical limit of this capability
is useful for testing, creating a single monolithic program from many microservices, so 
that the startup/shutdown of the entire _application_ is done without a need for scripts and with much less hassle.  

Further, running a debugger gets vastly simpler since there is only one process to attach to.
The former of these services can also be used in a debugging scenario to remove some parts 
of the program that are not needed.

A closely related feature would be one where a monolith is created for the purpose easing
testing--no more complex and unreliable setup scripts for running tests!  This is
even more crucial for integration tests.

## worlds
I think a collection of services that work together to reach a common
goal or "app" should be called a world.

When you wind all the services into a binary that is a single
process, or more precisely "WASM module", you can call it a "one world".

Somehow, I need to figure out how to work "one world government"
into all this.



