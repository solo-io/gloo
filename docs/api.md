# Protocol Documentation
<a name="top"/>

## Table of Contents

- [config.proto](#config.proto)
    - [Config](#v1.Config)
  
  
  
  

- [metadata.proto](#metadata.proto)
    - [Metadata](#v1.Metadata)
    - [Metadata.AnnotationsEntry](#v1.Metadata.AnnotationsEntry)
  
  
  
  

- [status.proto](#status.proto)
    - [Status](#v1.Status)
  
    - [Status.State](#v1.Status.State)
  
  
  

- [upstream.proto](#upstream.proto)
    - [Function](#v1.Function)
    - [Upstream](#v1.Upstream)
  
  
  
  

- [virtualhost.proto](#virtualhost.proto)
    - [Destination](#v1.Destination)
    - [EventMatcher](#v1.EventMatcher)
    - [FunctionDestination](#v1.FunctionDestination)
    - [RequestMatcher](#v1.RequestMatcher)
    - [RequestMatcher.HeadersEntry](#v1.RequestMatcher.HeadersEntry)
    - [RequestMatcher.QueryParamsEntry](#v1.RequestMatcher.QueryParamsEntry)
    - [Route](#v1.Route)
    - [SSLConfig](#v1.SSLConfig)
    - [UpstreamDestination](#v1.UpstreamDestination)
    - [VirtualHost](#v1.VirtualHost)
    - [WeightedDestination](#v1.WeightedDestination)
  
  
  
  

- [Scalar Value Types](#scalar-value-types)



<a name="config.proto"/>
<p align="right"><a href="#top">Top</a></p>

## config.proto



<a name="v1.Config"/>

### Config
Config is a top-level config object. It is used internally by gloo as a container for the entire user config.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| upstreams | [Upstream](#v1.Upstream) | repeated | The list of all upstreams defined by the user. |
| virtual_hosts | [VirtualHost](#v1.VirtualHost) | repeated | the list of all virtual hosts defined by the user. |





 

 

 

 



<a name="metadata.proto"/>
<p align="right"><a href="#top">Top</a></p>

## metadata.proto



<a name="v1.Metadata"/>

### Metadata



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_version | [string](#string) |  |  |
| namespace | [string](#string) |  |  |
| annotations | [Metadata.AnnotationsEntry](#v1.Metadata.AnnotationsEntry) | repeated | ignored by gloo but useful for clients |






<a name="v1.Metadata.AnnotationsEntry"/>

### Metadata.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 

 

 

 



<a name="status.proto"/>
<p align="right"><a href="#top">Top</a></p>

## status.proto



<a name="v1.Status"/>

### Status



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [Status.State](#v1.Status.State) |  |  |
| reason | [string](#string) |  |  |





 


<a name="v1.Status.State"/>

### Status.State


| Name | Number | Description |
| ---- | ------ | ----------- |
| Pending | 0 |  |
| Accepted | 1 |  |
| Rejected | 2 |  |


 

 

 



<a name="upstream.proto"/>
<p align="right"><a href="#top">Top</a></p>

## upstream.proto



<a name="v1.Function"/>

### Function



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| spec | [google.protobuf.Struct](#google.protobuf.Struct) |  |  |






<a name="v1.Upstream"/>

### Upstream



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| type | [string](#string) |  |  |
| connection_timeout | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| spec | [google.protobuf.Struct](#google.protobuf.Struct) |  |  |
| functions | [Function](#v1.Function) | repeated |  |
| status | [Status](#v1.Status) |  | read only |
| metadata | [Metadata](#v1.Metadata) |  |  |





 

 

 

 



<a name="virtualhost.proto"/>
<p align="right"><a href="#top">Top</a></p>

## virtualhost.proto



<a name="v1.Destination"/>

### Destination



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| function | [FunctionDestination](#v1.FunctionDestination) |  |  |
| upstream | [UpstreamDestination](#v1.UpstreamDestination) |  |  |






<a name="v1.EventMatcher"/>

### EventMatcher



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event_type | [string](#string) |  |  |






<a name="v1.FunctionDestination"/>

### FunctionDestination



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| upstream_name | [string](#string) |  |  |
| function_name | [string](#string) |  |  |






<a name="v1.RequestMatcher"/>

### RequestMatcher



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path_prefix | [string](#string) |  |  |
| path_regex | [string](#string) |  |  |
| path_exact | [string](#string) |  |  |
| headers | [RequestMatcher.HeadersEntry](#v1.RequestMatcher.HeadersEntry) | repeated |  |
| query_params | [RequestMatcher.QueryParamsEntry](#v1.RequestMatcher.QueryParamsEntry) | repeated |  |
| verbs | [string](#string) | repeated |  |






<a name="v1.RequestMatcher.HeadersEntry"/>

### RequestMatcher.HeadersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="v1.RequestMatcher.QueryParamsEntry"/>

### RequestMatcher.QueryParamsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="v1.Route"/>

### Route



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_matcher | [RequestMatcher](#v1.RequestMatcher) |  |  |
| event_matcher | [EventMatcher](#v1.EventMatcher) |  |  |
| multiple_destinations | [WeightedDestination](#v1.WeightedDestination) | repeated |  |
| single_destination | [Destination](#v1.Destination) |  |  |
| prefix_rewrite | [string](#string) |  |  |
| extensions | [google.protobuf.Struct](#google.protobuf.Struct) |  |  |






<a name="v1.SSLConfig"/>

### SSLConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| secret_ref | [string](#string) |  |  |






<a name="v1.UpstreamDestination"/>

### UpstreamDestination



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="v1.VirtualHost"/>

### VirtualHost



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | required, must be unique cannot be &#34;default&#34; unless it refers to the default vhost |
| domains | [string](#string) | repeated | if this is empty, this host will become / be merged with the default virtualhost who has domains = [&#34;*&#34;] |
| routes | [Route](#v1.Route) | repeated | require at least 1 route |
| ssl_config | [SSLConfig](#v1.SSLConfig) |  | optional |
| status | [Status](#v1.Status) |  | read only |
| metadata | [Metadata](#v1.Metadata) |  |  |






<a name="v1.WeightedDestination"/>

### WeightedDestination



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| destination | [Destination](#v1.Destination) |  |  |
| weight | [uint32](#uint32) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

