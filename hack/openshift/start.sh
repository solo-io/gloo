#!/usr/bin/env bash

set -ex

oc apply -f hack/openshift/template.yaml
oc process gloo-ee-installation-template \
 -p APISERVER_OPENSHIFT_MASTER_IP=$(minishift ip) \
  | oc apply -f -
kubectl port-forward deployment/apiserver-ui 8080:80