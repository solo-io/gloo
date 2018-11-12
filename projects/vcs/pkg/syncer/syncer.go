package syncer

import (
	"context"
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	gatewayv1 "github.com/solo-io/solo-projects/projects/gateway/pkg/api/v1"

	gloov1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/file"
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

	// gateway resources
	wg.Add(2)
	go updateGateways(ctx, dc, namespace, aggregatedErrs, &wg)
	go updateVirtualServices(ctx, dc, namespace, aggregatedErrs, &wg)

	// gloo resources
	wg.Add(2)
	go updateSettings(ctx, dc, namespace, aggregatedErrs, &wg)
	go updateProxies(ctx, dc, namespace, aggregatedErrs, &wg)

	// sqoop resources
	wg.Add(2)
	go updateResolverMaps(ctx, dc, namespace, aggregatedErrs, &wg)
	go updateSchemas(ctx, dc, namespace, aggregatedErrs, &wg)

	wg.Wait()
	return
}

////////////////////////////////////////////////////////////
// gateway resources
////////////////////////////////////////////////////////////

func updateGateways(ctx context.Context, dc file.DualClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := dc.File.GatewayClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gatewayv1.NewGatewayReconciler(dc.Kube.GatewayClient)
	err = rec.Reconcile(namespace, sourceResources, transitionGateways, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
}
func transitionGateways(original, desired *gatewayv1.Gateway) (bool, error) {
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
func transitionVirtualServices(original, desired *gatewayv1.VirtualService) (bool, error) {
	return true, nil
	// TODO - implement check
}

////////////////////////////////////////////////////////////
// gloo resources
////////////////////////////////////////////////////////////

func updateSettings(ctx context.Context, dc file.DualClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := dc.File.SettingsClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gloov1.NewSettingsReconciler(dc.Kube.SettingsClient)
	err = rec.Reconcile(namespace, sourceResources, transitionSettings, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
}
func transitionSettings(original, desired *gloov1.Settings) (bool, error) {
	return true, nil
	// TODO - implement check
}
func updateProxies(ctx context.Context, dc file.DualClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := dc.File.ProxyClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gloov1.NewProxyReconciler(dc.Kube.ProxyClient)
	err = rec.Reconcile(namespace, sourceResources, transitionProxies, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
}
func transitionProxies(original, desired *gloov1.Proxy) (bool, error) {
	return true, nil
	// TODO - implement check
}

////////////////////////////////////////////////////////////
// sqoop resources
////////////////////////////////////////////////////////////

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
func updateSchemas(ctx context.Context, dc file.DualClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := dc.File.SchemaClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := sqoopv1.NewSchemaReconciler(dc.Kube.SchemaClient)
	err = rec.Reconcile(namespace, sourceResources, transitionSchemas, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
}
func transitionSchemas(original, desired *sqoopv1.Schema) (bool, error) {
	return true, nil
	// TODO - implement check
}
