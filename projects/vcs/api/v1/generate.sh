#!/usr/bin/env bash

set -ex

PROJECTS="$( cd -P "$( dirname "$PROJECTS" )" >/dev/null && pwd )"/../../..

GLOO_IN=${PROJECTS}/gloo/api/v1/
SQOOP_IN=${PROJECTS}/sqoop/api/v1/
GATEWAY_IN=${PROJECTS}/gateway/api/v1/
VCS_IN=${PROJECTS}/vcs/api/v1/

OUT=${PROJECTS}/vcs/pkg/api/v1

mkdir -p ${OUT}

IMPORTS="-I=${VCS_IN} \
    -I=${GLOO_IN} \
    -I=${SQOOP_IN} \
    -I=${GATEWAY_IN} \
    -I=${GOPATH}/src \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/api/external"

GOGO_FLAG="--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:${GOPATH}/src/"

SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=${OUT} --solo-kit_opt=${PWD}/project.json"
INPUT_PROTOS="${VCS_IN}/*.proto"

protoc ${IMPORTS} \
    ${GOGO_FLAG} \
    ${SOLO_KIT_FLAG} \
    ${INPUT_PROTOS}