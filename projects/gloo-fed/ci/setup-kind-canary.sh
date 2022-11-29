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
CLUSTER_NODE_VERSION="${CLUSTER_NODE_VERSION:-v1.22.4}"
# The version used to tag images
VERSION="${VERSION:-1.0.0-ci}"
# The license key used to support enterprise features
GLOO_LICENSE_KEY="${GLOO_LICENSE_KEY:-}"
# Automatically (lazily) determine OS type
if [[ $OSTYPE == 'darwin'* ]]; then
  OS='darwin'
else
  OS='linux'
fi

# 1. Ensure that a license key is provided
if [ "$GLOO_LICENSE_KEY" == "" ]; then
  echo "please provide a license key"
  exit 0
fi

# 2. Build the gloo command line tool, ensuring we have one in the `_output` folder
make glooctl-$OS
shopt -s expand_aliases
alias glooctl=_output/glooctl-$OS-amd64

# 3. Create the kind clusters
# https://kind.sigs.k8s.io/docs/user/configuration/
kind create cluster --name "$MANAGEMENT_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" \
  --config="$SCRIPT_DIR/resources/management-kind-cluster.yaml"

kind create cluster --name "$REMOTE_RELEASE_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" \
  --config="$SCRIPT_DIR/resources/remote-kind-cluster.yaml"

kind create cluster --name "$REMOTE_CANARY_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" \
  --config="$SCRIPT_DIR/resources/remote-kind-cluster.yaml"

# 4. Build local federation and enterprise images and helm charts used in these clusters
# NOTE TO DEVELOPERS: This build step should only occur once, and we can load the images into separate clusters
VERSION=$VERSION make build-test-chart
VERSION=$VERSION make gloo-fed-docker gloo-fed-rbac-validating-webhook-docker
VERSION=$VERSION make gloo-ee-docker gloo-ee-envoy-wrapper-docker discovery-ee-docker

# 5. Seed the management cluster
kubectl config use-context kind-"$MANAGEMENT_CLUSTER"
CLUSTER_NAME=$MANAGEMENT_CLUSTER VERSION=$VERSION make kind-load-gloo-fed kind-load-gloo-fed-rbac-validating-webhook

# 6. Seed the remote-release cluster
kubectl config use-context kind-"$REMOTE_RELEASE_CLUSTER"
CLUSTER_NAME=$REMOTE_RELEASE_CLUSTER VERSION=$VERSION make kind-load-gloo-ee kind-load-gloo-ee-envoy-wrapper kind-load-discovery-ee

# 7. Seed the remote-canary cluster
kubectl config use-context kind-"$REMOTE_CANARY_CLUSTER"
CLUSTER_NAME=$REMOTE_CANARY_CLUSTER VERSION=$VERSION make kind-load-gloo-ee kind-load-gloo-ee-envoy-wrapper kind-load-discovery-ee