+++
title= "Marshaling and Unmarshaling"
description= """Marshaling packs an internal data structure in a program into a \
well-defined external format so the data structure can be given or transmitted \
to another system. Unmarshal does the reverse upon receiving."""
date='2023-07-04'
weight= 2
+++

### Protocol Buffers

Protocol Buffers ("protobufs") is a specification of a data serialization format, an interface
design language (IDL), and a code library.  The first of these is the most important,
with the other two a bit more ancillary.  By interchange format here we mean a
particular byte layout that is carefully speced out in terms of how each data
type should be formatted.  For example, if I want to send an integer value from one
system to another system, what *exactly* should I send? What if the receiver
has 32 bit integers and the sender 64 bit ones?  What if the sender and receiver have
different [Endianness](https://en.wikipedia.org/wiki/Endianness)?  Even this
trivial example is fraught with peril.

Protocol Buffers format has been widely used for 15 years, and for more than 20
years by google.  It is battle-tested.  There are numerous other, typicaly
newer, challengers to the data serialization throne, but none have managed to
disloge protobufs because it is well known, well tested, and reasonably good on nearly all
dimensions of goodness for a data serialization format.  Other challengers include
formats like Thrift, Avro, MessagePack, BSon, Hessian and (god help us) XML.  Many, but
not all of these, like protobufs, have an accompanying IDL to allow users to
specify their data structures of interest.

{{% alert title="Deep Cut" color="info" %}}
parigot is not really tied to the protobuf serialization format.  parigot would
operate exactly the same with another combination of an IDL and a serialization
format. Although protobufs is being used presently, this primarily because of the
wide acceptance of protobufs rather than some feature of it.
{{% /alert %}}

#### Protocol Buffers IDL

Here is a lightly edited example of a pretty trivial "service" definition.  This example
defines a service called `greeting` with a single method called `FetchGreeting`
which naturally takes an input of a `FetchGreetingRequest` and returns a
`FetchGreetingResponse.`

{{< tabpane  right=true >}}
  {{% tab text=true header="Golang" lang="go" highlight=true %}}

	// Greeting is a microservice with a very simple job, return a greeting in
	// language selected from the Tongue enum.
	service Greeting {
		// FetchGreeting returns a greeting in the language given by the
		// Request, field "tongue".
  		rpc FetchGreeting(FetchGreetingRequest) returns (FetchGreetingResponse);
	}
  {{% /tab %}}
  {{% tab header="Python" lang="python" disabled=true /%}}
  {{% tab header="Java" lang="java" disabled=true /%}}
{{< /tabpane >}}

Content like the above would be contained in the file `greeting.proto` or similar.
Although it looks like a programming language, and it is clearly quite similar
to one, this is a "specification","spec", or "schema" in that it only defines the
data to be transmitted and the functions for the data to be transmitted to, as
well as the reverse process for return values.

Let's take a look at the specification of the __messages__ which are the data
objects in a protobuf schema.  In our example, we have two of them, the matching
`FetchGreetingRequest` and `FetchGreetingResponse`.

{{< tabpane  right=true >}}
  {{% tab text=true header="Golang" lang="go" highlight=true %}}

	// FetchGreet is called to retreive a common greeting, like
	// Bonjuor in french.
	message FetchGreetingRequest {
    	Tongue tongue = 1;
	}

	// FetchGreetingResponse is returned to a caller who sent a request
	// to the FetchGreeting endpoint.
	message FetchGreetingResponse {
  	string greeting = 1;
	}                            
  {{% /tab %}}
  {{% tab header="Python" lang="python" disabled=true /%}}
  {{% tab header="Java" lang="java" disabled=true /%}}
{{< /tabpane >}}

So this example is largely what you would expect with the caller requesting the
greeting in a particular language, the "tongue", and the callee returning back
a response that contains the text like __bonjour__ or __guten tag__.

It is worth noticing that the definiton of `FetchGreetingRequest` is not finished
at this point because it references a different "type" called `Tongue`.  Let's
show the last two types.

[#enum-protobuf]
{{< tabpane  right=true >}}
  {{% tab text=true header="Golang" lang="go" highlight=true %}}

	// which language do you want?
	enum Tongue{
		Unspecified = 0;
		English = 1;
		French = 2;
		German = 3;
}

	// The first four values of any error enum are to be as shown below.
	enum GreetErr{
  		option (protosupport.v1.parigot_error) = true;
  	
		NoError = 0; // required
		// Dispatch error occurs when we are trying to call a service
		// implemented elsewhere.  This error indicates that the process
		// of the call itself had problems, not the execution of the
		// service's method.
		DispatchError = 1; // required
		// UnmarshalFailed is used to indicate that in unmarshaling
		// a request or result, the protobuf layer returned an error.
		UnmarshalFailed = 2; // required
		// MarshalFailed is used to indicate that in unmarshaling
		// a request or result, the protobuf layer returned an error.
		MarshalFailed = 3; // required

		// FetchGreeting returns this when the parameter presented to
		// it is not a language in its list.
		UnknownLang = 4;
	}

  {{% /tab %}}
  {{% tab header="Python" lang="python" disabled=true /%}}
  {{% tab header="Java" lang="java" disabled=true /%}}
{{< /tabpane >}}

An __enum__ in a protobuf specification is a collection of small integer values
with names to make them easier to remember and understand when seeing them.  Our
first enum here is a "normal" enum and the second one, `GreetErr` is a special
one for parigot. You will note the extra "option" that is used to inform parigot
about this special type.  A parigot error value is returned from every call to
a remote service.  It is expected that developers will specify all the possible
error values in their error types.

It is worth considering that `FetchGreetingRequests` references another type
as part of its definition--`Tongue` in this case.  What if `Tongue` referenced
one or two more types it __its__ definition?  This is where the protobuf IDL
works in concert with the data serialization format.  Any combination of things
that can be specified in the protobuf IDL will be marshaled to a well-known
sequence of bytes and when unmarshaled by another program will produce the
correct data structure.

{{% alert title="Deep Cut" color="info" %}}

parigot analyzes the definitions in
the protobuf spec and determines if there are request messages or result
messages that do not have any members.  These need to be in the specification
(in case you want to add something later) but are not useful in programs.
parigot removes these parameters from the code it generates so the method
`Bar()` takes no parameters if the corresponding input message has no members.
parigot behaves analogously for output; parigot, however, always returns an
error code.  

{{% /alert %}}
