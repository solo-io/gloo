package route

import (
	"fmt"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/cliutil"
	editOptions "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/edit/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/surveyutils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/constants"
	extauthpb "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type routeEditInput struct {
	Index   uint32
	Disable bool
}

func editRouteFlags(set *pflag.FlagSet, route *routeEditInput) {
	set.Uint32VarP(&route.Index, "index", "x", 0, "edit the route with this index in the virtual service "+
		"route list")
	set.BoolVarP(&route.Disable, "disable", "d", false, "set to true to disable authentication on this route")
}

func ExtAuthConfig(opts *editOptions.EditOptions, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {

	input := &routeEditInput{}

	cmd := &cobra.Command{
		// Use command constants to aid with replacement.
		Use:     constants.CONFIG_EXTAUTH_COMMAND.Use,
		Aliases: constants.CONFIG_EXTAUTH_COMMAND.Aliases,
		Short:   "Configure disable external auth on a route (Enterprise)",
		Long:    "Allows disabling external auth on specific routes. External auth is a gloo enterprise feature.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Top.Interactive {

				vsclient := helpers.MustVirtualServiceClient()
				vsvc, err := vsclient.Read(opts.Metadata.Namespace, opts.Metadata.Name, clients.ReadOpts{})
				if err != nil {
					return err
				}

				if idx, err := surveyutils.SelectRouteFromVirtualServiceInteractive(vsvc, "Choose the route you wish to change: "); err != nil {
					return err
				} else {
					input.Index = uint32(idx)
					opts.ResourceVersion = vsvc.Metadata.ResourceVersion
				}
				cliutil.ChooseBool("Disable auth on this route?", &input.Disable)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return editRoute(opts, input, args)
		},
	}

	pflags := cmd.PersistentFlags()
	flagutils.AddMetadataFlags(pflags, &opts.Metadata)

	editRouteFlags(cmd.Flags(), input)
	cliutils.ApplyOptions(cmd, optionsFunc)

	return cmd
}

func editRoute(opts *editOptions.EditOptions, input *routeEditInput, args []string) error {
	vsClient := helpers.MustVirtualServiceClient()
	vs, err := vsClient.Read(opts.Metadata.Namespace, opts.Metadata.Name, clients.ReadOpts{})
	if err != nil {
		return errors.Wrapf(err, "Error reading vhost")
	}

	if opts.ResourceVersion != "" {
		if vs.Metadata.ResourceVersion != opts.ResourceVersion {
			return fmt.Errorf("conflict - resource version does not match")
		}
	}

	if int(input.Index) >= len(vs.VirtualHost.Routes) {
		return fmt.Errorf("invalid route index")
	}

	route := vs.VirtualHost.Routes[input.Index]

	var extAuthRouteExtension extauthpb.RouteExtension
	err = utils.UnmarshalExtension(route.RoutePlugins, extauth.ExtensionName, &extAuthRouteExtension)
	if err != nil {
		if err != utils.NotFoundError {
			return err
		}
	}

	extAuthRouteExtension.Disable = input.Disable

	if route.RoutePlugins == nil {
		route.RoutePlugins = &gloov1.RoutePlugins{}
	}

	if route.RoutePlugins.Extensions == nil {
		route.RoutePlugins.Extensions = &gloov1.Extensions{}
	}

	if route.RoutePlugins.Extensions.Configs == nil {
		route.RoutePlugins.Extensions.Configs = make(map[string]*types.Struct)
	}

	extStruct, err := protoutils.MarshalStruct(&extAuthRouteExtension)
	if err != nil {
		return err
	}
	route.RoutePlugins.Extensions.Configs[extauth.ExtensionName] = extStruct

	_, err = vsClient.Write(vs, clients.WriteOpts{OverwriteExisting: true})
	return err
}
