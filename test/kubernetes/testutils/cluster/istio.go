package cluster

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rotisserie/eris"
	glooruntime "github.com/solo-io/gloo/test/kubernetes/testutils/runtime"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/go-utils/contextutils"
	"k8s.io/apimachinery/pkg/util/version"
)

const (
	// TODO(npolshak): Add support for other profiles (ambient, etc.)
	minimalProfile = "minimal"
)

var (
	// https://github.com/istio/istio/issues/55980
	// The `istioctl install` command does not automatically set the default revision tag starting in 1.24.0
	// We need to manually set/delete the default revision tag starting in 1.24.0
	applyDefaultRevisionTagAfterVersion = version.MustParseSemantic("1.24.0")
)

func GetIstioctl(ctx context.Context) (string, error) {
	// Download istioctl binary
	istioctlBinary, err := downloadIstio(ctx, getIstioVersion())
	if err != nil {
		return "", fmt.Errorf("failed to download istio: %w", err)
	}
	contextutils.LoggerFrom(ctx).Infof("Using Istio binary '%s'", istioctlBinary)

	return istioctlBinary, nil
}

func InstallMinimalIstio(
	ctx context.Context,
	istioctlBinary, kubeContext string,
) error {
	operatorFileContent := generateIstioOperatorFileContent("", minimalProfile)
	operatorFile := filepath.Join(os.TempDir(), "istio-operator.yaml")

	err := os.WriteFile(operatorFile, []byte(operatorFileContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write operator file: %w", err)
	}

	return installIstioOperator(ctx, istioctlBinary, kubeContext, "")
}

func InstallRevisionedIstio(
	ctx context.Context,
	istioctlBinary, kubeContext, revision, profile string,
) error {
	operatorFileContent := generateIstioOperatorFileContent(revision, profile)
	operatorFile := filepath.Join(os.TempDir(), "istio-operator.yaml")

	err := os.WriteFile(operatorFile, []byte(operatorFileContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write operator file: %w", err)
	}

	err = installIstioOperator(ctx, istioctlBinary, kubeContext, operatorFile)
	if err != nil {
		return err
	}

	if manuallyManageDefaultRevisionTag() {
		return setDefaultRevisionTag(ctx, istioctlBinary, kubeContext, revision)
	}

	return nil
}

func manuallyManageDefaultRevisionTag() bool {
	semver := version.MustParseSemantic(getIstioVersion())
	if semver.LessThan(applyDefaultRevisionTagAfterVersion) {
		return false
	}

	return true
}

func setDefaultRevisionTag(ctx context.Context, istioctlBinary, kubeContext, revision string) error {
	//  istioctl tag set default --revision <revision>
	cmd := exec.Command(istioctlBinary, "tag", "set", "default", "--revision", revision, "--context", kubeContext)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("setting default revision tag failed: %w, output: %s", err, string(out))
	}

	return ctx.Err()
}

func deleteDefaultRevisionTag(istioctlBinary, kubeContext string) error {
	//  istioctl tag delete default
	cmd := exec.Command(istioctlBinary, "tag", "delete", "default", "--context", kubeContext)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("deleting default revision tag failed: %w, output: %s", err, string(out))
	}

	return nil
}

// TODO(npolshak): Add Istio dependency to define operator in code instead of writing file
func generateIstioOperatorFileContent(revision, profile string) string {
	// use minimal as the default profile if none is provided
	if profile == "" {
		profile = minimalProfile
	}

	baseContent := `
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  namespace: istio-system
spec:
  profile: %s
`
	if revision != "" {
		return fmt.Sprintf(baseContent+"  revision: %s\n", profile, revision)
	}
	return fmt.Sprintf(baseContent, profile)
}

func installIstioOperator(
	ctx context.Context,
	istioctlBinary, kubeContext, operatorFile string) error {
	if testutils.ShouldSkipIstioInstall() {
		return nil
	}

	//  istioctl install -y --context <kube-context> -f <operator-file>
	cmd := exec.Command(istioctlBinary, "install", "-y", "--context", kubeContext, "-f", operatorFile)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("istioctl install failed: %w, output: %s", err, string(out))
	}

	return ctx.Err()
}

func getIstioVersion() string {
	if version := os.Getenv(glooruntime.IstioVersionEnv); version != "" {
		return version
	} else {
		// Fail loudly if ISTIO_VERSION is not set
		panic(fmt.Sprintf("%s environment variable must be specified to run", glooruntime.IstioVersionEnv))
	}
}

// Download istioctl binary from istio.io/downloadIstio and returns the path to the binary
func downloadIstio(ctx context.Context, version string) (string, error) {
	if version == "" {
		contextutils.LoggerFrom(ctx).Infof("ISTIO_VERSION not specified, using istioctl from PATH")
		binaryPath, err := exec.LookPath("istioctl")
		if err != nil {
			return "", eris.New("ISTIO_VERSION environment variable must be specified or istioctl must be installed")
		}

		contextutils.LoggerFrom(ctx).Infof("using istioctl path: %s", binaryPath)

		return binaryPath, nil
	}
	installLocation := filepath.Join(testutils.GitRootDirectory(), ".bin")
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

	req, err := http.NewRequest(http.MethodGet, "https://istio.io/downloadIstio", nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	cmd := exec.Command("sh", "-")

	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", glooruntime.IstioVersionEnv, version))
	cmd.Dir = installLocation

	cmd.Stdin = res.Body
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return binaryLocation, err
}

func UninstallIstio(istioctlBinary, kubeContext string) error {
	if manuallyManageDefaultRevisionTag() {
		err := deleteDefaultRevisionTag(istioctlBinary, kubeContext)
		if err != nil {
			return fmt.Errorf("failed to delete default revision tag: %w", err)
		}
	}

	// sh -c yes | istioctl uninstall —purge —context <kube-context>
	cmd := exec.Command("sh", "-c", fmt.Sprintf("yes | %s uninstall --purge --context %s", istioctlBinary, kubeContext))

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("istioctl uninstall failed: %w, output: %s", err, string(out))
	}
	return nil
}

// CreateIstioBugReport generates an istioctl bug report and moves it to the _output directory
func CreateIstioBugReport(ctx context.Context, istioctlBinary, kubeContext, artifactOutputDir string) {
	// Generate istioctl bug report
	if istioctlBinary == "" {
		contextutils.LoggerFrom(ctx).Panic("istioctl binary not set. Cannot generate istioctl bug report.")
	}

	bugReportCmd := exec.Command(istioctlBinary, "bug-report", "--full-secrets", "--context", kubeContext)
	bugReportErr := bugReportCmd.Run()
	if bugReportErr != nil {
		fmt.Println("Error generating bug report:", bugReportErr)
	}
	mvCmd := exec.Command("mv", "bug-report.tar.gz", artifactOutputDir)
	mvErr := mvCmd.Run()
	if mvErr != nil {
		fmt.Println("Error moving bug report file:", mvErr)
	}
}
