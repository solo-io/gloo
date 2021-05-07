package version

import (
	"fmt"

	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/constants"
	"github.com/spf13/cobra"
)

func RootCmd(_ *options.Options) *cobra.Command {

	cmd := &cobra.Command{
		Use:     constants.VERSION_COMMAND.Use,
		Aliases: constants.VERSION_COMMAND.Aliases,
		Short:   constants.VERSION_COMMAND.Short,
		Long:    constants.VERSION_COMMAND.Long,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("glooctl-fed %s", version.Version)
			return nil
		},
	}

	return cmd
}
