// Use Gloo Platform to install Istio ingress, egress, and east-west gateways in your workload clusters,
// as part of the Istio lifecycle management.
// In your `GatewayLifecycleManager` resource, you provide gateway settings in an `IstioOperator` configuration.
// When you create the `GatewayLifecycleManager` in your management cluster, Gloo translates the configuration
// into gateways in your registered workload clusters for you.
//
// For more information, see the [Install Istio by using the Istio Lifecycle Manager]({{% link path="/setup/install/gloo_mesh_managed/" %}}) guide.
//
// ## Examples
//
// **East-west**: This example creates an east-west gateway named `istio-eastwestgateway` in the `gloo-mesh-gateways`
// namespace of two workload clusters (`$REMOTE_CLUSTER1` and `$REMOTE_CLUSTER2`). You supply the repo key for the Solo distribution of Istio (`hub: $REPO`),
// image tag (`tag: $ISTIO_IMAGE`), and revision (`revision: $REVISION`). Note that for advanced east-west traffic routing across multiple clusters, you need a
// Gloo Mesh Enterprise license.
// ```yaml
// apiVersion: admin.gloo.solo.io/v2
// kind: GatewayLifecycleManager
// metadata:
//   name: istio-eastwestgateway
//   namespace: gloo-mesh
// spec:
//   installations:
//   # The revision for this installation
//   - gatewayRevision: $REVISION
//     # List all workload clusters to install Istio into
//     clusters:
//     - name: $REMOTE_CLUSTER1
//       # If set to true, the spec for this revision is applied in the cluster
//       activeGateway: true
//     - name: $REMOTE_CLUSTER2
//       activeGateway: true
//     istioOperatorSpec:
//       # No control plane components are installed
//       profile: empty
//       # Solo.io Istio distribution repository; required for the Solo distribution of Istio.
//       # You get the repo key from your Solo Account Representative.
//       hub: $REPO
//       # The Solo.io Gloo Istio version
//       tag: $ISTIO_IMAGE
//       components:
//         ingressGateways:
//         # Enable the default east-west gateway
//         - name: istio-eastwestgateway
//           # Deployed to gloo-mesh-gateways by default
//           namespace: gloo-mesh-gateways
//           enabled: true
//           label:
//             # Set a unique label for the gateway. This is required to
//             # ensure Gateways can select this workload.
//             istio: eastwestgateway
//             app: istio-eastwestgateway
//           k8s:
//             env:
//               # 'sni-dnat' enables AUTO_PASSTHROUGH mode for east-west communication through the gateway.
//               # The default value ('standard') does not set up a passthrough cluster.
//               # Required for multi-cluster communication and to preserve SNI.
//               - name: ISTIO_META_ROUTER_MODE
//                 value: "sni-dnat"
//             service:
//               type: LoadBalancer
//               selector:
//                 istio: eastwestgateway
//               # Default ports
//               ports:
//                 # Port for health checks on path /healthz/ready.
//                 # For AWS ELBs, this port must be listed first.
//                 - name: status-port
//                   port: 15021
//                   targetPort: 15021
//                 # Port for multicluster mTLS passthrough
//                 # Gloo looks for this default name 'tls' on a gateway
//                 # Required for Gloo east/west routing
//                 - name: tls
//                   port: 15443
//                   targetPort: 15443
//
// ```
//
// **Ingress**: This example creates an ingress gateway named `istio-ingressgateway` in the `gloo-mesh-gateways`
// namespace of two workload clusters (`$REMOTE_CLUSTER1` and `$REMOTE_CLUSTER2`). You supply the repo key for the Solo distribution of Istio (`hub: $REPO`),
// image tag (`tag: $ISTIO_IMAGE`), and revision (`revision: $REVISION`). Note that for advanced ingress routing features such as AWS Lambda, Portal, or ingress-specific
// policies, you need a Gloo Mesh Gateway license.
// ```yaml
// apiVersion: admin.gloo.solo.io/v2
// kind: GatewayLifecycleManager
// metadata:
//   name: istio-ingressgateway
//   namespace: gloo-mesh
// spec:
//   installations:
//   # The revision for this installation
//   - gatewayRevision: $REVISION
//     # List all workload clusters to install Istio into
//     clusters:
//     - name: $REMOTE_CLUSTER1
//       # If set to true, the spec for this revision is applied in the cluster
//       activeGateway: true
//     - name: $REMOTE_CLUSTER2
//       activeGateway: true
//     istioOperatorSpec:
//       # No control plane components are installed
//       profile: empty
//       # Solo.io Istio distribution repository; required for the Solo distribution of Istio.
//       # You get the repo key from your Solo Account Representative.
//       hub: $REPO
//       # Any tag for the Solo distribution of Istio
//       tag: $ISTIO_IMAGE
//       components:
//         ingressGateways:
//         # Enable the default ingress gateway
//         - name: istio-ingressgateway
//           # Deployed to gloo-mesh-gateways by default
//           namespace: gloo-mesh-gateways
//           enabled: true
//           label:
//             # Set a unique label for the gateway. This is required to
//             # ensure Gateways can select this workload
//             istio: ingressgateway
//             app: istio-ingressgateway
//           k8s:
//             service:
//               type: LoadBalancer
//               selector:
//                 istio: ingressgateway
//               # Default ports
//               ports:
//                 # Port for health checks on path /healthz/ready.
//                 # For AWS ELBs, this port must be listed first.
//                 - name: status-port
//                   port: 15021
//                   targetPort: 15021
//                 # Main HTTP ingress port
//                 - name: http2
//                   port: 80
//                   targetPort: 8080
//                 # Main HTTPS ingress port
//                 - name: https
//                   port: 443
//                   targetPort: 8443
//                 - name: tls
//                   port: 15443
//                   targetPort: 15443
// ```
//
// **Egress**: This example creates an egress gateway named `istio-egressgateway` in the `gloo-mesh-gateways` namespace of two workload clusters,
// (`$REMOTE_CLUSTER1` and `$REMOTE_CLUSTER2`). You supply the repo key for the Solo distribution of Istio (`hub: $REPO`),
// image tag (`tag: $ISTIO_IMAGE`), and revision (`revision: $REVISION`). For more information, see the
// [Block egress traffic with an egress gateway]({{% link path="/routing/forward-requests/external-service/egress-gateway/" %}}) guide.
// ```yaml
// apiVersion: admin.gloo.solo.io/v2
// kind: GatewayLifecycleManager
// metadata:
//   name: istio-egressgateway
//   namespace: gloo-mesh
// spec:
//   installations:
//       # The revision for this installation
//     - gatewayRevision: $REVISION
//       # List all workload clusters to install Istio into
//       clusters:
//       - name: $REMOTE_CLUSTER1
//         # If set to true, the spec for this revision is applied in the cluster
//         activeGateway: true
//       - name: $REMOTE_CLUSTER2
//         activeGateway: true
//       istioOperatorSpec:
//         # No control plane components are installed
//         profile: minimal
//         # Solo.io Istio distribution repository; required for Gloo Istio.
//         # You get the repo key from your Solo Account Representative.
//         hub: $REPO
//         # The Solo.io Gloo Istio version
//         tag: $ISTIO_IMAGE
//         meshConfig:
//           outboundTrafficPolicy:
//             mode: REGISTRY_ONLY
//             # Enable access logs
//           accessLogFile: /dev/stdout
//           defaultConfig:
//             proxyMetadata:
//               # For known hosts, enable the Istio agent to handle DNS requests
//               # for any custom ServiceEntry, such as non-Kubernetes services.
//               # Unknown hosts are automatically resolved using upstream DNS
//               # servers in resolv.conf (for proxy-dns)
//               ISTIO_META_DNS_CAPTURE: "true"
//         components:
//           egressGateways:
//           # Enable the egress gateway
//             - name: istio-egressgateway
//               # Deployed to gloo-mesh-gateways by default
//               namespace: gloo-mesh-gateways
//               enabled: true
//               label:
//                 # Set a unique label for the gateway. This is required to
//                 # ensure Gateways can select this workload.
//                 istio: egressgateway
//                 app: istio-egressgateway
//                 traffic: egress
//               k8s:
//                 affinity:
//                    nodeAffinity:
//                      requiredDuringSchedulingIgnoredDuringExecution:
//                        nodeSelectorTerms:
//                          - matchExpressions:
//                              - key: kubernetes.io/arch
//                                operator: In
//                                values:
//                                  - arm64
//                                  - amd64
//                 env:
//                   # 'sni-dnat' enables AUTO_PASSTHROUGH mode for east-west communication through the gateway.
//                   # The default value ('standard') does not set up a passthrough cluster.
//                   # Required for multi-cluster communication and to preserve SNI.
//                   - name: ISTIO_META_ROUTER_MODE
//                     value: "sni-dnat"
//                   - name: AUTO_RELOAD_PLUGIN_CERTS
//                     value: "true"
//                 podAnnotations:
//                   proxy.istio.io/config: |
//                     proxyStatsMatcher:
//                       inclusionRegexps:
//                       - .*ext_authz.*
//                 service:
//                   type: LoadBalancer
//                   selector:
//                     istio: egressgateway
//                   # Default ports
//                   ports:
//                     # Port for health checks on path /healthz/ready.
//                     # For AWS ELBs, this port must be listed first.
//                     - port: 15021
//                       targetPort: 15021
//                       name: status-port
//                     # Port for multicluster mTLS passthrough
//                     # Required for Gloo egress routing
//                     - port: 15443
//                       targetPort: 15443
//                       # Gloo looks for this default name 'tls' on a gateway
//                       name: tls
//                     # Required for Istio mutual TLS
//                     - port: 443
//                       targetPort: 8443
//                       name: https
// ```

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: github.com/solo-io/solo-apis/api/gloo.solo.io/admin/v2/gateway_lifecycle_manager.proto

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

// The current state of the gateway installation.
type GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State int32

const (
	// Waiting for resources to be installed or updated.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_PENDING GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 0
	// The Gloo management server encountered a problem while attempting
	// to install the gateway.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_FAILED GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 1
	// Could not select a istiod control plane.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_NO_CONTROL_PLANE_AVAILABLE GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 2
	// The gateway is currently being installed.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_INSTALLING_GATEWAY GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 3
	// All Istio components for the gateway are successfully installed and healthy.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_HEALTHY GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 4
	// The gateway installation is no longer healthy.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_UNHEALTHY GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 5
	// The gateway IstioOperator resource is in an 'ACTION_REQUIRED' state. Check the logs of the IstioOperator deployment for more info.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_ACTION_REQUIRED GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 6
	// The gateway IstioOperator resource is in an 'UPDATING' state.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_UPDATING_GATEWAY GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 7
	// The gateway IstioOperator resource is in a 'RECONCILING' state.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_RECONCILING_GATEWAY GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 8
	// The gateway installation state could not be determined.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_UNKNOWN GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 9
	// The gateway is currently being uninstalled.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_UNINSTALLING_GATEWAY GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 10
	// The gateway is uninstalled.
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_UNINSTALLED_GATEWAY GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 11
	// Successfully translated but not installing yet
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_INSTALL_PENDING GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State = 12
)

// Enum value maps for GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State.
var (
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State_name = map[int32]string{
		0:  "PENDING",
		1:  "FAILED",
		2:  "NO_CONTROL_PLANE_AVAILABLE",
		3:  "INSTALLING_GATEWAY",
		4:  "HEALTHY",
		5:  "UNHEALTHY",
		6:  "ACTION_REQUIRED",
		7:  "UPDATING_GATEWAY",
		8:  "RECONCILING_GATEWAY",
		9:  "UNKNOWN",
		10: "UNINSTALLING_GATEWAY",
		11: "UNINSTALLED_GATEWAY",
		12: "INSTALL_PENDING",
	}
	GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State_value = map[string]int32{
		"PENDING":                    0,
		"FAILED":                     1,
		"NO_CONTROL_PLANE_AVAILABLE": 2,
		"INSTALLING_GATEWAY":         3,
		"HEALTHY":                    4,
		"UNHEALTHY":                  5,
		"ACTION_REQUIRED":            6,
		"UPDATING_GATEWAY":           7,
		"RECONCILING_GATEWAY":        8,
		"UNKNOWN":                    9,
		"UNINSTALLING_GATEWAY":       10,
		"UNINSTALLED_GATEWAY":        11,
		"INSTALL_PENDING":            12,
	}
)

func (x GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State) Enum() *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State {
	p := new(GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State)
	*p = x
	return p
}

func (x GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_enumTypes[0].Descriptor()
}

func (GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_enumTypes[0]
}

func (x GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State.Descriptor instead.
func (GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{3, 1, 1, 0}
}

// Specifications for the `GatewayLifecycleManager` resource.
type GatewayLifecycleManagerSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of Istio gateway installations.
	Installations []*GatewayInstallation `protobuf:"bytes,1,rep,name=installations,proto3" json:"installations,omitempty"`
	// Optional default configuration applicable to all installations
	HelmGlobal *v2.IstioLifecycleHelmGlobals `protobuf:"bytes,2,opt,name=helm_global,json=helmGlobal,proto3" json:"helm_global,omitempty"`
}

func (x *GatewayLifecycleManagerSpec) Reset() {
	*x = GatewayLifecycleManagerSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayLifecycleManagerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayLifecycleManagerSpec) ProtoMessage() {}

func (x *GatewayLifecycleManagerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayLifecycleManagerSpec.ProtoReflect.Descriptor instead.
func (*GatewayLifecycleManagerSpec) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{0}
}

func (x *GatewayLifecycleManagerSpec) GetInstallations() []*GatewayInstallation {
	if x != nil {
		return x.Installations
	}
	return nil
}

func (x *GatewayLifecycleManagerSpec) GetHelmGlobal() *v2.IstioLifecycleHelmGlobals {
	if x != nil {
		return x.HelmGlobal
	}
	return nil
}

// Clusters to install the Istio gateways in.
type GatewayClusterSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the cluster to install the gateway into.
	// Must match the name of the cluster that you used when you registered the cluster with Gloo.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Optional: Defaults to false.
	// When set to true, the gateway installation for this revision is applied as the active gateway through which primary service traffic is routed in the cluster.
	// If the `istioOperatorSpec` defines a service, this field switches the service selectors to the revision specified in the `gatewayRevsion`.
	// You might change this setting for gateway installations during a canary upgrade.
	// For more info, see the [upgrade docs](https://docs.solo.io/gloo-mesh-enterprise/latest/istio/mesh/ilm-upgrade/).
	ActiveGateway bool `protobuf:"varint,2,opt,name=active_gateway,json=activeGateway,proto3" json:"active_gateway,omitempty"`
	// Optional: By default, the `trustDomain` value in the `meshConfig` section of the operator spec is automatically set by the Gloo to the name of each workload cluster.
	// To override the `trustDomain` for each cluster, you can instead specify the override value by using this `trustDomain` field,
	// and include the value in the list of cluster names. For example, if you specify `meshConfig.trustDomain: cluster1-trust-override` in your operator spec,
	// you then specify both the cluster name (`name: cluster1`) and the trust domain (`trustDomain: cluster1-trust-override`) in this `installations.clusters` section.
	// Additionally, because Gloo requires multiple trust domains for east-west routing, the `PILOT_SKIP_VALIDATE_TRUST_DOMAIN` field is set to `"true"` by default.
	// For more info, see the [Istio documentation](https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1).
	TrustDomain string `protobuf:"bytes,5,opt,name=trust_domain,json=trustDomain,proto3" json:"trust_domain,omitempty"`
}

func (x *GatewayClusterSelector) Reset() {
	*x = GatewayClusterSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayClusterSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayClusterSelector) ProtoMessage() {}

func (x *GatewayClusterSelector) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayClusterSelector.ProtoReflect.Descriptor instead.
func (*GatewayClusterSelector) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{1}
}

func (x *GatewayClusterSelector) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GatewayClusterSelector) GetActiveGateway() bool {
	if x != nil {
		return x.ActiveGateway
	}
	return false
}

func (x *GatewayClusterSelector) GetTrustDomain() string {
	if x != nil {
		return x.TrustDomain
	}
	return ""
}

// List of Istio gateway installations.
// Any components that are not related to the gateway are ignored.
// You can provide only one type of gateway installation per revision in a cluster.
// For example, in a workload cluster `cluster2`, you can install only one east-west
// gateway that runs revision `1-19-5`.
type GatewayInstallation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional: The revision of an Istio control plane in the cluster that this gateway should also use.
	// If a control plane installation of this revision is not found, no gateway is created.
	ControlPlaneRevision string `protobuf:"bytes,1,opt,name=control_plane_revision,json=controlPlaneRevision,proto3" json:"control_plane_revision,omitempty"`
	// Istio revision for this gateway installation.
	// When set to `auto`, Gloo installs the gateway with the default supported version of the Solo distribution of Istio.
	GatewayRevision string `protobuf:"bytes,2,opt,name=gateway_revision,json=gatewayRevision,proto3" json:"gateway_revision,omitempty"`
	// Clusters to install the Istio gateways in.
	Clusters []*GatewayClusterSelector `protobuf:"bytes,3,rep,name=clusters,proto3" json:"clusters,omitempty"`
	// IstioOperator specification for the gateway.
	// For more info, see the [Istio documentation](https://istio.io/latest/docs/reference/config/istio.operator.v1alpha1/).
	IstioOperatorSpec *v2.IstioOperatorSpec `protobuf:"bytes,4,opt,name=istio_operator_spec,json=istioOperatorSpec,proto3" json:"istio_operator_spec,omitempty"`
	// When set to true, the lifecycle manager allows you to perform
	// in-place upgrades by skipping checks that are required for canary upgrades.
	// In production environments, canary upgrades are recommended for
	// updating the minor version. To update the patch version or make
	// configuration changes within the same version, you can use in-place upgrades.
	// Be sure to test in-place upgrades in development or staging environments first.
	SkipUpgradeValidation bool `protobuf:"varint,5,opt,name=skip_upgrade_validation,json=skipUpgradeValidation,proto3" json:"skip_upgrade_validation,omitempty"`
}

func (x *GatewayInstallation) Reset() {
	*x = GatewayInstallation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayInstallation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayInstallation) ProtoMessage() {}

func (x *GatewayInstallation) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayInstallation.ProtoReflect.Descriptor instead.
func (*GatewayInstallation) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{2}
}

func (x *GatewayInstallation) GetControlPlaneRevision() string {
	if x != nil {
		return x.ControlPlaneRevision
	}
	return ""
}

func (x *GatewayInstallation) GetGatewayRevision() string {
	if x != nil {
		return x.GatewayRevision
	}
	return ""
}

func (x *GatewayInstallation) GetClusters() []*GatewayClusterSelector {
	if x != nil {
		return x.Clusters
	}
	return nil
}

func (x *GatewayInstallation) GetIstioOperatorSpec() *v2.IstioOperatorSpec {
	if x != nil {
		return x.IstioOperatorSpec
	}
	return nil
}

func (x *GatewayInstallation) GetSkipUpgradeValidation() bool {
	if x != nil {
		return x.SkipUpgradeValidation
	}
	return false
}

// The status of the `GatewayLifecycleManager` resource after you apply it to your Gloo environment.
type GatewayLifecycleManagerStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The list of clusters where Gloo manages Istio gateway installations.
	Clusters map[string]*GatewayLifecycleManagerStatus_ClusterStatuses `protobuf:"bytes,1,rep,name=clusters,proto3" json:"clusters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GatewayLifecycleManagerStatus) Reset() {
	*x = GatewayLifecycleManagerStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayLifecycleManagerStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayLifecycleManagerStatus) ProtoMessage() {}

func (x *GatewayLifecycleManagerStatus) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayLifecycleManagerStatus.ProtoReflect.Descriptor instead.
func (*GatewayLifecycleManagerStatus) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{3}
}

func (x *GatewayLifecycleManagerStatus) GetClusters() map[string]*GatewayLifecycleManagerStatus_ClusterStatuses {
	if x != nil {
		return x.Clusters
	}
	return nil
}

// $hide_from_docs
type GatewayLifecycleManagerNewStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GatewayLifecycleManagerNewStatus) Reset() {
	*x = GatewayLifecycleManagerNewStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayLifecycleManagerNewStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayLifecycleManagerNewStatus) ProtoMessage() {}

func (x *GatewayLifecycleManagerNewStatus) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayLifecycleManagerNewStatus.ProtoReflect.Descriptor instead.
func (*GatewayLifecycleManagerNewStatus) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{4}
}

// $hide_from_docs
type GatewayLifecycleManagerReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GatewayLifecycleManagerReport) Reset() {
	*x = GatewayLifecycleManagerReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayLifecycleManagerReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayLifecycleManagerReport) ProtoMessage() {}

func (x *GatewayLifecycleManagerReport) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayLifecycleManagerReport.ProtoReflect.Descriptor instead.
func (*GatewayLifecycleManagerReport) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{5}
}

// The list of clusters where Gloo manages Istio gateway installations.
type GatewayLifecycleManagerStatus_ClusterStatuses struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The Istio gateway installations in the cluster, listed by revision.
	Installations map[string]*GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus `protobuf:"bytes,1,rep,name=installations,proto3" json:"installations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses) Reset() {
	*x = GatewayLifecycleManagerStatus_ClusterStatuses{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayLifecycleManagerStatus_ClusterStatuses) ProtoMessage() {}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayLifecycleManagerStatus_ClusterStatuses.ProtoReflect.Descriptor instead.
func (*GatewayLifecycleManagerStatus_ClusterStatuses) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{3, 1}
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses) GetInstallations() map[string]*GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus {
	if x != nil {
		return x.Installations
	}
	return nil
}

// The status of the gateway installation.
type GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The current state of the gateway installation.
	State GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State `protobuf:"varint,1,opt,name=state,proto3,enum=admin.gloo.solo.io.GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State" json:"state,omitempty"`
	// A human-readable message about the current state of the installation.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// The observed revision of the gateway installation.
	ObservedRevision string `protobuf:"bytes,5,opt,name=observed_revision,json=observedRevision,proto3" json:"observed_revision,omitempty"`
	// The IstioOperator spec that is currently deployed for this revision.
	ObservedOperator *v2.IstioOperatorSpec `protobuf:"bytes,4,opt,name=observed_operator,json=observedOperator,proto3" json:"observed_operator,omitempty"`
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) Reset() {
	*x = GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) ProtoMessage() {}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus.ProtoReflect.Descriptor instead.
func (*GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP(), []int{3, 1, 1}
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) GetState() GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State {
	if x != nil {
		return x.State
	}
	return GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_PENDING
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) GetObservedRevision() string {
	if x != nil {
		return x.ObservedRevision
	}
	return ""
}

func (x *GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus) GetObservedOperator() *v2.IstioOperatorSpec {
	if x != nil {
		return x.ObservedOperator
	}
	return nil
}

var File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDesc = []byte{
	0x0a, 0x64, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x5f, 0x6c,
	0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x6c,
	0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x1a, 0x12, 0x65, 0x78, 0x74, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x5a,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d,
	0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74,
	0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x5f, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x56, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67,
	0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72,
	0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76,
	0x32, 0x2f, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x5f, 0x68, 0x65, 0x6c, 0x6d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xbd, 0x01, 0x0a, 0x1b, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x4c, 0x69,
	0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x70,
	0x65, 0x63, 0x12, 0x4d, 0x0a, 0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x4f, 0x0a, 0x0b, 0x68, 0x65, 0x6c, 0x6d, 0x5f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x49, 0x73, 0x74,
	0x69, 0x6f, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x48, 0x65, 0x6c, 0x6d, 0x47,
	0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x73, 0x52, 0x0a, 0x68, 0x65, 0x6c, 0x6d, 0x47, 0x6c, 0x6f, 0x62,
	0x61, 0x6c, 0x22, 0x76, 0x0a, 0x16, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x25, 0x0a, 0x0e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x72, 0x75, 0x73, 0x74,
	0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74,
	0x72, 0x75, 0x73, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0xce, 0x02, 0x0a, 0x13, 0x47,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x16, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x5f, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x14, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65,
	0x52, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x67, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x5f, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x52, 0x65, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x46, 0x0a, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x6c,
	0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x52, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x12, 0x56, 0x0a, 0x13, 0x69,
	0x73, 0x74, 0x69, 0x6f, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x70,
	0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x49,
	0x73, 0x74, 0x69, 0x6f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x70, 0x65, 0x63,
	0x52, 0x11, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x53,
	0x70, 0x65, 0x63, 0x12, 0x36, 0x0a, 0x17, 0x73, 0x6b, 0x69, 0x70, 0x5f, 0x75, 0x70, 0x67, 0x72,
	0x61, 0x64, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x15, 0x73, 0x6b, 0x69, 0x70, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64,
	0x65, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xda, 0x08, 0x0a, 0x1d,
	0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65,
	0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x5b, 0x0a,
	0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x3f, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c,
	0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x4c, 0x69, 0x66, 0x65,
	0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x7e, 0x0a, 0x0d, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x57, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x41, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69,
	0x6f, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63,
	0x6c, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0xdb, 0x06, 0x0a, 0x0f, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x12, 0x7a,
	0x0a, 0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x54, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x6c,
	0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0d, 0x69, 0x6e, 0x73,
	0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x96, 0x01, 0x0a, 0x12, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x6a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x54, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x4c,
	0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x1a, 0xb2, 0x04, 0x0a, 0x12, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x70, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x5a, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x64, 0x5f, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x10, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x52, 0x65, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x53, 0x0a, 0x11, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x5f,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c,
	0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x49, 0x73, 0x74, 0x69, 0x6f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x10, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x8d, 0x02, 0x0a, 0x05, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12,
	0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x4e,
	0x4f, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x4f, 0x4c, 0x5f, 0x50, 0x4c, 0x41, 0x4e, 0x45, 0x5f,
	0x41, 0x56, 0x41, 0x49, 0x4c, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x49,
	0x4e, 0x53, 0x54, 0x41, 0x4c, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x47, 0x41, 0x54, 0x45, 0x57, 0x41,
	0x59, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x04,
	0x12, 0x0d, 0x0a, 0x09, 0x55, 0x4e, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x05, 0x12,
	0x13, 0x0a, 0x0f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x49, 0x52,
	0x45, 0x44, 0x10, 0x06, 0x12, 0x14, 0x0a, 0x10, 0x55, 0x50, 0x44, 0x41, 0x54, 0x49, 0x4e, 0x47,
	0x5f, 0x47, 0x41, 0x54, 0x45, 0x57, 0x41, 0x59, 0x10, 0x07, 0x12, 0x17, 0x0a, 0x13, 0x52, 0x45,
	0x43, 0x4f, 0x4e, 0x43, 0x49, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x47, 0x41, 0x54, 0x45, 0x57, 0x41,
	0x59, 0x10, 0x08, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x09,
	0x12, 0x18, 0x0a, 0x14, 0x55, 0x4e, 0x49, 0x4e, 0x53, 0x54, 0x41, 0x4c, 0x4c, 0x49, 0x4e, 0x47,
	0x5f, 0x47, 0x41, 0x54, 0x45, 0x57, 0x41, 0x59, 0x10, 0x0a, 0x12, 0x17, 0x0a, 0x13, 0x55, 0x4e,
	0x49, 0x4e, 0x53, 0x54, 0x41, 0x4c, 0x4c, 0x45, 0x44, 0x5f, 0x47, 0x41, 0x54, 0x45, 0x57, 0x41,
	0x59, 0x10, 0x0b, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x4e, 0x53, 0x54, 0x41, 0x4c, 0x4c, 0x5f, 0x50,
	0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x0c, 0x22, 0x22, 0x0a, 0x20, 0x47, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x4d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x4e, 0x65, 0x77, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x1f, 0x0a, 0x1d,
	0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65,
	0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x56, 0x5a,
	0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f,
	0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e,
	0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x32, 0xc0, 0xf5, 0x04, 0x01, 0xb8, 0xf5, 0x04,
	0x01, 0xd0, 0xf5, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescData = file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDesc
)

func file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescData)
	})
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDescData
}

var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_goTypes = []interface{}{
	(GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus_State)(0), // 0: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationStatus.State
	(*GatewayLifecycleManagerSpec)(nil),                                         // 1: admin.gloo.solo.io.GatewayLifecycleManagerSpec
	(*GatewayClusterSelector)(nil),                                              // 2: admin.gloo.solo.io.GatewayClusterSelector
	(*GatewayInstallation)(nil),                                                 // 3: admin.gloo.solo.io.GatewayInstallation
	(*GatewayLifecycleManagerStatus)(nil),                                       // 4: admin.gloo.solo.io.GatewayLifecycleManagerStatus
	(*GatewayLifecycleManagerNewStatus)(nil),                                    // 5: admin.gloo.solo.io.GatewayLifecycleManagerNewStatus
	(*GatewayLifecycleManagerReport)(nil),                                       // 6: admin.gloo.solo.io.GatewayLifecycleManagerReport
	nil,                                                                         // 7: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClustersEntry
	(*GatewayLifecycleManagerStatus_ClusterStatuses)(nil),                       // 8: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses
	nil, // 9: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationsEntry
	(*GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus)(nil), // 10: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationStatus
	(*v2.IstioLifecycleHelmGlobals)(nil),                                     // 11: common.gloo.solo.io.IstioLifecycleHelmGlobals
	(*v2.IstioOperatorSpec)(nil),                                             // 12: common.gloo.solo.io.IstioOperatorSpec
}
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_depIdxs = []int32{
	3,  // 0: admin.gloo.solo.io.GatewayLifecycleManagerSpec.installations:type_name -> admin.gloo.solo.io.GatewayInstallation
	11, // 1: admin.gloo.solo.io.GatewayLifecycleManagerSpec.helm_global:type_name -> common.gloo.solo.io.IstioLifecycleHelmGlobals
	2,  // 2: admin.gloo.solo.io.GatewayInstallation.clusters:type_name -> admin.gloo.solo.io.GatewayClusterSelector
	12, // 3: admin.gloo.solo.io.GatewayInstallation.istio_operator_spec:type_name -> common.gloo.solo.io.IstioOperatorSpec
	7,  // 4: admin.gloo.solo.io.GatewayLifecycleManagerStatus.clusters:type_name -> admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClustersEntry
	8,  // 5: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClustersEntry.value:type_name -> admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses
	9,  // 6: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.installations:type_name -> admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationsEntry
	10, // 7: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationsEntry.value:type_name -> admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationStatus
	0,  // 8: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationStatus.state:type_name -> admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationStatus.State
	12, // 9: admin.gloo.solo.io.GatewayLifecycleManagerStatus.ClusterStatuses.InstallationStatus.observed_operator:type_name -> common.gloo.solo.io.IstioOperatorSpec
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_init()
}
func file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_init() {
	if File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayLifecycleManagerSpec); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayClusterSelector); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayInstallation); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayLifecycleManagerStatus); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayLifecycleManagerNewStatus); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayLifecycleManagerReport); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayLifecycleManagerStatus_ClusterStatuses); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayLifecycleManagerStatus_ClusterStatuses_InstallationStatus); i {
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
			RawDescriptor: file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_depIdxs,
		EnumInfos:         file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_enumTypes,
		MessageInfos:      file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto = out.File
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_rawDesc = nil
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_goTypes = nil
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_admin_v2_gateway_lifecycle_manager_proto_depIdxs = nil
}
