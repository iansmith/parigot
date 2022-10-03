---
title: "Is Os"
date: 2022-10-01T16:51:56-04:00
draft: false
---

I've been thinking of this project, Parigot, as an operating system.  In particular
the program that invokes wasm "apps" will be called the "kernel".  There
are some business reasons (fear) that one might not want to advertise an "os" or "kernel"
to the public at large.  That warning aside, I think about Parigot as an operating 
system because:
* It can run any binary, built from any source language, that is compatible with its
hardware platform. Since the hardware platform is WASM, there are already a number
of languages that can generate code for it.
* A key element in Parigot is the Application _Binary_ interface (ABI), a la linux. Again,
the language you program in doesn't matter if you can  call the correct binary entry
point.
* When a program does something bad like dividing by zero, uses a bad instruction,
or indexes outside the bounds of an array, a trap is generated to the controlling program,
the kernel. Typically, the kernel will kill the offending user process.
* The memory space of the kernel is protected from user programs.  
* The kernel provides services to user programs that, typically, involve
resources like the network and disk.  Naturally, resources like the CPU and memory
are managed by the kernel to insure that a user program does not exceed some limit.

## Properties
The kernel, parigot, has some strange, perhaps "unusual", properties:

* The kernel ABI does not offer synchronous calls.  All the return values of system calls
can return partial, incremental results.  This is modeled in Parigot's API as a
"channel".
* To enforce memory isolation of processes, the linker and loader conspire to
prevent access to memory outside the user's WASM Module, which is the equivalent of
an address space.
* Parigot only allows user-level access to services that are specifically named.  This includes
kernel services.
* Parigot's ABI does not offer access to threads proper.  It does offer access
to "Workers" for background tasks, but, broadly speaking, user processes should be
_though of_ as singly-threaded.  This development model means that "scaling up"
to multiple, concurrent service implementations is as easy [sic] as possible.
* The kernel itself is a collection of microservices, or a "world".  In development
situations, the world of the kernel is provided as, in linux terms, VDSOs. In production,
kernel services are provided by normal services accessed via, and secured at, the network.

