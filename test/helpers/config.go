package helpers

import (
	google_protobuf "github.com/golang/protobuf/ptypes/struct"
	"github.com/solo-io/glue/pkg/api/types"
)

func NewTestConfig() *types.Config {
	routes := []*types.Route{
		NewTestRoute1(),
		NewTestRoute2(),
	}
	upstreams := []*types.Upstream{
		{
			Name: "aws",
			Type: "lambda",
			Spec: &google_protobuf.Struct{
				Fields: map[string]*google_protobuf.Value{
					"region":         {Kind: &google_protobuf.Value_StringValue{StringValue: "us_east_1"}},
					"secret_key_ref": {Kind: &google_protobuf.Value_StringValue{StringValue: "my-aws-secret-key"}},
					"access_key_ref": {Kind: &google_protobuf.Value_StringValue{StringValue: "my-aws-access-key"}},
				},
			},
			Functions: []*types.Function{
				{
					Name: "my_lambda_function",
					Spec: &google_protobuf.Struct{
						Fields: map[string]*google_protobuf.Value{
							"key": {Kind: &google_protobuf.Value_StringValue{StringValue: "value"}},
						},
					},
				},
			},
		},
		{
			Name: "my_upstream",
			Type: "service",
			Spec: &google_protobuf.Struct{
				Fields: map[string]*google_protobuf.Value{
					"url": {Kind: &google_protobuf.Value_StringValue{StringValue: "http://www.example.com"}},
				},
			},
		},
	}
	virtualhosts := []*types.VirtualHost{
		{
			Domains: []string{"*.example.io"},
			SslConfig: &types.SSLConfig{
				CaCertPath: "/etc/my_crts/ca.crt",
			},
		},
	}
	return &types.Config{
		Routes:       routes,
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func NewTestRoute1() *types.Route {
	return &types.Route{
		Matcher: &types.Matcher{
			Path: &types.Matcher_Prefix{
				Prefix: "/foo",
			},
			Headers:     map[string]string{"x-foo-bar": ""},
			Verbs:       []string{"GET", "POST"},
			VirtualHost: "my_vhost",
		},
		Destination: &types.Route_FunctionName{
			FunctionName: &types.FunctionDestination{
				FunctionName: "foo",
				UpstreamName: "aws",
			},
		},
		Plugins: map[string]*google_protobuf.Struct{
			"auth": {
				Fields: map[string]*google_protobuf.Value{
					"username": {Kind: &google_protobuf.Value_StringValue{StringValue: "alice"}},
					"password": {Kind: &google_protobuf.Value_StringValue{StringValue: "bob"}},
				},
			},
		},
	}
}

func NewTestRoute2() *types.Route {
	return &types.Route{
		Matcher: &types.Matcher{
			Path: &types.Matcher_Exact{
				Exact: "/bar",
			},
			Verbs: []string{"GET", "POST"},
		},
		Destination: &types.Route_FunctionName{
			FunctionName: &types.FunctionDestination{
				FunctionName: "foo",
				UpstreamName: "aws",
			},
		},
		Plugins: map[string]*google_protobuf.Struct{
			"auth": {
				Fields: map[string]*google_protobuf.Value{
					"username": {Kind: &google_protobuf.Value_StringValue{StringValue: "alice"}},
					"password": {Kind: &google_protobuf.Value_StringValue{StringValue: "bob"}},
				},
			},
		},
	}
}
