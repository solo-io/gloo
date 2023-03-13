#!/bin/bash
set -e

curl -sL https://run.solo.io/gloo/install | sh
export PATH=$HOME/.gloo/bin:$PATH

 #todo: Update once branch release is supported
IFS='.' read -ra VERSION <<< "${GLOO_EE_VERSION}"
GLOO_CTL_RELEASE="v${VERSION[0]}.${VERSION[1]}.x"

echo "Set glooctl to $GLOO_CTL_RELEASE"
glooctl upgrade --release="$GLOO_CTL_RELEASE"
glooctl version