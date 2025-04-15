package query

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/schemes"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
)

var _ = DescribeTable("getDelegatedChildren",
	func(parentRef types.NamespacedName, backendRef gwv1.HTTPBackendRef, objects []client.Object, wantChildren int, wantErr error) {
		scheme := schemes.GatewayScheme()
		builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...)
		err := IterateIndices(func(o client.Object, f string, fun client.IndexerFunc) error {
			builder.WithIndex(o, f, fun)
			return nil
		})
		fakeClient := builder.Build()

		q := &gatewayQueries{
			client: fakeClient,
			scheme: scheme,
		}

		parentRoute := &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: parentRef.Namespace,
				Name:      parentRef.Name,
			},
			Spec: gwv1.HTTPRouteSpec{
				Rules: []gwv1.HTTPRouteRule{
					{
						BackendRefs: []gwv1.HTTPBackendRef{
							backendRef,
						},
					},
				},
			},
		}

		backendMap := q.getDelegatedChildren(context.TODO(), parentRoute, sets.New[types.NamespacedName]())
		backendKey := backendToRefKey(backendRef.BackendObjectReference)
		err = backendMap.errors[backendKey]
		children := backendMap.items[backendKey]

		if wantErr != nil {
			Expect(err).To(MatchError(wantErr))
		}
		Expect(children).To(HaveLen(wantChildren))
	},
	Entry(
		"with wildcard label selector",
		types.NamespacedName{Namespace: "parent-ns", Name: "parent"},
		gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Group:     ptr.To(gwv1.Group("delegation.gateway.solo.io")),
					Kind:      ptr.To(gwv1.Kind("label")),
					Name:      "delegate",
					Namespace: ptr.To(gwv1.Namespace(wellknown.RouteDelegationLabelSelectorWildcardNamespace)),
				},
			},
		},
		[]client.Object{
			&gwv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       wellknown.HTTPRouteKind,
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "r1",
					Labels: map[string]string{
						wellknown.RouteDelegationLabelSelector: "delegate",
					},
				},
			},
			&gwv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       wellknown.HTTPRouteKind,
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "r2",
					Labels: map[string]string{
						wellknown.RouteDelegationLabelSelector: "delegate",
					},
				},
			},
		},
		2,
		nil,
	),
	Entry(
		"with wildcard label selector when wildcard namespace exists",
		types.NamespacedName{Namespace: "parent-ns", Name: "parent"},
		gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Group:     ptr.To(gwv1.Group("delegation.gateway.solo.io")),
					Kind:      ptr.To(gwv1.Kind("label")),
					Name:      "delegate",
					Namespace: ptr.To(gwv1.Namespace(wellknown.RouteDelegationLabelSelectorWildcardNamespace)),
				},
			},
		},
		[]client.Object{
			&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: wellknown.RouteDelegationLabelSelectorWildcardNamespace,
				},
			},
			&gwv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       wellknown.HTTPRouteKind,
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test",
					Labels: map[string]string{
						wellknown.RouteDelegationLabelSelector: "delegate",
					},
				},
			},
		},
		0,
		ErrWildcardNamespaceDisallowed,
	),
	Entry(
		"filter self reference and mismatched parentRef",
		types.NamespacedName{Namespace: "parent-ns", Name: "parent"},
		gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
					Kind:      ptr.To(gwv1.Kind("HTTPRoute")),
					Name:      "*",
					Namespace: ptr.To(gwv1.Namespace("parent-ns")),
				},
			},
		},
		[]client.Object{
			&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: wellknown.RouteDelegationLabelSelectorWildcardNamespace,
				},
			},
			// self reference
			&gwv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       wellknown.HTTPRouteKind,
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "parent-ns",
					Name:      "parent",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Group: ptr.To(gwv1.Group("gateway.networking.k8s.io")),
								Kind:  ptr.To(gwv1.Kind("HTTPRoute")),
								Name:  "parent",
							},
						},
					},
				},
			},
			// ParentRef mismatch
			&gwv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       wellknown.HTTPRouteKind,
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "parent-ns",
					Name:      "invalid-ref",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Group: ptr.To(gwv1.Group("gateway.networking.k8s.io")),
								Kind:  ptr.To(gwv1.Kind("HTTPRoute")),
								Name:  "invalid", // mismatched parentRef
							},
						},
					},
				},
			},
			// valid child
			&gwv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       wellknown.HTTPRouteKind,
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "parent-ns",
					Name:      "child-1",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Group: ptr.To(gwv1.Group("gateway.networking.k8s.io")),
								Kind:  ptr.To(gwv1.Kind("HTTPRoute")),
								Name:  "parent",
							},
						},
					},
				},
			},
			// valid child
			&gwv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       wellknown.HTTPRouteKind,
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "parent-ns",
					Name:      "child-2",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Group: ptr.To(gwv1.Group("gateway.networking.k8s.io")),
								Kind:  ptr.To(gwv1.Kind("HTTPRoute")),
								Name:  "parent",
							},
						},
					},
				},
			},
		},
		2,
		nil,
	),
)
