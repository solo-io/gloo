package syncer

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/file"
)

func ApplyVcsToDeployment(ctx context.Context, dc file.DualClientSet) {
	// TODO - pass this as an arg
	namespace := "gloo-system"
	aggregatedErrs := make(chan error)

	go func() {
		err := updateVirtualServices(ctx, dc, namespace)
		if err != nil {
			aggregatedErrs <- err
		}
	}()

	return
}

func transitionVirtualServices(original, desired *v1.VirtualService) (bool, error) {
	return true, nil
	// TODO - implement check
}

func updateVirtualServices(ctx context.Context, dc file.DualClientSet, namespace string) error {
	sourceResources, err := dc.File.VirtualServiceClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		return err
	}
	virtualServiceReconciler := v1.NewVirtualServiceReconciler(dc.Kube.VirtualServiceClient)
	return virtualServiceReconciler.Reconcile(namespace, sourceResources, transitionVirtualServices, clients.ListOpts{
		Ctx: ctx,
	})

}
