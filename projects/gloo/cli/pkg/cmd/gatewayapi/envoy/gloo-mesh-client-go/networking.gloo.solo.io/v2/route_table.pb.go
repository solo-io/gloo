// A `RouteTable` resource defines one or more hosts and a set of traffic route rules that describe how to handle traffic for these hosts.
// Route tables support two types of routes: HTTP and TCP.
//
// You can delegate HTTP routes to other route tables based on one or more matching hosts and specific route paths.
// If your "parent" route table delegates some traffic rules to another "child" route table, the child route table must be in the same workspace
// or imported to the parent route table's workspace.
//
// You can match traffic that originates from an ingress gateway (north-south), Istio mesh gateway (east-west),
// or directly from the sidecars of workloads in your service mesh (east-west),
// depending on the configuration of the `virtualGateways` field.
//
// For more information, see the [Routing overview concept docs]({{< link path="/traffic_management/concepts/routes/" >}}).
//
// ## Examples
//
// The following example defines route configuration for the 'uk.bookinfo.com' and 'eu.bookinfo.com' hosts.
// Traffic arrives at the `my-gateway` virtual gateway in the `my-gateway-ws` workspace.
// The route table sets up several different matchers to direct HTTP traffic.
// * When the cookie in the header matches to `user=dev-123`, HTTP traffic is forwarded to the port `7777` of the `v1` of `reviews.qa` service.
// * When the path matches exactly to `/reviews/` HTTP traffic is forwarded to port 9080 of the `reviews.qa` service.
// * All other HTTP traffic is sent to the default destination, which is port 9080 of `reviews.prod` service in the `bookinfo` workspace.
// ```yaml
// apiVersion: networking.gloo.solo.io/v2
// kind: RouteTable
// metadata:
//   name: bookinfo-root-routes
//   namespace: bookinfo
// spec:
//   hosts:
//     - 'uk.bookinfo.com'
//     - 'eu.bookinfo.com'
//   virtualGateways:
//     - name: my-gateway
//       namespace: my-gateway-ws
//   defaultDestination:
//     ref:
//       name: reviews
//       namespace: prod
//     port:
//       number: 9080
//   http:
//     - name: reviews-qa
//       matchers:
//         - headers:
//             - name: cookie
//               value: 'user=dev-123'
//       forwardTo:
//         destinations:
//           - ref:
//               name: reviews
//               namespace: qa
//             subset:
//               version: v1
//             port:
//               number: 7777
//     - name: reviews
//       matchers:
//         - name: review-prefix
//           uri:
//             exact: /reviews/
//       forwardTo:
//         destinations:
//           - ref:
//               name: reviews
//               namespace: qa
//             port:
//               number: 9080
// ```
//
// The following example defines route configuration for the 'example.com' host.
// Traffic arrives at the `istio-ingressgateway` virtual gateway in the `global` workspace.
// The route table sets up one matcher to direct HTTP traffic to multiple versions of one app.
// For example, say each version contains the `app: global-app` label, and labels such as
// `version: v1` and `version: v2`. The route table references both versions of the app,
// and specifies the version labels in the `subset` field of each reference.
// Additionally, this example specifies an optional destination `weight` for each version of the app.
// 75% of traffic requests to the `/global-app` are directed to `v1`, and 25% of traffic
// requests are directed to `v2`. Weighted routing can be useful in scenarios such as
// controlled rollouts to slowly move traffic from an older to a newer version of your app.
// ```yaml
// apiVersion: networking.gloo.solo.io/v2
// kind: RouteTable
// metadata:
//   name: global-app-routes
//   namespace: global
// spec:
//   hosts:
//     - example.com
//   # Selects the virtual gateway you previously created
//   virtualGateways:
//     - name: istio-ingressgateway
//       namespace: global
//   http:
//     # Route for the global-app service
//     - name: global-app
//       # Prefix matching
//       matchers:
//       - uri:
//           prefix: /global-app
//       # Forwarding directive
//       forwardTo:
//         destinations:
//           # Reference to Kubernetes service for version 1 of the app in this cluster
//           - ref:
//               name: global-app
//               namespace: global
//             port:
//               number: 9080
//             # Label for v1
//             subset:
//               version: v1
//             # 75% of request traffic to /global-app
//             weight: 75
//           # Reference to Kubernetes service for version 2 of the app in this cluster
//           - ref:
//               name: global-app
//               namespace: global
//             port:
//               number: 9080
//             # Label for v2
//             subset:
//               version: v2
//             # 25% of request traffic to /global-app
//             weight: 25
// ```
//
// **AWS Lambda examples**: For more information, see the [AWS Lambda integration in the Gloo Mesh Gateway docs](https://docs.solo.io/gloo-mesh-gateway/latest/lambda/lambda_routing/).
//
// The following example defines route configuration for the 'uk.bookinfo.com' and 'eu.bookinfo.com' hosts.
// Traffic arrives at the `my-gateway` virtual gateway in the `my-gateway-ws` workspace. The route table sends traffic to an external cloud function.
// * When the HTTP route path matches the prefix `/lambda`, traffic is forwarded to the backing `aws-provider` CloudProvider.
// * The associated `aws-provider` CloudResources resource describes an AWS Lambda service named `logicalName: aws-dest`.
// * The `"SYNC"` option indicates that the AWS Lambda function is invoked synchronously, which is also the default behavior.
//
// ```yaml
// apiVersion: networking.gloo.solo.io/v2
// kind: RouteTable
// metadata:
//   name: bookinfo-root-routes
//   namespace: bookinfo
// spec:
//   hosts:
//     - 'uk.bookinfo.com'
//     - 'eu.bookinfo.com'
//   virtualGateways:
//     - name: my-gateway
//       namespace: my-gateway-ws
//   defaultDestination:
//     ref:
//       name: reviews
//       namespace: prod
//     port:
//       number: 9080
//   http:
//     - name: lambda
//       matchers:
//         - uri:
//             prefix: /lambda
//       labels:
//         route: lambda
//       forwardTo:
//         destinations:
//           - awsLambda:
//               cloudProvider:
//                 name: aws-provider
//                 namespace: bookinfo
//                 cluster: cluster-1
//               function: aws-dest
//               options:
//                 invocationStyle: SYNC
// ```
//
// The following example defines route configuration for the 'uk.bookinfo.com' and 'eu.bookinfo.com' hosts.
// Traffic arrives at the `my-gateway` virtual gateway in the `my-gateway-ws` workspace. The route table sends traffic to an external cloud function.
// * When the HTTP route path matches the prefix `/lambda`, traffic is forwarded to the delegated route table for handling requests to AWS Lambdas.
// * The `allowedRoutes` restrict the usage of CloudProvider functionality, which routes to cloud functions `backend-function-*` in region `us-east-2` and which assumes the `dev-team-B-*` IAM role in AWS to invoke the function.
//
// ```yaml
// apiVersion: networking.gloo.solo.io/v2
// kind: RouteTable
// metadata:
//   name: bookinfo-root-routes
//   namespace: bookinfo
// spec:
//   hosts:
//     - 'uk.bookinfo.com'
//     - 'eu.bookinfo.com'
//   virtualGateways:
//     - name: my-gateway
//       namespace: my-gateway-ws
//   defaultDestination:
//     ref:
//       name: reviews
//       namespace: prod
//     port:
//       number: 9080
//   http:
//     - name: lambda
//       matchers:
//         - uri:
//             prefix: /lambda
//       labels:
//         route: lambda
//       delegate:
//         allowedRoutes:
//           - cloudProvider:
//               aws:
//                 lambda_function:
//                   - backend-function-.*
//                 iam_roles:
//                   - dev-team-B-.*
//                 regions:
//                   - us-east-2
//         routeTables:
//           - labels:
//               table: lambda
// ```
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: github.com/solo-io/solo-apis/api/gloo.solo.io/networking/v2/route_table.proto

package v2

import (
	reflect "reflect"
	sync "sync"

	v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	_ "github.com/solo-io/protoc-gen-ext/extproto"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	_ "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/envoy/gloo-mesh-client-go/apimanagement.gloo.solo.io/v2"
	v2 "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/envoy/gloo-mesh-client-go/common.gloo.solo.io/v2"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The desired behavior when one or more routes in the route table are misconfigured.
//
// </br>**Configuration constraints**: For delegated child route tables, this field must be empty or unset.
// This setting is supported only for the parent route table,
// which controls the behavior for each child route table.
type RouteTableSpec_FailureMode int32

const (
	// The default behavior for handling misconfigured routes in a route table.
	// If you attempt to apply incorrect configuration to a route, the configuration is not applied in the cluster.
	// Instead, the misconfigured route is replaced with a 500 HTTP direct response until the error is resolved.
	// Other, correctly configured routes in the route table continue to route as configured and accept updates to their configuration.
	RouteTableSpec_ROUTE_REPLACEMENT RouteTableSpec_FailureMode = 0
	// If you attempt to apply incorrect configuration to a route, the configuration is not accepted
	// and the route table continues to serve the route configuration from the last accepted configuration.
	// Other, correctly configured routes in the route table also continue to serve the last accepted configuration.
	// Updates to correctly configured routes are ignored until the error is resolved.
	// Note that the same behavior applies when using route delegation;
	// any misconfigured route on the parent route table or any child route table freezes the configuration for all routes in the route table tree until the error is resolved.
	// Keep in mind that if you change the failure mode from `ROUTE_REPLACEMENT` to `FREEZE_CONFIG` while a route is in a misconfigured state,
	// any replaced routes will maintain their 500 HTTP direct response as that is their behavior in the last accepted configuration.
	RouteTableSpec_FREEZE_CONFIG RouteTableSpec_FailureMode = 1
)

// Enum value maps for RouteTableSpec_FailureMode.
var (
	RouteTableSpec_FailureMode_name = map[int32]string{
		0: "ROUTE_REPLACEMENT",
		1: "FREEZE_CONFIG",
	}
	RouteTableSpec_FailureMode_value = map[string]int32{
		"ROUTE_REPLACEMENT": 0,
		"FREEZE_CONFIG":     1,
	}
)

func (x RouteTableSpec_FailureMode) Enum() *RouteTableSpec_FailureMode {
	p := new(RouteTableSpec_FailureMode)
	*p = x
	return p
}

func (x RouteTableSpec_FailureMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RouteTableSpec_FailureMode) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes[0].Descriptor()
}

func (RouteTableSpec_FailureMode) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes[0]
}

func (x RouteTableSpec_FailureMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RouteTableSpec_FailureMode.Descriptor instead.
func (RouteTableSpec_FailureMode) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{0, 0}
}

type RedirectAction_RedirectResponseCode int32

const (
	// Moved Permanently HTTP Status Code - 301.
	RedirectAction_MOVED_PERMANENTLY RedirectAction_RedirectResponseCode = 0
	// Found HTTP Status Code - 302.
	RedirectAction_FOUND RedirectAction_RedirectResponseCode = 1
	// See Other HTTP Status Code - 303.
	RedirectAction_SEE_OTHER RedirectAction_RedirectResponseCode = 2
	// Temporary Redirect HTTP Status Code - 307.
	RedirectAction_TEMPORARY_REDIRECT RedirectAction_RedirectResponseCode = 3
	// Permanent Redirect HTTP Status Code - 308.
	RedirectAction_PERMANENT_REDIRECT RedirectAction_RedirectResponseCode = 4
)

// Enum value maps for RedirectAction_RedirectResponseCode.
var (
	RedirectAction_RedirectResponseCode_name = map[int32]string{
		0: "MOVED_PERMANENTLY",
		1: "FOUND",
		2: "SEE_OTHER",
		3: "TEMPORARY_REDIRECT",
		4: "PERMANENT_REDIRECT",
	}
	RedirectAction_RedirectResponseCode_value = map[string]int32{
		"MOVED_PERMANENTLY":  0,
		"FOUND":              1,
		"SEE_OTHER":          2,
		"TEMPORARY_REDIRECT": 3,
		"PERMANENT_REDIRECT": 4,
	}
)

func (x RedirectAction_RedirectResponseCode) Enum() *RedirectAction_RedirectResponseCode {
	p := new(RedirectAction_RedirectResponseCode)
	*p = x
	return p
}

func (x RedirectAction_RedirectResponseCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RedirectAction_RedirectResponseCode) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes[1].Descriptor()
}

func (RedirectAction_RedirectResponseCode) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes[1]
}

func (x RedirectAction_RedirectResponseCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RedirectAction_RedirectResponseCode.Descriptor instead.
func (RedirectAction_RedirectResponseCode) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{6, 0}
}

// The method by which routes across delegated route tables are sorted.
type DelegateAction_SortMethod int32

const (
	// Routes are kept in the order that they appear relative to their tables, but tables are sorted by weight.
	// Tables that have the same weight stay in the same order that they are listed in, which is the list
	// order when given as a reference and by creation timestamp when selected.
	DelegateAction_TABLE_WEIGHT DelegateAction_SortMethod = 0
	// After processing all routes, including additional route tables delegated to, the resulting routes are sorted
	// by specificity to reduce the chance that a more specific route will be short-circuited by a general route.
	// Matchers with exact path matchers are considered more specific than regex path patchers, which are more
	// specific than prefix path matchers. For prefix and exact, matchers of the same type are sorted by length of the path in descending
	// order. For regex matchers they are all treated equal when sorted. For sort ties, table weights are used across tables &
	// within tables user specified order is preserved. Only the most specific matcher on each route is used.
	//
	// For example, consider the following two sub-tables that are sorted by specificity and the resulting route list.
	//
	// Sub-table A, with a table weight of `1` in case of sort ties:<ul>
	// <li>`prefix: /foo`</li>
	// <li>`prefix: /foo/more/specific`</li>
	// <li>`prefix: /foo/even/more/specific`</li>
	// <li>`exact: /foo/exact`</li>
	// <li>`exact: /foo/another/exact`</li>
	// <li>`regex: /foo/*`</li>
	// <li>`regex: /fooo/*`</li></ul>
	// Sub-table B, with a table weight of `2` in case of sort ties:<ul>
	// <li>`prefix: /bar`</li>
	// <li>`prefix: /bar/more/specific`</li>
	// <li>`prefix: /bar/even/more/specific`</li>
	// <li>`exact: /bar/exact`</li>
	// <li>`regex: /bar/*`</li></ul>
	// The resulting routes are sorted in this order:<ul>
	// <li>`exact: /foo/another/exact`</li>
	// <li>`exact: /bar/exact`</li>
	// <li>`exact: /foo/exact`</li>
	// <li>`regex: /bar/*`</li>
	// <li>`regex: /foo/*`</li>
	// <li>`regex: /fooo/*`</li>
	// <li>`prefix: /bar/even/more/specific`</li>
	// <li>`prefix: /foo/even/more/specific`</li>
	// <li>`prefix: /bar/more/specific`</li>
	// <li>`prefix: /foo/more/specific`</li>
	// <li>`prefix: /bar`</li>
	// <li>`prefix: /foo`</li></ul>
	DelegateAction_ROUTE_SPECIFICITY DelegateAction_SortMethod = 1
)

// Enum value maps for DelegateAction_SortMethod.
var (
	DelegateAction_SortMethod_name = map[int32]string{
		0: "TABLE_WEIGHT",
		1: "ROUTE_SPECIFICITY",
	}
	DelegateAction_SortMethod_value = map[string]int32{
		"TABLE_WEIGHT":      0,
		"ROUTE_SPECIFICITY": 1,
	}
)

func (x DelegateAction_SortMethod) Enum() *DelegateAction_SortMethod {
	p := new(DelegateAction_SortMethod)
	*p = x
	return p
}

func (x DelegateAction_SortMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DelegateAction_SortMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes[2].Descriptor()
}

func (DelegateAction_SortMethod) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes[2]
}

func (x DelegateAction_SortMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DelegateAction_SortMethod.Descriptor instead.
func (DelegateAction_SortMethod) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{8, 0}
}

// Specifications for the `RouteTable` resource.
//
// +kubebuilder:validation:XValidation:rule="has(self.hosts) || !has(self.virtualGateways)",message="The virtualGateways field must be empty or unset for delegated RouteTables."
// +kubebuilder:validation:XValidation:rule="has(self.hosts) || !has(self.applyToDestinations)",message="The applyToDestinations field must be empty or unset for delegated RouteTables."
// +kubebuilder:validation:XValidation:rule="has(self.hosts) || !has(self.failureMode)",message="The failureMode field must be empty or unset for delegated RouteTables."
// +kubebuilder:validation:XValidation:rule="has(self.http) || has(self.tcp) || has(self.tls)",message="At least one of http, tcp, or tls keys must be set."
// +kubebuilder:validation:XValidation:rule="has(self.defaultDestination) ? has(self.http) && self.http.exists(r, has(r.forwardTo)) : true",message="defaultDestination can only be set for http routes with forwardTo actions."
// +kubebuilder:validation:XValidation:rule="has(self.http) && self.http.exists(r, has(r.forwardTo) && has(r.forwardTo.destinations) && r.forwardTo.destinations.size() == 0) ? has(self.defaultDestination) : true",message="defaultDestination must be set when http forwardTo actions without destinations are specified."
type RouteTableSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional: One or more hosts for which this route table routes traffic.
	// To avoid potential misconfigurations, fully
	// qualified domain names are recommended instead of short names.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>For regular (non-delegated) route tables, this field is required and must specify at least one host.</li>
	// <li>For delegated child route tables, this field must be empty or unset.</li>
	// <li>Wildcards (`*`) are supported only for the left-most segment. Full wildcards (`"*"`) are not supported. For example,
	// `*.foo.com`, `*bar.foo.com`, and `*-bar.foo.com` are valid; `bar.*.com`, `bar*.foo.com`, `bar.foo.*`, and `*` are invalid.</li>
	// <li>Each hostname must follow these requirements:<ul>
	//
	//	<li>Hostnames must be 1 - 255 characters in length.</li>
	//	<li>The hostname cannot be an empty string.</li>
	//	<li>Supported characters are `a-z`, `A-Z`, `0-9`, `-`, and `.`.</li>
	//	<li>Each segment separated by a period (`.`) must be 1 - 63 characters in length and cannot start with the `-` character.</li>
	//	<li>Each segment must meet the regex `^[a-zA-Z0-9](?:[-a-zA-Z0-9]*[a-zA-Z0-9])?$`.</li>
	//	<li>The top-level domain (last segment), such as `com` in `www.example.com`, cannot be all numeric characters.</li>
	//	<li>The top-level domain (last segment) can be empty, such as `"istio.io."`.</li></ul></li></ul>
	Hosts []string `protobuf:"bytes,1,rep,name=hosts,proto3" json:"hosts,omitempty"`
	// Optional: A list of virtual gateways that serve this route table.
	// When not specified, the route table applies either to all sidecars in the workspace
	// or only to sidecars for selected workloads (via the `workloadSelectors` field) in the workspace where
	// the route table is deployed or imported.
	//
	// </br>**Examples**:<ul>
	// <li>The following applies to sidecars of all the workloads for the workspace where the route table is
	// deployed or imported: set `virtualGateways` to `null` and `workloadSelectors` to `[]`.</li>
	// <li>The following applies to the `my-gateway` virtual gateway in the `gateway` workspace and
	// no sidecars: set `virtualGateways.name` to `my-gateway`, `virtualGateways.namespace` to `gateway`, and `workloadSelectors` to `[]`.</li>
	// <li>The following applies to the `my-gateway` virtual gateway in the `gateway` workspace and
	// sidecars of all the workloads for the workspace where the route table is
	// deployed or imported: set `virtualGateways.name` to `my-gateway`, `virtualGateways.namespace` to `gateway`, and `workloadSelectors` to `{}`.</li>
	// <li>The following applies to sidecars of all the `app: foo` workloads for the workspace where the route table
	// is deployed or imported: set `virtualGateways` to `null` and `workloadSelectors.selector.labels` to `app: foo`.</li>
	// <li>The following applies to the `my-gateway` virtual gateway in the `gateway` workspace and
	// sidecars of all the `app: foo` workloads for the workspace where the route table is deployed or imported:
	// set `virtualGateways.name` to `my-gateway`, `virtualGateways.namespace` to `gateway`, and `workloadSelectors.selector.labels` to `app: foo`.</li></ul>
	//
	// </br>**Configuration constraints**: For delegated child route tables, this field must be empty or unset.
	// This setting is supported only for the parent route table,
	// which controls the behavior for each child route table.
	VirtualGateways []*v2.ObjectReference `protobuf:"bytes,5,rep,name=virtual_gateways,json=virtualGateways,proto3" json:"virtual_gateways,omitempty"`
	// Optional: Selectors for source workloads (with sidecars) that route traffic for this route table.
	//
	// </br>**Implementation notes**:<ul>
	// <li>You can specify workload selectors only for east-west routing. Workload selectors do not apply to ingress (north-south) routing.</li>
	// <li>If do not specify virtual gateways or workload selectors, all workloads in the workspace are selected for east-west routing by default.</li>
	// <li>If you specify a virtual gateway, workloads are not automatically selected for east-west routing. To make them also available for east-west routing, configure this `workloadSelectors` option to select the workloads that you want, such as `{}` for all workloads in the workspace.</li>
	// <li>Selecting external workloads, such as VMs, is currently not supported.</li>
	// <li>You can select workloads by using labels only. Selecting workloads by using other references, such as the name, namespace, cluster or workspace, is not supported.</li>
	// <li>Delegated child route tables inherit the workload selectors of the parent route table, such as `value:foo`.
	// The delegated child route table can also have its own workload selectors, such as `env:prod`.
	// These workload selectors are logically AND'd together. As a result, the child route table routes traffic only to workloads
	// with both `value:foo` and `env:prod` labels. Note that the child route table cannot override the parent's workload selectors,
	// such as by setting `value:bar`. In that case, the child route gets an error until the conflict is resolved.</li></ul>
	//
	// **Configuration constraints**: If this field is set, `workloadSelectors.kind` must be set to `KUBE`,
	// and `workloadSelectors.selector.name`, `.namespace`, `.cluster`, and `.workspace` must be empty.
	//
	// +kubebuilder:validation:MaxItems=10000
	// +kubebuilder:validation:XValidation:rule="self.all(s, !has(s.kind) || s.kind == 'KUBE')",message="Selector kind must be KUBE or not set."
	WorkloadSelectors []*v2.WorkloadSelector `protobuf:"bytes,6,rep,name=workload_selectors,json=workloadSelectors,proto3" json:"workload_selectors,omitempty"`
	// Optional: Selectors for destinations that route traffic for this route table via a producer-side policy, such as on waypoint proxies.
	//
	// </br>**Implementation notes**: Selecting external workloads (such as VMs), external services, or destinations with sidecars is currently not supported.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>If this field is set, `applyToDestinations.kind` must be set to `KUBE`.</li>
	// <li>For delegated child route tables, this field must be empty or unset.
	// The values from the parent route table are always used for destination selection.</li></ul>
	//
	// +kubebuilder:validation:XValidation:rule="self.all(s, !has(s.kind) || s.kind == 'KUBE')",message="Selector kind must be KUBE or not set."
	ApplyToDestinations []*v2.DestinationSelector `protobuf:"bytes,10,rep,name=apply_to_destinations,json=applyToDestinations,proto3" json:"apply_to_destinations,omitempty"`
	// Optional: Routes that do not specify a destination forward traffic to this destination.
	// This field applies only to `forwardTo` routes.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>If you define a `http.forwardTo`, `tcp.forwardTo`, or `tls.forwardTo` action that does not specify at least one destination, you must set this field.</li>
	// <li>The `subset` must not be set to the empty object `{}`.</li></ul>
	//
	// +kubebuilder:validation:XValidation:rule="!has(self.subset) || self.subset.size() > 0",message="subset must not be an empty map."
	DefaultDestination *v2.DestinationReference `protobuf:"bytes,2,opt,name=default_destination,json=defaultDestination,proto3" json:"default_destination,omitempty"`
	// The HTTP routes that this route table serves. If no routes match the client request,
	// the client receives a 404 error code. For more information on supported HTTP features, see the
	// [Routing overview concept docs]({{< link path="/traffic_management/concepts/routes/" >}}).
	//
	// </br>**Configuration constraints**: At least one of `http`, `tcp`, or `tls` must be set.
	//
	// +kubebuilder:validation:MaxItems=10000
	Http []*HTTPRoute `protobuf:"bytes,3,rep,name=http,proto3" json:"http,omitempty"`
	// The TCP routes that this route table serves. TCP routes are available only for internal
	// traffic within the cluster, not for ingress gateway traffic. For more information on supported
	// TCP features, see the [Routing overview concept docs]({{< link path="/traffic_management/concepts/routes/" >}}).
	//
	// </br>**Configuration constraints**: At least one of `http`, `tcp`, or `tls` must be set.
	//
	// +kubebuilder:validation:MaxItems=10000
	Tcp []*TCPRoute `protobuf:"bytes,8,rep,name=tcp,proto3" json:"tcp,omitempty"`
	// The TLS routes that this route table serves. For more information on supported
	// TLS features, see the [Routing overview concept docs]({{< link path="/traffic_management/concepts/routes/" >}}).
	//
	// </br>**Configuration constraints**: At least one of `http`, `tcp`, or `tls` must be set.
	//
	// +kubebuilder:validation:MaxItems=10000
	Tls []*TLSRoute `protobuf:"bytes,9,rep,name=tls,proto3" json:"tls,omitempty"`
	// Weight is used to sort delegated route tables by priority.
	// Higher integer values indicate a higher priority.
	// Individual routes are kept in the order that they appear relative to their tables,
	// but tables are sorted by the weight that you assign to them.
	//
	// <br>When a request is sent to a service in your Gloo setup,
	// the request is matched against the routes in the highest-weighted route table first.
	// If the request doesn’t match a route in the first sub-table,
	// it is matched against the routes in the second-highest-weighted table, and so on.
	//
	// <br>For example, if you have two sub-tables with weights of 100 and 90,
	// Gloo will attempt to match a request against the routes in the sub-table with the weight of 100
	// first. If the request does not match, Gloo then attempts to match the request
	// against the routes in the sub-table with the weight of 90.
	//
	// <br>Note that tables of the same weight stay in the same order that you list them in the parent route table,
	// which is the list order when you specify sub-tables by name, or the creation timestamp when you select sub-tables by label.
	//
	// <br>The default value is 0.
	Weight int32 `protobuf:"varint,4,opt,name=weight,proto3" json:"weight,omitempty"`
	// Optional: If this route table bundles APIs that you want to expose in a developer portal, you can set portal metadata.
	// Portal metadata is a set of key-value pairs that describe your APIs.
	// Later, your developer portal displays this information in the end-user facing API documentation.
	PortalMetadata *v2.PortalMetadata `protobuf:"bytes,7,opt,name=portal_metadata,json=portalMetadata,proto3" json:"portal_metadata,omitempty"`
	// The desired behavior when one or more routes in the route table are misconfigured.
	//
	// </br>**Configuration constraints**: For delegated child route tables, this field must be empty or unset.
	// This setting is supported only for the parent route table,
	// which controls the behavior for each child route table.
	FailureMode RouteTableSpec_FailureMode `protobuf:"varint,11,opt,name=failure_mode,json=failureMode,proto3,enum=networking.gloo.solo.io.RouteTableSpec_FailureMode" json:"failure_mode,omitempty"`
	// Annotations to add to the metadata of the VirtualService generated by this RouteTable.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>For delegated child route tables, this field must be empty or unset.
	// This setting is supported only for the parent route table,
	// which controls the VirtualService that corresponds to all child route tables.</li>
	// <li>Designated annotation keys that might already be in use on the VirtualService, such as "cluster.solo.io/cluster", are not supported.</li></ul>
	VirtualServiceAnnotations map[string]string `protobuf:"bytes,12,rep,name=virtual_service_annotations,json=virtualServiceAnnotations,proto3" json:"virtual_service_annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *RouteTableSpec) Reset() {
	*x = RouteTableSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteTableSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteTableSpec) ProtoMessage() {}

func (x *RouteTableSpec) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteTableSpec.ProtoReflect.Descriptor instead.
func (*RouteTableSpec) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{0}
}

func (x *RouteTableSpec) GetHosts() []string {
	if x != nil {
		return x.Hosts
	}
	return nil
}

func (x *RouteTableSpec) GetVirtualGateways() []*v2.ObjectReference {
	if x != nil {
		return x.VirtualGateways
	}
	return nil
}

func (x *RouteTableSpec) GetWorkloadSelectors() []*v2.WorkloadSelector {
	if x != nil {
		return x.WorkloadSelectors
	}
	return nil
}

func (x *RouteTableSpec) GetApplyToDestinations() []*v2.DestinationSelector {
	if x != nil {
		return x.ApplyToDestinations
	}
	return nil
}

func (x *RouteTableSpec) GetDefaultDestination() *v2.DestinationReference {
	if x != nil {
		return x.DefaultDestination
	}
	return nil
}

func (x *RouteTableSpec) GetHttp() []*HTTPRoute {
	if x != nil {
		return x.Http
	}
	return nil
}

func (x *RouteTableSpec) GetTcp() []*TCPRoute {
	if x != nil {
		return x.Tcp
	}
	return nil
}

func (x *RouteTableSpec) GetTls() []*TLSRoute {
	if x != nil {
		return x.Tls
	}
	return nil
}

func (x *RouteTableSpec) GetWeight() int32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *RouteTableSpec) GetPortalMetadata() *v2.PortalMetadata {
	if x != nil {
		return x.PortalMetadata
	}
	return nil
}

func (x *RouteTableSpec) GetFailureMode() RouteTableSpec_FailureMode {
	if x != nil {
		return x.FailureMode
	}
	return RouteTableSpec_ROUTE_REPLACEMENT
}

func (x *RouteTableSpec) GetVirtualServiceAnnotations() map[string]string {
	if x != nil {
		return x.VirtualServiceAnnotations
	}
	return nil
}

// Use HTTP routes to control Layer 7 application level traffic to your services. To configure HTTP routes, you pair together
// HTTP request `matchers` with certain actions. Matchers are criteria such as a route name, port, header, or method to match
// with an incoming request. Actions describe what to do with a matching request, such as `forwardTo` a destination or `delegate`
// to another route table. You can add metadata such as names and labels to your HTTP routes so that you can apply policies,
// track metrics, and better manage the routes.
type HTTPRoute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique name of the route within the route table. This name is used to identify the route in metrics collection.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>If the value begins with the prefix `insecure-`, this prefix is trimmed.</li>
	// <li>The value must be unique to other routes that are listed in the `http` section of this route table.
	// If it is not unique, it is renamed to `duplicate-<previous-name>-<increment>`,
	// such as `duplicate-myname-1`.</li></ul>
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Labels for the route, which you can use to apply policies that support `routeSelectors`.
	//
	// For enhanced security, include the special label `gateway.gloo.solo.io/require_auth=true`
	// on the route. To activate this security feature, enable the `gatewayDefaultDenyAllHTTPRequests`
	// feature flag for your Gloo installation. When both the label and feature flag are in place, Gloo
	// requires an authentication policy, such as ExtAuthPolicy or JWTPolicy, to be applied to the route.
	// If the authentication policy is removed or has an error, Gloo rejects all requests to the route.
	//
	// For more information about the value format, see
	// [Syntax and character set in the Kubernetes docs](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set).
	//
	// </br>**Configuration constraints**:<ul>
	// <li>Key constraints:<ul>
	//
	//	<li>Cannot be empty</li>
	//	<li>Must have two segments separated by a slash (`/`)</li>
	//	<li>First segment constraints:<ul>
	//	  <li>Cannot be empty</li>
	//	  <li>Max length of 253 characters</li>
	//	  <li>Supported characters include `a-z`, `A-Z`, `0-9`, `-`, and `.`</li></ul></li>
	//	<li>Second segment constraints:<ul>
	//	  <li>Cannot be empty</li>
	//	  <li>Max length of 63 characters</li>
	//	  <li>Must begin and end with an alphanumeric character (`a-z`, `A-Z`, or `0-9`)</li>
	//	  <li>Supported characters include `a-z`, `A-Z`, `0-9`, `-`, `_`, and `.`</li></ul></li></ul></li>
	//
	// <li>Value constraints:<ul>
	//
	//	<li>Can be empty</li>
	//	<li>Max length of 63 characters</li>
	//	<li>Unless empty, must begin and end with an alphanumeric character (`a-z`, `A-Z`, or `0-9`)</li>
	//	<li>Supported characters include `a-z`, `A-Z`, `0-9`, `-`, `_`, and `.`</li></ul></li></ul>
	Labels map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Request matchers that this route matches on. If none are specified, the route matches any HTTP traffic.
	// For delegated child route tables, this route matches only traffic that includes both the parent and child's matchers.
	// If these matchers conflict, the delegating route on the parent table is replaced with a `directResponse` that indicates the misconfiguration.
	Matchers []*v2.HTTPRequestMatcher `protobuf:"bytes,3,rep,name=matchers,proto3" json:"matchers,omitempty"`
	// The action to take when a request matches this route.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>This field is required.</li>
	// <li>Exactly one action type can be specified per route.</li></ul>
	//
	// Types that are assignable to ActionType:
	//
	//	*HTTPRoute_ForwardTo
	//	*HTTPRoute_Delegate
	//	*HTTPRoute_Redirect
	//	*HTTPRoute_DirectResponse
	//	*HTTPRoute_Graphql
	ActionType isHTTPRoute_ActionType `protobuf_oneof:"action_type"`
}

func (x *HTTPRoute) Reset() {
	*x = HTTPRoute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HTTPRoute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HTTPRoute) ProtoMessage() {}

func (x *HTTPRoute) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HTTPRoute.ProtoReflect.Descriptor instead.
func (*HTTPRoute) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{1}
}

func (x *HTTPRoute) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *HTTPRoute) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *HTTPRoute) GetMatchers() []*v2.HTTPRequestMatcher {
	if x != nil {
		return x.Matchers
	}
	return nil
}

func (m *HTTPRoute) GetActionType() isHTTPRoute_ActionType {
	if m != nil {
		return m.ActionType
	}
	return nil
}

func (x *HTTPRoute) GetForwardTo() *ForwardToAction {
	if x, ok := x.GetActionType().(*HTTPRoute_ForwardTo); ok {
		return x.ForwardTo
	}
	return nil
}

func (x *HTTPRoute) GetDelegate() *DelegateAction {
	if x, ok := x.GetActionType().(*HTTPRoute_Delegate); ok {
		return x.Delegate
	}
	return nil
}

func (x *HTTPRoute) GetRedirect() *RedirectAction {
	if x, ok := x.GetActionType().(*HTTPRoute_Redirect); ok {
		return x.Redirect
	}
	return nil
}

func (x *HTTPRoute) GetDirectResponse() *DirectResponseAction {
	if x, ok := x.GetActionType().(*HTTPRoute_DirectResponse); ok {
		return x.DirectResponse
	}
	return nil
}

func (x *HTTPRoute) GetGraphql() *GraphQLAction {
	if x, ok := x.GetActionType().(*HTTPRoute_Graphql); ok {
		return x.Graphql
	}
	return nil
}

type isHTTPRoute_ActionType interface {
	isHTTPRoute_ActionType()
}

type HTTPRoute_ForwardTo struct {
	// Forward traffic to one or more destination services.
	ForwardTo *ForwardToAction `protobuf:"bytes,4,opt,name=forward_to,json=forwardTo,proto3,oneof"`
}

type HTTPRoute_Delegate struct {
	// Delegate routing decisions to one or more HTTP route tables.
	Delegate *DelegateAction `protobuf:"bytes,5,opt,name=delegate,proto3,oneof"`
}

type HTTPRoute_Redirect struct {
	// Return a redirect response to the downstream client.
	Redirect *RedirectAction `protobuf:"bytes,6,opt,name=redirect,proto3,oneof"`
}

type HTTPRoute_DirectResponse struct {
	// Respond directly to the client from the proxy.
	DirectResponse *DirectResponseAction `protobuf:"bytes,7,opt,name=direct_response,json=directResponse,proto3,oneof"`
}

type HTTPRoute_Graphql struct {
	// Handle the HTTP request as a GraphQL request, including query validation and execution of the GraphQL request.
	Graphql *GraphQLAction `protobuf:"bytes,8,opt,name=graphql,proto3,oneof"`
}

func (*HTTPRoute_ForwardTo) isHTTPRoute_ActionType() {}

func (*HTTPRoute_Delegate) isHTTPRoute_ActionType() {}

func (*HTTPRoute_Redirect) isHTTPRoute_ActionType() {}

func (*HTTPRoute_DirectResponse) isHTTPRoute_ActionType() {}

func (*HTTPRoute_Graphql) isHTTPRoute_ActionType() {}

// Use TCP routes to control lower-level, connection-based traffic to services such as a local database.
// TCP routes are available only for internal traffic within the cluster, not for ingress gateway traffic.
// To configure TCP routes, you pair together TCP request `matchers` with certain actions.
// Matchers are criteria, such as a port, to match with an incoming request.
// Actions describe what to do with a matching request, such as `forwardTo` a destination.
type TCPRoute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The set of request matchers for this route to match on.
	Matchers []*v2.TCPRequestMatcher `protobuf:"bytes,1,rep,name=matchers,proto3" json:"matchers,omitempty"`
	// The action to take when a request matches this route.
	//
	// Types that are assignable to ActionType:
	//
	//	*TCPRoute_ForwardTo
	ActionType isTCPRoute_ActionType `protobuf_oneof:"action_type"`
}

func (x *TCPRoute) Reset() {
	*x = TCPRoute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TCPRoute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TCPRoute) ProtoMessage() {}

func (x *TCPRoute) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TCPRoute.ProtoReflect.Descriptor instead.
func (*TCPRoute) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{2}
}

func (x *TCPRoute) GetMatchers() []*v2.TCPRequestMatcher {
	if x != nil {
		return x.Matchers
	}
	return nil
}

func (m *TCPRoute) GetActionType() isTCPRoute_ActionType {
	if m != nil {
		return m.ActionType
	}
	return nil
}

func (x *TCPRoute) GetForwardTo() *ForwardToAction {
	if x, ok := x.GetActionType().(*TCPRoute_ForwardTo); ok {
		return x.ForwardTo
	}
	return nil
}

type isTCPRoute_ActionType interface {
	isTCPRoute_ActionType()
}

type TCPRoute_ForwardTo struct {
	// Forward traffic to one or more destination services. Note that some `forwardTo` actions, such as path or host rewrite, are not
	// supported for TCP routes.
	//
	// </br>**Configuration constraints**: This field is required, and you must specify at least one destination.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="has(self.destinations) && 1 <= self.destinations.size() && self.destinations.size() <= 100",message="At least one destination must be specified."
	ForwardTo *ForwardToAction `protobuf:"bytes,2,opt,name=forward_to,json=forwardTo,proto3,oneof"`
}

func (*TCPRoute_ForwardTo) isTCPRoute_ActionType() {}

// Use TLS routes to route unterminated TLS traffic (TLS/HTTPS) through an ingress gateway or within the cluster, such as for pass-through SNI-routing.
// You must specify an SNI host in the matcher, and optionally a port on the host.
type TLSRoute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The set of request matchers for this route to match on.
	Matchers []*v2.TLSRequestMatcher `protobuf:"bytes,1,rep,name=matchers,proto3" json:"matchers,omitempty"`
	// The action to take when a request matches this route.
	//
	// Types that are assignable to ActionType:
	//
	//	*TLSRoute_ForwardTo
	ActionType isTLSRoute_ActionType `protobuf_oneof:"action_type"`
}

func (x *TLSRoute) Reset() {
	*x = TLSRoute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TLSRoute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSRoute) ProtoMessage() {}

func (x *TLSRoute) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSRoute.ProtoReflect.Descriptor instead.
func (*TLSRoute) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{3}
}

func (x *TLSRoute) GetMatchers() []*v2.TLSRequestMatcher {
	if x != nil {
		return x.Matchers
	}
	return nil
}

func (m *TLSRoute) GetActionType() isTLSRoute_ActionType {
	if m != nil {
		return m.ActionType
	}
	return nil
}

func (x *TLSRoute) GetForwardTo() *TLSRoute_TLSForwardToAction {
	if x, ok := x.GetActionType().(*TLSRoute_ForwardTo); ok {
		return x.ForwardTo
	}
	return nil
}

type isTLSRoute_ActionType interface {
	isTLSRoute_ActionType()
}

type TLSRoute_ForwardTo struct {
	// Forward traffic to one or more destination services.
	//
	// </br>**Configuration constraints**: This field is required.
	//
	// +kubebuilder:validation:Required
	ForwardTo *TLSRoute_TLSForwardToAction `protobuf:"bytes,2,opt,name=forward_to,json=forwardTo,proto3,oneof"`
}

func (*TLSRoute_ForwardTo) isTLSRoute_ActionType() {}

// Handle the HTTP request as a GraphQL request, including query validation and execution of the GraphQL request.
// The incoming GraphQL request must either be a GET or POST request. For more information, see
// [Serving over HTTP](https://graphql.org/learn/serving-over-http/) in the GraphQL docs.
type GraphQLAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Reference to a GraphQLSchema or GraphQLStitchedSchema resource that contains the configuration for this subschema.
	//
	// </br>**Configuration constraints**: One of `schema` or `stitchedSchema` must be set, but not both.
	//
	// Types that are assignable to GraphqlSchema:
	//
	//	*GraphQLAction_Schema
	//	*GraphQLAction_StitchedSchema
	GraphqlSchema isGraphQLAction_GraphqlSchema `protobuf_oneof:"graphql_schema"`
	// Options that apply to this GraphQL Schema.
	Options *GraphQLAction_Options `protobuf:"bytes,4,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *GraphQLAction) Reset() {
	*x = GraphQLAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GraphQLAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GraphQLAction) ProtoMessage() {}

func (x *GraphQLAction) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GraphQLAction.ProtoReflect.Descriptor instead.
func (*GraphQLAction) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{4}
}

func (m *GraphQLAction) GetGraphqlSchema() isGraphQLAction_GraphqlSchema {
	if m != nil {
		return m.GraphqlSchema
	}
	return nil
}

func (x *GraphQLAction) GetSchema() *v1.ClusterObjectRef {
	if x, ok := x.GetGraphqlSchema().(*GraphQLAction_Schema); ok {
		return x.Schema
	}
	return nil
}

func (x *GraphQLAction) GetStitchedSchema() *v1.ClusterObjectRef {
	if x, ok := x.GetGraphqlSchema().(*GraphQLAction_StitchedSchema); ok {
		return x.StitchedSchema
	}
	return nil
}

func (x *GraphQLAction) GetOptions() *GraphQLAction_Options {
	if x != nil {
		return x.Options
	}
	return nil
}

type isGraphQLAction_GraphqlSchema interface {
	isGraphQLAction_GraphqlSchema()
}

type GraphQLAction_Schema struct {
	// Reference to a GraphQLSchema resource that contains the configuration for this subschema.
	Schema *v1.ClusterObjectRef `protobuf:"bytes,1,opt,name=schema,proto3,oneof"`
}

type GraphQLAction_StitchedSchema struct {
	// Reference to a GraphQLStitchedSchema resource that contains the configuration for this subschema.
	StitchedSchema *v1.ClusterObjectRef `protobuf:"bytes,2,opt,name=stitched_schema,json=stitchedSchema,proto3,oneof"`
}

func (*GraphQLAction_Schema) isGraphQLAction_GraphqlSchema() {}

func (*GraphQLAction_StitchedSchema) isGraphQLAction_GraphqlSchema() {}

// When a client request matches a route, Gloo forwards the request to the destination that you specify in this `forwardTo` action.
type ForwardToAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Define the upstream destination to route the request to. Some destinations require additional configuration for
	// the route. For example, to forward requests to a CloudProvider for an AWS Lambda, you must also set a `function`.
	// HTTP routes support all destinations types. TCP routes support only Kubernetes services and Gloo VirtualDestinations.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>If `defaultDestination` is empty, you must specify at least one destination in this field.</li>
	// <li>A destination `subset` must not be set to the empty object `{}`.</li>
	// <li>You can optionally specify a destination `weight` to indicate the proportion of traffic
	// to forward to this destination. Weights across all destinations must sum to 100.
	// If the sum is less than 100, the remainder is distributed across destinations that do not specify a weight,
	// with a minimum of 1 weight per destination. Destination weight examples:<ul>
	// <li>Valid example: Port 80 specifies a weight of `50`, port 81 a weight of `25`, and port 82 a weight of `25`.
	// All weights equal 100. 50% of traffic is forwarded to port 80,
	// 25% to 81, and 25% to 82.</li>
	// <li>Valid example: Port 80 specifies a weight of `50`, port 81 a weight of `25`, and port 82 does not
	// specify a weight. All weights equal 75, and the remaining 25% is assigned to port 82.</li>
	// <li>Invalid example: Port 80 specifies a weight of `50`, port 81 a weight of `50`, and port 82 a weight
	// of `25`. All weights equal 125.</li>
	// <li>Invalid example: Port 80 specifies a weight of `50`, port 81 a weight of `50`, and port 82 does not
	// specify a weight. All weights equal 100, but no remainder exists for port 82.</li></ul></li></ul>
	//
	// +kubebuilder:validation:MaxItems=99
	// +kubebuilder:validation:XValidation:rule="self.all(d, !has(d.subset) || d.subset.size() > 0)",message="Destination subset must not be an empty map."
	Destinations []*v2.DestinationReference `protobuf:"bytes,1,rep,name=destinations,proto3" json:"destinations,omitempty"`
	// Types that are assignable to PathRewriteSpecifier:
	//
	//	*ForwardToAction_PathRewrite
	//	*ForwardToAction_RegexRewrite
	PathRewriteSpecifier isForwardToAction_PathRewriteSpecifier `protobuf_oneof:"path_rewrite_specifier"`
	// Types that are assignable to HostRewriteSpecifier:
	//
	//	*ForwardToAction_HostRewrite
	//	*ForwardToAction_AutoHostRewrite
	HostRewriteSpecifier isForwardToAction_HostRewriteSpecifier `protobuf_oneof:"host_rewrite_specifier"`
}

func (x *ForwardToAction) Reset() {
	*x = ForwardToAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardToAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardToAction) ProtoMessage() {}

func (x *ForwardToAction) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardToAction.ProtoReflect.Descriptor instead.
func (*ForwardToAction) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{5}
}

func (x *ForwardToAction) GetDestinations() []*v2.DestinationReference {
	if x != nil {
		return x.Destinations
	}
	return nil
}

func (m *ForwardToAction) GetPathRewriteSpecifier() isForwardToAction_PathRewriteSpecifier {
	if m != nil {
		return m.PathRewriteSpecifier
	}
	return nil
}

func (x *ForwardToAction) GetPathRewrite() string {
	if x, ok := x.GetPathRewriteSpecifier().(*ForwardToAction_PathRewrite); ok {
		return x.PathRewrite
	}
	return ""
}

func (x *ForwardToAction) GetRegexRewrite() *v3.RegexMatchAndSubstitute {
	if x, ok := x.GetPathRewriteSpecifier().(*ForwardToAction_RegexRewrite); ok {
		return x.RegexRewrite
	}
	return nil
}

func (m *ForwardToAction) GetHostRewriteSpecifier() isForwardToAction_HostRewriteSpecifier {
	if m != nil {
		return m.HostRewriteSpecifier
	}
	return nil
}

func (x *ForwardToAction) GetHostRewrite() string {
	if x, ok := x.GetHostRewriteSpecifier().(*ForwardToAction_HostRewrite); ok {
		return x.HostRewrite
	}
	return ""
}

func (x *ForwardToAction) GetAutoHostRewrite() bool {
	if x, ok := x.GetHostRewriteSpecifier().(*ForwardToAction_AutoHostRewrite); ok {
		return x.AutoHostRewrite
	}
	return false
}

type isForwardToAction_PathRewriteSpecifier interface {
	isForwardToAction_PathRewriteSpecifier()
}

type ForwardToAction_PathRewrite struct {
	// Replace the path specified in the matcher with this value before forwarding the request to the upstream destination.
	// When a prefix matcher is used, only the prefix portion of the path is rewritten. When an exact matcher is used,
	// the whole path is replaced. Rewriting the path when a regex matcher is used is currently unsupported. Note that path
	// rewrites are available for HTTP routes only and are not supported for TCP routes.
	PathRewrite string `protobuf:"bytes,2,opt,name=path_rewrite,json=pathRewrite,proto3,oneof"`
}

type ForwardToAction_RegexRewrite struct {
	// During forwarding, portions of the path that match the pattern are rewritten, even allowing the substitution
	// of capture groups from the pattern into the new path as specified by the rewrite substitution string. This substitution is useful
	// to allow application paths to be rewritten in a way that is aware of segments with variable content like identifiers.
	// Note that regex rewrites are available for RE2 syntax and HTTP routes only.
	//
	// </br>**Configuration constraints**: The value must follow a valid RE2 regex pattern.
	RegexRewrite *v3.RegexMatchAndSubstitute `protobuf:"bytes,5,opt,name=regex_rewrite,json=regexRewrite,proto3,oneof"`
}

func (*ForwardToAction_PathRewrite) isForwardToAction_PathRewriteSpecifier() {}

func (*ForwardToAction_RegexRewrite) isForwardToAction_PathRewriteSpecifier() {}

type isForwardToAction_HostRewriteSpecifier interface {
	isForwardToAction_HostRewriteSpecifier()
}

type ForwardToAction_HostRewrite struct {
	// Replace the Authority/Host header with this value before forwarding the request to the upstream destination.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>Supported for HTTP routes only. Unsupported for TCP routes.</li>
	// <li>Hostnames must be 1 - 255 characters in length.</li>
	// <li>Supported characters are `a-z`, `A-Z`, `0-9`, `-`, and `.`.</li>
	// <li>Each segment separated by a period (`.`) must be 1 - 63 characters in length and cannot start with the `-` character.</li></ul>
	HostRewrite string `protobuf:"bytes,3,opt,name=host_rewrite,json=hostRewrite,proto3,oneof"`
}

type ForwardToAction_AutoHostRewrite struct {
	// Automatically replace the Authority/Host header with the hostname of the upstream destination. Note
	// that host rewrites are available for HTTP routes only and are not supported for TCP routes.
	AutoHostRewrite bool `protobuf:"varint,4,opt,name=auto_host_rewrite,json=autoHostRewrite,proto3,oneof"`
}

func (*ForwardToAction_HostRewrite) isForwardToAction_HostRewriteSpecifier() {}

func (*ForwardToAction_AutoHostRewrite) isForwardToAction_HostRewriteSpecifier() {}

// <!-- This message needs to be at this level (rather than nested) due to cue restrictions.-->
// <!-- RedirectAction is copied directly from https://github.com/envoyproxy/envoy/blob/master/api/envoy/api/v2/route/route.proto-->
// Return a redirect response to the downstream client.
type RedirectAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The host portion of the URL is swapped with this value.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>Hostnames must be 1 - 255 characters in length</li>
	// <li>Supported characters are `a-z`, `A-Z`, `0-9`, `-`, and `.`.</li>
	// <li>Each segment separated by a period (`.`) must be 1 - 63 characters in length and cannot start with the `-` character.</li></ul>
	HostRedirect string `protobuf:"bytes,1,opt,name=host_redirect,json=hostRedirect,proto3" json:"host_redirect,omitempty"`
	// Defines whether and how the path portion of the URL is modified.
	//
	// Types that are assignable to PathRewriteSpecifier:
	//
	//	*RedirectAction_PathRedirect
	PathRewriteSpecifier isRedirectAction_PathRewriteSpecifier `protobuf_oneof:"path_rewrite_specifier"`
	// The HTTP status code to use in the redirect response. The default response
	// code is MOVED_PERMANENTLY (301).
	ResponseCode RedirectAction_RedirectResponseCode `protobuf:"varint,4,opt,name=response_code,json=responseCode,proto3,enum=networking.gloo.solo.io.RedirectAction_RedirectResponseCode" json:"response_code,omitempty"`
}

func (x *RedirectAction) Reset() {
	*x = RedirectAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedirectAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedirectAction) ProtoMessage() {}

func (x *RedirectAction) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedirectAction.ProtoReflect.Descriptor instead.
func (*RedirectAction) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{6}
}

func (x *RedirectAction) GetHostRedirect() string {
	if x != nil {
		return x.HostRedirect
	}
	return ""
}

func (m *RedirectAction) GetPathRewriteSpecifier() isRedirectAction_PathRewriteSpecifier {
	if m != nil {
		return m.PathRewriteSpecifier
	}
	return nil
}

func (x *RedirectAction) GetPathRedirect() string {
	if x, ok := x.GetPathRewriteSpecifier().(*RedirectAction_PathRedirect); ok {
		return x.PathRedirect
	}
	return ""
}

func (x *RedirectAction) GetResponseCode() RedirectAction_RedirectResponseCode {
	if x != nil {
		return x.ResponseCode
	}
	return RedirectAction_MOVED_PERMANENTLY
}

type isRedirectAction_PathRewriteSpecifier interface {
	isRedirectAction_PathRewriteSpecifier()
}

type RedirectAction_PathRedirect struct {
	// The entire path portion of the URL is overwritten with this value.
	PathRedirect string `protobuf:"bytes,2,opt,name=path_redirect,json=pathRedirect,proto3,oneof"`
}

func (*RedirectAction_PathRedirect) isRedirectAction_PathRewriteSpecifier() {}

// <!-- This message needs to be at this level (rather than nested) due to cue restrictions.-->
// <!-- DirectResponseAction is copied directly from https://github.com/envoyproxy/envoy/blob/master/api/envoy/api/v2/route/route.proto-->
// Respond directly to the client from the proxy.
type DirectResponseAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The HTTP response status code to return.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>This field is required.</li>
	// <li>The value must be 200 - 599, inclusive.</li>
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=100
	// +kubebuilder:validation:Maximum=599
	Status uint32 `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	// The content of the response body. If omitted,
	// no body is included in the generated response.
	//
	// </br>**Configuration constraints**: Must be less than 1MB in size.
	//
	// +kubebuilder:validation:MaxLength=1048576
	Body string `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *DirectResponseAction) Reset() {
	*x = DirectResponseAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirectResponseAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirectResponseAction) ProtoMessage() {}

func (x *DirectResponseAction) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirectResponseAction.ProtoReflect.Descriptor instead.
func (*DirectResponseAction) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{7}
}

func (x *DirectResponseAction) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *DirectResponseAction) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

// <!-- This message needs to be at this level (rather than nested) due to cue restrictions.-->
// Delegate routing decisions to one or more HTTP route tables.
// This can be used to delegate a subset of the route table's traffic to another route table, which may live
// in an imported workspace, or to separate routing concerns between objects.
type DelegateAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Delegate to the route tables that match the given selectors.
	// Selected route tables are ordered by creation time stamp in ascending order.
	// Route tables are selected from both the tables defined within the current workspace and any tables imported into the workspace.
	RouteTables []*v2.ObjectSelector `protobuf:"bytes,2,rep,name=route_tables,json=routeTables,proto3" json:"route_tables,omitempty"`
	// Optional: Restrict delegation to the route tables that match the set of route filter criteria specified.
	// If omitted, any route can be referenced by this route table.
	AllowedRoutes []*v2.RouteFilter `protobuf:"bytes,4,rep,name=allowed_routes,json=allowedRoutes,proto3" json:"allowed_routes,omitempty"`
	// The method by which routes across delegated route tables are sorted.
	SortMethod DelegateAction_SortMethod `protobuf:"varint,3,opt,name=sort_method,json=sortMethod,proto3,enum=networking.gloo.solo.io.DelegateAction_SortMethod" json:"sort_method,omitempty"`
}

func (x *DelegateAction) Reset() {
	*x = DelegateAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelegateAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelegateAction) ProtoMessage() {}

func (x *DelegateAction) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelegateAction.ProtoReflect.Descriptor instead.
func (*DelegateAction) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{8}
}

func (x *DelegateAction) GetRouteTables() []*v2.ObjectSelector {
	if x != nil {
		return x.RouteTables
	}
	return nil
}

func (x *DelegateAction) GetAllowedRoutes() []*v2.RouteFilter {
	if x != nil {
		return x.AllowedRoutes
	}
	return nil
}

func (x *DelegateAction) GetSortMethod() DelegateAction_SortMethod {
	if x != nil {
		return x.SortMethod
	}
	return DelegateAction_TABLE_WEIGHT
}

type RouteTableStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The state and workspace conditions of the applied resource.
	Common *v2.Status `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	// A map of policy GVK to the number of policies that are applied on this resource,
	// sorted by GVK.
	NumAppliedRoutePolicies map[string]uint32 `protobuf:"bytes,2,rep,name=num_applied_route_policies,json=numAppliedRoutePolicies,proto3" json:"num_applied_route_policies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// The number of parent route tables for this route table, if it is a delegated route table.
	NumParentRouteTables uint32 `protobuf:"varint,3,opt,name=num_parent_route_tables,json=numParentRouteTables,proto3" json:"num_parent_route_tables,omitempty"`
	// The name of the workspace that this route table belongs to.
	OwnedByWorkspace string `protobuf:"bytes,4,opt,name=owned_by_workspace,json=ownedByWorkspace,proto3" json:"owned_by_workspace,omitempty"`
	// The number virtual gateways that this route table can select.
	NumAllowedVirtualGateways uint32 `protobuf:"varint,5,opt,name=num_allowed_virtual_gateways,json=numAllowedVirtualGateways,proto3" json:"num_allowed_virtual_gateways,omitempty"`
}

func (x *RouteTableStatus) Reset() {
	*x = RouteTableStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteTableStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteTableStatus) ProtoMessage() {}

func (x *RouteTableStatus) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteTableStatus.ProtoReflect.Descriptor instead.
func (*RouteTableStatus) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{9}
}

func (x *RouteTableStatus) GetCommon() *v2.Status {
	if x != nil {
		return x.Common
	}
	return nil
}

func (x *RouteTableStatus) GetNumAppliedRoutePolicies() map[string]uint32 {
	if x != nil {
		return x.NumAppliedRoutePolicies
	}
	return nil
}

func (x *RouteTableStatus) GetNumParentRouteTables() uint32 {
	if x != nil {
		return x.NumParentRouteTables
	}
	return 0
}

func (x *RouteTableStatus) GetOwnedByWorkspace() string {
	if x != nil {
		return x.OwnedByWorkspace
	}
	return ""
}

func (x *RouteTableStatus) GetNumAllowedVirtualGateways() uint32 {
	if x != nil {
		return x.NumAllowedVirtualGateways
	}
	return 0
}

// The resources that the applied route table selects.
type RouteTableReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of workspaces in which the route table can be applied.
	Workspaces map[string]*v2.Report `protobuf:"bytes,1,rep,name=workspaces,proto3" json:"workspaces,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// A map of policy GVK to policy references for all policies that are applied on this
	// resource.
	AppliedRoutePolicies map[string]*v2.AppliedRoutePolicies `protobuf:"bytes,2,rep,name=applied_route_policies,json=appliedRoutePolicies,proto3" json:"applied_route_policies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// A list of the parents route tables for this route table, if it is a delegated route table.
	ParentRouteTables []*v2.ObjectReference `protobuf:"bytes,3,rep,name=parent_route_tables,json=parentRouteTables,proto3" json:"parent_route_tables,omitempty"`
	// The name of the workspace that owns the route table.
	OwnerWorkspace string `protobuf:"bytes,4,opt,name=owner_workspace,json=ownerWorkspace,proto3" json:"owner_workspace,omitempty"`
	// A list of allowed virtual gateways that this route table can select.
	AllowedVirtualGateways []*v2.ObjectReference `protobuf:"bytes,5,rep,name=allowed_virtual_gateways,json=allowedVirtualGateways,proto3" json:"allowed_virtual_gateways,omitempty"`
	// A list of routes delegated to by delegated routes in this route table.
	// Only tracks direct delegates of this route table; delegates of delegate routes are not included.
	DelegatedToRouteTables []*RouteTableReport_DelegatedRouteTableReference `protobuf:"bytes,6,rep,name=delegated_to_route_tables,json=delegatedToRouteTables,proto3" json:"delegated_to_route_tables,omitempty"`
}

func (x *RouteTableReport) Reset() {
	*x = RouteTableReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteTableReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteTableReport) ProtoMessage() {}

func (x *RouteTableReport) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteTableReport.ProtoReflect.Descriptor instead.
func (*RouteTableReport) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{10}
}

func (x *RouteTableReport) GetWorkspaces() map[string]*v2.Report {
	if x != nil {
		return x.Workspaces
	}
	return nil
}

func (x *RouteTableReport) GetAppliedRoutePolicies() map[string]*v2.AppliedRoutePolicies {
	if x != nil {
		return x.AppliedRoutePolicies
	}
	return nil
}

func (x *RouteTableReport) GetParentRouteTables() []*v2.ObjectReference {
	if x != nil {
		return x.ParentRouteTables
	}
	return nil
}

func (x *RouteTableReport) GetOwnerWorkspace() string {
	if x != nil {
		return x.OwnerWorkspace
	}
	return ""
}

func (x *RouteTableReport) GetAllowedVirtualGateways() []*v2.ObjectReference {
	if x != nil {
		return x.AllowedVirtualGateways
	}
	return nil
}

func (x *RouteTableReport) GetDelegatedToRouteTables() []*RouteTableReport_DelegatedRouteTableReference {
	if x != nil {
		return x.DelegatedToRouteTables
	}
	return nil
}

// When a client request matches a route, Gloo forwards the request to the destination that you specify in this `forwardTo` action.
type TLSRoute_TLSForwardToAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Define the upstream destination to route the request to.
	//
	// </br>**Configuration constraints**:<ul>
	// <li>If `defaultDestination` is empty, you must specify at least one destination in this field.</li>
	// <li>A destination `subset` must not be set to the empty object `{}`.</li>
	// <li>You can optionally specify a destination `weight` to indicate the proportion of traffic
	// to forward to this destination. Weights across all destinations must sum to 100.
	// If the sum is less than 100, the remainder is distributed across destinations that do not specify a weight,
	// with a minimum of 1 weight per destination. Destination weight examples:<ul>
	// <li>Valid example: Port 80 specifies a weight of `50`, port 81 a weight of `25`, and port 82 a weight of `25`.
	// All weights equal 100. 50% of traffic is forwarded to port 80,
	// 25% to 81, and 25% to 82.</li>
	// <li>Valid example: Port 80 specifies a weight of `50`, port 81 a weight of `25`, and port 82 does not
	// specify a weight. All weights equal 75, and the remaining 25% is assigned to port 82.</li>
	// <li>Invalid example: Port 80 specifies a weight of `50`, port 81 a weight of `50`, and port 82 a weight
	// of `25`. All weights equal 125.</li>
	// <li>Invalid example: Port 80 specifies a weight of `50`, port 81 a weight of `50`, and port 82 does not
	// specify a weight. All weights equal 100, but no remainder exists for port 82.</li></ul></li></ul>
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=99
	// +kubebuilder:validation:XValidation:rule="self.all(d, !has(d.subset) || d.subset.size() > 0)",message="Destination subset must not be an empty map."
	Destinations []*v2.DestinationReference `protobuf:"bytes,1,rep,name=destinations,proto3" json:"destinations,omitempty"`
}

func (x *TLSRoute_TLSForwardToAction) Reset() {
	*x = TLSRoute_TLSForwardToAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TLSRoute_TLSForwardToAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSRoute_TLSForwardToAction) ProtoMessage() {}

func (x *TLSRoute_TLSForwardToAction) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSRoute_TLSForwardToAction.ProtoReflect.Descriptor instead.
func (*TLSRoute_TLSForwardToAction) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{3, 0}
}

func (x *TLSRoute_TLSForwardToAction) GetDestinations() []*v2.DestinationReference {
	if x != nil {
		return x.Destinations
	}
	return nil
}

type GraphQLAction_Options struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Include information about request/response in the envoy debug logs.
	// This is helpful for debugging GraphQL.
	// Defaults to false.
	LogSensitiveInfo *wrapperspb.BoolValue `protobuf:"bytes,1,opt,name=log_sensitive_info,json=logSensitiveInfo,proto3" json:"log_sensitive_info,omitempty"`
}

func (x *GraphQLAction_Options) Reset() {
	*x = GraphQLAction_Options{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GraphQLAction_Options) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GraphQLAction_Options) ProtoMessage() {}

func (x *GraphQLAction_Options) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GraphQLAction_Options.ProtoReflect.Descriptor instead.
func (*GraphQLAction_Options) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{4, 0}
}

func (x *GraphQLAction_Options) GetLogSensitiveInfo() *wrapperspb.BoolValue {
	if x != nil {
		return x.LogSensitiveInfo
	}
	return nil
}

// A list of routes delegated to by delegated routes in this route table.
// Only tracks direct delegates of this route table; delegates of delegate routes are not included.
type RouteTableReport_DelegatedRouteTableReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The index of the route in the parent route table that delegates to the listed route table.
	RouteIndex int32 `protobuf:"varint,1,opt,name=route_index,json=routeIndex,proto3" json:"route_index,omitempty"`
	// The reference to the route table being delegated to by the parent route table.
	RouteTable *v2.ObjectReference `protobuf:"bytes,2,opt,name=route_table,json=routeTable,proto3" json:"route_table,omitempty"`
}

func (x *RouteTableReport_DelegatedRouteTableReference) Reset() {
	*x = RouteTableReport_DelegatedRouteTableReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[18]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteTableReport_DelegatedRouteTableReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteTableReport_DelegatedRouteTableReference) ProtoMessage() {}

func (x *RouteTableReport_DelegatedRouteTableReference) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[18]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteTableReport_DelegatedRouteTableReference.ProtoReflect.Descriptor instead.
func (*RouteTableReport_DelegatedRouteTableReference) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP(), []int{10, 2}
}

func (x *RouteTableReport_DelegatedRouteTableReference) GetRouteIndex() int32 {
	if x != nil {
		return x.RouteIndex
	}
	return 0
}

func (x *RouteTableReport_DelegatedRouteTableReference) GetRouteTable() *v2.ObjectReference {
	if x != nil {
		return x.RouteTable
	}
	return nil
}

var File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto protoreflect.FileDescriptor

var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDesc = []byte{
	0x0a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x32, 0x2f, 0x72, 0x6f, 0x75, 0x74,
	0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x1a, 0x12, 0x65, 0x78, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x76,
	0x33, 0x2f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x67, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69,
	0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65,
	0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c,
	0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x32, 0x2f, 0x67, 0x72, 0x61, 0x70,
	0x68, 0x71, 0x6c, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x5f, 0x6d, 0x61, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x59, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d,
	0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f,
	0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x68, 0x74,
	0x74, 0x70, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f,
	0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x55, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c,
	0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69,
	0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32,
	0x2f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x52, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x62, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65,
	0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69,
	0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f,
	0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73,
	0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f,
	0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x58, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d,
	0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76,
	0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x74, 0x63, 0x70,
	0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x58, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f,
	0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e,
	0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x74, 0x6c, 0x73, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6b,
	0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf7, 0x07, 0x0a, 0x0e, 0x52, 0x6f,
	0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12, 0x14, 0x0a, 0x05,
	0x68, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x68, 0x6f, 0x73,
	0x74, 0x73, 0x12, 0x4f, 0x0a, 0x10, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x52, 0x0f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x47, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x73, 0x12, 0x54, 0x0a, 0x12, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f,
	0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x11, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x5c, 0x0a, 0x15, 0x61, 0x70, 0x70,
	0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x13, 0x61, 0x70, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x44, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x5a, 0x0a, 0x13, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c,
	0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52,
	0x12, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x04, 0x68, 0x74, 0x74, 0x70, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x48, 0x54, 0x54, 0x50,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x04, 0x68, 0x74, 0x74, 0x70, 0x12, 0x33, 0x0a, 0x03, 0x74,
	0x63, 0x70, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x54, 0x43, 0x50, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x03, 0x74, 0x63, 0x70,
	0x12, 0x33, 0x0a, 0x03, 0x74, 0x6c, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x54, 0x4c, 0x53, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x52, 0x03, 0x74, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x4c, 0x0a,
	0x0f, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x50, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x0e, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x56, 0x0a, 0x0c, 0x66,
	0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x33, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x46, 0x61, 0x69, 0x6c, 0x75,
	0x72, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x4d,
	0x6f, 0x64, 0x65, 0x12, 0x86, 0x01, 0x0a, 0x1b, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x46, 0x2e, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x70,
	0x65, 0x63, 0x2e, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x19, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x4c, 0x0a, 0x1e,
	0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x37, 0x0a, 0x0b, 0x46, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x4f, 0x55,
	0x54, 0x45, 0x5f, 0x52, 0x45, 0x50, 0x4c, 0x41, 0x43, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x00,
	0x12, 0x11, 0x0a, 0x0d, 0x46, 0x52, 0x45, 0x45, 0x5a, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49,
	0x47, 0x10, 0x01, 0x22, 0xed, 0x04, 0x0a, 0x09, 0x48, 0x54, 0x54, 0x50, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x46, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69,
	0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e,
	0x48, 0x54, 0x54, 0x50, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x43, 0x0a,
	0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x72, 0x73, 0x12, 0x49, 0x0a, 0x0a, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x5f, 0x74, 0x6f,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x54, 0x6f, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x48, 0x00, 0x52, 0x09, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x54, 0x6f, 0x12, 0x45, 0x0a,
	0x08, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f,
	0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61,
	0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x08, 0x64, 0x65, 0x6c, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x12, 0x45, 0x0a, 0x08, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x00, 0x52, 0x08, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12, 0x58, 0x0a, 0x0f, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e,
	0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x07, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x51, 0x4c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00,
	0x52, 0x07, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x22, 0xa8, 0x01, 0x0a, 0x08, 0x54, 0x43, 0x50, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x12, 0x42, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x26, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f,
	0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x54, 0x43, 0x50, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x08, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x65, 0x72, 0x73, 0x12, 0x49, 0x0a, 0x0a, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x5f,
	0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x54, 0x6f, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x48, 0x00, 0x52, 0x09, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x54, 0x6f, 0x42,
	0x0d, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0x99,
	0x02, 0x0a, 0x08, 0x54, 0x4c, 0x53, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x42, 0x0a, 0x08, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x54, 0x4c, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x12,
	0x55, 0x0a, 0x0a, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x5f, 0x74, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67,
	0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x54, 0x4c,
	0x53, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x54, 0x4c, 0x53, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72,
	0x64, 0x54, 0x6f, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x09, 0x66, 0x6f, 0x72,
	0x77, 0x61, 0x72, 0x64, 0x54, 0x6f, 0x1a, 0x63, 0x0a, 0x12, 0x54, 0x4c, 0x53, 0x46, 0x6f, 0x72,
	0x77, 0x61, 0x72, 0x64, 0x54, 0x6f, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4d, 0x0a, 0x0c,
	0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f,
	0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0c, 0x64,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0xcf, 0x02, 0x0a, 0x0d, 0x47,
	0x72, 0x61, 0x70, 0x68, 0x51, 0x4c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a, 0x06,
	0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6b, 0x76, 0x32, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x66, 0x48, 0x00, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x4e, 0x0a, 0x0f, 0x73,
	0x74, 0x69, 0x74, 0x63, 0x68, 0x65, 0x64, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6b, 0x76, 0x32,
	0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x66, 0x48, 0x00, 0x52, 0x0e, 0x73, 0x74, 0x69,
	0x74, 0x63, 0x68, 0x65, 0x64, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x48, 0x0a, 0x07, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x51, 0x4c, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x53, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x48, 0x0a, 0x12, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x76,
	0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42,
	0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x10, 0x6c, 0x6f, 0x67, 0x53, 0x65, 0x6e,
	0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x10, 0x0a, 0x0e, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x71, 0x6c, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x22, 0xe3, 0x02, 0x0a,
	0x0f, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x54, 0x6f, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x4d, 0x0a, 0x0c, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x0c, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x23, 0x0a, 0x0c, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x72, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x70, 0x61, 0x74, 0x68, 0x52, 0x65, 0x77,
	0x72, 0x69, 0x74, 0x65, 0x12, 0x55, 0x0a, 0x0d, 0x72, 0x65, 0x67, 0x65, 0x78, 0x5f, 0x72, 0x65,
	0x77, 0x72, 0x69, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x65, 0x6e,
	0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72,
	0x2e, 0x76, 0x33, 0x2e, 0x52, 0x65, 0x67, 0x65, 0x78, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x41, 0x6e,
	0x64, 0x53, 0x75, 0x62, 0x73, 0x74, 0x69, 0x74, 0x75, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x72,
	0x65, 0x67, 0x65, 0x78, 0x52, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x68,
	0x6f, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x01, 0x52, 0x0b, 0x68, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65,
	0x12, 0x2c, 0x0a, 0x11, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x72, 0x65,
	0x77, 0x72, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x0f, 0x61,
	0x75, 0x74, 0x6f, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65, 0x42, 0x18,
	0x0a, 0x16, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x72, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x73,
	0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x72, 0x42, 0x18, 0x0a, 0x16, 0x68, 0x6f, 0x73, 0x74,
	0x5f, 0x72, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x65, 0x72, 0x22, 0xd2, 0x02, 0x0a, 0x0e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x72, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x68, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12, 0x25, 0x0a, 0x0d, 0x70, 0x61,
	0x74, 0x68, 0x5f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x0c, 0x70, 0x61, 0x74, 0x68, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x12, 0x61, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3c, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x43, 0x6f, 0x64, 0x65, 0x22, 0x77, 0x0a, 0x14, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x15, 0x0a, 0x11,
	0x4d, 0x4f, 0x56, 0x45, 0x44, 0x5f, 0x50, 0x45, 0x52, 0x4d, 0x41, 0x4e, 0x45, 0x4e, 0x54, 0x4c,
	0x59, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x53, 0x45, 0x45, 0x5f, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x02, 0x12, 0x16, 0x0a,
	0x12, 0x54, 0x45, 0x4d, 0x50, 0x4f, 0x52, 0x41, 0x52, 0x59, 0x5f, 0x52, 0x45, 0x44, 0x49, 0x52,
	0x45, 0x43, 0x54, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x50, 0x45, 0x52, 0x4d, 0x41, 0x4e, 0x45,
	0x4e, 0x54, 0x5f, 0x52, 0x45, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x10, 0x04, 0x42, 0x18, 0x0a,
	0x16, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x72, 0x65, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x73, 0x70,
	0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22, 0x42, 0x0a, 0x14, 0x44, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0xad, 0x02, 0x0a, 0x0e,
	0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x46,
	0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c,
	0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x47, 0x0a, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x64, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c,
	0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x12,
	0x53, 0x0a, 0x0b, 0x73, 0x6f, 0x72, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x32, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e,
	0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x6f,
	0x72, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x0a, 0x73, 0x6f, 0x72, 0x74, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x22, 0x35, 0x0a, 0x0a, 0x53, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x5f, 0x57, 0x45, 0x49, 0x47,
	0x48, 0x54, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x5f, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x43, 0x49, 0x54, 0x59, 0x10, 0x01, 0x22, 0xbf, 0x03, 0x0a, 0x10,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x33, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73,
	0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x12, 0x83, 0x01, 0x0a, 0x1a, 0x6e, 0x75, 0x6d, 0x5f, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x65, 0x64, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x46, 0x2e, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c,
	0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x4e, 0x75, 0x6d, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x64,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x17, 0x6e, 0x75, 0x6d, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x52, 0x6f,
	0x75, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x12, 0x35, 0x0a, 0x17, 0x6e,
	0x75, 0x6d, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x14, 0x6e, 0x75,
	0x6d, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x6f, 0x77, 0x6e, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x5f, 0x77,
	0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10,
	0x6f, 0x77, 0x6e, 0x65, 0x64, 0x42, 0x79, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x12, 0x3f, 0x0a, 0x1c, 0x6e, 0x75, 0x6d, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f,
	0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x19, 0x6e, 0x75, 0x6d, 0x41, 0x6c, 0x6c, 0x6f, 0x77,
	0x65, 0x64, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x73, 0x1a, 0x4a, 0x0a, 0x1c, 0x4e, 0x75, 0x6d, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xa4, 0x07,
	0x0a, 0x10, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x59, 0x0a, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x12, 0x79, 0x0a,
	0x16, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x43, 0x2e,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x64,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x14, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x12, 0x54, 0x0a, 0x13, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x11, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x27,
	0x0a, 0x0f, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x57, 0x6f,
	0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x5e, 0x0a, 0x18, 0x61, 0x6c, 0x6c, 0x6f, 0x77,
	0x65, 0x64, 0x5f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52,
	0x16, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x47,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x12, 0x81, 0x01, 0x0a, 0x19, 0x64, 0x65, 0x6c, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x6f, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x46, 0x2e, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x16, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x54, 0x6f,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a, 0x5a, 0x0a, 0x0f, 0x57,
	0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x31, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x72, 0x0a, 0x19, 0x41, 0x70, 0x70, 0x6c, 0x69,
	0x65, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x41, 0x70, 0x70, 0x6c,
	0x69, 0x65, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x86, 0x01, 0x0a, 0x1c,
	0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0a, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x45, 0x0a,
	0x0b, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x67, 0x6c, 0x6f, 0x6f,
	0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0a, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x42, 0x5b, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2d,
	0x6d, 0x65, 0x73, 0x68, 0x2d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2f,
	0x76, 0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2f, 0x76, 0x32, 0xc0, 0xf5, 0x04, 0x01, 0xb8, 0xf5, 0x04, 0x01, 0xd0, 0xf5, 0x04,
	0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescData = file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDesc
)

func file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescData)
	})
	return file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDescData
}

var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes = make([]protoimpl.MessageInfo, 19)
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_goTypes = []interface{}{
	(RouteTableSpec_FailureMode)(0),          // 0: networking.gloo.solo.io.RouteTableSpec.FailureMode
	(RedirectAction_RedirectResponseCode)(0), // 1: networking.gloo.solo.io.RedirectAction.RedirectResponseCode
	(DelegateAction_SortMethod)(0),           // 2: networking.gloo.solo.io.DelegateAction.SortMethod
	(*RouteTableSpec)(nil),                   // 3: networking.gloo.solo.io.RouteTableSpec
	(*HTTPRoute)(nil),                        // 4: networking.gloo.solo.io.HTTPRoute
	(*TCPRoute)(nil),                         // 5: networking.gloo.solo.io.TCPRoute
	(*TLSRoute)(nil),                         // 6: networking.gloo.solo.io.TLSRoute
	(*GraphQLAction)(nil),                    // 7: networking.gloo.solo.io.GraphQLAction
	(*ForwardToAction)(nil),                  // 8: networking.gloo.solo.io.ForwardToAction
	(*RedirectAction)(nil),                   // 9: networking.gloo.solo.io.RedirectAction
	(*DirectResponseAction)(nil),             // 10: networking.gloo.solo.io.DirectResponseAction
	(*DelegateAction)(nil),                   // 11: networking.gloo.solo.io.DelegateAction
	(*RouteTableStatus)(nil),                 // 12: networking.gloo.solo.io.RouteTableStatus
	(*RouteTableReport)(nil),                 // 13: networking.gloo.solo.io.RouteTableReport
	nil,                                      // 14: networking.gloo.solo.io.RouteTableSpec.VirtualServiceAnnotationsEntry
	nil,                                      // 15: networking.gloo.solo.io.HTTPRoute.LabelsEntry
	(*TLSRoute_TLSForwardToAction)(nil),      // 16: networking.gloo.solo.io.TLSRoute.TLSForwardToAction
	(*GraphQLAction_Options)(nil),            // 17: networking.gloo.solo.io.GraphQLAction.Options
	nil,                                      // 18: networking.gloo.solo.io.RouteTableStatus.NumAppliedRoutePoliciesEntry
	nil,                                      // 19: networking.gloo.solo.io.RouteTableReport.WorkspacesEntry
	nil,                                      // 20: networking.gloo.solo.io.RouteTableReport.AppliedRoutePoliciesEntry
	(*RouteTableReport_DelegatedRouteTableReference)(nil), // 21: networking.gloo.solo.io.RouteTableReport.DelegatedRouteTableReference
	(*v2.ObjectReference)(nil),                            // 22: common.gloo.solo.io.ObjectReference
	(*v2.WorkloadSelector)(nil),                           // 23: common.gloo.solo.io.WorkloadSelector
	(*v2.DestinationSelector)(nil),                        // 24: common.gloo.solo.io.DestinationSelector
	(*v2.DestinationReference)(nil),                       // 25: common.gloo.solo.io.DestinationReference
	(*v2.PortalMetadata)(nil),                             // 26: common.gloo.solo.io.PortalMetadata
	(*v2.HTTPRequestMatcher)(nil),                         // 27: common.gloo.solo.io.HTTPRequestMatcher
	(*v2.TCPRequestMatcher)(nil),                          // 28: common.gloo.solo.io.TCPRequestMatcher
	(*v2.TLSRequestMatcher)(nil),                          // 29: common.gloo.solo.io.TLSRequestMatcher
	(*v1.ClusterObjectRef)(nil),                           // 30: core.skv2.solo.io.ClusterObjectRef
	(*v3.RegexMatchAndSubstitute)(nil),                    // 31: envoy.type.matcher.v3.RegexMatchAndSubstitute
	(*v2.ObjectSelector)(nil),                             // 32: common.gloo.solo.io.ObjectSelector
	(*v2.RouteFilter)(nil),                                // 33: common.gloo.solo.io.RouteFilter
	(*v2.Status)(nil),                                     // 34: common.gloo.solo.io.Status
	(*wrapperspb.BoolValue)(nil),                          // 35: google.protobuf.BoolValue
	(*v2.Report)(nil),                                     // 36: common.gloo.solo.io.Report
	(*v2.AppliedRoutePolicies)(nil),                       // 37: common.gloo.solo.io.AppliedRoutePolicies
}
var file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_depIdxs = []int32{
	22, // 0: networking.gloo.solo.io.RouteTableSpec.virtual_gateways:type_name -> common.gloo.solo.io.ObjectReference
	23, // 1: networking.gloo.solo.io.RouteTableSpec.workload_selectors:type_name -> common.gloo.solo.io.WorkloadSelector
	24, // 2: networking.gloo.solo.io.RouteTableSpec.apply_to_destinations:type_name -> common.gloo.solo.io.DestinationSelector
	25, // 3: networking.gloo.solo.io.RouteTableSpec.default_destination:type_name -> common.gloo.solo.io.DestinationReference
	4,  // 4: networking.gloo.solo.io.RouteTableSpec.http:type_name -> networking.gloo.solo.io.HTTPRoute
	5,  // 5: networking.gloo.solo.io.RouteTableSpec.tcp:type_name -> networking.gloo.solo.io.TCPRoute
	6,  // 6: networking.gloo.solo.io.RouteTableSpec.tls:type_name -> networking.gloo.solo.io.TLSRoute
	26, // 7: networking.gloo.solo.io.RouteTableSpec.portal_metadata:type_name -> common.gloo.solo.io.PortalMetadata
	0,  // 8: networking.gloo.solo.io.RouteTableSpec.failure_mode:type_name -> networking.gloo.solo.io.RouteTableSpec.FailureMode
	14, // 9: networking.gloo.solo.io.RouteTableSpec.virtual_service_annotations:type_name -> networking.gloo.solo.io.RouteTableSpec.VirtualServiceAnnotationsEntry
	15, // 10: networking.gloo.solo.io.HTTPRoute.labels:type_name -> networking.gloo.solo.io.HTTPRoute.LabelsEntry
	27, // 11: networking.gloo.solo.io.HTTPRoute.matchers:type_name -> common.gloo.solo.io.HTTPRequestMatcher
	8,  // 12: networking.gloo.solo.io.HTTPRoute.forward_to:type_name -> networking.gloo.solo.io.ForwardToAction
	11, // 13: networking.gloo.solo.io.HTTPRoute.delegate:type_name -> networking.gloo.solo.io.DelegateAction
	9,  // 14: networking.gloo.solo.io.HTTPRoute.redirect:type_name -> networking.gloo.solo.io.RedirectAction
	10, // 15: networking.gloo.solo.io.HTTPRoute.direct_response:type_name -> networking.gloo.solo.io.DirectResponseAction
	7,  // 16: networking.gloo.solo.io.HTTPRoute.graphql:type_name -> networking.gloo.solo.io.GraphQLAction
	28, // 17: networking.gloo.solo.io.TCPRoute.matchers:type_name -> common.gloo.solo.io.TCPRequestMatcher
	8,  // 18: networking.gloo.solo.io.TCPRoute.forward_to:type_name -> networking.gloo.solo.io.ForwardToAction
	29, // 19: networking.gloo.solo.io.TLSRoute.matchers:type_name -> common.gloo.solo.io.TLSRequestMatcher
	16, // 20: networking.gloo.solo.io.TLSRoute.forward_to:type_name -> networking.gloo.solo.io.TLSRoute.TLSForwardToAction
	30, // 21: networking.gloo.solo.io.GraphQLAction.schema:type_name -> core.skv2.solo.io.ClusterObjectRef
	30, // 22: networking.gloo.solo.io.GraphQLAction.stitched_schema:type_name -> core.skv2.solo.io.ClusterObjectRef
	17, // 23: networking.gloo.solo.io.GraphQLAction.options:type_name -> networking.gloo.solo.io.GraphQLAction.Options
	25, // 24: networking.gloo.solo.io.ForwardToAction.destinations:type_name -> common.gloo.solo.io.DestinationReference
	31, // 25: networking.gloo.solo.io.ForwardToAction.regex_rewrite:type_name -> envoy.type.matcher.v3.RegexMatchAndSubstitute
	1,  // 26: networking.gloo.solo.io.RedirectAction.response_code:type_name -> networking.gloo.solo.io.RedirectAction.RedirectResponseCode
	32, // 27: networking.gloo.solo.io.DelegateAction.route_tables:type_name -> common.gloo.solo.io.ObjectSelector
	33, // 28: networking.gloo.solo.io.DelegateAction.allowed_routes:type_name -> common.gloo.solo.io.RouteFilter
	2,  // 29: networking.gloo.solo.io.DelegateAction.sort_method:type_name -> networking.gloo.solo.io.DelegateAction.SortMethod
	34, // 30: networking.gloo.solo.io.RouteTableStatus.common:type_name -> common.gloo.solo.io.Status
	18, // 31: networking.gloo.solo.io.RouteTableStatus.num_applied_route_policies:type_name -> networking.gloo.solo.io.RouteTableStatus.NumAppliedRoutePoliciesEntry
	19, // 32: networking.gloo.solo.io.RouteTableReport.workspaces:type_name -> networking.gloo.solo.io.RouteTableReport.WorkspacesEntry
	20, // 33: networking.gloo.solo.io.RouteTableReport.applied_route_policies:type_name -> networking.gloo.solo.io.RouteTableReport.AppliedRoutePoliciesEntry
	22, // 34: networking.gloo.solo.io.RouteTableReport.parent_route_tables:type_name -> common.gloo.solo.io.ObjectReference
	22, // 35: networking.gloo.solo.io.RouteTableReport.allowed_virtual_gateways:type_name -> common.gloo.solo.io.ObjectReference
	21, // 36: networking.gloo.solo.io.RouteTableReport.delegated_to_route_tables:type_name -> networking.gloo.solo.io.RouteTableReport.DelegatedRouteTableReference
	25, // 37: networking.gloo.solo.io.TLSRoute.TLSForwardToAction.destinations:type_name -> common.gloo.solo.io.DestinationReference
	35, // 38: networking.gloo.solo.io.GraphQLAction.Options.log_sensitive_info:type_name -> google.protobuf.BoolValue
	36, // 39: networking.gloo.solo.io.RouteTableReport.WorkspacesEntry.value:type_name -> common.gloo.solo.io.Report
	37, // 40: networking.gloo.solo.io.RouteTableReport.AppliedRoutePoliciesEntry.value:type_name -> common.gloo.solo.io.AppliedRoutePolicies
	22, // 41: networking.gloo.solo.io.RouteTableReport.DelegatedRouteTableReference.route_table:type_name -> common.gloo.solo.io.ObjectReference
	42, // [42:42] is the sub-list for method output_type
	42, // [42:42] is the sub-list for method input_type
	42, // [42:42] is the sub-list for extension type_name
	42, // [42:42] is the sub-list for extension extendee
	0,  // [0:42] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_init()
}
func file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_init() {
	if File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteTableSpec); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HTTPRoute); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TCPRoute); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TLSRoute); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GraphQLAction); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardToAction); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedirectAction); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirectResponseAction); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelegateAction); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteTableStatus); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteTableReport); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TLSRoute_TLSForwardToAction); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GraphQLAction_Options); i {
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
		file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[18].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteTableReport_DelegatedRouteTableReference); i {
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
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*HTTPRoute_ForwardTo)(nil),
		(*HTTPRoute_Delegate)(nil),
		(*HTTPRoute_Redirect)(nil),
		(*HTTPRoute_DirectResponse)(nil),
		(*HTTPRoute_Graphql)(nil),
	}
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*TCPRoute_ForwardTo)(nil),
	}
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*TLSRoute_ForwardTo)(nil),
	}
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*GraphQLAction_Schema)(nil),
		(*GraphQLAction_StitchedSchema)(nil),
	}
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*ForwardToAction_PathRewrite)(nil),
		(*ForwardToAction_RegexRewrite)(nil),
		(*ForwardToAction_HostRewrite)(nil),
		(*ForwardToAction_AutoHostRewrite)(nil),
	}
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*RedirectAction_PathRedirect)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   19,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_depIdxs,
		EnumInfos:         file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_enumTypes,
		MessageInfos:      file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto = out.File
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_rawDesc = nil
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_goTypes = nil
	file_github_com_solo_io_gloo_mesh_solo_apis_api_gloo_solo_io_networking_v2_route_table_proto_depIdxs = nil
}
