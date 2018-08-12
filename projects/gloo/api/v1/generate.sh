#!/usr/bin/env bash

GOGO_OUT_FLAG="--gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:${GOPATH}/src/"
SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out"

OUT=../../pkg/api/v1/
mkdir -p ${OUT}
protoc -I=./ \
    -I=${GOPATH}/src/github.com/gogo/protobuf/ \
    -I=${GOPATH}/src/github.com/solo-io/solo-kit/projects/gloo/ \
    -I=${GOPATH}/src \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}=${OUT} \
    *.proto

# AWS Plugin
cd ././../../pkg/plugins/aws/ && \
protoc -I=./ \
    -I=${GOPATH}/src/github.com/gogo/protobuf/ \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}=. \
    *.proto