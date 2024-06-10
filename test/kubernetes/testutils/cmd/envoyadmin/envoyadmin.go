package envoyadmin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/protoutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/testutils/cmd"
	"github.com/solo-io/go-utils/threadsafe"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConfigDumpPath     = admincli.ConfigDumpPath
	StatsPath          = admincli.StatsPath
	ClustersPath       = admincli.ClustersPath
	ListenersPath      = admincli.ListenersPath
	ModifyRuntimePath  = admincli.ModifyRuntimePath
	ShutdownServerPath = admincli.ShutdownServerPath
	HealthCheckPath    = admincli.HealthCheckPath
	LoggingPath        = admincli.LoggingPath
	ServerInfoPath     = admincli.ServerInfoPath
	ReadyPath          = admincli.ReadyPath

	DefaultAdminPort = admincli.DefaultAdminPort
)

type Client struct {
	// receiver is the default destination for the curl stdout and stderr
	receiver io.Writer

	// curlOptions is the set of default Option that the Client will use for curl commands
	curlOptions []curl.Option

	envoyMeta metav1.ObjectMeta

	// kubeContext is the context against which the kubectl under the hood will be executed
	kubeContext string
}

// NewClient returns an instance of the remote envoyadmin.Client
func NewClient() *Client {
	return &Client{
		receiver: io.Discard,
		curlOptions: []curl.Option{
			curl.WithScheme("http"),
			// We are using local host because we will exec from the pod in which envoy is running
			curl.WithHost("127.0.0.1"),
			curl.WithPort(DefaultAdminPort),
			// 3 retries, exponential back-off, 10 second max
			curl.WithRetries(3, 0, 10),
		},
		envoyMeta: metav1.ObjectMeta{
			Name:      "gateway-proxy",
			Namespace: "gloo-system",
		},
	}
}

// WithReceiver sets the io.Writer that will be used by default for the stdout and stderr
// of cmdutils.Cmd created by the Client
func (c *Client) WithReceiver(receiver io.Writer) *Client {
	c.receiver = receiver
	return c
}

// WithCurlOptions sets the default set of curl.Option that will be used by default with
// the cmdutils.Cmd created by the Client
func (c *Client) WithCurlOptions(options ...curl.Option) *Client {
	c.curlOptions = append(c.curlOptions, options...)
	return c
}

// WithEnvoyMeta sets the Envoy pod that will be used by default to exec
// the cmdutils.Cmd created by the Client
func (c *Client) WithEnvoyMeta(envoyMeta metav1.ObjectMeta) *Client {
	c.envoyMeta = envoyMeta
	return c
}

// WithKubeContext sets the Kubernetes context that will be used by default to exec
// the cmdutils.Cmd created by the Client
func (c *Client) WithKubeContext(kubeContext string) *Client {
	c.kubeContext = kubeContext
	return c
}

// Command returns a remote curl Command, using the provided curl.Option as well as the client.curlOptions.
// The returned command will be executed from within the envoy pod specified in client.envoyMeta, and
// therefore changing the curl options related to address should rarely be necessary.
func (c *Client) Command(ctx context.Context, options ...curl.Option) cmdutils.Cmd {
	commandCurlOptions := append(
		c.curlOptions,
		// Ensure any options defined for this command can override any defaults that the Client has defined
		options...)

	return cmd.RemoteCurlCmd(ctx, c.receiver, c.kubeContext, kubectl.PodExecOptions{
		Name:      c.envoyMeta.Name,
		Namespace: c.envoyMeta.Namespace,
		Container: "curl",
	}, commandCurlOptions...).
		// For convenience, we set the stdout and stderr to the receiver
		// This can still be overwritten by consumers who use the commands
		WithStdout(c.receiver).
		WithStderr(c.receiver)
}

// RunCommand executes a curl Command, using the provided curl.Option as well as the client.curlOptions
func (c *Client) RunCommand(ctx context.Context, options ...curl.Option) error {
	return c.Command(ctx, options...).Run().Cause()
}

// RequestPathCmd returns the cmdutils.Cmd that can be run, and will execute a request against the provided path
func (c *Client) RequestPathCmd(ctx context.Context, path string) cmdutils.Cmd {
	return c.Command(ctx, curl.WithPath(path))
}

// StatsCmd returns the cmdutils.Cmd that can be run to request data from the stats endpoint
func (c *Client) StatsCmd(ctx context.Context) cmdutils.Cmd {
	return c.RequestPathCmd(ctx, StatsPath)
}

// GetStats returns the data that is available at the stats endpoint
func (c *Client) GetStats(ctx context.Context) (string, error) {
	var outLocation threadsafe.Buffer

	err := c.StatsCmd(ctx).WithStdout(&outLocation).Run().Cause()
	if err != nil {
		return "", err
	}

	return outLocation.String(), nil
}

// ClustersCmd returns the cmdutils.Cmd that can be run to request data from the clusters endpoint
func (c *Client) ClustersCmd(ctx context.Context) cmdutils.Cmd {
	return c.RequestPathCmd(ctx, ClustersPath)
}

// ListenersCmd returns the cmdutils.Cmd that can be run to request data from the listeners endpoint
func (c *Client) ListenersCmd(ctx context.Context) cmdutils.Cmd {
	return c.RequestPathCmd(ctx, ListenersPath)
}

// ConfigDumpCmd returns the cmdutils.Cmd that can be run to request data from the config_dump endpoint
func (c *Client) ConfigDumpCmd(ctx context.Context, queryParams map[string]string) cmdutils.Cmd {
	return c.Command(ctx,
		curl.WithPath(ConfigDumpPath),
		curl.WithQueryParameters(queryParams),
	)
}

// ReadyCmd returns the cmdutils.Cmd that can be run to request data from the ready endpoint
func (c *Client) ReadyCmd(ctx context.Context) cmdutils.Cmd {
	return c.RequestPathCmd(ctx, ReadyPath)
}

// GetReady returns true if the Envoy server is responding 200 LIVE on the ready endpoint
func (c *Client) GetReady(ctx context.Context) bool {
	var (
		outLocation threadsafe.Buffer
	)

	err := c.ReadyCmd(ctx).WithStdout(&outLocation).Run().Cause()
	if err != nil {
		return false
	}
	return strings.Contains(outLocation.String(), "LIVE")
}

// GetConfigDump returns the structured data that is available at the config_dump endpoint
func (c *Client) GetConfigDump(ctx context.Context, queryParams map[string]string) (*adminv3.ConfigDump, error) {
	var (
		cfgDump     adminv3.ConfigDump
		outLocation threadsafe.Buffer
	)

	err := c.ConfigDumpCmd(ctx, queryParams).WithStdout(&outLocation).Run().Cause()
	if err != nil {
		return nil, err
	}

	// Ever since upgrading the go-control-plane to v0.10.1 the standard unmarshal fails with the following error:
	// unknown field \"hidden_envoy_deprecated_build_version\" in envoy.config.core.v3.Node"
	// To get around this, we rely on an unmarshaler with AllowUnknownFields set to true
	if err = protoutils.UnmarshalAllowUnknown(&outLocation, &cfgDump); err != nil {
		return nil, err
	}

	return &cfgDump, nil
}

// GetStaticClusters returns the map of static clusters available on a ConfigDump, indexed by their name
func (c *Client) GetStaticClusters(ctx context.Context) (map[string]*clusterv3.Cluster, error) {
	configDump, err := c.GetConfigDump(ctx, map[string]string{
		"resource": "static_clusters",
	})
	if err != nil {
		return nil, err
	}

	return admincli.GetStaticClustersByName(configDump)
}

// ModifyRuntimeConfiguration passes the queryParameters to the runtime_modify endpoint
func (c *Client) ModifyRuntimeConfiguration(ctx context.Context, queryParameters map[string]string) error {
	return c.RunCommand(ctx,
		curl.WithPath(ModifyRuntimePath),
		curl.WithQueryParameters(queryParameters),
		curl.WithMethod(http.MethodPost))
}

// ShutdownServer calls the shutdown server endpoint
func (c *Client) ShutdownServer(ctx context.Context) error {
	return c.RunCommand(ctx,
		curl.WithPath(ShutdownServerPath),
		curl.WithMethod(http.MethodPost))
}

// FailHealthCheck calls the endpoint to have the server start failing health checks
func (c *Client) FailHealthCheck(ctx context.Context) error {
	return c.RunCommand(ctx,
		curl.WithPath(fmt.Sprintf("%s/fail", HealthCheckPath)),
		curl.WithMethod(http.MethodPost))
}

// PassHealthCheck calls the endpoint to have the server start passing health checks
func (c *Client) PassHealthCheck(ctx context.Context) error {
	return c.RunCommand(ctx,
		curl.WithPath(fmt.Sprintf("%s/ok", HealthCheckPath)),
		curl.WithMethod(http.MethodPost))
}

// SetLogLevel calls the endpoint to change the log level for the server
func (c *Client) SetLogLevel(ctx context.Context, logLevel string) error {
	return c.RunCommand(ctx,
		curl.WithPath(LoggingPath),
		curl.WithQueryParameters(map[string]string{
			"level": logLevel,
		}),
		curl.WithMethod(http.MethodPost))
}

// GetServerInfo calls the endpoint to return the info for the server
func (c *Client) GetServerInfo(ctx context.Context) (*adminv3.ServerInfo, error) {
	var (
		serverInfo  adminv3.ServerInfo
		outLocation threadsafe.Buffer
	)

	err := c.RequestPathCmd(ctx, ServerInfoPath).
		WithStdout(&outLocation).
		Run().
		Cause()
	if err != nil {
		return nil, err
	}

	if err = protoutils.UnmarshalAllowUnknown(&outLocation, &serverInfo); err != nil {
		return nil, err
	}

	return &serverInfo, nil
}

// GetSingleListenerFromDynamicListeners queries for a single, active dynamic listener in the envoy config dump
// and returns it as an envoy v3.Listener. This helper will only work if the provided name_regex matches a single dynamic_listener
// but will always use the first set of configs returned regardless
func (c *Client) GetSingleListenerFromDynamicListeners(
	ctx context.Context,
	listenerNameRegex string,
) (*listenerv3.Listener, error) {
	queryParams := map[string]string{
		"resource":   "dynamic_listeners",
		"name_regex": listenerNameRegex,
	}
	cfgDump, err := c.GetConfigDump(ctx, queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not get envoy config_dump from adminClient: %w", err)
	}

	listenerDump := adminv3.ListenersConfigDump_DynamicListener{}
	err = cfgDump.GetConfigs()[0].UnmarshalTo(&listenerDump)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal envoy config_dump: %w", err)
	}

	listener := listenerv3.Listener{}
	err = listenerDump.GetActiveState().GetListener().UnmarshalTo(&listener)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal listener from listener dump: %w", err)
	}
	return &listener, nil
}
