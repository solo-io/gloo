package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/fgrosse/zaptest"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/pkg/utils"
	envoywaf "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/waf"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/waf"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	wafplugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/waf"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("waf", func() {
	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
	)

	const (
		rulesTemplate = `
            # Turn rule engine on
            SecRuleEngine On
            SecRule %s:User-Agent "nikto" "%s,id:107,%s,msg:'blocked nikto scammer'"
 `
	)

	var getRulesTemplate = func(deny, request, phase1 bool) *envoywaf.RuleSet {
		denialString := "deny,status:403"
		if deny == false {
			denialString = "redirect:'http://example.com'"
		}
		requestString := "REQUEST_HEADERS"
		if request == false {
			requestString = "RESPONSE_HEADERS"
		}
		phaseString := "phase:1"
		if phase1 == false {
			phaseString = "phase:3"
		}
		return &envoywaf.RuleSet{
			RuleStr: fmt.Sprintf(rulesTemplate, requestString, denialString, phaseString),
		}
	}

	var getProxyWaf = func(envoyPort uint32, upstream core.ResourceRef, wafSettings *gloov1.Extensions, vhostextensions *gloov1.Extensions, routeExtensions *gloov1.Extensions) *gloov1.Proxy {
		var vhosts []*gloov1.VirtualHost

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system.virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &gloov1.VirtualHostPlugins{
				Extensions: vhostextensions,
			},
			Routes: []*gloov1.Route{
				{
					RoutePlugins: &gloov1.RoutePlugins{
						Extensions: routeExtensions,
					},
					Matcher: &gloov1.Matcher{
						PathSpecifier: &gloov1.Matcher_Prefix{
							Prefix: "/hello",
						},
					},
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
					Matcher: &gloov1.Matcher{
						PathSpecifier: &gloov1.Matcher_Prefix{
							Prefix: "/world",
						},
					},
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
						ListenerPlugins: &gloov1.HttpListenerPlugins{
							Extensions: wafSettings,
						},
					},
				},
			}},
		}

		return p
	}

	var getProxyWafDisruptiveListener = func(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
		wafCfg := &waf.Settings{
			RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
		}
		cfg, err := envoyutil.MessageToStruct(wafCfg)
		Expect(err).NotTo(HaveOccurred())
		extensions := &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				wafplugin.ExtensionName: cfg,
			},
		}
		return getProxyWaf(envoyPort, upstream, extensions, nil, nil)
	}

	var getProxyWafDisruptiveVhost = func(envoyPort uint32, upstream core.ResourceRef, settings *waf.VhostSettings) *gloov1.Proxy {
		cfg, err := envoyutil.MessageToStruct(settings)
		Expect(err).NotTo(HaveOccurred())
		extensions := &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				wafplugin.ExtensionName: cfg,
			},
		}

		return getProxyWaf(envoyPort, upstream, nil, extensions, nil)
	}

	var getProxyWafDisruptiveRoute = func(envoyPort uint32, upstream core.ResourceRef, settings *waf.RouteSettings) *gloov1.Proxy {
		cfg, err := envoyutil.MessageToStruct(settings)
		Expect(err).NotTo(HaveOccurred())
		extensions := &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				wafplugin.ExtensionName: cfg,
			},
		}
		vhostSettings := &waf.VhostSettings{
			Disabled: true,
			Settings: nil,
		}
		cfg, err = envoyutil.MessageToStruct(vhostSettings)
		Expect(err).NotTo(HaveOccurred())
		vhostExtensions := &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				wafplugin.ExtensionName: cfg,
			},
		}
		return getProxyWaf(envoyPort, upstream, nil, vhostExtensions, extensions)
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

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil)
	})

	AfterEach(func() {
		cancel()
	})
	Context("With envoy", func() {
		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)
		)

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

			var opts clients.WriteOpts
			up := testUpstream.Upstream
			_, err = testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())

		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		Context("listener rules", func() {
			var (
				proxy *gloov1.Proxy
			)

			BeforeEach(func() {
				proxy = getProxyWafDisruptiveListener(envoyPort, testUpstream.Upstream.Metadata.Ref())

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
			})

			It("will get rejected by waf", func() {
				var resp *http.Response
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err = client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
						Header: map[string][]string{
							"user-agent": {"nikto"},
						},
					})
					if resp == nil {
						return 0, nil
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
				bodyStr, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(bodyStr).To(ContainSubstring("ModSecurity"))
			})

			It("will not get rejected by waf", func() {
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
					})
					if resp == nil {
						return 0, nil
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})

		})

		Context("vhost rules", func() {
			var (
				proxy *gloov1.Proxy
			)

			BeforeEach(func() {
				wafCfg := &waf.VhostSettings{
					Settings: &waf.Settings{
						RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					},
				}
				proxy = getProxyWafDisruptiveVhost(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg)

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
			})

			It("will get rejected by waf", func() {
				var resp *http.Response
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err = client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
						Header: map[string][]string{
							"user-agent": {"nikto"},
						},
					})
					if resp == nil {
						return 0, nil
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
				bodyStr, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(bodyStr).To(ContainSubstring("ModSecurity"))
			})

			It("will not get rejected by waf", func() {
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
					})
					if resp == nil {
						return 0, nil
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})

		})

		Context("route rules", func() {
			var (
				proxy *gloov1.Proxy
			)

			BeforeEach(func() {
				wafCfg := &waf.RouteSettings{
					Settings: &waf.Settings{
						RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					},
				}
				proxy = getProxyWafDisruptiveRoute(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg)

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
			})

			It("will get rejected by waf", func() {
				var resp *http.Response
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err = client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
						Header: map[string][]string{
							"user-agent": {"nikto"},
						},
					})
					if resp == nil {
						return 0, nil
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
				bodyStr, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(bodyStr).To(ContainSubstring("ModSecurity"))
			})

			It("will not get rejected by waf", func() {
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
					})
					if resp == nil {
						return 0, nil
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})

			It("will not get rejected by waf since it's on a different route", func() {
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/world/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
					})
					if resp == nil {
						return 0, nil
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})

		})
	})

})
