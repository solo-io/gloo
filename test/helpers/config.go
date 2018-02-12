package helpers

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
)

func NewTestConfig() v1.Config {
	upstreams := []v1.Upstream{
		{
			Name: "aws",
			Type: "lambda",
			Spec: map[string]interface{}{
				"region":         "us-east-1",
				"secret_key_ref": "my-secret-key",
				"access_key_ref": "my-access-key",
			},
			Functions: []v1.Function{
				{
					Name: "my_lambda_function",
					Spec: map[string]interface{}{
						"key": "value",
					},
				},
			},
		},
		{
			Name: "my-upstream",
			Type: "service",
			Spec: map[string]interface{}{
				"auth": map[string]interface{}{
					"url": "http://www.example.com",
				},
			},
		},
	}
	virtualhosts := []v1.VirtualHost{
		NewTestVirtualHost("my-vhost", NewTestRoute1(), NewTestRoute2()),
		NewTestVirtualHost("my-vhost-2", NewTestRoute1(), NewTestRoute2()),
	}
	return v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func NewTestVirtualHost(name string, routes ...v1.Route) v1.VirtualHost {
	return v1.VirtualHost{
		Name:   name,
		Routes: routes,
	}
}

func NewTestRoute1() v1.Route {
	return v1.Route{
		Matcher: v1.Matcher{
			Path: v1.Path{
				Prefix: "/foo",
			},
			Headers: map[string]string{"x-foo-bar": ""},
			Verbs:   []string{"GET", "POST"},
		},
		Destination: v1.Destination{
			SingleDestination: v1.SingleDestination{
				FunctionDestination: &v1.FunctionDestination{
					FunctionName: "foo",
					UpstreamName: "aws",
				},
			},
		},
		Plugins: v1.RoutePluginSpec{
			"auth": map[string]interface{}{
				"credentials": struct {
					Username, Password string
				}{
					Username: "alice",
					Password: "bob",
				},
				"token": "my-12345",
			},
		},
	}
}

func NewTestRoute2() v1.Route {
	return v1.Route{
		Matcher: v1.Matcher{
			Path: v1.Path{
				Exact: "/bar",
			},
			Verbs: []string{"GET", "POST"},
		},
		Destination: v1.Destination{
			SingleDestination: v1.SingleDestination{
				FunctionDestination: &v1.FunctionDestination{
					FunctionName: "foo",
					UpstreamName: "aws",
				},
			},
		},
		Plugins: v1.RoutePluginSpec{
			"auth": map[string]interface{}{
				"username": "alice",
				"password": "bob",
			},
		},
	}
}
