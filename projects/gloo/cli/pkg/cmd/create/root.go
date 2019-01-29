package create

import (
	"fmt"
	"log"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/create"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/create/secret"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/go-utils/cliutils"
	glooesecret "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/create/secret"
	"github.com/spf13/cobra"
)

/*
Implementation of new CLI extension format.
RootCmd is imported by the command which intends on extending the base command.
rootCmd is a pointer to the OS gloo create cmd, and MustReplaceCmd replaces it's
create virtual service command with the one located locally called VSCreate.
*/

func RootCmd(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	rootCmd := create.RootCmd(opts, optionsFunc...)
	// Replace old command with new one, replaces by name of command.
	cliutils.MustReplaceCmd(rootCmd, VSCreate(opts))

	// find the secret command and add to it
	mustAddChildCommand(rootCmd, secret.CreateCmd(opts), glooesecret.ExtAuthOathCmd(opts))

	return rootCmd
}

// This finds 'current' command in 'parent' command and adds the 'sub' command to it.
func mustAddChildCommand(parent *cobra.Command, current *cobra.Command, sub *cobra.Command) {
	parentCmds := parent.Commands()
	for _, old := range parentCmds {
		if old.Use == current.Use {
			old.AddCommand(sub)
			return
		}
	}
	err := fmt.Errorf("did not find child command to replace")
	log.Fatal(err)
}
