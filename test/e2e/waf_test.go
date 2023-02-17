package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	envoywaf "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/waf"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	gloossl "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/helpers"
	gloohelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/contextutils"
	envoy_type "github.com/solo-io/solo-kit/pkg/api/external/envoy/type"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

// NOTE: To run waf e2e tests locally, specify the
// env var ENVOY_IMAGE_TAG=v1.x.x (your gloo ee version)
var _ = Describe("waf", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		rulesTemplate = `
			# Turn rule engine on
			SecRuleEngine On
			SecAuditLogFormat %s
			SecRule %s:User-Agent "nikto" "%s,id:107,%s,msg:'blocked nikto'"
`
		rulesTemplate1 = `SecRule REQUEST_HEADERS:User-Agent "scammer" "deny,status:403,id:108,phase:1,msg:'blocked scammer'"`
		ignoredRule    = `SecRule REQUEST_HEADERS:User-Agent "skippedRule" "deny,status:403,id:109,phase:1,msg:'blocked skippedRule'"`

		configMapRuleSets = []*waf.RuleSetFromConfigMap{
			{
				ConfigMapRef: &core.ResourceRef{
					Namespace: "gloo-system",
					Name:      "configmapname",
				},
				DataMapKeys: []string{
					"key1", //order in artifact is key2, then key1 - validate order of passed keys is respected
				},
			},
		}
	)

	const (
		customInterventionMessage = "It's a custom intervention message"
	)

	var getRulesTemplate = func(deny, request, phase1, jsonLog bool) *envoywaf.RuleSet {
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
		logFormatString := "Native"
		if jsonLog {
			logFormatString = "JSON"
		}
		return &envoywaf.RuleSet{
			RuleStr: fmt.Sprintf(rulesTemplate, logFormatString, requestString, denialString, phaseString),
		}
	}

	var getProxyWaf = func(
		envoyPort uint32,
		upstream *core.ResourceRef,
		wafListenerSettings *waf.Settings,
		wafVhostSettings *waf.Settings,
		dlpVhostSettings *dlp.Config,
		dlpRouteSettings *dlp.Config,
		wafRouteSettings *waf.Settings,
		dlpListenerSettings *dlp.FilterConfig,
	) *gloov1.Proxy {
		var vhosts []*gloov1.VirtualHost

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system.virt1",
			Domains: []string{"*"},
			Options: &gloov1.VirtualHostOptions{
				Waf: wafVhostSettings,
				Dlp: dlpVhostSettings,
			},
			Routes: []*gloov1.Route{
				{
					Options: &gloov1.RouteOptions{
						Waf: wafRouteSettings,
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
				BindAddress: net.IPv4zero.String(),
				BindPort:    envoyPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: vhosts,
						Options: &gloov1.HttpListenerOptions{
							Waf: wafListenerSettings,
							Dlp: dlpListenerSettings,
						},
					},
				},
			}},
		}

		return p
	}

	var getProxyWafDisruptiveListener = func(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
		wafCfg := &waf.Settings{
			RuleSets:                  []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
			ConfigMapRuleSets:         configMapRuleSets,
			CustomInterventionMessage: customInterventionMessage,
		}
		return getProxyWaf(envoyPort, upstream, wafCfg, nil, nil, nil, nil, nil)
	}

	var getProxyWafDisruptiveVhost = func(
		envoyPort uint32,
		upstream *core.ResourceRef,
		wafVhostSettings *waf.Settings,
	) *gloov1.Proxy {
		return getProxyWaf(envoyPort, upstream, nil, wafVhostSettings, nil, nil, nil, nil)
	}

	var getProxyWafDisruptiveRoute = func(
		envoyPort uint32,
		upstream *core.ResourceRef,
		wafRouteSettings *waf.Settings,
	) *gloov1.Proxy {
		vhostSettings := &waf.Settings{
			Disabled: true,
		}
		return getProxyWaf(envoyPort, upstream, nil, vhostSettings, nil, nil, wafRouteSettings, nil)
	}

	var expectRequestBlockedByWaf = func(useragent string, envoyPort uint32) {
		var bodyStr string
		Eventually(func() (int, error) {
			client := http.DefaultClient
			reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
			Expect(err).NotTo(HaveOccurred())
			resp, err := client.Do(&http.Request{
				Method: http.MethodGet,
				URL:    reqUrl,
				Header: map[string][]string{
					"user-agent": {useragent},
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
		}, "10s", "0.5s").Should(Equal(http.StatusForbidden))
		Expect(bodyStr).To(ContainSubstring(customInterventionMessage))
	}

	var expectRequestAllowedByWaf = func(useragent string, envoyPort uint32) {
		Eventually(func() (int, error) {
			client := http.DefaultClient
			reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
			Expect(err).NotTo(HaveOccurred())
			var header map[string][]string
			if len(useragent) > 0 {
				header = map[string][]string{
					"user-agent": {useragent},
				}
			}
			resp, err := client.Do(&http.Request{
				Method: http.MethodGet,
				URL:    reqUrl,
				Header: header,
			})
			if err != nil {
				return 0, err
			}
			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body)
			return resp.StatusCode, nil
		}, "10s", "0.5s").Should(Equal(http.StatusOK))
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

			configmap := &gloov1.Artifact{
				Metadata: &core.Metadata{
					Namespace: "gloo-system",
					Name:      "configmapname",
				},
				Data: map[string]string{
					"key1": rulesTemplate1, //the only key passed
					"key2": ignoredRule,
				},
			}
			_, err = testClients.ArtifactClient.Write(configmap, clients.WriteOpts{})
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := testClients.ArtifactClient.Delete("gloo-system", "configmapname", clients.DeleteOpts{})
			Expect(err).ToNot(HaveOccurred())
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		Context("no body processing rules", func() {
			var (
				proxy *gloov1.Proxy
				route bool
				vhost bool
			)

			prep := func(request, bodypassthrough, configureTls bool) {
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
					SecRule REQUEST_BODY "@contains nikto" "deny,status:403,id:107,phase:2,msg:'blocked nikto'"
					`
					if bodypassthrough {
						wafCfg.RequestHeadersOnly = true
					}
				} else {
					rules += `
					SecRule RESPONSE_BODY "@contains nikto" "deny,status:403,id:107,phase:4,msg:'blocked nikto'"
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
					proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), nil, wafCfg, nil, nil, nil, nil)
				} else if route {
					proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), nil, nil, nil, nil, wafCfg, nil)
				} else {
					proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg, nil, nil, nil, nil, nil)
				}

				if configureTls {
					// write TLS secret
					secret := &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      "tls-secret",
							Namespace: "gloo-system",
						},
						Kind: &gloov1.Secret_Tls{
							Tls: &gloov1.TlsSecret{
								CertChain:  gloohelpers.Certificate(),
								PrivateKey: gloohelpers.PrivateKey(),
							},
						},
					}
					_, err := testClients.SecretClient.Write(secret, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					// configure TLS on the proxy
					proxy.Listeners[0].SslConfigurations = []*gloossl.SslConfig{{
						SslSecrets: &gloossl.SslConfig_SecretRef{
							SecretRef: secret.GetMetadata().Ref(),
						},
					}}
				}

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				}, "10s", "0.5s")
			}

			EventuallyWithBody := func() gomega.AsyncAssertion {
				return EventuallyWithOffset(1, func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodPost,
						URL:    reqUrl,
						Body:   io.NopCloser(bytes.NewBuffer([]byte("nikto"))),
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "10s", "0.5s")
			}

			EventuallyWithBodyHttps := func() gomega.AsyncAssertion {
				return EventuallyWithOffset(1, func() (int, error) {
					client := http.DefaultClient
					reqUrl, err := url.Parse(fmt.Sprintf("https://%s:%d/hello/1", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(&http.Request{
						Method: http.MethodPost,
						URL:    reqUrl,
						Body:   io.NopCloser(bytes.NewBuffer([]byte("nikto"))),
					})
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "10s", "0.5s")
			}

			EventuallyConnectRequest := func() gomega.AsyncAssertion {
				return EventuallyWithOffset(1, func() (int, error) {
					// run curl command and check output
					cmd := exec.Command("curl", "--insecure", "-X", "CONNECT", "--http1.1", "--request-target", fmt.Sprintf("localhost:%d", envoyPort), fmt.Sprintf("https://localhost:%d/hello/1", envoyPort), "-w", "%{http_code}")
					out, err := cmd.CombinedOutput()
					if err != nil {
						return 0, err
					}

					// get the status code from the last line of the output
					statusCodeString := strings.Split(string(out), "\n")[len(strings.Split(string(out), "\n"))-1]
					statusCode, err := strconv.Atoi(statusCodeString)
					if err != nil {
						return 0, err
					}
					return statusCode, nil
				}, "3s", "0.5s")
			}

			Context("on listener", func() {
				BeforeEach(func() {
					route = false
					vhost = false
				})

				It("will reject request body", func() {
					prep(true, false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})
				It("will reject response body", func() {
					prep(false, false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})

				It("will NOT reject request body", func() {
					prep(true, true, false)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})

				It("will handle CONNECT request", func() {
					prep(true, true, true)
					// validate TLS setup (send request to https port)
					EventuallyWithBodyHttps().Should(Equal(http.StatusOK))

					// we don't support connect requests in this service, but we should not crash
					EventuallyConnectRequest().Should(Equal(http.StatusNotFound))
				})

				It("will NOT reject response body", func() {
					prep(false, true, false)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})
			})

			Context("on vhost", func() {
				BeforeEach(func() {
					route = false
					vhost = true
				})

				It("will reject request body", func() {
					prep(true, false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})
				It("will reject response body", func() {
					prep(false, false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})

				It("will NOT reject request body", func() {
					prep(true, true, false)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})

				It("will NOT reject response body", func() {
					prep(false, true, false)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})
			})

			Context("on route", func() {
				BeforeEach(func() {
					route = true
					vhost = false
				})

				It("will reject request body", func() {
					prep(true, false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})
				It("will reject response body", func() {
					prep(false, false, false)
					EventuallyWithBody().Should(Equal(http.StatusForbidden))
				})

				It("will NOT reject request body", func() {
					prep(true, true, false)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})

				It("will NOT reject response body", func() {
					prep(false, true, false)
					EventuallyWithBody().Should(Equal(http.StatusOK))
				})
			})

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
				}, "10s", "0.5s")
			})

			It("will get rejected by waf", func() {
				expectRequestBlockedByWaf("nikto", envoyPort)
				expectRequestBlockedByWaf("scammer", envoyPort)
			})

			It("will not get rejected by waf", func() {
				expectRequestAllowedByWaf("", envoyPort)
				expectRequestAllowedByWaf("skippedRule", envoyPort)
			})

		})

		Context("vhost rules", func() {
			var (
				proxy *gloov1.Proxy
			)

			BeforeEach(func() {
				wafCfg := &waf.Settings{
					RuleSets:                  []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					ConfigMapRuleSets:         configMapRuleSets,
					CustomInterventionMessage: customInterventionMessage,
				}
				proxy = getProxyWafDisruptiveVhost(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg)

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				}, "10s", "0.5s")
			})

			It("will get rejected by waf", func() {
				expectRequestBlockedByWaf("nikto", envoyPort)
				expectRequestBlockedByWaf("scammer", envoyPort)
			})

			It("will not get rejected by waf", func() {
				expectRequestAllowedByWaf("", envoyPort)
				expectRequestAllowedByWaf("skippedRule", envoyPort)
			})

		})

		Context("route rules", func() {
			var (
				proxy *gloov1.Proxy
			)

			BeforeEach(func() {
				wafCfg := &waf.Settings{
					RuleSets:                  []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					CustomInterventionMessage: customInterventionMessage,
					ConfigMapRuleSets:         configMapRuleSets,
				}
				proxy = getProxyWafDisruptiveRoute(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafCfg)

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				}, "10s", "0.5s")
			})

			It("will get rejected by waf", func() {
				expectRequestBlockedByWaf("nikto", envoyPort)
				expectRequestBlockedByWaf("scammer", envoyPort)
			})

			It("will not get rejected by waf", func() {
				expectRequestAllowedByWaf("", envoyPort)
				expectRequestAllowedByWaf("skippedRule", envoyPort)
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
				}, "10s", "0.5s").Should(Equal(http.StatusOK))
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

			startProxy := func(wafListenerSettings, wafVhostSettings, wafRouteSettings *waf.Settings, dlpVhostSettings *dlp.Config, dlpRouteSettings *dlp.Config, dlpFilterSettings *dlp.FilterConfig) {

				By("tmp file " + tmpFileFSName + " " + tmpFileDMName)

				proxy = getProxyWaf(envoyPort, testUpstream.Upstream.Metadata.Ref(), wafListenerSettings, wafVhostSettings, dlpVhostSettings, dlpRouteSettings, wafRouteSettings, dlpFilterSettings)
				proxy.Listeners[0].Options = &gloov1.ListenerOptions{
					AccessLoggingService: &als.AccessLoggingService{
						AccessLog: []*als.AccessLog{
							{
								OutputDestination: &als.AccessLog_FileSink{
									FileSink: &als.FileSink{
										Path: tmpFileFSName,
										OutputFormat: &als.FileSink_StringFormat{
											StringFormat: "%FILTER_STATE(io.solo.modsecurity.audit_log)%",
										},
									},
								},
							}, {
								OutputDestination: &als.AccessLog_FileSink{
									FileSink: &als.FileSink{
										Path: tmpFileDMName,
										OutputFormat: &als.FileSink_StringFormat{
											StringFormat: "%DYNAMIC_METADATA(io.solo.filters.http.modsecurity:audit_log)%",
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
				}, "10s", "0.5s")
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
							"user-agent":  {"nikto"},
							"test-header": {"test-value"},
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
						Header: map[string][]string{
							"test-header": {"test-value"},
						},
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
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil, nil, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "10s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "10s", "1s").Should(Equal("-"))
			})

			It("auditlog listener dynamic meta", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_DYNAMIC_METADATA,
					},
				}, nil, nil, nil, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessDMLog, "10s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessFSLog, "10s", "1s").Should(Equal("-"))
			})
			It("auditlog listener fs - logs relevant", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_RELEVANT_ONLY,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil, nil, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "10s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "10s", "1s").Should(Equal("-"))
			})
			It("auditlog listener fs - not log not relevant", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_RELEVANT_ONLY,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil, nil, nil, nil)
				makeGoodRequest()
				// check the logs
				Eventually(getAccessFSLog, "10s", "1s").Should(Equal("-"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "10s", "1s").Should(Equal("-"))
			})
			It("auditlog listener dm - not log not relevant", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_RELEVANT_ONLY,
						Location: envoywaf.AuditLogging_DYNAMIC_METADATA,
					},
				}, nil, nil, nil, nil, nil)
				makeGoodRequest()
				// nothing written to dm log
				Eventually(getAccessDMLog, "10s", "1s").Should(Equal("-"))
				Eventually(getAccessFSLog, "10s", "1s").Should(Equal("-"))
			})
			It("auditlog listener dm - not log relevant if disabled", func() {
				startProxy(&waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
				}, nil, nil, nil, nil, nil)
				makeBadRequest()
				// nothing written to any logs
				Eventually(getAccessDMLog, "10s", "1s").Should(Equal("-"))
				Eventually(getAccessFSLog, "10s", "1s").Should(Equal("-"))
			})

			It("auditlog vhost filter state", func() {
				startProxy(nil, &waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "10s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "10s", "1s").Should(Equal("-"))
			})

			It("auditlog route filter state", func() {
				startProxy(nil, nil, &waf.Settings{
					RuleSets: []*envoywaf.RuleSet{getRulesTemplate(true, true, true, false)},
					AuditLogging: &envoywaf.AuditLogging{
						Action:   envoywaf.AuditLogging_ALWAYS,
						Location: envoywaf.AuditLogging_FILTER_STATE,
					},
				}, nil, nil, nil)
				makeBadRequest()
				// check the logs
				Eventually(getAccessFSLog, "10s", "1s").Should(ContainSubstring("nikto"))
				// nothing written to dm log
				Eventually(getAccessDMLog, "10s", "1s").Should(Equal("-"))
			})

			Context("DLP", func() {
				type setupOpts struct {
					routeDlp, vhostDlp, listenerDlp, logJSON bool
				}

				// Configure DLP to log to dynamic metadata and censor any instance of the text "test-value"
				var setupProxyDLP = func(opts setupOpts) {
					rule := getRulesTemplate(true, true, true, opts.logJSON)
					var dlpConfigVhost, dlpConfigRoute *dlp.Config
					var dlpConfigListener *dlp.FilterConfig
					if opts.routeDlp {
						dlpConfigRoute = &dlp.Config{
							Actions: []*dlp.Action{{
								ActionType: dlp.Action_CUSTOM,
								CustomAction: &dlp.CustomAction{
									Name:     "route-action",
									MaskChar: "R",
									Percent: &envoy_type.Percent{
										Value: 100,
									},
									RegexActions: []*transformation_ee.RegexAction{
										{
											Regex:    "(.*)(test-value)(.*)",
											Subgroup: 2,
										},
									},
								},
							}},
							EnabledFor: dlp.Config_ALL,
						}
					}
					if opts.vhostDlp {
						dlpConfigVhost = &dlp.Config{
							Actions: []*dlp.Action{{
								ActionType: dlp.Action_CUSTOM,
								CustomAction: &dlp.CustomAction{
									Name:     "vhost-action",
									MaskChar: "V",
									Percent: &envoy_type.Percent{
										Value: 100,
									},
									RegexActions: []*transformation_ee.RegexAction{
										{
											Regex:    "(.*)(test-value)(.*)",
											Subgroup: 2,
										},
									},
								},
							}},
							EnabledFor: dlp.Config_ALL,
						}
					}
					if opts.listenerDlp {
						dlpConfigListener = &dlp.FilterConfig{
							EnabledFor: dlp.FilterConfig_ALL,
							DlpRules: []*dlp.DlpRule{
								{
									Actions: []*dlp.Action{{
										ActionType: dlp.Action_CUSTOM,
										CustomAction: &dlp.CustomAction{
											Name:     "listener-action",
											MaskChar: "L",
											Percent: &envoy_type.Percent{
												Value: 100,
											},
											RegexActions: []*transformation_ee.RegexAction{
												{
													Regex:    "(.*)(test-value)(.*)",
													Subgroup: 2,
												},
											},
										},
									}},
								},
							},
						}
					}

					startProxy(&waf.Settings{
						RuleSets: []*envoywaf.RuleSet{rule},
						AuditLogging: &envoywaf.AuditLogging{
							Action:   envoywaf.AuditLogging_ALWAYS,
							Location: envoywaf.AuditLogging_DYNAMIC_METADATA,
						},
					}, nil, nil, dlpConfigVhost, dlpConfigRoute,
						dlpConfigListener,
					)
				}
				Context("route-level dlp", func() {
					Describe("string log format", func() {
						BeforeEach(func() {
							setupProxyDLP(setupOpts{routeDlp: true})
						})
						It("censors logs made to dynamic metadata when a good request is made", func() {
							makeGoodRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
					})
					Describe("json log format", func() {
						BeforeEach(func() {
							setupProxyDLP(setupOpts{routeDlp: true, logJSON: true})
						})
						It("censors logs made to dynamic metadata when a good request is made", func() {
							makeGoodRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
					})
				})
				Context("vhost-level dlp", func() {
					Describe("string log format", func() {
						BeforeEach(func() {
							setupProxyDLP(setupOpts{vhostDlp: true})
						})
						It("censors logs made to dynamic metadata when a good request is made", func() {
							makeGoodRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
					})
					Describe("json log format", func() {
						BeforeEach(func() {
							setupProxyDLP(setupOpts{vhostDlp: true, logJSON: true})
						})
						It("censors logs made to dynamic metadata when a good request is made", func() {
							makeGoodRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
					})
				})
				Context("listener-level dlp", func() {
					Describe("string log format", func() {
						BeforeEach(func() {
							setupProxyDLP(setupOpts{listenerDlp: true})
						})
						It("censors logs made to dynamic metadata when a good request is made", func() {
							makeGoodRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
					})
					Describe("json log format", func() {
						BeforeEach(func() {
							setupProxyDLP(setupOpts{listenerDlp: true, logJSON: true})
						})
						It("censors logs made to dynamic metadata when a good request is made", func() {
							makeGoodRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
						})
					})
				})
				Context("multi-level dlp", func() {
					Describe("route+vhost", func() {
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							setupProxyDLP(setupOpts{routeDlp: true, vhostDlp: true})
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
							Eventually(getAccessDMLog, "10s", "1s").Should(ContainSubstring("RRRR-RRRRR"))
						})
					})
					Describe("route+listener", func() {
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							setupProxyDLP(setupOpts{routeDlp: true, listenerDlp: true})
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
							Eventually(getAccessDMLog, "10s", "1s").Should(ContainSubstring("RRRR-RRRRR"))
						})
					})
					Describe("vhost+listener", func() {
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							setupProxyDLP(setupOpts{vhostDlp: true, listenerDlp: true})
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
							Eventually(getAccessDMLog, "10s", "1s").Should(ContainSubstring("VVVV-VVVVV"))
						})
					})
					Describe("route+vhost+listener", func() {
						It("censors logs made to dynamic metadata when a bad request is made", func() {
							setupProxyDLP(setupOpts{routeDlp: true, vhostDlp: true, listenerDlp: true})
							makeBadRequest()
							// should be no logs to Filter State
							Eventually(getAccessFSLog, "10s", "1s").Should(BeEquivalentTo("-"))

							// logs to dynamic metadata should not contain the masked substring
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(BeEquivalentTo("-"))
							Eventually(getAccessDMLog, "10s", "1s").ShouldNot(ContainSubstring("test-value"))
							Eventually(getAccessDMLog, "10s", "1s").Should(ContainSubstring("RRRR-RRRRR"))
						})
					})
				})
			})
		})
	})
})
