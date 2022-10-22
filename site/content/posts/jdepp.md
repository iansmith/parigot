---
title: "jdepp on trial"
date: 2022-10-21T19:32:33Z
draft: false
---

Yesterday, I would have to say that jdepp was on trial, and I think he won.  I spent a
good part of yesterday working on a tool I've christened `jdepp`.  All of the good
'dep' name were already taken.

jdepp solves a problem that has been a rock in my shoe for a long time: the problem of
go and generated code. The go compiler is blindingly fast and it is tempting, especially
early in a project, to just run the darn thing every time for every binary.  Given that it
both exact and uses all cores, it's hard to imagine that it could be faster.  The problem
is that it doesn't know when a source file of some kind that _produces_ go code is out
of date. So, it will happily do no work to recompile something because it sees the binary
and all the go source code as up to date.

You could try to use `go generate` but that is a lose too.  The command can, if used
carefully and correctly, regenerate files that are derived from some source file that
isn't in the domain of the compiler.  The problem with this is that then *you* have
to know when to run `go generate` to regenerate the appropriate files.  It won't be long
until you just run go generate every time you want to build and force all the binaries
to be recompiled and linked.

jdepp is primitive and idiosyncratic, but it solves this problem.  It works very much
like `makedepend` does for C, except that it has a lot of understanding of parigot in
addition to go source code and modules.  So, it generates a set of dependencies for a
given binary that _includes_ all the necessary dependent files like .proto files, template
files, and the like.  These are just dependencies, you still have to put the recipies
in the Makefile, but the hard part has been done for you.

Because jdepp knows about parigot, it lets you as a user of parigot (not a developer)
tell parigot where your source installation is and it will generate dependencies on that
as needed. So, you can do something like, `git pull` to update your copy of parigot
and then do `jdepp -ppath <path to parigot source> <path to my code>` and it will
updated the dependencies as needed. If you *are* developing on parigot, this is even
better because this means if you change the binary interface .proto file `abi.proto`
then everything that depends on it gets marked out of date and will be built by make
the next time. "Everything" here include sample programs and such.  jdepp even knows
about the template files that drive the code generator for `protoc-gen-parigot` and
it will include dependencies such that if the template(s) changes, the correct 
recompilations will be done.

I had spent a small, not not insignificant amount of time trying to maintain a Makefile
by hand for the project and it was becoming impossible.  Now, thanks to `jdepp` that's
not a problem I have anymore.
