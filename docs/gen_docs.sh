#!/usr/bin/env bash


for i in $(find $GOPATH/src/github.com/solo-io/gloo* -name *.go | sed "s@$GOPATH/src/github.com/solo-io/@@g"); do echo $i.md; done


godocdown pkg/api/types/v1/ > docs/types.md