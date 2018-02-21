#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..

FAIL=0

echo "Starting gloo"
DEBUG=1 ${DIR}/gloo --file.config.dir hack/gen-config-yaml/_gloo_config --file.secret.dir hack/gen-config-yaml/_gloo_secrets --kubeconfig $HOME/.kube/config &

echo "Starting envoy"
${DIR}/envoy -c ${DIR}/envoy.yaml --v2-config-only --service-cluster envoy --service-node envoy &

curl localhost:19000/logging?config=debug
curl localhost:19000/logging?router=debug
curl localhost:19000/logging?connection=debug

trap 'kill $(jobs -p)' EXIT

for job in `jobs -p`
do
echo ${job}
    wait ${job} || let "FAIL+=1"
done

echo ${FAIL} failed
