#!/bin/bash

# If gloo-fed and gloo-fed-console are already installed
GLOO_FED=$(kubectl -n gloo-fed rollout status deployment gloo-fed 2> /dev/null)
GLOO_FED_CONSOLE=$(kubectl -n gloo-fed rollout status deployment gloo-fed-console 2> /dev/null)
if [[ "$GLOO_FED" = "deployment \"gloo-fed\" successfully rolled out" ]] && [[ "$GLOO_FED_CONSOLE" = "deployment \"gloo-fed-console\" successfully rolled out" ]]; then
  # check if cluster is registered
  HAS_CLUSTER_REGISTERED=$(glooctl cluster list)
  if [[ "HAS_CLUSTER_REGISTERED" == "" ]]; then
    echo ""
    echo "Register a cluster with:"
    echo "KIND: glooctl cluster register --cluster-name [kind-local] --local-cluster-domain-override host.docker.internal --remote-context [kind-local] --federation-namespace [default: gloo-fed]"
    echo "GKE:  glooctl cluster register --cluster-name [gloo-fed-remote] --remote-context [gke-context-name] --federation-namespace [default: gloo-fed]"
  fi
  # check port forwarding is already setup
  PORT_FWD=$( pgrep -f "kubectl port-forward -n gloo-fed deploy/gloo-fed-console 8090")
  if [[ "$PORT_FWD" == "" ]]; then
    kubectl port-forward -n gloo-fed deploy/gloo-fed-console 8090 &
  fi
  exit 0
else
  LOCAL_APISERVER=$(pgrep -f run-apiserver)
  LOCAL_ENVOY=$(pgrep -f run-envoy)
  if [[ "$LOCAL_APISERVER" != "" ]] && [[ "$LOCAL_ENVOY" != "" ]]; then
    HAS_CLUSTER_REGISTERED=$(glooctl cluster list)
      if [[ "HAS_CLUSTER_REGISTERED" == "" ]]; then
        echo ""
        echo "Register a cluster with:"
        echo "KIND: glooctl cluster register --cluster-name [kind-local] --local-cluster-domain-override host.docker.internal --remote-context [kind-local] --federation-namespace [default: gloo-fed]"
        echo "GKE:  glooctl cluster register --cluster-name [gloo-fed-remote] --remote-context [gke-context-name] --federation-namespace [default: gloo-fed]"
      fi
    exit 0
  else
    # User needs to locally set up apiserver or install federation
    echo ""
    echo "Gloo UI requires apiserver to be running locally or gloo-fed to be installed."
    echo "Install gloo-fed with: glooctl install federation"
    echo "Run the apiserver locally with: make run-apiserver; make run-envoy"
    echo "Register a cluster with:"
    echo "KIND: glooctl cluster register --cluster-name [kind-local] --local-cluster-domain-override host.docker.internal --remote-context [kind-local] --federation-namespace [default: gloo-fed]"
    echo "GKE:  glooctl cluster register --cluster-name [gloo-fed-remote] --remote-context [gke-context-name] --federation-namespace [default: gloo-fed]"
    exit 1
  fi
fi