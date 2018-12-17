package del

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/spf13/cobra"
)

func Cmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "Delete a Gloo resource",
	}
	pflags := cmd.PersistentFlags()
	flagutils.AddMetadataFlags(pflags, &opts.Metadata)
	cmd.AddCommand(deleteUpstream(opts))
	cmd.AddCommand(deleteVirtualService(opts))
	cmd.AddCommand(deleteProxy(opts))
	return cmd
}

func getName(args []string, opts *options.Options) string {
	if len(args) > 0 {
		return args[0]
	}
	return opts.Metadata.Name
}

func deleteUpstream(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upstream",
		Aliases: []string{"u", "us", "upstreams"},
		Short:   "delete an upstream",
		Long:    "usage: glooctl get upstream [NAME] [--namespace=namespace]",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := getName(args, opts)
			if err := helpers.MustUpstreamClient().Delete(opts.Metadata.Namespace, name,
				clients.DeleteOpts{Ctx: opts.Top.Ctx}); err != nil {
				return err
			}
			fmt.Printf("uptream %v deleted", name)
			return nil
		},
	}
	return cmd
}

func deleteVirtualService(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "virtualservice",
		Aliases: []string{"v", "vs", "virtualservices"},
		Short:   "delete a virtualservice",
		Long:    "usage: glooctl delete virtualservice [NAME] [--namespace=namespace]",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := getName(args, opts)
			if err := helpers.MustVirtualServiceClient().Delete(opts.Metadata.Namespace, name,
				clients.DeleteOpts{Ctx: opts.Top.Ctx}); err != nil {
				return err
			}
			fmt.Printf("virtualservice %v deleted", name)
			return nil
		},
	}
	return cmd
}

func deleteProxy(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p", "px", "proxies"},
		Short:   "delete a proxy",
		Long:    "usage: glooctl delete proxy [NAME] [--namespace=namespace]",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := getName(args, opts)
			if err := helpers.MustProxyClient().Delete(opts.Metadata.Namespace, name,
				clients.DeleteOpts{Ctx: opts.Top.Ctx}); err != nil {
				return err
			}
			fmt.Printf("proxy %v deleted", name)
			return nil
		},
	}
	return cmd
}
