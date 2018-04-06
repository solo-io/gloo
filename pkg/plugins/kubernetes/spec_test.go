package kubernetes_test

import (
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/log"
	. "github.com/solo-io/gloo/pkg/plugins/kubernetes"
)

var _ = Describe("FromMap", func() {
	It("correctly deserializes Spec from a map[string]interface{}", func() {
		m := &types.Struct{
			Fields: map[string]*types.Value{
				"service_name":      {Kind: &types.Value_StringValue{StringValue: "svc"}},
				"service_namespace": {Kind: &types.Value_StringValue{StringValue: "my-ns"}},
				"service_port":      {Kind: &types.Value_NumberValue{NumberValue: 8080}},
			},
		}
		spec, err := DecodeUpstreamSpec(m)
		log.Debugf("%v", spec)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.ServiceName).To(Equal("svc"))
		Expect(spec.ServiceNamespace).To(Equal("my-ns"))
		Expect(spec.ServicePort).To(Equal(int32(8080)))
	})
})
