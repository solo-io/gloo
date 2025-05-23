// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit.proto

package ratelimit

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	local_ratelimit "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/local_ratelimit"
	_ "github.com/solo-io/protoc-gen-ext/extproto"
	v1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Basic rate-limiting API
type IngressRateLimit struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	AuthorizedLimits *v1alpha1.RateLimit    `protobuf:"bytes,1,opt,name=authorized_limits,json=authorizedLimits,proto3" json:"authorized_limits,omitempty"`
	AnonymousLimits  *v1alpha1.RateLimit    `protobuf:"bytes,2,opt,name=anonymous_limits,json=anonymousLimits,proto3" json:"anonymous_limits,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *IngressRateLimit) Reset() {
	*x = IngressRateLimit{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IngressRateLimit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IngressRateLimit) ProtoMessage() {}

func (x *IngressRateLimit) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IngressRateLimit.ProtoReflect.Descriptor instead.
func (*IngressRateLimit) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{0}
}

func (x *IngressRateLimit) GetAuthorizedLimits() *v1alpha1.RateLimit {
	if x != nil {
		return x.AuthorizedLimits
	}
	return nil
}

func (x *IngressRateLimit) GetAnonymousLimits() *v1alpha1.RateLimit {
	if x != nil {
		return x.AnonymousLimits
	}
	return nil
}

type Settings struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	RatelimitServerRef *core.ResourceRef      `protobuf:"bytes,1,opt,name=ratelimit_server_ref,json=ratelimitServerRef,proto3" json:"ratelimit_server_ref,omitempty"`
	RequestTimeout     *durationpb.Duration   `protobuf:"bytes,2,opt,name=request_timeout,json=requestTimeout,proto3" json:"request_timeout,omitempty"`
	DenyOnFail         bool                   `protobuf:"varint,3,opt,name=deny_on_fail,json=denyOnFail,proto3" json:"deny_on_fail,omitempty"`
	// Set this to true to return Envoy's X-RateLimit headers to the downstream.
	// reference docs here: https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ratelimit/v3/rate_limit.proto.html#envoy-v3-api-field-extensions-filters-http-ratelimit-v3-ratelimit-enable-x-ratelimit-headers
	EnableXRatelimitHeaders bool `protobuf:"varint,4,opt,name=enable_x_ratelimit_headers,json=enableXRatelimitHeaders,proto3" json:"enable_x_ratelimit_headers,omitempty"`
	// Set this is set to true if you would like to rate limit traffic before applying external auth to it.
	// *Note*: When this is true, you will lose some features like being able to rate limit a request based on its auth state
	RateLimitBeforeAuth bool `protobuf:"varint,9,opt,name=rate_limit_before_auth,json=rateLimitBeforeAuth,proto3" json:"rate_limit_before_auth,omitempty"`
	// Types that are valid to be assigned to ServiceType:
	//
	//	*Settings_GrpcService
	ServiceType   isSettings_ServiceType `protobuf_oneof:"service_type"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Settings) Reset() {
	*x = Settings{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Settings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Settings) ProtoMessage() {}

func (x *Settings) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Settings.ProtoReflect.Descriptor instead.
func (*Settings) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{1}
}

func (x *Settings) GetRatelimitServerRef() *core.ResourceRef {
	if x != nil {
		return x.RatelimitServerRef
	}
	return nil
}

func (x *Settings) GetRequestTimeout() *durationpb.Duration {
	if x != nil {
		return x.RequestTimeout
	}
	return nil
}

func (x *Settings) GetDenyOnFail() bool {
	if x != nil {
		return x.DenyOnFail
	}
	return false
}

func (x *Settings) GetEnableXRatelimitHeaders() bool {
	if x != nil {
		return x.EnableXRatelimitHeaders
	}
	return false
}

func (x *Settings) GetRateLimitBeforeAuth() bool {
	if x != nil {
		return x.RateLimitBeforeAuth
	}
	return false
}

func (x *Settings) GetServiceType() isSettings_ServiceType {
	if x != nil {
		return x.ServiceType
	}
	return nil
}

func (x *Settings) GetGrpcService() *GrpcService {
	if x != nil {
		if x, ok := x.ServiceType.(*Settings_GrpcService); ok {
			return x.GrpcService
		}
	}
	return nil
}

type isSettings_ServiceType interface {
	isSettings_ServiceType()
}

type Settings_GrpcService struct {
	// Optional gRPC settings used when calling the ratelimit server.
	GrpcService *GrpcService `protobuf:"bytes,10,opt,name=grpc_service,json=grpcService,proto3,oneof"`
}

func (*Settings_GrpcService) isSettings_ServiceType() {}

type GrpcService struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Set the authority header when calling the gRPC service.
	Authority     string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GrpcService) Reset() {
	*x = GrpcService{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GrpcService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcService) ProtoMessage() {}

func (x *GrpcService) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcService.ProtoReflect.Descriptor instead.
func (*GrpcService) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{2}
}

func (x *GrpcService) GetAuthority() string {
	if x != nil {
		return x.Authority
	}
	return ""
}

// API based on Envoy's rate-limit service API. (reference here: https://github.com/lyft/ratelimit#configuration)
// Sample configuration below:
//
// descriptors:
//   - key: account_id
//     descriptors:
//   - key: plan
//     value: BASIC
//     rateLimit:
//     requestsPerUnit: 1
//     unit: MINUTE
//   - key: plan
//     value: PLUS
//     rateLimit:
//     requestsPerUnit: 20
//     unit: MINUTE
type ServiceSettings struct {
	state          protoimpl.MessageState    `protogen:"open.v1"`
	Descriptors    []*v1alpha1.Descriptor    `protobuf:"bytes,1,rep,name=descriptors,proto3" json:"descriptors,omitempty"`
	SetDescriptors []*v1alpha1.SetDescriptor `protobuf:"bytes,2,rep,name=set_descriptors,json=setDescriptors,proto3" json:"set_descriptors,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *ServiceSettings) Reset() {
	*x = ServiceSettings{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServiceSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceSettings) ProtoMessage() {}

func (x *ServiceSettings) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceSettings.ProtoReflect.Descriptor instead.
func (*ServiceSettings) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{3}
}

func (x *ServiceSettings) GetDescriptors() []*v1alpha1.Descriptor {
	if x != nil {
		return x.Descriptors
	}
	return nil
}

func (x *ServiceSettings) GetSetDescriptors() []*v1alpha1.SetDescriptor {
	if x != nil {
		return x.SetDescriptors
	}
	return nil
}

// A list of references to `RateLimitConfig` resources.
// Each resource represents a rate limit policy that will be independently enforced.
type RateLimitConfigRefs struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Refs          []*RateLimitConfigRef  `protobuf:"bytes,1,rep,name=refs,proto3" json:"refs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RateLimitConfigRefs) Reset() {
	*x = RateLimitConfigRefs{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RateLimitConfigRefs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitConfigRefs) ProtoMessage() {}

func (x *RateLimitConfigRefs) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitConfigRefs.ProtoReflect.Descriptor instead.
func (*RateLimitConfigRefs) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{4}
}

func (x *RateLimitConfigRefs) GetRefs() []*RateLimitConfigRef {
	if x != nil {
		return x.Refs
	}
	return nil
}

// A reference to a `RateLimitConfig` resource.
type RateLimitConfigRef struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace     string                 `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RateLimitConfigRef) Reset() {
	*x = RateLimitConfigRef{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RateLimitConfigRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitConfigRef) ProtoMessage() {}

func (x *RateLimitConfigRef) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitConfigRef.ProtoReflect.Descriptor instead.
func (*RateLimitConfigRef) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{5}
}

func (x *RateLimitConfigRef) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RateLimitConfigRef) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

// Use this field if you want to inline the Envoy rate limits for this VirtualHost.
// Note that this does not configure the rate limit server. If you are running Gloo Enterprise, you need to
// specify the server configuration via the appropriate field in the Gloo `Settings` resource. If you are
// running a custom rate limit server you need to configure it yourself.
type RateLimitVhostExtension struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Define individual rate limits here. Each rate limit will be evaluated, if any rate limit
	// would be throttled, the entire request returns a 429 (gets throttled)
	RateLimits []*v1alpha1.RateLimitActions `protobuf:"bytes,1,rep,name=rate_limits,json=rateLimits,proto3" json:"rate_limits,omitempty"`
	// The token bucket configuration to use for local rate limiting requests.
	// These options provide the ability to locally rate limit the connections in envoy. Each request processed by the filter consumes a single token.
	// If the token is available, the request will be allowed. If no tokens are available, the request will receive the configured rate limit status.
	// This overrides any local rate limit configured on the gateway and requests to this vHost do not count against requests to the gateway's http local rate limit.
	// All routes that are part of this vHost will share this rate limit unless explicity configured with another limit.
	LocalRatelimit *local_ratelimit.TokenBucket `protobuf:"bytes,2,opt,name=local_ratelimit,json=localRatelimit,proto3" json:"local_ratelimit,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *RateLimitVhostExtension) Reset() {
	*x = RateLimitVhostExtension{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RateLimitVhostExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitVhostExtension) ProtoMessage() {}

func (x *RateLimitVhostExtension) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitVhostExtension.ProtoReflect.Descriptor instead.
func (*RateLimitVhostExtension) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{6}
}

func (x *RateLimitVhostExtension) GetRateLimits() []*v1alpha1.RateLimitActions {
	if x != nil {
		return x.RateLimits
	}
	return nil
}

func (x *RateLimitVhostExtension) GetLocalRatelimit() *local_ratelimit.TokenBucket {
	if x != nil {
		return x.LocalRatelimit
	}
	return nil
}

// Use this field if you want to inline the Envoy rate limits for this Route.
// Note that this does not configure the rate limit server. If you are running Gloo Enterprise, you need to
// specify the server configuration via the appropriate field in the Gloo `Settings` resource. If you are
// running a custom rate limit server you need to configure it yourself.
type RateLimitRouteExtension struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Whether or not to include rate limits as defined on the VirtualHost in addition to rate limits on the Route.
	IncludeVhRateLimits bool `protobuf:"varint,1,opt,name=include_vh_rate_limits,json=includeVhRateLimits,proto3" json:"include_vh_rate_limits,omitempty"`
	// Define individual rate limits here. Each rate limit will be evaluated, if any rate limit
	// would be throttled, the entire request returns a 429 (gets throttled)
	RateLimits []*v1alpha1.RateLimitActions `protobuf:"bytes,2,rep,name=rate_limits,json=rateLimits,proto3" json:"rate_limits,omitempty"`
	// The token bucket configuration to use for local rate limiting requests.
	// These options provide the ability to locally rate limit the connections in envoy. Each request processed by the filter consumes a single token.
	// If the token is available, the request will be allowed. If no tokens are available, the request will receive the configured rate limit status.
	// This overrides any local rate limit configured on the vHost or gateway and requests to this route do not count against requests to the vHost or gateway's http local rate limit.
	LocalRatelimit *local_ratelimit.TokenBucket `protobuf:"bytes,3,opt,name=local_ratelimit,json=localRatelimit,proto3" json:"local_ratelimit,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *RateLimitRouteExtension) Reset() {
	*x = RateLimitRouteExtension{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RateLimitRouteExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitRouteExtension) ProtoMessage() {}

func (x *RateLimitRouteExtension) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitRouteExtension.ProtoReflect.Descriptor instead.
func (*RateLimitRouteExtension) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP(), []int{7}
}

func (x *RateLimitRouteExtension) GetIncludeVhRateLimits() bool {
	if x != nil {
		return x.IncludeVhRateLimits
	}
	return false
}

func (x *RateLimitRouteExtension) GetRateLimits() []*v1alpha1.RateLimitActions {
	if x != nil {
		return x.RateLimits
	}
	return nil
}

func (x *RateLimitRouteExtension) GetLocalRatelimit() *local_ratelimit.TokenBucket {
	if x != nil {
		return x.LocalRatelimit
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDesc = "" +
	"\n" +
	"Ygithub.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit.proto\x12\x1eratelimit.options.gloo.solo.io\x1aFgithub.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/ratelimit.proto\x1aZgithub.com/solo-io/gloo/projects/gloo/api/v1/options/local_ratelimit/local_ratelimit.proto\x1a,github.com/solo-io/solo-kit/api/v1/ref.proto\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x1egoogle/protobuf/duration.proto\x1a\x12extproto/ext.proto\"\xae\x01\n" +
	"\x10IngressRateLimit\x12M\n" +
	"\x11authorized_limits\x18\x01 \x01(\v2 .ratelimit.api.solo.io.RateLimitR\x10authorizedLimits\x12K\n" +
	"\x10anonymous_limits\x18\x02 \x01(\v2 .ratelimit.api.solo.io.RateLimitR\x0fanonymousLimits\"\x91\x03\n" +
	"\bSettings\x12K\n" +
	"\x14ratelimit_server_ref\x18\x01 \x01(\v2\x19.core.solo.io.ResourceRefR\x12ratelimitServerRef\x12B\n" +
	"\x0frequest_timeout\x18\x02 \x01(\v2\x19.google.protobuf.DurationR\x0erequestTimeout\x12 \n" +
	"\fdeny_on_fail\x18\x03 \x01(\bR\n" +
	"denyOnFail\x12;\n" +
	"\x1aenable_x_ratelimit_headers\x18\x04 \x01(\bR\x17enableXRatelimitHeaders\x123\n" +
	"\x16rate_limit_before_auth\x18\t \x01(\bR\x13rateLimitBeforeAuth\x12P\n" +
	"\fgrpc_service\x18\n" +
	" \x01(\v2+.ratelimit.options.gloo.solo.io.GrpcServiceH\x00R\vgrpcServiceB\x0e\n" +
	"\fservice_type\"+\n" +
	"\vGrpcService\x12\x1c\n" +
	"\tauthority\x18\x01 \x01(\tR\tauthority\"\xa5\x01\n" +
	"\x0fServiceSettings\x12C\n" +
	"\vdescriptors\x18\x01 \x03(\v2!.ratelimit.api.solo.io.DescriptorR\vdescriptors\x12M\n" +
	"\x0fset_descriptors\x18\x02 \x03(\v2$.ratelimit.api.solo.io.SetDescriptorR\x0esetDescriptors\"]\n" +
	"\x13RateLimitConfigRefs\x12F\n" +
	"\x04refs\x18\x01 \x03(\v22.ratelimit.options.gloo.solo.io.RateLimitConfigRefR\x04refs\"F\n" +
	"\x12RateLimitConfigRef\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x1c\n" +
	"\tnamespace\x18\x02 \x01(\tR\tnamespace\"\xbf\x01\n" +
	"\x17RateLimitVhostExtension\x12H\n" +
	"\vrate_limits\x18\x01 \x03(\v2'.ratelimit.api.solo.io.RateLimitActionsR\n" +
	"rateLimits\x12Z\n" +
	"\x0flocal_ratelimit\x18\x02 \x01(\v21.local_ratelimit.options.gloo.solo.io.TokenBucketR\x0elocalRatelimit\"\xf4\x01\n" +
	"\x17RateLimitRouteExtension\x123\n" +
	"\x16include_vh_rate_limits\x18\x01 \x01(\bR\x13includeVhRateLimits\x12H\n" +
	"\vrate_limits\x18\x02 \x03(\v2'.ratelimit.api.solo.io.RateLimitActionsR\n" +
	"rateLimits\x12Z\n" +
	"\x0flocal_ratelimit\x18\x03 \x01(\v21.local_ratelimit.options.gloo.solo.io.TokenBucketR\x0elocalRatelimitB[\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01ZMgithub.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimitb\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_goTypes = []any{
	(*IngressRateLimit)(nil),            // 0: ratelimit.options.gloo.solo.io.IngressRateLimit
	(*Settings)(nil),                    // 1: ratelimit.options.gloo.solo.io.Settings
	(*GrpcService)(nil),                 // 2: ratelimit.options.gloo.solo.io.GrpcService
	(*ServiceSettings)(nil),             // 3: ratelimit.options.gloo.solo.io.ServiceSettings
	(*RateLimitConfigRefs)(nil),         // 4: ratelimit.options.gloo.solo.io.RateLimitConfigRefs
	(*RateLimitConfigRef)(nil),          // 5: ratelimit.options.gloo.solo.io.RateLimitConfigRef
	(*RateLimitVhostExtension)(nil),     // 6: ratelimit.options.gloo.solo.io.RateLimitVhostExtension
	(*RateLimitRouteExtension)(nil),     // 7: ratelimit.options.gloo.solo.io.RateLimitRouteExtension
	(*v1alpha1.RateLimit)(nil),          // 8: ratelimit.api.solo.io.RateLimit
	(*core.ResourceRef)(nil),            // 9: core.solo.io.ResourceRef
	(*durationpb.Duration)(nil),         // 10: google.protobuf.Duration
	(*v1alpha1.Descriptor)(nil),         // 11: ratelimit.api.solo.io.Descriptor
	(*v1alpha1.SetDescriptor)(nil),      // 12: ratelimit.api.solo.io.SetDescriptor
	(*v1alpha1.RateLimitActions)(nil),   // 13: ratelimit.api.solo.io.RateLimitActions
	(*local_ratelimit.TokenBucket)(nil), // 14: local_ratelimit.options.gloo.solo.io.TokenBucket
}
var file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_depIdxs = []int32{
	8,  // 0: ratelimit.options.gloo.solo.io.IngressRateLimit.authorized_limits:type_name -> ratelimit.api.solo.io.RateLimit
	8,  // 1: ratelimit.options.gloo.solo.io.IngressRateLimit.anonymous_limits:type_name -> ratelimit.api.solo.io.RateLimit
	9,  // 2: ratelimit.options.gloo.solo.io.Settings.ratelimit_server_ref:type_name -> core.solo.io.ResourceRef
	10, // 3: ratelimit.options.gloo.solo.io.Settings.request_timeout:type_name -> google.protobuf.Duration
	2,  // 4: ratelimit.options.gloo.solo.io.Settings.grpc_service:type_name -> ratelimit.options.gloo.solo.io.GrpcService
	11, // 5: ratelimit.options.gloo.solo.io.ServiceSettings.descriptors:type_name -> ratelimit.api.solo.io.Descriptor
	12, // 6: ratelimit.options.gloo.solo.io.ServiceSettings.set_descriptors:type_name -> ratelimit.api.solo.io.SetDescriptor
	5,  // 7: ratelimit.options.gloo.solo.io.RateLimitConfigRefs.refs:type_name -> ratelimit.options.gloo.solo.io.RateLimitConfigRef
	13, // 8: ratelimit.options.gloo.solo.io.RateLimitVhostExtension.rate_limits:type_name -> ratelimit.api.solo.io.RateLimitActions
	14, // 9: ratelimit.options.gloo.solo.io.RateLimitVhostExtension.local_ratelimit:type_name -> local_ratelimit.options.gloo.solo.io.TokenBucket
	13, // 10: ratelimit.options.gloo.solo.io.RateLimitRouteExtension.rate_limits:type_name -> ratelimit.api.solo.io.RateLimitActions
	14, // 11: ratelimit.options.gloo.solo.io.RateLimitRouteExtension.local_ratelimit:type_name -> local_ratelimit.options.gloo.solo.io.TokenBucket
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto != nil {
		return
	}
	file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes[1].OneofWrappers = []any{
		(*Settings_GrpcService)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_proto_depIdxs = nil
}
