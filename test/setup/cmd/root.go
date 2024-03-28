package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func NewCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		SilenceErrors: true,
	}

	cmd.AddCommand(
		setupCommand(ctx),
	)
	return cmd
}
