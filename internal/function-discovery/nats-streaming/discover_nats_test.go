package nats_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	. "github.com/solo-io/gloo/internal/function-discovery/nats-streaming"
	"github.com/solo-io/gloo/pkg/plugins/nats-streaming"
)

var _ = Describe("DiscoverNats", func() {
	Describe("happy path", func() {
		Context("upstream for a nats-streaming server", func() {
			It("returns service info for nats-streaming", func() {
				err = natsStreamingInstance.Run()
				Expect(err).NotTo(HaveOccurred())
				detector := NewNatsDetector(natsStreamingInstance.ClusterId())
				svcInfo, annotations, err := detector.DetectFunctionalService(&v1.Upstream{Name: "Test"}, fmt.Sprintf("localhost:%v", natsStreamingInstance.NatsPort()))
				Expect(err).To(BeNil())
				Expect(annotations).To(BeNil())
				Expect(svcInfo).To(Equal(&v1.ServiceInfo{
					Type: natsstreaming.ServiceTypeNatsStreaming,
					Properties: natsstreaming.EncodeServiceProperties(natsstreaming.ServiceProperties{
						ClusterID: natsStreamingInstance.ClusterId(),
					}),
				}))
			})
		})
	})
})
