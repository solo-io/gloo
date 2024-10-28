package gateway

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/solo-io/go-utils/cliutils"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/spf13/cobra"
)

func dumpCfgCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "dump Envoy config from one of the proxy instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getEnvoyCfgDump(opts)
		},
	}
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func dumpStatsCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "stats for one of the proxy instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getEnvoyStatsDump(opts)
		},
	}
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func writeSnapshotCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "snapshot complete proxy state for the given instance to an archive",
		RunE: func(cmd *cobra.Command, args []string) error {
			dumpFile, err := getEnvoyFullDumpToDisk(opts)
			if err != nil {
				// If we have an error writing zip (or fetching dump)
				// delete the file after it's flushed to clean up.
				_ = os.Remove(dumpFile)
				return err
			}
			return nil
		},
	}
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func buildEnvoyClient(ctx context.Context, proxySelector, namespace string) (*admincli.Client, func(), error) {
	var selector portforward.Option
	if sel := strings.Split(proxySelector, "/"); len(sel) == 2 {
		if strings.HasPrefix(sel[0], "deploy") {
			selector = portforward.WithDeployment(sel[1], namespace)
		} else if strings.HasPrefix(sel[0], "po") {
			selector = portforward.WithPod(sel[1], namespace)
		}
	} else {
		selector = portforward.WithPod(proxySelector, namespace)
	}

	// 1. Open a port-forward to the Kubernetes Deployment, so that we can query the Envoy Admin API directly
	portForwarder, err := kubectl.NewCli().StartPortForward(ctx,
		selector,
		portforward.WithRemotePort(int(defaults.EnvoyAdminPort)))
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

func writeEnvoyDumpToZip(ctx context.Context, proxySelector, namespace string, zip *zip.Writer) error {

	adminCli, deferFunc, err := buildEnvoyClient(ctx, proxySelector, namespace)
	if err != nil {
		return err
	}

	defer deferFunc()

	// zip writer has the benefit of not requiring tmpdirs or file ops (all in mem)
	// - but it can't support async writes, so do these sequentally
	// Also don't join errors, we want to fast-fail
	if err := adminCli.ConfigDumpCmd(ctx, nil).WithStdout(fileInArchive(zip, "config.log")).Run().Cause(); err != nil {
		return err
	}
	if err := adminCli.StatsCmd(ctx).WithStdout(fileInArchive(zip, "stats.log")).Run().Cause(); err != nil {
		return err
	}
	if err := adminCli.ClustersCmd(ctx).WithStdout(fileInArchive(zip, "clusters.log")).Run().Cause(); err != nil {
		return err
	}
	if err := adminCli.ListenersCmd(ctx).WithStdout(fileInArchive(zip, "listeners.log")).Run().Cause(); err != nil {
		return err
	}

	return nil
}

// GetEnvoyAdminData returns the response from the envoy admin interface identified by `proxySelector`.
// `proxySelector` can be any valid `kubectl` selection string,
// such as a podname (e.g `gloo-proxy-http-64c6746757-24cfx`) or a deployment (e.g. `deployment/gloo-proxy-http`).
//
// Admin endpoint data will be fetched from `path` using the defined timeout.
// Note that a `/` will be prepended to path if it does not exist.
func GetEnvoyAdminData(ctx context.Context, proxySelector, namespace, path string, timeout time.Duration) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	adminPort := strconv.Itoa(int(defaults.EnvoyAdminPort))
	// TODO	this should use a real Go kube client library someday
	portFwd := exec.Command("kubectl", "port-forward", "-n", namespace,
		proxySelector, adminPort)
	fwdOut, _ := portFwd.StdoutPipe()
	fwdErr, _ := portFwd.StderrPipe()
	if err := portFwd.Start(); err != nil {
		return "", errors.Wrapf(err, "failed to start port-forward")
	}

	// Because we are not using a real kube client but are spawning a long-running
	// subprocess that *attempts* to port-forward, we need to wait until the
	// port-forward actually completes (stdout scans) before trying to query the endpoint
	outScan := bufio.NewScanner(fwdOut)
	for {
		outScanned := outScan.Scan()
		if outScanned {
			if strings.Contains(outScan.Text(), "Forwarding from") {
				break
			} else {
				return "", errors.Errorf("failed to start port-forward")
			}
		} else {
			outErr := bufio.NewScanner(fwdErr)
			outErr.Scan()
			return "", errors.Errorf("failed to start port-forward: %s", outErr.Text())
		}
	}

	defer func() {
		if portFwd.Process != nil {
			portFwd.Process.Kill()
		}
	}()
	result := make(chan string)
	errs := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			res, err := http.Get("http://localhost:" + adminPort + path)
			if err != nil {
				errs <- err
				time.Sleep(time.Millisecond * 250)
				continue
			}
			if res.StatusCode != http.StatusOK {
				errs <- errors.Errorf("invalid status code: %v %v", res.StatusCode, res.Status)
				time.Sleep(time.Millisecond * 250)
				continue
			}
			b, err := io.ReadAll(res.Body)
			if err != nil {
				errs <- err
				time.Sleep(time.Millisecond * 250)
				continue
			}
			res.Body.Close()
			result <- string(b)
			return
		}
	}()

	timer := time.Tick(timeout)

	for {
		select {
		case <-ctx.Done():
			return "", errors.Errorf("cancelled")
		case err := <-errs:
			log.Printf("connecting to envoy failed with err %v", err.Error())
		case res := <-result:
			return res, nil
		case <-timer:
			return "", errors.Errorf("timed out trying to connect to Envoy admin port")
		}
	}
}

func getEnvoyCfgDump(opts *options.Options) error {
	adminCli, deferFunc, err := buildEnvoyClient(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace())
	if err != nil {
		return err
	}

	defer deferFunc()

	return adminCli.ConfigDumpCmd(opts.Top.Ctx, nil).WithStdout(os.Stdout).Run().Cause()
}

func getEnvoyStatsDump(opts *options.Options) error {
	adminCli, deferFunc, err := buildEnvoyClient(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace())
	if err != nil {
		return err
	}

	defer deferFunc()

	return adminCli.StatsCmd(opts.Top.Ctx).WithStdout(os.Stdout).Run().Cause()
}

func getEnvoyFullDumpToDisk(opts *options.Options) (string, error) {
	proxyOutArchiveFile, err := createArchiveFile()
	if err != nil {
		return proxyOutArchiveFile.Name(), err
	}
	proxyOutArchive := zip.NewWriter(proxyOutArchiveFile)
	defer proxyOutArchiveFile.Close()
	defer proxyOutArchive.Close()

	proxyName := opts.Proxy.Name
	proxyNamespace := opts.Metadata.GetNamespace()
	if proxyNamespace == "" {
		proxyNamespace = defaults.GlooSystem
	}

	writeErr := writeEnvoyDumpToZip(opts.Top.Ctx, proxyName, proxyNamespace, proxyOutArchive)

	if writeErr == nil {
		fmt.Println("proxy snapshot written to " + proxyOutArchiveFile.Name())
	} else {
		fmt.Printf("Error writing proxy snapshot: %s", writeErr)
	}

	return proxyOutArchiveFile.Name(), writeErr
}

// createArchive forcibly deletes/creates the output directory
func createArchiveFile() (*os.File, error) {
	f, err := os.Create(fmt.Sprintf("glooctl-proxy-snapshot-%s.zip", time.Now().Format("2006-01-02-T15.04.05")))
	if err != nil {
		fmt.Printf("error creating proxy snapshot archive: %f\n", err)
	}
	return f, err
}

// fileInArchive creates a file at the given path within the archive, and returns the file object for writing.
func fileInArchive(w *zip.Writer, path string) io.Writer {
	f, err := w.Create(path)
	if err != nil {
		fmt.Printf("unable to create file: %f\n", err)
	}
	return f
}
