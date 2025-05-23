// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/address.proto

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

type SocketAddress_Protocol int32

const (
	SocketAddress_TCP SocketAddress_Protocol = 0
	SocketAddress_UDP SocketAddress_Protocol = 1
)

// Enum value maps for SocketAddress_Protocol.
var (
	SocketAddress_Protocol_name = map[int32]string{
		0: "TCP",
		1: "UDP",
	}
	SocketAddress_Protocol_value = map[string]int32{
		"TCP": 0,
		"UDP": 1,
	}
)

func (x SocketAddress_Protocol) Enum() *SocketAddress_Protocol {
	p := new(SocketAddress_Protocol)
	*p = x
	return p
}

func (x SocketAddress_Protocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SocketAddress_Protocol) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_enumTypes[0].Descriptor()
}

func (SocketAddress_Protocol) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_enumTypes[0]
}

func (x SocketAddress_Protocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SocketAddress_Protocol.Descriptor instead.
func (SocketAddress_Protocol) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP(), []int{1, 0}
}

type Pipe struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unix Domain Socket path. On Linux, paths starting with '@' will use the
	// abstract namespace. The starting '@' is replaced by a null byte by Envoy.
	// Paths starting with '@' will result in an error in environments other than
	// Linux.
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// The mode for the Pipe. Not applicable for abstract sockets.
	Mode          uint32 `protobuf:"varint,2,opt,name=mode,proto3" json:"mode,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Pipe) Reset() {
	*x = Pipe{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pipe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pipe) ProtoMessage() {}

func (x *Pipe) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pipe.ProtoReflect.Descriptor instead.
func (*Pipe) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP(), []int{0}
}

func (x *Pipe) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Pipe) GetMode() uint32 {
	if x != nil {
		return x.Mode
	}
	return 0
}

// [#next-free-field: 7]
type SocketAddress struct {
	state    protoimpl.MessageState `protogen:"open.v1"`
	Protocol SocketAddress_Protocol `protobuf:"varint,1,opt,name=protocol,proto3,enum=solo.io.envoy.config.core.v3.SocketAddress_Protocol" json:"protocol,omitempty"`
	// The address for this socket. Listeners will bind
	// to the address. An empty address is not allowed. Specify `0.0.0.0` or `::`
	// to bind to any address. [#comment:TODO(zuercher) reinstate when implemented:
	// It is possible to distinguish a Listener address via the prefix/suffix matching
	// in FilterChainMatch.] When used
	// within an upstream BindConfig, the address
	// controls the source address of outbound connections. For :ref:`clusters
	// <envoy_api_msg_config.cluster.v3.Cluster>`, the cluster type determines whether the
	// address must be an IP (*STATIC* or *EDS* clusters) or a hostname resolved by DNS
	// (*STRICT_DNS* or *LOGICAL_DNS* clusters). Address resolution can be customized
	// via resolver_name.
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// Types that are valid to be assigned to PortSpecifier:
	//
	//	*SocketAddress_PortValue
	//	*SocketAddress_NamedPort
	PortSpecifier isSocketAddress_PortSpecifier `protobuf_oneof:"port_specifier"`
	// The name of the custom resolver. This must have been registered with Envoy. If
	// this is empty, a context dependent default applies. If the address is a concrete
	// IP address, no resolution will occur. If address is a hostname this
	// should be set for resolution other than DNS. Specifying a custom resolver with
	// *STRICT_DNS* or *LOGICAL_DNS* will generate an error at runtime.
	ResolverName string `protobuf:"bytes,5,opt,name=resolver_name,json=resolverName,proto3" json:"resolver_name,omitempty"`
	// When binding to an IPv6 address above, this enables [IPv4 compatibility](https://datatracker.ietf.org/doc/html/rfc3493#page-11). Binding to `::` will
	// allow both IPv4 and IPv6 connections, with peer IPv4 addresses mapped into
	// IPv6 space as `::FFFF:<IPv4-address>`.
	Ipv4Compat    bool `protobuf:"varint,6,opt,name=ipv4_compat,json=ipv4Compat,proto3" json:"ipv4_compat,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SocketAddress) Reset() {
	*x = SocketAddress{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SocketAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SocketAddress) ProtoMessage() {}

func (x *SocketAddress) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SocketAddress.ProtoReflect.Descriptor instead.
func (*SocketAddress) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP(), []int{1}
}

func (x *SocketAddress) GetProtocol() SocketAddress_Protocol {
	if x != nil {
		return x.Protocol
	}
	return SocketAddress_TCP
}

func (x *SocketAddress) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *SocketAddress) GetPortSpecifier() isSocketAddress_PortSpecifier {
	if x != nil {
		return x.PortSpecifier
	}
	return nil
}

func (x *SocketAddress) GetPortValue() uint32 {
	if x != nil {
		if x, ok := x.PortSpecifier.(*SocketAddress_PortValue); ok {
			return x.PortValue
		}
	}
	return 0
}

func (x *SocketAddress) GetNamedPort() string {
	if x != nil {
		if x, ok := x.PortSpecifier.(*SocketAddress_NamedPort); ok {
			return x.NamedPort
		}
	}
	return ""
}

func (x *SocketAddress) GetResolverName() string {
	if x != nil {
		return x.ResolverName
	}
	return ""
}

func (x *SocketAddress) GetIpv4Compat() bool {
	if x != nil {
		return x.Ipv4Compat
	}
	return false
}

type isSocketAddress_PortSpecifier interface {
	isSocketAddress_PortSpecifier()
}

type SocketAddress_PortValue struct {
	PortValue uint32 `protobuf:"varint,3,opt,name=port_value,json=portValue,proto3,oneof"`
}

type SocketAddress_NamedPort struct {
	// This is only valid if :ref:`resolver_name
	// <envoy_api_field_config.core.v3.SocketAddress.resolver_name>` is specified below and the
	// named resolver is capable of named port resolution.
	NamedPort string `protobuf:"bytes,4,opt,name=named_port,json=namedPort,proto3,oneof"`
}

func (*SocketAddress_PortValue) isSocketAddress_PortSpecifier() {}

func (*SocketAddress_NamedPort) isSocketAddress_PortSpecifier() {}

type TcpKeepalive struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Maximum number of keepalive probes to send without response before deciding
	// the connection is dead. Default is to use the OS level configuration (unless
	// overridden, Linux defaults to 9.)
	KeepaliveProbes *wrapperspb.UInt32Value `protobuf:"bytes,1,opt,name=keepalive_probes,json=keepaliveProbes,proto3" json:"keepalive_probes,omitempty"`
	// The number of seconds a connection needs to be idle before keep-alive probes
	// start being sent. Default is to use the OS level configuration (unless
	// overridden, Linux defaults to 7200s (i.e., 2 hours.)
	KeepaliveTime *wrapperspb.UInt32Value `protobuf:"bytes,2,opt,name=keepalive_time,json=keepaliveTime,proto3" json:"keepalive_time,omitempty"`
	// The number of seconds between keep-alive probes. Default is to use the OS
	// level configuration (unless overridden, Linux defaults to 75s.)
	KeepaliveInterval *wrapperspb.UInt32Value `protobuf:"bytes,3,opt,name=keepalive_interval,json=keepaliveInterval,proto3" json:"keepalive_interval,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *TcpKeepalive) Reset() {
	*x = TcpKeepalive{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TcpKeepalive) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TcpKeepalive) ProtoMessage() {}

func (x *TcpKeepalive) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TcpKeepalive.ProtoReflect.Descriptor instead.
func (*TcpKeepalive) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP(), []int{2}
}

func (x *TcpKeepalive) GetKeepaliveProbes() *wrapperspb.UInt32Value {
	if x != nil {
		return x.KeepaliveProbes
	}
	return nil
}

func (x *TcpKeepalive) GetKeepaliveTime() *wrapperspb.UInt32Value {
	if x != nil {
		return x.KeepaliveTime
	}
	return nil
}

func (x *TcpKeepalive) GetKeepaliveInterval() *wrapperspb.UInt32Value {
	if x != nil {
		return x.KeepaliveInterval
	}
	return nil
}

type BindConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The address to bind to when creating a socket.
	SourceAddress *SocketAddress `protobuf:"bytes,1,opt,name=source_address,json=sourceAddress,proto3" json:"source_address,omitempty"`
	// Whether to set the *IP_FREEBIND* option when creating the socket. When this
	// flag is set to true, allows the :ref:`source_address
	// <envoy_api_field_config.cluster.v3.UpstreamBindConfig.source_address>` to be an IP address
	// that is not configured on the system running Envoy. When this flag is set
	// to false, the option *IP_FREEBIND* is disabled on the socket. When this
	// flag is not set (default), the socket is not modified, i.e. the option is
	// neither enabled nor disabled.
	Freebind *wrapperspb.BoolValue `protobuf:"bytes,2,opt,name=freebind,proto3" json:"freebind,omitempty"`
	// Additional socket options that may not be present in Envoy source code or
	// precompiled binaries.
	SocketOptions []*SocketOption `protobuf:"bytes,3,rep,name=socket_options,json=socketOptions,proto3" json:"socket_options,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BindConfig) Reset() {
	*x = BindConfig{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BindConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BindConfig) ProtoMessage() {}

func (x *BindConfig) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BindConfig.ProtoReflect.Descriptor instead.
func (*BindConfig) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP(), []int{3}
}

func (x *BindConfig) GetSourceAddress() *SocketAddress {
	if x != nil {
		return x.SourceAddress
	}
	return nil
}

func (x *BindConfig) GetFreebind() *wrapperspb.BoolValue {
	if x != nil {
		return x.Freebind
	}
	return nil
}

func (x *BindConfig) GetSocketOptions() []*SocketOption {
	if x != nil {
		return x.SocketOptions
	}
	return nil
}

// Addresses specify either a logical or physical address and port, which are
// used to tell Envoy where to bind/listen, connect to upstream and find
// management servers.
type Address struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Address:
	//
	//	*Address_SocketAddress
	//	*Address_Pipe
	Address       isAddress_Address `protobuf_oneof:"address"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Address) Reset() {
	*x = Address{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP(), []int{4}
}

func (x *Address) GetAddress() isAddress_Address {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *Address) GetSocketAddress() *SocketAddress {
	if x != nil {
		if x, ok := x.Address.(*Address_SocketAddress); ok {
			return x.SocketAddress
		}
	}
	return nil
}

func (x *Address) GetPipe() *Pipe {
	if x != nil {
		if x, ok := x.Address.(*Address_Pipe); ok {
			return x.Pipe
		}
	}
	return nil
}

type isAddress_Address interface {
	isAddress_Address()
}

type Address_SocketAddress struct {
	SocketAddress *SocketAddress `protobuf:"bytes,1,opt,name=socket_address,json=socketAddress,proto3,oneof"`
}

type Address_Pipe struct {
	Pipe *Pipe `protobuf:"bytes,2,opt,name=pipe,proto3,oneof"`
}

func (*Address_SocketAddress) isAddress_Address() {}

func (*Address_Pipe) isAddress_Address() {}

// CidrRange specifies an IP Address and a prefix length to construct
// the subnet mask for a [CIDR](https://datatracker.ietf.org/doc/html/rfc4632) range.
type CidrRange struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// IPv4 or IPv6 address, e.g. `192.0.0.0` or `2001:db8::`.
	AddressPrefix string `protobuf:"bytes,1,opt,name=address_prefix,json=addressPrefix,proto3" json:"address_prefix,omitempty"`
	// Length of prefix, e.g. 0, 32.
	PrefixLen     *wrapperspb.UInt32Value `protobuf:"bytes,2,opt,name=prefix_len,json=prefixLen,proto3" json:"prefix_len,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CidrRange) Reset() {
	*x = CidrRange{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CidrRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CidrRange) ProtoMessage() {}

func (x *CidrRange) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CidrRange.ProtoReflect.Descriptor instead.
func (*CidrRange) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP(), []int{5}
}

func (x *CidrRange) GetAddressPrefix() string {
	if x != nil {
		return x.AddressPrefix
	}
	return ""
}

func (x *CidrRange) GetPrefixLen() *wrapperspb.UInt32Value {
	if x != nil {
		return x.PrefixLen
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDesc = "" +
	"\n" +
	"Ugithub.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/address.proto\x12\x1csolo.io.envoy.config.core.v3\x1a[github.com/solo-io/gloo/projects/gloo/api/external/envoy/config/core/v3/socket_option.proto\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x1dudpa/annotations/status.proto\x1a!udpa/annotations/versioning.proto\x1a\x17validate/validate.proto\x1a\x12extproto/ext.proto\"i\n" +
	"\x04Pipe\x12\x1b\n" +
	"\x04path\x18\x01 \x01(\tB\a\xfaB\x04r\x02 \x01R\x04path\x12\x1c\n" +
	"\x04mode\x18\x02 \x01(\rB\b\xfaB\x05*\x03\x18\xff\x03R\x04mode:&\x8a\xc8ގ\x04 \n" +
	"\x1esolo.io.envoy.api.v2.core.Pipe\"\x87\x03\n" +
	"\rSocketAddress\x12Z\n" +
	"\bprotocol\x18\x01 \x01(\x0e24.solo.io.envoy.config.core.v3.SocketAddress.ProtocolB\b\xfaB\x05\x82\x01\x02\x10\x01R\bprotocol\x12!\n" +
	"\aaddress\x18\x02 \x01(\tB\a\xfaB\x04r\x02 \x01R\aaddress\x12*\n" +
	"\n" +
	"port_value\x18\x03 \x01(\rB\t\xfaB\x06*\x04\x18\xff\xff\x03H\x00R\tportValue\x12\x1f\n" +
	"\n" +
	"named_port\x18\x04 \x01(\tH\x00R\tnamedPort\x12#\n" +
	"\rresolver_name\x18\x05 \x01(\tR\fresolverName\x12\x1f\n" +
	"\vipv4_compat\x18\x06 \x01(\bR\n" +
	"ipv4Compat\"\x1c\n" +
	"\bProtocol\x12\a\n" +
	"\x03TCP\x10\x00\x12\a\n" +
	"\x03UDP\x10\x01:/\x8a\xc8ގ\x04)\n" +
	"'solo.io.envoy.api.v2.core.SocketAddressB\x15\n" +
	"\x0eport_specifier\x12\x03\xf8B\x01\"\x99\x02\n" +
	"\fTcpKeepalive\x12G\n" +
	"\x10keepalive_probes\x18\x01 \x01(\v2\x1c.google.protobuf.UInt32ValueR\x0fkeepaliveProbes\x12C\n" +
	"\x0ekeepalive_time\x18\x02 \x01(\v2\x1c.google.protobuf.UInt32ValueR\rkeepaliveTime\x12K\n" +
	"\x12keepalive_interval\x18\x03 \x01(\v2\x1c.google.protobuf.UInt32ValueR\x11keepaliveInterval:.\x8a\xc8ގ\x04(\n" +
	"&solo.io.envoy.api.v2.core.TcpKeepalive\"\xa3\x02\n" +
	"\n" +
	"BindConfig\x12\\\n" +
	"\x0esource_address\x18\x01 \x01(\v2+.solo.io.envoy.config.core.v3.SocketAddressB\b\xfaB\x05\x8a\x01\x02\x10\x01R\rsourceAddress\x126\n" +
	"\bfreebind\x18\x02 \x01(\v2\x1a.google.protobuf.BoolValueR\bfreebind\x12Q\n" +
	"\x0esocket_options\x18\x03 \x03(\v2*.solo.io.envoy.config.core.v3.SocketOptionR\rsocketOptions:,\x8a\xc8ގ\x04&\n" +
	"$solo.io.envoy.api.v2.core.BindConfig\"\xd4\x01\n" +
	"\aAddress\x12T\n" +
	"\x0esocket_address\x18\x01 \x01(\v2+.solo.io.envoy.config.core.v3.SocketAddressH\x00R\rsocketAddress\x128\n" +
	"\x04pipe\x18\x02 \x01(\v2\".solo.io.envoy.config.core.v3.PipeH\x00R\x04pipe:)\x8a\xc8ގ\x04#\n" +
	"!solo.io.envoy.api.v2.core.AddressB\x0e\n" +
	"\aaddress\x12\x03\xf8B\x01\"\xaf\x01\n" +
	"\tCidrRange\x12.\n" +
	"\x0eaddress_prefix\x18\x01 \x01(\tB\a\xfaB\x04r\x02 \x01R\raddressPrefix\x12E\n" +
	"\n" +
	"prefix_len\x18\x02 \x01(\v2\x1c.google.protobuf.UInt32ValueB\b\xfaB\x05*\x03\x18\x80\x01R\tprefixLen:+\x8a\xc8ގ\x04%\n" +
	"#solo.io.envoy.api.v2.core.CidrRangeB\x9d\x01\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01\xe2\xb5\xdf\xcb\a\x02\x10\x02\n" +
	"*io.envoyproxy.solo.io.envoy.config.core.v3B\fAddressProtoP\x01ZKgithub.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3b\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_goTypes = []any{
	(SocketAddress_Protocol)(0),    // 0: solo.io.envoy.config.core.v3.SocketAddress.Protocol
	(*Pipe)(nil),                   // 1: solo.io.envoy.config.core.v3.Pipe
	(*SocketAddress)(nil),          // 2: solo.io.envoy.config.core.v3.SocketAddress
	(*TcpKeepalive)(nil),           // 3: solo.io.envoy.config.core.v3.TcpKeepalive
	(*BindConfig)(nil),             // 4: solo.io.envoy.config.core.v3.BindConfig
	(*Address)(nil),                // 5: solo.io.envoy.config.core.v3.Address
	(*CidrRange)(nil),              // 6: solo.io.envoy.config.core.v3.CidrRange
	(*wrapperspb.UInt32Value)(nil), // 7: google.protobuf.UInt32Value
	(*wrapperspb.BoolValue)(nil),   // 8: google.protobuf.BoolValue
	(*SocketOption)(nil),           // 9: solo.io.envoy.config.core.v3.SocketOption
}
var file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_depIdxs = []int32{
	0,  // 0: solo.io.envoy.config.core.v3.SocketAddress.protocol:type_name -> solo.io.envoy.config.core.v3.SocketAddress.Protocol
	7,  // 1: solo.io.envoy.config.core.v3.TcpKeepalive.keepalive_probes:type_name -> google.protobuf.UInt32Value
	7,  // 2: solo.io.envoy.config.core.v3.TcpKeepalive.keepalive_time:type_name -> google.protobuf.UInt32Value
	7,  // 3: solo.io.envoy.config.core.v3.TcpKeepalive.keepalive_interval:type_name -> google.protobuf.UInt32Value
	2,  // 4: solo.io.envoy.config.core.v3.BindConfig.source_address:type_name -> solo.io.envoy.config.core.v3.SocketAddress
	8,  // 5: solo.io.envoy.config.core.v3.BindConfig.freebind:type_name -> google.protobuf.BoolValue
	9,  // 6: solo.io.envoy.config.core.v3.BindConfig.socket_options:type_name -> solo.io.envoy.config.core.v3.SocketOption
	2,  // 7: solo.io.envoy.config.core.v3.Address.socket_address:type_name -> solo.io.envoy.config.core.v3.SocketAddress
	1,  // 8: solo.io.envoy.config.core.v3.Address.pipe:type_name -> solo.io.envoy.config.core.v3.Pipe
	7,  // 9: solo.io.envoy.config.core.v3.CidrRange.prefix_len:type_name -> google.protobuf.UInt32Value
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_init()
}
func file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto != nil {
		return
	}
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_socket_option_proto_init()
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[1].OneofWrappers = []any{
		(*SocketAddress_PortValue)(nil),
		(*SocketAddress_NamedPort)(nil),
	}
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes[4].OneofWrappers = []any{
		(*Address_SocketAddress)(nil),
		(*Address_Pipe)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_depIdxs,
		EnumInfos:         file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_enumTypes,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_external_envoy_config_core_v3_address_proto_depIdxs = nil
}
