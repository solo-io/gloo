#!/bin/bash
set -e -x
echo echo 
make install-plugin
echo echo
go generate github.com/solo-io/solo-kit/test/mocks
