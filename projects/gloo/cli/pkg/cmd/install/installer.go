package install

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"k8s.io/helm/pkg/manifest"
	"sigs.k8s.io/yaml"

	"github.com/solo-io/solo-projects/pkg/cliutil"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"k8s.io/client-go/rest"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	optionsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	v1 "k8s.io/api/core/v1"
	kubeerrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	Unknown               = "Unknown"
	PersistentVolumeClaim = "PersistentVolumeClaim"
)

type Installer struct {
	Namespace string
	Manifest  []byte
	Ctx       context.Context
	Cfg       *rest.Config
	Client    *kubernetes.Clientset
}

func newInstaller(namespace string, manifest []byte) (*Installer, error) {
	ctx := context.TODO()
	cfg, err := kubeutils.GetConfig("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	installer := &Installer{
		Ctx:       ctx,
		Cfg:       cfg,
		Manifest:  manifest,
		Namespace: namespace,
		Client:    kubeClient,
	}
	return installer, nil
}

func (installer *Installer) upgrade() error {
	manifests := parseYaml(installer.Manifest)
	updatedManifests := make([]string, 0, len(manifests))

	for _, spec := range manifests {
		switch spec.Head.Kind {
		case Unknown:
			continue
		case PersistentVolumeClaim:
			if found, err := installer.upgradePersistentVolumeClaim(spec); err != nil {
				return err
			} else if !found {
				updatedManifests = append(updatedManifests, spec.Content)
			}
		default:
			updatedManifests = append(updatedManifests, spec.Content)
		}
	}

	installer.Manifest = []byte(strings.Join(updatedManifests, "\n---\n"))
	return nil
}

func (installer *Installer) registerSettingsCrd() error {
	settingsClient, err := gloov1.NewSettingsClient(&factory.KubeResourceClientFactory{
		Crd:         gloov1.SettingsCrd,
		Cfg:         installer.Cfg,
		SharedCache: kube.NewKubeCache(context.TODO()),
	})
	if err != nil {
		return err
	}

	return settingsClient.Register()
}

func (installer *Installer) createImagePullSecretIfNeeded(install options.Install, installExt optionsExt.InstallExtended) error {
	dockerSecretDesired := installExt.DockerAuth.Username != "" ||
		installExt.DockerAuth.Password != "" ||
		installExt.DockerAuth.Email != ""

	if !dockerSecretDesired {
		return nil
	}

	validOpts := installExt.DockerAuth.Username != "" &&
		installExt.DockerAuth.Password != "" &&
		installExt.DockerAuth.Email != "" &&
		installExt.DockerAuth.Server != ""

	if !validOpts {
		return errors.Errorf("must provide one of each flag for docker authentication: \n" +
			"--docker-email \n" +
			"--docker-username \n" +
			"--docker-password \n")
	}

	if install.DryRun {
		return nil
	}

	kubectl := exec.Command("kubectl", "create", "secret", "docker-registry", "-n", installer.Namespace,
		"--docker-email", installExt.DockerAuth.Email,
		"--docker-username", installExt.DockerAuth.Username,
		"--docker-password", installExt.DockerAuth.Password,
		"--docker-server", installExt.DockerAuth.Server,
		imagePullSecretName,
	)
	kubectl.Stdout = cliutil.Logger
	kubectl.Stderr = cliutil.Logger
	return kubectl.Run()
}

func (installer *Installer) createNamespaceIfNotExist() (exists bool, err error) {
	installNamespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: InstallNamespace,
		},
	}
	if _, err := installer.Client.CoreV1().Namespaces().Create(installNamespace); err != nil {
		if kubeerrs.IsAlreadyExists(err) {
			return true, nil
		}
		return false, err
	}
	fmt.Printf("Installing glooe into namespace (%s)\n", installer.Namespace)
	return false, nil
}

func (installer *Installer) applyManifest() error {
	kubectl := exec.Command("kubectl", "apply", "-n", installer.Namespace, "-f", "-")
	kubectl.Stdin = bytes.NewBuffer(installer.Manifest)
	kubectl.Stdout = cliutil.Logger
	kubectl.Stderr = cliutil.Logger
	if err := kubectl.Run(); err != nil {
		return err
	}
	fmt.Println("Finished installing glooe")
	return nil
}

func (installer *Installer) upgradePersistentVolumeClaim(manifest manifest.Manifest) (bool, error) {
	var manifestPVC v1.PersistentVolumeClaim
	err := yaml.Unmarshal([]byte(manifest.Content), &manifestPVC)
	if err != nil {
		return false, err
	}

	pvcs, err := installer.Client.CoreV1().PersistentVolumeClaims(installer.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, pvc := range pvcs.Items {
		if pvc.Name == manifestPVC.Name {
			return true, nil
		}
	}
	return false, nil
}
