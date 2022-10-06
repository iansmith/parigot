---
title: "First_words"
date: 2022-10-05T20:04:47-04:00
draft: false
---

The first words uttered for real by parigot were 
> a=67332 b=5
> done with success!

It may not look like much, but it proves the basic idea that parigot needs to work.
I can make a call from wasm--a "client" program--to a "supervisor" ABI and the parameters reach the 
supervisory layer and it is free to implement the ABI how it sees fit. The implementation
of the supervisor level is not visible and not accessible to the client code.  This is
the equivalent of a kernel trap in Linux.

The two values there are pointer and a length--for a string.  That is how tinygo implements
a string.  However, it's critical to the functioning that the ABI _exposes_ 
`OutputString` but the linked implementation isn't really for a string at all.  The 
implementation inside the supervisor is based on the underlying "hardware" (WASM) model
of a string.  Soon, we'll add some things to check the values for correctness and
so forth.

So, the program that ran on parigot was:
```go
package main

import (
	"github.com/iansmith/parigot/abi/go/abi"
)

func main() {
	abi.OutputString("bleah")
	abi.Exit(107)
}
```
You'll notice that "bleah" is 5 characters, so that is why the b value that was output
above is 5.
