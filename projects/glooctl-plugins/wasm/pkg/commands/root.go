package commands

import (
	"context"

	ctxo "github.com/deislabs/oras/pkg/context"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/wasm/pkg/commands/deploy"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/wasm/pkg/commands/version"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/wasm/pkg/defaults"
	"github.com/solo-io/wasm-kit/pkg/abi"
	wasmcmd "github.com/solo-io/wasm-kit/pkg/commands"
	"github.com/solo-io/wasm-kit/pkg/commands/opts"
	"github.com/spf13/cobra"
)

func RootCommand(ctx context.Context) *cobra.Command {
	var (
		genOpts  opts.GeneralOptions
		authOpts opts.AuthOptions
	)

	commandGen := wasmcmd.NewCommandGen("glooctl wasm", defaults.WasmImageDir, defaults.WasmCredentialsFile)

	cmd := &cobra.Command{
		Use:   "wasm [command]",
		Short: "The interface for managing Gloo Edge WASM filters",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if genOpts.Verbose {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				ctx = ctxo.WithLoggerDiscarded(ctx)
			}

			if len(authOpts.CredentialsFiles) == 0 {
				authOpts.CredentialsFiles = []string{defaults.WasmCredentialsFile}
			}
		},
	}

	cmd.AddCommand(
		append(
			commandGen.CommandsWithAuth(ctx, &authOpts),
			commandGen.TagCmd(ctx),
			commandGen.InitCmdWithVersions(getAllGlooVersions()),
			commandGen.BuildCmd(ctx),
			commandGen.LoginCmd(),
			commandGen.ListCmd(),
			deploy.DeployCmd(ctx, cmd.PersistentPreRun),
			deploy.UndeployCmd(ctx),
			version.Command(),
		)...,
	)

	return cmd
}

func getAllGlooVersions() []abi.Platform {
	var platforms []abi.Platform
	glooVersions := abi.GetVersions(abi.GetAllSupportedPlatforms(), abi.PlatformNameGloo)
	for _, ver := range glooVersions {
		platforms = append(platforms, abi.Platform{Name: abi.PlatformNameGloo, Version: ver})
	}
	return platforms
}
