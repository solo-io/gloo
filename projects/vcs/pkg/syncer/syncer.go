package syncer

import (
	"context"
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/file"
)

// ApplyVcsToDeployment takes a dual client and applies the file client's configuration to the kubernetes client
// The expectation is that the file client represents the "master" branch in a git repo though any file client will work
// TODO - pass namespace as arg
// ToConsider - pass resources to sync as args (if we want to break out repos by resource type or if gitOps supports projects with differing resource sets)
func ApplyVcsToDeployment(ctx context.Context, fileClient, kubeClient file.ClientSet) {
	// TODO - pass this as an arg
	namespace := "gloo-system"
	aggregatedErrs := make(chan error)
	var wg sync.WaitGroup

	// gateway resources
	wg.Add(2)
	go updateGateways(ctx, fileClient, kubeClient, namespace, aggregatedErrs, &wg)
	go updateVirtualServices(ctx, fileClient, kubeClient, namespace, aggregatedErrs, &wg)

	// gloo resources
	wg.Add(2)
	go updateSettings(ctx, fileClient, kubeClient, namespace, aggregatedErrs, &wg)
	go updateProxies(ctx, fileClient, kubeClient, namespace, aggregatedErrs, &wg)

	// sqoop resources
	wg.Add(2)
	go updateResolverMaps(ctx, fileClient, kubeClient, namespace, aggregatedErrs, &wg)
	go updateSchemas(ctx, fileClient, kubeClient, namespace, aggregatedErrs, &wg)

	wg.Wait()
	return
}

////////////////////////////////////////////////////////////
// gateway resources
////////////////////////////////////////////////////////////

func updateGateways(ctx context.Context, fc, kc file.ClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := fc.GatewayClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gatewayv1.NewGatewayReconciler(kc.GatewayClient)
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

func updateVirtualServices(ctx context.Context, fc, kc file.ClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := fc.VirtualServiceClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gatewayv1.NewVirtualServiceReconciler(kc.VirtualServiceClient)
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

func updateSettings(ctx context.Context, fc, kc file.ClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := fc.SettingsClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gloov1.NewSettingsReconciler(kc.SettingsClient)
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
func updateProxies(ctx context.Context, fc, kc file.ClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := fc.ProxyClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := gloov1.NewProxyReconciler(kc.ProxyClient)
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

func updateResolverMaps(ctx context.Context, fc, kc file.ClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := fc.ResolverMapClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := sqoopv1.NewResolverMapReconciler(kc.ResolverMapClient)
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
func updateSchemas(ctx context.Context, fc, kc file.ClientSet, namespace string, aggErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sourceResources, err := fc.SchemaClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		aggErr <- err
	}
	rec := sqoopv1.NewSchemaReconciler(kc.SchemaClient)
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
