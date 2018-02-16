#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..

FAIL=0

echo "Starting glue"
DEBUG=1 ${DIR}/glue --file.config.dir hack/gen-config-yaml/_glue_config --file.secret.dir hack/gen-config-yaml/_glue_secrets &

echo "Starting envoy"
${DIR}/envoy -c ${DIR}/envoy.yaml --v2-config-only --service-cluster envoy --service-node envoy &

trap 'kill $(jobs -p)' EXIT

for job in `jobs -p`
do
echo ${job}
    wait ${job} || let "FAIL+=1"
done

echo ${FAIL} failed
