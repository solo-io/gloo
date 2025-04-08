package gatewayapi

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/convert"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway-api",
		Short: "Gateway API specific commands",
	}
	cmd.SilenceUsage = true
	cmd.AddCommand(convert.RootCmd(opts))
	return cmd
}
