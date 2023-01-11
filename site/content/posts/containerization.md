---
title: "Dev environment options"
date: 2023-01-10T16:26:25Z
draft: false
---

For the benefit of others, especially new folks, I have built four types of development containers:
* raw container
* local for use with VSCode devscontainer on your personal machine
* github codespaces 
* jetbrains spaces

The first one is for the hardcore do-it-yourselfer.  If you want to just have the tooling in a container and
do everything else yourself, _go for it_.  The name of the container is `iansmith/parigot-dev:atlanta`.
You'll have to git clone the repo and mount that into a container yourself.  You'll find the tools you need
in `~/tools`.

The second one is what I use.  I prefer to have everything running on my local machine so I am not dependent on
the internet if I should go to a cafe or be on an airplane.  This dev model is pretty well supported by
VSCode.  When you do `code parigot/parigot.code-workspace` in the repository, it will open vscode and in
the lower right you will see this:
<img src="/vscode-notification-container.jpg" alt="open in devcontainer" width="50%"/>

You can just hit the "open in container" and everything will be set up for you.  I'm not sure if
it remembers the local state, though, in this case.  There is a procedure for moving your dotfiles
into a container when it starts up inside vscode, and lots of other features of devcontainer
[in the docs](https://code.visualstudio.com/docs/devcontainers/containers).

The third and fourth options are for folks that just want to use the internet and don't want to clone
repos or otherwise sully their local disks.  Obviously, these need good internet connectivity.  

The third option, github codespaces, can be launched from the github page of the repo where
you see this:
<img src="/gh-codespaces-selection.jpg" alt="selecting a codespace" width="50%"/>

You have to click on the blue "Code" button on the upper right [of the repo page](github.com/iansmith/parigot)
and then click on the codespaces tab.  The name you see for the codespace on your screen might be
different, but you should take the latest one whatever it is.  You can use the three dots menu to the right
of the codespace name to choose an editor.  You can use one of:
* VSCode running on your machine, but using a container running on github
* Jetbrains IDE like goland, using container running on github
* Web, running in your browser

The latter is obviously the most convenient, but also the most bandwidth intensive.  If you want to use this
option seriously, you probably want to check out 
[chrome keymaps for vscode](https://marketplace.visualstudio.com/items?itemName=erikpukinskis.chrome-codespaces-keymap)
because the hardwired keymaps for chrome (and other browsers) interfere with the "normal" keybindings
of VSCode.  This is highly annoying and I asked someone on the Chrome team to look into it...

Finally, there is the spaces option.  Most people probably don't want this, but the hardcore jetbrains
fanbois might.  This option is triggered when you open the repository inside either the spaces local
desktop application or from the spaces website--which is the same content as the local desktop
application, naturally using electron.


