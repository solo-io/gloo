kind create cluster

SCRIPTPATH=$( cd "$(dirname "$0")" && pwd -P )

kubectl apply -f $SCRIPTPATH/crds/gateway-crds.yaml

# TODO: make this smarter if desired; simple relative path works for now
kubectl apply -f $SCRIPTPATH/../../install/helm/gloo/crds/gateway.solo.io_v1_RouteOption.yaml

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml

# Wait for MetalLB to become available.
kubectl rollout status -n metallb-system deployment/controller --timeout 2m
kubectl rollout status -n metallb-system daemonset/speaker --timeout 2m
kubectl wait -n metallb-system  pod -l app=metallb --for=condition=Ready --timeout=10s

SUBNET=$(docker network inspect  kind -f '{{(index .IPAM.Config 0).Subnet}}'| cut -d '.' -f1,2)
MIN=${SUBNET}.255.0
MAX=${SUBNET}.255.231

kubectl apply -f - <<EOF
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: address-pool
  namespace: metallb-system
spec:
  addresses:
  - ${MIN}-${MAX}
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: advertisement
  namespace: metallb-system
spec:
  ipAddressPools:
  - address-pool
EOF

ISTIO_VERSION="${ISTIO_VERSION:-1.18.2}"

TARGET_ARCH=x86_64
if [[ $ARCH == 'arm64' ]]; then
  TARGET_ARCH=arm64
fi

curl -L https://istio.io/downloadIstio | ISTIO_VERSION=$ISTIO_VERSION TARGET_ARCH=$TARGET_ARCH sh -
yes | "./istio-$ISTIO_VERSION/bin/istioctl" install --set profile=minimal

# apply gateway
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: http
spec:
  gatewayClassName: gloo-gateway
  listeners:
  - allowedRoutes:
      namespaces:
        from: All
    name: http
    port: 8080
    protocol: HTTP
EOF

# TODO: remove helper apps to test sds integration

kubectl create namespace bookinfo

kubectl label namespace bookinfo istio-injection=enabled

kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.20/samples/bookinfo/platform/kube/bookinfo.yaml -n bookinfo

cat <<EOF | kubectl apply -f -
apiVersion: "security.istio.io/v1beta1"
kind: "PeerAuthentication"
metadata:
  name: "productpage"
  namespace: "bookinfo"
spec:
  selector:
    matchLabels:
      app: productpage
  mtls:
    mode: STRICT
EOF

kubectl apply -f- <<EOF
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: productpage
  namespace: bookinfo
  labels:
    example: productpage-route
spec:
  parentRefs:
    - name: http
      namespace: default
  hostnames:
    - "www.example.com"
  rules:
    - backendRefs:
        - name: productpage
          port: 9080
EOF

# To test:
# kubectl port-forward deployment/gloo-proxy-http 8080:8080
# curl -I localhost:8080/productpage -H "host: www.example.com" -v