package local_e2e

import (
	"net/http"

	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-testing/helpers/local"

	"github.com/k0kubun/pp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HappyPath", func() {
	var (
		envoyInstance *localhelpers.EnvoyInstance
		glooInstance  *localhelpers.GlooInstance
	)
	BeforeEach(func() {
		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())
		glooInstance, err = glooFactory.NewGlooInstance()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if envoyInstance != nil {
			envoyInstance.Clean()
		}
		if glooInstance != nil {
			glooInstance.Clean()
		}
	})

	It("Receive proxied request", func() {
		err := envoyInstance.Run()
		Expect(err).NotTo(HaveOccurred())

		err = glooInstance.Run()
		Expect(err).NotTo(HaveOccurred())

		envoyPort := glooInstance.EnvoyPort()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tu := NewTestUpstream(ctx)
		err = glooInstance.AddUpstream(tu.Upstream)
		Expect(err).NotTo(HaveOccurred())

		v := &v1.VirtualHost{
			Name: "default",
			Routes: []*v1.Route{{
				Matcher: &v1.Route_RequestMatcher{
					RequestMatcher: &v1.RequestMatcher{
						Path: &v1.RequestMatcher_PathPrefix{PathPrefix: "/"},
					},
				},
				SingleDestination: &v1.Destination{
					DestinationType: &v1.Destination_Upstream{
						Upstream: &v1.UpstreamDestination{
							Name: tu.Upstream.Name,
						},
					},
				},
			}},
		}

		err = glooInstance.AddVhost(v)
		Expect(err).NotTo(HaveOccurred())

		body := []byte("solo.io test")
		timeout := time.After(1 * time.Minute)
		for {
			select {
			case request := <-tu.C:
				pp.Fprintf(GinkgoWriter, "%v", request)
				Expect(request.Body).NotTo(BeNil())
				Expect(request.Body).To(Equal(body))
				return
			case <-time.After(time.Second):
				// call the server again is it might not have initialized
				var buf bytes.Buffer
				buf.Write(body)
				_, err := http.Post(fmt.Sprintf("http://%s:%d", "localhost", envoyPort), "application/octet-stream", &buf)
				if err != nil {
					//	fmt.Println("post err " + err.Error())
				}
			case <-timeout:
				panic("timeout")
			}
		}

	})

})
