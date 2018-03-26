package local_e2e

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-testing/local_e2e/test_grpc_service/glootest/protos"
	"github.com/solo-io/gloo/pkg/log"
)

var _ = FDescribe("GRPC Plugin", func() {
	It("Routes to GRPC Functions", func() {
		err := envoyInstance.Run()
		Expect(err).NotTo(HaveOccurred())

		err = glooInstance.Run()
		Expect(err).NotTo(HaveOccurred())

		log.Printf("gloo: %v", glooInstance)

		envoyPort := glooInstance.EnvoyPort()

		tu := NewTestGRPCUpstream(glooInstance.FilesDir())
		err = glooInstance.AddUpstream(tu.Upstream)
		Expect(err).NotTo(HaveOccurred())

		v := &v1.VirtualHost{
			Name: "default",
			Routes: []*v1.Route{
				{
					Matcher: &v1.Route_RequestMatcher{
						RequestMatcher: &v1.RequestMatcher{
							Path: &v1.RequestMatcher_PathExact{PathExact: "/test"},
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Function{
							Function: &v1.FunctionDestination{
								UpstreamName: tu.Upstream.Name,
								FunctionName: "TestMethod",
							},
						},
					},
				},
			},
		}

		err = glooInstance.AddVhost(v)
		Expect(err).NotTo(HaveOccurred())

		body := []byte("solo.io test")

		time.Sleep(time.Second)

		// wait for envoy to start receiving request
		Eventually(func() error {
			// send a request with a body
			var buf bytes.Buffer
			buf.Write(body)
			res, err := http.Post(fmt.Sprintf("http://%s:%d/test", "localhost", envoyPort), "application/json", &buf)
			if err == nil {
				log.Printf("%v", res.Header)
				log.Printf("%v", res.Status)
			}
			return err
		}, 60, 1).Should(BeNil())

		expectedResponse := &ReceivedRequest{
			GRPCRequest: &glootest.TestResponse{Str: string(body)},
		}
		Eventually(tu.C).Should(Receive(Equal(expectedResponse)))

	})

})
