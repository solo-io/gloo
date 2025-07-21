package gatewayapi

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/migrate"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func RootCmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway-api",
		Short: "Gateway API specific commands",
	}
	cmd.SilenceUsage = true
	cmd.AddCommand(migrate.RootCmd(opts))
	return cmd
}
