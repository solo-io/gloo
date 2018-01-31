package xds_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/envoyproxy/go-control-plane/api/filter/network"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"

	. "github.com/solo-io/glue/internal/xds"
	"github.com/solo-io/glue/pkg/log"
	. "github.com/solo-io/glue/test/helpers"
)

var _ = Describe("Xds", func() {
	envoyBaseDir := os.Getenv("GOPATH") + "/src/github.com/solo-io/glue"
	envoyRunArgs := []string{
		filepath.Join(envoyBaseDir, "envoy"),
		"-c", filepath.Join(envoyBaseDir, "envoy.yaml"),
		"--v2-config-only",
		"--service-cluster", "envoy",
		"--service-node", "envoy",
	}
	var (
		envoyPid int
		buf      *bytes.Buffer
		srv      *grpc.Server

		routeConfigName = "xds-test-route-config"
		listenerName    = "xds-test-listener"
	)
	BeforeEach(func() {
		// setup
		envoyCmd := exec.Command("/bin/sh", "-c", strings.Join(envoyRunArgs, " "))
		buf = &bytes.Buffer{}
		envoyCmd.Stdout = buf
		envoyCmd.Stderr = buf
		go func() {
			if err := envoyCmd.Run(); err != nil {
				log.Fatalf(buf.String() + ": " + err.Error())
			}
		}()
		for envoyCmd.Process == nil {
			time.Sleep(time.Millisecond)
		}
		envoyPid = envoyCmd.Process.Pid
		cache, grpcSrv, err := RunXDS(8081)
		Must(err)
		srv = grpcSrv

		snapshot, err := createSnapshot("xds_cluster", routeConfigName, listenerName)
		if err != nil {
			log.Fatalf(err.Error())
		}
		cache.SetSnapshot(NodeKey, snapshot)
	})
	AfterEach(func() {
		srv.Stop()
		if err := syscall.Kill(envoyPid, syscall.SIGKILL); err != nil {
			log.Fatalf(err.Error())
		}
	})
	Describe("RunXDS Server", func() {
		It("successfully bootstraps the envoy proxy", func() {
			Eventually(func() string {
				str := buf.String()
				return str
			}, time.Second*10).Should(ContainSubstring("lds: add/update listener '" + listenerName))
		})
	})
})

func createSnapshot(xdsCluster, routeConfigName, listenerName string) (cache.Snapshot, error) {
	var (
		endpoints []proto.Message
		clusters  []proto.Message
		routes    []proto.Message
	)
	rdsSource := api.ConfigSource{
		ConfigSourceSpecifier: &api.ConfigSource_ApiConfigSource{
			ApiConfigSource: &api.ApiConfigSource{
				ApiType:      api.ApiConfigSource_GRPC,
				ClusterNames: []string{xdsCluster},
				GrpcServices: []*api.GrpcService{
					{
						TargetSpecifier: &api.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &api.GrpcService_EnvoyGrpc{
								ClusterName: xdsCluster,
							},
						},
					},
				},
			},
		},
	}
	manager := &network.HttpConnectionManager{
		CodecType:  network.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &network.HttpConnectionManager_Rds{
			Rds: &network.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: routeConfigName,
			},
		},
	}
	pbst, err := util.MessageToStruct(manager)
	if err != nil {
		panic(err)
	}

	listener := &api.Listener{
		Name: listenerName,
		Address: &api.Address{
			Address: &api.Address_SocketAddress{
				SocketAddress: &api.SocketAddress{
					Protocol: api.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &api.SocketAddress_PortValue{
						PortValue: 1234,
					},
				},
			},
		},
		FilterChains: []*api.FilterChain{{
			Filters: []*api.Filter{{
				Name:   "envoy.http_connection_manager",
				Config: pbst,
			}},
		}},
	}

	listeners := []proto.Message{
		listener,
	}
	return cache.NewSnapshot("foo", endpoints, clusters, routes, listeners), nil
}
