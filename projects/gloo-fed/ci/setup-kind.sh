#!/bin/bash

set -ex

# make all the docker images
# write the output to a temp file so that we can grab the image names out of it
# also ensure we clean up the file once we're done
TEMP_FILE=$(mktemp)
make TAGGED_VERSION=vkind gloofed-docker | tee ${TEMP_FILE}

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

# Install glooctl
GLOO_VERSION="$(echo $(go list -m github.com/solo-io/gloo) | cut -d' ' -f2)"
GLOO_VERSION="$GLOO_VERSION" curl -sL https://run.solo.io/gloo/install | sh
export PATH=$HOME/.gloo/bin:$PATH
glooctl upgrade --release="$GLOO_VERSION"

glooctl demo federation --license-key=${LICENSE_KEY}

# Use solo-projects gloo-fed
kubectl set image -n gloo-fed deployment/gloo-fed gloo-fed=quay.io/solo-io/gloo-fed:kind

# grab the image names out of the `make docker` output, load them into kind node
sed -nE 's|(\\x1b\[0m)?Successfully tagged (.*$)|\2|p' ${TEMP_FILE} | while read f; do kind load docker-image --name "$1" $f; done

# Apply the FederatedRateLimitConfigs until issue https://github.com/solo-io/solo-projects/issues/2027 is resolved
kubectl apply -f - <<EOF
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  labels:
    app: gloo-fed
    app.kubernetes.io/name: gloo-fed
  name: federatedratelimitconfigs.fed.ratelimit.solo.io
spec:
  group: fed.ratelimit.solo.io
  names:
    kind: FederatedRateLimitConfig
    listKind: FederatedRateLimitConfigList
    plural: federatedratelimitconfigs
    singular: federatedratelimitconfig
  scope: Namespaced
  subresources:
    status: {}
  versions:
  - name: v1alpha1
    served: true
    storage: true
EOF

# wait for setup to be complete
kubectl -n gloo-fed rollout status deployment gloo-fed --timeout=2m
kubectl rollout status deployment echo-blue-deployment --timeout=2m

# debug information
glooctl get us
glooctl get vs
kubectl get pods -A
kubectl get kubernetesclusters -A
kubectl get glooinstance -A
kubectl get failoverscheme -Aoyaml
kubectl describe pod echo-blue-deployment
kubectl describe pod echo-green-deployment --context kind-remote