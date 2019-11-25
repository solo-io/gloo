package syncer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"

	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template"
	templatemocks "github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template/mocks"

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
	mockCtrl                *gomock.Controller
	dashboardClient         *mocks.MockDashboardClient
	snapshotClient          *mocks.MockSnapshotClient
	templateGenerator       *templatemocks.MockTemplateGenerator
	testErr                 = errors.New("test err")
	dashboardsSnapshot      *v1.DashboardsSnapshot
	testGrafanaState        *grafanaState
	dashboardSyncer         *GrafanaDashboardsSyncer
	logger                  *zap.SugaredLogger
	upstreamOne             *gloov1.Upstream
	upstreamTwo             *gloov1.Upstream
	upstreamList            gloov1.UpstreamList
	snapshotResponse        *grafana.SnapshotListResponse
	dashboardSearchResponse []grafana.FoundBoard
	dashboardJsonTemplate   string
)

var _ = Describe("Grafana Syncer", func() {
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		dashboardClient = mocks.NewMockDashboardClient(mockCtrl)
		snapshotClient = mocks.NewMockSnapshotClient(mockCtrl)
		templateGenerator = templatemocks.NewMockTemplateGenerator(mockCtrl)
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
		dashboardJsonTemplate = "{\"test-content\": \"content\"}"
		dashboardSyncer = NewGrafanaDashboardSyncer(dashboardClient, snapshotClient, dashboardJsonTemplate)
		logger = contextutils.LoggerFrom(context.TODO())
		upstreamOne = &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Aws{
				Aws: &aws.UpstreamSpec{Region: "test"},
			},
			Status:   core.Status{},
			Metadata: core.Metadata{Name: "test-upstream-one", Namespace: "ns"},
		}
		upstreamTwo = &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &kubernetes.UpstreamSpec{
					ServiceName:      "upstream-service",
					ServiceNamespace: "ns",
					ServicePort:      80,
				},
			},
			Status:   core.Status{},
			Metadata: core.Metadata{Name: "test-upstream-two", Namespace: "ns"},
		}
		upstreamList = []*gloov1.Upstream{upstreamOne, upstreamTwo}

		snapshotResponse = &grafana.SnapshotListResponse{}
		dashboardSearchResponse = []grafana.FoundBoard{}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Dashboard Syncer", func() {
		var addIdToDashboardJson = func(dashboard []byte, id int) []byte {
			jsonDashboard := map[string]interface{}{}
			err := json.Unmarshal(dashboard, &jsonDashboard)
			Expect(err).NotTo(HaveOccurred())
			jsonDashboard["id"] = id
			jsonWithId, err := json.Marshal(jsonDashboard)
			Expect(err).NotTo(HaveOccurred())

			return jsonWithId
		}

		It("builds dashboards from scratch", func() {
			for _, upstream := range upstreamList {
				templateGenerator := template.NewTemplateGenerator(upstream, dashboardJsonTemplate)
				uid := templateGenerator.GenerateUid()

				dashboardClient.EXPECT().
					GetRawDashboard(uid).
					Return(nil, grafana.BoardProperties{}, grafana.DashboardNotFound(uid))

				dashBytes, err := templateGenerator.GenerateDashboard() // should just return "test-json"
				Expect(err).NotTo(HaveOccurred())
				snapshotBytes, err := templateGenerator.GenerateSnapshot() // should just return "test-json"
				Expect(err).NotTo(HaveOccurred())

				dashboardClient.EXPECT().
					SetRawDashboard(dashBytes).
					Return(nil)
				snapshotClient.EXPECT().
					SetRawSnapshot(snapshotBytes).
					Return(nil, nil) // we don't consume the snapshot response apart from checking the error
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return([]grafana.SnapshotListResponse{}, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return([]grafana.FoundBoard{}, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: upstreamList, // we just init'd two new upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("rebuilds dashboards when they exist already", func() {
			for index, upstream := range upstreamList {
				templateGenerator := template.NewTemplateGenerator(upstream, dashboardJsonTemplate)
				renderedDashboard, err := templateGenerator.GenerateDashboard()
				Expect(err).NotTo(HaveOccurred())

				uid := templateGenerator.GenerateUid()

				// have to set the key "id" in the json that gets returned by GetRawDashboard
				jsonWithId := addIdToDashboardJson(renderedDashboard, index)

				dashboardClient.EXPECT().
					GetRawDashboard(uid).
					Return(jsonWithId, grafana.BoardProperties{}, nil)
				dashboardClient.EXPECT().
					GetDashboardVersions(float64(index)).
					Return([]*grafana.Version{
						{
							VersionId: 1,
							Message:   template.DefaultCommitMessage,
						},
					}, nil)

				snapshotBytes, err := templateGenerator.GenerateSnapshot()
				Expect(err).NotTo(HaveOccurred())

				dashboardClient.EXPECT().
					SetRawDashboard(renderedDashboard).
					Return(nil)
				snapshotClient.EXPECT().
					SetRawSnapshot(snapshotBytes).
					Return(nil, nil) // we don't consume the snapshot response apart from checking the error
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return([]grafana.SnapshotListResponse{}, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return([]grafana.FoundBoard{}, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: upstreamList, // we just init'd two new upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not regenerate the dashboard when it has been edited by a human", func() {
			for index, upstream := range upstreamList {
				templateGenerator := template.NewTemplateGenerator(upstream, dashboardJsonTemplate)
				renderedDashboard, err := templateGenerator.GenerateDashboard()
				Expect(err).NotTo(HaveOccurred())

				uid := templateGenerator.GenerateUid()

				// have to set the key "id" in the json that gets returned by GetRawDashboard
				jsonWithId := addIdToDashboardJson(renderedDashboard, index)

				dashboardClient.EXPECT().
					GetRawDashboard(uid).
					Return(jsonWithId, grafana.BoardProperties{}, nil)
				dashboardClient.EXPECT().
					GetDashboardVersions(float64(index)).
					Return([]*grafana.Version{
						{
							VersionId: 1,
							Message:   "This is a message written by a human",
						},
					}, nil)
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return([]grafana.SnapshotListResponse{}, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return([]grafana.FoundBoard{}, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: upstreamList, // we just init'd two new upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes the dashboard for a deleted upstream", func() {
			var snapshots []grafana.SnapshotListResponse
			var dashboards []grafana.FoundBoard
			for _, upstream := range upstreamList {
				generator := template.NewTemplateGenerator(upstream, dashboardJsonTemplate)
				uid := generator.GenerateUid()
				snapshots = append(snapshots, grafana.SnapshotListResponse{
					Key: uid,
				})
				dashboards = append(dashboards, grafana.FoundBoard{
					UID: uid,
				})

				dashboardClient.EXPECT().
					DeleteDashboard(uid).
					Return(grafana.StatusMessage{}, nil)
				snapshotClient.EXPECT().
					DeleteSnapshot(uid).
					Return(nil)
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return(snapshots, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return(dashboards, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: []*gloov1.Upstream{}, // deleted both the upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
