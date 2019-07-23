#!/usr/bin/env bash

set -e

cd $SOLO_PROJECTS_DIR

rm -rf _output
export VERSION=dev

./hack/kube/safe_delete_namespace.sh

eval $(minikube docker-env)
make docker -B
make manifest -B

kubectl create namespace gloo-system
kubectl config set-context $(kubectl config current-context) --namespace=gloo-system
kubectl apply -f $SOLO_PROJECTS_DIR/install/manifest/glooe-release.yaml

