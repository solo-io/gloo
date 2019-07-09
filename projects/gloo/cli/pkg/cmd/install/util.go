package install

import (
	"bytes"
	"context"
	"time"

	"github.com/solo-io/gloo/pkg/cliutil/install"
	install2 "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/install"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/pkg/version"
	optionsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const PersistentVolumeClaim = "PersistentVolumeClaim"

func validateLicenseKey(extraOptions *optionsExt.ExtraOptions) error {
	if extraOptions.Install.LicenseKey == "" {
		return errors.Errorf("you must provide a valid license key to be able to install GlooE")
	}
	if err := license.IsLicenseValid(context.TODO(), extraOptions.Install.LicenseKey); err != nil {
		return errors.Wrapf(err, "the license key you provided is invalid")
	}
	return nil
}

func getGlooEVersion(opts *options.Options) (string, error) {
	if !version.IsReleaseVersion() && opts.Install.HelmChartOverride == "" {
		return "", errors.Errorf("you must provide a GlooE Helm chart URI via the 'file' option " +
			"when running an unreleased version of glooctl")
	}
	return version.Version, nil
}

func getExcludeExistingPVCs(namespace string) install.ResourceMatcherFunc {
	return func(resource install.ResourceType) (bool, error) {
		cfg, err := kubeutils.GetConfig("", "")
		if err != nil {
			return false, err
		}
		kubeClient, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return false, err
		}

		// If this is a PVC, check if it already exists. If so, exclude this resource from the manifest.
		// We don't want to overwrite existing PVCs.
		if resource.TypeMeta.Kind == PersistentVolumeClaim {

			_, err := kubeClient.CoreV1().PersistentVolumeClaims(namespace).Get(resource.Metadata.Name, v1.GetOptions{})
			if err != nil {
				if !kubeerrors.IsNotFound(err) {
					return false, errors.Wrapf(err, "retrieving %s: %s.%s", PersistentVolumeClaim, namespace, resource.Metadata.Name)
				}
			} else {
				// The PVC exists, exclude it from manifest
				return true, nil
			}
		}
		return false, nil
	}
}

type NamespacedGlooKubeInstallClient struct {
	namespace string
	delegate  install2.GlooKubeInstallClient
}

func (i *NamespacedGlooKubeInstallClient) KubectlApply(manifest []byte) error {
	if i.namespace == "" {
		return i.delegate.KubectlApply(manifest)
	}
	return install.Kubectl(bytes.NewBuffer(manifest), "apply", "-n", i.namespace, "-f", "-")
}

func (i *NamespacedGlooKubeInstallClient) WaitForCrdsToBeRegistered(crds []string, timeout, interval time.Duration) error {
	return i.delegate.WaitForCrdsToBeRegistered(crds, timeout, interval)
}
