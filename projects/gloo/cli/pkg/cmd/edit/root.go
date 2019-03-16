package edit

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/edit"

	editOptions "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/edit/options"
	gloovs "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/edit/virtualservice"
	glooOptions "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/edit/route"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/edit/settings"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/edit/virtualservice"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmdutils"
	"github.com/spf13/cobra"
)

func RootCmd(opts *glooOptions.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	editFlags := &editOptions.EditOptions{Options: opts}

	cmd := edit.RootCmdWithEditOpts(editFlags, optionsFunc...)

	cmdutils.MustAddChildCommand(cmd, gloovs.RootCmd(editFlags), virtualservice.RateLimitConfig(editFlags, optionsFunc...))
	cmd.AddCommand(settings.RootCmd(editFlags, optionsFunc...))
	cmd.AddCommand(route.RootCmd(editFlags, optionsFunc...))
	return cmd
}
