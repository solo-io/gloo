#!/usr/bin/env bash

cd $SOLO_PROJECTS_DIR

rm -rf _output
export VERSION=dev

kubectl delete namespace gloo-system

eval $(minikube docker-env)
make docker -B
make manifest -B

kubectl config set-context $(kubectl config current-context) --namespace=gloo-system


# external script for creds
# ask a teammate for the docker password and put this script in your ~/scripts/secret/ dir
kubectl apply -f $SOLO_PROJECTS_DIR/install/gloo-ee.yaml \
 && ~/scripts/secret/docker_credential.sh

