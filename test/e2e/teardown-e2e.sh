#!/usr/bin/env bash

set -e
# Assumes minikube is already running

echo "Destroying Kuberentes Resources"
kubectl delete -f kubernetes/namespace.yml

for pod in testcontainer helloservice envoy glue; do
    echo "Waiting for all ${pod} to terminate"
    while [ "$(kubectl get pods | grep ${pod})" != "" ]; do
        sleep 0.01
    done
done
