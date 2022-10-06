---
title: "IsOs() bool"
date: 2022-10-01T16:51:56-04:00
draft: false
---

>Updated October 6.

I've been thinking of this project, parigot, as an operating system.  In particular
the program that invokes and supervises wasm "apps" will be called the "kernel".  
There are some business reasons (fear) that one might not want to advertise an "os" or 
"kernel" to the public at large.  That warning aside, I think about parigot as an 
operating system because:

* It can run any binary, built from any source language, that is compatible with its
hardware platform. Since that hardware platform is WASM, there are already a number
of languages that can generate code for it.

* A key element in Parigot is the Application _Binary_ interface (ABI), a la linux. Again,
the language you program in doesn't matter if you can  call the correct binary entry
point.

* When a program does something bad like dividing by zero, uses a bad instruction,
or indexes outside the bounds of an array, a trap is generated to the controlling program,
the kernel. Typically, the kernel will kill the offending user process.

* The memory space of the kernel is protected from user programs.  In some modes, notably
development and test, the kernel implementation of API/ABI calls may be provided at
dynamic link-time, ala Linux VDSOs. 

* The kernel provides the machinery for programs (microservices) to talk to each other.
In production, this is networking.  However, in other modes the kernel may provide
other, faster, simpler implementations.  Notably by eliminating the boundary between
services altogether and making a "service call" a "function call".

* The kernel provides services to user programs that, typically, involve
resources like the network and disk.  Naturally, resources like the CPU and memory
are managed by the kernel to ensure that a user program does not exceed some limit.

* By managing the CPU resources that a user program can consume, the kernel can
effectively implement scheduling policy.

## Properties
The kernel, parigot, has some strange, perhaps "unusual", properties:

* The kernel ABI does not offer synchronous calls for any substantial resource like
networking or disk.  All the return values of system calls can return partial, 
incremental results.  This is necessary because a microservice needs to run its function
to completion and return.

* The kernel ABI/interface as well as all inter- and intra- process communication
is defined by protobuf definitions.  The kernel uses "real" grpc in production but
not in build/test/debug.

* To enforce memory isolation of processes, the linker and loader conspire to
prevent access to memory outside the user's WASM Module, which is the equivalent of
an address space.

* Parigot only allows user-level access to services that are specifically named.  
This includes kernel services.

* Parigot's ABI does not offer access to threads proper. (In golang, go routines
are also verboten.) It does offer access to "Workers" for background tasks, but, 
broadly speaking, user processes should be _thought of_ as singly-threaded and, thus,
unwilling to block.

* Parigot prevents users from user programs (services) from having writable global
variables.  This is enforced at link time.

* The kernel itself is a collection of microservices, or a "world".  In development
situations, the world of the kernel can be provided as, in linux terms, VDSOs. In 
production, kernel services are provided by normal network services accessed via, and 
secured at, the network.

