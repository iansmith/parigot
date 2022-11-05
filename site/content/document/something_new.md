+++
title = "Something New"
description = "White paper about microservices and operating systems."
weight = 2
disableReadmoreNav = true
toc = true
+++

In 1994 a young man from Finland posted some message, somewhere on the internet... 
oh, you know the story by now.  The message in question 
was the public genesis of Linux.  Linux  has spawned some 
multi-billion dollar companies (RedHat, a big chunk of IBM, 
Elastic NV, annoyingly Zynga), forced hardware vendors
to sell pre-installed copies of a free operating system, powers the
servers of nearly every big and small company, and completely
re-invented the notion of embedded computing.  Today, if you are 
building a new small, embedded device, you *start* with something that runs
linux.  It's been ported to every conceivable hardware platform and
operating system--&mdash;and probably some not conceivable ones (Cray).
Linux has completely revolutionized the idea of "having a server," which,
not that long ago, meant you had some box in a rack in your office.
Today, you have neither a box nor an office.

>I come here not to bury linux, but to praise it. 

Linux has proven to have unrivaled flexibilty, reliability, portablity 
and just about every other _-ability_. It has become _de rigeur_ to have 
entire Operating Systems conferences that feature only discussions about
Linux or things based on Linux.  Despite all these plaudits, it's time for
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
before. (Shocker!)  Unix is 50+ years and its prized child Linux is in late 20s.
One could argue, without significant exaggeration, that Open/Close/Read/Write
to small integers has "won" in every domain; the unix idea that "everything
is a file" is manifest in nearly every computer sold today.  Even the only 
big operating system that arguably didn't start
with unix (MS Windows) now exposes a unix-, actally linux-, emulation layer
because, well, FOMO.   The question of whether Mach, and thus "sort-of unix",
is the origin of today's MS Windows is left as an exercise for the reader.

One big thing that has changed is the notion of "files".  A multi-user PDP-11
with a 200MB disk the size of a clothes washer, doesn't have a whole lot in 
common with a single process Docker container on a modern laptop. 
In particular, all of the machinery created to  manage files in a 
multi-user setting is of a little value if you only run one process and have either
one or zero users, depending on how you count.  `chmod`, `chgrp`, `chown`,
and their brethren can quickly become more of a hindrance than a help
in a world of building an deploying microservices.  This is not to
say that there is __no__ value in all the file machinery of Unix, but the
"everyhing is a file" model can be in use and not be nearly as _important_
as befoe.

The second big change we want to point to concerns the way we specify and 
secure software today. We argue that these two activities are two aspects
of the same driving force--with apologies to Maxwell--similar to 
electricity and magnetism.  The driving force to day is the centrality
of the network.  We specify software at the network interface because it
is acknowleges the centrality of the netork without making a choice
about programming languages and programming models.  In this, it is
similar to the "trap to a syscall" model of Unix.  Further, the
ability to secure networks from threats means that in many cases
the "network spec of the software" becomes key to securing the
machines (and users!) that run the software.  In the simplest version,
this can be a firewall with only needed ports open, although the
list of approaches along this line is too large to detail here.
The network certainly
has a more central role in microservice building, testing, and
deploying than files!

We might not be on an island in deciding to deemphasize files in
our operating envirnoment, parigot, but we do feel kinda lonely.
Within the last couple 
of years, the WASM idea (and now standard) emerged.  Yet, what did 
the community
decide to do about IO? [WASI](https://wasi.dev), which is 
Open/Close/Read/Write on small
integers.  No doubt the WASI proponents and implementers will say, 
"if ain't broke, don't fix it."  We say, why not use the open field that 
is WASM to update, maybe even improve,
 the way we write software (esp. microservices) today? 

Our approach is to assume that everything is, and should be, specified at the network
interface; there are "traps" to the system level, the only option is to
make a network call.  The environment 
itself, user programs, third party
code you use, whatever it is, its network interface is what
matters.  To this end, we have a chosen a well-known interface
definition language (IDL), Protocal Buffers as our first cut
at these specifications.  This is less of a slight towards
thrift, hessian, avro and the million other IDLs that exist than just the
quickest way for us to make progress immediately.  We are
not saying it will stay protocol buffers forever.  (That
decision, though, is a slight against CORBA.)

We conclude with a few quotes from a man who provides us inspiration
not only on the way to build microservices, but also one's life.

>> The reasonable man adapts himself to the world: the unreasonable one persists in trying to adapt the world to himself. Therefore all progress depends on the unreasonable man. 

>>Attention and activity lead to mistakes as well as to successes; but a life spent in making mistakes is not only more honorable but more useful than a life spent doing nothing.

>> You see things; you say, 'Why?' But I dream things that never were; and I say 'Why not?

>> After all, the wrong road always leads somewhere.


All by George Bernard Shaw, Irish playwright and winner of the 1925 Nobel 
Prize for Literature.