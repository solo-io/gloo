
---
title: "Protocol"
weight: 5
---

<!-- Code generated by solo-kit. DO NOT EDIT. -->


### Package: `protocol.options.gloo.solo.io` 
**Types:**


- [HttpProtocolOptions](#httpprotocoloptions)
- [HeadersWithUnderscoresAction](#headerswithunderscoresaction)
- [Http1ProtocolOptions](#http1protocoloptions)
- [Http2ProtocolOptions](#http2protocoloptions)
  



**Source File: [github.com/solo-io/gloo/projects/gloo/api/v1/options/protocol/protocol.proto](https://github.com/solo-io/gloo/blob/main/projects/gloo/api/v1/options/protocol/protocol.proto)**





---
### HttpProtocolOptions



```yaml
"idleTimeout": .google.protobuf.Duration
"maxHeadersCount": int
"maxStreamDuration": .google.protobuf.Duration
"headersWithUnderscoresAction": .protocol.options.gloo.solo.io.HttpProtocolOptions.HeadersWithUnderscoresAction

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `idleTimeout` | [.google.protobuf.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration) | The idle timeout for connections. The idle timeout is defined as the period in which there are no active requests. When the idle timeout is reached the connection will be closed. If the connection is an HTTP/2 downstream connection a drain sequence will occur prior to closing the connection, see :ref:`drain_timeout <envoy_api_field_extensions.filters.network.http_connection_manager.v3.HttpConnectionManager.drain_timeout>`. Note that request based timeouts mean that HTTP/2 PINGs will not keep the connection alive. If not specified, this defaults to 1 hour. To disable idle timeouts explicitly set this to 0. **Warning**: Disabling this timeout has a highly likelihood of yielding connection leaks due to lost TCP FIN packets, etc. |
| `maxHeadersCount` | `int` | The maximum number of headers. If unconfigured, the default maximum number of request headers allowed is 100. Requests that exceed this limit will receive a 431 response for HTTP/1.x and cause a stream reset for HTTP/2. |
| `maxStreamDuration` | [.google.protobuf.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration) | Total duration to keep alive an HTTP request/response stream. If the time limit is reached the stream will be reset independent of any other timeouts. If not specified, this value is not set. |
| `headersWithUnderscoresAction` | [.protocol.options.gloo.solo.io.HttpProtocolOptions.HeadersWithUnderscoresAction](../protocol.proto.sk/#headerswithunderscoresaction) | Action to take when a client request with a header name containing underscore characters is received. If this setting is not specified, the value defaults to ALLOW. Note: upstream responses are not affected by this setting. |




---
### HeadersWithUnderscoresAction

 
Action to take when Envoy receives client request with header names containing underscore
characters.
Underscore character is allowed in header names by the RFC-7230 and this behavior is implemented
as a security measure due to systems that treat '_' and '-' as interchangeable. Envoy by default allows client request headers with underscore
characters.

| Name | Description |
| ----- | ----------- | 
| `ALLOW` | Allow headers with underscores. This is the default behavior. |
| `REJECT_REQUEST` | Reject client request. HTTP/1 requests are rejected with the 400 status. HTTP/2 requests end with the stream reset. The "httpN.requests_rejected_with_underscores_in_headers" counter is incremented for each rejected request. |
| `DROP_HEADER` | Drop the header with name containing underscores. The header is dropped before the filter chain is invoked and as such filters will not see dropped headers. The "httpN.dropped_headers_with_underscores" is incremented for each dropped header. |




---
### Http1ProtocolOptions



```yaml
"enableTrailers": bool
"properCaseHeaderKeyFormat": bool
"preserveCaseHeaderKeyFormat": bool
"overrideStreamErrorOnInvalidHttpMessage": .google.protobuf.BoolValue

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `enableTrailers` | `bool` | Enables trailers for HTTP/1. By default the HTTP/1 codec drops proxied trailers. Note: Trailers must also be enabled at the gateway level in order for this option to take effect. |
| `properCaseHeaderKeyFormat` | `bool` | Formats the RESPONSE HEADER by proper casing words: the first character and any character following a special character will be capitalized if it's an alpha character. For example, "content-type" becomes "Content-Type", and "foo$b#$are" becomes "Foo$B#$Are". Note that while this results in most headers following conventional casing, certain headers are not covered. For example, the "TE" header will be formatted as "Te". Only one of `properCaseHeaderKeyFormat` or `preserveCaseHeaderKeyFormat` can be set. |
| `preserveCaseHeaderKeyFormat` | `bool` | Generates configuration for a stateful formatter extension that allows using received headers to affect the output of encoding headers. Specifically: preserving RESPONSE HEADER case during proxying. Only one of `preserveCaseHeaderKeyFormat` or `properCaseHeaderKeyFormat` can be set. |
| `overrideStreamErrorOnInvalidHttpMessage` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | Allows invalid HTTP messaging. When this option is false, then Envoy will terminate HTTP/1.1 connections upon receiving an invalid HTTP message. However, when this option is true, then Envoy will leave the HTTP/1.1 connection open where possible. If set, this overrides any HCM :ref:`stream_error_on_invalid_http_messaging <envoy_v3_api_field_extensions.filters.network.http_connection_manager.v3.HttpConnectionManager.stream_error_on_invalid_http_message>`. |




---
### Http2ProtocolOptions



```yaml
"maxConcurrentStreams": .google.protobuf.UInt32Value
"initialStreamWindowSize": .google.protobuf.UInt32Value
"initialConnectionWindowSize": .google.protobuf.UInt32Value
"overrideStreamErrorOnInvalidHttpMessage": .google.protobuf.BoolValue

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `maxConcurrentStreams` | [.google.protobuf.UInt32Value](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/u-int-32-value) | [Maximum concurrent streams](https://httpwg.org/specs/rfc7540.html#rfc.section.5.1.2) allowed for peer on one HTTP/2 connection. Valid values range from 1 to 2147483647 (2^31 - 1) and defaults to 2147483647. For upstream connections, this also limits how many streams Envoy will initiate concurrently on a single connection. If the limit is reached, Envoy may queue requests or establish additional connections (as allowed per circuit breaker limits). This acts as an upper bound: Envoy will lower the max concurrent streams allowed on a given connection based on upstream settings. Config dumps will reflect the configured upper bound, not the per-connection negotiated limits. |
| `initialStreamWindowSize` | [.google.protobuf.UInt32Value](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/u-int-32-value) | [Initial stream-level flow-control window](https://httpwg.org/specs/rfc7540.html#rfc.section.6.9.2) size. Valid values range from 65535 (2^16 - 1, HTTP/2 default) to 2147483647 (2^31 - 1, HTTP/2 maximum) and defaults to 268435456 (256 * 1024 * 1024). NOTE: 65535 is the initial window size from HTTP/2 spec. We only support increasing the default window size now, so it's also the minimum. This field also acts as a soft limit on the number of bytes Envoy will buffer per-stream in the HTTP/2 codec buffers. Once the buffer reaches this pointer, watermark callbacks will fire to stop the flow of data to the codec buffers. |
| `initialConnectionWindowSize` | [.google.protobuf.UInt32Value](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/u-int-32-value) | Similar to *initial_stream_window_size*, but for connection-level flow-control window. Currently, this has the same minimum/maximum/default as *initial_stream_window_size*. |
| `overrideStreamErrorOnInvalidHttpMessage` | [.google.protobuf.BoolValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/bool-value) | Allows invalid HTTP messaging and headers. When this option is disabled (default), then the whole HTTP/2 connection is terminated upon receiving invalid HEADERS frame. However, when this option is enabled, only the offending stream is terminated. This overrides any HCM :ref:`stream_error_on_invalid_http_messaging <envoy_v3_api_field_extensions.filters.network.http_connection_manager.v3.HttpConnectionManager.stream_error_on_invalid_http_message>` See [RFC7540, sec. 8.1](https://datatracker.ietf.org/doc/html/rfc7540#section-8.1) for details. |





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
