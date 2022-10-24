---
title: "Callingconv"
date: 2022-10-24T08:46:09-04:00
draft: false
---
Over the course of yesterdayday afternoon and evening,
plus a little thought in bed this morning, I decided I was
going to have to abandon my idea (dream?) of having a simple,
language-neutral "core" for parigot.  

Under the original model, there would be a set of "porcelain" interfaces
to the kernal that would be simple functions.  These functions
would perform the lowest level tasks and the api would be built 
on top of these.  In principle, this would mean that you could scrap
my api and build a totally different one if you wanted to do so.

Ultimately, I decided this was going to be more trouble than it's worth. The reason is
that there is no standard calling conventions for any type of complex structure
in WASM--or not that I have found.  This means that these "lowest level"
functions would be extremely handcuffed in terms of the inbound parameters
and the outbound results.  These would be limited to the four basic
types of WASM and even within that restricted further by some limitations
of the JS<-->WASM marriage. Javascript cannot handle int64s as parameters, for
example, because they are not representable exactly in its "number" type.

Of course, all this is hidden from normal users by the compiler.  If you say you want to
pass an int64 to a function and the compiler knows about javascript's limitations it will generate
the necessary code to work around this. For example, tiny go uses the first parameter to a function
(at the implementation, not the user level)  as pointer to where results can be sent.  So if you are 
returning basically anything bigger than an int32, it will send this "extra, hidden" parameter
to your function and write the return result into the space provided at that pointer.
Similarly, for the int64 problem with javascript, it translates any parameter that is an int64 into
a pointer to an int64.  All of these shenanigans are fine because they don't impact the user...

Unless the user is me and I want to create a library that works across languages
and thus across _very different_ compilers.  I found myself having to reverse engineer what
the tinygo compiler was doing in these situations and decided that I didn't want to repeat that
for every possible WASM compiler to come up with a set of parameters and return
values that worked for each one.  

So, today I'm going to switch to a far simpler calling convention
based more strongly on the protobuf encoding.
