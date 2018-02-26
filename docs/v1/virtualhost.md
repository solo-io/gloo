<a name="top"></a>

## Contents
  - [VirtualHost](#v1.VirtualHost)
  - [Route](#v1.Route)
  - [RequestMatcher](#v1.RequestMatcher)
  - [EventMatcher](#v1.EventMatcher)
  - [WeightedDestination](#v1.WeightedDestination)
  - [Destination](#v1.Destination)
  - [FunctionDestination](#v1.FunctionDestination)
  - [UpstreamDestination](#v1.UpstreamDestination)
  - [SSLConfig](#v1.SSLConfig)



<a name="virtualhost"></a>
<p align="right"><a href="#top">Top</a></p>




<a name="v1.VirtualHost"></a>

### VirtualHost
Virtual Hosts represent a collection of routes for a set of domains.
Virtual Hosts can be compared to [virtual hosts](TODO) in [envoy](TODO) terminology.
A virtual host can be used to define &#34;apps&#34;; a collection of APIs that belong to a particular domain.
The Virtual Host concept allows configuration of per-virtualhost SSL certificates


```yaml
name: string
domains: [string]
routes: [{Route}]
ssl_config: {SSLConfig}
status: (read only)
metadata: (read only)

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | string |  | Name of the virtual host. Names must be unique and follow the following syntax rules: One or more lowercase rfc1035/rfc1123 labels separated by &#39;.&#39; with a maximum length of 253 characters. |
| domains | string | repeated | Domains represent the list of domains (host/authority header) that will match for all routes on this virtual host. As in [envoy](TODO): wildcard hosts are supported in the form of “*.foo.com” or “*-bar.foo.com”. If domains is empty, gloo will set the domain to &#34;*&#34;, making that virtual host the &#34;default&#34; virtualhost. The default virtualhost will be the fallback virtual host for all requests that do not match a domain on an existing virtual host. Only one default virtual host can be defined (either with an empty domain list, or a domain list that includes &#34;*&#34;) |
| routes | [Route](virtualhost.md#v1.Route) | repeated | Routes define the list of [routes](TODO) that live on this virtual host. |
| ssl_config | [SSLConfig](virtualhost.md#v1.SSLConfig) |  | SSL Config is optional for the virtual host. If provided, the virtual host will listen on the envoy HTTPS listener port (default :8443) If left empty, the virtual host will listen on the HTTP listener port (default :8080) |
| status | [Status](status.md#v1.Status) |  | Status indicates the validation status of the virtual host resource. Status is read-only by clients, and set by gloo during validation |
| metadata | [Metadata](metadata.md#v1.Metadata) |  | Metadata contains the resource metadata for the virtual host |






<a name="v1.Route"></a>

### Route
Routes declare the entrypoints on virtual hosts and the upstreams or functions they route requests to


```yaml
request_matcher: {RequestMatcher}
event_matcher: {EventMatcher}
multiple_destinations: [{WeightedDestination}]
single_destination: {Destination}
prefix_rewrite: string
extensions: {google.protobuf.Struct}

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_matcher | [RequestMatcher](virtualhost.md#v1.RequestMatcher) |  | request_matcher indicates this route should match requests according to the specification in the provided RequestMatcher only one of request_matcher or event_matcher can be set |
| event_matcher | [EventMatcher](virtualhost.md#v1.EventMatcher) |  | eventt_matcher indicates this route should match requests according to the specification in the provided EventMatcher only one of request_matcher or event_matcher can be set |
| multiple_destinations | [WeightedDestination](virtualhost.md#v1.WeightedDestination) | repeated | A route is only allowed to specify one of multiple_destinations or single_destination. Setting both will result in an error Multiple Destinations is used when a user wants a route to balance requests between multiple destinations Balancing is done by probability, where weights are specified for each destination |
| single_destination | [Destination](virtualhost.md#v1.Destination) |  | A single destination is specified when a route only routes to a single destination. |
| prefix_rewrite | string |  | PrefixRewrite can be specified to rewrite the matched path of the request path to a new prefix |
| extensions | [google.protobuf.Struct](virtualhost.md#google.protobuf.Struct) |  | Extensions provides a way to extend the behavior of a route. In addition to the [core route extensions](TODO), gloo provides the means for [route plugins](TODO) to be added to gloo which add new types of route extensions. See the [route extensions section](TODO) for a more detailed explanation |






<a name="v1.RequestMatcher"></a>

### RequestMatcher
Request Matcher is a route matcher for traditional http requests
Request Matchers stand in juxtoposition to Event Matchers, which match &#34;events&#34; rather than HTTP Requests


```yaml
path_prefix: string
path_regex: string
path_exact: string
headers: map<string,string>
query_params: map<string,string>
verbs: [string]

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path_prefix | string |  | Prefix will match any request whose path begins with this prefix Only one of path_prefix, path_regex, or path_exact can be set |
| path_regex | string |  | Regex will match any path that matches this regex string Only one of path_prefix, path_regex, or path_exact can be set |
| path_exact | string |  | Exact will match only requests with exactly this path Only one of path_prefix, path_regex, or path_exact can be set |
| headers | map&lt;string,string&gt; |  | Headers specify a list of request headers and their values the request must contain to match this route If a value is not specified (empty string) for a header, all values will match so long as the header is present on the request |
| query_params | map&lt;string,string&gt; |  | Query params work the same way as headers, but for query string parameters |
| verbs | string | repeated | HTTP Verb(s) to match on. If none specified, the matcher will match all verbs |






<a name="v1.EventMatcher"></a>

### EventMatcher
Event matcher is a special kind of matcher for CloudEvents
The CloudEvents API is described here: https://github.com/cloudevents/spec/blob/master/spec.md


```yaml
event_type: string

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event_type | string |  | Event Type indicates the event type or topic to match |






<a name="v1.WeightedDestination"></a>

### WeightedDestination
WeightedDestination attaches a weight to a destination
For use in routes with multiple destinations


```yaml
destination: {Destination}
weight: uint32

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| destination | [Destination](virtualhost.md#v1.Destination) |  |  |
| weight | uint32 |  | Weight must be greater than zero Routing to each destination will be balanced by the ratio of the destination&#39;s weight to the total weight on a route |






<a name="v1.Destination"></a>

### Destination
Destination is a destination that requests can be routed to.


```yaml
function: {FunctionDestination}
upstream: {UpstreamDestination}

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| function | [FunctionDestination](virtualhost.md#v1.FunctionDestination) |  | function indicates requests sent to this destination will invoke a function Only one of funtion or upstream should be set |
| upstream | [UpstreamDestination](virtualhost.md#v1.UpstreamDestination) |  | upstream indicates requests sent to this destination will be routed to an upstream Only one of funtion or upstream should be set |






<a name="v1.FunctionDestination"></a>

### FunctionDestination
FunctionDestination will route a request to a specific function defined for an upstream


```yaml
upstream_name: string
function_name: string

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| upstream_name | string |  | Upstream Name is the name of the upstream the function belongs to |
| function_name | string |  | Function Name is the name of the function as defined on the upstream |






<a name="v1.UpstreamDestination"></a>

### UpstreamDestination
Upstream Destination routes a request to an upstream


```yaml
name: string

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | string |  | Name of the upstream |






<a name="v1.SSLConfig"></a>

### SSLConfig
SSLConfig contains the options necessary to configure a virtualhost to use TLS


```yaml
secret_ref: string

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| secret_ref | string |  | SecretRef contains the [secret ref](TODO) to a [gloo secret](TODO) containing the following structure: { &#34;ca_chain&#34;: &lt;ca chain data...&gt;, &#34;private key&#34;: &lt;private key data...&gt; } |





 

 

 

