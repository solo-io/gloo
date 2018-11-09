package syncer

import (
	"context"
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"

	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/file"
)

// ApplyVcsToDeployment takes a dual client and applies the file client's configuration to the kubernetes client
// The expectation is that the file client represents the "master" branch in a git repo though any file client will work
// TODO - pass namespace as arg
// ToConsider - pass resources to sync as args (if we want to break out repos by resource type or if gitOps supports projects with differing resource sets)
func ApplyVcsToDeployment(ctx context.Context, dc file.DualClientSet) {
	// TODO - pass this as an arg
	namespace := "gloo-system"
	aggregatedErrs := make(chan error)
	var wg sync.WaitGroup

	wg.Add(2)
	go updateVirtualServices(ctx, dc, namespace, aggregatedErrs, &wg)
	go updateResolverMaps(ctx, dc, namespace, aggregatedErrs, &wg)
	wg.Wait()
	return
}

func transitionVirtualServices(original, desired *gatewayv1.VirtualService) (bool, error) {
	return true, nil
	// TODO - implement check
}

func updateVirtualServices(ctx context.Context, dc file.DualClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := dc.File.VirtualServiceClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gatewayv1.NewVirtualServiceReconciler(dc.Kube.VirtualServiceClient)
	err = rec.Reconcile(namespace, sourceResources, transitionVirtualServices, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
}

func updateResolverMaps(ctx context.Context, dc file.DualClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := dc.File.ResolverMapClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := sqoopv1.NewResolverMapReconciler(dc.Kube.ResolverMapClient)
	err = rec.Reconcile(namespace, sourceResources, transitionResolverMaps, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
}
func transitionResolverMaps(original, desired *sqoopv1.ResolverMap) (bool, error) {
	return true, nil
	// TODO - implement check
}
