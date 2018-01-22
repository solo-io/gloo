package module

import (
	"bytes"
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	var (
		yml = `- example_rule:
    timeout: 15s
    match:
      path: /foo
    upstream:
      name: foo_service
      address: 127.0.0.1
      port: 9090
- example_rule:
    timeout: 5s
    match:
      path: /bar
    upstream:
      name: bar_service
      address: 127.0.0.1
      port: 9091`
		jsn = `[
  {
    "example_rule": {
      "timeout": "15s",
      "match": {
        "path": "/foo"
      },
      "upstream": {
        "name": "foo_service",
        "address": "127.0.0.1",
        "port": 9090
      }
    }
  },
  {
    "example_rule": {
      "timeout": "5s",
      "match": {
        "path": "/bar"
      },
      "upstream": {
        "name": "bar_service",
        "address": "127.0.0.1",
        "port": 9091
      }
    }
  }
]`
		expected = `[{"match":{"path":"/foo"},"timeout":"15s","upstream":{"address":"127.0.0.1","name":"foo_service","port":9090}},{"match":{"path":"/bar"},"timeout":"5s","upstream":{"address":"127.0.0.1","name":"bar_service","port":9091}}]`
		key      = `example_rule`
	)
	Describe("getBlobsFromYml", func() {
		blob, err := getBlobsFromYml([]byte(yml), key)
		It("extracts and aggregates blobs by their key", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(string(blob)).To(Equal(expected))
		})
	})
	Describe("getBlobsFromJson", func() {
		blob := getBlobsFromJson([]byte(jsn), key)
		buffer := new(bytes.Buffer)
		err := json.Compact(buffer, blob)
		result := buffer.String()
		It("extracts and aggregates blobs by their key", func() {
			buffer = new(bytes.Buffer)
			err = json.Compact(buffer, blob)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(buffer.String()))
		})
	})
})
