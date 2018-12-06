package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

type Command struct {
	*cobra.Command
}

type ApplyOpts struct {
	File string
}

var RootOpts struct {
	Ctx       context.Context
	ApplyOpts ApplyOpts
}

var RootCmd = &Command{
	Command: &cobra.Command{
		Use:   "glooctl",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly apply a Cobra application.`,
	},
}

func init() {
	RootOpts.Ctx = context.Background()
	// global flags here
}
