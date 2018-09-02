#!/usr/bin/env bash 

set -e

# for symlink compatibility
# https://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done

IN="$( cd -P "$( dirname "$SOURCE" )" >/dev/null && pwd )"
OUT=${IN}/../../pkg/api/v1
GOGO_OUT_FLAG="--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:${GOPATH}/src/"
SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=${PWD}/project.json:${OUT}"

mkdir -p ${OUT}/plugins

PROTOC_FLAGS="-I=${GOPATH}/src \
    -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
    -I=${GOPATH}/src/github.com/gogo/protobuf \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}"

protoc -I=${IN} ${PROTOC_FLAGS} ${IN}/*.proto

IN=${IN}/plugins

# protoc made me do it
protoc -I=${IN} ${PROTOC_FLAGS} github.com/solo-io/solo-kit/projects/gloo/api/v1/plugins/service_spec.proto

for plugin in azure aws kubernetes rest transformation; do
mkdir -p ${OUT}/plugins/$plugin

# we need ${GOPATH}/src/github.com/gogo/protobuf/protobuf
# as the filter's protobufs use validate/validate.proto
protoc -I=${IN} ${PROTOC_FLAGS} ${IN}/$plugin/$plugin.proto
done
