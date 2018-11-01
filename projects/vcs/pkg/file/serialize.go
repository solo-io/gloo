package file

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
)

// GenerateFilesystem
func GenerateFilesystem(ctx context.Context, namespace string, dc DualClientSet) {
	writeVirtualServices(ctx, namespace, dc.Kube.VirtualServiceClient, dc.File.VirtualServiceClient)
}

func writeVirtualServices(ctx context.Context, namespace string, vsk gatewayv1.VirtualServiceClient, vsf gatewayv1.VirtualServiceClient) {
	listK, err := vsk.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if len(listK) == 0 {
		fmt.Printf("please make a virtual service\n")
		return
	}
	virtualServiceK, err := vsk.Read(listK[0].Metadata.Namespace, listK[0].Metadata.Name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	listF, err := vsf.List(namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if len(listF) > 0 {
		virtualServiceF, err := vsf.Read(listF[0].Metadata.Namespace, listF[0].Metadata.Name, clients.ReadOpts{Ctx: ctx})
		if err != nil {
			fmt.Printf("file read err: %v\n", err)
		}
		fmt.Printf("file rv: %v\n", virtualServiceF.Metadata.ResourceVersion)

		virtualServiceK.VirtualHost.Domains = append([]string{fmt.Sprintf("%v.co", len(listK[0].VirtualHost.Domains))}, listK[0].VirtualHost.Domains...)

		fmt.Println("writing to kubernetes")
		_, err = vsk.Write(virtualServiceK, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			fmt.Printf("kube write err: %v\n", err)
		}

		virtualServiceF, err = vsf.Read(listF[0].Metadata.Namespace, listF[0].Metadata.Name, clients.ReadOpts{Ctx: ctx})
		if err != nil {
			fmt.Printf("file read err: %v\n", err)
		}
		fmt.Printf("file rv: %v\n", virtualServiceF.Metadata.ResourceVersion)

		rvk := virtualServiceK.Metadata.ResourceVersion
		rvf := virtualServiceF.Metadata.ResourceVersion
		fmt.Printf("file rv: %v\n", rvf)
		fmt.Printf("%v, %v, (rvs)\n", rvk, rvf)
		virtualServiceK.Metadata.ResourceVersion = rvf
	}
	// copy the first VS to the file system (just for some sample data)
	fmt.Println("writing to file")
	_, err = vsf.Write(virtualServiceK, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		fmt.Printf("file write err: %v\n", err)
		// panic for faster dev iterations
		panic("ouch")
	}
}
