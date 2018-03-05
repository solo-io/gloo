package transformation_test

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-plugins/transformation"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/plugin"
)

var _ = Describe("Transformation", func() {
	FIt("processes response transformations", func() {
		p := &Plugin{CachedTransformations: make(map[string]*Transformation)}
		out := &envoyroute.Route{}
		params := &plugin.RoutePluginParams{}
		in := NewSingleDestRoute("nothing", "nothinf")
		in.Extensions = EncodeRouteExtension(RouteExtension{
			ResponseTemplate: &Template{
				Body: "{{body}}",
			},
		})
		err := p.ProcessRoute(params, in, out)
		Expect(err).NotTo(HaveOccurred())
		Expect(p.CachedTransformations).To(HaveLen(1))
		log.Printf("%v", out)
		log.Printf("%v", p.HttpFilters(nil))
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
					Regex:    "([_[:alnum:]]+)",
					Subgroup: 1,
				}))
			}
			// specific to the route
			Expect(trans.Extractors["id"]).To(Equal(&Extraction{
				Header:   ":path",
				Regex:    "/u\\(se\\)rs/([_[:alnum:]]+)/accounts/([_[:alnum:]]+)",
				Subgroup: 1,
			}))
			Expect(trans.Extractors["account"]).To(Equal(&Extraction{
				Header:   ":path",
				Regex:    "/u\\(se\\)rs/([_[:alnum:]]+)/accounts/([_[:alnum:]]+)",
				Subgroup: 2,
			}))
			Expect(trans.Extractors["type"]).To(Equal(&Extraction{
				Header:   "content-type",
				Regex:    "application/([_[:alnum:]]+)",
				Subgroup: 1,
			}))
			Expect(out.Metadata.FilterMetadata["io.solo.transformation"].Fields["transformation"].Kind).To(Equal(&types.Value_StringValue{StringValue: hash}))
		}
	})
})

func NonFunctionalUpstream(name string) *v1.Upstream {
	return &v1.Upstream{
		Name: name,
		Type: service.UpstreamTypeService,
		Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
			Hosts: []service.Host{
				{
					Addr: "127.0.0.1",
					Port: 9090,
				},
			},
		}),
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
				Path: "/{id}/why/{id}",
				Header: map[string]string{
					"x-user-id":      "{id}",
					"x-content-type": "{type}",
				},
				Body: `{"id":"{{id}}", "type":"{{type}}"}`,
			}),
		},
	}
	return us
}
