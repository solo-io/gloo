package proxy_syncer_test

import (
	"context"
	"testing"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway2/krtcollections"
	. "github.com/solo-io/gloo/projects/gateway2/proxy_syncer"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/krt/krttest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEndpoints(t *testing.T) {
	g := gomega.NewWithT(t)

	inputs := []any{
		&corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "name",
				Namespace: "ns",
			},
			Spec: corev1.PodSpec{
				NodeName: "node",
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				PodIP: "1.2.3.4",
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node",
				Labels: map[string]string{
					corev1.LabelTopologyRegion: "region",
					corev1.LabelTopologyZone:   "zone",
				},
			},
		},
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "svc",
				Namespace: "ns",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Ports: []corev1.EndpointPort{
						{
							Name: "http",
							Port: 8080,
						},
					},
					Addresses: []corev1.EndpointAddress{
						{
							IP: "1.2.3.4",
							TargetRef: &corev1.ObjectReference{
								Kind:      "Pod",
								Name:      "name",
								Namespace: "ns",
							},
						},
					},
				},
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "svc",
				Namespace: "ns",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: "http",
						Port: 8080,
					},
				},
			},
		},
	}
	// input
	us := UpstreamWrapper{
		Upstream: &gloov1.Upstream{
			Metadata: &core.Metadata{Name: "name", Namespace: "ns"},
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &kubernetes.UpstreamSpec{
					ServiceName:      "svc",
					ServiceNamespace: "ns",
					ServicePort:      8080,
				},
			},
		},
	}
	// output
	emd := EndpointWithMd{
		LbEndpoint: &endpointv3.LbEndpoint{
			HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
				Endpoint: &endpointv3.Endpoint{
					Address: &envoy_config_core_v3.Address{
						Address: &envoy_config_core_v3.Address_SocketAddress{
							SocketAddress: &envoy_config_core_v3.SocketAddress{
								Address: "1.2.3.4",
								PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
									PortValue: 8080,
								},
							},
						},
					},
				},
			},
		},
		EndpointMd: EndpointMetadata{
			Labels: map[string]string{
				corev1.LabelTopologyRegion: "region",
				corev1.LabelTopologyZone:   "zone",
			},
		},
	}
	result := NewEndpointsForUpstream(us)
	result.Add(krtcollections.PodLocality{
		Region: "region",
		Zone:   "zone",
	}, emd)

	mock := krttest.NewMock(t, inputs)
	nodes := krtcollections.NewNodeMetadataCollection(krttest.GetMockCollection[*corev1.Node](mock))
	pods := krtcollections.NewLocalityPodsCollection(nodes, krttest.GetMockCollection[*corev1.Pod](mock))
	pods.Synced().WaitUntilSynced(context.Background().Done())
	ei := EndpointsInputs{
		Upstreams:      krttest.GetMockCollection[UpstreamWrapper](mock),
		Endpoints:      krttest.GetMockCollection[*corev1.Endpoints](mock),
		Pods:           pods,
		EnableAutoMtls: false,
		Services:       krttest.GetMockCollection[*corev1.Service](mock),
	}

	builder := TransformUpstreamsBuilder(context.Background(), ei)

	eps := builder(krt.TestingDummyContext{}, us)
	g.Expect(eps).To(BeEquivalentTo(result))

}
