package failover_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloo_api_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	gloo_types "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	mock_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/mocks"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-kit/test/matchers"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/fields"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover"
	test_matchers "github.com/solo-io/solo-projects/projects/gloo-fed/test/matchers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Processor", func() {

	var (
		ctx  context.Context
		ctrl *gomock.Controller

		glooInstanceClient   *mock_fed_v1.MockGlooInstanceClient
		failoverSchemeClient *mock_fed_v1.MockFailoverSchemeClient
		glooMcClientset      *mock_gloo_v1.MockMulticlusterClientset
		glooClientset        *mock_gloo_v1.MockClientset
		upstreamClient       *mock_gloo_v1.MockUpstreamClient
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		glooInstanceClient = mock_fed_v1.NewMockGlooInstanceClient(ctrl)
		failoverSchemeClient = mock_fed_v1.NewMockFailoverSchemeClient(ctrl)
		glooMcClientset = mock_gloo_v1.NewMockMulticlusterClientset(ctrl)
		glooClientset = mock_gloo_v1.NewMockClientset(ctrl)
		upstreamClient = mock_gloo_v1.NewMockUpstreamClient(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ProcessFailoverUpdate", func() {

		Context("Failure", func() {

			It("will error if no primary target is given", func() {

				processor := failover.NewFailoverProcessor(glooMcClientset, glooInstanceClient, failoverSchemeClient)
				obj := &fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "one",
						Namespace: "two",
					},
				}
				_, status := processor.ProcessFailoverUpdate(ctx, obj)
				Expect(status).To(test_matchers.MatchFailoverStatus(&fed_types.FailoverSchemeStatus{
					State:   fed_types.FailoverSchemeStatus_INVALID,
					Message: failover.EmptyPrimaryTargetError.Error(),
				}))
			})

			It("will error if no primary is already in use", func() {

				processor := failover.NewFailoverProcessor(glooMcClientset, glooInstanceClient, failoverSchemeClient)
				obj := &fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "one",
						Namespace: "two",
					},
					Spec: fed_types.FailoverSchemeSpec{
						Primary: &skv2v1.ClusterObjectRef{
							Name:        "upstream",
							Namespace:   "gloo-system",
							ClusterName: "cluster",
						},
					},
				}

				glooMcClientset.EXPECT().
					Cluster(obj.Spec.GetPrimary().GetClusterName()).
					Return(glooClientset, nil)

				glooClientset.EXPECT().
					Upstreams().
					Return(upstreamClient)

				upstreamClient.EXPECT().
					GetUpstream(ctx, client.ObjectKey{
						Namespace: obj.Spec.GetPrimary().GetNamespace(),
						Name:      obj.Spec.GetPrimary().GetName(),
					}).
					Return(nil, nil)

				errored := obj.DeepCopy()
				errored.Name = "different"
				failoverSchemeClient.EXPECT().
					ListFailoverScheme(ctx).
					Return(&fedv1.FailoverSchemeList{Items: []fedv1.FailoverScheme{*errored}}, nil)

				_, status := processor.ProcessFailoverUpdate(ctx, obj)
				Expect(status).To(test_matchers.MatchFailoverStatus(&fed_types.FailoverSchemeStatus{
					State:   fed_types.FailoverSchemeStatus_INVALID,
					Message: failover.PrimaryTargetAlreadyInUseError(obj.Spec.GetPrimary(), errored).Error(),
				}))
			})

			It("will error if no failover groups are specified", func() {

				processor := failover.NewFailoverProcessor(glooMcClientset, glooInstanceClient, failoverSchemeClient)
				obj := &fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "one",
						Namespace: "two",
					},
					Spec: fed_types.FailoverSchemeSpec{
						Primary: &skv2v1.ClusterObjectRef{
							Name:        "upstream",
							Namespace:   "gloo-system",
							ClusterName: "cluster",
						},
					},
				}

				glooMcClientset.EXPECT().
					Cluster(obj.Spec.GetPrimary().GetClusterName()).
					Return(glooClientset, nil)

				glooClientset.EXPECT().
					Upstreams().
					Return(upstreamClient)

				upstreamClient.EXPECT().
					GetUpstream(ctx, client.ObjectKey{
						Namespace: obj.Spec.GetPrimary().GetNamespace(),
						Name:      obj.Spec.GetPrimary().GetName(),
					}).
					Return(nil, nil)

				failoverSchemeClient.EXPECT().
					ListFailoverScheme(ctx).
					Return(&fedv1.FailoverSchemeList{}, nil)

				_, status := processor.ProcessFailoverUpdate(ctx, obj)
				Expect(status).To(test_matchers.MatchFailoverStatus(&fed_types.FailoverSchemeStatus{
					State:   fed_types.FailoverSchemeStatus_INVALID,
					Message: failover.EmptyFailoverTargetsError.Error(),
				}))
			})

		})

		Context("Success", func() {

			It("can return a valid failover cfg", func() {

				processor := failover.NewFailoverProcessor(glooMcClientset, glooInstanceClient, failoverSchemeClient)

				instanceCluster1 := &fedv1.GlooInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "instance-1",
						Namespace: "gloo-system",
					},
					Spec: fed_types.GlooInstanceSpec{
						Cluster: "one",
						ControlPlane: &fed_types.GlooInstanceSpec_ControlPlane{
							Namespace: defaults.GlooSystem,
						},
						Proxies: []*fed_types.GlooInstanceSpec_Proxy{
							{
								Name:      "proxy-1",
								Namespace: defaults.GlooSystem,
								Zones:     []string{"zone-cluster-1"},
								IngressEndpoints: []*fed_types.GlooInstanceSpec_Proxy_IngressEndpoint{
									{
										Address: "address.cluster.1",
										Ports: []*fed_types.GlooInstanceSpec_Proxy_IngressEndpoint_Port{
											{
												Port: failover.PortNumber,
												Name: failover.PortName,
											},
										},
									},
								},
							},
						},
						Region: "region-cluster-1",
						Admin: &fed_types.GlooInstanceSpec_Admin{
							WriteNamespace: defaults.GlooSystem,
							ProxyId: &skv2v1.ObjectRef{
								Name:      "proxy-1",
								Namespace: defaults.GlooSystem,
							},
						},
					},
				}

				primaryUpstream := &gloov1.Upstream{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "primary",
						Namespace: defaults.GlooSystem,
					},
				}

				failoverUpstream := &gloov1.Upstream{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "failover-1",
						Namespace: defaults.GlooSystem,
					},
				}

				obj := &fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "one",
						Namespace: "two",
					},
					Spec: fed_types.FailoverSchemeSpec{
						Primary: &skv2v1.ClusterObjectRef{
							Name:        "upstream",
							Namespace:   defaults.GlooSystem,
							ClusterName: "one",
						},
						FailoverGroups: []*fed_types.FailoverSchemeSpec_FailoverEndpoints{
							{
								PriorityGroup: []*fed_types.FailoverSchemeSpec_FailoverEndpoints_LocalityLbTargets{
									{
										Cluster: instanceCluster1.Spec.GetCluster(),
										Upstreams: []*skv2v1.ObjectRef{
											{
												Name:      failoverUpstream.GetName(),
												Namespace: failoverUpstream.GetNamespace(),
											},
										},
										LocalityWeight: &wrappers.UInt32Value{Value: 1},
									},
								},
							},
						},
					},
				}

				glooMcClientset.EXPECT().
					Cluster(obj.Spec.GetPrimary().GetClusterName()).
					Return(glooClientset, nil).
					Times(2)

				glooClientset.EXPECT().
					Upstreams().
					Return(upstreamClient).
					Times(2)

				upstreamClient.EXPECT().
					GetUpstream(ctx, client.ObjectKey{
						Namespace: obj.Spec.GetPrimary().GetNamespace(),
						Name:      obj.Spec.GetPrimary().GetName(),
					}).
					Return(primaryUpstream, nil)

				failoverSchemeClient.EXPECT().
					ListFailoverScheme(ctx).
					Return(&fedv1.FailoverSchemeList{}, nil)

				glooInstanceClient.EXPECT().
					ListGlooInstance(ctx, fields.BuildClusterFieldMatcher(instanceCluster1.Spec.GetCluster())).
					Return(&fedv1.GlooInstanceList{Items: []fedv1.GlooInstance{*instanceCluster1}}, nil)

				upstreamClient.EXPECT().
					GetUpstream(ctx, client.ObjectKey{
						Namespace: failoverUpstream.GetNamespace(),
						Name:      failoverUpstream.GetName(),
					}).
					Return(failoverUpstream, nil)

				expected := primaryUpstream.DeepCopy()
				expected.Spec.Failover = &gloo_api_v1.Failover{
					PrioritizedLocalities: []*gloo_api_v1.Failover_PrioritizedLocality{
						{
							LocalityEndpoints: []*gloo_api_v1.LocalityLbEndpoints{
								{
									Locality: &gloo_api_v1.Locality{
										Region: instanceCluster1.Spec.GetRegion(),
										Zone:   instanceCluster1.Spec.Proxies[0].GetZones()[0],
									},
									LbEndpoints: []*gloo_api_v1.LbEndpoint{
										{
											Address: instanceCluster1.Spec.GetProxies()[0].GetIngressEndpoints()[0].GetAddress(),
											Port:    failover.PortNumber,
											UpstreamSslConfig: &gloo_api_v1.UpstreamSslConfig{
												SslSecrets: &gloo_api_v1.UpstreamSslConfig_SecretRef{
													SecretRef: &core.ResourceRef{
														Name:      failover.UpstreamSecretName,
														Namespace: defaults.GlooSystem,
													},
												},
												Sni: failover.UpstreamToClusterName(&skv2v1.ObjectRef{
													Name:      obj.Spec.GetFailoverGroups()[0].GetPriorityGroup()[0].GetUpstreams()[0].GetName(),
													Namespace: obj.Spec.GetFailoverGroups()[0].GetPriorityGroup()[0].GetUpstreams()[0].GetNamespace(),
												}),
											},
										},
									},
									LoadBalancingWeight: obj.Spec.GetFailoverGroups()[0].GetPriorityGroup()[0].GetLocalityWeight(),
								},
							},
						},
					},
				}

				returned, status := processor.ProcessFailoverUpdate(ctx, obj)
				Expect(status).To(BeNil())
				Expect(returned).To(MatchPublicFields(expected))
			})

		})
	})

	Context("ProcessFailoverDelete", func() {

		It("will get the primary upstream, and set the failover policy to nil", func() {
			processor := failover.NewFailoverProcessor(glooMcClientset, glooInstanceClient, failoverSchemeClient)
			obj := &fedv1.FailoverScheme{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "one",
					Namespace: "two",
				},
				Spec: fed_types.FailoverSchemeSpec{
					Primary: &skv2v1.ClusterObjectRef{
						Name:        "upstream",
						Namespace:   "gloo-system",
						ClusterName: "cluster",
					},
				},
			}

			glooMcClientset.EXPECT().
				Cluster(obj.Spec.GetPrimary().GetClusterName()).
				Return(glooClientset, nil)

			glooClientset.EXPECT().
				Upstreams().
				Return(upstreamClient)

			primaryUs := &gloov1.Upstream{
				Spec: gloo_types.UpstreamSpec{
					Failover: &gloo_api_v1.Failover{},
				},
			}
			upstreamClient.EXPECT().
				GetUpstream(ctx, client.ObjectKey{
					Namespace: obj.Spec.GetPrimary().GetNamespace(),
					Name:      obj.Spec.GetPrimary().GetName(),
				}).
				Return(primaryUs, nil)

			expected := primaryUs.DeepCopy()
			expected.Spec.Failover = nil
			result, err := processor.ProcessFailoverDelete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(MatchPublicFields(expected))
		})

	})

})
