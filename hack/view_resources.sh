#!/usr/bin/env bash
for ns in gloo-system; do
    for resource in upstream proxy gateway virtualservice schema resolvermap; do echo "**$resource**"; kubectl get -n $ns $resource; done;
done