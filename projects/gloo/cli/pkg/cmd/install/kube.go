package install

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-projects/pkg/version"
	v1 "k8s.io/api/core/v1"
	kubeerrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
			if err := createImagePullSecretIfNeeded(opts.Install, optsExt.Install); err != nil {
				return errors.Wrapf(err, "creating image pull secret")
			}
			if err := registerSettingsCrd(); err != nil {
				return errors.Wrapf(err, "registering settings crd")
			}
			if err := registerUpstreamCrd(); err != nil {
				return errors.Wrapf(err, "registering settings crd")
			}
			glooManifestBytes, err := readGlooManifest(opts)
			if err != nil {
				return errors.Wrapf(err, "reading gloo manifest")
			}
			if opts.Install.DryRun {
				fmt.Printf("%s", glooManifestBytes)
				return nil
			}
			return applyManifest(glooManifestBytes)
		},
	}

	return cmd
}

func readGlooManifest(opts *options.Options) ([]byte, error) {
	if opts.Install.GlooManifestOverride != "" {
		return readManifestFromFile(opts.Install.GlooManifestOverride)
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

func applyManifest(manifest []byte) error {
	kubectl := exec.Command("kubectl", "apply", "-n", InstallNamespace, "-f", "-")
	kubectl.Stdin = bytes.NewBuffer(manifest)
	kubectl.Stdout = os.Stdout
	kubectl.Stderr = os.Stderr
	return kubectl.Run()
}

func registerSettingsCrd() error {
	cfg, err := kubeutils.GetConfig("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	settingsClient, err := gloov1.NewSettingsClient(&factory.KubeResourceClientFactory{
		Crd:         gloov1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: kube.NewKubeCache(context.Background()),
	})

	return settingsClient.Register()
}

func registerUpstreamCrd() error {
	cfg, err := kubeutils.GetConfig("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	upstreamClient, err := gloov1.NewUpstreamClient(&factory.KubeResourceClientFactory{
		Crd:         gloov1.UpstreamCrd,
		Cfg:         cfg,
		SharedCache: kube.NewKubeCache(context.Background()),
	})

	return upstreamClient.Register()
}

func createImagePullSecretIfNeeded(install options.Install, installExt optionsExt.InstallExtended) error {
	if err := createNamespaceIfNotExist(); err != nil {
		return errors.Wrapf(err, "creating installation namespace")
	}
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

	kubectl := exec.Command("kubectl", "create", "secret", "docker-registry", "-n", InstallNamespace,
		"--docker-email", installExt.DockerAuth.Email,
		"--docker-username", installExt.DockerAuth.Username,
		"--docker-password", installExt.DockerAuth.Password,
		"--docker-server", installExt.DockerAuth.Server,
		imagePullSecretName,
	)
	kubectl.Stdout = os.Stdout
	kubectl.Stderr = os.Stderr
	return kubectl.Run()
}

func createNamespaceIfNotExist() error {
	restCfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return err
	}
	kube, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return err
	}
	installNamespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: InstallNamespace,
		},
	}
	if _, err := kube.CoreV1().Namespaces().Create(installNamespace); err != nil && !kubeerrs.IsAlreadyExists(err) {
		return err
	}
	return nil
}
