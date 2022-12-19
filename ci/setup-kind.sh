#!/bin/bash -ex

# using different architectures
# you can set GOARCH=amd64 to build amd64 images.


# 0. Assign default values to some of our environment variables
# The name of the kind cluster to deploy to
CLUSTER_NAME="${CLUSTER_NAME:-kind}"
# The version of the Node Docker image to use for booting the cluster
CLUSTER_NODE_VERSION="${CLUSTER_NODE_VERSION:-v1.22.4}"
# The version used to tag images
VERSION="${VERSION:-1.0.0-ci}"
# Whether or not to use fips compliant data plane images
USE_FIPS="${USE_FIPS:-false}"
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

# if user is running arm, these are configurations for the registry
REGISTRY_NAME='kind-registry'
# TODO(jake) change to 5001 because this is a better port, does not conflict with apple ear pods
REGISTRY_PORT="${REGISTRY_PORT:-5000}"

# set the image repo to the kind registry endpoint
IMAGE_REPO="${IMAGE_REPO:-}"
if [[ $ARCH == 'arm64' ]]; then
  IMAGE_REPO="${IMAGE_REPO:-localhost:${REGISTRY_PORT}}"
fi

function create_kind_registry() {
  # create registry container unless it already exists
  # 5000 is the port that the container receives images on docker push
  if [ "$(docker inspect -f '{{.State.Running}}' "${REGISTRY_NAME}" 2>/dev/null || true)" != 'true' ]; then
    docker run \
      -d --restart=always -p "127.0.0.1:${REGISTRY_PORT}:5000" --name "${REGISTRY_NAME}" \
      registry:2
  fi
}

function connectKindNetworkToRegistry() {
  # connect the registry to the cluster network if not already connected
  if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${REGISTRY_NAME}")" = 'null' ]; then
    docker network connect "kind" "${REGISTRY_NAME}"
  fi
}

function createKubeRegistryConfigMap() {
  # Document the local registry
  # https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
  cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${REGISTRY_PORT}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF
}

# 1. Create a kind cluster (or skip creation if a cluster with name=CLUSTER_NAME already exists)
# This config is roughly based on: https://kind.sigs.k8s.io/docs/user/ingress/
function create_kind_cluster_or_skip() {
  echo "creating cluster ${CLUSTER_NAME}"

  activeClusters=$(kind get clusters)

  # if the kind cluster exists already, return
  if [[ "$activeClusters" =~ .*"$CLUSTER_NAME".* ]]; then
    echo "cluster exists"
    return
  fi

  # create kind registry with endpoint
  if [[ $ARCH == "arm64" ]]; then
    create_kind_registry
    # 5000 is the port that the container receives images from docker push
    ARM_EXTENSION="containerdConfigPatches:
- |-
  [plugins.\"io.containerd.grpc.v1.cri\".registry.mirrors.\"localhost:${REGISTRY_PORT}\"]
    endpoint = [\"http://${REGISTRY_NAME}:5000\"]
"
  fi

  echo "creating cluster ${CLUSTER_NAME}"

  cat <<EOF | kind create cluster --name "$CLUSTER_NAME" --image "kindest/node:$CLUSTER_NODE_VERSION" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: ClusterConfiguration
  metadata:
    name: config
  apiServer:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  scheduler:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  controllerManager:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: InitConfiguration
  metadata:
    name: config
  nodeRegistration:
    kubeletExtraArgs:
      "feature-gates": "EphemeralContainers=true"
$ARM_EXTENSION
EOF

  # finish setting up registry for arm64
  if [[ $ARCH == "arm64" ]]; then
    connectKindNetworkToRegistry
    createKubeRegistryConfigMap
  fi
  echo "Finished setting up cluster $CLUSTER_NAME"
}
create_kind_cluster_or_skip

if [[ $JUST_KIND == "true" ]]; then
  exit
fi

# 2. Make all the docker images and load them to the kind cluster
if [[ $ARCH == 'arm64' ]]; then
  # if your local machine is arm64, push to the docker registry container, instead of kind
  # GOARCH allows you to support any type of image architecture
  GOARCH=$GOARCH VERSION=$VERSION IMAGE_REPO=${IMAGE_REPO:-} USE_FIPS=$USE_FIPS PUSH_TESTS_ARM=true make docker-push-local-arm -B
else
  VERSION=$VERSION CLUSTER_NAME=$CLUSTER_NAME USE_FIPS=$USE_FIPS make push-kind-images -B
fi

# 3. Build the test helm chart, ensuring we have a chart in the `_test` folder
# setting GOARCH here, because we have to set the helm registry values so that they pick up in the helm registry for arm64 support.
# if this is not done, then helm sets the registry as quay.io, which is only supported for amd64 machines.
GOARCH=$ARCH VERSION=$VERSION IMAGE_REPO=${IMAGE_REPO:-}  RUNNING_REGRESSION_TESTS=true make build-test-chart

# 4. Build the gloo command line tool, ensuring we have one in the `_output` folder
make glooctl-$OS-$GOARCH
