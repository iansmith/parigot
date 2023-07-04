---
date: 2023-07-04
title: Continuations
linkTitle: Continuations are the first option.
description: Continuations, sometimes called promises or futures, are a way of describing computation. In parigot, they are used to provide a way for a singly-threaded program to both send and receive messages.
category: "technical"
---

### Current continuation usage

The `atlanta-0.3` release uses continuations for providing the illusion
of a multi-threaded program on a single thread.  In parigot the mechanism
is called __futures__.  In particular, there are two kinds of futures
in the current release:

* `Method` futures.  These are futures that indicate that a method call
	to another service is in progress and we are waiting for the result.
	These have a method response type (defined in your protobuf spec) and
	an error code. 
* `Base` futures.  Base futures are used when there is only a simple value
	whose computation we are waiting on, like an int, a bool, or ServiceId.

If you think about it, any method call to another service could take a long
time to complete.  Further, it might _fail_ after a long period of time, say
because the call timed out.  That's where Method futures come in.  When you
call a method on another service, you are returned a `future.Method` that
has been parameterized by the return structure (a response) you expect from
the other service, and an error type that is specific to the service you
are calling. 

Here is a method call that prints out response from a "greeting service" that
has a single method `FetchGreeting` and our request to it is `req`.

{{< tabpane  right=true >}}
  {{% tab text=true header="Golang" lang="go" highlight=true %}}

	// Make the call to the greeting service.
	greetFuture := greetService.FetchGreeting(ctx, req)

	// Handle positive outcome.
	greetFuture.Method.Success(func(response *greeting.FetchGreetingResponse) {
		pcontext.Infof(ctx, "%s, world", response.Greeting)
	})

	//Handle negative outcome.
	greetFuture.Method.Failure(func(err greeting.GreetErr) {
		pcontext.Errorf(ctx, "failed to fetch greeting: %s", greeting.GreetErr_name[int32(err)])
	})

  {{% /tab %}}
  {{% tab header="Python" lang="python" disabled=true /%}}
  {{% tab header="Java" lang="java" disabled=true /%}}
{{< /tabpane >}}

An example of `Base` futures is the `Ready` method that you have to write
for any service of yours.  Base futures have a method called `Handle`
and one named `Set`. `Base` is parameterized by the type of the single
value it does or will contain. `Handle` is used to add your code that will be
executed when the value of the future is known, and the latter one is 
to indicate that the future is finished and we know the value.  Consider this
implementation of ready:

{{< tabpane >}}
  {{% tab text=true header="Golang" lang="go" highlight=true %}}

	func (m *myService) Ready(_ context.Context, _ id.ServiceId) *future.Base[bool] {
		fut := future.NewBase[bool]()
		fut.Set(true)
		return fut
	}

  {{% /tab %}}
  {{% tab header="Python" lang="python" disabled=true /%}}
  {{% tab header="Java" lang="java" disabled=true /%}}
{{< /tabpane >}}

This is the trivial implementation of `Ready` in that it returns a `Base`
future that is not only true but already completed.  When the `Set`
method is called, the future is marked completed.

A perhaps surprising consequence of the working with futures that are
(or might be) completed is that _later_ code that adds code
to be executed when the future's value is known, has its code executed 
immediately.  For example:

{{< tabpane  right=true >}}
  {{% tab text=true header="Golang" lang="go" highlight=true %}}

	func somefunc() {
		// ...
		fut := s.Ready() // shown above
		fut.Handle(func (ok bool){
			if ok {
				// we are ready to do something
			} else {
				// oops, ready call has failed
			}
		})
}
  {{% /tab %}}
  {{% tab header="Python" lang="python" disabled=true /%}}
  {{% tab header="Java" lang="java" disabled=true /%}}
{{< /tabpane >}}

If we consider the example above with the `Ready` function and the example with
`somefunc`, the code in the function literal will be executed immediately when
the call to `fut.Handle` is made.  There is nothing to be gained by waiting to
run the function literal, we already know the outcome!

Why, you may ask, is `Ready` expected by the parigot API to return a future? 
The reason is because it is common for the `Ready` function to perform actions
involving other services.  A good example is to lookup references to other
services that might be wanted later.  Since this lookup could take a long time
and/or fail, the evaluation of the `Ready` function has to be delayed until
we know what happened with the lookup.

### Golang and futures
Although these examples were shown in golang, the same basic structure can be
used in any programming language, provided it can create closures. In golang
(and soon in java) however, the futures mechanism is not really needed, golang
has its own, possibly superior tools for dealing with this issue.

Even though a program is singly threaded, golang can provide the program with
many goroutines.  It is not uncomming to have programs that create thousands
of goroutines, even running on only one thread! (In go, the same technique
is called __fibers__.)  In the case of go, it also has the channels mechanism
to allow the go routines to communicate as well as block waiting on "the
other party" when sending or receiving a message.  We expect to have a second
golang implementation of our API that uses goroutines and channels in the
near future (heh!).

### Why just one thread?
We can't get into all the things we have planned, but we will say that if you
write a program that is single threaded, we plan to offer you more capabilities
than a multi-threaded program.  Although the [WASM spec on multi-threaded
programs](https://github.com/WebAssembly/threads) is not finalized, it is true
that there are some implementations already out there.

parigot does plan to offer the ability to run multi-threaded programs if you
want to do so.  We are planning to implement this feature once the "ground
stabilizes" under the spec.