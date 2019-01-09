package get

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/get"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/spf13/cobra"
)

/*
Implementation of new CLI extension format.
RootCmd is imported by the command which intends on extending the base command.
rootCmd is a pointer to the OS gloo get cmd, and MustReplaceCmd replaces it's
get virtual service command with the one located locally called VSGet.
*/

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	rootCmd := get.RootCmd(opts, optionsFunc...)
	// replaces OS gloo get vs command with local get vs command.
	cliutils.MustReplaceCmd(rootCmd, VSGet(opts))
	return rootCmd
}
