package printers

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	core1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/grpc_json"
	"io"
	"os"
	"sort"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/xdsinspection"
	plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws/ec2"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/olekukonko/tablewriter"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/cluster"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
)

// PrintUpstreams
func PrintUpstreams(upstreams v1.UpstreamList, outputType OutputType, xdsDump *xdsinspection.XdsDump) error {
	if outputType == KUBE_YAML {
		return PrintKubeCrdList(upstreams.AsInputResources(), v1.UpstreamCrd)
	}

	var printList []interface{}
	for _, upstream := range upstreams {
		if upstream.GetKube().GetServiceSpec().GetGrpcJsonTranscoder() != nil {
			printList = append(printList, convertKubeGrpcJsonTranscoderUpstream(upstream))
			continue
		}
		printList = append(printList, upstream)
	}

	return cliutils.PrintList(outputType.String(), "", printList,
		func(data interface{}, w io.Writer) error {
			UpstreamTable(xdsDump, data.(v1.UpstreamList), w)
			return nil
		}, os.Stdout)
}

// PrintTable prints upstreams using tables to io.Writer
func UpstreamTable(xdsDump *xdsinspection.XdsDump, upstreams []*v1.Upstream, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Upstream", "type", "status", "details"})

	for _, us := range upstreams {
		name := us.GetMetadata().GetName()
		s := upstreamStatus(us)

		u := upstreamType(us)
		details := upstreamDetails(us, xdsDump)

		if len(details) == 0 {
			details = []string{""}
		}
		for i, line := range details {
			if i == 0 {
				table.Append([]string{name, u, s, line})
			} else {
				table.Append([]string{"", "", "", line})
			}
		}
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func upstreamStatus(us *v1.Upstream) string {
	return AggregateNamespacedStatuses(us.GetNamespacedStatuses(), func(status *core.Status) string {
		return status.GetState().String()
	})
}

func upstreamType(up *v1.Upstream) string {
	if up == nil {
		return "Invalid"
	}

	switch up.GetUpstreamType().(type) {
	case *v1.Upstream_Aws:
		return "AWS Lambda"
	case *v1.Upstream_Azure:
		return "Azure"
	case *v1.Upstream_Consul:
		return "Consul"
	case *v1.Upstream_AwsEc2:
		return "AWS EC2"
	case *v1.Upstream_Kube:
		return "Kubernetes"
	case *v1.Upstream_Static:
		return "Static"
	case *v1.Upstream_Gcp:
		return "GCP"
	default:
		return "Unknown"
	}
}

func upstreamDetails(up *v1.Upstream, xdsDump *xdsinspection.XdsDump) []string {
	if up == nil {
		return []string{"invalid: upstream was nil"}
	}

	var details []string
	add := func(s ...string) {
		details = append(details, s...)
	}
	switch usType := up.GetUpstreamType().(type) {
	case *v1.Upstream_Aws:
		var functions []string
		for _, fn := range usType.Aws.GetLambdaFunctions() {
			functions = append(functions, fn.GetLambdaFunctionName())
		}

		add(
			fmt.Sprintf("region: %v", usType.Aws.GetRegion()),
			fmt.Sprintf("secret: %s", stringifyKey(usType.Aws.GetSecretRef())),
		)
		for i := range functions {
			if i == 0 {
				add("functions:")
			}
			add(fmt.Sprintf("- %v", functions[i]))
		}
	case *v1.Upstream_AwsEc2:
		add(
			fmt.Sprintf("role:           %v", usType.AwsEc2.GetRoleArn()),
			fmt.Sprintf("uses public ip: %v", usType.AwsEc2.GetPublicIp()),
			fmt.Sprintf("port:           %v", usType.AwsEc2.GetPort()),
		)
		add(getEc2TagFiltersString(usType.AwsEc2.GetFilters())...)
		instances := xdsDump.GetEc2InstancesForUpstream(up.GetMetadata().Ref())
		add(
			"EC2 Instance Ids:",
		)
		add(
			instances...,
		)
	case *v1.Upstream_Azure:
		var functions []string
		for _, fn := range usType.Azure.GetFunctions() {
			functions = append(functions, fn.GetFunctionName())
		}
		add(
			fmt.Sprintf("function app name: %v", usType.Azure.GetFunctionAppName()),
			fmt.Sprintf("secret: %s", stringifyKey(usType.Azure.GetSecretRef())),
		)

		for i := range functions {
			if i == 0 {
				add("functions:")
			}
			add(fmt.Sprintf("- %v", functions[i]))
		}
	case *v1.Upstream_Consul:
		add(
			fmt.Sprintf("svc name: %v", usType.Consul.GetServiceName()),
			fmt.Sprintf("svc tags: %v", usType.Consul.GetServiceTags()),
		)
		if usType.Consul.GetServiceSpec() != nil {
			add(linesForServiceSpec(usType.Consul.GetServiceSpec())...)
		}
	case *v1.Upstream_Kube:
		add(
			fmt.Sprintf("svc name:      %v", usType.Kube.GetServiceName()),
			fmt.Sprintf("svc namespace: %v", usType.Kube.GetServiceNamespace()),
			fmt.Sprintf("port:          %v", usType.Kube.GetServicePort()),
		)
		if usType.Kube.GetServiceSpec() != nil {
			add(linesForServiceSpec(usType.Kube.GetServiceSpec())...)
		}
	case *v1.Upstream_Static:
		for i := range usType.Static.GetHosts() {
			if i == 0 {
				add("hosts:")
			}
			add(fmt.Sprintf("- %v:%v", usType.Static.GetHosts()[i].GetAddr(), usType.Static.GetHosts()[i].GetPort()))
		}
		if usType.Static.GetServiceSpec() != nil {
			add(linesForServiceSpec(usType.Static.GetServiceSpec())...)
		}
	case *v1.Upstream_Gcp:
		add(fmt.Sprintf("host: %v", usType.Gcp.GetHost()))
		if usType.Gcp.GetAudience() != "" {
			add(fmt.Sprintf("host: %v", usType.Gcp.GetAudience()))
		}

	}
	add("")
	return details
}

func linesForServiceSpec(serviceSpec *plugins.ServiceSpec) []string {
	var spec []string
	add := func(s ...string) {
		spec = append(spec, s...)
	}
	switch plug := serviceSpec.GetPluginType().(type) {
	case *plugins.ServiceSpec_Rest:
		add("REST service:")
		var functions []string
		for restFunc, transform := range plug.Rest.GetTransformations() {
			var method, path string
			methodP, ok := transform.GetHeaders()[":method"]
			if ok {
				method = methodP.GetText()
			}
			pathP, ok := transform.GetHeaders()[":path"]
			if ok {
				path = pathP.GetText()
			}
			if false {
				//TODO ilackarms: save it for -o wide
				functions = append(functions, fmt.Sprintf("- %v (%v %v)", restFunc, method, path))
			}
			functions = append(functions, fmt.Sprintf("- %v", restFunc))
		}
		// needed because map
		sort.Strings(functions)

		for i := range functions {
			if i == 0 {
				add("functions:")
			}
			add(functions[i])
		}
	case *plugins.ServiceSpec_Grpc:
		add("gRPC service:")
		for _, grpcService := range plug.Grpc.GetGrpcServices() {
			add(fmt.Sprintf("  %v", grpcService.GetServiceName()))
			for _, fn := range grpcService.GetFunctionNames() {
				add(fmt.Sprintf("  - %v", fn))
			}
		}
	case *plugins.ServiceSpec_GrpcJsonTranscoder:
		add("gRPC service:")
		descriptorBin := plug.GrpcJsonTranscoder.GetProtoDescriptorBin()
		for _, grpcService := range plug.GrpcJsonTranscoder.GetServices() {
			add(fmt.Sprintf("  %v", grpcService))
			methodDescriptors := getMethodDescriptors(grpcService, descriptorBin)
			for i := 0; i < methodDescriptors.Len(); i++ {
				add(fmt.Sprintf("  - %v", methodDescriptors.Get(i).Name()))
			}
		}
	}
	return spec
}

func getMethodDescriptors(service string, descriptorSet []byte) protoreflect.MethodDescriptors {
	fds := &descriptor.FileDescriptorSet{}
	err := proto.Unmarshal(descriptorSet, fds)
	if err != nil {
		fmt.Println("unable to unmarshal descriptor")
		return nil
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		fmt.Println("unable to create proto registry files")
		return nil
	}
	descriptor, err := files.FindDescriptorByName(protoreflect.FullName(service))
	if err != nil {
		fmt.Println("unable to fin descriptor")
		return nil
	}
	serviceDescriptor, ok := descriptor.(protoreflect.ServiceDescriptor)
	if !ok {
		fmt.Println("unable to decode service descriptor")
		return nil
	}
	return serviceDescriptor.Methods()
}

// stringifyKey for a resource likely could be done more nicely with spew
// or a better accessor but minimal this avoids panicing on nested references to nils
func stringifyKey(plausiblyNilRef *core.ResourceRef) string {

	if plausiblyNilRef == nil {
		return "<None>"
	}
	return plausiblyNilRef.Key()

}

func getEc2TagFiltersString(filters []*ec2.TagFilter) []string {
	var out []string
	add := func(s ...string) {
		out = append(out, s...)
	}

	var kFilters []*ec2.TagFilter_Key
	var kvFilters []*ec2.TagFilter_KvPair
	for _, f := range filters {
		switch x := f.GetSpec().(type) {
		case *ec2.TagFilter_Key:
			kFilters = append(kFilters, x)
		case *ec2.TagFilter_KvPair_:
			kvFilters = append(kvFilters, x.KvPair)
		}
	}
	if len(kFilters) == 0 {
		add(fmt.Sprintf("key filters: (none)"))
	} else {
		add(fmt.Sprintf("key filters:"))
		for _, f := range kFilters {
			add(fmt.Sprintf("- %v", f.Key))
		}
	}
	if len(kvFilters) == 0 {
		add(fmt.Sprintf("key-value filters: (none)"))
	} else {
		add(fmt.Sprintf("key-value filters:"))
		for _, f := range kvFilters {
			add(fmt.Sprintf("- %v: %v", f.GetKey(), f.GetValue()))
		}
	}
	return out
}

func convertKubeGrpcJsonTranscoderUpstream(us *v1.Upstream) *GrpcJsonTranscoderUpstream {
	return &GrpcJsonTranscoderUpstream{
		UpstreamType: &UpstreamType_GrpcJsonTranscoderKubeUpstream{
			Kube: &GrpcJsonTranscoderKubeUpstream{
				GrpcJsonTranscoder: &GrpcJsonTranscoderWithMethods{
					DescriptorSet: &GrpcJsonTranscoder_ProtoDescriptorBin{
						ProtoDescriptorBin: us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().GetProtoDescriptorBin(),
					},

					Services:                     us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().Services,
					PrintOptions:                 us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().PrintOptions,
					MatchIncomingRequestRoute:    us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().MatchIncomingRequestRoute,
					IgnoredQueryParameters:       us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().IgnoredQueryParameters,
					AutoMapping:                  us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().AutoMapping,
					IgnoreUnknownQueryParameters: us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().IgnoreUnknownQueryParameters,
					ConvertGrpcStatus:            us.GetKube().GetServiceSpec().GetGrpcJsonTranscoder().ConvertGrpcStatus,

					//Methods: map[string][]string{"foo": {"bar", "baz"}},
				},

				ServiceName:      us.GetKube().ServiceName,
				ServiceNamespace: us.GetKube().ServiceNamespace,
				ServicePort:      us.GetKube().ServicePort,
				Selector:         us.GetKube().Selector,
				SubsetSpec:       us.GetKube().SubsetSpec,
			},
		},
		NamespacedStatuses: us.NamespacedStatuses,
		Metadata:           us.Metadata,
		DiscoveryMetadata:  us.DiscoveryMetadata,
		SslConfig:          us.SslConfig,
		CircuitBreakers:    us.CircuitBreakers,
		LoadBalancerConfig: us.LoadBalancerConfig,
		HealthChecks:       us.HealthChecks,
		OutlierDetection:   us.OutlierDetection,
		Failover:           us.Failover,
		ConnectionConfig:   us.ConnectionConfig,
		ProtocolSelection:  us.ProtocolSelection,
	}
}

type GrpcJsonTranscoderUpstream struct {
	// Note to developers: new Upstream plugins must be added to this oneof field
	// to be usable by Gloo. (plugins currently need to be compiled into Gloo)
	//
	// Types that are assignable to UpstreamType:
	//
	//	*Upstream_Kube
	//	*Upstream_Static
	//	*Upstream_Consul
	UpstreamType serviceSpecUpstreamType `protobuf_oneof:"upstream_type"`

	NamespacedStatuses *core.NamespacedStatuses `protobuf:"bytes,23,opt,name=namespaced_statuses,json=namespacedStatuses,proto3" json:"namespaced_statuses,omitempty"`
	// Metadata contains the object metadata for this resource
	Metadata *core.Metadata `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Upstreams and their configuration can be automatically by Gloo Discovery
	// if this upstream is created or modified by Discovery, metadata about the operation will be placed here.
	DiscoveryMetadata *v1.DiscoveryMetadata `protobuf:"bytes,3,opt,name=discovery_metadata,json=discoveryMetadata,proto3" json:"discovery_metadata,omitempty"`
	// SslConfig contains the options necessary to configure envoy to originate TLS to an upstream.
	SslConfig *ssl.UpstreamSslConfig `protobuf:"bytes,4,opt,name=ssl_config,json=sslConfig,proto3" json:"ssl_config,omitempty"`
	// Circuit breakers for this upstream. if not set, the defaults ones from the Gloo settings will be used.
	// if those are not set, [envoy's defaults](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/circuit_breaker.proto#envoy-api-msg-cluster-circuitbreakers)
	// will be used.
	CircuitBreakers *v1.CircuitBreakerConfig `protobuf:"bytes,5,opt,name=circuit_breakers,json=circuitBreakers,proto3" json:"circuit_breakers,omitempty"`
	// Settings for the load balancer that sends requests to the Upstream. The load balancing method is set to round robin by default.
	LoadBalancerConfig *v1.LoadBalancerConfig    `protobuf:"bytes,6,opt,name=load_balancer_config,json=loadBalancerConfig,proto3" json:"load_balancer_config,omitempty"`
	HealthChecks       []*core1.HealthCheck      `protobuf:"bytes,8,rep,name=health_checks,json=healthChecks,proto3" json:"health_checks,omitempty"`
	OutlierDetection   *cluster.OutlierDetection `protobuf:"bytes,9,opt,name=outlier_detection,json=outlierDetection,proto3" json:"outlier_detection,omitempty"`
	Failover           *v1.Failover              `protobuf:"bytes,18,opt,name=failover,proto3" json:"failover,omitempty"`
	// HTTP/1 connection configurations
	ConnectionConfig *v1.ConnectionConfig `protobuf:"bytes,7,opt,name=connection_config,json=connectionConfig,proto3" json:"connection_config,omitempty"`
	// Determines how Envoy selects the protocol used to speak to upstream hosts.
	ProtocolSelection v1.Upstream_ClusterProtocolSelection `protobuf:"varint,25,opt,name=protocol_selection,json=protocolSelection,proto3,enum=gloo.solo.io.Upstream_ClusterProtocolSelection" json:"protocol_selection,omitempty"`
	// Use http2 when communicating with this upstream
	// this field is evaluated `true` for upstreams
	// with a grpc service spec. otherwise defaults to `false`
	UseHttp2 *wrappers.BoolValue `protobuf:"bytes,10,opt,name=use_http2,json=useHttp2,proto3" json:"use_http2,omitempty"`
	// (UInt32Value) Initial stream-level flow-control window size.
	// Valid values range from 65535 (2^16 - 1, HTTP/2 default) to 2147483647 (2^31 - 1, HTTP/2 maximum)
	// and defaults to 268435456 (256 * 1024 * 1024).
	// NOTE: 65535 is the initial window size from HTTP/2 spec.
	// We only support increasing the default window size now, so it’s also the minimum.
	// This field also acts as a soft limit on the number of bytes Envoy will buffer per-stream
	// in the HTTP/2 codec buffers. Once the buffer reaches this pointer,
	// watermark callbacks will fire to stop the flow of data to the codec buffers.
	// Requires UseHttp2 to be true to be acknowledged.
	InitialStreamWindowSize *wrappers.UInt32Value `protobuf:"bytes,19,opt,name=initial_stream_window_size,json=initialStreamWindowSize,proto3" json:"initial_stream_window_size,omitempty"`
	// (UInt32Value) Similar to initial_stream_window_size, but for connection-level flow-control window.
	// Currently, this has the same minimum/maximum/default as initial_stream_window_size.
	// Requires UseHttp2 to be true to be acknowledged.
	InitialConnectionWindowSize *wrappers.UInt32Value `protobuf:"bytes,20,opt,name=initial_connection_window_size,json=initialConnectionWindowSize,proto3" json:"initial_connection_window_size,omitempty"`
	// (UInt32Value) Maximum concurrent streams allowed for peer on one HTTP/2 connection.
	// Valid values range from 1 to 2147483647 (2^31 - 1) and defaults to 2147483647.
	// Requires UseHttp2 to be true to be acknowledged.
	MaxConcurrentStreams *wrappers.UInt32Value `protobuf:"bytes,24,opt,name=max_concurrent_streams,json=maxConcurrentStreams,proto3" json:"max_concurrent_streams,omitempty"`
	// Allows invalid HTTP messaging and headers. When this option is disabled (default), then
	// the whole HTTP/2 connection is terminated upon receiving invalid HEADERS frame. However,
	// when this option is enabled, only the offending stream is terminated.
	//
	// This overrides any HCM :ref:`stream_error_on_invalid_http_messaging
	// <envoy_v3_api_field_extensions.filters.network.http_connection_manager.v3.HttpConnectionManager.stream_error_on_invalid_http_message>`
	//
	// See `RFC7540, sec. 8.1 <https://tools.ietf.org/html/rfc7540#section-8.1>`_ for details.
	OverrideStreamErrorOnInvalidHttpMessage *wrappers.BoolValue `protobuf:"bytes,26,opt,name=override_stream_error_on_invalid_http_message,json=overrideStreamErrorOnInvalidHttpMessage,proto3" json:"override_stream_error_on_invalid_http_message,omitempty"`
	// Tells envoy that the upstream is an HTTP proxy (e.g., another proxy in a DMZ) that supports HTTP Connect.
	// This configuration sets the hostname used as part of the HTTP Connect request.
	// For example, setting to: host.com:443 and making a request routed to the upstream such as `curl <envoy>:<port>/v1`
	// would result in the following request:
	//
	//	CONNECT host.com:443 HTTP/1.1
	//	host: host.com:443
	//
	//	GET /v1 HTTP/1.1
	//	host: <envoy>:<port>
	//	user-agent: curl/7.64.1
	//	accept: */*
	//
	// Note: if setting this field to a hostname rather than IP:PORT, you may want to also set `host_rewrite` on the route
	HttpProxyHostname *wrappers.StringValue `protobuf:"bytes,21,opt,name=http_proxy_hostname,json=httpProxyHostname,proto3" json:"http_proxy_hostname,omitempty"`
	// HttpConnectSslConfig contains the options necessary to configure envoy to originate TLS to an HTTP Connect proxy.
	// If you also want to ensure the bytes proxied by the HTTP Connect proxy are encrypted, you should also
	// specify `ssl_config`.
	HttpConnectSslConfig *ssl.UpstreamSslConfig `protobuf:"bytes,27,opt,name=http_connect_ssl_config,json=httpConnectSslConfig,proto3" json:"http_connect_ssl_config,omitempty"`
	// HttpConnectHeaders specifies the headers sent with the initial HTTP Connect request.
	HttpConnectHeaders []*v1.HeaderValue `protobuf:"bytes,28,rep,name=http_connect_headers,json=httpConnectHeaders,proto3" json:"http_connect_headers,omitempty"`
	// (bool) If set to true, Envoy will ignore the health value of a host when processing its removal from service discovery.
	// This means that if active health checking is used, Envoy will not wait for the endpoint to go unhealthy before removing it.
	IgnoreHealthOnHostRemoval *wrappers.BoolValue `protobuf:"bytes,22,opt,name=ignore_health_on_host_removal,json=ignoreHealthOnHostRemoval,proto3" json:"ignore_health_on_host_removal,omitempty"`
	// If set to true, Service Discovery update period will be triggered once the TTL is expired.
	// If minimum TTL of all records is 0 then dns_refresh_rate will be used.
	RespectDnsTtl *wrappers.BoolValue `protobuf:"bytes,29,opt,name=respect_dns_ttl,json=respectDnsTtl,proto3" json:"respect_dns_ttl,omitempty"`
	// Service Discovery DNS Refresh Rate.
	// Minimum value is 1 ms. Values below the minimum are considered invalid.
	// Only valid for STRICT_DNS and LOGICAL_DNS cluster types. All other cluster types are considered invalid.
	DnsRefreshRate *duration.Duration `protobuf:"bytes,30,opt,name=dns_refresh_rate,json=dnsRefreshRate,proto3" json:"dns_refresh_rate,omitempty"`
	// Proxy Protocol Version to add when communicating with the upstream.
	// If unset will not wrap the transport socket.
	// These are of the format "V1" or "V2"
	ProxyProtocolVersion *wrappers.StringValue `protobuf:"bytes,31,opt,name=proxy_protocol_version,json=proxyProtocolVersion,proto3" json:"proxy_protocol_version,omitempty"`
	// Preconnect policy for the cluster
	// Aligns as closely as possible with https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/cluster.proto#envoy-v3-api-msg-config-cluster-v3-cluster-preconnectpolicy
	// This is not recommended for use unless you are sure you need it.
	// In most cases preconnect hurts more than it helps.
	PreconnectPolicy *v1.PreconnectPolicy `protobuf:"bytes,32,opt,name=preconnect_policy,json=preconnectPolicy,proto3" json:"preconnect_policy,omitempty"`
	// If set to true, the proxy will not allow automatic mTLS detection for Istio upstreams.
	// Defaults to false.
	DisableIstioAutoMtls *wrappers.BoolValue `protobuf:"bytes,33,opt,name=disable_istio_auto_mtls,json=disableIstioAutoMtls,proto3" json:"disable_istio_auto_mtls,omitempty"`
}

type GrpcJsonTranscoderServiceSpec struct {
	GrpcJsonTranscoder GrpcJsonTranscoderWithMethods `protobuf:"bytes,3,opt,name=grpc_json_transcoder,json=grpcJsonTranscoder,proto3,oneof"`
}

type UpstreamType_GrpcJsonTranscoderKubeUpstream struct {
	Kube *GrpcJsonTranscoderKubeUpstream `protobuf:"bytes,11,opt,name=kube,proto3,oneof"`
}

func (ut *UpstreamType_GrpcJsonTranscoderKubeUpstream) isServiceSpecUpstreamType() {}

type GrpcJsonTranscoderKubeUpstream struct {
	ServiceName string `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// The namespace where the Service lives
	ServiceNamespace string `protobuf:"bytes,2,opt,name=service_namespace,json=serviceNamespace,proto3" json:"service_namespace,omitempty"`
	// The access port of the kubernetes service is listening. This port is used by Gloo to look up the corresponding
	// port on the pod for routing.
	ServicePort uint32 `protobuf:"varint,3,opt,name=service_port,json=servicePort,proto3" json:"service_port,omitempty"`
	// Allows finer-grained filtering of pods for the Upstream. Gloo will select pods based on their labels if
	// any are provided here.
	// (see [Kubernetes labels and selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)
	Selector map[string]string `protobuf:"bytes,4,rep,name=selector,proto3" json:"selector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// An optional Service Spec describing the service listening at this address
	ServiceSpec *options.ServiceSpec `protobuf:"bytes,5,opt,name=service_spec,json=serviceSpec,proto3" json:"service_spec,omitempty"`
	// Subset configuration. For discovery sources that has labels (like kubernetes). this
	// configuration allows you to partition the upstream to a set of subsets.
	// for each unique set of keys and values, a subset will be created.
	SubsetSpec *options.SubsetSpec `protobuf:"bytes,6,opt,name=subset_spec,json=subsetSpec,proto3" json:"subset_spec,omitempty"`

	GrpcJsonTranscoder *GrpcJsonTranscoderWithMethods
}

type GrpcJsonTranscoderWithMethods struct {
	DescriptorSet isGrpcJsonTranscoder_DescriptorSet `protobuf_oneof:"descriptor_set"`

	Services                     []string                                   `protobuf:"bytes,2,rep,name=services,proto3" json:"services,omitempty"`
	PrintOptions                 *grpc_json.GrpcJsonTranscoder_PrintOptions `protobuf:"bytes,3,opt,name=print_options,json=printOptions,proto3" json:"print_options,omitempty"`
	MatchIncomingRequestRoute    bool                                       `protobuf:"varint,5,opt,name=match_incoming_request_route,json=matchIncomingRequestRoute,proto3" json:"match_incoming_request_route,omitempty"`
	IgnoredQueryParameters       []string                                   `protobuf:"bytes,6,rep,name=ignored_query_parameters,json=ignoredQueryParameters,proto3" json:"ignored_query_parameters,omitempty"`
	AutoMapping                  bool                                       `protobuf:"varint,7,opt,name=auto_mapping,json=autoMapping,proto3" json:"auto_mapping,omitempty"`
	IgnoreUnknownQueryParameters bool                                       `protobuf:"varint,8,opt,name=ignore_unknown_query_parameters,json=ignoreUnknownQueryParameters,proto3" json:"ignore_unknown_query_parameters,omitempty"`
	ConvertGrpcStatus            bool                                       `protobuf:"varint,9,opt,name=convert_grpc_status,json=convertGrpcStatus,proto3" json:"convert_grpc_status,omitempty"`

	Methods map[string][]string `protobuf:"bytes,20,opt,name=methods,json=methods,proto3" json:"methods,omitempty"`
}

type GrpcJsonTranscoder_ProtoDescriptorBin struct {
	// Supplies the binary content of the [proto descriptor set](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter#config-grpc-json-generate-proto-descriptor-set)
	// for the gRPC services.
	// Note: in yaml, this must be provided as a base64 standard encoded string; yaml can't handle binary bytes.
	ProtoDescriptorBin []byte `protobuf:"bytes,4,opt,name=proto_descriptor_bin,json=protoDescriptorBin,proto3,oneof"`
}

func (*GrpcJsonTranscoder_ProtoDescriptorBin) isGrpcJsonTranscoder_DescriptorSet() {}

type serviceSpecUpstreamType interface {
	isServiceSpecUpstreamType()
}

type isGrpcJsonTranscoder_DescriptorSet interface {
	isGrpcJsonTranscoder_DescriptorSet()
}
