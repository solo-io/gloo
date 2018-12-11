package install

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func Cmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install gloo on different platforms",
		Long:  "choose which system to install Gloo onto. options include: kubernetes",
	}
	cmd.AddCommand(kubeCmd(opts))
	return cmd
}
