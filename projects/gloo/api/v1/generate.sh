#!/usr/bin/env bash

GOGO_OUT_FLAG="--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:${GOPATH}/src/"
SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out"

ROOT=${GOPATH}/src/github.com/solo-io/solo-kit/projects/gloo

OUT=${ROOT}/pkg/api/v1/
mkdir -p ${OUT}
protoc -I=./ \
    -I=${GOPATH}/src/github.com/gogo/protobuf/ \
    -I=${ROOT}/ \
    -I=${GOPATH}/src \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}=${OUT} \
    *.proto

# Core Plugin protos
(cd ${ROOT}/pkg/plugins/core/ && \
protoc -I=./ \
    -I=${GOPATH}/src/github.com/gogo/protobuf/ \
    ${GOGO_OUT_FLAG} \
    *.proto)

# AWS Plugin
(cd ${ROOT}/pkg/plugins/aws/ && \
protoc -I=./ \
    -I=${GOPATH}/src/github.com/gogo/protobuf/ \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}=.   \
    *.proto)

# Kubernetes Plugin
(cd ${ROOT}/pkg/plugins/kubernetes/ && \
protoc -I=./ \
    -I=${GOPATH}/src/github.com/gogo/protobuf/ \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/projects/gloo/ \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}=.   \
    *.proto)