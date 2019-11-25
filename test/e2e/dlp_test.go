package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/pkg/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("dlp", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
	)

	var getProxyDlp = func(envoyPort uint32, upstream core.ResourceRef, dlpListenerSettings *dlp.FilterConfig,
		dlpVhostSettings *dlp.Config, dlpRouteSettings *dlp.Config) *gloov1.Proxy {

		var vhosts []*gloov1.VirtualHost

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system.virt1",
			Domains: []string{"*"},
			Options: &gloov1.VirtualHostOptions{
				Dlp: dlpVhostSettings,
			},
			Routes: []*gloov1.Route{
				{
					Options: &gloov1.RouteOptions{
						Dlp: dlpRouteSettings,
					},
					Matchers: []*matchers.Matcher{{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: "/hello",
						},
					}},
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(upstream),
									},
								},
							},
						},
					},
				},
				{
					Matchers: []*matchers.Matcher{{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: "/world",
						},
					}},
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(upstream),
									},
								},
							},
						},
					},
				},
			},
		}

		vhosts = append(vhosts, vhost)

		p := &gloov1.Proxy{
			Metadata: core.Metadata{
				Name:      "proxy",
				Namespace: "default",
			},
			Listeners: []*gloov1.Listener{{
				Name:        "listener",
				BindAddress: "0.0.0.0",
				BindPort:    envoyPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: vhosts,
						Options: &gloov1.HttpListenerOptions{
							Dlp: dlpListenerSettings,
						},
					},
				},
			}},
		}

		return p
	}

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, nil)
	})

	AfterEach(func() {
		cancel()
	})
	Context("With envoy", func() {
		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)

			proxy *gloov1.Proxy
		)

		var testRequest = func(result string) {
			var resp *http.Response
			Eventually(func() (int, error) {
				client := http.DefaultClient
				reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
				Expect(err).NotTo(HaveOccurred())
				resp, err = client.Do(&http.Request{
					Method: http.MethodGet,
					URL:    reqUrl,
				})
				if resp == nil {
					return 0, nil
				}
				return resp.StatusCode, nil
			}, "5s", "0.5s").Should(Equal(http.StatusOK))
			bodyStr, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(bodyStr).To(ContainSubstring(result))
		}

		var configureProxy = func() {
			Expect(proxy).NotTo(BeNil())
			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() (core.Status, error) {
				proxy, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				if err != nil {
					return core.Status{}, err
				}
				return proxy.Status, nil
			}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
				"Reason": BeEmpty(),
				"State":  Equal(core.Status_Accepted),
			}))
		}

		BeforeEach(func() {
			proxy = nil
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "hello")
			_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		Context("listener rules", func() {

			var configureListenerProxy = func(actions []*dlp.Action, matcher *matchers.Matcher) {
				if matcher == nil {
					matcher = &matchers.Matcher{
						PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"},
					}
				}
				dlpCfg := &dlp.FilterConfig{
					DlpRules: []*dlp.DlpRule{
						{
							Matcher: matcher,
							Actions: actions,
						},
					},
				}
				proxy = getProxyDlp(envoyPort, testUpstream.Upstream.Metadata.Ref(), dlpCfg, nil, nil)
				configureProxy()
			}

			It("simple dlp action", func() {
				configureListenerProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:  "test",
						Regex: []string{"hello"},
					},
				}}, nil)
				testRequest("XXXXo")
			})

			It("simple shadow action", func() {
				configureListenerProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:  "test",
						Regex: []string{"hello"},
					},
					Shadow: true,
				}}, nil)
				testRequest("hello")
			})

			It("more complex action", func() {
				configureListenerProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:     "test",
						Regex:    []string{"hello"},
						MaskChar: "Y",
						Percent: &envoy_type.Percent{
							Value: 60,
						},
					},
				}}, nil)
				testRequest("YYYlo")
			})

			It("no transform on route mismatch", func() {
				configureListenerProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:  "test",
						Regex: []string{"hello"},
					},
				}}, &matchers.Matcher{
					PathSpecifier: &matchers.Matcher_Exact{Exact: "/will/not/match"},
				})
				testRequest("hello")
			})

			Context("With SSN", func() {

				Context("Matches standalone SSN", func() {
					JustBeforeEach(func() {
						testUpstream = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "123-45-6789")
						_, err := testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())
					})

					It("matches if SSN provided alone", func() {
						configureListenerProxy([]*dlp.Action{{
							ActionType: dlp.Action_SSN,
						}}, nil)
						testRequest("XXX-XX-X789")
					})
				})

				Context("Matches SSN in JSON", func() {
					JustBeforeEach(func() {
						testUpstream = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "\"ssn\":\"123-45-6789\"")
						_, err := testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())
					})

					It("does not match boundary characters with standard regex", func() {
						configureListenerProxy([]*dlp.Action{{
							ActionType: dlp.Action_SSN,
						}}, nil)
						testRequest("\"ssn\":\"XXX-XX-X789\"")
					})
				})

			})
		})

		Context("vhost rules", func() {

			var configureDlpForProxy = func(actions []*dlp.Action) {

				dlpCfg := &dlp.Config{
					Actions: actions,
				}
				proxy = getProxyDlp(envoyPort, testUpstream.Upstream.Metadata.Ref(), nil, dlpCfg, nil)
				configureProxy()
			}

			It("will get mask the response by waf", func() {
				configureDlpForProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:  "test",
						Regex: []string{"hello"},
					},
				}})
				testRequest("XXXXo")
			})

			It("simple shadow action", func() {
				configureDlpForProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:  "test",
						Regex: []string{"hello"},
					},
					Shadow: true,
				}})
				testRequest("hello")
			})

			It("more complex action", func() {
				configureDlpForProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:     "test",
						Regex:    []string{"hello"},
						MaskChar: "Y",
						Percent: &envoy_type.Percent{
							Value: 60,
						},
					},
				}})
				testRequest("YYYlo")
			})

		})

		Context("route rules", func() {

			var configureDlpForProxy = func(actions []*dlp.Action) {

				dlpCfg := &dlp.Config{
					Actions: actions,
				}
				proxy = getProxyDlp(envoyPort, testUpstream.Upstream.Metadata.Ref(), nil, nil, dlpCfg)
				configureProxy()
			}

			It("will get mask the response by waf", func() {
				configureDlpForProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:  "test",
						Regex: []string{"hello"},
					},
				}})
				testRequest("XXXXo")
			})

			It("simple shadow action", func() {
				configureDlpForProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:  "test",
						Regex: []string{"hello"},
					},
					Shadow: true,
				}})
				testRequest("hello")
			})

			It("more complex action", func() {
				configureDlpForProxy([]*dlp.Action{{
					ActionType: dlp.Action_CUSTOM,
					CustomAction: &dlp.CustomAction{
						Name:     "test",
						Regex:    []string{"hello"},
						MaskChar: "Y",
						Percent: &envoy_type.Percent{
							Value: 60,
						},
					},
				}})
				testRequest("YYYlo")
			})

		})
	})
})
