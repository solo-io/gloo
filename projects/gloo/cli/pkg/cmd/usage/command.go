package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/convert"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Options struct {
	*options.Options
	ControlPlaneName      string
	ControlPlaneNamespace string
	GlooSnapshotFile      string
	ScanProxies           []string
	ProxyNamespaces       []string
	IncludeEndpointStats  bool
}

func RootCmd(op *options.Options) *cobra.Command {
	opts := &Options{
		Options: op,
	}
	cmd := &cobra.Command{
		Use:     "usage",
		Short:   "Scan Gloo for feature usage",
		Long:    "",
		Example: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.validate(); err != nil {
				return err
			}

			return run(opts)
		},
	}
	opts.addToFlags(cmd.PersistentFlags())
	cmd.SilenceUsage = true
	return cmd
}

func (opts *Options) validate() error {

	return nil
}

func run(opts *Options) error {
	tempDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // Clean up the directory when done
	// Create a temporary directory
	var filePath string

	if opts.GlooSnapshotFile == "" {

		filePath, err = LoadSnapshotFromGloo(opts, tempDir)
		if err != nil {
			return err
		}
	} else {
		filePath = opts.GlooSnapshotFile
	}

	// scan for gloo gateways
	if len(opts.ScanProxies) > 0 {
		clusters, err := findGlooProxyPods(opts, tempDir)
		if err != nil {
			return err
		}

		if err := printPodInfo(clusters); err != nil {
			return err
		}

	}

	output := convert.NewGatewayAPIOutput()

	inputSnapshot, err := snapshot.FromGlooSnapshot(filePath)
	if err != nil {
		return err
	}
	output.EdgeCache(inputSnapshot)

	// go through the edge snapshot
	usage, err := generateUsage(output)
	if err != nil {
		return err
	}
	usage.Print()

	return nil

}

func printPodInfo(proxies map[string]*ProxyInfo) error {

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	var noEndpointsPods []string
	for namespacePodName, info := range proxies {

		fmt.Printf("%s\n", namespacePodName)
		// Print upstream request stats first
		if err := printUpstreamRequestStats(namespacePodName, info); err != nil {
			return err
		}
		if info.Clusters != nil {
			noEndpointsPods = printEndpointInfo(headerFmt, columnFmt, info, noEndpointsPods, namespacePodName)
		}
	}
	for _, pod := range noEndpointsPods {
		fmt.Printf("\nNo active endpoints found for %s", pod)
	}

	return nil
}

func printEndpointInfo(headerFmt func(format string, a ...interface{}) string, columnFmt func(format string, a ...interface{}) string, info *ProxyInfo, noEndpointsPods []string, namespacePodName string) []string {
	tbl := table.New("Upstream", "Endpoint", "Port", "Rq Success", "Rq Error", "Cx Active", "Cx Connect Fail")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	var rows = 0
	for _, s := range info.Clusters.ClusterStatuses {
		for _, hs := range s.HostStatuses {
			statsMap := map[string]string{}
			for _, stat := range hs.Stats {
				if stat.Value != "" {
					value, err := strconv.Atoi(stat.Value)
					if err != nil {
						// just print the string
						if stat.Value != "" {
							statsMap[stat.Name] = stat.Value
						}
					} else {
						if value != 0 {
							statsMap[stat.Name] = stat.Value
						}
					}
				}
			}
			if len(statsMap) > 0 {
				var locality string
				if hs.Locality.Region != "" {
					locality = hs.Locality.Region
				}
				if hs.Locality.Zone != "" {
					locality += "/" + hs.Locality.Zone
				}
				if hs.Locality.SubZone != "" {
					locality += "/" + hs.Locality.SubZone
				}

				tbl.AddRow(s.Name, hs.Address.SocketAddress.Address, hs.Address.SocketAddress.PortValue, statsMap["rq_success"], statsMap["rq_error"], statsMap["cx_active"], statsMap["cx_connect_fail"])
				rows++
			}
		}
	}
	if rows > 0 {
		tbl.Print()
		fmt.Println("\n")
	} else {
		noEndpointsPods = append(noEndpointsPods, namespacePodName)
	}
	return noEndpointsPods
}

func printUpstreamRequestStats(namespacePodName string, info *ProxyInfo) error {
	if info.Stats == nil {
		return nil
	}

	var uptimeSeconds float64
	var totalRequests float64
	for _, stat := range info.Stats.Stats {
		// skip xds_cluster stats
		if strings.Contains(stat.Name, "xds_cluster") {
			continue
		}

		if strings.Contains(stat.Name, "server.uptime") {
			statValue, ok := stat.Value.(float64)
			if !ok {
				continue
			}
			uptimeSeconds = statValue
		}

		if strings.HasSuffix(stat.Name, "upstream_rq_2xx") || strings.HasSuffix(stat.Name, "upstream_rq_3xx") || strings.HasSuffix(stat.Name, "upstream_rq_4xx") {

			statValue, ok := stat.Value.(float64)
			if !ok {
				continue
			}
			if statValue > 0 {
				// fmt.Printf("\t%s: %d\n", stat.Name, int(statValue))
				totalRequests += statValue
			}
		}
	}

	fmt.Printf("\tUptime: %d minute | Total 2xx,3xx,4xx: %d | Requests/sec: %d\n", int(uptimeSeconds/60), int(totalRequests), int(totalRequests/uptimeSeconds))

	return nil
}

type ProxyInfo struct {
	Clusters *Clusters
	Stats    *Stats
}

func findGlooProxyPods(opts *Options, tempDir string) (map[string]*ProxyInfo, error) {
	podStats := map[string]*ProxyInfo{}
	// for each selector

	namespaces := opts.ProxyNamespaces
	if len(namespaces) == 0 {
		if opts.ControlPlaneNamespace != "" {
			namespaces = append(namespaces, opts.ControlPlaneNamespace)
		} else {
			namespaces = append(namespaces, "gloo-system")
		}
	}

	restCfg, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		return nil, err
	}
	kube, err := kubernetes.NewForConfig(restCfg)

	ctx := context.Background()
	for _, proxySelector := range opts.ScanProxies {
		for _, namespace := range namespaces {

			if strings.HasPrefix(proxySelector, "deploy/") {

				deployString := strings.Split(proxySelector, "/")

				deployment, err := kube.AppsV1().Deployments(namespace).Get(ctx, deployString[1], metav1.GetOptions{})
				if err != nil {
					return nil, err
				}
				if deployment == nil {
					return nil, fmt.Errorf("deployment %s/%s not found", namespace, deployString[1])
				}
				labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)
				deploymentPods, err := kube.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
					LabelSelector: labelSelector,
				})
				if err != nil {
					return nil, err
				}
				for _, pod := range deploymentPods.Items {
					var clusters *Clusters
					if opts.IncludeEndpointStats {
						clusters, err = getClusterStats(tempDir, &pod, opts)
						if err != nil {
							return nil, err
						}
					}

					stats, err := getProxyStats(tempDir, &pod, opts)
					if err != nil {
						return nil, err
					}
					proxyInfo := &ProxyInfo{
						Stats: stats,
					}
					if clusters != nil {
						proxyInfo.Clusters = clusters
					}

					podStats[pod.Name] = proxyInfo
				}
			} else {
				pod, err := kube.CoreV1().Pods(namespace).Get(ctx, proxySelector, metav1.GetOptions{})
				// assume its just a pod name
				clusters, err := getClusterStats(tempDir, pod, opts)
				if err != nil {
					return nil, err
				}
				stats, err := getProxyStats(tempDir, pod, opts)
				if err != nil {
					return nil, err
				}
				podStats[proxySelector] = &ProxyInfo{
					Clusters: clusters,
					Stats:    stats,
				}
			}
			// need to parse the clusters file

		}
	}
	return podStats, nil

}

// getProxyStats is used to get the stats for a proxy pod /stats endpoint
func getProxyStats(tempDir string, pod *v1.Pod, opts *Options) (*Stats, error) {

	// TODO we need to find all gloo proxies to grab the stats

	cli, shutdownFunc, err := NewPortForwardedClient(context.Background(), kubectl.NewCli().WithKubeContext(opts.Top.KubeContext), pod.Name, pod.Namespace, int(defaults.EnvoyAdminPort))
	if err != nil {
		return nil, err
	}
	defer shutdownFunc()
	filePath := filepath.Join(tempDir, fmt.Sprintf("%s_%s_stats.json", pod.Namespace, pod.Name))
	clustersFile, err := fileAtPath(filePath)
	if err != nil {
		return nil, err
	}
	// We couldn't ge wget to work within the envoy proxy pod so we will curl it from the control plane
	cmd := cli.Command(context.Background(), curl.WithPath("stats"), curl.WithQueryParameters(map[string]string{"format": "json"}))
	if err := cmd.WithStdout(clustersFile).Run().Cause(); err != nil {
		return nil, err
	}

	// read stats file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	stats, err := parseStatsIntoEndpointInfo(string(data))
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func parseStatsIntoEndpointInfo(file string) (*Stats, error) {
	var stats Stats
	err := json.Unmarshal([]byte(file), &stats)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling stats: %v", err)
	}
	return &stats, nil
}

func parseClusterStatsIntoEndpointInfo(file string) (*Clusters, error) {
	var clusters Clusters
	err := json.Unmarshal([]byte(file), &clusters)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling clusters: %v", err)
	}
	return &clusters, nil
}
func getClusterStats(tempDir string, pod *v1.Pod, opts *Options) (*Clusters, error) {

	// TODO we need to find all gloo proxies to grab the stats

	cli, shutdownFunc, err := NewPortForwardedClient(context.Background(), kubectl.NewCli().WithKubeContext(opts.Top.KubeContext), pod.Name, pod.Namespace, int(defaults.EnvoyAdminPort))
	if err != nil {
		return nil, err
	}
	defer shutdownFunc()
	filePath := filepath.Join(tempDir, fmt.Sprintf("%s_%s_clusters.json", pod.Namespace, pod.Name))
	clustersFile, err := fileAtPath(filePath)
	if err != nil {
		return nil, err
	}
	// We couldn't ge wget to work within the envoy proxy pod so we will curl it from the control plane
	cmd := cli.Command(context.Background(), curl.WithPath("clusters"), curl.WithQueryParameters(map[string]string{"format": "json"}))
	if err := cmd.WithStdout(clustersFile).Run().Cause(); err != nil {
		return nil, err
	}

	// read stats file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	clusterStats, err := parseClusterStatsIntoEndpointInfo(string(data))
	if err != nil {
		return nil, err
	}

	return clusterStats, nil
}

func LoadSnapshotFromGloo(opts *Options, tempDir string) (string, error) {

	cli, shutdownFunc, err := NewPortForwardedClient(context.Background(), kubectl.NewCli().WithKubeContext(opts.Top.KubeContext), opts.ControlPlaneName, opts.ControlPlaneNamespace, 9095)
	if err != nil {
		return "", err
	}
	defer shutdownFunc()
	filePath := filepath.Join(tempDir, "gg-input.json")
	inputSnapshotFile, err := fileAtPath(filePath)
	if err != nil {
		return "", err
	}
	err = cli.RequestPathCmd(context.Background(), "/snapshots/input").WithStdout(inputSnapshotFile).Run().Cause()
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.ControlPlaneName, "gloo-control-plane", "deploy/gloo", "Name of the Gloo control plane pod")
	flags.StringVarP(&o.ControlPlaneNamespace, "gloo-control-plane-namespace", "n", "gloo-system", "Namespace of the Gloo control plane pod")
	flags.StringVar(&o.GlooSnapshotFile, "input-snapshot", "", "Gloo input snapshot file location")
	flags.StringSliceVar(&o.ScanProxies, "scan-proxies", []string{}, "Scan for Gloo proxies and grab their routing information")
	flags.StringSliceVar(&o.ProxyNamespaces, "proxy-namespaces", []string{}, "Namespaces that contain gloo proxies (default gloo-system or gloo-control-plane-namespace)")
	flags.BoolVar(&o.IncludeEndpointStats, "include-endpoint-stats", false, "Include endpoint stats in the output")

}

func NewPortForwardedClient(ctx context.Context, kubectlCli *kubectl.Cli, podSelector, namespace string, port int) (*admincli.Client, func(), error) {
	selector := portforward.WithResourceSelector(podSelector, namespace)

	// 1. Open a port-forward to the Kubernetes Deployment, so that we can query the Envoy Admin API directly
	portForwarder, err := kubectlCli.StartPortForward(ctx,
		selector,
		portforward.WithRemotePort(port))
	if err != nil {
		return nil, nil, err
	}

	// 2. Close the port-forward when we're done accessing data
	deferFunc := func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}

	// 3. Create a CLI that connects to the Envoy Admin API
	adminCli := admincli.NewClient().
		WithCurlOptions(
			curl.WithHostPort(portForwarder.Address()),
		)

	return adminCli, deferFunc, err
}

func fileAtPath(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("unable to openfile: %s %v", path, err)
	}
	return f, nil
}
