
---
title: "Upstream"
weight: 5
---

<!-- Code generated by solo-kit. DO NOT EDIT. -->


### Package: `gloo.solo.io` 
**Types:**


- [Upstream](#upstream) **Top-Level Resource**
- [ClusterProtocolSelection](#clusterprotocolselection)
- [DiscoveryMetadata](#discoverymetadata)
- [HeaderValue](#headervalue)
- [PreconnectPolicy](#preconnectpolicy)
  



**Source File: [github.com/solo-io/gloo/projects/gloo/api/v1/upstream.proto](https://github.com/solo-io/gloo/blob/main/projects/gloo/api/v1/upstream.proto)**





---
### Upstream

 
Upstreams represent destination for routing HTTP requests. Upstreams can be compared to
[clusters](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/cluster.proto) in Envoy terminology.
Each upstream in Gloo has a type. Supported types include `static`, `kubernetes`, `aws`, `consul`, and more.
Each upstream type is handled by a corresponding Gloo plugin. (plugins currently need to be compiled into Gloo)

```yaml
"namespacedStatuses": .core.solo.io.NamespacedStatuses
"metadata": .core.solo.io.Metadata
"discoveryMetadata": .gloo.solo.io.DiscoveryMetadata
"sslConfig": .gloo.solo.io.UpstreamSslConfig
"circuitBreakers": .gloo.solo.io.CircuitBreakerConfig
"loadBalancerConfig": .gloo.solo.io.LoadBalancerConfig
"healthChecks": []solo.io.envoy.api.v2.core.HealthCheck
"outlierDetection": .solo.io.envoy.api.v2.cluster.OutlierDetection
"kube": .kubernetes.options.gloo.solo.io.UpstreamSpec
"static": .static.options.gloo.solo.io.UpstreamSpec
"pipe": .pipe.options.gloo.solo.io.UpstreamSpec
"aws": .aws.options.gloo.solo.io.UpstreamSpec
"azure": .azure.options.gloo.solo.io.UpstreamSpec
"consul": .consul.options.gloo.solo.io.UpstreamSpec
"awsEc2": .aws_ec2.options.gloo.solo.io.UpstreamSpec
"gcp": .gcp.options.gloo.solo.io.UpstreamSpec
"ai": .ai.options.gloo.solo.io.UpstreamSpec
"failover": .gloo.solo.io.Failover
"connectionConfig": .gloo.solo.io.ConnectionConfig
"protocolSelection": .gloo.solo.io.Upstream.ClusterProtocolSelection
"useHttp2": .google.protobuf.BoolValue
"initialStreamWindowSize": .google.protobuf.UInt32Value
"initialConnectionWindowSize": .google.protobuf.UInt32Value
"maxConcurrentStreams": .google.protobuf.UInt32Value
"overrideStreamErrorOnInvalidHttpMessage": .google.protobuf.BoolValue
"httpProxyHostname": .google.protobuf.StringValue
"httpConnectSslConfig": .gloo.solo.io.UpstreamSslConfig
"httpConnectHeaders": []gloo.solo.io.HeaderValue
"ignoreHealthOnHostRemoval": .google.protobuf.BoolValue
"respectDnsTtl": .google.protobuf.BoolValue
"dnsRefreshRate": .google.protobuf.Duration
"proxyProtocolVersion": .google.protobuf.StringValue
"preconnectPolicy": .gloo.solo.io.PreconnectPolicy
"disableIstioAutoMtls": .google.protobuf.BoolValue

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `namespacedStatuses` | [.core.solo.io.NamespacedStatuses](../../../../../../solo-kit/api/v1/status.proto.sk/#namespacedstatuses) | NamespacedStatuses indicates the validation status of this resource. NamespacedStatuses is read-only by clients, and set by gloo during validation. |
| `metadata` | [.core.solo.io.Metadata](../../../../../../solo-kit/api/v1/metadata.proto.sk/#metadata) | Metadata contains the object metadata for this resource. |
| `discoveryMetadata` | [.gloo.solo.io.DiscoveryMetadata](../upstream.proto.sk/#discoverymetadata) | Upstreams and their configuration can be automatically by Gloo Discovery if this upstream is created or modified by Discovery, metadata about the operation will be placed here. |
| `sslConfig` | [.gloo.solo.io.UpstreamSslConfig](../ssl/ssl.proto.sk/#upstreamsslconfig) | SslConfig contains the options necessary to configure envoy to originate TLS to an upstream. |
| `circuitBreakers` | [.gloo.solo.io.CircuitBreakerConfig](../circuit_breaker/circuit_breaker.proto.sk/#circuitbreakerconfig) | Circuit breakers for this upstream. if not set, the defaults ones from the Gloo settings will be used. if those are not set, [envoy's defaults](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/circuit_breaker.proto#envoy-api-msg-cluster-circuitbreakers) will be used. |
| `loadBalancerConfig` | [.gloo.solo.io.LoadBalancerConfig](../load_balancer.proto.sk/#loadbalancerconfig) | Settings for the load balancer that sends requests to the Upstream. The load balancing method is set to round robin by default. |
| `healthChecks` | []solo.io.envoy.api.v2.core.HealthCheck |  |
| `outlierDetection` | .solo.io.envoy.api.v2.cluster.OutlierDetection |  |
| `kube` | [.kubernetes.options.gloo.solo.io.UpstreamSpec](../options/kubernetes/kubernetes.proto.sk/#upstreamspec) |  Only one of `kube`, `static`, `pipe`, `aws`, `azure`, `consul`, `awsEc2`, `gcp`, or `ai` can be set. |
| `static` | [.static.options.gloo.solo.io.UpstreamSpec](../options/static/static.proto.sk/#upstreamspec) |  Only one of `static`, `kube`, `pipe`, `aws`, `azure`, `consul`, `awsEc2`, `gcp`, or `ai` can be set. |
| `pipe` | [.pipe.options.gloo.solo.io.UpstreamSpec](../options/pipe/pipe.proto.sk/#upstreamspec) |  Only one of `pipe`, `kube`, `static`, `aws`, `azure`, `consul`, `awsEc2`, `gcp`, or `ai` can be set. |
| `aws` | [.aws.options.gloo.solo.io.UpstreamSpec](../options/aws/aws.proto.sk/#upstreamspec) |  Only one of `aws`, `kube`, `static`, `pipe`, `azure`, `consul`, `awsEc2`, `gcp`, or `ai` can be set. |
| `azure` | [.azure.options.gloo.solo.io.UpstreamSpec](../options/azure/azure.proto.sk/#upstreamspec) |  Only one of `azure`, `kube`, `static`, `pipe`, `aws`, `consul`, `awsEc2`, `gcp`, or `ai` can be set. |
| `consul` | [.consul.options.gloo.solo.io.UpstreamSpec](../options/consul/consul.proto.sk/#upstreamspec) |  Only one of `consul`, `kube`, `static`, `pipe`, `aws`, `azure`, `awsEc2`, `gcp`, or `ai` can be set. |
| `awsEc2` | [.aws_ec2.options.gloo.solo.io.UpstreamSpec](../options/aws/ec2/aws_ec2.proto.sk/#upstreamspec) |  Only one of `awsEc2`, `kube`, `static`, `pipe`, `aws`, `azure`, `consul`, `gcp`, or `ai` can be set. |
| `gcp` | [.gcp.options.gloo.solo.io.UpstreamSpec](../enterprise/options/gcp/gcp.proto.sk/#upstreamspec) |  Only one of `gcp`, `kube`, `static`, `pipe`, `aws`, `azure`, `consul`, `awsEc2`, or `ai` can be set. |
| `ai` | [.ai.options.gloo.solo.io.UpstreamSpec](../enterprise/options/ai/ai.proto.sk/#upstreamspec) |  Only one of `ai`, `kube`, `static`, `pipe`, `aws`, `azure`, `consul`, `awsEc2`, or `gcp` can be set. |
| `failover` | [.gloo.solo.io.Failover](../failover.proto.sk/#failover) | Failover endpoints for this upstream. If omitted (the default) no failovers will be applied. |
| `connectionConfig` | [.gloo.solo.io.ConnectionConfig](../connection.proto.sk/#connectionconfig) | HTTP/1 connection configurations. |
| `protocolSelection` | [.gloo.solo.io.Upstream.ClusterProtocolSelection](../upstream.proto.sk/#clusterprotocolselection) | Determines how Envoy selects the protocol used to speak to upstream hosts. |
| `useHttp2` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | Use http2 when communicating with this upstream this field is evaluated `true` for upstreams with a grpc service spec. otherwise defaults to `false`. |
| `initialStreamWindowSize` | [.google.protobuf.UInt32Value](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/u-int-32-value) | (UInt32Value) Initial stream-level flow-control window size. Valid values range from 65535 (2^16 - 1, HTTP/2 default) to 2147483647 (2^31 - 1, HTTP/2 maximum) and defaults to 268435456 (256 * 1024 * 1024). NOTE: 65535 is the initial window size from HTTP/2 spec. We only support increasing the default window size now, so it’s also the minimum. This field also acts as a soft limit on the number of bytes Envoy will buffer per-stream in the HTTP/2 codec buffers. Once the buffer reaches this pointer, watermark callbacks will fire to stop the flow of data to the codec buffers. Requires UseHttp2 to be true to be acknowledged. |
| `initialConnectionWindowSize` | [.google.protobuf.UInt32Value](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/u-int-32-value) | (UInt32Value) Similar to initial_stream_window_size, but for connection-level flow-control window. Currently, this has the same minimum/maximum/default as initial_stream_window_size. Requires UseHttp2 to be true to be acknowledged. |
| `maxConcurrentStreams` | [.google.protobuf.UInt32Value](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/u-int-32-value) | (UInt32Value) Maximum concurrent streams allowed for peer on one HTTP/2 connection. Valid values range from 1 to 2147483647 (2^31 - 1) and defaults to 2147483647. Requires UseHttp2 to be true to be acknowledged. |
| `overrideStreamErrorOnInvalidHttpMessage` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | Allows invalid HTTP messaging and headers. When this option is disabled (default), then the whole HTTP/2 connection is terminated upon receiving invalid HEADERS frame. However, when this option is enabled, only the offending stream is terminated. This overrides any HCM :ref:`stream_error_on_invalid_http_messaging <envoy_v3_api_field_extensions.filters.network.http_connection_manager.v3.HttpConnectionManager.stream_error_on_invalid_http_message>` See [RFC7540, sec. 8.1](https://datatracker.ietf.org/doc/html/rfc7540#section-8.1) for details. |
| `httpProxyHostname` | [.google.protobuf.StringValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/string-value) | Tells Envoy that the upstream is an HTTP proxy that supports [HTTP CONNECT method](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT). The hostname is the destination of the tunnel established by the proxy. Some Envoy Command Operators (.e.g `%REQUESTED_SERVER_NAME%`) are supported allowing for dynamic destinations. For example, setting to: host.com:443 and making a request routed to the upstream such as `curl <envoy>:<port>/v1` would result in the following request: CONNECT host.com:443 HTTP/1.1 host: host.com:443 GET /v1 HTTP/1.1 host: <envoy>:<port> user-agent: curl/7.64.1 accept: */* Note: If setting this field to a hostname rather than IP:PORT, you may want to also set `host_rewrite` on the route. |
| `httpConnectSslConfig` | [.gloo.solo.io.UpstreamSslConfig](../ssl/ssl.proto.sk/#upstreamsslconfig) | HttpConnectSslConfig contains the options necessary to configure envoy to originate TLS to an HTTP Connect proxy. If you also want to ensure the bytes proxied by the HTTP Connect proxy are encrypted, you should also specify `ssl_config`. |
| `httpConnectHeaders` | [[]gloo.solo.io.HeaderValue](../upstream.proto.sk/#headervalue) | HttpConnectHeaders specifies the headers sent with the initial HTTP Connect request. |
| `ignoreHealthOnHostRemoval` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | (bool) If set to true, Envoy will ignore the health value of a host when processing its removal from service discovery. This means that if active health checking is used, Envoy will not wait for the endpoint to go unhealthy before removing it. |
| `respectDnsTtl` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | If set to true, Service Discovery update period will be triggered once the TTL is expired. If minimum TTL of all records is 0 then dns_refresh_rate will be used. |
| `dnsRefreshRate` | [.google.protobuf.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration) | Service Discovery DNS Refresh Rate. Minimum value is 1 ms. Values below the minimum are considered invalid. Only valid for STRICT_DNS and LOGICAL_DNS cluster types. All other cluster types are considered invalid. |
| `proxyProtocolVersion` | [.google.protobuf.StringValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/string-value) | Proxy Protocol Version to add when communicating with the upstream. If unset will not wrap the transport socket. These are of the format "V1" or "V2". |
| `preconnectPolicy` | [.gloo.solo.io.PreconnectPolicy](../upstream.proto.sk/#preconnectpolicy) | Preconnect policy for the cluster Aligns as closely as possible with https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/cluster.proto#envoy-v3-api-msg-config-cluster-v3-cluster-preconnectpolicy This is not recommended for use unless you are sure you need it. In most cases preconnect hurts more than it helps. |
| `disableIstioAutoMtls` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | If set to true, the proxy will not allow automatic mTLS detection for Istio upstreams. Defaults to false. |




---
### ClusterProtocolSelection



| Name | Description |
| ----- | ----------- | 
| `USE_CONFIGURED_PROTOCOL` | Cluster can only operate on one of the possible upstream protocols (HTTP1.1, HTTP2). If http2_protocol_options are present, HTTP2 will be used, otherwise HTTP1.1 will be used. |
| `USE_DOWNSTREAM_PROTOCOL` | Use HTTP1.1 or HTTP2, depending on which one is used on the downstream connection. |




---
### DiscoveryMetadata

 
created by discovery services

```yaml
"labels": map<string, string>

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `labels` | `map<string, string>` | Labels inherited from the original upstream (e.g. Kubernetes labels). |




---
### HeaderValue

 
Header name/value pair.

```yaml
"key": string
"value": string

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `key` | `string` | Header name. |
| `value` | `string` | Header value. |




---
### PreconnectPolicy



```yaml
"perUpstreamPreconnectRatio": .google.protobuf.DoubleValue
"predictivePreconnectRatio": .google.protobuf.DoubleValue

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `perUpstreamPreconnectRatio` | [.google.protobuf.DoubleValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/double-value) | Indicates how many streams (rounded up) can be anticipated per-upstream for each incoming stream. This is useful for high-QPS or latency-sensitive services. Preconnecting will only be done if the upstream is healthy and the cluster has traffic. For example if this is 2, for an incoming HTTP/1.1 stream, 2 connections will be established, one for the new incoming stream, and one for a presumed follow-up stream. For HTTP/2, only one connection would be established by default as one connection can serve both the original and presumed follow-up stream. In steady state for non-multiplexed connections a value of 1.5 would mean if there were 100 active streams, there would be 100 connections in use, and 50 connections preconnected. This might be a useful value for something like short lived single-use connections, for example proxying HTTP/1.1 if keep-alive were false and each stream resulted in connection termination. It would likely be overkill for long lived connections, such as TCP proxying SMTP or regular HTTP/1.1 with keep-alive. For long lived traffic, a value of 1.05 would be more reasonable, where for every 100 connections, 5 preconnected connections would be in the queue in case of unexpected disconnects where the connection could not be reused. If this value is not set, or set explicitly to one, Envoy will fetch as many connections as needed to serve streams in flight. This means in steady state if a connection is torn down, a subsequent streams will pay an upstream-rtt latency penalty waiting for a new connection. This is limited somewhat arbitrarily to 3 because preconnecting too aggressively can harm latency more than the preconnecting helps. |
| `predictivePreconnectRatio` | [.google.protobuf.DoubleValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/double-value) | Indicates how many streams (rounded up) can be anticipated across a cluster for each stream, useful for low QPS services. This is currently supported for a subset of deterministic non-hash-based load-balancing algorithms (weighted round robin, random). Unlike `per_upstream_preconnect_ratio` this preconnects across the upstream instances in a cluster, doing best effort predictions of what upstream would be picked next and pre-establishing a connection. Preconnecting will be limited to one preconnect per configured upstream in the cluster and will only be done if there are healthy upstreams and the cluster has traffic. For example if preconnecting is set to 2 for a round robin HTTP/2 cluster, on the first incoming stream, 2 connections will be preconnected - one to the first upstream for this cluster, one to the second on the assumption there will be a follow-up stream. If this value is not set, or set explicitly to one, Envoy will fetch as many connections as needed to serve streams in flight, so during warm up and in steady state if a connection is closed (and per_upstream_preconnect_ratio is not set), there will be a latency hit for connection establishment. If both this and preconnect_ratio are set, Envoy will make sure both predicted needs are met, basically preconnecting max(predictive-preconnect, per-upstream-preconnect), for each upstream. |





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
