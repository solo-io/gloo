#!/bin/bash

VERSION=$(protoc --version)
DESIRED_VERSION="libprotoc 3.6.1"

if [[ ${VERSION} != ${DESIRED_VERSION} ]]; then
    echo protoc --version must be ${DESIRED_VERSION}. Currently ${VERSION};
    exit 1;
fi
