package file

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
)

// NOTE: This file is wip/demo only
// The functionality will be replaced by calls to reconcilers

// GenerateFilesystem
func GenerateFilesystem(ctx context.Context, namespace string, fc, kc ClientSet) {
	transposeVirtualServices(ctx, namespace, fc.VirtualServiceClient, kc.VirtualServiceClient)
	writeSchemas(ctx, namespace, kc.SchemaClient, fc.SchemaClient)
	writeResolverMaps(ctx, namespace, kc.ResolverMapClient, fc.ResolverMapClient)
	writeGateways(ctx, namespace, kc.GatewayClient, fc.GatewayClient)
	writeProxies(ctx, namespace, kc.ProxyClient, fc.ProxyClient)
	writeSettings(ctx, namespace, kc.SettingsClient, fc.SettingsClient)
}

func UpdateKube(ctx context.Context, namespace string, fc, kc ClientSet) {
	transposeVirtualServices(ctx, namespace, kc.VirtualServiceClient, fc.VirtualServiceClient)
}

func writeSchemas(ctx context.Context, namespace string, vsk sqoopv1.SchemaClient, vsf sqoopv1.SchemaClient) {
	// TODO get resource type from reflection
	fmt.Printf("Writing schemas to file\n")
	listKube, err := vsk.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if len(listKube) == 0 {
		fmt.Printf("please make a schema\n")
		return
	}
	for _, kubeResource := range listKube {
		fmt.Printf("writing: %v\n", kubeResource.Metadata.Name)
		_, err = vsf.Write(kubeResource, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("file write err: %v\n", err)
		}
	}
}

func writeResolverMaps(ctx context.Context, namespace string, vsk sqoopv1.ResolverMapClient, vsf sqoopv1.ResolverMapClient) {
	// TODO get resource type from reflection
	fmt.Printf("Writing resolvermaps to file\n")
	listKube, err := vsk.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if len(listKube) == 0 {
		fmt.Printf("please make a resolvermap\n")
		return
	}
	for _, kubeResource := range listKube {
		fmt.Printf("writing: %v\n", kubeResource.Metadata.Name)
		_, err = vsf.Write(kubeResource, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("file write err: %v\n", err)
		}
	}
}

func writeGateways(ctx context.Context, namespace string, vsk gatewayv1.GatewayClient, vsf gatewayv1.GatewayClient) {
	// TODO get resource type from reflection
	fmt.Printf("Writing gateways to file\n")
	listKube, err := vsk.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if len(listKube) == 0 {
		fmt.Printf("please make a gateway\n")
		return
	}
	for _, kubeResource := range listKube {
		fmt.Printf("writing: %v\n", kubeResource.Metadata.Name)
		_, err = vsf.Write(kubeResource, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("file write err: %v\n", err)
		}
	}
}

func writeProxies(ctx context.Context, namespace string, vsk gloov1.ProxyClient, vsf gloov1.ProxyClient) {
	// TODO get resource type from reflection
	fmt.Printf("Writing proxies to file\n")
	listKube, err := vsk.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if len(listKube) == 0 {
		fmt.Printf("please make a proxy\n")
		return
	}
	for _, kubeResource := range listKube {
		fmt.Printf("writing: %v\n", kubeResource.Metadata.Name)
		_, err = vsf.Write(kubeResource, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("file write err: %v\n", err)
		}
	}
}

func writeSettings(ctx context.Context, namespace string, vsk gloov1.SettingsClient, vsf gloov1.SettingsClient) {
	// TODO get resource type from reflection
	fmt.Printf("Writing settings to file\n")
	listKube, err := vsk.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if len(listKube) == 0 {
		fmt.Printf("please make a setting\n")
		return
	}
	for _, kubeResource := range listKube {
		fmt.Printf("writing: %v\n", kubeResource.Metadata.Name)
		_, err = vsf.Write(kubeResource, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("file write err: %v\n", err)
		}
	}
}

func writeVirtualServicesToKube(ctx context.Context, namespace string, toClient gatewayv1.VirtualServiceClient, fromClient gatewayv1.VirtualServiceClient) {
	listK, err := fromClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		fmt.Printf("source list err: %v\n", err)
		return
	}
	if len(listK) == 0 {
		fmt.Printf("please make a virtual service\n")
		return
	}
	for _, vs := range listK {
		virtualServiceK, err := fromClient.Read(vs.Metadata.Namespace, listK[0].Metadata.Name, clients.ReadOpts{Ctx: ctx})
		if err != nil {
			fmt.Printf("source read err: %v\n", err)
		}
		listF, err := toClient.List(namespace, clients.ListOpts{
			Ctx: ctx,
		})
		if len(listF) > 0 {
			virtualServiceF, err := toClient.Read(vs.Metadata.Namespace, vs.Metadata.Name, clients.ReadOpts{Ctx: ctx})
			if err != nil {
				fmt.Printf("target read err: %v\n", err)
			}

			// overwrite the resource version
			rvf := virtualServiceF.Metadata.ResourceVersion
			virtualServiceK.Metadata.ResourceVersion = rvf
		}
		// copy the first VS to the file system (just for some sample data)
		fmt.Println("writing to kube")
		_, err = toClient.Write(virtualServiceK, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("file write err: %v\n", err)
			// panic for faster dev iterations
			panic("ouch")
		}
	}
}

func transposeVirtualServices(ctx context.Context, namespace string, toClient gatewayv1.VirtualServiceClient, fromClient gatewayv1.VirtualServiceClient) {
	sourceResources, err := fromClient.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		fmt.Printf("source list err: %v\n", err)
		return
	}
	if len(sourceResources) == 0 {
		fmt.Printf("no virtual services exist in source client\n")
		return
	}
	for _, sourceResource := range sourceResources {
		if err != nil {
			fmt.Printf("source read err: %v\n", err)
		}
		targetResource, err := toClient.Read(sourceResource.Metadata.Namespace, sourceResource.Metadata.Name, clients.ReadOpts{Ctx: ctx})
		if err != nil {
			fmt.Printf("target read err: %v\nres: %v\n", err, targetResource)
		}
		if targetResource != nil {
			// resource exists, overwrite the resource version
			targetRV := targetResource.Metadata.ResourceVersion
			sourceResource.Metadata.ResourceVersion = targetRV
		}

		fmt.Println("writing to target client")
		_, err = toClient.Write(sourceResource, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("file write err: %v\n", err)
			// panic for faster dev iterations
			panic("ouch")
		}
	}
}
