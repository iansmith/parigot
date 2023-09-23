#!bash

set -ex

if [ $# != 2 ] && [ $# != 3 ]; then
	echo not enough command line arguments
	echo usage: $0 \"your app directory\" \"toml deployment file\"  \[ \"parigot source directory\" \]
fi

if [ $# == 2 ]; then
	if [ "$PARIGOT_SOURCE" == "" ]; then
		PARIGOT_SOURCE='.'
	fi
else
	PARIGOT_SOURCE="$3"
fi 

if [ "$KOYEB_API_KEY" == "" ]; then
	echo no environment variable KOYEB_API_KEY found
fi

if [ \! -d "$1" ]; then
	echo directory "$1" not found
fi

if [ \! -f "$1/$2" ]; then
	echo directory "$1/$2" not found
fi

if [ \! -d "$PARIGOT_SOURCE" ]; then
	echo directory "$PARIGOT_SOURCE" not found
fi

TMPDIR=$(mktemp -d)
echo "$0"

mkdir -p $TMPDIR/app/build
BUILDDIR=$TMPDIR/app/build

MYDIR="$(dirname "$(readlink -f "$0")")"
echo $MYDIR

cp $MYDIR/Dockerfile.template $TMPDIR/Dockerfile
cp $PARIGOT_SOURCE/build/* $BUILDDIR
cp $1/build/* $BUILDDIR

cp "$1/$2" $TMPDIR/app/app.toml


docker build --no-cache $TMPDIR 
