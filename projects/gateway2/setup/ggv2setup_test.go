package setup_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	envoycluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyendpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoylistener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_service_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	envoy_service_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	envoy_service_route_v3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	envoy_service_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	"github.com/go-logr/zapr"
	"github.com/solo-io/gloo/projects/gateway2/extensions"
	"github.com/solo-io/gloo/projects/gateway2/krtcollections"
	ggv2setup "github.com/solo-io/gloo/projects/gateway2/setup"
	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	gloosetup "github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/slices"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

func getAssetsDir(t *testing.T) string {
	assets := ""
	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		// set default if not user provided
		out, err := exec.Command("sh", "-c", "make -sC $(dirname $(go env GOMOD))/projects/gateway2 envtest-path").CombinedOutput()
		t.Log("out:", string(out))
		if err != nil {
			t.Fatalf("failed to get assets dir: %v", err)
		}
		assets = strings.TrimSpace(string(out))
	}
	return assets
}

// testingWriter is a WriteSyncer that writes logs to testing.T.
type testingWriter struct {
	t atomic.Value
}

func (w *testingWriter) Write(p []byte) (n int, err error) {
	w.t.Load().(*testing.T).Log(string(p)) // Write the log to testing.T
	return len(p), nil
}

func (w *testingWriter) Sync() error {
	return nil
}

func (w *testingWriter) set(t *testing.T) {
	w.t.Store(t)
}

var (
	writer = &testingWriter{}
	logger = NewTestLogger()
)

// NewTestLogger creates a zap.Logger which can be used to write to *testing.T
// on each test, set the *testing.T on the writer.
func NewTestLogger() *zap.Logger {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(writer),
		// Adjust log level as needed
		// if a test assertion fails and logs or too noisy, change to zapcore.FatalLevel
		zapcore.DebugLevel,
	)

	return zap.New(core, zap.AddCaller())
}

func init() {
	log.SetLogger(zapr.NewLogger(logger))

}

func TestScenarios(t *testing.T) {
	writer.set(t)
	os.Setenv("POD_NAMESPACE", "gwtest")
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "crds"),
			filepath.Join("..", "..", "..", "install", "helm", "gloo", "crds"),
			filepath.Join("testdata", "istiocrds"),
		},
		ErrorIfCRDPathMissing: true,
		// set assets dir so we can run without the makefile
		BinaryAssetsDirectory: getAssetsDir(t),
		// web hook to add cluster ips to services
	}
	var wg sync.WaitGroup
	t.Cleanup(wg.Wait)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = contextutils.WithExistingLogger(ctx, logger.Sugar())

	cfg, err := testEnv.Start()
	if err != nil {
		t.Fatalf("failed to get assets dir: %v", err)
	}
	t.Cleanup(func() { testEnv.Stop() })

	kubeconfig := generateKubeConfiguration(t, cfg)
	t.Log("kubeconfig:", kubeconfig)

	client, err := istiokube.NewCLIClient(istiokube.NewClientConfigForRestConfig(cfg))
	if err != nil {
		t.Fatalf("failed to get init kube client: %v", err)
	}

	// apply settings/gwclass to the cluster
	err = client.ApplyYAMLFiles("default", "testdata/setupyaml/setup.yaml")
	if err != nil {
		t.Fatalf("failed to apply yaml: %v", err)
	}

	// create the test ns
	_, err = client.Kube().CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "gwtest"}}, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("failed to create namespace: %v", err)
	}

	err = client.ApplyYAMLFiles("gwtest", "testdata/setupyaml/pods.yaml")
	if err != nil {
		t.Fatalf("failed to apply yaml: %v", err)
	}

	// setup xDS server:
	uniqueClientCallbacks, builder := krtcollections.NewUniquelyConnectedClients()
	setupOpts := bootstrap.NewSetupOpts(xds.NewAdsSnapshotCache(ctx), uniqueClientCallbacks)
	addr := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: int(0),
	}
	controlPlane := gloosetup.NewControlPlane(ctx, setupOpts.Cache, grpc.NewServer(), addr, bootstrap.KubernetesControlPlaneConfig{}, uniqueClientCallbacks, true)
	xds.SetupEnvoyXds(controlPlane.GrpcServer, controlPlane.XDSServer, controlPlane.SnapshotCache)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("cant listen %v", err)
	}
	xdsPort := lis.Addr().(*net.TCPAddr).Port
	setupOpts.SetXdsAddress("localhost", int32(xdsPort))

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-controlPlane.GrpcService.Ctx.Done()
		controlPlane.GrpcServer.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		controlPlane.GrpcServer.Serve(lis)
		t.Log("grpc server stopped")
	}()

	// start ggv2
	// (note: we don't have gloo-edge working, so nothing will reconcile the proxies.
	// that's mostly ok, as we don't test the features that require these proxies in gloo-edge)
	setupOpts.ProxyReconcileQueue = ggv2utils.NewAsyncQueue[gloov1.ProxyList]()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ggv2setup.StartGGv2WithConfig(ctx, setupOpts, cfg, builder, extensions.NewK8sGatewayExtensions,
			registry.GetPluginRegistryFactory,
			types.NamespacedName{Name: "default", Namespace: "default"},
		)
	}()
	// give ggv2 time to initialize so we don't get
	// "ggv2 not initialized" error
	// this means that it attaches the pod collection to the unique client set collection.
	time.Sleep(time.Second)

	// list all yamls in test data
	files, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	for _, f := range files {
		// run tests with the yaml files (but not -out.yaml files)/s
		parentT := t
		if strings.HasSuffix(f.Name(), ".yaml") && !strings.HasSuffix(f.Name(), "-out.yaml") {
			fullpath := filepath.Join("testdata", f.Name())
			t.Run(strings.TrimSuffix(f.Name(), ".yaml"), func(t *testing.T) {
				writer.set(t)
				t.Cleanup(func() {
					writer.set(parentT)
				})
				t.Cleanup(func() {
					if t.Failed() {
						j, _ := setupOpts.KrtDebugger.MarshalJSON()
						t.Logf("krt state for failed test: %s %s", t.Name(), string(j))
					}
				})
				//sadly tests can't run yet in parallel, as ggv2 will add all the k8s services as clusters. this means
				// that we get test pollution.
				// once we change it to only include the ones in the proxy, we can re-enable this
				//				t.Parallel()
				testScenario(t, ctx, client, xdsPort, fullpath)

			})
		}
	}
}

func testScenario(t *testing.T, ctx context.Context, client istiokube.CLIClient, xdsPort int, f string) {
	fext := filepath.Ext(f)
	fpre := strings.TrimSuffix(f, fext)
	fout := fpre + "-out" + fext
	// read the out file
	write := false
	ya, err := os.ReadFile(fout)
	// if not exist
	if os.IsNotExist(err) {
		write = true
		err = nil
	}
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var expectedXdsDump xdsDump
	err = expectedXdsDump.FromYaml(ya)
	if err != nil {
		t.Fatalf("failed to read yaml: %v", err)
	}
	const gwname = "http-gw-for-test"
	testgwname := "http-" + filepath.Base(fpre)
	testyamlbytes, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	// change the gw name, so we could potentially run multiple tests in parallel (tough currently
	// it has other issues, so we don't run them in parallel)
	testyaml := strings.ReplaceAll(string(testyamlbytes), gwname, testgwname)

	yamlfile := filepath.Join(t.TempDir(), "test.yaml")
	os.WriteFile(yamlfile, []byte(testyaml), 0644)

	err = client.ApplyYAMLFiles("", yamlfile)
	defer func() {
		// always delete yamls, even if there was an error applying them; to prevent test pollution.
		err := client.DeleteYAMLFiles("", yamlfile)
		if err != nil {
			t.Fatalf("failed to delete yaml: %v", err)
		}
		t.Log("deleted yamls", t.Name())
	}()
	if err != nil {
		t.Fatalf("failed to apply yaml: %v", err)
	}
	t.Log("applied yamls", t.Name())
	// make sure all yamls reached the control plane
	time.Sleep(time.Second)

	dump := getXdsDump(t, ctx, xdsPort, testgwname)

	if write {
		t.Logf("writing out file")
		// serialize xdsDump to yaml
		d, err := dump.ToYaml()
		if err != nil {
			t.Fatalf("failed to serialize xdsDump: %v", err)
		}
		os.WriteFile(fout, d, 0644)
		t.Fatal("wrote out file - nothing to test")
	}
	dump.Compare(t, expectedXdsDump)
	fmt.Println("test done")
}

func getXdsDump(t *testing.T, ctx context.Context, xdsPort int, gwname string) xdsDump {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", xdsPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect to xds server: %v", err)
	}
	defer conn.Close()

	f := xdsFetcher{
		conn: conn,
		dr: &discovery_v3.DiscoveryRequest{Node: &envoycore.Node{
			Id: "gateway.gwtest",
			Metadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{"role": {Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("gloo-kube-gateway-api~%s~%s-%s", "gwtest", "gwtest", gwname)}}}},
		}},
	}
	clusters := f.getclusters(t, ctx)
	clusterServiceNames := slices.MapFilter(clusters, func(c *envoycluster.Cluster) *string {
		if c.GetEdsClusterConfig() != nil {
			if c.GetEdsClusterConfig().GetServiceName() != "" {
				s := c.GetEdsClusterConfig().GetServiceName()
				return &s
			}
			return &c.Name
		}
		return nil
	})

	listeners := f.getlisteners(t, ctx)
	var routenames []string
	for _, l := range listeners {
		routenames = append(routenames, getroutesnames(l)...)
	}
	return xdsDump{
		Clusters:  clusters,
		Listeners: listeners,
		Endpoints: f.getendpoints(t, ctx, clusterServiceNames),
		Routes:    f.getroutes(t, ctx, routenames),
	}
}

type xdsDump struct {
	Clusters  []*envoycluster.Cluster
	Listeners []*envoylistener.Listener
	Endpoints []*envoyendpoint.ClusterLoadAssignment
	Routes    []*envoy_config_route_v3.RouteConfiguration
}

func (x *xdsDump) Compare(t *testing.T, other xdsDump) {
	if len(x.Clusters) != len(other.Clusters) {
		t.Errorf("expected %v clusters, got %v", len(other.Clusters), len(x.Clusters))
	}

	if len(x.Listeners) != len(other.Listeners) {
		t.Errorf("expected %v listeners, got %v", len(other.Listeners), len(x.Listeners))
	}
	if len(x.Endpoints) != len(other.Endpoints) {
		t.Errorf("expected %v endpoints, got %v", len(other.Endpoints), len(x.Endpoints))
	}
	if len(x.Routes) != len(other.Routes) {
		t.Errorf("expected %v routes, got %v", len(other.Routes), len(x.Routes))
	}

	clusterset := map[string]*envoycluster.Cluster{}
	for _, c := range x.Clusters {
		clusterset[c.Name] = c
	}
	for _, c := range other.Clusters {
		otherc := clusterset[c.Name]
		if otherc == nil {
			t.Errorf("cluster %v not found", c.Name)
			continue
		}
		if !proto.Equal(c, otherc) {
			t.Errorf("cluster %v not equal", c.Name)
		}
	}
	listenerset := map[string]*envoylistener.Listener{}
	for _, c := range x.Listeners {
		listenerset[c.Name] = c
	}
	for _, c := range other.Listeners {
		otherc := listenerset[c.Name]
		if otherc == nil {
			t.Errorf("listener %v not found", c.Name)
			continue
		}
		if !proto.Equal(c, otherc) {
			t.Errorf("listener %v not equal", c.Name)
		}
	}
	routeset := map[string]*envoy_config_route_v3.RouteConfiguration{}
	for _, c := range x.Routes {
		routeset[c.Name] = c
	}
	for _, c := range other.Routes {
		otherc := routeset[c.Name]
		if otherc == nil {
			t.Errorf("route %v not found", c.Name)
			continue
		}
		if !proto.Equal(c, otherc) {
			t.Errorf("route %v not equal: %v vs %v", c.Name, c, otherc)
		}
	}

	epset := map[string]*envoyendpoint.ClusterLoadAssignment{}
	for _, c := range x.Endpoints {
		epset[c.ClusterName] = c
	}
	for _, c := range other.Endpoints {
		otherc := epset[c.ClusterName]
		if otherc == nil {
			t.Errorf("ep %v not found", c.ClusterName)
			continue
		}
		ep1 := flattenendpoints(c)
		ep2 := flattenendpoints(otherc)
		if !equalset(ep1, ep2) {
			t.Errorf("ep list %v not equal: %v %v", c.ClusterName, ep1, ep2)
		}
		ce := c.Endpoints
		ocd := otherc.Endpoints
		c.Endpoints = nil
		otherc.Endpoints = nil
		if !proto.Equal(c, otherc) {
			t.Errorf("ep %v not equal", c.ClusterName)
		}
		c.Endpoints = ce
		otherc.Endpoints = ocd

	}
}

func equalset(a, b []*envoyendpoint.LocalityLbEndpoints) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if slices.FindFunc(b, func(e *envoyendpoint.LocalityLbEndpoints) bool {
			return proto.Equal(v, e)
		}) == nil {
			return false
		}
	}
	return true
}

func flattenendpoints(v *envoyendpoint.ClusterLoadAssignment) []*envoyendpoint.LocalityLbEndpoints {
	var flat []*envoyendpoint.LocalityLbEndpoints
	for _, e := range v.Endpoints {
		for _, l := range e.LbEndpoints {
			flatbase := proto.Clone(e).(*envoyendpoint.LocalityLbEndpoints)
			flatbase.LbEndpoints = []*envoyendpoint.LbEndpoint{l}
			flat = append(flat, flatbase)
		}
	}
	return flat
}

func (x *xdsDump) FromYaml(ya []byte) error {
	ya, err := yaml.YAMLToJSON(ya)
	if err != nil {
		return err
	}

	jsonM := map[string][]any{}
	err = json.Unmarshal(ya, &jsonM)
	if err != nil {
		return err
	}
	for _, c := range jsonM["clusters"] {
		r, err := anyJsonRoundTrip[envoycluster.Cluster](c)
		if err != nil {
			return err
		}
		x.Clusters = append(x.Clusters, r)
	}
	for _, c := range jsonM["endpoints"] {
		r, err := anyJsonRoundTrip[envoyendpoint.ClusterLoadAssignment](c)
		if err != nil {
			return err
		}
		x.Endpoints = append(x.Endpoints, r)
	}
	for _, c := range jsonM["listeners"] {
		r, err := anyJsonRoundTrip[envoylistener.Listener](c)
		if err != nil {
			return err
		}
		x.Listeners = append(x.Listeners, r)
	}
	for _, c := range jsonM["routes"] {
		r, err := anyJsonRoundTrip[envoy_config_route_v3.RouteConfiguration](c)
		if err != nil {
			return err
		}
		x.Routes = append(x.Routes, r)
	}
	return nil
}

func anyJsonRoundTrip[T any, PT interface {
	proto.Message
	*T
}](c any) (PT, error) {
	var ju jsonpb.UnmarshalOptions
	jb, err := json.Marshal(c)
	var zero PT
	if err != nil {
		return zero, err
	}
	var r T
	var pr PT = &r
	err = ju.Unmarshal(jb, pr)
	return pr, err
}

func sortResource[T fmt.Stringer](resources []T) []T {
	// clone the slice
	resources = append([]T(nil), resources...)
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].String() < resources[j].String()
	})
	return resources
}

func (x *xdsDump) ToYaml() ([]byte, error) {
	jsonM := map[string][]any{}
	for _, c := range sortResource(x.Clusters) {
		roundtrip, err := protoJsonRoundTrip(c)
		if err != nil {
			return nil, err
		}
		jsonM["clusters"] = append(jsonM["clusters"], roundtrip)
	}
	for _, c := range sortResource(x.Listeners) {
		roundtrip, err := protoJsonRoundTrip(c)
		if err != nil {
			return nil, err
		}
		jsonM["listeners"] = append(jsonM["listeners"], roundtrip)
	}
	for _, c := range sortResource(x.Endpoints) {
		roundtrip, err := protoJsonRoundTrip(c)
		if err != nil {
			return nil, err
		}
		jsonM["endpoints"] = append(jsonM["endpoints"], roundtrip)
	}
	for _, c := range sortResource(x.Routes) {
		roundtrip, err := protoJsonRoundTrip(c)
		if err != nil {
			return nil, err
		}
		jsonM["routes"] = append(jsonM["routes"], roundtrip)
	}

	bytes, err := json.Marshal(jsonM)
	if err != nil {
		return nil, err
	}

	ya, err := yaml.JSONToYAML(bytes)
	if err != nil {
		return nil, err
	}
	return ya, nil
}

func protoJsonRoundTrip(c proto.Message) (any, error) {
	var j jsonpb.MarshalOptions
	s, err := j.Marshal(c)
	if err != nil {
		return nil, err
	}
	var roundtrip any
	err = json.Unmarshal(s, &roundtrip)
	if err != nil {
		return nil, err
	}
	return roundtrip, nil
}

type xdsFetcher struct {
	conn *grpc.ClientConn
	dr   *discovery_v3.DiscoveryRequest
}

func (x *xdsFetcher) getclusters(t *testing.T, ctx context.Context) []*envoycluster.Cluster {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cds := envoy_service_cluster_v3.NewClusterDiscoveryServiceClient(x.conn)

	//give ggv2 time to initialize so we don't get
	// "ggv2 not initialized" error

	epcli, err := cds.StreamClusters(ctx)
	if err != nil {
		t.Fatalf("failed to get cds client: %v", err)
	}
	defer epcli.CloseSend()
	epcli.Send(x.dr)
	dresp, err := epcli.Recv()
	if err != nil {
		t.Fatalf("failed to get response from xds server: %v", err)
	}
	t.Logf("got cds response. nonce, version: %s %s", dresp.GetNonce(), dresp.GetVersionInfo())
	var clusters []*envoycluster.Cluster
	for _, anyCluster := range dresp.GetResources() {

		var cluster envoycluster.Cluster
		if err := anyCluster.UnmarshalTo(&cluster); err != nil {
			t.Fatalf("failed to unmarshal cluster: %v", err)
		}
		clusters = append(clusters, &cluster)
	}
	return clusters
}

func getroutesnames(l *envoylistener.Listener) []string {
	var routes []string
	for _, fc := range l.GetFilterChains() {
		for _, filter := range fc.GetFilters() {
			suffix := string((&envoyhttp.HttpConnectionManager{}).ProtoReflect().Descriptor().FullName())
			if strings.HasSuffix(filter.GetTypedConfig().GetTypeUrl(), suffix) {
				var hcm envoyhttp.HttpConnectionManager
				switch config := filter.GetConfigType().(type) {
				case *envoylistener.Filter_TypedConfig:
					if err := config.TypedConfig.UnmarshalTo(&hcm); err == nil {
						rds := hcm.GetRds().GetRouteConfigName()
						if rds != "" {
							routes = append(routes, rds)
						}
					}
				}
			}
		}
	}
	return routes
}

func (x *xdsFetcher) getlisteners(t *testing.T, ctx context.Context) []*envoylistener.Listener {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ds := envoy_service_listener_v3.NewListenerDiscoveryServiceClient(x.conn)

	//give ggv2 time to initialize so we don't get
	// "ggv2 not initialized" error

	epcli, err := ds.StreamListeners(ctx)
	if err != nil {
		t.Fatalf("failed to get lds client: %v", err)
	}
	defer epcli.CloseSend()
	epcli.Send(x.dr)
	dresp, err := epcli.Recv()
	if err != nil {
		t.Fatalf("failed to get response from xds server: %v", err)
	}
	t.Logf("got lds response. nonce, version: %s %s", dresp.GetNonce(), dresp.GetVersionInfo())
	var resources []*envoylistener.Listener
	for _, anyResource := range dresp.GetResources() {

		var resource envoylistener.Listener
		if err := anyResource.UnmarshalTo(&resource); err != nil {
			t.Fatalf("failed to unmarshal resource: %v", err)
		}
		resources = append(resources, &resource)
	}
	return resources
}

func (x *xdsFetcher) getendpoints(t *testing.T, ctx context.Context, clusterServiceNames []string) []*envoyendpoint.ClusterLoadAssignment {

	eds := envoy_service_endpoint_v3.NewEndpointDiscoveryServiceClient(x.conn)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	epcli, err := eds.StreamEndpoints(ctx)
	if err != nil {
		t.Fatalf("failed to get eds client: %v", err)
	}
	defer epcli.CloseSend()
	dr := proto.Clone(x.dr).(*discovery_v3.DiscoveryRequest)
	dr.ResourceNames = clusterServiceNames
	epcli.Send(dr)
	dresp, err := epcli.Recv()
	if err != nil {
		t.Fatalf("failed to get response from xds server: %v", err)
	}
	t.Logf("got eds response. nonce, version: %s %s", dresp.GetNonce(), dresp.GetVersionInfo())
	var clas []*envoyendpoint.ClusterLoadAssignment
	for _, anyCluster := range dresp.GetResources() {

		var cla envoyendpoint.ClusterLoadAssignment
		if err := anyCluster.UnmarshalTo(&cla); err != nil {
			t.Fatalf("failed to unmarshal cluster: %v", err)
		}
		// remove kube endpoints, as with envtests we will get random ports, so we cant assert on them
		if !strings.Contains(cla.ClusterName, "kube-svc:default-kubernetes") {
			clas = append(clas, &cla)
		}
	}
	return clas
}

func (x *xdsFetcher) getroutes(t *testing.T, ctx context.Context, rosourceNames []string) []*envoy_config_route_v3.RouteConfiguration {

	eds := envoy_service_route_v3.NewRouteDiscoveryServiceClient(x.conn)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	epcli, err := eds.StreamRoutes(ctx)
	if err != nil {
		t.Fatalf("failed to get rds client: %v", err)
	}
	defer epcli.CloseSend()
	dr := proto.Clone(x.dr).(*discovery_v3.DiscoveryRequest)
	dr.ResourceNames = rosourceNames
	epcli.Send(dr)
	dresp, err := epcli.Recv()
	if err != nil {
		t.Fatalf("failed to get response from xds server: %v", err)
	}
	t.Logf("got rds response. nonce, version: %s %s", dresp.GetNonce(), dresp.GetVersionInfo())
	var clas []*envoy_config_route_v3.RouteConfiguration
	for _, anyCluster := range dresp.GetResources() {

		var cla envoy_config_route_v3.RouteConfiguration
		if err := anyCluster.UnmarshalTo(&cla); err != nil {
			t.Fatalf("failed to unmarshal cluster: %v", err)
		}
		clas = append(clas, &cla)
	}
	return clas
}

func generateKubeConfiguration(t *testing.T, restconfig *rest.Config) string {
	clusters := make(map[string]*clientcmdapi.Cluster)
	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	contexts := make(map[string]*clientcmdapi.Context)

	clusterName := "cluster"
	clusters[clusterName] = &clientcmdapi.Cluster{
		Server:                   restconfig.Host,
		CertificateAuthorityData: restconfig.CAData,
	}
	authinfos[clusterName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         restconfig.KeyData,
		ClientCertificateData: restconfig.CertData,
	}
	contexts[clusterName] = &clientcmdapi.Context{
		Cluster:   clusterName,
		Namespace: "default",
		AuthInfo:  clusterName,
	}

	clientConfig := clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters:   clusters,
		Contexts:   contexts,
		// current context must be mgmt cluster for now, as the api server doesn't have context configurable.
		CurrentContext: "cluster",
		AuthInfos:      authinfos,
	}
	// create temp file
	tmpfile := filepath.Join(t.TempDir(), "kubeconfig")
	err := clientcmd.WriteToFile(clientConfig, tmpfile)
	if err != nil {
		t.Fatalf("failed to write kubeconfig: %v", err)
	}

	return tmpfile
}
