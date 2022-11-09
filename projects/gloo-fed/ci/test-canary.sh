#!/bin/bash -ex

# This file contains a list of steps to perform an initial smoke test of the Gloo Fed canary upgrade feature.
# As canary upgrades are still a work-in-progress, for now this is just a list of steps that can
# be run and verified manually by a developer. Once canary upgrades are feature complete, these steps
# should ideally be moved into an e2e test that runs in ci.

# In this script, the "old" and "new" version are both the same locally-built version (since latest version is 
# guaranteed to have all the skv2 "allow unknown fields" fixes), but in a real upgrade scenario the "old" and "new"
# would be different gloo-ee versions.

# All of the `make` commands below should be run from root of the solo-projects repo.
# Ensure that the `GLOO_LICENSE_KEY` env var is set.

# Get latest glooctl
glooctl upgrade --release=experimental

# Clean up previous installations and built charts (if needed)
rm -rf _test _output
kind delete cluster --name mgmt
kind delete cluster --name remote1
kind delete cluster --name remote2

# Create clusters
# `mgmt` cluster will contain old and new versions of Gloo Fed (installed in namespace gloo-system-old and gloo-system-new)
kind create cluster --name mgmt
# `remote1` will contain the "old" version of Gloo Edge in namespace gloo-system, managed by old Gloo Fed
kind create cluster --name remote1
# `remote2` will contain the "new" version of Gloo Edge in namespace gloo-system, managed by new Gloo Fed
kind create cluster --name remote2

# Build chart and load images onto mgmt cluster
k config use-context kind-mgmt
make VERSION=0.0.0-kind CLUSTER_NAME=mgmt LOCAL_BUILD=1 build-kind-assets gloofed-docker gloofed-load-kind-images -B

# Load the images onto the remote clusters too (don't need to rebuild everything)
k config use-context kind-remote1
make VERSION=0.0.0-kind CLUSTER_NAME=remote1 load-kind-images-non-fips gloofed-load-kind-images -B
k config use-context kind-remote2
make VERSION=0.0.0-kind CLUSTER_NAME=remote2 load-kind-images-non-fips gloofed-load-kind-images -B

# Install old Gloo Fed on mgmt cluster. We are enabling multi-cluster rbac so we can make sure the rbac resources
# get generated successfully when there are 2 Gloo Fed instances:
# https://docs.solo.io/gloo-edge/master/guides/gloo_federation/multicluster_rbac/
k config use-context kind-mgmt
helm install -n gloo-system-old gloo-fed _test/gloo-fed --create-namespace --set-string license_key=$GLOO_LICENSE_KEY --set enableMultiClusterRbac=true

# This is needed so we have permission to apply federated resources on the kind cluster when multi-cluster rbac is enabled
kubectl apply -f - <<EOF
apiVersion: multicluster.solo.io/v1alpha1
kind: MultiClusterRoleBinding
metadata:
  name: kind-admin
  namespace: gloo-system-old
spec:
  roleRef:
    name: gloo-fed
    namespace: gloo-system-old
  subjects:
  - kind: User
    name: kubernetes-admin
EOF

# Install "old" Gloo Edge on remote1
k config use-context kind-remote1
helm install -n gloo-system gloo-ee _test/gloo-ee --create-namespace --set-string license_key=$GLOO_LICENSE_KEY --set gloo-fed.enabled=false

# Register the Gloo Edge instance on remote1
k config use-context kind-mgmt
# federation-namespace is the ns that gloo-fed is installed in
# remote-namespace is the ns that the remote gloo-ee is installed in
# cluster-name is the user-provided name for the KubernetesCluster CR that will be created
# remote-context is the k8s context of the cluster that gloo-ee is on
glooctl cluster register --federation-namespace gloo-system-old --remote-namespace gloo-system --cluster-name mycluster --remote-context kind-remote1 --local-cluster-domain-override host.docker.internal

# Create a federated resource in a separate (shared) namespace (just to illustrate that both Gloo Fed instances can read from the
# same shared ns using the same registered cluster name)
kubectl create ns gloo-fed
kubectl apply -f - <<EOF
apiVersion: fed.gloo.solo.io/v1
kind: FederatedUpstream
metadata:
  name: my-fed-upstream
  namespace: gloo-fed
spec:
  placement:
    clusters:
      - mycluster
    namespaces:
      - gloo-system
  template:
    metadata:
      name: my-upstream
    spec:
      nonExistentField: xyz # test for adding an unknown field, this should not cause errors
      kube:
        serviceName: my-service
        serviceNamespace: default
        servicePort: 10000
EOF

# Do some sanity checks:
# Make sure all pods are started up with no errors
kubectl get pod -A
# Check status of federated resource
k get federatedupstream -n gloo-fed my-fed-upstream -oyaml
# Check status of remote resource and make sure its spec matches the federated resource spec
k --context kind-remote1 get us -n gloo-system my-upstream -oyaml


# Install "new" Gloo Fed on mgmt cluster
helm install -n gloo-system-new gloo-fed _test/gloo-fed --create-namespace --set-string license_key=$GLOO_LICENSE_KEY --set enableMultiClusterRbac=true

kubectl apply -f - <<EOF
apiVersion: multicluster.solo.io/v1alpha1
kind: MultiClusterRoleBinding
metadata:
  name: kind-admin
  namespace: gloo-system-new
spec:
  roleRef:
    name: gloo-fed
    namespace: gloo-system-new
  subjects:
  - kind: User
    name: kubernetes-admin
EOF

# Install "new" Gloo Edge on remote2
k config use-context kind-remote2
helm install -n gloo-system gloo-ee _test/gloo-ee --create-namespace --set-string license_key=$GLOO_LICENSE_KEY --set gloo-fed.enabled=false

# Register the Gloo Edge instance on remote2
# We are using the same registered cluster name (mycluster) just to illustrate that this is possible, and that the federatedupstream
# placement (mycluster / gloo-system) does not need to be modified in order for the new Gloo Fed to pick it up)
k config use-context kind-mgmt
glooctl cluster register --federation-namespace gloo-system-new --remote-namespace gloo-system --cluster-name mycluster --remote-context kind-remote2 --local-cluster-domain-override host.docker.internal

# Do some sanity checks (again):
# Make sure all pods are started up with no errors
kubectl get pod -A
# Check status of federated resource
k get federatedupstream -n gloo-fed my-fed-upstream -oyaml
# Check status of remote resources and make sure their specs match the federated resource spec
k --context kind-remote1 get us -n gloo-system my-upstream -oyaml
k --context kind-remote2 get us -n gloo-system my-upstream -oyaml
