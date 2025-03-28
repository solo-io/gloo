// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gateway/api/v1/external_options.proto

package v1

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	_ "github.com/solo-io/protoc-gen-ext/extproto"
	v11 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The **VirtualHostOption** holds `options` configuration for a VirtualHost.
// VirtualHosts can inherit `options` config from `VirtualHostOption` objects by delegating to them.
//
// When a VirtualHost delegates to an external VirtualHostOptions object, any options config
// defined on the VirtualHost will override the external options config.
// Similarly, `VirtualHostOption`s which are delegated to first will get priority over following `VirtualHostOption`s.
//
// An example configuration with a VirtualService with its own options config and delegating to two VirtualHostOption objects:
//
// ```yaml
// apiVersion: gateway.solo.io/v1
// kind: VirtualService
// metadata:
//
//	name: http
//	namespace: gloo-system
//
// spec:
//
//	virtualHost:
//	  domains:
//	  - '*'
//	  options:
//	    headerManipulation:
//	      requestHeadersToRemove: ["header-from-vhost"]
//	  optionsConfigRefs:
//	    delegateOptions:
//	      - name: virtualhost-external-options-1
//	        namespace: opt-namespace
//	      - name: virtualhost-external-options-2
//	        namespace: opt-namespace
//
// ```
//
// ```yaml
// apiVersion: gateway.solo.io/v1
// kind: VirtualHostOption
// metadata:
//
//	name: virtualhost-external-options-1
//	namespace: opt-namespace
//
// spec:
//
//	options:
//	  headerManipulation:
//	    requestHeadersToRemove: ["header-from-external-options1"]
//	  cors:
//	    exposeHeaders:
//	      - header-from-extopt1
//	    allowOrigin:
//	      - 'https://solo.io'
//
// ```
//
// ```yaml
// apiVersion: gateway.solo.io/v1
// kind: VirtualHostOption
// metadata:
//
//	name: virtualhost-external-options-2
//	namespace: opt-namespace
//
// spec:
//
//	options:
//	  headerManipulation:
//	    requestHeadersToRemove: ["header-from-external-options2"]
//	  cors:
//	    exposeHeaders:
//	      - header-from-extopt2
//	    maxAge: 2s
//	    allowOrigin:
//	      - 'https://solo.io'
//	  transformations:
//	    requestTransformation:
//	      transformationTemplate:
//	        headers:
//	          x-header-added-in-opt2:
//	            text: this header was added in the VirtualHostOption object - #2
//
// ```
//
// The final virtual host options (visible in the Proxy CR) would be:
// ```yaml
// spec:
//
//	virtualHost:
//	  domains:
//	  - '*'
//	  options:
//	    # from Virtual host options
//	    headerManipulation:
//	      requestHeadersToRemove: ["header-from-vhost"]
//	    # from delegated virtualhost-external-options-1
//	    cors:
//	      exposeHeaders:
//	        - header-from-extopt1
//	      allowOrigin:
//	        - 'https://solo.io'
//	    # from delegated virtualhost-external-options-2
//	    transformations:
//	      requestTransformation:
//	        transformationTemplate:
//	          headers:
//	            x-header-added-in-opt2:
//	              text: this header was added in the VirtualHostOption object - #2
//
// ```
//
// Notice how the order of VirtualHostOption delegations matters, and that the VirtualHost-level config overrides all delegated configs.
type VirtualHostOption struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// NamespacedStatuses indicates the validation status of this resource.
	// NamespacedStatuses is read-only by clients, and set by gateway during validation
	NamespacedStatuses *core.NamespacedStatuses `protobuf:"bytes,4,opt,name=namespaced_statuses,json=namespacedStatuses,proto3" json:"namespaced_statuses,omitempty"`
	// Metadata contains the object metadata for this resource
	Metadata *core.Metadata `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// VirtualHost options. See VirtualHost for delegation behavior.
	Options *v1.VirtualHostOptions `protobuf:"bytes,3,opt,name=options,proto3" json:"options,omitempty"`
	// When using Kubernetes Gateway API mode, targetRefs can be used to attach this VirtualHostOption
	// to a gateway.networking.k8s.io/Gateway object. The `options` specified will then be used
	// for all `Listeners` in the `Gateway` unless `targetRef.sectionName` is specified.
	// NOTE: This is a repeated field but currently ONLY supports a single targetRef.
	// If multiple targetRefs are provided, only the first in the list will be used.
	TargetRefs    []*v11.PolicyTargetReferenceWithSectionName `protobuf:"bytes,5,rep,name=target_refs,json=targetRefs,proto3" json:"target_refs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VirtualHostOption) Reset() {
	*x = VirtualHostOption{}
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VirtualHostOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VirtualHostOption) ProtoMessage() {}

func (x *VirtualHostOption) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VirtualHostOption.ProtoReflect.Descriptor instead.
func (*VirtualHostOption) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescGZIP(), []int{0}
}

func (x *VirtualHostOption) GetNamespacedStatuses() *core.NamespacedStatuses {
	if x != nil {
		return x.NamespacedStatuses
	}
	return nil
}

func (x *VirtualHostOption) GetMetadata() *core.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *VirtualHostOption) GetOptions() *v1.VirtualHostOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *VirtualHostOption) GetTargetRefs() []*v11.PolicyTargetReferenceWithSectionName {
	if x != nil {
		return x.TargetRefs
	}
	return nil
}

// The **RouteOption** holds `options` configuration for a Route.
// Routes can inherit `options` config from `RouteOption` objects by delegating to them.
//
// When a Route delegates to an external RouteOptions object, any options config
// defined on the Route will override the external options config.
// Similarly, `RouteOption`s which are delegated to first will get priority over following `RouteOption`s.
//
// An example configuration with a Route with its own options config and delegating to two RouteOption objects:
//
// ```yaml
// apiVersion: gateway.solo.io/v1
// kind: VirtualService
// metadata:
//
//	name: http
//	namespace: gloo-system
//
// spec:
//
//	virtualHost:
//	  domains:
//	  - '*'
//	  routes:
//	  - matchers:
//	    - prefix: /
//	    options:
//	      headerManipulation:
//	        requestHeadersToRemove: ["header-from-route"]
//	    optionsConfigRefs:
//	      delegateOptions:
//	        - name: route-external-options-1
//	          namespace: opt-namespace
//	        - name: route-external-options-2
//	          namespace: opt-namespace
//
// ```
//
// ```yaml
// apiVersion: gateway.solo.io/v1
// kind: RouteOption
// metadata:
//
//	name: route-external-options-1
//	namespace: opt-namespace
//
// spec:
//
//	options:
//	  headerManipulation:
//	    requestHeadersToRemove: ["header-from-external-options1"]
//	  cors:
//	    exposeHeaders:
//	      - header-from-extopt1
//	    allowOrigin:
//	      - 'https://solo.io'
//
// ```
//
// ```yaml
// apiVersion: gateway.solo.io/v1
// kind: RouteOption
// metadata:
//
//	name: route-external-options-2
//	namespace: opt-namespace
//
// spec:
//
//	options:
//	  headerManipulation:
//	    requestHeadersToRemove: ["header-from-external-options2"]
//	  cors:
//	    exposeHeaders:
//	      - header-from-extopt2
//	    maxAge: 2s
//	    allowOrigin:
//	      - 'https://solo.io'
//	  transformations:
//	    requestTransformation:
//	      transformationTemplate:
//	        headers:
//	          x-header-added-in-opt2:
//	            text: this header was added in the RouteOption object - #2
//
// ```
//
// The final route options would bewould be:
// ```yaml
// routes:
//   - matchers:
//   - prefix: /
//     options:
//     # from Route options
//     headerManipulation:
//     requestHeadersToRemove: ["header-from-route"]
//     # from delegated route-external-options-1
//     cors:
//     exposeHeaders:
//   - header-from-extopt1
//     allowOrigin:
//   - 'https://solo.io'
//     # from delegated route-external-options-2
//     transformations:
//     requestTransformation:
//     transformationTemplate:
//     headers:
//     x-header-added-in-opt2:
//     text: this header was added in the Route object - #2
//
// ```
//
// Notice how the order of RouteOption delegations matters, and that the Route-level option config overrides all delegated option configs.
type RouteOption struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// NamespacedStatuses indicates the validation status of this resource.
	// NamespacedStatuses is read-only by clients, and set by gateway during validation
	NamespacedStatuses *core.NamespacedStatuses `protobuf:"bytes,4,opt,name=namespaced_statuses,json=namespacedStatuses,proto3" json:"namespaced_statuses,omitempty"`
	// Metadata contains the object metadata for this resource
	Metadata *core.Metadata `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Route options. See Route for delegation behavior.
	Options *v1.RouteOptions `protobuf:"bytes,3,opt,name=options,proto3" json:"options,omitempty"`
	// When using Kubernetes Gateway API mode, targetRefs can be used to attach this RouteOption
	// to a gateway.networking.k8s.io/HTTPRoute object. The `options` specified will then be used
	// for all `Rules` in the `HTTPRoute`.
	// NOTE: This is a repeated field but currently ONLY supports a single targetRef.
	// If multiple targetRefs are provided, only the first in the list will be used.
	TargetRefs    []*v11.PolicyTargetReference `protobuf:"bytes,5,rep,name=target_refs,json=targetRefs,proto3" json:"target_refs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RouteOption) Reset() {
	*x = RouteOption{}
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RouteOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteOption) ProtoMessage() {}

func (x *RouteOption) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteOption.ProtoReflect.Descriptor instead.
func (*RouteOption) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescGZIP(), []int{1}
}

func (x *RouteOption) GetNamespacedStatuses() *core.NamespacedStatuses {
	if x != nil {
		return x.NamespacedStatuses
	}
	return nil
}

func (x *RouteOption) GetMetadata() *core.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *RouteOption) GetOptions() *v1.RouteOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *RouteOption) GetTargetRefs() []*v11.PolicyTargetReference {
	if x != nil {
		return x.TargetRefs
	}
	return nil
}

type ListenerOption struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Metadata contains the object metadata for this resource
	Metadata *core.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Listener options
	Options *v1.ListenerOptions `protobuf:"bytes,2,opt,name=options,proto3" json:"options,omitempty"`
	// When using Kubernetes Gateway API mode, targetRefs can be used to attach this ListenerOption
	// to a gateway.networking.k8s.io/Gateway object. The `options` specified will then be used
	// for all `Listeners` in the `Gateway` unless `targetRef.sectionName` is specified.
	// NOTE: This is a repeated field but currently ONLY supports a single targetRef.
	// If multiple targetRefs are provided, only the first in the list will be used.
	TargetRefs    []*v11.PolicyTargetReferenceWithSectionName `protobuf:"bytes,3,rep,name=target_refs,json=targetRefs,proto3" json:"target_refs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListenerOption) Reset() {
	*x = ListenerOption{}
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListenerOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenerOption) ProtoMessage() {}

func (x *ListenerOption) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenerOption.ProtoReflect.Descriptor instead.
func (*ListenerOption) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescGZIP(), []int{2}
}

func (x *ListenerOption) GetMetadata() *core.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *ListenerOption) GetOptions() *v1.ListenerOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *ListenerOption) GetTargetRefs() []*v11.PolicyTargetReferenceWithSectionName {
	if x != nil {
		return x.TargetRefs
	}
	return nil
}

type HttpListenerOption struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Metadata contains the object metadata for this resource
	Metadata *core.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// HttpListener options
	Options *v1.HttpListenerOptions `protobuf:"bytes,2,opt,name=options,proto3" json:"options,omitempty"`
	// When using Kubernetes Gateway API mode, targetRef can be used to attach this VirtualHostOption
	// to a gateway.networking.k8s.io/Gateway object. The `options` specified will then be used
	// for all `Listeners` in the `Gateway` unless `targetRef.sectionName` is specified.
	// NOTE: This is a repeated field but currently ONLY supports a single targetRef.
	// If multiple targetRefs are provided, only the first in the list will be used.
	TargetRefs    []*v11.PolicyTargetReferenceWithSectionName `protobuf:"bytes,3,rep,name=target_refs,json=targetRefs,proto3" json:"target_refs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HttpListenerOption) Reset() {
	*x = HttpListenerOption{}
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HttpListenerOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpListenerOption) ProtoMessage() {}

func (x *HttpListenerOption) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpListenerOption.ProtoReflect.Descriptor instead.
func (*HttpListenerOption) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescGZIP(), []int{3}
}

func (x *HttpListenerOption) GetMetadata() *core.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *HttpListenerOption) GetOptions() *v1.HttpListenerOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *HttpListenerOption) GetTargetRefs() []*v11.PolicyTargetReferenceWithSectionName {
	if x != nil {
		return x.TargetRefs
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDesc = string([]byte{
	0x0a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x1a, 0x12, 0x65, 0x78, 0x74, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x31, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69,
	0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x6b, 0x69, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x6b, 0x69, 0x74, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f,
	0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x6b, 0x69, 0x74, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f,
	0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x6b, 0x69, 0x74, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x6b, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f,
	0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f,
	0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x40, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67,
	0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f,
	0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f,
	0x2f, 0x73, 0x6b, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe0, 0x02, 0x0a,
	0x11, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x48, 0x6f, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x57, 0x0a, 0x13, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x64,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65,
	0x73, 0x42, 0x04, 0xb8, 0xf5, 0x04, 0x01, 0x52, 0x12, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x3a, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e,
	0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x48, 0x6f, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x58, 0x0a, 0x0b, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x37, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6b, 0x76, 0x32, 0x2e, 0x73, 0x6f, 0x6c,
	0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x54, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x57, 0x69, 0x74, 0x68, 0x53, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x52, 0x65, 0x66, 0x73, 0x3a, 0x22, 0x82, 0xf1, 0x04, 0x1e, 0x0a, 0x06, 0x76, 0x68, 0x6f,
	0x70, 0x74, 0x73, 0x12, 0x14, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x68, 0x6f, 0x73,
	0x74, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x22,
	0xbe, 0x02, 0x0a, 0x0b, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x57, 0x0a, 0x13, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x64, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x42, 0x04,
	0xb8, 0xf5, 0x04, 0x01, 0x52, 0x12, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x64,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x34, 0x0a, 0x07,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x49, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x65, 0x66,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73,
	0x6b, 0x76, 0x32, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x66, 0x73, 0x3a, 0x1b, 0x82,
	0xf1, 0x04, 0x17, 0x0a, 0x06, 0x72, 0x74, 0x6f, 0x70, 0x74, 0x73, 0x12, 0x0d, 0x72, 0x6f, 0x75,
	0x74, 0x65, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02,
	0x22, 0xf8, 0x01, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6f, 0x6c,
	0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x37, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x58, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6b, 0x76,
	0x32, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x57,
	0x69, 0x74, 0x68, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x0a,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x66, 0x73, 0x3a, 0x1f, 0x82, 0xf1, 0x04, 0x1b,
	0x0a, 0x07, 0x6c, 0x69, 0x73, 0x6f, 0x70, 0x74, 0x73, 0x12, 0x10, 0x6c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x86, 0x02, 0x0a, 0x12,
	0x48, 0x74, 0x74, 0x70, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3b, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x4c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x58, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x65,
	0x66, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x73, 0x6b, 0x76, 0x32, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x57, 0x69, 0x74, 0x68, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d,
	0x65, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x66, 0x73, 0x3a, 0x25, 0x82,
	0xf1, 0x04, 0x21, 0x0a, 0x08, 0x68, 0x6c, 0x69, 0x73, 0x6f, 0x70, 0x74, 0x73, 0x12, 0x15, 0x68,
	0x74, 0x74, 0x70, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x42, 0x41, 0xb8, 0xf5, 0x04, 0x01, 0xc0, 0xf5, 0x04, 0x01, 0xd0, 0xf5,
	0x04, 0x01, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_goTypes = []any{
	(*VirtualHostOption)(nil),                        // 0: gateway.solo.io.VirtualHostOption
	(*RouteOption)(nil),                              // 1: gateway.solo.io.RouteOption
	(*ListenerOption)(nil),                           // 2: gateway.solo.io.ListenerOption
	(*HttpListenerOption)(nil),                       // 3: gateway.solo.io.HttpListenerOption
	(*core.NamespacedStatuses)(nil),                  // 4: core.solo.io.NamespacedStatuses
	(*core.Metadata)(nil),                            // 5: core.solo.io.Metadata
	(*v1.VirtualHostOptions)(nil),                    // 6: gloo.solo.io.VirtualHostOptions
	(*v11.PolicyTargetReferenceWithSectionName)(nil), // 7: core.skv2.solo.io.PolicyTargetReferenceWithSectionName
	(*v1.RouteOptions)(nil),                          // 8: gloo.solo.io.RouteOptions
	(*v11.PolicyTargetReference)(nil),                // 9: core.skv2.solo.io.PolicyTargetReference
	(*v1.ListenerOptions)(nil),                       // 10: gloo.solo.io.ListenerOptions
	(*v1.HttpListenerOptions)(nil),                   // 11: gloo.solo.io.HttpListenerOptions
}
var file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_depIdxs = []int32{
	4,  // 0: gateway.solo.io.VirtualHostOption.namespaced_statuses:type_name -> core.solo.io.NamespacedStatuses
	5,  // 1: gateway.solo.io.VirtualHostOption.metadata:type_name -> core.solo.io.Metadata
	6,  // 2: gateway.solo.io.VirtualHostOption.options:type_name -> gloo.solo.io.VirtualHostOptions
	7,  // 3: gateway.solo.io.VirtualHostOption.target_refs:type_name -> core.skv2.solo.io.PolicyTargetReferenceWithSectionName
	4,  // 4: gateway.solo.io.RouteOption.namespaced_statuses:type_name -> core.solo.io.NamespacedStatuses
	5,  // 5: gateway.solo.io.RouteOption.metadata:type_name -> core.solo.io.Metadata
	8,  // 6: gateway.solo.io.RouteOption.options:type_name -> gloo.solo.io.RouteOptions
	9,  // 7: gateway.solo.io.RouteOption.target_refs:type_name -> core.skv2.solo.io.PolicyTargetReference
	5,  // 8: gateway.solo.io.ListenerOption.metadata:type_name -> core.solo.io.Metadata
	10, // 9: gateway.solo.io.ListenerOption.options:type_name -> gloo.solo.io.ListenerOptions
	7,  // 10: gateway.solo.io.ListenerOption.target_refs:type_name -> core.skv2.solo.io.PolicyTargetReferenceWithSectionName
	5,  // 11: gateway.solo.io.HttpListenerOption.metadata:type_name -> core.solo.io.Metadata
	11, // 12: gateway.solo.io.HttpListenerOption.options:type_name -> gloo.solo.io.HttpListenerOptions
	7,  // 13: gateway.solo.io.HttpListenerOption.target_refs:type_name -> core.skv2.solo.io.PolicyTargetReferenceWithSectionName
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_init() }
func file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_init() {
	if File_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto = out.File
	file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gateway_api_v1_external_options_proto_depIdxs = nil
}
