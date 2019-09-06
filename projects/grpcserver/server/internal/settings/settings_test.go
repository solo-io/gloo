package settings_test

import (
	"context"
	"time"

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
	testErr        = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		settingsClient = mock_gloo.NewMockSettingsClient(mockCtrl)
		client = settings.NewSettingsValuesClient(context.TODO(), settingsClient)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetRefreshRate", func() {
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
				desc:         "with zero value in settings settings",
			},
			{
				readSettings: &v1.Settings{},
				expected:     settings.DefaultRefreshRate,
				desc:         "with nil refresh rate in settings settings",
			},
			{
				readErr:  testErr,
				expected: settings.DefaultRefreshRate,
				desc:     "with error on read",
			},
		}

		It("works", func() {
			for _, tc := range testCases {
				settingsClient.EXPECT().
					Read("gloo-system", "default", clients.ReadOpts{Ctx: context.TODO()}).
					Return(tc.readSettings, tc.readErr)

				actual := client.GetRefreshRate()
				Expect(actual).To(Equal(tc.expected), tc.desc)
			}
		})
	})
})
