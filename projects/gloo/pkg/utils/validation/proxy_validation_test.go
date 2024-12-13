package validation_test

import (
	"fmt"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	. "github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
)

var _ = Describe("validation utils", func() {

	numListeners := 3
	numVhosts := 5
	numRoutes := 8
	numTcpListeners := 4
	numHttpListeners := 2
	makeHttpProxy := func() *v1.Proxy {
		proxy := &v1.Proxy{}
		for range numListeners {
			httpListener := &v1.HttpListener{}
			proxy.Listeners = append(proxy.Listeners, &v1.Listener{
				ListenerType: &v1.Listener_HttpListener{
					HttpListener: httpListener,
				},
			})

			for range numVhosts {
				vh := &v1.VirtualHost{}
				httpListener.VirtualHosts = append(httpListener.VirtualHosts, vh)

				for range numRoutes {
					vh.Routes = append(vh.Routes, &v1.Route{})
				}
			}
		}
		return proxy
	}
	makeTcpProxy := func() *v1.Proxy {
		proxy := &v1.Proxy{}
		for range numListeners {
			tcpListener := &v1.TcpListener{}
			proxy.Listeners = append(proxy.Listeners, &v1.Listener{
				ListenerType: &v1.Listener_TcpListener{
					TcpListener: tcpListener,
				},
			})

			for range numVhosts {
				vh := &v1.TcpHost{}
				tcpListener.TcpHosts = append(tcpListener.TcpHosts, vh)
			}
		}
		return proxy
	}
	makeHybridProxy := func() *v1.Proxy {
		proxy := &v1.Proxy{}
		for range numListeners {
			hybridListener := &v1.HybridListener{}
			proxy.Listeners = append(proxy.Listeners, &v1.Listener{
				ListenerType: &v1.Listener_HybridListener{
					HybridListener: hybridListener,
				},
			})

			for l := range numTcpListeners {
				tcpListener := &v1.TcpListener{}

				for range numVhosts {
					vh := &v1.TcpHost{}
					tcpListener.TcpHosts = append(tcpListener.TcpHosts, vh)
				}
				hybridListener.MatchedListeners = append(hybridListener.MatchedListeners, &v1.MatchedListener{
					Matcher: &v1.Matcher{
						SourcePrefixRanges: []*v3.CidrRange{
							&v3.CidrRange{
								AddressPrefix: fmt.Sprintf("tcp-%d", l),
							},
						},
					},
					ListenerType: &v1.MatchedListener_TcpListener{
						TcpListener: tcpListener,
					},
				})
			}
			for l := range numHttpListeners {
				httpListener := &v1.HttpListener{}

				for range numVhosts {
					vh := &v1.VirtualHost{}
					httpListener.VirtualHosts = append(httpListener.VirtualHosts, vh)

					for range numRoutes {
						vh.Routes = append(vh.Routes, &v1.Route{})
					}
				}
				hybridListener.MatchedListeners = append(hybridListener.MatchedListeners, &v1.MatchedListener{
					Matcher: &v1.Matcher{
						SourcePrefixRanges: []*v3.CidrRange{
							&v3.CidrRange{
								AddressPrefix: fmt.Sprintf("http-%d", l),
							},
						},
					},
					ListenerType: &v1.MatchedListener_HttpListener{
						HttpListener: httpListener,
					},
				})
			}
		}
		return proxy
	}
	makeInvalidAggregateListenerProxyTcp := func() *v1.Proxy {
		proxy := makeHybridProxy()
		proxy.Listeners = []*v1.Listener{{
			Name:        "aggregate-listener",
			BindAddress: "::",
			BindPort:    defaults.HttpPort,
			ListenerType: &v1.Listener_AggregateListener{
				AggregateListener: &v1.AggregateListener{
					TcpListeners: []*v1.MatchedTcpListener{{
						Matcher: &v1.Matcher{
							SourcePrefixRanges: []*v3.CidrRange{{
								AddressPrefix: "tcp-0",
							}},
						},
						TcpListener: nil,
					}},
					HttpResources: &v1.AggregateListener_HttpResources{},
					HttpFilterChains: []*v1.AggregateListener_HttpFilterChain{{
						Matcher: &v1.Matcher{
							SourcePrefixRanges: []*v3.CidrRange{{
								AddressPrefix: "http-0",
							}},
						},
					}},
				},
			},
		}}
		return proxy
	}

	var _ = Describe("MakeReport", func() {
		It("generates a report which matches an http proxy", func() {

			proxy := makeHttpProxy()

			rpt := MakeReport(proxy)
			Expect(rpt.ListenerReports).To(HaveLen(len(proxy.Listeners)))
			for i := range rpt.ListenerReports {
				vhReports := rpt.ListenerReports[i].GetHttpListenerReport().VirtualHostReports
				Expect(vhReports).To(HaveLen(len(proxy.Listeners[i].GetHttpListener().VirtualHosts)))
				for j := range vhReports {
					Expect(vhReports[i].RouteReports).To(HaveLen(len(proxy.Listeners[i].GetHttpListener().VirtualHosts[j].Routes)))
				}
			}
		})

		It("generates a report which matches a tcp proxy", func() {

			proxy := makeTcpProxy()

			rpt := MakeReport(proxy)
			Expect(rpt.ListenerReports).To(HaveLen(len(proxy.Listeners)))
			for i := range rpt.ListenerReports {
				tcpHostReports := rpt.ListenerReports[i].GetTcpListenerReport().TcpHostReports
				Expect(tcpHostReports).To(HaveLen(len(proxy.Listeners[i].GetTcpListener().GetTcpHosts())))
			}

		})

		It("generates a report which matches a hybrid proxy", func() {

			proxy := makeHybridProxy()

			rpt := MakeReport(proxy)
			Expect(rpt.ListenerReports).To(HaveLen(len(proxy.Listeners)))
			for i := range rpt.ListenerReports {
				hybridListenerReports := rpt.ListenerReports[i].GetHybridListenerReport().GetMatchedListenerReports()
				Expect(hybridListenerReports).NotTo(BeNil())
				for j, matchedListener := range proxy.Listeners[i].GetHybridListener().GetMatchedListeners() {
					matchPrefix, matchNum := "tcp", j
					if j >= numTcpListeners {
						matchPrefix = "http"
						matchNum -= numTcpListeners
					}
					expMatcher := &v1.Matcher{
						SourcePrefixRanges: []*v3.CidrRange{
							&v3.CidrRange{
								AddressPrefix: fmt.Sprintf("%s-%d", matchPrefix, matchNum),
							},
						},
					}
					mlr, ok := hybridListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[i], expMatcher)]
					Expect(ok).To(BeTrue())
					if matchPrefix == "tcp" {
						tcpListenerReport := mlr.GetTcpListenerReport()
						Expect(tcpListenerReport).NotTo(BeNil())
						tcpHostReports := tcpListenerReport.TcpHostReports
						Expect(tcpHostReports).To(HaveLen(len(matchedListener.GetTcpListener().GetTcpHosts())))
					} else {
						httpListenerReport := mlr.GetHttpListenerReport()
						Expect(httpListenerReport).NotTo(BeNil())
						vhReports := httpListenerReport.VirtualHostReports
						Expect(vhReports).To(HaveLen(len(matchedListener.GetHttpListener().GetVirtualHosts())))
						for k := range vhReports {
							Expect(vhReports[i].RouteReports).To(HaveLen(len(matchedListener.GetHttpListener().VirtualHosts[k].Routes)))
						}
					}
				}
			}
		})

		It("properly instantiates TCP reports in aggregate listeners", func() {
			// fixes a crash reported in
			// https://github.com/solo-io/gloo/issues/8405
			proxy := makeInvalidAggregateListenerProxyTcp()
			proxyReports := MakeReport(proxy)
			Expect(proxyReports.GetListenerReports()).To(HaveLen(1))
			tcpListenerReports := proxyReports.GetListenerReports()[0].GetAggregateListenerReport().GetTcpListenerReports()
			Expect(tcpListenerReports).To(HaveLen(1))
		})
	})

	var _ = Describe("GetProxyError", func() {
		It("aggregates the errors at every level for http listener", func() {
			rpt := MakeReport(makeHttpProxy())

			rpt.ListenerReports[1].Errors = append(rpt.ListenerReports[1].Errors,
				&validation.ListenerReport_Error{
					Type:   validation.ListenerReport_Error_BindPortNotUniqueError,
					Reason: "bind port not unique",
				},
			)
			httpListenerReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HttpListenerReport).HttpListenerReport
			httpListenerReport.Errors = append(httpListenerReport.Errors, &validation.HttpListenerReport_Error{
				Type:   validation.HttpListenerReport_Error_ProcessingError,
				Reason: "bad http plugin",
			})

			virtualHostReport := rpt.ListenerReports[0].ListenerTypeReport.(*validation.ListenerReport_HttpListenerReport).HttpListenerReport.VirtualHostReports[2]

			virtualHostReport.Errors = append(virtualHostReport.Errors,
				&validation.VirtualHostReport_Error{
					Type:   validation.VirtualHostReport_Error_DomainsNotUniqueError,
					Reason: "domains not unique",
				},
			)
			routeReport := rpt.ListenerReports[1].ListenerTypeReport.(*validation.ListenerReport_HttpListenerReport).HttpListenerReport.VirtualHostReports[3].RouteReports[2]

			routeReport.Warnings = append(routeReport.Warnings,
				&validation.RouteReport_Warning{
					Type:   validation.RouteReport_Warning_InvalidDestinationWarning,
					Reason: "bad destination",
				},
			)

			err := GetProxyError(rpt)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("VirtualHost Error: DomainsNotUniqueError. Reason: domains not unique; Listener Error: BindPortNotUniqueError. Reason: bind port not unique; HttpListener Error: ProcessingError. Reason: bad http plugin"))
		})
		It("aggregates the errors at every level for aggregate listener", func() {
			proxy := makeInvalidAggregateListenerProxyTcp()
			rpt := MakeReport(proxy)

			rpt.ListenerReports[0].Errors = append(rpt.ListenerReports[0].Errors,
				&validation.ListenerReport_Error{
					Type:   validation.ListenerReport_Error_BindPortNotUniqueError,
					Reason: "bind port not unique",
				},
			)

			tcpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "tcp-0",
					},
				},
			}
			tcpListenerReport := rpt.ListenerReports[0].ListenerTypeReport.(*validation.ListenerReport_AggregateListenerReport).AggregateListenerReport.TcpListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[0], tcpMatcher)]

			// populate the errors - we should hit all cases in
			// getTcpListenerReportErrs with the errors we create here
			tcpListenerReport.Errors = append(tcpListenerReport.Errors,
				&validation.TcpListenerReport_Error{
					Type:   validation.TcpListenerReport_Error_SSLConfigError,
					Reason: "test SSLConfig Error",
				},
			)
			tcpListenerReport.TcpHostReports = append(tcpListenerReport.TcpHostReports,
				&validation.TcpHostReport{
					Errors: []*validation.TcpHostReport_Error{
						{
							Type:   validation.TcpHostReport_Error_InvalidDestinationError,
							Reason: "testing invalid destination error",
						},
						{
							Type:   validation.TcpHostReport_Error_ProcessingError,
							Reason: "testing processing error",
						},
					},
				})

			httpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "http-0",
					},
				},
			}
			httpListenerReport := rpt.ListenerReports[0].ListenerTypeReport.(*validation.ListenerReport_AggregateListenerReport).AggregateListenerReport.HttpListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[0], httpMatcher)]
			httpListenerReport.Errors = append(httpListenerReport.Errors, &validation.HttpListenerReport_Error{
				Type:   validation.HttpListenerReport_Error_ProcessingError,
				Reason: "bad http plugin",
			})

			err := GetProxyError(rpt)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("HttpListener Error: ProcessingError. Reason: bad http plugin"))
			Expect(err.Error()).To(ContainSubstring("Listener Error: BindPortNotUniqueError. Reason: bind port not unique"))
			Expect(err.Error()).To(ContainSubstring("TcpListener Error: SSLConfigError. Reason: test SSLConfig Error"))
			Expect(err.Error()).To(ContainSubstring("TcpHost Error: InvalidDestinationError. Reason: testing invalid destination error"))
			Expect(err.Error()).To(ContainSubstring("TcpHost Error: ProcessingError. Reason: testing processing error"))
		})
		It("aggregates the errors at every level for hybrid listener", func() {
			proxy := makeHybridProxy()
			rpt := MakeReport(proxy)

			rpt.ListenerReports[1].Errors = append(rpt.ListenerReports[1].Errors,
				&validation.ListenerReport_Error{
					Type:   validation.ListenerReport_Error_BindPortNotUniqueError,
					Reason: "bind port not unique",
				},
			)
			tcpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "tcp-0",
					},
				},
			}
			httpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "http-0",
					},
				},
			}

			httpListenerReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HybridListenerReport).HybridListenerReport.MatchedListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[2], httpMatcher)].GetHttpListenerReport()
			httpListenerReport.Errors = append(httpListenerReport.Errors, &validation.HttpListenerReport_Error{
				Type:   validation.HttpListenerReport_Error_ProcessingError,
				Reason: "bad http plugin",
			})

			virtualHostReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HybridListenerReport).HybridListenerReport.MatchedListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[2], httpMatcher)].GetHttpListenerReport().VirtualHostReports[2]

			virtualHostReport.Errors = append(virtualHostReport.Errors,
				&validation.VirtualHostReport_Error{
					Type:   validation.VirtualHostReport_Error_DomainsNotUniqueError,
					Reason: "domains not unique",
				},
			)
			routeReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HybridListenerReport).HybridListenerReport.MatchedListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[2], httpMatcher)].GetHttpListenerReport().VirtualHostReports[3].RouteReports[2]

			routeReport.Warnings = append(routeReport.Warnings,
				&validation.RouteReport_Warning{
					Type:   validation.RouteReport_Warning_InvalidDestinationWarning,
					Reason: "bad destination",
				},
			)

			tcpListenerReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HybridListenerReport).HybridListenerReport.MatchedListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[2], tcpMatcher)].GetTcpListenerReport()
			tcpListenerReport.Errors = append(tcpListenerReport.Errors,
				&validation.TcpListenerReport_Error{
					Type:   validation.TcpListenerReport_Error_SSLConfigError,
					Reason: "test SSLConfig Error",
				},
			)
			tcpListenerReport.TcpHostReports = append(tcpListenerReport.TcpHostReports,
				&validation.TcpHostReport{
					Errors: []*validation.TcpHostReport_Error{
						{
							Type:   validation.TcpHostReport_Error_InvalidDestinationError,
							Reason: "testing invalid destination error",
						},
						{
							Type:   validation.TcpHostReport_Error_ProcessingError,
							Reason: "testing processing error",
						},
					},
				})

			err := GetProxyError(rpt)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("VirtualHost Error: DomainsNotUniqueError. Reason: domains not unique"))
			Expect(err.Error()).To(ContainSubstring("Listener Error: BindPortNotUniqueError. Reason: bind port not unique"))
			Expect(err.Error()).To(ContainSubstring("HttpListener Error: ProcessingError. Reason: bad http plugin"))
			Expect(err.Error()).To(ContainSubstring("TcpListener Error: SSLConfigError. Reason: test SSLConfig Error"))
			Expect(err.Error()).To(ContainSubstring("TcpHost Error: InvalidDestinationError. Reason: testing invalid destination error"))
			Expect(err.Error()).To(ContainSubstring("TcpHost Error: ProcessingError. Reason: testing processing error"))
		})

	})

	var _ = Describe("GetProxyWarning", func() {
		It("aggregates the warnings at every level for http listener", func() {
			rpt := MakeReport(makeHttpProxy())

			rpt.ListenerReports[1].Warnings = append(rpt.ListenerReports[1].Warnings,
				&validation.ListenerReport_Warning{
					Type:   validation.ListenerReport_Warning_Type(validation.ListenerReport_Warning_SSLConfigWarning),
					Reason: "invalid ssl config",
				},
			)
			httpListenerReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HttpListenerReport).HttpListenerReport
			httpListenerReport.Warnings = append(httpListenerReport.Warnings, &validation.HttpListenerReport_Warning{
				Type:   validation.HttpListenerReport_Warning_UnknownWarning,
				Reason: "unknown warning",
			})

			routeReport := rpt.ListenerReports[1].ListenerTypeReport.(*validation.ListenerReport_HttpListenerReport).HttpListenerReport.VirtualHostReports[3].RouteReports[2]

			routeReport.Warnings = append(routeReport.Warnings,
				&validation.RouteReport_Warning{
					Type:   validation.RouteReport_Warning_InvalidDestinationWarning,
					Reason: "bad destination",
				},
			)

			warnings := GetProxyWarning(rpt)
			Expect(warnings).To(HaveLen(3))
			Expect(warnings[0]).To(ContainSubstring("Listener Warning: SSLConfigWarning. Reason: invalid ssl config"))
			Expect(warnings[1]).To(ContainSubstring("Route Warning: InvalidDestinationWarning. Reason: bad destination"))
			Expect(warnings[2]).To(ContainSubstring("HttpListener Warning: UnknownWarning. Reason: unknown warning"))
		})

		It("aggregates the warnings at every level for aggregate listener", func() {
			proxy := makeInvalidAggregateListenerProxyTcp()
			rpt := MakeReport(proxy)

			rpt.ListenerReports[0].Warnings = append(rpt.ListenerReports[0].Warnings,
				&validation.ListenerReport_Warning{
					Type:   validation.ListenerReport_Warning_SSLConfigWarning,
					Reason: "invalid ssl config",
				},
			)

			tcpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "tcp-0",
					},
				},
			}
			tcpListenerReport := rpt.ListenerReports[0].ListenerTypeReport.(*validation.ListenerReport_AggregateListenerReport).AggregateListenerReport.TcpListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[0], tcpMatcher)]

			// populate the warnings - we should hit all cases in
			// getTcpListenerReportWarns with the warnings we create here
			tcpListenerReport.Warnings = append(tcpListenerReport.Warnings,
				&validation.TcpListenerReport_Warning{
					Type:   validation.TcpListenerReport_Warning_UnknownWarning,
					Reason: "unknown warning",
				},
			)
			tcpListenerReport.TcpHostReports = append(tcpListenerReport.TcpHostReports,
				&validation.TcpHostReport{
					Warnings: []*validation.TcpHostReport_Warning{
						{
							Type:   validation.TcpHostReport_Warning_InvalidDestinationWarning,
							Reason: "testing invalid destination warning",
						},
						{
							Type:   validation.TcpHostReport_Warning_InvalidDestinationWarning,
							Reason: "testing invalid destination warning",
						},
					},
				})

			httpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "http-0",
					},
				},
			}
			httpListenerReport := rpt.ListenerReports[0].ListenerTypeReport.(*validation.ListenerReport_AggregateListenerReport).AggregateListenerReport.HttpListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[0], httpMatcher)]
			httpListenerReport.Warnings = append(httpListenerReport.Warnings, &validation.HttpListenerReport_Warning{
				Type:   validation.HttpListenerReport_Warning_InvalidDestinationWarning,
				Reason: "testing invalid destination warning",
			})

			warnings := GetProxyWarning(rpt)
			Expect(warnings).To(HaveLen(5))
			// Since the order of the warnings can't be guaranteed, convert it to a string and ensure
			// every warning has been added
			warningsStr := strings.Join(warnings, " ")
			Expect(warningsStr).To(ContainSubstring("Listener Warning: SSLConfigWarning. Reason: invalid ssl config"))
			Expect(warningsStr).To(ContainSubstring("HttpListener Warning: InvalidDestinationWarning. Reason: testing invalid destination warning"))
			Expect(warningsStr).To(ContainSubstring("TcpListener Warning: UnknownWarning. Reason: unknown warning"))
			Expect(warningsStr).To(ContainSubstring("TcpHost Warning: InvalidDestinationWarning. Reason: testing invalid destination warning"))
			Expect(warningsStr).To(ContainSubstring("TcpHost Warning: InvalidDestinationWarning. Reason: testing invalid destination warning"))
		})

		It("aggregates the warnings at every level for hybrid listener", func() {
			proxy := makeHybridProxy()
			rpt := MakeReport(proxy)

			rpt.ListenerReports[1].Warnings = append(rpt.ListenerReports[1].Warnings,
				&validation.ListenerReport_Warning{
					Type:   validation.ListenerReport_Warning_SSLConfigWarning,
					Reason: "invalid ssl config",
				},
			)
			tcpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "tcp-0",
					},
				},
			}
			httpMatcher := &v1.Matcher{
				SourcePrefixRanges: []*v3.CidrRange{
					&v3.CidrRange{
						AddressPrefix: "http-0",
					},
				},
			}

			httpListenerReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HybridListenerReport).HybridListenerReport.MatchedListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[2], httpMatcher)].GetHttpListenerReport()
			httpListenerReport.Warnings = append(httpListenerReport.Warnings, &validation.HttpListenerReport_Warning{
				Type:   validation.HttpListenerReport_Warning_InvalidDestinationWarning,
				Reason: "invalid destination",
			})

			routeReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HybridListenerReport).HybridListenerReport.MatchedListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[2], httpMatcher)].GetHttpListenerReport().VirtualHostReports[3].RouteReports[2]

			routeReport.Warnings = append(routeReport.Warnings,
				&validation.RouteReport_Warning{
					Type:   validation.RouteReport_Warning_InvalidDestinationWarning,
					Reason: "bad destination",
				},
			)

			tcpListenerReport := rpt.ListenerReports[2].ListenerTypeReport.(*validation.ListenerReport_HybridListenerReport).HybridListenerReport.MatchedListenerReports[utils.MatchedRouteConfigName(proxy.GetListeners()[2], tcpMatcher)].GetTcpListenerReport()
			tcpListenerReport.Warnings = append(tcpListenerReport.Warnings,
				&validation.TcpListenerReport_Warning{
					Type:   validation.TcpListenerReport_Warning_UnknownWarning,
					Reason: "unknown warning",
				},
			)
			tcpListenerReport.TcpHostReports = append(tcpListenerReport.TcpHostReports,
				&validation.TcpHostReport{
					Warnings: []*validation.TcpHostReport_Warning{
						{
							Type:   validation.TcpHostReport_Warning_InvalidDestinationWarning,
							Reason: "testing invalid destination warning",
						},
						{
							Type:   validation.TcpHostReport_Warning_InvalidDestinationWarning,
							Reason: "testing invalid destination warning",
						},
					},
				})

			warnings := GetProxyWarning(rpt)
			Expect(warnings).To(HaveLen(6))
			// Since the order of the warnings can't be guaranteed, convert it to a string and ensure
			// every warning has been added
			warningsStr := strings.Join(warnings, " ")
			Expect(warningsStr).To(ContainSubstring("Listener Warning: SSLConfigWarning. Reason: invalid ssl config"))
			Expect(warningsStr).To(ContainSubstring("TcpListener Warning: UnknownWarning. Reason: unknown warning"))
			Expect(warningsStr).To(ContainSubstring("TcpHost Warning: InvalidDestinationWarning. Reason: testing invalid destination warning"))
			Expect(warningsStr).To(ContainSubstring("TcpHost Warning: InvalidDestinationWarning. Reason: testing invalid destination warning"))
			Expect(warningsStr).To(ContainSubstring("HttpListener Warning: InvalidDestinationWarning. Reason: invalid destination"))
			Expect(warningsStr).To(ContainSubstring("Route Warning: InvalidDestinationWarning. Reason: bad destination"))
		})
	})
})
