package e2e_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/solo-io/gloo/test/services"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("grpcweb", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		t := services.RunGateway(ctx, true)
		testClients = t
		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())
		err = envoyInstance.Run(t.GlooPort)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		if envoyInstance != nil {
			envoyInstance.Clean()
		}
		cancel()
	})

	It("should translate grpc web", func() {

		tu := v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		var opts clients.WriteOpts
		up := tu.Upstream
		_, err := testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())

		proxycli := testClients.ProxyClient
		envoyPort := uint32(8080)
		proxy := &gloov1.Proxy{
			Metadata: core.Metadata{
				Name:      "proxy",
				Namespace: "default",
			},
			Listeners: []*gloov1.Listener{{
				Name:        "listener",
				BindAddress: "127.0.0.1",
				BindPort:    envoyPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: []*gloov1.VirtualHost{{
							Name:    "virt1",
							Domains: []string{"*"},
							Routes: []*gloov1.Route{{
								Matcher: &gloov1.Matcher{
									PathSpecifier: &gloov1.Matcher_Prefix{
										Prefix: "/",
									},
								},
								Action: &gloov1.Route_RouteAction{
									RouteAction: &gloov1.RouteAction{
										Destination: &gloov1.RouteAction_Single{
											Single: &gloov1.Destination{
												Upstream: up.Metadata.Ref(),
											},
										},
									},
								},
							}},
						}},
					},
				},
			}},
		}

		_, err = proxycli.Write(proxy, opts)
		Expect(err).NotTo(HaveOccurred())

		body := []byte("solo.io test")

		Eventually(func() error {
			// send a request with a body
			var buf bytes.Buffer
			buf.Write(body)

			res, err := http.Post(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), "application/octet-stream", &buf)
			if err != nil {
				return err
			}
			if res.StatusCode != http.StatusOK {
				return errors.New(fmt.Sprintf("%v is not OK", res.StatusCode))
			}
			return nil
		}, "5s", ".5s").Should(BeNil())

		Eventually(tu.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
			"Method": Equal("POST"),
			"Body":   Equal(body),
		}))))

	})
})
