package syncer

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/file"
)

func ApplyVcsToDeployment(ctx context.Context, dc file.DualClientSet) {
	// TODO - pass this as an arg
	namespace := "gloo-system"

	sourceResources, err := dc.File.VirtualServiceClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		fmt.Printf("source list err: %v\n", err)
		return
	}

	aggregatedErrs := make(chan error)
	virtualServiceReconciler := v1.NewVirtualServiceReconciler(dc.Kube.VirtualServiceClient)
	if err := virtualServiceReconciler.Reconcile(namespace, sourceResources, transitionVirtualServices, clients.ListOpts{
		Ctx: ctx,
	}); err != nil {
		aggregatedErrs <- err
	}

	fmt.Println("here")
	return
	// return aggregatedErrs, nil
}

func transitionVirtualServices(original, desired *v1.VirtualService) (bool, error) {
	return true, nil
	// TODO - implement check
}
