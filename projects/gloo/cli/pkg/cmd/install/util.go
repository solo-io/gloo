package install

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/solo-io/gloo/pkg/cliutil/install"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-projects/pkg/version"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const PersistentVolumeClaim = "PersistentVolumeClaim"

func getGlooEVersion(opts *options.Options) (string, error) {
	if !version.IsReleaseVersion() && opts.Install.HelmChartOverride == "" {
		return "", errors.Errorf("you must provide a GlooE Helm chart URI via the 'file' option " +
			"when running an unreleased version of glooctl")
	}
	return version.Version, nil

}

func removeExistingPVCs(manifestBytes []byte, namespace string) ([]byte, error) {

	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	var docs []string
	for _, doc := range strings.Split(string(manifestBytes), "---") {

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

		// If this is a PVC, check if it already exists. If so, exclude this resource from the manifest.
		// We don't want to overwrite existing PVCs.
		if resource.TypeMeta.Kind == PersistentVolumeClaim {

			_, err := kubeClient.CoreV1().PersistentVolumeClaims(namespace).Get(resource.Metadata.Name, v1.GetOptions{})
			if err != nil {
				if !kubeerrors.IsNotFound(err) {
					return nil, errors.Wrapf(err, "retrieving %s: %s.%s", PersistentVolumeClaim, namespace, resource.Metadata.Name)
				}
			} else {
				// The PVC exists, exclude it from manifest
				continue
			}
		}

		docs = append(docs, doc)
	}
	return []byte(strings.Join(docs, install.YamlDocumentSeparator)), nil
}

// TODO: copied over and modified for a quick fix, improve
//noinspection GoNameStartsWithPackageName
func installManifest(manifest []byte, isDryRun bool, namespace string) error {
	if isDryRun {
		fmt.Printf("%s", manifest)
		// For safety, print a YAML separator so multiple invocations of this function will produce valid output
		fmt.Println("\n---")
		return nil
	}
	if err := kubectlApply(manifest, namespace); err != nil {
		return errors.Wrapf(err, "running kubectl apply on manifest")
	}
	return nil
}

func kubectlApply(manifest []byte, namespace string) error {
	return kubectl(bytes.NewBuffer(manifest), "apply", "-n", namespace, "-f", "-")
}

func kubectl(stdin io.Reader, args ...string) error {
	kubectl := exec.Command("kubectl", args...)
	if stdin != nil {
		kubectl.Stdin = stdin
	}
	kubectl.Stdout = os.Stdout
	kubectl.Stderr = os.Stderr
	return kubectl.Run()
}
