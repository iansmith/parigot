#### The complexity boar.

I read an article today about why we don't ship software as fast as we
used to.

https://every.to/p/why-we-don-t-ship-software-as-fast-as-we-used-to?utm_source=tldrnewsletter

There is nothing in the article I disagree with (he's fairly careful to not
point fingers) but I think it could be strengthed.

I would make the larger claim that we are currently battling,
a complexity boar--and we are losing. No matter how much a few of us stab it with our
steely knives, we just can't kill the beast. (Apologies to Glenn Frey.)  If one
thinks about the backend, middle-end, and front-end together, software
development is probably battling a whole sounder of boars.

I'm not going to try to put too fine a point on it, but there are two easy to 
point to complexity boars (tm).  The first of these is the proliferation and
adoption of front-end frameworks.

Another article published today, coincidentally:
https://matt-rickard.com/why-is-the-frontend-stack-so-complicated?utm_source=tldrwebdev

In a previous job where I managed front end developers I would call this problem
the _"front end framework du jour"_.  That is only a mild exaggeration.   The
number of frameworks, the complexity and requirements of these frameworks, the
use of the npm package manager, and the complexity of the tooling that using so
many frameworks requires of the user is, frankly, silly.  I saw a developer with
20+ years of experience with many different tools, toolkits, programming
languages, etc run screaming away from the front end because of the complexity
and mental load to understand all these pieces-- naturally, which seem to change
frequently.  I'm recovered pretty well from that experience now.

The other complexity problem is microservices.  Microservices convert the 
simplicity of a single program with easy to understand causality into  a
distributed system that is much harder to reason about/program for.  This
complexity of the microservices is allegedly to buy some benefits like 
reliability/redundancy, horizontal scaling, easier development (because each
team only has one responsibility, their service), and easier deployment
for the same reason.  These alleged benefits are probably true for a company
who is already "at scale", whatever that meeans.  

If you are a company that has 100,000 concurrent users or one with extreme SLA requirements
that require no downtime, then the microservice approach probably buys you some things.  However,
the vast majority of businesses are not anywhere near that scale.  A simple
golang backend in a monolith with a single server that costs $1/hr can probably 
handle 60 to 70 thousand requests per second. (See https://medium.com/smsjunk/handling-1-million-requests-per-minute-with-golang-f70ac505fcaa) It is not 
out of the range of reason to suggest that with a bit of tuning that a single,
monolithic server that costs $1/hour could handle 200,000 requests per second.
Wait, if I have 100,000 concurrent users, how many requests per second would they
generate?

Microservices also ignore two key wins of a monolith, especially one that is
decently well built.  Because it is so much simpler to build _and deploy_ the
monolith gives you two *time* advantages.  If you can build it more easily, you'll
be able to get a version of it out to the world quicker.  The benefits first mover
advantage and fast failure/feedback are covered in detail in many hundreds of
books both in the software field and the business field.

I'll close with this.  If you are developer: How much time per week would you
save if you were not fighting the front end frameworks and/or the microservice
architecture on the back end?  You can DM your answers if you don't want them
to be posted publicly in the comments.


