package get

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/get"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/common"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	helpersExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/helpers"

	"github.com/spf13/cobra"
)

func VSGet(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		// Use command constants to aid with replacement.
		Use:     constants.VIRTUAL_SERVICE_COMMAND.Use,
		Aliases: constants.VIRTUAL_SERVICE_COMMAND.Aliases,
		Short:   "read a virtualservice or list virtualservices in a namespace",
		Long:    "usage: glooctl get virtualservice [NAME] [--namespace=namespace] [-o FORMAT]",
		RunE: func(cmd *cobra.Command, args []string) error {
			virtualServices, err := common.GetVirtualServices(common.GetName(args, opts), opts)
			if err != nil {
				return err
			}
			helpersExt.PrintVirtualServices(virtualServices, opts.Top.Output)
			return nil
		},
	}
	cmd.AddCommand(get.Routes(opts))
	return cmd
}
