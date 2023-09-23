
** Not Goto, Microservices Considered Harmful **
****** Ian Smith, Sep 2023 ******

This is the first of a three part series of artciles.  The first of these
explains the problems created by microservices in development.  The second
discusses two key problems, and the microservices' advocates claims about these
problems: dependencies and horizontal scaling.  The final article is explores
the cost, both operational and development, of typical microservice deployments.

> Ref: "Go To Statement Considered Harmful", March 1968 issue of Communications 
of the ACM, by Edsger Dijkstra.

***Building a distributed system is hard(er).***

You have two potential designs for the structure of a new product, a single
program and a distributed system.  If you ask a developer of even the most
modest experience which is "easier to build", you'll will certainly get the
answer that the single program is easier. (If you don't, then you have a
different problem.) There are certainly other constraints that need to be
considered when building a new software system other than the ease with which it
can be built, but frequently the ease of building is wrongly not given its
proper respect in the heirarchy of constraints.  Software developer cost is
always the most expensive part of the software build, so not giving ease of
construction the first priority in terms of the prioritization of goals is going
to be more expensive.  Second, the time to build--which is the brother or cousin
of ease of building--is nearly always a key driver of project success.  If it's
easier to build, you get customer feedback (and failure) faster.  If it's easier
to build, you have the opportunity to gain first-mover advantage in the
marketplace.

****Is a microservice architecture always a distributed system?****

Since distributed systems are harder to build, is it fair to characterize any
microservice-based architecture as a distributed system?  I argue that it is.
The problem is that if you remove the distributed system part of anything that
claims to be microservices, you are left with just the (much older) idea of
"modularization" software.   Few, if any, would argue against modularization
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

