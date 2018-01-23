#!/usr/bin/env bash

set -e

echo "Creating Kuberentes Resources"
kubectl create -f kubernetes/namespace.yml
kubectl config set-context minikube --namespace=glue-system

# glue
kubectl create -f kubernetes/glue-configmap.yml
kubectl create -f kubernetes/glue-deployment.yml
kubectl create -f kubernetes/glue-service.yml

# envoy
kubectl create -f kubernetes/envoy-configmap.yml
kubectl create -f kubernetes/envoy-deployment.yml
kubectl create -f kubernetes/envoy-service.yml

# helloservice
kubectl create -f kubernetes/helloservice-deployment.yml
kubectl create -f kubernetes/helloservice-service.yml

# test runner
kubectl create -f kubernetes/test-runner-pod.yml

for pod in testcontainer helloservice envoy glue; do
    echo "Waiting for all ${pod} to start"
    while [ "$(kubectl get pods | grep ${pod} | awk '{print $3}')" != "Running" ]; do
        sleep 0.01
    done
done

echo "Waiting 10s for envoy to get config"
sleep 10

ginkgo -r -v