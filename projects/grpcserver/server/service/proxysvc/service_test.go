package proxysvc_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	. "github.com/solo-io/solo-projects/projects/grpcserver/server/service/internal/testutils"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/proxysvc"
	"google.golang.org/grpc"
)

var (
	grpcServer  *grpc.Server
	conn        *grpc.ClientConn
	apiserver   v1.ProxyApiServer
	client      v1.ProxyApiClient
	mockCtrl    *gomock.Controller
	proxyClient *mocks.MockProxyClient
	rawGetter   *mock_rawgetter.MockRawGetter
	testErr     = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	getRaw := func(proxy *gloov1.Proxy) *v1.Raw {
		return &v1.Raw{FileName: proxy.GetMetadata().Name}
	}

	getProxyDetails := func(proxy *gloov1.Proxy) *v1.ProxyDetails {
		return &v1.ProxyDetails{
			Proxy: proxy,
			Raw:   getRaw(proxy),
		}
	}

	getProxyDetailsList := func(proxies ...*gloov1.Proxy) []*v1.ProxyDetails {
		list := make([]*v1.ProxyDetails, 0, len(proxies))
		for _, p := range proxies {
			list = append(list, &v1.ProxyDetails{
				Proxy: p,
				Raw:   getRaw(p),
			})
		}
		return list
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		proxyClient = mocks.NewMockProxyClient(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		apiserver = proxysvc.NewProxyGrpcService(context.TODO(), proxyClient, rawGetter)

		grpcServer, conn = MustRunGrpcServer(func(s *grpc.Server) { v1.RegisterProxyApiServer(s, apiserver) })
		client = v1.NewProxyApiClient(conn)
	})

	AfterEach(func() {
		grpcServer.Stop()
		mockCtrl.Finish()
	})

	Describe("GetProxy", func() {
		It("works when the proxy client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			proxy := &gloov1.Proxy{
				Metadata: metadata,
			}

			proxyClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(proxy, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), proxy, gloov1.ProxyCrd).
				Return(getRaw(proxy))

			request := &v1.GetProxyRequest{Ref: &ref}
			actual, err := client.GetProxy(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetProxyResponse{ProxyDetails: getProxyDetails(proxy)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the proxy client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			proxyClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetProxyRequest{Ref: &ref}
			_, err := client.GetProxy(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := proxysvc.FailedToGetProxyError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListProxies", func() {
		It("works when the proxy client works", func() {
			ns1, ns2 := "one", "two"
			proxy1 := &gloov1.Proxy{
				Metadata: core.Metadata{Namespace: ns1},
			}
			proxy2 := &gloov1.Proxy{
				Metadata: core.Metadata{Namespace: ns2},
			}

			proxyClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Proxy{proxy1}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), proxy1, gloov1.ProxyCrd).
				Return(getRaw(proxy1))
			proxyClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Proxy{proxy2}, nil)
			rawGetter.EXPECT().
				GetRaw(context.Background(), proxy2, gloov1.ProxyCrd).
				Return(getRaw(proxy2))

			request := &v1.ListProxiesRequest{Namespaces: []string{ns1, ns2}}
			actual, err := client.ListProxies(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListProxiesResponse{ProxyDetails: getProxyDetailsList(proxy1, proxy2)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the proxy client errors", func() {
			ns := "ns"

			proxyClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.ListProxiesRequest{Namespaces: []string{ns}}
			_, err := client.ListProxies(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := proxysvc.FailedToListProxiesError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
