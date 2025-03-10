// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/type/v3/http.proto

package v3

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/solo-io/gloo/projects/gloo/pkg/api/external/udpa/annotations"
	_ "github.com/solo-io/protoc-gen-ext/extproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CodecClientType int32

const (
	CodecClientType_HTTP1 CodecClientType = 0
	CodecClientType_HTTP2 CodecClientType = 1
	// [#not-implemented-hide:] QUIC implementation is not production ready yet. Use this enum with
	// caution to prevent accidental execution of QUIC code. I.e. `!= HTTP2` is no longer sufficient
	// to distinguish HTTP1 and HTTP2 traffic.
	CodecClientType_HTTP3 CodecClientType = 2
)

// Enum value maps for CodecClientType.
var (
	CodecClientType_name = map[int32]string{
		0: "HTTP1",
		1: "HTTP2",
		2: "HTTP3",
	}
	CodecClientType_value = map[string]int32{
		"HTTP1": 0,
		"HTTP2": 1,
		"HTTP3": 2,
	}
)

func (x CodecClientType) Enum() *CodecClientType {
	p := new(CodecClientType)
	*p = x
	return p
}

func (x CodecClientType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CodecClientType) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_enumTypes[0].Descriptor()
}

func (CodecClientType) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_enumTypes[0]
}

func (x CodecClientType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CodecClientType.Descriptor instead.
func (CodecClientType) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDescGZIP(), []int{0}
}

var File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDesc = string([]byte{
	0x0a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f,
	0x76, 0x33, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x2e, 0x76, 0x33, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x65, 0x78, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x32, 0x0a, 0x0f, 0x43, 0x6f, 0x64, 0x65, 0x63,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x54,
	0x54, 0x50, 0x31, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x54, 0x54, 0x50, 0x32, 0x10, 0x01,
	0x12, 0x09, 0x0a, 0x05, 0x48, 0x54, 0x54, 0x50, 0x33, 0x10, 0x02, 0x42, 0x8c, 0x01, 0xb8, 0xf5,
	0x04, 0x01, 0xc0, 0xf5, 0x04, 0x01, 0xd0, 0xf5, 0x04, 0x01, 0xe2, 0xb5, 0xdf, 0xcb, 0x07, 0x02,
	0x10, 0x02, 0x0a, 0x23, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78,
	0x79, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x33, 0x42, 0x09, 0x48, 0x74, 0x74, 0x70, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_goTypes = []any{
	(CodecClientType)(0), // 0: solo.io.envoy.type.v3.CodecClientType
}
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_init() }
func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_depIdxs,
		EnumInfos:         file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_enumTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_type_v3_http_proto_depIdxs = nil
}
