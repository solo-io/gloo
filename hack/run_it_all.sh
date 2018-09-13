#!/usr/bin/env bash

set -x -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..

# minikube should be running
go run projects/sqoop/cmd/main.go &
go run projects/aqpiserver/cmd/main.go &
# requires latest envoy in solo-kit/..
sudo ../envoy -c hack/envoy.yaml --disable-hot-restart &
docker run -ti --rm --net=host soloio/gloo-i:dev &

trap 'kill $(jobs -p)' EXIT

for job in `jobs -p`
do
echo ${job}
    wait ${job} || let "FAIL+=1"
done

echo ${FAIL} failed
