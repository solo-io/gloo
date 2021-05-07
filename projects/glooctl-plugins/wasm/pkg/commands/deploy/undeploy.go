package deploy

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func UndeployCmd(ctx context.Context) *cobra.Command {
	opts := &options{
		remove: true,
	}
	cmd := &cobra.Command{
		Use:   "undeploy --id=<unique id>",
		Args:  cobra.MinimumNArgs(0),
		Short: "Remove an Envoy WASM Filter from the Gloo Gateway Proxies (Envoy).",
		Long: `Uses the Gloo Gateway CR to pull and run wasm filters.

Use --namespaces to constrain the namespaces of Gateway CRs to update.

Use --labels to use a match Gateway CRs by label.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.filter.Id == "" {
				return errors.Errorf("--id cannot be empty")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.providerType = Provider_Gloo
			return runDeploy(ctx, opts)
		},
	}
	opts.addIdToFlags(cmd.PersistentFlags())

	return cmd
}
