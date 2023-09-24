
#### To build the base image 
1. `make prepdeploybase` from the root of the repo, this populates the
deploy/base with architecture specific binaries.
2. `docker build --progress=plain --no-cache  -f Dockerfile.buildbase .` from this directory which creates the image
3. `docker tag imageid yourusername/parigot-koyeb-base-0.3:arch` where imageid
is the id of the image  created by the build in #2 and arch is one of arm64 and
amd64.

4. `docker push yourusername/parigot-koyeb-base-0.3:arch`

Repeat 1-4 on the "other" arch, thus amd64 if you built the first image
on arm64, and vice versa.

At this point, you should have two images in the docker repository,
`yourusername/parigot-koyeb-base-0.3:arm64` and `yourusername/parigot-koyeb-base-0.3:amd64`

To create the final multi-arch manifest use

`docker manifest create yourusername/parigot-koyeb-base-0.3 \
--amend yourusername/parigot-koyeb-base-0.3:amd64 \
--amend yourusername/parigot-koyeb-base-0.3:arm64`

then push the manifest to the repository,
`docker manifest push yourusername/parigot-koyeb-base-0.3`


#### To build simple deployable
