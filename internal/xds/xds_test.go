package xds_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoylistener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoyhttpconnectionmanager "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/test/resource"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"

	"github.com/ghodss/yaml"
	. "github.com/solo-io/gloo-testing/helpers"
	. "github.com/solo-io/gloo/internal/xds"
	"github.com/solo-io/gloo/pkg/log"
	"google.golang.org/grpc"
)

var _ = Describe("Xds", func() {
	envoyBaseDir := os.Getenv("GOPATH") + "/src/github.com/solo-io/gloo"
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
		// fun times
		// write bootstrap file
		bootstrap := resource.MakeBootstrap(true, uint32(8081), 19000)
		buf = &bytes.Buffer{}
		err := (&jsonpb.Marshaler{OrigName: true}).Marshal(buf, bootstrap)
		Must(err)
		ym, err := yaml.JSONToYAML(buf.Bytes())
		Must(err)
		fmt.Printf("FOR MANUAL TESTING PURPOSES:\n%s\n%v\n", ym, envoyRunArgs)

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

		snapshot, err := createSnapshot(routeConfigName, listenerName)
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
				str := string(buf.Bytes())
				return str
			}, time.Second*10).Should(ContainSubstring("lds: add/update listener '" + listenerName))
		})
	})
})

func createSnapshot(routeConfigName, listenerName string) (cache.Snapshot, error) {
	var (
		endpoints []proto.Message
		clusters  []proto.Message
		routes    []proto.Message
	)
	adsSource := envoycore.ConfigSource{
		ConfigSourceSpecifier: &envoycore.ConfigSource_Ads{
			Ads: &envoycore.AggregatedConfigSource{},
		},
	}
	manager := &envoyhttpconnectionmanager.HttpConnectionManager{
		CodecType:  envoyhttpconnectionmanager.AUTO,
		StatPrefix: "http",
		RouteSpecifier: &envoyhttpconnectionmanager.HttpConnectionManager_Rds{
			Rds: &envoyhttpconnectionmanager.Rds{
				ConfigSource:    adsSource,
				RouteConfigName: routeConfigName,
			},
		},
	}
	pbst, err := util.MessageToStruct(manager)
	if err != nil {
		panic(err)
	}

	listener := &v2.Listener{
		Name: listenerName,
		Address: envoycore.Address{
			Address: &envoycore.Address_SocketAddress{
				SocketAddress: &envoycore.SocketAddress{
					Protocol: envoycore.TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &envoycore.SocketAddress_PortValue{
						PortValue: 1234,
					},
				},
			},
		},
		FilterChains: []envoylistener.FilterChain{{
			Filters: []envoylistener.Filter{{
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
