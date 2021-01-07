package failover_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloo_types "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	mock_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/mocks"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"
	. "github.com/solo-io/solo-kit/test/matchers"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover"
	mock_failover "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover/mocks"
	test_matchers "github.com/solo-io/solo-projects/projects/gloo-fed/test/matchers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Reconciler", func() {
	var (
		ctx                  context.Context
		ctrl                 *gomock.Controller
		processor            *mock_failover.MockFailoverProcessor
		failoverSchemeClient *mock_fed_v1.MockFailoverSchemeClient
		glooMcClientset      *mock_gloo_v1.MockMulticlusterClientset
		glooClientset        *mock_gloo_v1.MockClientset
		upstreamClient       *mock_gloo_v1.MockUpstreamClient
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		processor = mock_failover.NewMockFailoverProcessor(ctrl)
		failoverSchemeClient = mock_fed_v1.NewMockFailoverSchemeClient(ctrl)
		glooMcClientset = mock_gloo_v1.NewMockMulticlusterClientset(ctrl)
		glooClientset = mock_gloo_v1.NewMockClientset(ctrl)
		upstreamClient = mock_gloo_v1.NewMockUpstreamClient(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ReconcileFailoverScheme", func() {
		It("will skip reconcile", func() {

			When("generation is equal and status is invalid", func() {
				reconciler := failover.NewFailoverSchemeReconciler(ctx, processor, failoverSchemeClient, glooMcClientset)
				_, err := reconciler.ReconcileFailoverScheme(&fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Generation: 10,
					},
					Status: fed_types.FailoverSchemeStatus{
						State:              fed_types.FailoverSchemeStatus_INVALID,
						ObservedGeneration: 10,
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			When("generation is equal and status is accepted", func() {
				reconciler := failover.NewFailoverSchemeReconciler(ctx, processor, failoverSchemeClient, glooMcClientset)
				_, err := reconciler.ReconcileFailoverScheme(&fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Generation: 10,
					},
					Status: fed_types.FailoverSchemeStatus{
						State:              fed_types.FailoverSchemeStatus_ACCEPTED,
						ObservedGeneration: 10,
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

		})

		It("will return the returned status from a processor", func() {

			When("the status is non-nil", func() {
				reconciler := failover.NewFailoverSchemeReconciler(ctx, processor, failoverSchemeClient, glooMcClientset)

				obj := &fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Generation: 11,
					},
					Status: fed_types.FailoverSchemeStatus{
						State:              fed_types.FailoverSchemeStatus_ACCEPTED,
						ObservedGeneration: 10,
					},
				}
				status := &fed_types.FailoverSchemeStatus{
					State:              fed_types.FailoverSchemeStatus_FAILED,
					Message:            "im.a.status",
					ObservedGeneration: 8,
				}
				processor.EXPECT().ProcessFailoverUpdate(ctx, obj).Return(nil, status)

				expected := obj.DeepCopy()
				expected.Status = *status
				failoverSchemeClient.EXPECT().UpdateFailoverSchemeStatus(ctx, GomockMatchPublicFields(expected)).Return(nil)
				_, err := reconciler.ReconcileFailoverScheme(obj)
				Expect(obj).To(MatchPublicFields(expected))
				Expect(err).NotTo(HaveOccurred())

			})

		})

		It("can successfully upsert upstream", func() {
			reconciler := failover.NewFailoverSchemeReconciler(ctx, processor, failoverSchemeClient, glooMcClientset)

			obj := &fedv1.FailoverScheme{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
				Spec: fed_types.FailoverSchemeSpec{
					Primary: &skv2v1.ClusterObjectRef{
						Name:        "hello",
						Namespace:   "world",
						ClusterName: "cluster",
					},
				},
				Status: fed_types.FailoverSchemeStatus{
					ObservedGeneration: 0,
				},
			}

			us := &gloov1.Upstream{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route",
					Namespace: "#1",
				},
				Spec: gloo_types.UpstreamSpec{},
			}

			processor.EXPECT().ProcessFailoverUpdate(ctx, obj).Return(us, nil)

			glooMcClientset.EXPECT().
				Cluster(obj.Spec.GetPrimary().GetClusterName()).
				Return(glooClientset, nil)

			glooClientset.EXPECT().
				Upstreams().
				Return(upstreamClient)

			upstreamClient.EXPECT().
				UpsertUpstream(ctx, us).
				Return(nil)

			expected := obj.DeepCopy()
			expected.Status = fed_types.FailoverSchemeStatus{
				State:              fed_types.FailoverSchemeStatus_ACCEPTED,
				ObservedGeneration: 1,
				ProcessingTime:     prototime.Now(),
			}

			failoverSchemeClient.EXPECT().
				UpdateFailoverSchemeStatus(ctx, test_matchers.MatchesFailover(expected)).
				Return(nil)
			_, err := reconciler.ReconcileFailoverScheme(obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("FinalizeFailoverScheme", func() {

		It("will ignore a not found error", func() {

			When("one is returned from the processor", func() {
				reconciler := failover.NewFailoverSchemeReconciler(ctx, processor, failoverSchemeClient, glooMcClientset)

				obj := &fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Generation: 1,
					},
					Spec: fed_types.FailoverSchemeSpec{
						Primary: &skv2v1.ClusterObjectRef{
							Name:        "hello",
							Namespace:   "world",
							ClusterName: "cluster",
						},
					},
					Status: fed_types.FailoverSchemeStatus{
						ObservedGeneration: 0,
					},
				}

				processor.EXPECT().
					ProcessFailoverDelete(ctx, obj).
					Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))

				err := reconciler.FinalizeFailoverScheme(obj)
				Expect(err).NotTo(HaveOccurred())
			})

			When("one is returned from the upstream upsert", func() {
				reconciler := failover.NewFailoverSchemeReconciler(ctx, processor, failoverSchemeClient, glooMcClientset)

				obj := &fedv1.FailoverScheme{
					ObjectMeta: metav1.ObjectMeta{
						Generation: 1,
					},
					Spec: fed_types.FailoverSchemeSpec{
						Primary: &skv2v1.ClusterObjectRef{
							Name:        "hello",
							Namespace:   "world",
							ClusterName: "cluster",
						},
					},
					Status: fed_types.FailoverSchemeStatus{
						ObservedGeneration: 0,
					},
				}

				us := &gloov1.Upstream{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "route",
						Namespace: "#1",
					},
					Spec: gloo_types.UpstreamSpec{},
				}

				processor.EXPECT().
					ProcessFailoverDelete(ctx, obj).
					Return(us, nil)

				glooMcClientset.EXPECT().
					Cluster(obj.Spec.GetPrimary().GetClusterName()).
					Return(glooClientset, nil)

				glooClientset.EXPECT().
					Upstreams().
					Return(upstreamClient)

				upstreamClient.EXPECT().
					UpsertUpstream(ctx, us).
					Return(errors.NewNotFound(schema.GroupResource{}, ""))

				err := reconciler.FinalizeFailoverScheme(obj)
				Expect(err).NotTo(HaveOccurred())
			})

		})

	})

})
