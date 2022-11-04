---
title: "Four Milestones"
date: 2022-11-03T12:44:09Z
draft: false
---
I spent some time yesterday&mdash;in all honesty in the shower&mdash;trying to figure out
how I would know if parigot was "ready" to be revealed publicly.  I decide that there
were four key milestones that had to be achieved.  These are

1. Demonstrate that you can run a microservice based application in a single address
space, and make it deterministic in terms of startup and tests.
2. Demonstrate that you can patch together multiple services to create a deployment unit
or "wos" that has several services.
3. Demonstrate that you can use more than one language to write a single program.
4. Demonstrate that you can at least run with multiple separate processes talking "over the
wire" via localhost or similar.  This needs to use grpc and/or connect.

Beyond these four things, there are some features that I would consider critical to a
"launchable" version of parigot.  For the "nice things for people to build on" type of
feature, these are, with some degree of uncertainty:

* database, built around sqlc and sqlite3
* file, built around the wasmtime
* user support built into the pctx notion
* authentication, built around basic auth because it is simple enough to allow easy 
testing but allows real "logged in vs logged out" test
* some type of basic testing library, using `atlanta.test` variant or similar. This might
even include some type of mocking for services that is generated from the same input as the
normal code.
* fancy logging, since I've put a lot of resources into that for the test mode

In terms of the kernel portion, I think there are a few small features that need to be 
added as well:
* ability to clean up resources when processes die
* ability to know when all the "active" processes are dead and we should exit
* some type of background "timing" pulse that tells us 1) when to do the above and 2) when
user processes have been running too long and should be killed

In terms of the documentation, I'm not sure exactly what people are going to expect of a
new project.  I guess for myself I'd like to have two things:
* Documentation generators set up and auto-building documentation, even if that doesn't 
produce all that much right.
* One good example, probably the Vinny's Vintage Vinyl exampl, since I've been working on
that a lot as a test program.

So, how are we doing? Looking pretty poor. On the four key milestones, I expect to finish
#1 today.  I have the tooling in place to do #2 and I'm quite sure that my program 
`surgery` can be used to do it, but it's going to be some grungy work to get that to
actually happen.  I'm going to put off #2 for a bit until I have more code with more
services that it's worth testing on.  

In terms of #3 and number #4, I guess I'm going to take a stab at #4 because I am sort of
convinced it is the easier of the two and might be something that can be done with just 
link-time Makefile tricks or similar.  #3 is going to require me to get dirty, and possibly
bloody.


