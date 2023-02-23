package v1helpers

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/solo-io/gloo/test/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/proto"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	static_plugin_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ReceivedRequest struct {
	Method      string
	Body        []byte
	Host        string
	Headers     http.Header
	URL         *url.URL
	GRPCRequest proto.Message
}

func NewTestHttpUpstream(ctx context.Context, addr string) *TestUpstream {
	return NewTestHttpUpstreamWithHandler(ctx, addr, nil)
}

//func NewTestHttpUpstreamWithHandler(ctx context.Context, addr string, httpHandler http.Handler) *TestUpstream {
//	backendport, responses := RunTestServer(ctx, &HttpServer{}, handlerFunc)
//	return newTestUpstream(addr, backendport, responses)
//}

func NewTestHttpUpstreamWithHandler(ctx context.Context, addr string, handlerFunc ExtraHandlerFunc) *TestUpstream {
	backendport, responses := RunTestServer(ctx, &HttpServer{}, handlerFunc)
	return newTestUpstream(addr, backendport, responses)
}

func NewTestHttpsUpstream(ctx context.Context, addr string) ([]byte, *TestUpstream) {
	return NewTestHttpsUpstreamWithHandler(ctx, addr, nil)
}

func NewTestHttpsUpstreamWithHandler(ctx context.Context, addr string, handlerFunc ExtraHandlerFunc) ([]byte, *TestUpstream) {
	cert, privKey := helpers.GetCerts(helpers.Params{
		Hosts:      "127.0.0.1",
		IsCA:       true,
		EcdsaCurve: "P256",
	})

	serverCert, err := tls.X509KeyPair([]byte(cert), []byte(privKey))
	Expect(err).NotTo(HaveOccurred())

	server := &HttpsServer{
		tlsConfig: &tls.Config{
			Certificates: []tls.Certificate{serverCert},
		},
	}

	backendport, responses := RunTestServer(ctx, server, handlerFunc)
	// Return the ca cert to be used in https client tls
	return []byte(cert), newTestUpstream(addr, backendport, responses)
}

type TestUpstream struct {
	Upstream *gloov1.Upstream
	C        <-chan *ReceivedRequest
	Address  string
	Port     uint32
}

var id = 0

func newTestUpstream(addr string, port uint32, responses <-chan *ReceivedRequest) *TestUpstream {

	id += 1
	u := &gloov1.Upstream{
		Metadata: &core.Metadata{
			Name:      fmt.Sprintf("local-%d", id),
			Namespace: "default",
		},
		UpstreamType: &gloov1.Upstream_Static{
			Static: &static_plugin_gloo.UpstreamSpec{
				Hosts: []*static_plugin_gloo.Host{{
					Addr: addr,
					Port: port,
				}},
			},
		},
	}

	return &TestUpstream{
		Upstream: u,
		C:        responses,
		Address:  fmt.Sprintf("%s:%d", addr, port),
		Port:     port,
	}
}

// Handler func, if returns false, the test upstream will `not` echo the request body in the response
type ExtraHandlerFunc func(w http.ResponseWriter, req *http.Request) bool

func RunTestServer(ctx context.Context, server Server, extraHandler ExtraHandlerFunc) (uint32, <-chan *ReceivedRequest) {
	bodychan := make(chan *ReceivedRequest, 100)
	handlerfunc := func(rw http.ResponseWriter, r *http.Request) {
		var rr ReceivedRequest
		rr.Method = r.Method
		echoBody := true
		if extraHandler != nil {
			echoBody = extraHandler(rw, r)
		}
		if r.Body != nil {
			body, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			if len(body) != 0 {
				rr.Body = body
				if echoBody {
					rw.Write(body)
				}
			}
			rr.Headers = r.Header
			rr.URL = r.URL
		}

		rr.Host = r.Host

		bodychan <- &rr
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	addr := listener.Addr().String()
	_, portstr, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		panic(err)
	}

	handler := http.HandlerFunc(handlerfunc)
	go func() {
		defer GinkgoRecover()
		h := &http.Server{
			Handler: handler,
		}
		go server.Serve(h, listener)

		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		h.Shutdown(ctx)
		cancel()

		// drain and close channel
		// the http handler may panic but this should be caught by the http code.
		for range bodychan {
		}
		close(bodychan)
	}()
	return uint32(port), bodychan
}

type Server interface {
	Serve(h *http.Server, listener net.Listener)
}

type HttpServer struct {
}

func (s *HttpServer) Serve(h *http.Server, listener net.Listener) {
	defer GinkgoRecover()
	if err := h.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

type HttpsServer struct {
	tlsConfig *tls.Config
}

func (s *HttpsServer) Serve(h *http.Server, listener net.Listener) {
	defer GinkgoRecover()
	h.TLSConfig = s.tlsConfig
	if err := h.ServeTLS(listener, "", ""); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
