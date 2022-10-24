---
title: "Compiler status"
date: 2022-10-17T07:34:13-04:00
draft: false
---
I just refer to my protoc plugin now as "the compiler" because it is doing a lot of
compilerish stuff in compilerish ways.

I managed to get a fairly complete version of the compiler done yesterday.   It required
a second major rewrite, but sometimes you don't know what you want in the code until
you are dissatisfied with what you have.  The compiler correctly handles emitting go
code (the only supported language right now) for the ABI interface, parigot-provided
libraries that are defined in protos, and user created protos.  Doing all this 
correctly required a lot of work to get the type system right.

I implemented an abstraction for a "language" in the compiler to represent the
particulars of the text to emit for the language.  For example, what do variables
look like, camel or snake? How are the arguments to a function declared, etc?  This 
is mostly at the very smallest scale; the things that are usually one or two 
"words" like "x int32" for the declaration of a parameter to a function.

At the next level up, each language gets itself a set of go templates to "drive"
the translation.  This is better for things that are large scale like defining
functions, declaring variables and so forth.  The problem until a day or two ago
was the boundary between these two.  The templates were getting littered with
complicated cases--too complicated to reason about with the ugly notation of go's
text templates.  So I worked a way to have the various parts be driven from a table.
Entries in this table can produce values that are useful to the template, but these
generally do not include control flow.  Thus the template should do high level
control flow, but these variables from the table have the actual text.  Here is
definition of a go function.  The function has the same name as it had in the 
proto file plus an underscore.  The use of the $ notation indicates one of these
variables:
```
func {{.WasmMethodName}}_({{$methParamDeclWasm}}) {{$methodRet}} {
{{- if $needsRet}}
return impl.{{.GetName}}({{$methodCallWasm}})
{{- else }}
impl.{{.GetName}}({{$methodCallWasm}})
{{- end}}
}
```

It is far more manageable with this approach.





