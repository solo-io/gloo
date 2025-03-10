
---
title: "ListenerOptions"
weight: 5
---

<!-- Code generated by solo-kit. DO NOT EDIT. -->


### Package: `gloo.solo.io` 
**Types:**


- [ListenerOptions](#listeneroptions)
- [ConnectionBalanceConfig](#connectionbalanceconfig)
- [ExactBalance](#exactbalance)
  



**Source File: [github.com/solo-io/gloo/projects/gloo/api/v1/listener_options.proto](https://github.com/solo-io/gloo/blob/main/projects/gloo/api/v1/listener_options.proto)**





---
### ListenerOptions

 
Optional, feature-specific configuration that lives on gateways.
Each ListenerOption object contains configuration for a specific feature.
Note to developers: new Listener plugins must be added to this struct
to be usable by Gloo. (plugins currently need to be compiled into Gloo)

```yaml
"accessLoggingService": .als.options.gloo.solo.io.AccessLoggingService
"extensions": .gloo.solo.io.Extensions
"perConnectionBufferLimitBytes": .google.protobuf.UInt32Value
"socketOptions": []solo.io.envoy.api.v2.core.SocketOption
"proxyProtocol": .proxy_protocol.options.gloo.solo.io.ProxyProtocol
"connectionBalanceConfig": .gloo.solo.io.ConnectionBalanceConfig
"listenerAccessLoggingService": .als.options.gloo.solo.io.AccessLoggingService
"tcpStats": .google.protobuf.BoolValue

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `accessLoggingService` | [.als.options.gloo.solo.io.AccessLoggingService](../options/als/als.proto.sk/#accessloggingservice) | Configuration for access logging in a filter like the HttpConnectionManager. |
| `extensions` | [.gloo.solo.io.Extensions](../extensions.proto.sk/#extensions) | Extensions will be passed along from Listeners, Gateways, VirtualServices, Routes, and Route tables to the underlying Proxy, making them useful for controllers, validation tools, etc. which interact with kubernetes yaml. Some sample use cases: * controllers, deployment pipelines, helm charts, etc. which wish to use extensions as a kind of opaque metadata. * In the future, Gloo may support gRPC-based plugins which communicate with the Gloo translator out-of-process. Opaque Extensions enables development of out-of-process plugins without requiring recompiling & redeploying Gloo's API. |
| `perConnectionBufferLimitBytes` | [.google.protobuf.UInt32Value](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/u-int-32-value) | Soft limit on size of the listener's new connection read and write buffers. If unspecified, defaults to 1MiB For more info, check out the [Envoy docs](https://www.envoyproxy.io/docs/envoy/v1.14.1/api-v2/api/v2/listener.proto). |
| `socketOptions` | [[]solo.io.envoy.api.v2.core.SocketOption](../../../../../../solo-kit/api/external/envoy/api/v2/core/socket_option.proto.sk/#socketoption) | Additional socket options that may not be present in Envoy source code or precompiled binaries. |
| `proxyProtocol` | [.proxy_protocol.options.gloo.solo.io.ProxyProtocol](../options/proxy_protocol/proxy_protocol.proto.sk/#proxyprotocol) | Enable ProxyProtocol support for this listener. |
| `connectionBalanceConfig` | [.gloo.solo.io.ConnectionBalanceConfig](../listener_options.proto.sk/#connectionbalanceconfig) | Configuration for listener connection balancing. |
| `listenerAccessLoggingService` | [.als.options.gloo.solo.io.AccessLoggingService](../options/als/als.proto.sk/#accessloggingservice) | If enabled this sets up an early access logging service for the listener. Added initially to support listener level logging for HTTP listeners. For more info see https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/listener/v3/listener.proto#envoy-v3-api-field-config-listener-v3-listener-access-log. |
| `tcpStats` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | If true, will wrap all filter chains in the listener with a TCP stats transport socket, which is a passthrough listener that can report low-level Linux TCP stats, useful for diagnosis and triage. |




---
### ConnectionBalanceConfig

 
Configuration for listener connection balancing.

```yaml
"exactBalance": .gloo.solo.io.ConnectionBalanceConfig.ExactBalance

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `exactBalance` | [.gloo.solo.io.ConnectionBalanceConfig.ExactBalance](../listener_options.proto.sk/#exactbalance) |  |




---
### ExactBalance

 
A connection balancer implementation that does exact balancing. This means that a lock is
held during balancing so that connection counts are nearly exactly balanced between worker
threads. This is "nearly" exact in the sense that a connection might close in parallel thus
making the counts incorrect, but this should be rectified on the next accept. This balancer
sacrifices accept throughput for accuracy and should be used when there are a small number of
connections that rarely cycle (e.g., service mesh gRPC egress).

```yaml

```

| Field | Type | Description |
| ----- | ---- | ----------- | 





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
