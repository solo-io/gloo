package cmd

import (
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-function-discovery/pkg/util"
	"github.com/spf13/cobra"
)

// registerCmd is CLI command to register CRDs. This should only
// be necessary during development as Gloo would do this in
// production environment
func registerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register Gloo CRDs (only needed for developing)",
		RunE: func(c *cobra.Command, args []string) error {
			storageClient, err := util.GetStorageClient(c)
			if err != nil {
				return errors.Wrap(err, "unable to get storage client")
			}
			err = storageClient.V1().Register()
			if err != nil {
				return errors.Wrap(err, "unable to register resource definitions")
			}
			return nil
		},
	}
	return cmd
}
