package internal_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/testutils"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	v1sets "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/sets"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover/internal"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("DependentsCalculator", func() {
	var (
		ctx  context.Context
		ctrl *gomock.Controller

		failoverSchemeClient *mock_fed_v1.MockFailoverSchemeClient
		glooInstanceClient   *mock_fed_v1.MockGlooInstanceClient
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		failoverSchemeClient = mock_fed_v1.NewMockFailoverSchemeClient(ctrl)
		glooInstanceClient = mock_fed_v1.NewMockGlooInstanceClient(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("FailoverDependencyCalculator", func() {

		Context("ForUpstream", func() {
			It("will return all failovers who reference the upstream", func() {
				depCalc := internal.NewFailoverDependencyCalculator(failoverSchemeClient, glooInstanceClient)

				usRef := &skv2v1.ClusterObjectRef{
					Name:        "upstream",
					Namespace:   "gloo-system",
					ClusterName: "one",
				}
				list := &fedv1.FailoverSchemeList{
					Items: []fedv1.FailoverScheme{
						{
							Spec: fed_types.FailoverSchemeSpec{
								FailoverGroups: []*fed_types.FailoverSchemeSpec_FailoverEndpoints{
									{
										PriorityGroup: []*fed_types.FailoverSchemeSpec_FailoverEndpoints_LocalityLbTargets{
											{
												Cluster: usRef.GetClusterName(),
												Upstreams: []*skv2v1.ObjectRef{
													{
														Name:      usRef.GetName(),
														Namespace: usRef.GetNamespace(),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Spec: fed_types.FailoverSchemeSpec{
								Primary: &skv2v1.ClusterObjectRef{
									Name:        usRef.GetName(),
									Namespace:   usRef.GetNamespace(),
									ClusterName: usRef.GetClusterName(),
								},
							},
						},
						{
							Spec: fed_types.FailoverSchemeSpec{
								Primary: &skv2v1.ClusterObjectRef{
									Name:        "not",
									Namespace:   "valid",
									ClusterName: "cluster",
								},
							},
						},
						{
							Spec: fed_types.FailoverSchemeSpec{
								FailoverGroups: []*fed_types.FailoverSchemeSpec_FailoverEndpoints{
									{
										PriorityGroup: []*fed_types.FailoverSchemeSpec_FailoverEndpoints_LocalityLbTargets{
											{
												Cluster: "also",
												Upstreams: []*skv2v1.ObjectRef{
													{
														Name:      "not",
														Namespace: "valid",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}

				failoverSchemeClient.EXPECT().
					ListFailoverScheme(ctx).
					Return(list, nil)

				result, err := depCalc.ForUpstream(ctx, usRef)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal([]*fedv1.FailoverScheme{&list.Items[0], &list.Items[1]}))
			})
		})

		Context("ForGlooInstance", func() {

			It("will return whole list and error", func() {

				var (
					list = &fedv1.FailoverSchemeList{
						Items: []fedv1.FailoverScheme{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "one",
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "two",
								},
							},
						},
					}
				)

				When("gloo instance cannot be found", func() {
					depCalc := internal.NewFailoverDependencyCalculator(failoverSchemeClient, glooInstanceClient)

					failoverSchemeClient.EXPECT().
						ListFailoverScheme(ctx).
						Return(list, nil)

					instanceRef := &skv2v1.ObjectRef{
						Name:      "hello",
						Namespace: "world",
					}

					glooInstanceClient.EXPECT().
						GetGlooInstance(ctx, client.ObjectKey{
							Namespace: instanceRef.GetNamespace(),
							Name:      instanceRef.GetName(),
						}).
						Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))

					result, err := depCalc.ForGlooInstance(ctx, instanceRef)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal(v1sets.NewFailoverSchemeSetFromList(list).List()))
				})

				When("gloo instance fails for any other reason", func() {

					depCalc := internal.NewFailoverDependencyCalculator(failoverSchemeClient, glooInstanceClient)

					failoverSchemeClient.EXPECT().
						ListFailoverScheme(ctx).
						Return(list, nil)

					instanceRef := &skv2v1.ObjectRef{
						Name:      "hello",
						Namespace: "world",
					}

					testErr := eris.New("error")
					glooInstanceClient.EXPECT().
						GetGlooInstance(ctx, client.ObjectKey{
							Namespace: instanceRef.GetNamespace(),
							Name:      instanceRef.GetName(),
						}).
						Return(nil, testErr)

					result, err := depCalc.ForGlooInstance(ctx, instanceRef)
					Expect(err).To(testutils.HaveInErrorChain(testErr))
					Expect(result).To(Equal(v1sets.NewFailoverSchemeSetFromList(list).List()))
				})

			})

			It("will return all relevant schemes for the instance owned upstreams", func() {

				depCalc := internal.NewFailoverDependencyCalculator(failoverSchemeClient, glooInstanceClient)

				list := &fedv1.FailoverSchemeList{
					Items: []fedv1.FailoverScheme{
						{
							Spec: fed_types.FailoverSchemeSpec{
								Primary: &skv2v1.ClusterObjectRef{
									Name:        "upstream-two",
									Namespace:   "gloo-system",
									ClusterName: "one",
								},
							},
						},
						{
							Spec: fed_types.FailoverSchemeSpec{
								Primary: &skv2v1.ClusterObjectRef{
									Name:        "upstream-two",
									Namespace:   "gloo-system",
									ClusterName: "two",
								},
							},
						},
					},
				}
				failoverSchemeClient.EXPECT().
					ListFailoverScheme(ctx).
					Return(list, nil)

				instanceRef := &skv2v1.ObjectRef{
					Name:      "hello",
					Namespace: "world",
				}

				instance := &fedv1.GlooInstance{
					Spec: fed_types.GlooInstanceSpec{
						Cluster: "one",
					},
				}

				glooInstanceClient.EXPECT().
					GetGlooInstance(ctx, client.ObjectKey{
						Namespace: instanceRef.GetNamespace(),
						Name:      instanceRef.GetName(),
					}).
					Return(instance, nil)

				result, err := depCalc.ForGlooInstance(ctx, instanceRef)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal([]*fedv1.FailoverScheme{&list.Items[0]}))
			})

		})

	})
})
