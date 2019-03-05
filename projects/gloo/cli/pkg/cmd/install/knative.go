package install

import (
	"fmt"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func KnativeCmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knative",
		Short: "install Knative with GlooE on kubernetes",
		Long:  "requires kubectl to be installed",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("this feature will be available soon")
			return nil
		},
	}

	return cmd
}
