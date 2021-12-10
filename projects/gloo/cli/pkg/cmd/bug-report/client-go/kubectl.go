package client_go

import (
	"bytes"
	"context"
	"fmt"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/bug-report/utils"
	"golang.org/x/time/rate"

	//"golang.org/x/time/rate"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	// maxRequestsPerSecond is the max rate of requests to the API server.
	maxRequestsPerSecond = 10
	// maxLogFetchConcurrency is the max number of logs to fetch simultaneously.
	maxLogFetchConcurrency = 10

	// reportInterval controls how frequently to output progress reports on running tasks.
	reportInterval = 30 * time.Second
)

var (
	requestLimiter  = rate.NewLimiter(maxRequestsPerSecond, maxRequestsPerSecond)
	logFetchLimitCh = make(chan struct{}, maxLogFetchConcurrency)

	// runningTasks tracks the in-flight fetch operations for user feedback.
	runningTasks   = make(map[string]struct{})
	runningTasksMu sync.RWMutex

	// runningTasksTicker is the report interval for running tasks.
	runningTasksTicker = time.NewTicker(reportInterval)
)

// Options contains the Run options.
type Options struct {
	// Path to the kubeconfig file.
	Kubeconfig string
	// ComponentName of the kubeconfig context to use.
	Context string

	// namespace - k8s namespace for Run command
	Namespace string

	// DryRun performs all steps but only logs the Run command without running it.
	DryRun bool
	// Maximum amount of time to wait for resources to be ready after install when Wait=true.
	WaitTimeout time.Duration

	// output - output mode for Run i.e. --output.
	Output string

	// extraArgs - more args to be added to the Run command, which are appended to
	// the end of the Run command.
	ExtraArgs []string
}

// RunCmd runs the given command in kubectl, adding -n namespace if namespace is not empty.
func RunCmd(command, namespace, kubeConfig, kubeContext string, dryRun bool) (string, error) {
	return Run(strings.Split(command, " "),
		&Options{
			Namespace:  namespace,
			DryRun:     dryRun,
			Kubeconfig: kubeConfig,
			Context:    kubeContext,
		})
}

func ReportRunningTasks() {
	go func() {
		time.Sleep(reportInterval)
		for range runningTasksTicker.C {
			printRunningTasks()
		}
	}()
}

// Run runs the kubectl command by specifying subcommands in subcmds with opts.
func Run(subcmds []string, opts *Options) (string, error) {
	args := subcmds
	if opts.Kubeconfig != "" {
		args = append(args, "--kubeconfig", opts.Kubeconfig)
	}
	if opts.Context != "" {
		args = append(args, "--context", opts.Context)
	}
	if opts.Namespace != "" {
		args = append(args, "-n", opts.Namespace)
	}
	if opts.Output != "" {
		args = append(args, "-o", opts.Output)
	}
	args = append(args, opts.ExtraArgs...)

	cmd := exec.Command("kubectl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmdStr := strings.Join(args, " ")

	if opts.DryRun {
		utils.Log("dry run mode: would be running this cmd:\nkubectl %s\n", cmdStr)
		return "", nil
	}

	_ = requestLimiter.Wait(context.TODO())
	task := fmt.Sprintf("kubectl %s", cmdStr)
	addRunningTask(task)
	defer removeRunningTask(task)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("kubectl error: %s\n\nstderr:\n%s\n\nstdout:\n%s",
			err, utils.ConsolidateLog(stderr.String()), stdout.String())
	}

	return stdout.String(), nil
}

func printRunningTasks() {
	runningTasksMu.RLock()
	defer runningTasksMu.RUnlock()
	if len(runningTasks) == 0 {
		return
	}
	fmt.Printf("The following fetches are still running: \n")
	for t := range runningTasks {
		fmt.Print("  %s\n", t)
	}
	fmt.Print("\n")
}

func addRunningTask(task string) {
	runningTasksMu.Lock()
	defer runningTasksMu.Unlock()
	utils.Log("STARTING %s", task)
	runningTasks[task] = struct{}{}
}

func removeRunningTask(task string) {
	runningTasksMu.Lock()
	defer runningTasksMu.Unlock()
	utils.Log("COMPLETED %s", task)
	delete(runningTasks, task)
}

// EnvoyGet sends a GET request for the URL in the Envoy container in the given namespace/pod and returns the result.
func EnvoyGet(client Client, namespace, pod, url string, dryRun bool) (string, error) {
	if dryRun {
		return fmt.Sprintf("Dry run: would be running client.EnvoyDo(%s, %s, %s)", pod, namespace, url), nil
	}
	_ = requestLimiter.Wait(context.TODO())
	task := fmt.Sprintf("ProxyGet %s/%s:%s", namespace, pod, url)
	addRunningTask(task)
	defer removeRunningTask(task)
	out, err := client.EnvoyDo(context.TODO(), pod, namespace, "GET", url)
	return string(out), err
}
