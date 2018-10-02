#!/usr/bin/env bash

# TODO(ilackarms): refactor this out into setup-new-minishift, apply.sh, and rebuild.sh

set -ex

# set up user, only needed once
oc login -u system:admin
oc adm policy add-cluster-role-to-user cluster-admin gloo
oc login -u gloo -p gloo
