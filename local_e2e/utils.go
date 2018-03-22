package local_e2e

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"context"
	"io/ioutil"
	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
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
		h := &http.Server{Handler: handler}
		go func() {
			if err := h.Serve(listener); err != nil {
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
	return uint32(port), bodychan
}

var id = 0

func NewTestUpstream(ctx context.Context) *TestUpstream {

	backendport, responses := runServer(ctx)

	serviceSpec := service.UpstreamSpec{
		Hosts: []service.Host{{
			Addr: "localhost",
			Port: backendport,
		}},
	}
	id += 1
	u := &v1.Upstream{
		Name: fmt.Sprintf("local-%d", id),
		Type: "service",
		Spec: service.EncodeUpstreamSpec(serviceSpec),
	}

	return &TestUpstream{
		Upstream: u,
		C:        responses,
	}
}
