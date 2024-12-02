package mirror_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway2/query/mocks"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/mirror"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestSingleMirror(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	queries := mocks.NewMockGatewayQueries(ctrl)
	g.Expect(queries).ToNot(gomega.BeNil())

	filter := gwv1.HTTPRouteFilter{
		Type: gwv1.HTTPRouteFilterRequestMirror,
		RequestMirror: &gwv1.HTTPRequestMirrorFilter{
			BackendRef: gwv1.BackendObjectReference{
				Name: "foo",
				Port: ptr.To(gwv1.PortNumber(8080)),
			},
		},
	}
	rt := &gwv1.HTTPRoute{}
	routeCtx := &plugins.RouteContext{
		HTTPRoute: rt,
		Rule: &gwv1.HTTPRouteRule{
			Filters: []gwv1.HTTPRouteFilter{
				filter,
			},
		},
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
	}

	queries.EXPECT().ObjToFrom(rt).Return(nil)
	queries.EXPECT().GetBackendForRef(context.Background(), gomock.Any(), &filter.RequestMirror.BackendRef).Return(svc, nil)
	plugin := mirror.NewPlugin(queries)
	outputRoute := &v1.Route{
		Action:  &v1.Route_RouteAction{},
		Options: &v1.RouteOptions{},
	}
	plugin.ApplyRoutePlugin(context.Background(), routeCtx, outputRoute)

	shadowing := outputRoute.GetOptions().GetShadowing()
	g.Expect(shadowing).ToNot(gomega.BeNil())
	g.Expect(shadowing.Upstream).ToNot(gomega.BeNil())
	g.Expect(shadowing.Upstream.Name).To(gomega.Equal("kube-svc:bar-foo-8080"))
	g.Expect(shadowing.Upstream.Namespace).To(gomega.Equal("bar"))
	g.Expect(shadowing.Percentage).To(gomega.Equal(float32(100.0)))
}

func TestUpstreamMirror(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	queries := mocks.NewMockGatewayQueries(ctrl)

	backendKind := gwv1.Kind(v1.UpstreamGVK.Kind)
	backendGroup := gwv1.Group(v1.UpstreamGVK.Group)
	filter := gwv1.HTTPRouteFilter{
		Type: gwv1.HTTPRouteFilterRequestMirror,
		RequestMirror: &gwv1.HTTPRequestMirrorFilter{
			BackendRef: gwv1.BackendObjectReference{
				Name:  "foo",
				Port:  ptr.To(gwv1.PortNumber(8080)),
				Kind:  &backendKind,
				Group: &backendGroup,
			},
		},
	}
	rt := &gwv1.HTTPRoute{}
	routeCtx := &plugins.RouteContext{
		HTTPRoute: rt,
		Rule: &gwv1.HTTPRouteRule{
			Filters: []gwv1.HTTPRouteFilter{
				filter,
			},
		},
	}
	svc := &gloov1.Upstream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
	}

	queries.EXPECT().ObjToFrom(rt).Return(nil)
	queries.EXPECT().GetBackendForRef(context.Background(), gomock.Any(), &filter.RequestMirror.BackendRef).Return(svc, nil)
	plugin := mirror.NewPlugin(queries)
	outputRoute := &v1.Route{
		Action:  &v1.Route_RouteAction{},
		Options: &v1.RouteOptions{},
	}
	plugin.ApplyRoutePlugin(context.Background(), routeCtx, outputRoute)

	shadowing := outputRoute.GetOptions().GetShadowing()
	g.Expect(shadowing).ToNot(gomega.BeNil())
	g.Expect(shadowing.Upstream).ToNot(gomega.BeNil())
	g.Expect(shadowing.Upstream.Name).To(gomega.Equal("foo"))
	g.Expect(shadowing.Upstream.Namespace).To(gomega.Equal("bar"))
	g.Expect(shadowing.Percentage).To(gomega.Equal(float32(100.0)))
}

// NOTE: Gloo Edge Proxy IR doesn't support multiple mirror/shadow policies on the same route
// func TestMultipleMirrors(t *testing.T) {
// 	g := gomega.NewWithT(t)
// 	ctrl := gomock.NewController(t)
// 	queries := mocks.NewMockGatewayQueries(ctrl)
// 	rt := &gwv1.HTTPRoute{}
// 	ctx := &filterplugins.RouteContext{
// 		Ctx:     context.Background(),
// 		Queries: queries,
// 		Route:   rt,
// 	}
// 	filter1 := gwv1.HTTPRouteFilter{
// 		Type: gwv1.HTTPRouteFilterRequestMirror,
// 		RequestMirror: &gwv1.HTTPRequestMirrorFilter{
// 			BackendRef: gwv1.BackendObjectReference{
// 				Name: "foo",
// 				Port: ptr(gwv1.PortNumber(8080)),
// 			},
// 		},
// 	}
// 	filter2 := gwv1.HTTPRouteFilter{
// 		Type: gwv1.HTTPRouteFilterRequestMirror,
// 		RequestMirror: &gwv1.HTTPRequestMirrorFilter{
// 			BackendRef: gwv1.BackendObjectReference{
// 				Name: "bar",
// 				Port: ptr(gwv1.PortNumber(8080)),
// 			},
// 		},
// 	}
// 	svc1 := &corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "foo",
// 			Namespace: "bar",
// 		},
// 	}
// 	svc2 := &corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "bar",
// 			Namespace: "foo",
// 		},
// 	}
// 	queries.EXPECT().ObjToFrom(rt).Return(nil)
// 	queries.EXPECT().GetBackendForRef(ctx.Ctx, gomock.Any(), &filter1.RequestMirror.BackendRef).Return(svc1, nil)
// 	queries.EXPECT().ObjToFrom(rt).Return(nil)
// 	queries.EXPECT().GetBackendForRef(ctx.Ctx, gomock.Any(), &filter2.RequestMirror.BackendRef).Return(svc2, nil)

// 	plugin := mirror.NewPlugin()
// 	outputRoute := &routev3.Route{
// 		Action: &routev3.Route_Route{
// 			Route: &routev3.RouteAction{},
// 		},
// 	}
// 	plugin.ApplyFilter(ctx, filter1, outputRoute)
// 	plugin.ApplyFilter(ctx, filter2, outputRoute)

// 	g.Expect(outputRoute.GetRoute().RequestMirrorPolicies).To(gomega.HaveLen(2))
// 	g.Expect(outputRoute.GetRoute().RequestMirrorPolicies[0].Cluster).To(gomega.Equal("bar-foo-8080"))
// 	g.Expect(outputRoute.GetRoute().RequestMirrorPolicies[1].Cluster).To(gomega.Equal("foo-bar-8080"))
// }
