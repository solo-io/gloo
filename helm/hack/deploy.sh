#!/bin/bash

cat ../bootstrap.yaml | sed -e "s/{{ .Namespace }}/gloo-system/" | kubectl create -f -

helm install ../gloo --namespace gloo-system -n gloo-demo 