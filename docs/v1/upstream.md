<a name="top"/>

## Contents
  - [Upstream](#v1.Upstream)
  - [Function](#v1.Function)



<a name="upstream"/>
<p align="right"><a href="#top">Top</a></p>




<a name="v1.Upstream"/>

### Upstream
[]()Upstream represents a destination for routing. Upstreams can be compared to [clusters](TODO) in [envoy](TODO) terminology.
Upstreams can take a variety of [types](TODO) in gloo. Language extensions known as [plugins](TODO) allow the addition of new
types of upstreams. See [upstream types](TODO) for a detailed description of available upstream types.


```yaml
name: string
type: string
connection_timeout: {google.protobuf.Duration}
spec: {google.protobuf.Struct}
functions: [{Function}]
status: (read only)
metadata: (read only)

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the upstream. Names must be unique and follow the following syntax rules: One or more lowercase rfc1035/rfc1123 labels separated by &#39;.&#39; with a maximum length of 253 characters. |
| type | [string](#string) |  | Type indicates the type of the upstream. Examples include [service](TODO), [kubernetes](TODO), and [aws](TODO) Types are defined by the [plugin](TODO) associated with them. |
| connection_timeout | [google.protobuf.Duration](#google.protobuf.Duration) |  | Connection Timeout tells gloo to set a timeout for unresponsive connections created to this upstream. If not provided by the user, it will default to a [default value](TODO) |
| spec | [google.protobuf.Struct](#google.protobuf.Struct) |  | Spec contains properties that are specific to the upstream type. The spec is always required, but the expected content is specified by the [upstream plugin] for the given upstream type. Most often the upstream spec will be a map&lt;string, string&gt; |
| functions | [Function](#v1.Function) | repeated | Certain upstream types support (and may require) [functions](TODO). Functions allow function-level routing to be done. For example, the [aws lambda](TODO) upstream type Permits routing to [aws lambda functions]. [routes](TODO) on [virtualhosts] can specify [function destinations] to route to specific functions. |
| status | [Status](#v1.Status) |  | Status indicates the validation status of the upstream resource. Status is read-only by clients, and set by gloo during validation |
| metadata | [Metadata](#v1.Metadata) |  | Metadata contains the resource metadata for the upstream |






<a name="v1.Function"/>

### Function
[]()


```yaml
name: string
spec: {google.protobuf.Struct}

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the function. Functions are referenced by name from routes and therefore must be unique within an upstream |
| spec | [google.protobuf.Struct](#google.protobuf.Struct) |  | Spec for the function. Like [upstream specs](TODO), the content of function specs is specified by the [upstream plugin](TODO) for the upstream&#39;s type. |





 

 

 

