package gateway

import (
	"context"
	"archive/zip"
	"fmt"
	"strings"
	"strconv"
	"os"
	"time"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/solo-io/go-utils/cliutils"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/spf13/cobra"
	"github.com/solo-io/solo-kit/pkg/errors"
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

func getEnvoyCfgDump(opts *options.Options) error {
	adminCli, deferFunc, err := admincli.NewPortForwardedClient(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace())
	if err != nil {
		return err
	}

	defer deferFunc()

	return adminCli.ConfigDumpCmd(opts.Top.Ctx, nil).WithStdout(os.Stdout).Run().Cause()
}

func getEnvoyStatsDump(opts *options.Options) error {
	adminCli, deferFunc, err := admincli.NewPortForwardedClient(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace())
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

	proxyNamespace := opts.Metadata.GetNamespace()
	if proxyNamespace == "" {
		proxyNamespace = defaults.GlooSystem
	}

	adminCli, deferFunc, err := admincli.NewPortForwardedClient(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace())
	if err != nil {
		return proxyOutArchiveFile.Name(), err
	}

	defer deferFunc()

	writeErr := adminCli.WriteEnvoyDumpToZip(opts.Top.Ctx, proxyOutArchive)

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
// GetEnvoyAdminData returns the response from the envoy admin interface based on the `path` specified within the defined timeout.
// Note that a `/` will be prepended to path if it does not exist.
func GetEnvoyAdminData(ctx context.Context, proxyName, namespace, path string, timeout time.Duration) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	adminPort := strconv.Itoa(int(defaults.EnvoyAdminPort))
	portFwd := exec.Command("kubectl", "port-forward", "-n", namespace,
		"deployment/"+proxyName, adminPort)
	portFwd.Stdout = os.Stderr
	portFwd.Stderr = os.Stderr
	if err := portFwd.Start(); err != nil {
		return "", errors.Wrapf(err, "failed to start port-forward")
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
