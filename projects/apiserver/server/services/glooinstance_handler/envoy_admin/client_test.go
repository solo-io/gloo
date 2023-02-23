package envoy_admin_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	envoy_admin_v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/jsonpb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/envoy_admin"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	fakerest "k8s.io/client-go/rest/fake"
)

var _ = Describe("envoy admin client", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can get clusters", func() {
		// clusters for proxy1
		clusters1 := envoy_admin_v3.Clusters{
			ClusterStatuses: []*envoy_admin_v3.ClusterStatus{
				{
					Name: "default-httpbin1-8000_gloo-system",
					HostStatuses: []*envoy_admin_v3.HostStatus{
						{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{
										Address: "1.2.3.4",
										PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
											PortValue: 8000,
										},
									},
								},
							},
							Weight: 2,
						},
						{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{
										Address: "5.6.7.8",
										PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
											PortValue: 9000,
										},
									},
								},
							},
							Weight: 1,
						},
					},
				},
				{
					Name: "default-httpbin2-8000_gloo-system",
					HostStatuses: []*envoy_admin_v3.HostStatus{
						{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{
										Address: "9.8.7.6",
										PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
											PortValue: 3000,
										},
									},
								},
							},
							Weight: 1,
						},
					},
				},
			},
		}
		// clusters for proxy2
		clusters2 := envoy_admin_v3.Clusters{
			ClusterStatuses: []*envoy_admin_v3.ClusterStatus{
				{
					Name: "default-petstore-8080_gloo-system",
					HostStatuses: []*envoy_admin_v3.HostStatus{
						{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{
										Address: "2.2.2.2",
										PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
											PortValue: 3456,
										},
									},
								},
							},
							Weight: 3,
						},
						{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{
										Address: "5.6.7.8",
										PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
											PortValue: 9000,
										},
									},
								},
							},
							Weight: 1,
						},
					},
				},
				{
					Name: "default-httpbin1-8000_gloo-system",
					HostStatuses: []*envoy_admin_v3.HostStatus{
						{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{
										Address: "6.7.8.9",
										PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
											PortValue: 4000,
										},
									},
								},
							},
							Weight: 1,
						},
					},
				},
			},
		}
		// marshal them into bytes so they can be returned by the mock request handler
		var marshaller jsonpb.Marshaler
		var res1, res2 bytes.Buffer
		err := marshaller.Marshal(&res1, &clusters1)
		Expect(err).NotTo(HaveOccurred())
		err = marshaller.Marshal(&res2, &clusters2)
		Expect(err).NotTo(HaveOccurred())

		// have the request handler return clusters1 for proxy1 and clusters2 for proxy2
		fakeReqHandler := func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.Path, "proxy1") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     defaultHeaders(),
					Body:       bytesBody(res1.Bytes()),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     defaultHeaders(),
				Body:       bytesBody(res2.Bytes()),
			}, nil
		}

		restClient := &fakerest.RESTClient{
			GroupVersion:         schema.GroupVersion{},
			VersionedAPIPath:     "/not/a/real/path",
			NegotiatedSerializer: serializer.CodecFactory{},
			Client:               fakerest.CreateHTTPClient(fakeReqHandler),
		}

		// create a gloo instance with the 2 proxies
		glooInstance := &rpc_edge_v1.GlooInstance{
			Metadata: &rpc_edge_v1.ObjectMeta{
				Name:      "gloo",
				Namespace: "gloo-system",
			},
			Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
				Cluster:      "cluster",
				IsEnterprise: true,
				ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
					Version:           "1.2.3",
					Namespace:         "gloo-system",
					WatchedNamespaces: []string{"ns1", "ns2", "gloo-system"},
				},
				Proxies: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{
					{
						Name:                          "proxy1",
						Namespace:                     "ns1",
						AvailableReplicas:             1,
						ReadConfigMulticlusterEnabled: true,
					},
					{
						Name:                          "proxy2",
						Namespace:                     "ns2",
						AvailableReplicas:             1,
						ReadConfigMulticlusterEnabled: true,
					},
				},
			},
		}

		// the upstream hosts should contain all the clusters from the 2 proxies
		envoyAdminClient := envoy_admin.NewEnvoyAdminClient()
		upstreamHosts, err := envoyAdminClient.GetHostsByUpstream(ctx, glooInstance, restClient)
		Expect(err).NotTo(HaveOccurred())

		expectedResults := map[string]*rpc_edge_v1.HostList{
			"gloo-system.default-httpbin1-8000": {
				Hosts: []*rpc_edge_v1.Host{
					{Address: "1.2.3.4", Port: 8000, Weight: 2, ProxyRef: &skv2_v1.ObjectRef{Name: "proxy1", Namespace: "ns1"}},
					{Address: "5.6.7.8", Port: 9000, Weight: 1, ProxyRef: &skv2_v1.ObjectRef{Name: "proxy1", Namespace: "ns1"}},
					{Address: "6.7.8.9", Port: 4000, Weight: 1, ProxyRef: &skv2_v1.ObjectRef{Name: "proxy2", Namespace: "ns2"}},
				},
			},
			"gloo-system.default-httpbin2-8000": {
				Hosts: []*rpc_edge_v1.Host{
					{Address: "9.8.7.6", Port: 3000, Weight: 1, ProxyRef: &skv2_v1.ObjectRef{Name: "proxy1", Namespace: "ns1"}},
				},
			},
			"gloo-system.default-petstore-8080": {
				Hosts: []*rpc_edge_v1.Host{
					{Address: "2.2.2.2", Port: 3456, Weight: 3, ProxyRef: &skv2_v1.ObjectRef{Name: "proxy2", Namespace: "ns2"}},
					{Address: "5.6.7.8", Port: 9000, Weight: 1, ProxyRef: &skv2_v1.ObjectRef{Name: "proxy2", Namespace: "ns2"}},
				},
			},
		}
		Expect(upstreamHosts).To(HaveLen(3))
		for k, v := range upstreamHosts {
			Expect(v).To(Equal(expectedResults[k]))
		}
	})
})

func bytesBody(bodyBytes []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(bodyBytes))
}

func defaultHeaders() http.Header {
	header := http.Header{}
	header.Set("Content-Type", runtime.ContentTypeJSON)
	return header
}
