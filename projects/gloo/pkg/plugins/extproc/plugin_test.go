package extproc_test

import (
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloo_config_core_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	gloo_ext_proc_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/ext_proc/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/filters"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	extproc_plugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extproc"
	"github.com/solo-io/solo-projects/test/extproc/builders"
)

var _ = Describe("external processing plugin", func() {

	var (
		defaultUpstreamList = v1.UpstreamList{
			buildUpstream(builders.DefaultExtProcUpstreamName),
			buildUpstream(builders.OverrideExtProcUpstreamName),
		}
	)

	Context("listener", func() {
		type inputParams struct {
			// builds the extproc settings to put in Settings CR
			globalBuilder *builders.ExtProcBuilder
			// builds the extproc settings to put on the HttpListenerOptions
			listenerBuilder *builders.ExtProcBuilder
			// disabled flag on HttpListenerOptions
			extProcDisabled *wrappers.BoolValue
		}

		type expectedOutput struct {
			envoyConfig *envoy_ext_proc_v3.ExternalProcessor
			stage       *plugins.FilterStage
			err         string
		}

		var (
			// functions that get the builders used to generate global (Settings CR)
			// and listener-level (HttpListenerOptions) extproc settings.
			// (define these as functions because the builders get mutated in some of the tests)
			getDefaultGlobalBuilder = func() *builders.ExtProcBuilder {
				return builders.GetDefaultExtProcBuilder()
			}
			getDefaultListenerBuilder = func() *builders.ExtProcBuilder {
				return builders.GetDefaultExtProcBuilder().
					WithGrpcServiceBuilder(builders.GetDefaultGrpcServiceBuilder().
						WithUpstreamName(builders.OverrideExtProcUpstreamName).
						WithAuthority(nil).
						WithInitialMetadata([]*gloo_config_core_v3.HeaderValue{{Key: "E", Value: "F"}}),
					).
					WithStage(&filters.FilterStage{Stage: filters.FilterStage_CorsStage, Predicate: filters.FilterStage_After}).
					WithFailureModeAllow(nil).
					WithProcessingMode(&gloo_ext_proc_v3.ProcessingMode{
						ResponseHeaderMode: gloo_ext_proc_v3.ProcessingMode_SKIP,
					}).
					WithMessageTimeout(nil).
					WithMaxMessageTimeout(&duration.Duration{Seconds: 10}).
					WithRequestAttributes(nil).
					WithResponseAttributes([]string{"new1", "new2"})
			}

			defaultGlobalStage   = plugins.BeforeStage(plugins.AcceptedStage)
			defaultListenerStage = plugins.AfterStage(plugins.CorsStage)
		)

		DescribeTable("listener-level extproc configuration",
			func(testInput inputParams, expected expectedOutput) {
				// initialize plugin with the specified global extproc settings
				plugin := extproc_plugin.NewPlugin()
				initParams := plugins.InitParams{
					Settings: &v1.Settings{},
				}
				if testInput.globalBuilder != nil {
					initParams.Settings.ExtProc = testInput.globalBuilder.Build()
				}
				plugin.Init(initParams)

				// set the snapshot upstreams
				pluginParams := plugins.Params{
					Snapshot: &v1snap.ApiSnapshot{
						Upstreams: defaultUpstreamList,
					},
				}

				// create a listener with the given extproc options
				listener := &v1.HttpListener{
					Options: &v1.HttpListenerOptions{},
				}
				// at most one of disabled or extproc settings can be set
				if testInput.extProcDisabled != nil {
					listener.Options.ExtProcConfig = &v1.HttpListenerOptions_DisableExtProc{
						DisableExtProc: testInput.extProcDisabled,
					}
				} else if testInput.listenerBuilder != nil {
					listener.Options.ExtProcConfig = &v1.HttpListenerOptions_ExtProc{
						ExtProc: testInput.listenerBuilder.Build(),
					}
				}

				// translate
				filters, err := plugin.HttpFilters(pluginParams, listener)
				if expected.err != "" {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expected.err))
					return
				}

				Expect(err).NotTo(HaveOccurred())
				if expected.envoyConfig == nil {
					Expect(filters).To(HaveLen(0))
					return
				}

				// check that it's the expected filter and stage
				Expect(filters).To(HaveLen(1))
				extProcFilter := filters[0]
				Expect(extProcFilter.HttpFilter.Name).To(Equal(extproc_plugin.FilterName))
				Expect(&extProcFilter.Stage).To(Equal(expected.stage))

				// extract the config
				typedConfig := extProcFilter.HttpFilter.GetTypedConfig()
				Expect(typedConfig).NotTo(BeNil())
				Expect(typedConfig.GetTypeUrl()).To(Equal("type.googleapis.com/envoy.extensions.filters.http.ext_proc.v3.ExternalProcessor"))
				var extProc envoy_ext_proc_v3.ExternalProcessor
				err = typedConfig.UnmarshalTo(&extProc)
				Expect(err).NotTo(HaveOccurred())

				expectEqual(&extProc, expected.envoyConfig)
			},

			// The following entries test different combinations of:
			// - existence ("exists") or nonexistence ("none") of global extproc settings
			// - existence ("exists") or nonexistence ("none") of listener extproc settings
			// - disabled flag (true/false/not set) on listener extproc settings
			Entry("global=none, listener=none, disabled=not set -> should not add any filters",
				inputParams{}, expectedOutput{}),
			Entry("global=none, listener=exists -> should use listener settings",
				inputParams{
					listenerBuilder: getDefaultListenerBuilder(),
				},
				expectedOutput{
					stage:       &defaultListenerStage,
					envoyConfig: getExpectedListenerConfig(),
				}),
			Entry("global=exists, listener=none, disabled=not set -> should use global settings",
				inputParams{
					globalBuilder: getDefaultGlobalBuilder(),
				},
				expectedOutput{
					stage:       &defaultGlobalStage,
					envoyConfig: getExpectedGlobalConfig(),
				}),
			Entry("global=exists, listener=exists -> should merge",
				inputParams{
					globalBuilder:   getDefaultGlobalBuilder(),
					listenerBuilder: getDefaultListenerBuilder(),
				},
				expectedOutput{
					stage:       &defaultListenerStage,
					envoyConfig: getExpectedMergedConfig(),
				}),
			Entry("global=exists, disabled=true -> should not add any filters",
				inputParams{
					globalBuilder:   getDefaultGlobalBuilder(),
					extProcDisabled: &wrappers.BoolValue{Value: true},
				},
				expectedOutput{}),
			Entry("global=exists, disabled=false -> should use global settings",
				inputParams{
					globalBuilder:   getDefaultGlobalBuilder(),
					extProcDisabled: &wrappers.BoolValue{Value: false},
				},
				expectedOutput{
					stage:       &defaultGlobalStage,
					envoyConfig: getExpectedGlobalConfig(),
				}),
			Entry("global=none, disabled=true -> should not add any filters",
				inputParams{
					extProcDisabled: &wrappers.BoolValue{Value: true},
				},
				expectedOutput{}),
			Entry("global=none, disabled=false -> should not add any filters",
				inputParams{
					extProcDisabled: &wrappers.BoolValue{Value: false},
				},
				expectedOutput{}),

			// error conditions
			Entry("no upstream ref -> should error",
				inputParams{
					listenerBuilder: getDefaultListenerBuilder().WithGrpcServiceBuilder(
						builders.NewGrpcServiceBuilder().WithUpstreamName(""),
					),
				},
				expectedOutput{
					err: extproc_plugin.NoServerRefErr.Error(),
				}),
			Entry("nonexistent upstream -> should error",
				inputParams{
					listenerBuilder: getDefaultListenerBuilder().WithGrpcServiceBuilder(
						builders.NewGrpcServiceBuilder().WithUpstreamName("invalid"),
					),
				},
				expectedOutput{
					err: extproc_plugin.ServerNotFoundErr(&core.ResourceRef{Name: "invalid", Namespace: "gloo-system"}).Error(),
				}),
			Entry("no filter stage -> should error",
				inputParams{
					listenerBuilder: getDefaultListenerBuilder().WithStage(nil),
				},
				expectedOutput{
					err: extproc_plugin.NoFilterStageErr.Error(),
				}),
			Entry("timeout out of range -> should error",
				inputParams{
					listenerBuilder: getDefaultListenerBuilder().WithMessageTimeout(&duration.Duration{Seconds: 3700}),
				},
				expectedOutput{
					err: extproc_plugin.MessageTimeoutOutOfRangeErr(3700).Error(),
				}),
			Entry("max timeout < timeout -> should error",
				inputParams{
					listenerBuilder: getDefaultListenerBuilder().
						WithMessageTimeout(&duration.Duration{Seconds: 50}).
						WithMaxMessageTimeout(&duration.Duration{Seconds: 25}),
				},
				expectedOutput{
					err: extproc_plugin.MaxMessageTimeoutErr(50, 25).Error(),
				}),
		)
	})

	Context("virtual host and route", func() {
		type inputParams struct {
			builder *builders.ExtProcRouteBuilder
		}

		type expectedOutput struct {
			envoyConfig *envoy_ext_proc_v3.ExtProcPerRoute
			err         string
		}

		var (
			getDefaultBuilder = func() *builders.ExtProcRouteBuilder {
				return builders.GetDefaultExtProcRouteBuilder()
			}
		)

		DescribeTable("virtual host-level extproc configuration",
			func(testInput inputParams, expected expectedOutput) {
				plugin := extproc_plugin.NewPlugin()
				params := plugins.VirtualHostParams{
					Params: plugins.Params{
						Snapshot: &v1snap.ApiSnapshot{
							Upstreams: defaultUpstreamList,
						},
					},
				}

				// initialize virtual host with extproc options
				vh := &v1.VirtualHost{
					Options: &v1.VirtualHostOptions{},
				}
				if testInput.builder != nil {
					vh.Options.ExtProc = testInput.builder.Build()
				}

				// translate
				out := &envoy_config_route_v3.VirtualHost{}
				err := plugin.ProcessVirtualHost(params, vh, out)
				if expected.err != "" {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expected.err))
					return
				}

				Expect(err).NotTo(HaveOccurred())

				typedConfig := out.TypedPerFilterConfig[extproc_plugin.FilterName]
				if expected.envoyConfig == nil {
					Expect(typedConfig).To(BeNil())
					return
				}

				// extract the config
				Expect(typedConfig).NotTo(BeNil())
				Expect(typedConfig.GetTypeUrl()).To(Equal("type.googleapis.com/envoy.extensions.filters.http.ext_proc.v3.ExtProcPerRoute"))
				var extProcPerRoute envoy_ext_proc_v3.ExtProcPerRoute
				err = typedConfig.UnmarshalTo(&extProcPerRoute)
				Expect(err).NotTo(HaveOccurred())

				expectRouteEqual(&extProcPerRoute, expected.envoyConfig)
			},
			Entry("no extproc config -> should not set extproc config on virtual host",
				inputParams{}, expectedOutput{}),
			Entry("extproc disabled -> should set extproc disabled on virtual host",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithDisabled(&wrappers.BoolValue{Value: true}),
				},
				expectedOutput{
					envoyConfig: getExpectedDisabledRouteConfig(),
				}),
			Entry("extproc overrides -> should set extproc overrides on virtual host",
				inputParams{
					builder: getDefaultBuilder(),
				},
				expectedOutput{
					envoyConfig: getExpectedOverridesRouteConfig(),
				}),

			// error conditions
			Entry("disabled false -> should error",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithDisabled(&wrappers.BoolValue{Value: false}),
				},
				expectedOutput{
					err: extproc_plugin.DisabledErr.Error(),
				}),
			Entry("no upstream ref -> should error",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithGrpcServiceBuilder(
						builders.NewGrpcServiceBuilder().WithUpstreamName(""),
					),
				},
				expectedOutput{
					err: extproc_plugin.NoServerRefErr.Error(),
				}),
			Entry("nonexistent upstream -> should error",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithGrpcServiceBuilder(
						builders.NewGrpcServiceBuilder().WithUpstreamName("invalid"),
					),
				},
				expectedOutput{
					err: extproc_plugin.ServerNotFoundErr(&core.ResourceRef{Name: "invalid", Namespace: "gloo-system"}).Error(),
				}),
		)

		DescribeTable("route-level extproc configuration",
			func(testInput inputParams, expected expectedOutput) {
				plugin := extproc_plugin.NewPlugin()
				params := plugins.RouteParams{
					VirtualHostParams: plugins.VirtualHostParams{
						Params: plugins.Params{
							Snapshot: &v1snap.ApiSnapshot{
								Upstreams: defaultUpstreamList,
							},
						},
					},
				}

				// initialize route with extproc options
				route := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				if testInput.builder != nil {
					route.Options.ExtProc = testInput.builder.Build()
				}

				// translate
				out := &envoy_config_route_v3.Route{}
				err := plugin.ProcessRoute(params, route, out)
				if expected.err != "" {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expected.err))
					return
				}

				Expect(err).NotTo(HaveOccurred())

				typedConfig := out.TypedPerFilterConfig[extproc_plugin.FilterName]
				if expected.envoyConfig == nil {
					Expect(typedConfig).To(BeNil())
					return
				}

				// extract the config
				Expect(typedConfig).NotTo(BeNil())
				Expect(typedConfig.GetTypeUrl()).To(Equal("type.googleapis.com/envoy.extensions.filters.http.ext_proc.v3.ExtProcPerRoute"))
				var extProcPerRoute envoy_ext_proc_v3.ExtProcPerRoute
				err = typedConfig.UnmarshalTo(&extProcPerRoute)
				Expect(err).NotTo(HaveOccurred())

				expectRouteEqual(&extProcPerRoute, expected.envoyConfig)
			},
			Entry("no extproc config -> should not set extproc config on route",
				inputParams{}, expectedOutput{}),
			Entry("extproc disabled -> should set extproc disabled on route",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithDisabled(&wrappers.BoolValue{Value: true}),
				},
				expectedOutput{
					envoyConfig: getExpectedDisabledRouteConfig(),
				}),
			Entry("extproc overrides -> should set extproc overrides on route",
				inputParams{
					builder: getDefaultBuilder(),
				},
				expectedOutput{
					envoyConfig: getExpectedOverridesRouteConfig(),
				}),

			// error conditions
			Entry("disabled false -> should error",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithDisabled(&wrappers.BoolValue{Value: false}),
				},
				expectedOutput{
					err: extproc_plugin.DisabledErr.Error(),
				}),
			Entry("no upstream ref -> should error",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithGrpcServiceBuilder(
						builders.NewGrpcServiceBuilder().WithUpstreamName(""),
					),
				},
				expectedOutput{
					err: extproc_plugin.NoServerRefErr.Error(),
				}),
			Entry("nonexistent upstream -> should error",
				inputParams{
					builder: builders.NewExtProcRouteBuilder().WithGrpcServiceBuilder(
						builders.NewGrpcServiceBuilder().WithUpstreamName("invalid"),
					),
				},
				expectedOutput{
					err: extproc_plugin.ServerNotFoundErr(&core.ResourceRef{Name: "invalid", Namespace: "gloo-system"}).Error(),
				}),
		)
	})
})

// getExpectedGlobalConfig builds envoy extproc config with the values used in the default global ExtProcBuilder.
func getExpectedGlobalConfig() *envoy_ext_proc_v3.ExternalProcessor {
	return &envoy_ext_proc_v3.ExternalProcessor{
		GrpcService: &envoy_config_core_v3.GrpcService{
			TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
					ClusterName: translator.UpstreamToClusterName(&core.ResourceRef{Name: builders.DefaultExtProcUpstreamName, Namespace: "gloo-system"}),
					Authority:   "xyz",
					RetryPolicy: &envoy_config_core_v3.RetryPolicy{
						RetryBackOff: &envoy_config_core_v3.BackoffStrategy{
							BaseInterval: &duration.Duration{Seconds: 5},
							MaxInterval:  &duration.Duration{Seconds: 10},
						},
						NumRetries: &wrappers.UInt32Value{Value: 7},
					},
				},
			},
			Timeout: &duration.Duration{Seconds: 100},
			InitialMetadata: []*envoy_config_core_v3.HeaderValue{
				{Key: "A", Value: "B"},
				{Key: "C", Value: "D"},
			},
		},
		FailureModeAllow: true,
		ProcessingMode: &envoy_ext_proc_v3.ProcessingMode{
			RequestHeaderMode:   envoy_ext_proc_v3.ProcessingMode_SEND,
			ResponseHeaderMode:  envoy_ext_proc_v3.ProcessingMode_SEND,
			RequestBodyMode:     envoy_ext_proc_v3.ProcessingMode_BUFFERED,
			ResponseBodyMode:    envoy_ext_proc_v3.ProcessingMode_BUFFERED_PARTIAL,
			RequestTrailerMode:  envoy_ext_proc_v3.ProcessingMode_SKIP,
			ResponseTrailerMode: envoy_ext_proc_v3.ProcessingMode_DEFAULT,
		},
		MessageTimeout:     &duration.Duration{Seconds: 1},
		MaxMessageTimeout:  &duration.Duration{Seconds: 5},
		RequestAttributes:  []string{"req1", "req2"},
		ResponseAttributes: []string{"resp1", "resp2", "resp3"},
	}
}

// getExpectedListenerConfig builds envoy extproc config with the values used in the default listener ExtProcBuilder.
func getExpectedListenerConfig() *envoy_ext_proc_v3.ExternalProcessor {
	return &envoy_ext_proc_v3.ExternalProcessor{
		GrpcService: &envoy_config_core_v3.GrpcService{
			TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
					ClusterName: translator.UpstreamToClusterName(&core.ResourceRef{Name: builders.OverrideExtProcUpstreamName, Namespace: "gloo-system"}),
					RetryPolicy: &envoy_config_core_v3.RetryPolicy{
						RetryBackOff: &envoy_config_core_v3.BackoffStrategy{
							BaseInterval: &duration.Duration{Seconds: 5},
							MaxInterval:  &duration.Duration{Seconds: 10},
						},
						NumRetries: &wrappers.UInt32Value{Value: 7},
					},
				},
			},
			Timeout: &duration.Duration{Seconds: 100},
			InitialMetadata: []*envoy_config_core_v3.HeaderValue{
				{Key: "E", Value: "F"},
			},
		},
		ProcessingMode: &envoy_ext_proc_v3.ProcessingMode{
			ResponseHeaderMode: envoy_ext_proc_v3.ProcessingMode_SKIP,
		},
		MaxMessageTimeout:  &duration.Duration{Seconds: 10},
		ResponseAttributes: []string{"new1", "new2"},
	}
}

// getExpectedMergedConfig builds the expected config produced by merging the global and listener-level config
func getExpectedMergedConfig() *envoy_ext_proc_v3.ExternalProcessor {
	return &envoy_ext_proc_v3.ExternalProcessor{
		GrpcService: &envoy_config_core_v3.GrpcService{
			TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
					ClusterName: translator.UpstreamToClusterName(&core.ResourceRef{Name: builders.OverrideExtProcUpstreamName, Namespace: "gloo-system"}),
					RetryPolicy: &envoy_config_core_v3.RetryPolicy{
						RetryBackOff: &envoy_config_core_v3.BackoffStrategy{
							BaseInterval: &duration.Duration{Seconds: 5},
							MaxInterval:  &duration.Duration{Seconds: 10},
						},
						NumRetries: &wrappers.UInt32Value{Value: 7},
					},
				},
			},
			Timeout: &duration.Duration{Seconds: 100},
			InitialMetadata: []*envoy_config_core_v3.HeaderValue{
				{Key: "E", Value: "F"},
			},
		},
		FailureModeAllow: true,
		ProcessingMode: &envoy_ext_proc_v3.ProcessingMode{
			ResponseHeaderMode: envoy_ext_proc_v3.ProcessingMode_SKIP,
		},
		MessageTimeout:     &duration.Duration{Seconds: 1},
		MaxMessageTimeout:  &duration.Duration{Seconds: 10},
		RequestAttributes:  []string{"req1", "req2"},
		ResponseAttributes: []string{"new1", "new2"},
	}
}

// getExpectedOverridesRouteConfig builds envoy vhost/route extproc config with the values used in the default ExtProcRouteBuilder.
func getExpectedOverridesRouteConfig() *envoy_ext_proc_v3.ExtProcPerRoute {
	return &envoy_ext_proc_v3.ExtProcPerRoute{
		Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{
			Overrides: &envoy_ext_proc_v3.ExtProcOverrides{
				GrpcService: &envoy_config_core_v3.GrpcService{
					TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
							ClusterName: translator.UpstreamToClusterName(&core.ResourceRef{Name: builders.OverrideExtProcUpstreamName, Namespace: "gloo-system"}),
						},
					},
					InitialMetadata: []*envoy_config_core_v3.HeaderValue{
						{Key: "aaa", Value: "bbb"},
					},
				},
				ProcessingMode: &envoy_ext_proc_v3.ProcessingMode{
					RequestHeaderMode:  envoy_ext_proc_v3.ProcessingMode_SEND,
					ResponseHeaderMode: envoy_ext_proc_v3.ProcessingMode_SKIP,
					RequestBodyMode:    envoy_ext_proc_v3.ProcessingMode_STREAMED,
				},
				AsyncMode:          true,
				RequestAttributes:  []string{"x"},
				ResponseAttributes: []string{"y"},
			},
		},
	}
}

// getExpectedDisabledRouteConfig builds envoy vhost/route extproc config when extproc is disabled
func getExpectedDisabledRouteConfig() *envoy_ext_proc_v3.ExtProcPerRoute {
	return &envoy_ext_proc_v3.ExtProcPerRoute{
		Override: &envoy_ext_proc_v3.ExtProcPerRoute_Disabled{
			Disabled: true,
		},
	}
}

func buildUpstream(name string) *v1.Upstream {
	return &v1.Upstream{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: "gloo-system",
		},
		UpstreamType: &v1.Upstream_Static{
			Static: &static.UpstreamSpec{
				Hosts: []*static.Host{{
					Addr: "extproc-service",
					Port: 1234,
				}},
			},
		},
	}
}

// not sure why comparing the whole object fails, so compare each field individually
func expectEqual(a *envoy_ext_proc_v3.ExternalProcessor, b *envoy_ext_proc_v3.ExternalProcessor) {
	Expect(a.GetGrpcService()).To(Equal(b.GetGrpcService()))
	Expect(a.GetFailureModeAllow()).To(Equal(b.GetFailureModeAllow()))
	Expect(a.GetProcessingMode()).To(Equal(b.GetProcessingMode()))
	Expect(a.GetAsyncMode()).To(Equal(b.GetAsyncMode()))
	Expect(a.GetRequestAttributes()).To(Equal(b.GetRequestAttributes()))
	Expect(a.GetResponseAttributes()).To(Equal(b.GetResponseAttributes()))
	Expect(a.GetMessageTimeout()).To(Equal(b.GetMessageTimeout()))
	Expect(a.GetStatPrefix()).To(Equal(b.GetStatPrefix()))
	Expect(a.GetMutationRules()).To(Equal(b.GetMutationRules()))
	Expect(a.GetMaxMessageTimeout()).To(Equal(b.GetMaxMessageTimeout()))
	// Expect(a.GetDisableClearRouteCache()).To(Equal(b.GetDisableClearRouteCache()))
	// Expect(a.GetForwardRules()).To(Equal(b.GetForwardRules()))
	// Expect(a.GetFilterMetadata()).To(Equal(b.GetFilterMetadata()))
	// Expect(a.GetAllowModeOverride()).To(Equal(b.GetAllowModeOverride()))
}

func expectRouteEqual(a *envoy_ext_proc_v3.ExtProcPerRoute, b *envoy_ext_proc_v3.ExtProcPerRoute) {
	Expect(a.GetDisabled()).To(Equal(b.GetDisabled()))
	Expect(a.GetOverrides()).To(Equal(b.GetOverrides()))
}
