package proxy_syncer

import (
	"context"
	"fmt"
	"maps"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/solo-io/gloo/projects/gloo/constants"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"istio.io/api/label"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type CLA struct {
	*envoy_config_endpoint_v3.ClusterLoadAssignment
}

func (c CLA) ResourceName() string {
	return c.ClusterLoadAssignment.ClusterName
}
func (c CLA) Equals(in CLA) bool {
	return proto.Equal(c.ClusterLoadAssignment, in.ClusterLoadAssignment)
}

type nodeMetadata struct {
	name   string
	labels map[string]string
}

func (c nodeMetadata) ResourceName() string {
	return c.name
}
func (c nodeMetadata) Equals(in nodeMetadata) bool {
	return c.name == in.name && maps.Equal(c.labels, in.labels)
}

func NewNodeCollection(istioClient kube.Client) krt.Collection[nodeMetadata] {
	nodeClient := kclient.New[*corev1.Node](istioClient)
	nodes := krt.WrapClient(nodeClient, krt.WithName("Nodess"))
	return krt.NewCollection(nodes, func(kctx krt.HandlerContext, us *corev1.Node) *nodeMetadata {
		return &nodeMetadata{
			name:   us.Name,
			labels: us.Labels,
		}
	})
}

var _ krt.ResourceNamer = &upstream{}

type locality struct {
	region  string
	zone    string
	subzone string
}

func NewGlooK8sEndpoints(ctx context.Context, settings *v1.Settings, istioClient kube.Client, services krt.Collection[*corev1.Service], finalUpstreams krt.Collection[*upstream]) krt.Collection[CLA] {
	podClient := kclient.New[*corev1.Pod](istioClient)
	pods := krt.WrapClient(podClient, krt.WithName("Pods"))
	epClient := kclient.New[*corev1.Endpoints](istioClient)
	kubeEndpoints := krt.WrapClient(epClient, krt.WithName("Endpoints"))
	enableAutoMtls := settings.GetGloo().GetIstioOptions().GetEnableAutoMtls().GetValue()

	nodes := NewNodeCollection(istioClient)

	return krt.NewCollection(finalUpstreams, func(kctx krt.HandlerContext, us *upstream) *CLA {
		// TODO: log these
		var warnsToLog []string

		kubeUpstream, ok := us.GetUpstreamType().(*v1.Upstream_Kube)
		// only care about kube upstreams
		if !ok {
			return nil
		}
		spec := kubeUpstream.Kube
		kubeServicePort, singlePortService := findPortForService(kctx, services, spec)
		maybeEps := krt.FetchOne(kctx, kubeEndpoints, krt.FilterObjectName(types.NamespacedName{
			Namespace: spec.GetServiceNamespace(),
			Name:      spec.GetServiceName(),
		}))
		if maybeEps == nil {
			warnsToLog = append(warnsToLog, fmt.Sprintf("upstream %v: endpoints not found for service %v", us.GetMetadata().Ref().Key(), spec.GetServiceName()))
			return nil
		}
		eps := *maybeEps

		cla := createEndpoint(us.Upstream)
		lbEps := make(map[locality][]*envoy_config_endpoint_v3.LbEndpoint)
		for _, subset := range eps.Subsets {
			port := findFirstPortInEndpointSubsets(subset, singlePortService, kubeServicePort)
			if port == 0 {
				warnsToLog = append(warnsToLog, fmt.Sprintf("upstream %v: port %v not found for service %v in endpoint %v", us.Metadata.Ref().Key(), spec.GetServicePort(), spec.GetServiceName(), subset))
				continue
			}

			for _, addr := range subset.Addresses {
				var podName string
				podNamespace := eps.Namespace
				targetRef := addr.TargetRef
				if targetRef != nil {
					if targetRef.Kind == "Pod" {
						podName = targetRef.Name
						if targetRef.Namespace != "" {
							podNamespace = targetRef.Namespace
						}
					}
				}
				var podLabels map[string]string
				var nodeLabels map[string]string
				if podName != "" {
					maybePod := krt.FetchOne(kctx, pods, krt.FilterObjectName(types.NamespacedName{
						Namespace: podNamespace,
						Name:      podName,
					}))
					if maybePod != nil {
						pod := *maybePod
						podLabels = pod.Labels
						nodeName := pod.Spec.NodeName
						if nodeName != "" {
							maybeNode := krt.FetchOne(kctx, nodes, krt.FilterObjectName(types.NamespacedName{
								Name: nodeName,
							}))
							if maybeNode != nil {
								node := *maybeNode
								nodeLabels = node.labels
							}
						}
					}
				}
				ep, l := createLbEndpoint(addr.IP, port, podLabels, nodeLabels, enableAutoMtls)
				lbEps[l] = append(lbEps[l], ep)
			}
		}
		for locality, eps := range lbEps {
			var l *envoy_config_core_v3.Locality
			if locality.region != "" {
				l = &envoy_config_core_v3.Locality{
					Region:  locality.region,
					Zone:    locality.zone,
					SubZone: locality.subzone,
				}
			}

			cla.Endpoints = append(cla.Endpoints, &envoy_config_endpoint_v3.LocalityLbEndpoints{
				Locality:    l,
				LbEndpoints: eps,
			})
		}
		return &CLA{cla}

	}, krt.WithName("GlooEndpoints"))
}

func createLbEndpoint(address string, port uint32, podLabels, nodeLabels map[string]string, enableAutoMtls bool) (*envoy_config_endpoint_v3.LbEndpoint, locality) {
	// Don't get the metadata labels and filter metadata for the envoy load balancer based on the upstream, as this is not used
	// metadata := getLbMetadata(upstream, labels, "")
	// Get the metadata labels for the transport socket match if Istio auto mtls is enabled
	metadata := &envoy_config_core_v3.Metadata{
		FilterMetadata: map[string]*structpb.Struct{},
	}
	metadata = addIstioAutomtlsMetadata(metadata, podLabels, enableAutoMtls)
	// Don't add the annotations to the metadata - it's not documented so it's not coming
	// metadata = addAnnotations(metadata, addr.GetMetadata().GetAnnotations())

	region := nodeLabels[corev1.LabelTopologyRegion]
	zone := nodeLabels[corev1.LabelTopologyZone]
	subzone := nodeLabels[label.TopologySubzone.Name]
	l := locality{
		region:  region,
		zone:    zone,
		subzone: subzone,
	}

	if len(metadata.GetFilterMetadata()) == 0 {
		metadata = nil
	}

	return &envoy_config_endpoint_v3.LbEndpoint{
		Metadata: metadata,
		HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
			Endpoint: &envoy_config_endpoint_v3.Endpoint{
				Address: &envoy_config_core_v3.Address{
					Address: &envoy_config_core_v3.Address_SocketAddress{
						SocketAddress: &envoy_config_core_v3.SocketAddress{
							Protocol: envoy_config_core_v3.SocketAddress_TCP,
							Address:  address,
							PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
								PortValue: port,
							},
						},
					},
				},
			},
		},
	}, l
}

func addIstioAutomtlsMetadata(metadata *envoy_config_core_v3.Metadata, labels map[string]string, enableAutoMtls bool) *envoy_config_core_v3.Metadata {

	const EnvoyTransportSocketMatch = "envoy.transport_socket_match"
	if enableAutoMtls {
		if _, ok := labels[constants.IstioTlsModeLabel]; ok {
			metadata.GetFilterMetadata()[EnvoyTransportSocketMatch] = &structpb.Struct{
				Fields: map[string]*structpb.Value{
					constants.TLSModeLabelShortname: {
						Kind: &structpb.Value_StringValue{
							StringValue: constants.IstioMutualTLSModeLabel,
						},
					},
				},
			}
		}
	}
	return metadata
}

func createEndpoint(upstream *v1.Upstream) *envoy_config_endpoint_v3.ClusterLoadAssignment {
	clusterName := translator.UpstreamToClusterName(upstream.GetMetadata().Ref())
	return &envoy_config_endpoint_v3.ClusterLoadAssignment{
		ClusterName: getEndpointClusterName(clusterName, upstream),
	}
}

func findPortForService(kctx krt.HandlerContext, services krt.Collection[*corev1.Service], spec *kubeplugin.UpstreamSpec) (*corev1.ServicePort, bool) {
	maybeSvc := krt.FetchOne(kctx, services, krt.FilterObjectName(types.NamespacedName{
		Namespace: spec.GetServiceNamespace(),
		Name:      spec.GetServiceName(),
	}))
	if maybeSvc == nil {
		return nil, false
	}

	svc := *maybeSvc

	for _, port := range svc.Spec.Ports {
		if spec.GetServicePort() == uint32(port.Port) {
			return &port, len(svc.Spec.Ports) == 1
		}
	}

	return nil, false
}

func findFirstPortInEndpointSubsets(subset corev1.EndpointSubset, singlePortService bool, kubeServicePort *corev1.ServicePort) uint32 {
	var port uint32
	for _, p := range subset.Ports {
		// if the endpoint port is not named, it implies that
		// the kube service only has a single unnamed port as well.
		switch {
		case singlePortService:
			port = uint32(p.Port)
		case p.Name == kubeServicePort.Name:
			port = uint32(p.Port)
			break
		}
	}
	return port
}

// TODO: use exported version from translator?
func getEndpointClusterName(clusterName string, upstream *v1.Upstream) string {
	hash, err := upstream.Hash(nil)
	if err != nil {
		panic(err)
	}
	endpointClusterName := fmt.Sprintf("%s-%d", clusterName, hash)
	return endpointClusterName
}
