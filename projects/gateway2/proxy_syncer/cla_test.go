package proxy_syncer_test

import (
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	duration "github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway2/krtcollections"
	. "github.com/solo-io/gloo/projects/gateway2/proxy_syncer"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTransformsEndpoint(t *testing.T) {
	g := gomega.NewWithT(t)
	us := krtcollections.UpstreamWrapper{
		Inner: &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "name",
				Namespace: "ns",
			},
		},
	}
	efu := krtcollections.NewEndpointsForUpstream(us, nil)
	efu.Add(krtcollections.PodLocality{}, krtcollections.EndpointWithMd{
		LbEndpoint: &endpointv3.LbEndpoint{
			HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
				Endpoint: &endpointv3.Endpoint{
					Address: &corev3.Address{
						Address: &corev3.Address_Pipe{Pipe: &corev3.Pipe{Path: "a"}},
					},
				},
			},
		},
		EndpointMd: krtcollections.EndpointMetadata{},
	})

	envoyResources := TransformEndpointToResources(*efu)
	cla := envoyResources.Endpoints.ResourceProto().(*endpointv3.ClusterLoadAssignment)
	locality := cla.Endpoints[0].Locality
	g.Expect(locality).To(gomega.BeNil())
	g.Expect(cla.Endpoints[0].GetLbEndpoints()[0].GetEndpoint().GetAddress().GetPipe().GetPath()).To(gomega.Equal("a"))
}

func TestTransformsEndpointsWithLocality(t *testing.T) {
	g := gomega.NewWithT(t)
	us := krtcollections.UpstreamWrapper{
		Inner: &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "name",
				Namespace: "ns",
			},
		},
	}
	efu := krtcollections.NewEndpointsForUpstream(us, nil)
	efu.Add(krtcollections.PodLocality{Region: "R1"}, krtcollections.EndpointWithMd{
		LbEndpoint: &endpointv3.LbEndpoint{
			HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
				Endpoint: &endpointv3.Endpoint{
					Address: &corev3.Address{
						Address: &corev3.Address_Pipe{Pipe: &corev3.Pipe{Path: "a"}},
					},
				},
			},
		},
		EndpointMd: krtcollections.EndpointMetadata{},
	})
	efu.Add(krtcollections.PodLocality{Region: "R2"}, krtcollections.EndpointWithMd{
		LbEndpoint: &endpointv3.LbEndpoint{
			HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
				Endpoint: &endpointv3.Endpoint{
					Address: &corev3.Address{
						Address: &corev3.Address_Pipe{Pipe: &corev3.Pipe{Path: "a"}},
					},
				},
			},
		},
		EndpointMd: krtcollections.EndpointMetadata{},
	})

	envoyResources := TransformEndpointToResources(*efu)
	cla := envoyResources.Endpoints.ResourceProto().(*endpointv3.ClusterLoadAssignment)
	g.Expect(cla.Endpoints).To(gomega.HaveLen(2))
	locality1 := cla.Endpoints[0].Locality
	locality2 := cla.Endpoints[1].Locality
	regions := []string{locality1.Region, locality2.Region}
	g.Expect(regions).To(gomega.ContainElement("R1"))
	g.Expect(regions).To(gomega.ContainElement("R2"))
}

func TestTranslatesDestrulesFailoverPriority(t *testing.T) {
	testCases := []struct {
		name               string
		defaultNetwork     string
		clientLabels       map[string]string
		endpoints          map[string]map[string]string
		failoverPriority   []string
		expectedPriorities map[string]uint32
	}{
		{
			name:         "region failoverPriority",
			clientLabels: map[string]string{"topology.kubernetes.io/region": "R1"},
			endpoints: map[string]map[string]string{
				"a": {"topology.kubernetes.io/region": "R1"},
				"b": {"topology.kubernetes.io/region": "R2"},
			},
			failoverPriority: []string{
				"topology.kubernetes.io/region",
			},
			expectedPriorities: map[string]uint32{
				"a": 0,
				"b": 1,
			},
		},
		{
			name:         "network failoverPriority with explicit labels",
			clientLabels: map[string]string{"topology.istio.io/network": "local-network-1"},
			endpoints: map[string]map[string]string{
				"a": {"topology.istio.io/network": "local-network-1"},
				"b": {"topology.istio.io/network": "remote-network-2"},
			},
			failoverPriority: []string{
				"topology.istio.io/network",
			},
			expectedPriorities: map[string]uint32{
				"a": 0,
				"b": 1,
			},
		},
		{
			name:           "network failoverPriority infer client network",
			clientLabels:   map[string]string{}, // don't set anything here
			defaultNetwork: "local-network-1",   // instead it's inferred from istio-system
			endpoints: map[string]map[string]string{
				"a": {"topology.istio.io/network": "local-network-1"},
				"b": {"topology.istio.io/network": "remote-network-2"},
			},
			failoverPriority: []string{
				"topology.istio.io/network",
			},
			expectedPriorities: map[string]uint32{
				"a": 0,
				"b": 1,
			},
		},
		{
			name:           "network failoverPriority infer endpoint network",
			clientLabels:   map[string]string{"topology.istio.io/network": "local-network-1"},
			defaultNetwork: "local-network-1", // infer for missing label below
			endpoints: map[string]map[string]string{
				"a": {}, // inferred from defaultNetwork above
				"b": {"topology.istio.io/network": "remote-network-2"},
			},
			failoverPriority: []string{
				"topology.istio.io/network",
			},
			expectedPriorities: map[string]uint32{
				"a": 0,
				"b": 1,
			},
		},
		{
			name:         "network failoverPriority assume missing network is equal",
			clientLabels: map[string]string{"topology.istio.io/network": "local-network-1"},
			// there is no defaultNetwork to be inferred by "a"
			defaultNetwork: "",
			endpoints: map[string]map[string]string{
				"a": {}, // since a doesn't have any value for network, assume it's reachable
				"b": {"topology.istio.io/network": "remote-network-2"},
			},
			failoverPriority: []string{
				"topology.istio.io/network",
			},
			expectedPriorities: map[string]uint32{
				"a": 0,
				"b": 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			us := krtcollections.UpstreamWrapper{
				Inner: &gloov1.Upstream{
					Metadata: &core.Metadata{
						Name:      "name",
						Namespace: "ns",
					},
				},
			}
			efu := krtcollections.NewEndpointsForUpstream(us, nil)

			for endpointName, epLabels := range tc.endpoints {
				efu.Add(krtcollections.PodLocality{Region: epLabels[corev1.LabelTopologyRegion]}, krtcollections.EndpointWithMd{
					LbEndpoint: &endpointv3.LbEndpoint{
						HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
							Endpoint: &endpointv3.Endpoint{
								Address: &corev3.Address{
									Address: &corev3.Address_Pipe{Pipe: &corev3.Pipe{Path: string(endpointName)}},
								},
							},
						},
					},
					EndpointMd: krtcollections.EndpointMetadata{
						Labels: epLabels,
					},
				})
			}

			ucc := krtcollections.UniqlyConnectedClient{
				Namespace: "ns",
				Locality:  krtcollections.PodLocality{Region: tc.clientLabels[corev1.LabelTopologyRegion]},
				Labels:    tc.clientLabels,
			}

			destRule := &networkingclient.DestinationRule{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "networking.istio.io/v1alpha3",
					Kind:       "DestinationRule",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "do-failover",
				},
				Spec: networkingv1alpha3.DestinationRule{
					Host: "reviews.gwtest.svc.cluster.local",
					TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
						OutlierDetection: &networkingv1alpha3.OutlierDetection{
							Consecutive_5XxErrors: &wrappers.UInt32Value{Value: 7},
							Interval:              &duration.Duration{Seconds: 300}, // 5 minutes
							BaseEjectionTime:      &duration.Duration{Seconds: 900}, // 15 minutes
						},
						LoadBalancer: &networkingv1alpha3.LoadBalancerSettings{
							LocalityLbSetting: &networkingv1alpha3.LocalityLoadBalancerSetting{
								FailoverPriority: tc.failoverPriority,
							},
						},
					},
				},
			}

			uccWithEndpoints := PrioritizeEndpoints(nil, &DestinationRuleWrapper{destRule}, *efu, ucc, tc.defaultNetwork)
			cla := uccWithEndpoints.Endpoints.ResourceProto().(*endpointv3.ClusterLoadAssignment)
			// we expect each endpoint name to have its own priority
			g.Expect(cla.Endpoints).To(gomega.HaveLen(len(tc.endpoints)))

			epToPriority := make(map[string]uint32)
			for _, endpoint := range cla.Endpoints {
				epName := endpoint.LbEndpoints[0].GetEndpoint().Address.GetPipe().Path
				epToPriority[epName] = endpoint.Priority
			}

			// Verify priorities match expected values
			for endpointName, expectedPriority := range tc.expectedPriorities {
				g.Expect(epToPriority[endpointName]).To(gomega.Equal(expectedPriority), "wrong priority for "+endpointName)
			}
		})
	}
}

// similar to TestTranslatesDestrulesFailoverPriority but implicit
func TestTranslatesDestrulesFailover(t *testing.T) {
	g := gomega.NewWithT(t)
	us := krtcollections.UpstreamWrapper{
		Inner: &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "name",
				Namespace: "ns",
			},
		},
	}
	efu := krtcollections.NewEndpointsForUpstream(us, nil)
	efu.Add(krtcollections.PodLocality{Region: "R1"}, krtcollections.EndpointWithMd{
		LbEndpoint: &endpointv3.LbEndpoint{
			HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
				Endpoint: &endpointv3.Endpoint{
					Address: &corev3.Address{
						Address: &corev3.Address_Pipe{Pipe: &corev3.Pipe{Path: "a"}},
					},
				},
			},
		},
		EndpointMd: krtcollections.EndpointMetadata{
			Labels: map[string]string{corev1.LabelTopologyRegion: "R1"},
		},
	})
	efu.Add(krtcollections.PodLocality{Region: "R2"}, krtcollections.EndpointWithMd{
		LbEndpoint: &endpointv3.LbEndpoint{
			HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
				Endpoint: &endpointv3.Endpoint{
					Address: &corev3.Address{
						Address: &corev3.Address_Pipe{Pipe: &corev3.Pipe{Path: "b"}},
					},
				},
			},
		},
		EndpointMd: krtcollections.EndpointMetadata{
			Labels: map[string]string{corev1.LabelTopologyRegion: "R2"},
		},
	})
	ucc := krtcollections.UniqlyConnectedClient{
		Namespace: "ns",
		Locality:  krtcollections.PodLocality{Region: "R1"},
		Labels:    map[string]string{corev1.LabelTopologyRegion: "R1"},
	}

	destRule := &networkingclient.DestinationRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.istio.io/v1alpha3",
			Kind:       "DestinationRule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "do-failover",
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: "reviews.gwtest.svc.cluster.local",
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				OutlierDetection: &networkingv1alpha3.OutlierDetection{
					Consecutive_5XxErrors: &wrappers.UInt32Value{Value: 7},
					Interval:              &duration.Duration{Seconds: 300}, // 5 minutes
					BaseEjectionTime:      &duration.Duration{Seconds: 900}, // 15 minutes
				},
				LoadBalancer: &networkingv1alpha3.LoadBalancerSettings{
					LocalityLbSetting: &networkingv1alpha3.LocalityLoadBalancerSetting{},
				},
			},
		},
	}

	uccWithEndpoints := PrioritizeEndpoints(nil, &DestinationRuleWrapper{destRule}, *efu, ucc, "default-network")
	cla := uccWithEndpoints.Endpoints.ResourceProto().(*endpointv3.ClusterLoadAssignment)
	g.Expect(cla.Endpoints).To(gomega.HaveLen(2))

	remoteLocality := cla.Endpoints[0]
	localLocality := cla.Endpoints[1]
	if remoteLocality.Locality.Region == "R1" {
		remoteLocality = cla.Endpoints[1]
		localLocality = cla.Endpoints[0]
	}
	g.Expect(localLocality.Locality.Region).To(gomega.Equal("R1"))
	g.Expect(remoteLocality.Locality.Region).To(gomega.Equal("R2"))

	g.Expect(localLocality.Priority).To(gomega.Equal(uint32(0)))
	g.Expect(remoteLocality.Priority).To(gomega.Equal(uint32(1)))
}
