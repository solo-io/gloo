#!/usr/bin/env bash

set -ex

# set up user, only needed once
oc login -u system:admin
oc adm policy add-cluster-role-to-user cluster-admin gloo
oc login -u gloo -p gloo

cat << EOF | oc apply -f -
kind: OAuthClient
apiVersion: oauth.openshift.io/v1
metadata:
 name: gloo
secret: gloo
redirectURIs:
 - "http://localhost:8080"
 - "http://localhost:8082"
grantMethod: prompt
EOF