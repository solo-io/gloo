package failover_test

import (
	"context"

	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover"
	mock_failover "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/reconcile"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	test_matchers "github.com/solo-io/solo-projects/projects/gloo-fed/test/matchers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("DependentsReconciler", func() {
	var (
		ctx  context.Context
		ctrl *gomock.Controller

		failoverSchemeClient *mock_fed_v1.MockFailoverSchemeClient
		statusManager        *failover.StatusManager
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		failoverSchemeClient = mock_fed_v1.NewMockFailoverSchemeClient(ctrl)
		statusManager = failover.NewStatusManager(failoverSchemeClient, defaults.GlooFed)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("FailoverDependentReconciler", func() {
		var (
			depCalc *mock_failover.MockFailoverDependencyCalculator
			result  []*fedv1.FailoverScheme
		)

		BeforeEach(func() {
			depCalc = mock_failover.NewMockFailoverDependencyCalculator(ctrl)

			result = []*fedv1.FailoverScheme{
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
			}

			expected1 := result[0].DeepCopy()
			expectedStatus := createStatus(&failover.StatusImpl{
				State:              fed_types.FailoverSchemeStatus_PENDING,
				Message:            failover.DependentUpdateMessage,
				ObservedGeneration: expected1.GetGeneration(),
			})
			expected1.Status = expectedStatus
			expected2 := result[1].DeepCopy()
			expected2.Status = expectedStatus

			failoverSchemeClient.EXPECT().
				UpdateFailoverSchemeStatus(ctx, test_matchers.MatchesFailover(expected1)).
				Return(nil)
			failoverSchemeClient.EXPECT().
				UpdateFailoverSchemeStatus(ctx, test_matchers.MatchesFailover(expected2)).
				Return(errors.NewConflict(schema.GroupResource{}, "", nil))
		})

		It("can ReconcileGlooInstance", func() {
			reconciler := failover.NewFailoverDependentReconciler(ctx, depCalc, failoverSchemeClient, statusManager)
			obj := &fedv1.GlooInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "one",
					Namespace: "two",
				},
			}

			depCalc.EXPECT().
				ForGlooInstance(ctx, &skv2v1.ObjectRef{
					Name:      obj.GetName(),
					Namespace: obj.GetNamespace(),
				}).
				Return(result, nil)

			_, err := reconciler.ReconcileGlooInstance(obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can ReconcileGlooInstanceDeletion", func() {
			reconciler := failover.NewFailoverDependentReconciler(ctx, depCalc, failoverSchemeClient, statusManager)
			obj := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "one",
					Name:      "two",
				},
			}

			depCalc.EXPECT().
				ForGlooInstance(ctx, &skv2v1.ObjectRef{
					Name:      obj.Name,
					Namespace: obj.Namespace,
				}).
				Return(result, nil)

			err := reconciler.ReconcileGlooInstanceDeletion(obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can ReconcileUpstream", func() {
			reconciler := failover.NewFailoverDependentReconciler(ctx, depCalc, failoverSchemeClient, statusManager)
			obj := &gloov1.Upstream{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "one",
					Namespace: "two",
				},
			}
			clusterName := "cluster"
			depCalc.EXPECT().
				ForUpstream(ctx, &skv2v1.ClusterObjectRef{
					Name:        obj.GetName(),
					Namespace:   obj.GetNamespace(),
					ClusterName: clusterName,
				}).
				Return(result, nil)

			_, err := reconciler.ReconcileUpstream(clusterName, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can ReconcileUpstreamDeletion", func() {
			reconciler := failover.NewFailoverDependentReconciler(ctx, depCalc, failoverSchemeClient, statusManager)
			obj := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "one",
					Name:      "two",
				},
			}

			clusterName := "cluster"
			depCalc.EXPECT().
				ForUpstream(ctx, &skv2v1.ClusterObjectRef{
					Name:        obj.Name,
					Namespace:   obj.Namespace,
					ClusterName: clusterName,
				}).
				Return(result, nil)

			err := reconciler.ReconcileUpstreamDeletion(clusterName, obj)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
