package gatewayapi

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/envoy"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway-api",
		Short: "Gateway API specific commands",
	}
	cmd.SilenceUsage = true
	cmd.AddCommand(envoy.RootCmd())
	return cmd
}