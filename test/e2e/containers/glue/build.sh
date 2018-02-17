#!/usr/bin/env bash

set -e

CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o gloo .

docker build -t solo-io/gloo:v0.1 . && rm gloo