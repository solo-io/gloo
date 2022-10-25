package openApi_test

import (
	"fmt"
	"log"
	"os"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	openApi "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql"
	translate_graphql "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql/graphqlschematranslation"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("GraphQl Discovery Test", func() {
	var (
		opts bootstrap.Opts
	)
	BeforeEach(func() {
		opts = bootstrap.Opts{
			Settings: &v1.Settings{
				Discovery: &v1.Settings_DiscoveryOptions{
					FdsOptions: &v1.Settings_DiscoveryOptions_FdsOptions{
						GraphqlEnabled: &wrapperspb.BoolValue{Value: true},
					},
				},
			},
		}
	})

	Context("Graphql FDS Disabled", func() {
		BeforeEach(func() {
			opts.Settings.GetDiscovery().GetFdsOptions().GraphqlEnabled.Value = false
		})

		It("should return a nil Discovery Factory when disabled in settings", func() {
			Expect(openApi.NewFunctionDiscoveryFactory(opts)).To(BeNil())
		})

	})

	Context("GraphQL Discovery Enabled", func() {
		It("should enable GraphQL FDS when no setting is present", func() {
			opts.Settings.GetDiscovery().FdsOptions = nil
			Expect(openApi.NewFunctionDiscoveryFactory(opts)).NotTo(BeNil())
		})

		var (
			astDoc     *ast.Document
			schema     *graphql.Schema
			resolution map[string]*Resolution
			spec       *openapi.T
			DEBUG      = os.Getenv("DEBUG") == "1"
		)

		AfterEach(func() {
			astDoc = nil
			spec = nil
			schema = nil
			resolution = nil
		})

		translateToGraphqlSchemaAndResolutions := func(openapiSchema string) {
			var err error
			spec, err = openApi.GetOpenApi3Doc([]byte(openapiSchema))
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			oass := []*openapi.T{spec}
			t := translate_graphql.NewOasToGqlTranslator(&v1.Upstream{
				Metadata: &core.Metadata{
					Name:      "upstream-name",
					Namespace: "upstream-namespace",
				},
			})
			astDoc, schema, resolution, err = t.CreateGraphqlSchema(oass)

			// Set env var "DEBUG" to "1" to print out debug logs for this test.
			if DEBUG {
				schemaString := printer.Print(ast.Node(astDoc))
				fmt.Println(schemaString.(string))
				b, err := yaml.Marshal(resolution)
				if err != nil {
					log.Fatalf(err.Error())
				}
				fmt.Printf("\n%s\n", b)
			}
		}

		Context("Resolver path header", func() {

			testPath := func(openapischema, expectedPath string) {
				translateToGraphqlSchemaAndResolutions(openapischema)
				ExpectWithOffset(1, schema.QueryType().Fields()).To(HaveKey("getCat"))
				getCatField := schema.QueryType().Fields()["getCat"]
				ExpectWithOffset(1, getCatField.Description).To(Equal(`Return a cat.

Equivalent to OpenApiSpec 'Some title' GET /cat`))
				ExpectWithOffset(1, resolution).To(HaveLen(1))
				headers := resolution["upstream-name_upstream-namespace|Query|getCat"].GetRestResolver().GetRequest().GetHeaders()
				ExpectWithOffset(1, headers[":method"]).To(Equal("GET"))
				ExpectWithOffset(1, headers[":path"]).To(Equal(expectedPath))
			}

			It("Path is / when no servers are provided", func() {
				openapischema := `openapi: 3.0.0
info:
  title: Some title
paths:
  "/cat":
    get:
      description: Return a cat.`
				testPath(openapischema, "/cat")
			})

			It("Uses base url from first provided server", func() {
				openapischema := `openapi: 3.0.0
info:
  title: Some title
servers:
  - url: /v3/oas
  - url: /v3/shouldntusethis
paths:
  "/cat":
    get:
      description: Return a cat.`
				testPath(openapischema, "/v3/oas/cat")
			})

			It("Uses base url from first provided path (path overrides spec-level server)", func() {
				openapischema := `openapi: 3.0.0
info:
  title: Some title
paths:
  "/cat":
    get:
      servers:
        - url: /v3/pathitem
        - url: /v3/shouldntusethis
      description: Return a cat.`
				testPath(openapischema, "/v3/pathitem/cat")
			})
		})
	})
})
