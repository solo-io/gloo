#!/bin/bash

# Absolute path to this script, e.g. /home/user/go/.../gloo-fed/ci/setup-kind.sh
SCRIPT=$(readlink -f "$0")
# Absolute path this script is in, thus /home/user/go/.../gloo-fed/ci
CI_DIR=$(dirname "$SCRIPT")
GIT_SEMVER_SCRIPT="$CI_DIR/../../../git-semver.sh"

set -ex

kindClusterImage=kindest/node:v1.17.17@sha256:66f1d0d91a88b8a001811e2f1054af60eef3b669a9a74f9b6db871f2f1eeed00

if [ "$1" == "" ] || [ "$2" == "" ]; then
  echo "please provide a name for both the master and remote clusters"
  exit 0
fi

if [ "$GLOO_LICENSE_KEY" == "" ]; then
  echo "please provide a license key"
  exit 0
fi

# Ensure that dependencies are consistent with what's in go.mod.
go mod tidy

GLOO_VERSION="$(echo $(go list -m github.com/solo-io/gloo) | cut -d' ' -f2)"
# NOTE: If inter-PR dependency is needed, this must be changed to a hard-coded version (ex: v1.7.0-beta25).
GLOO_VERSION_TEST_INSTALL=$GLOO_VERSION
VERSION=$(. $GIT_SEMVER_SCRIPT)

# Install glooctl
if which glooctl;
then
    echo "Found glooctl installed already"
    glooctl upgrade --release="$GLOO_VERSION_TEST_INSTALL"
else
    echo "Installing glooctl"
    curl -sL https://run.solo.io/gloo/install | sh
    export PATH=$HOME/.gloo/bin:$PATH
    glooctl upgrade --release="$GLOO_VERSION_TEST_INSTALL"
fi

cat <<EOF | kind create cluster --name $1 --image $kindClusterImage --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: ClusterConfiguration
  metadata:
    name: config
  apiServer:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  scheduler:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  controllerManager:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: InitConfiguration
  metadata:
    name: config
  nodeRegistration:
    kubeletExtraArgs:
      "feature-gates": "EphemeralContainers=true"
EOF

# Add locality labels to remote kind cluster for discovery
(cat <<EOF | kind create cluster --name "$2" --image $kindClusterImage --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 32000
    hostPort: 32000
    protocol: TCP
# - role: worker
kubeadmConfigPatches:
- |
  kind: InitConfiguration
  nodeRegistration:
    kubeletExtraArgs:
      node-labels: "topology.kubernetes.io/region=us-east-1,topology.kubernetes.io/zone=us-east-1c"
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: ClusterConfiguration
  metadata:
    name: config
  apiServer:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  scheduler:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  controllerManager:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: InitConfiguration
  metadata:
    name: config
  nodeRegistration:
    kubeletExtraArgs:
      "feature-gates": "EphemeralContainers=true"
EOF
)

kubectl config use-context kind-"$1"

yarn --cwd projects/ui build
make CLUSTER_NAME=$1 VERSION=${VERSION} package-gloo-fed-chart package-gloo-edge-chart gloofed-load-kind-images
# Only build and load in the gloo-ee images used in this test
make VERSION=${VERSION} gloo-ee-docker gloo-ee-envoy-wrapper-docker rate-limit-ee-docker -B
kind load docker-image quay.io/solo-io/gloo-ee:${VERSION} --name $1
kind load docker-image quay.io/solo-io/gloo-ee-envoy-wrapper:${VERSION} --name $1
kind load docker-image quay.io/solo-io/rate-limit-ee:${VERSION} --name $1

# Install gloo-fed and gloo-ee to cluster $1

cat > basic-enterprise.yaml << EOF
rateLimit:
  enable: false
global:
  extensions:
    extAuth:
      enabled: false
observability:
  enabled: false
prometheus:
  enabled: false
grafana:
  defaultInstallationEnabled: false
gloo:
  gatewayProxies:
    gatewayProxy:
      readConfig: true
      readConfigMulticluster: true
      service:
        type: NodePort
gloo-fed:
  enabled: true
EOF

glooctl install gateway enterprise --license-key=$GLOO_LICENSE_KEY --file _output/helm/gloo-ee-${VERSION}.tgz --with-gloo-fed=false --values basic-enterprise.yaml

rm basic-enterprise.yaml

# gloo-system rollout
kubectl -n gloo-system rollout status deployment gloo-fed --timeout=1m || true
kubectl -n gloo-system rollout status deployment gloo --timeout=2m || true
kubectl -n gloo-system rollout status deployment discovery --timeout=2m || true
kubectl -n gloo-system rollout status deployment gateway-proxy --timeout=2m || true
kubectl -n gloo-system rollout status deployment gateway --timeout=2m || true


# Install gloo to cluster $2
kubectl config use-context kind-"$2"
cat > nodeport.yaml <<EOF
gatewayProxies:
  gatewayProxy:
    failover:
      enabled: true
    service:
      type: NodePort
EOF
glooctl install gateway --values nodeport.yaml
rm nodeport.yaml
kubectl -n gloo-system rollout status deployment gloo --timeout=2m || true
kubectl -n gloo-system rollout status deployment discovery --timeout=2m || true
kubectl -n gloo-system rollout status deployment gateway-proxy --timeout=2m || true
kubectl -n gloo-system rollout status deployment gateway --timeout=2m || true
kubectl patch settings -n gloo-system default --type=merge -p '{"spec":{"watchNamespaces":["gloo-system", "default"]}}'

# Generate downstream cert and key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
   -keyout tls.key -out tls.crt -subj "/CN=solo.io"

# Generate upstream ca cert and key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
   -keyout mtls.key -out mtls.crt -subj "/CN=solo.io"

glooctl create secret tls --name failover-downstream --certchain tls.crt --privatekey tls.key --rootca mtls.crt

# Revert back to cluster context $1
kubectl config use-context kind-"$1"

glooctl create secret tls --name failover-upstream --certchain mtls.crt --privatekey mtls.key
rm mtls.key mtls.crt tls.crt tls.key

case $(uname) in
  "Darwin")
  {
      CLUSTER_DOMAIN_MGMT=host.docker.internal
      CLUSTER_DOMAIN_REMOTE=host.docker.internal
  } ;;
  "Linux")
  {
      CLUSTER_DOMAIN_MGMT=$(docker exec $1-control-plane ip addr show dev eth0 | sed -nE 's|\s*inet\s+([0-9.]+).*|\1|p'):6443
      CLUSTER_DOMAIN_REMOTE=$(docker exec $2-control-plane ip addr show dev eth0 | sed -nE 's|\s*inet\s+([0-9.]+).*|\1|p'):6443
  } ;;
  *)
  {
      echo "Unsupported OS"
      exit 1
  } ;;
esac

# Register the gloo clusters
glooctl cluster register --cluster-name kind-$1 --remote-context kind-$1 --local-cluster-domain-override $CLUSTER_DOMAIN_MGMT --federation-namespace gloo-system
glooctl cluster register --cluster-name kind-$2 --remote-context kind-$2 --local-cluster-domain-override $CLUSTER_DOMAIN_REMOTE --federation-namespace gloo-system

echo "Registered gloo clusters kind-$1 and kind-$2"

# Set up resources for Failover demo
echo "Set up resources for Failover demo..."
# Apply blue deployment and service
kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  labels:
    app: bluegreen
  name: service-blue
  namespace: default
spec:
  clusterIP: 10.96.10.40
  ports:
    - name: color
      port: 10000
      protocol: TCP
      targetPort: 10000
  selector:
    app: bluegreen
    text: blue
  sessionAffinity: None
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: bluegreen
    text: blue
  name: echo-blue-deployment
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: bluegreen
      text: blue
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: bluegreen
        text: blue
    spec:
      containers:
        - args:
            - -text="blue-pod"
          image: hashicorp/http-echo@sha256:ba27d460cd1f22a1a4331bdf74f4fccbc025552357e8a3249c40ae216275de96
          imagePullPolicy: IfNotPresent
          name: echo
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
        - args:
            - --config-yaml
            - |2
              node:
               cluster: ingress
               id: "ingress~for-testing"
               metadata:
                role: "default~proxy"
              static_resources:
                listeners:
                - name: listener_0
                  address:
                    socket_address: { address: 0.0.0.0, port_value: 10000 }
                  filter_chains:
                  - filters:
                    - name: envoy.filters.network.http_connection_manager
                      typed_config:
                        "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                        stat_prefix: ingress_http
                        codec_type: AUTO
                        route_config:
                          name: local_route
                          virtual_hosts:
                          - name: local_service
                            domains: ["*"]
                            routes:
                            - match: { prefix: "/" }
                              route: { cluster: some_service }
                        http_filters:
                        - name: envoy.filters.http.health_check
                          typed_config:
                            "@type": type.googleapis.com/envoy.extensions.filters.http.health_check.v3.HealthCheck
                            pass_through_mode: true
                        - name: envoy.filters.http.router
                clusters:
                - name: some_service
                  connect_timeout: 0.25s
                  type: STATIC
                  lb_policy: ROUND_ROBIN
                  load_assignment:
                    cluster_name: some_service
                    endpoints:
                    - lb_endpoints:
                      - endpoint:
                          address:
                            socket_address:
                              address: 0.0.0.0
                              port_value: 5678
              admin:
                access_log_path: /dev/null
                address:
                  socket_address:
                    address: 0.0.0.0
                    port_value: 19000
            - --disable-hot-restart
            - --log-level
            - debug
            - --concurrency
            - "1"
            - --file-flush-interval-msec
            - "10"
          image: envoyproxy/envoy:v1.14.2
          imagePullPolicy: IfNotPresent
          name: envoy
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 0
EOF

kubectl config use-context kind-"$2"

# Apply green deployment and service
kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  labels:
    app: bluegreen
  name: service-green
  namespace: default
spec:
  clusterIP: 10.96.59.232
  ports:
    - name: color
      port: 10000
      protocol: TCP
      targetPort: 10000
  selector:
    app: bluegreen
    text: green
  sessionAffinity: None
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: bluegreen
    text: green
  name: echo-green-deployment
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: bluegreen
      text: green
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: bluegreen
        text: green
    spec:
      containers:
        - args:
            - -text="green-pod"
          image: hashicorp/http-echo@sha256:ba27d460cd1f22a1a4331bdf74f4fccbc025552357e8a3249c40ae216275de96
          imagePullPolicy: IfNotPresent
          name: echo
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
        - args:
            - --config-yaml
            - |2
              node:
               cluster: ingress
               id: "ingress~for-testing"
               metadata:
                role: "default~proxy"
              static_resources:
                listeners:
                - name: listener_0
                  address:
                    socket_address: { address: 0.0.0.0, port_value: 10000 }
                  filter_chains:
                  - filters:
                    - name: envoy.filters.network.http_connection_manager
                      typed_config:
                        "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                        stat_prefix: ingress_http
                        codec_type: AUTO
                        route_config:
                          name: local_route
                          virtual_hosts:
                          - name: local_service
                            domains: ["*"]
                            routes:
                            - match: { prefix: "/" }
                              route: { cluster: some_service }
                        http_filters:
                        - name: envoy.filters.http.health_check
                          typed_config:
                            "@type": type.googleapis.com/envoy.extensions.filters.http.health_check.v3.HealthCheck
                            pass_through_mode: true
                        - name: envoy.filters.http.router
                clusters:
                - name: some_service
                  connect_timeout: 0.25s
                  type: STATIC
                  lb_policy: ROUND_ROBIN
                  load_assignment:
                    cluster_name: some_service
                    endpoints:
                    - lb_endpoints:
                      - endpoint:
                          address:
                            socket_address:
                              address: 0.0.0.0
                              port_value: 5678
              admin:
                access_log_path: /dev/null
                address:
                  socket_address:
                    address: 0.0.0.0
                    port_value: 19000
            - --disable-hot-restart
            - --log-level
            - debug
            - --concurrency
            - "1"
            - --file-flush-interval-msec
            - "10"
          image: envoyproxy/envoy:v1.14.2
          imagePullPolicy: IfNotPresent
          name: envoy
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 0
EOF

kubectl config use-context kind-"$1"

kubectl apply -f - <<EOF
apiVersion: fed.gloo.solo.io/v1
kind: FederatedUpstream
metadata:
  name: default-service-blue
  namespace: gloo-system
spec:
  placement:
    clusters:
      - kind-$1
    namespaces:
      - gloo-system
  template:
    metadata:
      name: default-service-blue-10000
    spec:
      discoveryMetadata: {}
      healthChecks:
        - healthyThreshold: 1
          httpHealthCheck:
            path: /health
          interval: 1s
          noTrafficInterval: 1s
          timeout: 1s
          unhealthyThreshold: 1
      kube:
        selector:
          app: bluegreen
          text: blue
        serviceName: service-blue
        serviceNamespace: default
        servicePort: 10000
---
apiVersion: fed.gateway.solo.io/v1
kind: FederatedVirtualService
metadata:
  name: simple-route
  namespace: gloo-system
spec:
  placement:
    clusters:
      - kind-$1
    namespaces:
      - gloo-system
  template:
    spec:
      virtualHost:
        domains:
        - '*'
        routes:
        - matchers:
          - prefix: /
          routeAction:
            single:
              upstream:
                name: default-service-blue-10000
                namespace: gloo-system
    metadata:
      name: simple-route
---
apiVersion: fed.solo.io/v1
kind: FailoverScheme
metadata:
  name: failover-test-scheme
  namespace: gloo-system
spec:
  primary:
    clusterName: kind-$1
    name: default-service-blue-10000
    namespace: gloo-system
  failoverGroups:
  - priorityGroup:
    - cluster: kind-$2
      upstreams:
      - name: default-service-green-10000
        namespace: gloo-system
EOF

# wait for setup to be complete
kubectl -n gloo-system rollout status deployment gloo-fed --timeout=2m
kubectl rollout status deployment echo-blue-deployment --timeout=2m
