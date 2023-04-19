#!/bin/bash -ex

# 0. Assign default values to some of our environment variables
# Get directory this script is located in to access script local files
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
# The name of the management kind cluster to deploy to
MANAGEMENT_CLUSTER="${MANAGEMENT_CLUSTER:-management}"
# The name of the remote kind cluster to deploy to
REMOTE_RELEASE_CLUSTER="${REMOTE_RELEASE_CLUSTER_NAME:-remote-release}"
# The name of the remote kind cluster to deploy to
REMOTE_CANARY_CLUSTER="${REMOTE_CANARY_CLUSTER:-remote-canary}"
# The version of the Node Docker image to use for booting the clusters
CLUSTER_NODE_VERSION="${CLUSTER_NODE_VERSION:-v1.25.3}"
# The version used to tag images
VERSION="${VERSION:-1.0.0-ci}"
# The license key used to support enterprise features
GLOO_LICENSE_KEY="${GLOO_LICENSE_KEY:-}"
FROM_RELEASE="${FROM_RELEASE:-false}"
# Automatically (lazily) determine OS type
if [[ $OSTYPE == 'darwin'* ]]; then
  OS='darwin'
else
  OS='linux'
fi

# set the architecture of the machine (checking for arm64 and if not defaulting to amd64)
ARCH="amd64"
if [[ $(uname -m) == "arm64" ]]; then
  ARCH="arm64"
fi

# set the architecture of the images that you will be building, default to the machines architecture
if [[ -z "${GOARCH}" ]]; then
  GOARCH=$ARCH
fi


# 1. Ensure that a license key is provided
if [ "$GLOO_LICENSE_KEY" == "" ]; then
  echo "please provide a license key"
  exit 0
fi

# 2. Build the gloo command line tool, ensuring we have one in the `_output` folder
make glooctl-$OS-$GOARCH -B
shopt -s expand_aliases
alias glooctl=_output/glooctl-$OS-$GOARCH

# 3. Create the kind clusters
# https://kind.sigs.k8s.io/docs/user/configuration/
kind create cluster --name "$MANAGEMENT_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" \
  --config="$SCRIPT_DIR/resources/management-kind-cluster.yaml"

kind create cluster --name "$REMOTE_RELEASE_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" \
  --config="$SCRIPT_DIR/resources/remote-kind-cluster.yaml"

kind create cluster --name "$REMOTE_CANARY_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" \
  --config="$SCRIPT_DIR/resources/remote-kind-cluster.yaml"

if [[ "$FROM_RELEASE" == "true" ]]; then
  echo "skipping build because we will test a released version of gloo"
  exit;
fi
# 4. Build local federation and enterprise images and helm charts used in these clusters
# NOTE TO DEVELOPERS: This build step should only occur once, and we can load the images into separate clusters
VERSION=$VERSION make build-test-chart -B
VERSION=$VERSION make gloo-fed-docker gloo-fed-rbac-validating-webhook-docker -B
VERSION=$VERSION make gloo-ee-docker gloo-ee-envoy-wrapper-docker discovery-ee-docker -B

# 5. Seed the management cluster
kubectl config use-context kind-"$MANAGEMENT_CLUSTER"
CLUSTER_NAME=$MANAGEMENT_CLUSTER VERSION=$VERSION make kind-load-gloo-fed kind-load-gloo-fed-rbac-validating-webhook -B

# 6. Seed the remote-release cluster
kubectl config use-context kind-"$REMOTE_RELEASE_CLUSTER"
CLUSTER_NAME=$REMOTE_RELEASE_CLUSTER VERSION=$VERSION make kind-load-federation-control-plane-images -B

# 7. Seed the remote-canary cluster
kubectl config use-context kind-"$REMOTE_CANARY_CLUSTER"
CLUSTER_NAME=$REMOTE_CANARY_CLUSTER VERSION=$VERSION make kind-load-federation-control-plane-images -B