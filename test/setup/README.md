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

Istio components can be (optionally) defined via the `istioOperators` field. 
IstioOperators field is used to define one or more [Istio operators](https://github.com/istio/istio/blob/ccd8deedd9c735a11ebfa085e2fbe22be0ccd03a/operator/pkg/apis/istio/v1alpha1/types.go#L28-L46), that will be installed via istioctl.

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

If the `istioOperators` config is not provided, the tool will skip the Istio installation. 

#### Helm Charts

Field is used to define one or more helm charts. The charts can be locally built or remote charts can be provided from a registry.

Local helm charts will use a file path via `local` to target the chart.

```yaml
clusters:
- charts: 
  - name: gloo
    namespace: gloo-system
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

#### Examples

A make target is provided to run the setup. A config file needs to be provided via `CONFIG` variable.

To test the simple gloo gateway setup:

```bash
CONFIG=test/setup/example_configs/gloo-gateway-setup.yaml make setup-declarative-env
```

To test the Istio setup, make sure the required env values are set and then run: 

```bash
export ISTIOCTL_VERSION=1.18.2
export ISTIO_HUB=docker.io/istio
export ISTIO_VERSION=1.18.2
CONFIG=test/setup/example_configs/istio-setup.yaml make setup-declarative-env
```

The make target will build the images and package the helm chart by default. To skip these steps, you can use the 
`SKIP_BUILDING_IMAGES=true` and `SKIP_PACKAGE_HELM=true` env values.

```bash
export SKIP_BUILDING_IMAGES=true
export SKIP_PACKAGE_HELM=true
```

NOTE: The setup assumes you have `KUBERNETES_MASTER=$HOME/.kube` set in your environment. You can overwrite the version set
with `export VERSION=1.0.1-dev`. `1.0.1-dev` will be used by default, but make sure the helm chart in `_output` matches
the version. If the version does not match, remove the `_output` directory and rebuild the helm chart.
