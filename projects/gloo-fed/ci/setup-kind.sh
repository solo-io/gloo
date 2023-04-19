#!/bin/bash -ex

# 0. Assign default values to some of our environment variables
# The name of the management kind cluster to deploy to
MANAGEMENT_CLUSTER="${MANAGEMENT_CLUSTER:-management}"
# The name of the remote kind cluster to deploy to
REMOTE_CLUSTER="${REMOTE_CLUSTER:-remote}"
# The version of the Node Docker image to use for booting the clusters
CLUSTER_NODE_VERSION="${CLUSTER_NODE_VERSION:-v1.25.3}"
# The version used to tag images
VERSION="${VERSION:-1.0.0-ci}"
# If true, use a released chart with the version in $VERSION
FROM_RELEASE="${FROM_RELEASE:-false}"
# The license key used to support enterprise features
GLOO_LICENSE_KEY="${GLOO_LICENSE_KEY:-}"
# Automatically (lazily) determine OS type
if [[ $OSTYPE == 'darwin'* ]]; then
  OS='darwin'
else
  OS='linux'
fi

# set the architecture of the machine (checking for arm64 and if not defaulting to amd64)
ARCH="amd64"
if [[ $(uname -m) == "arm64" ]]; then
  ARCH="arm64"
fi

# set the architecture of the images that you will be building, default to the machines architecture
if [[ -z "${GOARCH}" ]]; then
  GOARCH=$ARCH
fi

# 1. Ensure that a license key is provided
if [ "$GLOO_LICENSE_KEY" == "" ]; then
  echo "please provide a license key"
  exit 0
fi
#Get the latest release from git
if [[ "$FROM_RELEASE" == "true" ]]; then
  VERSION=`git describe --abbrev=0 --tags`
fi
# 2. Build the gloo command line tool, ensuring we have one in the `_output` folder
make glooctl-$OS-$GOARCH -B
shopt -s expand_aliases
alias glooctl=_output/glooctl-$OS-$GOARCH

# 3. Create the management kind cluster
cat <<EOF | kind create cluster --name "$MANAGEMENT_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
EOF

# 4. Create the remote kind cluster
# Add locality labels to remote kind cluster for discovery
(cat <<EOF | kind create cluster --name "$REMOTE_CLUSTER" --image "kindest/node:$CLUSTER_NODE_VERSION" --config=-
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
  apiVersion: kubeadm.k8s.io/v1beta3
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
  apiVersion: kubeadm.k8s.io/v1beta3
  kind: InitConfiguration
  metadata:
    name: config
  nodeRegistration:
    kubeletExtraArgs:
      "feature-gates": "EphemeralContainers=true"
EOF
)

# 5. Set context to management cluster
kubectl config use-context kind-"$MANAGEMENT_CLUSTER"
if [[ $FROM_RELEASE != "true" ]]; then
  # 5. Build local federation and enterprise images used in these clusters and load them into the kind cluster
  VERSION=$VERSION make package-gloo-fed-chart package-gloo-edge-chart -B
  VERSION=$VERSION make kind-build-federation-images -B
  VERSION=$VERSION CLUSTER_NAME=$MANAGEMENT_CLUSTER make kind-load-federation-control-plane-images kind-load-federation-management-plane-images -B
fi

# 6a. Install gloo-fed and gloo-ee to management kind cluster
cat > management-helm-values.yaml << EOF
global:
  extensions:
    extAuth:
      enabled: false
    rateLimit:
      enabled: false
observability:
  enabled: false
prometheus:
  enabled: false
grafana:
  defaultInstallationEnabled: false
gloo-fed:
  enabled: true
gloo:
  gatewayProxies:
    gatewayProxy:
      readConfig: true
      readConfigMulticluster: true
      service:
        type: NodePort
EOF

if [[ $FROM_RELEASE == "true" ]]; then
  glooctl install gateway enterprise --license-key="$GLOO_LICENSE_KEY" --version "$VERSION" --values management-helm-values.yaml
else
  glooctl install gateway enterprise --license-key="$GLOO_LICENSE_KEY" --file _output/helm/gloo-ee-"$VERSION".tgz --values management-helm-values.yaml
fi
rm management-helm-values.yaml

# 6c. Wait for the installation to complete
kubectl -n gloo-system rollout status deployment gloo-fed --timeout=1m || true
kubectl -n gloo-system rollout status deployment gloo --timeout=2m || true
kubectl -n gloo-system rollout status deployment discovery --timeout=2m || true
kubectl -n gloo-system rollout status deployment gateway-proxy --timeout=2m || true
kubectl -n gloo-system rollout status deployment gateway --timeout=2m || true

# 7. Install Gloo in the remote kind cluster
kubectl config use-context kind-"$REMOTE_CLUSTER"


# 7a. Install gloo-ee to remote kind cluster
cat > remote-helm-values.yaml <<EOF
global:
  extensions:
    extAuth:
      enabled: false
    rateLimit:
      enabled: false
observability:
  enabled: false
prometheus:
  enabled: false
grafana:
  defaultInstallationEnabled: false
gloo-fed:
  enabled: false
  glooFedApiserver:
    enable: false
gloo:
  gatewayProxies:
    gatewayProxy:
      failover:
        enabled: true
      service:
        type: NodePort
EOF

if [[ $FROM_RELEASE == "true" ]]; then
  glooctl install gateway enterprise --license-key="$GLOO_LICENSE_KEY" --version "$VERSION" --values remote-helm-values.yaml
else
  # Load enterprise images used in this test
  CLUSTER_NAME=$REMOTE_CLUSTER VERSION=$VERSION make kind-load-federation-control-plane-images -B
  glooctl install gateway enterprise --license-key="$GLOO_LICENSE_KEY" --file _output/helm/gloo-ee-"$VERSION".tgz --values remote-helm-values.yaml
fi
rm remote-helm-values.yaml

# 7c. Wait for the installation to complete
kubectl -n gloo-system rollout status deployment gloo --timeout=2m || true
kubectl -n gloo-system rollout status deployment discovery --timeout=2m || true
kubectl -n gloo-system rollout status deployment gateway-proxy --timeout=2m || true
kubectl patch settings -n gloo-system default --type=merge -p '{"spec":{"watchNamespaces":["gloo-system", "default"]}}'

# 8. Generate certs and keys
# Generate downstream cert and key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
   -keyout tls.key -out tls.crt -subj "/CN=solo.io"

# Generate upstream ca cert and key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
   -keyout mtls.key -out mtls.crt -subj "/CN=solo.io"

glooctl create secret tls --name failover-downstream --certchain tls.crt --privatekey tls.key --rootca mtls.crt

# Revert back to management cluster context
kubectl config use-context kind-"$MANAGEMENT_CLUSTER"

glooctl create secret tls --name failover-upstream --certchain mtls.crt --privatekey mtls.key
rm mtls.key mtls.crt tls.crt tls.key

# Register the gloo clusters
# Automatically determine cluster domains based on OS
if [[ $OS == 'darwin' ]]; then
  MANAGEMENT_CLUSTER_DOMAIN=host.docker.internal
  REMOTE_CLUSTER_DOMAIN=host.docker.internal
else
  MANAGEMENT_CLUSTER_DOMAIN=$(docker exec "$MANAGEMENT_CLUSTER"-control-plane ip addr show dev eth0 | sed -nE 's|\s*inet\s+([0-9.]+).*|\1|p'):6443
  REMOTE_CLUSTER_DOMAIN=$(docker exec "$REMOTE_CLUSTER"-control-plane ip addr show dev eth0 | sed -nE 's|\s*inet\s+([0-9.]+).*|\1|p'):6443
fi
glooctl cluster register --cluster-name kind-"$MANAGEMENT_CLUSTER" --remote-context kind-"$MANAGEMENT_CLUSTER" --local-cluster-domain-override "$MANAGEMENT_CLUSTER_DOMAIN" --federation-namespace gloo-system
glooctl cluster register --cluster-name kind-"$REMOTE_CLUSTER" --remote-context kind-"$REMOTE_CLUSTER" --local-cluster-domain-override "$REMOTE_CLUSTER_DOMAIN" --federation-namespace gloo-system

echo "Registered gloo clusters kind-$MANAGEMENT_CLUSTER and kind-$REMOTE_CLUSTER"

# Set up resources for Failover demo
echo "Set up resources for Failover demo..."

# setup the correct echo image based off chipset type
HTTP_ECHO_IMAGE="hashicorp/http-echo@sha256:ba27d460cd1f22a1a4331bdf74f4fccbc025552357e8a3249c40ae216275de96"
if [[ $ARCH == "arm64" ]]; then
  HTTP_ECHO_IMAGE="gcr.io/solo-test-236622/http-echo@sha256:1efdc13babe46c9ff22154641e75e55400cb5fe1c0521259e6c24a223ccd9beb"
fi

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
          image: $HTTP_ECHO_IMAGE  
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
                          typed_config:
                            "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
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
          image: envoyproxy/envoy:v1.22.0
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

kubectl config use-context kind-"$REMOTE_CLUSTER"

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
          image: $HTTP_ECHO_IMAGE
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
                          typed_config:
                            "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
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
          image: envoyproxy/envoy:v1.22.0
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

kubectl config use-context kind-"$MANAGEMENT_CLUSTER"

kubectl apply -f - <<EOF
apiVersion: fed.gloo.solo.io/v1
kind: FederatedUpstream
metadata:
  name: default-service-blue
  namespace: gloo-system
spec:
  placement:
    clusters:
      - kind-$MANAGEMENT_CLUSTER
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
      - kind-$MANAGEMENT_CLUSTER
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
    clusterName: kind-$MANAGEMENT_CLUSTER
    name: default-service-blue-10000
    namespace: gloo-system
  failoverGroups:
  - priorityGroup:
    - cluster: kind-$REMOTE_CLUSTER
      upstreams:
      - name: default-service-green-10000
        namespace: gloo-system
EOF

# wait for setup to be complete
kubectl -n gloo-system rollout status deployment gloo-fed --timeout=2m
kubectl rollout status deployment echo-blue-deployment --timeout=2m
