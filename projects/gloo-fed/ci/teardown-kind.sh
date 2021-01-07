#!/bin/bash

set -ex

if [ "$1" == "" ] || [ "$2" == "" ]; then
  echo "please provide a name for both the master and remote clusters"
  exit 0
fi

#Delete given kind clusters
kind delete cluster --name "$1"
kind delete cluster --name "$2"