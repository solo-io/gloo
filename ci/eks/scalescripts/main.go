package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/ci/eks/scalescripts/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/maps"
)

var constantLabels = map[string]string{
	"scale": "testing",
	"new":   "test",
}

func main() {

	cmd := &cobra.Command{
		Use: "scale-test",
	}
	opts := defaultOptions()
	opts.addToFlags(cmd.PersistentFlags())
	cmd.AddCommand(
		scaleUp(opts),
		scaleDown(opts),
	)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func defaultOptions() *options {
	return &options{}
}

type options struct {
	testRequest bool
	testStatus  bool
	filename    string
	runs        uint32
	namespace   string
}

func (r *options) addToFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&r.testRequest, "test-request", false, "set to true to test request")
	fs.BoolVar(&r.testStatus, "test-status", false, "set to true to test status")
	fs.StringVarP(&r.filename, "filename", "f", "data.csv", "filename to save data to")
	fs.StringVarP(&r.namespace, "namespace", "n", "gloo-system", "namespace to look for gloo assets")
	fs.Uint32VarP(&r.runs, "runs", "r", 10, "number of runs")
}

func scaleUp(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "scaleup",
		RunE: func(cmd *cobra.Command, args []string) error {
			testClient, err := utils.NewTestClients(cmd.Context(), opts.namespace)
			if err != nil {
				return err
			}

			f, err := os.Create(opts.filename)
			if err != nil {
				return err
			}

			defer f.Close()

			for i := 0; i < int(opts.runs); i++ {
				if err := testCreationTime(cmd.Context(), constantLabels, testClient, f, opts); err != nil {
					return err
				}
				f.WriteString("\n")
			}

			return nil
		},
	}
	return cmd
}

func scaleDown(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "scaledown",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("calling scaledown func")
			testClient, err := utils.NewTestClients(cmd.Context(), opts.namespace)
			if err != nil {
				return err
			}

			f, err := os.Create(opts.filename)
			if err != nil {
				return err
			}

			defer f.Close()

			for i := 0; i < int(opts.runs); i++ {
				if err := testDeletionTime(cmd.Context(), constantLabels, testClient, f, opts); err != nil {
					return err
				}
				f.WriteString("\n")
			}

			return nil
		},
	}
	return cmd
}

func testDeletionTime(
	ctx context.Context,
	selectors map[string]string,
	testClients *utils.TestClients,
	w io.Writer,
	opts *options,
) error {
	fmt.Printf("Test Deletion Time\n")
	virtualServices, err := testClients.VirtualServiceClient.List(opts.namespace, clients.ListOpts{Ctx: ctx, Selector: selectors})
	if err != nil {
		return err
	}

	if len(virtualServices) == 0 {
		return errors.New("Cannot test deletion: No virtual services found with list")
	}
	//var domain string
	vsToDelete := virtualServices[0]
	fmt.Printf("Virtual Service to delete: %s\n", vsToDelete.GetMetadata().GetName())
	domain := vsToDelete.GetVirtualHost().GetDomains()[0]
	vsLabels := vsToDelete.GetMetadata().GetLabels()
	if err := testClients.VirtualServiceClient.Delete(vsToDelete.Metadata.Namespace, vsToDelete.Metadata.Name, clients.DeleteOpts{Ctx: ctx}); err != nil {
		return err
	}

	before := time.Now()
	upstreams, err := testClients.UpstreamClient.List(opts.namespace, clients.ListOpts{Ctx: ctx, Selector: vsLabels})
	if err != nil {
		return err
	}
	for _, us := range upstreams {
		fmt.Printf("Upstream to delete: %s\n", us.GetMetadata().GetName())
		if err := testClients.UpstreamClient.Delete(us.Metadata.Namespace, us.Metadata.Name, clients.DeleteOpts{Ctx: ctx}); err != nil {
			return err
		}
	}

	host, err := utils.GetIngressHost(ctx, opts.namespace)
	if err != nil {
		return err
	}
	requestUrl := fmt.Sprintf("%s/status/200", host)
	// create request
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return err
	}
	req.Host = domain
	for {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("connected")
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		// Looking for not found because we wanna make sure it was deleted.
		if res.StatusCode == http.StatusNotFound {
			break
		}
	}

	timeToDelete := time.Since(before).Seconds()
	fmt.Printf("Time to Delete %f\n", timeToDelete)
	w.Write([]byte(fmt.Sprintf(",%f", timeToDelete)))
	return nil
}

func testCreationTime(
	ctx context.Context,
	selectors map[string]string,
	testClient *utils.TestClients,
	w io.Writer,
	opts *options,
) error {
	if selectors == nil {
		selectors = map[string]string{}
	}

	// Always have twice as many upstreams
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Printf("created label: created=%s \n", now)
	usSelectors := map[string]string{
		"created": now,
	}
	maps.Copy(usSelectors, selectors)
	maps.Copy(usSelectors, constantLabels)
	usList, err := utils.CreateStaticUpstream(ctx, testClient.UpstreamClient, usSelectors, 2)
	if err != nil {
		return err
	}
	if err := doRequest(ctx, testClient, w, usSelectors, opts, usList...); err != nil {
		return err
	}

	return nil
}

func doRequest(
	ctx context.Context,
	testClients *utils.TestClients,
	w io.Writer,
	selectors map[string]string,
	opts *options,
	uss ...*gloov1.Upstream,
) error {
	host, err := utils.GetIngressHost(ctx, opts.namespace)
	if err != nil {
		return err
	}
	requestUrl := fmt.Sprintf("%s/status/200", host)

	var measurements []string
	var refs []*core.ResourceRef

	for _, us := range uss {
		refs = append(refs, us.GetMetadata().Ref())
	}

	before := time.Now()
	// create a VS that routes to the upstream
	vs, err := utils.CreateVirtualServiceForUpstream(ctx, opts.namespace, testClients.VirtualServiceClient, selectors, refs...)
	if err != nil {
		return err
	}
	timeToCreate := time.Since(before).Seconds()
	fmt.Printf("Time to create: %f\n", timeToCreate)
	measurements = append(measurements, fmt.Sprintf("%f", timeToCreate))

	if opts.testRequest {
		domain := vs.GetVirtualHost().GetDomains()[0]
		fmt.Printf("Starting request check for domain: %s\n", domain)
		before = time.Now()
		// create request
		req, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			return err
		}
		req.Host = domain
		for {
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				break
			}

			// Just a little break
			time.Sleep(50 * time.Millisecond)

		}
		timeToRequest := time.Since(before).Seconds()
		fmt.Printf("Time to request: %f\n", timeToRequest)
		measurements = append(measurements, fmt.Sprintf("%f", timeToRequest))
	}

	if _, err := w.Write([]byte(strings.Join(measurements, ","))); err != nil {
		log.Print(err)
	}

	if opts.testStatus {
		fmt.Printf("Waiting for statuses of resources (%+v) to be updated\n", selectors)
		if err := utils.WaitForUpstreams(ctx, testClients, 0, selectors, true, opts.namespace); err != nil {
			return err
		}
		// if err := utils.WaitForVirtualServices(ctx, testClients, 0, selectors, true); err != nil {
		// 	return err
		// }
	}
	return nil
}
