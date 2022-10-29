---
title: "First syscall"
date: 2022-10-29T11:00:05-04:00
draft: false
---
I guess I made an important decision this morning. I decided on the _tagline_ for parigot. I suppose shouldn't say this is permanent as all the decisions I make
related to pargot seem to change, but I think I'm going to use this:

{{< lead >}}  parigot: the operating environment for microservices. {{< /lead >}}

I'm going to use that capitalization too, if possible.  parigot is not capitalized, even if the first word of a sentence.  

Yesterday, I got the first system call (syscall) working, `register()`.  I've kept the name syscall, because it is like a syscall in that it a non-standard
way of calling a function and it gets wrapped by a library that hides most of that for callers, ala `libc` for unix.

A few days ago, I had [decided that having an ABI interface to parigot](../callingconv/) wasn't worth the trouble.
At that time I was calling the interface "the ABI" for application binary interface, as linux does.  This "decision" got made basically because I was spending
more time trying to generate code for the ABI than all the "normal" code and further it was becoming painfully apparent that said ABI it would end up being 
dependent [on the compiler](../compiler_change).  So, I decided that parigot's OE&ndash;_you see what I did there?_&ndash; would be called via the same mechanisms
that anything else would be called. Now this should have been a red flag but I guess I was tired and didn't think about it too clearly.  I went happily on my 
way and eventually was forced to come back to it.

The OE __cannot__ be called the same way as "normal" user code because it __implements the calling mechanism__. Like, duh.  So, that occured me yesterday
and I realized that there would have to be a different path for calling parigot, but I would try to hide that from users, the same way `libc` does.  You 
don't _think_ of `open()` actually being a system call because it looks and acts like a normal function call.

Then there was the problem of the compiler.  Because there aren't [calling conventions](../callingconv/) for WASM, I was, and still am, worried that I could 
end up in a situation where some languages can't make system calls.  So, I decided that somewhat similar to what I was trying to do with porcelain interfaces
before, I would create a very _dumb_ encoding of the parameters and return values for a system call.

<!--- 
Name, Number, Pointer, Direction,Points to, Purpose
package_name,0,YES,IN,*byte,points to the data for the name of the package (ascii, bytes)
package_len,1,NO,IN,,length of the package name in bytes
service_name,2,YES,IN,*byte,points the name of the service
service_len,3,NO,IN,,length of the service name
error_id,4,YES,OUT,*[2]int64,error id for register
service_id,5,YES,OUT,*[2]int64,service id result
https://arthursonzogni.com/Diagon/#Table
--->
```
┌────────────┬───────┬────────┬──────────┬─────────┬───────────────────────────────────────┬───────┐
│Name        │ Number│ Pointer│ Direction│Points to│ Purpose                               │       │
├────────────┼───────┼────────┼──────────┼─────────┼───────────────────────────────────────┼───────┤
│package_name│0      │YES     │ IN       │*byte    │data for the name of the package (ascii│ bytes)│
├────────────┼───────┼────────┼──────────┼─────────┼───────────────────────────────────────┼───────┤
│package_len │1      │NO      │IN        │         │length of the package name in bytes    │       │
├────────────┼───────┼────────┼──────────┼─────────┼───────────────────────────────────────┼───────┤
│service_name│2      │YES     │IN        │*byte    │points the name of the service         │       │
├────────────┼───────┼────────┼──────────┼─────────┼───────────────────────────────────────┼───────┤
│service_len │3      │NO      │IN        │         │length of the service name             │       │
├────────────┼───────┼────────┼──────────┼─────────┼───────────────────────────────────────┼───────┤
│error_id    │4      │YES     │OUT       │*[2]int64│error id for register                  │       │
├────────────┼───────┼────────┼──────────┼─────────┼───────────────────────────────────────┼───────┤
│service_id  │5      │YES     │OUT       │*[2]int64│service id result                      │       │
└────────────┴───────┴────────┴──────────┴─────────┴───────────────────────────────────────┴───────┘
```
So basically, this is going to be the standard, at least for now, for sending data to the OE.  All these slots are `INT64` in WASM parlance, or 64 bits
wide.  This is somewhat wasteful, but avoids the problem of the compiler trying to be cute with numbers actually placed _into_ the slots themselves,
such as our two lengths above. _gc_ when generating wasm won't put value into this structure, and always puts a pointer if the slots are 32 bits wide.
Since pointers in WASM in the browser are 32 bits wide, this is doubly wasteful, but KISS.

### The id type

My "id" type is encoded as two 64 bit numbers.  These ids have the first 16 bits shaved off for type information, so you can figure out what type of thing
it is even if you just get the bits. These ids are used for a few things, but notably for ids for a service and error codes from trying to `register` said
service.  The input in the table above is the input to `register`.   It is important that we return error codes and not require the OE to allocate memory, so 
all the pointers in the table above are pointers into the caller's address space.  Trying to allocate memory in somebody else's address space is a nasty
business if they have a GC running.  Ids can be printed with `Short()` to get the type and the last 3 bytes of the 14 of the content. So I service id looks 
like:
```
[s-000009]
```
when printed, meaning that service is number 9 in the OE's tables.  An error id from `register()` looks like this:
```
[r-000002]
```
to mean error code 2 from the register endpoint.  Obivously, in production ids will not always be small intergers, in particular service ids which will be 
112 random bits.  Service ids are the "token" needed to call somebody else's service. We _could_ try to enforce isolation through lack of information, since
an attacker is unlikely to be able to guess a 112 bit number.  (He might be able to guess _somebody's_ 112 bit number, though.) I have some plans to do
other isolation things as well (think: Firecracker).  A value of zero in an id means no error, no matter what service the error code is for. 
