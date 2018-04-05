package extensions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/types"
	. "github.com/solo-io/gloo/pkg/coreplugins/route-extensions"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/protoutil"
)

var _ = Describe("Spec", func() {
	It("decodes and uncodes from yaml", func() {
		spec := RouteExtensionSpec{
			MaxRetries: 1,
			Timeout:    time.Minute,
			AddRequestHeaders: []HeaderValue{
				{
					Key:   "FOO",
					Value: "BAR",
				},
			},
		}
		encoded := EncodeRouteExtensionSpec(spec)
		jsn, err := protoutil.Marshal(encoded)
		Expect(err).NotTo(HaveOccurred())
		yam, err := yaml.JSONToYAML(jsn)
		Expect(err).NotTo(HaveOccurred())
		fmt.Printf("%s\n", yam)
	})
	It("decodes from yaml", func() {
		yam := `
add_request_headers:
- key: FOO
  value: BAR
- key: BOO
  value: BAR
  append: true
max_retries: 1
prefix_rewrite: /foo
timeout: 60000000000
something_invalid: another_spec_maybe?`
		jsn, err := yaml.YAMLToJSON([]byte(yam))
		Expect(err).NotTo(HaveOccurred())
		var struc types.Struct
		err = protoutil.Unmarshal(jsn, &struc)
		Expect(err).NotTo(HaveOccurred())
		specc, err := DecodeRouteExtensions(&struc)
		Expect(err).NotTo(HaveOccurred())
		log.Printf("%v", specc)
	})
})
