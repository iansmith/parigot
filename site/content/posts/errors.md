---
title: "Errors"
date: 2023-01-26T17:39:02-04:00
draft: false
---
I realized that it been a while since I had written a post.  I suppose that is
good, because the reason is that things in the software itself have been coming
together.

After a lot of back and forth and trying various things to see how they "felt",
I finally made a decision on errors.  The question was in what way should errors
be signaled from the implementation of a system call to the client (WASM) side
or from the implementation of a service.  

I've decided that the go functions that return results to be consumed by the
WASM side, should return errors "out of band".  In this case, out of band means
that the errors are returned separately from the response object. So a function
that is implementing some part of a system call might look like this:

```go
func fooImplementation(req *FooRequest, resp *FooResponse) (lib.Id, string)
```

So if there are no errors the return values are `nil` and `""`.  The happy path
will consume the results of the function via the resp parameter.  

### Why lib.id and string?

The reason for using lib.Id and string as the return values is that this somewhat
simulates the "error" object in go, but in a language independent way. Normally,
the `lib.Id` return value is an error id, such `KernelErrorId` or `QueueErrorId`.
The string should be more details about the error and is intended to be read by
humans.   This setup makes it easy to write tests, becasue the automated tests
can just check the `lib.Id` and ignore the string with the error details. 

### Crossing the boundary

I said that the go implementations of services and the kernel would return errors
out of band. Well, that's really not possible when we have to travers the boundary
between go and WASM-land. 

For example, here is the structure used when we cross the boundary. The same
structure is used for all the calls, just the content of the in and out 
values (really the Request and Response objects) change.

```go
type SinglePayload struct {
	InPtr        int64
	InLen        int64
	OutPtr       int64
	OutLen       int64
	ErrPtr       [2]int64
	ErrDetailLen int64
	ErrDetail    int64
}
```

This is the structure that defines how a call is made from WASM-land to the 
kernel or a service implemented on the host system.  That structure has several
pointers in it, like `InPtr` and `OutPtr`, and these are set up by the caller
on the WASM side.  The caller may have any amount of shenanigans going on with
memory in his address space, such as multiple stacks and a running garbage 
collector.  So, the client has to allocate the memory for both the in and out
data (pointed to by the `InPtr` and `OutPtr`) and tell us how much memory is
used (input) or available for writing (output) in the `InLen` and `OutLen`.

The last few fields are for the error code and error detail if any.  Again, all
the memory pointed to be `ErrorDetail` has to be allocated by the caller and the 
caller informs the  kernel about the space allocated through the `Len` field. 
Strings are encoded as a pointer and a length field.

In the future, I hope to take this error handling stuff further and allow
the client side libraries to accept parameters, return values, and signal errors
in a way that is natural for their respective languages.  For example, python
should probably use `raise` to signal errors, not return them as part of the
return from a function.

