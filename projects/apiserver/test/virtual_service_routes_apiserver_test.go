package apiserver_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ghodss/yaml"
	"github.com/golang/mock/gomock"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gateway_solo_io_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	mock_gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1/mocks"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/rt_selector_handler"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("VirtualServiceRoutesApiServer", func() {

	var (
		ctrl                      *gomock.Controller
		ctx                       context.Context
		err                       error
		mockMCGatewayCRDClientset *mock_gateway_v1.MockMulticlusterClientset
		mockGatewayClientSet      *mock_gateway_v1.MockClientset

		mockVSClient                  *mock_gateway_v1.MockVirtualServiceClient
		mockRTClient                  *mock_gateway_v1.MockRouteTableClient
		virtualServiceRoutesApiServer rpc_edge_v1.VirtualServiceRoutesApiServer
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		mockMCGatewayCRDClientset = mock_gateway_v1.NewMockMulticlusterClientset(ctrl)
		mockGatewayClientSet = mock_gateway_v1.NewMockClientset(ctrl)
		mockVSClient = mock_gateway_v1.NewMockVirtualServiceClient(ctrl)
		mockRTClient = mock_gateway_v1.NewMockRouteTableClient(ctrl)

		mockGatewayClientSet.EXPECT().RouteTables().Return(mockRTClient).AnyTimes()
		mockGatewayClientSet.EXPECT().VirtualServices().Return(mockVSClient).AnyTimes()
		mockMCGatewayCRDClientset.EXPECT().Cluster("kind-local").Return(mockGatewayClientSet, nil).AnyTimes()

		virtualServiceRoutesApiServer = rt_selector_handler.NewVirtualServiceRoutesHandler(mockMCGatewayCRDClientset)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can GetVirtualServiceRoutes", func() {
		exampleVsYaml := `
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: 'example'
  namespace: 'gloo-system'
spec:
  virtualHost:
    domains:
    - 'example.com'
    routes:
    - matchers:
       - prefix: '/a'
      delegateAction:
        ref:
          name: 'a-routes'
          namespace: 'gloo-system'
    - matchers:
       - prefix: '/b'
      delegateAction:
        ref:
          name: 'b-routes'
          namespace: 'gloo-system'
`
		testVS := &gateway_solo_io_v1.VirtualService{}
		err = yaml.Unmarshal([]byte(exampleVsYaml), testVS)
		Expect(err).NotTo(HaveOccurred())

		routeTableAYaml := `
apiVersion: gateway.solo.io/v1
kind: RouteTable
metadata:
  name: 'a-routes'
  namespace: 'gloo-system'
spec:
  routes:
    - matchers:
       - prefix: '/a/1'
      routeAction:
        single:
          upstream:
            name: 'foo-upstream'
    - matchers:
       - prefix: '/a/2'
      routeAction:
        multi:
          destinations:
          - weight: 7
            destination:
              upstream:
                name: 'multi-upstream-1'
                namespace: 'gloo-system'
          - weight: 2
            destination:
              kube:
                ref:
                  name: 'kube-1'
                  namespace: 'default'
                port: 8080
          - weight: 1
            destination:
              consul:
                serviceName: 'consul-service'`

		testRTa := &gateway_solo_io_v1.RouteTable{}
		err = yaml.Unmarshal([]byte(routeTableAYaml), testRTa)
		Expect(err).NotTo(HaveOccurred())

		routeTableBYaml := `
apiVersion: gateway.solo.io/v1
kind: RouteTable
metadata:
  name: 'b-routes'
  namespace: 'gloo-system'
spec:
  routes:
    - matchers:
       - regex: '/b/3'
      routeAction:
        single:
          upstream:
            name: 'baz-upstream'
    - matchers:
       - regex: '/b/4'
      redirectAction:
          hostRedirect: "test.solo.io"
    - matchers:
       - prefix: '/b/c'
      delegateAction:
        ref:
          name: 'c-routes'
          namespace: 'gloo-system'`

		testRTb := &gateway_solo_io_v1.RouteTable{}
		err = yaml.Unmarshal([]byte(routeTableBYaml), testRTb)
		Expect(err).NotTo(HaveOccurred())

		routeTableCYaml := `
apiVersion: gateway.solo.io/v1
kind: RouteTable
metadata:
  name: 'c-routes'
  namespace: 'gloo-system'
spec:
  routes:
    - matchers:
       - exact: '/b/c/4'
      routeAction:
        single:
          upstream:
            name: 'qux-upstream'
    - matchers:
        - exact: /b/c/567
      directResponseAction:
        status: 200
        body: "Hello, world!"`

		testRTc := &gateway_solo_io_v1.RouteTable{}
		err = yaml.Unmarshal([]byte(routeTableCYaml), testRTc)
		Expect(err).NotTo(HaveOccurred())

		mockVSClient.EXPECT().GetVirtualService(ctx, client.ObjectKey{
			Namespace: "gloo-system",
			Name:      "example",
		}).Return(testVS, nil)
		mockRTClient.EXPECT().GetRouteTable(ctx, client.ObjectKey{
			Namespace: "gloo-system",
			Name:      "a-routes",
		}).Return(testRTa, nil)
		mockRTClient.EXPECT().GetRouteTable(ctx, client.ObjectKey{
			Namespace: "gloo-system",
			Name:      "b-routes",
		}).Return(testRTb, nil)
		mockRTClient.EXPECT().GetRouteTable(ctx, client.ObjectKey{
			Namespace: "gloo-system",
			Name:      "c-routes",
		}).Return(testRTc, nil)

		resp, err := virtualServiceRoutesApiServer.GetVirtualServiceRoutes(ctx, &rpc_edge_v1.GetVirtualServiceRoutesRequest{
			VirtualServiceRef: &v1.ClusterObjectRef{
				Name:        "example",
				Namespace:   "gloo-system",
				ClusterName: "kind-local",
			},
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(resp.VsRoutes)).To(Equal(2))
		Expect(resp.VsRoutes[0].GetDelegateAction().GetRef().GetName()).To(Equal("a-routes"))
		Expect(resp.VsRoutes[0].Matcher).To(Equal("/a"))
		Expect(resp.VsRoutes[0].MatchType).To(Equal("PREFIX"))
		Expect(len(resp.VsRoutes[0].RtRoutes)).To(Equal(2))
		rtA := resp.VsRoutes[0].RtRoutes
		Expect(rtA[0].GetRouteAction().GetSingle().GetUpstream().GetName()).To(Equal("foo-upstream"))
		Expect(rtA[0].Matcher).To(Equal("/a/1"))
		Expect(rtA[0].MatchType).To(Equal("PREFIX"))
		multiDest := rtA[1].GetRouteAction().GetMulti().GetDestinations()
		Expect(len(multiDest)).To(Equal(3))
		Expect(multiDest[0].GetWeight()).To(Equal(uint32(7)))
		Expect(multiDest[0].GetDestination().GetUpstream().GetName()).To(Equal("multi-upstream-1"))
		Expect(multiDest[1].GetWeight()).To(Equal(uint32(2)))
		kubeRef := multiDest[1].GetDestination().GetKube().GetRef()
		Expect(kubeRef.GetName()).To(Equal("kube-1"))
		Expect(multiDest[2].GetWeight()).To(Equal(uint32(1)))
		Expect(multiDest[2].GetDestination().GetConsul().GetServiceName()).To(Equal("consul-service"))
		Expect(rtA[1].Matcher).To(Equal("/a/2"))
		Expect(rtA[1].MatchType).To(Equal("PREFIX"))
		Expect(resp.VsRoutes[1].GetDelegateAction().GetRef().GetName()).To(Equal("b-routes"))
		Expect(resp.VsRoutes[1].Matcher).To(Equal("/b"))
		Expect(resp.VsRoutes[1].MatchType).To(Equal("PREFIX"))
		Expect(len(resp.VsRoutes[1].RtRoutes)).To(Equal(3))
		rtB := resp.VsRoutes[1].RtRoutes
		Expect(rtB[0].GetRouteAction().GetSingle().GetUpstream().GetName()).To(Equal("baz-upstream"))
		Expect(rtB[0].Matcher).To(Equal("/b/3"))
		Expect(rtB[0].MatchType).To(Equal("REGEX"))
		Expect(rtB[1].GetRedirectAction().GetHostRedirect()).To(Equal("test.solo.io"))
		Expect(rtB[1].Matcher).To(Equal("/b/4"))
		Expect(rtB[1].MatchType).To(Equal("REGEX"))
		Expect(rtB[2].GetDelegateAction().GetRef().GetName()).To(Equal("c-routes"))
		Expect(rtB[2].Matcher).To(Equal("/b/c"))
		Expect(rtB[2].MatchType).To(Equal("PREFIX"))
		Expect(len(rtB[2].RtRoutes)).To(Equal(2))
		rtC := rtB[2].RtRoutes
		Expect(rtC[0].GetRouteAction().GetSingle().GetUpstream().GetName()).To(Equal("qux-upstream"))
		Expect(rtC[0].Matcher).To(Equal("/b/c/4"))
		Expect(rtC[0].MatchType).To(Equal("EXACT"))
		Expect(rtC[1].GetDirectResponseAction().GetBody()).To(Equal("Hello, world!"))
		Expect(rtC[1].Matcher).To(Equal("/b/c/567"))
		Expect(rtC[1].MatchType).To(Equal("EXACT"))
	})

})
