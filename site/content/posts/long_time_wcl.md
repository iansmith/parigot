---
title: "Been a long time; WCL"
date: 2023-04-03T14:45:17-04:00
draft: false
---

To quote The Who from more than 50 years ago:

> It's been a long time since we rock and roll.

Yeah, it's been a long time since I blogged here, which frankly, is pretty poor since I've had my code editor on
this repository for about the last three months straight.  Anyway, I've been super busy with many non-code things,
but let's talk about the big progress.

### WCL (named "wickle")

As I discussed [months ago](./webcoordlang.md), I think the web coordiantion language is a great way to buid up
the wow factor.  That code is nearly done, with only a few more things to be worked out with the "views" that are
part of the new (since the previous post) ability to define a model in protobuf and then *use that model unchanged
inside the WCL.*  So there are now, of course, "models", "views" and "controllers" connected together such that
you should be able to define these in a language neutral way, and just dump models you get from the server to into
the WCL said to have it rendered and event handlers put in place.


