package listeneroptions

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/slices"

	"github.com/solo-io/gloo/pkg/utils/statusutils"
	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwquery "github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	lisoptquery "github.com/solo-io/gloo/projects/gateway2/translator/plugins/listeneroptions/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/proxy_protocol"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
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

var _ = Describe("ListenerOptions Plugin", func() {
	Describe("Attaching ListenerOptions via policy attachment", func() {
		var (
			ctx             context.Context
			listenerCtx     *plugins.ListenerContext
			outputListener  *v1.Listener
			expectedOptions *v1.ListenerOptions
			expectedSource  *v1.SourceMetadata_SourceRef
			statusCtx       plugins.StatusContext

			listenerOptionCollection krt.Collection[*solokubev1.ListenerOption]
			listenerOptionClientFull sologatewayv1.ListenerOptionClient
			statusReporter           reporter.StatusReporter

			listenerOptionClient interface {
				Read(namespace, name string, opts clients.ReadOpts) (*sologatewayv1.ListenerOption, error)
			}
		)

		initCollections := func(loptions ...*sologatewayv1.ListenerOption) {
			lokube := slices.Map(loptions, func(lo *sologatewayv1.ListenerOption) *solokubev1.ListenerOption {
				var ret solokubev1.ListenerOption
				ret.ObjectMeta.Name = lo.GetMetadata().GetName()
				ret.ObjectMeta.Namespace = lo.GetMetadata().GetNamespace()
				ret.Spec = *lo
				return &ret
			})
			listenerOptionCollection = krt.NewStaticCollection(nil, lokube)
			for _, lo := range loptions {
				listenerOptionClientFull.Write(lo, clients.WriteOpts{})
			}
		}

		BeforeEach(func() {
			ctx = context.Background()

			listenerCtx = &plugins.ListenerContext{
				Gateway: &gwv1.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gw",
						Namespace: "default",
					},
				},
				GwListener: &gwv1.Listener{
					Name: "test-listener",
				},
			}

			outputListener = &v1.Listener{
				Options: &v1.ListenerOptions{
					ProxyProtocol: &proxy_protocol.ProxyProtocol{},
				},
			}

			expectedOptions = &v1.ListenerOptions{
				// from config
				PerConnectionBufferLimitBytes: &wrapperspb.UInt32Value{
					Value: uint32(419),
				},
				// base
				ProxyProtocol: &proxy_protocol.ProxyProtocol{},
			}

			expectedSource = &v1.SourceMetadata_SourceRef{
				ResourceRef: &core.ResourceRef{
					Name:      "policy",
					Namespace: "default",
				},
			}

			statusCtx = plugins.StatusContext{
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

			resourceClientFactory := &factory.MemoryResourceClientFactory{
				Cache: memory.NewInMemoryResourceCache(),
			}

			var err error
			listenerOptionClientFull, err = sologatewayv1.NewListenerOptionClient(ctx, resourceClientFactory)
			Expect(err).ToNot(HaveOccurred())
			listenerOptionClient = listenerOptionClientFull

			listenerOptionCollection = krt.NewStatic[*solokubev1.ListenerOption](nil, true).AsCollection()
			statusClient := statusutils.GetStatusClientForNamespace("gloo-system")
			statusReporter = reporter.NewReporter(defaults.KubeGatewayReporter, statusClient,
				listenerOptionClientFull.BaseClient())
		})

		When("ListenerOptions exist in the same namespace and are attached correctly", func() {
			It("correctly adds buffer limit", func() {
				initCollections(attachedListenerOptionInternal())
				deps := []client.Object{attachedListenerOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, lisoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, listenerOptionCollection, statusReporter)

				err := plugin.ApplyListenerPlugin(ctx, listenerCtx, outputListener)
				Expect(err).ToNot(HaveOccurred())
				Expect(proto.Equal(outputListener.GetOptions(), expectedOptions)).To(BeTrue())
				Expect(outputListener.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputListener.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				err = plugin.ApplyStatusPlugin(ctx, &statusCtx)
				Expect(err).ToNot(HaveOccurred())

				loobj, err := listenerOptionClient.Read("default", "policy", clients.ReadOpts{Ctx: ctx})
				Expect(err).ToNot(HaveOccurred())
				status := loobj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted), "status should be accepted")
			})
		})

		When("ListenerOptions exist in the same namespace and are attached correctly with section name", func() {
			It("correctly adds buffer limit", func() {
				initCollections(attachedListenerOptionWithSectionNameInternal())
				deps := []client.Object{attachedListenerOptionWithSectionName()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, lisoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, listenerOptionCollection, statusReporter)

				err := plugin.ApplyListenerPlugin(ctx, listenerCtx, outputListener)
				Expect(err).ToNot(HaveOccurred())
				Expect(proto.Equal(outputListener.GetOptions(), expectedOptions)).To(BeTrue())
				Expect(outputListener.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputListener.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				err = plugin.ApplyStatusPlugin(ctx, &statusCtx)
				Expect(err).ToNot(HaveOccurred())

				loobj, err := listenerOptionClient.Read("default", "policy", clients.ReadOpts{Ctx: ctx})
				Expect(err).ToNot(HaveOccurred())
				status := loobj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted), "status should be accepted")
			})
		})

		When("ListenerOptions exist in the same namespace and are attached correctly but omit the namespace in targetRef", func() {
			It("correctly adds buffer limit", func() {
				initCollections(attachedListenerOptionOmitNamespaceInternal())
				deps := []client.Object{attachedListenerOptionOmitNamespace()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, lisoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, listenerOptionCollection, statusReporter)

				err := plugin.ApplyListenerPlugin(ctx, listenerCtx, outputListener)
				Expect(err).ToNot(HaveOccurred())
				Expect(proto.Equal(outputListener.GetOptions(), expectedOptions)).To(BeTrue())
				Expect(outputListener.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputListener.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				err = plugin.ApplyStatusPlugin(ctx, &statusCtx)
				Expect(err).ToNot(HaveOccurred())

				loobj, err := listenerOptionClient.Read("default", "policy", clients.ReadOpts{Ctx: ctx})
				Expect(err).ToNot(HaveOccurred())
				status := loobj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted), "status should be accepted")
			})
		})

		When("ListenerOptions exist in the same namespace but are not attached correctly", func() {
			It("does not add buffer limit", func() {
				initCollections(nonAttachedListenerOptionInternal())
				deps := []client.Object{nonAttachedListenerOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, lisoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, listenerOptionCollection, statusReporter)

				err := plugin.ApplyListenerPlugin(ctx, listenerCtx, outputListener)
				Expect(err).ToNot(HaveOccurred())
				Expect(outputListener.GetOptions().GetPerConnectionBufferLimitBytes()).To(BeNil())
				Expect(outputListener.GetMetadataStatic().GetSources()).To(BeEmpty())

				err = plugin.ApplyStatusPlugin(ctx, &statusCtx)
				Expect(err).ToNot(HaveOccurred())

				loobj, err := listenerOptionClient.Read("default", "bad-policy", clients.ReadOpts{Ctx: ctx})
				Expect(err).ToNot(HaveOccurred())
				Expect(loobj.GetNamespacedStatuses()).To(BeNil())
			})
		})

		When("ListenerOptions exist in a different namespace than the provided listenerCtx", func() {
			It("does not add buffer limit", func() {
				initCollections(attachedListenerOptionInternal())
				deps := []client.Object{attachedListenerOption()}
				listenerCtx.Gateway.SetNamespace("bad-namespace")
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, lisoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, listenerOptionCollection, statusReporter)

				err := plugin.ApplyListenerPlugin(ctx, listenerCtx, outputListener)
				Expect(err).ToNot(HaveOccurred())
				Expect(outputListener.GetOptions().GetPerConnectionBufferLimitBytes()).To(BeNil())
				Expect(outputListener.GetMetadataStatic().GetSources()).To(BeEmpty())

				err = plugin.ApplyStatusPlugin(ctx, &statusCtx)
				Expect(err).ToNot(HaveOccurred())

				loobj, err := listenerOptionClient.Read("default", "policy", clients.ReadOpts{Ctx: ctx})
				Expect(err).ToNot(HaveOccurred())
				Expect(loobj.GetNamespacedStatuses()).To(BeNil())
			})
		})

		When("Multiple ListenerOptions attaching to Gateway", func() {
			It("should mark the non-attached ListenerOption as not attached", func() {
				initCollections(attachedListenerOptionInternal(), unattachedListenerOptionInternal())
				deps := []client.Object{attachedListenerOption(), unattachedListenerOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, lisoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient, listenerOptionCollection, statusReporter)

				err := plugin.ApplyListenerPlugin(ctx, listenerCtx, outputListener)
				Expect(err).ToNot(HaveOccurred())
				Expect(proto.Equal(outputListener.GetOptions(), expectedOptions)).To(BeTrue())
				Expect(outputListener.GetMetadataStatic().GetSources()).To(HaveLen(1))
				Expect(proto.Equal(outputListener.GetMetadataStatic().GetSources()[0], expectedSource)).To(BeTrue())

				err = plugin.ApplyStatusPlugin(ctx, &statusCtx)
				Expect(err).ToNot(HaveOccurred())

				loobj, err := listenerOptionClient.Read("default", "policy", clients.ReadOpts{Ctx: ctx})
				Expect(err).ToNot(HaveOccurred())
				status := loobj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Accepted), "status should be accepted")

				loobj, err = listenerOptionClient.Read("default", "unattached-policy", clients.ReadOpts{Ctx: ctx})
				Expect(err).ToNot(HaveOccurred())
				status = loobj.GetNamespacedStatuses().Statuses["gloo-system"]
				Expect(status.State).To(Equal(core.Status_Warning), "status should be warning")
				Expect(status.Reason).To(ContainSubstring("istenerOption 'default/unattached-policy' not attached to Gateway 'default/gw' due to higher priority ListenerOption 'default/policy'"))
			})
		})
	})
})

func attachedListenerOptionInternal() *sologatewayv1.ListenerOption {
	return &sologatewayv1.ListenerOption{
		Metadata: &core.Metadata{
			Name:      "policy",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{
			{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.GatewayKind,
				Name:      "gw",
				Namespace: wrapperspb.String("default"),
			},
		},
		Options: &v1.ListenerOptions{
			PerConnectionBufferLimitBytes: &wrapperspb.UInt32Value{
				Value: uint32(419),
			},
		},
	}
}

func attachedListenerOption() *solokubev1.ListenerOption {
	return &solokubev1.ListenerOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.ListenerOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "default",
		},
		Spec: *attachedListenerOptionInternal(),
	}
}

func attachedListenerOptionWithSectionNameInternal() *sologatewayv1.ListenerOption {
	return &sologatewayv1.ListenerOption{
		Metadata: &core.Metadata{
			Name:      "policy",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{
			{
				Group:       gwv1.GroupVersion.Group,
				Kind:        wellknown.GatewayKind,
				Name:        "gw",
				Namespace:   wrapperspb.String("default"),
				SectionName: wrapperspb.String("test-listener"),
			},
		},
		Options: &v1.ListenerOptions{
			PerConnectionBufferLimitBytes: &wrapperspb.UInt32Value{
				Value: uint32(419),
			},
		},
	}
}

func attachedListenerOptionWithSectionName() *solokubev1.ListenerOption {
	return &solokubev1.ListenerOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.ListenerOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "default",
		},
		Spec: *attachedListenerOptionWithSectionNameInternal(),
	}
}

func attachedListenerOptionOmitNamespaceInternal() *sologatewayv1.ListenerOption {
	return &sologatewayv1.ListenerOption{
		Metadata: &core.Metadata{
			Name:      "policy",
			Namespace: "default",
		},
		TargetRefs: []*corev1.PolicyTargetReferenceWithSectionName{
			{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.GatewayKind,
				Name:      "gw",
				Namespace: nil,
			},
		},
		Options: &v1.ListenerOptions{
			PerConnectionBufferLimitBytes: &wrapperspb.UInt32Value{
				Value: uint32(419),
			},
		},
	}
}

func attachedListenerOptionOmitNamespace() *solokubev1.ListenerOption {
	return &solokubev1.ListenerOption{
		TypeMeta: metav1.TypeMeta{
			Kind: sologatewayv1.ListenerOptionGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "default",
		},
		Spec: *attachedListenerOptionOmitNamespaceInternal(),
	}
}

func nonAttachedListenerOptionInternal() *sologatewayv1.ListenerOption {
	listOpt := attachedListenerOptionInternal()
	listOpt.Metadata.Name = "bad-policy"
	listOpt.TargetRefs[0].Name = "bad-gw"
	return listOpt
}

func nonAttachedListenerOption() *solokubev1.ListenerOption {
	listOpt := attachedListenerOption()
	listOpt.ObjectMeta.Name = "bad-policy"
	listOpt.Spec = *nonAttachedListenerOptionInternal()
	return listOpt
}

func unattachedListenerOptionInternal() *sologatewayv1.ListenerOption {
	listOpt := attachedListenerOptionInternal()
	listOpt.Metadata.Name = "unattached-policy"
	return listOpt
}

func unattachedListenerOption() *solokubev1.ListenerOption {
	listOpt := attachedListenerOption()
	listOpt.ObjectMeta.Name = "unattached-policy"
	listOpt.Spec = *unattachedListenerOptionInternal()
	return listOpt
}
