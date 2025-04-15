package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestChildRouteCanAttachToParentRef(t *testing.T) {
	testCases := []struct {
		name      string
		route     *gwv1.HTTPRoute
		parentRef types.NamespacedName
		expected  bool
	}{
		{
			name: "no ParentRefs, should allow attachment",
			route: &gwv1.HTTPRoute{
				Spec: gwv1.HTTPRouteSpec{},
			},
			parentRef: types.NamespacedName{Name: "parent", Namespace: "default"},
			expected:  true,
		},
		{
			name: "ParentRefs match, should allow attachment",
			route: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
								Kind:      ptr.To(gwv1.Kind("HTTPRoute")),
								Name:      "parent",
								Namespace: ptr.To(gwv1.Namespace("default")),
							},
						},
					},
				},
			},
			parentRef: types.NamespacedName{Name: "parent", Namespace: "default"},
			expected:  true,
		},
		{
			name: "ParentRef doesn't match Name, should not allow attachment",
			route: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
								Kind:      ptr.To(gwv1.Kind("HTTPRoute")),
								Name:      "invalid",
								Namespace: ptr.To(gwv1.Namespace("default")),
							},
						},
					},
				},
			},
			parentRef: types.NamespacedName{Name: "parent", Namespace: "default"},
			expected:  false,
		},
		{
			name: "ParentRef doesn't match Namespace, should not allow attachment",
			route: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
								Kind:      ptr.To(gwv1.Kind("HTTPRoute")),
								Name:      "parent",
								Namespace: ptr.To(gwv1.Namespace("invalid")),
							},
						},
					},
				},
			},
			parentRef: types.NamespacedName{Name: "parent", Namespace: "default"},
			expected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			result := ChildRouteCanAttachToParentRef(tc.route, tc.parentRef)
			a.Equal(tc.expected, result)
		})
	}
}

func TestIsDelegatedRouteMatch(t *testing.T) {
	testCases := []struct {
		name       string
		parent     gwv1.HTTPRouteMatch
		parentRef  types.NamespacedName
		child      gwv1.HTTPRouteMatch
		childNs    string
		parentRefs []gwv1.ParentReference
		expected   bool
	}{
		{
			name: "child route without parentRef matches parent",
			parent: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
				},
				Method: ptr.To[gwv1.HTTPMethod]("GET"),
			},
			child: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo/baz"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header3"),
						Value: "val3",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query3"),
						Value: "val3.*foo",
					},
				},
				Method: ptr.To[gwv1.HTTPMethod]("GET"),
			},
			expected: true,
		},
		{
			name: "child route without parentRef doesn't match parent path",
			parent: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
				},
			},
			child: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/bar/baz"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header3"),
						Value: "val3",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query3"),
						Value: "val3.*foo",
					},
				},
			},
			expected: false,
		},
		{
			name: "child route without parentRef doesn't match parent headers",
			parent: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
				},
			},
			child: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo/baz"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header3"),
						Value: "val3",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query3"),
						Value: "val3.*foo",
					},
				},
			},
			expected: false,
		},
		{
			name: "child route without parentRef doesn't parent query params",
			parent: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
				},
			},
			child: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo/baz"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header3"),
						Value: "val3",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query3"),
						Value: "val3.*foo",
					},
				},
			},
			expected: false,
		},
		{
			name: "child route without parentRef doesn't match parent method",
			parent: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
				},
				Method: ptr.To[gwv1.HTTPMethod]("GET"),
			},
			child: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo/baz"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header3"),
						Value: "val3",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query3"),
						Value: "val3.*foo",
					},
				},
				Method: ptr.To[gwv1.HTTPMethod]("PUT"),
			},
			expected: false,
		},
		{
			name: "child route with parentRef matches parent",
			parent: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
				},
			},
			child: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo/baz"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header3"),
						Value: "val3",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query3"),
						Value: "val3.*foo",
					},
				},
			},
			expected: true,
		},
		{
			name: "child route with parentRef matches parent without parentRef.Namespace set",
			parent: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
				},
			},
			child: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/foo/baz"),
				},
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("header2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Name:  gwv1.HTTPHeaderName("header3"),
						Value: "val3",
					},
				},
				QueryParams: []gwv1.HTTPQueryParamMatch{
					{
						Type:  ptr.To(gwv1.QueryParamMatchExact),
						Name:  gwv1.HTTPHeaderName("query1"),
						Value: "val1",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query2"),
						Value: "val2.*foo",
					},
					{
						Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
						Name:  gwv1.HTTPHeaderName("query3"),
						Value: "val3.*foo",
					},
				},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			actual := IsDelegatedRouteMatch(tc.parent, tc.child)

			a.Equal(tc.expected, actual)
		})
	}
}
