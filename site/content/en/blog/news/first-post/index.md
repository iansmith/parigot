---
date: 2023-07-04
title: New Website
linkTitle: Announcing New Website
description: We had to do quite a number of changes to get our new website built out and have the automation we wanted.
category: "news"
---

After much gnashing of teeth we have deployed our new website... yay! This
website is powered by [hugo](gohugo.io), a static site generator.  This choice
was made primarily because it allows us to have this website checked into git
and the repository.  An additional bonus is that it takes
[markdown](https://www.markdownguide.org) as its input, which is convenient for
people that spend nearly their whole day in a text editor.

We were using hugo on the old website, but not very well and it was a quite old
version, the one that was obtained by homebrew. Now we've upgraded our hugo
game substantially, and added it to the dev container so that in VSCode you
can expose it with no probs.

We switched to [Docsy](https://www.docsy.dev) as our theme, after two other
failed attempts.  Even when the code for something *seems* like it shouldn't
be that complex--like a website theme--it turns out that the quality, effort, and the knowlege of the people creating it is crucial.  Big hats off to the google
folks that created Docsy.
In particular numerous other thems are so focused on adding the coolest features
that they fail to do the simple things well.  One theme, that shall remain
nameless, was actually so complex that I could not figure out how to __turn it
off__ and had to back the changes out with git. The google folks that did Docsy
have made something that is really nice.

We also switched to two new doc-generation tools:
* [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) for generating
  documentation of [protobuf schemas]({{< ref "docs/reference/api/proto" >}}).
* [gomarkdoc](https://github.com/princjef/gomarkdoc) for generating go 
  documentation from source files. 

Naturally, [usual godoc generator](https://pkg.go.dev/github.com/iansmith/parigot)
is also available.

