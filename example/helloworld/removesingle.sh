#!/bin/bash

# this horrible hack is necessary because there is no way for us to control
# buf and tell it not to generate a particular file.  We have configured
# two code generators into buf, "protoc-gen-go" and "protoc-gen-parigot".
# protoc-gen-parigot is used to generate parigot stubs for a particular proto
# file like foo.proto.  When we do generate the stubs, we *also* need the 
# protoc-gen-go's result, foo.pb.go.  However, we are careful to *not*
# generate stubs that are known system files--these should be imported via
# the parigot library.  However, when we "skip" a dependency such as that,
# we cannot tell buf to NOT run protoc-gen-go, and thus we get a bunch of
# files like queue.pb.go and a queue/v1 directory that contains it.  So we
# use this script AFTER generating stubs to remove these spurious files.
# 
# this is awful.

for dir in g/*
do
	entries="$(ls $dir/v1 | wc -l)"
	if [ $entries == "1" ]; then
		rm -rf $dir
	fi
done