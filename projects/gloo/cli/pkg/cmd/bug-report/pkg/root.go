package pkg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/bug-report/client-go"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/bug-report/content"
	. "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/bug-report/utils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/bug-report/utils/archive"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/bug-report/utils/cluster"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	version2 "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/version"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/printers"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	bugReportDefaultTimeout       = 30 * time.Minute
	bugReportDefaultGlooNamespace = "gloo-system"
)

var (
	// Logs, along with stats and importance metrics. Key is path (namespace/deployment/pod/cluster) which can be
	// parsed with ParsePath.
	logs = make(map[string]string)
	//stats      = make(map[string]*processlog.Stats)
	importance = make(map[string]int)
	// Aggregated errors for all fetch operations.
	gErrors Errors
	lock    = sync.RWMutex{}
)

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	bugReportConfig := &BugReportConfig{}
	cmd := &cobra.Command{
		Use:   constants.BUG_REPORT_COMMAND.Use,
		Short: constants.BUG_REPORT_COMMAND.Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBugReportCommand(cmd, args, bugReportConfig)
		},
	}
	addFlags(cmd, bugReportConfig)

	return cmd
}

func runBugReportCommand(cmd *cobra.Command, args []string, config *BugReportConfig) error {
	//todo(sai) - print running tasks?
	client_go.ReportRunningTasks()
	clusterCtxStr := ""
	if config.Context == "" {
		var err error
		clusterCtxStr, err = client_go.GetClusterContext(config.KubeConfigPath)
		if err != nil {
			return err
		} else {
			clusterCtxStr = config.Context
		}
	}

	Log("\nTarget cluster context: %s", clusterCtxStr)
	Log("Using following config: \n\n%+v\n\n", *config)

	clientConfig, clientset, err := client_go.New(config.KubeConfigPath, config.Context)
	if err != nil {
		return fmt.Errorf("could not initialize k8s client: %s ", err)
	}

	Log("Writing output to dir %s", tempDir)
	client, err := client_go.NewClient(clientConfig)
	if err != nil {
		return err
	}
	resources, err := cluster.GetClusterResources(context.Background(), clientset)
	if err != nil {
		return err
	}
	versionBuffer := bytes.NewBufferString("")
	err = version2.PrintVersion(version2.NewKube(config.GlooNamespace), versionBuffer, &options.Options{Top: options.Top{Output: printers.TABLE, Ctx: context.Background()}})
	if err != nil {
		return eris.Wrap(err, "error getting server versions")
	}
	writeFile(filepath.Join(archive.OutputRootDir(tempDir), "versions"), versionBuffer.String())
	Log("Cluster resource tree:\n\n%s\n\n", resources)
	gatherInfo(client, config, resources, nil)
	if len(gErrors) != 0 {
		Log(gErrors.ToError().Error())
	}
	return nil
}

func gatherInfo(client client_go.Client, config *BugReportConfig, resources *cluster.Resources, paths []string) {
	// no timeout on mandatoryWg.
	var mandatoryWg sync.WaitGroup
	cmdTimer := time.NewTimer(config.CommandTimeout)
	//beginTime := time.Now()

	clusterDir := archive.ClusterInfoPath(tempDir)
	params := &content.Params{
		Client:      client,
		DryRun:      config.DryRun,
		KubeConfig:  config.KubeConfigPath,
		KubeContext: config.Context,
	}
	Log("\nFetching Gloo Edge information from cluster.\n\n")
	getFromCluster(content.GetK8sResources, params, clusterDir, &mandatoryWg)
	getFromCluster(content.GetCRs, params, clusterDir, &mandatoryWg)
	getFromCluster(content.GetEvents, params, clusterDir, &mandatoryWg)
	getFromCluster(content.GetClusterInfo, params, clusterDir, &mandatoryWg)
	getFromCluster(content.GetNodeInfo, params, clusterDir, &mandatoryWg)
	getFromCluster(content.GetSecrets, params.SetVerbose(config.FullSecrets), clusterDir, &mandatoryWg)
	getFromCluster(content.GetDescribePods, params.SetGlooNamespace(config.GlooNamespace), clusterDir, &mandatoryWg)
	for _, proxyPod := range GetGatewayProxyPodNames(config.GlooNamespace, resources) {
		getFromCluster(content.GetProxyInfo, params.SetPod(proxyPod).SetNamespace(config.GlooNamespace), archive.ProxyOutputPath(tempDir, params.SetNamespace(config.GlooNamespace).Namespace, proxyPod), &mandatoryWg)
	}
	mandatoryWg.Wait()

	go func() {
		//optionalWg.Wait()
		cmdTimer.Reset(0)
	}()

	<-cmdTimer.C

}

func GetGatewayProxyPodNames(namespace string, resources *cluster.Resources) []string {
	var proxyPodNames []string
	if ns := resources.Root[namespace]; ns != nil {
		if ns, ok := ns.(map[string]interface{}); ok {
			if gp := ns["gateway-proxy"]; gp != nil {
				if gpns, ok := gp.(map[string]interface{}); ok {
					for podName, _ := range gpns {
						proxyPodNames = append(proxyPodNames, podName)
					}
				}
			}
		}
	}
	return proxyPodNames
}

// getFromCluster runs a cluster info fetching function f against the cluster and writes the results to fileName.
// Runs if a goroutine, with errors reported through gErrors.
func getFromCluster(f func(params *content.Params) (map[string]string, error), params *content.Params, dir string, wg *sync.WaitGroup) {
	wg.Add(1)
	Log("Waiting on %s", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
	go func() {
		defer wg.Done()
		out, err := f(params)
		appendGlobalErr(err)
		if err == nil {
			writeFiles(dir, out)
		}
		Log("Done with %s", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
	}()
}

func appendGlobalErr(err error) {
	if err == nil {
		return
	}
	lock.Lock()
	gErrors = AppendErr(gErrors, err)
	lock.Unlock()
}

func writeFiles(dir string, files map[string]string) {
	for fname, text := range files {
		writeFile(filepath.Join(dir, fname), text)
	}
}

func writeFile(path, text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	mkdirOrExit(path)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		Log(err.Error())
	}
}

func mkdirOrExit(fpath string) {
	if err := os.MkdirAll(path.Dir(fpath), 0o755); err != nil {
		fmt.Printf("Could not create output directories: %s", err)
		os.Exit(-1)
	}
}
