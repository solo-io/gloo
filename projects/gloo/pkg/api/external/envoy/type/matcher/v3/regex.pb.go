// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/regex.proto

package v3

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

// A regex matcher designed for safety when used with untrusted input.
type RegexMatcher struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to EngineType:
	//
	//	*RegexMatcher_GoogleRe2
	EngineType isRegexMatcher_EngineType `protobuf_oneof:"engine_type"`
	// The regex match string. The string must be supported by the configured engine.
	Regex         string `protobuf:"bytes,2,opt,name=regex,proto3" json:"regex,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegexMatcher) Reset() {
	*x = RegexMatcher{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegexMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegexMatcher) ProtoMessage() {}

func (x *RegexMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegexMatcher.ProtoReflect.Descriptor instead.
func (*RegexMatcher) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescGZIP(), []int{0}
}

func (x *RegexMatcher) GetEngineType() isRegexMatcher_EngineType {
	if x != nil {
		return x.EngineType
	}
	return nil
}

func (x *RegexMatcher) GetGoogleRe2() *RegexMatcher_GoogleRE2 {
	if x != nil {
		if x, ok := x.EngineType.(*RegexMatcher_GoogleRe2); ok {
			return x.GoogleRe2
		}
	}
	return nil
}

func (x *RegexMatcher) GetRegex() string {
	if x != nil {
		return x.Regex
	}
	return ""
}

type isRegexMatcher_EngineType interface {
	isRegexMatcher_EngineType()
}

type RegexMatcher_GoogleRe2 struct {
	// Google's RE2 regex engine.
	GoogleRe2 *RegexMatcher_GoogleRE2 `protobuf:"bytes,1,opt,name=google_re2,json=googleRe2,proto3,oneof"`
}

func (*RegexMatcher_GoogleRe2) isRegexMatcher_EngineType() {}

// Describes how to match a string and then produce a new string using a regular
// expression and a substitution string.
type RegexMatchAndSubstitute struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The regular expression used to find portions of a string (hereafter called
	// the "subject string") that should be replaced. When a new string is
	// produced during the substitution operation, the new string is initially
	// the same as the subject string, but then all matches in the subject string
	// are replaced by the substitution string. If replacing all matches isn't
	// desired, regular expression anchors can be used to ensure a single match,
	// so as to replace just one occurrence of a pattern. Capture groups can be
	// used in the pattern to extract portions of the subject string, and then
	// referenced in the substitution string.
	Pattern *RegexMatcher `protobuf:"bytes,1,opt,name=pattern,proto3" json:"pattern,omitempty"`
	// The string that should be substituted into matching portions of the
	// subject string during a substitution operation to produce a new string.
	// Capture groups in the pattern can be referenced in the substitution
	// string. Note, however, that the syntax for referring to capture groups is
	// defined by the chosen regular expression engine. Google's [RE2](https://github.com/google/re2) regular expression engine uses a
	// backslash followed by the capture group number to denote a numbered
	// capture group. E.g., `\1` refers to capture group 1, and `\2` refers
	// to capture group 2.
	Substitution  string `protobuf:"bytes,2,opt,name=substitution,proto3" json:"substitution,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegexMatchAndSubstitute) Reset() {
	*x = RegexMatchAndSubstitute{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegexMatchAndSubstitute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegexMatchAndSubstitute) ProtoMessage() {}

func (x *RegexMatchAndSubstitute) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegexMatchAndSubstitute.ProtoReflect.Descriptor instead.
func (*RegexMatchAndSubstitute) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescGZIP(), []int{1}
}

func (x *RegexMatchAndSubstitute) GetPattern() *RegexMatcher {
	if x != nil {
		return x.Pattern
	}
	return nil
}

func (x *RegexMatchAndSubstitute) GetSubstitution() string {
	if x != nil {
		return x.Substitution
	}
	return ""
}

// Google's [RE2](https://github.com/google/re2) regex engine. The regex string must adhere to
// the documented [syntax](https://github.com/google/re2/wiki/Syntax). The engine is designed
// to complete execution in linear time as well as limit the amount of memory used.
//
// Envoy supports program size checking via runtime. The runtime keys `re2.max_program_size.error_level`
// and `re2.max_program_size.warn_level` can be set to integers as the maximum program size or
// complexity that a compiled regex can have before an exception is thrown or a warning is
// logged, respectively. `re2.max_program_size.error_level` defaults to 100, and
// `re2.max_program_size.warn_level` has no default if unset (will not check/log a warning).
//
// Envoy emits two stats for tracking the program size of regexes: the histogram `re2.program_size`,
// which records the program size, and the counter `re2.exceeded_warn_level`, which is incremented
// each time the program size exceeds the warn level threshold.
type RegexMatcher_GoogleRE2 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// This field controls the RE2 "program size" which is a rough estimate of how complex a
	// compiled regex is to evaluate. A regex that has a program size greater than the configured
	// value will fail to compile. In this case, the configured max program size can be increased
	// or the regex can be simplified. If not specified, the default is 100.
	//
	// This field is deprecated; regexp validation should be performed on the management server
	// instead of being done by each individual client.
	//
	// Deprecated: Marked as deprecated in github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/regex.proto.
	MaxProgramSize *wrapperspb.UInt32Value `protobuf:"bytes,1,opt,name=max_program_size,json=maxProgramSize,proto3" json:"max_program_size,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *RegexMatcher_GoogleRE2) Reset() {
	*x = RegexMatcher_GoogleRE2{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegexMatcher_GoogleRE2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegexMatcher_GoogleRE2) ProtoMessage() {}

func (x *RegexMatcher_GoogleRE2) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegexMatcher_GoogleRE2.ProtoReflect.Descriptor instead.
func (*RegexMatcher_GoogleRE2) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescGZIP(), []int{0, 0}
}

// Deprecated: Marked as deprecated in github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/matcher/v3/regex.proto.
func (x *RegexMatcher_GoogleRE2) GetMaxProgramSize() *wrapperspb.UInt32Value {
	if x != nil {
		return x.MaxProgramSize
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDesc = string([]byte{
	0x0a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x76, 0x33, 0x2f, 0x72, 0x65, 0x67, 0x65, 0x78,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e,
	0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x65, 0x72, 0x2e, 0x76, 0x33, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x69, 0x6e,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x12, 0x65, 0x78, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd9, 0x02, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x65, 0x78, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x60, 0x0a, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x5f,
	0x72, 0x65, 0x32, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x76, 0x33, 0x2e, 0x52, 0x65, 0x67, 0x65, 0x78, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x52, 0x45, 0x32,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x48, 0x00, 0x52, 0x09, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x52, 0x65, 0x32, 0x12, 0x1d, 0x0a, 0x05, 0x72, 0x65, 0x67, 0x65, 0x78,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x20, 0x01, 0x52,
	0x05, 0x72, 0x65, 0x67, 0x65, 0x78, 0x1a, 0x8a, 0x01, 0x0a, 0x09, 0x47, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x52, 0x45, 0x32, 0x12, 0x4a, 0x0a, 0x10, 0x6d, 0x61, 0x78, 0x5f, 0x70, 0x72, 0x6f, 0x67,
	0x72, 0x61, 0x6d, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x02, 0x18, 0x01,
	0x52, 0x0e, 0x6d, 0x61, 0x78, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x53, 0x69, 0x7a, 0x65,
	0x3a, 0x31, 0x8a, 0xc8, 0xde, 0x8e, 0x04, 0x2b, 0x0a, 0x29, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x67,
	0x65, 0x78, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x52, 0x45, 0x32, 0x3a, 0x27, 0x8a, 0xc8, 0xde, 0x8e, 0x04, 0x21, 0x0a, 0x1f, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e,
	0x52, 0x65, 0x67, 0x65, 0x78, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x42, 0x12, 0x0a, 0x0b,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12, 0x03, 0xf8, 0x42, 0x01,
	0x22, 0xc2, 0x01, 0x0a, 0x17, 0x52, 0x65, 0x67, 0x65, 0x78, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x41,
	0x6e, 0x64, 0x53, 0x75, 0x62, 0x73, 0x74, 0x69, 0x74, 0x75, 0x74, 0x65, 0x12, 0x4f, 0x0a, 0x07,
	0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x76, 0x33, 0x2e, 0x52, 0x65,
	0x67, 0x65, 0x78, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a,
	0x01, 0x02, 0x10, 0x01, 0x52, 0x07, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x12, 0x22, 0x0a,
	0x0c, 0x73, 0x75, 0x62, 0x73, 0x74, 0x69, 0x74, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x75, 0x62, 0x73, 0x74, 0x69, 0x74, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x3a, 0x32, 0x8a, 0xc8, 0xde, 0x8e, 0x04, 0x2c, 0x0a, 0x2a, 0x65, 0x6e, 0x76, 0x6f, 0x79,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x52, 0x65,
	0x67, 0x65, 0x78, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x41, 0x6e, 0x64, 0x53, 0x75, 0x62, 0x73, 0x74,
	0x69, 0x74, 0x75, 0x74, 0x65, 0x42, 0x9d, 0x01, 0xb8, 0xf5, 0x04, 0x01, 0xc0, 0xf5, 0x04, 0x01,
	0xd0, 0xf5, 0x04, 0x01, 0xe2, 0xb5, 0xdf, 0xcb, 0x07, 0x02, 0x10, 0x02, 0x0a, 0x2b, 0x69, 0x6f,
	0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x76, 0x33, 0x42, 0x0a, 0x52, 0x65, 0x67, 0x65, 0x78,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f,
	0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x65, 0x72, 0x2f, 0x76, 0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_goTypes = []any{
	(*RegexMatcher)(nil),            // 0: solo.io.envoy.type.matcher.v3.RegexMatcher
	(*RegexMatchAndSubstitute)(nil), // 1: solo.io.envoy.type.matcher.v3.RegexMatchAndSubstitute
	(*RegexMatcher_GoogleRE2)(nil),  // 2: solo.io.envoy.type.matcher.v3.RegexMatcher.GoogleRE2
	(*wrapperspb.UInt32Value)(nil),  // 3: google.protobuf.UInt32Value
}
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_depIdxs = []int32{
	2, // 0: solo.io.envoy.type.matcher.v3.RegexMatcher.google_re2:type_name -> solo.io.envoy.type.matcher.v3.RegexMatcher.GoogleRE2
	0, // 1: solo.io.envoy.type.matcher.v3.RegexMatchAndSubstitute.pattern:type_name -> solo.io.envoy.type.matcher.v3.RegexMatcher
	3, // 2: solo.io.envoy.type.matcher.v3.RegexMatcher.GoogleRE2.max_program_size:type_name -> google.protobuf.UInt32Value
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto != nil {
		return
	}
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes[0].OneofWrappers = []any{
		(*RegexMatcher_GoogleRe2)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_matcher_v3_regex_proto_depIdxs = nil
}
