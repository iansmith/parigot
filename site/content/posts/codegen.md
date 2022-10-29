---
title: "Codegen"
date: 2022-10-13T03:05:34-04:00
draft: false
---

Yesterday I got the first big part of code generation done.  I feel really good about 
the investment in the "codegen toolkit". I found that yesterday I was able to
add functionality with this toolkit quickly and easily.  I think I could write a
code generator for a new language in an afterneen; this is almost certainly wrong,
since I haven't actually done the porting work yet, but it is a good feeling.

As part of doing that I had to actually deal with the _real_ type system of protobufs.
No more treating types just as a string.  I can now import other proto files from
my proto files as well as use their types and understand what they mean.  This also
is a huge win for the project overall and was certainly going to be necessary.

I worked out last night my plan for `Locators`, although everything ends up changing when
you start running for realsies.  In short, it's a "bootstrap" process
that has a sensible default:
```go
    loc:=net.BootstrapLocate("","","")
service:=loc.Locate(myOrg,myService)
```
Thus, we have only two lines added to the inevitable "hello world" example.  If a user
wants to use their own super-fancy locator they can just put one entry in the 
`BootstrapLocator` that points to their own locator service.

The mass of empty strings in the Bootstrap locator are for non-dev scenarios.  These
represent some things you can pass to the `BootstrapLocator` and as of now are:
```go
    func ... BootstrapLocate(organizationKey string, teamKey string, cluster string)
```

Again, it's very early on, but _roughly_ the ideas are:
* `organizationKey` is issued to you by the owner of the organization.
You can think of organization as the collection of units that are managed for by the
same person.
* `team` is just a layer of hierarchy among the people in the organization.  I didn't
want this to be too heavy so I only allowed one layer.  Teams may not have teams within
them.
* `cluster` is somewhat user defined, but for the default world it will be a collection
of constants like "production" and "staging".  I have some ideas about how to do this for
git branches so developers can easily test things or run integration tests on their own
cluster.

Let's think through a few scenarios. First, you are a solo dev and you sign up
for an account--maybe you are just "kicking the tires".  Signing up for an account
creates an organization with your id as the organization name and issues you an
organization key. Sign up also creates a default team with your name, and issues you a 
key for that team. Finally, sign up creates a cluster called "default".  With this 
in place, in local dev mode you can simply choose the default for 
everything by passing empty strings for all the values.

The only downside of this is that developers will have two keys in their environment
for doing tests that involve talking to something like production.  I didn't view this
as a big deal because most developers understand that this is part of "standard 
operating procedure" and have a scheme for dealing with it.  

The second scenario I thought of was a small development or consulting shop. In this
case the "team" is repurposed for "client".  Thus, the clients are insulated from each
other since they will have only their own key, issued by whoever manages the 
organization.  This is a bit of a headache because developers must keep a set of
keys that correspond one to one with the client projects they work on. 

The third scenario is a corporate environment.  In this case the organization has to
be created by _somebody_ and each team can correspond to a part of the organization
or a building.  Note that in all cases the keys are personal, they are not shared.

Each developer has a different set of `(organization,team)` pairs. 
Since these are individual they are easy to revoke and if there is someone who leaves
the organization, you don't have to worry about them sharing the organization key.
Similarly, if there is a mishap where a key gets checked into version control, you
can just revoke it and issue another.

In both the second and third scenarios, parigot should probably be doing some sort
of "roll up" accounting for each team and each resource.  An example might be number
of service calls made by all the teams' programs.   This account would be nice so 
folks can see which client or department  is using up the most resources within the 
organization.

