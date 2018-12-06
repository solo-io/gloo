#!/usr/bin/env bash 

set -ex

# for symlink compatibility
# https://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done

IN="$( cd -P "$( dirname "$SOURCE" )" >/dev/null && pwd )"

IN=.
OUT=${IN}/../../pkg/api/v1/
mkdir -p $OUT

SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=${PWD}/project.json:${OUT}"

GOGO_OUT_FLAG="--gogo_out=plugins=grpc,"
GOGO_OUT_FLAG+="Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,"
GOGO_OUT_FLAG+="Menvoy/api/v2/discovery.proto=github.com/envoyproxy/go-control-plane/envoy/api/v2,"


PROTOS=""
for proto in ${IN}/*.proto; do
  PROTOS+=" $proto"
done

for proto in ${IN}/plugins/*.proto; do
  PROTOS+=" $proto"
done

for plugindir in $(echo ${IN}/plugins/*/); do
  # remove trailing slash
  plugin=${plugindir%"/"}
  for proto in $plugin/*.proto; do
    PROTOS+=" $proto"
  done
done

GOGO_OUT_FLAG+="paths=source_relative"
GOGO_OUT_FLAG+=":$OUT"

PROTOC_FLAGS="-I=${GOPATH}/src \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/api/external"
    
# generate protos one by one, so protoc doesn't complain
for proto in $PROTOS; do
  protoc ${PROTOC_FLAGS} ${GOGO_OUT_FLAG} -I$IN $proto
done

# generate solokit stuff only for solokit protos
protoc ${PROTOC_FLAGS} ${SOLO_KIT_FLAG} -I$IN $(grep -l solo-kit *.proto)

gofmt -w ${OUT}
goimports -w ${OUT}
