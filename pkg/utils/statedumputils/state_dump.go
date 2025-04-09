package state_dump_utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/solo-io/gloo/pkg/cliutil/install"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	glooAdminCli "github.com/solo-io/gloo/pkg/utils/glooadminutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/pkg/servers/admin"
	"github.com/solo-io/go-utils/threadsafe"
)

// StandardCIDumpOnFail creates a dump of the CI system state, kubernetes state, gloo controller state,
// and certain envoy data from the admin interface when a test fails.
// Look at `CISystemDumpOnFail`, `KubeDumpOnFail`, `ControllerDumpOnFail`, && `EnvoyDumpOnFail` for more details
func StandardCIDumpOnFail(outLog io.Writer, outDir string, namespaces []string) func() {
	return func() {
		fmt.Printf("Test failed. Dumping state from %s...\n", strings.Join(namespaces, ", "))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		kubectlCli := kubectl.NewCli()

		// only wipe at the start of the dump
		wipeOutDir(outDir)

		CISystemDumpOnFail(ctx, kubectlCli, outLog, outDir, namespaces)()
		KubeDumpOnFail(ctx, kubectlCli, outLog, outDir, namespaces)()
		ControllerDumpOnFail(ctx, kubectlCli, outLog, outDir, namespaces)()
		EnvoyDumpOnFail(ctx, kubectlCli, outLog, outDir, namespaces)()

		fmt.Printf("Test failed. Logs and cluster state are available in %s\n", outDir)
	}
}

// CISystemDumpOnFail creates a small dump of the local docker and process state.
// This is useful for debugging test failures in CI.
// The dump includes:
// - docker state
// - process state
func CISystemDumpOnFail(_ context.Context, _ *kubectl.Cli, _ io.Writer, outDir string,
	_ []string) func() {
	return func() {
		setupOutDir(outDir)

		recordDockerState(fileAtPath(filepath.Join(outDir, "docker-state.log")))
		recordProcessState(fileAtPath(filepath.Join(outDir, "process-state.log")))

		fmt.Printf("Finished writing Docker and process state information to the \"%s\" directory.\n", outDir)
	}
}

// KubeDumpOnFail creates a small dump of the kubernetes state.
// This is useful for debugging failures.
// The dump includes:
// - kubernetes cluster state
// - logs from all pods in the given namespaces
// - yaml representations of all solo.io CRs in the given namespaces
func KubeDumpOnFail(ctx context.Context, _ *kubectl.Cli, _ io.Writer, outDir string,
	namespaces []string) func() {
	return func() {
		setupOutDir(outDir)

		recordKubeState(fileAtPath(filepath.Join(outDir, "kube-state.log")))

		recordKubeDump(ctx, outDir, namespaces...)

		fmt.Printf("Finished writing Kubernetes state information to the \"%s\" directory.\n", outDir)
	}
}

func recordDockerState(f *os.File) {
	defer f.Close()

	dockerCmd := exec.Command("docker", "ps")

	dockerState := &bytes.Buffer{}

	dockerCmd.Stdout = dockerState
	dockerCmd.Stderr = dockerState
	err := dockerCmd.Run()
	if err != nil {
		f.WriteString("*** Unable to get docker state ***. Reason: " + err.Error() + " \n")
		f.WriteString(dockerState.String() + "\n")
		return
	}
	f.WriteString("*** Docker state ***\n")
	f.WriteString(dockerState.String() + "\n")
	f.WriteString("*** End Docker state ***\n")
}

func recordProcessState(f *os.File) {
	defer f.Close()

	psCmd := exec.Command("ps", "-auxf")

	psState := &bytes.Buffer{}

	psCmd.Stdout = psState
	psCmd.Stderr = psState
	err := psCmd.Run()
	if err != nil {
		f.WriteString("unable to get process state. Reason: " + err.Error() + " \n")
		f.WriteString(psState.String() + "\n")
		return
	}
	f.WriteString("*** Process state ***\n")
	f.WriteString(psState.String() + "\n")
	f.WriteString("*** End Process state ***\n")
}

func recordKubeState(f *os.File) {
	defer f.Close()
	kubeCli := &install.CmdKubectl{}

	f.WriteString("*** Kube state ***\n")

	kubeState, err := kubeCli.KubectlOut(nil, "get", "all", "-A", "-o", "wide")
	if err != nil {
		f.WriteString("*** Unable to get kube state ***\n")
	}
	f.WriteString(string(kubeState) + "\n")

	resourcesToGet := []string{
		// Kubernetes resources
		"secrets",
		// Kube GW API resources
		"gateways.gateway.networking.k8s.io",
		"gatewayclasses.gateway.networking.k8s.io",
		"httproutes.gateway.networking.k8s.io",
		"referencegrants.gateway.networking.k8s.io",
		// GG Kube GW resources
		"gatewayparameters.gateway.gloo.solo.io",
		"listeneroptions.gateway.solo.io",     // only implemented for kube gw as of now
		"httplisteneroptions.gateway.solo.io", // only implemented for kube gw as of now
		// GG Gloo resources
		"graphqlapis.graphql.gloo.solo.io",
		"proxies.gloo.solo.io",
		"settings.gloo.solo.io",
		"upstreamgroups.gloo.solo.io",
		"upstreams.gloo.solo.io",
		// GG Edge GW resources
		"gateways.gateway.solo.io",
		"httpgateways.gateway.solo.io",
		"tcpgateways.gateway.solo.io",
		"virtualservices.gateway.solo.io",
		// Shared GW resources
		"routeoptions.gateway.solo.io",
		"virtualhostoptions.gateway.solo.io",
		// Dataplane extensions resources
		"authconfigs.enterprise.gloo.solo.io",
		"ratelimitconfigs.ratelimit.solo.io",
	}

	kubeResources, err := kubeCli.KubectlOut(nil, "get", strings.Join(resourcesToGet, ","), "-A", "-owide")
	if err != nil {
		f.WriteString("*** Unable to get kube resources ***. Reason: " + err.Error() + " \n")
	}
	f.WriteString(string(kubeResources) + "\n")

	// Describe everything to identify the reason for issues such as Pods, LoadBalancers stuck in pending state
	// (insufficient resources, unable to acquire an IP), etc.
	// Ie: More context around the output of the previous command `kubectl get all -A`
	kubeDescribe, err := kubeCli.KubectlOut(nil, "describe", "all", "-A")
	if err != nil {
		f.WriteString("*** Unable to get kube describe ***. Reason: " + err.Error() + " \n")
	}
	f.WriteString(string(kubeDescribe) + "\n")

	kubeEndpointsState, err := kubeCli.KubectlOut(nil, "get", "endpoints", "-A")
	if err != nil {
		f.WriteString("*** Unable to get endpoint state ***. Reason: " + err.Error() + " \n")
	}
	f.WriteString(string(kubeEndpointsState) + "\n")

	f.WriteString("*** End Kube state ***\n")
}

func recordKubeDump(ctx context.Context, outDir string, namespaces ...string) {
	// for each namespace, create a namespace directory that contains...
	for _, ns := range namespaces {
		// ...a pod logs subdirectory
		if err := recordPods(filepath.Join(outDir, ns, "_pods"), ns); err != nil {
			fmt.Printf("error recording pod logs: %f, \n", err)
		}

		// ...and a subdirectory for each solo.io CRD with non-zero resources
		if err := recordCRs(ctx, filepath.Join(outDir, ns), ns); err != nil {
			fmt.Printf("error recording pod logs: %f, \n", err)
		}
	}
}

// recordPods records logs from each pod to _output/kube2e-artifacts/$namespace/pods/$pod.log
func recordPods(podDir, namespace string) error {
	pods, _, err := kubeList(namespace, "pod")
	if err != nil {
		return err
	}

	outErr := &multierror.Error{}

	for _, pod := range pods {
		if err := os.MkdirAll(podDir, os.ModePerm); err != nil {
			return err
		}

		logs, errOutput, err := kubeLogs(namespace, pod)
		// store any error running the log command to return later
		// the error represents the cause of the failure, and should be bubbled up
		// we will still try to get logs for other pods even if this one returns an error
		if err != nil {
			outErr = multierror.Append(outErr, err)
		}
		// write any log output to the standard file
		if logs != "" {
			f := fileAtPath(filepath.Join(podDir, pod+".log"))
			f.WriteString(logs)
			f.Close()
		}
		// write any error output to the error file
		// this will consist of the combined stdout and stderr of the command
		if errOutput != "" {
			f := fileAtPath(filepath.Join(podDir, pod+"-error.log"))
			f.WriteString(errOutput)
			f.Close()
		}
	}

	return outErr.ErrorOrNil()
}

// recordCRs records all unique CRs floating about to _output/kube2e-artifacts/$namespace/$crd/$cr.yaml
func recordCRs(ctx context.Context, namespaceDir string, namespace string) error {
	crds, _, err := kubeList(namespace, "crd")
	if err != nil {
		return err
	}

	// record all unique CRs floating about
	for _, crd := range crds {
		// consider all installed CRDs that are solo-managed
		if !strings.Contains(crd, "solo.io") {
			continue
		}

		// if there are any existing CRs corresponding to this CRD
		crs, _, err := kubeList(namespace, crd)
		if err != nil {
			return err
		}
		if len(crs) == 0 {
			continue
		}
		crdDir := filepath.Join(namespaceDir, crd)
		if err := os.MkdirAll(crdDir, os.ModePerm); err != nil {
			return err
		}

		// we record each one in its own .yaml representation
		for _, cr := range crs {
			args := []string{"-n", namespace, crd, cr, "-oyaml"}
			stdout, stderr, err := kubectl.NewCli().Get(ctx, args...)
			if stdout != "" {
				f := fileAtPath(filepath.Join(crdDir, cr+".yaml"))
				f.WriteString(stdout)
				f.Close()
			}
			if stderr != "" {
				errF := fileAtPath(filepath.Join(crdDir, cr+"-error.log"))
				errF.WriteString(stderr)
				errF.Close()
			}
			// We don't expect an error to occur when executing this request.
			// If it does, we intentionally print stdout and stderr to the respective files so that we can more easily debug the error.
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// kubeLogs runs $(kubectl -n $namespace logs $pod --all-containers) and returns the string result
func kubeLogs(namespace string, pod string) (string, string, error) {
	args := []string{"-n", namespace, "logs", pod, "--all-containers=true", "--prefix=true"}
	return kubeExecute(args)
}

func kubeExecute(args []string) (string, string, error) {
	var outLocation threadsafe.Buffer
	runError := kubectl.NewCli().Command(context.Background(), args...).WithStdout(&outLocation).Run()
	if runError != nil {
		return outLocation.String(), runError.OutputString(), runError.Cause()
	}

	return outLocation.String(), "", nil
}

// kubeList runs $(kubectl -n $namespace $target) and returns a slice of kubernetes object names
func kubeList(namespace string, target string) ([]string, string, error) {
	args := []string{"-n", namespace, "get", target}
	lines, errContent, err := kubeExecute(args)
	if err != nil {
		return nil, errContent, err
	}

	var toReturn []string
	for _, line := range strings.Split(strings.TrimSuffix(lines, "\n"), "\n") {
		if strings.HasPrefix(line, "NAME") || strings.HasPrefix(line, "No resources found") {
			continue // skip header line and cases where there are no resources
		}
		if split := strings.Split(line, " "); len(split) > 1 {
			toReturn = append(toReturn, split[0])
		}
	}
	return toReturn, "", nil
}

// ControllerDumpOnFail creates a small dump of the gloo controller state.
// This is useful for debugging failures.
// The dump includes:
// - controller logs
// - controller metrics
// - controller xds snapshot
// - controller krt snapshot
func ControllerDumpOnFail(ctx context.Context, kubectlCli *kubectl.Cli, _ io.Writer,
	outDir string, namespaces []string) func() {
	return func() {
		for _, ns := range namespaces {
			glooPodNames, err := kubectlCli.GetPodsInNsWithLabel(ctx, ns, "gloo=gloo")
			if err != nil {
				fmt.Printf("error fetching controller pod names: %f\n", err)
				continue
			}

			if len(glooPodNames) == 0 {
				fmt.Printf("no gloo=gloo pods found in namespace %s\n", ns)
				continue
			}

			fmt.Printf("found controller pods: %s\n", strings.Join(glooPodNames, ", "))

			namespaceOutDir := filepath.Join(outDir, ns)
			setupOutDir(namespaceOutDir)

			for _, podName := range glooPodNames {
				writeControllerLog(ctx, namespaceOutDir, ns, podName, kubectlCli)
				writeMetricsLog(ctx, namespaceOutDir, ns, podName, kubectlCli)

				// Open a port-forward to the controller pod's admin port
				portForwarder, err := kubectlCli.StartPortForward(ctx,
					portforward.WithPod(podName, ns),
					portforward.WithPorts(admin.AdminPort, admin.AdminPort),
				)
				if err != nil {
					fmt.Printf("error starting port forward to controller admin port: %s\n", err.Error())
				}

				defer func() {
					portForwarder.Close()
					portForwarder.WaitForStop()
				}()

				adminClient := glooAdminCli.NewClient().
					WithReceiver(io.Discard).
					WithCurlOptions(
						curl.WithRetries(3, 0, 10),
						curl.WithPort(admin.AdminPort),
					)

				stderr := &bytes.Buffer{}
				krtSnapshotFile := fileAtPath(filepath.Join(namespaceOutDir, fmt.Sprintf("%s.krt_snapshot.log", podName)))
				err = adminClient.KrtSnapshotCmd(ctx).WithStdout(krtSnapshotFile).WithStderr(stderr).Run().Cause()
				if err != nil {
					fmt.Printf("error running krt snapshot command: %s\n", err.Error())
					fmt.Println(stderr)
				}

				stderr = &bytes.Buffer{}
				xdsSnapshotFile := fileAtPath(filepath.Join(namespaceOutDir, fmt.Sprintf("%s.xds_snapshot.log", podName)))
				err = adminClient.XdsSnapshotCmd(ctx).WithStdout(xdsSnapshotFile).WithStderr(stderr).Run().Cause()
				if err != nil {
					fmt.Printf("error running xds snapshot command: %s\n", err.Error())
					fmt.Println(stderr)
				}

				fmt.Printf("Finished writing Gloo Gateway controller state information to the \"%s\" directory.\n", outDir)
			}
		}
	}
}

// EnvoyDumpOnFail creates a small dump of the envoy admin interface.
// This is useful for debugging failures.
// The dump includes:
// - config dump
// - stats
// - clusters
// - listeners
func EnvoyDumpOnFail(ctx context.Context, kubectlCli *kubectl.Cli, _ io.Writer, outDir string, namespaces []string) func() {
	return func() {
		for _, ns := range namespaces {
			proxies := []string{}

			kubeGatewayProxies, err := kubectlCli.GetPodsInNsWithLabel(ctx, ns, "gloo=kube-gateway")
			if err != nil {
				fmt.Printf("error fetching kube-gateway proxies: %f\n", err)
			} else {
				proxies = append(proxies, kubeGatewayProxies...)
			}

			gatewayProxies, err := kubectlCli.GetPodsInNsWithLabel(ctx, ns, "gloo=gateway-proxy")
			if err != nil {
				fmt.Printf("error fetching gateway-proxy proxies: %f\n", err)
			} else {
				proxies = append(proxies, gatewayProxies...)
			}

			if len(proxies) == 0 {
				fmt.Printf("no proxies found in namespace %s\n", ns)
				continue
			}

			fmt.Printf("found proxies: %s\n", strings.Join(proxies, ", "))

			envoyOutDir := filepath.Join(outDir, ns)
			setupOutDir(envoyOutDir)

			for _, proxy := range proxies {
				adminCli, shutdown, err := admincli.NewPortForwardedClient(ctx, kubectlCli,
					fmt.Sprintf("pod/%s", proxy), ns)
				if err != nil {
					fmt.Printf("error creating admin cli: %f\n", err)
					continue
				}

				defer shutdown()

				configDumpFile := fileAtPath(filepath.Join(envoyOutDir, fmt.Sprintf("%s.config.log", proxy)))
				err = adminCli.ConfigDumpCmd(ctx, nil).WithStdout(configDumpFile).Run().Cause()
				if err != nil {
					fmt.Printf("error running config dump command: %f\n", err)
				}

				statsFile := fileAtPath(filepath.Join(envoyOutDir, fmt.Sprintf("%s.stats.log", proxy)))
				err = adminCli.StatsCmd(ctx, nil).WithStdout(statsFile).Run().Cause()
				if err != nil {
					fmt.Printf("error running stats command: %f\n", err)
				}

				clustersFile := fileAtPath(filepath.Join(envoyOutDir, fmt.Sprintf("%s.clusters.log", proxy)))
				err = adminCli.ClustersCmd(ctx).WithStdout(clustersFile).Run().Cause()
				if err != nil {
					fmt.Printf("error running clusters command: %f\n", err)
				}

				listenersFile := fileAtPath(filepath.Join(envoyOutDir, fmt.Sprintf("%s.listeners.log", proxy)))
				err = adminCli.ListenersCmd(ctx).WithStdout(listenersFile).Run().Cause()
				if err != nil {
					fmt.Printf("error running listeners command: %f\n", err)
				}

				fmt.Printf("Finished writing Envoy state information to the \"%s\" directory.\n", outDir)
			}
		}
	}
}

func wipeOutDir(outDir string) {
	err := os.RemoveAll(outDir)
	if err != nil {
		fmt.Printf("error wiping out directory: %f\n", err)
	}
}

// setupOutDir forcibly deletes/creates the output directory
func setupOutDir(outdir string) {
	err := os.MkdirAll(outdir, os.ModePerm)
	if err != nil {
		fmt.Printf("error creating log directory: %f\n", err)
	}
}

// fileAtPath creates a file at the given path, and returns the file object
func fileAtPath(path string) *os.File {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("unable to openfile: %f\n", err)
	}
	return f
}

func writeControllerLog(ctx context.Context, outDir string, ns string, podName string, kubectlCli *kubectl.Cli) {
	// Get the Gloo Gateway controller logs
	controllerLogsFile := fileAtPath(filepath.Join(outDir, fmt.Sprintf("%s.controller.log", podName)))
	controllerLogsCmd := kubectlCli.WithReceiver(controllerLogsFile).Command(ctx,
		"-n", ns, "logs", podName, "-c", "gloo", "--tail=1000")
	err := controllerLogsCmd.Run().Cause()
	if err != nil {
		fmt.Printf("error running controller logs for %s in %s command: %v\n", podName, ns, err)
	}
}

func writeMetricsLog(ctx context.Context, outDir string, ns string, podName string, kubectlCli *kubectl.Cli) {
	// Using an ephemeral debug pod fetch the metrics from the Gloo Gateway controller
	metricsFile := fileAtPath(filepath.Join(outDir, fmt.Sprintf("%s.metrics.log", podName)))
	metricsCmd := kubectlCli.Command(ctx, "debug", "-n", ns,
		"-it", "--image=curlimages/curl:7.83.1", podName, "--",
		"curl", "http://localhost:9091/metrics")
	err := metricsCmd.WithStdout(metricsFile).WithStderr(metricsFile).Run().Cause()
	if err != nil {
		fmt.Printf("error running metrics command: %f\n", err)
	}
}
