#!/bin/bash
set -e
export PATH=$HOME/.gloo/bin:$PATH

init_args() {
  FED_CONTEXT="mgmt-cluster"
  NAMESPACE="gloo-system"
  CLUSTER_CONTEXTS=()
  for ((i=1; i<=NUM_REMOTE_CLUSTERS; i++)); do
    CLUSTER_CONTEXTS[i]="cluster-${i}"
  done
}

init_helm() {
  helm repo add gloo-ee https://storage.googleapis.com/gloo-ee-helm
  helm repo add gloo-fed https://storage.googleapis.com/gloo-fed-helm
  helm repo update
}

install_released_edge() {
  echo "======== INSTALL RELEASED GLOO EDGE ========="
  for CLUSTER_CONTEXT in "${CLUSTER_CONTEXTS[@]}"; do
      echo "Installing edge on cluster ${CLUSTER_CONTEXT}..."
      DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

      echo "create namespace"
      kubectl --context "${CLUSTER_CONTEXT}" create namespace ${NAMESPACE}

      echo "create secret for TLS"
      kubectl --context "${CLUSTER_CONTEXT}" create secret tls client-tls-secret --key "$DIR/certs/key.pem" --cert "$DIR/certs/cert.pem" --namespace ${NAMESPACE}

      echo "install gloo"
      helm install -n ${NAMESPACE} gloo-ee gloo-ee/gloo-ee \
          --kube-context=${CLUSTER_CONTEXT} \
          --version ${GLOO_EE_VERSION} \
          --set-string license_key="${GLOO_EDGE_LICENSE_KEY}"\
          --set gloo-fed.enabled=false \
          --set gloo-fed.glooFedApiserver.enable=false
  done
}

install_fed() {
    helm install -n ${NAMESPACE} gloo-ee gloo-ee/gloo-ee \
        --kube-context=${FED_CONTEXT} \
        --create-namespace \
        --version "${GLOO_EE_VERSION}" \
        --set-string license_key="${GLOO_EDGE_LICENSE_KEY}"
}

register_clusters() {
  kubectl config use-context ${FED_CONTEXT}
  echo "Register Cluster"
  for CLUSTER_CONTEXT in "${CLUSTER_CONTEXTS[@]}"; do
    echo "Register Cluster ${CLUSTER_CONTEXT}"
    glooctl cluster register\
        --cluster-name ${CLUSTER_CONTEXT}\
        --remote-context ${CLUSTER_CONTEXT}
  done
}

main() {
  init_args
  init_helm
  if [[ ${GLOO_EE_VERSION} == "branch" ]]; then
      echo "install gloo edge from branch" #todo: currently a no op
    else
      echo "install edge version: ${GLOO_EE_VERSION}"
      install_released_edge
      install_fed
      register_clusters
  fi
}

main "$@"