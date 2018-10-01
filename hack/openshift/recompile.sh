#!/usr/bin/env bash

# TODO(ilackarms): refactor this out into setup-new-minishift, apply.sh, and rebuild.sh

set -ex

PROJECT=$1
export VERSION=$2
if [ "$#" -ne 2 ]; then
    echo "invalid, to run: ./hack/recompile.sh PROJECT VERSION"
    exit 1
fi

# set up user, only needed once
oc login -u system:admin
oc adm policy add-cluster-role-to-user cluster-admin gloo
oc login -u gloo -p gloo

eval $(minishift docker-env)
make $PROJECT

docker tag solo-io/$PROJECT-ee:$VERSION solo-io/$PROJECT-ee:$VERSION