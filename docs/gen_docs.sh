#!/usr/bin/env bash


#SOURCES=$(find $GOPATH/src/github.com/solo-io/gloo* -name *.go)
#PKGS=$(echo ${SOURCES} | sort | uniq)
#for i in ${PKGS}; do godocdown pkg/api/types/v1/ > docs/types.md
#for i in $(find $GOPATH/src/github.com/solo-io/gloo* -name *.go | sed "s@$GOPATH/src/github.com/solo-io/@@g"); do echo $i.md; done
#
#
#godocdown pkg/api/types/v1/ > docs/types.md
#

PKGS=$(dirname $(find $GOPATH/src/github.com/solo-io/gloo* -name *.go | grep -v "_test"))
PKGS=$(echo ${PKGS} | sort | uniq)

for i in ${PKGS}; do
 pkg=$(dirname $(echo $i | sed "s@$GOPATH/src/github.com/solo-io/@@g"))
 mkdir -p docs/godoc/$pkg
 set -x
 godocdown $i > docs/godoc/$pkg/$(basename $i).md
 set +x
done