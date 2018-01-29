package helpers

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
)

func NewTestConfig() v1.Config {
	routes := []v1.Route{
		NewTestRoute1(),
		NewTestRoute2(),
	}
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
			Name: "my_upstream",
			Type: "service",
			Spec: map[string]interface{}{
				"auth": map[string]interface{}{
					"url": "http://www.example.com",
				},
			},
		},
	}
	virtualhosts := []v1.VirtualHost{
		{
			Domains: []string{"*.example.io"},
			SSLConfig: v1.SSLConfig{
				CACertPath: "/etc/my_crts/ca.crt",
			},
		},
	}
	return v1.Config{
		Routes:       routes,
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func NewTestRoute1() v1.Route {
	return v1.Route{
		Matcher: v1.Matcher{
			Path: v1.Path{
				Prefix: "/foo",
			},
			Headers:     map[string]string{"x-foo-bar": ""},
			Verbs:       []string{"GET", "POST"},
			VirtualHost: "my_vhost",
		},
		Destination: v1.Destination{
			FunctionDestionation: &v1.FunctionDestination{
				FunctionName: "foo",
				UpstreamName: "aws",
			},
		},
		Plugins: map[string]interface{}{
			"auth": map[string]string{
				"username": "alice",
				"password": "bob",
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
			FunctionDestionation: &v1.FunctionDestination{
				FunctionName: "foo",
				UpstreamName: "aws",
			},
		},
		Plugins: map[string]interface{}{
			"auth": map[string]interface{}{
				"username": "alice",
				"password": "bob",
			},
		},
	}
}
