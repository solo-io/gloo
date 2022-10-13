#!/bin/bash

# Must be run from enclosing directory

set -ex

KIND_CLUSTER_NAME=management-plane

kind create cluster --name ${KIND_CLUSTER_NAME}

kubectl config use-context kind-${KIND_CLUSTER_NAME}

# make all the docker images
# write the output to a temp file so that we can grab the image names out of it
# also ensure we clean up the file once we're done
TEMP_FILE=$(mktemp)
make webhook-test-docker | tee ${TEMP_FILE}
cleanup() {
    echo ">> Removing ${TEMP_FILE}"
    rm ${TEMP_FILE}
}
trap "cleanup" EXIT SIGINT
echo ">> Temporary output file ${TEMP_FILE}"
# grab the image names out of the `make docker` output
sed -nE 's|Successfully tagged (.*$)|\1|p' ${TEMP_FILE} | while read f; do kind load docker-image --name ${KIND_CLUSTER_NAME} $f; done

kubectl create namespace multicluster-admission || true

make deploy-test-chart

kubectl -n  multicluster-admission rollout status deployment multicluster-admission-webhook-test --timeout=1m

