package graphql_handler_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/graphql_handler"
)

var _ = Describe("resolver yaml validation", func() {

	It("accepts empty yaml", func() {
		err := graphql_handler.ValidateResolverYaml("", rpc_edge_v1.ValidateResolverYamlRequest_REST_RESOLVER)
		Expect(err).NotTo(HaveOccurred())
	})

	It("throws error for unsupported resolver types", func() {
		err := graphql_handler.ValidateResolverYaml("yaml", rpc_edge_v1.ValidateResolverYamlRequest_RESOLVER_NOT_SET)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid resolver type"))
	})

	It("can parse valid rest resolver yaml", func() {
		yaml := `request:
  headers:
    :method: GET
    :path: /api/v1/products
response:
  resultRoot: "author"
upstreamRef:
  name: default-details-9080
  namespace: gloo-system
`
		err := graphql_handler.ValidateResolverYaml(yaml, rpc_edge_v1.ValidateResolverYamlRequest_REST_RESOLVER)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can parse valid grpc resolver yaml", func() {
		yaml := `upstreamRef:
  name: default-details-9080
  namespace: gloo-system
requestTransform:
  serviceName: my-service
  methodName: my-method
spanName: hello
`
		err := graphql_handler.ValidateResolverYaml(yaml, rpc_edge_v1.ValidateResolverYamlRequest_GRPC_RESOLVER)
		Expect(err).NotTo(HaveOccurred())
	})

	It("throws error for yaml with invalid field", func() {
		yaml := `request:
  headers:
    :method: GET
    :path: /api/v1/products
response:
  helloThere: "author"
upstreamRef:
  name: default-details-9080
  namespace: gloo-system
`
		err := graphql_handler.ValidateResolverYaml(yaml, rpc_edge_v1.ValidateResolverYamlRequest_REST_RESOLVER)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid resolver yaml"))
		Expect(err.Error()).To(ContainSubstring("unknown field \"helloThere\""))
	})

	It("throws error for yaml with invalid value", func() {
		yaml := `request:
  headers:
    :method: GET
    :path: /api/v1/products
upstreamRef: 123
`
		err := graphql_handler.ValidateResolverYaml(yaml, rpc_edge_v1.ValidateResolverYamlRequest_REST_RESOLVER)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid resolver yaml"))
		Expect(err.Error()).To(ContainSubstring("unexpected token 123"))
	})

	It("throws error for yaml with wrong indentation", func() {
		yaml := `request:
  headers:
    :method: GET
    :path: /api/v1/products
  upstreamRef:
    name: default-details-9080
    namespace: gloo-system
`
		err := graphql_handler.ValidateResolverYaml(yaml, rpc_edge_v1.ValidateResolverYamlRequest_REST_RESOLVER)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid resolver yaml"))
		Expect(err.Error()).To(ContainSubstring("unknown field \"upstreamRef\""))
	})
})
