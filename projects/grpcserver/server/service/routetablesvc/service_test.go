package routetablesvc_test

import (
	"context"

	clientmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"

	"google.golang.org/grpc/codes"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/routetablesvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/routetablesvc/mocks"
)

var (
	apiserver             v1.RouteTableApiServer
	mockCtrl              *gomock.Controller
	routeTableClient      *mocks.MockRouteTableClient
	licenseClient         *mock_license.MockClient
	settingsValues        *mock_settings.MockValuesClient
	clientCache           *clientmocks.MockClientCache
	detailsExpectation    *gomock.Call
	routeTable            *gatewayv1.RouteTable
	routeTableDetails     *v1.RouteTableDetails
	rawGetter             *mock_rawgetter.MockRawGetter
	testErr               = errors.Errorf("test-err")
	uint32Zero, uint32One = uint32(0), uint32(1)
	metadata              = core.Metadata{
		Namespace: "ns",
		Name:      "name",
	}
	ref = metadata.Ref()
)

var _ = Describe("ServiceTest", func() {
	getRaw := func(name string) *v1.Raw {
		return &v1.Raw{
			FileName: name + ".yaml",
		}
	}
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		routeTableClient = mocks.NewMockRouteTableClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		clientCache = clientmocks.NewMockClientCache(mockCtrl)
		clientCache.EXPECT().GetRouteTableClient().Return(routeTableClient).AnyTimes()
		apiserver = routetablesvc.NewRouteTableGrpcService(
			context.TODO(),
			clientCache,
			licenseClient,
			settingsValues,
			rawGetter,
		)

		routeTable = &gatewayv1.RouteTable{Metadata: metadata}
		routeTableDetails = &v1.RouteTableDetails{RouteTable: routeTable, Raw: getRaw(routeTable.Metadata.Name)}

	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetRouteTable", func() {
		It("works when the route table client works", func() {
			routeTableClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(routeTable, nil)
			detailsExpectation = rawGetter.EXPECT().
				GetRaw(context.TODO(), routeTable, gatewayv1.RouteTableCrd).
				Return(getRaw(routeTable.Metadata.Name))
			detailsExpectation.Times(1)

			request := &v1.GetRouteTableRequest{Ref: &ref}
			actual, err := apiserver.GetRouteTable(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetRouteTableResponse{RouteTableDetails: routeTableDetails}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the route table client errors", func() {
			routeTableClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetRouteTableRequest{Ref: &ref}
			_, err := apiserver.GetRouteTable(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := routetablesvc.FailedToReadRouteTableError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListRouteTables", func() {
		It("works", func() {
			ns1, ns2 := "one", "two"
			name1, name2 := "a", "b"
			routeTable1 := &gatewayv1.RouteTable{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: core.Metadata{Namespace: ns1, Name: name1},
			}
			routeTable2 := &gatewayv1.RouteTable{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: core.Metadata{Namespace: ns2, Name: name2},
			}
			routeTableDetails1 := &v1.RouteTableDetails{
				RouteTable: routeTable1,
				Raw:        getRaw(routeTable1.Metadata.Name),
			}
			routeTableDetails2 := &v1.RouteTableDetails{
				RouteTable: routeTable2,
				Raw:        getRaw(routeTable2.Metadata.Name),
			}

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns1, ns2})
			routeTableClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv1.RouteTable{routeTable1}, nil)
			routeTableClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv1.RouteTable{routeTable2}, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), routeTable1, gatewayv1.RouteTableCrd).
				Return(getRaw(routeTable1.Metadata.Name))
			rawGetter.EXPECT().
				GetRaw(context.TODO(), routeTable2, gatewayv1.RouteTableCrd).
				Return(getRaw(routeTable2.Metadata.Name))

			actual, err := apiserver.ListRouteTables(context.TODO(), &v1.ListRouteTablesRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListRouteTablesResponse{
				RouteTableDetails: []*v1.RouteTableDetails{routeTableDetails1, routeTableDetails2},
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the route table client errors", func() {
			ns := "ns"

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns})
			routeTableClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)
			detailsExpectation.Times(0)

			_, err := apiserver.ListRouteTables(context.TODO(), &v1.ListRouteTablesRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := routetablesvc.FailedToListRouteTablesError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateRouteTable", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		It("works", func() {
			routeTableClient.EXPECT().
				Write(gomock.Any(), gomock.Any()).
				Return(routeTable, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), routeTable, gatewayv1.RouteTableCrd).
				Return(getRaw(routeTable.Metadata.Name))
			routetableDetails := &v1.RouteTableDetails{
				RouteTable: routeTable,
				Raw:        getRaw(routeTable.Metadata.Name),
			}

			actual, err := apiserver.CreateRouteTable(context.TODO(), &v1.CreateRouteTableRequest{RouteTable: routeTable})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.CreateRouteTableResponse{RouteTableDetails: routetableDetails}
			ExpectEqualProtoMessages(actual, expected)
		})
	})

	Describe("UpdateRouteTableYaml", func() {
		It("works on valid input", func() {
			yamlString := "totally-valid-yaml"

			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}
			labels := map[string]string{"test-label": "test-value"}

			action := func(ctx context.Context,
				yamlString string,
				refToValidate *core.ResourceRef,
				emptyInputResource resources.InputResource,
			) error {
				routeTable.Metadata.Labels = labels
				return nil
			}

			rawGetter.EXPECT().
				GetRaw(context.TODO(), routeTable, gatewayv1.RouteTableCrd).
				Return(getRaw(routeTable.Metadata.Name))
			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, ref, gomock.Any()).
				DoAndReturn(action)
			routeTableClient.EXPECT().
				Write(gomock.Any(), gomock.Any()).
				Return(routeTable, nil)
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			response, err := apiserver.UpdateRouteTableYaml(context.TODO(), &v1.UpdateRouteTableYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(err).NotTo(HaveOccurred())
			updatedRouteTableResponce := &v1.UpdateRouteTableResponse{
				RouteTableDetails: routeTableDetails,
			}
			ExpectEqualProtoMessages(response, updatedRouteTableResponce)
		})

		It("fails when the license is invalid", func() {
			yamlString := "valid-yaml"
			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}

			licenseClient.EXPECT().IsLicenseValid().Return(testErr)

			response, err := apiserver.UpdateRouteTableYaml(context.TODO(), &v1.UpdateRouteTableYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(response).To(BeNil())
			Expect(err.Error()).To(ContainSubstring(codes.PermissionDenied.String()))
		})

		It("fails when the client fails", func() {
			yamlString := "totally-valid-yaml"

			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}
			labels := map[string]string{"test-label": "test-value"}

			action := func(ctx context.Context,
				yamlString string,
				refToValidate *core.ResourceRef,
				emptyInputResource resources.InputResource,
			) error {
				routeTable.Metadata.Labels = labels
				return nil
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, ref, gomock.Any()).
				DoAndReturn(action)
			routeTableClient.EXPECT().
				Write(gomock.Any(), gomock.Any()).
				Return(nil, testErr)
			detailsExpectation.Times(0)
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			response, err := apiserver.UpdateRouteTableYaml(context.TODO(), &v1.UpdateRouteTableYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(response).To(BeNil())
			Expect(err.Error()).To(ContainSubstring(routetablesvc.FailedToUpdateRouteTableError(
				testErr, ref.Namespace, ref.Name).Error()))
		})

		It("fails on invalid input", func() {
			yamlString := "totally-broken-yaml"
			ref := &core.ResourceRef{
				Name:      "service",
				Namespace: "gloo-system",
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, ref, gomock.Any()).
				Return(testErr)
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			// should not expect details
			detailsExpectation = rawGetter.EXPECT().
				GetRaw(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(getRaw(""))
			detailsExpectation.Times(0)

			response, err := apiserver.UpdateRouteTableYaml(context.TODO(), &v1.UpdateRouteTableYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        ref,
					EditedYaml: yamlString,
				},
			})

			Expect(err.Error()).To(Equal(routetablesvc.FailedToParseRouteTableFromYamlError(testErr, ref).Error()))
			Expect(response).To(BeNil())
		})
	})

	Describe("UpdateRouteTable", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		It("works", func() {
			routeTableClient.EXPECT().
				Write(gomock.Any(), gomock.Any()).
				Return(routeTable, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), routeTable, gatewayv1.RouteTableCrd).
				Return(getRaw(routeTable.Metadata.Name))
			actual, err := apiserver.UpdateRouteTable(context.TODO(), &v1.UpdateRouteTableRequest{RouteTable: routeTable})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateRouteTableResponse{RouteTableDetails: routeTableDetails}
			ExpectEqualProtoMessages(actual, expected)

		})
	})

	Describe("DeleteRouteTable", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		It("works when the route table client works", func() {
			routeTableClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)
			detailsExpectation.Times(0)

			request := &v1.DeleteRouteTableRequest{Ref: &ref}
			actual, err := apiserver.DeleteRouteTable(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteRouteTableResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the route table client errors", func() {
			routeTableClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)
			detailsExpectation.Times(0)

			request := &v1.DeleteRouteTableRequest{Ref: &ref}
			_, err := apiserver.DeleteRouteTable(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := routetablesvc.FailedToDeleteRouteTableError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
