// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/tap/output_sink/v3/http_output_sink.proto

package v3

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// HTTP output sink definition
type HttpOutputSink struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// URI of the HTTP server to which output traces should be submitted
	ServerUri     *v3.HttpUri `protobuf:"bytes,1,opt,name=server_uri,json=serverUri,proto3" json:"server_uri,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HttpOutputSink) Reset() {
	*x = HttpOutputSink{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HttpOutputSink) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpOutputSink) ProtoMessage() {}

func (x *HttpOutputSink) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpOutputSink.ProtoReflect.Descriptor instead.
func (*HttpOutputSink) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDescGZIP(), []int{0}
}

func (x *HttpOutputSink) GetServerUri() *v3.HttpUri {
	if x != nil {
		return x.ServerUri
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDesc = string([]byte{
	0x0a, 0x69, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x74, 0x61, 0x70, 0x2f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x5f, 0x73, 0x69, 0x6e,
	0x6b, 0x2f, 0x76, 0x33, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x5f, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1f, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x74, 0x61, 0x70, 0x2e, 0x6f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x5f, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x33, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x33, 0x2f, 0x68,
	0x74, 0x74, 0x70, 0x5f, 0x75, 0x72, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x60, 0x0a,
	0x0e, 0x48, 0x74, 0x74, 0x70, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x53, 0x69, 0x6e, 0x6b, 0x12,
	0x4e, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x65, 0x6e,
	0x76, 0x6f, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x33, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x55, 0x72, 0x69, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a,
	0x01, 0x02, 0x10, 0x01, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x55, 0x72, 0x69, 0x42,
	0x8f, 0x01, 0x0a, 0x2d, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78,
	0x79, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6f,
	0x75, 0x74, 0x70, 0x75, 0x74, 0x5f, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x74, 0x61, 0x70, 0x2e, 0x76,
	0x33, 0x42, 0x13, 0x48, 0x74, 0x74, 0x70, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x53, 0x69, 0x6e,
	0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x74,
	0x61, 0x70, 0x2f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x5f, 0x73, 0x69, 0x6e, 0x6b, 0x2f, 0x76,
	0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_goTypes = []any{
	(*HttpOutputSink)(nil), // 0: envoy.config.tap.output_sink.v3.HttpOutputSink
	(*v3.HttpUri)(nil),     // 1: solo.io.envoy.config.core.v3.HttpUri
}
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_depIdxs = []int32{
	1, // 0: envoy.config.tap.output_sink.v3.HttpOutputSink.server_uri:type_name -> solo.io.envoy.config.core.v3.HttpUri
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_tap_output_sink_v3_http_output_sink_proto_depIdxs = nil
}
