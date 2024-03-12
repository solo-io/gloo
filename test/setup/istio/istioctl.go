package istio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/test/setup/helpers"
	"github.com/solo-io/gloo/test/setup/kubernetes"
	"istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/go-utils/contextutils"
)

func Install(
	ctx context.Context,
	istioctlBinary, clusterName string,
	operator *v1alpha1.IstioOperator,
	info *kubernetes.Cluster,
) error {
	if operator == nil {
		return nil
	}

	if istioctlBinary == "" {
		contextutils.LoggerFrom(ctx).Infof("Istio binary not specified, skipping installation")
		return nil
	}

	var (
		client     = info.GetKubernetes()
		controller = info.GetController()
	)

	uniqueNS := map[string]bool{"istio-system": true}

	for _, gateway := range append(operator.Spec.Components.IngressGateways, operator.Spec.Components.EgressGateways...) {
		if !gateway.Enabled.Value {
			continue
		}

		if !uniqueNS[gateway.Namespace] {
			if err := createNamespace(ctx, client, controller, gateway); err != nil {
				return err
			}
			uniqueNS[gateway.Namespace] = true
		}
	}

	timerFn := helpers.TimerFunc(fmt.Sprintf("[%s] istio installation", clusterName))
	defer timerFn()

	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(operator); err != nil {
		return err
	}

	cmd := exec.Command(istioctlBinary, "install", "-y", "--context", info.GetKubeContext(), "--kubeconfig", info.GetKubeConfig(), "-f", "-")
	cmd.Stdin = buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("istioctl install failed: %w", err)
	}

	for ns := range uniqueNS {
		if err := kubernetes.RolloutStatus(ns, info); err != nil {
			return fmt.Errorf("istio rollout failed: %w", err)
		}
	}

	return ctx.Err()
}

func createNamespace[Object interface{ GetNamespace() string }](
	ctx context.Context,
	k8sclient k8s.Interface,
	client client.Client,
	objs ...Object,
) error {
	for _, obj := range objs {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: obj.GetNamespace(),
			},
		}

		if err := kubernetes.CreateOrUpdate[*corev1.Namespace](ctx, namespace, client); err != nil {
			return err
		}
	}
	return nil
}

// Download istioctl binary from istio.io/downloadIstio and returns the path to the binary
func DownloadIstio(ctx context.Context, version string) (string, error) {
	if version == "" {
		contextutils.LoggerFrom(ctx).Infof("ISTIOCTL_VERSION not specified, using istioctl from PATH")
		binaryPath, err := exec.LookPath("istioctl")
		if err != nil {
			return "", errors.New("ISTIOCTL_VERSION environment variable must be specified or istioctl must be installed")
		}

		contextutils.LoggerFrom(ctx).Infof("using istioctl path: %s", binaryPath)

		return binaryPath, nil
	}
	installLocation := filepath.Join(helpers.GlooDirectory(), ".bin")
	binaryDir := filepath.Join(installLocation, fmt.Sprintf("istio-%s", version), "bin")
	binaryLocation := filepath.Join(binaryDir, "istioctl")

	fileInfo, _ := os.Stat(binaryLocation)
	if fileInfo != nil {
		return binaryLocation, nil
	}
	if err := os.MkdirAll(binaryDir, 0755); err != nil {
		return "", eris.Wrap(err, "create directory")
	}

	if istioctlDownloadFrom := os.Getenv("ISTIOCTL_DOWNLOAD_FROM"); istioctlDownloadFrom != "" {
		osName := "linux"
		if runtime.GOOS == "darwin" {
			osName = "osx"
		}

		arch := runtime.GOARCH
		archModifier := fmt.Sprintf("-%s", arch)

		if osName == "osx" && arch != "arm64" {
			archModifier = ""
		}

		url := fmt.Sprintf("%s/%s/istioctl-%s-%s%s.tar.gz", istioctlDownloadFrom, version, version, osName, archModifier)

		// Use curl and tar to download and extract the file
		cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -sSL %s | tar -xz -C %s", url, binaryDir))
		if err := cmd.Run(); err != nil {
			return "", eris.Wrapf(err, "download and extract istioctl, cmd: %s", cmd.Args)
		}
		// Change permissions
		if err := os.Chmod(binaryLocation, 0755); err != nil {
			return "", eris.Wrap(err, "change permissions")
		}
		return binaryLocation, nil
	}

	req, err := http.NewRequest("GET", "https://istio.io/downloadIstio", nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	cmd := exec.Command("sh", "-")

	cmd.Env = append(cmd.Env, fmt.Sprintf("ISTIO_VERSION=%s", version))
	cmd.Dir = installLocation

	cmd.Stdin = res.Body
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return binaryLocation, err
}
