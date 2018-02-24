# v1
--
    import "github.com/solo-io/gloo-api/pkg/api/types/v1"

Package v1 is a generated protocol buffer package.

It is generated from these files:

    config.proto
    metadata.proto
    status.proto
    upstream.proto
    virtualhost.proto

It has these top-level messages:

    Config
    Metadata
    Status
    Upstream
    Function
    VirtualHost
    Route
    RequestMatcher
    EventMatcher
    WeightedDestination
    Destination
    FunctionDestination
    UpstreamDestination
    SSLConfig

## Usage

```go
var Status_State_name = map[int32]string{
	0: "Pending",
	1: "Accepted",
	2: "Rejected",
}
```

```go
var Status_State_value = map[string]int32{
	"Pending":  0,
	"Accepted": 1,
	"Rejected": 2,
}
```

#### type Config

```go
type Config struct {
	Upstreams    []*Upstream    `protobuf:"bytes,1,rep,name=upstreams" json:"upstreams,omitempty"`
	VirtualHosts []*VirtualHost `protobuf:"bytes,2,rep,name=virtual_hosts,json=virtualHosts" json:"virtual_hosts,omitempty"`
}
```

* Config is a top-level config object. It is used internally by gloo as a
container for the entire user config.

#### func (*Config) Descriptor

```go
func (*Config) Descriptor() ([]byte, []int)
```

#### func (*Config) Equal

```go
func (this *Config) Equal(that interface{}) bool
```

#### func (*Config) GetUpstreams

```go
func (m *Config) GetUpstreams() []*Upstream
```

#### func (*Config) GetVirtualHosts

```go
func (m *Config) GetVirtualHosts() []*VirtualHost
```

#### func (*Config) ProtoMessage

```go
func (*Config) ProtoMessage()
```

#### func (*Config) Reset

```go
func (m *Config) Reset()
```

#### func (*Config) String

```go
func (m *Config) String() string
```

#### type ConfigObject

```go
type ConfigObject interface {
	proto.Message
	GetName() string
	GetMetadata() *Metadata
}
```


#### type Destination

```go
type Destination struct {
	// Types that are valid to be assigned to DestinationType:
	//	*Destination_Function
	//	*Destination_Upstream
	DestinationType isDestination_DestinationType `protobuf_oneof:"destination_type"`
}
```


#### func (*Destination) Descriptor

```go
func (*Destination) Descriptor() ([]byte, []int)
```

#### func (*Destination) Equal

```go
func (this *Destination) Equal(that interface{}) bool
```

#### func (*Destination) GetDestinationType

```go
func (m *Destination) GetDestinationType() isDestination_DestinationType
```

#### func (*Destination) GetFunction

```go
func (m *Destination) GetFunction() *FunctionDestination
```

#### func (*Destination) GetUpstream

```go
func (m *Destination) GetUpstream() *UpstreamDestination
```

#### func (*Destination) ProtoMessage

```go
func (*Destination) ProtoMessage()
```

#### func (*Destination) Reset

```go
func (m *Destination) Reset()
```

#### func (*Destination) String

```go
func (m *Destination) String() string
```

#### func (*Destination) XXX_OneofFuncs

```go
func (*Destination) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{})
```
XXX_OneofFuncs is for the internal use of the proto package.

#### type Destination_Function

```go
type Destination_Function struct {
	Function *FunctionDestination `protobuf:"bytes,1,opt,name=function,oneof"`
}
```


#### func (*Destination_Function) Equal

```go
func (this *Destination_Function) Equal(that interface{}) bool
```

#### type Destination_Upstream

```go
type Destination_Upstream struct {
	Upstream *UpstreamDestination `protobuf:"bytes,2,opt,name=upstream,oneof"`
}
```


#### func (*Destination_Upstream) Equal

```go
func (this *Destination_Upstream) Equal(that interface{}) bool
```

#### type EventMatcher

```go
type EventMatcher struct {
	EventType string `protobuf:"bytes,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
}
```


#### func (*EventMatcher) Descriptor

```go
func (*EventMatcher) Descriptor() ([]byte, []int)
```

#### func (*EventMatcher) Equal

```go
func (this *EventMatcher) Equal(that interface{}) bool
```

#### func (*EventMatcher) GetEventType

```go
func (m *EventMatcher) GetEventType() string
```

#### func (*EventMatcher) ProtoMessage

```go
func (*EventMatcher) ProtoMessage()
```

#### func (*EventMatcher) Reset

```go
func (m *EventMatcher) Reset()
```

#### func (*EventMatcher) String

```go
func (m *EventMatcher) String() string
```

#### type Function

```go
type Function struct {
	Name string                  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Spec *google_protobuf.Struct `protobuf:"bytes,4,opt,name=spec" json:"spec,omitempty"`
}
```


#### func (*Function) Descriptor

```go
func (*Function) Descriptor() ([]byte, []int)
```

#### func (*Function) Equal

```go
func (this *Function) Equal(that interface{}) bool
```

#### func (*Function) GetName

```go
func (m *Function) GetName() string
```

#### func (*Function) GetSpec

```go
func (m *Function) GetSpec() *google_protobuf.Struct
```

#### func (*Function) ProtoMessage

```go
func (*Function) ProtoMessage()
```

#### func (*Function) Reset

```go
func (m *Function) Reset()
```

#### func (*Function) String

```go
func (m *Function) String() string
```

#### type FunctionDestination

```go
type FunctionDestination struct {
	UpstreamName string `protobuf:"bytes,1,opt,name=upstream_name,json=upstreamName,proto3" json:"upstream_name,omitempty"`
	FunctionName string `protobuf:"bytes,2,opt,name=function_name,json=functionName,proto3" json:"function_name,omitempty"`
}
```


#### func (*FunctionDestination) Descriptor

```go
func (*FunctionDestination) Descriptor() ([]byte, []int)
```

#### func (*FunctionDestination) Equal

```go
func (this *FunctionDestination) Equal(that interface{}) bool
```

#### func (*FunctionDestination) GetFunctionName

```go
func (m *FunctionDestination) GetFunctionName() string
```

#### func (*FunctionDestination) GetUpstreamName

```go
func (m *FunctionDestination) GetUpstreamName() string
```

#### func (*FunctionDestination) ProtoMessage

```go
func (*FunctionDestination) ProtoMessage()
```

#### func (*FunctionDestination) Reset

```go
func (m *FunctionDestination) Reset()
```

#### func (*FunctionDestination) String

```go
func (m *FunctionDestination) String() string
```

#### type FunctionSpec

```go
type FunctionSpec *types.Struct
```


#### type Metadata

```go
type Metadata struct {
	ResourceVersion string `protobuf:"bytes,1,opt,name=resource_version,json=resourceVersion,proto3" json:"resource_version,omitempty"`
	Namespace       string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// ignored by gloo but useful for clients
	Annotations map[string]string `protobuf:"bytes,3,rep,name=annotations" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}
```


#### func (*Metadata) Descriptor

```go
func (*Metadata) Descriptor() ([]byte, []int)
```

#### func (*Metadata) Equal

```go
func (this *Metadata) Equal(that interface{}) bool
```

#### func (*Metadata) GetAnnotations

```go
func (m *Metadata) GetAnnotations() map[string]string
```

#### func (*Metadata) GetNamespace

```go
func (m *Metadata) GetNamespace() string
```

#### func (*Metadata) GetResourceVersion

```go
func (m *Metadata) GetResourceVersion() string
```

#### func (*Metadata) ProtoMessage

```go
func (*Metadata) ProtoMessage()
```

#### func (*Metadata) Reset

```go
func (m *Metadata) Reset()
```

#### func (*Metadata) String

```go
func (m *Metadata) String() string
```

#### type RequestMatcher

```go
type RequestMatcher struct {
	// Types that are valid to be assigned to Path:
	//	*RequestMatcher_PathPrefix
	//	*RequestMatcher_PathRegex
	//	*RequestMatcher_PathExact
	Path        isRequestMatcher_Path `protobuf_oneof:"path"`
	Headers     map[string]string     `protobuf:"bytes,4,rep,name=headers" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	QueryParams map[string]string     `protobuf:"bytes,5,rep,name=query_params,json=queryParams" json:"query_params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Verbs       []string              `protobuf:"bytes,6,rep,name=verbs" json:"verbs,omitempty"`
}
```


#### func (*RequestMatcher) Descriptor

```go
func (*RequestMatcher) Descriptor() ([]byte, []int)
```

#### func (*RequestMatcher) Equal

```go
func (this *RequestMatcher) Equal(that interface{}) bool
```

#### func (*RequestMatcher) GetHeaders

```go
func (m *RequestMatcher) GetHeaders() map[string]string
```

#### func (*RequestMatcher) GetPath

```go
func (m *RequestMatcher) GetPath() isRequestMatcher_Path
```

#### func (*RequestMatcher) GetPathExact

```go
func (m *RequestMatcher) GetPathExact() string
```

#### func (*RequestMatcher) GetPathPrefix

```go
func (m *RequestMatcher) GetPathPrefix() string
```

#### func (*RequestMatcher) GetPathRegex

```go
func (m *RequestMatcher) GetPathRegex() string
```

#### func (*RequestMatcher) GetQueryParams

```go
func (m *RequestMatcher) GetQueryParams() map[string]string
```

#### func (*RequestMatcher) GetVerbs

```go
func (m *RequestMatcher) GetVerbs() []string
```

#### func (*RequestMatcher) ProtoMessage

```go
func (*RequestMatcher) ProtoMessage()
```

#### func (*RequestMatcher) Reset

```go
func (m *RequestMatcher) Reset()
```

#### func (*RequestMatcher) String

```go
func (m *RequestMatcher) String() string
```

#### func (*RequestMatcher) XXX_OneofFuncs

```go
func (*RequestMatcher) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{})
```
XXX_OneofFuncs is for the internal use of the proto package.

#### type RequestMatcher_PathExact

```go
type RequestMatcher_PathExact struct {
	PathExact string `protobuf:"bytes,3,opt,name=path_exact,json=pathExact,proto3,oneof"`
}
```


#### func (*RequestMatcher_PathExact) Equal

```go
func (this *RequestMatcher_PathExact) Equal(that interface{}) bool
```

#### type RequestMatcher_PathPrefix

```go
type RequestMatcher_PathPrefix struct {
	PathPrefix string `protobuf:"bytes,1,opt,name=path_prefix,json=pathPrefix,proto3,oneof"`
}
```


#### func (*RequestMatcher_PathPrefix) Equal

```go
func (this *RequestMatcher_PathPrefix) Equal(that interface{}) bool
```

#### type RequestMatcher_PathRegex

```go
type RequestMatcher_PathRegex struct {
	PathRegex string `protobuf:"bytes,2,opt,name=path_regex,json=pathRegex,proto3,oneof"`
}
```


#### func (*RequestMatcher_PathRegex) Equal

```go
func (this *RequestMatcher_PathRegex) Equal(that interface{}) bool
```

#### type Route

```go
type Route struct {
	// Types that are valid to be assigned to Matcher:
	//	*Route_RequestMatcher
	//	*Route_EventMatcher
	Matcher              isRoute_Matcher         `protobuf_oneof:"matcher"`
	MultipleDestinations []*WeightedDestination  `protobuf:"bytes,3,rep,name=multiple_destinations,json=multipleDestinations" json:"multiple_destinations,omitempty"`
	SingleDestination    *Destination            `protobuf:"bytes,4,opt,name=single_destination,json=singleDestination" json:"single_destination,omitempty"`
	PrefixRewrite        string                  `protobuf:"bytes,5,opt,name=prefix_rewrite,json=prefixRewrite,proto3" json:"prefix_rewrite,omitempty"`
	Extensions           *google_protobuf.Struct `protobuf:"bytes,6,opt,name=extensions" json:"extensions,omitempty"`
}
```


#### func (*Route) Descriptor

```go
func (*Route) Descriptor() ([]byte, []int)
```

#### func (*Route) Equal

```go
func (this *Route) Equal(that interface{}) bool
```

#### func (*Route) GetEventMatcher

```go
func (m *Route) GetEventMatcher() *EventMatcher
```

#### func (*Route) GetExtensions

```go
func (m *Route) GetExtensions() *google_protobuf.Struct
```

#### func (*Route) GetMatcher

```go
func (m *Route) GetMatcher() isRoute_Matcher
```

#### func (*Route) GetMultipleDestinations

```go
func (m *Route) GetMultipleDestinations() []*WeightedDestination
```

#### func (*Route) GetPrefixRewrite

```go
func (m *Route) GetPrefixRewrite() string
```

#### func (*Route) GetRequestMatcher

```go
func (m *Route) GetRequestMatcher() *RequestMatcher
```

#### func (*Route) GetSingleDestination

```go
func (m *Route) GetSingleDestination() *Destination
```

#### func (*Route) ProtoMessage

```go
func (*Route) ProtoMessage()
```

#### func (*Route) Reset

```go
func (m *Route) Reset()
```

#### func (*Route) String

```go
func (m *Route) String() string
```

#### func (*Route) XXX_OneofFuncs

```go
func (*Route) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{})
```
XXX_OneofFuncs is for the internal use of the proto package.

#### type Route_EventMatcher

```go
type Route_EventMatcher struct {
	EventMatcher *EventMatcher `protobuf:"bytes,2,opt,name=event_matcher,json=eventMatcher,oneof"`
}
```


#### func (*Route_EventMatcher) Equal

```go
func (this *Route_EventMatcher) Equal(that interface{}) bool
```

#### type Route_RequestMatcher

```go
type Route_RequestMatcher struct {
	RequestMatcher *RequestMatcher `protobuf:"bytes,1,opt,name=request_matcher,json=requestMatcher,oneof"`
}
```


#### func (*Route_RequestMatcher) Equal

```go
func (this *Route_RequestMatcher) Equal(that interface{}) bool
```

#### type SSLConfig

```go
type SSLConfig struct {
	SecretRef string `protobuf:"bytes,1,opt,name=secret_ref,json=secretRef,proto3" json:"secret_ref,omitempty"`
}
```


#### func (*SSLConfig) Descriptor

```go
func (*SSLConfig) Descriptor() ([]byte, []int)
```

#### func (*SSLConfig) Equal

```go
func (this *SSLConfig) Equal(that interface{}) bool
```

#### func (*SSLConfig) GetSecretRef

```go
func (m *SSLConfig) GetSecretRef() string
```

#### func (*SSLConfig) ProtoMessage

```go
func (*SSLConfig) ProtoMessage()
```

#### func (*SSLConfig) Reset

```go
func (m *SSLConfig) Reset()
```

#### func (*SSLConfig) String

```go
func (m *SSLConfig) String() string
```

#### type Status

```go
type Status struct {
	State  Status_State `protobuf:"varint,1,opt,name=state,proto3,enum=v1.Status_State" json:"state,omitempty"`
	Reason string       `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}
```


#### func (*Status) Descriptor

```go
func (*Status) Descriptor() ([]byte, []int)
```

#### func (*Status) Equal

```go
func (this *Status) Equal(that interface{}) bool
```

#### func (*Status) GetReason

```go
func (m *Status) GetReason() string
```

#### func (*Status) GetState

```go
func (m *Status) GetState() Status_State
```

#### func (*Status) ProtoMessage

```go
func (*Status) ProtoMessage()
```

#### func (*Status) Reset

```go
func (m *Status) Reset()
```

#### func (*Status) String

```go
func (m *Status) String() string
```

#### type Status_State

```go
type Status_State int32
```


```go
const (
	Status_Pending  Status_State = 0
	Status_Accepted Status_State = 1
	Status_Rejected Status_State = 2
)
```

#### func (Status_State) EnumDescriptor

```go
func (Status_State) EnumDescriptor() ([]byte, []int)
```

#### func (Status_State) String

```go
func (x Status_State) String() string
```

#### type Upstream

```go
type Upstream struct {
	Name              string                  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type              string                  `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	ConnectionTimeout time.Duration           `protobuf:"bytes,3,opt,name=connection_timeout,json=connectionTimeout,stdduration" json:"connection_timeout"`
	Spec              *google_protobuf.Struct `protobuf:"bytes,4,opt,name=spec" json:"spec,omitempty"`
	Functions         []*Function             `protobuf:"bytes,5,rep,name=functions" json:"functions,omitempty"`
	// read only
	Status   *Status   `protobuf:"bytes,6,opt,name=status" json:"status,omitempty"`
	Metadata *Metadata `protobuf:"bytes,7,opt,name=metadata" json:"metadata,omitempty"`
}
```


#### func (*Upstream) Descriptor

```go
func (*Upstream) Descriptor() ([]byte, []int)
```

#### func (*Upstream) Equal

```go
func (this *Upstream) Equal(that interface{}) bool
```

#### func (*Upstream) GetConnectionTimeout

```go
func (m *Upstream) GetConnectionTimeout() time.Duration
```

#### func (*Upstream) GetFunctions

```go
func (m *Upstream) GetFunctions() []*Function
```

#### func (*Upstream) GetMetadata

```go
func (m *Upstream) GetMetadata() *Metadata
```

#### func (*Upstream) GetName

```go
func (m *Upstream) GetName() string
```

#### func (*Upstream) GetSpec

```go
func (m *Upstream) GetSpec() *google_protobuf.Struct
```

#### func (*Upstream) GetStatus

```go
func (m *Upstream) GetStatus() *Status
```

#### func (*Upstream) GetType

```go
func (m *Upstream) GetType() string
```

#### func (*Upstream) ProtoMessage

```go
func (*Upstream) ProtoMessage()
```

#### func (*Upstream) Reset

```go
func (m *Upstream) Reset()
```

#### func (*Upstream) String

```go
func (m *Upstream) String() string
```

#### type UpstreamDestination

```go
type UpstreamDestination struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}
```


#### func (*UpstreamDestination) Descriptor

```go
func (*UpstreamDestination) Descriptor() ([]byte, []int)
```

#### func (*UpstreamDestination) Equal

```go
func (this *UpstreamDestination) Equal(that interface{}) bool
```

#### func (*UpstreamDestination) GetName

```go
func (m *UpstreamDestination) GetName() string
```

#### func (*UpstreamDestination) ProtoMessage

```go
func (*UpstreamDestination) ProtoMessage()
```

#### func (*UpstreamDestination) Reset

```go
func (m *UpstreamDestination) Reset()
```

#### func (*UpstreamDestination) String

```go
func (m *UpstreamDestination) String() string
```

#### type UpstreamSpec

```go
type UpstreamSpec *types.Struct
```


#### type VirtualHost

```go
type VirtualHost struct {
	// required, must be unique
	// cannot be "default" unless it refers to the default vhost
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// if this is empty, this host will become / be merged with
	// the default virtualhost who has domains = ["*"]
	Domains []string `protobuf:"bytes,2,rep,name=domains" json:"domains,omitempty"`
	// require at least 1 route
	Routes []*Route `protobuf:"bytes,3,rep,name=routes" json:"routes,omitempty"`
	// optional
	SslConfig *SSLConfig `protobuf:"bytes,4,opt,name=ssl_config,json=sslConfig" json:"ssl_config,omitempty"`
	// read only
	Status   *Status   `protobuf:"bytes,5,opt,name=status" json:"status,omitempty"`
	Metadata *Metadata `protobuf:"bytes,6,opt,name=metadata" json:"metadata,omitempty"`
}
```


#### func (*VirtualHost) Descriptor

```go
func (*VirtualHost) Descriptor() ([]byte, []int)
```

#### func (*VirtualHost) Equal

```go
func (this *VirtualHost) Equal(that interface{}) bool
```

#### func (*VirtualHost) GetDomains

```go
func (m *VirtualHost) GetDomains() []string
```

#### func (*VirtualHost) GetMetadata

```go
func (m *VirtualHost) GetMetadata() *Metadata
```

#### func (*VirtualHost) GetName

```go
func (m *VirtualHost) GetName() string
```

#### func (*VirtualHost) GetRoutes

```go
func (m *VirtualHost) GetRoutes() []*Route
```

#### func (*VirtualHost) GetSslConfig

```go
func (m *VirtualHost) GetSslConfig() *SSLConfig
```

#### func (*VirtualHost) GetStatus

```go
func (m *VirtualHost) GetStatus() *Status
```

#### func (*VirtualHost) ProtoMessage

```go
func (*VirtualHost) ProtoMessage()
```

#### func (*VirtualHost) Reset

```go
func (m *VirtualHost) Reset()
```

#### func (*VirtualHost) String

```go
func (m *VirtualHost) String() string
```

#### type WeightedDestination

```go
type WeightedDestination struct {
	*Destination `protobuf:"bytes,1,opt,name=destination,embedded=destination" json:"destination,omitempty"`
	Weight       uint32 `protobuf:"varint,2,opt,name=weight,proto3" json:"weight,omitempty"`
}
```


#### func (*WeightedDestination) Descriptor

```go
func (*WeightedDestination) Descriptor() ([]byte, []int)
```

#### func (*WeightedDestination) Equal

```go
func (this *WeightedDestination) Equal(that interface{}) bool
```

#### func (*WeightedDestination) GetWeight

```go
func (m *WeightedDestination) GetWeight() uint32
```

#### func (*WeightedDestination) ProtoMessage

```go
func (*WeightedDestination) ProtoMessage()
```

#### func (*WeightedDestination) Reset

```go
func (m *WeightedDestination) Reset()
```

#### func (*WeightedDestination) String

```go
func (m *WeightedDestination) String() string
```
