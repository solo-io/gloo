#!/usr/bin/env bash

set -e

cd $SOLO_PROJECTS_DIR

rm -rf _output
export VERSION=dev

./hack/kube/safe_delete_namespace.sh

eval $(minikube docker-env)
make docker -B
make manifest -B

kubectl config set-context $(kubectl config current-context) --namespace=gloo-system


# external script for creds
# ask a teammate for the docker password and put this script in your ~/scripts/secret/ dir
kubectl apply -f $SOLO_PROJECTS_DIR/install/gloo-ee.yaml \
 && ~/scripts/secret/docker_credential.sh

