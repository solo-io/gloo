package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ApplyCmd = &Command{
	Command: &cobra.Command{
		Use:   "apply",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly apply a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply called")
		},
	},
}

func init() {
	RootCmd.AddCommand(ApplyCmd.Command)
	ApplyCmd.Flags().StringVarP(&RootOpts.ApplyOpts.File, "file", "f", "", "yaml file to read when applying resources")
}
