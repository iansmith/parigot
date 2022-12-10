---
title: "Product Brief v0"
date: 2022-10-03T04:42:48-04:00
draft: false
---

# Parigot
Parigot is an open source suite of tools that make the construction of microservice-based applications
10x better than "doing everything by hand" and probably 20x better than the combination of
dubious Makefiles and finicky shell scripts you are probably used to.

If you need to position this set of tools in your head: "Heroku for microservices, with some
'Firebaseness' rolled in".  Less sober alternatives might be "making microservices suck less"
or "it's a lot easier to test and debug a monolith than a swarm of services.".

Parigot supports only microservices based on WASM.  Given the broad support available now
for WASM, your can probably program in your favorite language and compile to WASM.  Parigot
does not care (and does not know) what tools you used to build your service.  Any properly
formed WASM file is a Parigot microservice, although the WASM file can't do much unless it is
linked against the client side library of parigot.

# You should probably ignore or take with a grain of salt all the information below this point. It was written well before we knew what we were doing.


# ParigotApp
ParigotApp is an app that supports some parts of normal development workflow.  The  
development workflows that are supported by parigotApp are:

* deployment
* repeatable integration testing
* simple autoscaling
* simple time-based processes (cron)
* timing testing

If you don't why it was necessary to put "repeatable" in front of integration testing,
you probably haven't tried to integration test your microservices-based app.

The last element, "timing testing", may need some explanation.  Timing testing is designed
to help you find some of the most insidious bugs in a microservice architecture, bugs that
depend on the exact timing/ordering of the concurrent activities in a cooperating set of 
services.  Parigot timing testing, perhaps analogously to golang's fuzz testing, allows you
to specify that you believe there is a timing issue involving these N services.  Parigot will
run integration tests provided by you, and try to force different orderings by choosing to schedule
the services on the CPU in _legal but horrifically evil_ ways.  A correct program should pass
all the integration tests no matter what scheduling policies are in use.

parigotApp is for-pay service.  Nothing requires you to use parigotApp and the basic
libraries and tools will work fine with the free, open-source version.  

# Parigot tooling
Within the open-source parigot project, numerous tools and libraries are provided for
you to use.  These tools and libraries differ in many cases between what parigot calls
"situations".  The supported situations for developers of a service are:

* local
* demo
* production

"Local" is for local development and testing.  Local development of a parigot-based application,
despite being composed of many microservices, is typically done with a single program and no networking
between the services or database. We will discuss the tools to achieve this below.  Having a single 
program/process that is _all_ the code/services makes it easy to startup and shutdown.  Similarly,
it is far easier to use a debugger on a single program than a collection of processes, or, even
worse, some array of containers.  Finally, this monolithic version of the application has
no networking associated with it, service calls are _literally_ function calls.

For apps with a large number of microservices, or collections of apps, 
the "local situation" can be applied to a subset of the services (the ones you are 
currently working on) with the others being provided production  or staging via parigotApp.   

## program: parigot-wind
"Wind" is pronounced like the act of causing a watch's spring to tighten, not like the 
movement of air.  Parigot-wind is a program that takes a collection of microservices
that were compiled to use parigot and "winds" them together.  The result of a wind is:
* Some collection of services previously contained in many WASM files are now contained in a 
single, "wound" WASM binary.  The resulting monolithic application is suitable for local 
development as it is far easier to manage than a set of services and also it will have 
exactly one code path that can be unit tested.  Parigot-wind does its job via binary editing.
* Parigot-wind knows how to find parigot library calls that don't make sense when all the services
are in the same binary. A wound binary will have all the networking calls between services
that are wound together to be replaced with simple function calls.

## deploy=git push
Heroku has shown the way here.  You push your code to a git repository of some kind
and a post-commit hook in the repository uses the API of parigot app to update the 
running microservices.  Parigot can detect which microservices have changed and in
a production situation can be told to only update the changed services.




