#!/usr/bin/env bash
set -ex

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

oc project gloo-system
oc process -f ${BASEDIR}/template.yaml \
 -p APISERVER_OPENSHIFT_MASTER_IP=$(minishift ip) \
  | oc apply -f -

cat << EOF | oc apply -f -
kind: OAuthClient
apiVersion: oauth.openshift.io/v1
metadata:
 name: gloo
secret: gloo
redirectURIs:
 - "http://localhost:8080"
grantMethod: prompt
EOF

oc port-forward deployment/apiserver-ui 8080