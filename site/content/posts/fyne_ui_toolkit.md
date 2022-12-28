---
title: "Fyne UI toolkit and logging"
date: 2022-12-27T15:53:41Z
draft: false
---

Prior to the Christmas holidays, I spent an afternoon learning about the [Fyne
UI toolkit](http://fyne.io). This toolkit's claims to fame are that 1) it is
written in go 2) it expects client code to be writteN in go and 3) it can run
on mobile device as well as desktops. It further claims that it can run "the
same" on windows, linux and mac--but I only tested the mac version. I have not
tried the use of fyne

As UI toolkits go, it is pretty standard in terms of the way it handles most
features. It does not support constraints for doing UI layout, which was kind
of a bummer. For doing layout, it follows the Java idea of doing "layout
managers" although it calls them just "layouts".  Probably the only surprising thing
was that Fyne has an abstraction for pixels that is resolution independent, device
independant, or maybe declaration of independant. It wasn't clear to me that this
mapped nicely to something physical, but maybe I missed that in the docs.  It would
make sense if just said "size in millimeters" and then you would understand what is
happening when you do "SetSize" or use "size" with respect to a font.

After a bit of learning, I built a really simple logging display that runs on the
host system but the `log.Log` service of parigot can talk it over the network.  This
is kind of hacky interface now using TCP, but it does work.  So you can run the
display on your local machine with `go run command/logviewer/main.go` and that will
open a display that shows the log messages as a series of labels.  You can see how
this works in the [New API Mechanism](new_api_interface) post.  

Some fun things to do and/or fix with the log viewer would be:

* Fix the width to be the size the user set with the size of window on the desktop
    rather than using the width of the longest log message.
* Demarcate messages that are from the go side of a service implementation (the 
    "true server" in the language of the New API Mechanism post).  The 
    log implementation for the terminal already does this. Perhaps demarcate with a 
    color or font.
* Demarcate messages that are from the kernel.  The log implementation
    for the terminal already does this. Again, perhaps with a color or font.
* Have a special section for messages from the log viewer itself.  The code already
    distinguishes this.  The most important ones are "new logging client connected"
    and "client disconnected".
* Add UI elements that correspond to filtering choices for the level of log message 
    (DEBUG, WARN, ERROR, etc) you want to see.  Similarly, one could add filtering
    that restricted the display to only a particular log client or (perhaps better)
    service name.

A big effort that is partially in place now is the use of a Call Id.  This id is created
an passed around when services are handling a call.  The entire path from the true
client program all the way through all parts of the implementation should be checked
to insure that are receiving and passing along the Call Id as needed.  If this could
be sent all the way to the logging client, it would give us a way to achieve the dream
of showing the client and server logs unified by the Call Id.  If you did this, you
could probaly have a list of colors like blue, green, yellow, orange, and white that
had no semantic meaning themselves, but caused messages from the same CallId to be
visually similar.

