package graphql_test

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	graphql "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/graphql-graphql"
)

var _ = Describe("GraphQL introspection Discovery tests", func() {

	Context("run with a graphql server docker container", func() {

		It("should make a query, and match the expected schema results", func() {
			// this test will call the server created in the suite.
			var introspectionResults []byte
			var err error

			u := url.URL{Host: fmt.Sprintf("%s:%s", "localhost", graphqlPort), Scheme: "http", Path: "graphql"}
			introspectionResults, err = graphql.GetIntrospectionResultsFromHost(&u)
			Expect(err).ToNot(HaveOccurred())
			expectedIntrospectionResult, err := os.ReadFile("test/introspection_from_query.json")
			Expect(err).ToNot(HaveOccurred())
			actualSchema, err := graphql.GetSchemaFromIntrospectionJSON(introspectionResults)
			Expect(err).ToNot(HaveOccurred())
			expectedSchema, err := graphql.GetSchemaFromIntrospectionJSON(expectedIntrospectionResult)
			Expect(err).ToNot(HaveOccurred())
			// the order of the lines always appears to be random, so it is better to just look at each line seperately
			lines := strings.Split(actualSchema, "\n")
			for _, l := range lines {
				// because updateTodo has multiple parameters, the order can be change
				if strings.Contains(l, "updateTodo(") {
					Expect(expectedSchema).To(ContainSubstring("updateTodo("))
					Expect(expectedSchema).To(ContainSubstring("done: Boolean"))
					Expect(expectedSchema).To(ContainSubstring("id: String!"))
					Expect(expectedSchema).To(ContainSubstring("): Todo"))
				} else {
					Expect(expectedSchema).To(ContainSubstring(l))
				}
			}
		})

	})

	Context("GetSchemaFromIntrospectionJSON", func() {
		It("Should produce a matching schema", func() {
			instrospectionJSON, err := os.ReadFile("./test/introspection.json")
			Expect(err).ToNot(HaveOccurred())
			schema, err := graphql.GetSchemaFromIntrospectionJSON(instrospectionJSON)
			Expect(err).ToNot(HaveOccurred())
			expectedSchemabytes, err := os.ReadFile("./test/expected_schema.graphql")
			Expect(err).ToNot(HaveOccurred())
			expectedSchema := string(expectedSchemabytes)
			Expect(schema).To(Equal(string(expectedSchema)))
		})

		It("should be able to differentiate between a schema and not", func() {
			invalidJSON := []byte("this is invalid JSON for introspection")
			_, err := graphql.GetSchemaFromIntrospectionJSON(invalidJSON)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error converting introspection JSON to graphql schema"))
		})
	})
})
