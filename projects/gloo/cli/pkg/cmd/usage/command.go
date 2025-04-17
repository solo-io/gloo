package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
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
	OutputFormat          string
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
	// Go fetch all the data needed for usage
	inputs, err := gatherUsageInformation(opts)
	if err != nil {
		return err
	}
	usageStats := &UsageStats{}

	// calculate all the stats for each input
	proxyData, err := generateProxyData(inputs.ProxyStats)
	if err != nil {
		return err
	}
	usageStats.ProxyData = proxyData

	// go through the edge snapshot and count feature usage
	usage, err := generateGlooFeatureUsage(inputs.GlooEdgeConfigs)
	if err != nil {
		return err
	}
	usageStats.GlooFeatureUsage = usage

	if inputs.ProxyStats != nil {
		proxyStats, err := gatherProxyPodInformation(inputs.ProxyStats)
		if err != nil {
			return err
		}
		usageStats.GlooProxyStats = proxyStats
	}
	// Calculate node resources
	nodeResources, err := calculateNodeResources(inputs.K8sClusterInfo.Nodes)
	if err != nil {
		return err
	}
	usageStats.KubernetesStats = &KubernetesStats{
		Pods:          len(inputs.K8sClusterInfo.Pods),
		NodeResources: nodeResources,
		Services:      len(inputs.K8sClusterInfo.Services),
	}

	usage.Print(opts.OutputFormat)

	return nil

}

func generateProxyData(proxyInfo map[string]*ProxyInfo) ([]*ProxyData, error) {

	var proxyData []*ProxyData
	for nameNamespace, info := range proxyInfo {
		// for each proxy organize its information
		envoyMetrics, err := getEnvoyMetrics(info)
		if err != nil {
			return nil, err
		}
		nnsSplit := strings.Split(nameNamespace, "/")

		proxyData = append(proxyData, &ProxyData{
			Name:         nnsSplit[1],
			Namespace:    nnsSplit[0],
			EnvoyMetrics: envoyMetrics,
		})
	}
	return proxyData, nil
}

type ProxyData struct {
	Name         string
	Namespace    string
	EnvoyMetrics *EnvoyMetrics
}

// gatherUsageInformation reads data from multiple sources and returns it
func gatherUsageInformation(opts *Options) (*Inputs, error) {

	tempDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir) // Clean up the directory when done

	inputs := &Inputs{}

	// Get cluster info
	clusterInfo, err := getK8sClusterInfo(opts)
	if err != nil {
		return nil, err
	}
	inputs.K8sClusterInfo = clusterInfo

	// scan for gloo gateways
	if len(opts.ScanProxies) > 0 {
		clusters, err := findGlooProxyPods(opts, tempDir)
		if err != nil {
			return nil, err
		}
		inputs.ProxyStats = clusters
	}

	// Create a temporary directory
	var filePath string

	if opts.GlooSnapshotFile == "" {
		filePath, err = LoadSnapshotFromGloo(opts, tempDir)
		if err != nil {
			return nil, err
		}
	} else {
		filePath = opts.GlooSnapshotFile
	}

	inputSnapshot, err := snapshot.FromGlooSnapshot(filePath)
	if err != nil {
		return nil, err
	}
	inputs.GlooEdgeConfigs = inputSnapshot

	return inputs, nil
}

func gatherProxyPodInformation(proxies map[string]*ProxyInfo) (map[string]*GlooProxyStats, error) {

	proxyStats := make(map[string]*GlooProxyStats)

	var envoyMetrics *EnvoyMetrics
	var upstreamStats []*UpstreamStat
	var err error
	for proxyName, proxy := range proxies {
		// Print upstream request stats first
		envoyMetrics, err = getEnvoyMetrics(proxy)
		if err != nil {
			return nil, err
		}
		if proxy.Clusters != nil {
			upstreamStats = generateUpstreamStats(proxy)
		}
		proxyStats[proxyName] = &GlooProxyStats{
			GlooProxyMetrics: envoyMetrics,
			GlooProxyStats:   upstreamStats,
		}
	}
	return proxyStats, nil
}

func generateUpstreamStats(proxy *ProxyInfo) []*UpstreamStat {
	var stats []*UpstreamStat
	for _, s := range proxy.Clusters.ClusterStatuses {
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
				upstreamStat := &UpstreamStat{}
				if hs.Locality.Region != "" {
					upstreamStat.Region = hs.Locality.Region
				}
				if hs.Locality.Zone != "" {
					upstreamStat.Zone = hs.Locality.Zone
				}
				if hs.Locality.SubZone != "" {
					upstreamStat.SubZone = hs.Locality.SubZone
				}

				rqSuccess, exists := statsMap["rq_success"]
				if exists {
					c, err := strconv.Atoi(rqSuccess)
					if err != nil {
						fmt.Printf("Error converting rq_success(%s) to int: %v\n", rqSuccess, err)
					}
					upstreamStat.RqSuccess = c
				}
				rqError, exists := statsMap["rq_error"]
				if exists {
					c, err := strconv.Atoi(rqError)
					if err != nil {
						fmt.Printf("Error converting rq_error(%s) to int: %v\n", rqError, err)
					}
					upstreamStat.RqError = c
				}
				cxActive, exists := statsMap["cx_active"]
				if exists {
					c, err := strconv.Atoi(cxActive)
					if err != nil {
						fmt.Printf("Error converting cx_active(%s) to int: %v\n", cxActive, err)
					}
					upstreamStat.CxActive = c
				}
				cxConnectFail, exists := statsMap["cx_connect_fail"]
				if exists {
					ccf, err := strconv.Atoi(cxConnectFail)
					if err != nil {
						fmt.Printf("Error converting cx_connect_fail(%s) to int: %v\n", cxConnectFail, err)
					}
					upstreamStat.CxConnectFail = ccf
				}
				stats = append(stats, upstreamStat)
			}
		}
	}
	return stats
}

func getEnvoyMetrics(info *ProxyInfo) (*EnvoyMetrics, error) {
	if info.Stats == nil {
		return nil, nil
	}
	envoyMetrics := &EnvoyMetrics{}
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
			envoyMetrics.UptimeSeconds = statValue
		}
		if strings.HasSuffix(stat.Name, "upstream_rq_2xx") {
			statValue, ok := stat.Value.(float64)
			if !ok {
				continue
			}
			if statValue > 0 {
				envoyMetrics.Total2xxRequests += statValue
			}
		}
		if strings.HasSuffix(stat.Name, "upstream_rq_3xx") {
			statValue, ok := stat.Value.(float64)
			if !ok {
				continue
			}
			if statValue > 0 {
				envoyMetrics.Total2xxRequests += statValue
			}
		}
		if strings.HasSuffix(stat.Name, "upstream_rq_4xx") {
			statValue, ok := stat.Value.(float64)
			if !ok {
				continue
			}
			if statValue > 0 {
				envoyMetrics.Total2xxRequests += statValue
			}
		}
		if strings.HasSuffix(stat.Name, "upstream_rq_5xx") {
			statValue, ok := stat.Value.(float64)
			if !ok {
				continue
			}
			if statValue > 0 {
				envoyMetrics.Total2xxRequests += statValue
			}
		}
	}

	return envoyMetrics, nil
}

type EnvoyMetrics struct {
	Total2xxRequests float64
	Total3xxRequests float64
	Total4xxRequests float64
	Total5xxRequests float64
	UptimeSeconds    float64
}

type UpstreamStat struct {
	Host          string
	Region        string
	Zone          string
	SubZone       string
	IPAddress     string
	Port          int
	RqSuccess     int
	RqError       int
	CxActive      int
	CxConnectFail int
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
	flags.StringVar(&o.OutputFormat, "output-format", "text", "Output format (text, json, yaml)")
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
