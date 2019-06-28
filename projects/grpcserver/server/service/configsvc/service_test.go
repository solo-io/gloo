package configsvc_test

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_namespace "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/configsvc"
	. "github.com/solo-io/solo-projects/projects/grpcserver/server/service/internal/testutils"
	"google.golang.org/grpc"
)

var (
	grpcServer        *grpc.Server
	conn              *grpc.ClientConn
	apiserver         v1.ConfigApiServer
	client            v1.ConfigApiClient
	mockCtrl          *gomock.Controller
	settingsClient    *mock_gloo.MockSettingsClient
	licenseClient     *mock_license.MockClient
	namespaceClient   *mock_namespace.MockNamespaceClient
	testVersion       = "test-version"
	testOAuthEndpoint = v1.OAuthEndpoint{Url: "test", ClientName: "name"}
	testErr           = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		settingsClient = mock_gloo.NewMockSettingsClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		namespaceClient = mock_namespace.NewMockNamespaceClient(mockCtrl)
		apiserver = configsvc.NewConfigGrpcService(context.TODO(), settingsClient, licenseClient, namespaceClient, testOAuthEndpoint, testVersion)

		grpcServer, conn = MustRunGrpcServer(func(s *grpc.Server) { v1.RegisterConfigApiServer(s, apiserver) })
		client = v1.NewConfigApiClient(conn)
	})

	AfterEach(func() {
		grpcServer.Stop()
		mockCtrl.Finish()
	})

	Describe("GetVersion", func() {
		It("works", func() {
			actual, err := client.GetVersion(context.TODO(), &v1.GetVersionRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetVersionResponse{Version: testVersion}
			ExpectEqualProtoMessages(actual, expected)
		})
	})

	Describe("GetOAuthEndpoint", func() {
		It("works", func() {
			actual, err := client.GetOAuthEndpoint(context.TODO(), &v1.GetOAuthEndpointRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetOAuthEndpointResponse{OAuthEndpoint: &testOAuthEndpoint}
			ExpectEqualProtoMessages(actual, expected)
		})
	})

	Describe("GetIsLicenseValid", func() {
		It("works when the license client works", func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			actual, err := client.GetIsLicenseValid(context.TODO(), &v1.GetIsLicenseValidRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetIsLicenseValidResponse{IsLicenseValid: true}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the license client errors", func() {
			licenseClient.EXPECT().IsLicenseValid().Return(testErr)

			_, err := client.GetIsLicenseValid(context.TODO(), &v1.GetIsLicenseValidRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.LicenseIsInvalidError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("GetSettings", func() {
		It("works when the settings client works", func() {
			settings := &gloov1.Settings{RefreshRate: &types.Duration{Seconds: 1}}

			settingsClient.EXPECT().
				Read(defaults.GlooSystem, defaults.SettingsName, clients.ReadOpts{Ctx: context.TODO()}).
				Return(settings, nil)

			actual, err := client.GetSettings(context.TODO(), &v1.GetSettingsRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetSettingsResponse{Settings: settings}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the settings client errors", func() {
			settingsClient.EXPECT().
				Read(defaults.GlooSystem, defaults.SettingsName, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := client.GetSettings(context.TODO(), &v1.GetSettingsRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToReadSettingsError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateSettings", func() {
		It("works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}
			refreshRate := types.Duration{Seconds: 1}
			watchNamespaces := []string{"a", "b"}
			request := &v1.UpdateSettingsRequest{
				Ref:                &ref,
				RefreshRate:        &refreshRate,
				WatchNamespaceList: watchNamespaces,
			}
			readSettings := &gloov1.Settings{
				Metadata: core.Metadata{
					Namespace: ref.Namespace,
					Name:      ref.Name,
				},
				RefreshRate:     &types.Duration{Seconds: 50},
				WatchNamespaces: []string{"x", "y"},
			}
			writeSettings := &gloov1.Settings{
				Metadata: core.Metadata{
					Namespace: ref.Namespace,
					Name:      ref.Name,
				},
				RefreshRate:     &refreshRate,
				WatchNamespaces: watchNamespaces,
			}

			settingsClient.EXPECT().
				Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(readSettings, nil)
			settingsClient.EXPECT().
				Write(writeSettings, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(writeSettings, nil)

			actual, err := client.UpdateSettings(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateSettingsResponse{Settings: writeSettings}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("is a no-op when only ref is provided", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}
			request := &v1.UpdateSettingsRequest{Ref: &ref}
			settings := &gloov1.Settings{
				Metadata: core.Metadata{
					Namespace: ref.Namespace,
					Name:      ref.Name,
				},
				RefreshRate:     &types.Duration{Seconds: 50},
				WatchNamespaces: []string{"x", "y"},
			}

			settingsClient.EXPECT().
				Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(settings, nil)
			settingsClient.EXPECT().
				Write(settings, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(settings, nil)

			actual, err := client.UpdateSettings(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateSettingsResponse{Settings: settings}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the provided refresh rate is invalid", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}
			refreshRate := types.Duration{Nanos: 1}
			request := &v1.UpdateSettingsRequest{
				Ref:         &ref,
				RefreshRate: &refreshRate,
			}
			readSettings := &gloov1.Settings{}

			settingsClient.EXPECT().
				Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(readSettings, nil)

			_, err := client.UpdateSettings(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.InvalidRefreshRateError(1 * time.Nanosecond)
			wrapped := configsvc.FailedToUpdateSettingsError(expectedErr)
			Expect(err.Error()).To(ContainSubstring(wrapped.Error()))
		})

		It("errors when the settings client fails to read", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}
			request := &v1.UpdateSettingsRequest{
				Ref: &ref,
			}

			settingsClient.EXPECT().
				Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := client.UpdateSettings(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToUpdateSettingsError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the settings client fails to write", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}
			request := &v1.UpdateSettingsRequest{Ref: &ref}
			settings := &gloov1.Settings{}

			settingsClient.EXPECT().
				Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(settings, nil)
			settingsClient.EXPECT().
				Write(settings, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			_, err := client.UpdateSettings(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToUpdateSettingsError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListNamespaces", func() {
		It("works when the namespace client works", func() {
			namespaceList := []string{"a", "b", "c"}

			namespaceClient.EXPECT().
				ListNamespaces().
				Return(namespaceList, nil)

			actual, err := client.ListNamespaces(context.TODO(), &v1.ListNamespacesRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListNamespacesResponse{NamespaceList: namespaceList}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the namespace client errors", func() {
			namespaceClient.EXPECT().
				ListNamespaces().
				Return(nil, testErr)

			_, err := client.ListNamespaces(context.TODO(), &v1.ListNamespacesRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToListNamespacesError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
