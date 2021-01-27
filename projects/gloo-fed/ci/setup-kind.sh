#!/bin/bash

set -ex

GLOO_FED_VERSION=$(git describe --tags --dirty --always | sed -e "s/^refs\/tags\///" | cut -c 2-)
# make all the docker images
# write the output to a temp file so that we can grab the image names out of it
# also ensure we clean up the file once we're done
TEMP_FILE=$(mktemp)
make VERSION=${GLOO_FED_VERSION} package-gloo-fed-charts
make VERSION=${GLOO_FED_VERSION} gloofed-docker | tee ${TEMP_FILE}

cleanup() {
    echo ">> Removing ${TEMP_FILE}"
    rm ${TEMP_FILE}
}
trap "cleanup" EXIT SIGINT

echo ">> Temporary output file ${TEMP_FILE}"

if [ "$1" == "" ] || [ "$2" == "" ]; then
  echo "please provide a name for both the master and remote clusters"
  exit 0
fi

if [ "$LICENSE_KEY" == "" ]; then
  echo "please provide a license key"
  exit 0
fi

# Ensure that dependencies are consistent with what's in go.mod.
go mod tidy

# Install glooctl
if which glooctl;
then
    echo "Found glooctl installed already"
    GLOO_VERSION="$(echo $(go list -m github.com/solo-io/gloo) | cut -d' ' -f2)"
    glooctl upgrade --release="$GLOO_VERSION"
else
    echo "Installing glooctl"
    GLOO_VERSION="$(echo $(go list -m github.com/solo-io/gloo) | cut -d' ' -f2)"
    GLOO_VERSION="$GLOO_VERSION" curl -sL https://run.solo.io/gloo/install | sh
    export PATH=$HOME/.gloo/bin:$PATH
    glooctl upgrade --release="$GLOO_VERSION"
fi

glooctl demo federation --license-key=${LICENSE_KEY} --file _output/helm_gloo_fed/gloo-fed-${GLOO_FED_VERSION}.tgz

# grab the image names out of the `make docker` output, load them into kind node
sed -nE 's|(\\x1b\[0m)?Successfully tagged (.*$)|\2|p' ${TEMP_FILE} | while read f; do kind load docker-image --name "$1" $f; done

# wait for setup to be complete
kubectl -n gloo-fed rollout status deployment gloo-fed --timeout=2m
kubectl rollout status deployment echo-blue-deployment --timeout=2m
