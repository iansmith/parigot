# parigot-example/apptempl

This is the  template and scaffolding for building golang 
microservices app using parigot.

#### Setup your own app
* Clone the repo github.com/iansmith/parigot on the master branch.  It must reside in a directory called "parigot", because 
of the use of VSCode dev containers.
* In the new parigot directory, open vscode on the file 
`parigot.code-workspace`.  
    * Vscode will give you the option to 
    open in a dev  container and you should accept that with the dialog in the lower right of your VS code window.
    * Either you will get a VSCode shell that is _inside_ the 
    container or you should create one.
* Do `make` in the shell to build all the parigot code and 
tools.  **If this fails, there is a bug that should be reported.**
* Quit VSCode.
* Now move the directory `parigot` that you created with the clone
operation to another name, we'll use `parigot-src`.
* Make  a new `parigot`` directory (again, the name is required because of the way VSCode handles dev containers) and go into it.
* With VSCode edit the Makefile at the top level to have a correct
path to the PARIGOT_SRC directory, in the case above we would 
probably use `../parigot_src`.


#### The most important files
* `myapp.toml`: deployment descriptor
* `g/`: generated code goes here, contents can be deleted at will.

#### Makefile targets
All Makefile targets are `.PHONY`, they contain no dependencies.
* `make generate`: convert the `proto` files into go source code
    in the g/ directory.
    * You need to use this when you change any `.proto` file.
* `make` builds the source.
    * This uses golang's build tooling, which is exact, so it doesn't hurt to run it anytime.  No changes, nothing done.

#### Run it
* `runner myapp.toml`



