package e2e_test

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/gloo/test/services/envoy"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/resource"
)

// An EDS subscription is satisfied by name. A proxy asks for the names its
// clusters reference and matches what comes back by the name each resource
// gives for itself, so a response carrying a resource under the right map key
// while identifying itself as something else satisfies the cache and not the
// proxy. It is delivered, ignored, and never retried, because from the cache's
// side the subscription has been answered.
//
// The unit specs in projects/gloo/pkg/xds assert that MakeConsistent produces a
// correctly named placeholder. They cannot assert the consequence, because the
// consequence belongs to Envoy: whether it applies the assignment and drops the
// endpoints it was holding. That is what this spec measures, against a real
// Envoy, by reading the runtime host count out of the admin interface.
//
// Reverting the fix -- naming the placeholder after the cluster rather than the
// resource it is filed under -- leaves Envoy on its previous endpoints and
// fails the final assertion.
var _ = Describe("EDS placeholder assignments", Label(), func() {

	const (
		role        = "placeholder~test~proxy"
		clusterName = "placeholder-test-cluster"
		// Gloo serves EDS under "<cluster>-<upstream hash>", never under the
		// cluster's own name. The two differing is the entire point: a
		// placeholder that names the cluster names something nobody subscribed
		// to. See translator.GetEndpointClusterName.
		serviceName = clusterName + "-12345"
	)

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		envoyInst   *envoy.Instance
		snapshots   cache.SnapshotCache
		grpcServer  *grpc.Server
		membersRgx  = regexp.MustCompile(`membership_total:\s*(\d+)`)
		snapshotFor func(withAssignment bool) cache.Snapshot
	)

	// membership reads the runtime host count Envoy holds for the cluster.
	// This is deliberately read from /clusters rather than /config_dump: the
	// endpoint config dump is built from the same host sets, so it cannot
	// distinguish what was delivered from what is merely still held.
	membership := func() int {
		stats, err := envoyInst.AdminClient().GetStats(ctx, map[string]string{
			"filter": fmt.Sprintf("^cluster\\.%s\\.membership_total$", regexp.QuoteMeta(clusterName)),
		})
		if err != nil {
			return -1
		}
		match := membersRgx.FindStringSubmatch(stats)
		if match == nil {
			return -1
		}
		count, err := strconv.Atoi(match[1])
		if err != nil {
			return -1
		}
		return count
	}

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		snapshotFor = func(withAssignment bool) cache.Snapshot {
			cluster := resource.NewEnvoyResource(&envoy_config_cluster_v3.Cluster{
				Name:                 clusterName,
				ConnectTimeout:       durationpb.New(5 * time.Second),
				ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{Type: envoy_config_cluster_v3.Cluster_EDS},
				EdsClusterConfig: &envoy_config_cluster_v3.Cluster_EdsClusterConfig{
					EdsConfig: &envoy_config_core_v3.ConfigSource{
						ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
						ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
							Ads: &envoy_config_core_v3.AggregatedConfigSource{},
						},
					},
					ServiceName: serviceName,
				},
			})

			endpoints := []cache.Resource{}
			if withAssignment {
				endpoints = append(endpoints, resource.NewEnvoyResource(
					&envoy_config_endpoint_v3.ClusterLoadAssignment{
						ClusterName: serviceName,
						Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{{
							LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
								lbEndpoint("10.0.0.1", 8080),
								lbEndpoint("10.0.0.2", 8080),
							},
						}},
					}))
			}

			version := "with-assignment"
			if !withAssignment {
				version = "assignment-withdrawn"
			}
			snap := xds.NewSnapshot(version,
				endpoints,
				[]cache.Resource{cluster},
				[]cache.Resource{},
				[]cache.Resource{})
			// The control plane repairs a snapshot before publishing it; this
			// is the call under test.
			snap.MakeConsistent()
			return snap
		}

		snapshots = xds.NewAdsSnapshotCache(ctx)
		grpcServer = grpc.NewServer()

		// Bind first and read the port back, so parallel specs cannot race for
		// the same number.
		listener, err := net.Listen("tcp", ":0")
		Expect(err).NotTo(HaveOccurred())
		port := listener.Addr().(*net.TCPAddr).Port

		controlPlane := setup.NewControlPlane(ctx, snapshots, grpcServer, &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: port,
		}, bootstrap.KubernetesControlPlaneConfig{}, nil, true)
		// NewControlPlane builds the xDS server but does not attach it; without
		// this the proxy connects to a gRPC server that serves no discovery
		// service and silently receives nothing.
		xds.SetupEnvoyXds(grpcServer, controlPlane.XDSServer, snapshots, false)
		go func() {
			defer GinkgoRecover()
			_ = grpcServer.Serve(listener)
		}()

		snapshots.SetSnapshot(role, snapshotFor(true))

		envoyInst = envoyFactory.NewInstance()
		Expect(envoyInst.RunWithRole(role, port)).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		envoyInst.Clean()
		grpcServer.Stop()
		cancel()
	})

	It("drops the endpoints when the control plane withdraws the assignment", func() {
		// Control. If this does not hold, the rest of the spec says nothing:
		// it would mean Envoy never accepted the cluster or never connected.
		Eventually(membership, "20s", "500ms").Should(Equal(2),
			"Envoy should hold both published endpoints before anything is withdrawn")

		// Withdraw the assignment. MakeConsistent sees a cluster referencing an
		// EDS name nothing publishes and substitutes an empty placeholder.
		snapshots.SetSnapshot(role, snapshotFor(false))

		// The placeholder has to name the resource it is filed under, or Envoy
		// cannot match it to the subscription and silently keeps what it had.
		Eventually(membership, "20s", "500ms").Should(Equal(0),
			"Envoy should apply the withdrawal; holding endpoints here means the "+
				"assignment named something the proxy never subscribed to")
	})
})

func lbEndpoint(address string, port uint32) *envoy_config_endpoint_v3.LbEndpoint {
	return &envoy_config_endpoint_v3.LbEndpoint{
		HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
			Endpoint: &envoy_config_endpoint_v3.Endpoint{
				Address: &envoy_config_core_v3.Address{
					Address: &envoy_config_core_v3.Address_SocketAddress{
						SocketAddress: &envoy_config_core_v3.SocketAddress{
							Address:       address,
							PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{PortValue: port},
						},
					},
				},
			},
		},
	}
}
