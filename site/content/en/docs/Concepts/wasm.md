---
title: WASM
date: 2023-07-04
weight: 1
---

The name WASM, which is pronounced "wahzem", is an acronym of WebASeMbly, if a bit
inelegantly. The name is used to refer to both a file format and a processor
specification.  

### Processor Spec

The details of the [WASM processor
spec](https://webassembly.github.io/spec/core/syntax/instructions.html) is
beyond the scope of this document.  In summary, it is specification for a
virtual CPU, similar to the spec produced by Intel, AMD or other chip vendors
when describing their processor offerings.  If you are familiar with other
virtual hardware specs, like the Java Virtual Machine (JVM) or the Microsoft
Common Language Runtim (CLR), the WASM spec is __far__ lower level than JVM
spec. The WASM spec was designed to make the job of a Just-In-Time compiler easy
compared to other virtual machines.

### File Format

A "WASM file" refers to a file, usually with the extension ".wasm" that
describes a sequence of instructions that run on the processor definition of the
same name.  This is a binary file format, analogous to the `.exe` or `.o` files
used with a "real" processor.  

A WASM file is primarily a sequence of instructions but there are other parts to
a WASM file. The [full spec is
available](https://webassembly.github.io/spec/core/binary/index.html) but for
most people the [WABT toolkit](https://github.com/WebAssembly/wabt) is useful
for exploring and tinkering.  A primary tool in WABT is `wasm2wat` which
converts a WASM file into a (large) textual representation that shows all the
instructions in the file and all the metadata use by loaders, JIT compilers,
linkers, etc.  You can play with this sequence of text s-expressions and then
use `wat2wasm` to convert the file back to binary and run it on a WASM
processor.

Two key points about WASM files are important to parigot users here. First, the
WASM file format is designed to be [dynamically
linked](https://en.wikipedia.org/wiki/Dynamic_linker) and a WASM file
corresponds roughly to a the wasm notion of a __module__. Modules have
interfaces that they export as functions, and similarly they have interfaces
they __import__ (consume) that are not resolved in the WASM binary file.  This
point is critical to how parigot works as parigot exports an interface, to be
imported by your code in WASM, so you can access its functionality.  

Second, as we will see in more detail later, the _binary_ file format of WASM
says nothing about the _origin_ programming language of the program.  The WASM
program you "write" could be written in nearly any programming language and,
assuming the implementation follows the WASM spec correctly, it should work with
parigot.

{{% alert title="Limitation" color="warning" %}}
At the current time, only golang is supported as a guest programming language
for parigot.  There is some amount of effort required to make parigot's functions--
which are independent of language--idomatic and pleasant for use a particular
programming language.  By way of analogy, many programming languages wrap or
adapt `libc.a` on unix systems to make their particular language have access
to the operating system functionality provided by libc, but feel "normal" for
that language's users.
{{% /alert %}}

### WASM Host

At a somewhat naïve level, a __Host__ is a program that simulates the processor
spec WASM such that the Host can run a program specified in a .wasm file.  There
are [many hosts](https://github.com/appcypher/awesome-wasm-runtimes) out there
to choose from.  Hosts have been built for nearly every conceivable dimension of
programming whether that is speed of execution of the WASM program, size of the
host program, ease of installation on tiny hardware, simplicity of the host
program, conformance to the WASM specification or extensions, and startup speed.

parigot requires the use of a host called [wazero](https://wazero.io) which has
chosen the axes of simplicity and a lack of dependencies to optimize on.  
Although a Host can be thought of as an interpreter for the WASM binary
format, in practice most Hosts compile WASM programs to some degree before
or while the program is executing. Wazero is no exception, and does a
complete compilation of the WASM program to native hardware instructions
at the time the WASM program is loaded.  This is trading start-up time off
against simplicity.

{{% alert title="Deep Cut" color="info" %}}
Earlier versions of parigot ran on the [WASMtime](https://wasmtime.dev) Host.
Since then, the Host-specific code in parigot has been isolated and it seems
that "running parigot on top of different Hosts" should be only modestly difficult.
{{% /alert %}}

### Guest

A __Guest__ program, in WASM parlance is an executable sequence of intsructions
for the virtual processor wrapped in the WASM file format.  This allows just
about any compiled programming language
like golang, C, C++, Rust, or Zig to be improved to emit this instruction
type and file format to work with WASM.  Once this "port" has been done,
the programs run on WASM hosts just like they do on a "real" computer.

If your "WASM program" includes the interpreter for your language and
(optionally) the text of your program that will be interpreted, then the language 
does not matter here either for languages like python, javascript, or R.  There is
another type of languages that themselves are based on a virtual machine like
the JVM or CLR and these programs can be transpiled at the virtual instruction
level to work on WASM.  This enables any of the programming languages based on
the JVM like Java, Scala, or Kotlin to work as a WASM programs, plus all the CLR
languages like C# and F#.

{{% alert title="Deep Cut" color="info" %}}
There are many things about the WASM processor specification related to security
and isolation of WASM programs.   While these particulars are important to
the parigot ecosystem as a whole, they are not discussed in this document.
{{% /alert %}}


