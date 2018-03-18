package transformation_test

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/common/annotations"
	. "github.com/solo-io/gloo-plugins/transformation"
	"github.com/solo-io/gloo/pkg/plugin"
)

var _ = Describe("Transformation", func() {
	It("processes response transformations", func() {
		p := &Plugin{CachedTransformations: make(map[string]*Transformation)}
		out := &envoyroute.Route{}

		upstreamName := "nothing"
		params := &plugin.RoutePluginParams{}

		in := NewNonFunctionSingleDestRoute(upstreamName)
		in.Extensions = EncodeRouteExtension(RouteExtension{
			ResponseTemplate: &Template{
				Body: "{{body}}",
			},
		})
		err := p.ProcessRoute(params, in, out)
		Expect(err).NotTo(HaveOccurred())
		Expect(p.CachedTransformations).To(HaveLen(1))
	})
	It("process route for functional upstream", func() {
		upstreamName := "users-svc"
		funcName := "get_user"
		params := &plugin.RoutePluginParams{
			Upstreams: []*v1.Upstream{
				NewFunctionalUpstream(upstreamName, funcName),
			},
		}
		in := NewSingleDestRoute(upstreamName, funcName)
		p := &Plugin{CachedTransformations: make(map[string]*Transformation)}
		out := &envoyroute.Route{}
		err := p.ProcessRoute(params, in, out)
		Expect(err).NotTo(HaveOccurred())
		Expect(p.CachedTransformations).To(HaveLen(1))
		// gives the free headers automatically
		for hash, trans := range p.CachedTransformations {
			for _, header := range []string{"path",
				"method",
				"scheme",
				"authority",
			} {
				Expect(trans.Extractors[header]).To(Equal(&Extraction{
					Header:   ":" + header,
					Regex:    `([\.\-_[:alnum:]]+)`,
					Subgroup: 1,
				}))
			}
			// specific to the route
			Expect(trans.Extractors["id"]).To(Equal(&Extraction{
				Header:   ":path",
				Regex:    `/u\(se\)rs/([\.\-_[:alnum:]]+)/accounts/([\.\-_[:alnum:]]+)`,
				Subgroup: 1,
			}))
			Expect(trans.Extractors["account"]).To(Equal(&Extraction{
				Header:   ":path",
				Regex:    `/u\(se\)rs/([\.\-_[:alnum:]]+)/accounts/([\.\-_[:alnum:]]+)`,
				Subgroup: 2,
			}))
			Expect(trans.Extractors["type"]).To(Equal(&Extraction{
				Header:   "content-type",
				Regex:    `application/([\.\-_[:alnum:]]+)`,
				Subgroup: 1,
			}))
			Expect(trans.Extractors["foo.bar"]).To(Equal(&Extraction{
				Header:   "x-foo-bar",
				Regex:    `([\.\-_[:alnum:]]+)`,
				Subgroup: 1,
			}))
			Expect(trans.RequestTemplate.Body).To(Equal(&InjaTemplate{
				Text: `{"id":"{{id}}", "type":"{{type}}"}`,
			}))
			Expect(trans.RequestTemplate.Headers).To(Equal(map[string]*InjaTemplate{
				":path":          {Text: "/{{id}}/why/{{id}}"},
				"x-user-id":      {Text: "{{id}}"},
				"x-content-type": {Text: "{{type}}"},
			}))
			Expect(out.Metadata.FilterMetadata["io.solo.transformation"].Fields["request-transformation"].Kind).To(Equal(&types.Value_StringValue{StringValue: hash}))
			Expect(out.Metadata.FilterMetadata["io.solo.transformation"].Fields["request-transformation"].Kind).To(Equal(&types.Value_StringValue{StringValue: hash}))
		}
	})
	It("errors when user provides invalid parameters", func() {
		upstreamName := "users-svc"
		funcName := "get_user"
		params := &plugin.RoutePluginParams{
			Upstreams: []*v1.Upstream{
				NewFunctionalUpstream(upstreamName, funcName),
			},
		}
		in := NewBadExtractorsRoute(upstreamName, funcName)
		p := &Plugin{CachedTransformations: make(map[string]*Transformation)}
		out := &envoyroute.Route{}
		err := p.ProcessRoute(params, in, out)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("f{foo/bar} is not valid syntax. {} braces must be closed and variable names must satisfy regex ([\\.\\-_[:alnum:]]+)"))
	})
})

func NonFunctionalUpstream(name string) *v1.Upstream {
	return &v1.Upstream{
		Name: name,
		Type: "test",
		Metadata: &v1.Metadata{
			Annotations: map[string]string{
				annotations.ServiceType: ServiceTypeHttpFunctions,
			},
		},
	}
}

func NewSingleDestRoute(upstreamName, functionName string) *v1.Route {
	return &v1.Route{
		Matcher: &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path: &v1.RequestMatcher_PathRegex{
					PathRegex: "/users/.*/accounts/.*",
				},
				Verbs: []string{"GET"},
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					FunctionName: functionName,
					UpstreamName: upstreamName,
				},
			},
		},
		Extensions: EncodeRouteExtension(RouteExtension{
			Parameters: &Parameters{
				Path: "/u(se)rs/{id}/accounts/{account}",
				Headers: map[string]string{
					"content-type": "application/{type}",
					"x-foo-bar":    "{foo.bar}",
				},
			},
		}),
	}
}

func NewBadExtractorsRoute(upstreamName, functionName string) *v1.Route {
	return &v1.Route{
		Matcher: &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path: &v1.RequestMatcher_PathRegex{
					PathRegex: "/users/.*/accounts/.*",
				},
				Verbs: []string{"GET"},
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					FunctionName: functionName,
					UpstreamName: upstreamName,
				},
			},
		},
		Extensions: EncodeRouteExtension(RouteExtension{
			Parameters: &Parameters{
				Path: "/u(se)rs/{id}/accounts/{account}",
				Headers: map[string]string{
					"bad-extractor": "f{foo/bar}",
				},
			},
		}),
	}
}

func NewNonFunctionSingleDestRoute(upstreamName string) *v1.Route {
	return &v1.Route{
		Matcher: &v1.Route_RequestMatcher{
			RequestMatcher: &v1.RequestMatcher{
				Path: &v1.RequestMatcher_PathRegex{
					PathRegex: "/users/.*/accounts/.*",
				},
				Verbs: []string{"GET"},
			},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Upstream{
				Upstream: &v1.UpstreamDestination{
					Name: upstreamName,
				},
			},
		},
		Extensions: EncodeRouteExtension(RouteExtension{
			Parameters: &Parameters{
				Path:    "/u(se)rs/{id}/accounts/{account}",
				Headers: map[string]string{"content-type": "application/{type}"},
			},
		}),
	}
}

func NewFunctionalUpstream(name, funcName string) *v1.Upstream {
	us := NonFunctionalUpstream(name)
	us.Functions = []*v1.Function{
		{
			Name: funcName,
			Spec: EncodeFunctionSpec(Template{
				Path: "/{{id}}/why/{{id}}",
				Header: map[string]string{
					"x-user-id":      "{{id}}",
					"x-content-type": "{{type}}",
				},
				Body: `{"id":"{{id}}", "type":"{{type}}"}`,
			}),
		},
	}
	return us
}
