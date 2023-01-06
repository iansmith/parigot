---
title: "Big reorg"
date: 2023-01-06T13:35:28Z
draft: false
---
It had to be done. I did a grand reorg of the code over the last two days in my effort to make
testing easier, and add some functional tests.  These tests will test that the core functionality
of RPC works properly.  It's hard to test kernel features with unit tests because the interactions
are so complex.

This reorg had a number of good effects because of simplification, and one sad effect.

* I am going to mothball the `jdepp` dependency analyzer I built.  I'm pretty sad about this, but the
    new Makefile is much simpler than any prior version and is only slighly less accurate than the
    automatically generated ones from jdepp.  The cost here is that targets can be rebuilt now that
    would not have been with a jdepp-generated Makefile. Go compiles sufficiently fast that this is
    a minor annoyance.  Goodbye jdepp, at least for now.
* I got rid of all the "in the path" version names and variants.  Although this likely would have some
    nice effects later, it was adding a lot of complexity.  Without it, I can now have a single `go.mod`
    that applies to the whole repo.  Similarly, the `go.work` is simplied to two lines, one for the 
    parigot system and one for my only current example:
    ```go
    go 1.19

    use (
        .
        ./example/vvv
    )
    ```
    The example is separated so it can have it's own Makefile and packaging to work as an example
    of how to build without doing it in the parigot tree.
* I started work on a new file type that specifies how to run a particular set of wasm services.  I guess
    this will be referred to as a "deployment file".  It's [TOML](https://toml.io/en/) of course.  I'll
    write a post about this in the near future.  It should simplify all the configuration of "runner"
    that was previously necessary.

Courtesy of git, here is the summary of the changes:
```
> git diff --shortstat  master origin/master
 113 files changed, 1906 insertions(+), 2936 deletions(-)
```
