package settings_test

import (
	"context"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"
	"k8s.io/kubernetes/pkg/apis/core"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
)

var (
	client         settings.ValuesClient
	mockCtrl       *gomock.Controller
	settingsClient *mock_gloo.MockSettingsClient
	clientCache    *mocks.MockClientCache
	testErr        = errors.Errorf("test-err")
	podNamespace   = "test-ns"
)

var _ = Describe("ValuesClientTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		settingsClient = mock_gloo.NewMockSettingsClient(mockCtrl)
		clientCache = mocks.NewMockClientCache(mockCtrl)
		clientCache.EXPECT().GetSettingsClient().Return(settingsClient).AnyTimes()
		client = settings.NewSettingsValuesClient(context.TODO(), clientCache, podNamespace)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetRefreshRate", func() {
		It("works", func() {
			testCases := []struct {
				readSettings *v1.Settings
				readErr      error
				expected     time.Duration
				desc         string
			}{
				{
					readSettings: &v1.Settings{RefreshRate: types.DurationProto(5 * time.Second)},
					expected:     5 * time.Second,
					desc:         "with valid settings",
				},
				{
					readSettings: &v1.Settings{RefreshRate: types.DurationProto(0)},
					expected:     settings.DefaultRefreshRate,
					desc:         "with zero value in settings",
				},
				{
					readSettings: &v1.Settings{},
					expected:     settings.DefaultRefreshRate,
					desc:         "with nil refresh rate in settings",
				},
				{
					readErr:  testErr,
					expected: settings.DefaultRefreshRate,
					desc:     "with error on read",
				},
			}

			for _, tc := range testCases {
				settingsClient.EXPECT().
					Read(podNamespace, defaults.SettingsName, clients.ReadOpts{Ctx: context.Background()}).
					Return(tc.readSettings, tc.readErr)

				actual := client.GetRefreshRate()
				Expect(actual).To(Equal(tc.expected), tc.desc)
			}
		})
	})

	Describe("GetWatchNamespaces", func() {
		It("works", func() {
			testCases := []struct {
				readSettings *v1.Settings
				readErr      error
				expected     []string
				desc         string
			}{
				{
					readSettings: &v1.Settings{WatchNamespaces: []string{core.NamespaceAll}, DiscoveryNamespace: "disc"},
					expected:     []string{core.NamespaceAll},
					desc:         "with NamespaceAll in settings",
				},
				{
					readSettings: &v1.Settings{WatchNamespaces: []string{"one", "two"}, DiscoveryNamespace: "disc"},
					expected:     []string{"one", "two", "disc"},
					desc:         "with specific namespaces and discovery namespace excluded",
				},
				{
					readSettings: &v1.Settings{WatchNamespaces: []string{"one", "two", "disc"}, DiscoveryNamespace: "disc"},
					expected:     []string{"one", "two", "disc"},
					desc:         "with specific namespaces and discovery namespace included",
				},
				{
					readErr:  testErr,
					expected: []string{core.NamespaceAll},
					desc:     "defaults to all namespaces if the settings client errors",
				},
			}
			for _, tc := range testCases {
				settingsClient.EXPECT().
					Read(podNamespace, defaults.SettingsName, clients.ReadOpts{Ctx: context.Background()}).
					Return(tc.readSettings, tc.readErr)

				actual := client.GetWatchNamespaces()
				Expect(actual).To(Equal(tc.expected), tc.desc)
			}
		})
	})
})
