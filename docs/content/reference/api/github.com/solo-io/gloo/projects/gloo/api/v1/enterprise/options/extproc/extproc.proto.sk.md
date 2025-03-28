
---
title: "Extproc"
weight: 5
---

<!-- Code generated by solo-kit. DO NOT EDIT. -->


### Package: `extproc.options.gloo.solo.io` 
**Types:**


- [Settings](#settings)
- [RouteSettings](#routesettings)
- [GrpcService](#grpcservice)
- [Overrides](#overrides)
- [HeaderForwardingRules](#headerforwardingrules)
  



**Source File: [github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extproc/extproc.proto](https://github.com/solo-io/gloo/blob/main/projects/gloo/api/v1/enterprise/options/extproc/extproc.proto)**





---
### Settings

 
Enterprise-only: Configuration for Envoy's [External Processing Filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_proc_filter).
The External Processing filter allows for calling out to an external gRPC service at a specified
point within a HTTP filter chain. The external service may access and modify various parts of the
request or response, and may terminate processing.
Envoy's External Processing Filter is considered a work in progress and has an unknown security posture.
Users should take care to understand the risks of using this extension before proceeding.

```yaml
"grpcService": .extproc.options.gloo.solo.io.GrpcService
"filterStage": .filters.gloo.solo.io.FilterStage
"failureModeAllow": .google.protobuf.BoolValue
"processingMode": .solo.io.envoy.extensions.filters.http.ext_proc.v3.ProcessingMode
"asyncMode": .google.protobuf.BoolValue
"requestAttributes": []string
"responseAttributes": []string
"messageTimeout": .google.protobuf.Duration
"statPrefix": .google.protobuf.StringValue
"mutationRules": .solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules
"maxMessageTimeout": .google.protobuf.Duration
"disableClearRouteCache": .google.protobuf.BoolValue
"forwardRules": .extproc.options.gloo.solo.io.HeaderForwardingRules
"filterMetadata": .google.protobuf.Struct
"allowModeOverride": .google.protobuf.BoolValue
"metadataContextNamespaces": []string
"typedMetadataContextNamespaces": []string

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `grpcService` | [.extproc.options.gloo.solo.io.GrpcService](../extproc.proto.sk/#grpcservice) | Required. Configuration for the gRPC service that the filter will communicate with. |
| `filterStage` | [.filters.gloo.solo.io.FilterStage](../../../../filters/stages.proto.sk/#filterstage) | Required. Where in the HTTP filter chain to insert the filter. |
| `failureModeAllow` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | By default, if the gRPC stream cannot be established, or if it is closed prematurely with an error, the filter will fail. Specifically, if the response headers have not yet been delivered, then it will return a 500 error downstream. If they have been delivered, then instead the HTTP stream to the downstream client will be reset. With this parameter set to true, however, then if the gRPC stream is prematurely closed or could not be opened, processing continues without error. |
| `processingMode` | .solo.io.envoy.extensions.filters.http.ext_proc.v3.ProcessingMode | Specifies default options for how HTTP headers, trailers, and bodies are sent. |
| `asyncMode` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | NOT CURRENTLY IMPLEMENTED. If true, send each part of the HTTP request or response specified by ProcessingMode asynchronously -- in other words, send the message on the gRPC stream and then continue filter processing. If false, which is the default, suspend filter execution after each message is sent to the remote service and wait up to "message_timeout" for a reply. |
| `requestAttributes` | `[]string` | NOT CURRENTLY IMPLEMENTED. Envoy provides a number of [attributes](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#arch-overview-attributes) for expressive policies. Each attribute name provided in this field will be matched against that list and populated in the request_headers message. See the [request attribute documentation](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#request-attributes) for the list of supported attributes and their types. |
| `responseAttributes` | `[]string` | NOT CURRENTLY IMPLEMENTED. Envoy provides a number of [attributes](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#arch-overview-attributes) for expressive policies. Each attribute name provided in this field will be matched against that list and populated in the response_headers message. See the [response attribute documentation](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#response-attributes) for the list of supported attributes and their types. |
| `messageTimeout` | [.google.protobuf.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration) | Specifies the timeout for each individual message sent on the stream when the filter is running in synchronous mode. Whenever the proxy sends a message on the stream that requires a response, it will reset this timer, and will stop processing and return an error (subject to the processing mode) if the timer expires before a matching response is received. There is no timeout when the filter is running in asynchronous mode. Value must be at least 0 seconds, and less than or equal to 3600 seconds. Zero is a valid value which means the timer will be triggered immediately. If not configured, default is 200 milliseconds. |
| `statPrefix` | [.google.protobuf.StringValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/string-value) | Optional additional prefix to use when emitting statistics. This allows distinguishing between statistics emitted by multiple *ext_proc* filters in an HTTP filter chain. |
| `mutationRules` | .solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules | Rules that determine what modifications an external processing server may make to message headers. If not set, all headers may be modified except for "host", ":authority", ":scheme", ":method", and headers that start with the header prefix set via [header_prefix](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/bootstrap/v3/bootstrap.proto#envoy-v3-api-field-config-bootstrap-v3-bootstrap-header-prefix) (which is usually "x-envoy"). Note that changing headers such as "host" or ":authority" may not in itself change Envoy's routing decision, as routes can be cached. To also force the route to be recomputed, set the [clear_route_cache](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/ext_proc/v3/external_processor.proto#envoy-v3-api-field-service-ext-proc-v3-commonresponse-clear-route-cache) field to true in the same response. |
| `maxMessageTimeout` | [.google.protobuf.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration) | Specify the upper bound of [override_message_timeout](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/ext_proc/v3/external_processor.proto#envoy-v3-api-field-service-ext-proc-v3-processingresponse-override-message-timeout). If not specified, by default it is 0, which will effectively disable the `override_message_timeout` API. Value must be greater than or equal to the `messageTimeout` and less than or equal to 3600 seconds. |
| `disableClearRouteCache` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | Prevents clearing the route-cache when the [clear_route_cache](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/ext_proc/v3/external_processor.proto#envoy-v3-api-field-service-ext-proc-v3-commonresponse-clear-route-cache) field is set in an external processor response. |
| `forwardRules` | [.extproc.options.gloo.solo.io.HeaderForwardingRules](../extproc.proto.sk/#headerforwardingrules) | Allow headers matching the `forward_rules` to be forwarded to the external processing server. If not set, all headers are forwarded to the external processing server. |
| `filterMetadata` | [.google.protobuf.Struct](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/struct) | Additional metadata to be added to the filter state for logging purposes. The metadata will be added to StreamInfo's filter state under the namespace corresponding to the ext_proc filter name. |
| `allowModeOverride` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | If `allow_mode_override` is set to true, the filter config [processing_mode](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_proc/v3/ext_proc.proto#envoy-v3-api-field-extensions-filters-http-ext-proc-v3-externalprocessor-processing-mode) can be overridden by the response message from the external processing server [mode_override](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/ext_proc/v3/external_processor.proto#envoy-v3-api-field-service-ext-proc-v3-processingresponse-mode-override). If not set, `mode_override` API in the response message will be ignored. |
| `metadataContextNamespaces` | `[]string` | Specifies a list of metadata namespaces whose values, if present, will be passed to the ext_proc service as an opaque *protobuf::Struct*. |
| `typedMetadataContextNamespaces` | `[]string` | Specifies a list of metadata namespaces whose values, if present, will be passed to the ext_proc service. typed_filter_metadata is passed as an `protobuf::Any`. It works in a way similar to `metadata_context_namespaces` but allows envoy and external processing server to share the protobuf message definition in order to do a safe parsing. |




---
### RouteSettings

 
External processor settings that can be configured on a virtual host or route.

```yaml
"disabled": .google.protobuf.BoolValue
"overrides": .extproc.options.gloo.solo.io.Overrides

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `disabled` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | Set to true to disable the External Processing filter for this virtual host or route. Setting this value to false is not supported. Only one of `disabled` or `overrides` can be set. |
| `overrides` | [.extproc.options.gloo.solo.io.Overrides](../extproc.proto.sk/#overrides) | Override specific configuration for this virtual host or route. If a route specifies overrides, it will override the disabled flag of its parent virtual host. Only one of `overrides` or `disabled` can be set. |




---
### GrpcService



```yaml
"extProcServerRef": .core.solo.io.ResourceRef
"authority": .google.protobuf.StringValue
"retryPolicy": .solo.io.envoy.config.core.v3.RetryPolicy
"timeout": .google.protobuf.Duration
"initialMetadata": []solo.io.envoy.config.core.v3.HeaderValue

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `extProcServerRef` | [.core.solo.io.ResourceRef](../../../../../../../../../solo-kit/api/v1/ref.proto.sk/#resourceref) | A reference to the Upstream representing the external processor gRPC server. See https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/ext_proc/v3/external_processor.proto for details. |
| `authority` | [.google.protobuf.StringValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/string-value) | The `:authority` header in the grpc request. If this field is not set, the authority header value will be the cluster name. Note that this authority does not override the SNI. The SNI is provided by the transport socket of the cluster. |
| `retryPolicy` | .solo.io.envoy.config.core.v3.RetryPolicy | Indicates the retry policy for re-establishing the gRPC stream This field is optional. If max interval is not provided, it will be set to ten times the provided base interval. Currently only supported for xDS gRPC streams. If not set, xDS gRPC streams default base interval:500ms, maximum interval:30s will be applied. |
| `timeout` | [.google.protobuf.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration) | The timeout for the gRPC request. This is the timeout for a specific request. |
| `initialMetadata` | []solo.io.envoy.config.core.v3.HeaderValue | Additional metadata to include in streams initiated to the GrpcService. This can be used for scenarios in which additional ad hoc authorization headers (e.g. `x-foo-bar: baz-key`) are to be injected. For more information, including details on header value syntax, see the documentation on [custom request headers](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#config-http-conn-man-headers-custom-request-headers). |




---
### Overrides



```yaml
"processingMode": .solo.io.envoy.extensions.filters.http.ext_proc.v3.ProcessingMode
"asyncMode": .google.protobuf.BoolValue
"requestAttributes": []string
"responseAttributes": []string
"grpcService": .extproc.options.gloo.solo.io.GrpcService
"metadataContextNamespaces": []string
"typedMetadataContextNamespaces": []string

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `processingMode` | .solo.io.envoy.extensions.filters.http.ext_proc.v3.ProcessingMode | Set a different processing mode for this virtual host or route than the default. |
| `asyncMode` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | NOT CURRENTLY IMPLEMENTED. Set a different asynchronous processing option than the default. |
| `requestAttributes` | `[]string` | NOT FINALIZED UPSTREAM use at your own upgrade risk Set different optional attributes than the default setting of the `request_attributes` field. |
| `responseAttributes` | `[]string` | NOT FINALIZED UPSTREAM use at your own upgrade risk Set different optional properties than the default setting of the `response_attributes` field. |
| `grpcService` | [.extproc.options.gloo.solo.io.GrpcService](../extproc.proto.sk/#grpcservice) | Set a different gRPC service for this virtual host or route than the default. |
| `metadataContextNamespaces` | `[]string` | Specifies a list of metadata namespaces whose values, if present, will be passed to the ext_proc service as an opaque *protobuf::Struct*. |
| `typedMetadataContextNamespaces` | `[]string` | Specifies a list of metadata namespaces whose values, if present, will be passed to the ext_proc service. typed_filter_metadata is passed as an `protobuf::Any`. It works in a way similar to `metadata_context_namespaces` but allows envoy and external processing server to share the protobuf message definition in order to do a safe parsing. |




---
### HeaderForwardingRules

 
The HeaderForwardingRules structure specifies what headers are
allowed to be forwarded to the external processing server.
See https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_proc/v3/ext_proc.proto#extensions-filters-http-ext-proc-v3-headerforwardingrules
for details.

```yaml
"allowedHeaders": .solo.io.envoy.type.matcher.v3.ListStringMatcher
"disallowedHeaders": .solo.io.envoy.type.matcher.v3.ListStringMatcher

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `allowedHeaders` | .solo.io.envoy.type.matcher.v3.ListStringMatcher | If set, specifically allow any header in this list to be forwarded to the external processing server. This can be overridden by the below `disallowed_headers`. |
| `disallowedHeaders` | .solo.io.envoy.type.matcher.v3.ListStringMatcher | If set, specifically disallow any header in this list to be forwarded to the external processing server. This overrides the above `allowed_headers` if a header matches both. NOT CURRENTLY IMPLEMENTED. |





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
