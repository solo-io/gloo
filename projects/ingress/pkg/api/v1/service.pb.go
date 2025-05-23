// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/ingress/api/v1/service.proto

package v1

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/solo-io/protoc-gen-ext/extproto"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A simple wrapper for a Kubernetes Service Object.
type KubeService struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// a raw byte representation of the kubernetes service this resource wraps
	KubeServiceSpec *anypb.Any `protobuf:"bytes,1,opt,name=kube_service_spec,json=kubeServiceSpec,proto3" json:"kube_service_spec,omitempty"`
	// a raw byte representation of the service status of the kubernetes service object
	KubeServiceStatus *anypb.Any `protobuf:"bytes,2,opt,name=kube_service_status,json=kubeServiceStatus,proto3" json:"kube_service_status,omitempty"`
	// Metadata contains the object metadata for this resource
	Metadata      *core.Metadata `protobuf:"bytes,7,opt,name=metadata,proto3" json:"metadata,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KubeService) Reset() {
	*x = KubeService{}
	mi := &file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KubeService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KubeService) ProtoMessage() {}

func (x *KubeService) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KubeService.ProtoReflect.Descriptor instead.
func (*KubeService) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *KubeService) GetKubeServiceSpec() *anypb.Any {
	if x != nil {
		return x.KubeServiceSpec
	}
	return nil
}

func (x *KubeService) GetKubeServiceStatus() *anypb.Any {
	if x != nil {
		return x.KubeServiceStatus
	}
	return nil
}

func (x *KubeService) GetMetadata() *core.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDesc = "" +
	"\n" +
	"=github.com/solo-io/gloo/projects/ingress/api/v1/service.proto\x12\x0fingress.solo.io\x1a\x19google/protobuf/any.proto\x1a1github.com/solo-io/solo-kit/api/v1/metadata.proto\x1a1github.com/solo-io/solo-kit/api/v1/solo-kit.proto\x1a\x12extproto/ext.proto\"\xdd\x01\n" +
	"\vKubeService\x12@\n" +
	"\x11kube_service_spec\x18\x01 \x01(\v2\x14.google.protobuf.AnyR\x0fkubeServiceSpec\x12D\n" +
	"\x13kube_service_status\x18\x02 \x01(\v2\x14.google.protobuf.AnyR\x11kubeServiceStatus\x122\n" +
	"\bmetadata\x18\a \x01(\v2\x16.core.solo.io.MetadataR\bmetadata:\x12\x82\xf1\x04\x0e\n" +
	"\x02sv\x12\bservicesBA\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01Z3github.com/solo-io/gloo/projects/ingress/pkg/api/v1b\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_goTypes = []any{
	(*KubeService)(nil),   // 0: ingress.solo.io.KubeService
	(*anypb.Any)(nil),     // 1: google.protobuf.Any
	(*core.Metadata)(nil), // 2: core.solo.io.Metadata
}
var file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_depIdxs = []int32{
	1, // 0: ingress.solo.io.KubeService.kube_service_spec:type_name -> google.protobuf.Any
	1, // 1: ingress.solo.io.KubeService.kube_service_status:type_name -> google.protobuf.Any
	2, // 2: ingress.solo.io.KubeService.metadata:type_name -> core.solo.io.Metadata
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_init() }
func file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_init() {
	if File_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto = out.File
	file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_ingress_api_v1_service_proto_depIdxs = nil
}
