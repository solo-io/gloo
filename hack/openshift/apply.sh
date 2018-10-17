#!/usr/bin/env bash
set -ex

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

pullPolicy=IfNotPresent

oc project gloo-system
oc process -f ${BASEDIR}/template.yaml \
 -p APISERVER_OPENSHIFT_MASTER_IP=$(minishift ip) \
 -p IMAGE_PULL_POLICY=${pullPolicy} \
  | oc apply -f -

if [[ "$1" = "-w" ]]; then
    echo "attempting to attach to ui"
    sleep 1 # TODO implement actual wait (do we really care?)
    kubectl port-forward deployment/apiserver-ui 8080
fi