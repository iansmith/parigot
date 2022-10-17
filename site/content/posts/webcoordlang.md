---
title: "Web coordination language"
date: 2022-10-17T07:12:19-04:00
draft: false
---

As I thought about the future last night, I had an interesting idea that could result
in a webassembly-based website (in browser) that didn't have the usual size problems.

I began to think about what I wanted, and I was pretty sure of two things I did
not want: typescript and the usual suspects of JS tooling.  This combination results
in the huge masses of code in _node_modules_ and the complete inability to do
partial compilation via make because of things like webpack.  Basically, that whole
ecosystem is a giant pain in the backside.  The resulting binaries are not that
small either.

So I think what I really want is a "web coordination language", a term that I happily
just invented.  This language would combine four things, none of which are terribly
complex, in a nice strictly enforced whole.

1. CSS.  At first, this would simply be a declaration of the css classes that your
library provides  Think of all the classes provided by Bootstrap, for example.
You would define your own names for these rather than use the fully qualified name.
Could be as simple as something like:
 ```
css Bootstrap {
    primary
    light
    dark
    ...
}
```

2.  Named chunks of text.  These are going to be used when there is a large amount
of text to be displayed inside some element.  Further, you could parameterize it
with the browser's language setting to get internationalization.  It might look
like this:
```
text _foo {{ This is my foo, there are many like it but this one is mine. }}
text _bar {{ Now is the time for all good men to come to the aid of their country. }}
```
Clearly these should take parameters as well so you can reuse _most_ of the text
with a small substitution or two.

3. A nice, small way to create segments of HTML, perhaps in s-expression form. For
example:
```
(h2 "first example" (p "hello world") (ul (li "hello") (li "bienvenue") (li "guss gott"))
```
You can use this mechanism with named text chunks like this:
```
(h2 _bar (p "hello world") (ul (li "hello") (li "bienvenue") (li "guss gott"))
```
This notation isn't necessary the final one, but it shows the idea.  These chunks
of html would be named and could combine via names, something like:
```
(:name threehellos (ul (li "hello") (li "bienvenue") (li "guss gott"))
(h2 "first example" (p "hello world") threehellos )
```

4. A simple notation to indicate elements in the html and then an event, or set of
events to handle on that element. You would use something like `#` to define a node's
id so it could referenced later.  The definition of an event should just use
element select expressions ala jquery.

```
(:name threehellos #hello3 (ul (li "hello") (li "bienvenue") (li "guss gott"))
(h2 "first example" (p "hello world") threehellos )

event #hello3 mouse-down {
...
}  
```
5. In terms of the code to respond to that event, the tool would assume that
you wanted some standard type of behavior such as "add css class", "remove class",
"duplicate section of html with name blah and put it at foo", etc.  These would
implemented in AssemblyScript to minimize the size of the result.  If you wanted
some really custom behavior you could have the tool simply spit out the code inside
the braces without interpreting it so you could use it with JS, TS, or AS.

By combining numbers 1-5 in some way, you could create "components" functionality
for your site--without buying into some giant framework or meta-framework...

The key advantage of this is that as a website gets more complex, the interactions
of these parts get more complex to manage in your head and with stupid strings. 
This tool offers a way to cross-check that everything is "lined up": you didn't
reference a non-existent node to put your CSS or actions on. Or you wrote an action
that requires exactly one place to operate, but your code makes copies of the HTML
in question so you probably have trouble.
