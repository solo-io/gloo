package syncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana/mocks"
	"go.uber.org/zap"
)

var (
	mockCtrl           *gomock.Controller
	dashboardClient    *mocks.MockDashboardClient
	snapshotClient     *mocks.MockSnapshotClient
	templateGenerator  *mocks.MockTemplateGenerator
	testErr            = errors.New("test err")
	dashboardsSnapshot *v1.DashboardsSnapshot
	testGrafanaState   *grafanaState
	dashboardSyncer    *GrafanaDashboardsSyncer
	logger             *zap.SugaredLogger
)

var _ = Describe("Grafana Syncer", func() {
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		dashboardClient = mocks.NewMockDashboardClient(mockCtrl)
		snapshotClient = mocks.NewMockSnapshotClient(mockCtrl)
		templateGenerator = mocks.NewMockTemplateGenerator(mockCtrl)
		dashboardsSnapshot = &v1.DashboardsSnapshot{
			Upstreams: []*gloov1.Upstream{
				{
					Metadata: core.Metadata{Name: "us1", Namespace: "ns"},
				},
				{
					Metadata: core.Metadata{Name: "us2", Namespace: "ns"},
				},
			},
		}
		testGrafanaState = &grafanaState{
			boards:    nil,
			snapshots: nil,
		}
		dashboardSyncer = NewGrafanaDashboardSyncer(dashboardClient, snapshotClient)
		logger = contextutils.LoggerFrom(context.TODO())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("shouldRegenDashboard", func() {
		It("answers yes when the dashboard is un-edited", func() {
			upstreamUid := "test-uid"
			var dashboardId float64 = 123
			rawDashboard := []byte(fmt.Sprintf("{\"id\": %f}", dashboardId))
			versions := []*grafana.Version{
				{
					Message: grafana.DefaultCommitMessage,
				},
			}

			dashboardClient.EXPECT().
				GetRawDashboard(upstreamUid).
				Return(rawDashboard, grafana.BoardProperties{}, nil)
			dashboardClient.EXPECT().
				GetDashboardVersions(dashboardId).
				Return(versions, nil)

			shouldRegen, err := dashboardSyncer.shouldRegenDashboard(logger, upstreamUid)

			Expect(err).NotTo(HaveOccurred())
			Expect(shouldRegen).To(BeTrue())
		})

		It("answers no when the dashboard is edited", func() {
			upstreamUid := "test-uid"
			var dashboardId float64 = 123
			rawDashboard := []byte(fmt.Sprintf("{\"id\": %f}", dashboardId))
			versions := []*grafana.Version{
				{
					Message: "my message written by a user",
				},
				{
					Message: grafana.DefaultCommitMessage,
				},
			}

			dashboardClient.EXPECT().
				GetRawDashboard(upstreamUid).
				Return(rawDashboard, grafana.BoardProperties{}, nil)
			dashboardClient.EXPECT().
				GetDashboardVersions(dashboardId).
				Return(versions, nil)

			shouldRegen, err := dashboardSyncer.shouldRegenDashboard(logger, upstreamUid)

			Expect(err).NotTo(HaveOccurred())
			Expect(shouldRegen).To(BeFalse())
		})

		It("answers yes when the dashboard has been deleted", func() {
			upstreamUid := "test-uid"

			dashboardClient.EXPECT().
				GetRawDashboard(upstreamUid).
				Return(nil, grafana.BoardProperties{}, grafana.DashboardNotFound(upstreamUid))

			shouldRegen, err := dashboardSyncer.shouldRegenDashboard(logger, upstreamUid)

			Expect(err).NotTo(HaveOccurred())
			Expect(shouldRegen).To(BeTrue())
		})

		It("answers no when the dashboard client errors", func() {
			upstreamUid := "test-uid"

			dashboardClient.EXPECT().
				GetRawDashboard(upstreamUid).
				Return(nil, grafana.BoardProperties{}, testErr)

			shouldRegen, err := dashboardSyncer.shouldRegenDashboard(logger, upstreamUid)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(testErr))
			Expect(shouldRegen).To(BeFalse())
		})
	})
})
