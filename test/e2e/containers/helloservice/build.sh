#!/usr/bin/env bash

set -e

CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o helloservice .

docker build -t solo-io/helloservice:v0.1 . && rm helloservice