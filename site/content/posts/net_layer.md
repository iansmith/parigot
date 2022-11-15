---
title: "Networking layers"
date: 2022-11-15T06:29:37-05:00
draft: false
---
After a couple of rewrites of the networking layers to support the
distributed case, I finally have a setup that I can live with for 
some period of time.

There are three layers to the way we handle networking on the server
side:

* Logical layer: This is implemented with go channels, so when waiting for
a call from the network, you use a channel to block on. 
* Netio layer: This layer implements a simple network protocol for sending
and receiving bundles. A bundle is a protobuf-encoded blob of content. This
protocol has a 8 byte preamble and a 4 byte crc at the end.
* Quic layer: This is the networking that actually talks to the wire or its
representatives.  This is implemented primarily through goroutines that
block on channels, the same ones that are visible in the logical layer.

As an example, consider an incoming request, an RPC, from some client 
across the network.  This call is initially fielded by the quic layer 
as a connection, and a "stream" in quic
parlance.  These are accepted indepedently by the quic layer.  Once 
a stream has been established to the calling program,
the quic layer then uses the netio layer to read a bundle from
the caller.  Once the bundle has been read successfully, its contents
are pushed through a channel to the logical layer.  In addition to the
content, also pushed to the logical layer is a response channel.  The quic
layer blocks on that channel waiting for the logical layer to send a 
response.  In both directions, the data (content) is in the form of
an `anypb.Any` that can hold any protobuf type.  Once the quic layer
has received the result through the response channel, it uses the 
netio layer to send a response over the wire.  

To call a remote service, the entire process is the same but with the
directional sense reversed.  In the server case, the quic layer accepted
the remote data then pushed that data to a blocked goroutine in the
logical layer.  In the client case, the quic layer blocks on a channel
and the logical layer pushes an `anypb.Any` through the channel as well
as a response channel. This is the "input" to a remote procedure call.
The quic layer implements the call, and then  responds through the 
response channel to the blocked logical layer.  

At the moment we use channels to implement a serialization of the 
requests on user services, so there is no way for concurrent calls to
user services.  This is likely to be relaxed later, but I'm trying to
get the simple case working first.

For now, we are being pretty aggressive about closing streams that are
established to remote clients or servers.  We are doing this now to try
to insure "sync" across the connection, even though it implies that 
frequently clients must reconnect to the server to make progress.  This,
I hope, will be mitigated in the future by the use of the RTT0 reconnects
that is offered by quic.