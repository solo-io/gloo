package search_test

import (
	"context"
	"encoding/json"

	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	gatewaymocks "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search/mocks"

	clientmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search"
	vsmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mocks"
)

var (
	mockCtrl             *gomock.Controller
	upstreamGroupClient  *mocks.MockUpstreamGroupClient
	virtualServiceClient *vsmocks.MockVirtualServiceClient
	routeTableClient     *gatewaymocks.MockRouteTableClient
	clientCache          *clientmocks.MockClientCache
	settingsValues       *mock_settings.MockValuesClient
	upstreamSearcher     search.UpstreamSearcher
	testErr              = errors.Errorf("test-err")
	listOpts             = clients.ListOpts{Ctx: context.TODO()}
	allVirtualServices   gatewayv1.VirtualServiceList
	allUpstreamGroups    gloov1.UpstreamGroupList
	allRouteTables       gatewayv1.RouteTableList
	testNamespace        = "test-ns"
)

var _ = Describe("Upstream Search Test", func() {
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		upstreamGroupClient = mocks.NewMockUpstreamGroupClient(mockCtrl)
		virtualServiceClient = vsmocks.NewMockVirtualServiceClient(mockCtrl)
		routeTableClient = gatewaymocks.NewMockRouteTableClient(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		clientCache = clientmocks.NewMockClientCache(mockCtrl)

		settingsValues.EXPECT().GetWatchNamespaces().Return([]string{testNamespace}).AnyTimes()

		clientCache.EXPECT().GetVirtualServiceClient().Return(virtualServiceClient).AnyTimes()
		clientCache.EXPECT().GetUpstreamGroupClient().Return(upstreamGroupClient).AnyTimes()
		clientCache.EXPECT().GetRouteTableClient().Return(routeTableClient).AnyTimes()

		upstreamSearcher = search.NewUpstreamSearcher(clientCache, settingsValues)
		allVirtualServices = nil
		allUpstreamGroups = nil
		allRouteTables = nil

		var untypedVirtualServices []interface{}
		virtualServiceListRaw := helpers.MustReadFile("fixtures/test_virtual_service_list.json")
		err := json.Unmarshal(virtualServiceListRaw, &untypedVirtualServices)

		Expect(err).NotTo(HaveOccurred())

		for _, rawVirtualService := range untypedVirtualServices {
			vsMap := rawVirtualService.(map[string]interface{})
			var virtualService gatewayv1.VirtualService
			err := protoutils.UnmarshalMap(vsMap, &virtualService)
			Expect(err).NotTo(HaveOccurred())

			allVirtualServices = append(allVirtualServices, &virtualService)
		}

		var untypedUpstreamGroups []interface{}
		upstreamGroupListRaw := helpers.MustReadFile("fixtures/test_upstream_group_list.json")
		err = json.Unmarshal(upstreamGroupListRaw, &untypedUpstreamGroups)

		Expect(err).NotTo(HaveOccurred())

		for _, rawUpstreamGroup := range untypedUpstreamGroups {
			ugMap := rawUpstreamGroup.(map[string]interface{})
			var upstreamGroup gloov1.UpstreamGroup
			err := protoutils.UnmarshalMap(ugMap, &upstreamGroup)
			Expect(err).NotTo(HaveOccurred())

			allUpstreamGroups = append(allUpstreamGroups, &upstreamGroup)
		}

		var untypedRouteTables []interface{}
		routeTableListRaw := helpers.MustReadFile("fixtures/test_route_table_list.json")
		err = json.Unmarshal(routeTableListRaw, &untypedRouteTables)

		Expect(err).NotTo(HaveOccurred())

		for _, rawRouteTable := range untypedRouteTables {
			rtMap := rawRouteTable.(map[string]interface{})
			var routeTable gatewayv1.RouteTable
			err := protoutils.UnmarshalMap(rtMap, &routeTable)
			Expect(err).NotTo(HaveOccurred())

			allRouteTables = append(allRouteTables, &routeTable)
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("FindContainingVirtualServices", func() {
		It("does not find upstreams that are unreferenced in any virtual service", func() {
			upstreamRef := &core.ResourceRef{
				Name:      "my-unreferenced-upstream-name",
				Namespace: "ns",
			}

			virtualServiceClient.EXPECT().
				List(testNamespace, listOpts).
				Return(allVirtualServices, nil)

			upstreamGroupClient.EXPECT().
				List(testNamespace, listOpts).
				Return(allUpstreamGroups, nil)

			routeTableClient.EXPECT().
				List(testNamespace, listOpts).
				Return(allRouteTables, nil)

			foundVirtualServiceRef, err := upstreamSearcher.FindContainingVirtualServices(context.TODO(), upstreamRef)

			Expect(err).NotTo(HaveOccurred())
			Expect(foundVirtualServiceRef).To(BeNil())
		})

		It("errors when the virtual service client errors", func() {
			upstreamRef := &core.ResourceRef{
				Name:      "my-unreferenced-upstream-name",
				Namespace: "ns",
			}

			virtualServiceClient.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, testErr)
			upstreamGroupClient.EXPECT().List(testNamespace, gomock.Any()).Times(0)

			foundVirtualServiceRef, err := upstreamSearcher.FindContainingVirtualServices(context.TODO(), upstreamRef)

			Expect(err).To(Equal(testErr))
			Expect(foundVirtualServiceRef).To(BeNil())
		})

		It("errors when the upstream group client errors", func() {
			upstreamRef := &core.ResourceRef{
				Name:      "my-unreferenced-upstream-name",
				Namespace: "ns",
			}

			virtualServiceClient.EXPECT().
				List(testNamespace, listOpts).
				Return(allVirtualServices, nil)
			upstreamGroupClient.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, testErr)
			upstreamGroupClient.EXPECT().List(testNamespace, gomock.Any()).Times(0)

			foundVirtualServiceRef, err := upstreamSearcher.FindContainingVirtualServices(context.TODO(), upstreamRef)

			Expect(err).To(Equal(testErr))
			Expect(foundVirtualServiceRef).To(BeNil())
		})

		It("finds upstreams that are referenced in a virtual service", func() {
			upstreamOneRef := &core.ResourceRef{
				Name:      "upstream1",
				Namespace: "ns",
			}
			upstreamTwoRef := &core.ResourceRef{
				Name:      "upstream2",
				Namespace: "ns",
			}
			upstreamThreeRef := &core.ResourceRef{
				Name:      "upstream3",
				Namespace: "ns",
			}
			upstreamFourRef := &core.ResourceRef{
				Name:      "upstream4",
				Namespace: "ns",
			}
			virtualServiceOneMetadata := &core.Metadata{
				Name:      "vs1",
				Namespace: "ns",
			}
			virtualServiceTwoMetadata := &core.Metadata{
				Name:      "vs2",
				Namespace: "ns",
			}
			virtualServiceThreeMetadata := &core.Metadata{
				Name:      "vs3",
				Namespace: "ns",
			}
			virtualServiceFourMetadata := &core.Metadata{
				Name:      "vs4",
				Namespace: "ns",
			}
			testCases := []struct {
				upstreamRef                      *core.ResourceRef
				containingVirtualServiceMetadata *core.Metadata
			}{
				{
					upstreamRef:                      upstreamOneRef,
					containingVirtualServiceMetadata: virtualServiceOneMetadata,
				},
				{
					upstreamRef:                      upstreamTwoRef,
					containingVirtualServiceMetadata: virtualServiceTwoMetadata,
				},
				{
					upstreamRef:                      upstreamThreeRef,
					containingVirtualServiceMetadata: virtualServiceThreeMetadata,
				},
				{
					upstreamRef:                      upstreamFourRef,
					containingVirtualServiceMetadata: virtualServiceFourMetadata,
				},
			}

			virtualServiceClient.EXPECT().
				List(testNamespace, listOpts).
				Return(allVirtualServices, nil).
				AnyTimes()

			upstreamGroupClient.EXPECT().
				List(testNamespace, listOpts).
				Return(allUpstreamGroups, nil).
				AnyTimes()

			routeTableClient.EXPECT().
				List(testNamespace, listOpts).
				Return(allRouteTables, nil).
				AnyTimes()

			for _, testCase := range testCases {
				upstreamRef := testCase.upstreamRef
				expectedVirtualServiceRef := testCase.containingVirtualServiceMetadata.Ref()

				foundVirtualServiceRef, err := upstreamSearcher.FindContainingVirtualServices(context.TODO(), upstreamRef)

				Expect(err).NotTo(HaveOccurred())
				expectedArr := []*core.ResourceRef{&expectedVirtualServiceRef}
				Expect(foundVirtualServiceRef).To(Equal(expectedArr))
			}
		})
	})
})
