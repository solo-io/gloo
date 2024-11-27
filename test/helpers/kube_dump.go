package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/solo-io/go-utils/threadsafe"

	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"

	"github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/pkg/cliutil/install"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"

	gateway_defaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/servers/admin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StandardGlooDumpOnFail creates adump of the kubernetes state and certain envoy data from
// the admin interface when a test fails.
// Look at `KubeDumpOnFail` && `EnvoyDumpOnFail` for more details
func StandardGlooDumpOnFail(outLog io.Writer, outDir string, proxies ...metav1.ObjectMeta) func() {
	return func() {
		var namespaces []string
		for _, proxy := range proxies {
			if proxy.GetNamespace() != "" {
				namespaces = append(namespaces, proxy.Namespace)
			}
		}

		KubeDumpOnFail(outLog, outDir, namespaces...)()
		ControllerDumpOnFail(outLog, outDir, namespaces...)()
		EnvoyDumpOnFail(outLog, outDir, proxies...)()

		fmt.Printf("Test failed. Logs and cluster state are available in %s\n", outDir)
	}
}

// KubeDumpOnFail creates a small dump of the kubernetes state when a test fails.
// This is useful for debugging test failures.
// The dump includes:
// - docker state
// - process state
// - kubernetes state
// - logs from all pods in the given namespaces
// - yaml representations of all solo.io CRs in the given namespaces
func KubeDumpOnFail(outLog io.Writer, outDir string, namespaces ...string) func() {
	return func() {
		setupOutDir(outDir)

		recordDockerState(fileAtPath(filepath.Join(outDir, "docker-state.log")))
		recordProcessState(fileAtPath(filepath.Join(outDir, "process-state.log")))
		recordKubeState(fileAtPath(filepath.Join(outDir, "kube-state.log")))

		recordKubeDump(outDir, namespaces...)
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
		return
	}
	f.WriteString("*** Process state ***\n")
	f.WriteString(psState.String() + "\n")
	f.WriteString("*** End Process state ***\n")
}

func recordKubeState(f *os.File) {
	defer f.Close()
	kubeCli := &install.CmdKubectl{}

	kubeState, err := kubeCli.KubectlOut(nil, "get", "all", "-A", "-o", "wide")
	if err != nil {
		f.WriteString("*** Unable to get kube state ***\n")
		return
	}

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
		return
	}

	// Describe everything to identify the reason for issues such as Pods, LoadBalancers stuck in pending state
	// (insufficient resources, unable to acquire an IP), etc.
	// Ie: More context around the output of the previous command `kubectl get all -A`
	kubeDescribe, err := kubeCli.KubectlOut(nil, "describe", "all", "-A")
	if err != nil {
		f.WriteString("*** Unable to get kube describe ***. Reason: " + err.Error() + " \n")
		return
	}

	kubeEndpointsState, err := kubeCli.KubectlOut(nil, "get", "endpoints", "-A")
	if err != nil {
		f.WriteString("*** Unable to get endpoint state ***. Reason: " + err.Error() + " \n")
		return
	}

	f.WriteString("*** Kube state ***\n")
	f.WriteString(string(kubeState) + "\n")
	f.WriteString(string(kubeResources) + "\n")
	f.WriteString(string(kubeDescribe) + "\n")
	f.WriteString(string(kubeEndpointsState) + "\n")

	f.WriteString("*** End Kube state ***\n")
}

func recordKubeDump(outDir string, namespaces ...string) {
	// for each namespace, create a namespace directory that contains...
	for _, ns := range namespaces {
		// ...a pod logs subdirectoy
		if err := recordPods(filepath.Join(outDir, ns, "_pods"), ns); err != nil {
			fmt.Printf("error recording pod logs: %f, \n", err)
		}

		// ...and a subdirectory for each solo.io CRD with non-zero resources
		if err := recordCRs(filepath.Join(outDir, ns), ns); err != nil {
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
func recordCRs(namespaceDir string, namespace string) error {
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
			f := fileAtPath(filepath.Join(crdDir, cr+".yaml"))
			errF := fileAtPath(filepath.Join(crdDir, cr+"-error.log"))

			crDetails, errOutput, err := kubeGet(namespace, crd, cr)

			if crDetails != "" {
				f.WriteString(crDetails)
				f.Close()
			}
			if errOutput != "" {
				errF.WriteString(errOutput)
				errF.Close()
			}

			return err
		}
	}

	return nil
}

// kubeLogs runs $(kubectl -n $namespace logs $pod --all-containers) and returns the string result
func kubeLogs(namespace string, pod string) (string, string, error) {
	args := []string{"-n", namespace, "logs", pod, "--all-containers"}
	return kubeExecute(args)
}

// kubeGet runs $(kubectl -n $namespace get $kubeType $name -oyaml) and returns the string result
func kubeGet(namespace string, kubeType string, name string) (string, string, error) {
	args := []string{"-n", namespace, "get", kubeType, name, "-oyaml"}
	return kubeExecute(args)
}

func kubeExecute(args []string) (string, string, error) {
	cli := kubectl.NewCli().WithReceiver(ginkgo.GinkgoWriter)

	var outLocation threadsafe.Buffer
	runError := cli.Command(context.Background(), args...).WithStdout(&outLocation).Run()
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

// ControllerDumpOnFail creates a small dump of the controller state when a test fails.
// This is useful for debugging test failures.
func ControllerDumpOnFail(outLog io.Writer, outDir string, namespaces ...string) func() {
	return func() {
		for _, ns := range namespaces {
			namespaceOutDir := filepath.Join(outDir, ns)
			setupOutDir(namespaceOutDir)

			// Get the Gloo Gateway controller logs
			controllerLogsFilePath := filepath.Join(namespaceOutDir, "controller.log")
			controllerLogsFile, err := os.OpenFile(controllerLogsFilePath,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				fmt.Printf("error opening controller log file: %f\n", err)
			}

			controllerLogsCmd := kubectl.NewCli().WithReceiver(controllerLogsFile).Command(context.Background(),
				"-n", ns, "logs", "deployment/gloo", "-c", "gloo", "--tail=1000")
			err = controllerLogsCmd.Run().Cause()
			if err != nil {
				fmt.Printf("error running controller logs command: %f\n", err)
			}

			// podStdOut := bytes.NewBuffer(nil)
			// podStdErr := bytes.NewBuffer(nil)

			// Fetch the name of the Gloo Gateway controller pod
			getGlooPodNameCmd := i.Actions.Kubectl().Command(ctx, "get", "pod", "-n", i.Metadata.InstallNamespace,
				"--selector", "gloo=gloo", "--output", "jsonpath='{.items[0].metadata.name}'")
			cmdErr := getGlooPodNameCmd.WithStdout(podStdOut).WithStderr(podStdErr).Run()
			if cmdErr != nil {
				i.Assertions.Require.NoError(cmdErr)
			}

			// Clean up and check the output
			glooPodName := strings.Trim(podStdOut.String(), "'")
			if glooPodName == "" {
				i.Assertions.Require.NoError(fmt.Errorf("failed to get the Gloo Gateway controller pod name: %s",
					podStdErr.String()))
			}

			// Get the metrics from the Gloo Gateway controller pod and write them to a file
			metricsFilePath := filepath.Join(failureDir, "metrics.log")
			metricsFile, err := os.OpenFile(metricsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			i.Assertions.Require.NoError(err)

			// Using an ephemeral debug pod fetch the metrics from the Gloo Gateway controller
			metricsCmd := i.Actions.Kubectl().Command(ctx, "debug", "-n", i.Metadata.InstallNamespace,
				"-it", "--image=curlimages/curl:7.83.1", glooPodName, "--",
				"curl", "http://localhost:9091/metrics")
			cmdErr = metricsCmd.WithStdout(metricsFile).WithStderr(metricsFile).Run()
			if cmdErr != nil {
				i.Assertions.Require.NoError(cmdErr)
			}

			// Open a port-forward to the Gloo Gateway controller pod's admin port
			portForwarder, err := i.Actions.Kubectl().StartPortForward(ctx,
				portforward.WithDeployment("gloo", i.Metadata.InstallNamespace),
				portforward.WithPorts(int(admin.AdminPort), int(admin.AdminPort)),
			)
			i.Assertions.Require.NoError(err)

			defer func() {
				portForwarder.Close()
				portForwarder.WaitForStop()
			}()

			adminClient := admincli.NewClient().
				WithReceiver(io.Discard).
				WithCurlOptions(
					curl.WithRetries(3, 0, 10),
					curl.WithPort(int(admin.AdminPort)),
				)

			// Get krt snapshot from the Gloo Gateway controller pod and write it to a file
			krtSnapshotFilePath := filepath.Join(failureDir, "krt_snapshot.log")
			krtSnapshotFile, err := os.OpenFile(krtSnapshotFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			i.Assertions.Require.NoError(err)

			cmdErr = adminClient.KrtSnapshotCmd(ctx).
				WithStdout(krtSnapshotFile).
				WithStderr(krtSnapshotFile).
				Run()
			if cmdErr != nil {
				i.Assertions.Require.NoError(cmdErr)
			}

			// Get xds snapshot from the Gloo Gateway controller pod and write it to a file
			xdsSnapshotFilePath := filepath.Join(failureDir, "xds_snapshot.log")
			xdsSnapshotFile, err := os.OpenFile(xdsSnapshotFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			i.Assertions.Require.NoError(err)

			cmdErr = adminClient.XdsSnapshotCmd(ctx).
				WithStdout(xdsSnapshotFile).
				WithStderr(xdsSnapshotFile).
				Run()
			if cmdErr != nil {
				i.Assertions.Require.NoError(cmdErr)
			}

		}
	}
}

// EnvoyDumpOnFail creates a small dump of the envoy admin interface when a test fails.
// This is useful for debugging test failures.
// The dump includes:
// - config dump
// - stats
// - clusters
// - listeners
func EnvoyDumpOnFail(_ io.Writer, outDir string, proxies ...metav1.ObjectMeta) func() {
	return func() {
		for _, proxy := range proxies {
			envoyOutDir := filepath.Join(outDir, proxy.Namespace, proxy.Name)
			setupOutDir(envoyOutDir)

			proxyName := proxy.GetName()
			if proxyName == "" {
				proxyName = gateway_defaults.GatewayProxyName
			}
			proxyNamespace := proxy.GetNamespace()
			if proxyNamespace == "" {
				proxyNamespace = defaults.GlooSystem
			}

			adminCli, shutdown, err := admincli.NewPortForwardedClient(context.Background(),
				fmt.Sprintf("deployment/%s", proxyName), proxyNamespace)
			if err != nil {
				fmt.Printf("error creating admin cli: %f\n", err)
				return
			}

			defer shutdown()

			adminCli.ConfigDumpCmd(context.Background(), nil).
				WithStdout(fileAtPath(filepath.Join(envoyOutDir, "config.log"))).Run().Cause()
			adminCli.StatsCmd(context.Background()).
				WithStdout(fileAtPath(filepath.Join(envoyOutDir, "stats.log"))).Run().Cause()
			adminCli.ClustersCmd(context.Background()).
				WithStdout(fileAtPath(filepath.Join(envoyOutDir, "clusters.log"))).Run().Cause()
			adminCli.ListenersCmd(context.Background()).
				WithStdout(fileAtPath(filepath.Join(envoyOutDir, "listeners.log"))).Run().Cause()
		}
	}
}

// setupOutDir forcibly deletes/creates the output directory
func setupOutDir(outdir string) {
	err := os.RemoveAll(outdir)
	if err != nil {
		fmt.Printf("error removing log directory: %f\n", err)
	}
	err = os.MkdirAll(outdir, os.ModePerm)
	if err != nil {
		fmt.Printf("error creating log directory: %f\n", err)
	}

	fmt.Println("kube dump artifacts will be stored at: " + outdir)
}

// fileAtPath creates a file at the given path, and returns the file object
func fileAtPath(path string) *os.File {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("unable to openfile: %f\n", err)
	}
	return f
}
