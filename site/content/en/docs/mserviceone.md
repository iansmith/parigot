
### Not Goto, Microservices Considered Harmful 
#### Ian Smith, Oct 2023 
##### Episode 1 

This is the first of a three part series of articles.  The first of these
explains the problems created by microservices in development.  The second
discusses two key problems, and the microservices' advocates claims about these
problems: dependencies and horizontal scaling.  The final article is explores
the cost, both operational and development, of typical microservice deployments.

> Ref: "Go To Statement Considered Harmful", March 1968 issue of Communications 
of the ACM, by Edsger Dijkstra.

***Building a distributed system is hard(er).***

You have two potential designs for the structure of a new product, a single
program and a distributed system.  If you ask a developer of even the most
modest experience which is "easier to build", you will certainly get the answer
that the single program is easier. (If you don't, then you have a different
problem.) There are certainly other constraints that need to be considered when
building a new software system other than the ease with which it can be built,
but frequently the ease of building is wrongly not given its proper respect in
the heirarchy of constraints.  Software developer cost is always the most
expensive part of the software build, so not giving ease of construction the
primary importance in terms of the prioritization of goals is going to be more
expensive.  Second, the time to build--which is the brother or cousin of ease of
building--is nearly always a key driver of project success.  If it's easier to
build, you get customer feedback (and failure) faster.  If it's easier to build,
you have the opportunity to gain first-mover advantage in the marketplace.

****Is a microservice architecture always a distributed system?****

Since distributed systems are harder to build, is it fair to characterize any
microservice-based architecture as a distributed system?  I argue that it is.
The problem is that if you remove the distributed system part of anything that
claims to be microservices, you are left with just the much older idea of
"modularization".   Few, if any, would argue against modularization
in software, so the term "microservice" must imply a particular, distributed 
type of modularization, or there is no point in such a word.

***The complexity disease and its effects on the body software***

There is not sufficient space in a short article to include all the problems
that complexity causes in software building.  There are dozens of books on the
topic from authors both estemeed and not.  Internet rants about complexity are
rife and a google search reveals many thousands of such articles
and discussions\[sic\]. For our purpose we'll limit ourselves to three,
easy-to-describe ailments that flow from the complexity disease.  These are
debugging, failure handling (or not), and testing.

****Symptom:Debugging****

All of the debugging problems flow from the distributed nature of the
microservices architecture.  The most obvious problem is the name of this
symptom: debuggers.  Modern debuggers are designed to work on a single program.
Although it's possible to run several copies of a debugger on the different
microservices that you need to debug, this is a big hassle and likely not
pssible for larger microservice deployments.  If you have 300 microservices
running, which ones do you run the debugger on?  How would one determine what
those are, due to the modularity (concealment) of the microservices'
dependencies?   Even if you can get debuggers running in the needed places, the
debuggers are likely to have a substantial effect on the ordering of operations
when considering the distributed system as a whole, since a debugger changes the
CPU profile of an application.

The usual step before trying to run multiple debuggers is to use logs.  If you
consider multiple log files being generated on different machines,  the
difficulty is obvious. These have to be collected and then "merged in your
head" or with tools to see the complete picture of what is happening.  With
clock drift between machines, it's possible that this "merge" cannot be done at
all, especially on the timescales that modern servers operate on.  What is the
cost of not knowing the causal ordering of different services?  With the use
of centralized logging tool of some type, this problem can be largely be
obviated.  However, in practice several well-known tech companies do not use
such a tool/system _for developers._  The operational complexity and cost seems
to cause some companies to spend the money on a centralized logging system only
for production.  While this is certainly a great thing for finding production
problems, it's little help with debugging during normal development.

To help address this problem, there are numerous vendors and open source projects
trying to increase visibilty into the complete path that a given request follows
during its lifecycle.  This is clearly a major benefit, but again it is frequently
not deployed in organzations _for developers_.  In both this case and centralized
logging, one can make a reasonable argument that the money/effort spent on both
centralized logging and request path visibility would be better spent on developers
rather than production.

****Symptom:Failure Handling (or not)****

Broadly speaking, failure/error handling in a distributed system is known to be
hard. Like, predicting the stock market kind of hard!  For those not familiar
with this problem, it hinges on a quite simple fact: The (separate) program that
notices/receives the failure may not (does not) have the global information
necessary to make a good decision about error handling.  If you think about a
caller and callee where a request flows from caller to callee, and the callee is
seriously impaired, making a "good" decision about errors that occur in the
broken receives is difficult.  This difficulty comes from the fact the broken
service doesn't have even have knowlege of the caller, much less global
knowlege.  Even if the error is simply propagated back  to the caller, the same
issue results, with caller now being the callee of some other service.

The "right thing to do" about many errors requires global knowlege about the
state of the system (to say nothing of the user!)  The simplest possible example
here the difference between product and testing; in one case a system failure is
expected/desired, in the other case it is not.  In fairness, there are some 
languages (go) and frameworks that try to provide some amount of context to
the callee, but these are not use anywhere near enough to be considered standard.

****Symptom:Testing****

Developers love tests.  These are the "backstop" that allows a developer to have
confidence in changes they are making.  In a microservice architecture, testing
becomes vastly more difficult, and many organizations abandon testing in areas
that they would never allow in a non-microservice architecture.  

As it is indicative of related choices organizations make, we will digress
briefly into why organizations give up on testing when using microservices.
The primary driver of this behavior is that as the number of microservices gets
large, ever larger efforts are required to give developers a "reasonable
facisimile" of the production environment.  I spoke to two major tech companies
who did *not* make this choice, and both admitted that they had an entire team
dedicated to the issue of making environments available to developers that were
very similar to the production system.  For those without the resources to spend
on multiple engineering heads on development tooling, things look far grimmer.
It is not uncommon for even mid-sized technology companies to have no easy way
to create a complete copy of the production system for a developers' personal
use.  Many of these types of organizions (several of which I spoke to
personally) had various "dodges" to avoid needing complete copies of the running
system.  The two most common dodges were "partial production" or "shared staging
environment."  In the former case the developer declares that some services
(usually on the scale of 3-5) are to be run locally and the remainder are to use
the existing production system.  The latter is simply the creation of a fixed
number of copies of the production system, and a plausible data set, that
developers can/must use for their testing. The coordination problem with this
latter approach here is obvious.

Returning to testing more directly, the problems above with respect to logging
in a distributed system reappear here.  Because of the network latencies,
scheduler choices, and generally just randomness, a distributed system is hard
to get into a "particular state" for a test.  This is true of nearly any test
where multiple services are involved, because the interactions of services
matter.  Getting the world into a particular state is critical to the very
concept of testing, particularly when the world state is a bad one. If you used
the "partial production" approach above, are you going to be willing to force
production services into known bad states so that a developer can run a test?
No.   The worst outcome of trying to test a microservice based system can be
heard in many engineering offices, "Oh, try running the tests again."  This
means that the developers no longer believe in the test results they are
receiving, which is the dark and direct route to system failures.

****Bonus Symptom:Startup****

For most developers, the way to create a workable integration test structure is
to start up a fresh copy of the "system" (see above) and run specific requests
to test the results.  The key here is "fresh system".  Because systems typically
start in a well-known and predictable state, that is the easiest way to control
things for a test.  However, if your system has, say 20 or 30 services that need
to be started, variance in the startup timing can be significant.  If you are
trying to test service A, how are you to know that service A's behavior depends
materially on whether B is "already up" when A starts?   In effect, your test
now has a hidden dependency on the startup ordering.  A distributed system can
be 71% up, a single program cannot.

With the needed apologies to The Bard and Marc Anthony, I have written this article to bury
microservices, not to praise them.  For a developer, the symptoms above are
sufficient to make him/her question if a new software project **really** wants
to start with the complexity disease.   The disease may be something that is worth
tolerating for some specific benefits--but it's not something that a couple of
aspirin and a morning phone call will fix.

> Brutus was a microservice advocate,
> Brutus was an honorable man.