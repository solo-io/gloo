#!/usr/bin/env bash

set -x -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..

# minikube should be running
go run projects/sqoop/cmd/main.go &
go run projects/apiserver/cmd/main.go &
docker run -i --rm --net=host soloio/gloo-i:dev &
docker run -i -p 9091:8080 --rm soloio/petstore-example:latest &

## requires latest envoy in solo-kit/..
sudo ../envoy -c hack/envoy.yaml --disable-hot-restart &

trap 'sudo kill $(jobs -p)' EXIT

for job in `jobs -p`
do
echo ${job}
    wait ${job} || let "FAIL+=1"
done

sudo killall envoy

echo ${FAIL} failed
