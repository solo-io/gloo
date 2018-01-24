#!/usr/bin/env bash

set -e

CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o glue .

docker build -t solo-io/glue:v0.1 . && rm glue