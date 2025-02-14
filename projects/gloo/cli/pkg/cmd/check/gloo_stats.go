package check

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/printers"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	rateLimitDeployment = "rate-limit"
	glooStatsPath       = "/metrics"

	glooeTotalEntites   = "glooe_solo_io_xds_total_entities"
	glooeInSyncEntities = "glooe_solo_io_xds_insync"

	GlooeRateLimitConnectedState = "glooe_ratelimit_connected_state"
)

var (
	resourcesOutOfSyncMessage = func(resourceNames []string) string {
		return fmt.Sprintf("Gloo has detected that the data plane is out of sync. The following types of resources have not been accepted: %v. "+
			"Gloo will not be able to process any other configuration updates until these errors are resolved.", resourceNames)
	}

	// Initialize the custom deployment name that is overwritten later on
	customGlooDeploymentName = helpers.GlooDeploymentName
)

func ResourcesSyncedOverXds(printer printers.P, stats, deploymentName string) bool {
	var outOfSyncResources []string
	metrics := parseMetrics(stats, []string{glooeTotalEntites, glooeInSyncEntities}, deploymentName)
	for metric, val := range metrics {
		if strings.HasPrefix(metric, glooeTotalEntites) {
			inSyncMetric := strings.ReplaceAll(metric, glooeTotalEntites, glooeInSyncEntities)
			if inSyncVal, ok := metrics[inSyncMetric]; ok {
				if val > inSyncVal {
					outOfSyncResources = append(outOfSyncResources, strings.ReplaceAll(metric, glooeTotalEntites, ""))
				}
			}
		}
	}
	if len(outOfSyncResources) > 0 {
		fmt.Println(resourcesOutOfSyncMessage(outOfSyncResources))
		return false
	}

	if len(metrics) == 0 {
		printer.AppendStatus("xds metrics", "No xds metrics to check")
	} else {
		printer.AppendStatus("xds metrics", "OK")
	}
	return true
}

func RateLimitIsConnected(stats string) bool {
	connectedStateErrMessage := "The rate limit server is out of sync with the Gloo control plane and is not receiving valid gloo config.\n" +
		"You may want to try using the `glooctl debug logs --errors-only` command to find any relevant error logs."

	metrics := parseMetrics(stats, []string{GlooeRateLimitConnectedState}, "gloo")

	if val, ok := metrics[GlooeRateLimitConnectedState]; ok && val == 0 {
		fmt.Println(connectedStateErrMessage)
		return false
	}

	return true
}

func checkXdsMetrics(ctx context.Context, printer printers.P, opts *options.Options, deployments *appsv1.DeploymentList) error {
	printer.AppendCheck("Checking xds metrics... ")
	errMessage := "Problem while checking for gloo xds errors"
	if deployments == nil {
		fmt.Println("Skipping due to an error in checking deployments")
		return fmt.Errorf("xds metrics check was skipped due to an error in checking deployments")
	}
	// port-forward proxy deployment and get prometheus metrics
	freePort, err := cliutil.GetFreePort()
	if err != nil {
		fmt.Println(errMessage)
		return err
	}
	localPort := strconv.Itoa(freePort)
	adminPort := strconv.Itoa(int(defaults.GlooAdminPort))
	// stats is the string containing all stats from /stats/prometheus
	if opts.Top.ReadOnly {
		printer.AppendCheck("Warning: checking xds with port forwarding is disabled\n")
		return nil
	}
	stats, portForwarder, err := cliutil.PortForwardGet(ctx, opts.Metadata.GetNamespace(), "deploy/"+customGlooDeploymentName,
		localPort, adminPort, false, glooStatsPath)
	if err != nil {
		return err
	}
	defer func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}()

	if strings.TrimSpace(stats) == "" {
		err := fmt.Sprint(errMessage+": could not find any metrics at", glooStatsPath, "endpoint of the "+customGlooDeploymentName+" deployment")
		fmt.Println(err)
		return fmt.Errorf(err)
	}

	if !ResourcesSyncedOverXds(printer, stats, customGlooDeploymentName) {
		fmt.Println(errMessage)
		return fmt.Errorf(errMessage)
	}

	for _, deployment := range deployments.Items {
		if deployment.Name == rateLimitDeployment {
			printer.AppendCheck("Checking rate limit server... ")
			if !RateLimitIsConnected(stats) {
				return fmt.Errorf("rate limit server is not connected")
			}
			printer.AppendStatus("rate limit server", "OK")
		}
	}

	return nil
}
