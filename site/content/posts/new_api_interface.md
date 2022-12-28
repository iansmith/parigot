---
title: "New API Mechanism"
date: 2022-12-27T15:02:30Z
draft: false
---

I had spent some time before Christmas doing a "standardish" interface to dropping from user code, in WASM
into system code, in go.  This was done primarily to support some future work by the S3 team.  This
interface is designed to function much like an RPC call with protobufs... in fact, it is quite similar to 
an RPC call with protobufs, just using memory instead of a wire.

For clarity, it is worth saying that this is not _exactly_ the same thing as an RPC call that you would
make on another service.  This interface is for folks that want to implement backends in go that will 
expose some service to the **parigot** and WASM side of the house.  An example here might be that you
want to expose a service like Redis.  Redis has a go client side, but not a **parigot** client side.  So,
if you want to expose Redis in **parigot** you need to implement the "back side" of the service as go
code running natively on the _true_ platform. The caller of your new service can treat this as a
normal **parigot** service--and it should have a API like one--but in reality most of the heavy lifting
is delegated to some go code that is available already to talk to the Redis server.

This new interface has a simplified structure so any service that wants to use go code for its backend
can do so easily.  I have already built the `log.Log` and `file.File` service using this mechanism so 
folks will have examples.  This new API, like an RPC, hinges on have a protobuf-based request and 
response message.  The true client running in WASM creates a normal request on (say) the Redis service
and that service responds as normal, as if the server was also running in WASM.  However, in fact it is
receiving the request from the true client in WASM code, but then it "drops through" to the native side
which is written in go which can access _anything_ on the host system.  (In a production environment the
go code is sandboxed so it can't do harmful things to the host system, but it certainly could open a
network connection to a redis server.) The "drops through" is the critical bit that previously was not 
simple to implement.

So, we have three bits of code to consider:
1. The true client code, somewhere in the world of microservices written using **parigot** and WASM.
2. The fake server code, somewhere in the world of microservices written using **parigot** and WASM 
    to receive the call from the true client code.  It then becomes the client of the true server.
3. The true server code written in go that is linked into the binary of the WASM engine and thus has
    none of the constraints of WASM code.

The interface from (2) to (3) is what is being discussed here.  In a perfect world, the fake server 
code (2) could simply take the request (unmodified) that it received from the true client (1) and pass it
to the true server code (3).  Then, the fake server (2) code would receive a response from the true 
server (3)  code and it could respond back to the true client (1) the response it received from the 
true server (3). In this perfect world, the job of the fake server is literally just a middle man 
that hands what it gets from (1) to (3) and from (3) to (1).  This real world, though, is not perfect.

In practice, the fake server will likely have to modify or add parameters to the message it receives from
the true client before talking to the true server.  This is because the fake server may have its 
own internal data structures that are not visible to the true client but are needed by calls to the 
true server.  An example of this might be a Redis connection id. Similarly, the fake server may not 
want to expose all of the parts of the reponse, or may want to add to the response,  from the true server
to the true client. A good example here might be information about errors that a
true server sends back to a fake server; the fake server may want to decorate its response to the 
true client with extra information about the error, making the error more friendly for the true client.
Analagously, the fake server may want to turn a multitude of errors from the true server into an error
to the true client like "not available right now, please try again and the request might succeed".

With this new machinery in place, I realized that the worst example of the complexity between a fake and
true server was **parigot's** kernel itself.  Before I realize how much I could lean into protobufs, I 
had spent a lot of time and effort trying to (badly) emulate the job of protobuf's serialization
and deserialization.  I was doing a simple version of "convert this to some bytes" and "pull the result
out of the bytes" and it was very complex and had tons of tricky things embedded in it.

So, as of today, I'm trying to rip out that whole pile of cruft and replace it with exactly the same
mechanism used by `log.Log` and `file.File` to "drop through" to go code.  This surely will make kernel
changes less difficult for others, since it follows a pattern for which there are several examples.
When this is completed, there will be a major upgrade in the "once and only once" score for
**parigot** itself!


