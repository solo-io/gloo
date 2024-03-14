# Overview

The declarative environment builder tool can setup Kubernetes test cluster(s) with IstioOperators, Helm Charts, and any Applications required
for testing. This is based on the declarative setup tool used in GME. 

## Configuration

Configuration is typically provided via yaml file.  Usage is usually coupled with `make` targets
since a lot of default environment variables are defined there.

### Clusters

Field is used to define one or more kubernetes clusters. The app supports creating local KinD clusters as well as using
existing clusters deployed via public cloud providers.

When using local clusters you can provide a [KinD configuration](https://github.com/kubernetes-sigs/kind/blob/v0.20.0/pkg/apis/config/v1alpha4/types.go) via `kindConfig` field to specify how you want the cluster setup.

When using existing clusters deployed in a public cloud, you can reference the cluster in the configuration via the name field. The name field should exist in the kubeconfig file as the context value (e.g: `kubectl config get-contexts -o name`)

#### Management
One of the clusters should be marked as management via `management` field.

#### Istio Components

Istio components can be defined via the `istioOperators` field. IstioOperators field is used to define one or more [Istio operators](https://github.com/istio/istio/blob/ccd8deedd9c735a11ebfa085e2fbe22be0ccd03a/operator/pkg/apis/istio/v1alpha1/types.go#L28-L46), that will be installed via istioctl.

```yaml
clusters:
- istioOperators: 
  - apiVersion: install.istio.io/v1alpha1
    kind: IstioOperator
    metadata:
      name: ingress-gateway
      namespace: istio-system
    ...
```

#### Helm Charts

Field is used to define one or more helm charts. The charts can be locally built or remote charts can be provided from a registry.

Local helm charts will use a file path via `local` to target the chart.

```yaml
clusters:
- charts: 
  - name: gloo
    namespace: gloo-mesh
    local: install/helm/gloo
    values:
      gloo:
        deployment:
          image:
            repository: gloo
    ...
```

Remote helm charts will use a url via `remote` field to target the registry.

```yaml
clusters:
- charts: 
  - name: metallb
    namespace: kube-system
    remote: https://metallb.github.io/metallb
    version: 0.13.12
    ...
```

#### Test Applications

Deploying test applications can be done via `apps` field. Each app is composed of
of ServiceAccount, Service, and Deployment every field is optional expect the
deployment.

To help with file size the versions field exist to version a deployment. When using
versions the ServiceAccount, Service and Deployment will be versioned.

```yaml
clusters:
- apps:
  - serviceAccount:
      apiVersion: v1
      kind: ServiceAccount
      metadata:
        name: httpbin
        namespace: httpbin
    service:
      apiVersion: v1
      kind: Service
      metadata:
        name: httpbin
        namespace: httpbin
        labels:
          app: httpbin
          service: httpbin
      spec:
        ports:
          - name: http
            port: 8000
            targetPort: 8080
          - name: tcp
            port: 9000
        selector:
          app: httpbin
    deployment:
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: httpbin
        namespace: httpbin
        labels:
          app: httpbin
      spec:
        replicas: 1
        selector:
          matchLabels:
            app: httpbin
        template:
          metadata:
            labels:
              app: httpbin
```
