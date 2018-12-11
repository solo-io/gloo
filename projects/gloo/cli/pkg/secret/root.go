package secret

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func CreateCmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secret",
		Aliases: []string{"s", "secret"},
		Short:   "Create a secret",
		Long:    "Create a secret",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(awsCmd(opts))
	cmd.AddCommand(tlsCmd(opts))

	return cmd
}
