package setup

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ghodss/yaml"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/file"
)

func Setup(ctx context.Context, port int) error {
	gatewayOptsKube, gatewayOptsFile, err := DefaultGatewayOpts()
	vsClientKube, err := gatewayv1.NewVirtualServiceClient(gatewayOptsKube.VirtualServices)
	if err != nil {
		return err
	}
	if err := vsClientKube.Register(); err != nil {
		return err
	}

	vsClientFile, err := gatewayv1.NewVirtualServiceClient(gatewayOptsFile.VirtualServices)
	if err != nil {
		return err
	}
	if err := vsClientFile.Register(); err != nil {
		return err
	}

	dc, err := file.NewDualClient("kube")
	if err != nil {
		return err
	}

	// go ListCo(ctx, vsClientKube, vsClientFile)
	// go ListCo(ctx, dc.Kube.VirtualServiceClient, dc.File.VirtualServiceClient)
	go ListCo(ctx, dc)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func ListCo(ctx context.Context, dc file.DualClientSet) {
	// make a few edits
	for i := 0; i < 2; i++ {
		time.Sleep(1000 * time.Millisecond)
		// List(ctx, vsClientK, vsClientF)
		// file.WriteVirtualServices(ctx, "gloo-system", vsClientK, vsClientF)
		file.GenerateFilesystem(ctx, "gloo-system", dc)
	}
}

func DefaultGatewayOpts() (gatewaysetup.Opts, gatewaysetup.Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return gatewaysetup.Opts{}, gatewaysetup.Opts{}, err
	}
	cache := kube.NewKubeCache()
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return gatewaysetup.Opts{
			WriteNamespace: defaults.GlooSystem,
			Gateways: &factory.KubeResourceClientFactory{
				Crd:         gatewayv1.GatewayCrd,
				Cfg:         cfg,
				SharedCache: cache,
			},
			VirtualServices: &factory.KubeResourceClientFactory{
				Crd:         gatewayv1.VirtualServiceCrd,
				Cfg:         cfg,
				SharedCache: cache,
			},
			Proxies: &factory.KubeResourceClientFactory{
				Crd:         gloov1.ProxyCrd,
				Cfg:         cfg,
				SharedCache: cache,
			},
			WatchNamespaces: []string{"default", defaults.GlooSystem},
			WatchOpts: clients.WatchOpts{
				Ctx:         ctx,
				RefreshRate: defaults.RefreshRate,
			},
			DevMode: true,
		},
		gatewaysetup.Opts{
			WriteNamespace: defaults.GlooSystem,
			Gateways: &factory.KubeResourceClientFactory{
				Crd:         gatewayv1.GatewayCrd,
				Cfg:         cfg,
				SharedCache: cache,
			},
			VirtualServices: &factory.FileResourceClientFactory{
				RootDir: "gloo/virtualServices",
			},
			Proxies: &factory.KubeResourceClientFactory{
				Crd:         gloov1.ProxyCrd,
				Cfg:         cfg,
				SharedCache: cache,
			},
			WatchNamespaces: []string{"default", defaults.GlooSystem},
			WatchOpts: clients.WatchOpts{
				Ctx:         ctx,
				RefreshRate: defaults.RefreshRate,
			},
			DevMode: true,
		}, nil
}

func List(ctx context.Context, vsk gatewayv1.VirtualServiceClient, vsf gatewayv1.VirtualServiceClient) {
	var convertedSelector map[string]string
	list, err := vsk.List("gloo-system", clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	// print it go's default format
	fmt.Printf("List is %v\n", list)
	// print it in yaml
	yml, err := yaml.Marshal(list)
	fmt.Println(string(yml))

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	Write(ctx, vsk, vsf)
}

func Write(ctx context.Context, vsk gatewayv1.VirtualServiceClient, vsf gatewayv1.VirtualServiceClient) {
	var convertedSelector map[string]string
	listK, err := vsk.List("gloo-system", clients.ListOpts{
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
	listF, err := vsf.List("gloo-system", clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
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
