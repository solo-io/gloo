#!/usr/bin/env bash

set -ex

minishift addons enable admin-user

oc login -u system:admin
oc adm policy add-cluster-role-to-user cluster-admin gloo
oc adm policy add-cluster-role-to-user cluster-admin default
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
 - "http://localhost:3000"
 - "http://localhost:3001"
 - "http://localhost:3002"
grantMethod: prompt
EOF
