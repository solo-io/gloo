#!/bin/bash

set -ex

# For each argument, delete the kind cluster
for arg in "$@"
do
    kind delete cluster --name "$arg"
done