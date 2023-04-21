#!/bin/bash
set -e
export PATH=$HOME/.gloo/bin:$PATH

export AWS_ACCESS_KEY_ID="ASIA3VU3PCIYCZ6YCE3Z"
export AWS_SECRET_ACCESS_KEY="ZYU5/cR2WhCnC0TKfM0aAaqSrw17e3wumTZQ0RVn"
export AWS_SESSION_TOKEN="IQoJb3JpZ2luX2VjEJr//////////wEaCXVzLWVhc3QtMSJHMEUCIEehUX+0RMaQmCYPIzVn6rUJJfzDMBpJb588yxNDZAONAiEA8c74/9sceD03ixh7nZtdNLb8rvD5CozBtQv6EOoT4zgqnQMIk///////////ARAAGgw4MDI0MTExODg3ODQiDEBdYe3L20K2o8C//yrxAgnNqNG3pz6O/3+qDWQiAKGQ/gTA0uM1ArUUosLRLOjEI2qQKF1Y8Oo7T3p3Vi/pFt16V6JEcPecboi4f1hqcRHa11u9F7p52liiz8EemsetLPRVHrGYkhxlVzJ/ssygX++LJ5q4lm99QmfTpbI9h08j8RBVTvIejkbGObCyoK+lQBgvZ/ybzAyq1ZmBIH5fHTt2Xr1ZgFFBU79il0E+iTXDcynXLQfjUWO5BGRgtio7Sii2GqVm9tX/Ynj2j58Qc2zti0xtzDtSOK/HyNomNw/wwZewrQ8ayV5zkp12ZtFKnWhOFGp+i0yvrQVXXJJkum4f7sr2LojYCJjUVXk0PELEGEMZsOIrHA2yS47OXktY1MN9DzM3Fyu0QcB9le272dKI3TZJH58fAMahVJxeYvGRpU1AJe60t2NtkGinGUHwCzckVV5QVbv7GWFgykW4j9r9V+aglFbN0gmesAvFcyIA1ttvbY4TZDzCBZc+JLUz4zDF+oWiBjqmAVfivamGH2v4RTrUnELuUbOKo0BKAUlGx2tat0ju+puJqg2aLC+imONfh+AAAT+Yy82AKKX4xdyaMm37ldpuIwbHdFGFeA3fNdZK6/vd2RhxqhPOsTa9ZT+OEl7KSxvpMv3nVCEWltEHuu0ZbTTxfwI8HgZ8ScLEtSfGSiPCBTeVdlwZ8iduLd9MnHmPpvZ5r5mBKpelA+geFjEnBh4soxzNg2Mdzhk="
export KUBECONFIG=/Users/ianmacclancy/Downloads/cicd-testNewScripts
NUM_REMOTE_CLUSTERS=1
GLOO_EDGE_LICENSE_KEY="eyJhZGRPbnMiOiIiLCJleHAiOjE2ODM4OTc3NzksImlhdCI6MTY1MjM2MTc3OSwiayI6IjhhL1pqUSIsImx0IjoiZW50IiwicHJvZHVjdCI6Imdsb28ifQ.M33cSptpUqE5R4hZUIm8RHwwQSZGvbiEYpXshPUpEAk"
GLOO_EE_VERSION=1.13.15

init_args() {
  MGMT_CONTEXT="mgmt-cluster"
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




# install edge on the specified cluster call with install_edge_on_cluster <cluster-name> <override file location>
# will not be used for the fed management cluster
install_edge_on_cluster() {
  cluster=$1
  overrideLocation=$2
  echo "Installing edge on cluster $cluster..."
  DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

  echo "create namespace"
  kubectl --context "$cluster" create namespace ${NAMESPACE}

  echo "create secret for TLS"
  kubectl --context "$cluster" create secret tls client-tls-secret --key "$DIR/certs/key.pem" --cert "$DIR/certs/cert.pem" --namespace ${NAMESPACE}

  # Build out helm command by each value -makes adding and removing as we change the tests easier
  HELM_OPTIONS="-n ${NAMESPACE} gloo-ee gloo-ee/gloo-ee"
  HELM_OPTIONS="$HELM_OPTIONS --kube-context=$cluster"

  HELM_OPTIONS="$HELM_OPTIONS --version ${GLOO_EE_VERSION}"
  HELM_OPTIONS="$HELM_OPTIONS --set-string license_key=${GLOO_EDGE_LICENSE_KEY}"
  HELM_OPTIONS="$HELM_OPTIONS --set gloo-fed.enabled=false"
  HELM_OPTIONS="$HELM_OPTIONS --set gloo-fed.glooFedApiserver.enable=false"

  # Used by elb to correctly connect to the proxy
  HELM_OPTIONS="$HELM_OPTIONS --set gloo.gatewayProxies.gatewayProxy.service.extraAnnotations.service.beta.kubernetes.io/aws-load-balancer-backend-protocol=https"
  HELM_OPTIONS="$HELM_OPTIONS --set gloo.gatewayProxies.gatewayProxy.service.extraAnnotations.service.beta.kubernetes.io/aws-load-balancer-ssl-ports=https"

  if [ "$WEBHOOK_VALIDATION" == 'false' ]; then
    HELM_OPTIONS="$HELM_OPTIONS --set gloo.gateway.validation.enabled=false"
  fi

  if [[ -n $overrideLocation ]]; then
    HELM_OPTIONS="$HELM_OPTIONS --values $overrideLocation"
  fi
  # Split the HELM_OPTIONS variable into an array of arguments
  IFS=' ' read -r -a HELM_ARGS <<< "${HELM_OPTIONS}"

  # Pass the array of arguments to helm install
  helm install "${HELM_ARGS[@]}"
}

install_fed_on_management_cluster() {
  echo "Installing edge on cluster ${MGMT_CONTEXT}..."
  DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
  TMOBILE_HELM_OVERRIDES="$DIR/assets/tmobileOverrides.yaml"

  echo "create namespace"
  kubectl --context "${MGMT_CONTEXT}" create namespace ${NAMESPACE}

  echo "create secret for TLS"
  kubectl --context "${CLUSTER_CONTEXT}" create secret tls client-tls-secret --key "$DIR/certs/key.pem" --cert "$DIR/certs/cert.pem" --namespace ${NAMESPACE}

  HELM_OPTIONS="-n ${NAMESPACE} gloo-ee gloo-ee/gloo-ee"
  HELM_OPTIONS="$HELM_OPTIONS --kube-context=${MGMT_CONTEXT}"
  HELM_OPTIONS="$HELM_OPTIONS --version ${GLOO_EE_VERSION}"
  HELM_OPTIONS="$HELM_OPTIONS --set-string license_key=${GLOO_EDGE_LICENSE_KEY}"

  # Used by elb to correctly connect to the proxy
  HELM_OPTIONS="$HELM_OPTIONS --set gloo.gatewayProxies.gatewayProxy.service.extraAnnotations.service.beta.kubernetes.io/aws-load-balancer-backend-protocol=https"
  HELM_OPTIONS="$HELM_OPTIONS --set gloo.gatewayProxies.gatewayProxy.service.extraAnnotations.service.beta.kubernetes.io/aws-load-balancer-ssl-ports=https"

  if [ $WEBHOOK_VALIDATION == 'false' ]; then
    HELM_OPTIONS="$HELM_OPTIONS --set gloo.gateway.validation.enabled=false"
  fi

  helm install $HELM_OPTIONS
}

register_clusters() {
  kubectl config use-context ${MGMT_CONTEXT}
  echo "Register Cluster"
  for CLUSTER_CONTEXT in "${CLUSTER_CONTEXTS[@]}"; do
    echo "Register Cluster ${CLUSTER_CONTEXT}"
    glooctl cluster register\
        --cluster-name ${CLUSTER_CONTEXT}\
        --remote-context ${CLUSTER_CONTEXT}
  done
}

install_edge_fed() {
  # fed cluster
  install_fed_on_management_cluster

  # managed clusters
  for CLUSTER_CONTEXT in "${CLUSTER_CONTEXTS[@]}"; do
      install_edge_on_cluster $CLUSTER_CONTEXT $FED_HELM_OVERRIDES
  done

  register_clusters
}

install_edge_tmobile() {
  DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
  TMOBILE_HELM_OVERRIDES="$DIR/assets/tmobileOverrides.yaml"

  # Management cluster is the most likely case for tmobile as we generally run with only 1 cluster
  install_edge_on_cluster $MGMT_CONTEXT $TMOBILE_HELM_OVERRIDES

  for CLUSTER_CONTEXT in "${CLUSTER_CONTEXTS[@]}"; do
        install_edge_on_cluster $CLUSTER_CONTEXT $TMOBILE_HELM_OVERRIDES
  done
}

install_edge_no_overrides() {
  install_edge_on_cluster $MGMT_CONTEXT
  for CLUSTER_CONTEXT in "${CLUSTER_CONTEXTS[@]}"; do
        install_edge_on_cluster $CLUSTER_CONTEXT
  done
}

main() {
  init_args
  init_helm
  if [[ ${FED} == "true" ]]; then
    echo Fed install - Installing Fed on mgmt context and edge on clusters
    if [[ ${GLOO_EE_VERSION} == "branch" ]]; then
          echo "install gloo edge from branch" #todo: currently a no op
    else
      install_edge_fed
    fi
  elif [[ ${TMOBILE_HELM_VALUES} == "true" ]]; then
    echo Tmobile install - Installing edge with tmobile overrides
    if [[ ${GLOO_EE_VERSION} == "branch" ]]; then
      echo "install gloo edge from branch" #todo: currently a no op
    else
      install_edge_tmobile
    fi
  else
    echo standard edge install
    echo "install edge version: ${GLOO_EE_VERSION}"
    install_edge_no_overrides
  fi
}

main "$@"