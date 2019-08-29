#!/usr/bin/env bash

####################################################################################################
# This script is used to deploy an LDAP server with sample user/group configuration to Kubernetes #
####################################################################################################
set -e

if [ -z "$1" ]; then
  echo "No namespace provided, using default namespace"
  NAMESPACE='default'
else
  NAMESPACE=$1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "Creating configmap with LDAP server bootstrap config..."
kubectl create configmap ldap -n "${NAMESPACE}" --from-file="$DIR/ldif"

echo "Creating LDAP service and deployment..."
kubectl apply -n "${NAMESPACE}" -f "${DIR}/ldap-server-manifest.yaml"