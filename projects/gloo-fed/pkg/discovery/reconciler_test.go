package discovery_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/go-utils/contextutils"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/input"
	mock_input "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/input/mocks"
	mock_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery"
	mock_translator "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/translator/mocks"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Deployment Reconciler Test", func() {
	var (
		ctrl               *gomock.Controller
		glooInstanceClient *mock_v1.MockGlooInstanceClient
		translator         *mock_translator.MockGlooInstanceTranslator
		snapshotBuilder    *mock_input.MockBuilder

		ctx            = contextutils.WithExistingLogger(context.Background(), zap.NewNop().Sugar())
		writeNamespace = "foo-ns"
		cluster        = "remote-cluster"
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		glooInstanceClient = mock_v1.NewMockGlooInstanceClient(ctrl)
		translator = mock_translator.NewMockGlooInstanceTranslator(ctrl)
		snapshotBuilder = mock_input.NewMockBuilder(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ReconcileAll", func() {

		It("will properly upsert all instances", func() {
			reconciler := discovery.NewGlooResourceReconciler(ctx, glooInstanceClient, snapshotBuilder, translator)

			instance1 := &fedv1.GlooInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance-1",
					Namespace: writeNamespace,
				},
			}
			instance2 := &fedv1.GlooInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance-2",
					Namespace: writeNamespace,
				},
			}

			snap := input.NewSnapshot(discovery.SnapshotName, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			snapshotBuilder.EXPECT().
				BuildSnapshot(ctx, discovery.SnapshotName, input.BuildOptions{}).
				Return(snap, nil)

			translator.EXPECT().FromSnapshot(ctx, snap).
				Return([]*fedv1.GlooInstance{instance1, instance2})

			glooInstanceClient.EXPECT().
				ListGlooInstance(ctx).
				Return(&fedv1.GlooInstanceList{}, nil)

			glooInstanceClient.EXPECT().
				UpsertGlooInstance(ctx, instance1).
				Return(nil)

			glooInstanceClient.EXPECT().
				UpsertGlooInstance(ctx, instance2).
				Return(nil)

			_, err := reconciler.ReconcileAll(&skv2v1.ClusterObjectRef{
				Name:        "",
				Namespace:   "",
				ClusterName: cluster,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("will delete old instances", func() {
			reconciler := discovery.NewGlooResourceReconciler(ctx, glooInstanceClient, snapshotBuilder, translator)

			instance1 := fedv1.GlooInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance-1",
					Namespace: writeNamespace,
				},
			}

			snap := input.NewSnapshot(discovery.SnapshotName, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			snapshotBuilder.EXPECT().
				BuildSnapshot(ctx, discovery.SnapshotName, input.BuildOptions{}).
				Return(snap, nil)

			translator.EXPECT().FromSnapshot(ctx, snap).
				Return([]*fedv1.GlooInstance{})

			glooInstanceClient.EXPECT().
				ListGlooInstance(ctx).
				Return(&fedv1.GlooInstanceList{Items: []fedv1.GlooInstance{instance1}}, nil)

			glooInstanceClient.EXPECT().
				DeleteGlooInstance(ctx, client.ObjectKey{
					Namespace: instance1.GetNamespace(),
					Name:      instance1.GetName(),
				}).
				Return(nil)

			_, err := reconciler.ReconcileAll(&skv2v1.ClusterObjectRef{
				Name:        "",
				Namespace:   "",
				ClusterName: cluster,
			})
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
