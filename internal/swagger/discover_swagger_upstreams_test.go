package swagger_test

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"net/http"

	"strconv"
	"strings"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-function-discovery/internal/swagger"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-plugins/common/annotations"
	"github.com/solo-io/gloo-plugins/rest"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

var _ = Describe("DiscoverSwaggerUpstreams", func() {
	Describe("happy path", func() {
		Context("upstream for a service that serves a swagger spec", func() {
			// create a test swagger server
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, swaggerDoc)
			}))
			BeforeEach(func() {
				//srv.Start()
			})
			AfterEach(func() {
				srv.Close()
			})
			It("marks the upstream annotations with swagger", func() {
				u := strings.Split(srv.URL, ":")
				addr := strings.TrimPrefix(u[1], "//")
				port, _ := strconv.Atoi(u[2])

				us := &v1.Upstream{
					Name:     "something",
					Type:     service.UpstreamTypeService,
					Metadata: &v1.Metadata{Annotations: make(map[string]string)},
					Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
						Hosts: []service.Host{
							{
								Addr: addr,
								Port: uint32(port),
							},
						},
					}),
				}
				DiscoverSwaggerUpstream(&resolver.Resolver{}, []string{"/test/path"}, us)
				Expect(us.Metadata.Annotations).To(HaveKey(annotations.ServiceType))
				Expect(us.Metadata.Annotations[annotations.ServiceType]).To(Equal(rest.ServiceTypeTransformation))
				Expect(us.Metadata.Annotations).To(HaveKey(AnnotationKeySwaggerURL))
				Expect(us.Metadata.Annotations[AnnotationKeySwaggerURL]).To(Equal(srv.URL + "/test/path"))
			})
		})
	})
})

const swaggerDoc = `{
  "swagger": "2.0",
  "info": {
    "version": "1.0.0",
    "title": "Swagger Petstore",
    "description": "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "name": "Swagger API Team"
    },
    "license": {
      "name": "MIT"
    }
  },
  "host": "petstore.swagger.io",
  "basePath": "/api",
  "schemes": [
    "http"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/pets": {
      "get": {
        "description": "Returns all pets from the system that the user has access to",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "A list of pets.",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Pet"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Pet": {
      "type": "object",
      "required": [
        "id",
        "name"
      ],
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      }
    }
  }
}`
