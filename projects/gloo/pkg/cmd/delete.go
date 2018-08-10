package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var DeleteCmd = &Command{
	Command: &cobra.Command{
		Use:   "delete",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly delete a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
		},
	},
}

func init() {
	RootCmd.AddCommand(DeleteCmd.Command)
}
