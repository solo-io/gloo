package constants

import "github.com/spf13/cobra"

var (
	OBSERVABILITY_COMMAND = cobra.Command{
		Use:     "observability",
		Aliases: []string{"o", "obs", "observe"},
		Short:   "root command for observability functionality",
	}

	CONFIG_COMMAND = cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg", "settings", "set"},
		Short:   "root command for settings",
	}

	CONFIG_EXTAUTH_COMMAND = cobra.Command{
		Use:     "externalauth",
		Aliases: []string{"extauth"},
		Short:   "root command for external auth functionality",
	}
)
