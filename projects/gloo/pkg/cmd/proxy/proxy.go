package apply

import (
	"fmt"

	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/fileutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/cmd"
	"github.com/spf13/cobra"
)

var applyCmd = &cmd.Command{
	Command: &cobra.Command{
		Use:   "proxy",
		Short: "create or update proxy resources for gloo",
		RunE: func(_ *cobra.Command, _ []string) error {
			proxy, err := ProxyFromFile(cmd.RootOpts.ApplyOpts.File)
			if err != nil {
				return err
			}
			client, err := v1.NewProxyClient()
			changed, err := ApplyProxy(cmd.RootOpts.Ctx, client, proxy)

		},
	},
}

/*
if err != nil {
	return err
}
if err != nil {
	return err
}
if err != nil {
	return err
}


*/

var deleteCmd = &cmd.Command{
	Command: &cobra.Command{
		Use: "proxy",
		//		Short: "A brief description of your command",
		//		Long: `A longer description that spans multiple lines and likely contains examples
		//and usage of using your command. For example:
		//
		//Cobra is a CLI library for Go that empowers applications.
		//This application is a tool to generate the needed files
		//to quickly apply a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete.proxy called")
		},
	},
}

func init() {
	cmd.ApplyCmd.AddCommand(applyCmd.Command)
	cmd.DeleteCmd.AddCommand(deleteCmd.Command)
}

func ProxyFromFile(filename string) (*v1.Proxy, error) {
	var proxy v1.Proxy
	return &proxy, fileutils.ReadFileInto(filename, &proxy)
}

func ApplyProxy(ctx context.Context, client v1.ProxyClient, proxy *v1.Proxy) (bool, error) {
	existing, err := client.Read(proxy.GetMetadata().GetNamespace(), proxy.GetMetadata().GetName(), clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil && errors.IsNotExist(err) {
		return false, err
	}
	if proxy.Equal(existing) {
		return false, nil
	}
	if existing != nil {
		proxy.Metadata.ResourceVersion = existing.Metadata.ResourceVersion
	}
	return true, writeProxy(ctx, client, proxy)
}

func writeProxy(ctx context.Context, client v1.ProxyClient, proxy *v1.Proxy) error {
	_, err := client.Write(proxy, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	return err
}

func DeleteProxy(ctx context.Context, client v1.ProxyClient, namespace, name string) error {
	return client.Delete(namespace, name, clients.DeleteOpts{
		Ctx: ctx,
	})
}
