package constants

import "github.com/spf13/cobra"

var (
	OBSERVABILITY_COMMAND = cobra.Command{
		Use:     "observability",
		Aliases: []string{"o", "obs", "observe"},
		Short:   "root command for observability functionality",
	}

	EDIT_COMMAND = cobra.Command{
		Use:     "edit",
		Aliases: []string{"ed", "settings", "set"},
		Short:   "root command for editing",
	}

	SETTINGS_COMMAND = cobra.Command{
		Use:     "settings",
		Aliases: []string{"st", "set"},
		Short:   "root command for settings",
	}

	CONFIG_EXTAUTH_COMMAND = cobra.Command{
		Use:     "externalauth",
		Aliases: []string{"extauth"},
		Short:   "root command for external auth functionality",
	}

	CONFIG_RATELIMIT_COMMAND = cobra.Command{
		Use:   "ratelimit",
		Short: "root command for rate limit functionality",
	}
)
