#!/bin/bash

set -ex

protoc --version

if [ ! -f .gitignore ]; then
  echo "_output" > .gitignore
fi

make update-all-deps

set +e

make generate-all -B  > /dev/null

if [[ $? -ne 0 ]]; then
  echo "Code generation failed"
  exit 1;
fi
if [[ $(git status --porcelain | wc -l) -ne 0 ]]; then
  echo "Generating code produced a non-empty diff."
  echo "Try running 'make update-all-deps generate-all -B' then re-pushing."
  git status --porcelain
  git diff | cat
  exit 1;
fi