#!/bin/bash

helm delete --purge gloo-demo
cat ../bootstrap.yaml | sed -e "s/{{ .Namespace }}/gloo-system/" | kubectl delete -f -