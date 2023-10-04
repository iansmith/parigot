
### Not Goto, Microservices Considered Harmful 
### Ian Smith, Oct 2023 
#### Episode 2 

This is the second of a three part series of articles.  The [first of
these](/docs/mserviceone) explains the problems created by microservices in
development.  This article discusses two key problems, and the microservices'
advocates claims about these problems: dependencies and horizontal scaling.  The
final article is explores the cost, both in operations and development, of
typical microservice deployments.

Microservice advocates are quick to claim there are many paths in the software
process that are easier with microservices. Here is an example exerpted from 
the websie of a company that sells a microservice-based product.

1. Microservices Are Easier to Build and Enhance
2. Microservices Are Easier to Deploy
3. Microservices Are Easier to Maintain, Troubleshoot, and Extend
4. Microservices Simplify Cross-Team Coordination
5. Microservices Deliver Performance and Scale

As was promised earlier, we are only going to discuss #4 and #5.  Thus, some
small nits right now: #1) This is clearly not targeted at developers who actually
do the work, because anyone can see that a distributed system has more
complexity than a monolith or similar.  Further, he problem of a distributed
system makes debugging and testing much more difficult (see previous article).
#2) As was mentioned in episode 1, this is true, but only for small,
isolated changes.  Interesting changes typically are *at least as difficult*
with microservices because the deployment advantage is only present when a
single service has to be updated. #3) I fail to comprehend this claim as these
"benefits" are, at best, a wash when a microservice-based system is compared to
a simpler deployment strategy.

For all of these claims, I will clarify something that has been hinted at in
many of my comments about microservices to this point:

> The way you deploy it need not be the same as, or even related to, the way you
> build it.

Suppose one had a "unoservice" architecture that had these properties:

* A single process
* Each service gets an OS-level thread
* Each service uses various `mprotect` tricks to given the illusion of
isolated memory spaces within a single process
* Services communicate through unix pipes

In what way does such a "unoservice" approach differ from one in which each
microservice is deployed as independent processes running on different machines
and communicating over a network?  I would argue that does not differ at all in
terms of *how you build it*.  The difference here is a deployment-time issue.  I
am not saying that deployment doesn't matter, it can be critical in many
organizations.  However, deployment matters to a different set people at a
different time than development-time.  Such a unoservice design is certainly not
harder build than a microservice one, since the abstractions provided are
effectively the same and the topology is the simpler.

So, number 4 in the marketing list above makes the claim that *at even earlier
time* than development-time is somehow (magically?) affected by the a
deployment time decision.  Cross-team coordination occurs before the software is
changed or constructed--in a planning stage.  (If you don't think this must be
true, you may want to discuss with your project manager.)   So, as was said just
above, why should how we deploy matter to the development process?  If that's
the case, how can one argue that things that happen before the software is built
are affected?  Put another way, if one switched from a unoservice to a
microservice deployment (or vice versa), why  should the cross-team
dependencies be mitigated?

While the previous may be logically convincing or not to you, any microservice
developer (not the management or the consultants) will be happy to talk to
you--or rant at you--about the pain of coordinating changes across teams.
Their comments likely can summarized in this simple sentence:

> Pain is conserved.

Experience with the microservice approach has shown that is true that the
problem of coordinating software changes between teams is *moved* but it is not
any less of problem. We will leave as an exercise for the reader if indeed the
coordination problem might be harder. The modularization that a microservice
deployment forces can be a benefit because there are far fewer _hidden_
dependencies between teams.  (In a deployment topology of microservices, there
can be many hidden timing dependencies, so use a race detector!)  It seems hard
to credit the deployment strategy for this benefit, as opposed to the
modularization.  Further, when there are cross-team dependencies, the
discussion/planning when using microservices tends to be about interfaces.
(These are called "contracts" by a few people, and sadly that name is more
descriptive in my opinion.) The software teams, independent of deployment, need
to hash out things like what the new/improved interfaces will look like and how
they will/could/should be used.  How or why should this materially different
than if they were running in a monolith?

If you are careful reader, you will notice I did not mention backward
compatibility. If service A & B are changing their interfaces in response to
desired change, what to do about service C, an already existing customeer of the
(prior) interface to A?  I simply argue that this makes the coordination problem
between A, B, and C--whether the teams behind A & B know about C's situation or
not.

Objections to claim 5 above are difficult to make in a general sesnse.  The
particulars of the performance of a system (distributed or not) depend a great
deal on the details of the system implementation, its goals, and its non-goals.
It may vary even by who is asking questions about the performance or who is
answering them.  It is more than common that different individuals & teams have
different desired outcomes in terms of performance, and these outcomes are
frequently unstated and or in opposition.

We are going to make the simplifying and clarifying assumption that there is a
desire for horizontal scalability.  Horizontal scalability here means that there
is some part of the system where the designers and builders believe that having
multiple copies of a portion of the software will improve performance.  It must
be acknowleged here that it is not a given that such desired copies will fall
nicely along the boundaries of some existing service (or team).  This implies
copies of more than one service when a system is under load.

Typically, the desired for horizontal scalability is because of some well known
part of the software that is or must be "slow".  If responding to a particular
web request will require a highly complex ("slow") computation, it may be the
case that having two or more copies of a section of the softmare (say, an app
server) can improve overall system performance because other copies that are not
doing the expensive computation can handle incoming requests, even while one
instance of the software is calculating a prime number that has 1 million
digits.  

Let's return to our prior thought-expirement the unoservice.  If the desire for
horizontal scaling of a service is one where the service in question is I/O
bound and not CPU or memory bound, then the unoservice can simply spawn more
threads that are running the service in question's code--at least until the
number of cores in the CPU is equal to the number of OS threads. This same
technique can be used if a CPU bound program is not using multiple cores on a
server-class machine.  So, at least in some particular cases, we don't need
microservices to get our desired horizontal scaling; it can still be far
simpler.

In general however, we need a solution that can deal with the situations where
the limiting factor is either the CPU or memory.  In these cases the addition
of more machines, and the accompanying network, makes a significant diffence.
Our unoservice won't help us here so we shall invent the "duoservice".  The
duoservice model is one in which one has a unoservice for all but one of
the microservices, and then a "standard" microservice approach for the
one service that is the bottleneck.  The stardard service-oriented additions
are also welcome here, such as load-balancing and autoscaling, because their 
advantages are now relevant to our overall design.  

Even in the complexity department, this solution is *still* simpler than a
standard microservice approach, although not as much as the pure unoservice.  In
the duoservice design, yes, there is a lot of complexity both for developers and
operations for the single service that must be replicate. However, price of this
complexity in this case is clearly worth paying because the overall system
throughput is much improved.  In standard microservice deployment strategy you
pay this complexity penalty for every service, whether it needs multiple
(horizontal) copies or not.

Many readers I'm sure are asking, "But what about the reliability and redundancy
offered by horizontal scaling?" This is a completely valid concern that I have
delayed discussing until after the duoservice was explained.  In short, the
microservice approach can offer significantly approved availabilty by virtue
of having multiple instances of a any service and thus it is possible for the 
system in total to tolerate some or many services crashing.  (The ability to
handle the failover and restart are being ignored here but these have some
complexity cost.) 

Let us consider the duoservice approach again.  Does this approach yield better
redundancy and reliability? Only for the special service that can be
horizontally scaled as explained above.  So, if you are only concerned with the
reliability of this service, the problem is solved at a substantial complexity
savings for the system in total.  One can simply generalize the duoservice
approach to those services for which the business objectives require the
additional redundancy.  If indeed *all* of your services need this redundancy,
the result is a microservice architecture. Before you decide you "need" this,
you should consider the costs detailed in part 3 of this series.  If you deploy
everything as horizontally scalable because it is truly needed, at least the
microservice complexity is countered by an accompanyingly large business win.

For readers that haven't seen the punch line of part 2 coming at this point, I
will let the Marlon Brando of computer science, Donald Knuth, have the last
word. 

> Premature optimization is the root of all evil.
