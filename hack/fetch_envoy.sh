#!/bin/bash

set -e


ENVOY_GLOO_EE=${ENVOY_GLOO_EE:-$HOME/Projects/solo/envoy-gloo-ee}
ENVOY_GLOO=${ENVOY_GLOO:-$HOME/Projects/solo/envoy-gloo}

# gets you envoy with symbols for commit
COMMIT=$1

if [ -z "$COMMIT" ]; then
  echo usage $0 COMMIT_OR_TAG
  exit 1
fi

# this works for either tag or sha
SHA=$(git --no-pager log -1 --format=format:"%H" $COMMIT)

# get the Dockerfile at that version

ENVOY_GLOO_EE_VERSION=$(git show ${SHA}:cmd/envoyinit/Dockerfile | head -n1 | cut -d: -f2)

ENVOY_GLOO_EE_SHA=$(git ls-remote git@github.com:solo-io/envoy-gloo-ee v$ENVOY_GLOO_EE_VERSION | cut -f 1)

if [ -z "$ENVOY_GLOO_EE_SHA" ]; then
  echo "invalid version. this shouldn't happen"
  exit 1
fi

ENOVY_BINARY=./envoy-$ENVOY_GLOO_EE_VERSION

# get file from envoy gloo-ee
echo envoy-gloo-ee $ENVOY_GLOO_EE_SHA $ENVOY_GLOO_EE_VERSION

if [ -d $ENVOY_GLOO_EE ]; then 
    ENVOY_GLOO_SHA=$( (cd $ENVOY_GLOO_EE; git show ${ENVOY_GLOO_EE_SHA}:bazel/repository_locations.bzl) | python -c "import sys;exec(sys.stdin.read()); print REPOSITORY_LOCATIONS['envoy_gloo']['commit']")
    echo envoy-gloo $ENVOY_GLOO_SHA

    if [ -d $ENVOY_GLOO ]; then 
        ENVOY_SHA=$( (cd $ENVOY_GLOO; git show ${ENVOY_GLOO_SHA}:bazel/repository_locations.bzl) | python -c "import sys;exec(sys.stdin.read()); print REPOSITORY_LOCATIONS['envoy']['commit']")
        echo envoy $ENVOY_SHA
    fi
fi


# get the tag for the remote one

gsutil cp gs://artifacts.solo.io/envoy/$ENVOY_GLOO_EE_SHA/envoy $ENOVY_BINARY
chmod +x $ENOVY_BINARY
