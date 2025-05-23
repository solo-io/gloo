// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/transformation/parameters.proto

package transformation

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

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

type Parameters struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// headers that will be used to extract data for processing output templates
	// Gloo will search for parameters by their name in header value strings, enclosed in single
	// curly braces
	// Example:
	//
	//	extensions:
	//	  parameters:
	//	      headers:
	//	        x-user-id: '{userId}'
	Headers map[string]string `protobuf:"bytes,1,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// part of the (or the entire) path that will be used extract data for processing output templates
	// Gloo will search for parameters by their name in header value strings, enclosed in single
	// curly braces
	// Example:
	//
	//	extensions:
	//	  parameters:
	//	      path: /users/{ userId }
	Path          *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Parameters) Reset() {
	*x = Parameters{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Parameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Parameters) ProtoMessage() {}

func (x *Parameters) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Parameters.ProtoReflect.Descriptor instead.
func (*Parameters) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDescGZIP(), []int{0}
}

func (x *Parameters) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Parameters) GetPath() *wrapperspb.StringValue {
	if x != nil {
		return x.Path
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDesc = "" +
	"\n" +
	"Tgithub.com/solo-io/gloo/projects/gloo/api/v1/options/transformation/parameters.proto\x12#transformation.options.gloo.solo.io\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x12extproto/ext.proto\"\xd2\x01\n" +
	"\n" +
	"Parameters\x12V\n" +
	"\aheaders\x18\x01 \x03(\v2<.transformation.options.gloo.solo.io.Parameters.HeadersEntryR\aheaders\x120\n" +
	"\x04path\x18\x02 \x01(\v2\x1c.google.protobuf.StringValueR\x04path\x1a:\n" +
	"\fHeadersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01BU\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01ZGgithub.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformationb\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_goTypes = []any{
	(*Parameters)(nil),             // 0: transformation.options.gloo.solo.io.Parameters
	nil,                            // 1: transformation.options.gloo.solo.io.Parameters.HeadersEntry
	(*wrapperspb.StringValue)(nil), // 2: google.protobuf.StringValue
}
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_depIdxs = []int32{
	1, // 0: transformation.options.gloo.solo.io.Parameters.headers:type_name -> transformation.options.gloo.solo.io.Parameters.HeadersEntry
	2, // 1: transformation.options.gloo.solo.io.Parameters.path:type_name -> google.protobuf.StringValue
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_transformation_parameters_proto_depIdxs = nil
}
