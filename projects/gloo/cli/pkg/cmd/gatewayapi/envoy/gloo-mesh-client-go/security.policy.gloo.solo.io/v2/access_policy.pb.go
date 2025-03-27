// {{% reuse "conrefs/snippets/policies/ov_access.md" %}}
// An access policy describes how clients should be authenticated and authorized
// to access a service. For more information about cross-origin resource sharing,
// see [this article](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS).
// AccessPolicies can be applied at the workload or destination level.
//
// **Note**:
// Workload selectors are considered more secure than destination selectors
// and are therefore recommended to be used whenever possible. While workload
// selectors are applied to the translated Istio AuthorizationPolicy resource directly,
// destination selectors require the translation of the selected destination first
// before the access policy is enforced. This can lead to a window where traffic is unsecured
// if a new destination is added to the cluster. However, keep in mind that workload selector
// cannot be used when service isolation is enabled in your workspace.
// If service isolation is enabled, you must use destination selectors instead.
// Note that virtual destinations are not supported as destinations with this policy.
//
// ## Examples
// The following example is for a simple access policy that allows
// the productpage app to access the ratings app.
// ```yaml
// apiVersion: security.policy.gloo.solo.io/v2
// kind: AccessPolicy
// metadata:
//   name: ratings-access
//   namespace: bookinfo
// spec:
//   applyToWorkloads:
//   - selector:
//       labels:
//         app: ratings
//   config:
//     authn:
//       tlsMode: STRICT
//     authzList:
//     - allowedClients:
//       - serviceAccountSelector:
//           labels:
//             app: productpage
//       allowedPaths:
//       - /ratings*
// ```
//
// The following example specifies the IP address that you want to
// allow access to the ratings app. When a client sends a request to ratings,
// the client's IP address is matched against the IP addresses that are defined
// in the access policy. If the IP address matches, the request is forwarded
// to the ratings app. If the IP does not match, access to ratings is denied.
// ```yaml
// apiVersion: security.policy.gloo.solo.io/v2
// kind: AccessPolicy
// metadata:
//   name: ratings-access
//   namespace: bookinfo
// spec:
//   applyToWorkloads:
//   - selector:
//       labels:
//         app: ratings
//   config:
//     authn:
//       tlsMode: STRICT
//     authzList:
//     - allowedIpBlocks:
//       - 112.114.230.1
//       allowedPaths:
//       - /ratings*
// ```
//
// The following access policy uses a destination selector and specifies the request
// headers that must be sent to allow or deny the communication between the productpage
// and ratings apps. For example, if you send a request with the `X-Test-Header: match`
// header from the productpage app to the ratings app, the request is matched and the
// communication between productpage and ratings is allowed. If you send the same request
// without a header or with the `X-Test-Header: noMatch` header, the request is not matched
// and the communication between the apps is denied.
// ```yaml
// apiVersion: security.policy.gloo.solo.io/v2
// kind: AccessPolicy
// metadata:
//   name: ratings-access
//   namespace: bookinfo
// spec:
//   applyToDestinations:
//   - port:
//       number: 9080
//     selector:
//       labels:
//         app: ratings
//   config:
//     authn:
//       tlsMode: STRICT
//     authzList:
//     - allowedClients:
//       - serviceAccountSelector:
//           labels:
//             app: productpage
//       allowedPaths:
//       - /ratings*
//       match:
//         request:
//           headers:
//             X-Test-Header:
//               notValues:
//               - noMatch
//               - partial-blocked
//               values:
//               - match
//               - partial*
// ```
//
// You can have multiple `authzList` entries to control access to workloads.
// A request is allowed when it matches at least one of the `authzList` entries (logically OR'd together).
//
// For each entry, you can specify different requirements for allowed clients, paths, methods, IP blocks, and other configuration settings.
// Then, a request is allowed only when ALL of the requirements are met (logically AND'd together).
//
// In the following example:
// * The product page app is allowed to send GET requests to the ratings app along the `/ratings*` wildcard path.
// * The product page app is allowed to send PATCH requests to the ratings app along the `/ratings/2*` wildcard path.
// * The reviews app is allowed to access the ratings app.
//
// ```yaml
// apiVersion: security.policy.gloo.solo.io/v2
// kind: AccessPolicy
// metadata:
//   name: ratings-access
//   namespace: bookinfo
// spec:
//   applyToWorkloads:
//   - selector:
//       labels:
//         app: ratings
//   config:
//     authn:
//       tlsMode: STRICT
//     authzList:
//     - allowedClients:
//       - serviceAccountSelector:
//           labels:
//             app: productpage
//       allowedPaths:
//       - /ratings*
//       allowedMethods:
//       - GET
//     - allowedClients:
//       - serviceAccountSelector:
//           labels:
//             app: productpage
//       allowedPaths:
//       - /ratings/2*
//       allowedMethods:
//       - PATCH
//     - allowedClients:
//       - serviceAccountSelector:
//           labels:
//             app: reviews
// ```

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: github.com/solo-io/solo-apis/api/gloo.solo.io/policy/v2/security/access_policy.proto

package v2

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/solo-io/protoc-gen-ext/extproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	v2 "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/envoy/gloo-mesh-client-go/common.gloo.solo.io/v2"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The mutual TLS (mTLS) connection mode. The following enums correspond to the
// [modes defined by Istio](https://github.com/istio/api/blob/master/security/v1beta1/peer_authentication.proto#L129).
type AccessPolicySpec_Config_Authentication_TLSmode int32

const (
	// Do not originate mTLS connections to the upstream workload,
	// and instead use unencrypted plaintext.
	AccessPolicySpec_Config_Authentication_DISABLE AccessPolicySpec_Config_Authentication_TLSmode = 0
	// Permit both unencrypted plaintext and mTLS-secured connections to the upstream workload.
	// Use this mode only when you migrate workloads to your service mesh.
	// After the workload is onboarded to the service mesh, using the `STRICT` mode is recommended.
	AccessPolicySpec_Config_Authentication_PERMISSIVE AccessPolicySpec_Config_Authentication_TLSmode = 1
	// Secure connections to the upstream workload with mTLS by presenting
	// client certificates for authentication.
	// This mode uses certificates generated
	// automatically by Istio for mTLS authentication. When you use
	// this mode, keep all other fields in `ClientTLSSettings` empty.
	AccessPolicySpec_Config_Authentication_STRICT AccessPolicySpec_Config_Authentication_TLSmode = 2
)

// Enum value maps for AccessPolicySpec_Config_Authentication_TLSmode.
var (
	AccessPolicySpec_Config_Authentication_TLSmode_name = map[int32]string{
		0: "DISABLE",
		1: "PERMISSIVE",
		2: "STRICT",
	}
	AccessPolicySpec_Config_Authentication_TLSmode_value = map[string]int32{
		"DISABLE":    0,
		"PERMISSIVE": 1,
		"STRICT":     2,
	}
)

func (x AccessPolicySpec_Config_Authentication_TLSmode) Enum() *AccessPolicySpec_Config_Authentication_TLSmode {
	p := new(AccessPolicySpec_Config_Authentication_TLSmode)
	*p = x
	return p
}

func (x AccessPolicySpec_Config_Authentication_TLSmode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccessPolicySpec_Config_Authentication_TLSmode) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_enumTypes[0].Descriptor()
}

func (AccessPolicySpec_Config_Authentication_TLSmode) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_enumTypes[0]
}

func (x AccessPolicySpec_Config_Authentication_TLSmode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccessPolicySpec_Config_Authentication_TLSmode.Descriptor instead.
func (AccessPolicySpec_Config_Authentication_TLSmode) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 1, 0, 0}
}

// Specifications for the policy.
type AccessPolicySpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Destinations to apply the policy to.
	// Note that virtual destinations are not supported as destinations with this policy.
	// If `applyToWorkloads` is non-empty, this field is ignored.
	// If this field and `applyToWorkloads` are both empty,
	// the policy applies to all ports on all destinations in the workspace.
	// {{< alert context="info" >}}
	// For security reasons, <code>applyToWorkloads</code> is preferred.
	// {{< /alert >}}
	ApplyToDestinations []*v2.DestinationSelector `protobuf:"bytes,1,rep,name=apply_to_destinations,json=applyToDestinations,proto3" json:"apply_to_destinations,omitempty"`
	// Workloads to apply the policy to. For security reasons,
	// this field is preferred over `applyToDestinations`. If an empty selector is
	// provided in the list, the policy applies to all workloads in a namespace, cluster,
	// and workspace that are available in the parent object's workspace.
	ApplyToWorkloads []*AccessPolicySpec_NamespaceWorkloadSelector `protobuf:"bytes,3,rep,name=apply_to_workloads,json=applyToWorkloads,proto3" json:"apply_to_workloads,omitempty"`
	// Details of the access policy to apply to the selected workloads.
	Config *AccessPolicySpec_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *AccessPolicySpec) Reset() {
	*x = AccessPolicySpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec) ProtoMessage() {}

func (x *AccessPolicySpec) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0}
}

func (x *AccessPolicySpec) GetApplyToDestinations() []*v2.DestinationSelector {
	if x != nil {
		return x.ApplyToDestinations
	}
	return nil
}

func (x *AccessPolicySpec) GetApplyToWorkloads() []*AccessPolicySpec_NamespaceWorkloadSelector {
	if x != nil {
		return x.ApplyToWorkloads
	}
	return nil
}

func (x *AccessPolicySpec) GetConfig() *AccessPolicySpec_Config {
	if x != nil {
		return x.Config
	}
	return nil
}

// The status of the policy after it is applied to your Gloo environment.
type AccessPolicyStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The state and workspace conditions of the applied policy.
	Common *v2.Status `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	// The number of destination ports selected by the policy.
	NumSelectedDestinationPorts uint32 `protobuf:"varint,5,opt,name=num_selected_destination_ports,json=numSelectedDestinationPorts,proto3" json:"num_selected_destination_ports,omitempty"`
	// The number of namespaces containing selected workloads by the policy.
	NumSelectedNamespaces uint32 `protobuf:"varint,2,opt,name=num_selected_namespaces,json=numSelectedNamespaces,proto3" json:"num_selected_namespaces,omitempty"`
	// The number of service accounts allowed to access the selected destinations.
	NumAllowedServiceAccounts uint32 `protobuf:"varint,3,opt,name=num_allowed_service_accounts,json=numAllowedServiceAccounts,proto3" json:"num_allowed_service_accounts,omitempty"`
}

func (x *AccessPolicyStatus) Reset() {
	*x = AccessPolicyStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicyStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicyStatus) ProtoMessage() {}

func (x *AccessPolicyStatus) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicyStatus.ProtoReflect.Descriptor instead.
func (*AccessPolicyStatus) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{1}
}

func (x *AccessPolicyStatus) GetCommon() *v2.Status {
	if x != nil {
		return x.Common
	}
	return nil
}

func (x *AccessPolicyStatus) GetNumSelectedDestinationPorts() uint32 {
	if x != nil {
		return x.NumSelectedDestinationPorts
	}
	return 0
}

func (x *AccessPolicyStatus) GetNumSelectedNamespaces() uint32 {
	if x != nil {
		return x.NumSelectedNamespaces
	}
	return 0
}

func (x *AccessPolicyStatus) GetNumAllowedServiceAccounts() uint32 {
	if x != nil {
		return x.NumAllowedServiceAccounts
	}
	return 0
}

// The report shows the resources that the policy selects after the policy is successfully applied.
type AccessPolicyReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of workspaces in which the policy can apply to destinations.
	Workspaces map[string]*v2.Report `protobuf:"bytes,1,rep,name=workspaces,proto3" json:"workspaces,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// A list of destination ports selected by the policy.
	SelectedDestinationPorts []*v2.DestinationReference `protobuf:"bytes,2,rep,name=selected_destination_ports,json=selectedDestinationPorts,proto3" json:"selected_destination_ports,omitempty"`
	// A list of the service accounts whose workloads are allowed
	// to send requests to the selected destinations.
	AllowedServiceAccounts []*AccessPolicyReport_IdentityReference `protobuf:"bytes,3,rep,name=allowed_service_accounts,json=allowedServiceAccounts,proto3" json:"allowed_service_accounts,omitempty"`
	// A list of namespaces that contain workloads selected by the policy.
	SelectedNamespaces []*v2.ObjectReference `protobuf:"bytes,4,rep,name=selected_namespaces,json=selectedNamespaces,proto3" json:"selected_namespaces,omitempty"`
}

func (x *AccessPolicyReport) Reset() {
	*x = AccessPolicyReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicyReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicyReport) ProtoMessage() {}

func (x *AccessPolicyReport) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicyReport.ProtoReflect.Descriptor instead.
func (*AccessPolicyReport) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{2}
}

func (x *AccessPolicyReport) GetWorkspaces() map[string]*v2.Report {
	if x != nil {
		return x.Workspaces
	}
	return nil
}

func (x *AccessPolicyReport) GetSelectedDestinationPorts() []*v2.DestinationReference {
	if x != nil {
		return x.SelectedDestinationPorts
	}
	return nil
}

func (x *AccessPolicyReport) GetAllowedServiceAccounts() []*AccessPolicyReport_IdentityReference {
	if x != nil {
		return x.AllowedServiceAccounts
	}
	return nil
}

func (x *AccessPolicyReport) GetSelectedNamespaces() []*v2.ObjectReference {
	if x != nil {
		return x.SelectedNamespaces
	}
	return nil
}

// Select individual namespaces and workloads within the namespaces by label.
// Workloads must have injected (sidecars) or be standalone proxies (gateways)
// to be selected by Gloo policies.
type AccessPolicySpec_NamespaceWorkloadSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Selector to match workload objects by their metadata.
	Selector *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector) Reset() {
	*x = AccessPolicySpec_NamespaceWorkloadSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_NamespaceWorkloadSelector) ProtoMessage() {}

func (x *AccessPolicySpec_NamespaceWorkloadSelector) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_NamespaceWorkloadSelector.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_NamespaceWorkloadSelector) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 0}
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector) GetSelector() *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector {
	if x != nil {
		return x.Selector
	}
	return nil
}

// Details of the access policy to apply to the selected workloads.
type AccessPolicySpec_Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// How clients are authenticated to the workload.
	Authn *AccessPolicySpec_Config_Authentication `protobuf:"bytes,1,opt,name=authn,proto3" json:"authn,omitempty"`
	// Deprecated; Use authzList instead. If authzList is set, it takes precedence and this field is ignored.
	// How clients are authorized to access the workload.
	Authz *AccessPolicySpec_Config_Authorization `protobuf:"bytes,2,opt,name=authz,proto3" json:"authz,omitempty"`
	// Optional: When NetworkPolicy translation is enabled, all available layers are used to enforce AccessPolicies by default.
	// If you want to explicitly define which layers to use to enforce this AccessPolicy, you can set them by using this field.
	// Note that the layer that you define in this field must be available to be configured.
	EnforcementLayers *v2.EnforcementLayers `protobuf:"bytes,3,opt,name=enforcement_layers,json=enforcementLayers,proto3" json:"enforcement_layers,omitempty"`
	// How clients are authorized to access the workload.
	// A request is allowed when it matches at least one authz entry in the list (logically OR'd together).
	AuthzList []*AccessPolicySpec_Config_Authorization `protobuf:"bytes,4,rep,name=authz_list,json=authzList,proto3" json:"authz_list,omitempty"`
}

func (x *AccessPolicySpec_Config) Reset() {
	*x = AccessPolicySpec_Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_Config) ProtoMessage() {}

func (x *AccessPolicySpec_Config) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_Config.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_Config) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 1}
}

func (x *AccessPolicySpec_Config) GetAuthn() *AccessPolicySpec_Config_Authentication {
	if x != nil {
		return x.Authn
	}
	return nil
}

func (x *AccessPolicySpec_Config) GetAuthz() *AccessPolicySpec_Config_Authorization {
	if x != nil {
		return x.Authz
	}
	return nil
}

func (x *AccessPolicySpec_Config) GetEnforcementLayers() *v2.EnforcementLayers {
	if x != nil {
		return x.EnforcementLayers
	}
	return nil
}

func (x *AccessPolicySpec_Config) GetAuthzList() []*AccessPolicySpec_Config_Authorization {
	if x != nil {
		return x.AuthzList
	}
	return nil
}

// Selects zero or more Kubernetes workloads by matching on labels, namespace, cluster, and workspace.
type AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Only select workloads with matching labels.
	Labels map[string]string `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Only select objects in the matching namespace. If empty, Gloo selects
	// matching objects across all namespaces available in the parent object's workspace.
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Only select objects in the matching cluster. If empty, Gloo selects
	// matching objects across all clusters available in the parent object's workspace.
	Cluster string `protobuf:"bytes,3,opt,name=cluster,proto3" json:"cluster,omitempty"`
	// Only select objects in the given workspace. If empty, Gloo selects
	// matching objects across all workspaces available in the parent object's workspace.
	Workspace string `protobuf:"bytes,4,opt,name=workspace,proto3" json:"workspace,omitempty"`
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) Reset() {
	*x = AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) ProtoMessage() {}

func (x *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 0, 0}
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) GetCluster() string {
	if x != nil {
		return x.Cluster
	}
	return ""
}

func (x *AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector) GetWorkspace() string {
	if x != nil {
		return x.Workspace
	}
	return ""
}

// How clients are authenticated to the workload.
type AccessPolicySpec_Config_Authentication struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Type of TLS policy that is enforced on clients connecting to the workload.
	// If service isolation is enabled for the workspace, this field is treated as 'STRICT'.
	TlsMode AccessPolicySpec_Config_Authentication_TLSmode `protobuf:"varint,1,opt,name=tls_mode,json=tlsMode,proto3,enum=security.policy.gloo.solo.io.AccessPolicySpec_Config_Authentication_TLSmode" json:"tls_mode,omitempty"`
}

func (x *AccessPolicySpec_Config_Authentication) Reset() {
	*x = AccessPolicySpec_Config_Authentication{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_Config_Authentication) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_Config_Authentication) ProtoMessage() {}

func (x *AccessPolicySpec_Config_Authentication) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_Config_Authentication.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_Config_Authentication) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 1, 0}
}

func (x *AccessPolicySpec_Config_Authentication) GetTlsMode() AccessPolicySpec_Config_Authentication_TLSmode {
	if x != nil {
		return x.TlsMode
	}
	return AccessPolicySpec_Config_Authentication_DISABLE
}

// Configure access to workloads.
//
// You can have multiple `authzList` entries.
// A request is allowed when it matches at least one of the `authzList` entries (logically OR'd together).
//
// For each entry, you can specify different requirements for allowed clients, paths, methods, IP blocks, and other configuration settings.
// Then, a request is allowed only when ALL of the requirements are met (logically AND'd together).
//
// If the policy uses `applyToWorkloads`, you can also allow NO requests by setting this value to the empty object `{}`,
// which will serve as a fallback when requests do not match another `authz` case for the given workload.
type AccessPolicySpec_Config_Authorization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Client identities that are permitted to access the workload.
	// To allow access for all client identities, provide a single empty selector.
	AllowedClients []*v2.IdentitySelector `protobuf:"bytes,1,rep,name=allowed_clients,json=allowedClients,proto3" json:"allowed_clients,omitempty"`
	// Optional: A list of HTTP paths or gRPC methods to allow.
	// gRPC methods must be presented as fully-qualified name in the form of
	// "/packageName.serviceName/methodName", and are case sensitive.
	// Exact match, prefix match, and suffix match are supported for paths.
	// For example, the path `/books/review` matches
	// `/books/review` (exact match), `*books/` (suffix match), or `/books*` (prefix match).
	//
	// If empty, any path is allowed.
	AllowedPaths []string `protobuf:"bytes,2,rep,name=allowed_paths,json=allowedPaths,proto3" json:"allowed_paths,omitempty"`
	// Optional: A list of HTTP methods to allow (e.g., "GET", "POST").
	// If empty, any method is allowed.
	// This field is ignored for gRPC, because the value is always "POST".
	AllowedMethods []string `protobuf:"bytes,3,rep,name=allowed_methods,json=allowedMethods,proto3" json:"allowed_methods,omitempty"`
	// Optional: Additional request matching conditions.
	Match *AccessPolicySpec_Config_Authorization_MatchSpec `protobuf:"bytes,4,opt,name=match,proto3" json:"match,omitempty"`
	// Optional: A list of IP blocks, populated from the source address of the IP packet.
	// Single IP addresses (e.g. “1.2.3.4”) and CIDRs (e.g. “1.2.3.0/24”) are supported. If empty,
	// any IP address is allowed.
	AllowedIpBlocks []string `protobuf:"bytes,5,rep,name=allowed_ip_blocks,json=allowedIpBlocks,proto3" json:"allowed_ip_blocks,omitempty"`
	// Optional: A list of IP blocks, populated from X-Forwarded-For header or proxy protocol.
	// Single IP addresses (e.g. “1.2.3.4”) and CIDRs (e.g. “1.2.3.0/24”) are supported. This field
	// is equivalent to the remote.ip attribute. If empty, any IP address is allowed.
	// {{< alert >}}
	// To use this field, you must configure the <code>meshConfig.defaultConfig.gatewayTopology.numTrustedProxies</code>
	// field in your Istio installation. For more info, see the
	// <a href="https://istio.io/latest/docs/ops/configuration/traffic-management/network-topologies/#configuring-network-topologies">Istio documentation</a>.
	// {{< /alert >}}
	AllowedRemoteIpBlocks []string `protobuf:"bytes,6,rep,name=allowed_remote_ip_blocks,json=allowedRemoteIpBlocks,proto3" json:"allowed_remote_ip_blocks,omitempty"`
	// Set to true to enable a dry run of the access policy for L7 Istio service mesh authorization only. Then, you can check the sidecar proxy logs, metrics, and tracing to determine if traffic would be allowed or denied. However, the authorization is not enforced until you disable the dry run and re-apply the access policy.
	// Note that when there are both dry run and enforced policies, dry run policies are considered independently of enforced policies;
	// i.e. the logs, metrics, and tracing results indicating if traffic would be allowed or denied is based on the behavior if all dry run policies were enforced but all currently enforced policies were deleted.
	// Note that dry run cannot be used to review allow or deny decisions for L4 traffic. Even if you enable the dry run feature with a Gloo Network setup, no Cilium network policy and decision logs are created or enforced.
	DryRun bool `protobuf:"varint,7,opt,name=dry_run,json=dryRun,proto3" json:"dry_run,omitempty"`
}

func (x *AccessPolicySpec_Config_Authorization) Reset() {
	*x = AccessPolicySpec_Config_Authorization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_Config_Authorization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_Config_Authorization) ProtoMessage() {}

func (x *AccessPolicySpec_Config_Authorization) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_Config_Authorization.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_Config_Authorization) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 1, 1}
}

func (x *AccessPolicySpec_Config_Authorization) GetAllowedClients() []*v2.IdentitySelector {
	if x != nil {
		return x.AllowedClients
	}
	return nil
}

func (x *AccessPolicySpec_Config_Authorization) GetAllowedPaths() []string {
	if x != nil {
		return x.AllowedPaths
	}
	return nil
}

func (x *AccessPolicySpec_Config_Authorization) GetAllowedMethods() []string {
	if x != nil {
		return x.AllowedMethods
	}
	return nil
}

func (x *AccessPolicySpec_Config_Authorization) GetMatch() *AccessPolicySpec_Config_Authorization_MatchSpec {
	if x != nil {
		return x.Match
	}
	return nil
}

func (x *AccessPolicySpec_Config_Authorization) GetAllowedIpBlocks() []string {
	if x != nil {
		return x.AllowedIpBlocks
	}
	return nil
}

func (x *AccessPolicySpec_Config_Authorization) GetAllowedRemoteIpBlocks() []string {
	if x != nil {
		return x.AllowedRemoteIpBlocks
	}
	return nil
}

func (x *AccessPolicySpec_Config_Authorization) GetDryRun() bool {
	if x != nil {
		return x.DryRun
	}
	return false
}

// Optional: Additional request matching conditions.
type AccessPolicySpec_Config_Authorization_MatchSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional: HTTP request header matching conditions.
	Request *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec) Reset() {
	*x = AccessPolicySpec_Config_Authorization_MatchSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_Config_Authorization_MatchSpec) ProtoMessage() {}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_Config_Authorization_MatchSpec.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_Config_Authorization_MatchSpec) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 1, 1, 0}
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec) GetRequest() *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec {
	if x != nil {
		return x.Request
	}
	return nil
}

type AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional: HTTP request header matching conditions.
	Headers map[string]*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues `protobuf:"bytes,1,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec) Reset() {
	*x = AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec) ProtoMessage() {}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 1, 1, 0, 0}
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec) GetHeaders() map[string]*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues {
	if x != nil {
		return x.Headers
	}
	return nil
}

// A value matching condition for HTTP request headers.
// At least one of `values` or `notValues` must be set.
type AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of match values for the header. A request must match at least one value.
	// Supports wildcards. For example, to
	// match a request with header values containing `exact-books` OR `partial-matched-books`,
	// set `values` to `exact-books` and `partial-*-books`.
	Values []string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	// A list of negative match values for the header. A request must not match any values.
	// Supports wildcards. For example, to
	// _not_ match a request with header values containing `ignore-books` or `partial-ignored-books`,
	// set `notValues` to `ignore-books` and `partial-ig*-books`.
	NotValues []string `protobuf:"bytes,2,rep,name=not_values,json=notValues,proto3" json:"not_values,omitempty"`
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues) Reset() {
	*x = AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues) ProtoMessage() {}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues.ProtoReflect.Descriptor instead.
func (*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{0, 1, 1, 0, 0, 1}
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues) GetNotValues() []string {
	if x != nil {
		return x.NotValues
	}
	return nil
}

// A list of the service accounts whose workloads are allowed
// to send requests to the selected destinations.
type AccessPolicyReport_IdentityReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The index of the identity in the list of identity selectors.
	IdentityIndex int32 `protobuf:"varint,1,opt,name=identity_index,json=identityIndex,proto3" json:"identity_index,omitempty"`
	// The reference to the service account backing the identity.
	ServiceAccount *v2.ObjectReference `protobuf:"bytes,2,opt,name=service_account,json=serviceAccount,proto3" json:"service_account,omitempty"`
	// The index of the authz in the authzList.
	AuthzIndex int32 `protobuf:"varint,3,opt,name=authz_index,json=authzIndex,proto3" json:"authz_index,omitempty"`
}

func (x *AccessPolicyReport_IdentityReference) Reset() {
	*x = AccessPolicyReport_IdentityReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPolicyReport_IdentityReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPolicyReport_IdentityReference) ProtoMessage() {}

func (x *AccessPolicyReport_IdentityReference) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPolicyReport_IdentityReference.ProtoReflect.Descriptor instead.
func (*AccessPolicyReport_IdentityReference) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP(), []int{2, 1}
}

func (x *AccessPolicyReport_IdentityReference) GetIdentityIndex() int32 {
	if x != nil {
		return x.IdentityIndex
	}
	return 0
}

func (x *AccessPolicyReport_IdentityReference) GetServiceAccount() *v2.ObjectReference {
	if x != nil {
		return x.ServiceAccount
	}
	return nil
}

func (x *AccessPolicyReport_IdentityReference) GetAuthzIndex() int32 {
	if x != nil {
		return x.AuthzIndex
	}
	return 0
}

var File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDesc = []byte{
	0x0a, 0x62, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x2f, 0x76, 0x32, 0x2f, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79,
	0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x1a, 0x12, 0x65, 0x78, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d,
	0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f,
	0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x55,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d,
	0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74,
	0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x5e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d,
	0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76,
	0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x65, 0x6e, 0x66,
	0x6f, 0x72, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x52, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d,
	0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76,
	0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad, 0x11, 0x0a, 0x10, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x12, 0x5c,
	0x0a, 0x15, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x13, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x54, 0x6f,
	0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x76, 0x0a, 0x12,
	0x61, 0x70, 0x70, 0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x48, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72,
	0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x57, 0x6f, 0x72, 0x6b, 0x6c,
	0x6f, 0x61, 0x64, 0x73, 0x12, 0x4d, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x53, 0x70, 0x65, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x1a, 0xb1, 0x03, 0x0a, 0x19, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x12, 0x73, 0x0a, 0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x57, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53,
	0x70, 0x65, 0x63, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x57, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x73, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a, 0x9e, 0x02, 0x0a, 0x0e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x7b, 0x0a, 0x06, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x63, 0x2e, 0x73, 0x65, 0x63, 0x75,
	0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f,
	0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06,
	0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1c,
	0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x1a, 0x39, 0x0a, 0x0b,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0xbf, 0x0b, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x5a, 0x0a, 0x05, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x44, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65,
	0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x12, 0x59,
	0x0a, 0x05, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x43, 0x2e,
	0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x05, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x12, 0x55, 0x0a, 0x12, 0x65, 0x6e, 0x66,
	0x6f, 0x72, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x45, 0x6e, 0x66, 0x6f,
	0x72, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x52, 0x11, 0x65,
	0x6e, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x73,
	0x12, 0x62, 0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x43, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x53, 0x70, 0x65, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x7a,
	0x4c, 0x69, 0x73, 0x74, 0x1a, 0xad, 0x01, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x67, 0x0a, 0x08, 0x74, 0x6c, 0x73, 0x5f, 0x6d,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x4c, 0x2e, 0x73, 0x65, 0x63, 0x75,
	0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f,
	0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x54, 0x4c, 0x53, 0x6d, 0x6f, 0x64, 0x65, 0x52, 0x07, 0x74, 0x6c, 0x73, 0x4d, 0x6f, 0x64, 0x65,
	0x22, 0x32, 0x0a, 0x07, 0x54, 0x4c, 0x53, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x44,
	0x49, 0x53, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x45, 0x52, 0x4d,
	0x49, 0x53, 0x53, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x52, 0x49,
	0x43, 0x54, 0x10, 0x02, 0x1a, 0x92, 0x07, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4e, 0x0a, 0x0f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x64, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x64, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x61,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x50, 0x61, 0x74, 0x68, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x61,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x73, 0x12, 0x63, 0x0a, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x4d, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53,
	0x70, 0x65, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x53, 0x70,
	0x65, 0x63, 0x52, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x2a, 0x0a, 0x11, 0x61, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x64, 0x5f, 0x69, 0x70, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x49, 0x70, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x37, 0x0a, 0x18, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64,
	0x5f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x5f, 0x69, 0x70, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x15, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64,
	0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x49, 0x70, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x17,
	0x0a, 0x07, 0x64, 0x72, 0x79, 0x5f, 0x72, 0x75, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x64, 0x72, 0x79, 0x52, 0x75, 0x6e, 0x1a, 0xff, 0x03, 0x0a, 0x09, 0x4d, 0x61, 0x74, 0x63,
	0x68, 0x53, 0x70, 0x65, 0x63, 0x12, 0x73, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x59, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74,
	0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x74, 0x63,
	0x68, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x70, 0x65,
	0x63, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0xfc, 0x02, 0x0a, 0x0b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x70, 0x65, 0x63, 0x12, 0x80, 0x01, 0x0a, 0x07, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x66, 0x2e, 0x73,
	0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x1a, 0xa2, 0x01,
	0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x7c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x66, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x70, 0x65, 0x63, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x1a, 0x45, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x6f,
	0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09,
	0x6e, 0x6f, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x87, 0x02, 0x0a, 0x12, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x33, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x12, 0x43, 0x0a, 0x1e, 0x6e, 0x75, 0x6d, 0x5f, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x1b, 0x6e,
	0x75, 0x6d, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x36, 0x0a, 0x17, 0x6e, 0x75,
	0x6d, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x15, 0x6e, 0x75, 0x6d,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x73, 0x12, 0x3f, 0x0a, 0x1c, 0x6e, 0x75, 0x6d, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x19, 0x6e, 0x75, 0x6d, 0x41, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x22, 0xbd, 0x05, 0x0a, 0x12, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x60, 0x0a, 0x0a, 0x77, 0x6f,
	0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x40,
	0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x12, 0x67, 0x0a, 0x1a,
	0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x29, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x18, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x7c, 0x0a, 0x18, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69,
	0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x16, 0x61, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x12, 0x55, 0x0a, 0x13, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x24, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x12, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64,
	0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x1a, 0x5a, 0x0a, 0x0f, 0x57, 0x6f,
	0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x31, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c,
	0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0xaa, 0x01, 0x0a, 0x11, 0x49, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x0e,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x4d, 0x0a, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x52, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x5f, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x42, 0x60, 0x5a, 0x52, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d,
	0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76,
	0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69,
	0x74, 0x79, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x32, 0xc0, 0xf5, 0x04, 0x01, 0xb8, 0xf5, 0x04,
	0x01, 0xd0, 0xf5, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescData = file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDesc
)

func file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescData)
	})
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDescData
}

var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_goTypes = []interface{}{
	(AccessPolicySpec_Config_Authentication_TLSmode)(0),               // 0: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authentication.TLSmode
	(*AccessPolicySpec)(nil),                                          // 1: security.policy.gloo.solo.io.AccessPolicySpec
	(*AccessPolicyStatus)(nil),                                        // 2: security.policy.gloo.solo.io.AccessPolicyStatus
	(*AccessPolicyReport)(nil),                                        // 3: security.policy.gloo.solo.io.AccessPolicyReport
	(*AccessPolicySpec_NamespaceWorkloadSelector)(nil),                // 4: security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector
	(*AccessPolicySpec_Config)(nil),                                   // 5: security.policy.gloo.solo.io.AccessPolicySpec.Config
	(*AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector)(nil), // 6: security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector.ObjectSelector
	nil, // 7: security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector.ObjectSelector.LabelsEntry
	(*AccessPolicySpec_Config_Authentication)(nil),                      // 8: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authentication
	(*AccessPolicySpec_Config_Authorization)(nil),                       // 9: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization
	(*AccessPolicySpec_Config_Authorization_MatchSpec)(nil),             // 10: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec
	(*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec)(nil), // 11: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec
	nil, // 12: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec.HeadersEntry
	(*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues)(nil), // 13: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec.HeaderValues
	nil, // 14: security.policy.gloo.solo.io.AccessPolicyReport.WorkspacesEntry
	(*AccessPolicyReport_IdentityReference)(nil), // 15: security.policy.gloo.solo.io.AccessPolicyReport.IdentityReference
	(*v2.DestinationSelector)(nil),               // 16: common.gloo.solo.io.DestinationSelector
	(*v2.Status)(nil),                            // 17: common.gloo.solo.io.Status
	(*v2.DestinationReference)(nil),              // 18: common.gloo.solo.io.DestinationReference
	(*v2.ObjectReference)(nil),                   // 19: common.gloo.solo.io.ObjectReference
	(*v2.EnforcementLayers)(nil),                 // 20: common.gloo.solo.io.EnforcementLayers
	(*v2.IdentitySelector)(nil),                  // 21: common.gloo.solo.io.IdentitySelector
	(*v2.Report)(nil),                            // 22: common.gloo.solo.io.Report
}
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_depIdxs = []int32{
	16, // 0: security.policy.gloo.solo.io.AccessPolicySpec.apply_to_destinations:type_name -> common.gloo.solo.io.DestinationSelector
	4,  // 1: security.policy.gloo.solo.io.AccessPolicySpec.apply_to_workloads:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector
	5,  // 2: security.policy.gloo.solo.io.AccessPolicySpec.config:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config
	17, // 3: security.policy.gloo.solo.io.AccessPolicyStatus.common:type_name -> common.gloo.solo.io.Status
	14, // 4: security.policy.gloo.solo.io.AccessPolicyReport.workspaces:type_name -> security.policy.gloo.solo.io.AccessPolicyReport.WorkspacesEntry
	18, // 5: security.policy.gloo.solo.io.AccessPolicyReport.selected_destination_ports:type_name -> common.gloo.solo.io.DestinationReference
	15, // 6: security.policy.gloo.solo.io.AccessPolicyReport.allowed_service_accounts:type_name -> security.policy.gloo.solo.io.AccessPolicyReport.IdentityReference
	19, // 7: security.policy.gloo.solo.io.AccessPolicyReport.selected_namespaces:type_name -> common.gloo.solo.io.ObjectReference
	6,  // 8: security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector.selector:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector.ObjectSelector
	8,  // 9: security.policy.gloo.solo.io.AccessPolicySpec.Config.authn:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authentication
	9,  // 10: security.policy.gloo.solo.io.AccessPolicySpec.Config.authz:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization
	20, // 11: security.policy.gloo.solo.io.AccessPolicySpec.Config.enforcement_layers:type_name -> common.gloo.solo.io.EnforcementLayers
	9,  // 12: security.policy.gloo.solo.io.AccessPolicySpec.Config.authz_list:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization
	7,  // 13: security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector.ObjectSelector.labels:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.NamespaceWorkloadSelector.ObjectSelector.LabelsEntry
	0,  // 14: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authentication.tls_mode:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authentication.TLSmode
	21, // 15: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.allowed_clients:type_name -> common.gloo.solo.io.IdentitySelector
	10, // 16: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.match:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec
	11, // 17: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.request:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec
	12, // 18: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec.headers:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec.HeadersEntry
	13, // 19: security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec.HeadersEntry.value:type_name -> security.policy.gloo.solo.io.AccessPolicySpec.Config.Authorization.MatchSpec.RequestSpec.HeaderValues
	22, // 20: security.policy.gloo.solo.io.AccessPolicyReport.WorkspacesEntry.value:type_name -> common.gloo.solo.io.Report
	19, // 21: security.policy.gloo.solo.io.AccessPolicyReport.IdentityReference.service_account:type_name -> common.gloo.solo.io.ObjectReference
	22, // [22:22] is the sub-list for method output_type
	22, // [22:22] is the sub-list for method input_type
	22, // [22:22] is the sub-list for extension type_name
	22, // [22:22] is the sub-list for extension extendee
	0,  // [0:22] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_init()
}
func file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_init() {
	if File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicyStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicyReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_NamespaceWorkloadSelector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_Config); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_NamespaceWorkloadSelector_ObjectSelector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_Config_Authentication); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_Config_Authorization); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_Config_Authorization_MatchSpec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicySpec_Config_Authorization_MatchSpec_RequestSpec_HeaderValues); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPolicyReport_IdentityReference); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_depIdxs,
		EnumInfos:         file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_enumTypes,
		MessageInfos:      file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto = out.File
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_rawDesc = nil
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_goTypes = nil
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_policy_v2_security_access_policy_proto_depIdxs = nil
}
