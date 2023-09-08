---
title: Getting Started
description: Tooling, editors, etc
weight: 2
---

{{% pageinfo %}}

This page relates to the tooling you need for the `atlanta-0.3` release (pronounced 'at-lan-ta-oh-three').

{{% /pageinfo %}}

### Tooling

This project makes heavy use of __dev containers__ which is a relatively new
tactic for what is an age-old problem.  A dev container is a
[container](https://www.docker.com/resources/what-container/) whose image (files,
tools, and environment) is already packaged so that development work can proceed
without needing a lot of "set up your machine" processes and time.  

With [VSCode](https://www.google.com/search?client=safari&rls=en&q=download+vscode&ie=UTF-8&oe=UTF-8) the container is actually run by the editor (really "dev
environment" or "IDE") and the editor knows how to work with the dev container
for doing things like running programs, launching the debugger, using the
shell *inside* the container.

### What to try

Here is a [video](https://youtu.be/gJEEHfl-n6I) that can show you some things 
to try in the dev container, like running the hello world example.


### Github codespaces

You can launch a complete editor and working shell in your browser with parigot's 
tooling already configured.  You can do this with the `codespaces` tab under
the code button on the github repository.

[Example screencap](/codespaces-scap.png)

Althought this is convenient, the machines that back codespace are quite slow, 
and it can take several minutes to do the launch of a codespace.

### Github classic way

* Clone the repo (`http://github.com/iansmith/parigot`) using git
* Launch vscode on our configured workspace, `parigot.code-workspace` at the
root of the repo.  If you have vscode installed in the normal way, you
can use `code parigot.code-workspace` at the root of the repo.
* When you get the popup in the lower right of vscode that says "reopen in
dev container" click "Reopen in Container" or other affirmative button.
* You will be given your copy of the repo's code, plus a shell that is
pre-configured to work with parigot, as in the video above.

[Screencap of VS Code startup](/vscode-scap.png)

### Other editors

First, VSCode is the recommended way to edit code in parigot simply because it
is so much easier.

If you are wanting to use another editor that understands containers, the
Dockerfile that builds the dev container is in `.devcontainer/Dockerfile` and
should build on your local copy of Docker. From this point, you'll have a properly
set up image you can use in your editor.

In the past, we have demonstrated that this method works with [goland](https://www.jetbrains.com/go/) and  that one can have an environment much like the VSCode one.

