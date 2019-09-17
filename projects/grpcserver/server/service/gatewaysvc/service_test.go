package gatewaysvc_test

import (
	"context"

	clientmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"

	"google.golang.org/grpc/codes"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gatewayv2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status"
	mock_status_converter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/gatewaysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/gatewaysvc/mocks"
)

var (
	apiserver       v1.GatewayApiServer
	mockCtrl        *gomock.Controller
	gatewayClient   *mocks.MockGatewayClient
	licenseClient   *mock_license.MockClient
	rawGetter       *mock_rawgetter.MockRawGetter
	settingsValues  *mock_settings.MockValuesClient
	clientCache     *clientmocks.MockClientCache
	statusConverter *mock_status_converter.MockInputResourceStatusGetter
	testErr         = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	getRaw := func(gateway *gatewayv2.Gateway) *v1.Raw {
		return &v1.Raw{FileName: gateway.GetMetadata().Name}
	}

	getStatus := func(code v1.Status_Code, message string) *v1.Status {
		return &v1.Status{
			Code:    code,
			Message: message,
		}
	}

	getGatewayDetails := func(gateway *gatewayv2.Gateway, status *v1.Status) *v1.GatewayDetails {
		return &v1.GatewayDetails{
			Gateway: gateway,
			Raw:     getRaw(gateway),
			Status:  status,
		}
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		gatewayClient = mocks.NewMockGatewayClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		clientCache = clientmocks.NewMockClientCache(mockCtrl)
		clientCache.EXPECT().GetGatewayClient().Return(gatewayClient).AnyTimes()
		statusConverter = mock_status_converter.NewMockInputResourceStatusGetter(mockCtrl)
		apiserver = gatewaysvc.NewGatewayGrpcService(context.TODO(), clientCache, rawGetter, statusConverter, licenseClient, settingsValues)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetGateway", func() {
		It("works when the gateway client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			gateway := &gatewayv2.Gateway{
				Metadata: metadata,
				Status: core.Status{
					State: core.Status_Accepted,
				},
			}

			gatewayClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(gateway, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), gateway, gatewayv2.GatewayCrd).
				Return(getRaw(gateway))
			statusConverter.EXPECT().
				GetApiStatusFromResource(gateway).
				Return(getStatus(v1.Status_OK, ""))

			request := &v1.GetGatewayRequest{Ref: &ref}
			actual, err := apiserver.GetGateway(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetGatewayResponse{
				GatewayDetails: getGatewayDetails(gateway, getStatus(v1.Status_OK, "")),
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the gateway client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			gatewayClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetGatewayRequest{Ref: &ref}
			_, err := apiserver.GetGateway(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := gatewaysvc.FailedToGetGatewayError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListGateways", func() {
		It("works when the gateway client works", func() {
			ns1, ns2 := "one", "two"
			gateway1 := &gatewayv2.Gateway{
				Metadata: core.Metadata{Namespace: ns1},
				Status: core.Status{
					State: core.Status_Accepted,
				},
			}
			gateway2 := &gatewayv2.Gateway{
				Metadata: core.Metadata{Namespace: ns2},
				Status: core.Status{
					State: core.Status_Pending,
				},
			}

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns1, ns2})
			gatewayClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv2.Gateway{gateway1}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), gateway1, gatewayv2.GatewayCrd).
				Return(getRaw(gateway1))
			gatewayClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv2.Gateway{gateway2}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), gateway2, gatewayv2.GatewayCrd).
				Return(getRaw(gateway2))
			statusConverter.EXPECT().
				GetApiStatusFromResource(gateway1).
				Return(getStatus(v1.Status_OK, ""))
			statusConverter.EXPECT().
				GetApiStatusFromResource(gateway2).
				Return(getStatus(v1.Status_WARNING, status.ResourcePending(ns2, "")))

			actual, err := apiserver.ListGateways(context.TODO(), &v1.ListGatewaysRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListGatewaysResponse{
				GatewayDetails: []*v1.GatewayDetails{
					getGatewayDetails(gateway1, getStatus(v1.Status_OK, "")),
					getGatewayDetails(gateway2, getStatus(v1.Status_WARNING, status.ResourcePending(ns2, ""))),
				},
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the gateway client errors", func() {
			ns := "ns"

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns})
			gatewayClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := apiserver.ListGateways(context.TODO(), &v1.ListGatewaysRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := gatewaysvc.FailedToListGatewaysError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateGatewayYaml", func() {
		It("works on valid input", func() {
			yamlString := "totally-valid-yaml"

			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			gateway := &gatewayv2.Gateway{
				Metadata: metadata,
				Status: core.Status{
					State: core.Status_Accepted,
				},
			}
			labels := map[string]string{"test-label": "test-value"}

			action := func(ctx context.Context,
				yamlString string,
				refToValidate *core.ResourceRef,
				emptyInputResource resources.InputResource,
			) error {
				gateway.Metadata.Labels = labels
				return nil
			}

			licenseClient.EXPECT().IsLicenseValid().Return(nil)
			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				DoAndReturn(action)

			gatewayClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(gateway, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), gateway, gatewayv2.GatewayCrd).
				Return(getRaw(gateway))
			statusConverter.EXPECT().
				GetApiStatusFromResource(gateway).
				Return(getStatus(v1.Status_WARNING, status.ResourcePending("ns", "name")))

			response, err := apiserver.UpdateGatewayYaml(context.TODO(), &v1.UpdateGatewayYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			})

			Expect(err).NotTo(HaveOccurred())
			ExpectEqualProtoMessages(response.GatewayDetails.Gateway, gateway)
		})

		It("fails when the client fails", func() {
			yamlString := "totally-valid-yaml"

			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			gateway := &gatewayv2.Gateway{
				Metadata: metadata,
				Status: core.Status{
					State: core.Status_Accepted,
				},
			}
			labels := map[string]string{"test-label": "test-value"}

			action := func(ctx context.Context,
				yamlString string,
				refToValidate *core.ResourceRef,
				emptyInputResource resources.InputResource,
			) error {
				gateway.Metadata.Labels = labels
				return nil
			}

			licenseClient.EXPECT().IsLicenseValid().Return(nil)
			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				DoAndReturn(action)
			gatewayClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			response, err := apiserver.UpdateGatewayYaml(context.TODO(), &v1.UpdateGatewayYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			})

			Expect(response).To(BeNil())
			Expect(err.Error()).To(ContainSubstring(gatewaysvc.FailedToUpdateGatewayError(testErr, &ref).Error()))
		})

		It("fails on invalid license", func() {
			yamlString := "totally-valid-yaml"

			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			licenseClient.EXPECT().IsLicenseValid().Return(testErr)

			response, err := apiserver.UpdateGatewayYaml(context.TODO(), &v1.UpdateGatewayYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			})

			Expect(response).To(BeNil())
			Expect(err.Error()).To(ContainSubstring(codes.PermissionDenied.String()))
		})

		It("fails on invalid input", func() {
			yamlString := "totally-broken-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(testErr)

			response, err := apiserver.UpdateGatewayYaml(context.TODO(), &v1.UpdateGatewayYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					EditedYaml: yamlString,
					Ref:        &ref,
				},
			})

			Expect(response).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(gatewaysvc.FailedToParseGatewayFromYamlError(testErr, &ref)))
		})
	})

	Describe("UpdateGateway", func() {
		var metadata core.Metadata
		var ref core.ResourceRef
		var existing, input, toWrite *gatewayv2.Gateway

		BeforeEach(func() {
			metadata = core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref = metadata.Ref()
			existing = &gatewayv2.Gateway{
				Metadata: core.Metadata{
					Namespace:       "ns",
					Name:            "name",
					ResourceVersion: "10",
				},
				BindAddress: "test-old-value",
			}
			input = &gatewayv2.Gateway{
				Metadata:    metadata,
				BindAddress: "test-new-value",
				Status:      core.Status{State: 1},
			}
			toWrite = &gatewayv2.Gateway{
				Metadata: core.Metadata{
					Namespace:       "ns",
					Name:            "name",
					ResourceVersion: "10",
				},
				Status:      core.Status{State: core.Status_Pending},
				BindAddress: "test-new-value",
			}
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})

		It("works when the gateway client works", func() {
			gatewayClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(existing, nil)
			gatewayClient.EXPECT().
				Write(toWrite, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(toWrite, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), toWrite, gatewayv2.GatewayCrd).
				Return(getRaw(toWrite))
			statusConverter.EXPECT().
				GetApiStatusFromResource(toWrite).
				Return(getStatus(v1.Status_WARNING, status.ResourcePending("ns", "name")))

			request := &v1.UpdateGatewayRequest{Gateway: input}
			actual, err := apiserver.UpdateGateway(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateGatewayResponse{GatewayDetails: getGatewayDetails(toWrite, getStatus(v1.Status_WARNING, status.ResourcePending("ns", "name")))}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the gateway client errors on read", func() {
			gatewayClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.UpdateGatewayRequest{Gateway: input}
			_, err := apiserver.UpdateGateway(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := gatewaysvc.FailedToUpdateGatewayError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the gateway client errors on write", func() {
			gatewayClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(existing, nil)
			gatewayClient.EXPECT().
				Write(toWrite, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			request := &v1.UpdateGatewayRequest{Gateway: input}
			_, err := apiserver.UpdateGateway(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := gatewaysvc.FailedToUpdateGatewayError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
