package local_e2e

import (
	"net/http"

	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-testing/helpers/local"
	"github.com/solo-io/gloo/pkg/coreplugins/service"

	"github.com/k0kubun/pp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// . "github.com/solo-io/gloo-testing/local_e2e"
)

type ReceivedRequest struct {
	Method string
	Body   []byte
}
type TestUpstream struct {
	Upstream *v1.Upstream
	C        <-chan *ReceivedRequest
}

func runServer(ctx context.Context) (uint32, <-chan *ReceivedRequest) {
	bodychan := make(chan *ReceivedRequest, 100)
	handlerfunc := func(rw http.ResponseWriter, r *http.Request) {
		var rr ReceivedRequest
		rr.Method = r.Method
		if r.Body != nil {
			body, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			rr.Body = body
		}
		bodychan <- &rr
	}

	port := uint32(1334)
	handler := http.HandlerFunc(handlerfunc)
	go func() {
		h := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
		go func() {
			if err := h.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					panic(err)
				}
			}
		}()

		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		h.Shutdown(ctx)
		cancel()

	}()
	return port, bodychan
}

func NewTestUpstream(ctx context.Context) *TestUpstream {

	backendport, responses := runServer(ctx)

	serviceSpec := service.UpstreamSpec{
		Hosts: []service.Host{{
			Addr: "localhost",
			Port: backendport,
		}},
	}

	u := &v1.Upstream{
		Name: "local", // TODO: randomize
		Type: "service",
		Spec: service.EncodeUpstreamSpec(serviceSpec),
	}

	return &TestUpstream{
		Upstream: u,
		C:        responses,
	}
}

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
