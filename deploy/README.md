
#### To build the base image 
1. `make prepdeploybase` from the root of the repo, this populates the
deploy/base with architecture specific binaries.
2. `docker build --progress=plain --no-cache  -f Dockerfile.buildbase .` from this directory which creates the image
3. `docker tag imageid your username/parigot-base-build-X.Y:arch` where 
imageid is the id of the image  created by the build in #2 and arch is
one of arm64 and amd64.
4. `docker push username/parigot-base-build-X.Y:arch`

Repeat 1-4 on the "other" arch, thus amd64 if you built the first image
on arm64, and vice versa.

At this point, you should have two images in the docker repository,
`yourusername/parigot-base-build-X.Y:arm64` and `yourusername/parigot-base-build-X.Y:amd64`


The command to build the image, once you have cleaned up the go mod files go to the root of the parigot repo and do:

```
docker build --progress=plain --no-cache  -f Dockerfile.buildbase .
```

After that, you probably want to run docker-squash to get rid of the 
extra layers.  It's usually installed by pip3 in BLAH/bin when you use
a virtualenv of BLAH.

To build the container ready for 
