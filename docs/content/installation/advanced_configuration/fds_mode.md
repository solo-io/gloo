---
title: Configuring Discovery
weight: 10
description: Set up automatic discovery for upstreams and functions. automatic function discovery.
---

Gloo Gateway supports automatically discovering upstreams and functions. 

{{% notice note %}}
Because discovery is a resource-intensive process, both discovery services are disabled by default.
{{% /notice %}}

As a quick reference, you can enable both UDS and FDS with the following Helm values. For more options, continue to the following sections.

```yaml
settings:
  create: true

discovery:
  enabled: true
  fdsMode: WHITELIST
```

## Upstream Discovery Service (UDS) {#uds}

Gloo Gateway can automatically create Upstream resources for Kubernetes services that it detects in your cluster. The Upstream is created in the namespace that you configure in the `discoveryNamespace` setting, which is `gloo-system` by default.

### Upstream names

The name of the Upstream is based on the name, namespace, and port of the service that it refers to, in the following format: 

```
<service-namespace>-<service-name>-<service-port>
```

For example, if you have a Pet Store service in the `default` namespace that listens on port `8080`, the discovered Upstream gets the follwing name:

```
default-petstore-8080
```

To list discovered Upstreams, run the following command:

```
glooctl get upstreams -n gloo-system
```

Example output:

```
+-----------------------+------------+----------+-------------------------+
|       UPSTREAM        |    TYPE    |  STATUS  |         DETAILS         |
+-----------------------+------------+----------+-------------------------+
| default-petstore-8080 | Kubernetes | Accepted | svc name:      petstore |
|                       |            |          | svc namespace: default  |
|                       |            |          | port:          8080     |
|                       |            |          | REST service:           |
|                       |            |          | functions:              |
|                       |            |          | - addPet                |
|                       |            |          | - deletePet             |
|                       |            |          | - findPetById           |
|                       |            |          | - findPets              |
|                       |            |          |                         |
+-----------------------+------------+----------+-------------------------+
```


### Enable UDS mode

You can enable UDS mode with Helm or the Settings CR in the `gloo-system` namespace.

{{< tabs >}} 
{{% tab name="Helm values file" %}}
Add the `discovery.enabled=true` setting to your Helm overrides file:

```yaml
settings:
  create: true

discovery:
  enabled: true
```
{{% /tab %}}

{{% tab name="Helm CLI" %}}
Add the following CLI flags to `helm install|template` commands:

```bash
helm install|template ... --set settings.create=true --set discovery.enabled=true
```
{{% /tab %}}

{{% tab name="Settings CR" %}}
Enable `fdsMode` by editing the `gloo.solo.io/v1.Settings` custom resource to add the `discovery.enabled=true` setting.

```bash
kubectl edit -n gloo-system settings.gloo.solo.io
```

Example Settings CR:
```yaml
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: gloo-system
spec:
  bindAddr: 0.0.0.0:9977
  discoveryNamespace: gloo-system
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  # add the following lines
  discovery:
    enabled: true
```
{{% /tab %}} 
{{< /tabs >}}

### More information about UDS {#more-info}

For more information, check out the following guides:

* [Discovered Upstreams]({{< versioned_link_path fromRoot="/guides/traffic_management/destination_types/discovered_upstream/" >}})
* [Configuration via Annotations]({{< versioned_link_path fromRoot="/guides/traffic_management/destination_types/discovered_upstream/discovered-upstream-configuration/" >}})

---

## Function Discovery Service (FDS) {#fds}

The **Function Discovery Service** (FDS) discovers and understands the OpenAPI Spec (OAS, formerly known as Swagger) or gRPC endpoints of the Upstreams in your cluster.

### Evaluated endpoints

When enabled, FDS attempts to poll endpoints for:

* A path serving an [OpenAPI (Swagger) document](https://swagger.io/specification/).
* gRPC Services with [gRPC Reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) enabled.

The default endpoints evaluated for `swagger` or `OpenAPISpec` docs are:

```
"/openapi.json"
"/swagger.json"
"/swagger/docs/v1"
"/swagger/docs/v2"
"/v1/swagger"
"/v2/swagger"
```

If you have an OpenAPI definition on a different endpoint, you can customize the location by configuring it in the `serviceSpec.rest.swaggerInfo.url` field. For example, for a given Upstream, you can add the following including an explicit location for the OpenAPI document:


{{< highlight yaml "hl_lines=17-20" >}}
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  labels:
    discovered_by: kubernetesplugin
    service: petstore
  name: default-petstore-8080
  namespace: gloo-system
spec:
  discoveryMetadata: {}
  kube:
    selector:
      app: petstore
    serviceName: petstore
    serviceNamespace: default
    servicePort: 8080
    serviceSpec:
      rest:
        swaggerInfo:
          url: http://petstore.default:8080/foo/bar/swagger.json

{{< /highlight >}}

### How FDS works

Using FDS means that the Gloo Gateway `discovery` component will make HTTP requests to all `Upstreams` known to Gloo Gateway trying to discover functions. This behavior causes increased network traffic and may be undesirable if it causes unexpected behavior or logs to appear in the services Gloo Gateway is attempting to poll. For this reason, we may want to restrict the manner in which FDS polls services.

Gloo Gateway allows whitelisting/blacklisting services, either by namespace or on the individual service level.

We can use these configuration settings to restrict FDS to discover only the namespaces or individual services we choose.

### Enable FDS mode

You can enable FDS mode with Helm or the Settings CR in the `gloo-system` namespace.

**Enterprise-only**: Automated schema generation for GraphQL is enabled by default. This can be disabled by modifying the `gloo.solo.io/v1.Settings` custom resource as shown in the following tabs.

{{< tabs >}} 
{{% tab name="Helm values file" %}}
Add the `discovery.fdsMode` setting to your Helm overrides file:

```yaml
settings:
  create: true

discovery:
  # set to either WHITELIST, BLACKLIST, or DISABLED
  # WHITELIST is the default value
  fdsMode: BLACKLIST
```
{{% /tab %}}

{{% tab name="Helm CLI" %}}
Add the following CLI flags to `helm install|template` commands:

```bash
helm install|template ... --set settings.create=true --set discovery.fdsMode=BLACKLIST
```
{{% /tab %}}

{{% tab name="Settings CR" %}}
Enable `fdsMode` by editing the `gloo.solo.io/v1.Settings` custom resource to add the `discovery.fdsMode` setting.

```bash
kubectl edit -n gloo-system settings.gloo.solo.io
```

Example Settings CR:
```yaml
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: gloo-system
spec:
  bindAddr: 0.0.0.0:9977
  discoveryNamespace: gloo-system
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  # add the following lines
  discovery:
    # set to either WHITELIST, BLACKLIST, or DISABLED
    # WHITELIST is the default value
    fdsMode: WHITELIST
    fdsOptions:
      # Enterprise-only: set to false to disable automated GraphQL schema generation as part of FDS.
      # true is the default value (enabled)
      graphqlEnabled: true
```
{{% /tab %}} 
{{< /tabs >}}

### FDS modes

FDS can run in one of three modes:

* `BLACKLIST`: The most liberal FDS polling policy. FDS polls all services, unless the service or the service's namespace is blacklisted.
* `WHITELIST` (default mode): A more restrictive FDS polling policy. FDS polls only services that are either whitelisted or exist in a whitelisted namespace.
* `DISABLED`: FDS does not run (the default mode). UDS continues to run as normal, if enabled.

### Blacklisting namespaces and upstreams

When running in `BLACKLIST` mode, blacklist upstreams by adding the disabled label. This label can be applied to namespaces and upstreams.

Label:

`discovery.solo.io/function_discovery=disabled`.

Example commands:

```bash
kubectl label namespace default discovery.solo.io/function_discovery=disabled
kubectl label upstream -n myapp myupstream discovery.solo.io/function_discovery=disabled

# if the Upstream was discovered and is managed by UDS
# then you can add the label to the service, which propagates to the Upstream
kubectl label service -n myapp myservice discovery.solo.io/function_discovery=disabled
```

To enable FDS for specific upstreams in a blacklisted namespace, use the enabled label:

`discovery.solo.io/function_discovery=enabled`

Example command:

```bash
kubectl label upstream -n default upstream-fds discovery.solo.io/function_discovery=enabled
```

### Whitelisting namespaces and upstreams

When running in `WHITELIST` mode, whitelist `Upstreams` by adding the enabled label. This label can be applied to namespaces and upstreams.

`discovery.solo.io/function_discovery=enabled`.

Example commands:

```bash
kubectl label namespace default discovery.solo.io/function_discovery=enabled
kubectl label upstream -n myapp myupstream discovery.solo.io/function_discovery=enabled

# if the Upstream was discovered and is managed by UDS
# then you can add the label to the service, which propagates to the Upstream
kubectl label service -n myapp myservice discovery.solo.io/function_discovery=enabled
```

To disable FDS for specific services/upstreams in a whitelisted namespace, use the disabled label:

`discovery.solo.io/function_discovery=disabled`

Example command:

```bash
kubectl label upstream -n default upstream-fds discovery.solo.io/function_discovery=disabled
```
