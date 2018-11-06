#!/usr/bin/env bash

set -ex

PROJECTS="$( cd -P "$( dirname "$PROJECTS" )" >/dev/null && pwd )"/../../..

OUT=${PROJECTS}/supergloo/pkg/api/v1/

mkdir -p ${OUT}

GLOO_IN=${PROJECTS}/gloo/api/v1/

SUPERGLOO_IN=${PROJECTS}/supergloo/api/v1/

IMPORTS="-I=${GLOO_IN} \
    -I=${SUPERGLOO_IN} \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/projects/gloo/api/v1 \
    -I=${GOPATH}/src \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/api/external/proto"

# Run protoc once for gogo
GOGO_FLAG="--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:${GOPATH}/src/"

INPUT_PROTOS="${SUPERGLOO_IN}/*.proto"

protoc ${IMPORTS} \
    ${GOGO_FLAG} \
    ${INPUT_PROTOS}

# Run protoc once for solo kit
SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=${PWD}/project.json:${OUT}"

INPUT_PROTOS="${SUPERGLOO_IN}/*.proto ${GLOO_IN}/upstream.proto"

protoc ${IMPORTS} \
    ${SOLO_KIT_FLAG} \
    ${INPUT_PROTOS}

