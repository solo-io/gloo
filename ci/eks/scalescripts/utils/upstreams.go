package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/labels"
)

// Waits for the services specified by the selector to be discovered as upstreams.
// Returns when `count` upstreams that match the given selector have been found, or a timeout occurs.
// If `mustBeAccepted` is true, then only upstreams with accepted status are counted.
func WaitForUpstreams(
	ctx context.Context,
	testClients *TestClients,
	count int,
	selector map[string]string,
	mustBeAccepted bool,
	ns string,
) error {
	if count == 0 {
		upstreamList, err := testClients.UpstreamClient.List(ns, clients.ListOpts{Ctx: ctx, Selector: selector})
		if err != nil {
			return err
		}
		count = len(upstreamList)
	}
	// watch for upstream changes
	upstreamWatch, watchErrors, err := testClients.UpstreamClient.Watch(ns, clients.WatchOpts{Ctx: ctx})
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			// timeout
			return eris.New("timed out waiting for upstreams to be discovered")
		case err := <-watchErrors:
			// a watch error occurred
			return err
		case newUpstreamList := <-upstreamWatch:

			// received new upstreams; count how many match the selector and are accepted
			matchingUpstreams := 0
			for _, us := range newUpstreamList {
				if IsMatchingUpstream(us, selector, mustBeAccepted) {
					matchingUpstreams++

					// if all `count` upstreams have been created, we can return
					if matchingUpstreams == count {
						return nil
					}
				}
			}
		}
	}
}

// Returns true if this is an upstream with the given selector. If `mustBeAccepted` is true, the upstream must also
// have accepted status.
func IsMatchingUpstream(us *gloov1.Upstream, selector map[string]string, mustBeAccepted bool) bool {
	if us.Metadata != nil && us.Metadata.Labels != nil &&
		labels.SelectorFromSet(selector).Matches(labels.Set(us.Metadata.Labels)) {
		if !mustBeAccepted {
			return true
		}
		if us.NamespacedStatuses != nil {
			for _, status := range us.NamespacedStatuses.GetStatuses() {
				if status.State == core.Status_Accepted {
					return true
				}
			}
		}
	}
	return false
}

func CreateStaticUpstream(
	ctx context.Context,
	upstreamClient gloov1.UpstreamClient,
	labels map[string]string,
	count int,
) ([]*gloov1.Upstream, error) {
	hosts := make([]*static.Host, 1)
	hosts[0] = &static.Host{
		Addr:    "34.193.132.77",
		Port:    80,
		SniAddr: "httpbin.org",
	}

	var results []*gloov1.Upstream
	for i := 0; i < count; i++ {

		us := &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "test-" + strconv.FormatInt(time.Now().UnixNano(), 10),
				Namespace: "gloo-system",
				Labels:    labels,
			},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &static.UpstreamSpec{
					Hosts: hosts,
				},
			},
		}

		resultUs, err := upstreamClient.Write(us, clients.WriteOpts{Ctx: ctx})
		if err != nil {
			return nil, err
		}
		fmt.Printf("Created us %s/%s\n", resultUs.Metadata.Name, resultUs.Metadata.Namespace)
		results = append(results, resultUs)

	}
	return results, nil
}
