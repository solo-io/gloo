package constants

import "github.com/spf13/cobra"

var (
	OBSERVABILITY_COMMAND = cobra.Command{
		Use:     "observability",
		Aliases: []string{"o", "obs", "observe"},
		Short:   "root command for observability functionality",
	}
)
