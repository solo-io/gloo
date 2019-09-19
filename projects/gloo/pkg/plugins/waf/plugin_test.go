package waf

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	envoywaf "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/waf"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("waf plugin", func() {
	var (
		plugin       *Plugin
		params       plugins.Params
		vhostParams  plugins.VirtualHostParams
		virtualHost  *v1.VirtualHost
		route        *v1.Route
		httpListener *v1.HttpListener
		wafVhost     *waf.VhostSettings
		wafRoute     *waf.RouteSettings
		wafListener  *waf.Settings
	)

	const (
		rulesString    = "rules rules rules"
		crsRulesString = "crs rules rules rules"
	)

	JustBeforeEach(func() {
		if wafRoute == nil {
			wafRoute = &waf.RouteSettings{}
		}
		wafRouteSt, err := util.MessageToStruct(wafRoute)
		Expect(err).NotTo(HaveOccurred())
		route = &v1.Route{
			Matcher: &v1.Matcher{
				PathSpecifier: &v1.Matcher_Prefix{
					Prefix: "/",
				},
			},
			Action: &v1.Route_DirectResponseAction{
				DirectResponseAction: &v1.DirectResponseAction{
					Status: 200,
					Body:   "test",
				},
			},
			RoutePlugins: &v1.RoutePlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: wafRouteSt,
					},
				},
			},
		}

		if wafVhost == nil {
			wafVhost = &waf.VhostSettings{}
		}
		wafVhostSt, err := util.MessageToStruct(wafVhost)
		Expect(err).NotTo(HaveOccurred())

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &v1.VirtualHostPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: wafVhostSt,
					},
				},
			},
			Routes: []*v1.Route{route},
		}

		if wafListener == nil {
			wafListener = &waf.Settings{}
		}
		wafListenerSt, err := util.MessageToStruct(wafListener)
		Expect(err).NotTo(HaveOccurred())

		httpListener = &v1.HttpListener{
			VirtualHosts: []*v1.VirtualHost{virtualHost},
			ListenerPlugins: &v1.HttpListenerPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: wafListenerSt,
					},
				},
			},
		}
		proxy := &v1.Proxy{
			Metadata: core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Listeners: []*v1.Listener{{
				Name: "default",
				ListenerType: &v1.Listener_HttpListener{
					HttpListener: httpListener,
				},
			}},
		}

		params.Snapshot = &v1.ApiSnapshot{
			Proxies: v1.ProxyList{proxy},
		}
		vhostParams = plugins.VirtualHostParams{
			Params:   params,
			Proxy:    proxy,
			Listener: proxy.Listeners[0],
		}

	})

	Context("process snapshot", func() {
		var (
			outRoute   envoyroute.Route
			outVhost   envoyroute.VirtualHost
			outFilters []plugins.StagedHttpFilter
		)

		var checkRuleSets = func(rs []*envoywaf.RuleSet) {
			Expect(rs).To(HaveLen(3))
			Expect(rs[0].Files).To(BeNil())
			Expect(rs[0].RuleStr).To(Equal(rulesString))
			Expect(rs[1].Files).To(BeNil())
			Expect(rs[1].RuleStr).To(Equal(crsRulesString))
			Expect(rs[2].Files).To(Equal(getCoreRuleSetFiles()))
			Expect(rs[2].RuleStr).To(Equal(""))
		}

		JustBeforeEach(func() {
			outVhost = envoyroute.VirtualHost{
				Name: "test",
			}
			outRoute = envoyroute.Route{}
			routesParams := plugins.RouteParams{
				VirtualHostParams: vhostParams,
				VirtualHost:       virtualHost,
			}
			// run it like the translator:
			err := plugin.ProcessRoute(routesParams, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			err = plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.HttpFilters(params, httpListener)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			plugin = NewPlugin()
			plugin.Init(plugins.InitParams{})
		})

		Context("empty extensions", func() {

			It("can create the proper filters", func() {
				Expect(outFilters).To(HaveLen(1))
				wafFilter := outFilters[0]
				Expect(wafFilter.HttpFilter.Name).To(Equal(FilterName))
				Expect(wafFilter.Stage).To(Equal(plugins.DuringStage(plugins.WafStage)))
				st := wafFilter.HttpFilter.GetConfig()
				Expect(st).NotTo(BeNil())
				var filterWaf envoywaf.ModSecurity
				err := util.StructToMessage(st, &filterWaf)
				Expect(err).NotTo(HaveOccurred())
				Expect(filterWaf.Disabled).To(BeTrue())
			})

		})

		Context("http filters", func() {
			BeforeEach(func() {
				ruleSets := []*envoywaf.RuleSet{
					{
						RuleStr: rulesString,
					},
				}
				wafListener = &waf.Settings{
					CoreRuleSet: &waf.CoreRuleSet{
						CustomSettingsType: &waf.CoreRuleSet_CustomSettingsString{
							CustomSettingsString: crsRulesString,
						},
					},
					RuleSets: ruleSets,
				}
			})

			It("can create the proper http filters", func() {
				Expect(outFilters).To(HaveLen(1))
				wafFilter := outFilters[0]
				Expect(wafFilter.HttpFilter.Name).To(Equal(FilterName))
				Expect(wafFilter.Stage).To(Equal(plugins.DuringStage(plugins.WafStage)))
				st := wafFilter.HttpFilter.GetConfig()
				Expect(st).NotTo(BeNil())
				var filterWaf envoywaf.ModSecurity
				err := util.StructToMessage(st, &filterWaf)
				Expect(err).NotTo(HaveOccurred())
				checkRuleSets(filterWaf.RuleSets)
			})
		})

		Context("per route/vhost", func() {
			Context("nil", func() {
				BeforeEach(func() {
					wafRoute = &waf.RouteSettings{
						Disabled: true,
						Settings: &waf.Settings{
							CoreRuleSet: nil,
							RuleSets:    nil,
						},
					}

					wafVhost = &waf.VhostSettings{
						Disabled: true,
						Settings: &waf.Settings{
							CoreRuleSet: nil,
							RuleSets:    nil,
						},
					}
				})

				It("sets disabled on route", func() {
					pfc := outRoute.PerFilterConfig[FilterName]
					Expect(pfc).NotTo(BeNil())
					var perRouteWaf envoywaf.ModSecurityPerRoute
					err := util.StructToMessage(pfc, &perRouteWaf)
					Expect(err).NotTo(HaveOccurred())
					Expect(perRouteWaf.Disabled).To(BeTrue())
				})

				It("sets disabled on vhost", func() {
					pfc := outVhost.PerFilterConfig[FilterName]
					Expect(pfc).NotTo(BeNil())
					var perVhostWaf envoywaf.ModSecurityPerRoute
					err := util.StructToMessage(pfc, &perVhostWaf)
					Expect(err).NotTo(HaveOccurred())
					Expect(perVhostWaf.Disabled).To(BeTrue())
				})
			})

			Context("filled in", func() {

				BeforeEach(func() {
					ruleSets := []*envoywaf.RuleSet{
						{
							RuleStr: rulesString,
						},
					}
					wafRoute = &waf.RouteSettings{
						Settings: &waf.Settings{
							CoreRuleSet: &waf.CoreRuleSet{
								CustomSettingsType: &waf.CoreRuleSet_CustomSettingsString{
									CustomSettingsString: crsRulesString,
								},
							},
							RuleSets: ruleSets,
						},
					}

					wafVhost = &waf.VhostSettings{
						Settings: &waf.Settings{
							CoreRuleSet: &waf.CoreRuleSet{
								CustomSettingsType: &waf.CoreRuleSet_CustomSettingsString{
									CustomSettingsString: crsRulesString,
								},
							},
							RuleSets: ruleSets,
						},
					}
				})

				It("sets disabled on route", func() {
					pfc := outRoute.PerFilterConfig[FilterName]
					Expect(pfc).NotTo(BeNil())
					var perRouteWaf envoywaf.ModSecurityPerRoute
					err := util.StructToMessage(pfc, &perRouteWaf)
					Expect(err).NotTo(HaveOccurred())
					Expect(perRouteWaf.Disabled).To(BeFalse())
					checkRuleSets(perRouteWaf.RuleSets)
				})

				It("sets disabled on vhost", func() {
					pfc := outVhost.PerFilterConfig[FilterName]
					Expect(pfc).NotTo(BeNil())
					var perVhostWaf envoywaf.ModSecurityPerRoute
					err := util.StructToMessage(pfc, &perVhostWaf)
					Expect(err).NotTo(HaveOccurred())
					Expect(perVhostWaf.Disabled).To(BeFalse())
					checkRuleSets(perVhostWaf.RuleSets)
				})
			})
		})

	})

})
