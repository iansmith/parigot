
To build the base image, you have to make sure that the example programs
do not have "replace" directives that point to /workspaces/parigot.

The command to build the image, once you have cleaned up the go mod files go to the root of the parigot repo and do:

```
docker build --progress=plain --no-cache  -f Dockerfile.buildbase .
```

After that, you probably want to run docker-squash to get rid of the 
extra layers.  It's usually installed by pip3 in BLAH/bin when you use
a virtualenv of BLAH.

To build the container ready for 
