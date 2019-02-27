package install

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-projects/pkg/version"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	optionsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

// TODO: support configuring install namespace (blocked on grafana and prometheus pods allowing namespace to be configurable
// requires changing a few places in the yaml as well
const (
	InstallNamespace    = "gloo-system"
	imagePullSecretName = "solo-io-docker-secret"
)

func KubeCmd(opts *options.Options, optsExt *optionsExt.ExtraOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube",
		Short: fmt.Sprintf("install Gloo on kubernetes to the %s namespace", InstallNamespace),
		Long:  "requires kubectl to be installed",
		RunE: func(cmd *cobra.Command, args []string) error {
			glooManifestBytes, err := readGlooManifest(opts)
			if err != nil {
				return errors.Wrapf(err, "reading gloo manifest")
			}
			installer, err := newInstaller(InstallNamespace, glooManifestBytes)
			if err != nil {
				return err
			}
			return installGlooe(installer, opts, optsExt)
		},
	}

	return cmd
}

func installGlooe(installer *Installer, opts *options.Options, optsExt *optionsExt.ExtraOptions) error {
	exists, err := installer.createNamespaceIfNotExist()
	if err != nil {
		return err
	}
	if exists {
		fmt.Printf("Former installation of glooe found in (%s)\n", installer.Namespace)
		fmt.Printf("Upgrading current installation to version (%s)\n", version.Version)
		if err := installer.upgrade(); err != nil {
			return err
		}
	}

	if err := installer.createImagePullSecretIfNeeded(opts.Install, optsExt.Install); err != nil {
		return errors.Wrapf(err, "creating image pull secret")
	}
	if err := installer.registerSettingsCrd(); err != nil {
		return errors.Wrapf(err, "registering settings crd")
	}
	if opts.Install.DryRun {
		fmt.Printf("%s", installer.Manifest)
		return nil
	}
	return installer.applyManifest()
}

func readGlooManifest(opts *options.Options) ([]byte, error) {
	if opts.Install.HelmChartOverride != "" {
		return readManifestFromFile(opts.Install.HelmChartOverride)
	}
	if version.Version == version.UndefinedVersion || version.Version == version.DevVersion {
		return nil, errors.Errorf("You must provide a file containing the gloo manifest when running an unreleased version of glooctl.")
	}
	return readManifest(version.Version)
}

func readManifestFromFile(path string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading file %s", path)
	}
	return bytes, nil
}
