#!/bin/bash

set -ex

protoc --version

if [ ! -f .gitignore ]; then
  echo "_output" > .gitignore
fi

git config user.name "bot"
git config user.email "bot@solo.io"

make update-deps


set +e

make generated-code -B  > /dev/null
make generated-ui
if [[ $? -ne 0 ]]; then
  echo "Code generation failed"
  exit 1;
fi
if [[ $(git status --porcelain | wc -l) -ne 0 ]]; then
  echo "Generating code produced a non-empty diff."
  echo "Try running 'dep ensure && make update-deps update-ui-deps generated-code generated-ui -B' then re-pushing."
  git status --porcelain
  git diff | cat
  exit 1;
fi
