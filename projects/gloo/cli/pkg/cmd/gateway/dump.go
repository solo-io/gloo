package gateway

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	goerr "errors"
	"archive/zip"

	"github.com/solo-io/go-utils/cliutils"

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
			cfgDump, err := GetEnvoyCfgDump(opts)
			if err != nil {
				return err
			}
			fmt.Printf("%v", cfgDump)
			return nil
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
			cfgDump, err := getEnvoyStatsDump(opts)
			if err != nil {
				return err
			}
			fmt.Printf("%v", cfgDump)
			return nil
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
			err := getEnvoyFullDumpToDisk(opts)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
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
	if err := portFwd.Start(); err != nil {
		return "", errors.Wrapf(err, "failed to start port-forward")
	}

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
			return "", errors.Errorf("failed to start port-forward")
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
		case _ = <-errs:
			// TODO because we are just forking `kubectl` here, we don't know when the port-forward
			// is actually ready, so we basically can't stop ourselves from trying before its ready,
			// leading to spurious console errors if we log here.
			// This should be fixed with a real kube client. Until then the timeout error is sufficient.
		case res := <-result:
			return res, nil
		case <-timer:
			return "", errors.Errorf("timed out trying to connect to Envoy admin port")
		}
	}
}

func GetEnvoyCfgDump(opts *options.Options) (string, error) {
	return GetEnvoyAdminData(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace(), "/config_dump", 5*time.Second)
}

func getEnvoyStatsDump(opts *options.Options) (string, error) {
	return GetEnvoyAdminData(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace(), "/stats", 30*time.Second)
}

func getEnvoyFullDumpToDisk(opts *options.Options) error {
	proxyOutArchiveFile, err := createArchiveFile()
	if err != nil {
		return err
	}
	proxyOutArchive := zip.NewWriter(proxyOutArchiveFile)
	defer proxyOutArchiveFile.Close()
	defer proxyOutArchive.Close()

	proxyName := opts.Proxy.Name
	proxyNamespace := opts.Metadata.GetNamespace()
	if proxyNamespace == "" {
		proxyNamespace = defaults.GlooSystem
	}
    var errs []error

	errs = append(errs, recordEnvoyAdminData(fileInArchive(proxyOutArchive, "config.log"), opts.Top.Ctx, "/config_dump", proxyName, proxyNamespace))
	errs = append(errs, recordEnvoyAdminData(fileInArchive(proxyOutArchive, "stats.log"), opts.Top.Ctx, "/stats", proxyName, proxyNamespace))
	errs = append(errs, recordEnvoyAdminData(fileInArchive(proxyOutArchive, "clusters.log"), opts.Top.Ctx, "/clusters", proxyName, proxyNamespace))
	errs = append(errs, recordEnvoyAdminData(fileInArchive(proxyOutArchive, "listeners.log"), opts.Top.Ctx, "/listeners", proxyName, proxyNamespace))

	combinedErr := goerr.Join(errs...)
	if combinedErr == nil {
		fmt.Println("proxy snapshot written to " + proxyOutArchiveFile.Name())
	}
	return combinedErr
}

func recordEnvoyAdminData(w io.Writer, ctx context.Context, path, proxyName, namespace string) error {
	cfg, err := GetEnvoyAdminData(ctx, proxyName, namespace, path, 30*time.Second)
	if err != nil {
		io.WriteString(w, "*** Unable to get envoy " + path + " dump ***. Reason: " + err.Error() + " \n")
		return err
	}
	fmt.Printf("Snapshotting envoy state for %s from proxy instance %s.%s\n", path, proxyName, namespace)
	io.WriteString(w, "*** Envoy " + path + " dump ***\n")
	io.WriteString(w, cfg + "\n")
	io.WriteString(w, "*** End Envoy " + path + " dump ***\n")
	return nil
}

// createArchive forcibly deletes/creates the output directory
func createArchiveFile() (*os.File, error) {
	f, err := os.Create(fmt.Sprintf("glooctl-proxy-snapshot-%s.zip", time.Now().Format("2006-01-02-T15.04.05")))
	if err != nil {
		fmt.Printf("error creating proxy snapshot archive: %f\n", err)
	}
	return f, err
}

// fileInArchive creates a file at the given path, and returns the file object
func fileInArchive(w *zip.Writer, path string) io.Writer {
	f, err := w.Create(path)
	if err != nil {
		fmt.Printf("unable to create file: %f\n", err)
	}
	return f
}
