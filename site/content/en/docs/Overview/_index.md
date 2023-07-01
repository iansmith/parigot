---
title: Overview
description: parigot makes building microservices fun again.
weight: 1
---

{{% pageinfo %}} 

All of this information applies the current version of parigot, `atlanta-0.3.0`.
Versions up to `some_city-1.0.0` carry no backward compatibility guarantees,
promises or anything else to protect you from change.  One may assume that
there will be significant changes in any and all publicly facing APIs and
internal code as well.  

{{% /pageinfo %}}


## What is it?
parigot is a programming language agnostic software development tool.  It offers
benefits to developers building real applications with microservices.  Although
microservices have benefits, there is simply no doubt that building,
reliably testing, debugging, and getting logging information is more 
complicated in the case of a microservice-based designed.  parigot offers a
mechanism where the application is built as a monolith, and then can be 
split apart for deployment with no code changes. 

parigot calls your set of microservices that accomplish a business task an
__application__ or the __guest__ program.  Guest programs are given to parigot
in WASM format.  Since [many programming languages](https://github.com/appcypher/awesome-wasm-langs)
are compatible with WASM, parigot neither knows nor cares what programming 
language produced or interprets the guest program.




## Why do I want parigot?

* **What is it good for?**: Reliable startup ordering.  Less flaky tests.  Debug your collection of microservices in a single debugger from your IDE.  Get your logs from a _single_ central location. Avoid dependencies-between-services surprises.

* **What is it not good for?**: parigot aims to provide many "built-in" services that meet common needs such as logging, queues, files, databases, http, etc, but it does not offer a "standard" API.  The programming model of parigot is different than most other systems, and the particulars of the programming model enable many of the advanced capabilities that parigot offers.  _For-profit entities will provide the ability to have a  more "standard" type of programming models such as a `libc.a` workalike, connections to legacy services, multi-threading of services, etc._  These types of features will, thus, be available but the free, open-source project will focus on parigot's unique programming model.

* **What is it *not yet* good for?**: Python support will be available in a foreseeable, upcoming version.  In addition, a second set of guest binding for golang will be added that makes golang programming with pargiot more idiomatic golang.  

## Where should I go next?

Give your users next steps from the Overview. For example:

* [Getting Started](/docs/getting-started/): Get started with $project
* [Examples](/docs/examples/): Check out some example code!

