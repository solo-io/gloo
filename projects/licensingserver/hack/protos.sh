#!/usr/bin/env bash

set -ex

API_DIR=/api/v1

cd ${GOPATH}/src/github.com/solo-io/solo-projects/projects/licensingserver

protoc --proto_path=.${API_DIR}  --gogo_out=plugins=grpc:pkg${API_DIR} .${API_DIR}/*.proto

