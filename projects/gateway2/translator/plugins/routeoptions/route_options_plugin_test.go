package routeoptions

import (
	"context"
	"errors"
	"time"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/slices"

	"github.com/solo-io/gloo/pkg/utils/statusutils"
	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwquery "github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	rtoptquery "github.com/solo-io/gloo/projects/gateway2/translator/plugins/routeoptions/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/faultinjection"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = Describe("RouteOptionsPlugin", func() {
	var (
		ctx                   context.Context
		cancel                context.CancelFunc
		routeOptionCollection krt.Collection[*solokubev1.RouteOption]
		statusReporter        reporter.StatusReporter
		routeOptionClient     interface {
			Read(namespace, name string, opts clients.ReadOpts) (*sologatewayv1.RouteOption, error)
		}
		routeOptionClientFull sologatewayv1.RouteOptionClient
		statusCtx             *plugins.StatusContext
	)

	initCollections := func(rtoptions ...*sologatewayv1.RouteOption) {
		rtkube := slices.Map(rtoptions, func(rto *sologatewayv1.RouteOption) *solokubev1.RouteOption {
			var ret solokubev1.RouteOption
			ret.ObjectMeta.Name = rto.GetMetadata().GetName()
			ret.ObjectMeta.Namespace = rto.GetMetadata().GetNamespace()
			ret.Spec = *rto
			return &ret
		})
		routeOptionCollection = krt.NewStaticCollection(rtkube)
		for _, rto := range rtoptions {
			routeOptionClientFull.Write(rto, clients.WriteOpts{})
		}
	}

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		resourceClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}

		routeOptionClientFull, _ = sologatewayv1.NewRouteOptionClient(ctx, resourceClientFactory)
		routeOptionClient = routeOptionClientFull

		routeOptionCollection = krt.NewStatic[*solokubev1.RouteOption](nil, true).AsCollection()
		statusClient := statusutils.GetStatusClientForNamespace("gloo-system")
		statusReporter = reporter.NewReporter(defaults.KubeGatewayReporter, statusClient, routeOptionClientFull.BaseClient())
	})

	AfterEach(func() {
		cancel()
	})

	When("applying RouteOptions as Filter", func() {
		It("applies fault injecton RouteOptions directly from resource to output route", func() {
			deps := []client.Object{routeOption()}
			fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
			gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
			plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

			rtCtx := &plugins.RouteContext{
				HTTPRoute: &gwv1.HTTPRoute{},
				Rule: &gwv1.HTTPRouteRule{
					Filters: []gwv1.HTTPRouteFilter{{
						Type: gwv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gwv1.LocalObjectReference{
							Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
							Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
							Name:  "filter-policy",
						},
					}},
				},
			}

			outputRoute := &v1.Route{
				Options: &v1.RouteOptions{},
			}
			plugin.ApplyRoutePlugin(context.Background(), rtCtx, outputRoute)

			expectedRoute := &v1.Route{
				Options: &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 1.00,
							HttpStatus: 500,
						},
					},
				},
			}
			// Expect(proto.Equal(outputRoute, expectedRoute)).To(BeTrue())
			// proto.Equal on the top-level route object doesn't work
			Expect(proto.Equal(outputRoute.GetOptions(), expectedRoute.GetOptions())).To(BeTrue())
		})

		It("reports an error and does not apply any RouteOptions when the referenced obj doesn't exist", func() {
			deps := []client.Object{}
			fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
			gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
			plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

			route := routeWithFilter()
			reportsMap := reports.NewReportMap()
			reporter := reports.NewReporter(&reportsMap)
			parentRefReporter := reporter.Route(route).ParentRef(parentRef())

			rtCtx := &plugins.RouteContext{
				HTTPRoute: route,
				Rule:      routeRuleWithExtRef(),
				Reporter:  parentRefReporter,
			}

			outputRoute := &v1.Route{
				Options: &v1.RouteOptions{},
			}
			err := plugin.ApplyRoutePlugin(context.Background(), rtCtx, outputRoute)

			Expect(err).To(HaveOccurred())
			Expect(proto.Equal(outputRoute.GetOptions(), &v1.RouteOptions{})).To(BeTrue())
		})
	})

	Describe("Attaching RouteOptions via policy attachemnt", func() {
		When("RouteOptions exist in the same namespace and are attached correctly", func() {
			It("correctly adds faultinjection", func() {
				initCollections(attachedInternal())
				deps := []client.Object{attachedRouteOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					HTTPRoute: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "route",
							Namespace: "default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				expectedOptions := &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 4.19,
							HttpStatus: 500,
						},
					},
				}
				Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())

				expectedSource := &v1.SourceMetadata_SourceRef{
					ResourceRef: &core.ResourceRef{
						Name:      "policy",
						Namespace: "default",
					},
					ResourceKind: "RouteOption",
				}

				Expect(outputRoute.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputRoute.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				px := &v1.Proxy{}
				statusCtx := plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: px,
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &validation.ProxyReport{},
								ResourceReports: reporter.ResourceReports{},
							},
						},
					},
				}

				plugin.ApplyStatusPlugin(ctx, &statusCtx)

				robj, _ := routeOptionClient.Read("default", "policy", clients.ReadOpts{Ctx: ctx})
				status := robj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted))
			})
		})

		When("RouteOptions exist in the same namespace and are attached correctly but omit the namespace in targetRef", func() {
			It("correctly adds faultinjection", func() {
				initCollections(attachedOmitNamespaceInternal())
				deps := []client.Object{attachedRouteOptionOmitNamespace()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					HTTPRoute: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "route",
							Namespace: "default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				expectedOptions := &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 4.19,
							HttpStatus: 500,
						},
					},
				}
				Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())

				expectedSource := &v1.SourceMetadata_SourceRef{
					ResourceRef: &core.ResourceRef{
						Name:      "policy",
						Namespace: "default",
					},
					ResourceKind: "RouteOption",
				}

				Expect(outputRoute.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputRoute.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				px := &v1.Proxy{}
				statusCtx := plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: px,
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &validation.ProxyReport{},
								ResourceReports: reporter.ResourceReports{},
							},
						},
					},
				}

				plugin.ApplyStatusPlugin(ctx, &statusCtx)

				robj, _ := routeOptionClient.Read("default", "policy", clients.ReadOpts{Ctx: ctx})
				status := robj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted))
			})
		})

		Context("There is an error reading the RouteOptions", func() {
			BeforeEach(func() {
				statusCtx = &plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: &v1.Proxy{},
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &validation.ProxyReport{},
								ResourceReports: reporter.ResourceReports{},
							},
						},
					},
				}
			})

			When("The RouteOption has a TargetRef", func() {
				It("errors out", func() {
					deps := []client.Object{attachedRouteOption()}
					fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
					gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
					plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

					ctx := context.Background()
					routeCtx := &plugins.RouteContext{
						HTTPRoute: &gwv1.HTTPRoute{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "route",
								Namespace: "default",
							},
						},
					}

					outputRoute := &v1.Route{
						Options: &v1.RouteOptions{},
					}
					plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

					err := plugin.ApplyStatusPlugin(ctx, statusCtx)
					Expect(err).To(MatchError(ContainSubstring(ReadingRouteOptionErrStr)))
				})
			})

			When("The HTTPRoute has an ExtensionRef", func() {
				It("errors out", func() {
					deps := []client.Object{routeOption()}
					fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
					gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
					plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

					rtCtx := &plugins.RouteContext{
						HTTPRoute: &gwv1.HTTPRoute{},
						Rule: &gwv1.HTTPRouteRule{
							Filters: []gwv1.HTTPRouteFilter{{
								Type: gwv1.HTTPRouteFilterExtensionRef,
								ExtensionRef: &gwv1.LocalObjectReference{
									Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
									Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
									Name:  "filter-policy",
								},
							}},
						},
					}

					outputRoute := &v1.Route{
						Options: &v1.RouteOptions{},
					}
					plugin.ApplyRoutePlugin(context.Background(), rtCtx, outputRoute)

					err := plugin.ApplyStatusPlugin(ctx, statusCtx)
					Expect(err).To(MatchError(ContainSubstring(ReadingRouteOptionErrStr)))
				})
			})
		})

		When("Two RouteOptions are attached correctly with different creation timestamps", func() {
			It("correctly adds faultinjection from the earliest created object", func() {
				initCollections(attachedInternal(), attachedBeforeInternal())
				deps := []client.Object{attachedRouteOption(), attachedRouteOptionBefore()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					HTTPRoute: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "route",
							Namespace: "default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				expectedOptions := &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 6.55,
							HttpStatus: 500,
						},
					},
				}
				Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())

				expectedSource := &v1.SourceMetadata_SourceRef{
					ResourceRef: &core.ResourceRef{
						Name:      "policy-older",
						Namespace: "default",
					},
					ResourceKind: "RouteOption",
				}

				Expect(outputRoute.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputRoute.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				px := &v1.Proxy{}
				statusCtx := plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: px,
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &validation.ProxyReport{},
								ResourceReports: reporter.ResourceReports{},
							},
						},
					},
				}

				plugin.ApplyStatusPlugin(ctx, &statusCtx)

				robj, _ := routeOptionClient.Read("default", "policy-older", clients.ReadOpts{Ctx: ctx})
				status := robj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted))
			})
		})

		When("Multiple RouteOptions using targetRef attach to the same route", func() {
			It("correctly merges in priority order from oldest to newest", func() {
				first := newBaseRouteOption("first")

				second := first.Clone().(*sologatewayv1.RouteOption)
				second.Metadata.Name = "second"
				second.Options.Faults.Abort.Percentage = 5
				second.Options.PrefixRewrite = &wrapperspb.StringValue{Value: "/prefix2"}

				third := second.Clone().(*sologatewayv1.RouteOption)
				third.Metadata.Name = "third"
				third.Options.PrefixRewrite = &wrapperspb.StringValue{Value: "/prefix3"}
				third.Options.IdleTimeout = &durationpb.Duration{Seconds: 5}

				initCollections(first, second, third)

				firstOpt := attachedRouteOptionAfterT("first", 0, first)
				secondOpt := attachedRouteOptionAfterT("second", 1*time.Hour, second)
				thirdOpt := attachedRouteOptionAfterT("third", 2*time.Hour, third)

				deps := []client.Object{firstOpt, secondOpt, thirdOpt}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					HTTPRoute: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "route",
							Namespace: "default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				// Confirm that the merged values are as expected
				Expect(outputRoute.GetOptions().GetFaults().GetAbort().GetPercentage()).To(BeNumerically("==", first.GetOptions().GetFaults().GetAbort().GetPercentage()))
				Expect(outputRoute.GetOptions().GetPrefixRewrite().GetValue()).To(Equal(second.GetOptions().GetPrefixRewrite().GetValue()))
				Expect(outputRoute.GetOptions().GetIdleTimeout()).To(Equal(third.GetOptions().GetIdleTimeout()))

				Expect(outputRoute.GetMetadataStatic().GetSources()).To(HaveLen(3))
				Expect(outputRoute.GetMetadataStatic().GetSources()[0].GetResourceRef().GetName()).To(Equal("first"))
				Expect(outputRoute.GetMetadataStatic().GetSources()[1].GetResourceRef().GetName()).To(Equal("second"))
				Expect(outputRoute.GetMetadataStatic().GetSources()[2].GetResourceRef().GetName()).To(Equal("third"))

				px := &v1.Proxy{}
				statusCtx := plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: px,
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &validation.ProxyReport{},
								ResourceReports: reporter.ResourceReports{},
							},
						},
					},
				}

				plugin.ApplyStatusPlugin(ctx, &statusCtx)

				firstObj, _ := routeOptionClient.Read("default", "first", clients.ReadOpts{Ctx: ctx})
				status := firstObj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted))

				secondObj, _ := routeOptionClient.Read("default", "second", clients.ReadOpts{Ctx: ctx})
				status = secondObj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted))

				thirdObj, _ := routeOptionClient.Read("default", "third", clients.ReadOpts{Ctx: ctx})
				status = thirdObj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted))
			})
		})

		When("RouteOptions exist in the same namespace but are not attached correctly", func() {
			It("does not add faultinjection", func() {
				initCollections(nonAttachedInternal())
				deps := []client.Object{nonAttachedRouteOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					HTTPRoute: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "route",
							Namespace: "default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				Expect(proto.Equal(outputRoute.GetOptions(), &v1.RouteOptions{})).To(BeTrue())

				Expect(outputRoute.GetMetadataStatic().GetSources()).To(BeEmpty())

				px := &v1.Proxy{}
				statusCtx := plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: px,
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &validation.ProxyReport{},
								ResourceReports: reporter.ResourceReports{},
							},
						},
					},
				}

				plugin.ApplyStatusPlugin(ctx, &statusCtx)

				robj, _ := routeOptionClient.Read("default", "bad-policy", clients.ReadOpts{Ctx: ctx})
				Expect(robj.GetNamespacedStatuses()).To(BeNil())
			})
		})

		When("RouteOptions exist in a different namespace than the provided routeCtx", func() {
			It("does not add faultinjection", func() {
				initCollections(attachedInternal())
				deps := []client.Object{attachedRouteOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					HTTPRoute: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "route",
							Namespace: "non-default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				Expect(proto.Equal(outputRoute.GetOptions(), &v1.RouteOptions{})).To(BeTrue())

				Expect(outputRoute.GetMetadataStatic().GetSources()).To(BeEmpty())

				px := &v1.Proxy{}
				statusCtx := plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: px,
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &validation.ProxyReport{},
								ResourceReports: reporter.ResourceReports{},
							},
						},
					},
				}

				plugin.ApplyStatusPlugin(ctx, &statusCtx)

				robj, _ := routeOptionClient.Read("default", "bad-policy", clients.ReadOpts{Ctx: ctx})
				Expect(robj.GetNamespacedStatuses()).To(BeNil())
			})
		})

		When("RouteOptions exist in the same namespace and are attached correctly but have processing errors during xds translation", func() {
			It("propagates faultinjection config but reports the processing error on resource status", func() {
				initCollections(attachedInvalidInternal())
				deps := []client.Object{attachedInvalidRouteOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					HTTPRoute: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "route",
							Namespace: "default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				expectedOptions := &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 4.19,
						},
					},
				}
				Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())

				expectedSource := &v1.SourceMetadata_SourceRef{
					ResourceRef: &core.ResourceRef{
						Name:      "invalid-policy",
						Namespace: "default",
					},
					ResourceKind: "RouteOption",
				}

				Expect(outputRoute.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputRoute.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				px := &v1.Proxy{}

				proxyResourceReport := reporter.ResourceReports{}
				proxyResourceReport.AddError(px, errors.New("route processing error"))

				proxyValidationReport := validation.ProxyReport{}
				proxyValidationReport.ListenerReports = []*validation.ListenerReport{}
				// this is ugly; should be replaced by an actual report from translation
				// but should that go in a higher-level test?
				// e.g. an "integration" test for route translation / one that isn't manually running the plugin?
				proxyValidationReport.ListenerReports = append(proxyValidationReport.ListenerReports, &validation.ListenerReport{
					ListenerTypeReport: &validation.ListenerReport_AggregateListenerReport{
						AggregateListenerReport: &validation.AggregateListenerReport{
							HttpListenerReports: map[string]*validation.HttpListenerReport{
								"test": {
									VirtualHostReports: []*validation.VirtualHostReport{
										{
											RouteReports: []*validation.RouteReport{{
												Errors: []*validation.RouteReport_Error{
													{
														Type:   validation.RouteReport_Error_ProcessingError,
														Reason: "route processing error",
														Metadata: &v1.SourceMetadata{
															Sources: []*v1.SourceMetadata_SourceRef{{
																ResourceRef: &core.ResourceRef{
																	Name:      "invalid-policy",
																	Namespace: "default",
																},
																ResourceKind: sologatewayv1.RouteOptionGVK.Kind,
															}},
														},
													},
												},
											}},
										},
									},
								},
							},
						},
					},
				})

				statusCtx := plugins.StatusContext{
					ProxiesWithReports: []translatorutils.ProxyWithReports{
						{
							Proxy: px,
							Reports: translatorutils.TranslationReports{
								ProxyReport:     &proxyValidationReport,
								ResourceReports: proxyResourceReport,
							},
						},
					},
				}

				plugin.ApplyStatusPlugin(ctx, &statusCtx)

				robj, _ := routeOptionClient.Read("default", "invalid-policy", clients.ReadOpts{Ctx: ctx})
				status := robj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Rejected))
			})
		})
	})

	Describe("HTTPRoute with RouteOptions filter AND attached RouteOptions", func() {
		It("Only applies RouteOptions from filter", func() {
			initCollections(attachedInternal())
			deps := []client.Object{routeOption(), attachedRouteOption()}
			fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
			gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
			plugin := NewPlugin(gwQueries, fakeClient, routeOptionCollection, statusReporter)

			ctx := context.Background()
			routeCtx := &plugins.RouteContext{
				HTTPRoute: routeWithFilter(),
				Rule:      routeRuleWithExtRef(),
			}

			outputRoute := &v1.Route{
				Options: &v1.RouteOptions{},
			}
			plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

			expectedOptions := &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 1.00,
						HttpStatus: 500,
					},
				},
			}
			Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())
		})
	})
})

var _ = DescribeTable("mergeOptionsForRoute",
	func(route *gwv1.HTTPRoute, dst, src *v1.RouteOptions, expectedOptions *v1.RouteOptions, expectedResult glooutils.OptionsMergeResult) {
		mergedOptions, result := mergeOptionsForRoute(context.TODO(), route, dst, src)
		Expect(cmp.Diff(mergedOptions, expectedOptions, protocmp.Transform())).To(BeEmpty())
		Expect(result).To(Equal(expectedResult))
	},
	Entry("prefer dst options by default",
		&gwv1.HTTPRoute{},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
		},
		&v1.RouteOptions{
			PrefixRewrite: &wrapperspb.StringValue{Value: "/prefix"},
			// Faults will be ignored because it is set in dst
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
			PrefixRewrite: &wrapperspb.StringValue{Value: "/prefix"},
		},
		glooutils.OptionsMergedPartial,
	),
	Entry("override dst options with annotation: full override",
		&gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{policyOverrideAnnotation: "*"},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
		},
		&v1.RouteOptions{
			PrefixRewrite: &wrapperspb.StringValue{Value: "/prefix"},
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
			PrefixRewrite: &wrapperspb.StringValue{Value: "/prefix"},
		},
		glooutils.OptionsMergedFull,
	),
	Entry("override dst options with annotation: partial override",
		&gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{policyOverrideAnnotation: "*"},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
			Timeout: durationpb.New(5 * time.Second),
		},
		&v1.RouteOptions{
			PrefixRewrite: &wrapperspb.StringValue{Value: "/prefix"},
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
			PrefixRewrite: &wrapperspb.StringValue{Value: "/prefix"},
			Timeout:       durationpb.New(5 * time.Second),
		},
		glooutils.OptionsMergedPartial,
	),
	Entry("override dst options with annotation: no override",
		&gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{policyOverrideAnnotation: "*"},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
		},
		&v1.RouteOptions{},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
		},
		glooutils.OptionsMergedNone,
	),
	Entry("override dst options with annotation: specific fields",
		&gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{policyOverrideAnnotation: "faults,timeout"},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
			Timeout:       durationpb.New(5 * time.Second),
			PrefixRewrite: &wrapperspb.StringValue{Value: "/dst"},
		},
		&v1.RouteOptions{
			PrefixRewrite: &wrapperspb.StringValue{Value: "/src"},
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
			Timeout: durationpb.New(10 * time.Second),
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
			PrefixRewrite: &wrapperspb.StringValue{Value: "/dst"},
			Timeout:       durationpb.New(10 * time.Second),
		},
		glooutils.OptionsMergedFull,
	),
	Entry("override and augment dst options with annotation: specific fields",
		&gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{policyOverrideAnnotation: "faults,timeout"},
			},
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 500,
				},
			},
			Timeout: durationpb.New(5 * time.Second),
		},
		&v1.RouteOptions{
			PrefixRewrite: &wrapperspb.StringValue{Value: "/src"},
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
			Timeout: durationpb.New(10 * time.Second),
		},
		&v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					HttpStatus: 100,
				},
			},
			PrefixRewrite: &wrapperspb.StringValue{Value: "/src"},
			Timeout:       durationpb.New(10 * time.Second),
		},
		glooutils.OptionsMergedFull,
	),
)

func routeOption() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "filter-policy",
			Namespace: "default",
		},
		Spec: sologatewayv1.RouteOption{
			Options: &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 1.00,
						HttpStatus: 500,
					},
				},
			},
		},
	}
}

func routeRuleWithExtRef() *gwv1.HTTPRouteRule {
	return &gwv1.HTTPRouteRule{
		Filters: []gwv1.HTTPRouteFilter{
			{
				Type: "ExtensionRef",
				ExtensionRef: &gwv1.LocalObjectReference{
					Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
					Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
					Name:  "filter-policy",
				},
			},
		},
	}
}

func parentRef() *gwv1.ParentReference {
	return &gwv1.ParentReference{
		Name: "my-gw",
	}
}

func route() *gwv1.HTTPRoute {
	return &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route",
			Namespace: "default",
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{
				ParentRefs: []gwv1.ParentReference{
					*parentRef(),
				},
			},
		},
	}
}

func routeWithFilter() *gwv1.HTTPRoute {
	rwf := route()
	rwf.Spec.Rules = []gwv1.HTTPRouteRule{
		*routeRuleWithExtRef(),
	}
	return rwf
}

func newBaseRouteOption(name string) *sologatewayv1.RouteOption {
	return &sologatewayv1.RouteOption{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReference{
			{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.HTTPRouteKind,
				Name:      "route",
				Namespace: wrapperspb.String("default"),
			},
		},
		Options: &v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					Percentage: 4.19,
					HttpStatus: 500,
				},
			},
		},
	}
}

func attachedInternal() *sologatewayv1.RouteOption {
	return &sologatewayv1.RouteOption{
		Metadata: &core.Metadata{
			Name:      "policy",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReference{
			{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.HTTPRouteKind,
				Name:      "route",
				Namespace: wrapperspb.String("default"),
			},
		},
		Options: &v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					Percentage: 4.19,
					HttpStatus: 500,
				},
			},
		},
	}
}

func attachedRouteOption() *solokubev1.RouteOption {
	now := metav1.Now()
	return &solokubev1.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.RouteOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "policy",
			Namespace:         "default",
			CreationTimestamp: now,
		},
		Spec: *attachedInternal(),
	}
}

func attachedBeforeInternal() *sologatewayv1.RouteOption {
	return &sologatewayv1.RouteOption{
		Metadata: &core.Metadata{
			Name:      "policy-older",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReference{
			{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.HTTPRouteKind,
				Name:      "route",
				Namespace: wrapperspb.String("default"),
			},
		},
		Options: &v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					Percentage: 6.55,
					HttpStatus: 500,
				},
			},
		},
	}
}

func attachedRouteOptionBefore() *solokubev1.RouteOption {
	anHourAgo := metav1.NewTime(time.Now().Add(-1 * time.Hour))
	return &solokubev1.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.RouteOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "policy-older",
			Namespace:         "default",
			CreationTimestamp: anHourAgo,
		},
		Spec: *attachedBeforeInternal(),
	}
}

func attachedRouteOptionAfterT(name string, d time.Duration, spec *sologatewayv1.RouteOption) *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.RouteOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         "default",
			CreationTimestamp: metav1.NewTime(time.Now().Add(d)),
		},
		Spec: *spec,
	}
}

func attachedOmitNamespaceInternal() *sologatewayv1.RouteOption {
	return &sologatewayv1.RouteOption{
		Metadata: &core.Metadata{
			Name:      "policy",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReference{
			{
				Group: gwv1.GroupVersion.Group,
				Kind:  wellknown.HTTPRouteKind,
				Name:  "route",
			},
		},
		Options: &v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					Percentage: 4.19,
					HttpStatus: 500,
				},
			},
		},
	}
}

func attachedRouteOptionOmitNamespace() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.RouteOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "default",
		},
		Spec: *attachedOmitNamespaceInternal(),
	}
}

func nonAttachedInternal() *sologatewayv1.RouteOption {
	return &sologatewayv1.RouteOption{
		Metadata: &core.Metadata{
			Name:      "bad-policy",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReference{
			{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.HTTPRouteKind,
				Name:      "bad-route",
				Namespace: wrapperspb.String("default"),
			},
		},
		Options: &v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					Percentage: 1.00,
					HttpStatus: 500,
				},
			},
		},
	}
}

func nonAttachedRouteOption() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.RouteOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bad-policy",
			Namespace: "default",
		},
		Spec: *nonAttachedInternal(),
	}
}

func attachedInvalidInternal() *sologatewayv1.RouteOption {
	return &sologatewayv1.RouteOption{
		Metadata: &core.Metadata{
			Name:      "invalid-policy",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReference{
			{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.HTTPRouteKind,
				Name:      "route",
				Namespace: wrapperspb.String("default"),
			},
		},
		Options: &v1.RouteOptions{
			Faults: &faultinjection.RouteFaults{
				Abort: &faultinjection.RouteAbort{
					Percentage: 4.19,
				},
			},
		},
	}
}

func attachedInvalidRouteOption() *solokubev1.RouteOption {
	now := metav1.Now()
	return &solokubev1.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.RouteOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "invalid-policy",
			Namespace:         "default",
			CreationTimestamp: now,
		},
		Spec: *attachedInvalidInternal(),
	}
}
