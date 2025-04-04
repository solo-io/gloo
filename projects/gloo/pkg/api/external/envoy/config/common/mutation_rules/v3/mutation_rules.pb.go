// copied from https://github.com/envoyproxy/envoy/blob/ad89a587aa0177bfdad6b5c968a6aead5d9be7a4/api/envoy/config/common/mutation_rules/v3/mutation_rules.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/common/mutation_rules/v3/mutation_rules.proto

// manually updated pkg

package v3

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	v31 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	_ "github.com/solo-io/gloo/projects/gloo/pkg/api/external/udpa/annotations"
	_ "github.com/solo-io/protoc-gen-ext/extproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The HeaderMutationRules structure specifies what headers may be
// manipulated by a processing filter. This set of rules makes it
// possible to control which modifications a filter may make.
//
// By default, an external processing server may add, modify, or remove
// any header except for an "Envoy internal" header (which is typically
// denoted by an x-envoy prefix) or specific headers that may affect
// further filter processing:
//
// * `host`
// * `:authority`
// * `:scheme`
// * `:method`
//
// Every attempt to add, change, append, or remove a header will be
// tested against the rules here. Disallowed header mutations will be
// ignored unless `disallow_is_error` is set to true.
//
// Attempts to remove headers are further constrained -- regardless of the
// settings, system-defined headers (that start with `:`) and the `host`
// header may never be removed.
//
// In addition, a counter will be incremented whenever a mutation is
// rejected. In the ext_proc filter, that counter is named
// `rejected_header_mutations`.
// [#next-free-field: 8]
type HeaderMutationRules struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// By default, certain headers that could affect processing of subsequent
	// filters or request routing cannot be modified. These headers are
	// `host`, `:authority`, `:scheme`, and `:method`. Setting this parameter
	// to true allows these headers to be modified as well.
	AllowAllRouting *wrapperspb.BoolValue `protobuf:"bytes,1,opt,name=allow_all_routing,json=allowAllRouting,proto3" json:"allow_all_routing,omitempty"`
	// If true, allow modification of envoy internal headers. By default, these
	// start with `x-envoy` but this may be overridden in the `Bootstrap`
	// configuration using the
	// :ref:`header_prefix <envoy_v3_api_field_config.bootstrap.v3.Bootstrap.header_prefix>`
	// field. Default is false.
	AllowEnvoy *wrapperspb.BoolValue `protobuf:"bytes,2,opt,name=allow_envoy,json=allowEnvoy,proto3" json:"allow_envoy,omitempty"`
	// If true, prevent modification of any system header, defined as a header
	// that starts with a `:` character, regardless of any other settings.
	// A processing server may still override the `:status` of an HTTP response
	// using an `ImmediateResponse` message. Default is false.
	DisallowSystem *wrapperspb.BoolValue `protobuf:"bytes,3,opt,name=disallow_system,json=disallowSystem,proto3" json:"disallow_system,omitempty"`
	// If true, prevent modifications of all header values, regardless of any
	// other settings. A processing server may still override the `:status`
	// of an HTTP response using an `ImmediateResponse` message. Default is false.
	DisallowAll *wrapperspb.BoolValue `protobuf:"bytes,4,opt,name=disallow_all,json=disallowAll,proto3" json:"disallow_all,omitempty"`
	// If set, specifically allow any header that matches this regular
	// expression. This overrides all other settings except for
	// `disallow_expression`.
	AllowExpression *v3.RegexMatcher `protobuf:"bytes,5,opt,name=allow_expression,json=allowExpression,proto3" json:"allow_expression,omitempty"`
	// If set, specifically disallow any header that matches this regular
	// expression regardless of any other settings.
	DisallowExpression *v3.RegexMatcher `protobuf:"bytes,6,opt,name=disallow_expression,json=disallowExpression,proto3" json:"disallow_expression,omitempty"`
	// If true, and if the rules in this list cause a header mutation to be
	// disallowed, then the filter using this configuration will terminate the
	// request with a 500 error. In addition, regardless of the setting of this
	// parameter, any attempt to set, add, or modify a disallowed header will
	// cause the `rejected_header_mutations` counter to be incremented.
	// Default is false.
	DisallowIsError *wrapperspb.BoolValue `protobuf:"bytes,7,opt,name=disallow_is_error,json=disallowIsError,proto3" json:"disallow_is_error,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *HeaderMutationRules) Reset() {
	*x = HeaderMutationRules{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeaderMutationRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeaderMutationRules) ProtoMessage() {}

func (x *HeaderMutationRules) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeaderMutationRules.ProtoReflect.Descriptor instead.
func (*HeaderMutationRules) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescGZIP(), []int{0}
}

func (x *HeaderMutationRules) GetAllowAllRouting() *wrapperspb.BoolValue {
	if x != nil {
		return x.AllowAllRouting
	}
	return nil
}

func (x *HeaderMutationRules) GetAllowEnvoy() *wrapperspb.BoolValue {
	if x != nil {
		return x.AllowEnvoy
	}
	return nil
}

func (x *HeaderMutationRules) GetDisallowSystem() *wrapperspb.BoolValue {
	if x != nil {
		return x.DisallowSystem
	}
	return nil
}

func (x *HeaderMutationRules) GetDisallowAll() *wrapperspb.BoolValue {
	if x != nil {
		return x.DisallowAll
	}
	return nil
}

func (x *HeaderMutationRules) GetAllowExpression() *v3.RegexMatcher {
	if x != nil {
		return x.AllowExpression
	}
	return nil
}

func (x *HeaderMutationRules) GetDisallowExpression() *v3.RegexMatcher {
	if x != nil {
		return x.DisallowExpression
	}
	return nil
}

func (x *HeaderMutationRules) GetDisallowIsError() *wrapperspb.BoolValue {
	if x != nil {
		return x.DisallowIsError
	}
	return nil
}

// The HeaderMutation structure specifies an action that may be taken on HTTP
// headers.
type HeaderMutation struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Action:
	//
	//	*HeaderMutation_Remove
	//	*HeaderMutation_Append
	Action        isHeaderMutation_Action `protobuf_oneof:"action"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HeaderMutation) Reset() {
	*x = HeaderMutation{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeaderMutation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeaderMutation) ProtoMessage() {}

func (x *HeaderMutation) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeaderMutation.ProtoReflect.Descriptor instead.
func (*HeaderMutation) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescGZIP(), []int{1}
}

func (x *HeaderMutation) GetAction() isHeaderMutation_Action {
	if x != nil {
		return x.Action
	}
	return nil
}

func (x *HeaderMutation) GetRemove() string {
	if x != nil {
		if x, ok := x.Action.(*HeaderMutation_Remove); ok {
			return x.Remove
		}
	}
	return ""
}

func (x *HeaderMutation) GetAppend() *v31.HeaderValueOption {
	if x != nil {
		if x, ok := x.Action.(*HeaderMutation_Append); ok {
			return x.Append
		}
	}
	return nil
}

type isHeaderMutation_Action interface {
	isHeaderMutation_Action()
}

type HeaderMutation_Remove struct {
	// Remove the specified header if it exists.
	Remove string `protobuf:"bytes,1,opt,name=remove,proto3,oneof"`
}

type HeaderMutation_Append struct {
	// Append new header by the specified HeaderValueOption.
	Append *v31.HeaderValueOption `protobuf:"bytes,2,opt,name=append,proto3,oneof"`
}

func (*HeaderMutation_Remove) isHeaderMutation_Action() {}

func (*HeaderMutation_Append) isHeaderMutation_Action() {}

var File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDesc = "" +
	"\n" +
	"mgithub.com/solo-io/gloo/projects/gloo/api/external/envoy/config/common/mutation_rules/v3/mutation_rules.proto\x12-solo.io.envoy.config.common.mutation_rules.v3\x1aRgithub.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/base.proto\x1aTgithub.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/regex.proto\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x1dudpa/annotations/status.proto\x1a\x17validate/validate.proto\x1a\x12extproto/ext.proto\"\x9c\x04\n" +
	"\x13HeaderMutationRules\x12F\n" +
	"\x11allow_all_routing\x18\x01 \x01(\v2\x1a.google.protobuf.BoolValueR\x0fallowAllRouting\x12;\n" +
	"\vallow_envoy\x18\x02 \x01(\v2\x1a.google.protobuf.BoolValueR\n" +
	"allowEnvoy\x12C\n" +
	"\x0fdisallow_system\x18\x03 \x01(\v2\x1a.google.protobuf.BoolValueR\x0edisallowSystem\x12=\n" +
	"\fdisallow_all\x18\x04 \x01(\v2\x1a.google.protobuf.BoolValueR\vdisallowAll\x12V\n" +
	"\x10allow_expression\x18\x05 \x01(\v2+.solo.io.envoy.type.matcher.v3.RegexMatcherR\x0fallowExpression\x12\\\n" +
	"\x13disallow_expression\x18\x06 \x01(\v2+.solo.io.envoy.type.matcher.v3.RegexMatcherR\x12disallowExpression\x12F\n" +
	"\x11disallow_is_error\x18\a \x01(\v2\x1a.google.protobuf.BoolValueR\x0fdisallowIsError\"\x91\x01\n" +
	"\x0eHeaderMutation\x12%\n" +
	"\x06remove\x18\x01 \x01(\tB\v\xfaB\br\x06\xc8\x01\x00\xc0\x01\x02H\x00R\x06remove\x12I\n" +
	"\x06append\x18\x02 \x01(\v2/.solo.io.envoy.config.core.v3.HeaderValueOptionH\x00R\x06appendB\r\n" +
	"\x06action\x12\x03\xf8B\x01B\xbd\x01\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01\xe2\xb5\xdf\xcb\a\x02\x10\x02\n" +
	"3io.envoyproxy.envoy.config.common.mutation_rules.v3B\x12MutationRulesProtoP\x01Z\\github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/common/mutation_rules/v3b\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_goTypes = []any{
	(*HeaderMutationRules)(nil),   // 0: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules
	(*HeaderMutation)(nil),        // 1: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutation
	(*wrapperspb.BoolValue)(nil),  // 2: google.protobuf.BoolValue
	(*v3.RegexMatcher)(nil),       // 3: solo.io.envoy.type.matcher.v3.RegexMatcher
	(*v31.HeaderValueOption)(nil), // 4: solo.io.envoy.config.core.v3.HeaderValueOption
}
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_depIdxs = []int32{
	2, // 0: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules.allow_all_routing:type_name -> google.protobuf.BoolValue
	2, // 1: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules.allow_envoy:type_name -> google.protobuf.BoolValue
	2, // 2: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules.disallow_system:type_name -> google.protobuf.BoolValue
	2, // 3: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules.disallow_all:type_name -> google.protobuf.BoolValue
	3, // 4: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules.allow_expression:type_name -> solo.io.envoy.type.matcher.v3.RegexMatcher
	3, // 5: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules.disallow_expression:type_name -> solo.io.envoy.type.matcher.v3.RegexMatcher
	2, // 6: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutationRules.disallow_is_error:type_name -> google.protobuf.BoolValue
	4, // 7: solo.io.envoy.config.common.mutation_rules.v3.HeaderMutation.append:type_name -> solo.io.envoy.config.core.v3.HeaderValueOption
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto != nil {
		return
	}
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_msgTypes[1].OneofWrappers = []any{
		(*HeaderMutation_Remove)(nil),
		(*HeaderMutation_Append)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_common_mutation_rules_v3_mutation_rules_proto_depIdxs = nil
}
