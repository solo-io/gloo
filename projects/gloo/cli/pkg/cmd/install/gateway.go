package install

import (
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/solo-io/solo-projects/pkg/cliutil"

	"github.com/ghodss/yaml"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/manifest"

	"github.com/solo-io/gloo/pkg/cliutil/install"
	"k8s.io/helm/pkg/chartutil"
	helmhooks "k8s.io/helm/pkg/hooks"
	"k8s.io/helm/pkg/renderutil"

	"github.com/pkg/errors"
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

			glooEVersion, err := getGlooEVersion(opts)
			if err != nil {
				return err
			}

			// Get location of Gloo helm chart
			helmChartArchiveUri := fmt.Sprintf(GlooEHelmRepoTemplate, glooEVersion)
			if helmChartOverride := opts.Install.HelmChartOverride; helmChartOverride != "" {
				helmChartArchiveUri = helmChartOverride
			}

			if path.Ext(helmChartArchiveUri) != ".tgz" && !strings.HasSuffix(helmChartArchiveUri, ".tar.gz") {
				return errors.Errorf("unsupported file extension for Helm chart URI: [%s]. Extension must "+
					"either be .tgz or .tar.gz", helmChartArchiveUri)
			}

			chart, err := install.GetHelmArchive(helmChartArchiveUri)
			if err != nil {
				return errors.Wrapf(err, "retrieving gloo helm chart archive")
			}

			// Passing fileName="" means "use default values".
			// TODO: change when we move to only one value file in the OS Gloo chart
			values, err := getValuesFromFile(chart, "", optsExt.Install.LicenseKey)
			if err != nil {
				return errors.Wrapf(err, "retrieving values")
			}

			// These are the .Release.* variables used during rendering
			renderOpts := renderutil.Options{
				ReleaseOptions: chartutil.ReleaseOptions{
					Namespace: opts.Install.Namespace,
					Name:      "glooe",
				},
			}

			/**************************************************************
			 *******************	Filter functions    *******************
			 **************************************************************/

			filterKnativeResources, err := install.GetKnativeResourceFilterFunction()
			if err != nil {
				return err
			}

			// Keep only CRDs and collect the names
			var crdNames []string
			keepOnlyCrds := func(input []manifest.Manifest) ([]manifest.Manifest, error) {

				var crdManifests []manifest.Manifest
				for _, man := range input {

					// Split manifest into individual YAML docs
					crdDocs := make([]string, 0)
					for _, doc := range strings.Split(man.Content, "---") {

						// We need to define this ourselves, because if we unmarshal into `apiextensions.CustomResourceDefinition`
						// we don't get the TypeMeta (in the yaml they are nested under `metadata`, but the k8s struct has
						// them as top level fields...)
						var resource struct {
							Metadata v1.ObjectMeta
							v1.TypeMeta
						}
						if err := yaml.Unmarshal([]byte(doc), &resource); err != nil {
							return nil, errors.Wrapf(err, "parsing resource: %s", doc)
						}

						// Skip non-CRD resources
						if resource.TypeMeta.Kind != install.CrdKindName {
							continue
						}

						// Check whether the CRD is a Helm "crd-install" hook.
						// If not, throw an error, because this will cause race conditions when installing with Helm (which is
						// not the case here, but we want to validate the manifests whenever we have the chance)
						helmCrdInstallHookAnnotation, ok := resource.Metadata.Annotations[helmhooks.HookAnno]
						if !ok || helmCrdInstallHookAnnotation != helmhooks.CRDInstall {
							return nil, errors.Errorf("CRD [%s] must be annotated as a Helm '%s' hook", resource.Metadata.Name, helmhooks.CRDInstall)
						}

						// Keep track of the CRD name
						crdNames = append(crdNames, resource.Metadata.Name)

						crdDocs = append(crdDocs, doc)
					}

					crdManifests = append(crdManifests, manifest.Manifest{
						Name:    man.Name,
						Head:    man.Head,
						Content: strings.Join(crdDocs, install.YamlDocumentSeparator),
					})
				}

				return crdManifests, nil
			}

			/**************************************************************
			 **************   End of filter functions   *******************
			 **************************************************************/

			// Helm uses the standard go log package. Redirect its output to the debug.log file  so that we don't
			// expose useless warnings to the user.
			log.SetOutput(cliutil.Logger)

			// Render and install CRD manifests
			crdManifestBytes, err := install.RenderChart(chart, values, renderOpts,
				install.ExcludeNotes,
				filterKnativeResources,
				keepOnlyCrds,
				install.ExcludeEmptyManifests)
			if err != nil {
				return errors.Wrapf(err, "rendering crd manifests")
			}
			if err := installManifest(crdManifestBytes, opts.Install.DryRun, ""); err != nil {
				return errors.Wrapf(err, "installing crd manifests")
			}
			// Only run if this is not a dry run
			if !opts.Install.DryRun {
				if err := waitForCrdsToBeRegistered(crdNames, time.Second*5, time.Millisecond*500); err != nil {
					return errors.Wrapf(err, "waiting for crds to be registered")
				}
			}

			// Render the rest of the GlooE manifest
			glooEManifestBytes, err := install.RenderChart(chart, values, renderOpts,
				install.ExcludeNotes,
				filterKnativeResources,
				install.ExcludeCrds,
				install.ExcludeEmptyManifests)
			if err != nil {
				return err
			}

			// Remove existing PVCs
			glooEManifestBytes, err = removeExistingPVCs(glooEManifestBytes, opts.Install.Namespace)
			if err != nil {
				return errors.Wrapf(err, "checking for existing PVCs to remove")
			}

			if err := installManifest(glooEManifestBytes, opts.Install.DryRun, opts.Install.Namespace); err != nil {
				return err
			}

			fmt.Printf("\nGlooE was successfully installed! 🎉\n")

			return nil

		},
	}

	return cmd
}
