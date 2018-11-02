#!/usr/bin/env bash

set -ex

PROJECTS="$( cd -P "$( dirname "$PROJECTS" )" >/dev/null && pwd )"/../../..

GOGO_OUT_FLAG="--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:${GOPATH}/src/"

# Step 1: Gloo Protos
GLOO_IN=${PROJECTS}/gloo/api/v1/

SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=project_file=${PWD}/project.json,collection_run=true:."
PROTOC_FLAGS="-I=${GOPATH}/src \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/api/external/proto \
    ${SOLO_KIT_FLAG}"

protoc -I=${GLOO_IN} \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/projects/gloo/api/v1 \
    ${PROTOC_FLAGS} \
    ${GLOO_IN}/*.proto

# Step 2: Supergloo Protos
SUPERGLOO_IN=${PROJECTS}/supergloo/api/v1/
OUT=${PROJECTS}/supergloo/pkg/api/v1/

SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=project_file=${PWD}/project.json:${OUT}"

PROTOC_FLAGS="-I=${GOPATH}/src \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/api/external/proto \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}"

mkdir -p ${OUT}
protoc -I=${SUPERGLOO_IN} \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/projects/gloo/api/v1 \
    ${PROTOC_FLAGS} \
    ${SUPERGLOO_IN}/*.proto
