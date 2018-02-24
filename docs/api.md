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
Metadata contains general properties of config resources useful to clients and the gloo control plane for purposes of versioning, annotating, and namespacing resources.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_version | [string](#string) |  | ResourceVersion keeps track of the resource version of a config resource. This mechanism is used by [gloo-storage](TODO) to ensure safety with concurrent writes/updates to a resource in storage. |
| namespace | [string](#string) |  | Namespace is used for the namespacing of resources. Currently unused by gloo internally. |
| annotations | [Metadata.AnnotationsEntry](#v1.Metadata.AnnotationsEntry) | repeated | Annotations allow clients to tag resources for special use cases. gloo ignores annotations but preserved them on read/write from/to storage. |






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
Status indicates whether a config resource (currently only [virtualhosts](TODO) and [upstreams](TODO)) has been (in)validated by gloo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [Status.State](#v1.Status.State) |  | State is the enum indicating the state of the resource |
| reason | [string](#string) |  | Reason is a description of the error for Rejected resources. If the resource is pending or accepted, this field will be empty |





 


<a name="v1.Status.State"/>

### Status.State


| Name | Number | Description |
| ---- | ------ | ----------- |
| Pending | 0 | Pending status indicates the resource has not yet been validated |
| Accepted | 1 | Accepted indicates the resource has been validated |
| Rejected | 2 | Rejected indicates an invalid configuration by the user Rejected resources may be propagated to the xDS server depending on their severity |


 

 

 



<a name="upstream.proto"/>
<p align="right"><a href="#top">Top</a></p>

## upstream.proto



<a name="v1.Function"/>

### Function



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the function. Functions are referenced by name from routes and therefore must be unique within an upstream |
| spec | [google.protobuf.Struct](#google.protobuf.Struct) |  | Spec for the function. Like [upstream specs](TODO), the content of function specs is specified by the [upstream plugin](TODO) for the upstream&#39;s type. |






<a name="v1.Upstream"/>

### Upstream
Upstream represents a destination for routing. Upstreams can be compared to [clusters](TODO) in [envoy](TODO) terminology.
Upstreams can take a variety of [types](TODO) in gloo. Language extensions known as [plugins](TODO) allow the addition of new
types of upstreams. See [upstream types](TODO) for a detailed description of available upstream types.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the upstream. Names must be unique and follow the following syntax rules: One or more lowercase rfc1035/rfc1123 labels separated by &#39;.&#39; with a maximum length of 253 characters. |
| type | [string](#string) |  | Type indicates the type of the upstream. Examples include [service](TODO), [kubernetes](TODO), and [aws](TODO) Types are defined by the [plugin](TODO) associated with them. |
| connection_timeout | [google.protobuf.Duration](#google.protobuf.Duration) |  | Connection Timeout tells gloo to set a timeout for unresponsive connections created to this upstream. If not provided by the user, it will default to a [default value](TODO) |
| spec | [google.protobuf.Struct](#google.protobuf.Struct) |  | Spec contains properties that are specific to the upstream type. The spec is always required, but the expected content is specified by the [upstream plugin] for the given upstream type. Most often the upstream spec will be a map&lt;string, string&gt; |
| functions | [Function](#v1.Function) | repeated | Certain upstream types support (and may require) [functions](TODO). Functions allow function-level routing to be done. For example, the [aws lambda](TODO) upstream type Permits routing to [aws lambda functions]. [routes](TODO) on [virtualhosts] can specify [function destinations] to route to specific functions. |
| status | [Status](#v1.Status) |  | Status indicates the validation status of the upstream resource. Status is read-only by clients, and set by gloo during validation |
| metadata | [Metadata](#v1.Metadata) |  | Metadata contains the resource metadata for the upstream |





 

 

 

 



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
Virtual Hosts represent a collection of routes for a set of domains.
Virtual Hosts can be compared to [virtual hosts](TODO) in [envoy](TODO) terminology.
A virtual host can be used to define &#34;apps&#34;; a collection of APIs that belong to a particular domain.
The Virtual Host concept allows configuration of per-virtualhost SSL certificates


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
Virtual Hosts represent a collection of routes for a set of domains.
Virtual Hosts can be compared to [virtual hosts](TODO) in [envoy](TODO) terminology.
A virtual host can be used to define &#34;apps&#34;; a collection of APIs that belong to a particular domain.
The Virtual Host concept allows configuration of per-virtualhost SSL certificates


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the virtual host. Names must be unique and follow the following syntax rules: One or more lowercase rfc1035/rfc1123 labels separated by &#39;.&#39; with a maximum length of 253 characters. |
| domains | [string](#string) | repeated | Domains represent the list of domains (host/authority header) that will match for all routes on this virtual host. As in [envoy](TODO): wildcard hosts are supported in the form of “*.foo.com” or “*-bar.foo.com”. If domains is empty, gloo will set the domain to &#34;*&#34;, making that virtual host the &#34;default&#34; virtualhost. The default virtualhost will be the fallback virtual host for all requests that do not match a domain on an existing virtual host. Only one default virtual host can be defined (either with an empty domain list, or a domain list that includes &#34;*&#34;) |
| routes | [Route](#v1.Route) | repeated | Routes define the list of [routes](TODO) that live on this virtual host. |
| ssl_config | [SSLConfig](#v1.SSLConfig) |  | SSL Config is optional for the virtual host. If provided, the virtual host will listen on the envoy HTTPS listener port (default :8443) If left empty, the virtual host will listen on the HTTP listener port (default :8080) |
| status | [Status](#v1.Status) |  | Status indicates the validation status of the virtual host resource. Status is read-only by clients, and set by gloo during validation |
| metadata | [Metadata](#v1.Metadata) |  | Metadata contains the resource metadata for the virtual host |






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

