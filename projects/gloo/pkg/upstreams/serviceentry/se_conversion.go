package serviceentry

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/types/known/wrapperspb"

	networking "istio.io/api/networking/v1beta1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
)

const UpstreamNamePrefix = "istio-se:"

func ConvertUpstreamList(servicEntries []*networkingclient.ServiceEntry) v1.UpstreamList {
	var out v1.UpstreamList
	for _, se := range servicEntries {
		out = append(out, ConvertUpstreams(se)...)
	}
	return out
}

func ConvertUpstreams(se *networkingclient.ServiceEntry) []*v1.Upstream {
	// TODO How to handle these fields:
	// - Reslution: DNS; will gloo automatically use an envoy DNS cluster when it sees non IP endpoints?
	// - Reslution: None; I don't think we can do original_dst/passthrough.
	// - Location: MESH_INTERNAL/MESH_EXTERNAL. Decide whether to use mTLS
	// - exportTo: we're not respecting this but we should
	var out []*v1.Upstream
	for _, port := range se.Spec.Ports {
		us := &v1.Upstream{
			Metadata: &core.Metadata{
				Name:        UpstreamNamePrefix + se.Name,
				Namespace:   se.Namespace,
				Cluster:     "", // can we populate this for the local cluster?
				Labels:      se.Labels,
				Annotations: se.Annotations,
			},
		}

		if se.Spec.WorkloadSelector != nil {
			us.UpstreamType = buildKube(se, port)
		} else {
			us.UpstreamType = buildStatic(se, port)
		}

		out = append(out, us)
	}
	return out
}

func buildKube(se *networkingclient.ServiceEntry, port *networking.ServicePort) *v1.Upstream_Kube {
	return &v1.Upstream_Kube{
		Kube: &kubernetes.UpstreamSpec{
			ServiceName:      se.Name,
			ServiceNamespace: se.Namespace,
			ServicePort:      port.Number,
			Selector:         se.Spec.WorkloadSelector.Labels,
		},
	}
}

func buildStatic(se *networkingclient.ServiceEntry, port *networking.ServicePort) *v1.Upstream_Static {
	var hosts []*static.Host
	for _, ep := range se.Spec.Endpoints {
		hosts = append(hosts, &static.Host{
			Addr: ep.Address,
			Port: targetPort(port.Number, port.TargetPort, wePort(ep, port)),
			// SniAddr:             "",
			// LoadBalancingWeight: &wrapperspb.UInt32Value{},
			// HealthCheckConfig:   &static.Host_HealthCheckConfig{},
			// Metadata:            map[string]*structpb.Struct{},
		})
	}
	useTLS := false
	if port.Protocol == "HTTPS" || port.Protocol == "TLS" {
		useTLS = true
	}
	return &v1.Upstream_Static{
		Static: &static.UpstreamSpec{
			Hosts:  hosts,
			UseTls: &wrapperspb.BoolValue{Value: useTLS},
			// AutoSniRewrite: &wrapperspb.BoolValue{},
		},
	}
}

func wePort(we *networking.WorkloadEntry, svcPort *networking.ServicePort) uint32 {
	if we.Ports == nil {
		return 0
	}
	return we.Ports[svcPort.Name]
}

func targetPort(base, svc, ep uint32) uint32 {
	if ep > 0 {
		return ep
	}
	if svc > 0 {
		return svc
	}
	return base
}
