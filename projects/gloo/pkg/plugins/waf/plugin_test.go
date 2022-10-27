package waf

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	envoywaf "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/waf"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	envoy_type "github.com/solo-io/solo-kit/pkg/api/external/envoy/type"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("waf plugin", func() {
	var (
		plugin            plugins.Plugin
		params            plugins.Params
		vhostParams       plugins.VirtualHostParams
		virtualHost       *v1.VirtualHost
		route             *v1.Route
		httpListener      *v1.HttpListener
		wafVhost          *waf.Settings
		wafRoute          *waf.Settings
		dlpSettings       *dlp.Config
		wafListener       *waf.Settings
		configMapRuleSets []*waf.RuleSetFromConfigMap
	)

	const (
		rulesString               = "rules rules rules"
		crsRulesString            = "crs rules rules rules"
		ruleStringFromArtifact1   = "rule from artifact 1"
		ruleStringFromArtifact2   = "rule from artifact 2"
		ruleStringFromArtifact3   = "rule from artifact 3"
		ruleStringFromArtifact4   = "rule from artifact 4"
		customInterventionMessage = "custom intervention message"
	)

	BeforeEach(func() {
		wafListener = &waf.Settings{}
	})

	JustBeforeEach(func() {
		if wafRoute == nil {
			wafRoute = &waf.Settings{}
		}
		if dlpSettings == nil {
			dlpSettings = &dlp.Config{}
		}
		route = &v1.Route{
			Matchers: []*matchers.Matcher{{
				PathSpecifier: &matchers.Matcher_Prefix{
					Prefix: "/",
				},
			}},
			Action: &v1.Route_DirectResponseAction{
				DirectResponseAction: &v1.DirectResponseAction{
					Status: 200,
					Body:   "test",
				},
			},
			Options: &v1.RouteOptions{
				Waf: wafRoute,
				Dlp: dlpSettings,
			},
		}

		if wafVhost == nil {
			wafVhost = &waf.Settings{}
		}

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			Options: &v1.VirtualHostOptions{
				Waf: wafVhost,
				Dlp: dlpSettings,
			},
			Routes: []*v1.Route{route},
		}

		httpListener = &v1.HttpListener{
			VirtualHosts: []*v1.VirtualHost{virtualHost},
			Options: &v1.HttpListenerOptions{
				Waf: wafListener,
			},
		}
		proxy := &v1.Proxy{
			Metadata: &core.Metadata{
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

		params.Snapshot = &v1snap.ApiSnapshot{
			Proxies: v1.ProxyList{proxy},
		}
		vhostParams = plugins.VirtualHostParams{
			Params:   params,
			Proxy:    proxy,
			Listener: proxy.Listeners[0],
		}

		configMapRuleSets = []*waf.RuleSetFromConfigMap{
			{
				ConfigMapRef: &core.ResourceRef{
					Namespace: "namespace",
					Name:      "name",
				},
				DataMapKeys: []string{
					"key1", //order in artifact is key2, then key1 - validate order of passed keys is respected
					"key2",
				},
			},
			{
				ConfigMapRef: &core.ResourceRef{
					Namespace: "namespace",
					Name:      "name2",
				},
			},
		}

	})

	allTests := func() {
		Context("process snapshot", func() {
			var (
				outRoute   envoy_config_route_v3.Route
				outVhost   envoy_config_route_v3.VirtualHost
				outFilters []plugins.StagedHttpFilter
			)

			var checkRuleSets = func(rs []*envoywaf.RuleSet) {
				Expect(rs).To(HaveLen(7))
				Expect(rs[0].Files).To(BeNil())
				Expect(rs[0].RuleStr).To(Equal(rulesString))
				Expect(rs[1].Files).To(BeNil())
				Expect(rs[1].RuleStr).To(Equal(crsRulesString))
				Expect(rs[2].Directory).To(Equal(crsPathPrefix))
				Expect(rs[2].RuleStr).To(Equal(""))
				//rules from configmap
				Expect(rs[3].Files).To(BeNil())
				Expect(rs[3].RuleStr).To(Equal(ruleStringFromArtifact1))
				Expect(rs[4].Files).To(BeNil())
				Expect(rs[4].RuleStr).To(Equal(ruleStringFromArtifact2))
				Expect(rs[5].Files).To(BeNil())
				Expect(rs[5].RuleStr).To(Equal(ruleStringFromArtifact3))
				Expect(rs[6].Files).To(BeNil())
				Expect(rs[6].RuleStr).To(Equal(ruleStringFromArtifact4))
			}

			JustBeforeEach(func() {
				outVhost = envoy_config_route_v3.VirtualHost{
					Name: "test",
				}
				outRoute = envoy_config_route_v3.Route{}
				routesParams := plugins.RouteParams{
					VirtualHostParams: vhostParams,
					VirtualHost:       virtualHost,
				}

				params.Snapshot.Artifacts = v1.ArtifactList{
					{
						Metadata: &core.Metadata{
							Name:      "name",
							Namespace: "namespace",
						},
						Data: map[string]string{
							"key2": ruleStringFromArtifact2, //rule will come after rule from source.conf - tests sorting of keys
							"key1": ruleStringFromArtifact1,
						},
					},
					{
						Metadata: &core.Metadata{
							Name:      "name2",
							Namespace: "namespace",
						},
						Data: map[string]string{
							"key4": ruleStringFromArtifact4, //should be sorted to key3, key4 as no keys are passed
							"key3": ruleStringFromArtifact3,
						},
					},
				}

				// run it like the translator:
				err := plugin.(plugins.RoutePlugin).ProcessRoute(routesParams, route, &outRoute)
				Expect(err).NotTo(HaveOccurred())
				err = plugin.(plugins.VirtualHostPlugin).ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
				Expect(err).NotTo(HaveOccurred())
				outFilters, err = plugin.(plugins.HttpFilterPlugin).HttpFilters(params, httpListener)
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
					goTypedConfig := wafFilter.HttpFilter.GetTypedConfig()
					Expect(goTypedConfig).NotTo(BeNil())
					var filterWaf envoywaf.ModSecurity
					err := ptypes.UnmarshalAny(goTypedConfig, &filterWaf)
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
						RuleSets:                  ruleSets,
						ConfigMapRuleSets:         configMapRuleSets,
						CustomInterventionMessage: customInterventionMessage,
						RequestHeadersOnly:        true,
						ResponseHeadersOnly:       true,
					}
				})

				It("can create the proper http filters", func() {
					Expect(outFilters).To(HaveLen(1))
					wafFilter := outFilters[0]
					Expect(wafFilter.HttpFilter.Name).To(Equal(FilterName))
					Expect(wafFilter.Stage).To(Equal(plugins.DuringStage(plugins.WafStage)))
					goTypedConfig := wafFilter.HttpFilter.GetTypedConfig()
					Expect(goTypedConfig).NotTo(BeNil())
					var filterWaf envoywaf.ModSecurity
					err := ptypes.UnmarshalAny(goTypedConfig, &filterWaf)
					Expect(err).NotTo(HaveOccurred())
					checkRuleSets(filterWaf.RuleSets)
					Expect(filterWaf.RequestHeadersOnly).To(BeTrue())
					Expect(filterWaf.ResponseHeadersOnly).To(BeTrue())
				})
			})

			Context("per route/vhost", func() {
				Context("nil", func() {
					BeforeEach(func() {
						wafRoute = &waf.Settings{
							Disabled:    true,
							CoreRuleSet: nil,
							RuleSets:    nil,
						}

						wafVhost = &waf.Settings{
							Disabled:    true,
							CoreRuleSet: nil,
							RuleSets:    nil,
						}
					})

					It("sets disabled on route", func() {
						goTpfc := outRoute.TypedPerFilterConfig[FilterName]
						Expect(goTpfc).NotTo(BeNil())
						var perRouteWaf envoywaf.ModSecurityPerRoute
						err := ptypes.UnmarshalAny(goTpfc, &perRouteWaf)
						Expect(err).NotTo(HaveOccurred())
						Expect(perRouteWaf.Disabled).To(BeTrue())
					})

					It("sets disabled on vhost", func() {
						goTpfc := outVhost.TypedPerFilterConfig[FilterName]
						Expect(goTpfc).NotTo(BeNil())
						var perVhostWaf envoywaf.ModSecurityPerRoute
						err := ptypes.UnmarshalAny(goTpfc, &perVhostWaf)
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

						wafRoute = &waf.Settings{
							CoreRuleSet: &waf.CoreRuleSet{
								CustomSettingsType: &waf.CoreRuleSet_CustomSettingsString{
									CustomSettingsString: crsRulesString,
								},
							},
							RuleSets:                  ruleSets,
							ConfigMapRuleSets:         configMapRuleSets,
							CustomInterventionMessage: customInterventionMessage,
							RequestHeadersOnly:        true,
							ResponseHeadersOnly:       true,
						}

						wafVhost = &waf.Settings{
							CoreRuleSet: &waf.CoreRuleSet{
								CustomSettingsType: &waf.CoreRuleSet_CustomSettingsString{
									CustomSettingsString: crsRulesString,
								},
							},
							RuleSets:                  ruleSets,
							ConfigMapRuleSets:         configMapRuleSets,
							CustomInterventionMessage: customInterventionMessage,
							RequestHeadersOnly:        true,
							ResponseHeadersOnly:       true,
						}
					})

					It("sets properties on route", func() {
						goTpfc := outRoute.TypedPerFilterConfig[FilterName]
						Expect(goTpfc).NotTo(BeNil())
						var perRouteWaf envoywaf.ModSecurityPerRoute
						err := ptypes.UnmarshalAny(goTpfc, &perRouteWaf)
						Expect(err).NotTo(HaveOccurred())
						Expect(perRouteWaf.Disabled).To(BeFalse())
						Expect(perRouteWaf.RequestHeadersOnly).To(BeTrue())
						Expect(perRouteWaf.ResponseHeadersOnly).To(BeTrue())
						checkRuleSets(perRouteWaf.RuleSets)
					})

					It("sets properties on vhost", func() {
						goTpfc := outVhost.TypedPerFilterConfig[FilterName]
						Expect(goTpfc).NotTo(BeNil())
						var perVhostWaf envoywaf.ModSecurityPerRoute
						err := ptypes.UnmarshalAny(goTpfc, &perVhostWaf)
						Expect(err).NotTo(HaveOccurred())
						Expect(perVhostWaf.Disabled).To(BeFalse())
						Expect(perVhostWaf.RequestHeadersOnly).To(BeTrue())
						Expect(perVhostWaf.ResponseHeadersOnly).To(BeTrue())
						checkRuleSets(perVhostWaf.RuleSets)
					})
				})
				Context("with DLP", func() {
					BeforeEach(func() {
						customTestAction := &dlp.Action{
							ActionType: dlp.Action_CUSTOM,
							Shadow:     true,
							CustomAction: &dlp.CustomAction{
								Name:  "test",
								Regex: []string{"regex"},
								Percent: &envoy_type.Percent{
									Value: 75,
								},
								MaskChar: "Z",
								RegexActions: []*transformation_ee.RegexAction{
									{Regex: "actionRegex", Subgroup: 1},
								},
							},
						}

						dlpSettings = &dlp.Config{
							Actions:    []*dlp.Action{customTestAction},
							EnabledFor: dlp.Config_ALL,
						}
					})

					It("sets properties on route", func() {
						goTpfc := outRoute.TypedPerFilterConfig[FilterName]
						Expect(goTpfc).NotTo(BeNil())
						var perRouteWaf envoywaf.ModSecurityPerRoute
						err := ptypes.UnmarshalAny(goTpfc, &perRouteWaf)
						Expect(err).NotTo(HaveOccurred())
						Expect(perRouteWaf.Disabled).To(BeFalse())

						Expect(perRouteWaf.DlpTransformation).ToNot(BeNil())
						Expect(perRouteWaf.DlpTransformation.Actions).To(HaveLen(1))
						action := perRouteWaf.DlpTransformation.Actions[0]
						Expect(action.MaskChar).To(BeEquivalentTo("Z"))
						Expect(action.Shadow).To(BeEquivalentTo(true))
						Expect(action.Name).To(BeEquivalentTo("test"))
						Expect(action.Regex).To(BeEquivalentTo([]string{"regex"}))
						Expect(action.GetMatcher().GetRegexMatcher().GetRegexActions()).To(HaveLen(1))
						regexAction := action.GetMatcher().GetRegexMatcher().GetRegexActions()[0]
						Expect(regexAction.Regex).To(BeEquivalentTo("actionRegex"))
						Expect(regexAction.Subgroup).To(BeEquivalentTo(1))
					})

					It("sets properties on vhost", func() {
						goTpfc := outVhost.TypedPerFilterConfig[FilterName]
						Expect(goTpfc).NotTo(BeNil())
						var perVhostWaf envoywaf.ModSecurityPerRoute
						err := ptypes.UnmarshalAny(goTpfc, &perVhostWaf)
						Expect(err).NotTo(HaveOccurred())
						Expect(perVhostWaf.Disabled).To(BeFalse())

						Expect(perVhostWaf.DlpTransformation).ToNot(BeNil())
						Expect(perVhostWaf.DlpTransformation.Actions).To(HaveLen(1))
						action := perVhostWaf.DlpTransformation.Actions[0]
						Expect(action.MaskChar).To(BeEquivalentTo("Z"))
						Expect(action.Shadow).To(BeEquivalentTo(true))
						Expect(action.Name).To(BeEquivalentTo("test"))
						Expect(action.Regex).To(BeEquivalentTo([]string{"regex"}))
						Expect(action.GetMatcher().GetRegexMatcher().GetRegexActions()).To(HaveLen(1))
						regexAction := action.GetMatcher().GetRegexMatcher().GetRegexActions()[0]
						Expect(regexAction.Regex).To(BeEquivalentTo("actionRegex"))
						Expect(regexAction.Subgroup).To(BeEquivalentTo(1))
					})
				})
			})

		})
	}
	allTests()

})
