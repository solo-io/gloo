package v1helpers

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
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

func NewTestHttpUpstreamWithHandler(ctx context.Context, addr string, handlerFunc ExtraHandlerFunc) *TestUpstream {
	backendport, responses := RunTestServer(ctx, &HttpServer{}, handlerFunc)
	return newTestUpstream(addr, backendport, responses)
}

func NewTestHttpsUpstream(ctx context.Context, addr string) ([]byte, *TestUpstream) {
	return NewTestHttpsUpstreamWithHandler(ctx, addr, nil)
}

func NewTestHttpsUpstreamWithHandler(ctx context.Context, addr string, handlerFunc ExtraHandlerFunc) ([]byte, *TestUpstream) {
	certSubject := pkix.Name{
		Organization: []string{"solo.io"},
		Country:      []string{"US"},
		Province:     []string{""},
		Locality:     []string{"Cambridge"},
	}
	// Generate Certificates for TLS
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               certSubject,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	Expect(err).NotTo(HaveOccurred())
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	Expect(err).NotTo(HaveOccurred())

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	Expect(err).NotTo(HaveOccurred())
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject:      certSubject,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	Expect(err).NotTo(HaveOccurred())

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	Expect(err).NotTo(HaveOccurred())

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	Expect(err).NotTo(HaveOccurred())

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	Expect(err).NotTo(HaveOccurred())

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	Expect(err).NotTo(HaveOccurred())

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPEM.Bytes())
	server := &HttpsServer{
		tlsConfig: serverTLSConf,
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caPEM.Bytes())
	Expect(ok).To(BeTrue())
	backendport, responses := RunTestServer(ctx, server, handlerFunc)
	// Return the ca cert to be used in https client tls
	return caPEM.Bytes(), newTestUpstream(addr, backendport, responses)
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
type ExtraHandlerFunc func(rw http.ResponseWriter, r *http.Request) bool

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
