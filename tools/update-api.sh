#!/usr/bin/env bash

# copied from here https://github.com/envoyproxy/java-control-plane/blob/ccc28659aa7473233e5f0ab602c690a99f2e23f2/tools/update-api.sh
# and adapted for linux

set -o errexit
set -o pipefail
set -o nounset

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "${__dir}/API_SHAS"

protodir="${__dir}/../api/external/proto"
tmpdir=`mktemp -d 2>/dev/null || mktemp -d -t 'tmpdir'`

# Check if the temp dir was created.
if [[ ! "${tmpdir}" || ! -d "${tmpdir}" ]]; then
  echo "Could not create temp dir"
  exit 1
fi

# Clean up the temp directory that we created.
function cleanup {
  rm -rf "${tmpdir}"
}

# Register the cleanup function to be called on the EXIT signal.
trap cleanup EXIT

pushd "${tmpdir}" >/dev/null

rm -rf "${protodir}"

curl -sL https://github.com/envoyproxy/envoy/archive/${ENVOY_SHA}.tar.gz | tar xz --wildcards '*.proto'
mkdir -p "${protodir}/envoy"
cp -r envoy-*/api/envoy/* "${protodir}/envoy"

curl -sL https://github.com/gogo/protobuf/archive/${GOGOPROTO_SHA}.tar.gz | tar xz --wildcards '*.proto'
mkdir -p "${protodir}/gogoproto"
cp protobuf-*/gogoproto/gogo.proto "${protodir}/gogoproto"
# remove protobuf as it may be reused
rm -rf protobuf

curl -sL https://github.com/protocolbuffers/protobuf/archive/${PROTOBUF_SHA}.tar.gz | tar xz --wildcards '*.proto'
mkdir -p "${protodir}/google/protobuf"
cp protobuf-*/src/google/protobuf/{timestamp,descriptor,duration,struct}.proto "${protodir}/google/protobuf"

curl -sL https://github.com/googleapis/googleapis/archive/${GOOGLEAPIS_SHA}.tar.gz | tar xz --wildcards '*.proto'
mkdir -p "${protodir}/google/api"
mkdir -p "${protodir}/google/rpc"
cp googleapis-*/google/api/annotations.proto googleapis-*/google/api/http.proto "${protodir}/google/api"
cp googleapis-*/google/rpc/status.proto "${protodir}/google/rpc"

curl -sL https://github.com/lyft/protoc-gen-validate/archive/${PGV_GIT_SHA}.tar.gz | tar xz --wildcards '*.proto'
mkdir -p "${protodir}/validate"
cp -r protoc-gen-validate-*/validate/* "${protodir}/validate"

curl -sL https://github.com/census-instrumentation/opencensus-proto/archive/${OPENCENSUS_SHA}.tar.gz | tar xz --wildcards '*.proto'
cp opencensus-proto-*/opencensus/proto/trace/trace.proto "${protodir}"

curl -sL https://github.com/prometheus/client_model/archive/${PROMETHEUS_SHA}.tar.gz | tar xz --wildcards '*.proto'
cp client_model-*/metrics.proto "${protodir}"

popd >/dev/null
