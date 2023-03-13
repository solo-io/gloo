#!/bin/bash
set -e

if which helm > /dev/null; then
    echo "helm is already installed"
    exit 0
fi
VERSION=$HELM_VERSION
IS_VERSION_2="true"
if [ -n "${VERSION}" ]; then
    set -- "$@" --version "${VERSION}"
    if [ "${VERSION}" == "${VERSION#v2.}" ]; then
        IS_VERSION_2="false"
    fi
fi
INSTALL_SCRIPT="https://raw.githubusercontent.com/helm/helm/master/scripts/get"
if [ "${IS_VERSION_2}" == "false" ]; then
    INSTALL_SCRIPT="https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3"
fi

curl "${INSTALL_SCRIPT}" > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh "$@"
if [ "${IS_VERSION_2}" == "true" ]; then
    helm init --client-only
else
    helm repo add stable https://charts.helm.sh/stable
fi

