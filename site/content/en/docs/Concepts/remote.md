+++
title= "Remote Procedure Calls"
description= """Remote Procedure calls are function calls that take place over a network."""
date='2023-07-04'
weight= 3
+++

Remote Procedure Calls (RPCs) have been around since the earliest days of
computer networks, dating from about 1980. It seems natural that one computer
might want to send a request for computation to be accomplished by another
computer and then for the remote computer to send back the result.  These
systems have mostly followed the same basic design.

The design is that a specification is written to specify what Procedure Calls
can be made between the two computers, what parameters the first will send
to the second, and what return result will be sent from the second to the first.
This "specification" requires that the two computers send compatible data formats
between them when making the request and the response.

It should be clear that the [previous section]({{< ref "marshal" >}}) section explained a particular
specification language (the protobuf IDL) and the data format to use when the
communcation channel (the protobuf data serialization format).

The second part of the common design of RPC systems is that the RPC system generates
"boilerplate" code that handles the packing up of parameters, sending them, waiting
for the response, then unpacking the result and presenting back to other layers
of the program in a way that makes sense for the programming language.  This generated,
repetitive code is ofter referred to as the __stubs__.  There is a "stub" for each
method call and response that makes the calling of a Remote Procedure Call look
either mostly or completely like a "normal" procedure call that does not utilize
the network.  

This ability to "hide" the use of the network in an RPC is of arguable value.  Consider
this procedure call where the method `Bar` is being called on the object `Foo` and
the result placed in the variable `result`.

{{< tabpane  right=true >}}
  {{% tab text=true header="Golang" lang="go" highlight=true %}}

	result=Foo.Bar()
	
  {{% /tab %}}
  {{% tab header="Python" lang="python" disabled=true /%}}
  {{% tab header="Java" lang="java" disabled=true /%}}
{{< /tabpane >}}

If this is a normal, non-networked function call that our current processors
can do billions of times per second, the odds of this function truly failing
are very low.  "Failing" here might be a situation such that the program is out
of memory, memory corruption by cosmic rays has been detected, or the processor is
powering down so as to not overheat.  These are all failures, it is true, but outside
the first one these types of failures are so rare that a typical developer may never
see them in their who career.  Running out of memory is not outside the experience
of most developers, but that usally is a catastrophic development for the running
program not one that most programming languages provide much way to do anything
about (the program just crashes).

However, if by the generation of stubs the call in the example above uses a network,
the class of failures that can happen is not only much larger, but the liklihood
of a failure is vastly larger. Some reasonable failures might be:
	* The machine that was expected to do the computation `Bar` on behalf of `Foo` is
currently overloaded and cannot accept the request, although it might be able
to later!
	* The network connecting the caller and receiver of the mesasge about `Bar` are
not connected via a network right now (the network is down, the plug was pulled
out by the dog).
	* The remote machine that would normally do the computation of `Bar` is offline
for maintainence.

... and there are many more.  This is the reason that parigot *always* returns
an error code from any method call, because we do not want to have a situation
where the parigot developer "forgets" that the are using a network when making 
a procedure call.  The case of `Foo.Bar()` above makes it very pleasant to ignore
the networking involved in the computation, until it doesn't.

### parigot and generated code

Previously, we mentioned that RPC systems for 40+ years have been generating
code to make RPC calls easier and more pleasant to use.  We also mentioned that
there are risks with *completely* hiding the network from a user trying to do
what appears to be a simple procedure call.  parigot generates a large amount of
code based on the `.proto` files that specify the interfaces between services in
your system.  parigot tries to strike a balance between convenience of notation
and exposing the multitude ways that a network can fail.  In the case of the
current golang support parigot generates code that is strongly typed such that
the developer must use the correct types when interacting between a caller and
receiver (there are no loopholes in the type system).   Further, parigot is
careful to expose to the user the return value that would be expected and might
include information in it about the details of the network failure. Finally,
when using the continuations sytle of development with parigot, parigot has
strongly typed notions called `Futures` that express that a network call is in
progress and may yet fail.

{{% alert title="Limitation" color="warning" %}}
In the example [previously with the enums]({{< ref "/marshal#enum-protobuf" >}}), the reader may have noticed that
some of the enum values of an error are "mandatory".  These values are
reserved for errors that parigot's generated code calls.  At the current
time, the additional values that need to be added to the "reserved" list
of enum error values are missing.  These additional values are to allow
parigot to propagate network failures it detects to your code.  At the
moment, all the various network errors are conflated with the Marshal and
Unmarshal errors.
{{% /alert %}}
