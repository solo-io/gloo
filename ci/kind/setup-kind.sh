#!/bin/bash -ex

# 0. Assign default values to some of our environment variables
# Get directory this script is located in to access script local files
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
# The name of the kind cluster to deploy to
CLUSTER_NAME="${CLUSTER_NAME:-kind}"
# The version of the Node Docker image to use for booting the cluster
CLUSTER_NODE_VERSION="${CLUSTER_NODE_VERSION:-v1.25.3}"
# The version used to tag images
VERSION="${VERSION:-1.0.0-ci}"
# Whether or not to use fips compliant data plane images
USE_FIPS="${USE_FIPS:-false}"
# If true, use a released chart with the version in $VERSION
FROM_RELEASE="${FROM_RELEASE:-false}"

# Create a kind cluster (or skip creation if a cluster with name=CLUSTER_NAME already exists)
# This config is roughly based on: https://kind.sigs.k8s.io/docs/user/ingress/
function create_kind_cluster_or_skip() {
  activeClusters=$(kind get clusters)

  # if the kind cluster exists already, return
  if [[ "$activeClusters" =~ .*"$CLUSTER_NAME".* ]]; then
    echo "cluster exists, skipping cluster creation"
    return
  fi

  echo "creating cluster ${CLUSTER_NAME}"
  kind create cluster \
  --name "$CLUSTER_NAME" \
  --image "kindest/node:$CLUSTER_NODE_VERSION" \
  --config="$SCRIPT_DIR/cluster.yaml"
  echo "Finished setting up cluster $CLUSTER_NAME"
}
create_kind_cluster_or_skip

if [[ $FROM_RELEASE == "true" ]]; then
  echo "FROM_RELEASE=true: not building docker images, helm chart, or gloo cli"
fi

# Make all the docker images and load them to the kind cluster
VERSION=$VERSION CLUSTER_NAME=$CLUSTER_NAME USE_FIPS=$USE_FIPS make kind-build-and-load -B

# Build the test helm chart, ensuring we have a chart in the `_test` folder
VERSION=$VERSION make build-test-chart

# Build the gloo command line tool, ensuring we have one in the `_output` folder
VERSION=$VERSION make build-cli-local
