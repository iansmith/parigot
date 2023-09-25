#!bash

set -e

if [ $# != 3 ] && [ $# != 4 ]; then
	echo not enough command line arguments
	echo usage: $0 \"your app directory\" \"toml deployment file\" \"docker tag for result\" \[ \"parigot source directory\" \]
	exit 1
fi

if [ $# == 3 ]; then
	if [ "$PARIGOT_SOURCE" == "" ]; then
		PARIGOT_SOURCE='.'
	fi
else
	PARIGOT_SOURCE="$4"
fi 

if [ "$KOYEB_API_KEY" == "" ]; then
	echo no environment variable KOYEB_API_KEY found
	exit 1
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

### create a temporary dir
TMPDIR=$(mktemp -d)

### set it up for use by the Dockerfile
mkdir -p $TMPDIR/app/build
BUILDDIR=$TMPDIR/app/build

### get OUR dircetory name (place where this script is)
MYDIR="$(dirname "$(readlink -f "$0")")"

### copy dockerfile to tmp dir
cp $MYDIR/Dockerfile.template $TMPDIR/Dockerfile
### copy arch-independent parts of parigot built-ins to tmp dir
cp $PARIGOT_SOURCE/build/*.p.wasm $BUILDDIR
### copy arch-independent of this particular app
cp $1/build/* $BUILDDIR

### copy in the toml file that was specified on the cmd line
cp "$1/$2" $TMPDIR/app/app.toml

### build it
echo building image with docker: $3, this can take a few minutes
docker build -t "$3" --no-cache $TMPDIR 
