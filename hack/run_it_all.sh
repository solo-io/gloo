#!/usr/bin/env bash

set -x -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..

cat << EOF | oc apply -f -
kind: OAuthClient
apiVersion: oauth.openshift.io/v1
metadata:
 name: yuval
secret: yuval
redirectURIs:
 - "http://localhost"
 - "http://localhost:80"
 - "http://localhost:8000"
 - "http://localhost:8082"
grantMethod: prompt
EOF

# minikube should be running
echo RUNNING HYPERGLOO
go run projects/hypergloo/main.go &

echo RUNNING APISERVER
export OAUTH_SERVER="https://$(minishift ip):8443/oauth/authorize"
export OAUTH_CLIENT=yuval

go run projects/apiserver/cmd/main.go &

echo RUNNING UI
docker run -i --rm --net=host soloio/gloo-i:dev &

echo RUNNING Petstore example
docker run -i -p 1234:8080 --rm soloio/petstore-example:latest &

echo "Don't forget to add a static upstream for 127.0.0.1:1234"

## requires latest envoy in solo-kit/..
../envoy -c hack/envoy-sqoop.yaml --disable-hot-restart &
../envoy -c hack/envoy-gateway.yaml --disable-hot-restart &

trap 'sudo kill $(jobs -p)' EXIT

for job in `jobs -p`
do
echo ${job}
    wait ${job} || let "FAIL+=1"
done

echo ${FAIL} failed
