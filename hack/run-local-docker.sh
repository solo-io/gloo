#!/usr/bin/env bash

set -x -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..

#DIR=${PWD}
CONFIG_DIR=${CONFIG_DIR:-${DIR}/_gloo_config/}
SECRETS_DIR=${CONFIG_DIR}/secrets/
FILES_DIR=${CONFIG_DIR}/files

mkdir -p ${CONFIG_DIR}/upstreams
mkdir -p ${CONFIG_DIR}/virtualhosts
mkdir -p ${SECRETS_DIR}
mkdir -p ${FILES_DIR}

FAIL=0

echo "Starting control plane"
docker run --rm -i \
    -e DEBUG=1 \
    -v ${CONFIG_DIR}:/config \
    -v ${SECRETS_DIR}:/secrets \
    -v ${files_dir}:/files \
    --name control-plane \
    soloio/control-plane:0.2.1 \
    --file.config.dir /config \
    --file.secret.dir /secrets \
    --file.files.dir /files &

docker run --rm -i \
    -e DEBUG=1 \
    -v ${CONFIG_DIR}:/config \
    -v ${SECRETS_DIR}:/secrets \
    -v ${files_dir}:/files \
    --name function-discovery \
    soloio/function-discovery:0.2.0 \
    --file.config.dir /config \
    --file.secret.dir /secrets \
    --file.files.dir /files &

sleep 1

echo "Starting envoy"
#GLOO_IP=$(docker inspect gloo -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
GLOO_IP=localhost
cat > ${CONFIG_DIR}/envoy.yaml <<EOF
node:
  cluster: ingress
  id: ingress

static_resources:
  clusters:

  - name: xds_cluster
    connect_timeout: 5.000s
    hosts:
    - socket_address:
        address: ${GLOO_IP}
        port_value: 8081
    http2_protocol_options: {}
    type: STRICT_DNS

dynamic_resources:
  ads_config:
    api_type: GRPC
    grpc_services:
    - envoy_grpc: {cluster_name: xds_cluster}
  cds_config:
    ads: {}
  lds_config:
    ads: {}

admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
EOF

#docker run --rm -i \
#    -v ${CONFIG_DIR}:/config \
#    --name envoy \
#    --net=host \
#    soloio/envoy:be0d5f72 \
#    envoy \
#    -c /config/envoy.yaml \
#    --service-cluster envoy \
#    --service-node envoy &

#./envoy -c ./hack/gen-config-yaml/_gloo_config/envoy.yaml --v2-config-only

./envoy -c ${CONFIG_DIR}/envoy.yaml --v2-config-only

sleep 1

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
