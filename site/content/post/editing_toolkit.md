---
title: "Editing_toolkit"
date: 2022-10-05T12:17:02-04:00
draft: false
---

# command/transform
I worked for about 3-4 days in an effort to build a binary editing toolkit for WASM. 
Honestly, I was starting to get worried on Monday (day 4) because I had actually wondered
to myself, "Should I try to do this right with an editing toolkit, or just hack some
script together and move on to other, likely more important, thinngs?"  By Monday morning
I was feeling like maybe I had gone the wrong way and I should have just done something
cheep ande cheerful.

The original problem that motivated this effort was that there is a bunch of gunk related
to javascript that is generated into WASM binaries--from both go and tinygo.  I wanted
to either remove this stuff or, better, connect these calls to some kind of "unimplemented
function" error/abort so I could see more easily what they are being used for.  What was
not clear at the start of this effort was that the compiler itself will generate some 
calls into these functions, in much the same way that it might call `malloc` or 
similar. 

Well, the editing toolkit has proven to be a worthwhile investment!  I have nearly 
gotten my first "kernel trap" work today (Wednesday) and I've spent a _lot_ of time
using the toolkit to accomplish things synthesizing functions, changing parameters to
functions, and changing functions to "point to" an unimplemented function implementation
and so forth. Without the toolkit, I would never have been able to do this, and would have
probably either had to use much dumber/simpler test programs at this or just given up!

