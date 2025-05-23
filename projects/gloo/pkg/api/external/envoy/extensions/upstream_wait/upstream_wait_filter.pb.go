// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/upstream_wait/upstream_wait_filter.proto

package upstream_wait

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpstreamWaitFilterConfig struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpstreamWaitFilterConfig) Reset() {
	*x = UpstreamWaitFilterConfig{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpstreamWaitFilterConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpstreamWaitFilterConfig) ProtoMessage() {}

func (x *UpstreamWaitFilterConfig) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpstreamWaitFilterConfig.ProtoReflect.Descriptor instead.
func (*UpstreamWaitFilterConfig) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDescGZIP(), []int{0}
}

var File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDesc = "" +
	"\n" +
	"lgithub.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/upstream_wait/upstream_wait_filter.proto\x12)envoy.config.filter.http.upstream_wait.v2\"\x1a\n" +
	"\x18UpstreamWaitFilterConfigB\xab\x01\n" +
	"7io.envoyproxy.envoy.config.filter.http.upstream_wait.v2B\x17UpstreamWaitFilterProtoP\x01ZUgithub.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/upstream_waitb\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_goTypes = []any{
	(*UpstreamWaitFilterConfig)(nil), // 0: envoy.config.filter.http.upstream_wait.v2.UpstreamWaitFilterConfig
}
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_extensions_upstream_wait_upstream_wait_filter_proto_depIdxs = nil
}
