// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation_ee/transformation.proto

package transformation_ee

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	route "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/route"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/route/v3"
	_type "github.com/solo-io/solo-kit/pkg/api/external/envoy/type"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FilterTransformations struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Specifies transformations based on the route matches. The first matched transformation will be
	// applied. If there are overlapped match conditions, please put the most specific match first.
	Transformations []*TransformationRule `protobuf:"bytes,1,rep,name=transformations,proto3" json:"transformations,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *FilterTransformations) Reset() {
	*x = FilterTransformations{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilterTransformations) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterTransformations) ProtoMessage() {}

func (x *FilterTransformations) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterTransformations.ProtoReflect.Descriptor instead.
func (*FilterTransformations) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{0}
}

func (x *FilterTransformations) GetTransformations() []*TransformationRule {
	if x != nil {
		return x.Transformations
	}
	return nil
}

type TransformationRule struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The route matching parameter. Only when the match is satisfied, the "requires" field will
	// apply.
	//
	// For example: following match will match all requests.
	//
	// .. code-block:: yaml
	//
	//	match:
	//	  prefix: /
	Match   *route.RouteMatch `protobuf:"bytes,1,opt,name=match,proto3" json:"match,omitempty"`
	MatchV3 *v3.RouteMatch    `protobuf:"bytes,3,opt,name=match_v3,json=matchV3,proto3" json:"match_v3,omitempty"`
	// transformation to perform
	RouteTransformations *RouteTransformations `protobuf:"bytes,2,opt,name=route_transformations,json=routeTransformations,proto3" json:"route_transformations,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *TransformationRule) Reset() {
	*x = TransformationRule{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransformationRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformationRule) ProtoMessage() {}

func (x *TransformationRule) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransformationRule.ProtoReflect.Descriptor instead.
func (*TransformationRule) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{1}
}

func (x *TransformationRule) GetMatch() *route.RouteMatch {
	if x != nil {
		return x.Match
	}
	return nil
}

func (x *TransformationRule) GetMatchV3() *v3.RouteMatch {
	if x != nil {
		return x.MatchV3
	}
	return nil
}

func (x *TransformationRule) GetRouteTransformations() *RouteTransformations {
	if x != nil {
		return x.RouteTransformations
	}
	return nil
}

type RouteTransformations struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	RequestTransformation *Transformation        `protobuf:"bytes,1,opt,name=request_transformation,json=requestTransformation,proto3" json:"request_transformation,omitempty"`
	// clear the route cache if the request transformation was applied
	ClearRouteCache        bool            `protobuf:"varint,3,opt,name=clear_route_cache,json=clearRouteCache,proto3" json:"clear_route_cache,omitempty"`
	ResponseTransformation *Transformation `protobuf:"bytes,2,opt,name=response_transformation,json=responseTransformation,proto3" json:"response_transformation,omitempty"`
	// Apply a transformation in the onStreamComplete callback
	// (for modifying headers and dynamic metadata for access logs)
	OnStreamCompletionTransformation *Transformation `protobuf:"bytes,4,opt,name=on_stream_completion_transformation,json=onStreamCompletionTransformation,proto3" json:"on_stream_completion_transformation,omitempty"`
	unknownFields                    protoimpl.UnknownFields
	sizeCache                        protoimpl.SizeCache
}

func (x *RouteTransformations) Reset() {
	*x = RouteTransformations{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RouteTransformations) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteTransformations) ProtoMessage() {}

func (x *RouteTransformations) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteTransformations.ProtoReflect.Descriptor instead.
func (*RouteTransformations) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{2}
}

func (x *RouteTransformations) GetRequestTransformation() *Transformation {
	if x != nil {
		return x.RequestTransformation
	}
	return nil
}

func (x *RouteTransformations) GetClearRouteCache() bool {
	if x != nil {
		return x.ClearRouteCache
	}
	return false
}

func (x *RouteTransformations) GetResponseTransformation() *Transformation {
	if x != nil {
		return x.ResponseTransformation
	}
	return nil
}

func (x *RouteTransformations) GetOnStreamCompletionTransformation() *Transformation {
	if x != nil {
		return x.OnStreamCompletionTransformation
	}
	return nil
}

type Transformation struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Template is in the transformed request language domain
	//
	// Types that are valid to be assigned to TransformationType:
	//
	//	*Transformation_DlpTransformation
	TransformationType isTransformation_TransformationType `protobuf_oneof:"transformation_type"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Transformation) Reset() {
	*x = Transformation{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transformation) ProtoMessage() {}

func (x *Transformation) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transformation.ProtoReflect.Descriptor instead.
func (*Transformation) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{3}
}

func (x *Transformation) GetTransformationType() isTransformation_TransformationType {
	if x != nil {
		return x.TransformationType
	}
	return nil
}

func (x *Transformation) GetDlpTransformation() *DlpTransformation {
	if x != nil {
		if x, ok := x.TransformationType.(*Transformation_DlpTransformation); ok {
			return x.DlpTransformation
		}
	}
	return nil
}

type isTransformation_TransformationType interface {
	isTransformation_TransformationType()
}

type Transformation_DlpTransformation struct {
	DlpTransformation *DlpTransformation `protobuf:"bytes,1,opt,name=dlp_transformation,json=dlpTransformation,proto3,oneof"`
}

func (*Transformation_DlpTransformation) isTransformation_TransformationType() {}

type DlpTransformation struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// list of actions to apply
	Actions []*Action `protobuf:"bytes,1,rep,name=actions,proto3" json:"actions,omitempty"`
	// If true, headers will be transformed. Should only be true for the
	// on_stream_complete_transformation route transformation type.
	EnableHeaderTransformation bool `protobuf:"varint,2,opt,name=enable_header_transformation,json=enableHeaderTransformation,proto3" json:"enable_header_transformation,omitempty"`
	// If true, dynamic metadata will be transformed. Should only be used for the
	// on_stream_complete_transformation route transformation type.
	EnableDynamicMetadataTransformation bool `protobuf:"varint,3,opt,name=enable_dynamic_metadata_transformation,json=enableDynamicMetadataTransformation,proto3" json:"enable_dynamic_metadata_transformation,omitempty"`
	unknownFields                       protoimpl.UnknownFields
	sizeCache                           protoimpl.SizeCache
}

func (x *DlpTransformation) Reset() {
	*x = DlpTransformation{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DlpTransformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DlpTransformation) ProtoMessage() {}

func (x *DlpTransformation) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DlpTransformation.ProtoReflect.Descriptor instead.
func (*DlpTransformation) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{4}
}

func (x *DlpTransformation) GetActions() []*Action {
	if x != nil {
		return x.Actions
	}
	return nil
}

func (x *DlpTransformation) GetEnableHeaderTransformation() bool {
	if x != nil {
		return x.EnableHeaderTransformation
	}
	return false
}

func (x *DlpTransformation) GetEnableDynamicMetadataTransformation() bool {
	if x != nil {
		return x.EnableDynamicMetadataTransformation
	}
	return false
}

type Action struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Identifier for this action.
	// Used mostly to help ID specific actions in logs.
	// If left null will default to unknown
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Deprecated in favor of DlpMatcher
	// List of regexes to apply to the response body to match data which should be masked
	// They will be applied iteratively in the order which they are specified
	//
	// Deprecated: Marked as deprecated in github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation_ee/transformation.proto.
	Regex []string `protobuf:"bytes,2,rep,name=regex,proto3" json:"regex,omitempty"`
	// Deprecated in favor of DlpMatcher
	// List of regexes to apply to the response body to match data which should be
	// masked. They will be applied iteratively in the order which they are
	// specified. If this field and `regex` are both provided, all the regexes will
	// be applied iteratively in the order provided, starting with the ones from `regex`
	//
	// Deprecated: Marked as deprecated in github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation_ee/transformation.proto.
	RegexActions []*RegexAction `protobuf:"bytes,6,rep,name=regex_actions,json=regexActions,proto3" json:"regex_actions,omitempty"`
	// If specified, this rule will not actually be applied, but only logged.
	Shadow bool `protobuf:"varint,3,opt,name=shadow,proto3" json:"shadow,omitempty"`
	// The percent of the string which should be masked.
	// If not set, defaults to 75%
	Percent *_type.Percent `protobuf:"bytes,4,opt,name=percent,proto3" json:"percent,omitempty"`
	// The character which should overwrite the masked data
	// If left empty, defaults to "X"
	MaskChar string `protobuf:"bytes,5,opt,name=mask_char,json=maskChar,proto3" json:"mask_char,omitempty"`
	// The matcher used to determine which values will be masked by this action.
	Matcher       *Action_DlpMatcher `protobuf:"bytes,7,opt,name=matcher,proto3" json:"matcher,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Action) Reset() {
	*x = Action{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{5}
}

func (x *Action) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Deprecated: Marked as deprecated in github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation_ee/transformation.proto.
func (x *Action) GetRegex() []string {
	if x != nil {
		return x.Regex
	}
	return nil
}

// Deprecated: Marked as deprecated in github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation_ee/transformation.proto.
func (x *Action) GetRegexActions() []*RegexAction {
	if x != nil {
		return x.RegexActions
	}
	return nil
}

func (x *Action) GetShadow() bool {
	if x != nil {
		return x.Shadow
	}
	return false
}

func (x *Action) GetPercent() *_type.Percent {
	if x != nil {
		return x.Percent
	}
	return nil
}

func (x *Action) GetMaskChar() string {
	if x != nil {
		return x.MaskChar
	}
	return ""
}

func (x *Action) GetMatcher() *Action_DlpMatcher {
	if x != nil {
		return x.Matcher
	}
	return nil
}

type RegexAction struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The regex to match for masking.
	Regex string `protobuf:"bytes,1,opt,name=regex,proto3" json:"regex,omitempty"`
	// If provided and not 0, only this specific subgroup of the regex will be masked.
	Subgroup      uint32 `protobuf:"varint,2,opt,name=subgroup,proto3" json:"subgroup,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegexAction) Reset() {
	*x = RegexAction{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegexAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegexAction) ProtoMessage() {}

func (x *RegexAction) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegexAction.ProtoReflect.Descriptor instead.
func (*RegexAction) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{6}
}

func (x *RegexAction) GetRegex() string {
	if x != nil {
		return x.Regex
	}
	return ""
}

func (x *RegexAction) GetSubgroup() uint32 {
	if x != nil {
		return x.Subgroup
	}
	return 0
}

// List of regexes to apply to the response body to match data which should be
// masked. They will be applied iteratively in the order which they are
// specified.
type Action_RegexMatcher struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RegexActions  []*RegexAction         `protobuf:"bytes,1,rep,name=regex_actions,json=regexActions,proto3" json:"regex_actions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Action_RegexMatcher) Reset() {
	*x = Action_RegexMatcher{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Action_RegexMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action_RegexMatcher) ProtoMessage() {}

func (x *Action_RegexMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action_RegexMatcher.ProtoReflect.Descriptor instead.
func (*Action_RegexMatcher) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{5, 0}
}

func (x *Action_RegexMatcher) GetRegexActions() []*RegexAction {
	if x != nil {
		return x.RegexActions
	}
	return nil
}

// List of headers for which associated values will be masked.
// Note that enable_header_transformation must be set for this to take effect.
// Note that if enable_dynamic_metadata_transformation is set, proto struct dynamic metadata
// (i.e., the values matching any JSON keys specified in `keys`; primarily for json-formatted WAF audit logs) will also be masked accordingly.
type Action_KeyValueMatcher struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Keys          []string               `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Action_KeyValueMatcher) Reset() {
	*x = Action_KeyValueMatcher{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Action_KeyValueMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action_KeyValueMatcher) ProtoMessage() {}

func (x *Action_KeyValueMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action_KeyValueMatcher.ProtoReflect.Descriptor instead.
func (*Action_KeyValueMatcher) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{5, 1}
}

func (x *Action_KeyValueMatcher) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type Action_DlpMatcher struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Matcher:
	//
	//	*Action_DlpMatcher_RegexMatcher
	//	*Action_DlpMatcher_KeyValueMatcher
	Matcher       isAction_DlpMatcher_Matcher `protobuf_oneof:"matcher"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Action_DlpMatcher) Reset() {
	*x = Action_DlpMatcher{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Action_DlpMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action_DlpMatcher) ProtoMessage() {}

func (x *Action_DlpMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action_DlpMatcher.ProtoReflect.Descriptor instead.
func (*Action_DlpMatcher) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP(), []int{5, 2}
}

func (x *Action_DlpMatcher) GetMatcher() isAction_DlpMatcher_Matcher {
	if x != nil {
		return x.Matcher
	}
	return nil
}

func (x *Action_DlpMatcher) GetRegexMatcher() *Action_RegexMatcher {
	if x != nil {
		if x, ok := x.Matcher.(*Action_DlpMatcher_RegexMatcher); ok {
			return x.RegexMatcher
		}
	}
	return nil
}

func (x *Action_DlpMatcher) GetKeyValueMatcher() *Action_KeyValueMatcher {
	if x != nil {
		if x, ok := x.Matcher.(*Action_DlpMatcher_KeyValueMatcher); ok {
			return x.KeyValueMatcher
		}
	}
	return nil
}

type isAction_DlpMatcher_Matcher interface {
	isAction_DlpMatcher_Matcher()
}

type Action_DlpMatcher_RegexMatcher struct {
	RegexMatcher *Action_RegexMatcher `protobuf:"bytes,1,opt,name=regex_matcher,json=regexMatcher,proto3,oneof"`
}

type Action_DlpMatcher_KeyValueMatcher struct {
	KeyValueMatcher *Action_KeyValueMatcher `protobuf:"bytes,2,opt,name=key_value_matcher,json=keyValueMatcher,proto3,oneof"`
}

func (*Action_DlpMatcher_RegexMatcher) isAction_DlpMatcher_Matcher() {}

func (*Action_DlpMatcher_KeyValueMatcher) isAction_DlpMatcher_Matcher() {}

var File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDesc = "" +
	"\n" +
	"jgithub.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation_ee/transformation.proto\x12-envoy.config.filter.http.transformation_ee.v2\x1a\x17validate/validate.proto\x1a\x1eenvoy/api/v2/route/route.proto\x1aAgithub.com/solo-io/solo-kit/api/external/envoy/type/percent.proto\x1a_github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/route/v3/route_components.proto\"\x84\x01\n" +
	"\x15FilterTransformations\x12k\n" +
	"\x0ftransformations\x18\x01 \x03(\v2A.envoy.config.filter.http.transformation_ee.v2.TransformationRuleR\x0ftransformations\"\x92\x02\n" +
	"\x12TransformationRule\x12<\n" +
	"\x05match\x18\x01 \x01(\v2&.solo.io.envoy.api.v2.route.RouteMatchR\x05match\x12D\n" +
	"\bmatch_v3\x18\x03 \x01(\v2).solo.io.envoy.config.route.v3.RouteMatchR\amatchV3\x12x\n" +
	"\x15route_transformations\x18\x02 \x01(\v2C.envoy.config.filter.http.transformation_ee.v2.RouteTransformationsR\x14routeTransformations\"\xbf\x03\n" +
	"\x14RouteTransformations\x12t\n" +
	"\x16request_transformation\x18\x01 \x01(\v2=.envoy.config.filter.http.transformation_ee.v2.TransformationR\x15requestTransformation\x12*\n" +
	"\x11clear_route_cache\x18\x03 \x01(\bR\x0fclearRouteCache\x12v\n" +
	"\x17response_transformation\x18\x02 \x01(\v2=.envoy.config.filter.http.transformation_ee.v2.TransformationR\x16responseTransformation\x12\x8c\x01\n" +
	"#on_stream_completion_transformation\x18\x04 \x01(\v2=.envoy.config.filter.http.transformation_ee.v2.TransformationR onStreamCompletionTransformation\"\x9a\x01\n" +
	"\x0eTransformation\x12q\n" +
	"\x12dlp_transformation\x18\x01 \x01(\v2@.envoy.config.filter.http.transformation_ee.v2.DlpTransformationH\x00R\x11dlpTransformationB\x15\n" +
	"\x13transformation_type\"\xfb\x01\n" +
	"\x11DlpTransformation\x12O\n" +
	"\aactions\x18\x01 \x03(\v25.envoy.config.filter.http.transformation_ee.v2.ActionR\aactions\x12@\n" +
	"\x1cenable_header_transformation\x18\x02 \x01(\bR\x1aenableHeaderTransformation\x12S\n" +
	"&enable_dynamic_metadata_transformation\x18\x03 \x01(\bR#enableDynamicMetadataTransformation\"\x8a\x06\n" +
	"\x06Action\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12$\n" +
	"\x05regex\x18\x02 \x03(\tB\x0e\xfaB\t\x92\x01\x06\"\x04r\x02 \x01\x18\x01R\x05regex\x12c\n" +
	"\rregex_actions\x18\x06 \x03(\v2:.envoy.config.filter.http.transformation_ee.v2.RegexActionB\x02\x18\x01R\fregexActions\x12\x16\n" +
	"\x06shadow\x18\x03 \x01(\bR\x06shadow\x125\n" +
	"\apercent\x18\x04 \x01(\v2\x1b.solo.io.envoy.type.PercentR\apercent\x12$\n" +
	"\tmask_char\x18\x05 \x01(\tB\a\xfaB\x04r\x02(\x01R\bmaskChar\x12Z\n" +
	"\amatcher\x18\a \x01(\v2@.envoy.config.filter.http.transformation_ee.v2.Action.DlpMatcherR\amatcher\x1ao\n" +
	"\fRegexMatcher\x12_\n" +
	"\rregex_actions\x18\x01 \x03(\v2:.envoy.config.filter.http.transformation_ee.v2.RegexActionR\fregexActions\x1a%\n" +
	"\x0fKeyValueMatcher\x12\x12\n" +
	"\x04keys\x18\x01 \x03(\tR\x04keys\x1a\xf7\x01\n" +
	"\n" +
	"DlpMatcher\x12i\n" +
	"\rregex_matcher\x18\x01 \x01(\v2B.envoy.config.filter.http.transformation_ee.v2.Action.RegexMatcherH\x00R\fregexMatcher\x12s\n" +
	"\x11key_value_matcher\x18\x02 \x01(\v2E.envoy.config.filter.http.transformation_ee.v2.Action.KeyValueMatcherH\x00R\x0fkeyValueMatcherB\t\n" +
	"\amatcher\"H\n" +
	"\vRegexAction\x12\x1d\n" +
	"\x05regex\x18\x01 \x01(\tB\a\xfaB\x04r\x02 \x01R\x05regex\x12\x1a\n" +
	"\bsubgroup\x18\x02 \x01(\rR\bsubgroupB\xb7\x01\n" +
	";io.envoyproxy.envoy.config.filter.http.transformation_ee.v2B\x1bTransformationEEFilterProtoP\x01ZYgithub.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_eeb\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_goTypes = []any{
	(*FilterTransformations)(nil),  // 0: envoy.config.filter.http.transformation_ee.v2.FilterTransformations
	(*TransformationRule)(nil),     // 1: envoy.config.filter.http.transformation_ee.v2.TransformationRule
	(*RouteTransformations)(nil),   // 2: envoy.config.filter.http.transformation_ee.v2.RouteTransformations
	(*Transformation)(nil),         // 3: envoy.config.filter.http.transformation_ee.v2.Transformation
	(*DlpTransformation)(nil),      // 4: envoy.config.filter.http.transformation_ee.v2.DlpTransformation
	(*Action)(nil),                 // 5: envoy.config.filter.http.transformation_ee.v2.Action
	(*RegexAction)(nil),            // 6: envoy.config.filter.http.transformation_ee.v2.RegexAction
	(*Action_RegexMatcher)(nil),    // 7: envoy.config.filter.http.transformation_ee.v2.Action.RegexMatcher
	(*Action_KeyValueMatcher)(nil), // 8: envoy.config.filter.http.transformation_ee.v2.Action.KeyValueMatcher
	(*Action_DlpMatcher)(nil),      // 9: envoy.config.filter.http.transformation_ee.v2.Action.DlpMatcher
	(*route.RouteMatch)(nil),       // 10: solo.io.envoy.api.v2.route.RouteMatch
	(*v3.RouteMatch)(nil),          // 11: solo.io.envoy.config.route.v3.RouteMatch
	(*_type.Percent)(nil),          // 12: solo.io.envoy.type.Percent
}
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_depIdxs = []int32{
	1,  // 0: envoy.config.filter.http.transformation_ee.v2.FilterTransformations.transformations:type_name -> envoy.config.filter.http.transformation_ee.v2.TransformationRule
	10, // 1: envoy.config.filter.http.transformation_ee.v2.TransformationRule.match:type_name -> solo.io.envoy.api.v2.route.RouteMatch
	11, // 2: envoy.config.filter.http.transformation_ee.v2.TransformationRule.match_v3:type_name -> solo.io.envoy.config.route.v3.RouteMatch
	2,  // 3: envoy.config.filter.http.transformation_ee.v2.TransformationRule.route_transformations:type_name -> envoy.config.filter.http.transformation_ee.v2.RouteTransformations
	3,  // 4: envoy.config.filter.http.transformation_ee.v2.RouteTransformations.request_transformation:type_name -> envoy.config.filter.http.transformation_ee.v2.Transformation
	3,  // 5: envoy.config.filter.http.transformation_ee.v2.RouteTransformations.response_transformation:type_name -> envoy.config.filter.http.transformation_ee.v2.Transformation
	3,  // 6: envoy.config.filter.http.transformation_ee.v2.RouteTransformations.on_stream_completion_transformation:type_name -> envoy.config.filter.http.transformation_ee.v2.Transformation
	4,  // 7: envoy.config.filter.http.transformation_ee.v2.Transformation.dlp_transformation:type_name -> envoy.config.filter.http.transformation_ee.v2.DlpTransformation
	5,  // 8: envoy.config.filter.http.transformation_ee.v2.DlpTransformation.actions:type_name -> envoy.config.filter.http.transformation_ee.v2.Action
	6,  // 9: envoy.config.filter.http.transformation_ee.v2.Action.regex_actions:type_name -> envoy.config.filter.http.transformation_ee.v2.RegexAction
	12, // 10: envoy.config.filter.http.transformation_ee.v2.Action.percent:type_name -> solo.io.envoy.type.Percent
	9,  // 11: envoy.config.filter.http.transformation_ee.v2.Action.matcher:type_name -> envoy.config.filter.http.transformation_ee.v2.Action.DlpMatcher
	6,  // 12: envoy.config.filter.http.transformation_ee.v2.Action.RegexMatcher.regex_actions:type_name -> envoy.config.filter.http.transformation_ee.v2.RegexAction
	7,  // 13: envoy.config.filter.http.transformation_ee.v2.Action.DlpMatcher.regex_matcher:type_name -> envoy.config.filter.http.transformation_ee.v2.Action.RegexMatcher
	8,  // 14: envoy.config.filter.http.transformation_ee.v2.Action.DlpMatcher.key_value_matcher:type_name -> envoy.config.filter.http.transformation_ee.v2.Action.KeyValueMatcher
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	15, // [15:15] is the sub-list for extension type_name
	15, // [15:15] is the sub-list for extension extendee
	0,  // [0:15] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto != nil {
		return
	}
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[3].OneofWrappers = []any{
		(*Transformation_DlpTransformation)(nil),
	}
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes[9].OneofWrappers = []any{
		(*Action_DlpMatcher_RegexMatcher)(nil),
		(*Action_DlpMatcher_KeyValueMatcher)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_transformation_ee_transformation_proto_depIdxs = nil
}
