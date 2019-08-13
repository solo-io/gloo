package install

import (
	"fmt"

	glooInstall "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/install"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	optionsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

const (
	GlooEHelmRepoTemplate = "https://storage.googleapis.com/gloo-ee-helm/charts/gloo-ee-%s.tgz"
)

func GatewayCmd(opts *options.Options, optsExt *optionsExt.ExtraOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "install the GlooE Gateway on kubernetes",
		Long:  "requires kubectl to be installed",
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := validateLicenseKey(optsExt); err != nil {
				return err
			}

			if !opts.Install.DryRun {
				fmt.Printf("Starting GlooE installation...\n")
			}

			installSpec, err := GetInstallSpec(opts, optsExt)
			if err != nil {
				return err
			}

			kubeInstallClient := NamespacedGlooKubeInstallClient{
				namespace: opts.Install.Namespace,
				delegate:  &glooInstall.DefaultGlooKubeInstallClient{},
			}

			if err := glooInstall.InstallGloo(opts, *installSpec, &kubeInstallClient); err != nil {
				return err
			}

			if !opts.Install.DryRun {
				fmt.Printf("\nGlooE was successfully installed!\n")
			}

			return nil
		},
	}
	return cmd
}

func GetInstallSpec(opts *options.Options, optsExt *optionsExt.ExtraOptions) (*glooInstall.GlooInstallSpec, error) {
	glooEVersion, err := getGlooEVersion(opts)
	if err != nil {
		return nil, err
	}

	// Get location of Gloo helm chart
	helmChartArchiveUri := fmt.Sprintf(GlooEHelmRepoTemplate, glooEVersion)
	if helmChartOverride := opts.Install.HelmChartOverride; helmChartOverride != "" {
		helmChartArchiveUri = helmChartOverride
	}

	extraValues := map[string]interface{}{
		"license_key": optsExt.Install.LicenseKey,
	}

	if opts.Install.Upgrade {
		extraValues["gloo"] = map[string]interface{}{
			"gateway": map[string]interface{}{
				"upgrade": "true",
			},
		}
	} else {
		extraValues["gloo"] = map[string]interface{}{
			"namespace": map[string]interface{}{
				"create": "true",
			},
		}
	}

	return &glooInstall.GlooInstallSpec{
		HelmArchiveUri:   helmChartArchiveUri,
		ProductName:      "glooe",
		ValueFileName:    "",
		ExtraValues:      extraValues,
		ExcludeResources: getExcludeExistingPVCs(opts.Install.Namespace),
	}, nil
}
