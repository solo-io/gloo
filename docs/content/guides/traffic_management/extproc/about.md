---
title: About external processing
weight: 40
description: Learn about what external processing is, how it works, and how to enable it in Gloo Gateway.
---

Envoy offers multiple filters that you can use to manage, monitor, and secure traffic to your apps. Although Envoy is extensible via C++ and WebAssembly modules, it might not be practical to implement these extensions for all of your apps. You might also have very specific requirements for how to process a request or response to allow traffic routing between different types of apps, such as adding specific headers to new and legacy apps.

With external processing, you can implement an external processing server that can read and modify all aspects of an HTTP request or response, such as headers, body, and trailers. This gives you the flexibility to apply your requirements to all types of apps, without the need to run WebAssembly or other custom scripts.

{{% notice note %}}
External processing is an Enterprise-only feature.
{{% /notice %}}

### How it works

The following diagram shows an example for how request header manipulation works when an external processing server is used.

<figure><img src="{{% versioned_link_path fromRoot="/img/extproc.svg" %}}">
<figcaption style="text-align:center;font-style:italic">External processing for request headers</figcaption></figure>

1. The downstream service sends a request with headers to the Envoy gateway.
2. The gateway extracts the header information and sends it to the external processing server.
3. The external processing server modifies, adds, or removes the request headers.
4. The modified request headers are sent back to the gateway.
5. The modified headers are added to the request.
6. The request is forwarded to the upstream application.

## ExtProc server considerations

The ExtProc server is a gRPC interface that must be able to respond to events in the lifecycle of an HTTP request. When the ExtProc filter is enabled in Gloo Gateway and a request or response is received on the gateway, the filter communicates with the ExtProc server by using bidirectional gRPC streams.

To implement your own ExtProc server, make sure that you follow [Envoy's technical specification for an external processor](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_proc/v3/ext_proc.proto#extensions-filters-http-ext-proc-v3-externalprocessor). You can also follow the [Header manipulation example]({{% versioned_link_path fromRoot="/guides/traffic_management/extproc/header-manipulation/" %}}) to try out ExtProc in Gloo Gateway with a sample ExtProc server.

{{% notice warning %}}
In Gloo Gateway version 1.17.0, the Gloo Gateway extProc filter implementation was changed to comply with the latest extProc implementation in Envoy. Previously, request and response attributes were included only in a [header processing request](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/ext_proc/v3/external_processor.proto#service-ext-proc-v3-httpheaders), and were therefore sent to the extProc server only when request header processing messages were configured to be sent. Starting in Gloo Gateway version 1.17.0, the Gloo extProc filter sends request and response attributes as part of the top-level [processing request](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/ext_proc/v3/external_processor.proto#service-ext-proc-v3-processingrequest). That way, attributes can be processed on the first processing request regardless of its type.  </br></br>

If you implemented your extProc server to expect request and response attributes as part of the HTTP header processing request, you must change this implementation to read attributes from the top-level processing request instead. </br></br>

For more information, see the [extProc proto definition](https://github.com/envoyproxy/envoy/blob/main/api/envoy/service/ext_proc/v3/external_processor.proto) in Envoy.
{{% /notice %}}

## ExtProc filter variants

Gloo Gateway supports three ExtProc filter variants that run at different positions in the Envoy filter chain. You can configure one or more variants simultaneously to process a request at multiple stages. For an overview of the filters in the Envoy filter chain, see [Filter flow description]({{% versioned_link_path fromRoot="/introduction/traffic_filter/#filter-flow-description" %}}).

| Field | Position in filter chain | Notes |
|---|---|---|
| `extProcEarly` | Early stage | Stage is configurable via the `filterStage` field. For example, choose `FaultStage` to execute the variant before or after the fault injection stage, which is the first filter in the Envoy filter chain. Depending on your filter stage setting, this variant might be executed before or after the `extProc` variant.  |
| `extProc` | Middle stage | Stage is configurable via the `filterStage` field. For example, choose `AuthNStage` to execute the variant before or after the external authentication stage. Depending on your filter stage setting, this variant might be executed before or after the `extProcEarly` variant.  |
| `extProcLate` | Final filter before a request leaves Envoy; first filter when a response enters Envoy | You must provide a `filterStage` setting. However, this setting is ignored as this variant is always executed as part of the `upstream_http_filter`. This filter is part of the Router phase, which is the last filter in the Envoy filter chain. |

Using multiple variants lets you observe what Envoy modifies between stages. For example, you can configure `extProcEarly` and `extProcLate` to log a request as it enters the filter chain and again just before it leaves Envoy, and compare the two to identify changes that Envoy made in between.

## Enable ExtProc in Gloo Gateway

You can enable any of the extProc filter variants globally for all requests and responses that the gateway processes by using the [Settings]({{% versioned_link_path fromRoot="/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/settings.proto.sk/" %}}) custom resource. Alternatively, you can enable extProc for a specific gateway listener, virtual host, or route.

The following table summarizes which fields are available at each configuration level.

| Configuration level | Available fields |
|---|---|
| Settings (global) | `extProcEarly`, `extProc`, `extProcLate` |
| HttpListenerOptions | `extProcEarly` / `disableExtProcEarly`, `extProc` / `disableExtProc`, `extProcLate` / `disableExtProcLate` |
| VirtualHostOptions | `extProcEarly`, `extProc`, `extProcLate` |
| RouteOptions | `extProcEarly`, `extProc`, `extProcLate` |

Settings defined at a lower level (listener, virtual host, or route) override the global Settings defaults via a shallow merge.

### Global settings configuration

The following example shows how to configure all three extProc filter variants in the `default` Settings resource. 

```yaml
extProcEarly:
  grpcService:
    extProcServerRef:
      name: early-ext-proc-grpc-4444
      namespace: gloo-system
  filterStage:
    stage: AuthNStage
    predicate: Before
  processingMode:
    requestHeaderMode: SEND
    responseHeaderMode: SKIP
extProc:
  grpcService:
    extProcServerRef:
      name: default-ext-proc-grpc-4444
      namespace: gloo-system
  filterStage:
    stage: AuthZStage
    predicate: After
  failureModeAllow: false
  allowModeOverride: false
  processingMode:
    requestHeaderMode: SEND
    responseHeaderMode: SKIP
extProcLate:
  grpcService:
    extProcServerRef:
      name: late-ext-proc-grpc-4444
      namespace: gloo-system
  # Filter stage must be set, but has no effect. Variant is always executed in the upstream_http_filter.
  filterStage:
    stage: AuthZStage
    predicate: After
  processingMode:
    requestHeaderMode: SEND
    responseHeaderMode: SEND
```

{{% notice note %}}
The `filterStage` field has no effect on `extProcLate`. The late filter always runs as the final filter before a request leaves Envoy and as the first filter when a response enters Envoy.
{{% /notice %}}

### Disable a variant at the listener level

To disable a specific extProc variant for a gateway listener while leaving other variants enabled globally, create an `HttpListenerOption` resource and disable the specific variant in the `options` section. The following example disables `extProcEarly` and `extProcLate` for the `gateway-proxy`, while leaving the `extProc` variant active. 

```yaml
apiVersion: gateway.solo.io/v1
kind: HttpListenerOption
metadata:
  name: disable-extproc-variants
  namespace: gloo-system
spec:
  targetRefs:
  - group: gateway.networking.k8s.io
    kind: Gateway
    name: gateway-proxy
    namespace: gloo-system
  options:
    disableExtProcEarly: true
    disableExtProcLate: true
```

### Override global settings at route level

You can override the global extProc settings, such as the processing mode or request and response attributes for a route by using the `overrides` fields as shown in the following example. 

```yaml
kubectl apply -f- <<EOF
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: vs
  namespace: gloo-system
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
            name: default-httpbin-8000
            namespace: gloo-system
      options:
        extProcEarly:
          overrides:
            processingMode:
              requestHeaderMode: SEND
              responseHeaderMode: SEND
        extProcLate:
          overrides:
            processingMode:
              requestHeaderMode: SEND
              responseHeaderMode: SEND
EOF
```

