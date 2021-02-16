#!/bin/bash

set -ex

protoc --version

if [ ! -f .gitignore ]; then
  echo "_output" > .gitignore
fi

if [[ $(git status --porcelain | wc -l) -ne 0 ]]; then
  echo "Need to run go mod tidy before committing"
  git diff
  exit 1;
fi

make install-go-tools

set +e

# Ensure that the gloo and solo-apis dependencies are in lockstep
GLOO_VERSION="$(echo $(go list -m github.com/solo-io/gloo) | cut -d' ' -f2)"
go get github.com/solo-io/solo-apis@gloo-$GLOO_VERSION

make generate-all -B  > /dev/null

if [[ $? -ne 0 ]]; then
  echo "Code generation failed"
  exit 1;
fi
if [[ $(git status --porcelain | wc -l) -ne 0 ]]; then
  echo "Generating code produced a non-empty diff."
  echo "Try running 'go mod tidy; make update-all-deps generate-all -B' then re-pushing."
  git status --porcelain
  git diff | cat
  exit 1;
fi