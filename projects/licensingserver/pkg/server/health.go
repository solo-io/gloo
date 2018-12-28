package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthChecker struct {
	grpc     *health.Server
	ok       uint32
	router   *mux.Router
	listener net.Listener
	port     int
	ctx      context.Context
}

func newHealthChecker(grpcServer *grpc.Server, port int, ctx context.Context) *healthChecker {
	ret := &healthChecker{}
	ret.ok = 1
	ret.port = port

	ret.ctx = ctx

	ret.grpc = health.NewServer()
	ret.grpc.SetServingStatus(service_name, healthpb.HealthCheckResponse_SERVING)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)

	go func() {
		<-sigterm
		atomic.StoreUint32(&ret.ok, 0)
		ret.grpc.SetServingStatus(service_name, healthpb.HealthCheckResponse_NOT_SERVING)
	}()

	// setup http router
	ret.router = mux.NewRouter()
	ret.router.Path(health_check_endpoint).Handler(ret)

	healthpb.RegisterHealthServer(grpcServer, ret.grpc)

	return ret
}

func (hc *healthChecker) start() error {
	// start healthcheck server
	addr := fmt.Sprintf(":%d", hc.port)
	contextutils.LoggerFrom(hc.ctx).Infof("Listening for HTTP on '%s'", addr)
	var err error
	hc.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to open HTTP listener: '%+v'", err)
	}
	return http.Serve(hc.listener, hc.router)
}

func (hc *healthChecker) close() {
	if hc.listener != nil {
		hc.listener.Close()
	}
}

func (hc *healthChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ok := atomic.LoadUint32(&hc.ok)
	if ok == 1 {
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(500)
	}
}

func (hc *healthChecker) Fail() {
	atomic.StoreUint32(&hc.ok, 0)
	hc.grpc.SetServingStatus(service_name, healthpb.HealthCheckResponse_NOT_SERVING)
}
