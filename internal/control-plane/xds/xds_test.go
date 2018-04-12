package xds_test

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoylistener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoyhttpconnectionmanager "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/jsonpb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"

	bootstrap "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	"github.com/ghodss/yaml"
	. "github.com/solo-io/gloo/internal/control-plane/xds"
	"github.com/solo-io/gloo/pkg/log"
	. "github.com/solo-io/gloo/test/helpers"
	"google.golang.org/grpc"
)

const staticEnvoyConfig = `
node:
  cluster: ingress
  id: some-id

static_resources:
  clusters:

  - name: xds_cluster
    connect_timeout: 5.000s
    hosts:
    - socket_address:
        address: 127.0.0.1
        port_value: 8081
    http2_protocol_options: {}
    type: STRICT_DNS

dynamic_resources:
  ads_config:
    api_type: GRPC
    grpc_services:
    - envoy_grpc: {cluster_name: xds_cluster}
  cds_config:
    ads: {}
  lds_config:
    ads: {}

admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
`

var _ = Describe("Xds", func() {
	envoyRunArgs := []string{
		"docker", "run", "--rm",
		"--name", "one-at-a-time",
		"--network", "host",
		"soloio/envoy:v0.1.2",
		"/bin/sh", "-c",
		"\"echo", "'" + staticEnvoyConfig + "'", ">", "/envoy.yaml", "&&",
		"envoy",
		"-c", "/envoy.yaml",
		"--v2-config-only\"",
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
		bootstrap := makeBootstrap(true, uint32(8081), 19000)
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
		exec.Command("docker", "kill", "one-at-a-time").Run()
	})
	Describe("RunXDS Server", func() {
		It("successfully bootstraps the envoy proxy", func() {
			Eventually(func() string {
				str := string(buf.Bytes())
				return str
			}, time.Second*45).Should(ContainSubstring("lds: add/update listener '" + listenerName))
		})
	})
})

func createSnapshot(routeConfigName, listenerName string) (cache.Snapshot, error) {
	var (
		endpoints []cache.Resource
		clusters  []cache.Resource
		routes    []cache.Resource
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

	listeners := []cache.Resource{
		listener,
	}
	return cache.NewSnapshot("foo", endpoints, clusters, routes, listeners), nil
}

// MakeBootstrap creates a bootstrap envoy configuration
func makeBootstrap(ads bool, controlPort, adminPort uint32) *bootstrap.Bootstrap {
	source := &envoycore.ApiConfigSource{
		ApiType:      envoycore.ApiConfigSource_GRPC,
		ClusterNames: []string{"xds_cluster"},
	}

	var dynamic *bootstrap.Bootstrap_DynamicResources
	if ads {
		dynamic = &bootstrap.Bootstrap_DynamicResources{
			LdsConfig: &envoycore.ConfigSource{
				ConfigSourceSpecifier: &envoycore.ConfigSource_Ads{Ads: &envoycore.AggregatedConfigSource{}},
			},
			CdsConfig: &envoycore.ConfigSource{
				ConfigSourceSpecifier: &envoycore.ConfigSource_Ads{Ads: &envoycore.AggregatedConfigSource{}},
			},
			AdsConfig: source,
		}
	} else {
		dynamic = &bootstrap.Bootstrap_DynamicResources{
			LdsConfig: &envoycore.ConfigSource{
				ConfigSourceSpecifier: &envoycore.ConfigSource_ApiConfigSource{ApiConfigSource: source},
			},
			CdsConfig: &envoycore.ConfigSource{
				ConfigSourceSpecifier: &envoycore.ConfigSource_ApiConfigSource{ApiConfigSource: source},
			},
		}
	}

	return &bootstrap.Bootstrap{
		Node: &envoycore.Node{
			Id:      "test-id",
			Cluster: "test-cluster",
		},
		Admin: bootstrap.Admin{
			AccessLogPath: "/dev/null",
			Address: envoycore.Address{
				Address: &envoycore.Address_SocketAddress{
					SocketAddress: &envoycore.SocketAddress{
						Address: "127.0.0.1",
						PortSpecifier: &envoycore.SocketAddress_PortValue{
							PortValue: adminPort,
						},
					},
				},
			},
		},
		StaticResources: &bootstrap.Bootstrap_StaticResources{
			Clusters: []v2.Cluster{{
				Name:           "xds_cluster",
				ConnectTimeout: 5 * time.Second,
				Type:           v2.Cluster_STATIC,
				Hosts: []*envoycore.Address{{
					Address: &envoycore.Address_SocketAddress{
						SocketAddress: &envoycore.SocketAddress{
							Address: "127.0.0.1",
							PortSpecifier: &envoycore.SocketAddress_PortValue{
								PortValue: controlPort,
							},
						},
					},
				}},
				LbPolicy:             v2.Cluster_ROUND_ROBIN,
				Http2ProtocolOptions: &envoycore.Http2ProtocolOptions{},
			}},
		},
		DynamicResources: dynamic,
	}
}
