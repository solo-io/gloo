package deploy

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/proto"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.StandardLogger()

func DeployCmd(ctx context.Context, parentPreRun func(cmd *cobra.Command, args []string)) *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:   "deploy <image> --id=<unique id> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--labels <key1=val1,key2=val2>]",
		Short: "Deploy an Envoy WASM Filter to the Gloo Gateway Proxies (Envoy).",
		Long: `Deploys an Envoy WASM Filter to Gloo Gateway Proxies.
	
This CLI uses the Gloo Gateway CR to pull and run wasm filters.
	
Use --namespaces to constrain the namespaces of Gateway CRs to update.
	
Use --labels to use a match Gateway CRs by label.
	`,
		Args: cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			parentPreRun(cmd, args)
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.filter.Image = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.filter.Id == "" {
				return errors.Errorf("--id cannot be empty")
			}
			opts.providerType = Provider_Gloo
			// If we were passed a config via CLI flag, default config type to StringValue
			if opts.filterConfig != "" {
				wrappedConfig := &wrappers.StringValue{
					Value: opts.filterConfig,
				}
				marshaledAny, err := proto.Marshal(wrappedConfig)
				if err != nil {
					return errors.Errorf("--config value could not be parsed")
				}
				opts.filter.Config = &any.Any{
					TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
					Value:   marshaledAny,
				}
			}
			return runDeploy(ctx, opts)
		},
	}

	opts.addToFlags(cmd.PersistentFlags())

	return cmd
}

func runDeploy(ctx context.Context, opts *options) error {
	deployer, err := makeDeployer(ctx, opts)
	if err != nil {
		return err
	}

	if opts.remove {
		return deployer.RemoveFilter(&opts.filter)
	}

	return deployer.ApplyFilter(&opts.filter)
}
