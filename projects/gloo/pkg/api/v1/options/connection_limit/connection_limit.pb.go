// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/connection_limit/connection_limit.proto

package connection_limit

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/solo-io/protoc-gen-ext/extproto"
	_ "github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// These options provide the ability to limit the active connections in envoy.
// Ref. https://www.envoyproxy.io/docs/envoy/latest/configuration/listeners/network_filters/connection_limit_filter
type ConnectionLimit struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The maximum number of active connections for this gateway. When this limit is reached, any incoming connection
	// will be closed after delay duration.
	// Must be greater than or equal to one.
	MaxActiveConnections *wrapperspb.UInt32Value `protobuf:"bytes,1,opt,name=max_active_connections,json=maxActiveConnections,proto3" json:"max_active_connections,omitempty"`
	// The time to wait before a connection is dropped. Useful for DoS prevention.
	// Defaults to zero and the connection will be closed immediately.
	DelayBeforeClose *durationpb.Duration `protobuf:"bytes,2,opt,name=delay_before_close,json=delayBeforeClose,proto3" json:"delay_before_close,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *ConnectionLimit) Reset() {
	*x = ConnectionLimit{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConnectionLimit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionLimit) ProtoMessage() {}

func (x *ConnectionLimit) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionLimit.ProtoReflect.Descriptor instead.
func (*ConnectionLimit) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDescGZIP(), []int{0}
}

func (x *ConnectionLimit) GetMaxActiveConnections() *wrapperspb.UInt32Value {
	if x != nil {
		return x.MaxActiveConnections
	}
	return nil
}

func (x *ConnectionLimit) GetDelayBeforeClose() *durationpb.Duration {
	if x != nil {
		return x.DelayBeforeClose
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDesc = "" +
	"\n" +
	"\\github.com/solo-io/gloo/projects/gloo/api/v1/options/connection_limit/connection_limit.proto\x12%connection_limit.options.gloo.solo.io\x1a\x12extproto/ext.proto\x1aEgithub.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base.proto\x1a\x1egoogle/protobuf/duration.proto\x1a\x1egoogle/protobuf/wrappers.proto\"\xae\x01\n" +
	"\x0fConnectionLimit\x12R\n" +
	"\x16max_active_connections\x18\x01 \x01(\v2\x1c.google.protobuf.UInt32ValueR\x14maxActiveConnections\x12G\n" +
	"\x12delay_before_close\x18\x02 \x01(\v2\x19.google.protobuf.DurationR\x10delayBeforeCloseBW\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01ZIgithub.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/connection_limitb\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_goTypes = []any{
	(*ConnectionLimit)(nil),        // 0: connection_limit.options.gloo.solo.io.ConnectionLimit
	(*wrapperspb.UInt32Value)(nil), // 1: google.protobuf.UInt32Value
	(*durationpb.Duration)(nil),    // 2: google.protobuf.Duration
}
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_depIdxs = []int32{
	1, // 0: connection_limit.options.gloo.solo.io.ConnectionLimit.max_active_connections:type_name -> google.protobuf.UInt32Value
	2, // 1: connection_limit.options.gloo.solo.io.ConnectionLimit.delay_before_close:type_name -> google.protobuf.Duration
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_connection_limit_connection_limit_proto_depIdxs = nil
}
