package helpers

import "github.com/solo-io/glue/pkg/api/types"

func NewTestConfig() types.Config {
	routes := []types.Route{
		{
			Matcher: types.Matcher{
				Path: types.PrefixPathMatcher{
					Prefix: "/foo",
				},
				Headers:     map[string]string{"x-foo-bar": ""},
				Verbs:       []string{"GET", "POST"},
				VirtualHost: "my_vhost",
			},
			Destination: types.FunctionDestination{
				FunctionName: "aws.foo",
			},
			Plugins: map[string]types.Spec{
				"auth": {"username": "alice", "password": "bob"},
			},
		},
		{
			Matcher: types.Matcher{
				Path: types.ExactPathMatcher{
					Exact: "/bar",
				},
				Verbs: []string{"GET", "POST"},
			},
			Destination: types.UpstreamDestination{
				UpstreamName:  "my_upstream",
				RewritePrefix: "/baz",
			},
			Plugins: map[string]types.Spec{
				"auth": {"username": "alice", "password": "bob"},
			},
		},
	}
	upstreams := []types.Upstream{
		{
			Name: "aws",
			Type: "lambda",
			Spec: types.Spec{"region": "us_east_1", "secret_key_ref": "my-aws-secret-key", "access_key_ref": "my0aws-access-key"},
			Functions: []types.Function{
				{
					Name: "my_lambda_function",
					Spec: types.Spec{
						"context_parameter": map[string]string{"KEY": "VAL"},
					},
				},
			},
		},
		{
			Name: "my_upstream",
			Type: "service",
			Spec: types.Spec{"url": "https://myapi.example.com"},
		},
	}
	virtualhosts := []types.VirtualHost{
		{
			Domains:   []string{"*.example.io"},
			SSLConfig: types.SSLConfig{},
		},
	}
	return types.Config{
		Routes:       routes,
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}
