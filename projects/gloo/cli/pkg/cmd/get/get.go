package get

import (
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/common"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/spf13/cobra"
)

func Cmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Display one or a list of Gloo resources",
	}
	pflags := cmd.PersistentFlags()
	flagutils.AddMetadataFlags(pflags, &opts.Metadata)
	flagutils.AddOutputFlag(pflags, &opts.Top.Output)
	cmd.AddCommand(getUpstreams(opts))
	cmd.AddCommand(getVirtualServices(opts))
	cmd.AddCommand(getProxies(opts))
	return cmd
}

func getName(args []string, opts *options.Options) string {
	if len(args) > 0 {
		return args[0]
	}
	return opts.Metadata.Name
}

func getUpstreams(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upstream",
		Aliases: []string{"u", "us", "upstreams"},
		Short:   "read an upstream or list upstreams in a namespace",
		Long:    "usage: glooctl get upstream [NAME] [--namespace=namespace] [-o FORMAT] [-o FORMAT]",
		RunE: func(cmd *cobra.Command, args []string) error {
			upstreams, err := common.GetUpstreams(getName(args, opts), opts)
			if err != nil {
				return err
			}
			helpers.PrintUpstreams(upstreams, opts.Top.Output)
			return nil
		},
	}
	return cmd
}

func getVirtualServices(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "virtualservice",
		Aliases: []string{"vs", "virtualservices"},
		Short:   "read a virtualservice or list virtualservices in a namespace",
		Long:    "usage: glooctl get virtualservice [NAME] [--namespace=namespace] [-o FORMAT]",
		RunE: func(cmd *cobra.Command, args []string) error {
			virtualServices, err := common.GetVirtualServices(getName(args, opts), opts)
			if err != nil {
				return err
			}
			helpers.PrintVirtualServices(virtualServices, opts.Top.Output)
			return nil
		},
	}
	cmd.AddCommand(getRoutes(opts))
	return cmd
}

func getProxies(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p", "proxies"},
		Short:   "read a proxy or list proxies in a namespace",
		Long:    "usage: glooctl get proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			proxyList, err := common.GetProxies(getName(args, opts), opts)
			if err != nil {
				return err
			}
			helpers.PrintProxies(proxyList, opts.Top.Output)
			return nil
		},
	}
	return cmd
}

func getRoutes(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "route",
		Aliases: []string{"r", "routes"},
		Short:   "get a list of routes for a given virtual service",
		Long:    "usage: glooctl get virtualservice route",
		RunE: func(cmd *cobra.Command, args []string) error {
			var vsName string
			if len(args) > 0 {
				vsName = args[0]
			}
			virtualServices, err := common.GetVirtualServices(vsName, opts)
			if err != nil {
				return err
			}
			if len(virtualServices.Names()) != 1 {
				return errors.Errorf("no virtualservice id provided")
			}
			vs, err := virtualServices.Find(opts.Metadata.Namespace, opts.Metadata.Name)
			if err != nil {
				return errors.Errorf("virtualservice id provided was incorrect")
			}
			helpers.PrintRoutes(vs.VirtualHost.Routes, opts.Top.Output)
			return nil
		},
	}
	return cmd
}
