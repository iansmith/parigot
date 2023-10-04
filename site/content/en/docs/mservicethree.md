### Not Goto, Microservices Considered Harmful 
#### Ian Smith, Oct 2023 
##### Episode 3 

This is the third of a three part series of articles.  The 
[first of these](/docs/mserviceone) explains the problems 
created by microservices in
development.  [The second](/docs/mservicetwo) discusses two key problems, and the
microservices' advocates claims about these problems: dependencies and
horizontal scaling.  This, final article covers the cost, both operational and
development, of typical microservice deployments.

Few authors who write about microservices discuss the operational costs of
microservices.  This is not a vast conspiracy on the part of the cloud
providers, but yet one has to assume that the executives at Google Cloud, AWS,
and Azure are well aware of the benefits they have received from this deployment
strategy.  This deployment strategy (neé architecture) suffers from a quite
particular problem: overprovisioning; not something Azure et al are excited to
point to, much less correct.

This overprovisioning of resources, notably cloud servers, when deploying as a
microservice architecture comes from two sources, hopes and dreams.  A "full"
microservice deployment with one server per service has enourmous scaling
possibilities, since each service can be scaled, at minimum, on its own machine
to the limits of that server's resources.  With proper programming (and even
more complexity) services can be scaled and deployed on additional servers--
*if they are needed*.  That last qualification is doing a great deal of lifting.

Every startup or small company wants to grow (and dreams of it!), at least to
being a medium sized company.  Thus, the people putting the microservice
architecture in place, often the full variant, are "planning for the future" and
"making sure we don't paint ourselves into a corner" and so forth.  This sounds
completely reasonable until you think about the cost-savings on cloud servers
that could be had with something like a single server deployment.  For a small
company, the ability save $20K or $40K per year on operational costs might be
something the management should be thinking about trading off against the risk
of "not being able to service all the customers that want to pay us".  Having
more cash now, a well-tested alerting system, and a plan for expanding capacity
quickly can be a good choice compared paying for more cloud resources based on
hope.

Continuing our example of the small company that hopes to get "big", what is the
cost of the additional complexity and developer difficulty when using a
microservice architecture? A simple, single machine monolith is easier to
construct and operate... and what is the time-to-market worth to the firm?
([See episode one for more on this](mserviceone.md)) It's hard to make general
statements about the time to market advantage of doing something simple.  Even
if you believe the simple, monolithic approach has some probability of failing
to scale up as needed, is it still worth it because of faster feedback from
customers and the cash in the bank?  Time-to-market value is hard to quantify,
developer cost is not.
 
If you have a developer making $50K (almost certainly not in the USA!) and you
expect them to work 37 hours per week (some meetings intrude) and 50 weeks a
year, saving an hour per week is worth $1.4K/year.  Obviously, it's triple that
for a developer making $150K.  Don't forget to add about %10 in the EU because
of more vacation!  With even a modest five person engineering group based in the
USA, it's easy to realize $20K/year in savings at only a savings one hour per
week.  I would argue, and I am sure many developers would agree, that a simpler
architecture is worth a few hours per week. If your lower complexity solution
with your 5 person team, *also* is worth $30K by avoiding over-provsioning
suddenly you have created $50K in extra cash! As I stated in episode 1, I am not
naive to the idea that there are cases where other factors become more important
than the complexity of the software and the ease of its construction.  However,
it seems that any prudent manager should be considering time to market,
operational costs, and developer expense when considering a more complex
architecture.

I have been hard on microservices on this article and I believe with good
reasons that I have elucidated as well as one could in a short amount of time.
The truly awful problem of microservices are the people collecting
large consulting fees to promote an architecture that *might* be good for your
business, but said consultant does not fully explain the potential problems.  I
prefer to think of microservices precisely opposite way of such a consultant:
Start with the presumption of doing the simplest thing that could possibly work
and consider the costs of microservices as alternative approach.  Surprisingly
often, ease of development, time to market, and operational costs are enough to
convince decision makers to avoid microservices, especially at the beginning of
a project.  You also might want to ask the developers.

I will finish with this screen capture taken from the website of a fortune 500
company.  Is anything missing?

![screencap](/screencap-mservice.png)


> The life of the dead company is placed in the memory of the still alive 
> companies. 