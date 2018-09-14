#!/usr/bin/env bash

set -x -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..

# minikube should be running
echo RUNNING SQOOP+GLOO+GATEWAY+DISCOVERY
#go run projects/sqoop/cmd/main.go &

echo RUNNING APISERVER
go run projects/apiserver/cmd/main.go &

echo RUNNING UI
docker run -i --rm --net=host soloio/gloo-i:dev &

echo RUNNING Petstore example
docker run -i -p 1234:8080 --rm soloio/petstore-example:latest &

echo "Don't forget to add a static upstream for 127.0.0.1:1234"

## requires latest envoy in solo-kit/..
../envoy -c hack/envoy.yaml --disable-hot-restart &

trap 'sudo kill $(jobs -p)' EXIT

for job in `jobs -p`
do
echo ${job}
    wait ${job} || let "FAIL+=1"
done

echo ${FAIL} failed
