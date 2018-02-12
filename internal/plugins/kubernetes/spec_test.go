package kubernetes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/glue/internal/pkg/kube/upstream"
	"github.com/solo-io/glue/pkg/log"
)

var _ = Describe("FromMap", func() {
	It("correctly deserializes Spec from a map[string]interface{}", func() {
		m := map[string]interface{}{
			"service_name":      "svc",
			"service_namespace": "my-ns",
			"service_port_name": "svc-port",
		}
		spec, err := FromMap(m)
		log.Debugf("%v", spec)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.ServiceName).To(Equal(m["service_name"]))
		Expect(spec.ServiceNamespace).To(Equal(m["service_namespace"]))
		Expect(spec.ServicePortName).To(Equal(m["service_port_name"]))
	})
})
