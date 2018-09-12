#!/usr/bin/env bash
for ns in gloo-system; do
    for resource in upstream proxy gateway virtualservice schema resolvermap; do kubectl delete -n $ns $resource --all; done;
done