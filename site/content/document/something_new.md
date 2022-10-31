+++
title = "Something New"
description = "White paper about microservices and operating systems."
weight = 2
disableReadmoreNav = true
toc = true
+++

In 1994 a young man from Finland posted some message, somewhere on the internet... oh, you know the story by now.  The message in question 
was the public genesis of Linux.  Linux is now has spawned some 
multi-billion dollar companies (RedHat, a big chunk of IBM, 
Elastic NV, annoyingly Zynga), forced hardware vendors
to sell pre-installed copies of a free operating system, powers the
servers of nearly every big and small company, and completely
re-invented the notion of embedded computing.  Today, if you are 
building a new small device, you *start* with something that runs
linux.  It's been ported to every conceivable hardware platform and
operating system--&mdash;and probably some not conceivable ones (Cray).
Linux has completely revolutionized the idea of "having a server," which,
not that long ago, meant you had some box in a rack in your office.
Today, you have neither a box nor an office.

>I come here not to bury linux, but to praise it. 

Linux has proven to have unrivaled flexibilty, reliability, portablity 
and just about every other _-ability_. It has become _de rigeur_ to have 
entire Operating Systems conferences that feature only discussions about
Linux or things based on Linux.  Despite all the plaudits, it's time for
something new.

The approach we are proposing as "the new thing" comes from our 
experience&ndash;on linux, naturally&ndash;building microservices. This is
by no means to claim "we got it right" or that microservice oriented
system software is "the next big thing&copy;."  Predicting the future 
is known to be hard.  The next big thing in system software might come 
Malaysia (did you guess Finland last time?), might trace its origin
to somebody's Dad's facebook post, or might make a big impact in the 
food and beverage market.  All we can do is say publicly that we are
going to give it a try and we are going to hit the microservices
angle hard.

We have chosen to put our flag in the ground of "newness" in the area of
microservice development for three reasons, not even one of them 
revolutionary.  First and foremost, building microservices is a common
problem.  The number of companies that either have already have changed to
microservice development/deployment or are in the process of doing same is
probably not countable.  The number of those companies' engineers
that are  working on their microservies are also uncountable, but larger.
Our  reasoning here is pure self-interest; one seems more likely to have
an impact when addressing a problem many people have versus a niche.

Second, there is the fact that the world is different now than it was
before.  (How could it
not be?) The posix standard is old enough to have children in high
school (34).  Unix is 50+ years old.  Linux is seeing it's 20s fade into
the review mirror. One
could argue, without significant exaggeration, that Open/Close/Read/Write
to small integers has "won" in every domain; yet we argue that sometimes
blindly following in the time worn path of your elders can make your life 
harder.  We are proposing that the connection between file descriptors and
the technical problems of today is quite tenuous and we can at least 
diminish the importance of these descriptors and their consequences, 
if not dump them entirely.

Within the last couple 
of years, the WASM idea (and now standard) emerged.  Yet, what did 
the community
decide to do about IO? [WASI](https://wasi.dev), which is 
Open/Close/Read/Write on small
integers.  No doubt the WASI proponents and implementers will say, 
"if ain't broke, don't fix it."  We say, why not use the open field that 
is WASM to update, maybe even improve,
 the way we write software (esp. microservices) today? 

It may help to think about the effect (or lack of it) that 
Open/Close/Read/Write have on something as simple as logging in
modern applications.  If you are worried about logging in a microservice
application today, are you concerned about the file descriptors in
use?  (Sockets being the network equivalent of file descriptors is its
own east coast versus west coast story.) Certainly not.  You are
concerned with issues like:
* How do I get my data to our logging service (e.g. Datadog)? What will the delay be if we are trying to fix an urgent problem?
* How do I guarantee causal ourdering of the items logged?
* How do we format our logs (json?) such that our search provider can 
  index them properly for fast access?  Can we guarantee that format is always there in _all_ log messages?
* What should we be doing to reduce the cost of the previous three
  questions? What tradeoffs in our logging can we accept to get lower
  costs?

Do any of these problems sound like ones where Open/Close/Read/Write
is the critical piece? We argue that the network interface (the API and 
the protocol) to logging is vastly more important in those questions 
than file descriptors.
What options should system software provide, given those questions,
to make chosing your point in the trade-offs natural and easy? We intend
to find out if, indeed, those little integers can be omitted.

As an aside, we argue that the entire model of logging in a microservices
needs a new approach. The logging of the past was designed for a single
process that is singly threaded (1970s unix processes).  
Given that, just stuffing lines to the terminal
works great and is quite informative.  In a system with 20 microservices,
each of which is running many threads (goroutines), do you want to 
do the same thing? We argue that any logging approach in a 
microservices-based system has to insure that the logical idea of a 
"request" is preserved in the logs.  We don't want to try to "stitch 
together" the behavior of a given
request based on some strings (request id) in the log results of 20 
different services. Shouldn't one log to the _request_ not the file 
descriptor connected to your terminal or file? The notion of a 
"log message" has to change.

Our third reason for suggesting something new is needed comes from the
fact that we _specify_ software a great deal differently that the 
"old" unix  world did.  Because memory was limited, almost every
program needed to be able to take its input from, and produce its output
to, files.  A good way of specifying what a program did was "it takes this
kind of file in and produces that format of file as output." This was 
so pervasive that Unix reduced user input to a very slow filesystem 
(stdin).

Today, we specify how microservice systems work, as well as what they
should and should not do, based on their network interfaces.  

The analogy continues to work when thinking about the files themselves.
Unix was and
is very concerned about who
could and could not access or write files (chmod, chgrp, chown, and 
ls -l); we are not suggesting that this control is misguided or undesirable! However, today we control access to programs and data at the 
network interface.  We have many, many tools (firewalls, containers) and 
even  programming languages (BPF filtering) that allow us to specify how 
to  control access to programs and data.  We argue that in today's world,
the first priority of specifying a program is at the network interface
and endpoint.  In a world of containerized microservices, it is quite 
common that only a single process running as a single user occupies
the entire "machine", rendering unix's permissions/users/group notions
irrelevant and sometimes even counterproductive.

We conclude with a few quotes from a man who provides us inspiration
not only on the way to build microservices, but also one's life.

>> The reasonable man adapts himself to the world: the unreasonable one persists in trying to adapt the world to himself. Therefore all progress depends on the unreasonable man. 

>>Attention and activity lead to mistakes as well as to successes; but a life spent in making mistakes is not only more honorable but more useful than a life spent doing nothing.

>> You see things; you say, 'Why?' But I dream things that never were; and I say 'Why not?

>> After all, the wrong road always leads somewhere.


All by George Bernard Shaw, Irish playwright and winner of the 1925 Nobel 
Prize for Literature.