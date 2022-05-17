package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	envoywaf "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/waf"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

// NOTE: To run waf e2e tests locally, specify the
// env var ENVOY_IMAGE_TAG=v1.x.x (your gloo ee version)
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
		customInterventionMessage = "It's a custom intervention message"
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

	var getProxyWaf = func(
		envoyPort uint32,
		upstream *core.ResourceRef,
		wafListenerSettings *waf.Settings,
		wafVhostSettings *waf.Settings,
		wafRouteSettings *waf.Settings,
	) *gloov1.Proxy {
		var vhosts []*gloov1.VirtualHost

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system.virt1",
			Domains: []string{"*"},
			Options: &gloov1.VirtualHostOptions{
				Waf: wafVhostSettings,
			},
			Routes: []*gloov1.Route{
				{
					Options: &gloov1.RouteOptions{
						Waf: wafRouteSettings,
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
										Upstream: upstream,
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
										Upstream: upstream,
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
			Metadata: &core.Metadata{
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
							Waf: wafListenerSettings,
						},
					},
				},
			}},
		}

		return p
	}

	var getProxyWafDisruptiveListener = func(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
		wafCfg := &waf.Settings{
			RuleSets:                  []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
			CustomInterventionMessage: customInterventionMessage,
		}
		return getProxyWaf(envoyPort, upstream, wafCfg, nil, nil)
	}

	var getProxyWafDisruptiveVhost = func(
		envoyPort uint32,
		upstream *core.ResourceRef,
		wafVhostSettings *waf.Settings,
	) *gloov1.Proxy {
		return getProxyWaf(envoyPort, upstream, nil, wafVhostSettings, nil)
	}

	var getProxyWafDisruptiveRoute = func(
		envoyPort uint32,
		upstream *core.ResourceRef,
		wafRouteSettings *waf.Settings,
	) *gloov1.Proxy {
		vhostSettings := &waf.Settings{
			Disabled: true,
		}
		return getProxyWaf(envoyPort, upstream, nil, vhostSettings, wafRouteSettings)
	}

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem})
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

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				})
			})

			It("will get rejected by waf", func() {
				var bodyStr string
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
						Header: map[string][]string{
							"user-agent": {"nikto"},
						},
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return 0, err
					}
					bodyStr = string(body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
				Expect(bodyStr).To(ContainSubstring(customInterventionMessage))
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
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})

		})

		Context("no body processing rules", func() {
			var (
				proxy *gloov1.Proxy
				route bool
				vhost bool
			)

			prep := func(request, bodypassthrough bool) {
				// make sure body rules are passed through
				rules := `
				# Turn rule engine on
				SecRuleEngine On
				SecResponseBodyAccess On
	 `
				wafCfg := &waf.Settings{
					CustomInterventionMessage: customInterventionMessage,
				}
				if request {
					rules += `
					SecRule REQUEST_BODY "@contains nikto" "deny,status:403,id:107,phase:2,msg:'blocked nikto scammer'"
					`
					if bodypassthrough {
						wafCfg.RequestHeadersOnly = true
					}
				} else {
					rules += `
					SecRule RESPONSE_BODY "@contains nikto" "deny,status:403,id:107,phase:4,msg:'blocked nikto scammer'"
					`
					if bodypassthrough {
						wafCfg.ResponseHeadersOnly = true
					}
				}

				ruleset := &envoywaf.RuleSet{
					RuleStr: rules,
				}
				wafCfg.RuleSets = []*envoywaf.RuleSet{ruleset}
				if vhost {
					proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), nil, wafCfg, nil)
				} else if route {
					proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), nil, nil, wafCfg)
				} else {
					proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg, nil, nil)
				}

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				})
			}

			EventuallyWithBody := func() gomega.AsyncAssertion {
				return EventuallyWithOffset(1, func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodPost,
						URL:    reqUrl,
						Body:   ioutil.NopCloser(bytes.NewBuffer([]byte("nikto"))),
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s")
			}
			Context("on listener", func() {
				BeforeEach(func() {
					route = false
					vhost = false
				})

				It("will reject request body", func() {
					prep(true, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})
				It("will reject response body", func() {
					prep(false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})

				It("will NOT reject request body", func() {
					prep(true, true)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})

				It("will NOT reject response body", func() {
					prep(false, true)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})
			})

			Context("on vhost", func() {
				BeforeEach(func() {
					route = false
					vhost = true
				})

				It("will reject request body", func() {
					prep(true, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})
				It("will reject response body", func() {
					prep(false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})

				It("will NOT reject request body", func() {
					prep(true, true)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})

				It("will NOT reject response body", func() {
					prep(false, true)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})
			})

			Context("on route", func() {
				BeforeEach(func() {
					route = true
					vhost = false
				})

				It("will reject request body", func() {
					prep(true, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})
				It("will reject response body", func() {
					prep(false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})

				It("will NOT reject request body", func() {
					prep(true, true)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})

				It("will NOT reject response body", func() {
					prep(false, true)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})
			})

		})

		Context("vhost rules", func() {
			var (
				proxy *gloov1.Proxy
			)

			BeforeEach(func() {
				wafCfg := &waf.Settings{
					RuleSets:                  []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					CustomInterventionMessage: customInterventionMessage,
				}
				proxy = getProxyWafDisruptiveVhost(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg)

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				})
			})

			It("will get rejected by waf", func() {
				var bodyStr string
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
						Header: map[string][]string{
							"user-agent": {"nikto"},
						},
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return 0, err
					}
					bodyStr = string(body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
				Expect(bodyStr).To(ContainSubstring(customInterventionMessage))
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
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})

		})

		Context("route rules", func() {
			var (
				proxy *gloov1.Proxy
			)

			BeforeEach(func() {
				wafCfg := &waf.Settings{
					RuleSets:                  []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					CustomInterventionMessage: customInterventionMessage,
				}
				proxy = getProxyWafDisruptiveRoute(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg)

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				})
			})

			It("will get rejected by waf", func() {
				var bodyStr string
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
						Header: map[string][]string{
							"user-agent": {"nikto"},
						},
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return 0, err
					}
					bodyStr = string(body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
				Expect(bodyStr).To(ContainSubstring(customInterventionMessage))
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
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
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
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})

		})

		Context("audit logs", func() {
			var (
				proxy         *gloov1.Proxy
				tmpFileFSName string
				tmpFileDMName string
			)
			BeforeEach(func() {
				tmpFile, err := ioutil.TempFile("", "envoy-access-fs-log-*.txt")
				Expect(err).NotTo(HaveOccurred())
				tmpFileFSName, err = filepath.Abs(tmpFile.Name())
				tmpFile.Close()
				Expect(err).NotTo(HaveOccurred())
				tmpFile, err = ioutil.TempFile("", "envoy-access-dm-log-*.txt")
				Expect(err).NotTo(HaveOccurred())
				tmpFileDMName, err = filepath.Abs(tmpFile.Name())
				tmpFile.Close()
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				if tmpFileFSName != "" {
					os.Remove(tmpFileFSName)
				}
				if tmpFileDMName != "" {
					os.Remove(tmpFileDMName)
				}
			})

			getAccessFSLog := func() string {
				b, err := ioutil.ReadFile(tmpFileFSName)
				Expect(err).NotTo(HaveOccurred())
				return string(b)
			}
			getAccessDMLog := func() string {
				b, err := ioutil.ReadFile(tmpFileDMName)
				Expect(err).NotTo(HaveOccurred())
				return string(b)
			}

			startProxy := func(wafListenerSettings, wafVhostSettings, wafRouteSettings *waf.Settings) {

				By("tmp file " + tmpFileFSName + " " + tmpFileDMName)

				proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafListenerSettings, wafVhostSettings, wafRouteSettings)
				proxy.Listeners[0].Options = &gloov1.ListenerOptions{
					AccessLoggingService: &als.AccessLoggingService{
						AccessLog: []*als.AccessLog{
							{
								OutputDestination: &als.AccessLog_FileSink{
									FileSink: &als.FileSink{
										Path: tmpFileFSName,
										OutputFormat: &als.FileSink_StringFormat{
											StringFormat: "%FILTER_STATE(io.solo.modsecurity.audit_log)%\n",
										},
									},
								},
							}, {
								OutputDestination: &als.AccessLog_FileSink{
									FileSink: &als.FileSink{
										Path: tmpFileDMName,
										OutputFormat: &als.FileSink_StringFormat{
											StringFormat: "%DYNAMIC_METADATA(io.solo.filters.http.modsecurity:audit_log)%\n",
										},
									},
								},
							},
						},
					},
				}
				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				})
			}
			makeBadRequest := func() {
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
						Header: map[string][]string{
							"user-agent": {"nikto"},
						},
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "10s", "0.5s").Should(Equal(http.StatusForbidden))
			}
			makeGoodRequest := func() {
				Eventually(func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodGet,
						URL:    reqUrl,
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "10s", "0.5s").Should(Equal(http.StatusOK))
			}

			It("auditlog listener filter state", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "5s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "5s", "1s").Should(Equal("-\n"))
			})

			It("auditlog listener dynamic meta", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_DYNAMIC_METADATA,
					},
				}, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessDMLog, "5s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessFSLog, "5s", "1s").Should(Equal("-\n"))
			})
			It("auditlog listener fs - logs relevant", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_RELEVANT_ONLY,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "5s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "5s", "1s").Should(Equal("-\n"))
			})
			It("auditlog listener fs - not log not relevant", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_RELEVANT_ONLY,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil)
				makeGoodRequest()
				// check the logs
				Eventually(getAccessFSLog, "5s", "1s").Should(Equal("-\n"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "5s", "1s").Should(Equal("-\n"))
			})
			It("auditlog listener dm - not log not relevant", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_RELEVANT_ONLY,
						Location: envoywaf.AuditLogging_DYNAMIC_METADATA,
					},
				}, nil, nil)
				makeGoodRequest()
				// nothing written to dm log
				Eventually(getAccessDMLog, "5s", "1s").Should(Equal("-\n"))
				Eventually(getAccessFSLog, "5s", "1s").Should(Equal("-\n"))
			})
			It("auditlog listener dm - not log relevant if disabled", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
				}, nil, nil)
				makeBadRequest()
				// nothing written to any logs
				Eventually(getAccessDMLog, "5s", "1s").Should(Equal("-\n"))
				Eventually(getAccessFSLog, "5s", "1s").Should(Equal("-\n"))
			})

			It("auditlog vhost filter state", func() {
				startProxy(nil, &waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "5s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "5s", "1s").Should(Equal("-\n"))
			})

			It("auditlog route filter state", func() {
				startProxy(nil, nil, &waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				})
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "5s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "5s", "1s").Should(Equal("-\n"))
			})

		})
	})
})
