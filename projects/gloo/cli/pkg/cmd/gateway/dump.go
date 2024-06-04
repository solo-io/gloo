package gateway

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/solo-io/go-utils/cliutils"

	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/spf13/cobra"
)

func dumpCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
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

// GetEnvoyAdminData returns the response from the envoy admin interface based on the `path` specified within the defined timeout.
// Note that a `/` will be prepended to path if it does not exist.
func GetEnvoyAdminData(ctx context.Context, proxyName, namespace, path string, timeout time.Duration) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	adminPort := int(defaults.EnvoyAdminPort)

	pf := portforward.NewPortForwarder(portforward.WithDeployment(proxyName, namespace), portforward.WithPorts(adminPort, adminPort))
	err := pf.Start(ctx,
		retry.LastErrorOnly(true),
		retry.Delay(100*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(5))
	if err != nil {
		return "", errors.Wrapf(err, "failed to start port-forward")
	}
	defer func() {
		pf.Close()
		pf.WaitForStop()
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
			res, err := http.Get(pf.Address() + path)
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

func GetEnvoyCfgDump(opts *options.Options) (string, error) {
	return GetEnvoyAdminData(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace(), "/config_dump", 5*time.Second)
}

func statsCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
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

func getEnvoyStatsDump(opts *options.Options) (string, error) {
	return GetEnvoyAdminData(opts.Top.Ctx, opts.Proxy.Name, opts.Metadata.GetNamespace(), "/stats", 30*time.Second)
}
