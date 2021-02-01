package version

import (
	"fmt"

	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version of glooctl wasm",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("glooctl-wasm %s\n", version.Version)
			return nil
		},
	}
	cmd.SilenceUsage = true
	return cmd
}
