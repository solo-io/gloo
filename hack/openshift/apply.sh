#!/usr/bin/env bash
set -ex

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

oc project gloo-system
oc process -f ${BASEDIR}/template.yaml \
 -p APISERVER_OPENSHIFT_MASTER_IP=$(minishift ip) \
  | oc apply -f -

if [ "$1" -eq "-w" ]; then
    echo "attempting to attach to ui"
    sleep 1 # TODO implement actual wait (do we really care?)
    oc port-forward deployment/apiserver-ui 8080
fi